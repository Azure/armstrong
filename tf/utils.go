package tf

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/ms-henglu/azurerm-restapi-testing-tool/types"
	"strings"
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

func NewReports(plan *tfjson.Plan) []types.Report {
	res := make([]types.Report, 0)
	if plan != nil {
		for _, resourceChange := range plan.ResourceChanges {
			if !strings.HasPrefix(resourceChange.Address, "azapi_") {
				continue
			}
			if resourceChange == nil || resourceChange.Change == nil || resourceChange.Change.Before == nil || resourceChange.Change.After == nil {
				continue
			}
			beforeMap, beforeMapOk := resourceChange.Change.Before.(map[string]interface{})
			afterMap, afterMapOk := resourceChange.Change.After.(map[string]interface{})
			if !beforeMapOk || !afterMapOk {
				continue
			}
			res = append(res, types.Report{
				Id:      afterMap["id"].(string),
				Type:    afterMap["type"].(string),
				Address: resourceChange.Address,
				Change: types.Diff{
					Before: beforeMap["body"].(string),
					After:  afterMap["body"].(string),
				},
			})

		}
	}
	return res
}
