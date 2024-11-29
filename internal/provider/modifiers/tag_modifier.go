package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the interface
var _ planmodifier.Map = TagPlanModifier{}

// TagPlanModifier is a custom plan modifier that ensures the `tags` attribute is always initialized.
type TagPlanModifier struct{}

// Description provides a plain text description of the plan modifier's behavior.
func (m TagPlanModifier) Description(ctx context.Context) string {
	return "Sets an empty map for tags if not explicitly set."
}

// MarkdownDescription provides a markdown-formatted description of the plan modifier's behavior.
func (m TagPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyMap adjusts the planned value of the `tags` attribute during the planning phase.
func (m TagPlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// If the `tags` attribute already has a value in the plan, do nothing
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// Set an empty map as the plan value
	resp.PlanValue = types.MapValueMust(types.StringType, map[string]attr.Value{})
}
