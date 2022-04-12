package tf

import tfjson "github.com/hashicorp/terraform-json"

const TestResourceAddress = "azapi_resource.test"

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

func GetBodyChange(plan *tfjson.Plan) (string, string) {
	before := "{}"
	after := "{}"
	if plan != nil {
		for _, resourceChange := range plan.ResourceChanges {
			if resourceChange != nil && resourceChange.Address == TestResourceAddress && resourceChange.Change != nil {
				if resourceChange.Change.Before != nil {
					if beforeMap, ok := resourceChange.Change.Before.(map[string]interface{}); ok && beforeMap["body"] != nil {
						if value, ok := beforeMap["body"].(string); ok {
							before = value
						}
					}
				}
				if resourceChange.Change.After != nil {
					if afterMap, ok := resourceChange.Change.After.(map[string]interface{}); ok && afterMap["body"] != nil {
						if value, ok := afterMap["body"].(string); ok {
							after = value
						}
					}
				}
				break
			}
		}
	}
	return before, after
}
