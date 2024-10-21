package cpln

import (
	"context"
	"encoding/json"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMk8s() *schema.Resource {

	var mk8sProviders = []string{
		"generic_provider", "hetzner_provider", "aws_provider",
		"linode_provider", "oblivus_provider", "lambdalabs_provider",
		"paperspace_provider", "ephemeral_provider", "triton_provider",
		"digital_ocean_provider",
	}

	return &schema.Resource{
		CreateContext: resourceMk8sCreate,
		ReadContext:   resourceMk8sRead,
		UpdateContext: resourceMk8sUpdate,
		DeleteContext: resourceMk8sDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the Mk8s.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the Mk8s.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"alias": {
				Type:        schema.TypeString,
				Description: "The alias name of the Mk8s.",
				Computed:    true,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the Mk8s.",
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:         schema.TypeMap,
				Description:  "Key-value map of resource tags.",
				Optional:     true,
				Elem:         StringSchema(),
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "",
				Required:    true,
			},
			"firewall": {
				Type:        schema.TypeList,
				Description: "Allow-list.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_cidr": {
							Type:        schema.TypeString,
							Description: "",
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "",
							Optional:    true,
						},
					},
				},
			},
			"generic_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeString,
							Description: "Control Plane location that will host the K8s components. Prefer one that is closest to where the nodes are running.",
							Required:    true,
						},
						"networking": Mk8sNetworkingSchema(),
						"node_pool":  Mk8sGenericNodePoolSchema("List of node pools."),
					},
				},
			},
			"hetzner_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Hetzner region to deploy nodes to.",
							Required:    true,
						},
						"hetzner_labels": {
							Type:        schema.TypeMap,
							Description: "Extra labels to attach to servers.",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"networking":         Mk8sNetworkingSchema(),
						"pre_install_script": Mk8sPreInstallScriptSchema(),
						"token_secret_link": {
							Type:        schema.TypeString,
							Description: "Link to a secret holding Hetzner access key.",
							Required:    true,
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "ID of the Hetzner network to deploy nodes to.",
							Required:    true,
						},
						"firewall_id": {
							Type:        schema.TypeString,
							Description: "Optional firewall rule to attach to all nodes.",
							Optional:    true,
						},
						"node_pool":                  Mk8sHetznerNodePoolSchema(),
						"dedicated_server_node_pool": Mk8sGenericNodePoolSchema("Node pools that can configure dedicated Hetzner servers."),
						"image": {
							Type:        schema.TypeString,
							Description: "Default image for all nodes.",
							Optional:    true,
							Default:     "ubuntu-20.04",
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Description: "SSH key name for accessing deployed nodes.",
							Optional:    true,
						},
						"autoscaler": Mk8sAutoscalerSchema(),
						"floating_ip_selector": {
							Type:        schema.TypeMap,
							Description: "If supplied, nodes will get assigned a random floating ip matching the selector.",
							Optional:    true,
							Elem:        StringSchema(),
						},
					},
				},
			},
			"aws_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"aws_tags": {
							Type:        schema.TypeMap,
							Description: "Extra tags to attach to all created objects.",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"skip_create_roles": {
							Type:        schema.TypeBool,
							Description: "If true, Control Plane will not create any roles.",
							Optional:    true,
							Default:     false,
						},
						"networking":         Mk8sNetworkingSchema(),
						"pre_install_script": Mk8sPreInstallScriptSchema(),
						"image":              Mk8sAwsAmiSchema(),
						"deploy_role_arn": {
							Type:        schema.TypeString,
							Description: "Control Plane will set up the cluster by assuming this role.",
							Required:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "The vpc where nodes will be deployed. Supports SSM.",
							Required:    true,
						},
						"key_pair": {
							Type:        schema.TypeString,
							Description: "Name of keyPair. Supports SSM",
							Optional:    true,
						},
						"disk_encryption_key_arn": {
							Type:        schema.TypeString,
							Description: "KMS key used to encrypt volumes. Supports SSM.",
							Optional:    true,
						},
						"security_group_ids": {
							Type:        schema.TypeSet,
							Description: "Security groups to deploy nodes to. Security groups control if the cluster is multi-zone or single-zon.",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"node_pool":  Mk8sAwsNodePoolSchema(),
						"autoscaler": Mk8sAutoscalerSchema(),
					},
				},
			},
			"linode_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"token_secret_link": {
							Type:        schema.TypeString,
							Description: "Link to a secret holding Linode access key.",
							Required:    true,
						},
						"firewall_id": {
							Type:        schema.TypeString,
							Description: "Optional firewall rule to attach to all nodes.",
							Optional:    true,
						},
						"node_pool": Mk8sLinodeNodePoolSchema(),
						"image": {
							Type:        schema.TypeString,
							Description: "Default image for all nodes.",
							Required:    true,
						},
						"authorized_users": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"authorized_keys": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "The vpc where nodes will be deployed. Supports SSM.",
							Required:    true,
						},
						"pre_install_script": Mk8sPreInstallScriptSchema(),
						"networking":         Mk8sNetworkingSchema(),
						"autoscaler":         Mk8sAutoscalerSchema(),
					},
				},
			},
			"oblivus_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter": {
							Type:        schema.TypeString,
							Description: "",
							Required:    true,
						},
						"token_secret_link": {
							Type:        schema.TypeString,
							Description: "Link to a secret holding Oblivus access key.",
							Required:    true,
						},
						"node_pool": Mk8sOblivusNodePoolSchema(),
						"ssh_keys": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"unmanaged_node_pool": Mk8sGenericNodePoolSchema(""),
						"autoscaler":          Mk8sAutoscalerSchema(),
						"pre_install_script":  Mk8sPreInstallScriptSchema(),
					},
				},
			},
			"lambdalabs_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"token_secret_link": {
							Type:        schema.TypeString,
							Description: "Link to a secret holding Lambdalabs access key.",
							Required:    true,
						},
						"node_pool": Mk8sLambdalabsNodePoolSchema(),
						"ssh_key": {
							Type:        schema.TypeString,
							Description: "SSH key name for accessing deployed nodes.",
							Required:    true,
						},
						"unmanaged_node_pool": Mk8sGenericNodePoolSchema(""),
						"autoscaler":          Mk8sAutoscalerSchema(),
						"pre_install_script":  Mk8sPreInstallScriptSchema(),
					},
				},
			},
			"paperspace_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"token_secret_link": {
							Type:        schema.TypeString,
							Description: "Link to a secret holding Paperspace access key.",
							Required:    true,
						},
						"shared_drives": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"node_pool":           Mk8sPaperspaceNodePoolSchema(),
						"autoscaler":          Mk8sAutoscalerSchema(),
						"unmanaged_node_pool": Mk8sGenericNodePoolSchema(""),
						"pre_install_script":  Mk8sPreInstallScriptSchema(),
						"user_ids": {
							Type:        schema.TypeSet,
							Description: "",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "",
							Required:    true,
						},
					},
				},
			},
			"ephemeral_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeString,
							Description: "Control Plane location that will host the K8s components. Prefer one that is closest to where the nodes are running.",
							Required:    true,
						},
						"node_pool": Mk8sEphemeralNodePoolSchema(),
					},
				},
			},
			"triton_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection": {
							Type:        schema.TypeList,
							Description: "",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Description: "",
										Required:    true,
									},
									"account": {
										Type:        schema.TypeString,
										Description: "",
										Required:    true,
									},
									"user": {
										Type:        schema.TypeString,
										Description: "",
										Optional:    true,
									},
									"private_key_secret_link": {
										Type:        schema.TypeString,
										Description: "Link to a SSH or opaque secret.",
										Required:    true,
									},
								},
							},
						},
						"networking":         Mk8sNetworkingSchema(),
						"pre_install_script": Mk8sPreInstallScriptSchema(),
						"location": {
							Type:        schema.TypeString,
							Description: "Control Plane location that will host the K8s components. Prefer one that is closest to the Triton datacenter.",
							Required:    true,
						},
						"private_network_id": {
							Type:        schema.TypeString,
							Description: "ID of the private Fabric/Network.",
							Required:    true,
						},
						"firewall_enabled": {
							Type:        schema.TypeBool,
							Description: "Enable firewall for the instances deployed.",
							Optional:    true,
						},
						"node_pool": Mk8sTritonNodePoolSchema(),
						"image_id": {
							Type:        schema.TypeString,
							Description: "Default image for all nodes.",
							Required:    true,
						},
						"ssh_keys": {
							Type:        schema.TypeSet,
							Description: "Extra SSH keys to provision for user root.",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"autoscaler": Mk8sAutoscalerSchema(),
					},
				},
			},
			"digital_ocean_provider": {
				Type:         schema.TypeList,
				Description:  "",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: mk8sProviders,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Region to deploy nodes to.",
							Required:    true,
						},
						"digital_ocean_tags": {
							Type:        schema.TypeSet,
							Description: "Extra tags to attach to droplets.",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"networking":         Mk8sNetworkingSchema(),
						"pre_install_script": Mk8sPreInstallScriptSchema(),
						"token_secret_link": {
							Type:        schema.TypeString,
							Description: "Link to a secret holding personal access token.",
							Required:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "ID of the Hetzner network to deploy nodes to.",
							Required:    true,
						},
						"node_pool": Mk8sDigitalOceanNodePoolSchema(),
						"image": {
							Type:        schema.TypeString,
							Description: "Default image for all nodes.",
							Required:    true,
						},
						"ssh_keys": {
							Type:        schema.TypeSet,
							Description: "SSH key name for accessing deployed nodes.",
							Required:    true,
							Elem:        StringSchema(),
						},
						"extra_ssh_keys": {
							Type:        schema.TypeSet,
							Description: "Extra SSH keys to provision for user root that are not registered in the DigitalOcean.",
							Optional:    true,
							Elem:        StringSchema(),
						},
						"autoscaler": Mk8sAutoscalerSchema(),
						"reserved_ips": {
							Type:        schema.TypeSet,
							Description: "Optional set of IPs to assign as extra IPs for nodes of the cluster.",
							Optional:    true,
							Elem:        StringSchema(),
						},
					},
				},
			},
			"add_ons": {
				Type:        schema.TypeList,
				Description: "",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dashboard": {
							Type:        schema.TypeBool,
							Description: "",
							Optional:    true,
						},
						"azure_workload_identity": {
							Type:        schema.TypeList,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tenant_id": {
										Type:        schema.TypeString,
										Description: "Tenant ID to use for workload identity.",
										Optional:    true,
									},
									"_sentinel": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
						"aws_workload_identity": {
							Type:        schema.TypeBool,
							Description: "",
							Optional:    true,
						},
						"local_path_storage": {
							Type:        schema.TypeBool,
							Description: "",
							Optional:    true,
						},
						"metrics": {
							Type:        schema.TypeList,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kube_state": {
										Type:        schema.TypeBool,
										Description: "Enable kube-state metrics.",
										Optional:    true,
									},
									"core_dns": {
										Type:        schema.TypeBool,
										Description: "Enable scraping of core-dns service.",
										Optional:    true,
									},
									"kubelet": {
										Type:        schema.TypeBool,
										Description: "Enable scraping kubelet stats.",
										Optional:    true,
									},
									"api_server": {
										Type:        schema.TypeBool,
										Description: "Enable scraping apiserver stats.",
										Optional:    true,
									},
									"node_exporter": {
										Type:        schema.TypeBool,
										Description: "Enable collecting node-level stats (disk, network, filesystem, etc).",
										Optional:    true,
									},
									"cadvisor": {
										Type:        schema.TypeBool,
										Description: "Enable CNI-level container stats.",
										Optional:    true,
									},
									"scrape_annotated": {
										Type:        schema.TypeList,
										Description: "Scrape pods annotated with prometheus.io/scrape=true.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval_seconds": {
													Type:        schema.TypeInt,
													Description: "",
													Optional:    true,
												},
												"include_namespaces": {
													Type:        schema.TypeString,
													Description: "",
													Optional:    true,
												},
												"exclude_namespaces": {
													Type:        schema.TypeString,
													Description: "",
													Optional:    true,
												},
												"retain_labels": {
													Type:        schema.TypeString,
													Description: "",
													Optional:    true,
												},
												"_sentinel": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  true,
												},
											},
										},
									},
									"_sentinel": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
						"logs": {
							Type:        schema.TypeList,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"audit_enabled": {
										Type:        schema.TypeBool,
										Description: "Collect k8s audit log as log events.",
										Optional:    true,
									},
									"include_namespaces": {
										Type:        schema.TypeString,
										Description: "",
										Optional:    true,
									},
									"exclude_namespaces": {
										Type:        schema.TypeString,
										Description: "",
										Optional:    true,
									},
									"_sentinel": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
						"nvidia": {
							Type:        schema.TypeList,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"taint_gpu_nodes": {
										Type:        schema.TypeBool,
										Description: "",
										Optional:    true,
									},
									"_sentinel": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
						"aws_efs": Mk8sHasRoleArnSchema("Use this role for EFS interaction."),
						"aws_ecr": Mk8sHasRoleArnSchema("Role to use when authorizing ECR pulls. Optional on AWS, in which case it will use the instance role to pull."),
						"aws_elb": Mk8sHasRoleArnSchema("Role to use when authorizing calls to EC2 ELB. Optional on AWS, when not provided it will create the recommended role."),
						"azure_acr": {
							Type:        schema.TypeList,
							Description: "",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Type:        schema.TypeString,
										Description: "",
										Required:    true,
									},
								},
							},
						},
						"sysbox": {
							Type:        schema.TypeBool,
							Description: "",
							Optional:    true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeList,
				Description: "Status of the mk8s.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oidc_provider_url": {
							Type:        schema.TypeString,
							Description: "",
							Computed:    true,
						},
						"server_url": {
							Type:        schema.TypeString,
							Description: "",
							Computed:    true,
						},
						"home_location": {
							Type:        schema.TypeString,
							Description: "",
							Computed:    true,
						},
						"add_ons": {
							Type:        schema.TypeList,
							Description: "",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dashboard": {
										Type:        schema.TypeList,
										Description: "",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Description: "Access to dashboard.",
													Computed:    true,
												},
											},
										},
									},
									"aws_workload_identity": {
										Type:        schema.TypeList,
										Description: "",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"oidc_provider_config": {
													Type:        schema.TypeList,
													Description: "",
													Computed:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"provider_url": {
																Type:        schema.TypeString,
																Description: "",
																Computed:    true,
															},
															"audience": {
																Type:        schema.TypeString,
																Description: "",
																Computed:    true,
															},
														},
													},
												},
												"trust_policy": Mk8sObjectUnknownStatusSchema(),
											},
										},
									},
									"metrics": {
										Type:        schema.TypeList,
										Description: "",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prometheus_endpoint": {
													Type:        schema.TypeString,
													Description: "",
													Computed:    true,
												},
												"remote_write_config": Mk8sObjectUnknownStatusSchema(),
											},
										},
									},
									"logs": {
										Type:        schema.TypeList,
										Description: "",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"loki_address": {
													Type:        schema.TypeString,
													Description: "Loki endpoint to query logs from.",
													Computed:    true,
												},
											},
										},
									},
									"aws_ecr": Mk8sAwsAddOnStatusSchema(),
									"aws_efs": Mk8sAwsAddOnStatusSchema(),
									"aws_elb": Mk8sAwsAddOnStatusSchema(),
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceMk8sCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	// Define & Build
	mk8s := client.Mk8s{
		Spec: &client.Mk8sSpec{},
	}

	mk8s.Name = GetString(d.Get("name"))
	mk8s.Description = GetString(d.Get("description"))
	mk8s.Tags = GetStringMap(d.Get("tags"))

	mk8s.Spec.Version = GetString(d.Get("version"))

	if d.Get("firewall") != nil {
		mk8s.Spec.Firewall = buildMk8sFirewall(d.Get("firewall").([]interface{}))
	}

	mk8s.Spec.Provider = buildMk8sProvider(d)

	if d.Get("add_ons") != nil {
		mk8s.Spec.AddOns = buildMk8sAddOns(d.Get("add_ons").([]interface{}))
	}

	// Create
	newMk8s, code, err := c.CreateMk8s(mk8s)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setMk8s(d, newMk8s)
}

func resourceMk8sRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	mk8s, code, err := c.GetMk8s(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setMk8s(d, mk8s)
}

func resourceMk8sUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "version", "firewall", "generic_provider", "hetzner_provider", "aws_provider", "oblivus_provider", "linode_provider", "lambdalabs_provider", "paperspace_provider", "ephemeral_provider", "triton_provider", "digital_ocean_provider", "add_ons") {
		c := m.(*client.Client)

		// Define & Build
		mk8sToUpdate := client.Mk8s{
			SpecReplace: &client.Mk8sSpec{
				Version:  GetString(d.Get("version")),
				Firewall: buildMk8sFirewall(d.Get("firewall").([]interface{})),
				Provider: buildMk8sProvider(d),
			},
		}

		mk8sToUpdate.Name = GetString(d.Get("name"))
		mk8sToUpdate.Description = GetDescriptionString(d.Get("description"), *mk8sToUpdate.Name)
		mk8sToUpdate.Tags = GetTagChanges(d)

		if d.Get("add_ons") != nil {
			mk8sToUpdate.SpecReplace.AddOns = buildMk8sAddOns(d.Get("add_ons").([]interface{}))
		}

		// Update
		updatedMk8s, _, err := c.UpdateMk8s(mk8sToUpdate)

		if err != nil {
			return diag.FromErr(err)
		}

		return setMk8s(d, updatedMk8s)
	}

	return nil
}

func resourceMk8sDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	err := c.DeleteMk8s(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setMk8s(d *schema.ResourceData, mk8s *client.Mk8s) diag.Diagnostics {

	if mk8s == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*mk8s.Name)

	if err := SetBase(d, mk8s.Base); err != nil {
		return diag.FromErr(err)
	}

	if mk8s.Spec != nil {
		if err := d.Set("version", mk8s.Spec.Version); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("firewall", flattenMk8sFirewall(mk8s.Spec.Firewall)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("generic_provider", flattenMk8sGenericProvider(mk8s.Spec.Provider.Generic)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("hetzner_provider", flattenMk8sHetznerProvider(mk8s.Spec.Provider.Hetzner)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("aws_provider", flattenMk8sAwsProvider(mk8s.Spec.Provider.Aws)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("linode_provider", flattenMk8sLinodeProvider(mk8s.Spec.Provider.Linode)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("oblivus_provider", flattenMk8sOblivusProvider(mk8s.Spec.Provider.Oblivus)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("lambdalabs_provider", flattenMk8sLambdalabsProvider(mk8s.Spec.Provider.Lambdalabs)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("paperspace_provider", flattenMk8sPaperspaceProvider(mk8s.Spec.Provider.Paperspace)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("ephemeral_provider", flattenMk8sEphemeralProvider(mk8s.Spec.Provider.Ephemeral)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("triton_provider", flattenMk8sTritonProvider(mk8s.Spec.Provider.Triton)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("digital_ocean_provider", flattenMk8sDigitalOceanProvider(mk8s.Spec.Provider.DigitalOcean)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("add_ons", flattenMk8sAddOns(mk8s.Spec.AddOns)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("status", flattenMk8sStatus(mk8s.Status)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

/*** Build ***/

func buildMk8sFirewall(specs []interface{}) *[]client.Mk8sFirewallRule {

	if len(specs) == 0 {
		return nil
	}

	output := []client.Mk8sFirewallRule{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		firewallRule := client.Mk8sFirewallRule{
			SourceCIDR: GetString(spec["source_cidr"]),
		}

		if spec["description"] != nil {
			firewallRule.Description = GetString(spec["description"])
		}

		output = append(output, firewallRule)
	}

	return &output
}

func buildMk8sProvider(d *schema.ResourceData) *client.Mk8sProvider {

	output := client.Mk8sProvider{}

	if d.Get("generic_provider") != nil {
		output.Generic = buildMk8sGenericProvider(d.Get("generic_provider").([]interface{}))
	}

	if d.Get("hetzner_provider") != nil {
		output.Hetzner = buildMk8sHetznerProvider(d.Get("hetzner_provider").([]interface{}))
	}

	if d.Get("aws_provider") != nil {
		output.Aws = buildMk8sAwsProvider(d.Get("aws_provider").([]interface{}))
	}

	if d.Get("linode_provider") != nil {
		output.Linode = buildMk8sLinodeProvider(d.Get("linode_provider").([]interface{}))
	}

	if d.Get("oblivus_provider") != nil {
		output.Oblivus = buildMk8sOblivusProvider(d.Get("oblivus_provider").([]interface{}))
	}

	if d.Get("lambdalabs_provider") != nil {
		output.Lambdalabs = buildMk8sLambdalabsProvider(d.Get("lambdalabs_provider").([]interface{}))
	}

	if d.Get("paperspace_provider") != nil {
		output.Paperspace = buildMk8sPaperspaceProvider(d.Get("paperspace_provider").([]interface{}))
	}

	if d.Get("ephemeral_provider") != nil {
		output.Ephemeral = buildMk8sEphemeralProvider(d.Get("ephemeral_provider").([]interface{}))
	}

	if d.Get("triton_provider") != nil {
		output.Triton = buildMk8sTritonProvider(d.Get("triton_provider").([]interface{}))
	}

	if d.Get("digital_ocean_provider") != nil {
		output.DigitalOcean = buildMk8sDigitalOceanProvider(d.Get("digital_ocean_provider").([]interface{}))
	}

	return &output
}

func buildMk8sAddOns(specs []interface{}) *client.Mk8sSpecAddOns {

	if len(specs) == 0 && specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sSpecAddOns{}

	if spec["dashboard"] != nil && spec["dashboard"].(bool) {
		output.Dashboard = &client.Mk8sNonCustomizableAddonConfig{}
	}

	if spec["azure_workload_identity"] != nil {
		output.AzureWorkloadIdentity = buildMk8sAzureWorkloadIdentityAddOn(spec["azure_workload_identity"].([]interface{}))
	}

	if spec["aws_workload_identity"] != nil && spec["aws_workload_identity"].(bool) {
		output.AwsWorkloadIdentity = &client.Mk8sNonCustomizableAddonConfig{}
	}

	if spec["local_path_storage"] != nil && spec["local_path_storage"].(bool) {
		output.LocalPathStorage = &client.Mk8sNonCustomizableAddonConfig{}
	}

	if spec["metrics"] != nil {
		output.Metrics = buildMk8sMetricsAddOn(spec["metrics"].([]interface{}))
	}

	if spec["logs"] != nil {
		output.Logs = buildMk8sLogsAddOn(spec["logs"].([]interface{}))
	}

	if spec["nvidia"] != nil {
		output.Nvidia = buildMk8sNvidiaAddOn(spec["nvidia"].([]interface{}))
	}

	if spec["aws_efs"] != nil {
		output.AwsEFS = buildMk8sAwsAddOn(spec["aws_efs"].([]interface{}))
	}

	if spec["aws_ecr"] != nil {
		output.AwsECR = buildMk8sAwsAddOn(spec["aws_ecr"].([]interface{}))
	}

	if spec["aws_elb"] != nil {
		output.AwsELB = buildMk8sAwsAddOn(spec["aws_elb"].([]interface{}))
	}

	if spec["azure_acr"] != nil {
		output.AzureACR = buildMk8sAzureAcrAddOn(spec["azure_acr"].([]interface{}))
	}

	if spec["sysbox"] != nil && spec["sysbox"].(bool) {
		output.Sysbox = &client.Mk8sNonCustomizableAddonConfig{}
	}

	return &output
}

// Providers //

func buildMk8sGenericProvider(specs []interface{}) *client.Mk8sGenericProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sGenericProvider{}
	spec := specs[0].(map[string]interface{})

	output.Location = GetString(spec["location"])

	if spec["networking"] != nil {
		output.Networking = buildMk8sNetworking(spec["networking"].([]interface{}))
	}

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sGenericNodePools(spec["node_pool"].([]interface{}))
	}

	return &output
}

func buildMk8sHetznerProvider(specs []interface{}) *client.Mk8sHetznerProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sHetznerProvider{}
	spec := specs[0].(map[string]interface{})

	output.Region = GetString(spec["region"])

	if spec["hetzner_labels"] != nil {
		output.HetznerLabels = GetStringMap(spec["hetzner_labels"])
	}

	if spec["networking"] != nil {
		output.Networking = buildMk8sNetworking(spec["networking"].([]interface{}))
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	output.TokenSecretLink = GetString(spec["token_secret_link"])
	output.NetworkId = GetString(spec["network_id"])

	if spec["firewall_id"] != nil {
		output.FirewallId = GetString(spec["firewall_id"])
	}

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sHetznerNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["dedicated_server_node_pool"] != nil {
		output.DedicatedServerNodePools = buildMk8sGenericNodePools(spec["dedicated_server_node_pool"].([]interface{}))
	}

	if spec["image"] != nil {
		output.Image = GetString(spec["image"])
	}

	if spec["ssh_key"] != nil {
		output.SshKey = GetString(spec["ssh_key"])
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	if spec["floating_ip_selector"] != nil {
		output.FloatingIpSelector = GetStringMap(spec["floating_ip_selector"])
	}

	return &output
}

func buildMk8sAwsProvider(specs []interface{}) *client.Mk8sAwsProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sAwsProvider{}
	spec := specs[0].(map[string]interface{})

	output.Region = GetString(spec["region"])

	if spec["aws_tags"] != nil {
		output.AwsTags = GetStringMap(spec["aws_tags"])
	}

	if spec["skip_create_roles"] != nil {
		output.SkipCreateRoles = GetBool(spec["skip_create_roles"])
	}

	if spec["networking"] != nil {
		output.Networking = buildMk8sNetworking(spec["networking"].([]interface{}))
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	if spec["image"] != nil {
		output.Image = buildMk8sAwsAmi(spec["image"].([]interface{}))
	}

	output.DeployRoleArn = GetString(spec["deploy_role_arn"])
	output.VpcId = GetString(spec["vpc_id"])

	if spec["key_pair"] != nil {
		output.KeyPair = GetString(spec["key_pair"])
	}

	if spec["disk_encryption_key_arn"] != nil {
		output.DiskEncryptionKeyArn = GetString(spec["disk_encryption_key_arn"])
	}

	if spec["security_group_ids"] != nil {
		output.SecurityGroupIds = BuildStringTypeSet(spec["security_group_ids"])
	}

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sAwsNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	return &output
}

func buildMk8sLinodeProvider(specs []interface{}) *client.Mk8sLinodeProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sLinodeProvider{}
	spec := specs[0].(map[string]interface{})

	output.Region = GetString(spec["region"])
	output.TokenSecretLink = GetString(spec["token_secret_link"])
	output.Image = GetString(spec["image"])
	output.VpcId = GetString(spec["vpc_id"])

	if spec["firewall_id"] != nil {
		output.FirewallId = GetString(spec["firewall_id"])
	}

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sLinodeNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["authorized_users"] != nil {
		output.AuthorizedUsers = BuildStringTypeSet(spec["authorized_users"])
	}

	if spec["authorized_keys"] != nil {
		output.AuthorizedKeys = BuildStringTypeSet(spec["authorized_keys"])
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	if spec["networking"] != nil {
		output.Networking = buildMk8sNetworking(spec["networking"].([]interface{}))
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	return &output
}

func buildMk8sOblivusProvider(specs []interface{}) *client.Mk8sOblivusProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sOblivusProvider{}
	spec := specs[0].(map[string]interface{})

	output.Datacenter = GetString(spec["datacenter"])
	output.TokenSecretLink = GetString(spec["token_secret_link"])

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sOblivusNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["ssh_keys"] != nil {
		output.SshKeys = BuildStringTypeSet(spec["ssh_keys"])
	}

	if spec["unmanaged_node_pool"] != nil {
		output.UnmanagedNodePools = buildMk8sGenericNodePools(spec["unmanaged_node_pool"].([]interface{}))
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	return &output
}

func buildMk8sLambdalabsProvider(specs []interface{}) *client.Mk8sLambdalabsProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sLambdalabsProvider{}
	spec := specs[0].(map[string]interface{})

	output.Region = GetString(spec["region"])
	output.TokenSecretLink = GetString(spec["token_secret_link"])
	output.SshKey = GetString(spec["ssh_key"])

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sLambdalabsNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["unmanaged_node_pool"] != nil {
		output.UnmanagedNodePools = buildMk8sGenericNodePools(spec["unmanaged_node_pool"].([]interface{}))
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	return &output
}

func buildMk8sPaperspaceProvider(specs []interface{}) *client.Mk8sPaperspaceProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sPaperspaceProvider{}
	spec := specs[0].(map[string]interface{})

	output.Region = GetString(spec["region"])
	output.TokenSecretLink = GetString(spec["token_secret_link"])
	output.NetworkId = GetString(spec["network_id"])

	if spec["shared_drives"] != nil {
		output.SharedDrives = BuildStringTypeSet(spec["shared_drives"])
	}

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sPaperspaceNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	if spec["unmanaged_node_pool"] != nil {
		output.UnmanagedNodePools = buildMk8sGenericNodePools(spec["unmanaged_node_pool"].([]interface{}))
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	if spec["user_ids"] != nil {
		output.UserIds = BuildStringTypeSet(spec["user_ids"])
	}

	return &output
}

func buildMk8sEphemeralProvider(specs []interface{}) *client.Mk8sEphemeralProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sEphemeralProvider{}
	spec := specs[0].(map[string]interface{})

	output.Location = GetString(spec["location"])

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sEphemeralNodePools(spec["node_pool"].([]interface{}))
	}

	return &output
}

func buildMk8sTritonProvider(specs []interface{}) *client.Mk8sTritonProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sTritonProvider{}
	spec := specs[0].(map[string]interface{})

	output.Connection = buildMk8sTritonConnection(spec["connection"].([]interface{}))
	output.Location = GetString(spec["location"])
	output.PrivateNetworkId = GetString(spec["private_network_id"])
	output.ImageId = GetString(spec["image_id"])

	if spec["networking"] != nil {
		output.Networking = buildMk8sNetworking(spec["networking"].([]interface{}))
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	if spec["firewall_enabled"] != nil {
		output.FirewallEnabled = GetBool(spec["firewall_enabled"])
	}

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sTritonNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["ssh_keys"] != nil {
		output.SshKeys = BuildStringTypeSet(spec["ssh_keys"])
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	return &output
}

func buildMk8sDigitalOceanProvider(specs []interface{}) *client.Mk8sDigitalOceanProvider {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sDigitalOceanProvider{}
	spec := specs[0].(map[string]interface{})

	output.Region = GetString(spec["region"])
	output.Networking = buildMk8sNetworking(spec["networking"].([]interface{}))
	output.TokenSecretLink = GetString(spec["token_secret_link"])
	output.VpcId = GetString(spec["vpc_id"])
	output.Image = GetString(spec["image"])
	output.SshKeys = BuildStringTypeSet(spec["ssh_keys"])

	if spec["digital_ocean_tags"] != nil {
		output.DigitalOceanTags = BuildStringTypeSet(spec["digital_ocean_tags"])
	}

	if spec["pre_install_script"] != nil {
		output.PreInstallScript = GetString(spec["pre_install_script"])
	}

	if spec["node_pool"] != nil {
		output.NodePools = buildMk8sDigitalOceanNodePools(spec["node_pool"].([]interface{}))
	}

	if spec["extra_ssh_keys"] != nil {
		output.ExtraSshKeys = BuildStringTypeSet(spec["extra_ssh_keys"])
	}

	if spec["autoscaler"] != nil {
		output.Autoscaler = buildMk8sAutoscaler(spec["autoscaler"].([]interface{}))
	}

	if spec["reserved_ips"] != nil {
		output.ReservedIps = BuildStringTypeSet(spec["reserved_ips"])
	}

	return &output
}

// Provider Helpers //

// Node Pools

func buildMk8sGenericNodePools(specs []interface{}) *[]client.Mk8sGenericPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sGenericPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sGenericPool{
			Name: GetString(spec["name"]),
		}

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sHetznerNodePools(specs []interface{}) *[]client.Mk8sHetznerPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sHetznerPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sHetznerPool{}
		nodePool.Name = GetString(spec["name"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		nodePool.ServerType = GetString(spec["server_type"])

		if spec["override_image"] != nil {
			nodePool.OverrideImage = GetString(spec["override_image"])
		}

		if spec["min_size"] != nil {
			nodePool.MinSize = GetInt(spec["min_size"])
		}

		if spec["max_size"] != nil {
			nodePool.MaxSize = GetInt(spec["max_size"])
		}

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sAwsNodePools(specs []interface{}) *[]client.Mk8sAwsPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sAwsPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sAwsPool{}
		nodePool.Name = GetString(spec["name"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		nodePool.InstanceTypes = BuildStringTypeSet(spec["instance_types"])

		if spec["override_image"] != nil {
			nodePool.OverrideImage = buildMk8sAwsAmi(spec["override_image"].([]interface{}))
		}

		nodePool.BootDiskSize = GetInt(spec["boot_disk_size"])
		nodePool.MinSize = GetInt(spec["min_size"])
		nodePool.MaxSize = GetInt(spec["max_size"])
		nodePool.OnDemandBaseCapacity = GetInt(spec["on_demand_base_capacity"])
		nodePool.OnDemandPercentageAboveBaseCapacity = GetInt(spec["on_demand_percentage_above_base_capacity"])

		if spec["spot_allocation_strategy"] != nil {
			nodePool.SpotAllocationStrategy = GetString(spec["spot_allocation_strategy"])
		}

		nodePool.SubnetIds = BuildStringTypeSet(spec["subnet_ids"])

		if spec["extra_security_group_ids"] != nil {
			nodePool.ExtraSecurityGroupIds = BuildStringTypeSet(spec["extra_security_group_ids"])
		}

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sLinodeNodePools(specs []interface{}) *[]client.Mk8sLinodePool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sLinodePool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sLinodePool{}
		nodePool.Name = GetString(spec["name"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		if spec["override_image"] != nil {
			nodePool.OverrideImage = GetString(spec["override_image"])
		}

		nodePool.ServerType = GetString(spec["server_type"])
		nodePool.SubnetId = GetString(spec["subnet_id"])
		nodePool.MinSize = GetInt(spec["min_size"])
		nodePool.MaxSize = GetInt(spec["max_size"])

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sOblivusNodePools(specs []interface{}) *[]client.Mk8sOblivusPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sOblivusPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sOblivusPool{}
		nodePool.Name = GetString(spec["name"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		nodePool.MinSize = GetInt(spec["min_size"])
		nodePool.MaxSize = GetInt(spec["max_size"])
		nodePool.Flavor = GetString(spec["flavor"])

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sLambdalabsNodePools(specs []interface{}) *[]client.Mk8sLambdalabsPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sLambdalabsPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sLambdalabsPool{}
		nodePool.Name = GetString(spec["name"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		nodePool.MinSize = GetInt(spec["min_size"])
		nodePool.MaxSize = GetInt(spec["max_size"])
		nodePool.InstanceType = GetString(spec["instance_type"])

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sPaperspaceNodePools(specs []interface{}) *[]client.Mk8sPaperspacePool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sPaperspacePool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sPaperspacePool{}
		nodePool.Name = GetString(spec["name"])
		nodePool.MinSize = GetInt(spec["min_size"])
		nodePool.MaxSize = GetInt(spec["max_size"])
		nodePool.PublicIpType = GetString(spec["public_ip_type"])
		nodePool.MachineType = GetString(spec["machine_type"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		if spec["boot_disk_size"] != nil {
			nodePool.BootDiskSize = GetInt(spec["boot_disk_size"])
		}

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sEphemeralNodePools(specs []interface{}) *[]client.Mk8sEphemeralPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sEphemeralPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sEphemeralPool{
			Name:   GetString(spec["name"]),
			Count:  GetInt(spec["count"]),
			Arch:   GetString(spec["arch"]),
			Flavor: GetString(spec["flavor"]),
			Cpu:    GetString(spec["cpu"]),
			Memory: GetString(spec["memory"]),
		}

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sTritonNodePools(specs []interface{}) *[]client.Mk8sTritonPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sTritonPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sTritonPool{}
		nodePool.Name = GetString(spec["name"])
		nodePool.PackageId = GetString(spec["package_id"])
		nodePool.MinSize = GetInt(spec["min_size"])
		nodePool.MaxSize = GetInt(spec["max_size"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		if spec["override_image_id"] != nil {
			nodePool.OverrideImageId = GetString(spec["override_image_id"])
		}

		if spec["public_network_id"] != nil {
			nodePool.PublicNetworkId = GetString(spec["public_network_id"])
		}

		if spec["private_network_ids"] != nil {
			nodePool.PrivateNetworkIds = BuildStringTypeSet(spec["private_network_ids"])
		}

		if spec["triton_tags"] != nil {
			nodePool.TritonTags = GetStringMap(spec["triton_tags"])
		}

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sDigitalOceanNodePools(specs []interface{}) *[]client.Mk8sDigitalOceanPool {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sDigitalOceanPool{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sDigitalOceanPool{}
		nodePool.Name = GetString(spec["name"])
		nodePool.DropletSize = GetString(spec["droplet_size"])
		nodePool.MinSize = GetInt(spec["min_size"])
		nodePool.MaxSize = GetInt(spec["max_size"])

		if spec["labels"] != nil {
			nodePool.Labels = GetStringMap(spec["labels"])
		}

		if spec["taint"] != nil {
			nodePool.Taints = buildMk8sTaints(spec["taint"].([]interface{}))
		}

		if spec["override_image"] != nil {
			nodePool.OverrideImage = GetString(spec["override_image"])
		}

		output = append(output, nodePool)
	}

	return &output
}

// AWS

func buildMk8sAwsAmi(specs []interface{}) *client.Mk8sAwsAmi {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sAwsAmi{}
	spec := specs[0].(map[string]interface{})

	if spec["recommended"] != nil {
		output.Recommended = GetString(spec["recommended"])
	}

	if spec["exact"] != nil {
		output.Exact = GetString(spec["exact"])
	}

	return &output
}

// Triton

func buildMk8sTritonConnection(specs []interface{}) *client.Mk8sTritonConnection {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sTritonConnection{}
	spec := specs[0].(map[string]interface{})

	output.Url = GetString(spec["url"])
	output.Account = GetString(spec["account"])
	output.PrivateKeySecretLink = GetString(spec["private_key_secret_link"])

	if spec["user"] != nil {
		output.User = GetString(spec["user"])
	}

	return &output
}

// Common

func buildMk8sNetworking(specs []interface{}) *client.Mk8sNetworkingConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sNetworkingConfig{}
	spec := specs[0].(map[string]interface{})

	if spec["service_network"] != nil {
		output.ServiceNetwork = GetString(spec["service_network"])
	}

	if spec["pod_network"] != nil {
		output.PodNetwork = GetString(spec["pod_network"])
	}

	return &output
}

func buildMk8sTaints(specs []interface{}) *[]client.Mk8sTaint {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := []client.Mk8sTaint{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		nodePool := client.Mk8sTaint{}

		if spec["key"] != nil {
			nodePool.Key = GetString(spec["key"])
		}

		if spec["value"] != nil {
			nodePool.Value = GetString(spec["value"])
		}

		if spec["effect"] != nil {
			nodePool.Effect = GetString(spec["effect"])
		}

		output = append(output, nodePool)
	}

	return &output
}

func buildMk8sAutoscaler(specs []interface{}) *client.Mk8sAutoscalerConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	output := client.Mk8sAutoscalerConfig{}
	spec := specs[0].(map[string]interface{})

	output.Expander = BuildStringTypeSet(spec["expander"])
	output.UnneededTime = GetString(spec["unneeded_time"])
	output.UnreadyTime = GetString(spec["unready_time"])
	output.UtilizationThreshold = GetFloat64(spec["utilization_threshold"])

	return &output
}

// Add On Helpers //

func buildMk8sAzureWorkloadIdentityAddOn(specs []interface{}) *client.Mk8sAzureWorkloadIdentityAddOnConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sAzureWorkloadIdentityAddOnConfig{}

	if spec["tenant_id"] != nil {
		output.TenantId = GetString(spec["tenant_id"])
	}

	return &output
}

func buildMk8sMetricsAddOn(specs []interface{}) *client.Mk8sMetricsAddOnConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sMetricsAddOnConfig{}

	if spec["kube_state"] != nil {
		output.KubeState = GetBool(spec["kube_state"])
	}

	if spec["core_dns"] != nil {
		output.CoreDns = GetBool(spec["core_dns"])
	}

	if spec["kubelet"] != nil {
		output.Kubelet = GetBool(spec["kubelet"])
	}

	if spec["api_server"] != nil {
		output.Apiserver = GetBool(spec["api_server"])
	}

	if spec["node_exporter"] != nil {
		output.NodeExporter = GetBool(spec["node_exporter"])
	}

	if spec["cadvisor"] != nil {
		output.Cadvisor = GetBool(spec["cadvisor"])
	}

	if spec["scrape_annotated"] != nil {
		output.ScrapeAnnotated = buildMk8sMetricsScrapeAnnotated(spec["scrape_annotated"].([]interface{}))
	}

	return &output
}

func buildMk8sMetricsScrapeAnnotated(specs []interface{}) *client.Mk8sMetricsScrapeAnnotated {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sMetricsScrapeAnnotated{}

	if spec["interval_seconds"] != nil {
		output.IntervalSeconds = GetInt(spec["interval_seconds"])
	}

	if spec["include_namespaces"] != nil {
		output.IncludeNamespaces = GetString(spec["include_namespaces"])
	}

	if spec["exclude_namespaces"] != nil {
		output.ExcludeNamespaces = GetString(spec["exclude_namespaces"])
	}

	if spec["retain_labels"] != nil {
		output.RetainLabels = GetString(spec["retain_labels"])
	}

	return &output
}

func buildMk8sLogsAddOn(specs []interface{}) *client.Mk8sLogsAddOnConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sLogsAddOnConfig{}

	if spec["audit_enabled"] != nil {
		output.AuditEnabled = GetBool(spec["audit_enabled"])
	}

	if spec["include_namespaces"] != nil {
		output.IncludeNamespaces = GetString(spec["include_namespaces"])
	}

	if spec["exclude_namespaces"] != nil {
		output.ExcludeNamespaces = GetString(spec["exclude_namespaces"])
	}

	return &output
}

func buildMk8sNvidiaAddOn(specs []interface{}) *client.Mk8sNvidiaAddOnConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sNvidiaAddOnConfig{}

	if spec["taint_gpu_nodes"] != nil {
		output.TaintGPUNodes = GetBool(spec["taint_gpu_nodes"])
	}

	return &output
}

func buildMk8sAwsAddOn(specs []interface{}) *client.Mk8sAwsAddOnConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sAwsAddOnConfig{}

	if spec["role_arn"] != nil {
		output.RoleArn = GetString(spec["role_arn"])
	}

	return &output
}

func buildMk8sAzureAcrAddOn(specs []interface{}) *client.Mk8sAzureACRAddOnConfig {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Mk8sAzureACRAddOnConfig{
		ClientId: GetString(spec["client_id"]),
	}

	return &output
}

/*** Flatten ***/

func flattenMk8sFirewall(firewalls *[]client.Mk8sFirewallRule) []interface{} {

	if firewalls == nil {
		return nil
	}

	specs := []interface{}{}

	for _, firewall := range *firewalls {

		spec := map[string]interface{}{
			"source_cidr": *firewall.SourceCIDR,
		}

		if firewall.Description != nil {
			spec["description"] = *firewall.Description
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sAddOns(addOns *client.Mk8sSpecAddOns) []interface{} {

	if addOns == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if addOns.Dashboard != nil {
		spec["dashboard"] = true
	}

	if addOns.AzureWorkloadIdentity != nil {
		spec["azure_workload_identity"] = flattenMk8sAzureWorkloadIdentityAddOn(addOns.AzureWorkloadIdentity)
	}

	if addOns.AwsWorkloadIdentity != nil {
		spec["aws_workload_identity"] = true
	}

	if addOns.LocalPathStorage != nil {
		spec["local_path_storage"] = true
	}

	if addOns.Metrics != nil {
		spec["metrics"] = flattenMk8sMetricsAddOn(addOns.Metrics)
	}

	if addOns.Logs != nil {
		spec["logs"] = flattenMk8sLogsAddOn(addOns.Logs)
	}

	if addOns.Nvidia != nil {
		spec["nvidia"] = flattenMk8sNvidiaAddOn(addOns.Nvidia)
	}

	if addOns.AwsEFS != nil {
		spec["aws_efs"] = flattenMk8sAwsAddOn(addOns.AwsEFS)
	}

	if addOns.AwsECR != nil {
		spec["aws_ecr"] = flattenMk8sAwsAddOn(addOns.AwsECR)
	}

	if addOns.AwsELB != nil {
		spec["aws_elb"] = flattenMk8sAwsAddOn(addOns.AwsELB)
	}

	if addOns.AzureACR != nil {
		spec["azure_acr"] = flattenMk8sAzureAcrAddOn(addOns.AzureACR)
	}

	if addOns.Sysbox != nil {
		spec["sysbox"] = true
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sStatus(status *client.Mk8sStatus) []interface{} {

	if status == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if status.OidcProviderUrl != nil {
		spec["oidc_provider_url"] = *status.OidcProviderUrl
	}

	if status.ServerUrl != nil {
		spec["server_url"] = *status.ServerUrl
	}

	if status.HomeLocation != nil {
		spec["home_location"] = *status.HomeLocation
	}

	if status.AddOns != nil {
		spec["add_ons"] = flattenMk8sAddOnsStatus(status.AddOns)
	}

	return []interface{}{
		spec,
	}
}

// Providers //

func flattenMk8sGenericProvider(generic *client.Mk8sGenericProvider) []interface{} {

	if generic == nil {
		return nil
	}

	spec := map[string]interface{}{
		"location": *generic.Location,
	}

	if generic.Networking != nil {
		spec["networking"] = flattenMk8sNetworking(generic.Networking)
	}

	if generic.NodePools != nil {
		spec["node_pool"] = flattenMk8sGenericNodePools(generic.NodePools)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sHetznerProvider(hetzner *client.Mk8sHetznerProvider) []interface{} {

	if hetzner == nil {
		return nil
	}

	spec := map[string]interface{}{
		"region": *hetzner.Region,
	}

	if hetzner.HetznerLabels != nil {
		spec["hetzner_labels"] = *hetzner.HetznerLabels
	}

	if hetzner.Networking != nil {
		spec["networking"] = flattenMk8sNetworking(hetzner.Networking)
	}

	if hetzner.PreInstallScript != nil {
		spec["pre_install_script"] = *hetzner.PreInstallScript
	}

	spec["token_secret_link"] = *hetzner.TokenSecretLink
	spec["network_id"] = *hetzner.NetworkId

	if hetzner.FirewallId != nil {
		spec["firewall_id"] = *hetzner.FirewallId
	}

	if hetzner.NodePools != nil {
		spec["node_pool"] = flattenMk8sHetznerNodePools(hetzner.NodePools)
	}

	if hetzner.DedicatedServerNodePools != nil {
		spec["dedicated_server_node_pool"] = flattenMk8sGenericNodePools(hetzner.DedicatedServerNodePools)
	}

	if hetzner.Image != nil {
		spec["image"] = *hetzner.Image
	}

	if hetzner.SshKey != nil {
		spec["ssh_key"] = *hetzner.SshKey
	}

	if hetzner.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(hetzner.Autoscaler)
	}

	if hetzner.FloatingIpSelector != nil {
		spec["floating_ip_selector"] = *hetzner.FloatingIpSelector
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sAwsProvider(aws *client.Mk8sAwsProvider) []interface{} {

	if aws == nil {
		return nil
	}

	spec := map[string]interface{}{
		"region": *aws.Region,
	}

	if aws.AwsTags != nil {
		spec["aws_tags"] = *aws.AwsTags
	}

	if aws.SkipCreateRoles != nil {
		spec["skip_create_roles"] = *aws.SkipCreateRoles
	}

	if aws.Networking != nil {
		spec["networking"] = flattenMk8sNetworking(aws.Networking)
	}

	if aws.PreInstallScript != nil {
		spec["pre_install_script"] = *aws.PreInstallScript
	}

	spec["image"] = flattenMk8sAwsAmi(aws.Image)
	spec["deploy_role_arn"] = *aws.DeployRoleArn
	spec["vpc_id"] = *aws.VpcId

	if aws.KeyPair != nil {
		spec["key_pair"] = *aws.KeyPair
	}

	if aws.DiskEncryptionKeyArn != nil {
		spec["disk_encryption_key_arn"] = *aws.DiskEncryptionKeyArn
	}

	if aws.SecurityGroupIds != nil {
		spec["security_group_ids"] = FlattenStringTypeSet(aws.SecurityGroupIds)
	}

	if aws.NodePools != nil {
		spec["node_pool"] = flattenMk8sAwsNodePools(aws.NodePools)
	}

	if aws.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(aws.Autoscaler)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sLinodeProvider(linode *client.Mk8sLinodeProvider) []interface{} {

	if linode == nil {
		return nil
	}

	spec := map[string]interface{}{
		"region":            *linode.Region,
		"token_secret_link": *linode.TokenSecretLink,
		"image":             *linode.Image,
		"vpc_id":            *linode.VpcId,
	}

	if linode.FirewallId != nil {
		spec["firewall_id"] = *linode.FirewallId
	}

	if linode.NodePools != nil {
		spec["node_pool"] = flattenMk8sLinodeNodePools(linode.NodePools)
	}

	if linode.AuthorizedUsers != nil {
		spec["authorized_users"] = FlattenStringTypeSet(linode.AuthorizedUsers)
	}

	if linode.AuthorizedKeys != nil {
		spec["authorized_keys"] = FlattenStringTypeSet(linode.AuthorizedKeys)
	}

	if linode.PreInstallScript != nil {
		spec["pre_install_script"] = *linode.PreInstallScript
	}

	if linode.Networking != nil {
		spec["networking"] = flattenMk8sNetworking(linode.Networking)
	}

	if linode.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(linode.Autoscaler)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sOblivusProvider(oblivus *client.Mk8sOblivusProvider) []interface{} {

	if oblivus == nil {
		return nil
	}

	spec := map[string]interface{}{
		"datacenter":        *oblivus.Datacenter,
		"token_secret_link": *oblivus.TokenSecretLink,
	}

	if oblivus.NodePools != nil {
		spec["node_pool"] = flattenMk8sOblivusNodePools(oblivus.NodePools)
	}

	if oblivus.SshKeys != nil {
		spec["ssh_keys"] = FlattenStringTypeSet(oblivus.SshKeys)
	}

	if oblivus.UnmanagedNodePools != nil {
		spec["unmanaged_node_pool"] = flattenMk8sGenericNodePools(oblivus.UnmanagedNodePools)
	}

	if oblivus.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(oblivus.Autoscaler)
	}

	if oblivus.PreInstallScript != nil {
		spec["pre_install_script"] = *oblivus.PreInstallScript
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sLambdalabsProvider(lambdalabs *client.Mk8sLambdalabsProvider) []interface{} {

	if lambdalabs == nil {
		return nil
	}

	spec := map[string]interface{}{
		"region":            *lambdalabs.Region,
		"token_secret_link": *lambdalabs.TokenSecretLink,
		"ssh_key":           *lambdalabs.SshKey,
	}

	if lambdalabs.NodePools != nil {
		spec["node_pool"] = flattenMk8sLambdalabsNodePools(lambdalabs.NodePools)
	}

	if lambdalabs.UnmanagedNodePools != nil {
		spec["unmanaged_node_pool"] = flattenMk8sGenericNodePools(lambdalabs.UnmanagedNodePools)
	}

	if lambdalabs.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(lambdalabs.Autoscaler)
	}

	if lambdalabs.PreInstallScript != nil {
		spec["pre_install_script"] = *lambdalabs.PreInstallScript
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sPaperspaceProvider(paperspace *client.Mk8sPaperspaceProvider) []interface{} {

	if paperspace == nil {
		return nil
	}

	spec := map[string]interface{}{
		"region":            *paperspace.Region,
		"token_secret_link": *paperspace.TokenSecretLink,
		"network_id":        *paperspace.NetworkId,
	}

	if paperspace.SharedDrives != nil {
		spec["shared_drives"] = FlattenStringTypeSet(paperspace.SharedDrives)
	}

	if paperspace.NodePools != nil {
		spec["node_pool"] = flattenMk8sPaperspaceNodePools(paperspace.NodePools)
	}

	if paperspace.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(paperspace.Autoscaler)
	}

	if paperspace.UnmanagedNodePools != nil {
		spec["unmanaged_node_pool"] = flattenMk8sGenericNodePools(paperspace.UnmanagedNodePools)
	}

	if paperspace.PreInstallScript != nil {
		spec["pre_install_script"] = *paperspace.PreInstallScript
	}

	if paperspace.UserIds != nil {
		spec["user_ids"] = FlattenStringTypeSet(paperspace.UserIds)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sEphemeralProvider(ephemeral *client.Mk8sEphemeralProvider) []interface{} {

	if ephemeral == nil {
		return nil
	}

	spec := map[string]interface{}{
		"location": *ephemeral.Location,
	}

	if ephemeral.NodePools != nil {
		spec["node_pool"] = flattenMk8sEphemeralNodePools(ephemeral.NodePools)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sTritonProvider(triton *client.Mk8sTritonProvider) []interface{} {

	if triton == nil {
		return nil
	}

	spec := map[string]interface{}{
		"connection":         flattenMk8sTritonConnection(triton.Connection),
		"location":           *triton.Location,
		"private_network_id": *triton.PrivateNetworkId,
		"image_id":           *triton.ImageId,
	}

	if triton.Networking != nil {
		spec["networking"] = flattenMk8sNetworking(triton.Networking)
	}

	if triton.PreInstallScript != nil {
		spec["pre_install_script"] = *triton.PreInstallScript
	}

	if triton.FirewallEnabled != nil {
		spec["firewall_enabled"] = *triton.FirewallEnabled
	}

	if triton.NodePools != nil {
		spec["node_pool"] = flattenMk8sTritonNodePools(triton.NodePools)
	}

	if triton.SshKeys != nil {
		spec["ssh_keys"] = FlattenStringTypeSet(triton.SshKeys)
	}

	if triton.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(triton.Autoscaler)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sDigitalOceanProvider(digitalOcean *client.Mk8sDigitalOceanProvider) []interface{} {

	if digitalOcean == nil {
		return nil
	}

	spec := map[string]interface{}{
		"region":            *digitalOcean.Region,
		"token_secret_link": *digitalOcean.TokenSecretLink,
		"vpc_id":            *digitalOcean.VpcId,
		"image":             *digitalOcean.Image,
		"ssh_keys":          FlattenStringTypeSet(digitalOcean.SshKeys),
	}

	if digitalOcean.DigitalOceanTags != nil {
		spec["digital_ocean_tags"] = FlattenStringTypeSet(digitalOcean.DigitalOceanTags)
	}

	if digitalOcean.Networking != nil {
		spec["networking"] = flattenMk8sNetworking(digitalOcean.Networking)
	}

	if digitalOcean.PreInstallScript != nil {
		spec["pre_install_script"] = *digitalOcean.PreInstallScript
	}

	if digitalOcean.NodePools != nil {
		spec["node_pool"] = flattenMk8sDigitalOceanNodePools(digitalOcean.NodePools)
	}

	if digitalOcean.ExtraSshKeys != nil {
		spec["extra_ssh_keys"] = FlattenStringTypeSet(digitalOcean.ExtraSshKeys)
	}

	if digitalOcean.Autoscaler != nil {
		spec["autoscaler"] = flattenMk8sAutoscaler(digitalOcean.Autoscaler)
	}

	if digitalOcean.ReservedIps != nil {
		spec["reserved_ips"] = FlattenStringTypeSet(digitalOcean.ReservedIps)
	}

	return []interface{}{
		spec,
	}
}

// Provider Helpers //

// Node Pools

func flattenMk8sGenericNodePools(nodePools *[]client.Mk8sGenericPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name": *pool.Name,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sHetznerNodePools(nodePools *[]client.Mk8sHetznerPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name": *pool.Name,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		spec["server_type"] = *pool.ServerType

		if pool.OverrideImage != nil {
			spec["override_image"] = *pool.OverrideImage
		}

		spec["min_size"] = *pool.MinSize
		spec["max_size"] = *pool.MaxSize

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sAwsNodePools(nodePools *[]client.Mk8sAwsPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name": *pool.Name,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		spec["instance_types"] = FlattenStringTypeSet(pool.InstanceTypes)

		if pool.OverrideImage != nil {
			spec["override_image"] = flattenMk8sAwsAmi(pool.OverrideImage)
		}

		spec["boot_disk_size"] = *pool.BootDiskSize
		spec["min_size"] = *pool.MinSize
		spec["max_size"] = *pool.MaxSize
		spec["on_demand_base_capacity"] = *pool.OnDemandBaseCapacity
		spec["on_demand_percentage_above_base_capacity"] = *pool.OnDemandPercentageAboveBaseCapacity

		if pool.SpotAllocationStrategy != nil {
			spec["spot_allocation_strategy"] = *pool.SpotAllocationStrategy
		}

		spec["subnet_ids"] = FlattenStringTypeSet(pool.SubnetIds)

		if pool.ExtraSecurityGroupIds != nil {
			spec["extra_security_group_ids"] = FlattenStringTypeSet(pool.ExtraSecurityGroupIds)
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sLinodeNodePools(nodePools *[]client.Mk8sLinodePool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name": *pool.Name,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		if pool.OverrideImage != nil {
			spec["override_image"] = *pool.OverrideImage
		}

		spec["server_type"] = *pool.ServerType
		spec["subnet_id"] = *pool.SubnetId
		spec["min_size"] = *pool.MinSize
		spec["max_size"] = *pool.MaxSize

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sOblivusNodePools(nodePools *[]client.Mk8sOblivusPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name": *pool.Name,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		spec["min_size"] = *pool.MinSize
		spec["max_size"] = *pool.MaxSize
		spec["flavor"] = *pool.Flavor

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sLambdalabsNodePools(nodePools *[]client.Mk8sLambdalabsPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name": *pool.Name,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		spec["min_size"] = *pool.MinSize
		spec["max_size"] = *pool.MaxSize
		spec["instance_type"] = *pool.InstanceType

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sPaperspaceNodePools(nodePools *[]client.Mk8sPaperspacePool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name":           *pool.Name,
			"min_size":       *pool.MinSize,
			"max_size":       *pool.MaxSize,
			"public_ip_type": *pool.PublicIpType,
			"machine_type":   *pool.MachineType,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		if pool.BootDiskSize != nil {
			spec["boot_disk_size"] = *pool.BootDiskSize
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sEphemeralNodePools(nodePools *[]client.Mk8sEphemeralPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name":   *pool.Name,
			"count":  *pool.Count,
			"arch":   *pool.Arch,
			"flavor": *pool.Flavor,
			"cpu":    *pool.Cpu,
			"memory": *pool.Memory,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sTritonNodePools(nodePools *[]client.Mk8sTritonPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name":       *pool.Name,
			"package_id": *pool.PackageId,
			"min_size":   *pool.MinSize,
			"max_size":   *pool.MaxSize,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		if pool.OverrideImageId != nil {
			spec["override_image_id"] = *pool.OverrideImageId
		}

		if pool.PublicNetworkId != nil {
			spec["public_network_id"] = *pool.PublicNetworkId
		}

		if pool.PrivateNetworkIds != nil {
			spec["private_network_ids"] = FlattenStringTypeSet(pool.PrivateNetworkIds)
		}

		if pool.TritonTags != nil {
			spec["triton_tags"] = *pool.TritonTags
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sDigitalOceanNodePools(nodePools *[]client.Mk8sDigitalOceanPool) []interface{} {

	if nodePools == nil {
		return nil
	}

	specs := []interface{}{}

	for _, pool := range *nodePools {

		spec := map[string]interface{}{
			"name":         *pool.Name,
			"droplet_size": *pool.DropletSize,
			"min_size":     *pool.MinSize,
			"max_size":     *pool.MaxSize,
		}

		if pool.Labels != nil {
			spec["labels"] = *pool.Labels
		}

		if pool.Taints != nil {
			spec["taint"] = flattenMk8sTaints(pool.Taints)
		}

		if pool.OverrideImage != nil {
			spec["override_image"] = *pool.OverrideImage
		}

		specs = append(specs, spec)
	}

	return specs
}

// AWS

func flattenMk8sAwsAmi(ami *client.Mk8sAwsAmi) []interface{} {

	if ami == nil {
		return nil
	}

	spec := make(map[string]interface{})

	if ami.Recommended != nil {
		spec["recommended"] = *ami.Recommended
	}

	if ami.Exact != nil {
		spec["exact"] = *ami.Exact
	}

	return []interface{}{
		spec,
	}
}

// Triton

func flattenMk8sTritonConnection(connection *client.Mk8sTritonConnection) []interface{} {

	if connection == nil {
		return nil
	}

	spec := make(map[string]interface{})

	spec["url"] = *connection.Url
	spec["account"] = *connection.Account
	spec["private_key_secret_link"] = *connection.PrivateKeySecretLink

	if connection.User != nil {
		spec["user"] = *connection.User
	}

	return []interface{}{
		spec,
	}
}

// Common

func flattenMk8sNetworking(networking *client.Mk8sNetworkingConfig) []interface{} {

	if networking == nil {
		return nil
	}

	spec := make(map[string]interface{})

	if networking.ServiceNetwork != nil {
		spec["service_network"] = *networking.ServiceNetwork
	}

	if networking.PodNetwork != nil {
		spec["pod_network"] = *networking.PodNetwork
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sTaints(taints *[]client.Mk8sTaint) []interface{} {

	if taints == nil {
		return nil
	}

	specs := []interface{}{}

	for _, taint := range *taints {

		spec := make(map[string]interface{})

		if taint.Key != nil {
			spec["key"] = *taint.Key
		}

		if taint.Value != nil {
			spec["value"] = *taint.Value
		}

		if taint.Effect != nil {
			spec["effect"] = *taint.Effect
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenMk8sAutoscaler(autoscaler *client.Mk8sAutoscalerConfig) []interface{} {

	if autoscaler == nil {
		return nil
	}

	spec := make(map[string]interface{})

	if autoscaler.Expander != nil {
		spec["expander"] = FlattenStringTypeSet(autoscaler.Expander)
	}

	if autoscaler.UnneededTime != nil {
		spec["unneeded_time"] = *autoscaler.UnneededTime
	}

	if autoscaler.UnreadyTime != nil {
		spec["unready_time"] = *autoscaler.UnreadyTime
	}

	if autoscaler.UtilizationThreshold != nil {
		spec["utilization_threshold"] = *autoscaler.UtilizationThreshold
	}

	return []interface{}{
		spec,
	}
}

// Add On Helpers //

func flattenMk8sAzureWorkloadIdentityAddOn(azureWorkloadIdentity *client.Mk8sAzureWorkloadIdentityAddOnConfig) []interface{} {

	if azureWorkloadIdentity == nil {
		return nil
	}

	spec := map[string]interface{}{
		"_sentinel": true,
	}

	if azureWorkloadIdentity.TenantId != nil {
		spec["tenant_id"] = *azureWorkloadIdentity.TenantId
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sMetricsAddOn(metrics *client.Mk8sMetricsAddOnConfig) []interface{} {

	if metrics == nil {
		return nil
	}

	spec := map[string]interface{}{
		"_sentinel": true,
	}

	if metrics.KubeState != nil {
		spec["kube_state"] = *metrics.KubeState
	}

	if metrics.CoreDns != nil {
		spec["core_dns"] = *metrics.CoreDns
	}

	if metrics.Kubelet != nil {
		spec["kubelet"] = *metrics.Kubelet
	}

	if metrics.Apiserver != nil {
		spec["api_server"] = *metrics.Apiserver
	}

	if metrics.NodeExporter != nil {
		spec["node_exporter"] = *metrics.NodeExporter
	}

	if metrics.Cadvisor != nil {
		spec["cadvisor"] = *metrics.Cadvisor
	}

	if metrics.ScrapeAnnotated != nil {
		spec["scrape_annotated"] = flattenMk8sMetricsScrapeAnnotated(metrics.ScrapeAnnotated)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sMetricsScrapeAnnotated(scrapeAnnotated *client.Mk8sMetricsScrapeAnnotated) []interface{} {

	if scrapeAnnotated == nil {
		return nil
	}

	spec := map[string]interface{}{
		"_sentinel": true,
	}

	if scrapeAnnotated.IntervalSeconds != nil {
		spec["interval_seconds"] = *scrapeAnnotated.IntervalSeconds
	}

	if scrapeAnnotated.IncludeNamespaces != nil {
		spec["include_namespaces"] = *scrapeAnnotated.IncludeNamespaces
	}

	if scrapeAnnotated.ExcludeNamespaces != nil {
		spec["exclude_namespaces"] = *scrapeAnnotated.ExcludeNamespaces
	}

	if scrapeAnnotated.RetainLabels != nil {
		spec["retain_labels"] = *scrapeAnnotated.RetainLabels
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sLogsAddOn(logs *client.Mk8sLogsAddOnConfig) []interface{} {

	if logs == nil {
		return nil
	}

	spec := map[string]interface{}{
		"_sentinel": true,
	}

	if logs.AuditEnabled != nil {
		spec["audit_enabled"] = *logs.AuditEnabled
	}

	if logs.IncludeNamespaces != nil {
		spec["include_namespaces"] = *logs.IncludeNamespaces
	}

	if logs.ExcludeNamespaces != nil {
		spec["exclude_namespaces"] = *logs.ExcludeNamespaces
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sNvidiaAddOn(nvidia *client.Mk8sNvidiaAddOnConfig) []interface{} {

	if nvidia == nil {
		return nil
	}

	spec := map[string]interface{}{
		"_sentinel": true,
	}

	if nvidia.TaintGPUNodes != nil {
		spec["taint_gpu_nodes"] = *nvidia.TaintGPUNodes
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sAwsAddOn(aws *client.Mk8sAwsAddOnConfig) []interface{} {

	if aws == nil {
		return nil
	}

	spec := map[string]interface{}{
		"_sentinel": true,
	}

	if aws.RoleArn != nil {
		spec["role_arn"] = *aws.RoleArn
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sAzureAcrAddOn(azureAcr *client.Mk8sAzureACRAddOnConfig) []interface{} {

	if azureAcr == nil {
		return nil
	}

	spec := map[string]interface{}{
		"client_id": *azureAcr.ClientId,
	}

	return []interface{}{
		spec,
	}
}

// Status Helpers //

// Add Ons

func flattenMk8sAddOnsStatus(addOns *client.Mk8sStatusAddOns) []interface{} {

	if addOns == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if addOns.Dashboard != nil {
		spec["dashboard"] = flattenMk8sDashboardAddOnStatus(addOns.Dashboard)
	}

	if addOns.AwsWorkloadIdentity != nil {
		spec["aws_workload_identity"] = flattenMk8sAwsWorkloadIdentityAddOnStatus(addOns.AwsWorkloadIdentity)
	}

	if addOns.Metrics != nil {
		spec["metrics"] = flattenMk8sMetricsAddOnStatus(addOns.Metrics)
	}

	if addOns.Logs != nil {
		spec["logs"] = flattenMk8sLogsAddOnStatus(addOns.Logs)
	}

	if addOns.AwsECR != nil {
		spec["aws_ecr"] = flattenMk8sAwsAddOnStatus(addOns.AwsECR)
	}

	if addOns.AwsEFS != nil {
		spec["aws_efs"] = flattenMk8sAwsAddOnStatus(addOns.AwsEFS)
	}

	if addOns.AwsELB != nil {
		spec["aws_elb"] = flattenMk8sAwsAddOnStatus(addOns.AwsELB)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sDashboardAddOnStatus(dashboard *client.Mk8sDashboardAddOnStatus) []interface{} {

	if dashboard == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if dashboard.Url != nil {
		spec["url"] = *dashboard.Url
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sAwsWorkloadIdentityAddOnStatus(awsWorkloadIdentity *client.Mk8sAwsWorkloadIdentityAddOnStatus) []interface{} {

	if awsWorkloadIdentity == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if awsWorkloadIdentity.OidcProviderConfig != nil {
		spec["oidc_provider_config"] = flattenMk8sAwsOidcProviderConfigStatus(awsWorkloadIdentity.OidcProviderConfig)
	}

	if awsWorkloadIdentity.TrustPolicy != nil {
		spec["trust_policy"] = flattenObjectUnknown(awsWorkloadIdentity.TrustPolicy)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sMetricsAddOnStatus(metrics *client.Mk8sMetricsAddOnStatus) []interface{} {

	if metrics == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if metrics.PrometheusEndpoint != nil {
		spec["prometheus_endpoint"] = *metrics.PrometheusEndpoint
	}

	if metrics.RemoteWriteConfig != nil {
		spec["remote_write_config"] = flattenObjectUnknown(metrics.RemoteWriteConfig)
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sLogsAddOnStatus(logs *client.Mk8sLogsAddOnStatus) []interface{} {

	if logs == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if logs.LokiAddress != nil {
		spec["loki_address"] = *logs.LokiAddress
	}

	return []interface{}{
		spec,
	}
}

func flattenMk8sAwsAddOnStatus(aws *client.Mk8sAwsAddOnStatus) []interface{} {

	if aws == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if aws.TrustPolicy != nil {
		spec["trust_policy"] = flattenObjectUnknown(aws.TrustPolicy)
	}

	return []interface{}{
		spec,
	}
}

// Other

func flattenMk8sAwsOidcProviderConfigStatus(oidcProviderConfig *client.Mk8sOidcProviderConfig) []interface{} {

	if oidcProviderConfig == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if oidcProviderConfig.ProviderUrl != nil {
		spec["provider_url"] = *oidcProviderConfig.ProviderUrl
	}

	if oidcProviderConfig.Audience != nil {
		spec["audience"] = *oidcProviderConfig.Audience
	}

	return []interface{}{
		spec,
	}
}

func flattenObjectUnknown(unknown *map[string]interface{}) interface{} {

	if unknown == nil {
		return nil
	}

	// Convert map to JSON
	jsonData, _ := json.Marshal(*unknown)

	// Convert byte array to string
	return string(jsonData)
}

/*** Schema Helpers ***/

// Node Pools //

func Mk8sGenericNodePoolSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: description,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":   Mk8sGenericNodePoolNameSchema(),
				"labels": Mk8sGenericNodePoolLabelsSchema(),
				"taint":  Mk8sGenericNodePoolTaintsSchema(),
			},
		},
	}
}

func Mk8sHetznerNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":   Mk8sGenericNodePoolNameSchema(),
				"labels": Mk8sGenericNodePoolLabelsSchema(),
				"taint":  Mk8sGenericNodePoolTaintsSchema(),
				"server_type": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"override_image": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
			},
		},
	}
}

func Mk8sAwsNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":   Mk8sGenericNodePoolNameSchema(),
				"labels": Mk8sGenericNodePoolLabelsSchema(),
				"taint":  Mk8sGenericNodePoolTaintsSchema(),
				"instance_types": {
					Type:        schema.TypeSet,
					Description: "",
					Required:    true,
					Elem:        StringSchema(),
				},
				"override_image": Mk8sAwsAmiSchema(),
				"boot_disk_size": {
					Type:        schema.TypeInt,
					Description: "Size in GB.",
					Optional:    true,
					Default:     20,
				},
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
				"on_demand_base_capacity": {
					Type:        schema.TypeInt,
					Description: "",
					Optional:    true,
					Default:     0,
				},
				"on_demand_percentage_above_base_capacity": {
					Type:        schema.TypeInt,
					Description: "",
					Optional:    true,
					Default:     0,
				},
				"spot_allocation_strategy": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
					Default:     "lowest-price",
				},
				"subnet_ids": {
					Type:        schema.TypeSet,
					Description: "",
					Required:    true,
					Elem:        StringSchema(),
				},
				"extra_security_group_ids": {
					Type:        schema.TypeSet,
					Description: "",
					Optional:    true,
					Elem:        StringSchema(),
				},
			},
		},
	}
}

func Mk8sLinodeNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":   Mk8sGenericNodePoolNameSchema(),
				"labels": Mk8sGenericNodePoolLabelsSchema(),
				"taint":  Mk8sGenericNodePoolTaintsSchema(),
				"server_type": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"override_image": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
				"subnet_id": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
			},
		},
	}
}

func Mk8sOblivusNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":     Mk8sGenericNodePoolNameSchema(),
				"labels":   Mk8sGenericNodePoolLabelsSchema(),
				"taint":    Mk8sGenericNodePoolTaintsSchema(),
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
				"flavor": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
			},
		},
	}
}

func Mk8sPaperspaceNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":     Mk8sGenericNodePoolNameSchema(),
				"labels":   Mk8sGenericNodePoolLabelsSchema(),
				"taint":    Mk8sGenericNodePoolTaintsSchema(),
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
				"public_ip_type": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"boot_disk_size": {
					Type:        schema.TypeInt,
					Description: "",
					Optional:    true,
				},
				"machine_type": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
			},
		},
	}
}

func Mk8sLambdalabsNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":     Mk8sGenericNodePoolNameSchema(),
				"labels":   Mk8sGenericNodePoolLabelsSchema(),
				"taint":    Mk8sGenericNodePoolTaintsSchema(),
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
				"instance_type": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
			},
		},
	}
}

func Mk8sEphemeralNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":   Mk8sGenericNodePoolNameSchema(),
				"labels": Mk8sGenericNodePoolLabelsSchema(),
				"taint":  Mk8sGenericNodePoolTaintsSchema(),
				"count": {
					Type:        schema.TypeInt,
					Description: "Number of nodes to deploy.",
					Required:    true,
				},
				"arch": {
					Type:        schema.TypeString,
					Description: "CPU architecture of the nodes.",
					Required:    true,
				},
				"flavor": {
					Type:        schema.TypeString,
					Description: "Linux distro to use for ephemeral nodes.",
					Required:    true,
				},
				"cpu": {
					Type:        schema.TypeString,
					Description: "Allocated CPU.",
					Required:    true,
				},
				"memory": {
					Type:        schema.TypeString,
					Description: "Allocated memory.",
					Required:    true,
				},
			},
		},
	}
}

func Mk8sTritonNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":   Mk8sGenericNodePoolNameSchema(),
				"labels": Mk8sGenericNodePoolLabelsSchema(),
				"taint":  Mk8sGenericNodePoolTaintsSchema(),
				"package_id": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"override_image_id": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
				"public_network_id": {
					Type:        schema.TypeString,
					Description: "If set, machine will also get a public IP.",
					Optional:    true,
				},
				"private_network_ids": {
					Type:        schema.TypeSet,
					Description: "More private networks to join.",
					Optional:    true,
					Elem:        StringSchema(),
				},
				"triton_tags": {
					Type:        schema.TypeMap,
					Description: "Extra tags to attach to instances from a node pool.",
					Optional:    true,
					Elem:        StringSchema(),
				},
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
			},
		},
	}
}

func Mk8sDigitalOceanNodePoolSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of node pools.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":   Mk8sGenericNodePoolNameSchema(),
				"labels": Mk8sGenericNodePoolLabelsSchema(),
				"taint":  Mk8sGenericNodePoolTaintsSchema(),
				"droplet_size": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"override_image": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
				"min_size": Mk8sGenericNodePoolMinSizeSchema(),
				"max_size": Mk8sGenericNodePoolMaxSizeSchema(),
			},
		},
	}
}

// Node Pools Helpers //

func Mk8sGenericNodePoolNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	}
}

func Mk8sGenericNodePoolLabelsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Labels to attach to nodes of a node pool.",
		Optional:    true,
		Elem:        StringSchema(),
	}
}

func Mk8sGenericNodePoolTaintsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Taint for the nodes of a pool.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
				"effect": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
			},
		},
	}
}

func Mk8sGenericNodePoolMinSizeSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeInt,
		Description: "",
		Optional:    true,
		Default:     0,
	}
}

func Mk8sGenericNodePoolMaxSizeSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeInt,
		Description: "",
		Optional:    true,
		Default:     0,
	}
}

// AWS Helpers //

func Mk8sAwsAmiSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Default image for all nodes.",
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"recommended": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
				},
				"exact": {
					Type:        schema.TypeString,
					Description: "Support SSM.",
					Optional:    true,
				},
			},
		},
	}
}

func Mk8sHasRoleArnSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Description: description,
					Optional:    true,
				},
				"_sentinel": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
			},
		},
	}
}

// Common //

func Mk8sNetworkingSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "",
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_network": {
					Type:        schema.TypeString,
					Description: "The CIDR of the service network.",
					Optional:    true,
					Default:     "10.43.0.0/16",
				},
				"pod_network": {
					Type:        schema.TypeString,
					Description: "The CIDR of the pod network.",
					Optional:    true,
					Default:     "10.42.0.0/16",
				},
			},
		},
	}
}

func Mk8sAutoscalerSchema() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expander": {
					Type:        schema.TypeSet,
					Description: "",
					Required:    true,
					Elem:        StringSchema(),
				},
				"unneeded_time": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
					Default:     "10m",
				},
				"unready_time": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
					Default:     "20m",
				},
				"utilization_threshold": {
					Type:        schema.TypeFloat,
					Description: "",
					Optional:    true,
					Default:     0.7,
				},
			},
		},
	}
}

func Mk8sPreInstallScriptSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Optional shell script that will be run before K8s is installed. Supports SSM.",
		Optional:    true,
	}
}

// Status //

func Mk8sObjectUnknownStatusSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "",
		Computed:    true,
	}
}

func Mk8sAwsAddOnStatusSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"trust_policy": Mk8sObjectUnknownStatusSchema(),
			},
		},
	}
}
