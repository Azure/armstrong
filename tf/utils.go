package tf

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/ms-henglu/armstrong/coverage"
	"github.com/ms-henglu/armstrong/resource/utils"
	"github.com/ms-henglu/armstrong/types"
)

type Action string

const (
	ActionCreate  Action = "create"
	ActionReplace Action = "replace"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
)

// Actions denotes a valid change type.
type Actions []Action

func GetChanges(plan *tfjson.Plan) []Action {
	if plan == nil {
		return []Action{}
	}
	actions := make([]Action, 0)
	for _, change := range plan.ResourceChanges {
		if change.Change != nil {
			if len(change.Change.Actions) == 0 {
				continue
			}
			if len(change.Change.Actions) == 1 {
				switch change.Change.Actions[0] {
				case tfjson.ActionCreate:
					actions = append(actions, ActionCreate)
				case tfjson.ActionDelete:
					actions = append(actions, ActionDelete)
				case tfjson.ActionUpdate:
					actions = append(actions, ActionUpdate)
				case tfjson.ActionNoop:
				case tfjson.ActionRead:
				}
			} else {
				actions = append(actions, ActionReplace)
			}
		}
	}
	return actions
}

func NewDiffReport(plan *tfjson.Plan, logs []types.RequestTrace) types.DiffReport {
	out := types.DiffReport{
		Diffs: make([]types.Diff, 0),
		Logs:  logs,
	}
	if plan == nil {
		return out
	}

	for _, resourceChange := range plan.ResourceChanges {
		if !strings.HasPrefix(resourceChange.Address, "azapi_") {
			continue
		}
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Change.Before == nil || resourceChange.Change.After == nil {
			continue
		}
		if len(resourceChange.Change.Actions) == 1 && resourceChange.Change.Actions[0] == tfjson.ActionNoop {
			continue
		}
		beforeMap, beforeMapOk := resourceChange.Change.Before.(map[string]interface{})
		afterMap, afterMapOk := resourceChange.Change.After.(map[string]interface{})
		if !beforeMapOk || !afterMapOk {
			continue
		}
		out.Diffs = append(out.Diffs, types.Diff{
			Id:      afterMap["id"].(string),
			Type:    afterMap["type"].(string),
			Address: resourceChange.Address,
			Change: types.Change{
				Before: beforeMap["body"].(string),
				After:  afterMap["body"].(string),
			},
		})
	}

	return out
}

func NewPassReportFromState(state *tfjson.State) types.PassReport {
	out := types.PassReport{
		Resources: make([]types.Resource, 0),
	}
	if state == nil || state.Values == nil || state.Values.RootModule == nil || state.Values.RootModule.Resources == nil {
		log.Printf("[WARN] new pass report from state: state is nil")
		return out
	}
	for _, res := range state.Values.RootModule.Resources {
		if !strings.HasPrefix(res.Address, "azapi_") {
			continue
		}
		resourceType := ""
		if v, ok := res.AttributeValues["type"]; ok {
			resourceType = v.(string)
		}
		out.Resources = append(out.Resources, types.Resource{
			Type:    resourceType,
			Address: res.Address,
		})
	}
	return out
}

func NewPassReport(plan *tfjson.Plan) types.PassReport {
	out := types.PassReport{
		Resources: make([]types.Resource, 0),
	}
	if plan == nil {
		return out
	}

	for _, resourceChange := range plan.ResourceChanges {
		if !strings.HasPrefix(resourceChange.Address, "azapi_") {
			continue
		}
		if resourceChange == nil || resourceChange.Change == nil {
			continue
		}
		if len(resourceChange.Change.Actions) == 1 && resourceChange.Change.Actions[0] == tfjson.ActionNoop {
			beforeMap, beforeMapOk := resourceChange.Change.Before.(map[string]interface{})
			if !beforeMapOk {
				continue
			}
			out.Resources = append(out.Resources, types.Resource{
				Type:    beforeMap["type"].(string),
				Address: resourceChange.Address,
			})
		}
	}

	return out
}

func NewCoverageReportFromState(state *tfjson.State) (coverage.CoverageReport, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] panic when producing coverage report from state: %+v", r)
		}
	}()

	out := coverage.CoverageReport{
		Coverages: make(map[coverage.ArmResource]*coverage.Model, 0),
	}
	if state == nil || state.Values == nil || state.Values.RootModule == nil || state.Values.RootModule.Resources == nil {
		log.Print("[WARN] new coverage report from state: state is nil")
		return out, nil
	}
	for _, res := range state.Values.RootModule.Resources {
		if !strings.HasPrefix(res.Address, "azapi_") {
			continue
		}

		id := ""
		if v, ok := res.AttributeValues["id"]; ok {
			id = v.(string)
		}
		resourceType := ""
		if v, ok := res.AttributeValues["type"]; ok {
			resourceType = v.(string)
		}

		body, err := getBody(res.AttributeValues)
		if err != nil {
			return out, err
		}

		err = out.AddCoverageFromState(id, resourceType, body)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func NewCoverageReport(plan *tfjson.Plan) (coverage.CoverageReport, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] panic when producing coverage report: %+v", r)
		}
	}()

	out := coverage.CoverageReport{
		Coverages: make(map[coverage.ArmResource]*coverage.Model, 0),
	}
	if plan == nil {
		return out, nil
	}

	for _, resourceChange := range plan.ResourceChanges {
		if !strings.HasPrefix(resourceChange.Address, "azapi_") {
			continue
		}
		if resourceChange == nil || resourceChange.Change == nil {
			continue
		}
		if len(resourceChange.Change.Actions) == 1 && resourceChange.Change.Actions[0] == tfjson.ActionNoop {
			beforeMap, beforeMapOk := resourceChange.Change.Before.(map[string]interface{})
			if !beforeMapOk {
				continue
			}

			id := ""
			if v, ok := beforeMap["id"]; ok {
				id = v.(string)
			}

			resourceType := ""
			if v, ok := beforeMap["type"]; ok {
				resourceType = v.(string)
			}

			body, err := getBody(beforeMap)
			if err != nil {
				return out, err
			}

			err = out.AddCoverageFromState(id, resourceType, body)
			if err != nil {
				return out, err
			}

		}
	}

	return out, nil
}

func getBody(input map[string]interface{}) (map[string]interface{}, error) {
	output := map[string]interface{}{}
	if bodyRaw, ok := input["body"]; ok && bodyRaw != nil && bodyRaw.(string) != "" {
		if value, ok := input["tags"]; ok && value != nil && len(value.(map[string]interface{})) > 0 {
			output["tags"] = value.(map[string]interface{})
		}

		if value, ok := input["location"]; ok && value != nil && value.(string) != "" {
			output["location"] = value.(string)
		}

		if value, ok := input["identity"]; ok && value != nil && len(value.([]interface{})) > 0 {
			output["identity"] = expandIdentity(value.([]interface{}))
		}

		err := json.Unmarshal([]byte(bodyRaw.(string)), &output)
		if err != nil {
			return output, err
		}
	}

	return output, nil
}

func expandIdentity(input []interface{}) map[string]interface{} {
	config := map[string]interface{}{}
	if len(input) == 0 {
		return config
	}
	v := input[0].(map[string]interface{})

	if identityTypeRaw, ok := v["type"]; ok && identityTypeRaw != nil && identityTypeRaw.(string) != "" {
		config["type"] = identityTypeRaw.(string)
	}

	if identityIdsRaw, ok := v["identity_ids"]; ok && identityIdsRaw != nil && len(identityIdsRaw.([]interface{})) > 0 {
		identityIds := identityIdsRaw.([]interface{})
		userAssignedIdentities := make(map[string]interface{}, len(identityIds))
		for _, id := range identityIds {
			userAssignedIdentities[id.(string)] = make(map[string]interface{})
		}
		config["userAssignedIdentities"] = userAssignedIdentities
	}

	return config
}

func NewErrorReport(applyErr error, logs []types.RequestTrace) types.ErrorReport {
	out := types.ErrorReport{
		Errors: make([]types.Error, 0),
		Logs:   logs,
	}
	res := strings.Split(applyErr.Error(), "Error: creating/updating")
	for _, e := range res {
		var id, apiVersion, label string
		errorMessage := e
		if lastIndex := strings.LastIndex(e, "------"); lastIndex != -1 {
			errorMessage = errorMessage[0:lastIndex]
		}
		if matches := regexp.MustCompile(`ResourceId \\"(.+)\\" / Api Version \\"(.+)\\"\)`).FindAllStringSubmatch(e, -1); len(matches) == 1 {
			id = matches[0][1]
			apiVersion = matches[0][2]
		}
		if matches := regexp.MustCompile(`resource "azapi_resource" "(.+)"`).FindAllStringSubmatch(e, -1); len(matches) != 0 {
			label = matches[0][1]
		}
		if len(label) == 0 {
			continue
		}
		out.Errors = append(out.Errors, types.Error{
			Id:      id,
			Type:    fmt.Sprintf("%s@%s", utils.GetResourceType(id), apiVersion),
			Label:   label,
			Message: errorMessage,
		})
	}
	return out
}
