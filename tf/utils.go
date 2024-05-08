package tf

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/azure/armstrong/coverage"
	"github.com/azure/armstrong/types"
	"github.com/azure/armstrong/utils"
	tfjson "github.com/hashicorp/terraform-json"
	paltypes "github.com/ms-henglu/pal/types"
	"github.com/sirupsen/logrus"
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

func NewDiffReport(plan *tfjson.Plan, logs []paltypes.RequestTrace) types.DiffReport {
	out := types.DiffReport{
		Diffs: make([]types.Diff, 0),
		Logs:  logs,
	}
	if plan == nil {
		return out
	}

	for _, resourceChange := range plan.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Change.Before == nil || resourceChange.Change.After == nil {
			continue
		}
		if !strings.HasPrefix(resourceChange.Address, "azapi_") {
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
		if afterMap["id"] == nil {
			logrus.Errorf("resource %s has no id", resourceChange.Address)
			continue
		}

		change := types.Change{}

		if _, ok := beforeMap["body"].(string); ok {
			change = types.Change{
				Before: beforeMap["body"].(string),
				After:  afterMap["body"].(string),
			}
		} else {
			payloadBefore, _ := json.Marshal(beforeMap["body"])
			payloadAfter, _ := json.Marshal(afterMap["body"])
			change = types.Change{
				Before: string(payloadBefore),
				After:  string(payloadAfter),
			}
		}

		out.Diffs = append(out.Diffs, types.Diff{
			Id:      afterMap["id"].(string),
			Type:    afterMap["type"].(string),
			Address: resourceChange.Address,
			Change:  change,
		})
	}

	return out
}

func NewPassReportFromState(state *tfjson.State) types.PassReport {
	out := types.PassReport{
		Resources: make([]types.Resource, 0),
	}
	if state == nil || state.Values == nil || state.Values.RootModule == nil || state.Values.RootModule.Resources == nil {
		logrus.Warnf("new pass report from state: state is nil")
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
		if resourceChange == nil || resourceChange.Change == nil {
			continue
		}
		if !strings.HasPrefix(resourceChange.Address, "azapi_") {
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

func NewCoverageReportFromState(state *tfjson.State, swaggerPath string) (coverage.CoverageReport, error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("panic when producing coverage report from state: %+v", r)
		}
	}()

	out := coverage.CoverageReport{
		Coverages: make(map[coverage.ArmResource]*coverage.Model, 0),
	}
	if state == nil || state.Values == nil || state.Values.RootModule == nil || state.Values.RootModule.Resources == nil {
		logrus.Warnf("new coverage report from state: state is nil")
		return out, nil
	}
	for _, res := range state.Values.RootModule.Resources {
		if res.Type != "azapi_resource" || res.Mode != tfjson.ManagedResourceMode {
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

		err = out.AddCoverageFromState(id, resourceType, body, swaggerPath)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func NewCoverageReport(plan *tfjson.Plan, swaggerPath string) (coverage.CoverageReport, error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("panic when producing coverage report: %+v", r)
		}
	}()

	out := coverage.CoverageReport{
		Coverages: make(map[coverage.ArmResource]*coverage.Model, 0),
	}
	if plan == nil {
		return out, nil
	}

	for _, resourceChange := range plan.ResourceChanges {
		if resourceChange.Type != "azapi_resource" {
			continue
		}
		if resourceChange == nil || resourceChange.Change == nil {
			continue
		}
		if actions := resourceChange.Change.Actions; len(actions) == 1 && (actions[0] == tfjson.ActionNoop || actions[0] == tfjson.ActionUpdate) {
			outMap, beforeMapOk := resourceChange.Change.Before.(map[string]interface{})
			if !beforeMapOk {
				continue
			}

			beforeMap := DeepCopy(outMap).(map[string]interface{})

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

			err = out.AddCoverageFromState(id, resourceType, body, swaggerPath)
			if err != nil {
				return out, err
			}

		}
	}

	return out, nil
}

func getBody(input map[string]interface{}) (map[string]interface{}, error) {
	output := map[string]interface{}{}
	bodyRaw, ok := input["body"]
	if !ok || bodyRaw == nil {
		return output, nil
	}
	if bodyStr, ok := bodyRaw.(string); ok && bodyStr != "" {
		if value, ok := input["tags"]; ok && value != nil && len(value.(map[string]interface{})) > 0 {
			output["tags"] = value.(map[string]interface{})
		}

		if value, ok := input["location"]; ok && value != nil && value.(string) != "" {
			output["location"] = value.(string)
		}

		if value, ok := input["identity"]; ok && value != nil && len(value.([]interface{})) > 0 {
			output["identity"] = expandIdentity(value.([]interface{}))
		}

		err := json.Unmarshal([]byte(bodyStr), &output)
		return output, err
	}
	if bodyMap, ok := bodyRaw.(map[string]interface{}); ok {
		if value, ok := input["tags"]; ok && value != nil && len(value.(map[string]interface{})) > 0 {
			bodyMap["tags"] = value.(map[string]interface{})
		}

		if value, ok := input["location"]; ok && value != nil && value.(string) != "" {
			bodyMap["location"] = value.(string)
		}

		if value, ok := input["identity"]; ok && value != nil && len(value.([]interface{})) > 0 {
			bodyMap["identity"] = expandIdentity(value.([]interface{}))
		}
		return bodyMap, nil
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

func NewErrorReport(applyErr error, logs []paltypes.RequestTrace) types.ErrorReport {
	out := types.ErrorReport{
		Errors: make([]types.Error, 0),
		Logs:   logs,
	}
	if applyErr == nil {
		return out
	}
	res := make([]string, 0)
	if strings.Contains(applyErr.Error(), "Error: Failed to create/update resource") {
		res = strings.Split(applyErr.Error(), "Error: Failed to create/update resource")
	} else {
		res = strings.Split(applyErr.Error(), "Error: creating/updating")
	}
	for _, e := range res {
		var id, apiVersion, label string
		errorMessage := e
		if lastIndex := strings.LastIndex(e, "------"); lastIndex != -1 {
			errorMessage = errorMessage[0:lastIndex]
		}
		if matches := regexp.MustCompile(`ResourceId\s+\\?"([^\\]+)\\?"\s+/\s+Api Version \\?"([^\\]+)\\?"\)`).FindAllStringSubmatch(e, -1); len(matches) == 1 {
			id = matches[0][1]
			apiVersion = matches[0][2]
		}
		if matches := regexp.MustCompile(`resource "azapi_.+" "(.+)"`).FindAllStringSubmatch(e, -1); len(matches) != 0 {
			label = matches[0][1]
		}
		if len(label) == 0 {
			continue
		}
		out.Errors = append(out.Errors, types.Error{
			Id:      id,
			Type:    fmt.Sprintf("%s@%s", utils.ResourceTypeOfResourceId(id), apiVersion),
			Label:   label,
			Message: errorMessage,
		})
	}
	return out
}

func NewCleanupErrorReport(applyErr error, logs []paltypes.RequestTrace) types.ErrorReport {
	out := types.ErrorReport{
		Errors: make([]types.Error, 0),
		Logs:   logs,
	}
	if applyErr == nil {
		return out
	}
	res := make([]string, 0)
	if strings.Contains(applyErr.Error(), "Error: Failed to delete resource") {
		res = strings.Split(applyErr.Error(), "Error: Failed to delete resource")
	} else {
		res = strings.Split(applyErr.Error(), "Error: deleting")
	}
	for _, e := range res {
		var id, apiVersion string
		errorMessage := e
		if lastIndex := strings.LastIndex(e, "------"); lastIndex != -1 {
			errorMessage = errorMessage[0:lastIndex]
		}
		if matches := regexp.MustCompile(`ResourceId\s+\\?"([^\\]+)\\?"\s+/\s+Api Version \\?"([^\\]+)\\?"\)`).FindAllStringSubmatch(e, -1); len(matches) == 1 {
			id = matches[0][1]
			apiVersion = matches[0][2]
		} else {
			continue
		}

		out.Errors = append(out.Errors, types.Error{
			Id:      id,
			Type:    fmt.Sprintf("%s@%s", utils.ResourceTypeOfResourceId(id), apiVersion),
			Message: errorMessage,
		})
	}
	return out
}

func NewIdAddressFromState(state *tfjson.State) map[string]string {
	out := map[string]string{}
	if state == nil || state.Values == nil || state.Values.RootModule == nil || state.Values.RootModule.Resources == nil {
		logrus.Warnf("new id address mapping from state: state is nil")
		return out
	}
	for _, res := range state.Values.RootModule.Resources {
		id := ""
		if v, ok := res.AttributeValues["id"]; ok {
			id = v.(string)
		}
		out[id] = res.Address
	}
	return out
}

func DeepCopy(input interface{}) interface{} {
	if input == nil {
		return nil
	}
	switch v := input.(type) {
	case map[string]interface{}:
		out := map[string]interface{}{}
		for key, value := range v {
			out[key] = DeepCopy(value)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, value := range v {
			out[i] = DeepCopy(value)
		}
		return out
	default:
		return input
	}
}
