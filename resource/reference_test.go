package resource_test

import (
	"testing"

	"github.com/ms-henglu/armstrong/resource"
)

func TestReference_NewReferenceFromAddress(t *testing.T) {
	testcases := []struct {
		address  string
		expected *resource.Reference
	}{
		{
			address: "data.azapi_resource.test.id",
			expected: &resource.Reference{
				Type:         "data",
				ResourceType: "azapi_resource",
				Label:        "test",
				PropertyName: "id",
			},
		},
		{
			address: "azapi_resource.test.id",
			expected: &resource.Reference{
				Type:         "resource",
				ResourceType: "azapi_resource",
				Label:        "test",
				PropertyName: "id",
			},
		},
		{
			address:  "azapi_resource.test",
			expected: nil,
		},
		{
			address:  "azapi_resource.test.id.name.value",
			expected: nil,
		},
	}

	for _, tc := range testcases {
		t.Logf("testing address: %s", tc.address)
		got := resource.NewReferenceFromAddress(tc.address)
		if tc.expected == nil || got == nil {
			if tc.expected != got {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
			continue
		}
		if got.Type != tc.expected.Type {
			t.Fatalf("expected %s, got %s", tc.expected.Type, got.Type)
		}
		if got.ResourceType != tc.expected.ResourceType {
			t.Fatalf("expected %s, got %s", tc.expected.ResourceType, got.ResourceType)
		}
		if got.Label != tc.expected.Label {
			t.Fatalf("expected %s, got %s", tc.expected.Label, got.Label)
		}
		if got.PropertyName != tc.expected.PropertyName {
			t.Fatalf("expected %s, got %s", tc.expected.PropertyName, got.PropertyName)
		}
	}
}

func TestReference_IsKnown(t *testing.T) {
	testcases := []struct {
		ref      *resource.Reference
		expected bool
	}{
		{
			ref: &resource.Reference{
				Type:         "data",
				ResourceType: "azapi_resource",
				Label:        "test",
				PropertyName: "id",
			},
			expected: true,
		},
		{
			ref: &resource.Reference{
				Type:         "resource",
				ResourceType: "azapi_resource",
				Label:        "test",
				PropertyName: "id",
			},
			expected: true,
		},
		{
			ref: &resource.Reference{
				Type:         "resource",
				ResourceType: "azapi_resource",
				Label:        "",
				PropertyName: "id",
			},
			expected: false,
		},
		{
			ref: &resource.Reference{
				Type:         "resource",
				ResourceType: "azapi_resource",
				Label:        "test",
				PropertyName: "",
			},
			expected: false,
		},
		{
			ref:      nil,
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Logf("testing reference: %v", tc.ref)
		got := tc.ref.IsKnown()
		if got != tc.expected {
			t.Fatalf("expected %v, got %v", tc.expected, got)
		}
	}
}