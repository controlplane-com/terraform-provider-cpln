package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/workload"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                   = &WorkloadResource{}
	_ resource.ResourceWithImportState    = &WorkloadResource{}
	_ resource.ResourceWithValidateConfig = &WorkloadResource{}
)

/*** Resource Model ***/

// WorkloadResourceModel holds the Terraform state for the resource.
type WorkloadResourceModel struct {
	EntityBaseModel
	Gvc                types.String                     `tfsdk:"gvc"`
	Type               types.String                     `tfsdk:"type"`
	IdentityLink       types.String                     `tfsdk:"identity_link"`
	Containers         []models.ContainerModel          `tfsdk:"container"`
	Firewall           []models.FirewallModel           `tfsdk:"firewall_spec"`
	Options            []models.OptionsModel            `tfsdk:"options"`
	LocalOptions       []models.LocalOptionsModel       `tfsdk:"local_options"`
	Job                []models.JobModel                `tfsdk:"job"`
	Sidecar            []models.SidecarModel            `tfsdk:"sidecar"`
	SupportDynamicTags types.Bool                       `tfsdk:"support_dynamic_tags"`
	RolloutOptions     []models.RolloutOptionsModel     `tfsdk:"rollout_options"`
	SecurityOptions    []models.SecurityOptionsModel    `tfsdk:"security_options"`
	LoadBalancer       []models.LoadBalancerModel       `tfsdk:"load_balancer"`
	Extras             types.String                     `tfsdk:"extras"`
	RequestRetryPolicy []models.RequestRetryPolicyModel `tfsdk:"request_retry_policy"`
	Status             types.List                       `tfsdk:"status"`
}

/*** Resource Configuration ***/

// WorkloadResource is the resource implementation.
type WorkloadResource struct {
	EntityBase
	Operations EntityOperations[WorkloadResourceModel, client.Workload]
}

// NewWorkloadResource returns a new instance of the resource implementation.
func NewWorkloadResource() resource.Resource {
	return &WorkloadResource{}
}

// Configure configures the resource before use.
func (wr *WorkloadResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	wr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	wr.Operations = NewEntityOperations(wr.client, &WorkloadResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (wr *WorkloadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the import ID
	parts := strings.SplitN(req.ID, ":", 2)

	// Validate that ID has exactly three non-empty segments
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		// Report error when import identifier format is unexpected
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: "+
					"'gvc:workload_name'. Got: %q", req.ID,
			),
		)

		// Abort import operation on error
		return
	}

	// Extract gvc and workloadName from parts
	gvc, workloadName := parts[0], parts[1]

	// Set the generated ID attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(workloadName))...,
	)

	// Set the GVC attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("gvc"), types.StringValue(gvc))...,
	)
}

// Metadata provides the resource type name.
func (wr *WorkloadResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_workload"
}

// Schema defines the schema for the resource.
func (wr *WorkloadResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(wr.EntityBaseAttributes("workload"), map[string]schema.Attribute{
			"gvc": schema.StringAttribute{
				Description: "Name of the associated GVC.",
				Required:    true,
				Validators: []validator.String{
					validators.NameValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Workload Type. Either `serverless`, `standard`, `stateful`, or `cron`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("serverless", "standard", "stateful", "cron"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identity_link": schema.StringAttribute{
				Description: "The identityLink is used as the access scope for 3rd party cloud resources. A single identity can provide access to multiple cloud providers.",
				Optional:    true,
				Validators: []validator.String{
					validators.LinkValidator{},
				},
			},
			"support_dynamic_tags": schema.BoolAttribute{
				Description: "Workload will automatically redeploy when one of the container images is updated in the container registry. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"extras": schema.StringAttribute{
				Description: "Extra Kubernetes modifications. Only used for BYOK.",
				Optional:    true,
			},
			"status": schema.ListNestedAttribute{
				Description: "Status of the workload.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"parent_id": schema.StringAttribute{
							Description: "ID of the parent object.",
							Computed:    true,
						},
						"canonical_endpoint": schema.StringAttribute{
							Description: "Canonical endpoint for the workload.",
							Computed:    true,
						},
						"endpoint": schema.StringAttribute{
							Description: "Endpoint for the workload.",
							Computed:    true,
						},
						"internal_name": schema.StringAttribute{
							Description: "Internal hostname for the workload. Used for service-to-service requests.",
							Computed:    true,
						},
						"health_check": schema.ListNestedAttribute{
							Description: "Current health status.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"active": schema.BoolAttribute{
										Description: "Active boolean for the associated workload.",
										Computed:    true,
									},
									"success": schema.BoolAttribute{
										Description: "Success boolean for the associated workload.",
										Computed:    true,
									},
									"code": schema.Int32Attribute{
										Description: "Current output code for the associated workload.",
										Computed:    true,
									},
									"message": schema.StringAttribute{
										Description: "Current health status for the associated workload.",
										Computed:    true,
									},
									"failures": schema.Int32Attribute{
										Description: "Failure integer for the associated workload.",
										Computed:    true,
									},
									"successes": schema.Int32Attribute{
										Description: "Success integer for the associated workload.",
										Computed:    true,
									},
									"last_checked": schema.StringAttribute{
										Description: "Timestamp in UTC of the last health check.",
										Computed:    true,
									},
								},
							},
						},
						"current_replica_count": schema.Int32Attribute{
							Description: "Current amount of replicas deployed.",
							Computed:    true,
						},
						"resolved_images": schema.ListNestedAttribute{
							Description: "Resolved images for workloads with dynamic tags enabled.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"resolved_for_version": schema.Int32Attribute{
										Description: "Workload version the images were resolved for.",
										Computed:    true,
									},
									"resolved_at": schema.StringAttribute{
										Description: "UTC Time when the images were resolved.",
										Computed:    true,
									},
									"error_messages": schema.SetAttribute{
										Description: "",
										ElementType: types.StringType,
										Computed:    true,
									},
									"images": schema.ListNestedAttribute{
										Description: "A list of images that were resolved.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"digest": schema.StringAttribute{
													Description: "A unique SHA256 hash value that identifies a specific image content. This digest serves as a fingerprint of the image's content, ensuring the image you pull or run is exactly what you expect, without any modifications or corruptions.",
													Computed:    true,
												},
												"manifests": schema.ListNestedAttribute{
													Description: "",
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"image": schema.StringAttribute{
																Description: "The name and tag of the resolved image.",
																Computed:    true,
															},
															"media_type": schema.StringAttribute{
																Description: "The MIME type used in the Docker Registry HTTP API to specify the format of the data being sent or received. Docker uses media types to distinguish between different kinds of JSON objects and binary data formats within the registry protocol, enabling the Docker client and registry to understand and process different components of Docker images correctly.",
																Computed:    true,
															},
															"digest": schema.StringAttribute{
																Description: "A SHA256 hash that uniquely identifies the specific image manifest.",
																Computed:    true,
															},
															"platform": schema.MapAttribute{
																Description: "Key-value map of strings. The combination of the operating system and architecture for which the image is built.",
																ElementType: types.StringType,
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"load_balancer": schema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"origin": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"url": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"container": schema.ListNestedBlock{
				Description: "An isolated and lightweight runtime environment that encapsulates an application and its dependencies.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the container.",
							Required:    true,
							Validators: []validator.String{
								validators.NameValidator{},
								validators.DisallowPrefixValidator{Prefix: "cpln-"},
								validators.DisallowListValidator{
									Forbidden: []string{
										"istio-proxy",
										"queue-proxy",
										"istio-validation",
										"cpln-envoy-assassin",
										"cpln-writer-proxy",
										"cpln-reader-proxy",
										"cpln-dbaas-config",
									},
								},
							},
						},
						"image": schema.StringAttribute{
							Description: "The full image and tag path.",
							Required:    true,
						},
						"working_directory": schema.StringAttribute{
							Description: "Override the working directory. Must be an absolute path.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^/.*`),
									"must be an absolute (slash-prefixed) path",
								),
							},
						},
						"port": schema.Int32Attribute{
							Description:        "The port the container exposes. Only one container is allowed to specify a port. Min: `80`. Max: `65535`. Used by `serverless` Workload type. **DEPRECATED - Use `ports`.**",
							DeprecationMessage: "The 'port' attribute will be deprecated in the next major version. Use the 'ports' attribute instead.",
							Optional:           true,
							Validators:         wr.GetPortValidators(),
						},
						"memory": schema.StringAttribute{
							Description: "Reserved memory of the workload when capacityAI is disabled. Maximum memory when CapacityAI is enabled. Default: \"128Mi\".",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("128Mi"),
							Validators:  wr.GetCpuMemoryValidators("must be a valid memory quantity"),
						},
						"cpu": schema.StringAttribute{
							Description: "Reserved CPU of the workload when capacityAI is disabled. Maximum CPU when CapacityAI is enabled. Default: \"50m\".",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("50m"),
							Validators:  wr.GetCpuMemoryValidators("must be a valid CPU quantity"),
						},
						"min_cpu": schema.StringAttribute{
							Description: "Minimum CPU when capacity AI is enabled.",
							Optional:    true,
							Validators:  wr.GetCpuMemoryValidators("must be a valid CPU quantity"),
						},
						"min_memory": schema.StringAttribute{
							Description: "Minimum memory when capacity AI is enabled.",
							Optional:    true,
							Validators:  wr.GetCpuMemoryValidators("must be a valid memory quantity"),
						},
						"env": schema.MapAttribute{
							Description: "Name-Value list of environment variables.",
							ElementType: types.StringType,
							Optional:    true,
							Validators: []validator.Map{
								mapvalidator.KeysAre(
									stringvalidator.NoneOfCaseInsensitive(
										"k_service",
										"k_configuration",
										"k_revision",
									),
								),
								mapvalidator.ValueStringsAre(
									stringvalidator.LengthAtMost(4 * 1024),
								),
							},
						},
						"inherit_env": schema.BoolAttribute{
							Description: "Enables inheritance of GVC environment variables. A variable in spec.env will override a GVC variable with the same name.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"command": schema.StringAttribute{
							Description: "Override the entry point.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(256),
							},
						},
						"args": schema.SetAttribute{
							Description: "Command line arguments passed to the container at runtime. Replaces the CMD arguments of the running container. It is an ordered list.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"metrics": schema.ListNestedBlock{
							MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/workload#metrics).",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"port": schema.Int32Attribute{
										Description: "Port from container emitting custom metrics.",
										Required:    true,
										Validators:  wr.GetPortValidators(),
									},
									"path": schema.StringAttribute{
										Description: "Path from container emitting custom metrics.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(128),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"ports": schema.ListNestedBlock{
							Description: "Communication endpoints used by the workload to send and receive network traffic.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"protocol": schema.StringAttribute{
										Description: "Protocol. Choice of: `http`, `http2`, `tcp`, or `grpc`.",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("http"),
										Validators: []validator.String{
											stringvalidator.OneOf("http", "http2", "tcp", "grpc"),
										},
									},
									"number": schema.Int32Attribute{
										Description: "Port to expose.",
										Required:    true,
										Validators:  wr.GetPortValidators(),
									},
								},
							},
						},
						"readiness_probe": wr.HealthCheckSchema("readiness_probe", "Readiness Probe"),
						"liveness_probe":  wr.HealthCheckSchema("liveness_probe", "Liveness Probe"),
						"gpu_nvidia": schema.ListNestedBlock{
							Description: "GPUs manufactured by NVIDIA, which are specialized hardware accelerators used to offload and accelerate computationally intensive tasks within the workload.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"model": schema.StringAttribute{
										Description: "GPU Model (i.e.: t4)",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("t4", "a10g"),
										},
									},
									"quantity": schema.Int32Attribute{
										Description: "Number of GPUs.",
										Required:    true,
										Validators: []validator.Int32{
											int32validator.Between(0, 4),
										},
									},
								},
								Validators: []validator.Object{
									objectvalidator.ConflictsWith(
										path.MatchRelative().AtParent().AtParent().AtName("gpu_custom"),
									),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"gpu_custom": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"resource": schema.StringAttribute{
										Description: "",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(64),
											stringvalidator.RegexMatches(
												regexp.MustCompile(`^[a-zA-Z0-9./]*$`),
												"must be a valid resource name",
											),
										},
									},
									"runtime_class": schema.StringAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(64),
											stringvalidator.RegexMatches(
												regexp.MustCompile(`^[a-zA-Z0-9./]*$`),
												"must be a valid runtime class",
											),
										},
									},
									"quantity": schema.Int32Attribute{
										Description: "Number of GPUs.",
										Required:    true,
										Validators: []validator.Int32{
											int32validator.Between(0, 8),
										},
									},
								},
								Validators: []validator.Object{
									objectvalidator.ConflictsWith(
										path.MatchRelative().AtParent().AtParent().AtName("gpu_nvidia"),
									),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"lifecycle": schema.ListNestedBlock{
							MarkdownDescription: "Lifecycle [Reference Page](https://docs.controlplane.com/reference/workload#lifecycle).",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"post_start": wr.LifecycleSpecSchema("Command and arguments executed immediately after the container is created."),
									"pre_stop":   wr.LifecycleSpecSchema("Command and arguments executed immediately before the container is stopped."),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"volume": schema.ListNestedBlock{
							MarkdownDescription: "Mount Object Store (S3, GCS, AzureBlob) buckets as file system.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"uri": schema.StringAttribute{
										Description: "URI of a volume hosted at Control Plane (Volume Set) or at a cloud provider (AWS, Azure, GCP).",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(`^(s3|gs|azureblob|azurefs|cpln|scratch):\/\/.+`),
												"must be in the form s3://bucket, gs://bucket, azureblob://storageAccount/container, azurefs://storageAccount/share, cpln://, or scratch://",
											),
										},
									},
									"recovery_policy": schema.StringAttribute{
										Description: "Only applicable to persistent volumes, this determines what Control Plane will do when creating a new workload replica if a corresponding volume exists. Available Values: `retain`, `recycle`. Default: `retain`. **DEPRECATED - No longer being used.**",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("retain"),
										Validators: []validator.String{
											stringvalidator.OneOf("retain", "recycle"),
										},
									},
									"path": schema.StringAttribute{
										Description: "File path added to workload pointing to the volume.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(`^/.*`),
												"must be an absolute path (start with \"/\")",
											),
											stringvalidator.NoneOf(
												"/dev", "/dev/",
												"/dev/log", "/dev/log/",
												"/tmp", "/tmp/",
												"/var", "/var/",
												"/var/log", "/var/log/",
											),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(15),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(8),
				},
			},
			"firewall_spec": schema.ListNestedBlock{
				Description: "Control of inbound and outbound access to the workload for external (public) and internal (service to service) traffic. Access is restricted by default.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"external": schema.ListNestedBlock{
							Description: "The external firewall is used to control inbound and outbound access to the workload for public-facing traffic.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"inbound_allow_cidr": schema.SetAttribute{
										Description: "The list of ipv4/ipv6 addresses or cidr blocks that are allowed to access this workload. No external access is allowed by default. Specify '0.0.0.0/0' to allow access to the public internet.",
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
										Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
									},
									"inbound_blocked_cidr": schema.SetAttribute{
										Description: "The list of ipv4/ipv6 addresses or cidr blocks that are NOT allowed to access this workload. Addresses in the allow list will only be allowed if they do not exist in this list.",
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
										Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
									},
									"outbound_allow_hostname": schema.SetAttribute{
										Description: "The list of public hostnames that this workload is allowed to reach. No outbound access is allowed by default. A wildcard `*` is allowed on the prefix of the hostname only, ex: `*.amazonaws.com`. Use `outboundAllowCIDR` to allow access to all external websites.",
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
										Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
									},
									"outbound_allow_cidr": schema.SetAttribute{
										Description: "The list of ipv4/ipv6 addresses or cidr blocks that this workload is allowed reach. No outbound access is allowed by default. Specify '0.0.0.0/0' to allow outbound access to the public internet.",
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
										Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
									},
									"outbound_blocked_cidr": schema.SetAttribute{
										Description: "The list of ipv4/ipv6 addresses or cidr blocks that this workload is NOT allowed to reach. Addresses in the allow list will only be allowed if they do not exist in this list.",
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
										Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
									},
								},
								Blocks: map[string]schema.Block{
									"outbound_allow_port": schema.ListNestedBlock{
										Description: "Allow outbound access to specific ports and protocols. When not specified, communication to address ranges in outboundAllowCIDR is allowed on all ports and communication to names in outboundAllowHostname is allowed on ports 80/443.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"protocol": schema.StringAttribute{
													Description: "Either `http`, `https` or `tcp`.",
													Required:    true,
													Validators: []validator.String{
														stringvalidator.OneOf("http", "https", "tcp"),
													},
												},
												"number": schema.Int32Attribute{
													Description: "Port number. Max: 65000",
													Required:    true,
													Validators: []validator.Int32{
														int32validator.AtMost(65000),
													},
												},
											},
										},
									},
									"http": schema.ListNestedBlock{
										Description: "Firewall options for HTTP workloads.",
										NestedObject: schema.NestedBlockObject{
											Blocks: map[string]schema.Block{
												"inbound_header_filter": schema.ListNestedBlock{
													Description: "A list of header filters for HTTP workloads.",
													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"key": schema.StringAttribute{
																Description: "The header to match for.",
																Required:    true,
																Validators: []validator.String{
																	stringvalidator.LengthAtMost(128),
																},
															},
															"allowed_values": schema.SetAttribute{
																Description: "A list of regular expressions to match for allowed header values. Headers that do not match ANY of these values will be filtered and will not reach the workload.",
																ElementType: types.StringType,
																Optional:    true,
																Validators: []validator.Set{
																	setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("blocked_values")),
																},
															},
															"blocked_values": schema.SetAttribute{
																Description: "A list of regular expressions to match for blocked header values. Headers that match ANY of these values will be filtered and will not reach the workload.",
																ElementType: types.StringType,
																Optional:    true,
																Validators: []validator.Set{
																	setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("allowed_values")),
																},
															},
														},
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
						"internal": schema.ListNestedBlock{
							Description: "The internal firewall is used to control access between workloads.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"inbound_allow_type": schema.StringAttribute{
										Description: "Used to control the internal firewall configuration and mutual tls. Allowed Values: \"none\", \"same-gvc\", \"same-org\", \"workload-list\".",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("none"),
										Validators: []validator.String{
											stringvalidator.OneOf("none", "same-gvc", "same-org", "workload-list"),
										},
									},
									"inbound_allow_workload": schema.SetAttribute{
										Description: "A list of specific workloads which are allowed to access this workload internally. This list is only used if the 'inboundAllowType' is set to 'workload-list'.",
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
										Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
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
			"options":       wr.OptionsSchema("Configurable settings or parameters that allow fine-tuning and customization of the behavior, performance, and characteristics of the workload."),
			"local_options": wr.LocalOptionsSchema(),
			"job": schema.ListNestedBlock{
				MarkdownDescription: "[Cron Job Reference Page](https://docs.controlplane.com/reference/workload#cron).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"schedule": schema.StringAttribute{
							Description: "A standard cron [schedule expression](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#schedule-syntax) used to determine when your job should execute.",
							Required:    true,
						},
						"concurrency_policy": schema.StringAttribute{
							Description: "Either 'Forbid' or 'Replace'. This determines what Control Plane will do when the schedule requires a job to start, while a prior instance of the job is still running. Enum: [ Forbid, Replace ] Default: `Forbid`.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("Forbid"),
							Validators: []validator.String{
								stringvalidator.OneOf("Forbid", "Replace"),
							},
						},
						"history_limit": schema.Int32Attribute{
							Description: "The maximum number of completed job instances to display. This should be an integer between 1 and 10. Default: `5`.",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(5),
							Validators: []validator.Int32{
								int32validator.Between(1, 10),
							},
						},
						"restart_policy": schema.StringAttribute{
							Description: "Either 'OnFailure' or 'Never'. This determines what Control Plane will do when a job instance fails. Enum: [ OnFailure, Never ] Default: `Never`.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("Never"),
							Validators: []validator.String{
								stringvalidator.OneOf("OnFailure", "Never"),
							},
						},
						"active_deadline_seconds": schema.Int32Attribute{
							Description: "The maximum number of seconds Control Plane will wait for the job to complete. If a job does not succeed or fail in the allotted time, Control Plane will stop the job, moving it into the Removed status.",
							Optional:    true,
							Validators: []validator.Int32{
								int32validator.Between(1, 86400),
							},
						},
					},
				},
			},
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
			"rollout_options": schema.ListNestedBlock{
				Description: "Defines the parameters for updating applications and services, including settings for minimum readiness, unavailable replicas, surge replicas, and scaling policies.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"min_ready_seconds": schema.Int32Attribute{
							Description: "The minimum number of seconds a container must run without crashing to be considered available.",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(0),
						},
						"max_unavailable_replicas": schema.StringAttribute{
							Description: "The number of replicas that can be unavailable during the update process.",
							Optional:    true,
						},
						"max_surge_replicas": schema.StringAttribute{
							Description: "The number of replicas that can be created above the desired amount of replicas during an update.",
							Optional:    true,
						},
						"scaling_policy": schema.StringAttribute{
							Description: "The strategies used to update applications and services deployed. Valid values: `OrderedReady` (Updates workloads in a rolling fashion, taking down old ones and bringing up new ones incrementally, ensuring that the service remains available during the update.), `Parallel` (Causes all pods affected by a scaling operation to be created or destroyed simultaneously. This does not affect update operations.). Default: `OrderedReady`.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("OrderedReady"),
							Validators: []validator.String{
								stringvalidator.OneOf("OrderedReady", "Parallel"),
							},
						},
						"termination_grace_period_seconds": schema.Int32Attribute{
							Description: "The amount of time in seconds a workload has to gracefully terminate before forcefully terminating it. This includes the time it takes for the preStop hook to run.",
							Optional:    true,
							Validators: []validator.Int32{
								int32validator.Between(0, 900),
							},
						},
					},
				},
				Validators: []validator.List{},
			},
			"security_options": schema.ListNestedBlock{
				Description: "Allows for the configuration of the `file system group id` and `geo location`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"file_system_group_id": schema.Int32Attribute{
							Description: "The group id assigned to any mounted volume.",
							Optional:    true,
							Validators: []validator.Int32{
								int32validator.Between(1, 65534),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"load_balancer": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"replica_direct": schema.BoolAttribute{
							Description: "When enabled, individual replicas of the workload can be reached directly using the subdomain prefix replica-<index>. For example, replica-0.my-workload.my-gvc.cpln.local or replica-0.my-workload-<gvc-alias>.cpln.app - Can only be used with stateful workloads.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
					Blocks: map[string]schema.Block{
						"direct": schema.ListNestedBlock{
							Description: "Direct load balancers are created in each location that a workload is running in and are configured for the standard endpoints of the workload. Customers are responsible for configuring the workload with certificates if TLS is required.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "When disabled, this load balancer will be stopped.",
										Required:    true,
									},
									"ipset": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"port": schema.ListNestedBlock{
										Description: "List of ports that will be exposed by this load balancer.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"external_port": schema.Int32Attribute{
													Description: "The port that is available publicly.",
													Required:    true,
													Validators: []validator.Int32{
														int32validator.Between(22, 32768),
													},
												},
												"protocol": schema.StringAttribute{
													Description: "The protocol that is exposed publicly.",
													Required:    true,
													Validators: []validator.String{
														stringvalidator.OneOf("TCP", "UDP"),
													},
												},
												"scheme": schema.StringAttribute{
													Description: "Overrides the default `https` url scheme that will be used for links in the UI and status.",
													Optional:    true,
													Validators: []validator.String{
														stringvalidator.OneOf("http", "tcp", "https", "ws", "wss"),
													},
												},
												"container_port": schema.Int32Attribute{
													Description: "The port on the container tha will receive this traffic.",
													Optional:    true,
													Validators:  wr.GetPortValidators(),
												},
											},
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"geo_location": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "When enabled, geo location headers will be included on inbound http requests. Existing headers will be replaced.",
										Optional:    true,
										Computed:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
								Blocks: map[string]schema.Block{
									"headers": schema.ListNestedBlock{
										Description: "",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"asn": schema.StringAttribute{
													Description: "The geo asn header.",
													Optional:    true,
													Validators: []validator.String{
														stringvalidator.LengthAtMost(128),
													},
												},
												"city": schema.StringAttribute{
													Description: "The geo city header.",
													Optional:    true,
													Validators: []validator.String{
														stringvalidator.LengthAtMost(128),
													},
												},
												"country": schema.StringAttribute{
													Description: "The geo country header.",
													Optional:    true,
													Validators: []validator.String{
														stringvalidator.LengthAtMost(128),
													},
												},
												"region": schema.StringAttribute{
													Description: "The geo region header.",
													Optional:    true,
													Validators: []validator.String{
														stringvalidator.LengthAtMost(128),
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
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"request_retry_policy": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"attempts": schema.Int32Attribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(2),
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
							},
						},
						"retry_on": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
								types.StringValue("connect-failure"),
								types.StringValue("refused-stream"),
								types.StringValue("unavailable"),
								types.StringValue("cancelled"),
								types.StringValue("resource-exhausted"),
								types.StringValue("retriable-status-codes"),
							})),
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

// ModifyPlan modifies the plan for the resource.
func (wr *WorkloadResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If this is a destroy plan, leave everything null and return immediately
	if req.Plan.Raw.IsNull() {
		return
	}

	// Declare variable to store desired resource plan
	var plan WorkloadResourceModel

	// Populate plan variable from request and capture diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	// Abort if any diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Modify autoscaling in options if specified
	if len(plan.Options) != 0 && len(plan.Options[0].Autoscaling) != 0 {
		wr.ModifyAutoscaling(&plan.Options[0].Autoscaling[0])
	}

	// Modify autoscaling in local options if specified
	if len(plan.LocalOptions) != 0 {
		// Iterate over local options and modify autoscaling
		for _, l := range plan.LocalOptions {
			// Skip if autoscaling is not specified
			if len(l.Autoscaling) == 0 {
				continue
			}

			// Modify autoscaling for this local options block
			wr.ModifyAutoscaling(&l.Autoscaling[0])
		}
	}

	// Modify containers if specified
	if len(plan.Containers) != 0 {
		// Iterate over containers and modify each
		for _, c := range plan.Containers {
			wr.ModifyContainers(&c)
		}
	}

	// Persist new plan into Terraform
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// ValidateConfig validates the configuration of the resource.
func (wr *WorkloadResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Declare variable to store desired resource plan
	var plan WorkloadResourceModel

	// Populate plan variable from config and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	// Halt further processing if plan retrieval failed
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize the validator
	validator := WorkloadResourceValidator{Diags: &resp.Diagnostics, Plan: plan}

	// Call the validate method
	validator.Validate()
}

// Create creates the resource.
func (wr *WorkloadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, wr.Operations)
}

// Read fetches the current state of the resource.
func (wr *WorkloadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, wr.Operations)
}

// Update modifies the resource.
func (wr *WorkloadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, wr.Operations)
}

// Delete removes the resource.
func (wr *WorkloadResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, wr.Operations)
}

/*** Schemas ***/

// HealthCheckSchema returns a nested block list schema for configuring workload health checks.
func (wr *WorkloadResource) HealthCheckSchema(attributeName string, description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"initial_delay_seconds": schema.Int32Attribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     int32default.StaticInt32(10),
					Validators: []validator.Int32{
						int32validator.Between(0, 600),
					},
				},
				"period_seconds": schema.Int32Attribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     int32default.StaticInt32(10),
					Validators: []validator.Int32{
						int32validator.Between(1, 600),
					},
				},
				"timeout_seconds": schema.Int32Attribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     int32default.StaticInt32(1),
					Validators: []validator.Int32{
						int32validator.Between(1, 600),
					},
				},
				"success_threshold": schema.Int32Attribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     int32default.StaticInt32(1),
					Validators: []validator.Int32{
						int32validator.Between(1, 20),
					},
				},
				"failure_threshold": schema.Int32Attribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     int32default.StaticInt32(3),
					Validators: []validator.Int32{
						int32validator.Between(1, 20),
					},
				},
			},
			Blocks: map[string]schema.Block{
				"exec": wr.ExecSchema(""),
				"grpc": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"port": schema.Int32Attribute{
								Description: "",
								Optional:    true,
								Validators:  wr.GetPortValidators(),
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
				},
				"tcp_socket": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"port": schema.Int32Attribute{
								Description: "",
								Optional:    true,
								Validators:  wr.GetPortValidators(),
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
				},
				"http_get": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("/"),
								Validators: []validator.String{
									stringvalidator.LengthAtMost(256),
								},
							},
							"port": schema.Int32Attribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								Validators:  wr.GetPortValidators(),
							},
							"http_headers": schema.MapAttribute{
								Description: "",
								ElementType: types.StringType,
								Optional:    true,
							},
							"scheme": schema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("HTTP"),
								Validators: []validator.String{
									stringvalidator.OneOf("HTTP", "HTTPS"),
								},
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
				},
			},
			Validators: []validator.Object{
				// objectvalidator.ConflictsWith(
				// 	path.MatchRoot("container").AtAnyListIndex().AtName(attributeName).AtAnyListIndex().AtName("exec"),
				// 	path.MatchRoot("container").AtAnyListIndex().AtName(attributeName).AtAnyListIndex().AtName("grpc"),
				// 	path.MatchRoot("container").AtAnyListIndex().AtName(attributeName).AtAnyListIndex().AtName("tcp_socket"),
				// 	path.MatchRoot("container").AtAnyListIndex().AtName(attributeName).AtAnyListIndex().AtName("http_get"),
				// ),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// ExecSchema returns a nested block list schema for configuring exec-based probes.
func (wr *WorkloadResource) ExecSchema(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"command": schema.SetAttribute{
					Description: description,
					ElementType: types.StringType,
					Required:    true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// LifecycleSpecSchema returns a nested block list schema for workload lifecycle specifications.
func (wr *WorkloadResource) LifecycleSpecSchema(commandDescription string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"exec": wr.ExecSchema(commandDescription),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// OptionsSchema returns a nested block list schema for workload options such as AI capacity, debug mode, and auto-scaling.
func (wr *WorkloadResource) OptionsSchema(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"timeout_seconds": schema.Int32Attribute{
					Description: "Timeout in seconds. Default: `5`.",
					Optional:    true,
					Computed:    true,
					Default:     int32default.StaticInt32(5),
					Validators: []validator.Int32{
						int32validator.Between(1, 3600),
					},
				},
				"capacity_ai": schema.BoolAttribute{
					Description: "Capacity AI. Default: `true`.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
				},
				"debug": schema.BoolAttribute{
					Description: "Debug mode. Default: `false`.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"suspend": schema.BoolAttribute{
					Description: "Workload suspend. Default: `false`.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
			},
			Blocks: map[string]schema.Block{
				"autoscaling": schema.ListNestedBlock{
					Description: "Auto-scaling adjusts horizontal scaling based on a set strategy, target value, and possibly a metric percentile.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"metric": schema.StringAttribute{
								Description: "Valid values: `disabled`, `concurrency`, `cpu`, `memory`, `latency`, or `rps`.",
								Optional:    true,
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("concurrency", "cpu", "memory", "rps", "latency", "disabled"),
								},
							},
							"metric_percentile": schema.StringAttribute{
								Description: "For metrics represented as a distribution (e.g. latency) a percentile within the distribution must be chosen as the target.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("p50", "p75", "p99"),
								},
							},
							"target": schema.Int32Attribute{
								Description: "Control Plane will scale the number of replicas for this deployment up/down in order to be as close as possible to the target metric across all replicas of a deployment. Min: `1`. Max: `20000`. Default: `95`.",
								Optional:    true,
								Computed:    true,
								Validators: []validator.Int32{
									int32validator.Between(1, 20000),
								},
							},
							"max_scale": schema.Int32Attribute{
								Description: "The maximum allowed number of replicas. Min: `0`. Default `5`.",
								Optional:    true,
								Computed:    true,
								Default:     int32default.StaticInt32(5),
								Validators: []validator.Int32{
									int32validator.AtLeast(0),
								},
							},
							"min_scale": schema.Int32Attribute{
								Description: "The minimum allowed number of replicas. Control Plane can scale the workload down to 0 when there is no traffic and scale up immediately to fulfill new requests. Min: `0`. Max: `max_scale`. Default `1`.",
								Optional:    true,
								Computed:    true,
								Default:     int32default.StaticInt32(1),
								Validators: []validator.Int32{
									int32validator.AtLeast(0),
								},
							},
							"max_concurrency": schema.Int32Attribute{
								Description: "A hard maximum for the number of concurrent requests allowed to a replica. If no replicas are available to fulfill the request then it will be queued until a replica with capacity is available and delivered as soon as one is available again. Capacity can be available from requests completing or when a new replica is available from scale out.Min: `0`. Max: `1000`. Default `0`.",
								Optional:    true,
								Computed:    true,
								Default:     int32default.StaticInt32(0),
								Validators: []validator.Int32{
									int32validator.Between(0, 30000),
								},
							},
							"scale_to_zero_delay": schema.Int32Attribute{
								Description: "The amount of time (in seconds) with no requests received before a workload is scaled to 0. Min: `30`. Max: `3600`. Default: `300`.",
								Optional:    true,
								Computed:    true,
								Default:     int32default.StaticInt32(300),
								Validators: []validator.Int32{
									int32validator.Between(30, 3600),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"multi": schema.ListNestedBlock{
								Description: "",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"metric": schema.StringAttribute{
											Description: "Valid values: `cpu` or `memory`.",
											Optional:    true,
											Validators: []validator.String{
												stringvalidator.OneOf("cpu", "memory", "rps"),
											},
										},
										"target": schema.Int32Attribute{
											Description: "Control Plane will scale the number of replicas for this deployment up/down in order to be as close as possible to the target metric across all replicas of a deployment. Min: `1`. Max: `20000`.",
											Optional:    true,
											Validators: []validator.Int32{
												int32validator.Between(1, 20000),
											},
										},
									},
								},
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
				},
				"multi_zone": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
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
	}
}

// LocalOptionsSchema returns a nested block list schema for overriding default options per specific Control Plane location.
func (wr *WorkloadResource) LocalOptionsSchema() schema.ListNestedBlock {
	// Build base options schema with override description
	options := wr.OptionsSchema("Override defaultOptions for the workload in specific Control Plane Locations.")

	// Merge in location attribute to override options per location
	options.NestedObject.Attributes = MergeAttributes(options.NestedObject.Attributes, map[string]schema.Attribute{
		"location": schema.StringAttribute{
			Description: "Valid only for `local_options`. Override options for a specific location.",
			Required:    true,
		},
	})

	// Remove any list validators so multiple blocks are allowed
	options.Validators = []validator.List{}

	// Return the configured local options schema
	return options
}

/*** Shared Modifiers ***/

// ModifyAutoscaling sets default values for metric and target if autoscaling is single-target and values are not provided.
func (wr *WorkloadResource) ModifyAutoscaling(autoscaling *models.OptionsAutoscalingModel) {
	// If there are no multi, then set target and metric
	if len(autoscaling.Multi) == 0 {
		// Only modify if metric is not specified by the user
		if autoscaling.Metric.IsNull() || autoscaling.Metric.IsUnknown() {
			autoscaling.Metric = types.StringValue("concurrency")
		}

		// Only modify if target is not specified by the user
		if autoscaling.Target.IsNull() || autoscaling.Target.IsUnknown() {
			autoscaling.Target = types.Int32Value(95)
		}
	}
}

// ModifyContainers updates container health checks using the first available port if port is not explicitly set.
func (wr *WorkloadResource) ModifyContainers(container *models.ContainerModel) {
	// Declare a variable to hold the port number
	var firstPortNumber *int

	// Attempt to retrieve the port number from the port attribute
	if !container.Port.IsNull() && !container.Port.IsUnknown() {
		firstPortNumber = BuildInt(container.Port)
	}

	// If the port number is still nil, extract the port number from the very first item of container ports
	if len(container.Ports) != 0 {
		firstPortNumber = BuildInt(container.Ports[0].Number)
	}

	// Modify liveness probe if specified
	if len(container.LivenessProbe) != 0 {
		wr.ModifyHealthCheck(&container.LivenessProbe[0], firstPortNumber)
	}

	// Modify readiness probe if specified
	if len(container.ReadinessProbe) != 0 {
		wr.ModifyHealthCheck(&container.ReadinessProbe[0], firstPortNumber)
	}
}

// ModifyHealthCheck sets a default port for the HTTP health check if not explicitly defined.
func (wr *WorkloadResource) ModifyHealthCheck(healthCheck *models.ContainerHealthCheckModel, port *int) {
	// Modify the port of the httpGet health check if it hasn't been specified
	if port != nil && len(healthCheck.HttpGet) != 0 && healthCheck.HttpGet[0].Port.IsUnknown() {
		healthCheck.HttpGet[0].Port = types.Int32Value(int32(*port))
	}
}

/*** Shared Validators ***/

// GetCpuMemoryValidators returns a list of validators to ensure CPU/memory values follow proper format and size limits.
func (wr *WorkloadResource) GetCpuMemoryValidators(regexMessage string) []validator.String {
	return []validator.String{
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$`),
			regexMessage,
		),
		stringvalidator.LengthAtMost(20),
	}
}

/*** Resource Validator ***/

// WorkloadResourceValidator is used to validate the resource configuration.
type WorkloadResourceValidator struct {
	Diags *diag.Diagnostics
	Plan  WorkloadResourceModel
}

// Validate validates the workload resource plan against defined constraints.
func (wrv *WorkloadResourceValidator) Validate() {
	// Extract the workload type from the plan
	workloadType := wrv.Plan.Type.ValueString()

	// Job must be defined for cron workloads
	if workloadType == "cron" && len(wrv.Plan.Job) == 0 {
		wrv.Diags.AddAttributeError(
			path.Root("type"),
			"Missing Cron Job Configuration",
			"The 'job' block must be defined when the workload type is set to 'cron'.",
		)
		return
	}

	// Only standard and stateful workload can have rollout options
	if workloadType != "standard" && workloadType != "stateful" && len(wrv.Plan.RolloutOptions) > 0 {
		wrv.Diags.AddAttributeError(
			path.Root("type"),
			"Invalid Rollout Options",
			"The 'rollout_options' block can only be defined for 'standard' or 'stateful' workload types.",
		)
	}

	// Initialize flags for GPU, Min CPU, and Min Memory usage
	isUsingGpu := false
	isUsingMinCpu := false
	isUsingMinMemory := false

	// Iterate over each container and validate it
	for _, container := range wrv.Plan.Containers {
		// Probes are not allowed for cron workloads
		if workloadType == "cron" && (len(container.ReadinessProbe) > 0 || len(container.LivenessProbe) > 0) {
			wrv.Diags.AddAttributeError(
				path.Root("containers"),
				"Invalid Probes for Cron Workload",
				"Health checks are not supported for cron workloads. Remove 'readiness_probe' and 'liveness_probe' blocks.",
			)
			return
		}

		// Validate CPU and Memory values for workloads with GPU
		if len(container.GpuNvidia) > 0 {
			isUsingGpu = true
			cpuAmount, cpuUnit := ParseValueAndUnit(container.Cpu.ValueString())
			memoryAmount, memoryUnit := ParseValueAndUnit(container.Memory.ValueString())

			// Check if CPU and Memory meet the minimum requirements for GPU workloads
			if (cpuUnit == "" && cpuAmount < 2) ||
				(cpuUnit == "m" && cpuAmount < 2000) ||
				(memoryUnit == "Gi" && memoryAmount < 7) ||
				(memoryUnit == "Mi" && memoryAmount < 7000) {
				wrv.Diags.AddAttributeError(
					path.Root("containers"),
					"Insufficient CPU or Memory",
					"The GPU requires this container to have at least 2 CPU cores (or 2000m) and 7 Gi (or 7000 Mi) of RAM.",
				)
			}
		}

		// Check if Min CPU is being used
		if !container.MinCpu.IsNull() && !container.MinCpu.IsUnknown() {
			isUsingMinCpu = true
		}

		// Check if Min Memory is being used
		if !container.MinMemory.IsNull() && !container.MinMemory.IsUnknown() {
			isUsingMinMemory = true
		}
	}

	// Validate Default Options
	wrv.validateOptions(path.Root("options").AtListIndex(0), workloadType, wrv.Plan.Options, isUsingGpu, isUsingMinCpu, isUsingMinMemory)

	// Validate local options
	for i, localOption := range wrv.Plan.LocalOptions {
		wrv.validateOptions(path.Root("local_options").AtListIndex(i), workloadType, []models.OptionsModel{localOption.OptionsModel}, isUsingGpu, isUsingMinCpu, isUsingMinMemory)
	}
}

// ValidateOptions validates the options for different workload types.
func (wrv *WorkloadResourceValidator) validateOptions(
	basePath path.Path,
	workloadType string,
	options []models.OptionsModel,
	isUsingGpu bool,
	isUsingMinCpu bool,
	isUsingMinMemory bool,
) {
	// Return early if no options provided
	if len(options) == 0 {
		return
	}

	// Select the first options model
	opt := options[0]

	// Return early if autoscaling configuration is absent
	if len(opt.Autoscaling) == 0 {
		wrv.Diags.AddAttributeError(
			basePath.AtName("autoscaling"),
			"Autoscaling Block Is Required",
			"Add an empty autoscaling block under your options/local_options block. Example: 'autoscaling {}'",
		)
		return
	}

	// Retrieve the first autoscaling configuration and its path
	asc := opt.Autoscaling[0]
	ascPath := basePath.AtName("autoscaling")

	// Apply cron-specific validation rules
	if workloadType == "cron" {
		// Report error if min CPU is used for cron workloads
		if isUsingMinCpu {
			wrv.Diags.AddAttributeError(
				basePath.AtName("min_cpu"),
				"Min CPU not allowed for cron workloads",
				"'min_cpu' is not allowed for workload of type 'cron'",
			)
		}

		// Report error if min memory is used for cron workloads
		if isUsingMinMemory {
			wrv.Diags.AddAttributeError(
				basePath.AtName("min_memory"),
				"Min Memory not allowed for cron workloads",
				"'min_memory' is not allowed for workload of type 'cron'",
			)
		}
	} else {
		// Apply non-cron validation rules
		if len(asc.Multi) > 0 {
			// Report error if metric is set alongside multiple metrics strategy
			if !asc.Metric.IsNull() && !asc.Metric.IsUnknown() {
				wrv.Diags.AddAttributeError(
					ascPath.AtName("metric"),
					"Metric conflicts with Multi",
					"'metric' must not exist simultaneously with 'multi'",
				)
			}

			// Report error if target is set alongside multiple metrics strategy
			if !asc.Target.IsNull() && !asc.Target.IsUnknown() {
				wrv.Diags.AddAttributeError(
					ascPath.AtName("target"),
					"Target conflicts with Multi",
					"'target' must not exist simultaneously with 'multi'",
				)
			}
		}

		// Validate Capacity AI settings
		if opt.CapacityAI.ValueBool() {
			// Report error if GPU is used with Capacity AI enabled
			if isUsingGpu {
				wrv.Diags.AddAttributeError(
					basePath.AtName("capacity_ai"),
					"Invalid Capacity AI for Workload With GPU",
					"Capacity AI cannot be enabled for workloads with GPU. Please disable it",
				)
			}
		}
	}
}

/*** Resource Operator ***/

// WorkloadResourceOperator is the operator for managing the state.
type WorkloadResourceOperator struct {
	EntityOperator[WorkloadResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (wro *WorkloadResourceOperator) NewAPIRequest(isUpdate bool) client.Workload {
	// Initialize a new request payload
	requestPayload := client.Workload{}

	// Initialize the GVC spec struct
	var spec *client.WorkloadSpec = &client.WorkloadSpec{}

	// Populate Base fields from state
	wro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Map planned state attributes to the API struct
	if isUpdate {
		requestPayload.SpecReplace = spec
	} else {
		requestPayload.Spec = spec
	}

	// Set specific attributes
	spec.Type = BuildString(wro.Plan.Type)
	spec.IdentityLink = BuildString(wro.Plan.IdentityLink)
	spec.Containers = wro.buildContainers(wro.Plan.Containers)
	spec.FirewallConfig = wro.buildFirewall(wro.Plan.Firewall)
	spec.DefaultOptions = wro.buildOptions(wro.Plan.Options)
	spec.LocalOptions = wro.buildLocalOptions(wro.Plan.LocalOptions)
	spec.Job = wro.buildJob(wro.Plan.Job)
	spec.Sidecar = wro.buildSidecar(wro.Plan.Sidecar)
	spec.SupportDynamicTags = BuildBool(wro.Plan.SupportDynamicTags)
	spec.RolloutOptions = wro.buildRolloutOptions(wro.Plan.RolloutOptions)
	spec.SecurityOptions = wro.buildSecurityOptions(wro.Plan.SecurityOptions)
	spec.LoadBalancer = wro.buildLoadBalancer(wro.Plan.LoadBalancer)
	spec.Extras = wro.buildExtras(wro.Plan.Extras)
	spec.RequestRetryPolicy = wro.buildRequestRetryPolicy(wro.Plan.RequestRetryPolicy)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (wro *WorkloadResourceOperator) MapResponseToState(apiResp *client.Workload, isCreate bool) WorkloadResourceModel {
	// Initialize empty state model
	state := WorkloadResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Set specific attributes
	state.Gvc = types.StringPointerValue(BuildString(wro.Plan.Gvc))
	state.Status = wro.flattenStatus(apiResp.Status)

	// Just in case the spec is nil
	if apiResp.Spec == nil {
		state.Type = types.StringNull()
		state.IdentityLink = types.StringNull()
		state.Containers = nil
		state.Firewall = nil
		state.Options = nil
		state.LocalOptions = nil
		state.Job = nil
		state.Sidecar = nil
		state.SupportDynamicTags = types.BoolNull()
		state.RolloutOptions = nil
		state.SecurityOptions = nil
		state.LoadBalancer = nil
		state.Extras = types.StringNull()
		state.RequestRetryPolicy = nil
	} else {
		state.Type = types.StringPointerValue(apiResp.Spec.Type)
		state.IdentityLink = types.StringPointerValue(apiResp.Spec.IdentityLink)
		state.Containers = wro.flattenContainers(apiResp.Spec.Containers)
		state.Firewall = wro.flattenFirewall(apiResp.Spec.FirewallConfig)
		state.Options = wro.flattenOptions(apiResp.Spec.DefaultOptions)
		state.LocalOptions = wro.flattenLocalOptions(apiResp.Spec.LocalOptions)
		state.Job = wro.flattenJob(apiResp.Spec.Job)
		state.Sidecar = wro.flattenSidecar(apiResp.Spec.Sidecar)
		state.SupportDynamicTags = types.BoolPointerValue(apiResp.Spec.SupportDynamicTags)
		state.RolloutOptions = wro.flattenRolloutOptions(apiResp.Spec.RolloutOptions)
		state.SecurityOptions = wro.flattenSecurityOptions(apiResp.Spec.SecurityOptions)
		state.LoadBalancer = wro.flattenLoadBalancer(wro.Plan.LoadBalancer, apiResp.Spec.LoadBalancer)
		state.Extras = wro.flattenExtras(apiResp.Spec.Extras)
		state.RequestRetryPolicy = wro.flattenRequestRetryPolicy(apiResp.Spec.RequestRetryPolicy)
	}

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (wro *WorkloadResourceOperator) InvokeCreate(req client.Workload) (*client.Workload, int, error) {
	return wro.Client.CreateWorkload(req, wro.Plan.Gvc.ValueString())
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (wro *WorkloadResourceOperator) InvokeRead(name string) (*client.Workload, int, error) {
	return wro.Client.GetWorkload(name, wro.Plan.Gvc.ValueString())
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (wro *WorkloadResourceOperator) InvokeUpdate(req client.Workload) (*client.Workload, int, error) {
	return wro.Client.UpdateWorkload(req, wro.Plan.Gvc.ValueString())
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (wro *WorkloadResourceOperator) InvokeDelete(name string) error {
	return wro.Client.DeleteWorkload(name, wro.Plan.Gvc.ValueString())
}

// Builders //

// buildContainers constructs a []client.WorkloadContainer from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainers(state []models.ContainerModel) *[]client.WorkloadContainer {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadContainer{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Initialize a slice for container port models copying existing ports
		blockPorts := make([]models.ContainerPortModel, len(block.Ports))

		// Copy existing ports into the new slice
		copy(blockPorts, block.Ports)

		// Append legacy port attribute if it is specified
		if !block.Port.IsNull() && !block.Port.IsUnknown() {
			blockPorts = append(blockPorts, models.ContainerPortModel{
				Protocol: types.StringValue("http"),
				Number:   block.Port,
			})
		}

		// Construct the item
		item := client.WorkloadContainer{
			Name:             BuildString(block.Name),
			Image:            BuildString(block.Image),
			WorkingDirectory: BuildString(block.WorkingDirectory),
			Metrics:          wro.buildContainerMetrics(block.Metrics),
			Ports:            wro.buildContainerPort(blockPorts),
			Memory:           BuildString(block.Memory),
			ReadinessProbe:   wro.buildHealthCheck(block.ReadinessProbe),
			LivenessProbe:    wro.buildHealthCheck(block.LivenessProbe),
			CPU:              BuildString(block.Cpu),
			MinCPU:           BuildString(block.MinCpu),
			MinMemory:        BuildString(block.MinMemory),
			Env:              wro.buildNameValue(block.Env),
			GPU:              wro.buildContainerGpu(block.GpuNvidia, block.GpuCustom),
			InheritEnv:       BuildBool(block.InheritEnv),
			Command:          BuildString(block.Command),
			Args:             wro.BuildSetString(block.Args),
			LifeCycle:        wro.buildContainerLifecycle(block.Lifecycle),
			Volumes:          wro.buildContainerVolume(block.Volumes),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildContainerMetrics constructs a WorkloadContainerMetrics from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerMetrics(state []models.ContainerMetricsModel) *client.WorkloadContainerMetrics {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadContainerMetrics{
		Port: BuildInt(block.Port),
		Path: BuildString(block.Path),
	}
}

// buildContainerPort constructs a []client.WorkloadContainerPort from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerPort(state []models.ContainerPortModel) *[]client.WorkloadContainerPort {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadContainerPort{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.WorkloadContainerPort{
			Protocol: BuildString(block.Protocol),
			Number:   BuildInt(block.Number),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildHealthCheck constructs a WorkloadHealthCheck from the given Terraform state.
func (wro *WorkloadResourceOperator) buildHealthCheck(state []models.ContainerHealthCheckModel) *client.WorkloadHealthCheck {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadHealthCheck{
		Exec:                wro.buildExec(block.Exec),
		GRPC:                wro.buildHealthCheckGrpc(block.Grpc),
		TCPSocket:           wro.buildHealthCheckTcpSocket(block.TcpSocket),
		HTTPGet:             wro.buildHealthCheckHttpGet(block.HttpGet),
		InitialDelaySeconds: BuildInt(block.InitialDelaySeconds),
		PeriodSeconds:       BuildInt(block.PeriodSeconds),
		TimeoutSeconds:      BuildInt(block.TimeoutSeconds),
		SuccessThreshold:    BuildInt(block.SuccessThreshold),
		FailureThreshold:    BuildInt(block.FailureThreshold),
	}
}

// buildExec constructs a WorkloadExec from the given Terraform state.
func (wro *WorkloadResourceOperator) buildExec(state []models.ContainerExecModel) *client.WorkloadExec {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadExec{
		Command: wro.BuildSetString(block.Command),
	}
}

// buildHealthCheckGrpc constructs a WorkloadHealthCheckGrpc from the given Terraform state.
func (wro *WorkloadResourceOperator) buildHealthCheckGrpc(state []models.ContainerHealthCheckGrpcModel) *client.WorkloadHealthCheckGrpc {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadHealthCheckGrpc{
		Port: BuildInt(block.Port),
	}
}

// buildHealthCheckTcpSocket constructs a WorkloadHealthCheckTcpSocket from the given Terraform state.
func (wro *WorkloadResourceOperator) buildHealthCheckTcpSocket(state []models.ContainerHealthCheckTcpSocketModel) *client.WorkloadHealthCheckTcpSocket {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadHealthCheckTcpSocket{
		Port: BuildInt(block.Port),
	}
}

// buildHealthCheckHttpGet constructs a WorkloadHealthCheckHttpGet from the given Terraform state.
func (wro *WorkloadResourceOperator) buildHealthCheckHttpGet(state []models.ContainerHealthCheckHttpGetModel) *client.WorkloadHealthCheckHttpGet {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadHealthCheckHttpGet{
		Path:        BuildString(block.Path),
		Port:        BuildInt(block.Port),
		HttpHeaders: wro.buildNameValue(block.HttpHeaders),
		Scheme:      BuildString(block.Scheme),
	}
}

// buildNameValue constructs a *[]client.WorkloadContainerNameValue from the given Terraform state.
func (wro *WorkloadResourceOperator) buildNameValue(state types.Map) *[]client.WorkloadContainerNameValue {
	// Convert Terraform HTTP headers to a map
	headersMap := wro.BuildMapString(state)

	// Return nil if the map is nil
	if headersMap == nil {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadContainerNameValue{}

	// Iterate over the map to convert it to a slice of NameValue structs
	for name, value := range *headersMap {
		// Append each header as a NameValue struct
		output = append(output, client.WorkloadContainerNameValue{
			Name:  &name,
			Value: StringPointerFromInterface(value),
		})
	}

	// Return a pointer to the output
	return &output
}

// buildContainerGpu constructs a WorkloadContainerGpu from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerGpu(nvidiaState []models.ContainerGpuNvidiaModel, customState []models.ContainerGpuCustomModel) *client.WorkloadContainerGpu {
	// Build the GPU models from the provided states
	nvidia := wro.buildContainerGpuNvidia(nvidiaState)
	custom := wro.buildContainerGpuCustom(customState)

	// Return nil if both nvidia and custom are nil
	if nvidia == nil && custom == nil {
		return nil
	}

	// Construct and return the output
	return &client.WorkloadContainerGpu{
		Nvidia: nvidia,
		Custom: custom,
	}
}

// buildContainerGpuNvidia constructs a WorkloadContainerGpuNvidia from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerGpuNvidia(state []models.ContainerGpuNvidiaModel) *client.WorkloadContainerGpuNvidia {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadContainerGpuNvidia{
		Model:    BuildString(block.Model),
		Quantity: BuildInt(block.Quantity),
	}
}

// buildContainerGpuCustom constructs a WorkloadContainerGpuCustom from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerGpuCustom(state []models.ContainerGpuCustomModel) *client.WorkloadContainerGpuCustom {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadContainerGpuCustom{
		Resource:     BuildString(block.Resource),
		RuntimeClass: BuildString(block.RuntimeClass),
		Quantity:     BuildInt(block.Quantity),
	}
}

// buildContainerLifecycle constructs a WorkloadLifeCycle from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerLifecycle(state []models.ContainerLifecycleModel) *client.WorkloadLifeCycle {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadLifeCycle{
		PostStart: wro.buildContainerLifecycleSpec(block.PostStart),
		PreStop:   wro.buildContainerLifecycleSpec(block.PreStop),
	}
}

// buildContainerLifecycleSpec constructs a WorkloadLifeCycleSpec from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerLifecycleSpec(state []models.ContainerLifecycleSpecModel) *client.WorkloadLifeCycleSpec {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadLifeCycleSpec{
		Exec: wro.buildExec(block.Exec),
	}
}

// buildContainerVolume constructs a []client.WorkloadContainerVolume from the given Terraform state.
func (wro *WorkloadResourceOperator) buildContainerVolume(state []models.ContainerVolumeModel) *[]client.WorkloadContainerVolume {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadContainerVolume{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.WorkloadContainerVolume{
			Uri:            BuildString(block.Uri),
			RecoveryPolicy: BuildString(block.RecoveryPolicy),
			Path:           BuildString(block.Path),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildFirewall constructs a WorkloadFirewall from the given Terraform state.
func (wro *WorkloadResourceOperator) buildFirewall(state []models.FirewallModel) *client.WorkloadFirewall {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadFirewall{
		External: wro.buildFirewallExternal(block.External),
		Internal: wro.buildFirewallInternal(block.Internal),
	}
}

// buildFirewallExternal constructs a WorkloadFirewallExternal from the given Terraform state.
func (wro *WorkloadResourceOperator) buildFirewallExternal(state []models.FirewallExternalModel) *client.WorkloadFirewallExternal {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadFirewallExternal{
		InboundAllowCidr:      wro.BuildSetString(block.InboundAllowCidr),
		InboundBlockedCidr:    wro.BuildSetString(block.InboundBlockedCidr),
		OutboundAllowHostname: wro.BuildSetString(block.OutboundAllowHostname),
		OutboundAllowPort:     wro.buildFirewallExternalOutboundAllowPort(block.OutboundAllowPort),
		OutboundAllowCidr:     wro.BuildSetString(block.OutboundAllowCidr),
		OutboundBlockedCidr:   wro.BuildSetString(block.OutboundBlockedCidr),
		Http:                  wro.buildFirewallExternalHttp(block.Http),
	}
}

// buildFirewallExternalOutboundAllowPort constructs a []client.WorkloadFirewallOutboundAllowPort from the given Terraform state.
func (wro *WorkloadResourceOperator) buildFirewallExternalOutboundAllowPort(state []models.FirewallExternalOutboundAllowPortModel) *[]client.WorkloadFirewallOutboundAllowPort {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadFirewallOutboundAllowPort{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.WorkloadFirewallOutboundAllowPort{
			Protocol: BuildString(block.Protocol),
			Number:   BuildInt(block.Number),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildFirewallExternalHttp constructs a WorkloadFirewallExternalHttp from the given Terraform state.
func (wro *WorkloadResourceOperator) buildFirewallExternalHttp(state []models.FirewallExternalHttpModel) *client.WorkloadFirewallExternalHttp {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadFirewallExternalHttp{
		InboundHeaderFilter: wro.buildFirewallExternalHttpHeaderFilter(block.InboundHeaderFilter),
	}
}

// buildFirewallExternalHttpHeaderFilter constructs a []client.WorkloadFirewallExternalHttpHeaderFilter from the given Terraform state.
func (wro *WorkloadResourceOperator) buildFirewallExternalHttpHeaderFilter(state []models.FirewallExternalHttpHeaderFilterModel) *[]client.WorkloadFirewallExternalHttpHeaderFilter {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadFirewallExternalHttpHeaderFilter{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.WorkloadFirewallExternalHttpHeaderFilter{
			Key:           BuildString(block.Key),
			AllowedValues: wro.BuildSetString(block.AllowedValues),
			BlockedValues: wro.BuildSetString(block.BlockedValues),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildFirewallInternal constructs a WorkloadFirewallInternal from the given Terraform state.
func (wro *WorkloadResourceOperator) buildFirewallInternal(state []models.FirewallInternalModel) *client.WorkloadFirewallInternal {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadFirewallInternal{
		InboundAllowType:     BuildString(block.InboundAllowType),
		InboundAllowWorkload: wro.BuildSetString(block.InboundAllowWorkload),
	}
}

// buildOptions constructs a WorkloadOptions from the given Terraform state.
func (wro *WorkloadResourceOperator) buildOptions(state []models.OptionsModel) *client.WorkloadOptions {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadOptions{
		AutoScaling:    wro.buildOptionsAutoscaling(block.Autoscaling),
		TimeoutSeconds: BuildInt(block.TimeoutSeconds),
		CapacityAI:     BuildBool(block.CapacityAI),
		Debug:          BuildBool(block.Debug),
		Suspend:        BuildBool(block.Suspend),
		MultiZone:      wro.buildOptionsMultiZone(block.MultiZone),
	}
}

// buildOptionsAutoscaling constructs a WorkloadOptionsAutoscaling from the given Terraform state.
func (wro *WorkloadResourceOperator) buildOptionsAutoscaling(state []models.OptionsAutoscalingModel) *client.WorkloadOptionsAutoscaling {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadOptionsAutoscaling{
		Metric:           BuildString(block.Metric),
		Multi:            wro.buildOptionsAutoscalingMulti(block.Multi),
		MetricPercentile: BuildString(block.MetricPercentile),
		Target:           BuildInt(block.Target),
		MaxScale:         BuildInt(block.MaxScale),
		MinScale:         BuildInt(block.MinScale),
		MaxConcurrency:   BuildInt(block.MaxConcurrency),
		ScaleToZeroDelay: BuildInt(block.ScaleToZeroDelay),
	}
}

// buildOptionsAutoscalingMulti constructs a []client.WorkloadOptionsAutoscalingMulti from the given Terraform state.
func (wro *WorkloadResourceOperator) buildOptionsAutoscalingMulti(state []models.OptionsAutoscalingMultiModel) *[]client.WorkloadOptionsAutoscalingMulti {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadOptionsAutoscalingMulti{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.WorkloadOptionsAutoscalingMulti{
			Metric: BuildString(block.Metric),
			Target: BuildInt(block.Target),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildOptionsMultiZone constructs a WorkloadOptionsMultiZone from the given Terraform state.
func (wro *WorkloadResourceOperator) buildOptionsMultiZone(state []models.OptionsMultiZoneModel) *client.WorkloadOptionsMultiZone {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadOptionsMultiZone{
		Enabled: BuildBool(block.Enabled),
	}
}

// buildLocalOptions constructs a []client.WorkloadOptions from the given Terraform state.
func (wro *WorkloadResourceOperator) buildLocalOptions(state []models.LocalOptionsModel) *[]client.WorkloadOptions {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadOptions{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := wro.buildOptions([]models.OptionsModel{block.OptionsModel})

		// If the item is nil, skip it
		if item == nil {
			continue
		}

		// Build the location for the item
		location := BuildString(block.Location)

		// Set the location for the item
		if location != nil {
			item.Location = StringPointer(GetSelfLink(wro.Client.Org, "location", *location))
		} else {
			item.Location = nil
		}

		// Add the item to the output slice
		output = append(output, *item)
	}

	// Return a pointer to the output
	return &output
}

// buildJob constructs a WorkloadJob from the given Terraform state.
func (wro *WorkloadResourceOperator) buildJob(state []models.JobModel) *client.WorkloadJob {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadJob{
		Schedule:              BuildString(block.Schedule),
		ConcurrencyPolicy:     BuildString(block.ConcurrencyPolicy),
		HistoryLimit:          BuildInt(block.HistoryLimit),
		RestartPolicy:         BuildString(block.RestartPolicy),
		ActiveDeadlineSeconds: BuildInt(block.ActiveDeadlineSeconds),
	}
}

// buildSidecar constructs a WorkloadSidecar from the given Terraform state.
func (wro *WorkloadResourceOperator) buildSidecar(state []models.SidecarModel) *client.WorkloadSidecar {
	// Return nil if no sidecar configuration is provided
	if len(state) == 0 {
		return nil
	}

	// Select the first sidecar model
	block := state[0]

	// Build JSON string for Envoy configuration
	envoyJSON := BuildString(block.Envoy)

	// Return empty sidecar if Envoy config is missing
	if envoyJSON == nil {
		return &client.WorkloadSidecar{}
	}

	// Initialize an empty interface to hold the parsed Envoy configuration
	var envoyConfig interface{}

	// Attempt to unmarshal the JSON into the envoyConfig interface
	if err := json.Unmarshal([]byte(*envoyJSON), &envoyConfig); err != nil {
		wro.Diags.AddError("Envoy Unmarshal Error", fmt.Sprintf("unable to parse Envoy JSON: %s", err))
		return nil
	}

	// Return constructed WorkloadSidecar with parsed Envoy configuration
	return &client.WorkloadSidecar{
		Envoy: &envoyConfig,
	}
}

// buildRolloutOptions constructs a WorkloadRolloutOptions from the given Terraform state.
func (wro *WorkloadResourceOperator) buildRolloutOptions(state []models.RolloutOptionsModel) *client.WorkloadRolloutOptions {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadRolloutOptions{
		MinReadySeconds:               BuildInt(block.MinReadySeconds),
		MaxUnavailableReplicas:        BuildString(block.MaxUnavailableReplicas),
		MaxSurgeReplicas:              BuildString(block.MaxSurgeReplicas),
		ScalingPolicy:                 BuildString(block.ScalingPolicy),
		TerminationGracePeriodSeconds: BuildInt(block.TerminationGracePeriodSeconds),
	}
}

// buildSecurityOptions constructs a WorkloadSecurityOptions from the given Terraform state.
func (wro *WorkloadResourceOperator) buildSecurityOptions(state []models.SecurityOptionsModel) *client.WorkloadSecurityOptions {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadSecurityOptions{
		FileSystemGroupId: BuildInt(block.FileSystemGroupId),
	}
}

// buildLoadBalancer constructs a WorkloadLoadBalancer from the given Terraform state.
func (wro *WorkloadResourceOperator) buildLoadBalancer(state []models.LoadBalancerModel) *client.WorkloadLoadBalancer {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadLoadBalancer{
		Direct:        wro.buildLoadBalancerDirect(block.Direct),
		GeoLocation:   wro.buildLoadBalancerGeoLocation(block.GeoLocation),
		ReplicaDirect: BuildBool(block.ReplicaDirect),
	}
}

// buildLoadBalancerDirect constructs a WorkloadLoadBalancerDirect from the given Terraform state.
func (wro *WorkloadResourceOperator) buildLoadBalancerDirect(state []models.LoadBalancerDirectModel) *client.WorkloadLoadBalancerDirect {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadLoadBalancerDirect{
		Enabled: BuildBool(block.Enabled),
		Ports:   wro.buildLoadBalancerDirectPort(block.Ports),
		IpSet:   wro.BuildLoadBalancerIpSet(block.IpSet, wro.Client.Org),
	}
}

// buildLoadBalancerDirectPort constructs a []client.WorkloadLoadBalancerDirectPort from the given Terraform state.
func (wro *WorkloadResourceOperator) buildLoadBalancerDirectPort(state []models.LoadBalancerDirectPortModel) *[]client.WorkloadLoadBalancerDirectPort {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.WorkloadLoadBalancerDirectPort{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.WorkloadLoadBalancerDirectPort{
			ExternalPort:  BuildInt(block.ExternalPort),
			Protocol:      BuildString(block.Protocol),
			Scheme:        BuildString(block.Scheme),
			ContainerPort: BuildInt(block.ContainerPort),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildLoadBalancerGeoLocation constructs a WorkloadLoadBalancerGeoLocation from the given Terraform state.
func (wro *WorkloadResourceOperator) buildLoadBalancerGeoLocation(state []models.LoadBalancerGeoLocationModel) *client.WorkloadLoadBalancerGeoLocation {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadLoadBalancerGeoLocation{
		Enabled: BuildBool(block.Enabled),
		Headers: wro.buildLoadBalancerGeoLocationHeaders(block.Headers),
	}
}

// buildLoadBalancerGeoLocationHeaders constructs a WorkloadLoadBalancerGeoLocationHeaders from the given Terraform state.
func (wro *WorkloadResourceOperator) buildLoadBalancerGeoLocationHeaders(state []models.LoadBalancerGeoLocationHeadersModel) *client.WorkloadLoadBalancerGeoLocationHeaders {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadLoadBalancerGeoLocationHeaders{
		Asn:     BuildString(block.Asn),
		City:    BuildString(block.City),
		Country: BuildString(block.Country),
		Region:  BuildString(block.Region),
	}
}

// BuildExtras unmarshals a JSON string from Terraform state into a generic extras configuration interface.
func (wro *WorkloadResourceOperator) buildExtras(state types.String) *interface{} {
	// Retrieve raw JSON string from Terraform state
	extrasJSON := BuildString(state)

	// Return nil when no extras JSON is provided
	if extrasJSON == nil {
		return nil
	}

	// Prepare a container for the parsed extras JSON
	var extrasConfig interface{}

	// Unmarshal the JSON into the generic extrasConfig
	if err := json.Unmarshal([]byte(*extrasJSON), &extrasConfig); err != nil {
		wro.Diags.AddError("Extras Unmarshal Error", fmt.Sprintf("unable to parse Extras JSON: %s", err))
		return nil
	}

	// Return a pointer to the parsed extras configuration
	return &extrasConfig
}

// buildRequestRetryPolicy constructs a WorkloadRequestRetryPolicy from the given Terraform state.
func (wro *WorkloadResourceOperator) buildRequestRetryPolicy(state []models.RequestRetryPolicyModel) *client.WorkloadRequestRetryPolicy {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.WorkloadRequestRetryPolicy{
		Attempts: BuildInt(block.Attempts),
		RetryOn:  wro.BuildSetString(block.RetryOn),
	}
}

// Flatteners //

// flattenContainers transforms *[]client.WorkloadContainer into a []models.ContainerModel.
func (wro *WorkloadResourceOperator) flattenContainers(input *[]client.WorkloadContainer) []models.ContainerModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.ContainerModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Initialize legacyPort to null by default
		legacyPort := types.Int32Null()

		// Declare planPort to capture potential legacy port from the plan
		var planPort types.Int32

		// Flag to indicate if a legacy port was found
		var hasLegacy bool

		// Iterate through plan containers to find matching container by name
		for _, c := range wro.Plan.Containers {
			// Check if this plan container corresponds to the current item
			if c.Name.ValueString() == *item.Name {
				// If a valid legacy port is set in the plan, capture it
				if !c.Port.IsNull() && !c.Port.IsUnknown() {
					planPort = c.Port
					hasLegacy = true
				}

				// Stop searching once the matching container is handled
				break
			}
		}

		// Prepare slice to hold ports excluding any legacy port
		var filtered *[]client.WorkloadContainerPort = nil

		// Only proceed if the container has defined ports
		if item.Ports != nil {
			// Initialize a filtered slice to hold ports excluding legacy port
			filtered = &[]client.WorkloadContainerPort{}

			// Loop through each port in the container
			for _, p := range *item.Ports {
				// If this port matches the identified legacy port criteria, record it
				if hasLegacy && *p.Protocol == "http" && *p.Number == *BuildInt(planPort) {
					legacyPort = planPort
				} else {
					// Otherwise, retain the port in the filtered slice
					*filtered = append(*filtered, p)
				}
			}
		}

		// Check if the legacy port is set within the input container and use it
		if item.Port != nil {
			legacyPort = types.Int32Value(int32(*item.Port))
		}

		// Construct a block
		block := models.ContainerModel{
			Name:             types.StringPointerValue(item.Name),
			Image:            types.StringPointerValue(item.Image),
			WorkingDirectory: types.StringPointerValue(item.WorkingDirectory),
			Metrics:          wro.flattenContainerMetrics(item.Metrics),
			Port:             legacyPort,
			Ports:            wro.flattenContainerPort(filtered),
			Memory:           types.StringPointerValue(item.Memory),
			ReadinessProbe:   wro.flattenHealthCheck(item.ReadinessProbe),
			LivenessProbe:    wro.flattenHealthCheck(item.LivenessProbe),
			Cpu:              types.StringPointerValue(item.CPU),
			MinCpu:           types.StringPointerValue(item.MinCPU),
			MinMemory:        types.StringPointerValue(item.MinMemory),
			Env:              FlattenMapString(wro.flattenNameValue(item.Env)),
			GpuNvidia:        wro.flattenContainerGpuNvidia(item.GPU),
			GpuCustom:        wro.flattenContainerGpuCustom(item.GPU),
			InheritEnv:       types.BoolPointerValue(item.InheritEnv),
			Command:          types.StringPointerValue(item.Command),
			Args:             FlattenSetString(item.Args),
			Lifecycle:        wro.flattenContainerLifecycle(item.LifeCycle),
			Volumes:          wro.flattenContainerVolume(item.Volumes),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenContainerMetrics transforms *client.WorkloadContainerMetrics into a []models.ContainerMetricsModel.
func (wro *WorkloadResourceOperator) flattenContainerMetrics(input *client.WorkloadContainerMetrics) []models.ContainerMetricsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerMetricsModel{
		Port: FlattenInt(input.Port),
		Path: types.StringPointerValue(input.Path),
	}

	// Return a slice containing the single block
	return []models.ContainerMetricsModel{block}
}

// flattenContainerPort transforms *[]client.WorkloadContainerPort into a []models.ContainerPortModel.
func (wro *WorkloadResourceOperator) flattenContainerPort(input *[]client.WorkloadContainerPort) []models.ContainerPortModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.ContainerPortModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.ContainerPortModel{
			Protocol: types.StringPointerValue(item.Protocol),
			Number:   FlattenInt(item.Number),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenHealthCheck transforms *client.WorkloadHealthCheck into a []models.ContainerHealthCheckModel.
func (wro *WorkloadResourceOperator) flattenHealthCheck(input *client.WorkloadHealthCheck) []models.ContainerHealthCheckModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerHealthCheckModel{
		Exec:                wro.flattenExec(input.Exec),
		Grpc:                wro.flattenHealthCheckGrpc(input.GRPC),
		TcpSocket:           wro.flattenHealthCheckTcpSocket(input.TCPSocket),
		HttpGet:             wro.flattenHealthCheckHttpGet(input.HTTPGet),
		InitialDelaySeconds: FlattenInt(input.InitialDelaySeconds),
		PeriodSeconds:       FlattenInt(input.PeriodSeconds),
		TimeoutSeconds:      FlattenInt(input.TimeoutSeconds),
		SuccessThreshold:    FlattenInt(input.SuccessThreshold),
		FailureThreshold:    FlattenInt(input.FailureThreshold),
	}

	// Return a slice containing the single block
	return []models.ContainerHealthCheckModel{block}
}

// flattenExec transforms *client.WorkloadExec into a []models.ContainerExecModel.
func (wro *WorkloadResourceOperator) flattenExec(input *client.WorkloadExec) []models.ContainerExecModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerExecModel{
		Command: FlattenSetString(input.Command),
	}

	// Return a slice containing the single block
	return []models.ContainerExecModel{block}
}

// flattenHealthCheckGrpc transforms *client.WorkloadHealthCheckGrpc into a []models.ContainerHealthCheckGrpcModel.
func (wro *WorkloadResourceOperator) flattenHealthCheckGrpc(input *client.WorkloadHealthCheckGrpc) []models.ContainerHealthCheckGrpcModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerHealthCheckGrpcModel{
		Port: FlattenInt(input.Port),
	}

	// Return a slice containing the single block
	return []models.ContainerHealthCheckGrpcModel{block}
}

// flattenHealthCheckTcpSocket transforms *client.WorkloadHealthCheckTcpSocket into a []models.ContainerHealthCheckTcpSocketModel.
func (wro *WorkloadResourceOperator) flattenHealthCheckTcpSocket(input *client.WorkloadHealthCheckTcpSocket) []models.ContainerHealthCheckTcpSocketModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerHealthCheckTcpSocketModel{
		Port: FlattenInt(input.Port),
	}

	// Return a slice containing the single block
	return []models.ContainerHealthCheckTcpSocketModel{block}
}

// flattenHealthCheckHttpGet transforms *client.WorkloadHealthCheckHttpGet into a []models.ContainerHealthCheckHttpGetModel.
func (wro *WorkloadResourceOperator) flattenHealthCheckHttpGet(input *client.WorkloadHealthCheckHttpGet) []models.ContainerHealthCheckHttpGetModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerHealthCheckHttpGetModel{
		Path:        types.StringPointerValue(input.Path),
		Port:        FlattenInt(input.Port),
		HttpHeaders: FlattenMapString(wro.flattenNameValue(input.HttpHeaders)),
		Scheme:      types.StringPointerValue(input.Scheme),
	}

	// Return a slice containing the single block
	return []models.ContainerHealthCheckHttpGetModel{block}
}

// flattenNameValue transforms *[]client.WorkloadContainerNameValue into a map[string]interface{}.
func (wro *WorkloadResourceOperator) flattenNameValue(input *[]client.WorkloadContainerNameValue) *map[string]interface{} {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Prepare the output slice
	output := map[string]interface{}{}

	// Iterate over each block and map the name to its value
	for _, item := range *input {
		// Skip this record just in case the name was nil
		if item.Name == nil {
			continue
		}

		// Dereference the record name
		key := *item.Name

		// Initialize with a nil value
		output[key] = nil

		// If the value is not nil, update the output key
		if item.Value != nil {
			output[key] = *item.Value
		}
	}

	// Return the constructed output slice
	return &output
}

// flattenContainerGpuNvidia transforms *client.WorkloadContainerGpu into a []models.ContainerGpuNvidiaModel.
func (wro *WorkloadResourceOperator) flattenContainerGpuNvidia(input *client.WorkloadContainerGpu) []models.ContainerGpuNvidiaModel {
	// Check if the input is nil
	if input == nil || input.Nvidia == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerGpuNvidiaModel{
		Model:    types.StringPointerValue(input.Nvidia.Model),
		Quantity: FlattenInt(input.Nvidia.Quantity),
	}

	// Return a slice containing the single block
	return []models.ContainerGpuNvidiaModel{block}
}

// flattenContainerGpuCustom transforms *client.WorkloadContainerGpu into a []models.ContainerGpuCustomModel.
func (wro *WorkloadResourceOperator) flattenContainerGpuCustom(input *client.WorkloadContainerGpu) []models.ContainerGpuCustomModel {
	// Check if the input is nil
	if input == nil || input.Custom == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerGpuCustomModel{
		Resource:     types.StringPointerValue(input.Custom.Resource),
		RuntimeClass: types.StringPointerValue(input.Custom.RuntimeClass),
		Quantity:     FlattenInt(input.Custom.Quantity),
	}

	// Return a slice containing the single block
	return []models.ContainerGpuCustomModel{block}
}

// flattenContainerLifecycle transforms *client.WorkloadLifeCycle into a []models.ContainerLifecycleModel.
func (wro *WorkloadResourceOperator) flattenContainerLifecycle(input *client.WorkloadLifeCycle) []models.ContainerLifecycleModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerLifecycleModel{
		PostStart: wro.flattenContainerLifecycleSpec(input.PostStart),
		PreStop:   wro.flattenContainerLifecycleSpec(input.PreStop),
	}

	// Return a slice containing the single block
	return []models.ContainerLifecycleModel{block}
}

// flattenContainerLifecycleSpec transforms *client.WorkloadLifeCycleSpec into a []models.ContainerLifecycleSpecModel.
func (wro *WorkloadResourceOperator) flattenContainerLifecycleSpec(input *client.WorkloadLifeCycleSpec) []models.ContainerLifecycleSpecModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ContainerLifecycleSpecModel{
		Exec: wro.flattenExec(input.Exec),
	}

	// Return a slice containing the single block
	return []models.ContainerLifecycleSpecModel{block}
}

// flattenContainerVolume transforms *[]client.WorkloadContainerVolume into a []models.ContainerVolumeModel.
func (wro *WorkloadResourceOperator) flattenContainerVolume(input *[]client.WorkloadContainerVolume) []models.ContainerVolumeModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.ContainerVolumeModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.ContainerVolumeModel{
			Uri:            types.StringPointerValue(item.Uri),
			RecoveryPolicy: types.StringPointerValue(item.RecoveryPolicy),
			Path:           types.StringPointerValue(item.Path),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenFirewall transforms *client.WorkloadFirewall into a []models.FirewallModel.
func (wro *WorkloadResourceOperator) flattenFirewall(input *client.WorkloadFirewall) []models.FirewallModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.FirewallModel{
		External: wro.flattenFirewallExternal(input.External),
		Internal: wro.flattenFirewallInternal(input.Internal),
	}

	// Return a slice containing the single block
	return []models.FirewallModel{block}
}

// flattenFirewallExternal transforms *client.WorkloadFirewallExternal into a []models.FirewallExternalModel.
func (wro *WorkloadResourceOperator) flattenFirewallExternal(input *client.WorkloadFirewallExternal) []models.FirewallExternalModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.FirewallExternalModel{
		InboundAllowCidr:      FlattenSetString(input.InboundAllowCidr),
		InboundBlockedCidr:    FlattenSetString(input.InboundBlockedCidr),
		OutboundAllowHostname: FlattenSetString(input.OutboundAllowHostname),
		OutboundAllowPort:     wro.flattenFirewallExternalOutboundAllowPort(input.OutboundAllowPort),
		OutboundAllowCidr:     FlattenSetString(input.OutboundAllowCidr),
		OutboundBlockedCidr:   FlattenSetString(input.OutboundBlockedCidr),
		Http:                  wro.flattenFirewallExternalHttp(input.Http),
	}

	// Return a slice containing the single block
	return []models.FirewallExternalModel{block}
}

// flattenFirewallExternalOutboundAllowPort transforms *[]client.WorkloadFirewallOutboundAllowPort into a []models.FirewallExternalOutboundAllowPortModel.
func (wro *WorkloadResourceOperator) flattenFirewallExternalOutboundAllowPort(input *[]client.WorkloadFirewallOutboundAllowPort) []models.FirewallExternalOutboundAllowPortModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.FirewallExternalOutboundAllowPortModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.FirewallExternalOutboundAllowPortModel{
			Protocol: types.StringPointerValue(item.Protocol),
			Number:   FlattenInt(item.Number),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenFirewallExternalHttp transforms *client.WorkloadFirewallExternalHttp into a []models.FirewallExternalHttpModel.
func (wro *WorkloadResourceOperator) flattenFirewallExternalHttp(input *client.WorkloadFirewallExternalHttp) []models.FirewallExternalHttpModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.FirewallExternalHttpModel{
		InboundHeaderFilter: wro.flattenFirewallExternalHttpHeaderFilter(input.InboundHeaderFilter),
	}

	// Return a slice containing the single block
	return []models.FirewallExternalHttpModel{block}
}

// flattenFirewallExternalHttpHeaderFilter transforms *[]client.WorkloadFirewallExternalHttpHeaderFilter into a []models.FirewallExternalHttpHeaderFilterModel.
func (wro *WorkloadResourceOperator) flattenFirewallExternalHttpHeaderFilter(input *[]client.WorkloadFirewallExternalHttpHeaderFilter) []models.FirewallExternalHttpHeaderFilterModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.FirewallExternalHttpHeaderFilterModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.FirewallExternalHttpHeaderFilterModel{
			Key:           types.StringPointerValue(item.Key),
			AllowedValues: FlattenSetString(item.AllowedValues),
			BlockedValues: FlattenSetString(item.BlockedValues),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenFirewallInternal transforms *client.WorkloadFirewallInternal into a []models.FirewallInternalModel.
func (wro *WorkloadResourceOperator) flattenFirewallInternal(input *client.WorkloadFirewallInternal) []models.FirewallInternalModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.FirewallInternalModel{
		InboundAllowType:     types.StringPointerValue(input.InboundAllowType),
		InboundAllowWorkload: FlattenSetString(input.InboundAllowWorkload),
	}

	// Return a slice containing the single block
	return []models.FirewallInternalModel{block}
}

// flattenOptions transforms *client.WorkloadOptions into a []models.OptionsModel.
func (wro *WorkloadResourceOperator) flattenOptions(input *client.WorkloadOptions) []models.OptionsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.OptionsModel{
		Autoscaling:    wro.flattenOptionsAutoscaling(input.AutoScaling),
		TimeoutSeconds: FlattenInt(input.TimeoutSeconds),
		CapacityAI:     types.BoolPointerValue(input.CapacityAI),
		Debug:          types.BoolPointerValue(input.Debug),
		Suspend:        types.BoolPointerValue(input.Suspend),
		MultiZone:      wro.flattenOptionsMultiZone(input.MultiZone),
	}

	// Return a slice containing the single block
	return []models.OptionsModel{block}
}

// flattenOptionsAutoscaling transforms *client.WorkloadOptionsAutoscaling into a []models.OptionsAutoscalingModel.
func (wro *WorkloadResourceOperator) flattenOptionsAutoscaling(input *client.WorkloadOptionsAutoscaling) []models.OptionsAutoscalingModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.OptionsAutoscalingModel{
		Metric:           types.StringPointerValue(input.Metric),
		Multi:            wro.flattenOptionsAutoscalingMulti(input.Multi),
		MetricPercentile: types.StringPointerValue(input.MetricPercentile),
		Target:           FlattenInt(input.Target),
		MaxScale:         FlattenInt(input.MaxScale),
		MinScale:         FlattenInt(input.MinScale),
		MaxConcurrency:   FlattenInt(input.MaxConcurrency),
		ScaleToZeroDelay: FlattenInt(input.ScaleToZeroDelay),
	}

	// Return a slice containing the single block
	return []models.OptionsAutoscalingModel{block}
}

// flattenOptionsAutoscalingMulti transforms *[]client.WorkloadOptionsAutoscalingMulti into a []models.OptionsAutoscalingMultiModel.
func (wro *WorkloadResourceOperator) flattenOptionsAutoscalingMulti(input *[]client.WorkloadOptionsAutoscalingMulti) []models.OptionsAutoscalingMultiModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.OptionsAutoscalingMultiModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.OptionsAutoscalingMultiModel{
			Metric: types.StringPointerValue(item.Metric),
			Target: FlattenInt(item.Target),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenOptionsMultiZone transforms *client.WorkloadOptionsMultiZone into a []models.OptionsMultiZoneModel.
func (wro *WorkloadResourceOperator) flattenOptionsMultiZone(input *client.WorkloadOptionsMultiZone) []models.OptionsMultiZoneModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.OptionsMultiZoneModel{
		Enabled: types.BoolPointerValue(input.Enabled),
	}

	// Return a slice containing the single block
	return []models.OptionsMultiZoneModel{block}
}

// flattenLocalOptions transforms *[]client.WorkloadOptions into a []models.LocalOptionsModel.
func (wro *WorkloadResourceOperator) flattenLocalOptions(input *[]client.WorkloadOptions) []models.LocalOptionsModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.LocalOptionsModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		options := wro.flattenOptions(&item)

		// If the block is nil, skip it
		if len(options) == 0 {
			continue
		}

		// Construct the local options block
		block := models.LocalOptionsModel{
			Location:     types.StringNull(),
			OptionsModel: options[0],
		}

		// Flatten the location
		if item.Location != nil {
			block.Location = types.StringValue(GetNameFromSelfLink(*item.Location))
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenJob transforms *client.WorkloadJob into a []models.JobModel.
func (wro *WorkloadResourceOperator) flattenJob(input *client.WorkloadJob) []models.JobModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.JobModel{
		Schedule:              types.StringPointerValue(input.Schedule),
		ConcurrencyPolicy:     types.StringPointerValue(input.ConcurrencyPolicy),
		HistoryLimit:          FlattenInt(input.HistoryLimit),
		RestartPolicy:         types.StringPointerValue(input.RestartPolicy),
		ActiveDeadlineSeconds: FlattenInt(input.ActiveDeadlineSeconds),
	}

	// Return a slice containing the single block
	return []models.JobModel{block}
}

// flattenSidecar transforms *client.WorkloadSidecar into a []models.SidecarModel.
func (wro *WorkloadResourceOperator) flattenSidecar(input *client.WorkloadSidecar) []models.SidecarModel {
	// Return nil when no sidecar is present
	if input == nil {
		return nil
	}

	// Initialize an empty SidecarModel block
	block := models.SidecarModel{}

	// Return a slice with an empty block if Envoy config is missing
	if input.Envoy == nil {
		return []models.SidecarModel{block}
	}

	// Marshal the Envoy configuration back to JSON
	jsonOut, err := json.Marshal(*input.Envoy)

	// Handle any errors that occur during marshaling
	if err != nil {
		wro.Diags.AddError("Envoy Marshaling Error", fmt.Sprintf("error occurred during marshaling 'envoy' attribute. Error: %s", err.Error()))
		return nil
	}

	// Assign the JSON string to the Envoy field on the block
	block.Envoy = types.StringValue(string(jsonOut))

	// Return a slice containing the populated block
	return []models.SidecarModel{block}
}

// flattenRolloutOptions transforms *client.WorkloadRolloutOptions into a []models.RolloutOptionsModel.
func (wro *WorkloadResourceOperator) flattenRolloutOptions(input *client.WorkloadRolloutOptions) []models.RolloutOptionsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.RolloutOptionsModel{
		MinReadySeconds:               FlattenInt(input.MinReadySeconds),
		MaxUnavailableReplicas:        types.StringPointerValue(input.MaxUnavailableReplicas),
		MaxSurgeReplicas:              types.StringPointerValue(input.MaxSurgeReplicas),
		ScalingPolicy:                 types.StringPointerValue(input.ScalingPolicy),
		TerminationGracePeriodSeconds: FlattenInt(input.TerminationGracePeriodSeconds),
	}

	// Return a slice containing the single block
	return []models.RolloutOptionsModel{block}
}

// flattenSecurityOptions transforms *client.WorkloadSecurityOptions into a []models.SecurityOptionsModel.
func (wro *WorkloadResourceOperator) flattenSecurityOptions(input *client.WorkloadSecurityOptions) []models.SecurityOptionsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.SecurityOptionsModel{
		FileSystemGroupId: FlattenInt(input.FileSystemGroupId),
	}

	// Return a slice containing the single block
	return []models.SecurityOptionsModel{block}
}

// flattenLoadBalancer transforms *client.WorkloadLoadBalancer into a []models.LoadBalancerModel.
func (wro *WorkloadResourceOperator) flattenLoadBalancer(state []models.LoadBalancerModel, input *client.WorkloadLoadBalancer) []models.LoadBalancerModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Initialize direct list from existing state
	direct := []models.LoadBalancerDirectModel{}

	// Preserve previous direct configuration if present
	if len(state) > 0 {
		direct = state[0].Direct
	}

	// Build a single block
	block := models.LoadBalancerModel{
		Direct:        wro.flattenLoadBalancerDirect(direct, input.Direct),
		GeoLocation:   wro.flattenLoadBalancerGeoLocation(input.GeoLocation),
		ReplicaDirect: types.BoolPointerValue(input.ReplicaDirect),
	}

	// Return a slice containing the single block
	return []models.LoadBalancerModel{block}
}

// flattenLoadBalancerDirect transforms *client.WorkloadLoadBalancerDirect into a []models.LoadBalancerDirectModel.
func (wro *WorkloadResourceOperator) flattenLoadBalancerDirect(state []models.LoadBalancerDirectModel, input *client.WorkloadLoadBalancerDirect) []models.LoadBalancerDirectModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Initialize the ipSetState to null
	ipSetState := types.StringNull()

	// If the state contains blocks, extract the IP Set from the first block
	if len(state) > 0 {
		ipSetState = state[0].IpSet
	}

	// Build a single block
	block := models.LoadBalancerDirectModel{
		Enabled: types.BoolPointerValue(input.Enabled),
		Ports:   wro.flattenLoadBalancerDirectPort(input.Ports),
		IpSet:   wro.FlattenLoadBalancerIpSet(ipSetState, input.IpSet, wro.Client.Org),
	}

	// Return a slice containing the single block
	return []models.LoadBalancerDirectModel{block}
}

// flattenLoadBalancerDirectPort transforms *[]client.WorkloadLoadBalancerDirectPort into a []models.LoadBalancerDirectPortModel.
func (wro *WorkloadResourceOperator) flattenLoadBalancerDirectPort(input *[]client.WorkloadLoadBalancerDirectPort) []models.LoadBalancerDirectPortModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.LoadBalancerDirectPortModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.LoadBalancerDirectPortModel{
			ExternalPort:  FlattenInt(item.ExternalPort),
			Protocol:      types.StringPointerValue(item.Protocol),
			Scheme:        types.StringPointerValue(item.Scheme),
			ContainerPort: FlattenInt(item.ContainerPort),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenLoadBalancerGeoLocation transforms *client.WorkloadLoadBalancerGeoLocation into a []models.LoadBalancerGeoLocationModel.
func (wro *WorkloadResourceOperator) flattenLoadBalancerGeoLocation(input *client.WorkloadLoadBalancerGeoLocation) []models.LoadBalancerGeoLocationModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.LoadBalancerGeoLocationModel{
		Enabled: types.BoolPointerValue(input.Enabled),
		Headers: wro.flattenLoadBalancerGeoLocationHeaders(input.Headers),
	}

	// Return a slice containing the single block
	return []models.LoadBalancerGeoLocationModel{block}
}

// flattenLoadBalancerGeoLocationHeaders transforms *client.WorkloadLoadBalancerGeoLocationHeaders into a []models.LoadBalancerGeoLocationHeadersModel.
func (wro *WorkloadResourceOperator) flattenLoadBalancerGeoLocationHeaders(input *client.WorkloadLoadBalancerGeoLocationHeaders) []models.LoadBalancerGeoLocationHeadersModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.LoadBalancerGeoLocationHeadersModel{
		Asn:     types.StringPointerValue(input.Asn),
		City:    types.StringPointerValue(input.City),
		Country: types.StringPointerValue(input.Country),
		Region:  types.StringPointerValue(input.Region),
	}

	// Return a slice containing the single block
	return []models.LoadBalancerGeoLocationHeadersModel{block}
}

// FlattenExtras marshals extras into a JSON string or returns null when input is nil.
func (wro *WorkloadResourceOperator) flattenExtras(input *interface{}) types.String {
	// Return null string when extras is not provided
	if input == nil {
		return types.StringNull()
	}

	// Marshal the extras configuration into JSON
	jsonOut, err := json.Marshal(*input)
	// Report error and return null when marshaling fails
	if err != nil {
		wro.Diags.AddError("Extras Marshaling Error", fmt.Sprintf("error occurred during marshaling 'extras' attribute. Error: %s", err.Error()))
		return types.StringNull()
	}

	// Return the JSON string as a Terraform string value
	return types.StringValue(string(jsonOut))
}

// flattenRequestRetryPolicy transforms *client.WorkloadRequestRetryPolicy into a []models.RequestRetryPolicyModel.
func (wro *WorkloadResourceOperator) flattenRequestRetryPolicy(input *client.WorkloadRequestRetryPolicy) []models.RequestRetryPolicyModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.RequestRetryPolicyModel{
		Attempts: FlattenInt(input.Attempts),
		RetryOn:  FlattenSetString(input.RetryOn),
	}

	// Return a slice containing the single block
	return []models.RequestRetryPolicyModel{block}
}

// flattenStatus transforms *client.WorkloadStatus into a Terraform types.List.
func (wro *WorkloadResourceOperator) flattenStatus(input *client.WorkloadStatus) types.List {
	// Get attribute types
	elementType := models.StatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusModel{
		ParentId:            types.StringPointerValue(input.ParentId),
		CanonicalEndpoint:   types.StringPointerValue(input.CanonicalEndpoint),
		Endpoint:            types.StringPointerValue(input.Endpoint),
		InternalName:        types.StringPointerValue(input.InternalName),
		HealthCheck:         wro.flattenStatusHealthCheck(input.HealthCheck),
		CurrentReplicaCount: FlattenInt(input.CurrentReplicaCount),
		ResolvedImages:      wro.flattenStatusResolvedImages(input.ResolvedImages),
		LoadBalancer:        wro.flattenStatusLoadBalancer(input.LoadBalancer),
	}

	// Return the successfully created types.List
	return FlattenList(wro.Ctx, wro.Diags, []models.StatusModel{block})
}

// flattenStatusHealthCheck transforms *client.WorkloadStatusHealthCheck into a Terraform types.List.
func (wro *WorkloadResourceOperator) flattenStatusHealthCheck(input *client.WorkloadStatusHealthCheck) types.List {
	// Get attribute types
	elementType := models.StatusHealthCheckModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusHealthCheckModel{
		Active:      types.BoolPointerValue(input.Active),
		Success:     types.BoolPointerValue(input.Success),
		Code:        FlattenInt(input.Code),
		Message:     types.StringPointerValue(input.Message),
		Failures:    FlattenInt(input.Failures),
		Successes:   FlattenInt(input.Successes),
		LastChecked: types.StringPointerValue(input.LastChecked),
	}

	// Return the successfully created types.List
	return FlattenList(wro.Ctx, wro.Diags, []models.StatusHealthCheckModel{block})
}

// flattenStatusResolvedImages transforms *client.WorkloadStatusResolvedImages into a Terraform types.List.
func (wro *WorkloadResourceOperator) flattenStatusResolvedImages(input *client.WorkloadStatusResolvedImages) types.List {
	// Get attribute types
	elementType := models.StatusResolvedImagesModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusResolvedImagesModel{
		ResolvedForVersion: FlattenInt(input.ResolvedForVersion),
		ResolvedAt:         types.StringPointerValue(input.ResolvedAt),
		ErrorMessages:      FlattenSetString(input.ErrorMessages),
		Images:             wro.flattenStatusResolvedImage(input.Images),
	}

	// Return the successfully created types.List
	return FlattenList(wro.Ctx, wro.Diags, []models.StatusResolvedImagesModel{block})
}

// flattenStatusResolvedImage transforms *[]client.WorkloadStatusResolvedImage into a Terraform types.List.
func (wro *WorkloadResourceOperator) flattenStatusResolvedImage(input *[]client.WorkloadStatusResolvedImage) types.List {
	// Get attribute types
	elementType := models.StatusResolvedImageModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.StatusResolvedImageModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.StatusResolvedImageModel{
			Digest:    types.StringPointerValue(item.Digest),
			Manifests: wro.flattenStatusResolvedImageManifest(item.Manifests),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(wro.Ctx, wro.Diags, blocks)
}

// flattenStatusResolvedImageManifest transforms *[]client.WorkloadStatusResolvedImageManifest into a Terraform types.List.
func (wro *WorkloadResourceOperator) flattenStatusResolvedImageManifest(input *[]client.WorkloadStatusResolvedImageManifest) types.List {
	// Get attribute types
	elementType := models.StatusResolvedImageManifestModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.StatusResolvedImageManifestModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.StatusResolvedImageManifestModel{
			Image:     types.StringPointerValue(item.Image),
			MediaType: types.StringPointerValue(item.MediaType),
			Digest:    types.StringPointerValue(item.Digest),
			Platform:  FlattenMapString(item.Platform),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(wro.Ctx, wro.Diags, blocks)
}

// flattenStatusLoadBalancer transforms *[]client.WorkloadStatusLoadBalancer into a Terraform types.List.
func (wro *WorkloadResourceOperator) flattenStatusLoadBalancer(input *[]client.WorkloadStatusLoadBalancer) types.List {
	// Get attribute types
	elementType := models.StatusLoadBalancerModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.StatusLoadBalancerModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.StatusLoadBalancerModel{
			Origin: types.StringPointerValue(item.Origin),
			Url:    types.StringPointerValue(item.Url),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(wro.Ctx, wro.Diags, blocks)
}
