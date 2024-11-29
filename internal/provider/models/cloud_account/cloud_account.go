package cloud_account

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AWS //

type Aws struct {
	RoleArn types.String `tfsdk:"role_arn"`
}

func (c Aws) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"role_arn": types.StringType,
		},
	}
}

// Azure //

type Azure struct {
	SecretLink types.String `tfsdk:"secret_link"`
}

func (c Azure) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"secret_link": types.StringType,
		},
	}
}

// GCP //

type Gcp struct {
	ProjectId types.String `tfsdk:"project_id"`
}

func (c Gcp) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"project_id": types.StringType,
		},
	}
}

// NGS //

type Ngs struct {
	SecretLink types.String `tfsdk:"secret_link"`
}

func (c Ngs) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"secret_link": types.StringType,
		},
	}
}
