package coverage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openapispec "github.com/go-openapi/spec"
	"github.com/ms-henglu/armstrong/utils"
	"github.com/sirupsen/logrus"
)

func GetModelInfoFromLocalDir(resourceId, apiVersion, swaggerPath string) (*SwaggerModel, error) {
	swaggerPath, err := filepath.Abs(swaggerPath)
	if err != nil {
		return nil, err
	}
	file, err := os.Stat(swaggerPath)
	if err != nil {
		return nil, err
	}
	if !file.IsDir() {
		return GetModelInfoFromLocalSpecFile(resourceId, apiVersion, swaggerPath)
	}
	files, err := utils.ListFiles(swaggerPath, ".json", 1)
	if err != nil {
		return nil, err
	}
	for _, filename := range files {
		model, err := GetModelInfoFromLocalSpecFile(resourceId, apiVersion, filename)
		if err != nil {
			logrus.Warnf("failed to get model info from local spec file %v: %+v", filename, err)
		}
		if model != nil {
			return model, nil
		}
	}
	return nil, nil
}

func GetModelInfoFromLocalSpecFile(resourceId, apiVersion, swaggerPath string) (*SwaggerModel, error) {
	doc, err := loadSwagger(swaggerPath)
	if err != nil {
		return nil, err
	}

	paths := doc.Spec().Paths
	if paths == nil {
		return nil, fmt.Errorf("paths is nil, swagger path: %v", swaggerPath)
	}

	for pathKey, pathItem := range paths.Paths {
		if !isPathKeyMatchWithResourceId(pathKey, resourceId) {
			continue
		}

		operation := pathItem.Put
		if operation == nil {
			// should not happen
			logrus.Warnf("no PUT operation found for path %v", pathKey)
			continue
		}

		var modelName string
		for _, param := range operation.Parameters {
			paramRef := param.Ref
			if paramRef.String() != "" {
				refParam, err := openapispec.ResolveParameterWithBase(nil, param.Ref, &openapispec.ExpandOptions{RelativeBase: swaggerPath})
				if err != nil {
					return nil, fmt.Errorf("resolve param ref %q: %+v", param.Ref.String(), err)
				}

				// Update the param
				param = *refParam
			}
			if param.In == "body" {
				if paramRef.String() != "" {
					modelName, swaggerPath = SchemaNamePathFromRef(swaggerPath, paramRef)
				}

				if param.Schema.Ref.String() != "" {
					modelName, swaggerPath = SchemaNamePathFromRef(swaggerPath, param.Schema.Ref)
				}
				break
			}
		}

		return &SwaggerModel{
			ApiPath:     pathKey,
			ModelName:   modelName,
			SwaggerPath: swaggerPath,
		}, nil
	}
	return nil, nil
}

func isPathKeyMatchWithResourceId(pathKey, resourceId string) bool {
	pathParts := strings.Split(strings.Trim(pathKey, "/"), "/")
	resourceIdParts := strings.Split(strings.Trim(resourceId, "/"), "/")
	i := len(pathParts) - 1
	j := len(resourceIdParts) - 1
	for i >= 0 && j >= 0 {
		if i == 0 && (strings.EqualFold(pathParts[i], "{resourceid}") || strings.EqualFold(pathParts[i], "{scope}") || strings.EqualFold(pathParts[i], "{resourceUri}")) {
			return true
		}
		if strings.EqualFold(pathParts[i], resourceIdParts[j]) ||
			strings.HasPrefix(pathParts[i], "{") && strings.HasSuffix(pathParts[i], "}") {
			i--
			j--
		} else {
			break
		}
	}

	return i < 0 && j < 0
}
