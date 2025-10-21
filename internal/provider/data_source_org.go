package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &OrgDataSource{}
	_ datasource.DataSourceWithConfigure = &OrgDataSource{}
)

// OrgDataSource is the data source implementation.
type OrgDataSource struct {
	EntityBase
	Operations EntityOperations[OrgResourceModel, client.Org]
}

// NewOrgDataSource returns a new instance of the data source implementation.
func NewOrgDataSource() datasource.DataSource {
	return &OrgDataSource{}
}

// Metadata provides the data source type name.
func (d *OrgDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_org"
}

// Configure configures the data source before use.
func (d *OrgDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	d.Operations = NewEntityOperations(d.client, &OrgResourceOperator{})
}

// Schema defines the schema for the data source.
func (d *OrgDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this org.",
				Computed:    true,
			},
			"cpln_id": schema.StringAttribute{
				Description: "The ID, in GUID format, of the org.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of this org.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of this org.",
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
			"account_id": schema.StringAttribute{
				Description: "The associated account ID that will be used when creating the org. Only used on org creation. The account ID can be obtained from the `Org Management & Billing` page.",
				Computed:    true,
			},
			"invitees": schema.SetAttribute{
				Description: "When an org is created, the list of email addresses which will receive an invitation to join the org and be assigned to the `superusers` group. The user account used when creating the org will be included in this list.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"session_timeout_seconds": schema.Int32Attribute{
				Description: "The idle time (in seconds) in which the console UI will automatically sign-out the user. Default: 900 (15 minutes)",
				Computed:    true,
			},
			"status": schema.ListNestedAttribute{
				Description: "Status of the org.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_link": schema.StringAttribute{
							Description: "The link of the account the org belongs to.",
							Computed:    true,
						},
						"active": schema.BoolAttribute{
							Description: "Indicates whether the org is active or not.",
							Computed:    true,
						},
						"endpoint_prefix": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"auth_config": schema.ListNestedBlock{
				Description: "The configuration settings and parameters related to authentication within the org.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"domain_auto_members": schema.SetAttribute{
							Description: "List of domains which will auto-provision users when authenticating using SAML.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"saml_only": schema.BoolAttribute{
							Description: "Enforce SAML only authentication.",
							Computed:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"observability": schema.ListNestedBlock{
				Description: "The retention period (in days) for logs, metrics, and traces. Charges apply for storage beyond the 30 day default.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"logs_retention_days": schema.Int32Attribute{
							Description: "Log retention days. Default: 30",
							Computed:    true,
						},
						"metrics_retention_days": schema.Int32Attribute{
							Description: "Metrics retention days. Default: 30",
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
							},
						},
						"traces_retention_days": schema.Int32Attribute{
							Description: "Traces retention days. Default: 30",
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
							},
						},
						"default_alert_emails": schema.SetAttribute{
							Description: "These emails are configured as alert recipients in Grafana when the 'grafana-default-email' contact delivery type is 'Email'.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"security": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"threat_detection": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "Indicates whether threat detection should be forwarded or not.",
										Computed:    true,
									},
									"minimum_severity": schema.StringAttribute{
										Description: "Any threats with this severity and more severe will be sent. Others will be ignored. Valid values: `warning`, `error`, or `critical`.",
										Optional:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"syslog": schema.ListNestedBlock{
										Description: "Configuration for syslog forwarding.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"transport": schema.StringAttribute{
													Description: "The transport-layer protocol to send the syslog messages over. If TCP is chosen, messages will be sent with TLS. Default: `tcp`.",
													Computed:    true,
												},
												"host": schema.StringAttribute{
													Description: "The hostname to send syslog messages to.",
													Computed:    true,
												},
												"port": schema.Int32Attribute{
													Description: "The port to send syslog messages to.",
													Computed:    true,
												},
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtMost(1),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

// Read fetches the current state of the resource.
func (d *OrgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state OrgResourceModel

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
	apiResp, code, err := operator.InvokeRead(d.client.Org)

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
