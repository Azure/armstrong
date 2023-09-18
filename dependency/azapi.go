package dependency

import (
	"embed"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Dependency struct {
	AzureResourceType    string
	ApiVersion           string
	ExampleConfiguration string
	ResourceKind         string // resource or data
	ReferredProperty     string // only supports "id" for now
	ResourceName         string
	ResourceLabel        string
}

//go:embed azapi_examples
var StaticFiles embed.FS

var azapiMutex = sync.Mutex{}

var azapiDeps = make([]Dependency, 0)

func LoadAzapiDependencies() ([]Dependency, error) {
	azapiMutex.Lock()
	defer azapiMutex.Unlock()
	if len(azapiDeps) != 0 {
		return azapiDeps, nil
	}

	dir := "azapi_examples"
	entries, err := StaticFiles.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		filename := path.Join(dir, entry.Name(), "main.tf")
		data, err := StaticFiles.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		f, diags := hclwrite.ParseConfig(data, filename, hcl.InitialPos)
		if diags.HasErrors() {
			return nil, diags
		}
		blockTotal := len(f.Body().Blocks())
		lastBlock := f.Body().Blocks()[blockTotal-1]
		typeValue := string(lastBlock.Body().GetAttribute("type").Expr().BuildTokens(nil).Bytes())
		typeValue = strings.Trim(typeValue, ` "`)
		parts := strings.Split(typeValue, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("resource type is invalid: %s, filename: %s", typeValue, filename)
		}
		dep := Dependency{
			AzureResourceType:    parts[0],
			ApiVersion:           parts[1],
			ExampleConfiguration: string(data),
			ReferredProperty:     "id",
			ResourceKind:         lastBlock.Type(),
			ResourceName:         lastBlock.Labels()[0],
			ResourceLabel:        lastBlock.Labels()[1],
		}
		azapiDeps = append(azapiDeps, dep)
	}

	// add a special case for Microsoft.Resources/subscriptions
	azapiDeps = append(azapiDeps, Dependency{
		AzureResourceType: "Microsoft.Resources/subscriptions",
		ApiVersion:        "2020-06-01",
		ReferredProperty:  "id",
		ResourceKind:      "data",
		ResourceName:      "azapi_resource",
		ResourceLabel:     "subscription",
		ExampleConfiguration: `
data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2020-06-01"
  response_export_values = ["*"]
}
`,
	})
	return azapiDeps, nil
}
