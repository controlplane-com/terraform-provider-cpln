package modifiers

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// Ensure implementation satisfies the planmodifier.String interface
var _ planmodifier.String = SuppressDiffOnEqualJSON{}

// SuppressDiffOnEqualJSON is a custom plan modifier that suppresses diffs when JSON strings are logically equivalent.
type SuppressDiffOnEqualJSON struct{}

// Returns a plain text description of the plan modifier's behavior
func (m SuppressDiffOnEqualJSON) Description(ctx context.Context) string {
	return "Suppresses diff if the planned and prior state JSON strings are equivalent"
}

// Returns a markdown-formatted description of the plan modifier's behavior
func (m SuppressDiffOnEqualJSON) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// Applies the plan modification by suppressing the diff if both JSON strings are semantically equal
func (m SuppressDiffOnEqualJSON) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if either the current state or planned value is null
	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		return
	}

	// Declare maps to hold unmarshaled JSON objects
	var stateObj, planObj map[string]interface{}

	// Unmarshal the current state JSON string into a map
	if err := json.Unmarshal([]byte(req.StateValue.ValueString()), &stateObj); err != nil {
		return
	}

	// Unmarshal the planned JSON string into a map
	if err := json.Unmarshal([]byte(req.PlanValue.ValueString()), &planObj); err != nil {
		return
	}

	// If the two JSON objects are semantically equal, suppress the diff by using the state value
	if equalJSON(stateObj, planObj) {
		resp.PlanValue = req.StateValue
	}
}

// Compares two JSON objects by marshaling and checking byte-level equality
func equalJSON(a, b map[string]interface{}) bool {
	// Marshal both maps into JSON byte slices
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	// Compare the marshaled byte slices as strings
	return string(aJSON) == string(bJSON)
}
