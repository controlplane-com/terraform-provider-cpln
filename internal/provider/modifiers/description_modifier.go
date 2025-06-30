package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the interface.
var _ planmodifier.String = DescriptionPlanModifier{}

// DescriptionPlanModifier is a custom plan modifier that sets the `description` attribute
// to match the `name` attribute if `description` is not explicitly set by the user.
type DescriptionPlanModifier struct{}

// Description provides a plain text description of the plan modifier's behavior.
func (m DescriptionPlanModifier) Description(ctx context.Context) string {
	return "Sets the description to match the name if not explicitly set"
}

// MarkdownDescription a markdown-formatted description of the plan modifier's behavior.
func (m DescriptionPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyString adjusts the `description` attribute during the planning phase.
func (m DescriptionPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the `description` attribute already has a value, no modification is needed
	if !req.ConfigValue.IsNull() {
		return
	}

	// Retrieve the `name` attribute from the configuration to use as a fallback value
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)

	// Exit if there were errors retrieving the `name` attribute
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the `description` attribute in the planned value to match the `name` attribute
	resp.PlanValue = name
}
