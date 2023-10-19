package dependency_test

import (
	"testing"

	"github.com/ms-henglu/armstrong/dependency"
)

func Test_LoadAzapiDependencies(t *testing.T) {
	res, err := dependency.LoadAzapiDependencies()
	if err != nil {
		t.Error(err)
	}
	if len(res) == 0 {
		t.Error("No dependencies loaded")
	}

	for _, dep := range res {
		if dep.AzureResourceType == "" {
			t.Errorf("AzureResourceType is empty for %s", dep.ResourceName)
		}
		if dep.ApiVersion == "" {
			t.Errorf("ApiVersion is empty for %s", dep.ResourceName)
		}
		if dep.ExampleConfiguration == "" {
			t.Errorf("ExampleConfiguration is empty for %s", dep.ResourceName)
		}
		if dep.ResourceKind == "" {
			t.Errorf("ResourceKind is empty for %s", dep.ResourceName)
		}
		if dep.ResourceName == "" {
			t.Errorf("ResourceName is empty for %s", dep.ResourceName)
		}
		if dep.ResourceLabel == "" {
			t.Errorf("ResourceLabel is empty for %s", dep.ResourceName)
		}

	}
}
