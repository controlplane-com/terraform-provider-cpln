package cpln

import (
	"fmt"
	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/modifiers"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BaseResourceModel holds the shared attributes for Terraform resources.
type BaseResourceModel struct {
	ID          types.String `tfsdk:"id"`
	CplnID      types.String `tfsdk:"cpln_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.Map    `tfsdk:"tags"`
	SelfLink    types.String `tfsdk:"self_link"`
}

// BaseResourceAttributes returns a map of attributes for a given resource name.
func BaseResourceAttributes(resourceName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The unique identifier for this %s.", resourceName),
		},
		"cpln_id": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The ID, in GUID format, of the %s.", resourceName),
		},
		"name": schema.StringAttribute{
			Required:    true,
			Description: fmt.Sprintf("Name of the %s.", resourceName),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				validators.NameValidator{},
			},
		},
		"description": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: fmt.Sprintf("Description of the %s.", resourceName),
			PlanModifiers: []planmodifier.String{
				modifiers.DescriptionPlanModifier{},
			},
			Validators: []validator.String{
				validators.DescriptionValidator{},
			},
		},
		"tags": schema.MapAttribute{
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			Description: "Key-value map of resource tags.",
			PlanModifiers: []planmodifier.Map{
				modifiers.TagPlanModifier{},
			},
			Validators: []validator.Map{
				validators.TagValidator{},
			},
		},
		"self_link": schema.StringAttribute{
			Computed:    true,
			Description: "Full link to this resource. Can be referenced by other resources.",
		},
	}
}

// UpdateBaseClientFromState updates the base client from BaseResourceModel.
func UpdateBaseClientFromState(base *client.Base, state BaseResourceModel) {
	base.Name = BuildString(state.Name)
	base.Description = BuildString(state.Description)
	base.Tags = BuildTags(state.Tags)
}

// UpdateReplaceBaseClientFromState updates the base client from BaseResourceModel.
func UpdateReplaceBaseClientFromState(base *client.Base, state BaseResourceModel) {
	base.Name = BuildString(state.Name)
	base.Description = BuildString(state.Description)
	base.TagsReplace = BuildTags(state.Tags)
}

// UpdateStateFromBaseClient updates the state attributes from the base client.
func UpdateStateFromBaseClient(state *BaseResourceModel, base client.Base) {
	state.ID = types.StringPointerValue(base.Name)
	state.CplnID = types.StringPointerValue(base.ID)
	state.Name = types.StringPointerValue(base.Name)
	state.Description = types.StringPointerValue(base.Description)
	state.Tags = SchemaTags(base.Tags)
	state.SelfLink = SchemaSelfLink(base.Links)
}
