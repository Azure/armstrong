package coverage

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	openapispec "github.com/go-openapi/spec"
)

type Model struct {
	Bool                    *map[string]bool   `json:"Bool,omitempty"`
	Discriminator           *string            `json:"Discriminator,omitempty"`
	Enum                    *map[string]bool   `json:"Enum,omitempty"`
	Format                  *string            `json:"Format,omitempty"`
	Identifier              string             `json:"Identifier,omitempty"`
	IsAnyCovered            bool               `json:"IsAnyCovered"`
	IsFullyCovered          bool               `json:"IsFullyCovered,omitempty"`
	HasAdditionalProperties bool               `json:"HasAdditionalProperties,omitempty"`
	CoveredCount            int                `json:"CoveredCount,omitempty"`
	TotalCount              int                `json:"TotalCount,omitempty"`
	IsReadOnly              bool               `json:"IsReadOnly,omitempty"`
	IsRequired              bool               `json:"IsRequired,omitempty"`
	Item                    *Model             `json:"Item,omitempty"`
	Properties              *map[string]*Model `json:"Properties,omitempty"`
	Type                    *string            `json:"Type,omitempty"`
	Variants                *map[string]*Model `json:"Variants,omitempty"`
}

func Expand(modelName, swaggerPath string) (*Model, error) {
	doc, err := loads.JSONSpec(swaggerPath)
	if err != nil {
		return nil, err
	}

	if modelName == "" {
		return nil, nil
	}

	spec := doc.Spec()

	modelSchema, ok := spec.Definitions[modelName]
	if !ok {
		return nil, fmt.Errorf("%s not found in the definition of %s", modelName, swaggerPath)
	}

	variantsTable := map[string][]string{}
	for k, v := range spec.Definitions {
		if v.Extensions["x-ms-discriminator-value"] != nil && len(v.AllOf) > 0 {
			for _, variant := range v.AllOf {
				if variant.Ref.String() != "" {
					resolved, err := openapispec.ResolveRefWithBase(spec, &variant.Ref, &openapispec.ExpandOptions{RelativeBase: swaggerPath})
					if err != nil {
						panic(err)
					}
					if resolved.Extensions["x-ms-discriminator-value"] != nil || resolved.Discriminator != "" {
						modelName, _ := SchemaInfoFromRef(variant.Ref)
						if variantsTable[modelName] == nil {
							variantsTable[modelName] = []string{k}
						} else {
							variantsTable[modelName] = append(variantsTable[modelName], k)
						}
					}
				}
			}
		}
	}

	output := expandSchema(modelSchema, swaggerPath, modelName, "#", spec, variantsTable, map[string]interface{}{}, map[string]interface{}{})

	return output, nil
}

func expandSchema(input openapispec.Schema, swaggerPath, modelName, identifier string, root interface{}, variantsTable map[string][]string, resolvedDiscriminator map[string]interface{}, resolvedModel map[string]interface{}) *Model {
	output := Model{Identifier: identifier}

	//fmt.Println("expand schema for", swaggerPath, modelName)
	if _, hasResolvedModel := resolvedModel[modelName]; hasResolvedModel {
		log.Printf("[WARN]circular reference detected for %s %s", swaggerPath, modelName)
		return &output
	}
	resolvedModel[modelName] = nil

	if len(input.Type) > 0 {
		output.Type = &input.Type[0]
		if *output.Type == "boolean" {
			boolMap := make(map[string]bool)
			boolMap["true"] = false
			boolMap["false"] = false

			output.Bool = &boolMap
		}
	}

	if input.AdditionalProperties != nil {
		output.HasAdditionalProperties = true
	}

	if input.Format != "" {
		output.Format = &input.Format
	}

	if input.ReadOnly {
		output.IsReadOnly = input.ReadOnly
	}

	if input.Enum != nil {
		enumMap := make(map[string]bool)
		for _, v := range input.Enum {
			switch t := v.(type) {
			case string:
				enumMap[t] = false
			case float64:
				enumMap[fmt.Sprintf("%v", t)] = false
			case int:
				enumMap[fmt.Sprintf("%v", t)] = false
			default:
				panic(fmt.Sprintf("unknown enum type %T", t))
			}
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
			swaggerPath = strings.Replace(swaggerPath, "https:/", "https://", 1)

			doc, err := loads.JSONSpec(swaggerPath)
			if err != nil {
				panic(err)
			}

			root = doc.Spec()
		}

		output = *expandSchema(*resolved, swaggerPath, modelName, identifier, root, variantsTable, resolvedDiscriminator, resolvedModel)
	}

	// expand properties
	properties := make(map[string]*Model)
	for k, v := range input.Properties {
		//fmt.Println("expand properties", k)
		properties[k] = expandSchema(v, swaggerPath, fmt.Sprintf("%s.%s", modelName, k), identifier+"."+k, root, variantsTable, resolvedDiscriminator, resolvedModel)
	}

	// expand composition
	for _, v := range input.AllOf {
		//fmt.Println("expand composition", v.Ref.String())
		allOf := expandSchema(v, swaggerPath, fmt.Sprintf("%s.allOf", modelName), identifier, root, variantsTable, resolvedDiscriminator, resolvedModel)
		if allOf.Properties != nil {
			for k, v := range *allOf.Properties {
				properties[k] = v
			}
		}
	}

	if len(properties) > 0 {
		for _, v := range input.Required {
			p, ok := properties[v]
			if !ok {
				panic(fmt.Sprintf("required property %s not found in %s", v, modelName))
			}
			p.IsRequired = true
		}
		output.Properties = &properties
	}

	// expand items
	if input.Items != nil {
		//fmt.Println("expand items", input.Items.Schema.Ref.String())
		item := expandSchema(*input.Items.Schema, swaggerPath, fmt.Sprintf("%s[]", modelName), identifier+"[]", root, variantsTable, resolvedDiscriminator, resolvedModel)
		output.Item = item
	}

	delete(resolvedModel, modelName)

	// variants have circular reference
	// expand variants
	if input.Discriminator != "" {
		_, hasResolvedDiscriminator := resolvedDiscriminator[modelName]
		if !hasResolvedDiscriminator {
			resolvedDiscriminator[modelName] = nil
			//fmt.Println("expand variants", modelName)
			variants := make(map[string]*Model)

			vars, ok := variantsTable[modelName]
			// level order traverse to find all variants
			for ok && len(vars) > 0 {
				vars2 := []string{}
				for _, v := range vars {
					schema := root.(*openapispec.Swagger).Definitions[v]
					variantName := schema.Extensions["x-ms-discriminator-value"].(string)
					resolved := expandSchema(schema, swaggerPath, v, identifier+"{"+variantName+"}", root, variantsTable, resolvedDiscriminator, resolvedModel)
					variants[variantName] = resolved
					if vv, ok := variantsTable[v]; ok {
						vars2 = append(vars2, vv...)
					}
				}
				vars = vars2
			}
			output.Discriminator = &input.Discriminator
			output.Variants = &variants
		}
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
				if operation == nil {
					operation = item.Post
				}
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
