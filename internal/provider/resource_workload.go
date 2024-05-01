package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Main ***/
func resourceWorkload() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceWorkloadCreate,
		ReadContext:   resourceWorkloadRead,
		UpdateContext: resourceWorkloadUpdate,
		DeleteContext: resourceWorkloadDelete,
		Schema: map[string]*schema.Schema{
			"gvc": {
				Type:         schema.TypeString,
				Description:  "Name of the associated GVC.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the Workload.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the Workload.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the Workload.",
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:        schema.TypeMap,
				Description: "Key-value map of resource tags.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"identity_link": {
				Type:         schema.TypeString,
				Description:  "Full link to an Identity.",
				Optional:     true,
				ValidateFunc: LinkValidator,
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "Workload Type. Either `serverless`, `standard`, `stateful`, or `cron`.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: WorkloadTypeValidator,
			},
			"container": {
				Type:     schema.TypeList,
				Description: "An isolated and lightweight runtime environment that encapsulates an application and its dependencies.",
				Required: true,
				MinItems: 1,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the container.",
							Required:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								warns, errs = NameValidator(val, key)

								v := val.(string)

								if strings.HasPrefix(v, "cpln-") {
									errs = append(errs, fmt.Errorf("%q cannot start with 'cpln-', got: %s", key, v))
									return
								}

								if v == "istio-proxy" || v == "queue-proxy" || v == "istio-validation" || v == "cpln-envoy-assassin" || v == "cpln-writer-proxy" || v == "cpln-reader-proxy" || v == "cpln-dbaas-config" {
									errs = append(errs, fmt.Errorf("%q cannot be set to 'istio-proxy', 'queue-proxy', 'istio-validation', 'cpln-envoy-assassin', 'cpln-writer-proxy', 'cpln-reader-proxy', 'cpln-dbaas-config', got: %s", key, v))
								}

								return
							},
						},
						"image": {
							Type:        schema.TypeString,
							Description: "The full image and tag path.",
							Required:    true,
						},
						"port": {
							Type:         schema.TypeInt,
							Description:  "The port the container exposes. Only one container is allowed to specify a port. Min: `80`. Max: `65535`. Used by `serverless` Workload type. **DEPRECATED - Use `ports`.**",
							Optional:     true,
							ValidateFunc: PortValidator,
							Deprecated:   "The 'port' attribute will be deprecated in the next major version. Use the 'ports' attribute instead.",
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Description:  "Protocol. Choice of: `http`, `http2`, `tcp`, or `grpc`.",
										Optional:     true,
										ValidateFunc: PortProtocolValidator,
										Default:      "http",
									},
									"number": {
										Type:        schema.TypeInt,
										Description: "Port to expose.",
										Required:    true,
									},
								},
							},
						},
						"cpu": {
							Type:         schema.TypeString,
							Description:  "Reserved CPU of the workload when capacityAI is disabled. Maximum CPU when CapacityAI is enabled. Default: \"50m\".",
							Optional:     true,
							Default:      "50m",
							ValidateFunc: CpuMemoryValidator,
						},
						"gpu_nvidia": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:        schema.TypeString,
										Description: "GPU Model (i.e.: t4)",
										Required:    true,
									},
									"quantity": {
										Type:        schema.TypeInt,
										Description: "Number of GPUs.",
										Required:    true,
									},
								},
							},
						},
						"memory": {
							Type:         schema.TypeString,
							Description:  "Reserved memory of the workload when capacityAI is disabled. Maximum memory when CapacityAI is enabled. Default: \"128Mi\".",
							Optional:     true,
							Default:      "128Mi",
							ValidateFunc: CpuMemoryValidator,
						},
						"min_cpu": {
							Type:         schema.TypeString,
							Description:  "Minimum CPU when capacity AI is enabled.",
							Optional:     true,
							ValidateFunc: CpuMemoryValidator,
						},
						"min_memory": {
							Type:         schema.TypeString,
							Description:  "Minimum memory when capacity AI is enabled.",
							Optional:     true,
							ValidateFunc: CpuMemoryValidator,
						},
						"working_directory": {
							Type:        schema.TypeString,
							Description: "Override the working directory. Must be an absolute path.",
							Optional:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								v := val.(string)
								v = path.Clean(v)

								if !path.IsAbs(v) {
									errs = append(errs, fmt.Errorf("%q must be an absolute path, got: %s", key, v))
								}

								return
							},
						},
						"command": {
							Type:        schema.TypeString,
							Description: "Override the entry point.",
							Optional:    true,
						},
						"env": {
							Type:        schema.TypeMap,
							Description: "Name-Value list of environment variables.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								v := val.(map[string]interface{})

								for name, value := range v {

									nameLower := strings.ToLower(name)

									if nameLower == "k_service" || nameLower == "k_configuration" || nameLower == "k_revision" {
										errs = append(errs, fmt.Errorf("%q cannot be 'K_SERVICE', 'K_CONFIGURATION', 'K_REVISION', got: %s", key, nameLower))
									}

									maxValueLength := 4 * 1024

									if len(value.(string)) > maxValueLength {
										errs = append(errs, fmt.Errorf("%q length cannot be > %d, got: %d", key, maxValueLength, len(value.(string))))
									}
								}

								return
							},
						},
						"inherit_env": {
							Type:        schema.TypeBool,
							Description: "Enables inheritance of GVC environment variables. A variable in spec.env will override a GVC variable with the same name.",
							Optional:    true,
							Default:     false,
						},
						"args": {
							Type:        schema.TypeList,
							Description: "Command line arguments passed to the container at runtime.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"liveness_probe": {
							Type:        schema.TypeList,
							Description: "Liveness Probe",
							Optional:    true,
							MaxItems:    1,
							Elem:        healthCheckSpec(),
						},
						"readiness_probe": {
							Type:        schema.TypeList,
							Description: "Readiness Probe",
							Optional:    true,
							MaxItems:    1,
							Elem:        healthCheckSpec(),
						},
						"volume": {
							Type:        schema.TypeList,
							Description: "[Reference Page](https://docs.controlplane.com/reference/workload#volumes).",
							Optional:    true,
							MaxItems:    5,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uri": {
										Type:        schema.TypeString,
										Description: "URI of a volume hosted at Control Plane (Volume Set) or at a cloud provider (AWS, Azure, GCP).",
										Required:    true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

											v := val.(string)

											re := regexp.MustCompile(`^(s3|gs|azureblob|azurefs|cpln|scratch):\/\/.+`)

											if !re.MatchString(v) {
												errs = append(errs, fmt.Errorf("%q must be in the form s3://bucket, gs://bucket, azureblob://storageAccount/container, azurefs://storageAccount/share, cpln://, scratch://, got: %s", key, v))
											}

											return
										},
									},
									"recovery_policy": {
										Type:        schema.TypeString,
										Description: "Only applicable to persistent volumes, this determines what Control Plane will do when creating a new workload replica if a corresponding volume exists. Available Values: `retain`, `recycle`. Default: `retain`. **DEPRECATED - No longer being used.**",
										Optional:    true,
										Default:     "retain",
									},
									"path": {
										Type:        schema.TypeString,
										Description: "File path added to workload pointing to the volume.",
										Required:    true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

											v := val.(string)

											if !path.IsAbs(v) {
												errs = append(errs, fmt.Errorf("%q must be an absolute path, got: %s", key, v))
												return
											}

											v = path.Clean(v)
											v = strings.TrimRight(v, "/")

											if v == "/dev/log" || v == "/dev" || v == "/tmp" || v == "/var" || v == "/var/log" {
												errs = append(errs, fmt.Errorf("%q is set to a reserved path, got: %s", key, v))
											}

											return
										},
									},
								},
							},
						},
						"metrics": {
							Type:        schema.TypeList,
							Description: "[Reference Page](https://docs.controlplane.com/reference/workload#metrics).",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:         schema.TypeInt,
										Description:  "Port from container emitting custom metrics",
										Required:     true,
										ValidateFunc: PortValidator,
									},
									"path": {
										Type:        schema.TypeString,
										Description: "Path from container emitting custom metrics",
										Required:    true,
									},
								},
							},
						},
						"lifecycle": {
							Type:        schema.TypeList,
							Description: "Lifecycle [Reference Page](https://docs.controlplane.com/reference/workload#lifecycle).",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"post_start": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem:     lifeCycleSpec(),
									},
									"pre_stop": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem:     lifeCycleSpec(),
									},
								},
							},
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity_ai": {
							Type:        schema.TypeBool,
							Description: "Capacity AI. Default: `true`.",
							Optional:    true,
							Default:     true,
						},
						"debug": {
							Type:        schema.TypeBool,
							Description: "Debug mode. Default: `false`",
							Optional:    true,
							Default:     false,
						},
						"suspend": {
							Type:        schema.TypeBool,
							Description: "Workload suspend. Default: `false`",
							Optional:    true,
							Default:     false,
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Description: "Timeout in seconds. Default: `5`.",
							Optional:    true,
							Default:     5,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 3600 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 3600 inclusive, got: %d", key, v))
								}
								return
							},
						},
						"autoscaling": {
							Type: schema.TypeList,
							// Required: true,
							Optional: true,
							MaxItems: 1,
							Elem:     AutoScalingResource(),
						},
					},
				},
			},
			"local_options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeString,
							Description: "Valid only for `local_options`. Override options for a specific location.",
							Required:    true,
						},
						"capacity_ai": {
							Type:        schema.TypeBool,
							Description: "Capacity AI. Default: `true`.",
							Optional:    true,
							Default:     true,
						},
						"debug": {
							Type:        schema.TypeBool,
							Description: "Debug mode. Default: `false`",
							Optional:    true,
							Default:     false,
						},
						"suspend": {
							Type:        schema.TypeBool,
							Description: "Workload suspend. Default: `false`",
							Optional:    true,
							Default:     false,
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Description: "Timeout in seconds. Default: `5`.",
							Optional:    true,
							Default:     5,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 3600 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 3600 inclusive, got: %d", key, v))
								}
								return
							},
						},
						"autoscaling": {
							Type: schema.TypeList,
							// Required: true,
							Optional: true,
							MaxItems: 1,
							Elem:     AutoScalingResource(),
						},
					},
				},
			},
			"firewall_spec": {
				Type:        schema.TypeList,
				Description: "Control of inbound and outbound access to the workload for external (public) and internal (service to service) traffic. Access is restricted by default.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     ExternalFirewallResource(),
						},
						"internal": {
							Type:        schema.TypeList,
							Description: "The internal firewall is used to control access between workloads.",
							Optional:    true,
							MaxItems:    1,
							Elem:        InternalFirewallResource(),
						},
					},
				},
			},
			"job": {
				Type:        schema.TypeList,
				Description: "[Cron Job Reference Page](https://docs.controlplane.com/reference/workload#cron).",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule": {
							Type:        schema.TypeString,
							Description: "A standard cron [schedule expression](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#schedule-syntax) used to determine when your job should execute.",
							Required:    true,
						},
						"concurrency_policy": {
							Type:        schema.TypeString,
							Description: "Either 'Forbid' or 'Replace'. This determines what Control Plane will do when the schedule requires a job to start, while a prior instance of the job is still running. Enum: [ Forbid, Replace ] Default: `Forbid`",
							Optional:    true,
							Default:     "Forbid",
						},
						"history_limit": {
							Type:        schema.TypeInt,
							Description: "The maximum number of completed job instances to display. This should be an integer between 1 and 10. Default: `5`",
							Optional:    true,
							Default:     5,
						},
						"restart_policy": {
							Type:        schema.TypeString,
							Description: "Either 'OnFailure' or 'Never'. This determines what Control Plane will do when a job instance fails. Enum: [ OnFailure, Never ] Default: `Never`",
							Optional:    true,
							Default:     "Never",
						},
						"active_deadline_seconds": {
							Type:        schema.TypeInt,
							Description: "The maximum number of seconds Control Plane will wait for the job to complete. If a job does not succeed or fail in the allotted time, Control Plane will stop the job, moving it into the Removed status.",
							Optional:    true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeList,
				Description: "Status of the workload.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_id": {
							Type:        schema.TypeString,
							Description: "ID of the parent object.",
							Optional:    true,
						},
						"canonical_endpoint": {
							Type:        schema.TypeString,
							Description: "Canonical endpoint for the workload.",
							Optional:    true,
						},
						"endpoint": {
							Type:        schema.TypeString,
							Description: "Endpoint for the workload.",
							Optional:    true,
						},
						"internal_name": {
							Type:        schema.TypeString,
							Description: "Internal hostname for the workload. Used for service-to-service requests.",
							Optional:    true,
						},
						"health_check": {
							Type:        schema.TypeList,
							Description: "Current health status.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"active": {
										Type:        schema.TypeBool,
										Description: "Active boolean for the associated workload.",
										Required:    true,
									},
									"success": {
										Type:        schema.TypeBool,
										Description: "Success boolean for the associated workload.",
										Optional:    true,
									},
									"code": {
										Type:        schema.TypeInt,
										Description: "Current output code for the associated workload.",
										Optional:    true,
									},
									"message": {
										Type:        schema.TypeString,
										Description: "Current health status for the associated workload.",
										Optional:    true,
									},
									"failures": {
										Type:        schema.TypeInt,
										Description: "Failure integer for the associated workload.",
										Optional:    true,
									},
									"successes": {
										Type:        schema.TypeInt,
										Description: "Success integer for the associated workload.",
										Optional:    true,
									},
									"last_checked": {
										Type:        schema.TypeString,
										Description: "Timestamp in UTC of the last health check.",
										Optional:    true,
									},
								},
							},
						},
						"current_replica_count": {
							Type:        schema.TypeInt,
							Description: "Current amount of replicas deployed.",
							Optional:    true,
						},
						"resolved_images": {
							Type:        schema.TypeList,
							Description: "Resolved images for workloads with dynamic tags enabled.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resolved_for_version": {
										Type:        schema.TypeInt,
										Description: "Workload version the images were resolved for.",
										Optional:    true,
									},
									"resolved_at": {
										Type:        schema.TypeString,
										Description: "UTC Time when the images were resolved.",
										Optional:    true,
									},
									"images": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"digest": {
													Type:        schema.TypeString,
													Description: "A unique SHA256 hash value that identifies a specific image content. This digest serves as a fingerprint of the image's content, ensuring the image you pull or run is exactly what you expect, without any modifications or corruptions.",
													Optional:    true,
												},
												"manifests": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"image": {
																Type:        schema.TypeString,
																Description: "The name and tag of the resolved image.",
																Optional:    true,
															},
															"media_type": {
																Type:        schema.TypeString,
																Description: "The MIME type used in the Docker Registry HTTP API to specify the format of the data being sent or received. Docker uses media types to distinguish between different kinds of JSON objects and binary data formats within the registry protocol, enabling the Docker client and registry to understand and process different components of Docker images correctly.",
																Optional:    true,
															},
															"digest": {
																Type:        schema.TypeString,
																Description: "A SHA256 hash that uniquely identifies the specific image manifest.",
																Optional:    true,
															},
															"platform": {
																Type:        schema.TypeMap,
																Description: "Key-value map of strings. The combination of the operating system and architecture for which the image is built.",
																Optional:    true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
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
						},
					},
				},
			},
			"rollout_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_ready_seconds": {
							Type:        schema.TypeInt,
							Description: "The minimum number of seconds a container must run without crashing to be considered available",
							Optional:    true,
							Default:     0,
						},
						"max_unavailable_replicas": {
							Type:        schema.TypeString,
							Description: "The number of replicas that can be unavailable during the update process.",
							Optional:    true,
						},
						"max_surge_replicas": {
							Type:        schema.TypeString,
							Description: "The number of replicas that can be created above the desired amount of replicas during an update.",
							Optional:    true,
						},
						"scaling_policy": {
							Type:        schema.TypeString,
							Description: "The strategies used to update applications and services deployed. Valid values: `OrderedReady` (Updates workloads in a rolling fashion, taking down old ones and bringing up new ones incrementally, ensuring that the service remains available during the update.), `Parallel` (Causes all pods affected by a scaling operation to be created or destroyed simultaneously. This does not affect update operations.). Default: `OrderedReady`.",
							Optional:    true,
							Default:     "OrderedReady",
						},
					},
				},
			},
			"security_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_system_group_id": {
							Type:        schema.TypeInt,
							Description: "The group id assigned to any mounted volume.",
							Required:    true,
						},
					},
				},
			},
			"support_dynamic_tags": {
				Type:        schema.TypeBool,
				Description: "Workload will automatically redeploy when one of the container images is updated in the container registry. Default: false.",
				Optional:    true,
				Default:     false,
			},
			"sidecar": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"envoy": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateWorkload,
		},
	}
}

func importStateWorkload(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.SplitN(d.Id(), ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected ID syntax: 'gvc:workload'. Example: 'terraform import cpln_workload.RESOURCE_NAME GVC_NAME:WORKLOAD_NAME'", d.Id())
	}

	d.Set("gvc", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceWorkloadCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceWorkloadCreate")

	if checkLegacyPort(d.Get("container").([]interface{})) {
		var diags diag.Diagnostics

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "port and ports are both defined",
			Detail:   "Use only the ports attributes",
		})

		return diags
	}

	c := m.(*client.Client)

	gvcName := d.Get("gvc").(string)

	workload := client.Workload{}
	workload.Spec = &client.WorkloadSpec{}
	workload.Name = GetString(d.Get("name"))
	workload.Description = GetString(d.Get("description"))
	workload.Tags = GetStringMap(d.Get("tags"))

	legacyPort := buildContainers(d.Get("container").([]interface{}), workload.Spec)
	buildFirewallSpec(d.Get("firewall_spec").([]interface{}), workload.Spec)
	buildOptions(d.Get("options").([]interface{}), workload.Spec, false, c.Org)
	buildOptions(d.Get("local_options").([]interface{}), workload.Spec, true, c.Org)
	workload.Spec.Job = buildJobSpec(d.Get("job").([]interface{}))

	workload.Spec.Type = GetString(strings.TrimSpace(d.Get("type").(string)))

	if d.Get("identity_link") != nil {

		identityLink := strings.TrimSpace(d.Get("identity_link").(string))

		if identityLink != "" {

			workload.Spec.IdentityLink = GetString(identityLink)
		}
	}

	if d.Get("rollout_options") != nil {
		rolloutOptions := buildRolloutOptions(d.Get("rollout_options").([]interface{}))

		workload.Spec.RolloutOptions = rolloutOptions
	}

	if d.Get("security_options") != nil {
		securityOptions := buildSecurityOptions(d.Get("security_options").([]interface{}))

		workload.Spec.SecurityOptions = securityOptions
	}

	if d.Get("support_dynamic_tags") != nil {
		workload.Spec.SupportDynamicTags = GetBool(d.Get("support_dynamic_tags"))
	}

	if d.Get("sidecar") != nil {
		workload.Spec.Sidecar = buildWorkloadSidecar(d.Get("sidecar").([]interface{}))
	}

	if e := workloadSpecValidate(workload.Spec); e != nil {
		return e
	}

	newWorkload, code, err := c.CreateWorkload(workload, gvcName)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setWorkload(d, newWorkload, c.Org, legacyPort, nil)
}

func resourceWorkloadRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceWorkloadRead")

	workloadName := d.Id()
	gvcName := d.Get("gvc").(string)

	workloadTemp := client.Workload{}
	workloadTemp.Spec = &client.WorkloadSpec{}
	legacyPort := buildContainers(d.Get("container").([]interface{}), workloadTemp.Spec)

	c := m.(*client.Client)
	workload, code, err := c.GetWorkload(workloadName, gvcName)

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if tfTest := os.Getenv("CPLN_TF_TEST"); tfTest == "true" {

		if workload.Status == nil {
			workload.Status = &client.WorkloadStatus{}
		}

		testEndpoint := "http://tf-test"

		workload.Status.Endpoint = &testEndpoint
		workload.Status.CanonicalEndpoint = &testEndpoint
	}

	var diags diag.Diagnostics
	count := 0

	for workload.Status == nil || workload.Status.Endpoint == nil || workload.Status.CanonicalEndpoint == nil || strings.TrimSpace(*workload.Status.Endpoint) == "" || strings.TrimSpace(*workload.Status.CanonicalEndpoint) == "" {

		if count++; count > 8 {
			// Exit loop after 120 seconds

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to obtain current status",
				Detail:   "Workload status is not available. Run 'terraform apply' again.",
			})

			break
		}

		// log.Printf("Waiting For Valid Status. Count: %d", count)

		time.Sleep(15 * time.Second)

		workload, _, err = c.GetWorkload(workloadName, gvcName)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	// log.Printf("Before Calling SET: Endpoint: %s. Canonical: %s", workload.Status.Endpoint, workload.Status.CanonicalEndpoint)

	return setWorkload(d, workload, c.Org, legacyPort, diags)
}

func resourceWorkloadUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceWorkloadUpdate")

	if d.HasChanges("description", "tags", "type", "container", "options", "local_options", "firewall_spec", "job", "identity_link", "rollout_options", "security_options", "support_dynamic_tags", "sidecar") {

		if checkLegacyPort(d.Get("container").([]interface{})) {
			var diags diag.Diagnostics

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "port and ports are both defined",
				Detail:   "Use only the ports attributes",
			})

			return diags
		}

		c := m.(*client.Client)

		gvcName := d.Get("gvc").(string)

		workloadToUpdate := client.Workload{}
		workloadToUpdate.Name = GetString(d.Get("name"))
		workloadToUpdate.Description = GetDescriptionString(d.Get("description"), *workloadToUpdate.Name)
		workloadToUpdate.Tags = GetTagChanges(d)

		workloadToUpdate.SpecReplace = &client.WorkloadSpec{}
		workloadToUpdate.SpecReplace.Type = GetString(d.Get("type"))

		legacyPort := buildContainers(d.Get("container").([]interface{}), workloadToUpdate.SpecReplace)
		buildOptions(d.Get("options").([]interface{}), workloadToUpdate.SpecReplace, false, c.Org)
		buildOptions(d.Get("local_options").([]interface{}), workloadToUpdate.SpecReplace, true, c.Org)
		buildFirewallSpec(d.Get("firewall_spec").([]interface{}), workloadToUpdate.SpecReplace)
		workloadToUpdate.SpecReplace.Job = buildJobSpec(d.Get("job").([]interface{}))
		workloadToUpdate.SpecReplace.RolloutOptions = buildRolloutOptions(d.Get("rollout_options").([]interface{}))
		workloadToUpdate.SpecReplace.SecurityOptions = buildSecurityOptions(d.Get("security_options").([]interface{}))
		workloadToUpdate.SpecReplace.SupportDynamicTags = GetBool(d.Get("support_dynamic_tags"))
		workloadToUpdate.SpecReplace.Sidecar = buildWorkloadSidecar(d.Get("sidecar").([]interface{}))

		if d.Get("identity_link") != nil {

			if identityLink := strings.TrimSpace(d.Get("identity_link").(string)); identityLink != "" {
				workloadToUpdate.SpecReplace.IdentityLink = GetString(identityLink)
			}
		}

		if e := workloadSpecValidate(workloadToUpdate.SpecReplace); e != nil {
			return e
		}

		updatedWorkload, _, err := c.UpdateWorkload(workloadToUpdate, gvcName)
		if err != nil {
			return diag.FromErr(err)
		}

		if tfTest := os.Getenv("CPLN_TF_TEST"); tfTest == "true" {

			if updatedWorkload.Status == nil {
				updatedWorkload.Status = &client.WorkloadStatus{}
			}

			testEndpoint := "http://tf-test"

			updatedWorkload.Status.Endpoint = &testEndpoint
			updatedWorkload.Status.CanonicalEndpoint = &testEndpoint
		}

		return setWorkload(d, updatedWorkload, c.Org, legacyPort, nil)
	}

	return nil
}

func resourceWorkloadDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceWorkloadDelete")

	c := m.(*client.Client)
	err := c.DeleteWorkload(d.Id(), d.Get("gvc").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setWorkload(d *schema.ResourceData, workload *client.Workload, org string, legacyPort bool, diags diag.Diagnostics) diag.Diagnostics {

	if workload == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*workload.Name)

	if err := SetBase(d, workload.Base); err != nil {
		return diag.FromErr(err)
	}

	if workload.Spec != nil {
		if err := d.Set("container", flattenContainer(workload.Spec.Containers, legacyPort)); err != nil {
			return diag.FromErr(err)
		}

		if workload.Spec.DefaultOptions != nil {
			if err := d.Set("options", flattenOptions([]client.Options{*workload.Spec.DefaultOptions}, false, org)); err != nil {
				return diag.FromErr(err)
			}
		}

		if workload.Spec.LocalOptions != nil {
			if err := d.Set("local_options", flattenOptions(*workload.Spec.LocalOptions, true, org)); err != nil {
				return diag.FromErr(err)
			}
		}

		if err := d.Set("firewall_spec", flattenFirewallSpec(workload.Spec.FirewallConfig)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("job", flattenJobSpec(workload.Spec.Job)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("identity_link", workload.Spec.IdentityLink); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("type", workload.Spec.Type); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("rollout_options", flattenRolloutOptions(workload.Spec.RolloutOptions)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("security_options", flattenSecurityOptions(workload.Spec.SecurityOptions)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("support_dynamic_tags", workload.Spec.SupportDynamicTags); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("sidecar", flattenWorkloadSidecar(workload.Spec.Sidecar)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("status", flattenWorkloadStatus(workload.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(workload.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func checkLegacyPort(containers []interface{}) bool {

	if containers == nil {
		return false
	}

	for _, container := range containers {

		c := container.(map[string]interface{})

		if (c["port"] != nil && c["port"].(int) > 0) && (c["ports"] != nil && len(c["ports"].([]interface{})) > 0) {
			return true
		}

	}

	return false
}

func buildContainers(containers []interface{}, workloadSpec *client.WorkloadSpec) bool {

	output := false

	if containers == nil {
		return output
	}

	newContainers := []client.ContainerSpec{}

	for _, container := range containers {

		c := container.(map[string]interface{})

		newContainer := client.ContainerSpec{
			Name:             GetString(c["name"].(string)),
			Image:            GetString(c["image"].(string)),
			Memory:           GetString(c["memory"].(string)),
			CPU:              GetString(c["cpu"].(string)),
			Command:          GetString(c["command"].(string)),
			InheritEnv:       GetBool(c["inherit_env"].(bool)),
			WorkingDirectory: GetString(c["working_directory"].(string)),
		}

		if c["gpu_nvidia"] != nil {
			newContainer.GPU = buildGpuNvidia(c["gpu_nvidia"].([]interface{}))
		}

		if c["min_cpu"] != nil {
			newContainer.MinCPU = GetString(c["min_cpu"].(string))
		}

		if c["min_memory"] != nil {
			newContainer.MinMemory = GetString(c["min_memory"].(string))
		}

		if c["port"] != nil && c["port"].(int) > 0 {
			// newContainer.Port = GetPortInt(c["port"])

			newPorts := map[string]interface{}{
				"protocol": "http",
				"number":   c["port"],
			}

			newContainer.Ports = buildPortSpec([]interface{}{newPorts})

			output = true
		}

		if c["ports"] != nil && len(c["ports"].([]interface{})) > 0 {
			newContainer.Ports = buildPortSpec(c["ports"].([]interface{}))
		}

		argArray := []string{}

		for _, value := range c["args"].([]interface{}) {
			if value != nil {
				argArray = append(argArray, value.(string))
			}
		}

		if len(argArray) > 0 {
			newContainer.Args = &argArray
		}

		envArray := []client.NameValue{}

		keys, m := MapSortHelper(c["env"])

		for _, k := range keys {

			name := k
			value := m[k].(string)

			newEnv := client.NameValue{
				Name:  &name,
				Value: &value,
			}

			envArray = append(envArray, newEnv)
		}

		if len(envArray) > 0 {
			newContainer.Env = &envArray
		}

		if c["readiness_probe"] != nil {
			newContainer.ReadinessProbe = buildHealthCheckSpec(c["readiness_probe"].([]interface{}))
		}

		if c["liveness_probe"] != nil {
			newContainer.LivenessProbe = buildHealthCheckSpec(c["liveness_probe"].([]interface{}))
		}

		if c["volume"] != nil {
			newContainer.Volumes = buildVolumeSpec(c["volume"].([]interface{}))
		}

		if c["metrics"] != nil {
			newContainer.Metrics = buildMetrics(c["metrics"].([]interface{}))
		}

		if c["lifecycle"] != nil {
			buildLifeCycleSpec(c["lifecycle"].([]interface{}), &newContainer)
		}

		newContainers = append(newContainers, newContainer)
	}

	workloadSpec.Containers = &newContainers

	return output
}

func buildPortSpec(ports []interface{}) *[]client.PortSpec {

	if len(ports) > 0 {
		output := []client.PortSpec{}

		for _, value := range ports {

			v := value.(map[string]interface{})

			protocol := v["protocol"].(string)
			number := v["number"].(int)

			localPort := client.PortSpec{
				Protocol: &protocol,
				Number:   &number,
			}

			output = append(output, localPort)
		}

		return &output
	}

	return nil
}

func buildGpuNvidia(specs []interface{}) *client.GpuResource {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})

	gpuResource := client.GpuResource{
		Nvidia: &client.Nvidia{
			Model:    GetString(spec["model"].(string)),
			Quantity: GetInt(spec["quantity"].(int)),
		},
	}

	return &gpuResource
}

func buildVolumeSpec(volumes []interface{}) *[]client.VolumeSpec {

	if len(volumes) > 0 {
		output := []client.VolumeSpec{}

		for _, value := range volumes {

			v := value.(map[string]interface{})

			uri := v["uri"].(string)
			recoveryPolicy := v["recovery_policy"].(string)
			path := v["path"].(string)

			localVolume := client.VolumeSpec{
				Uri:            &uri,
				RecoveryPolicy: &recoveryPolicy,
				Path:           &path,
			}

			output = append(output, localVolume)
		}

		return &output
	}

	return nil
}

func buildMetrics(metrics []interface{}) *client.Metrics {

	if len(metrics) == 1 {

		v := metrics[0].(map[string]interface{})

		path := v["path"].(string)
		port := v["port"].(int)

		localMetric := client.Metrics{
			Path: &path,
			Port: &port,
		}

		return &localMetric
	}

	return nil
}

func buildHealthCheckSpec(healthCheck []interface{}) *client.HealthCheckSpec {

	if len(healthCheck) > 0 {

		output := client.HealthCheckSpec{}

		hc := healthCheck[0].(map[string]interface{})

		initDelaySeconds := hc["initial_delay_seconds"].(int)
		periodSeconds := hc["period_seconds"].(int)
		timeoutSeconds := hc["timeout_seconds"].(int)
		successThreshold := hc["success_threshold"].(int)
		failureThreshold := hc["failure_threshold"].(int)

		output.InitialDelaySeconds = &initDelaySeconds
		output.PeriodSeconds = &periodSeconds
		output.TimeoutSeconds = &timeoutSeconds
		output.SuccessThreshold = &successThreshold
		output.FailureThreshold = &failureThreshold

		if hc["exec"] != nil {

			exec := hc["exec"].([]interface{})

			if len(exec) > 0 && exec[0] != nil {
				e := exec[0].(map[string]interface{})
				commands := []string{}

				for _, k := range e["command"].([]interface{}) {
					if k != nil {
						commands = append(commands, k.(string))
					} else {
						commands = append(commands, "")
					}
				}

				if len(commands) > 0 {
					output.Exec = &client.Exec{}
					output.Exec.Command = &commands
				}
			}
		}

		if hc["grpc"] != nil {
			grpc := hc["grpc"].([]interface{})

			if len(grpc) > 0 {
				output.GRPC = &client.GRPC{}

				if grpc[0] != nil {
					output.GRPC.Port = GetPortInt(grpc[0].(map[string]interface{})["port"].(int))
				}
			}
		}

		if hc["tcp_socket"] != nil {
			tcp := hc["tcp_socket"].([]interface{})

			if len(tcp) > 0 {

				output.TCPSocket = &client.TCPSocket{}

				if tcp[0] != nil {
					t := tcp[0].(map[string]interface{})
					port := t["port"].(int)
					output.TCPSocket.Port = GetPortInt(port)
				}
			}
		}

		if hc["http_get"] != nil {

			http := hc["http_get"].([]interface{})

			if len(http) > 0 {

				output.HTTPGet = &client.HTTPGet{}

				h := http[0].(map[string]interface{})

				path := h["path"].(string)
				port := h["port"].(int)
				scheme := h["scheme"].(string)

				output.HTTPGet.Path = &path
				output.HTTPGet.Port = GetPortInt(port)
				output.HTTPGet.Scheme = &scheme

				keys, m := MapSortHelper(h["http_headers"])

				httpHeaders := []client.NameValue{}

				for _, k := range keys {

					name := k
					value := m[k].(string)

					newHeader := client.NameValue{
						Name:  &name,
						Value: &value,
					}

					httpHeaders = append(httpHeaders, newHeader)
				}

				if len(httpHeaders) > 0 {
					output.HTTPGet.HTTPHeaders = &httpHeaders
				}

			}
		}

		return &output
	}

	return nil
}

func buildLifeCycleSpec(lifecycle []interface{}, containerSpec *client.ContainerSpec) {

	if len(lifecycle) > 0 {

		containerSpec.LifeCycle = &client.LifeCycleSpec{}

		if lifecycle[0] != nil {

			lc := lifecycle[0].(map[string]interface{})

			// Set struct fields
			if lc["post_start"] != nil {

				ps := lc["post_start"].([]interface{})

				if len(ps) > 0 && ps[0] != nil {

					psMap := ps[0].(map[string]interface{})
					exec := psMap["exec"].([]interface{})

					if len(exec) > 0 {
						containerSpec.LifeCycle.PostStart = &client.LifeCycleInner{}
						containerSpec.LifeCycle.PostStart.Exec = &client.Exec{}
						containerSpec.LifeCycle.PostStart.Exec.Command = buildCommand(exec)
					}
				}
			}

			if lc["pre_stop"] != nil {

				ps := lc["pre_stop"].([]interface{})

				if len(ps) > 0 && ps[0] != nil {

					psMap := ps[0].(map[string]interface{})
					exec := psMap["exec"].([]interface{})

					if len(exec) > 0 {
						containerSpec.LifeCycle.PreStop = &client.LifeCycleInner{}
						containerSpec.LifeCycle.PreStop.Exec = &client.Exec{}
						containerSpec.LifeCycle.PreStop.Exec.Command = buildCommand(exec)
					}
				}
			}
		}

	}
}

func buildOptions(options []interface{}, workloadSpec *client.WorkloadSpec, localOptions bool, org string) {

	if options == nil {
		return
	}

	output := []client.Options{}

	if len(options) > 0 {

		for _, o := range options {

			option := o.(map[string]interface{})

			newOptions := client.Options{}

			if localOptions {
				newOptions.Location = GetString(fmt.Sprintf("/org/%s/location/%s", org, option["location"].(string)))
			}

			newOptions.CapacityAI = GetBool(option["capacity_ai"])
			newOptions.TimeoutSeconds = GetInt(option["timeout_seconds"])
			newOptions.Debug = GetBool(option["debug"])
			newOptions.Suspend = GetBool(option["suspend"])

			autoScaling := option["autoscaling"].([]interface{})

			if len(autoScaling) > 0 {

				as := autoScaling[0].(map[string]interface{})

				cas := client.AutoScaling{

					Metric:           GetString(as["metric"]),
					MetricPercentile: GetString(as["metric_percentile"]),
					Target:           GetInt(as["target"]),
					MaxScale:         GetInt(as["max_scale"]),
					MinScale:         GetInt(as["min_scale"]),
					MaxConcurrency:   GetInt(as["max_concurrency"]),
					ScaleToZeroDelay: GetInt(as["scale_to_zero_delay"]),
				}

				newOptions.AutoScaling = &cas
			}

			output = append(output, newOptions)
		}
	}

	if workloadSpec == nil {
		workloadSpec = &client.WorkloadSpec{}
	}

	if localOptions {
		workloadSpec.LocalOptions = &output
	} else {
		workloadSpec.DefaultOptions = &output[0]
	}
}

func buildFirewallSpec(specs []interface{}, workloadSpec *client.WorkloadSpec) {

	if len(specs) > 0 {

		newSpec := client.FirewallSpec{}
		workloadSpec.FirewallConfig = &newSpec

		if specs[0] == nil {
			return
		}

		spec := specs[0].(map[string]interface{})
		external := spec["external"].([]interface{})

		if len(external) > 0 {

			we := client.FirewallSpecExternal{}
			newSpec.External = &we

			if external[0] != nil {

				e := external[0].(map[string]interface{})

				if e["inbound_allow_cidr"] != nil {
					inboundAllowCIDR := []string{}

					for _, value := range e["inbound_allow_cidr"].(*schema.Set).List() {
						inboundAllowCIDR = append(inboundAllowCIDR, value.(string))
					}

					we.InboundAllowCIDR = &inboundAllowCIDR

				}

				if e["outbound_allow_cidr"] != nil {
					outboundAllowCIDR := []string{}

					for _, value := range e["outbound_allow_cidr"].(*schema.Set).List() {
						outboundAllowCIDR = append(outboundAllowCIDR, value.(string))
					}

					we.OutboundAllowCIDR = &outboundAllowCIDR

				}

				if e["outbound_allow_hostname"] != nil {
					outboundAllowHostname := []string{}

					for _, value := range e["outbound_allow_hostname"].(*schema.Set).List() {
						outboundAllowHostname = append(outboundAllowHostname, value.(string))
					}

					we.OutboundAllowHostname = &outboundAllowHostname

				}

				if e["outbound_allow_port"] != nil {

					we.OutboundAllowPort = buildFirewallOutboundAllowPort(e["outbound_allow_port"].([]interface{}))
				}
			}

		}

		internal := spec["internal"].([]interface{})

		if len(internal) > 0 {

			wi := client.FirewallSpecInternal{}
			newSpec.Internal = &wi

			if internal[0] != nil {

				i := internal[0].(map[string]interface{})

				wi.InboundAllowType = GetString(i["inbound_allow_type"])

				if i["inbound_allow_workload"] != nil {
					inboundAllowWorkload := []string{}

					for _, value := range i["inbound_allow_workload"].(*schema.Set).List() {
						inboundAllowWorkload = append(inboundAllowWorkload, value.(string))
					}

					// if len(inboundAllowWorkload) > 0 {
					wi.InboundAllowWorkload = &inboundAllowWorkload
					// }
				}
			}
		}
	}
}

func buildFirewallOutboundAllowPort(specs []interface{}) *[]client.FirewallOutboundAllowPort {

	if len(specs) == 0 {
		return nil
	}

	output := []client.FirewallOutboundAllowPort{}

	for _, spec := range specs {

		specMap := spec.(map[string]interface{})
		outboundAllowPort := client.FirewallOutboundAllowPort{
			Protocol: GetString(specMap["protocol"]),
			Number:   GetInt(specMap["number"]),
		}

		output = append(output, outboundAllowPort)
	}

	return &output
}

func buildJobSpec(specs []interface{}) *client.JobSpec {

	if len(specs) > 0 && specs[0] != nil {

		result := &client.JobSpec{}

		spec := specs[0].(map[string]interface{})

		if spec["schedule"] != nil {
			result.Schedule = GetString(spec["schedule"].(string))
		}

		if spec["concurrency_policy"] != nil {
			result.ConcurrencyPolicy = GetString(spec["concurrency_policy"].(string))
		}

		if spec["history_limit"] != nil {
			result.HistoryLimit = GetInt(spec["history_limit"].(int))
		}

		if spec["restart_policy"] != nil {
			result.RestartPolicy = GetString(spec["restart_policy"].(string))
		}

		if spec["active_deadline_seconds"] != nil {

			if spec["active_deadline_seconds"].(int) == 0 {
				result.ActiveDeadlineSeconds = nil
			} else {
				result.ActiveDeadlineSeconds = GetInt(spec["active_deadline_seconds"].(int))
			}
		}

		return result
	}

	return nil
}

func buildRolloutOptions(specs []interface{}) *client.RolloutOptions {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.RolloutOptions{}

	if spec["min_ready_seconds"] != nil {
		output.MinReadySeconds = GetInt(spec["min_ready_seconds"].(int))
	}

	if spec["max_unavailable_replicas"] != nil {
		output.MaxUnavailableReplicas = GetString(spec["max_unavailable_replicas"].(string))
	}

	if spec["max_surge_replicas"] != nil {
		output.MaxSurgeReplicas = GetString(spec["max_surge_replicas"].(string))
	}

	if spec["scaling_policy"] != nil {
		output.ScalingPolicy = GetString(spec["scaling_policy"].(string))
	}

	return &output
}

func buildSecurityOptions(specs []interface{}) *client.SecurityOptions {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.SecurityOptions{
		FileSystemGroupID: GetInt(spec["file_system_group_id"].(int)),
	}

	return &output
}

func buildCommand(exec []interface{}) *[]string {

	if len(exec) > 0 {

		output := []string{}

		if exec[0] == nil {
			return &output
		}

		e := exec[0].(map[string]interface{})

		if e["command"] != nil {

			for _, k := range e["command"].([]interface{}) {
				if k != nil {
					output = append(output, k.(string))
				} else {
					output = append(output, "")
				}
			}

			return &output
		}
	}

	return nil
}

func buildWorkloadSidecar(specs []interface{}) *client.WorkloadSidecar {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.WorkloadSidecar{}

	// Attempt to unmarshal `envoy`
	var envoy interface{}
	err := json.Unmarshal([]byte(spec["envoy"].(string)), &envoy)
	if err != nil {
		log.Fatalf("Error occurred during unmarshaling 'envoy' value. Error: %s", err.Error())
	}

	// Set envoy
	output.Envoy = &envoy

	return &output
}

/*** Flatten ***/

func flattenWorkloadStatus(status *client.WorkloadStatus) []interface{} {

	if status != nil {

		fs := make(map[string]interface{})

		if status.ParentID != nil {
			fs["parent_id"] = *status.ParentID
		}

		if status.Endpoint != nil {
			fs["endpoint"] = *status.Endpoint
		}

		if status.CanonicalEndpoint != nil {
			fs["canonical_endpoint"] = *status.CanonicalEndpoint
		}

		if status.InternalName != nil {
			fs["internal_name"] = *status.InternalName
		}

		if status.CurrentReplicaCount != nil {
			fs["current_replica_count"] = *status.CurrentReplicaCount
		}

		if status.HealthCheck != nil {
			healthCheck := make(map[string]interface{})

			if status.HealthCheck.Active != nil {
				healthCheck["active"] = *status.HealthCheck.Active
			}

			if status.HealthCheck.Success != nil {
				healthCheck["success"] = *status.HealthCheck.Success
			}

			if status.HealthCheck.Code != nil {
				healthCheck["code"] = *status.HealthCheck.Code
			}

			if status.HealthCheck.Message != nil {
				healthCheck["message"] = *status.HealthCheck.Message
			}

			if status.HealthCheck.Failures != nil {
				healthCheck["failures"] = *status.HealthCheck.Failures
			}

			if status.HealthCheck.Successes != nil {
				healthCheck["successes"] = *status.HealthCheck.Successes
			}

			if status.HealthCheck.LastChecked != nil {
				healthCheck["last_checked"] = *status.HealthCheck.LastChecked
			}

			fs["health_check"] = healthCheck
		}

		resolvedImages := flattenWorkloadStatusResolvedImages(status.ResolvedImages)
		if resolvedImages != nil {
			fs["resolved_images"] = resolvedImages
		}

		output := []interface{}{
			fs,
		}

		return output
	}

	return nil
}

func flattenWorkloadStatusResolvedImages(resolvedImages *client.ResolvedImages) []interface{} {
	if resolvedImages == nil {
		return nil
	}

	output := make(map[string]interface{})

	if resolvedImages.ResolvedForVersion != nil {
		output["resolved_for_version"] = *resolvedImages.ResolvedForVersion
	}

	if resolvedImages.ResolvedAt != nil {
		output["resolved_at"] = *resolvedImages.ResolvedAt
	}

	if resolvedImages.Images != nil {
		output["images"] = flattenWorkloadStatusImages(resolvedImages.Images)
	}

	return []interface{}{
		output,
	}
}

func flattenWorkloadStatusImages(images *[]client.ResolvedImage) []interface{} {
	if images == nil || len(*images) == 0 {
		return nil
	}

	specs := []interface{}{}

	for _, image := range *images {

		spec := make(map[string]interface{})

		if image.Digest != nil {
			spec["digest"] = *image.Digest
		}

		if image.Manifests != nil {
			spec["manifests"] = flattenWorkloadStatusManifest(image.Manifests)
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenWorkloadStatusManifest(manifests *[]client.ResolvedImageManifest) []interface{} {
	if manifests == nil || len(*manifests) == 0 {
		return nil
	}

	specs := []interface{}{}

	for _, manifest := range *manifests {

		spec := make(map[string]interface{})

		if manifest.Image != nil {
			spec["image"] = *manifest.Image
		}

		if manifest.MediaType != nil {
			spec["media_type"] = *manifest.MediaType
		}

		if manifest.Digest != nil {
			spec["digest"] = *manifest.Digest
		}

		if manifest.Platform != nil {
			platform := make(map[string]interface{})

			for key, value := range *manifest.Platform {
				platform[key] = fmt.Sprintf("%v", value)
			}

			spec["platform"] = platform
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenContainer(containers *[]client.ContainerSpec, legacyPort bool) []interface{} {

	if containers != nil && len(*containers) > 0 {

		cs := make([]interface{}, len(*containers))

		for i, container := range *containers {

			c := make(map[string]interface{})

			c["name"] = *container.Name
			c["image"] = *container.Image

			// if container.Port != nil && *container.Port > 0 {
			// 	c["port"] = *container.Port
			// }

			if container.Ports != nil {

				if legacyPort {
					c["port"] = *((*container.Ports)[0].Number)
				} else {
					c["ports"] = flattenPortSpec(container.Ports)
				}
			} else if container.Port != nil {
				c["port"] = *container.Port
			}

			c["cpu"] = *container.CPU
			c["memory"] = *container.Memory

			if container.GPU != nil && container.GPU.Nvidia != nil {
				c["gpu_nvidia"] = flattenGpuNvidia(container.GPU)
			}

			if container.MinCPU != nil {
				c["min_cpu"] = *container.MinCPU
			}

			if container.MinMemory != nil {
				c["min_memory"] = *container.MinMemory
			}

			if container.Command != nil {
				c["command"] = *container.Command
			}

			if container.InheritEnv != nil {
				c["inherit_env"] = *container.InheritEnv
			}

			if container.WorkingDirectory != nil {
				c["working_directory"] = *container.WorkingDirectory
			}

			if container.Args != nil && len(*container.Args) > 0 {
				c["args"] = []interface{}{}

				for _, arg := range *container.Args {
					c["args"] = append(c["args"].([]interface{}), arg)
				}
			}

			if container.Env != nil && len(*container.Env) > 0 {
				envs := make(map[string]interface{})

				for _, env := range *container.Env {
					envs[*env.Name] = *env.Value
				}

				c["env"] = envs
			}

			if container.LivenessProbe != nil {
				c["liveness_probe"] = flattenHealthCheckSpec(container.LivenessProbe)
			}

			if container.ReadinessProbe != nil {
				c["readiness_probe"] = flattenHealthCheckSpec(container.ReadinessProbe)
			}

			if container.Volumes != nil {
				c["volume"] = flattenVolumeSpec(container.Volumes)
			}

			if container.Metrics != nil {
				c["metrics"] = flattenMetrics(container.Metrics)
			}

			if container.LifeCycle != nil {
				c["lifecycle"] = flattenLifeCycle(container.LifeCycle)
			}

			cs[i] = c
		}

		return cs
	}

	return nil
}

func flattenVolumeSpec(volumes *[]client.VolumeSpec) []interface{} {

	if volumes != nil && len(*volumes) > 0 {

		output := []interface{}{}

		for _, volume := range *volumes {

			v := map[string]interface{}{}

			if volume.Uri != nil {
				v["uri"] = *volume.Uri
			}

			if volume.RecoveryPolicy != nil {
				v["recovery_policy"] = *volume.RecoveryPolicy
			}

			if volume.Path != nil {
				v["path"] = *volume.Path
			}

			output = append(output, v)
		}

		return output
	}

	return nil
}

func flattenPortSpec(ports *[]client.PortSpec) []interface{} {

	if ports != nil && len(*ports) > 0 {

		output := []interface{}{}

		for _, port := range *ports {

			v := map[string]interface{}{}

			if port.Protocol != nil {
				v["protocol"] = *port.Protocol
			}

			if port.Number != nil {
				v["number"] = *port.Number
			}

			output = append(output, v)
		}

		return output
	}

	return nil
}

func flattenGpuNvidia(spec *client.GpuResource) []interface{} {
	if spec == nil || spec.Nvidia == nil {
		return nil
	}

	gpu := map[string]interface{}{
		"model":    *spec.Nvidia.Model,
		"quantity": *spec.Nvidia.Quantity,
	}

	return []interface{}{
		gpu,
	}
}

func flattenMetrics(metrics *client.Metrics) []interface{} {

	if metrics != nil {

		output := []interface{}{}

		m := map[string]interface{}{}

		if metrics.Path != nil {
			m["path"] = *metrics.Path
		}

		if metrics.Port != nil {
			m["port"] = *metrics.Port
		}

		output = append(output, m)

		return output
	}

	return nil
}

func flattenHealthCheckSpec(healthCheck *client.HealthCheckSpec) []interface{} {

	if healthCheck != nil {

		hc := map[string]interface{}{}

		if healthCheck.InitialDelaySeconds != nil {
			hc["initial_delay_seconds"] = *healthCheck.InitialDelaySeconds
		}

		if healthCheck.PeriodSeconds != nil {
			hc["period_seconds"] = *healthCheck.PeriodSeconds
		}

		if healthCheck.TimeoutSeconds != nil {
			hc["timeout_seconds"] = *healthCheck.TimeoutSeconds
		}

		if healthCheck.SuccessThreshold != nil {
			hc["success_threshold"] = *healthCheck.SuccessThreshold
		}

		if healthCheck.FailureThreshold != nil {
			hc["failure_threshold"] = *healthCheck.FailureThreshold
		}

		if healthCheck.Exec != nil && len(*healthCheck.Exec.Command) > 0 {
			e := make(map[string]interface{})
			e["command"] = *healthCheck.Exec.Command
			hc["exec"] = []interface{}{e}
		}

		if healthCheck.GRPC != nil {
			g := make(map[string]interface{})

			if healthCheck.GRPC.Port != nil && *healthCheck.GRPC.Port > 0 {
				g["port"] = *healthCheck.GRPC.Port
			}

			hc["grpc"] = []interface{}{g}
		}

		if healthCheck.TCPSocket != nil {
			t := make(map[string]interface{})

			if healthCheck.TCPSocket.Port != nil && *healthCheck.TCPSocket.Port > 0 {
				t["port"] = *healthCheck.TCPSocket.Port
			}

			ti := []interface{}{t}
			hc["tcp_socket"] = ti
		}

		if healthCheck.HTTPGet != nil {
			h := make(map[string]interface{})
			h["path"] = *healthCheck.HTTPGet.Path

			if healthCheck.HTTPGet.Port != nil && *healthCheck.HTTPGet.Port > 0 {
				h["port"] = *healthCheck.HTTPGet.Port
			}

			if healthCheck.HTTPGet.Scheme != nil {
				h["scheme"] = *healthCheck.HTTPGet.Scheme
			}

			headers := make(map[string]interface{})

			if healthCheck.HTTPGet.HTTPHeaders != nil {
				for _, header := range *healthCheck.HTTPGet.HTTPHeaders {
					if header.Value != nil {
						headers[*header.Name] = *header.Value
					} else {
						headers[*header.Name] = ""
					}
				}
			}

			h["http_headers"] = headers
			hi := []interface{}{h}
			hc["http_get"] = hi
		}

		return []interface{}{hc}
	}

	return nil
}

func flattenOptions(options []client.Options, localOptions bool, org string) []interface{} {

	if len(options) > 0 {

		output := []interface{}{}

		for _, o := range options {

			option := make(map[string]interface{})

			if localOptions && o.Location != nil {
				option["location"] = strings.TrimPrefix(*o.Location, fmt.Sprintf("/org/%s/location/", org))
			}

			if o.CapacityAI != nil {
				option["capacity_ai"] = *o.CapacityAI
			}

			if o.TimeoutSeconds != nil {
				option["timeout_seconds"] = *o.TimeoutSeconds
			}

			if o.Debug != nil {
				option["debug"] = *o.Debug
			}

			if o.Suspend != nil {
				option["suspend"] = *o.Suspend
			}

			as := make(map[string]interface{})

			if o.AutoScaling != nil {

				if o.AutoScaling.Metric != nil {
					as["metric"] = *o.AutoScaling.Metric
				}

				if o.AutoScaling.MetricPercentile != nil {
					as["metric_percentile"] = *o.AutoScaling.MetricPercentile
				}

				if o.AutoScaling.Target != nil {
					as["target"] = *o.AutoScaling.Target
				}

				if o.AutoScaling.MaxScale != nil {
					as["max_scale"] = *o.AutoScaling.MaxScale
				}

				if o.AutoScaling.MinScale != nil {
					as["min_scale"] = *o.AutoScaling.MinScale
				}

				if o.AutoScaling.MaxConcurrency != nil {
					as["max_concurrency"] = *o.AutoScaling.MaxConcurrency
				}

				if o.AutoScaling.ScaleToZeroDelay != nil {
					as["scale_to_zero_delay"] = *o.AutoScaling.ScaleToZeroDelay
				}
				autoScaling := make([]interface{}, 1)
				autoScaling[0] = as
				option["autoscaling"] = autoScaling
			}

			output = append(output, option)
		}

		return output
	}

	return nil
}

func flattenFirewallSpec(spec *client.FirewallSpec) []interface{} {

	if spec != nil {

		localSpec := make(map[string]interface{})

		if spec.External != nil {

			external := make(map[string]interface{})

			if spec.External.InboundAllowCIDR != nil && len(*spec.External.InboundAllowCIDR) > 0 {
				external["inbound_allow_cidr"] = []interface{}{}

				for _, arg := range *spec.External.InboundAllowCIDR {
					external["inbound_allow_cidr"] = append(external["inbound_allow_cidr"].([]interface{}), arg)
				}
			}

			if spec.External.OutboundAllowCIDR != nil && len(*spec.External.OutboundAllowCIDR) > 0 {
				external["outbound_allow_cidr"] = []interface{}{}

				for _, arg := range *spec.External.OutboundAllowCIDR {
					external["outbound_allow_cidr"] = append(external["outbound_allow_cidr"].([]interface{}), arg)
				}
			}

			if spec.External.OutboundAllowHostname != nil && len(*spec.External.OutboundAllowHostname) > 0 {
				external["outbound_allow_hostname"] = []interface{}{}

				for _, arg := range *spec.External.OutboundAllowHostname {
					external["outbound_allow_hostname"] = append(external["outbound_allow_hostname"].([]interface{}), arg)
				}
			}

			if spec.External.OutboundAllowPort != nil && len(*spec.External.OutboundAllowPort) > 0 {
				external["outbound_allow_port"] = flattenFirewallOutboundAllowPort(spec.External.OutboundAllowPort)
			}

			e := make([]interface{}, 1)
			e[0] = external
			localSpec["external"] = e
		}

		if spec.Internal != nil {

			internal := make(map[string]interface{})

			if spec.Internal.InboundAllowType != nil {
				internal["inbound_allow_type"] = *spec.Internal.InboundAllowType
			}

			if spec.Internal.InboundAllowWorkload != nil && len(*spec.Internal.InboundAllowWorkload) > 0 {
				internal["inbound_allow_workload"] = []interface{}{}

				for _, arg := range *spec.Internal.InboundAllowWorkload {
					internal["inbound_allow_workload"] = append(internal["inbound_allow_workload"].([]interface{}), arg)
				}
			}

			i := make([]interface{}, 1)
			i[0] = internal
			localSpec["internal"] = i
		}

		c := make([]interface{}, 1)
		c[0] = localSpec

		return c
	}

	return nil
}

func flattenFirewallOutboundAllowPort(outboundAllowPorts *[]client.FirewallOutboundAllowPort) []interface{} {

	if outboundAllowPorts == nil || len(*outboundAllowPorts) == 0 {
		return []interface{}{}
	}

	specs := []interface{}{}

	for _, outboundAllowPort := range *outboundAllowPorts {

		spec := map[string]interface{}{
			"protocol": *outboundAllowPort.Protocol,
			"number":   *outboundAllowPort.Number,
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenJobSpec(spec *client.JobSpec) []interface{} {

	if spec == nil {
		return nil
	}

	result := make(map[string]interface{})

	if spec.Schedule != nil {
		result["schedule"] = *spec.Schedule
	}

	if spec.ConcurrencyPolicy != nil {
		result["concurrency_policy"] = *spec.ConcurrencyPolicy
	}

	if spec.HistoryLimit != nil {
		result["history_limit"] = *spec.HistoryLimit
	}

	if spec.RestartPolicy != nil {
		result["restart_policy"] = *spec.RestartPolicy
	}

	if spec.ActiveDeadlineSeconds != nil {
		result["active_deadline_seconds"] = *spec.ActiveDeadlineSeconds
	}

	return []interface{}{
		result,
	}
}

func flattenLifeCycle(spec *client.LifeCycleSpec) []interface{} {

	if spec != nil {

		lc := map[string]interface{}{}

		if spec.PostStart != nil {

			postStart := make(map[string]interface{})

			if spec.PostStart.Exec != nil {

				exec := make(map[string]interface{})

				if spec.PostStart.Exec.Command != nil && len(*spec.PostStart.Exec.Command) > 0 {
					exec["command"] = []interface{}{}

					for _, command := range *spec.PostStart.Exec.Command {
						exec["command"] = append(exec["command"].([]interface{}), command)
					}
				}

				postStart["exec"] = []interface{}{exec}
			}

			lc["post_start"] = []interface{}{postStart}
		}

		if spec.PreStop != nil {

			preStop := make(map[string]interface{})

			if spec.PreStop.Exec != nil {

				exec := make(map[string]interface{})

				if spec.PreStop.Exec.Command != nil && len(*spec.PreStop.Exec.Command) > 0 {
					exec["command"] = []interface{}{}

					for _, command := range *spec.PreStop.Exec.Command {
						exec["command"] = append(exec["command"].([]interface{}), command)
					}

					preStop["exec"] = []interface{}{exec}
				}

				lc["pre_stop"] = []interface{}{preStop}
			}

			lc["pre_stop"] = []interface{}{preStop}
		}

		return []interface{}{lc}
	}

	return nil
}

func flattenRolloutOptions(spec *client.RolloutOptions) []interface{} {

	if spec == nil {
		return nil
	}

	rolloutOptions := map[string]interface{}{}

	if spec.MinReadySeconds != nil {
		rolloutOptions["min_ready_seconds"] = *spec.MinReadySeconds
	}

	if spec.MaxUnavailableReplicas != nil {
		rolloutOptions["max_unavailable_replicas"] = *spec.MaxUnavailableReplicas
	}

	if spec.MaxSurgeReplicas != nil {
		rolloutOptions["max_surge_replicas"] = *spec.MaxSurgeReplicas
	}

	if spec.ScalingPolicy != nil {
		rolloutOptions["scaling_policy"] = *spec.ScalingPolicy
	}

	return []interface{}{
		rolloutOptions,
	}
}

func flattenSecurityOptions(spec *client.SecurityOptions) []interface{} {
	if spec == nil {
		return nil
	}

	securityOptions := map[string]interface{}{
		"file_system_group_id": *spec.FileSystemGroupID,
	}

	return []interface{}{
		securityOptions,
	}
}

func flattenWorkloadSidecar(spec *client.WorkloadSidecar) []interface{} {
	if spec == nil {
		return nil
	}

	// Attempt to marshal `envoy`
	jsonOut, err := json.Marshal(*spec.Envoy)
	if err != nil {
		log.Fatalf("Error occurred during marshaling 'envoy' value. Error: %s", err.Error())
	}

	sidecar := map[string]interface{}{
		"envoy": string(jsonOut),
	}

	return []interface{}{
		sidecar,
	}
}

/*** Helpers ***/
func AutoScalingResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:        schema.TypeString,
				Description: "Valid values: `disabled`, `concurrency`, `cpu`, `latency`, or `rps`.",
				Optional:    true,
				Default:     "disabled",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "concurrency" && v != "cpu" && v != "rps" && v != "latency" && v != "disabled" {
						errs = append(errs, fmt.Errorf("%q must be 'concurrency', 'cpu', 'rps', 'latency' or 'disabled', got: %s", key, v))
					}

					return
				},
			},
			"metric_percentile": {
				Type:        schema.TypeString,
				Description: "For metrics represented as a distribution (e.g. latency) a percentile within the distribution must be chosen as the target.",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "p50" && v != "p75" && v != "p99" {
						errs = append(errs, fmt.Errorf("%q must be 'p50', 'p75' or 'p99', got: %s", key, v))
					}

					return
				},
			},
			"target": {
				Type:        schema.TypeInt,
				Description: "Control Plane will scale the number of replicas for this deployment up/down in order to be as close as possible to the target metric across all replicas of a deployment. Min: `0`. Max: `20000`. Default: `95`.",
				Optional:    true,
				Default:     95,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 20000 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 20000 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"max_scale": {
				Type:        schema.TypeInt,
				Description: "The maximum allowed number of replicas. Min: `0`. Default `5`.",
				Optional:    true,
				Default:     5,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be >= 0, got: %d", key, v))
					}
					return
				},
			},
			"min_scale": {
				Type:        schema.TypeInt,
				Description: "The minimum allowed number of replicas. Control Plane can scale the workload down to 0 when there is no traffic and scale up immediately to fulfill new requests. Min: `0`. Max: `max_scale`. Default `1`.",
				Optional:    true,
				Default:     1,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be >= 0, got: %d", key, v))
					}
					return
				},
			},
			"max_concurrency": {
				Type:        schema.TypeInt,
				Description: "A hard maximum for the number of concurrent requests allowed to a replica. If no replicas are available to fulfill the request then it will be queued until a replica with capacity is available and delivered as soon as one is available again. Capacity can be available from requests completing or when a new replica is available from scale out.Min: `0`. Max: `1000`. Default `0`.",
				Optional:    true,
				Default:     0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 30000 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 30000 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"scale_to_zero_delay": {
				Type:        schema.TypeInt,
				Description: "The amount of time (in seconds) with no requests received before a workload is scaled to 0. Min: `30`. Max: `3600`. Default: `300`.",
				Optional:    true,
				Default:     300,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 30 || v > 3600 {
						errs = append(errs, fmt.Errorf("%q must be between 30 and 3600 inclusive, got: %d", key, v))
					}
					return
				},
			},
		},
	}
}

func ExternalFirewallResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inbound_allow_cidr": {
				Type:        schema.TypeSet,
				Description: "The list of ipv4/ipv6 addresses or cidr blocks that are allowed to access this workload. No external access is allowed by default. Specify '0.0.0.0/0' to allow access to the public internet.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"outbound_allow_cidr": {
				Type:        schema.TypeSet,
				Description: "The list of ipv4/ipv6 addresses or cidr blocks that this workload is allowed reach. No outbound access is allowed by default. Specify '0.0.0.0/0' to allow outbound access to the public internet.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"outbound_allow_hostname": {
				Type:        schema.TypeSet,
				Description: "The list of public hostnames that this workload is allowed to reach. No outbound access is allowed by default. A wildcard `*` is allowed on the prefix of the hostname only, ex: `*.amazonaws.com`. Use `outboundAllowCIDR` to allow access to all external websites.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"outbound_allow_port": {
				Type:        schema.TypeList,
				Description: "Allow outbound access to specific ports and protocols. When not specified, communication to address ranges in outboundAllowCIDR is allowed on all ports and communication to names in outboundAllowHostname is allowed on ports 80/443.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Description: "Either `http`, `https` or `tcp`. Default: `tcp`.",
							Required:    true,
						},
						"number": {
							Type:        schema.TypeInt,
							Description: "Port number. Max: 65000",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func InternalFirewallResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inbound_allow_type": {
				Type:        schema.TypeString,
				Description: "Used to control the internal firewall configuration and mutual tls. Allowed Values: \"none\", \"same-gvc\", \"same-org\", \"workload-list\".",
				Optional:    true,
				Default:     "none",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "none" && v != "same-gvc" && v != "same-org" && v != "workload-list" {
						errs = append(errs, fmt.Errorf("%q must be 'none', 'same-gvc', 'same-org', or 'workload-list', got: %s", key, v))
					}

					return
				},
			},
			"inbound_allow_workload": {
				Type:        schema.TypeSet,
				Description: "A list of specific workloads which are allowed to access this workload internally. This list is only used if the 'inboundAllowType' is set to 'workload-list'.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func healthCheckSpec() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"initial_delay_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 600 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 600 inclusive, got: %d", key, v))
					}

					return
				},
			},
			"period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 60 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 60 inclusive, got: %d", key, v))
					}

					return
				},
			},
			"timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 60 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 60 inclusive, got: %d", key, v))
					}

					return
				},
			},
			"success_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: ThresholdValidator,
			},
			"failure_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: ThresholdValidator,
			},
			"exec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// ExactlyOneOf: []string{"http_get", "tcp_socket", "exec"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem:     StringSchema(),
						},
					},
				},
			},
			"grpc": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"tcp_socket": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// ExactlyOneOf: []string{"http_get", "tcp_socket", "exec"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 80 || v > 65535 {
									errs = append(errs, fmt.Errorf("%q must be between 80 and 65535 inclusive, got: %d", key, v))
									return
								}

								if v == 8012 || v == 8022 || v == 9090 || v == 9091 || v == 15000 || v == 15001 || v == 15006 || v == 15020 || v == 15021 || v == 15090 || v == 41000 {
									errs = append(errs, fmt.Errorf("%q cannot be 8012, 8022, 9090, 9091, 15000, 15001, 15006, 15020, 15021, 15090, 41000, got: %d", key, v))
								}

								return
							},
						},
					},
				},
			},
			"http_get": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// ExactlyOneOf: []string{"http_get", "tcp_socket", "exec"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/",
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: PortValidator,
						},
						"http_headers": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "HTTP",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								vLower := strings.ToLower(v)

								if vLower != "http" && vLower != "https" {
									errs = append(errs, fmt.Errorf("%q must be either HTTP or HTTPS: %s", key, v))
								}

								return
							},
						},
					},
				},
			},
		},
	}
}

func lifeCycleSpec() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem:     StringSchema(),
						},
					},
				},
			},
		},
	}
}

func workloadSpecValidate(workloadSpec *client.WorkloadSpec) diag.Diagnostics {

	if workloadSpec != nil {

		if *workloadSpec.Type == "cron" && workloadSpec.Job == nil {
			return diag.FromErr(fmt.Errorf("'job' section is required when workload type is 'cron'"))
		}

		if (*workloadSpec.Type != "standard" && *workloadSpec.Type != "stateful") && workloadSpec.RolloutOptions != nil {
			return diag.FromErr(fmt.Errorf("rollout options are only available when workload type is 'standard'"))
		}

		hasMinCpu := false
		hasMinMemory := false

		for _, c := range *workloadSpec.Containers {

			if *workloadSpec.Type == "cron" {
				if c.ReadinessProbe != nil || c.LivenessProbe != nil {
					return diag.FromErr(fmt.Errorf("probes are not allowed when workload type is 'cron'"))
				}
			}

			if c.GPU != nil && c.GPU.Nvidia != nil {
				cpuAmount, cpuUnit := ExtractNumberAndCharactersFromString(*c.CPU)
				memoryAmount, memoryUnit := ExtractNumberAndCharactersFromString(*c.Memory)

				// Return an error if the CPU amount is less than 2 and memory is less than 7Gi RAM
				if (cpuUnit == "" && cpuAmount < 2) ||
					(cpuUnit == "m" && cpuAmount < 2000) ||
					(memoryUnit == "Gi" && memoryAmount < 7) ||
					(memoryUnit == "Mi" && memoryAmount < 7000) {
					return diag.FromErr(fmt.Errorf("the GPU requires this container to have at least 2 CPU Cores and 7 Gi RAM"))
				}

				if *workloadSpec.DefaultOptions.CapacityAI {
					return diag.FromErr(fmt.Errorf("capacity AI must be disabled when using GPUs. Please remove the GPU selection from the containers or disable Capacity AI"))
				}
			}

			if c.MinCPU != nil {
				hasMinCpu = true
			}

			if c.MinMemory != nil {
				hasMinMemory = true
			}
		}

		if workloadSpec.DefaultOptions != nil {
			if e := validateOptions(*workloadSpec.Type, "", workloadSpec.DefaultOptions, hasMinCpu, hasMinMemory); e != nil {
				return e
			}
		}

		if workloadSpec.LocalOptions != nil && len(*workloadSpec.LocalOptions) > 0 {
			for _, o := range *workloadSpec.LocalOptions {
				if e := validateOptions(*workloadSpec.Type, "local_options - ", &o, hasMinCpu, hasMinMemory); e != nil {
					return e
				}
			}
		}
	}

	return nil
}

func validateOptions(workloadType, errorMsg string, options *client.Options, hasMinCpu bool, hasMinMemory bool) diag.Diagnostics {

	if options != nil && options.AutoScaling != nil {
		if workloadType == "cron" {
			if options.CapacityAI != nil && *options.CapacityAI {
				return diag.FromErr(fmt.Errorf(errorMsg + "capacity AI must be false when workload type is 'cron'"))
			}

			if hasMinCpu {
				return diag.FromErr(fmt.Errorf("min_cpu is not allowed for workload of type cron"))
			}

			if hasMinMemory {
				return diag.FromErr(fmt.Errorf("min_memory is not allowed for workload of type cron"))
			}

			if options.AutoScaling.MinScale != nil && *options.AutoScaling.MinScale != 1 {
				return diag.FromErr(fmt.Errorf(errorMsg + "min scale must be set to 1 when workload type is 'cron'"))
			}

			if options.AutoScaling.MaxScale != nil && *options.AutoScaling.MaxScale != 1 {
				return diag.FromErr(fmt.Errorf(errorMsg + "max scale must be set to 1 when workload type is 'cron'"))
			}
		} else {
			if options.AutoScaling.Metric == nil || strings.TrimSpace(*options.AutoScaling.Metric) == "" {
				return diag.FromErr(fmt.Errorf(errorMsg + "scaling strategy metric is required"))
			}

			if options.CapacityAI == nil || !*options.CapacityAI {
				if hasMinCpu {
					return diag.FromErr(fmt.Errorf("capacity AI must be enabled to include minimum CPU value"))
				}

				if hasMinMemory {
					return diag.FromErr(fmt.Errorf("capacity AI must be enabled to include minimum memory value"))
				}
			}
		}
	}

	return nil
}
