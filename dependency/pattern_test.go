package dependency_test

import (
	"log"
	"testing"

	"github.com/ms-henglu/armstrong/dependency"
)

func Test_NewPattern(t *testing.T) {
	testcases := []struct {
		input    string
		expected dependency.Pattern
	}{
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subnets/test",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Network/virtualNetworks/subnets",
				Scope:             dependency.ScopeResource,
			},
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Network/virtualNetworks",
				Scope:             dependency.ScopeResourceGroup,
			},
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Resources/resourceGroups",
				Scope:             dependency.ScopeSubscription,
			},
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Resources/subscriptions",
				Scope:             dependency.ScopeTenant,
			},
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subnets/test/providers/Microsoft.Network/networkSecurityGroups/test",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Network/networkSecurityGroups",
				Scope:             dependency.ScopeResource,
			},
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Network",
				Scope:             dependency.ScopeResourceGroup,
			},
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Network",
				Scope:             dependency.ScopeSubscription,
			},
		},
		{
			input: "/providers/Microsoft.Network",
			expected: dependency.Pattern{
				AzureResourceType: "Microsoft.Network",
				Scope:             dependency.ScopeTenant,
			},
		},
	}

	for _, testcase := range testcases {
		log.Printf("[DEBUG] input: %s", testcase.input)
		actual := dependency.NewPattern(testcase.input)
		if actual.AzureResourceType != testcase.expected.AzureResourceType {
			t.Errorf("expected %s, got %s", testcase.expected.AzureResourceType, actual.AzureResourceType)
		}
		if actual.Scope != testcase.expected.Scope {
			t.Errorf("expected %s, got %s", testcase.expected.Scope, actual.Scope)
		}
	}
}

func Test_IsMatch(t *testing.T) {
	testcases := []struct {
		Input   string
		Pattern dependency.Pattern
		Match   bool
	}{
		{
			Input:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subnets/test",
			Pattern: dependency.NewPattern("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subnets/test"),
			Match:   true,
		},
		{
			Input:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test",
			Pattern: dependency.NewPattern("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test"),
			Match:   true,
		},
		{
			Input:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test",
			Pattern: dependency.NewPattern("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subntes/test"),
			Match:   false,
		},
	}

	for _, testcase := range testcases {
		t.Logf("input: %s", testcase.Input)
		actual := testcase.Pattern.IsMatch(testcase.Input)
		if actual != testcase.Match {
			t.Errorf("expected %v, got %v", testcase.Match, actual)
		}
	}
}

func Test_String(t *testing.T) {
	testcases := []struct {
		Input  string
		Expect string
	}{
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subnets/test",
			Expect: "resource:microsoft.network/virtualnetworks/subnets",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test",
			Expect: "resource_group:microsoft.network/virtualnetworks",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test",
			Expect: "subscription:microsoft.resources/resourcegroups",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expect: "tenant:microsoft.resources/subscriptions",
		},
	}

	for _, testcase := range testcases {
		t.Logf("input: %s", testcase.Input)
		actual := dependency.NewPattern(testcase.Input).String()
		if actual != testcase.Expect {
			t.Errorf("expected %s, got %s", testcase.Expect, actual)
		}
	}
}
