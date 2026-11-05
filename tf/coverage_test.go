package tf_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azure/armstrong/coverage"
	"github.com/azure/armstrong/tf"
	tfjson "github.com/hashicorp/terraform-json"
	paltypes "github.com/ms-henglu/pal/types"
	"github.com/sirupsen/logrus"
)

const (
	resourceType = "Microsoft.Automation/automationAccounts@2022-08-08"
	resourceID   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Automation/automationAccounts/test"
)

var swaggerPath = filepath.Join("..", "coverage", "testdata", "Microsoft.Automation", "stable", "2022-08-08", "account.json")

// assertCovered checks whether the property at path, such as
// "#.properties.sku.name", is covered in markdown.
func assertCovered(t *testing.T, markdown, path string, covered bool) {
	t.Helper()
	style := ""
	if !covered {
		style = ` style="color:red"`
	}
	label := path[strings.LastIndex(path, ".")+1:]
	if !strings.Contains(markdown, "<!-- "+path+" -->\n<details>\n<summary><span"+style+">"+label+"</span>") {
		t.Errorf("expected %s covered=%v", path, covered)
	}
}

func Test_NewCoverageReports(t *testing.T) {

	newResource := func(kind, name string, body interface{}) *tfjson.StateResource {
		return &tfjson.StateResource{
			Address: kind + "." + name,
			Type:    kind,
			Mode:    tfjson.ManagedResourceMode,
			AttributeValues: map[string]interface{}{
				"id":   resourceID,
				"type": resourceType,
				"body": body,
			},
		}
	}
	updateBody := map[string]interface{}{
		"properties": map[string]interface{}{"disableLocalAuth": true},
	}
	create := newResource("azapi_resource", "create", map[string]interface{}{})
	create.AttributeValues["location"] = "westus"
	nullBody := newResource("azapi_resource", "null", "null")
	nullBody.AttributeValues["location"] = "westus"
	update := newResource("azapi_update_resource", "update", updateBody)
	secondUpdate := newResource("azapi_update_resource", "second", map[string]interface{}{
		"properties": map[string]interface{}{"publicNetworkAccess": false},
	})
	data := newResource("azapi_resource", "data", updateBody)
	data.Mode = tfjson.DataResourceMode
	action := newResource("azapi_resource_action", "action", updateBody)

	testcases := []struct {
		name      string
		resources []*tfjson.StateResource
		covered   []string
	}{
		{
			name:      "resource with native body",
			resources: []*tfjson.StateResource{create},
			covered:   []string{"#.location"},
		},
		{
			name:      "resource with null JSON body",
			resources: []*tfjson.StateResource{nullBody},
			covered:   []string{"#.location"},
		},
		{
			name:      "resource with JSON body",
			resources: []*tfjson.StateResource{newResource("azapi_resource", "json", `{"location":"westus"}`)},
			covered:   []string{"#.location"},
		},
		{
			name:      "update with native body",
			resources: []*tfjson.StateResource{update},
			covered:   []string{"#.properties.disableLocalAuth"},
		},
		{
			name:      "update with JSON body",
			resources: []*tfjson.StateResource{newResource("azapi_update_resource", "json", `{"properties":{"disableLocalAuth":true}}`)},
			covered:   []string{"#.properties.disableLocalAuth"},
		},
		{
			name:      "create and updates aggregate",
			resources: []*tfjson.StateResource{create, update, secondUpdate, update},
			covered:   []string{"#.location", "#.properties.disableLocalAuth", "#.properties.publicNetworkAccess"},
		},
		{
			name:      "updates without a create resource aggregate",
			resources: []*tfjson.StateResource{update, secondUpdate},
			covered:   []string{"#.properties.disableLocalAuth", "#.properties.publicNetworkAccess"},
		},
		{
			name:      "ignore data sources",
			resources: []*tfjson.StateResource{data},
		},
		{
			name:      "ignore resource actions",
			resources: []*tfjson.StateResource{action},
		},
	}
	for _, testcase := range testcases {
		for _, source := range []string{"state", "no-op plan", "update plan"} {
			t.Run(testcase.name+"/"+source, func(t *testing.T) {
				var actual coverage.CoverageReport
				var err error
				if source == "state" {
					actual, err = tf.NewCoverageReportFromState(&tfjson.State{
						Values: &tfjson.StateValues{
							RootModule: &tfjson.StateModule{Resources: testcase.resources},
						},
					}, swaggerPath, nil)
				} else {
					action := tfjson.ActionNoop
					if source == "update plan" {
						action = tfjson.ActionUpdate
					}
					plan := &tfjson.Plan{}
					for _, resource := range testcase.resources {
						plan.ResourceChanges = append(plan.ResourceChanges, &tfjson.ResourceChange{
							Address: resource.Address,
							Type:    resource.Type,
							Mode:    resource.Mode,
							Change: &tfjson.Change{
								Actions: []tfjson.Action{action},
								Before:  resource.AttributeValues,
								// Pending values must not contribute coverage.
								After: map[string]interface{}{
									"body": map[string]interface{}{"name": "not-applied"},
								},
							},
						})
					}
					actual, err = tf.NewCoverageReport(plan, swaggerPath, nil)
				}
				if err != nil {
					t.Fatal(err)
				}
				if len(testcase.covered) == 0 {
					if len(actual.Coverages) != 0 {
						t.Fatalf("expected no coverage, got %d items", len(actual.Coverages))
					}
					return
				}
				if len(actual.Coverages) != 1 {
					t.Fatalf("expected one aggregated coverage item, got %d", len(actual.Coverages))
				}
				markdown := actual.MarkdownContent()
				for _, item := range actual.Coverages {
					if item.DisplayName != resourceType {
						t.Errorf("expected display name %s, got %s", resourceType, item.DisplayName)
					}
					if item.Model.RootCoveredCount != len(testcase.covered) {
						t.Errorf("expected %d covered properties, got %d", len(testcase.covered), item.Model.RootCoveredCount)
					}
					for _, property := range testcase.covered {
						assertCovered(t, markdown, property, true)
					}
				}
				if !strings.Contains(markdown, "##### <!-- "+resourceType+" -->") {
					t.Error("coverage report is missing the resource heading")
				}
			})
		}
	}
}

func Test_NewCoverageReportSkipsUnsupportedActions(t *testing.T) {
	for _, action := range []tfjson.Action{tfjson.ActionCreate, tfjson.ActionDelete, tfjson.ActionRead} {
		t.Run(string(action), func(t *testing.T) {
			actual, err := tf.NewCoverageReport(&tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Type: "azapi_update_resource",
						Mode: tfjson.ManagedResourceMode,
						Change: &tfjson.Change{
							Actions: []tfjson.Action{action},
							Before: map[string]interface{}{
								"body": map[string]interface{}{"name": "not-applied"},
							},
						},
					},
				},
			}, "", nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(actual.Coverages) != 0 {
				t.Fatalf("expected no coverage, got %d items", len(actual.Coverages))
			}
		})
	}
}

// Terraform state and plans are decoded with json.Number, which coverage
// expects as float64 like a string body.
func Test_CoverageReportFromStateWithJSONNumbers(t *testing.T) {
	var state tfjson.State
	state.UseJSONNumber(true)
	err := json.Unmarshal([]byte(`{
  "format_version": "1.0",
  "values": {"root_module": {"resources": [{
    "address": "azapi_resource.example", "mode": "managed", "type": "azapi_resource", "name": "example",
    "values": {
      "id": "`+resourceID+`", "type": "`+resourceType+`", "location": "westus",
      "body": {"properties": {"sku": {"capacity": 3}}}
    }
  }]}}
}`), &state)
	if err != nil {
		t.Fatal(err)
	}
	body := state.Values.RootModule.Resources[0].AttributeValues["body"].(map[string]interface{})
	sku := body["properties"].(map[string]interface{})["sku"].(map[string]interface{})
	if _, ok := sku["capacity"].(json.Number); !ok {
		t.Fatal("expected state to be decoded with json.Number")
	}

	output := captureLogs(t)
	report, err := tf.NewCoverageReportFromState(&state, swaggerPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "unexpect type") {
		t.Fatalf("unexpected type error: %s", output.String())
	}
	assertCovered(t, report.MarkdownContent(), "#.properties.sku.capacity", true)
	if _, ok := body["location"]; ok {
		t.Fatal("coverage mutated Terraform state")
	}
}

// Requests include write-only properties, such as sensitive_body, that are
// absent from Terraform state.
func Test_NewCoverageReportsFromRequests(t *testing.T) {
	const secret = "synthetic-secret-not-a-credential"
	request := func(method, url string, statusCode int, body string) paltypes.RequestTrace {
		return paltypes.RequestTrace{
			Provider: "azapi", Method: method, Url: url + "?api-version=2022-08-08", StatusCode: statusCode,
			Request: &paltypes.HttpRequest{Body: body},
		}
	}
	azurerm := request("PUT", resourceID, 200, `{"properties":{"sku":{"name":"Basic"}}}`)
	azurerm.Provider = "azurerm"
	logs := []paltypes.RequestTrace{
		request("PUT", resourceID, 201, `{"properties":{"sku":{"family":"`+secret+`"}}}`),
		request("PUT", resourceID, 200, `{"properties":{"publicNetworkAccess":false}}`),
		request("PUT", resourceID, 400, `{"properties":{"sku":{"capacity":1}}}`),
		request("GET", resourceID, 200, `{"properties":{"encryption":{"keySource":"Microsoft.Automation"}}}`),
		request("PUT", resourceID+"/runbooks/example", 200, `{"location":"westus"}`),
		request("PATCH", resourceID, 200, `{"name":"action"}`),
		azurerm,
		request("PUT", resourceID, 200, `not json`),
		request("PUT", resourceID, 200, ""),
		{Provider: "azapi", Method: "PUT", Url: resourceID + "?api-version=2022-08-08", StatusCode: 200},
	}
	attributes := map[string]interface{}{
		"id": resourceID, "type": resourceType,
		"body": map[string]interface{}{"properties": map[string]interface{}{"disableLocalAuth": true}},
	}
	output := captureLogs(t)

	for _, kind := range []string{"azapi_resource", "azapi_update_resource"} {
		for _, source := range []string{"state", "no-op", "update"} {
			t.Run(kind+"/"+source, func(t *testing.T) {
				var report coverage.CoverageReport
				var err error
				if source == "state" {
					report, err = tf.NewCoverageReportFromState(&tfjson.State{
						Values: &tfjson.StateValues{RootModule: &tfjson.StateModule{
							Resources: []*tfjson.StateResource{{Type: kind, Mode: tfjson.ManagedResourceMode, AttributeValues: attributes}},
						}},
					}, swaggerPath, logs)
				} else {
					report, err = tf.NewCoverageReport(&tfjson.Plan{
						ResourceChanges: []*tfjson.ResourceChange{{
							Type: kind, Mode: tfjson.ManagedResourceMode,
							Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.Action(source)}, Before: attributes},
						}},
					}, swaggerPath, logs)
				}
				if err != nil {
					t.Fatal(err)
				}
				markdown := report.MarkdownContent()
				for path, covered := range map[string]bool{
					"#.properties.disableLocalAuth":     true,
					"#.properties.sku.family":           true,
					"#.properties.publicNetworkAccess":  true,
					"#.properties.sku.capacity":         false,
					"#.properties.encryption.keySource": false,
					"#.location":                        false,
					"#.name":                            false,
					"#.properties.sku.name":             false,
				} {
					assertCovered(t, markdown, path, covered)
				}
				if strings.Contains(markdown, secret) {
					t.Fatal("secret value leaked into Markdown")
				}
			})
		}
	}
	if !strings.Contains(output.String(), "invalid body") {
		t.Errorf("expected a warning for the invalid request body, got %s", output.String())
	}
}

func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var output bytes.Buffer
	original := logrus.StandardLogger().Out
	logrus.SetOutput(&output)
	t.Cleanup(func() { logrus.SetOutput(original) })
	return &output
}
