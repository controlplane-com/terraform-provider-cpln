package cpln

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	client "terraform-provider-cpln/internal/provider/client"

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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"cpln_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"identity_link": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: LinkValidator,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: WorkloadTypeValidator,
			},
			"container": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: PortValidator,
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: PortProtocolValidator,
										Default:      "http",
									},
									"number": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"cpu": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "50m",
							ValidateFunc: CpuMemoryValidator,
						},
						"memory": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "128Mi",
							ValidateFunc: CpuMemoryValidator,
						},
						"working_directory": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeString,
							Optional: true,
						},
						"env": {
							Type:     schema.TypeMap,
							Optional: true,
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
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"liveness_probe": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     healthCheckSpec(),
						},
						"readiness_probe": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     healthCheckSpec(),
						},
						"volume": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 5,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uri": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

											v := val.(string)

											re := regexp.MustCompile(`^(s3|gs|azureblob|azurefs|cpln|scratch):\/\/.+`)

											if !re.MatchString(v) {
												errs = append(errs, fmt.Errorf("%q must be in the form s3://bucket, gs://bucket, azureblob://storageAccount/container, azurefs://storageAccount/share, cpln://, scratch://, got: %s", key, v))
											}

											return
										},
									},
									"path": {
										Type:     schema.TypeString,
										Required: true,
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
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: PortValidator,
									},
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"lifecycle": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
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
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"debug": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"suspend": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 3600 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 3600 inclusive, got: %d", key, v))
								}
								return
							},
						},
						"autoscaling": {
							Type:     schema.TypeList,
							Required: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"capacity_ai": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"debug": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"suspend": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 3600 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 3600 inclusive, got: %d", key, v))
								}
								return
							},
						},
						"autoscaling": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     AutoScalingResource(),
						},
					},
				},
			},
			"firewall_spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     ExternalFirewallResource(),
						},
						"internal": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     InternalFirewallResource(),
						},
					},
				},
			},
			"job": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule": {
							Type:     schema.TypeString,
							Required: true,
						},
						"concurrency_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Forbid",
						},
						"history_limit": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
						"restart_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Never",
						},
						"active_deadline_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"canonical_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"internal_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"current_replica_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"health_check": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"active": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"success": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"code": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"message": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"failures": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"successes": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"last_checked": {
										Type:     schema.TypeString,
										Optional: true,
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
							Optional:    true,
							Default:     0,
							Description: "The minimum number of seconds a container must run without crashing to be considered available",
						},
						"max_unavailable_replicas": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_surge_replicas": {
							Type:     schema.TypeString,
							Optional: true,
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
							Required:    true,
							Description: "The group id assigned to any mounted volume",
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
		return nil, fmt.Errorf("unexpected format of ID (%s), expected gvc:workload", d.Id())
	}

	d.Set("gvc", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceWorkloadCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceWorkloadCreate")

	c := m.(*client.Client)

	gvcName := d.Get("gvc").(string)

	workload := client.Workload{}
	workload.Spec = &client.WorkloadSpec{}
	workload.Name = GetString(d.Get("name"))
	workload.Description = GetString(d.Get("description"))
	workload.Tags = GetStringMap(d.Get("tags"))

	buildContainers(d.Get("container").([]interface{}), workload.Spec)
	buildFirewallSpec(d.Get("firewall_spec").([]interface{}), workload.Spec)
	buildOptions(d.Get("options").([]interface{}), workload.Spec, false, c.Org)
	buildOptions(d.Get("local_options").([]interface{}), workload.Spec, true, c.Org)
	workload.Spec.Job = buildJobSpec(d.Get("job").([]interface{}))

	workload.Spec.Type = GetString(strings.TrimSpace(d.Get("type").(string)))

	if e := workloadSpecValidate(workload.Spec); e != nil {
		return e
	}

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

	newWorkload, code, err := c.CreateWorkload(workload, gvcName)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setWorkload(d, newWorkload, gvcName, c.Org, nil)
}

func resourceWorkloadRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceWorkloadRead")

	workloadName := d.Id()
	gvcName := d.Get("gvc").(string)

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

	return setWorkload(d, workload, gvcName, c.Org, diags)
}

func resourceWorkloadUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceWorkloadUpdate")

	if d.HasChanges("description", "tags", "type", "container", "options", "local_options", "firewall_spec", "job", "identity_link", "rollout_options", "security_options") {

		c := m.(*client.Client)

		gvcName := d.Get("gvc").(string)

		workloadToUpdate := client.Workload{}
		workloadToUpdate.Name = GetString(d.Get("name"))
		workloadToUpdate.Description = GetDescriptionString(d.Get("description"), *workloadToUpdate.Name)
		workloadToUpdate.Tags = GetTagChanges(d)

		workloadToUpdate.SpecReplace = &client.WorkloadSpec{}
		workloadToUpdate.SpecReplace.Type = GetString(d.Get("type"))

		buildContainers(d.Get("container").([]interface{}), workloadToUpdate.SpecReplace)
		buildOptions(d.Get("options").([]interface{}), workloadToUpdate.SpecReplace, false, c.Org)
		buildOptions(d.Get("local_options").([]interface{}), workloadToUpdate.SpecReplace, true, c.Org)
		buildFirewallSpec(d.Get("firewall_spec").([]interface{}), workloadToUpdate.SpecReplace)
		workloadToUpdate.SpecReplace.Job = buildJobSpec(d.Get("job").([]interface{}))
		workloadToUpdate.SpecReplace.RolloutOptions = buildRolloutOptions(d.Get("rollout_options").([]interface{}))
		workloadToUpdate.SpecReplace.SecurityOptions = buildSecurityOptions(d.Get("security_options").([]interface{}))
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

		return setWorkload(d, updatedWorkload, gvcName, c.Org, nil)
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

func setWorkload(d *schema.ResourceData, workload *client.Workload, gvcName, org string, diags diag.Diagnostics) diag.Diagnostics {

	if workload == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*workload.Name)

	if err := SetBase(d, workload.Base); err != nil {
		return diag.FromErr(err)
	}

	if workload.Spec != nil {
		if err := d.Set("container", flattenContainer(workload.Spec.Containers)); err != nil {
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
	}

	if err := d.Set("status", flattenWorkloadStatus(workload.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(workload.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func buildContainers(containers []interface{}, workloadSpec *client.WorkloadSpec) {

	if containers == nil {
		return
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

		if c["port"] != nil {
			newContainer.Port = GetPortInt(c["port"])
		}

		if c["ports"] != nil {
			newContainer.Ports = buildPortSpec(c["ports"].([]interface{}))
		}

		argArray := []string{}

		for _, value := range c["args"].([]interface{}) {
			argArray = append(argArray, value.(string))
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

func buildVolumeSpec(volumes []interface{}) *[]client.VolumeSpec {

	if len(volumes) > 0 {
		output := []client.VolumeSpec{}

		for _, value := range volumes {

			v := value.(map[string]interface{})

			uri := v["uri"].(string)
			path := v["path"].(string)

			localVolume := client.VolumeSpec{
				Uri:  &uri,
				Path: &path,
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

		output := []interface{}{
			fs,
		}

		return output
	}

	return nil
}

func flattenContainer(containers *[]client.ContainerSpec) []interface{} {

	if containers != nil && len(*containers) > 0 {

		cs := make([]interface{}, len(*containers))

		for i, container := range *containers {

			c := make(map[string]interface{})

			c["name"] = *container.Name
			c["image"] = *container.Image

			if container.Port != nil && *container.Port > 0 {
				c["port"] = *container.Port
			}

			if container.Ports != nil {
				c["ports"] = flattenPortSpec(container.Ports)
			}

			c["memory"] = *container.Memory
			c["cpu"] = *container.CPU

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

/*** Helpers ***/
func AutoScalingResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "concurrency" && v != "cpu" && v != "rps" && v != "latency" && v != "disabled" {
						errs = append(errs, fmt.Errorf("%q must be 'concurrency', 'cpu', 'rps', 'latency' or 'disabled', got: %s", key, v))
					}

					return
				},
			},
			"metric_percentile": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "p50" && v != "p75" && v != "p99" {
						errs = append(errs, fmt.Errorf("%q must be 'p50', 'p75' or 'p99', got: %s", key, v))
					}

					return
				},
			},
			"target": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  95,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 20000 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 20000 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"max_scale": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be >= 0, got: %d", key, v))
					}
					return
				},
			},
			"min_scale": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be >= 0, got: %d", key, v))
					}
					return
				},
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 30000 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 30000 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"scale_to_zero_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"outbound_allow_cidr": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"outbound_allow_hostname": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func InternalFirewallResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inbound_allow_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "none" && v != "same-gvc" && v != "same-org" && v != "workload-list" {
						errs = append(errs, fmt.Errorf("%q must be 'none', 'same-gvc', 'same-org', or 'workload-list', got: %s", key, v))
					}

					return
				},
			},
			"inbound_allow_workload": {
				Type:     schema.TypeSet,
				Optional: true,
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

		if *workloadSpec.Type != "standard" && workloadSpec.RolloutOptions != nil {
			return diag.FromErr(fmt.Errorf("rollout options are only available when workload type is 'standard'"))
		}

		for _, c := range *workloadSpec.Containers {

			if *workloadSpec.Type == "cron" {
				if c.ReadinessProbe != nil || c.LivenessProbe != nil {
					return diag.FromErr(fmt.Errorf("probes are not allowed when workload type is 'cron'"))
				}
			}
		}

		if workloadSpec.DefaultOptions != nil {
			if e := validateOptions(*workloadSpec.Type, "", workloadSpec.DefaultOptions); e != nil {
				return e
			}
		}

		if workloadSpec.LocalOptions != nil && len(*workloadSpec.LocalOptions) > 0 {
			for _, o := range *workloadSpec.LocalOptions {
				if e := validateOptions(*workloadSpec.Type, "local_options - ", &o); e != nil {
					return e
				}
			}
		}
	}

	return nil
}

func validateOptions(workloadType, errorMsg string, options *client.Options) diag.Diagnostics {

	if options != nil && options.AutoScaling != nil {
		if workloadType == "cron" {
			if options.CapacityAI != nil && *options.CapacityAI {
				return diag.FromErr(fmt.Errorf(errorMsg + "capacity AI must be false when workload type is 'cron'"))
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
		}
	}

	return nil
}
