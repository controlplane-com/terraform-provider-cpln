package gvc

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Sidecar //

type SidecarModel struct {
	Envoy types.String `tfsdk:"envoy"`
}

func (s SidecarModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"envoy": types.StringType,
		},
	}
}

// Load Balancer //

type LoadBalancerModel struct {
	Dedicated      types.Bool   `tfsdk:"dedicated"`
	MultiZone      types.List   `tfsdk:"multi_zone"`
	TrustedProxies types.Int32  `tfsdk:"trusted_proxies"`
	Redirect       types.List   `tfsdk:"redirect"`
	IpSet          types.String `tfsdk:"ipset"`
}

func (l LoadBalancerModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"dedicated":       types.BoolType,
			"multi_zone":      types.ListType{ElemType: LoadBalancerMultiZoneModel{}.AttributeTypes()},
			"trusted_proxies": types.Int32Type,
			"redirect":        types.ListType{ElemType: LoadBalancerRedirectModel{}.AttributeTypes()},
			"ipset":           types.StringType,
		},
	}
}

// Load Balancer -> Multi Zone //

type LoadBalancerMultiZoneModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

func (l LoadBalancerMultiZoneModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
		},
	}
}

// Load Balancer -> Redirect //

type LoadBalancerRedirectModel struct {
	Class types.List `tfsdk:"class"`
}

func (l LoadBalancerRedirectModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"class": types.ListType{ElemType: LoadBalancerRedirectClassModel{}.AttributeTypes()},
		},
	}
}

// Load Balancer -> Redirect -> Class //

type LoadBalancerRedirectClassModel struct {
	Status5xx types.String `tfsdk:"status_5xx"`
	Status401 types.String `tfsdk:"status_401"`
}

func (l LoadBalancerRedirectClassModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"status_5xx": types.StringType,
			"status_401": types.StringType,
		},
	}
}

// KEDA //

type KedaModel struct {
	Enabled      types.Bool   `tfsdk:"enabled"`
	IdentityLink types.String `tfsdk:"identity_link"`
	Secrets      types.Set    `tfsdk:"secrets"`
}

func (k KedaModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled":       types.BoolType,
			"identity_link": types.StringType,
			"secrets":       types.SetType{ElemType: types.StringType},
		},
	}
}
