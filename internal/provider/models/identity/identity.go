package identity

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AWS Access Policy //

type AwsAccessPolicyModel struct {
	CloudAccountLink types.String `tfsdk:"cloud_account_link"`
	PolicyRefs       types.Set    `tfsdk:"policy_refs"`
	RoleName         types.String `tfsdk:"role_name"`
	TrustPolicy      types.List   `tfsdk:"trust_policy"`
}

func (a AwsAccessPolicyModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cloud_account_link": types.StringType,
			"policy_refs":        types.SetType{ElemType: types.StringType},
			"role_name":          types.StringType,
			"trust_policy":       types.ListType{ElemType: AwsAccessPolicyTrustPolicyModel{}.AttributeTypes()},
		},
	}
}

// AWS Access Policy -> Trust Policy //

type AwsAccessPolicyTrustPolicyModel struct {
	Version   types.String `tfsdk:"version"`
	Statement types.Set    `tfsdk:"statement"`
}

func (a AwsAccessPolicyTrustPolicyModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"version":   types.StringType,
			"statement": types.SetType{ElemType: types.MapType{ElemType: types.StringType}},
		},
	}
}

// GCP Access Policy //

type GcpAccessPolicyModel struct {
	CloudAccountLink types.String `tfsdk:"cloud_account_link"`
	Scopes           types.String `tfsdk:"scopes"`
	ServiceAccount   types.String `tfsdk:"service_account"`
	Binding          types.List   `tfsdk:"binding"`
}

func (g GcpAccessPolicyModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cloud_account_link": types.StringType,
			"scopes":             types.StringType,
			"service_account":    types.StringType,
			"binding":            types.ListType{ElemType: GcpAccessPolicyBindingModel{}.AttributeTypes()},
		},
	}
}

// GCP Access Policy -> Binding //

type GcpAccessPolicyBindingModel struct {
	Resource types.String `tfsdk:"resource"`
	Roles    types.Set    `tfsdk:"roles"`
}

func (g GcpAccessPolicyBindingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource": types.StringType,
			"roles":    types.SetType{ElemType: types.StringType},
		},
	}
}

// Azure Access Policy //

type AzureAccessPolicyModel struct {
	CloudAccountLink types.String `tfsdk:"cloud_account_link"`
	RoleAssignment   types.List   `tfsdk:"role_assignment"`
}

func (a AzureAccessPolicyModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cloud_account_link": types.StringType,
			"role_assignment":    types.ListType{ElemType: AzureAccessPolicyRoleAssignmentModel{}.AttributeTypes()},
		},
	}
}

// Azure Access Polic -> Role Assignment //

type AzureAccessPolicyRoleAssignmentModel struct {
	Scope types.String `tfsdk:"scope"`
	Roles types.Set    `tfsdk:"roles"`
}

func (g AzureAccessPolicyRoleAssignmentModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"scope": types.StringType,
			"roles": types.SetType{ElemType: types.StringType},
		},
	}
}

// NGS Access Policy //

type NgsAccessPolicyModel struct {
	CloudAccountLink types.String `tfsdk:"cloud_account_link"`
	Subs             types.Int32  `tfsdk:"subs"`
	Data             types.Int32  `tfsdk:"data"`
	Payload          types.Int32  `tfsdk:"payload"`
	Pub              types.List   `tfsdk:"pub"`
	Sub              types.List   `tfsdk:"sub"`
	Resp             types.List   `tfsdk:"resp"`
}

func (n NgsAccessPolicyModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cloud_account_link": types.StringType,
			"subs":               types.Int32Type,
			"data":               types.Int32Type,
			"payload":            types.Int32Type,
			"pub":                types.ListType{ElemType: NgsAccessPolicyPermissionModel{}.AttributeTypes()},
			"sub":                types.ListType{ElemType: NgsAccessPolicyPermissionModel{}.AttributeTypes()},
			"resp":               types.ListType{ElemType: NgsAccessPolicyResponsesModel{}.AttributeTypes()},
		},
	}
}

// NGS Access Policy -> Permission //

type NgsAccessPolicyPermissionModel struct {
	Allow types.Set `tfsdk:"allow"`
	Deny  types.Set `tfsdk:"deny"`
}

func (n NgsAccessPolicyPermissionModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"allow": types.SetType{ElemType: types.StringType},
			"deny":  types.SetType{ElemType: types.StringType},
		},
	}
}

// NGS Access Policy -> Responses //

type NgsAccessPolicyResponsesModel struct {
	Max types.Int32  `tfsdk:"max"`
	TTL types.String `tfsdk:"ttl"`
}

func (n NgsAccessPolicyResponsesModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"max": types.Int32Type,
			"ttl": types.StringType,
		},
	}
}

// Network Resource //

type NetworkResourceModel struct {
	Name       types.String `tfsdk:"name"`
	AgentLink  types.String `tfsdk:"agent_link"`
	IPs        types.Set    `tfsdk:"ips"`
	FQDN       types.String `tfsdk:"fqdn"`
	ResolverIP types.String `tfsdk:"resolver_ip"`
	Ports      types.Set    `tfsdk:"ports"`
}

func (n NetworkResourceModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":        types.StringType,
			"agent_link":  types.StringType,
			"ips":         types.SetType{ElemType: types.StringType},
			"fqdn":        types.StringType,
			"resolver_ip": types.StringType,
			"ports":       types.SetType{ElemType: types.Int32Type},
		},
	}
}

// Native Network Resource //

type NativeNetworkResourceModel struct {
	Name              types.String `tfsdk:"name"`
	FQDN              types.String `tfsdk:"fqdn"`
	Ports             types.Set    `tfsdk:"ports"`
	AwsPrivateLink    types.List   `tfsdk:"aws_private_link"`
	GcpServiceConnect types.List   `tfsdk:"gcp_service_connect"`
}

func (n NativeNetworkResourceModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":                types.StringType,
			"fqdn":                types.StringType,
			"ports":               types.SetType{ElemType: types.Int32Type},
			"aws_private_link":    types.ListType{ElemType: NativeNetworkResourceAwsPrivateLinkModel{}.AttributeTypes()},
			"gcp_service_connect": types.ListType{ElemType: NativeNetworkResourceGcpServiceConnectModel{}.AttributeTypes()},
		},
	}
}

// Native Network Resource -> AWS Private Link //

type NativeNetworkResourceAwsPrivateLinkModel struct {
	EndpointServiceName types.String `tfsdk:"endpoint_service_name"`
}

func (n NativeNetworkResourceAwsPrivateLinkModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"endpoint_service_name": types.StringType,
		},
	}
}

// Native Network Resource -> GCP Service Connect //

type NativeNetworkResourceGcpServiceConnectModel struct {
	TargetService types.String `tfsdk:"target_service"`
}

func (n NativeNetworkResourceGcpServiceConnectModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"target_service": types.StringType,
		},
	}
}
