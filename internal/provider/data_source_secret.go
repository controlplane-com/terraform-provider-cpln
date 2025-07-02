package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &SecretDataSource{}
	_ datasource.DataSourceWithConfigure = &SecretDataSource{}
)

// SecretDataSource is the data source implementation.
type SecretDataSource struct {
	EntityBase
	Operations EntityOperations[SecretResourceModel, client.Secret]
}

// NewSecretDataSource returns a new instance of the data source implementation.
func NewSecretDataSource() datasource.DataSource {
	return &SecretDataSource{}
}

// Metadata provides the data source type name.
func (d *SecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_secret"
}

// Configure configures the data source before use.
func (d *SecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	d.Operations = NewEntityOperations(d.client, &SecretResourceOperator{})
}

// Schema defines the schema for the data source.
func (d *SecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this secret.",
				Computed:    true,
			},
			"cpln_id": schema.StringAttribute{
				Description: "The ID, in GUID format, of the secret.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of this secret.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of this secret.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value map of resource tags.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"self_link": schema.StringAttribute{
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"gcp": schema.StringAttribute{
				MarkdownDescription: "JSON string containing the GCP secret. [Reference Page](https://docs.controlplane.com/reference/secret#gcp)",
				Computed:            true,
				Sensitive:           true,
			},
			"docker": schema.StringAttribute{
				MarkdownDescription: "JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#docker).",
				Computed:            true,
				Sensitive:           true,
			},
			"dictionary": schema.MapAttribute{
				MarkdownDescription: "List of unique key-value pairs. [Reference Page](https://docs.controlplane.com/reference/secret#dictionary).",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"dictionary_as_envs": schema.MapAttribute{
				MarkdownDescription: "If a dictionary secret is defined, this output will be a key-value map in the following format: `key = cpln://secret/SECRET_NAME.key`.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"azure_sdk": schema.StringAttribute{
				MarkdownDescription: "JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#azure).",
				Computed:            true,
				Sensitive:           true,
			},
			"secret_link": schema.StringAttribute{
				Description: "Output used when linking a secret to an environment variable or volume.",
				Computed:    true,
			},
		},

		Blocks: map[string]schema.Block{
			"opaque": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#opaque).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"payload": schema.StringAttribute{
							Description: "Plain text or base64 encoded string. Use `encoding` attribute to specify encoding.",
							Computed:    true,
							Sensitive:   true,
						},
						"encoding": schema.StringAttribute{
							Description: "Available encodings: `plain`, `base64`. Default: `plain`.",
							Computed:    true,
						},
					},
				},
			},
			"tls": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#tls).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "Private Certificate.",
							Computed:    true,
							Sensitive:   true,
						},
						"cert": schema.StringAttribute{
							Description: "Public Certificate.",
							Computed:    true,
						},
						"chain": schema.StringAttribute{
							Description: "Chain Certificate.",
							Computed:    true,
						},
					},
				},
			},
			"aws": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#aws).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"access_key": schema.StringAttribute{
							Description: "Access Key provided by AWS.",
							Computed:    true,
							Sensitive:   true,
						},
						"secret_key": schema.StringAttribute{
							Description: "Secret Key provided by AWS.",
							Computed:    true,
							Sensitive:   true,
						},
						"role_arn": schema.StringAttribute{
							Description: "Role ARN provided by AWS.",
							Computed:    true,
						},
						"external_id": schema.StringAttribute{
							Description: "AWS IAM Role External ID.",
							Computed:    true,
						},
					},
				},
			},
			"ecr": schema.ListNestedBlock{
				MarkdownDescription: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"access_key": schema.StringAttribute{
							Description: "Access Key provided by AWS.",
							Computed:    true,
							Sensitive:   true,
						},
						"secret_key": schema.StringAttribute{
							Description: "Secret Key provided by AWS.",
							Computed:    true,
							Sensitive:   true,
						},
						"role_arn": schema.StringAttribute{
							Description: "Role ARN provided by AWS.",
							Computed:    true,
						},
						"external_id": schema.StringAttribute{
							Description: "AWS IAM Role External ID. Used when setting up cross-account access to your ECR repositories.",
							Optional:    true,
						},
						"repos": schema.SetAttribute{
							Description: "List of ECR repositories.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
			"userpass": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#username).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							Description: "Username.",
							Computed:    true,
						},
						"password": schema.StringAttribute{
							Description: "Password.",
							Computed:    true,
							Sensitive:   true,
						},
						"encoding": schema.StringAttribute{
							Description: "Available encodings: `plain`, `base64`. Default: `plain`.",
							Computed:    true,
						},
					},
				},
			},
			"keypair": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#keypair).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"secret_key": schema.StringAttribute{
							Description: "Secret/Private Key.",
							Computed:    true,
							Sensitive:   true,
						},
						"public_key": schema.StringAttribute{
							Description: "Public Key.",
							Computed:    true,
						},
						"passphrase": schema.StringAttribute{
							Description: "Passphrase for private key.",
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"azure_connector": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#azure-connector).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: "Deployment URL.",
							Computed:    true,
							Sensitive:   true,
						},
						"code": schema.StringAttribute{
							Description: "Code/Key to authenticate to deployment URL.",
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"nats_account": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#nats-account).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Description: "Account ID.",
							Computed:    true,
						},
						"private_key": schema.StringAttribute{
							Description: "Private Key.",
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

// Read fetches the current state of the resource.
func (d *SecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state SecretResourceModel

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Create a new operator instance
	operator := d.Operations.NewOperator(ctx, &resp.Diagnostics, state)

	// Invoke API to read resource details
	apiResp, code, err := operator.InvokeRead(state.Name.ValueString())

	// Remove resource from state if not found
	if code == 404 {
		// Drop resource from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Build new state from API response
	newState := operator.MapResponseToState(apiResp, true)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
