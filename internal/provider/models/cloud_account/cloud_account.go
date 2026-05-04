package cloud_account

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AWS //

type AwsModel struct {
	RoleArn types.String `tfsdk:"role_arn"`
}

func (c AwsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"role_arn": types.StringType,
		},
	}
}

// Azure //

type AzureModel struct {
	SecretLink types.String `tfsdk:"secret_link"`
}

func (c AzureModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"secret_link": types.StringType,
		},
	}
}

// GCP //

type GcpModel struct {
	ProjectId types.String `tfsdk:"project_id"`
}

func (c GcpModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"project_id": types.StringType,
		},
	}
}

// NGS //

type NgsModel struct {
	SecretLink types.String `tfsdk:"secret_link"`
}

func (c NgsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"secret_link": types.StringType,
		},
	}
}

// Status //

type StatusModel struct {
	Usable      types.Bool   `tfsdk:"usable"`
	LastChecked types.String `tfsdk:"last_checked"`
	LastError   types.String `tfsdk:"last_error"`
}

func (c StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"usable":       types.BoolType,
			"last_checked": types.StringType,
			"last_error":   types.StringType,
		},
	}
}
