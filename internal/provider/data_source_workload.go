package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &WorkloadDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkloadDataSource{}
)

// WorkloadDataSource is the data source implementation.
type WorkloadDataSource struct {
	EntityBase
	Operations EntityOperations[WorkloadResourceModel, client.Workload]
}

// NewWorkloadDataSource returns a new instance of the data source implementation.
func NewWorkloadDataSource() datasource.DataSource {
	return &WorkloadDataSource{}
}

// Metadata provides the data source type name.
func (d *WorkloadDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_workload"
}

// Configure configures the data source before use.
func (d *WorkloadDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	d.Operations = NewEntityOperations(d.client, &WorkloadResourceOperator{})
}

// Schema defines the schema for the data source.
func (d *WorkloadDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"gvc": schema.StringAttribute{
				Description: "Name of the associated GVC.",
				Required:    true,
				Validators: []validator.String{
					validators.NameValidator{},
				},
			},
			"type": schema.StringAttribute{
				Description: "Workload Type. Either `serverless`, `standard`, `stateful`, or `cron`.",
				Computed:    true,
			},
			"identity_link": schema.StringAttribute{
				Description: "The identityLink is used as the access scope for 3rd party cloud resources. A single identity can provide access to multiple cloud providers.",
				Computed:    true,
			},
			"support_dynamic_tags": schema.BoolAttribute{
				Description: "Workload will automatically redeploy when one of the container images is updated in the container registry. Default: false.",
				Computed:    true,
			},
			"extras": schema.StringAttribute{
				Description: "Extra Kubernetes modifications. Only used for BYOK.",
				Computed:    true,
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
						"replica_internal_names": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
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
									"next_retry_at": schema.StringAttribute{
										Description: "",
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
		},
		Blocks: map[string]schema.Block{
			"container": schema.ListNestedBlock{
				Description: "An isolated and lightweight runtime environment that encapsulates an application and its dependencies.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the container.",
							Computed:    true,
						},
						"image": schema.StringAttribute{
							Description: "The full image and tag path.",
							Computed:    true,
						},
						"working_directory": schema.StringAttribute{
							Description: "Override the working directory. Must be an absolute path.",
							Computed:    true,
						},
						"port": schema.Int32Attribute{
							Description:        "The port the container exposes. Only one container is allowed to specify a port. Min: `80`. Max: `65535`. Used by `serverless` Workload type. **DEPRECATED - Use `ports`.**",
							DeprecationMessage: "The 'port' attribute will be deprecated in the next major version. Use the 'ports' attribute instead.",
							Computed:           true,
						},
						"memory": schema.StringAttribute{
							Description: "Reserved memory of the workload when capacityAI is disabled. Maximum memory when CapacityAI is enabled. Default: \"128Mi\".",
							Computed:    true,
						},
						"cpu": schema.StringAttribute{
							Description: "Reserved CPU of the workload when capacityAI is disabled. Maximum CPU when CapacityAI is enabled. Default: \"50m\".",
							Computed:    true,
						},
						"min_cpu": schema.StringAttribute{
							Description: "Minimum CPU when capacity AI is enabled.",
							Computed:    true,
						},
						"min_memory": schema.StringAttribute{
							Description: "Minimum memory when capacity AI is enabled.",
							Computed:    true,
						},
						"env": schema.MapAttribute{
							Description: "Name-Value list of environment variables.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"inherit_env": schema.BoolAttribute{
							Description: "Enables inheritance of GVC environment variables. A variable in spec.env will override a GVC variable with the same name.",
							Computed:    true,
						},
						"command": schema.StringAttribute{
							Description: "Override the entry point.",
							Computed:    true,
						},
						"args": schema.ListAttribute{
							Description: "Command line arguments passed to the container at runtime. Replaces the CMD arguments of the running container. It is an ordered list.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"metrics": schema.ListNestedBlock{
							MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/workload#metrics).",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"port": schema.Int32Attribute{
										Description: "Port from container emitting custom metrics.",
										Computed:    true,
									},
									"path": schema.StringAttribute{
										Description: "Path from container emitting custom metrics.",
										Computed:    true,
									},
									"drop_metrics": schema.SetAttribute{
										Description: "Drop metrics that match given patterns.",
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
						"ports": schema.ListNestedBlock{
							Description: "Communication endpoints used by the workload to send and receive network traffic.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"protocol": schema.StringAttribute{
										Description: "Protocol. Choice of: `http`, `http2`, `tcp`, or `grpc`.",
										Computed:    true,
									},
									"number": schema.Int32Attribute{
										Description: "Port to expose.",
										Computed:    true,
									},
								},
							},
						},
						"readiness_probe": d.HealthCheckSchema("readiness_probe", "Readiness Probe"),
						"liveness_probe":  d.HealthCheckSchema("liveness_probe", "Liveness Probe"),
						"gpu_nvidia": schema.ListNestedBlock{
							Description: "GPUs manufactured by NVIDIA, which are specialized hardware accelerators used to offload and accelerate computationally intensive tasks within the workload.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"model": schema.StringAttribute{
										Description: "GPU Model (i.e.: t4)",
										Computed:    true,
									},
									"quantity": schema.Int32Attribute{
										Description: "Number of GPUs.",
										Computed:    true,
									},
								},
							},
						},
						"gpu_custom": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"resource": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"runtime_class": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"quantity": schema.Int32Attribute{
										Description: "Number of GPUs.",
										Computed:    true,
									},
								},
							},
						},
						"lifecycle": schema.ListNestedBlock{
							MarkdownDescription: "Lifecycle [Reference Page](https://docs.controlplane.com/reference/workload#lifecycle).",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"post_start": d.LifecycleSpecSchema("Command and arguments executed immediately after the container is created."),
									"pre_stop":   d.LifecycleSpecSchema("Command and arguments executed immediately before the container is stopped."),
								},
							},
						},
						"volume": schema.SetNestedBlock{
							MarkdownDescription: "Mount Object Store (S3, GCS, AzureBlob) buckets as file system.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"uri": schema.StringAttribute{
										Description: "URI of a volume hosted at Control Plane (Volume Set) or at a cloud provider (AWS, Azure, GCP).",
										Computed:    true,
									},
									"recovery_policy": schema.StringAttribute{
										Description: "Only applicable to persistent volumes, this determines what Control Plane will do when creating a new workload replica if a corresponding volume exists. Available Values: `retain`, `recycle`. Default: `retain`. **DEPRECATED - No longer being used.**",
										Computed:    true,
									},
									"path": schema.StringAttribute{
										Description: "File path added to workload pointing to the volume.",
										Computed:    true,
									},
								},
							},
						},
					},
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
										Computed:    true,
									},
									"inbound_blocked_cidr": schema.SetAttribute{
										Description: "The list of ipv4/ipv6 addresses or cidr blocks that are NOT allowed to access this workload. Addresses in the allow list will only be allowed if they do not exist in this list.",
										ElementType: types.StringType,
										Computed:    true,
									},
									"outbound_allow_hostname": schema.SetAttribute{
										Description: "The list of public hostnames that this workload is allowed to reach. No outbound access is allowed by default. A wildcard `*` is allowed on the prefix of the hostname only, ex: `*.amazonaws.com`. Use `outboundAllowCIDR` to allow access to all external websites.",
										ElementType: types.StringType,
										Computed:    true,
									},
									"outbound_allow_cidr": schema.SetAttribute{
										Description: "The list of ipv4/ipv6 addresses or cidr blocks that this workload is allowed reach. No outbound access is allowed by default. Specify '0.0.0.0/0' to allow outbound access to the public internet.",
										ElementType: types.StringType,
										Computed:    true,
									},
									"outbound_blocked_cidr": schema.SetAttribute{
										Description: "The list of ipv4/ipv6 addresses or cidr blocks that this workload is NOT allowed to reach. Addresses in the allow list will only be allowed if they do not exist in this list.",
										ElementType: types.StringType,
										Computed:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"outbound_allow_port": schema.SetNestedBlock{
										Description: "Allow outbound access to specific ports and protocols. When not specified, communication to address ranges in outboundAllowCIDR is allowed on all ports and communication to names in outboundAllowHostname is allowed on ports 80/443.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"protocol": schema.StringAttribute{
													Description: "Either `http`, `https` or `tcp`.",
													Computed:    true,
												},
												"number": schema.Int32Attribute{
													Description: "Port number. Max: 65000",
													Computed:    true,
												},
											},
										},
									},
									"http": schema.ListNestedBlock{
										Description: "Firewall options for HTTP workloads.",
										NestedObject: schema.NestedBlockObject{
											Blocks: map[string]schema.Block{
												"inbound_header_filter": schema.SetNestedBlock{
													Description: "A list of header filters for HTTP workloads.",
													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"key": schema.StringAttribute{
																Description: "The header to match for.",
																Computed:    true,
															},
															"allowed_values": schema.SetAttribute{
																Description: "A list of regular expressions to match for allowed header values. Headers that do not match ANY of these values will be filtered and will not reach the workload.",
																ElementType: types.StringType,
																Computed:    true,
															},
															"blocked_values": schema.SetAttribute{
																Description: "A list of regular expressions to match for blocked header values. Headers that match ANY of these values will be filtered and will not reach the workload.",
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
						"internal": schema.ListNestedBlock{
							Description: "The internal firewall is used to control access between workloads.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"inbound_allow_type": schema.StringAttribute{
										Description: "Used to control the internal firewall configuration and mutual tls. Allowed Values: \"none\", \"same-gvc\", \"same-org\", \"workload-list\".",
										Computed:    true,
									},
									"inbound_allow_workload": schema.SetAttribute{
										Description: "A list of specific workloads which are allowed to access this workload internally. This list is only used if the 'inboundAllowType' is set to 'workload-list'.",
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"options":       d.OptionsSchema("Configurable settings or parameters that allow fine-tuning and customization of the behavior, performance, and characteristics of the workload."),
			"local_options": d.LocalOptionsSchema(),
			"job": schema.ListNestedBlock{
				MarkdownDescription: "[Cron Job Reference Page](https://docs.controlplane.com/reference/workload#cron).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"schedule": schema.StringAttribute{
							Description: "A standard cron [schedule expression](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#schedule-syntax) used to determine when your job should execute.",
							Computed:    true,
						},
						"concurrency_policy": schema.StringAttribute{
							Description: "Either 'Forbid', 'Replace', or 'Allow'. This determines what Control Plane will do when the schedule requires a job to start, while a prior instance of the job is still running.",
							Computed:    true,
						},
						"history_limit": schema.Int32Attribute{
							Description: "The maximum number of completed job instances to display. This should be an integer between 1 and 10. Default: `5`.",
							Computed:    true,
						},
						"restart_policy": schema.StringAttribute{
							Description: "Either 'OnFailure' or 'Never'. This determines what Control Plane will do when a job instance fails. Enum: [ OnFailure, Never ] Default: `Never`.",
							Computed:    true,
						},
						"active_deadline_seconds": schema.Int32Attribute{
							Description: "The maximum number of seconds Control Plane will wait for the job to complete. If a job does not succeed or fail in the allotted time, Control Plane will stop the job, moving it into the Removed status.",
							Computed:    true,
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
							Computed:    true,
						},
					},
				},
			},
			"rollout_options": schema.ListNestedBlock{
				Description: "Defines the parameters for updating applications and services, including settings for minimum readiness, unavailable replicas, surge replicas, and scaling policies.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"min_ready_seconds": schema.Int32Attribute{
							Description: "The minimum number of seconds a container must run without crashing to be considered available.",
							Computed:    true,
						},
						"max_unavailable_replicas": schema.StringAttribute{
							Description: "The number of replicas that can be unavailable during the update process.",
							Computed:    true,
						},
						"max_surge_replicas": schema.StringAttribute{
							Description: "The number of replicas that can be created above the desired amount of replicas during an update.",
							Computed:    true,
						},
						"scaling_policy": schema.StringAttribute{
							Description: "The strategies used to update applications and services deployed. Valid values: `OrderedReady` (Updates workloads in a rolling fashion, taking down old ones and bringing up new ones incrementally, ensuring that the service remains available during the update.), `Parallel` (Causes all pods affected by a scaling operation to be created or destroyed simultaneously. This does not affect update operations.). Default: `OrderedReady`.",
							Computed:    true,
						},
						"termination_grace_period_seconds": schema.Int32Attribute{
							Description: "The amount of time in seconds a workload has to gracefully terminate before forcefully terminating it. This includes the time it takes for the preStop hook to run.",
							Computed:    true,
						},
					},
				},
			},
			"security_options": schema.ListNestedBlock{
				Description: "Allows for the configuration of the `file system group id` and `geo location`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"file_system_group_id": schema.Int32Attribute{
							Description: "The group id assigned to any mounted volume.",
							Computed:    true,
						},
					},
				},
			},
			"load_balancer": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"replica_direct": schema.BoolAttribute{
							Description: "When enabled, individual replicas of the workload can be reached directly using the subdomain prefix replica-<index>. For example, replica-0.my-workload.my-gvc.cpln.local or replica-0.my-workload-<gvc-alias>.cpln.app - Can only be used with stateful workloads.",
							Computed:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"direct": schema.ListNestedBlock{
							Description: "Direct load balancers are created in each location that a workload is running in and are configured for the standard endpoints of the workload. Customers are responsible for configuring the workload with certificates if TLS is required.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "When disabled, this load balancer will be stopped.",
										Computed:    true,
									},
									"ipset": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"port": schema.SetNestedBlock{
										Description: "List of ports that will be exposed by this load balancer.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"external_port": schema.Int32Attribute{
													Description: "The port that is available publicly.",
													Computed:    true,
												},
												"protocol": schema.StringAttribute{
													Description: "The protocol that is exposed publicly.",
													Computed:    true,
												},
												"scheme": schema.StringAttribute{
													Description: "Overrides the default `https` url scheme that will be used for links in the UI and status.",
													Computed:    true,
												},
												"container_port": schema.Int32Attribute{
													Description: "The port on the container tha will receive this traffic.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"geo_location": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "When enabled, geo location headers will be included on inbound http requests. Existing headers will be replaced.",
										Computed:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"headers": schema.ListNestedBlock{
										Description: "",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"asn": schema.StringAttribute{
													Description: "The geo asn header.",
													Computed:    true,
												},
												"city": schema.StringAttribute{
													Description: "The geo city header.",
													Computed:    true,
												},
												"country": schema.StringAttribute{
													Description: "The geo country header.",
													Computed:    true,
												},
												"region": schema.StringAttribute{
													Description: "The geo region header.",
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
			"request_retry_policy": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"attempts": schema.Int32Attribute{
							Description: "",
							Computed:    true,
						},
						"retry_on": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read fetches the current state of the resource.
func (d *WorkloadDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state WorkloadResourceModel

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

/*** Schemas ***/

// HealthCheckSchema returns a nested block list schema for configuring workload health checks.
func (d *WorkloadDataSource) HealthCheckSchema(attributeName string, description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"initial_delay_seconds": schema.Int32Attribute{
					Description: "",
					Computed:    true,
				},
				"period_seconds": schema.Int32Attribute{
					Description: "",
					Computed:    true,
				},
				"timeout_seconds": schema.Int32Attribute{
					Description: "",
					Computed:    true,
				},
				"success_threshold": schema.Int32Attribute{
					Description: "",
					Computed:    true,
				},
				"failure_threshold": schema.Int32Attribute{
					Description: "",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"exec": d.ExecSchema(""),
				"grpc": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"port": schema.Int32Attribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
				"tcp_socket": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"port": schema.Int32Attribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
				"http_get": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"port": schema.Int32Attribute{
								Description: "",
								Computed:    true,
							},
							"http_headers": schema.MapAttribute{
								Description: "",
								ElementType: types.StringType,
								Computed:    true,
							},
							"scheme": schema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// ExecSchema returns a nested block list schema for configuring exec-based probes.
func (d *WorkloadDataSource) ExecSchema(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"command": schema.ListAttribute{
					Description: description,
					ElementType: types.StringType,
					Computed:    true,
				},
			},
		},
	}
}

// LifecycleSpecSchema returns a nested block list schema for workload lifecycle specifications.
func (d *WorkloadDataSource) LifecycleSpecSchema(commandDescription string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"exec": d.ExecSchema(commandDescription),
			},
		},
	}
}

// OptionsSchema returns a nested block list schema for workload options such as AI capacity, debug mode, and auto-scaling.
func (d *WorkloadDataSource) OptionsSchema(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"timeout_seconds": schema.Int32Attribute{
					Description: "Timeout in seconds. Default: `5`.",
					Computed:    true,
				},
				"capacity_ai": schema.BoolAttribute{
					Description: "Capacity AI. Default: `true`.",
					Computed:    true,
				},
				"capacity_ai_update_minutes": schema.Int32Attribute{
					Description: "The highest frequency capacity AI is allowed to update resource reservations when CapacityAI is enabled.",
					Computed:    true,
				},
				"debug": schema.BoolAttribute{
					Description: "Debug mode. Default: `false`.",
					Computed:    true,
				},
				"suspend": schema.BoolAttribute{
					Description: "Workload suspend. Default: `false`.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"autoscaling": schema.ListNestedBlock{
					Description: "Auto-scaling adjusts horizontal scaling based on a set strategy, target value, and possibly a metric percentile.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"metric": schema.StringAttribute{
								Description: "Valid values: `concurrency`, `cpu`, `memory`, `rps`, `latency`, `keda` or `disabled`.",
								Computed:    true,
							},
							"metric_percentile": schema.StringAttribute{
								Description: "For metrics represented as a distribution (e.g. latency) a percentile within the distribution must be chosen as the target.",
								Computed:    true,
							},
							"target": schema.Int32Attribute{
								Description: "Control Plane will scale the number of replicas for this deployment up/down in order to be as close as possible to the target metric across all replicas of a deployment. Min: `1`. Max: `20000`. Default: `95`.",
								Computed:    true,
							},
							"min_scale": schema.Int32Attribute{
								Description: "The minimum allowed number of replicas. Control Plane can scale the workload down to 0 when there is no traffic and scale up immediately to fulfill new requests. Min: `0`. Max: `max_scale`. Default `1`.",
								Computed:    true,
							},
							"max_scale": schema.Int32Attribute{
								Description: "The maximum allowed number of replicas. Min: `0`. Default `5`.",
								Computed:    true,
							},
							"scale_to_zero_delay": schema.Int32Attribute{
								Description: "The amount of time (in seconds) with no requests received before a workload is scaled to 0. Min: `30`. Max: `3600`. Default: `300`.",
								Computed:    true,
							},
							"max_concurrency": schema.Int32Attribute{
								Description: "A hard maximum for the number of concurrent requests allowed to a replica. If no replicas are available to fulfill the request then it will be queued until a replica with capacity is available and delivered as soon as one is available again. Capacity can be available from requests completing or when a new replica is available from scale out.Min: `0`. Max: `1000`. Default `0`.",
								Computed:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"multi": schema.ListNestedBlock{
								Description: "",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"metric": schema.StringAttribute{
											Description: "Valid values: `cpu` or `memory`.",
											Computed:    true,
										},
										"target": schema.Int32Attribute{
											Description: "Control Plane will scale the number of replicas for this deployment up/down in order to be as close as possible to the target metric across all replicas of a deployment. Min: `1`. Max: `20000`.",
											Computed:    true,
										},
									},
								},
							},
							"keda": schema.ListNestedBlock{
								Description: "KEDA (Kubernetes-based Event Driven Autoscaling) allows for advanced autoscaling based on external metrics and triggers.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"polling_interval": schema.Int32Attribute{
											Description: "The interval in seconds at which KEDA will poll the external metrics to determine if scaling is required.",
											Computed:    true,
										},
										"cooldown_period": schema.Int32Attribute{
											Description: "The cooldown period in seconds after scaling down to 0 replicas before KEDA will allow scaling up again.",
											Computed:    true,
										},
										"initial_cooldown_period": schema.Int32Attribute{
											Description: "The initial cooldown period in seconds after scaling down to 0 replicas before KEDA will allow scaling up again.",
											Computed:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"trigger": schema.ListNestedBlock{
											Description: "An array of KEDA triggers to be used for scaling workloads in this GVC. This is used to define how KEDA will scale workloads in the GVC based on external metrics or events. Each trigger type may have its own specific configuration options.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"type": schema.StringAttribute{
														Description: `The type of KEDA trigger, e.g "prometheus", "aws-sqs", etc.`,
														Computed:    true,
													},
													"metadata": schema.MapAttribute{
														Description: "The configuration parameters that the trigger requires.",
														ElementType: types.StringType,
														Computed:    true,
													},
													"name": schema.StringAttribute{
														Description: "An optional name for the trigger. If not provided, a default name will be generated based on the trigger type.",
														Computed:    true,
													},
													"use_cached_metrics": schema.BoolAttribute{
														Description: "Enables caching of metric values during polling interval.",
														Computed:    true,
													},
													"metric_type": schema.StringAttribute{
														Description: "The type of metric to be used for scaling.",
														Computed:    true,
													},
												},
												Blocks: map[string]schema.Block{
													"authentication_ref": schema.ListNestedBlock{
														Description: "Reference to a KEDA authentication object for secure access to external systems.",
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"name": schema.StringAttribute{
																	Description: "The name of secret listed in the GVC spec.keda.secrets.",
																	Computed:    true,
																},
															},
														},
													},
												},
											},
										},
										"advanced": schema.ListNestedBlock{
											Description: "Advanced configuration options for KEDA.",
											NestedObject: schema.NestedBlockObject{
												Blocks: map[string]schema.Block{
													"scaling_modifiers": schema.ListNestedBlock{
														Description: "Scaling modifiers allow for fine-tuning the scaling behavior of KEDA.",
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"target": schema.StringAttribute{
																	Description: "Defines new target value to scale on for the composed metric.",
																	Computed:    true,
																},
																"activation_target": schema.StringAttribute{
																	Description: "Defines the new activation target value to scale on for the composed metric.",
																	Computed:    true,
																},
																"metric_type": schema.StringAttribute{
																	Description: "Defines metric type used for this new composite-metric.",
																	Computed:    true,
																},
																"formula": schema.StringAttribute{
																	Description: "Composes metrics together and allows them to be modified/manipulated. It accepts mathematical/conditional statements.",
																	Computed:    true,
																},
															},
														},
													},
												},
											},
										},
										"fallback": schema.ListNestedBlock{
											Description: "Fallback configuration for KEDA.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"failure_threshold": schema.Int32Attribute{
														Description: "Number of consecutive failures required to trigger fallback behavior.",
														Computed:    true,
													},
													"replicas": schema.Int32Attribute{
														Description: "Number of replicas to scale to when fallback is triggered.",
														Computed:    true,
													},
													"behavior": schema.StringAttribute{
														Description: "Behavior to apply when fallback is triggered.",
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
				"multi_zone": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// LocalOptionsSchema returns a nested block list schema for overriding default options per specific Control Plane location
func (d *WorkloadDataSource) LocalOptionsSchema() schema.ListNestedBlock {
	// Build base options schema with an override-focused description
	options := d.OptionsSchema("Override defaultOptions for the workload in specific Control Plane Locations.")

	// Define local-options-specific attributes to be merged into the nested object attributes
	localOptions := map[string]schema.Attribute{
		"location": schema.StringAttribute{
			Description: "Valid only for `local_options`. Override options for a specific location.",
			Required:    true,
		},
	}

	// Create a fresh map with capacity for existing attributes plus local overrides
	merged := make(map[string]schema.Attribute, len(options.NestedObject.Attributes)+len(localOptions))

	// Copy all existing attributes from the base options into the merged map
	for k, v := range options.NestedObject.Attributes {
		merged[k] = v
	}

	// Overlay local-options attributes so they add new keys or override duplicates by design
	for k, v := range localOptions {
		merged[k] = v
	}

	// Assign the merged attribute map back to the nested object schema
	options.NestedObject.Attributes = merged

	// Clear list validators so multiple local_options blocks are permitted
	options.Validators = []validator.List{}

	// Return the fully configured local options schema
	return options
}
