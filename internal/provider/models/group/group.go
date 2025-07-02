package group

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Identity Matcher //

type IdentityMatcherModel struct {
	Expression types.String `tfsdk:"expression"`
	Language   types.String `tfsdk:"language"`
}

func (i IdentityMatcherModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"expression": types.StringType,
			"language":   types.StringType,
		},
	}
}
