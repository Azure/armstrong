package commands

import (
	"strings"
	"testing"

	"github.com/azure/armstrong/hcl"
)

func TestCheckAzureProviderSecretDoesNotExposeCredentialInput(t *testing.T) {
	target := credScanTarget{
		FileName:   "main.tf",
		LineNumber: 10,
		Name:       "azurerm",
		Type:       "provider",
	}

	testCases := []struct {
		name          string
		value         string
		variables     map[string]hcl.Variable
		expectedCount int
		forbidden     string
	}{
		{
			name:          "literal secret",
			value:         "super-secret-password",
			expectedCount: 1,
			forbidden:     "super-secret-password",
		},
		{
			name:          "unknown variable",
			value:         "$var.secret_variable_from_input",
			expectedCount: 1,
			forbidden:     "secret_variable_from_input",
		},
		{
			name:  "declared variable",
			value: "$var.certificate_password",
			variables: map[string]hcl.Variable{
				"certificate_password": {
					Name:       "certificate_password",
					HasDefault: true,
					FileName:   "variables.tf",
					LineNumber: 3,
				},
			},
			expectedCount: 2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			errors := checkAzureProviderSecret(target, "client_certificate_password", testCase.value, testCase.variables)
			if len(errors) != testCase.expectedCount {
				t.Fatalf("expected %d errors, got %d", testCase.expectedCount, len(errors))
			}
			for _, scanError := range errors {
				if testCase.forbidden != "" && strings.Contains(scanError.Error(), testCase.forbidden) {
					t.Errorf("error exposes credential input %q: %s", testCase.forbidden, scanError.Error())
				}
			}
		})
	}
}
