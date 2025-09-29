package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &GvcDataSource{}
	_ datasource.DataSourceWithConfigure = &GvcDataSource{}
)

// GvcDataSource is the data source implementation.
type GvcDataSource struct {
	EntityBase
	Operations EntityOperations[GvcResourceModel, client.Gvc]
}

// NewGvcDataSource returns a new instance of the data source implementation.
func NewGvcDataSource() datasource.DataSource {
	return &GvcDataSource{}
}

// Metadata provides the data source type name.
func (d *GvcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_gvc"
}

// Configure configures the data source before use.
func (d *GvcDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	d.Operations = NewEntityOperations(d.client, &GvcResourceOperator{})
}

// Schema defines the schema for the data source.
func (d *GvcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this GVC.",
				Computed:    true,
			},
			"cpln_id": schema.StringAttribute{
				Description: "The ID, in GUID format, of the GVC.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the GVC.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the GVC.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value map of resource tags.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"self_link": schema.StringAttribute{
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"alias": schema.StringAttribute{
				Description: "The alias name of the GVC.",
				Computed:    true,
			},
			"locations": schema.SetAttribute{
				MarkdownDescription: "A list of [locations](https://docs.controlplane.com/reference/location#current) making up the Global Virtual Cloud.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"pull_secrets": schema.SetAttribute{
				MarkdownDescription: "A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"domain": schema.StringAttribute{
				Description:        "Custom domain name used by associated workloads.",
				DeprecationMessage: "Selecting a domain on a GVC will be deprecated in the future. Use the 'cpln_domain resource' instead.",
				Optional:           true,
			},
			"endpoint_naming_format": schema.StringAttribute{
				Description: "Customizes the subdomain format for the canonical workload endpoint. `default` leaves it as '${workloadName}-${gvcName}.cpln.app'. `org` follows the scheme '${workloadName}-${gvcName}.${org}.cpln.app'.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("default", "org"),
				},
			},
			"env": schema.MapAttribute{
				Description: "Key-value array of resource environment variables.",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"lightstep_tracing":    d.LightstepTracingSchema(),
			"otel_tracing":         d.OtelTracingSchema(),
			"controlplane_tracing": d.ControlPlaneTracingSchema(),
			"sidecar": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"envoy": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"load_balancer": schema.ListNestedBlock{
				Description: "Dedicated load balancer configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"dedicated": schema.BoolAttribute{
							Description: "Creates a dedicated load balancer in each location and enables additional Domain features: custom ports, protocols and wildcard hostnames. Charges apply for each location.",
							Optional:    true,
						},
						"trusted_proxies": schema.Int32Attribute{
							Description: "Controls the address used for request logging and for setting the X-Envoy-External-Address header. If set to 1, then the last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If set to 2, then the second to last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If the XFF header does not have at least two addresses or does not exist then the source client IP address will be used instead.",
							Optional:    true,
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
								int32validator.AtMost(2),
							},
						},
						"ipset": schema.StringAttribute{
							Description: "The link or the name of the IP Set that will be used for this load balancer.",
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"multi_zone": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"redirect": schema.ListNestedBlock{
							Description: "Specify the url to be redirected to for different http status codes.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"class": schema.ListNestedBlock{
										Description: "Specify the redirect url for all status codes in a class.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"status_5xx": schema.StringAttribute{
													Description: "Specify the redirect url for any 500 level status code.",
													Optional:    true,
												},
												"status_401": schema.StringAttribute{
													Description: "An optional url redirect for 401 responses. Supports envoy format strings to include request information. E.g. https://your-oauth-server/oauth2/authorize?return_to=%REQ(:path)%&client_id=your-client-id",
													Optional:    true,
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
			"keda": schema.ListNestedBlock{
				Description: "KEDA configuration for the GVC.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: "Enable KEDA for this GVC. KEDA is a Kubernetes-based event-driven autoscaler that allows you to scale workloads based on external events. When enabled, a keda operator will be deployed in the GVC and workloads in the GVC can use KEDA to scale based on external metrics.",
							Optional:    true,
							Computed:    true,
						},
						"identity_link": schema.StringAttribute{
							Description: "A link to an Identity resource that will be used for KEDA. This will allow the keda operator to access cloud and network resources.",
							Optional:    true,
							Validators: []validator.String{
								validators.LinkValidator{},
							},
						},
						"secrets": schema.SetAttribute{
							Description: "A list of secrets to be used as TriggerAuthentication objects. The TriggerAuthentication object will be named after the secret and can be used by triggers on workloads in this GVC.",
							ElementType: types.StringType,
							Optional:    true,
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
func (d *GvcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state GvcResourceModel

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
