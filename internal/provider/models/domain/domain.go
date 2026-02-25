package domain

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Domain ***/

// Spec //

type SpecModel struct {
	DnsMode             types.String `tfsdk:"dns_mode"`
	GvcLink             types.String `tfsdk:"gvc_link"`
	CertChallengeType   types.String `tfsdk:"cert_challenge_type"`
	WorkloadLink        types.String `tfsdk:"workload_link"`
	AcceptAllHosts      types.Bool   `tfsdk:"accept_all_hosts"`
	AcceptAllSubdomains types.Bool   `tfsdk:"accept_all_subdomains"`
	Ports               types.List   `tfsdk:"ports"`
}

func (s SpecModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"dns_mode":              types.StringType,
			"gvc_link":              types.StringType,
			"cert_challenge_type":   types.StringType,
			"workload_link":         types.StringType,
			"accept_all_hosts":      types.BoolType,
			"accept_all_subdomains": types.BoolType,
			"ports":                 types.ListType{ElemType: SpecPortsModel{}.AttributeTypes()},
		},
	}
}

// Spec -> Ports //

type SpecPortsModel struct {
	Number   types.Int32  `tfsdk:"number"`
	Protocol types.String `tfsdk:"protocol"`
	Cors     types.List   `tfsdk:"cors"`
	TLS      types.List   `tfsdk:"tls"`
	Route    types.List   `tfsdk:"route"`
}

func (s SpecPortsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"number":   types.Int32Type,
			"protocol": types.StringType,
			"cors":     types.ListType{ElemType: SpecPortsCorsModel{}.AttributeTypes()},
			"tls":      types.ListType{ElemType: SpecPortsTlsModel{}.AttributeTypes()},
			"route":    types.ListType{ElemType: RouteModel{}.AttributeTypes()},
		},
	}
}

// Spec -> Ports -> Cors //

type SpecPortsCorsModel struct {
	AllowOrigins     types.Set    `tfsdk:"allow_origins"`
	AllowMethods     types.Set    `tfsdk:"allow_methods"`
	AllowHeaders     types.Set    `tfsdk:"allow_headers"`
	ExposeHeaders    types.Set    `tfsdk:"expose_headers"`
	MaxAge           types.String `tfsdk:"max_age"`
	AllowCredentials types.Bool   `tfsdk:"allow_credentials"`
}

func (s SpecPortsCorsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"allow_origins":     types.SetType{ElemType: SpecPortsCorsAllowOriginsModel{}.AttributeTypes()},
			"allow_methods":     types.SetType{ElemType: types.StringType},
			"allow_headers":     types.SetType{ElemType: types.StringType},
			"expose_headers":    types.SetType{ElemType: types.StringType},
			"max_age":           types.StringType,
			"allow_credentials": types.BoolType,
		},
	}
}

// Spec -> Ports -> Cors -> AllowOrigins //

type SpecPortsCorsAllowOriginsModel struct {
	Exact types.String `tfsdk:"exact"`
	Regex types.String `tfsdk:"regex"`
}

func (s SpecPortsCorsAllowOriginsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"exact": types.StringType,
			"regex": types.StringType,
		},
	}
}

// Spec -> Ports -> TLS //

type SpecPortsTlsModel struct {
	MinProtocolVersion types.String `tfsdk:"min_protocol_version"`
	CipherSuites       types.Set    `tfsdk:"cipher_suites"`
	ClientCertificate  types.List   `tfsdk:"client_certificate"`
	ServerCertificate  types.List   `tfsdk:"server_certificate"`
}

func (s SpecPortsTlsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_protocol_version": types.StringType,
			"cipher_suites":        types.SetType{ElemType: types.StringType},
			"client_certificate":   types.ListType{ElemType: SpecPortsTlsCertificateModel{}.AttributeTypes()},
			"server_certificate":   types.ListType{ElemType: SpecPortsTlsCertificateModel{}.AttributeTypes()},
		},
	}
}

// Spec -> Ports -> TLS -> Certificate //

type SpecPortsTlsCertificateModel struct {
	SecretLink types.String `tfsdk:"secret_link"`
}

func (s SpecPortsTlsCertificateModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"secret_link": types.StringType,
		},
	}
}

// Status //

type StatusModel struct {
	Endpoints   types.List   `tfsdk:"endpoints"`
	Status      types.String `tfsdk:"status"`
	Warning     types.String `tfsdk:"warning"`
	Locations   types.List   `tfsdk:"locations"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	DnsConfig   types.List   `tfsdk:"dns_config"`
}

func (s StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"status":      types.StringType,
			"warning":     types.StringType,
			"fingerprint": types.StringType,
			"endpoints":   types.ListType{ElemType: StatusEndpointModel{}.AttributeTypes()},
			"locations":   types.ListType{ElemType: StatusLocationModel{}.AttributeTypes()},
			"dns_config":  types.ListType{ElemType: StatusDnsConfigModel{}.AttributeTypes()},
		},
	}
}

// Status -> Endpoints //

type StatusEndpointModel struct {
	URL          types.String `tfsdk:"url"`
	WorkloadLink types.String `tfsdk:"workload_link"`
}

func (s StatusEndpointModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"url":           types.StringType,
			"workload_link": types.StringType,
		},
	}
}

// Status -> Locations //

type StatusLocationModel struct {
	Name              types.String `tfsdk:"name"`
	CertificateStatus types.String `tfsdk:"certificate_status"`
}

func (s StatusLocationModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":               types.StringType,
			"certificate_status": types.StringType,
		},
	}
}

// Status -> DnsConfig //

type StatusDnsConfigModel struct {
	Type  types.String `tfsdk:"type"`
	TTL   types.Int32  `tfsdk:"ttl"`
	Host  types.String `tfsdk:"host"`
	Value types.String `tfsdk:"value"`
}

func (s StatusDnsConfigModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":  types.StringType,
			"ttl":   types.Int32Type,
			"host":  types.StringType,
			"value": types.StringType,
		},
	}
}

/*** Domain Route ***/

// Headers //

type RouteHeadersModel struct {
	Request types.List `tfsdk:"request"`
}

func (r RouteHeadersModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"request": types.ListType{ElemType: RouteHeadersRequestModel{}.AttributeTypes()},
		},
	}
}

// Headers -> Request //

type RouteHeadersRequestModel struct {
	Set types.Map `tfsdk:"set"`
}

func (r RouteHeadersRequestModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"set": types.MapType{ElemType: types.StringType},
		},
	}
}

// Route //

type RouteModel struct {
	Prefix        types.String `tfsdk:"prefix"`
	ReplacePrefix types.String `tfsdk:"replace_prefix"`
	Regex         types.String `tfsdk:"regex"`
	WorkloadLink  types.String `tfsdk:"workload_link"`
	Port          types.Int32  `tfsdk:"port"`
	HostPrefix    types.String `tfsdk:"host_prefix"`
	HostRegex     types.String `tfsdk:"host_regex"`
	Headers       types.List   `tfsdk:"headers"`
	Replica       types.Int32  `tfsdk:"replica"`
}

func (r RouteModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"prefix":         types.StringType,
			"replace_prefix": types.StringType,
			"regex":          types.StringType,
			"workload_link":  types.StringType,
			"port":           types.Int32Type,
			"host_prefix":    types.StringType,
			"host_regex":     types.StringType,
			"headers":        types.ListType{ElemType: RouteHeadersModel{}.AttributeTypes()},
			"replica":        types.Int32Type,
		},
	}
}
