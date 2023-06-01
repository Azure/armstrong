package coverage

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	openapispec "github.com/go-openapi/spec"
)

type Model struct {
	Bool           *map[bool]bool     `json:"Bool,omitempty"`
	Discriminator  *string            `json:"Discriminator,omitempty"`
	Enum           *map[string]bool   `json:"Enum,omitempty"`
	Format         *string            `json:"Format,omitempty"`
	Identifier     string             `json:"Identifier,omitempty"`
	IsAnyCovered   bool               `json:"IsAnyCovered"`
	IsFullyCovered bool               `json:"IsFullyCovered,omitempty"`
	CoveredCount   int                `json:"CoveredCount,omitempty"`
	TotalCount     int                `json:"TotalCount,omitempty"`
	IsReadOnly     bool               `json:"IsReadOnly,omitempty"`
	IsRequired     bool               `json:"IsRequired,omitempty"`
	Item           *Model             `json:"Item,omitempty"`
	Properties     *map[string]*Model `json:"Properties,omitempty"`
	Type           *string            `json:"Type,omitempty"`
	Variants       *map[string]*Model `json:"Variants,omitempty"`
}

func Expand(modelName, swaggerPath string) (*Model, error) {
	doc, err := loads.JSONSpec(swaggerPath)
	if err != nil {
		return nil, err
	}

	spec := doc.Spec()

	modelSchema, ok := spec.Definitions[modelName]
	if !ok {
		return nil, fmt.Errorf("%s not found in the definition of %s", modelName, swaggerPath)
	}

	allOfs := map[string][]string{}
	for k, v := range spec.Definitions {
		if v.VendorExtensible.Extensions["x-ms-discriminator-value"] != nil && len(v.SchemaProps.AllOf) > 0 {
			for _, v2 := range v.SchemaProps.AllOf {
				if v2.Ref.String() != "" {
					resolved, err := openapispec.ResolveRefWithBase(spec, &v2.Ref, &openapispec.ExpandOptions{RelativeBase: swaggerPath})
					if err != nil {
						panic(err)
					}
					if resolved.VendorExtensible.Extensions["x-ms-discriminator-value"] != nil || resolved.SwaggerSchemaProps.Discriminator != "" {
						modelName, _ := SchemaInfoFromRef(v2.Ref)
						if allOfs[modelName] == nil {
							allOfs[modelName] = []string{k}
						} else {
							allOfs[modelName] = append(allOfs[modelName], k)
						}
					}
				}
			}
		}
	}

	output := expandSchema(modelSchema, swaggerPath, modelName, "#", spec, allOfs, map[string]interface{}{})

	return output, nil
}

func expandSchema(input openapispec.Schema, swaggerPath, modelName, identifier string, root interface{}, allOfs map[string][]string, resolvedDiscriminator map[string]interface{}) *Model {
	output := Model{Identifier: identifier}

	//fmt.Println("expand schema for", swaggerPath, modelName)

	if len(input.SchemaProps.Type) > 0 {
		output.Type = &input.SchemaProps.Type[0]
		if *output.Type == "bool" {
			boolMap := make(map[bool]bool)
			boolMap[true] = false
			boolMap[false] = false

			output.Bool = &boolMap
		}
	}
	if input.SchemaProps.Format != "" {
		output.Format = &input.SchemaProps.Format
	}
	if input.SwaggerSchemaProps.ReadOnly {
		output.IsReadOnly = input.SwaggerSchemaProps.ReadOnly
	}
	if input.SchemaProps.Enum != nil {
		enumMap := make(map[string]bool)
		for _, v := range input.SchemaProps.Enum {
			enumMap[v.(string)] = false
		}

		output.Enum = &enumMap
	}

	// expand ref
	if input.Ref.String() != "" {
		//fmt.Println("expand ref", input.Ref.String())
		resolved, err := openapispec.ResolveRefWithBase(root, &input.Ref, &openapispec.ExpandOptions{RelativeBase: swaggerPath})
		if err != nil {
			panic(err)
		}

		modelName, relativPath := SchemaInfoFromRef(input.Ref)
		if relativPath != "" {
			swaggerPath = filepath.Join(filepath.Dir(swaggerPath), relativPath)
			swaggerPath = strings.Replace(swaggerPath, ":/", "://", 1)

			doc, err := loads.JSONSpec(swaggerPath)
			if err != nil {
				panic(err)
			}

			root = doc.Spec()
		}

		return expandSchema(*resolved, swaggerPath, modelName, identifier, root, allOfs, resolvedDiscriminator)
	}

	// expand variants
	if input.SwaggerSchemaProps.Discriminator != "" {
		_, hasResolvedDiscriminator := resolvedDiscriminator[modelName]
		if !hasResolvedDiscriminator {
			resolvedDiscriminator[modelName] = nil
			//fmt.Println("expand variants", modelName)
			variants := make(map[string]*Model)

			vars, ok := allOfs[modelName]
			for ok && len(vars) > 0 {
				vars2 := []string{}
				for _, v2 := range vars {
					schema2 := root.(*openapispec.Swagger).Definitions[v2]
					variantName := schema2.VendorExtensible.Extensions["x-ms-discriminator-value"].(string)
					resolved := expandSchema(schema2, swaggerPath, v2, identifier+"{"+variantName+"}", root, allOfs, resolvedDiscriminator)
					variants[variantName] = resolved
					if vv, ok := allOfs[v2]; ok {
						vars2 = append(vars2, vv...)
					}
				}
				vars = vars2
			}
			output.Discriminator = &input.SwaggerSchemaProps.Discriminator
			output.Variants = &variants
		}
	}

	// expand properties
	properties := make(map[string]*Model)
	for k, v := range input.Properties {
		//fmt.Println("expand properties", k)
		properties[k] = expandSchema(v, swaggerPath, fmt.Sprintf("%s.%s", modelName, k), identifier+"."+k, root, allOfs, resolvedDiscriminator)
	}

	// expand composition
	for _, v := range input.AllOf {
		//fmt.Println("expand composition", v.Ref.String())
		allOf := expandSchema(v, swaggerPath, fmt.Sprintf("%s.allOf", modelName), identifier, root, allOfs, resolvedDiscriminator)
		if allOf.Properties != nil {
			for k, v := range *allOf.Properties {
				properties[k] = v
			}
		}
	}

	if len(properties) > 0 {
		for _, v := range input.SchemaProps.Required {
			p := properties[v]
			p.IsRequired = true
		}
		output.Properties = &properties
	}

	// expand items
	if input.Items != nil {
		//fmt.Println("expand items", input.Items.Schema.Ref.String())
		item := expandSchema(*input.Items.Schema, swaggerPath, fmt.Sprintf("%s[]", modelName), identifier+"[]", root, allOfs, resolvedDiscriminator)
		output.Item = item
	}

	return &output
}

func Flatten(property *Model, path string, lookupTable map[string]bool, discriminatorTable map[string]string) {
	if property.IsReadOnly {
		return
	}

	lookupTable[path] = false

	if property.Properties != nil {
		for k, v := range *property.Properties {
			if strings.Contains(k, ".") {
				k = "\"" + k + "\""
			}
			Flatten(v, strings.TrimLeft(path+"."+k, "."), lookupTable, discriminatorTable)
		}
	}

	if property.Variants != nil && property.Discriminator != nil {
		discriminatorTable[path] = *property.Discriminator
		for k, v := range *property.Variants {
			Flatten(v, path+"{"+*property.Discriminator+"("+k+")}", lookupTable, discriminatorTable)
		}
	}

	if property.Item != nil {
		Flatten(property.Item, strings.TrimLeft(path+"[]", "."), lookupTable, discriminatorTable)
	}

	if property.Enum != nil {
		lookupTable[path+"()"] = false
		for k, v := range *property.Enum {
			lookupTable[path+"("+k+")"] = v
		}
	}

	if property.Type != nil && *property.Type == "boolean" {
		lookupTable[path+"()"] = false
		lookupTable[path+"(true)"] = false
		lookupTable[path+"(false)"] = false
	}

	return
}

func SchemaInfoFromRef(ref openapispec.Ref) (name string, path string) {
	if ref.GetURL() == nil {
		return "", ""
	}
	fragments := strings.Split(ref.GetURL().Fragment, "/")
	return fragments[len(fragments)-1], ref.GetURL().Path
}

func toRegex(path string) string {
	segs := strings.Split(path, "/")
	out := make([]string, 0, len(segs))
	for _, seg := range segs {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			out = append(out, ".+")
			continue
		}
		out = append(out, seg)
	}
	return "^" + strings.Join(out, "/") + "$"
}

func PathPatternFromId(resourceId, swaggerPath string) (*string, *string, *string, error) {
	doc, err := loads.JSONSpec(swaggerPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading swagger spec: %+v", err)
	}

	spec := doc.Spec()

	var apiPath, modelName string

	if spec.Paths != nil {
	pathLoop:
		for p, item := range spec.Paths.Paths {
			regex := toRegex(p)
			re := regexp.MustCompile(regex)
			if re.MatchString(resourceId) {
				apiPath = p
				operation := item.Put
				for _, param := range operation.Parameters {
					if param.In == "body" {
						var modelRelativePath string
						modelName, modelRelativePath = SchemaInfoFromRef(param.Schema.Ref)
						if modelRelativePath != "" {
							//fmt.Println("modelRelativePath", modelRelativePath)
							swaggerPath = filepath.Join(filepath.Dir(swaggerPath), modelRelativePath)
						}

						break pathLoop
					}
				}
			}
		}
	}
	strings.Replace(swaggerPath, "https:/", "https://", 1)

	return &apiPath, &modelName, &swaggerPath, nil
}
