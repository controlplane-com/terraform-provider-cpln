package catalog_template

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Helm Release Resource //

type HelmReleaseResourceModel struct {
	Kind types.String `tfsdk:"kind"`
	Name types.String `tfsdk:"name"`
	Link types.String `tfsdk:"link"`
}

func (r HelmReleaseResourceModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"kind": types.StringType,
			"name": types.StringType,
			"link": types.StringType,
		},
	}
}
