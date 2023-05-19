package coverage

import (
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	apiPath, modelName, modelSwaggerPath, err := PathPatternFromIdFromIndex("/subscriptions/85b3dbca-5974-4067-9669-67a141095a76/resourceGroups/wangta-ex3-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR", "2022-06-01")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*apiPath, *modelName, *modelSwaggerPath)
}
