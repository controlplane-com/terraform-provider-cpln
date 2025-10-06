package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the planmodifier.Map interface
var _ planmodifier.Map = DictionaryAsEnvsPlanModifier{}

// DictionaryAsEnvsPlanModifier plans `dictionary_as_envs` deterministically when possible,
// leaves it unknown when the `dictionary` changes (so Terraform expects new keys),
// and otherwise stabilizes by reusing state to avoid noisy diffs.
type DictionaryAsEnvsPlanModifier struct{}

// Returns a plain text description of the plan modifier's behavior
func (DictionaryAsEnvsPlanModifier) Description(context.Context) string {
	return "Plans dictionary_as_envs from dictionary/name when possible; otherwise stabilizes or leaves unknown on change."
}

// Returns a markdown-formatted description of the plan modifier's behavior
func (DictionaryAsEnvsPlanModifier) MarkdownDescription(context.Context) string {
	return "Plans `dictionary_as_envs` from `dictionary`/`name` when possible; otherwise stabilizes or leaves unknown on change."
}

// Applies the planning rules for `dictionary_as_envs` to avoid churn and handle shape changes safely
func (DictionaryAsEnvsPlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// Exit early since there is nothing to compute or stabilize
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Declare a variable to hold the planned value of the `dictionary` attribute
	var plannedDict types.Map

	// Load the planned `dictionary` from the plan for decision-making
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("dictionary"), &plannedDict)...)

	// Exit to avoid compounding errors
	if resp.Diagnostics.HasError() {
		return
	}

	// Declare a variable to hold the planned value of the `name` attribute
	var plannedName types.String

	// Load the planned `name` (used to synthesize the `cpln://secret/<name>.<key>` values)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plannedName)...)

	// Exit to avoid compounding errors
	if resp.Diagnostics.HasError() {
		return
	}

	// If this is not a dictionary secret (plannedDict is null), keep dictionary_as_envs null to prevent churn
	if plannedDict.IsNull() {
		// Set a null map of strings to represent the absence of env links
		resp.PlanValue = types.MapNull(types.StringType)

		// We are done for non-dictionary types
		return
	}

	// Declare a variable to hold the prior state's `dictionary` for change detection
	var priorDict types.Map

	// Load the prior `dictionary` from state so we can detect shape/content changes
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("dictionary"), &priorDict)...)

	// Exit to avoid compounding errors
	if resp.Diagnostics.HasError() {
		return
	}

	// Leaving as unknown informs Terraform that new keys may appear at apply, preventing "inconsistent result"
	if plannedDict.IsUnknown() || !equalMaps(ctx, plannedDict, priorDict) {
		return
	}

	// If both the `name` and `dictionary` are known and unchanged, compute the exact env map now
	if !plannedName.IsNull() && !plannedName.IsUnknown() {
		// Prepare a map to hydrate the planned dictionary into Go values
		dictKV := map[string]string{}

		// Extract the planned dictionary entries as plain strings
		resp.Diagnostics.Append(plannedDict.ElementsAs(ctx, &dictKV, false)...)

		// Exit to avoid compounding errors
		if resp.Diagnostics.HasError() {
			return
		}

		// Allocate an output map sized to the number of dictionary keys
		out := make(map[string]string, len(dictKV))

		// For each key in the dictionary, synthesize the corresponding secret link value
		for k := range dictKV {
			// Compose the cpln://secret/<name>.<key> link
			out[k] = fmt.Sprintf("cpln://secret/%s.%s", plannedName.ValueString(), k)
		}

		// Convert the synthesized map into a Terraform types.Map of strings
		mv, diags := types.MapValueFrom(ctx, types.StringType, out)

		// Append any diagnostics that occurred during conversion
		resp.Diagnostics.Append(diags...)

		// Exit to avoid compounding errors
		if resp.Diagnostics.HasError() {
			return
		}

		// Set the fully-computed map as the planned value
		resp.PlanValue = mv

		// We are done after computing the precise value
		return
	}

	// As a fallback when we cannot compute or confirm stability, reuse the state's value to avoid churn
	resp.PlanValue = req.StateValue
}

// Compares two Terraform maps (string -> string) for equality without treating unknown as equal
func equalMaps(ctx context.Context, a, b types.Map) bool {
	// Return equality for the null case
	if a.IsNull() && b.IsNull() {
		return true
	}

	// Return inequality when unknowns are involved
	if a.IsUnknown() || b.IsUnknown() {
		return false
	}

	// Declare containers for materializing map A as native Go values
	var am map[string]string

	// On extraction error, consider maps unequal
	if diags := a.ElementsAs(ctx, &am, false); diags.HasError() {
		return false
	}

	// Declare containers for materializing map B as native Go values
	var bm map[string]string

	// On extraction error, consider maps unequal
	if diags := b.ElementsAs(ctx, &bm, false); diags.HasError() {
		return false
	}

	// Return inequality due to length mismatch
	if len(am) != len(bm) {
		return false
	}

	// Iterate over all keys of the first map to compare values
	for k, av := range am {
		// Fetch the counterpart value from the second map
		bv, ok := bm[k]

		// Return inequality due to missing key or value mismatch
		if !ok || av != bv {
			return false
		}
	}

	// If all keys and values match, the maps are equal
	return true
}
