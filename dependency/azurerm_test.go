package dependency_test

import (
	"testing"

	"github.com/azure/armstrong/dependency"
)

func Test_LoadAzurermDependencies(t *testing.T) {
	res := dependency.LoadAzurermDependencies()
	if len(res) == 0 {
		t.Error("No dependencies loaded")
	}
	for _, dep := range res {
		if dep.AzureResourceType == "" {
			t.Errorf("AzureResourceType is empty for %s", dep.ResourceName)
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
	}
}
