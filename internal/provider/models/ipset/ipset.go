package ipset

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Location //

type LocationModel struct {
	Name            types.String `tfsdk:"name"`
	RetentionPolicy types.String `tfsdk:"retention_policy"`
}

func (l LocationModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":             types.StringType,
			"retention_policy": types.StringType,
		},
	}
}

// Status //

type StatusModel struct {
	IpAddresses types.List   `tfsdk:"ip_address"`
	Error       types.String `tfsdk:"error"`
	Warning     types.String `tfsdk:"warning"`
}

func (l StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ip_address": types.ListType{ElemType: IpAddressModel{}.AttributeTypes()},
			"error":      types.StringType,
			"warning":    types.StringType,
		},
	}
}

// Status -> IP Address //

type IpAddressModel struct {
	Name    types.String `tfsdk:"name"`
	IP      types.String `tfsdk:"ip"`
	ID      types.String `tfsdk:"id"`
	State   types.String `tfsdk:"state"`
	Created types.String `tfsdk:"created"`
}

func (i IpAddressModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":    types.StringType,
			"ip":      types.StringType,
			"id":      types.StringType,
			"state":   types.StringType,
			"created": types.StringType,
		},
	}
}
