package coverage_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ms-henglu/armstrong/coverage"
)

func TestLoad(t *testing.T) {
	resourceId := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ex-resources/providers/Microsoft.Media/mediaServices/mediatest/transforms/transform1"
	swaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"

	apiPath, modelName, modelSwaggerPath, err := coverage.PathPatternFromId(resourceId, swaggerPath)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*apiPath, *modelName, *modelSwaggerPath)

	model, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("model", string(out))

	//lookupTable := map[string]bool{}
	//discriminatorTable := map[string]string{}
	//coverage.Flatten(*model, "", lookupTable, discriminatorTable)
	//
	//out, err = json.MarshalIndent(lookupTable, "", "\t")
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println("lookupTable", string(out))
	//
	//out, err = json.MarshalIndent(discriminatorTable, "", "\t")
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println("discriminatorTable", string(out))
}
