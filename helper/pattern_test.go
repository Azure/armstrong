package helper_test

import (
	"testing"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/helper"
)

func Test_IsValueMatchPattern(t *testing.T) {
	testcases := []struct {
		Value   string
		Pattern string
		Expect  bool
	}{
		{
			Value:   "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123/computes/compute123",
			Pattern: "/subscriptions/resourceGroups/providers/Microsoft.MachineLearningServices/workspaces/computes",
			Expect:  true,
		},
		{
			Value:   "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123",
			Pattern: "/subscriptions/resourceGroups/providers/Microsoft.MachineLearningServices/workspaces",
			Expect:  true,
		},
	}

	for _, testcase := range testcases {
		if output := helper.IsValueMatchPattern(testcase.Value, testcase.Pattern); output != testcase.Expect {
			t.Fatalf("expect %v but got %v", testcase.Expect, output)
		}
	}

}

func Test_GetIdPattern(t *testing.T) {
	testcases := []struct {
		Value  string
		Expect string
	}{
		{
			Value:  "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123/computes/compute123",
			Expect: "/subscriptions/resourceGroups/providers/Microsoft.MachineLearningServices/workspaces/computes",
		},
		{
			Value:  "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123",
			Expect: "/subscriptions/resourceGroups/providers/Microsoft.MachineLearningServices/workspaces",
		},
	}

	for _, testcase := range testcases {
		if output, err := helper.GetIdPattern(testcase.Value); err != nil || output != testcase.Expect {
			t.Fatalf("expect %v but got %v, err: %+v", testcase.Expect, output, err)
		}
	}

}

func Test_GetResourceType(t *testing.T) {
	testcases := []struct {
		Value  string
		Expect string
	}{
		{
			Value:  "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123/computes/compute123",
			Expect: "Microsoft.MachineLearningServices/workspaces/computes",
		},
		{
			Value:  "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123",
			Expect: "Microsoft.MachineLearningServices/workspaces",
		},
	}

	for _, testcase := range testcases {
		if output := helper.GetResourceType(testcase.Value); output != testcase.Expect {
			t.Fatalf("expect %v but got %v", testcase.Expect, output)
		}
	}

}
