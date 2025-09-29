package policy

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Binding //

type BindingModel struct {
	Permissions    types.Set `tfsdk:"permissions"`
	PrincipalLinks types.Set `tfsdk:"principal_links"`
}

func (b BindingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"permissions":     types.SetType{ElemType: types.StringType},
			"principal_links": types.SetType{ElemType: types.StringType},
		},
	}
}
