package secret

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Opaque //

type OpaqueModel struct {
	Payload  types.String `tfsdk:"payload"`
	Encoding types.String `tfsdk:"encoding"`
}

// TLS //

type TlsModel struct {
	Key   types.String `tfsdk:"key"`
	Cert  types.String `tfsdk:"cert"`
	Chain types.String `tfsdk:"chain"`
}

// AWS //

type AwsModel struct {
	AccessKey  types.String `tfsdk:"access_key"`
	SecretKey  types.String `tfsdk:"secret_key"`
	RoleArn    types.String `tfsdk:"role_arn"`
	ExternalId types.String `tfsdk:"external_id"`
}

// ECR //

type EcrModel struct {
	AccessKey  types.String `tfsdk:"access_key"`
	SecretKey  types.String `tfsdk:"secret_key"`
	RoleArn    types.String `tfsdk:"role_arn"`
	ExternalId types.String `tfsdk:"external_id"`
	Repos      types.Set    `tfsdk:"repos"`
}

// Userpass //

type UserpassModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Encoding types.String `tfsdk:"encoding"`
}

// Key Pair //

type KeyPairModel struct {
	SecretKey  types.String `tfsdk:"secret_key"`
	PublicKey  types.String `tfsdk:"public_key"`
	Passphrase types.String `tfsdk:"passphrase"`
}

// Azure Connector //

type AzureConnectorModel struct {
	Url  types.String `tfsdk:"url"`
	Code types.String `tfsdk:"code"`
}

// NATS Account //

type NatsAccountModel struct {
	AccountId  types.String `tfsdk:"account_id"`
	PrivateKey types.String `tfsdk:"private_key"`
}
