package tf_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/azure/armstrong/coverage"
	"github.com/azure/armstrong/tf"
	tfjson "github.com/hashicorp/terraform-json"
)

func Test_NewCoverageReports(t *testing.T) {
	const resourceType = "Microsoft.Automation/automationAccounts@2022-08-08"
	const resourceID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Automation/automationAccounts/test"
	swaggerPath := filepath.Join("..", "coverage", "testdata", "Microsoft.Automation", "stable", "2022-08-08", "account.json")

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
					}, swaggerPath)
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
					actual, err = tf.NewCoverageReport(plan, swaggerPath)
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
						name := property[strings.LastIndex(property, ".")+1:]
						expected := "<!-- " + property + " -->\n<details>\n<summary><span>" + name + "</span>"
						if !strings.Contains(markdown, expected) {
							t.Errorf("property %s was not marked fully covered in Markdown", property)
						}
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
			}, "")
			if err != nil {
				t.Fatal(err)
			}
			if len(actual.Coverages) != 0 {
				t.Fatalf("expected no coverage, got %d items", len(actual.Coverages))
			}
		})
	}
}
