package tf

import tfjson "github.com/hashicorp/terraform-json"

func GetChanges(plan *tfjson.Plan) int {
	if plan == nil {
		return 0
	}
	count := 0
	for _, change := range plan.ResourceChanges {
		if change.Change != nil {
			if len(change.Change.Actions) != 1 || change.Change.Actions[0] != tfjson.ActionNoop {
				count++
			}
		}
	}
	return count
}
