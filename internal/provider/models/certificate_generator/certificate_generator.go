package certificate_generator

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Acme //

type AcmeModel struct {
	Server                 types.String `tfsdk:"server"`
	Email                  types.String `tfsdk:"email"`
	ExternalAccountBinding types.Object `tfsdk:"external_account_binding"`
	Challenge              types.Object `tfsdk:"challenge"`
}

func (a AcmeModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"server":                   types.StringType,
			"email":                    types.StringType,
			"external_account_binding": ExternalAccountBindingModel{}.AttributeTypes(),
			"challenge":                ChallengeModel{}.AttributeTypes(),
		},
	}
}

// Acme -> External Account Binding //

type ExternalAccountBindingModel struct {
	KeyId      types.String `tfsdk:"key_id"`
	SecretLink types.String `tfsdk:"secret_link"`
}

func (e ExternalAccountBindingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key_id":      types.StringType,
			"secret_link": types.StringType,
		},
	}
}

// Acme -> Challenge //

type ChallengeModel struct {
	Http01 types.Object `tfsdk:"http01"`
	Dns01  types.Object `tfsdk:"dns01"`
}

func (c ChallengeModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"http01": Http01Model{}.AttributeTypes(),
			"dns01":  Dns01Model{}.AttributeTypes(),
		},
	}
}

// Acme -> Challenge -> Http01 //

type Http01Model struct{}

func (h Http01Model) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{},
	}
}

// Acme -> Challenge -> Dns01 //

type Dns01Model struct {
	Provider   types.String `tfsdk:"provider"`
	Route53    types.Object `tfsdk:"route53"`
	CloudDns   types.Object `tfsdk:"cloud_dns"`
	AzureDns   types.Object `tfsdk:"azure_dns"`
	Cloudflare types.Object `tfsdk:"cloudflare"`
}

func (d Dns01Model) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"provider":   types.StringType,
			"route53":    Route53Model{}.AttributeTypes(),
			"cloud_dns":  CloudDnsModel{}.AttributeTypes(),
			"azure_dns":  AzureDnsModel{}.AttributeTypes(),
			"cloudflare": CloudflareModel{}.AttributeTypes(),
		},
	}
}

// Acme -> Challenge -> Dns01 -> Route53 //

type Route53Model struct {
	Region       types.String `tfsdk:"region"`
	HostedZoneId types.String `tfsdk:"hosted_zone_id"`
	SecretLink   types.String `tfsdk:"secret_link"`
}

func (r Route53Model) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":         types.StringType,
			"hosted_zone_id": types.StringType,
			"secret_link":    types.StringType,
		},
	}
}

// Acme -> Challenge -> Dns01 -> Cloud DNS //

type CloudDnsModel struct {
	Project    types.String `tfsdk:"project"`
	SecretLink types.String `tfsdk:"secret_link"`
}

func (c CloudDnsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"project":     types.StringType,
			"secret_link": types.StringType,
		},
	}
}

// Acme -> Challenge -> Dns01 -> Azure DNS //

type AzureDnsModel struct {
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	HostedZoneName    types.String `tfsdk:"hosted_zone_name"`
	SecretLink        types.String `tfsdk:"secret_link"`
}

func (a AzureDnsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource_group_name": types.StringType,
			"hosted_zone_name":    types.StringType,
			"secret_link":         types.StringType,
		},
	}
}

// Acme -> Challenge -> Dns01 -> Cloudflare //

type CloudflareModel struct {
	SecretLink types.String `tfsdk:"secret_link"`
}

func (c CloudflareModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"secret_link": types.StringType,
		},
	}
}

// Status //

type StatusModel struct {
	Ready       types.Bool   `tfsdk:"ready"`
	Message     types.String `tfsdk:"message"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (s StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ready":        types.BoolType,
			"message":      types.StringType,
			"last_updated": types.StringType,
		},
	}
}
