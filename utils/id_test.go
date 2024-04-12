package utils_test

import (
	"testing"

	"github.com/azure/armstrong/utils"
)

func Test_IsResourceId(t *testing.T) {
	testcases := []struct {
		Input  string
		Expect bool
	}{
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expect: true,
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
			Expect: true,
		},
		{
			Input:  "/",
			Expect: false,
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb/providers/Microsoft.Resources/locks/myLock",
			Expect: true,
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb",
			Expect: true,
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb/providers/Microsoft.Resources/locks/myLock/foo",
			Expect: false,
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb/foo",
			Expect: false,
		},
	}

	for _, testcase := range testcases {
		t.Logf("[DEBUG] testcase: %s", testcase.Input)
		actual := utils.IsResourceId(testcase.Input)
		if actual != testcase.Expect {
			t.Fatalf("[ERROR] expect %v, actual %v", testcase.Expect, actual)
		}
	}
}

func Test_ResourceTypeOfResourceId(t *testing.T) {
	testcases := []struct {
		Input  string
		Expect string
	}{
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expect: "Microsoft.Resources/subscriptions",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
			Expect: "Microsoft.Resources/resourceGroups",
		},
		{
			Input:  "/",
			Expect: "Microsoft.Resources/tenants",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb/providers/Microsoft.Resources/locks/myLock",
			Expect: "Microsoft.Resources/locks",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb",
			Expect: "Microsoft.Automation/automationAccounts/runbooks",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa",
			Expect: "Microsoft.Automation/automationAccounts",
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts",
			// though it would be better to return Microsoft.Automation/automationAccounts, but there's no way to know it's a list API or a provider action, see below case
			Expect: "Microsoft.Automation",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/checkNameAvailability",
			Expect: "Microsoft.Automation",
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Automation/automationAccounts",
			// though it would be better to return Microsoft.Automation/automationAccounts, but there's no way to know it's a list API or a provider action, see below case
			Expect: "Microsoft.Automation",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Automation/checkNameAvailability",
			Expect: "Microsoft.Automation",
		},
	}

	for _, testcase := range testcases {
		t.Logf("[DEBUG] testcase: %s", testcase.Input)
		actual := utils.ResourceTypeOfResourceId(testcase.Input)
		if actual != testcase.Expect {
			t.Fatalf("[ERROR] expect %s, actual %s", testcase.Expect, actual)
		}
	}
}

func Test_PrentIdOfResourceId(t *testing.T) {
	testcases := []struct {
		Input  string
		Expect string
	}{
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expect: "/",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
			Expect: "/subscriptions/00000000-0000-0000-0000-000000000000",
		},
		{
			Input:  "/",
			Expect: "",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb/providers/Microsoft.Resources/locks/myLock",
			Expect: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa/runbooks/rb",
			Expect: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation/automationAccounts/aa",
			Expect: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Automation",
			Expect: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg",
		},
		{
			Input:  "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Automation",
			Expect: "/subscriptions/00000000-0000-0000-0000-000000000000",
		},
		{
			Input:  "/providers/Microsoft.Automation",
			Expect: "/",
		},
	}

	for _, testcase := range testcases {
		t.Logf("[DEBUG] testcase: %s", testcase.Input)
		actual := utils.ParentIdOfResourceId(testcase.Input)
		if actual != testcase.Expect {
			t.Fatalf("[ERROR] expect %s, actual %s", testcase.Expect, actual)
		}
	}
}
