package commands

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ms-henglu/armstrong/autorest"
	"github.com/ms-henglu/armstrong/resource"
	"github.com/ms-henglu/armstrong/resource/resolver"
	"github.com/ms-henglu/armstrong/resource/types"
	"github.com/ms-henglu/armstrong/swagger"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type GenerateCommand struct {
	// common options
	verbose           bool
	workingDir        string
	useRawJsonPayload bool

	// create with example path
	path         string
	resourceType string
	overwrite    bool

	// create with swagger path
	swaggerPath string

	// create with autorest config, TODO: remove them? because the tag contains swaggers from different api-versions
	readmePath string
	tag        string
}

func (c *GenerateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("generate")

	// common options
	fs.StringVar(&c.workingDir, "working-dir", "", "output path to Terraform configuration files")
	fs.BoolVar(&c.useRawJsonPayload, "raw", false, "whether use raw json payload in 'body'")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")

	// generate with example options
	fs.StringVar(&c.path, "path", "", "path to a swagger 'Create' example")
	fs.StringVar(&c.resourceType, "type", "resource", "the type of the resource to be generated, allowed values: 'resource'(supports CRUD) and 'data'(read-only). Defaults to 'resource'")
	fs.BoolVar(&c.overwrite, "overwrite", false, "whether overwrite existing terraform configurations")

	// generate with swagger options
	fs.StringVar(&c.swaggerPath, "swagger", "", "path or directory to swagger.json files")

	// generate with autorest config
	fs.StringVar(&c.readmePath, "readme", "", "path to the autorest config file(readme.md)")
	fs.StringVar(&c.tag, "tag", "", "tag in the autorest config file(readme.md)")

	fs.Usage = func() { logrus.Error(c.Help()) }

	return fs
}

func (c GenerateCommand) Help() string {
	helpText := `
Usage:
	armstrong generate -path <path to a swagger 'Create' example> [-working-dir <output path to Terraform configuration files>]
	armstrong generate -swagger <path/dir to the swagger files> [-working-dir <output path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c GenerateCommand) Synopsis() string {
	return "Generate testing files including terraform configuration for dependencies and testing resource."
}

func (c GenerateCommand) Run(args []string) int {
	logrus.Debugf("args: %+v", args)
	f := c.flags()
	if err := f.Parse(args); err != nil {
		logrus.Errorf("Error parsing command-line flags: %+v", err)
		return 1
	}
	logrus.Debugf("flags: %+v", f)
	if c.verbose {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Infof("verbose mode enabled")
	}
	if c.swaggerPath != "" && c.path != "" && c.readmePath != "" {
		logrus.Error("only one of 'swagger', 'path' and 'readme' can be specified")
		return 1
	}
	if c.path == "" && c.swaggerPath == "" && c.readmePath == "" {
		logrus.Error(c.Help())
		return 1
	}
	if c.readmePath != "" && c.tag == "" {
		logrus.Error("tag must be specified when 'readme' is specified")
		return 1
	}
	if c.readmePath == "" && c.tag != "" {
		logrus.Errorf("tag can only be specified when 'readme' is specified")
		return 1
	}
	return c.Execute()
}

func (c GenerateCommand) Execute() int {
	wd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Error getting working directory: %+v", err)
		return 1
	}
	if c.workingDir != "" {
		wd, err = filepath.Abs(c.workingDir)
		if err != nil {
			logrus.Errorf("Error getting absolute path of working directory: %+v", err)
			return 1
		}
	}
	c.workingDir = wd
	logrus.Infof("working directory: %s", c.workingDir)

	switch {
	case c.swaggerPath != "":
		return c.fromSwaggerPath()
	case c.path != "":
		return c.fromExamplePath()
	case c.readmePath != "":
		return c.fromAutorestConfig()
	}

	// should not reach here
	logrus.Println(c.Help())
	return 1
}

func (c GenerateCommand) fromExamplePath() int {
	wd := c.workingDir
	if c.overwrite {
		logrus.Infof("overwriting existing terraform configurations...")
		_ = os.RemoveAll(path.Join(wd, "testing.tf"))
		_ = os.RemoveAll(path.Join(wd, "dependency.tf"))
	}
	err := os.WriteFile(path.Join(wd, "provider.tf"), hclwrite.Format([]byte(resource.DefaultProviderConfig)), 0644)
	if err != nil {
		logrus.Errorf("writing provider.tf: %+v", err)
	}
	logrus.Infof("provider configuration is written to %s", path.Join(wd, "provider.tf"))

	// load example
	logrus.Infof("loading example: %s", c.path)
	example, err := resource.NewAzapiDefinitionFromExample(c.path, c.resourceType)
	if err != nil {
		logrus.Fatalf("loading example: %+v", err)
	}
	if c.useRawJsonPayload {
		logrus.Infof("using raw json payload in 'body'...")
		example.BodyFormat = types.BodyFormatJson
	}

	// load dependencies
	logrus.Infof("loading dependencies...")
	referenceResolvers := []resolver.ReferenceResolver{
		resolver.NewExistingDependencyResolver(wd),
		resolver.NewAzapiDependencyResolver(),
		resolver.NewAzurermDependencyResolver(),
		resolver.NewProviderIDResolver(),
		resolver.NewLocationIDResolver(),
		resolver.NewAzapiResourcePlaceholderResolver(),
	}
	context := resource.NewContext(referenceResolvers)
	err = context.InitFile(allTerraformConfig(wd))
	if err != nil {
		logrus.Errorf("initializing terraform configurations: %+v", err)
		return 1
	}

	logrus.Infof("generating terraform configurations...")
	err = context.AddAzapiDefinition(example)
	if err != nil {
		return 0
	}

	logrus.Infof("writing terraform configurations...")
	blockMap := blockFileMap(wd)
	contentToAppend := make(map[string]string)
	len := len(context.File.Body().Blocks())
	for i, block := range context.File.Body().Blocks() {
		switch block.Type() {
		case "terraform", "provider", "variable":
			continue
		default:
			key := fmt.Sprintf("%s.%s", block.Type(), strings.Join(block.Labels(), "."))
			if _, ok := blockMap[key]; ok {
				continue
			}
			outputFilename := "dependency.tf"
			if i == len-1 {
				outputFilename = "testing.tf"
			}
			contentToAppend[outputFilename] = contentToAppend[outputFilename] + "\n" + string(block.BuildTokens(nil).Bytes())
		}
	}

	for filename, content := range contentToAppend {
		err := appendContent(path.Join(wd, filename), content)
		if err != nil {
			logrus.Errorf("writing %s: %+v", filename, err)
		}
		logrus.Infof("configuration is written to %s", path.Join(wd, filename))
	}
	return 0
}

func (c GenerateCommand) fromSwaggerPath() int {
	swaggerPath, err := filepath.Abs(c.swaggerPath)
	if err == nil {
		c.swaggerPath = swaggerPath
	}
	logrus.Infof("loading swagger spec: %s...", c.swaggerPath)
	file, err := os.Stat(c.swaggerPath)
	if err != nil {
		logrus.Fatalf("loading swagger spec: %+v", err)
	}
	apiPathsAll := make([]swagger.ApiPath, 0)
	if file.IsDir() {
		logrus.Infof("swagger spec is a directory")
		logrus.Infof("loading swagger spec directory: %s...", c.swaggerPath)
		files, err := os.ReadDir(c.swaggerPath)
		if err != nil {
			logrus.Fatalf("reading swagger spec directory: %+v", err)
		}
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".json") || file.IsDir() {
				continue
			}
			filename := path.Join(c.swaggerPath, file.Name())
			logrus.Infof("parsing swagger spec: %s...", filename)
			apiPaths, err := swagger.Load(filename)
			if err != nil {
				logrus.Fatalf("parsing swagger spec: %+v", err)
			}
			apiPathsAll = append(apiPathsAll, apiPaths...)
		}
	} else {
		logrus.Infof("parsing swagger spec: %s...", c.swaggerPath)
		apiPaths, err := swagger.Load(c.swaggerPath)
		if err != nil {
			logrus.Fatalf("parsing swagger spec: %+v", err)
		}
		apiPathsAll = append(apiPathsAll, apiPaths...)
	}

	logrus.Infof("found %d api paths", len(apiPathsAll))
	return c.generate(apiPathsAll)
}

func (c *GenerateCommand) fromAutorestConfig() int {
	logrus.Infof("parsing autorest config: %s...", c.readmePath)
	packages := autorest.ParseAutoRestConfig(c.readmePath)
	logrus.Debugf("found %d packages", len(packages))
	var targetPackage *autorest.Package
	for _, pkg := range packages {
		if pkg.Tag == c.tag {
			targetPackage = &pkg
			break
		}
	}
	if targetPackage == nil {
		logrus.Fatalf("package with tag %s not found in %s", c.tag, c.readmePath)
	}

	apiPathsAll := make([]swagger.ApiPath, 0)
	for _, swaggerPath := range targetPackage.InputFiles {
		logrus.Infof("parsing swagger spec: %s...", swaggerPath)
		azapiPaths, err := swagger.Load(swaggerPath)
		if err != nil {
			logrus.Fatalf("parsing swagger spec: %+v", err)
		}
		apiPathsAll = append(apiPathsAll, azapiPaths...)
	}

	return c.generate(apiPathsAll)
}

func (c *GenerateCommand) generate(apiPaths []swagger.ApiPath) int {
	wd := c.workingDir
	azapiDefinitionsAll := make([]types.AzapiDefinition, 0)
	for _, apiPath := range apiPaths {
		azapiDefinitionsAll = append(azapiDefinitionsAll, resource.NewAzapiDefinitionsFromSwagger(apiPath)...)
	}

	for i := range azapiDefinitionsAll {
		if c.useRawJsonPayload {
			azapiDefinitionsAll[i].BodyFormat = types.BodyFormatJson
		}
	}

	azapiDefinitionByResourceType := make(map[string][]types.AzapiDefinition)
	for _, azapiDefinition := range azapiDefinitionsAll {
		azapiDefinitionByResourceType[azapiDefinition.AzureResourceType] = append(azapiDefinitionByResourceType[azapiDefinition.AzureResourceType], azapiDefinition)
	}

	resourceTypes := make([]string, 0)
	for resourceType := range azapiDefinitionByResourceType {
		slices.SortFunc(azapiDefinitionByResourceType[resourceType], func(i, j types.AzapiDefinition) int {
			return azapiDefinitionOrder(i) - azapiDefinitionOrder(j)
		})
		resourceTypes = append(resourceTypes, resourceType)
	}

	sort.Strings(resourceTypes)

	referenceResolvers := []resolver.ReferenceResolver{
		resolver.NewAzapiDependencyResolver(),
		resolver.NewAzapiDefinitionResolver(azapiDefinitionsAll),
		resolver.NewProviderIDResolver(),
		resolver.NewLocationIDResolver(),
		resolver.NewAzapiResourceIdResolver(),
	}

	for _, resourceType := range resourceTypes {
		logrus.Infof("generating terraform configurations for %s...", resourceType)
		azapiDefinitions := azapiDefinitionByResourceType[resourceType]
		// remove existing folders by default
		folderName := strings.ReplaceAll(resourceType, "/", "_")
		err := os.RemoveAll(path.Join(wd, folderName))
		if err != nil {
			logrus.Errorf("removing existing folder: %+v", err)
		}
		err = os.MkdirAll(path.Join(wd, folderName), 0755)
		if err != nil {
			logrus.Fatalf("creating folder: %+v", err)
		}

		context := resource.NewContext(referenceResolvers)

		for _, azapiDefinition := range azapiDefinitions {
			logrus.Debugf("generating terraform configurations for %s...", azapiDefinition.Id)
			err = context.AddAzapiDefinition(azapiDefinition)
			if err != nil {
				logrus.Warnf("adding azapi definition for %s: %+v", azapiDefinition.Id, err)
			}
		}

		filename := path.Join(wd, folderName, "main.tf")
		err = os.WriteFile(filename, hclwrite.Format([]byte(context.String())), 0644)
		if err != nil {
			logrus.Errorf("writing %s: %+v", filename, err)
		}
	}
	return 0
}

func azapiDefinitionOrder(azapiDefinition types.AzapiDefinition) int {
	switch azapiDefinition.ResourceName {
	case "azapi_resource":
		return 0
	case "azapi_update_resource":
		return 2
	case "azapi_resource_action":
		if actionField := azapiDefinition.AdditionalFields["action"]; actionField == nil || actionField.String() == "" {
			return 1
		}
		return 3
	case "azapi_resource_list":
		return 4
	}
	return 5
}

func appendContent(filename string, hclContent string) error {
	content := hclContent
	if _, err := os.Stat(filename); err == nil {
		existingHcl, err := os.ReadFile(filename)
		if err != nil {
			logrus.Warnf("reading existing file: %+v", err)
		}
		content = string(existingHcl) + "\n" + content
	}
	return os.WriteFile(filename, hclwrite.Format([]byte(content)), 0644)
}

func blockFileMap(workingDirectory string) map[string]string {
	files, err := os.ReadDir(workingDirectory)
	if err != nil {
		logrus.Warnf("reading dir %s: %+v", workingDirectory, err)
		return nil
	}
	out := make(map[string]string)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}
		src, err := os.ReadFile(path.Join(workingDirectory, file.Name()))
		if err != nil {
			logrus.Warnf("reading file %s: %+v", file.Name(), err)
			continue
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if diag.HasErrors() {
			logrus.Warnf("parsing file %s: %+v", file.Name(), diag.Error())
			continue
		}
		if f == nil || f.Body() == nil {
			continue
		}
		for _, block := range f.Body().Blocks() {
			key := fmt.Sprintf("%s.%s", block.Type(), strings.Join(block.Labels(), "."))
			out[key] = file.Name()
		}
	}

	return out
}

func allTerraformConfig(workingDirectory string) string {
	out := ""
	files, err := os.ReadDir(workingDirectory)
	if err != nil {
		logrus.Warnf("reading dir %s: %+v", workingDirectory, err)
		return out
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}
		src, err := os.ReadFile(path.Join(workingDirectory, file.Name()))
		if err != nil {
			logrus.Warnf("reading file %s: %+v", file.Name(), err)
			continue
		}
		out += string(src)
	}

	return out
}
