package cpln

import (
	"context"
	"encoding/json"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/mk8s"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
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
	_ resource.Resource                = &Mk8sResource{}
	_ resource.ResourceWithImportState = &Mk8sResource{}
)

/*** Resource Model ***/

// Mk8sResourceModel holds the Terraform state for the resource.
type Mk8sResourceModel struct {
	EntityBaseModel
	Alias                types.String                       `tfsdk:"alias"`
	Version              types.String                       `tfsdk:"version"`
	Firewall             []models.FirewallModel             `tfsdk:"firewall"`
	GenericProvider      []models.GenericProviderModel      `tfsdk:"generic_provider"`
	HetznerProvider      []models.HetznerProviderModel      `tfsdk:"hetzner_provider"`
	AwsProvider          []models.AwsProviderModel          `tfsdk:"aws_provider"`
	LinodeProvider       []models.LinodeProviderModel       `tfsdk:"linode_provider"`
	OblivusProvider      []models.OblivusProviderModel      `tfsdk:"oblivus_provider"`
	LambdalabsProvider   []models.LambdalabsProviderModel   `tfsdk:"lambdalabs_provider"`
	PaperspaceProvider   []models.PaperspaceProviderModel   `tfsdk:"paperspace_provider"`
	EphemeralProvider    []models.EphemeralProviderModel    `tfsdk:"ephemeral_provider"`
	TritonProvider       []models.TritonProviderModel       `tfsdk:"triton_provider"`
	AzureProvider        []models.AzureProviderModel        `tfsdk:"azure_provider"`
	DigitalOceanProvider []models.DigitalOceanProviderModel `tfsdk:"digital_ocean_provider"`
	AddOns               []models.AddOnsModel               `tfsdk:"add_ons"`
	Status               types.List                         `tfsdk:"status"`
}

/*** Resource Configuration ***/

// Mk8sResource is the resource implementation.
type Mk8sResource struct {
	EntityBase
	Operations EntityOperations[Mk8sResourceModel, client.Mk8s]
}

// NewMk8sResource returns a new instance of the resource implementation.
func NewMk8sResource() resource.Resource {
	return &Mk8sResource{}
}

// Configure configures the resource before use.
func (mr *Mk8sResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	mr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	mr.Operations = NewEntityOperations(mr.client, &Mk8sResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (mr *Mk8sResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (mr *Mk8sResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_mk8s"
}

// Schema defines the schema for the resource.
func (mr *Mk8sResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(mr.EntityBaseAttributes("mk8s"), map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description: "The alias name of the Mk8s.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"status": schema.ListNestedAttribute{
				Description: "Status of the mk8s.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"oidc_provider_url": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"server_url": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"home_location": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"add_ons": schema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"dashboard": schema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"url": schema.StringAttribute{
													Description: "Access to dashboard.",
													Computed:    true,
												},
											},
										},
									},
									"aws_workload_identity": schema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"oidc_provider_config": schema.ListNestedAttribute{
													Description: "",
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"provider_url": schema.StringAttribute{
																Description: "",
																Computed:    true,
															},
															"audience": schema.StringAttribute{
																Description: "",
																Computed:    true,
															},
														},
													},
												},
												"trust_policy": mr.ObjectUnknownStatusSchema(),
											},
										},
									},
									"metrics": schema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"prometheus_endpoint": schema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"remote_write_config": mr.ObjectUnknownStatusSchema(),
											},
										},
									},
									"logs": schema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"loki_address": schema.StringAttribute{
													Description: "Loki endpoint to query logs from.",
													Computed:    true,
												},
											},
										},
									},
									"aws_ecr": mr.AwsAddOnStatusSchema(),
									"aws_efs": mr.AwsAddOnStatusSchema(),
									"aws_elb": mr.AwsAddOnStatusSchema(),
								},
							},
						},
					},
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"firewall": schema.ListNestedBlock{
				Description: "Allow-list.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"source_cidr": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "",
							Optional:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.IsRequired(),
				},
			},
			"generic_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"location": schema.StringAttribute{
							Description: "Control Plane location that will host the K8s components. Prefer one that is closest to where the nodes are running.",
							Required:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"networking": mr.NetworkingSchema(),
						"node_pool":  mr.GenericNodePoolSchema("List of node pools."),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"hetzner_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "Hetzner region to deploy nodes to.",
							Required:    true,
						},
						"hetzner_labels": schema.MapAttribute{
							Description: "Extra labels to attach to servers.",
							ElementType: types.StringType,
							Optional:    true,
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
						"token_secret_link": schema.StringAttribute{
							Description: "Link to a secret holding Hetzner access key.",
							Required:    true,
						},
						"network_id": schema.StringAttribute{
							Description: "ID of the Hetzner network to deploy nodes to.",
							Required:    true,
						},
						"firewall_id": schema.StringAttribute{
							Description: "Optional firewall rule to attach to all nodes.",
							Optional:    true,
						},
						"image": schema.StringAttribute{
							Description: "Default image for all nodes.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("ubuntu-20.04"),
						},
						"ssh_key": schema.StringAttribute{
							Description: "SSH key name for accessing deployed nodes.",
							Optional:    true,
						},
						"floating_ip_selector": schema.MapAttribute{
							Description: "If supplied, nodes will get assigned a random floating ip matching the selector.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"networking": mr.NetworkingSchema(),
						"node_pool": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   mr.GenericNodePoolNameSchema(),
									"labels": mr.GenericNodePoolLabelsSchema(),
									"server_type": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"override_image": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
						"dedicated_server_node_pool": mr.GenericNodePoolSchema("Node pools that can configure dedicated Hetzner servers."),
						"autoscaler":                 mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"aws_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"aws_tags": schema.MapAttribute{
							Description: "Extra tags to attach to all created objects.",
							ElementType: types.StringType,
							Optional:    true,
						},
						"skip_create_roles": schema.BoolAttribute{
							Description: "If true, Control Plane will not create any roles.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
						"deploy_role_arn": schema.StringAttribute{
							Description: "Control Plane will set up the cluster by assuming this role.",
							Required:    true,
						},
						"vpc_id": schema.StringAttribute{
							Description: "The vpc where nodes will be deployed. Supports SSM.",
							Required:    true,
						},
						"key_pair": schema.StringAttribute{
							Description: "Name of keyPair. Supports SSM",
							Optional:    true,
						},
						"disk_encryption_key_arn": schema.StringAttribute{
							Description: "KMS key used to encrypt volumes. Supports SSM.",
							Optional:    true,
						},
						"security_group_ids": schema.SetAttribute{
							Description: "Security groups to deploy nodes to. Security groups control if the cluster is multi-zone or single-zon.",
							ElementType: types.StringType,
							Optional:    true,
						},
						"extra_node_policies": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"networking": mr.NetworkingSchema(),
						"image":      mr.AwsAmiSchema(),
						"deploy_role_chain": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"role_arn": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"external_id": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"session_name_prefix": schema.StringAttribute{
										Description: "Control Plane will set up the cluster by assuming this role.",
										Optional:    true,
									},
								},
							},
						},
						"node_pool": schema.ListNestedBlock{
							Description: "List of node pools.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   mr.GenericNodePoolNameSchema(),
									"labels": mr.GenericNodePoolLabelsSchema(),
									"instance_types": schema.SetAttribute{
										Description: "",
										ElementType: types.StringType,
										Required:    true,
									},
									"boot_disk_size": schema.Int32Attribute{
										Description: "Size in GB.",
										Optional:    true,
										Computed:    true,
										Default:     int32default.StaticInt32(20),
									},
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
									"on_demand_base_capacity": schema.Int32Attribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										Default:     int32default.StaticInt32(0),
									},
									"on_demand_percentage_above_base_capacity": schema.Int32Attribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										Default:     int32default.StaticInt32(0),
									},
									"spot_allocation_strategy": schema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("lowest-price"),
									},
									"subnet_ids": schema.SetAttribute{
										Description: "",
										ElementType: types.StringType,
										Required:    true,
									},
									"extra_security_group_ids": schema.SetAttribute{
										Description: "Security groups to deploy nodes to. Security groups control if the cluster is multi-zone or single-zon.",
										ElementType: types.StringType,
										Optional:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"taint":          mr.GenericNodePoolTaintsSchema(),
									"override_image": mr.AwsAmiSchema(),
								},
							},
						},
						"autoscaler": mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"linode_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"token_secret_link": schema.StringAttribute{
							Description: "Link to a secret holding Linode access key.",
							Required:    true,
						},
						"firewall_id": schema.StringAttribute{
							Description: "Optional firewall rule to attach to all nodes.",
							Optional:    true,
						},
						"image": schema.StringAttribute{
							Description: "Default image for all nodes.",
							Required:    true,
						},
						"authorized_users": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
						},
						"authorized_keys": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
						},
						"vpc_id": schema.StringAttribute{
							Description: "The vpc where nodes will be deployed. Supports SSM.",
							Required:    true,
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
					},
					Blocks: map[string]schema.Block{
						"node_pool": schema.ListNestedBlock{
							Description: "List of node pools.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   mr.GenericNodePoolNameSchema(),
									"labels": mr.GenericNodePoolLabelsSchema(),
									"server_type": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"override_image": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"subnet_id": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
						"networking": mr.NetworkingSchema(),
						"autoscaler": mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"oblivus_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"datacenter": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"token_secret_link": schema.StringAttribute{
							Description: "Link to a secret holding Oblivus access key.",
							Required:    true,
						},
						"ssh_keys": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
					},
					Blocks: map[string]schema.Block{
						"node_pool": schema.ListNestedBlock{
							Description: "List of node pools.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":     mr.GenericNodePoolNameSchema(),
									"labels":   mr.GenericNodePoolLabelsSchema(),
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
									"flavor": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
						"unmanaged_node_pool": mr.GenericNodePoolSchema(""),
						"autoscaler":          mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"lambdalabs_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"token_secret_link": schema.StringAttribute{
							Description: "Link to a secret holding Lambdalabs access key.",
							Required:    true,
						},
						"ssh_key": schema.StringAttribute{
							Description: "SSH key name for accessing deployed nodes.",
							Required:    true,
						},
						"file_systems": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
					},
					Blocks: map[string]schema.Block{
						"node_pool": schema.ListNestedBlock{
							Description: "List of node pools.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":     mr.GenericNodePoolNameSchema(),
									"labels":   mr.GenericNodePoolLabelsSchema(),
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
									"instance_type": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
						"unmanaged_node_pool": mr.GenericNodePoolSchema(""),
						"autoscaler":          mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"paperspace_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"token_secret_link": schema.StringAttribute{
							Description: "Link to a secret holding Paperspace access key.",
							Required:    true,
						},
						"shared_drives": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
						"user_ids": schema.SetAttribute{
							Description: "",
							ElementType: types.StringType,
							Optional:    true,
						},
						"network_id": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"node_pool": schema.ListNestedBlock{
							Description: "List of node pools.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":     mr.GenericNodePoolNameSchema(),
									"labels":   mr.GenericNodePoolLabelsSchema(),
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
									"public_ip_type": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"boot_disk_size": schema.Int32Attribute{
										Description: "",
										Optional:    true,
									},
									"machine_type": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
						"autoscaler":          mr.AutoscalerSchema(),
						"unmanaged_node_pool": mr.GenericNodePoolSchema(""),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"ephemeral_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"location": schema.StringAttribute{
							Description: "Control Plane location that will host the K8s components. Prefer one that is closest to where the nodes are running.",
							Required:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"node_pool": schema.ListNestedBlock{
							Description: "List of node pools.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   mr.GenericNodePoolNameSchema(),
									"labels": mr.GenericNodePoolLabelsSchema(),
									"count": schema.Int32Attribute{
										Description: "Number of nodes to deploy.",
										Required:    true,
									},
									"arch": schema.StringAttribute{
										Description: "CPU architecture of the nodes.",
										Required:    true,
									},
									"flavor": schema.StringAttribute{
										Description: "Linux distro to use for ephemeral nodes.",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("debian"),
									},
									"cpu": schema.StringAttribute{
										Description: "Allocated CPU.",
										Required:    true,
									},
									"memory": schema.StringAttribute{
										Description: "Allocated memory.",
										Required:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"triton_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"pre_install_script": mr.PreInstallScriptSchema(),
						"location": schema.StringAttribute{
							Description: "Control Plane location that will host the K8s components. Prefer one that is closest to the Triton datacenter.",
							Required:    true,
						},
						"private_network_id": schema.StringAttribute{
							Description: "ID of the private Fabric/Network.",
							Required:    true,
						},
						"firewall_enabled": schema.BoolAttribute{
							Description: "Enable firewall for the instances deployed.",
							Optional:    true,
						},
						"image_id": schema.StringAttribute{
							Description: "Default image for all nodes.",
							Required:    true,
						},
						"ssh_keys": schema.SetAttribute{
							Description: "Extra SSH keys to provision for user root.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"connection": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"url": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"account": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"user": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"private_key_secret_link": schema.StringAttribute{
										Description: "Link to a SSH or opaque secret.",
										Required:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								listvalidator.SizeAtMost(1),
							},
						},
						"networking": mr.NetworkingSchema(),
						"load_balancer": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"manual": schema.ListNestedBlock{
										Description: "",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"package_id": schema.StringAttribute{
													Description: "",
													Required:    true,
												},
												"image_id": schema.StringAttribute{
													Description: "",
													Required:    true,
												},
												"public_network_id": schema.StringAttribute{
													Description: "If set, machine will also get a public IP.",
													Required:    true,
												},
												"private_network_ids": schema.SetAttribute{
													Description: "If set, machine will also get a public IP.",
													ElementType: types.StringType,
													Required:    true,
												},
												"metadata": schema.MapAttribute{
													Description: "Extra tags to attach to instances from a node pool.",
													ElementType: types.StringType,
													Optional:    true,
												},
												"tags": schema.MapAttribute{
													Description: "Extra tags to attach to instances from a node pool.",
													ElementType: types.StringType,
													Optional:    true,
												},
												"count": schema.Int32Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													Default:     int32default.StaticInt32(1),
													Validators: []validator.Int32{
														int32validator.Between(1, 3),
													},
												},
												"cns_internal_domain": schema.StringAttribute{
													Description: "",
													Required:    true,
												},
												"cns_public_domain": schema.StringAttribute{
													Description: "",
													Required:    true,
												},
											},
											Validators: []validator.Object{
												objectvalidator.ConflictsWith(
													path.MatchRelative().AtParent().AtParent().AtName("gateway"),
												),
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtMost(1),
										},
									},
									"gateway": schema.ListNestedBlock{
										Description: "",
										NestedObject: schema.NestedBlockObject{
											Validators: []validator.Object{
												objectvalidator.ConflictsWith(
													path.MatchRelative().AtParent().AtParent().AtName("manual"),
												),
											}},
										Validators: []validator.List{
											listvalidator.SizeAtMost(1),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								listvalidator.SizeAtMost(1),
							},
						},
						"node_pool": schema.ListNestedBlock{
							Description: "List of node pools.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   mr.GenericNodePoolNameSchema(),
									"labels": mr.GenericNodePoolLabelsSchema(),
									"package_id": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"override_image_id": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"public_network_id": schema.StringAttribute{
										Description: "If set, machine will also get a public IP.",
										Optional:    true,
									},
									"private_network_ids": schema.SetAttribute{
										Description: "More private networks to join.",
										ElementType: types.StringType,
										Optional:    true,
									},
									"triton_tags": schema.MapAttribute{
										Description: "Extra tags to attach to instances from a node pool.",
										ElementType: types.StringType,
										Optional:    true,
									},
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
						"autoscaler": mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"azure_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"location": schema.StringAttribute{
							Description: "Region where the cluster nodes will live.",
							Required:    true,
						},
						"subscription_id": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"sdk_secret_link": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"resource_group": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
						"ssh_keys": schema.SetAttribute{
							Description: `SSH keys to install for "azureuser" linux user`,
							ElementType: types.StringType,
							Required:    true,
						},
						"network_id": schema.StringAttribute{
							Description: "The vpc where nodes will be deployed.",
							Required:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Extra tags to attach to all created objects.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"networking": mr.NetworkingSchema(),
						"image":      mr.AzureImageSchema("Default image for all nodes."),
						"node_pool": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   mr.GenericNodePoolNameSchema(),
									"labels": mr.GenericNodePoolLabelsSchema(),
									"size": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"subnet_id": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"zones": schema.SetAttribute{
										Description: "",
										ElementType: types.Int32Type,
										Required:    true,
									},
									"boot_disk_size": schema.Int32Attribute{
										Description: "",
										Required:    true,
									},
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
								},
								Blocks: map[string]schema.Block{
									"taint":          mr.GenericNodePoolTaintsSchema(),
									"override_image": mr.AzureImageSchema(""),
								},
							},
						},
						"autoscaler": mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"digital_ocean_provider": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "Region to deploy nodes to.",
							Required:    true,
						},
						"digital_ocean_tags": schema.SetAttribute{
							Description: "Extra tags to attach to droplets.",
							ElementType: types.StringType,
							Optional:    true,
						},
						"pre_install_script": mr.PreInstallScriptSchema(),
						"token_secret_link": schema.StringAttribute{
							Description: "Link to a secret holding personal access token.",
							Required:    true,
						},
						"vpc_id": schema.StringAttribute{
							Description: "ID of the Hetzner network to deploy nodes to.",
							Required:    true,
						},
						"image": schema.StringAttribute{
							Description: "Default image for all nodes.",
							Required:    true,
						},
						"ssh_keys": schema.SetAttribute{
							Description: "SSH key name for accessing deployed nodes.",
							ElementType: types.StringType,
							Required:    true,
						},
						"extra_ssh_keys": schema.SetAttribute{
							Description: "Extra SSH keys to provision for user root that are not registered in the DigitalOcean.",
							ElementType: types.StringType,
							Optional:    true,
						},
						"reserved_ips": schema.SetAttribute{
							Description: "Optional set of IPs to assign as extra IPs for nodes of the cluster.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"networking": mr.NetworkingSchema(),
						"node_pool": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   mr.GenericNodePoolNameSchema(),
									"labels": mr.GenericNodePoolLabelsSchema(),
									"droplet_size": schema.StringAttribute{
										Description: "",
										Required:    true,
									},
									"override_image": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"min_size": mr.GenericNodePoolMinSizeSchema(),
									"max_size": mr.GenericNodePoolMaxSizeSchema(),
								},
								Blocks: map[string]schema.Block{
									"taint": mr.GenericNodePoolTaintsSchema(),
								},
							},
						},
						"autoscaler": mr.AutoscalerSchema(),
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"add_ons": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"dashboard": schema.BoolAttribute{
							Description: "",
							Optional:    true,
						},
						"aws_workload_identity": schema.BoolAttribute{
							Description: "",
							Optional:    true,
						},
						"local_path_storage": schema.BoolAttribute{
							Description: "",
							Optional:    true,
						},
						"sysbox": schema.BoolAttribute{
							Description: "",
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"azure_workload_identity": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"tenant_id": schema.StringAttribute{
										Description: "Tenant ID to use for workload identity.",
										Optional:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"metrics": schema.ListNestedBlock{
							Description: "Scrape pods annotated with prometheus.io/scrape=true",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"kube_state": schema.BoolAttribute{
										Description: "Enable kube-state metrics.",
										Optional:    true,
									},
									"core_dns": schema.BoolAttribute{
										Description: "Enable scraping of core-dns service.",
										Optional:    true,
									},
									"kubelet": schema.BoolAttribute{
										Description: "Enable scraping kubelet stats.",
										Optional:    true,
									},
									"api_server": schema.BoolAttribute{
										Description: "Enable scraping apiserver stats.",
										Optional:    true,
									},
									"node_exporter": schema.BoolAttribute{
										Description: "Enable collecting node-level stats (disk, network, filesystem, etc).",
										Optional:    true,
									},
									"cadvisor": schema.BoolAttribute{
										Description: "Enable CNI-level container stats.",
										Optional:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"scrape_annotated": schema.ListNestedBlock{
										Description: "",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"interval_seconds": schema.Int32Attribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													Default:     int32default.StaticInt32(30),
												},
												"include_namespaces": schema.StringAttribute{
													Description: "",
													Optional:    true,
												},
												"exclude_namespaces": schema.StringAttribute{
													Description: "",
													Optional:    true,
												},
												"retain_labels": schema.StringAttribute{
													Description: "",
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
						"logs": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"audit_enabled": schema.BoolAttribute{
										Description: "Collect k8s audit log as log events.",
										Optional:    true,
									},
									"include_namespaces": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"exclude_namespaces": schema.StringAttribute{
										Description: "",
										Optional:    true,
									},
									"docker": schema.BoolAttribute{
										Description: "Collect docker logs if docker is also running.",
										Optional:    true,
									},
									"kubelet": schema.BoolAttribute{
										Description: "Collect kubelet logs from journald.",
										Optional:    true,
									},
									"kernel": schema.BoolAttribute{
										Description: "Collect kernel logs.",
										Optional:    true,
									},
									"events": schema.BoolAttribute{
										Description: "Collect K8S events from all namespaces.",
										Optional:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"nvidia": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"taint_gpu_nodes": schema.BoolAttribute{
										Description: "",
										Optional:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"aws_efs": mr.HasRoleArnSchema("Use this role for EFS interaction."),
						"aws_ecr": mr.HasRoleArnSchema("Role to use when authorizing ECR pulls. Optional on AWS, in which case it will use the instance role to pull."),
						"aws_elb": mr.HasRoleArnSchema("Role to use when authorizing calls to EC2 ELB. Optional on AWS, when not provided it will create the recommended role."),
						"azure_acr": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"client_id": schema.StringAttribute{
										Description: "",
										Required:    true,
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

// ConfigValidators enforces mutual exclusivity between attributes.
func (mr *Mk8sResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	expressions := []path.Expression{
		path.MatchRoot("generic_provider"),
		path.MatchRoot("hetzner_provider"),
		path.MatchRoot("aws_provider"),
		path.MatchRoot("linode_provider"),
		path.MatchRoot("oblivus_provider"),
		path.MatchRoot("lambdalabs_provider"),
		path.MatchRoot("paperspace_provider"),
		path.MatchRoot("ephemeral_provider"),
		path.MatchRoot("triton_provider"),
		path.MatchRoot("azure_provider"),
		path.MatchRoot("digital_ocean_provider"),
	}

	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(expressions...),
	}
}

// Create creates the resource.
func (mr *Mk8sResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, mr.Operations)
}

// Read fetches the current state of the resource.
func (mr *Mk8sResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, mr.Operations)
}

// Update modifies the resource.
func (mr *Mk8sResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, mr.Operations)
}

// Delete removes the resource.
func (mr *Mk8sResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, mr.Operations)
}

/*** Schemas ***/

// NetworkingSchema returns the schema for the networking nested block.
func (mr *Mk8sResource) NetworkingSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"service_network": schema.StringAttribute{
					Description: "The CIDR of the service network.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("10.43.0.0/16"),
				},
				"pod_network": schema.StringAttribute{
					Description: "The CIDR of the pod network.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("10.42.0.0/16"),
				},
				"dns_forwarder": schema.StringAttribute{
					Description: "DNS forwarder used by the cluster. Can be a space-delimited list of dns servers. Default is /etc/resolv.conf when not specified.",
					Optional:    true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.SizeAtMost(1),
			listvalidator.IsRequired(),
		},
	}
}

// GenericNodePoolSchema returns the schema for a generic node pool nested block.
func (mr *Mk8sResource) GenericNodePoolSchema(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name":   mr.GenericNodePoolNameSchema(),
				"labels": mr.GenericNodePoolLabelsSchema(),
			},
			Blocks: map[string]schema.Block{
				"taint": mr.GenericNodePoolTaintsSchema(),
			},
		},
	}
}

// GenericNodePoolNameSchema returns the schema for the generic node pool name attribute.
func (mr *Mk8sResource) GenericNodePoolNameSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "",
		Required:    true,
	}
}

// GenericNodePoolLabelsSchema returns the schema for the generic node pool labels attribute.
func (mr *Mk8sResource) GenericNodePoolLabelsSchema() schema.MapAttribute {
	return schema.MapAttribute{
		Description: "Labels to attach to nodes of a node pool.",
		ElementType: types.StringType,
		Optional:    true,
	}
}

// GenericNodePoolTaintsSchema returns the schema for the generic node pool taints nested block.
func (mr *Mk8sResource) GenericNodePoolTaintsSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Taint for the nodes of a pool.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"value": schema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"effect": schema.StringAttribute{
					Description: "",
					Optional:    true,
				},
			},
		},
	}
}

// GenericNodePoolMinSizeSchema returns the schema for the generic node pool minimum size attribute.
func (mr *Mk8sResource) GenericNodePoolMinSizeSchema() schema.Int32Attribute {
	return schema.Int32Attribute{
		Description: "",
		Optional:    true,
		Computed:    true,
		Default:     int32default.StaticInt32(0),
	}
}

// GenericNodePoolMaxSizeSchema returns the schema for the generic node pool maximum size attribute.
func (mr *Mk8sResource) GenericNodePoolMaxSizeSchema() schema.Int32Attribute {
	return schema.Int32Attribute{
		Description: "",
		Optional:    true,
		Computed:    true,
		Default:     int32default.StaticInt32(0),
	}
}

// PreInstallScriptSchema returns the schema for the pre-installation script attribute.
func (mr *Mk8sResource) PreInstallScriptSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "Optional shell script that will be run before K8s is installed. Supports SSM.",
		Optional:    true,
	}
}

// AutoscalerSchema returns the schema for the cluster autoscaler nested block.
func (mr *Mk8sResource) AutoscalerSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"expander": schema.SetAttribute{
					Description: "",
					ElementType: types.StringType,
					Optional:    true,
					Computed:    true,
					Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("most-pods")})),
				},
				"unneeded_time": schema.StringAttribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("10m"),
				},
				"unready_time": schema.StringAttribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("20m"),
				},
				"utilization_threshold": schema.Float64Attribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Default:     float64default.StaticFloat64(0.7),
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// AwsAmiSchema returns the schema for the AWS AMI nested block.
func (mr *Mk8sResource) AwsAmiSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Default image for all nodes.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"recommended": schema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"exact": schema.StringAttribute{
					Description: "Support SSM.",
					Optional:    true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.SizeAtMost(1),
		},
	}
}

// HasRoleArnSchema returns the schema for the nested block that specifies a role ARN.
func (mr *Mk8sResource) HasRoleArnSchema(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"role_arn": schema.StringAttribute{
					Description: description,
					Optional:    true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// ObjectUnknownStatusSchema returns the schema for an object’s unknown status attribute.
func (mr *Mk8sResource) ObjectUnknownStatusSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "",
		Computed:    true,
	}
}

// AwsAddOnStatusSchema returns the schema for the AWS add-on status nested attribute.
func (mr *Mk8sResource) AwsAddOnStatusSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: "",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"trust_policy": mr.ObjectUnknownStatusSchema(),
			},
		},
	}
}

// AzureImageSchema returns a ListNestedBlock describing Azure VM image configuration with either a recommended image or a specific reference.
func (mr *Mk8sResource) AzureImageSchema(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"recommended": schema.StringAttribute{
					Description: "",
					Optional:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"reference": schema.ListNestedBlock{
					Description: "",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"publisher": schema.StringAttribute{
								Description: "",
								Required:    true,
							},
							"offer": schema.StringAttribute{
								Description: "",
								Required:    true,
							},
							"sku": schema.StringAttribute{
								Description: "",
								Required:    true,
							},
							"version": schema.StringAttribute{
								Description: "",
								Required:    true,
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
			listvalidator.SizeAtLeast(1),
			listvalidator.SizeAtMost(1),
			listvalidator.ExactlyOneOf(
				path.MatchRelative().AtName("recommended"),
				path.MatchRelative().AtName("reference"),
			),
		},
	}
}

/*** Resource Operator ***/

// Mk8sResourceOperator is the operator for managing the state.
type Mk8sResourceOperator struct {
	EntityOperator[Mk8sResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (mro *Mk8sResourceOperator) NewAPIRequest(isUpdate bool) client.Mk8s {
	// Initialize a new request payload
	requestPayload := client.Mk8s{}

	// Initialize the Mk8s spec struct
	var spec *client.Mk8sSpec = &client.Mk8sSpec{
		Provider: &client.Mk8sProvider{},
	}

	// Populate Base fields from state
	mro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Assignt he spec to the appropriate attribute
	if isUpdate {
		requestPayload.SpecReplace = spec
	} else {
		requestPayload.Spec = spec
	}

	// Set specific attributes
	spec.Version = BuildString(mro.Plan.Version)
	spec.Firewall = mro.buildFirewall(mro.Plan.Firewall)
	spec.Provider.Generic = mro.buildGenericProvider(mro.Plan.GenericProvider)
	spec.Provider.Hetzner = mro.buildHetznerProvider(mro.Plan.HetznerProvider)
	spec.Provider.Aws = mro.buildAwsProvider(mro.Plan.AwsProvider)
	spec.Provider.Linode = mro.buildLinodeProvider(mro.Plan.LinodeProvider)
	spec.Provider.Oblivus = mro.buildOblivusProvider(mro.Plan.OblivusProvider)
	spec.Provider.Lambdalabs = mro.buildLambdalabsProvider(mro.Plan.LambdalabsProvider)
	spec.Provider.Paperspace = mro.buildPaperspaceProvider(mro.Plan.PaperspaceProvider)
	spec.Provider.Ephemeral = mro.buildEphemeralProvider(mro.Plan.EphemeralProvider)
	spec.Provider.Triton = mro.buildTritonProvider(mro.Plan.TritonProvider)
	spec.Provider.Azure = mro.buildAzureProvider(mro.Plan.AzureProvider)
	spec.Provider.DigitalOcean = mro.buildDigitalOceanProvider(mro.Plan.DigitalOceanProvider)
	spec.AddOns = mro.buildAddOns(mro.Plan.AddOns)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (mro *Mk8sResourceOperator) MapResponseToState(apiResp *client.Mk8s, isCreate bool) Mk8sResourceModel {
	// Initialize empty state model
	state := Mk8sResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Set specific attributes
	state.Alias = types.StringPointerValue(apiResp.Alias)
	state.Version = types.StringPointerValue(apiResp.Spec.Version)
	state.Firewall = mro.flattenFirewall(apiResp.Spec.Firewall)
	state.GenericProvider = mro.flattenGenericProvider(apiResp.Spec.Provider.Generic)
	state.HetznerProvider = mro.flattenHetznerProvider(apiResp.Spec.Provider.Hetzner)
	state.AwsProvider = mro.flattenAwsProvider(apiResp.Spec.Provider.Aws)
	state.LinodeProvider = mro.flattenLinodeProvider(apiResp.Spec.Provider.Linode)
	state.OblivusProvider = mro.flattenOblivusProvider(apiResp.Spec.Provider.Oblivus)
	state.LambdalabsProvider = mro.flattenLambdalabsProvider(apiResp.Spec.Provider.Lambdalabs)
	state.PaperspaceProvider = mro.flattenPaperspaceProvider(apiResp.Spec.Provider.Paperspace)
	state.EphemeralProvider = mro.flattenEphemeralProvider(apiResp.Spec.Provider.Ephemeral)
	state.TritonProvider = mro.flattenTritonProvider(apiResp.Spec.Provider.Triton)
	state.AzureProvider = mro.flattenAzureProvider(apiResp.Spec.Provider.Azure)
	state.DigitalOceanProvider = mro.flattenDigitalOceanProvider(apiResp.Spec.Provider.DigitalOcean)
	state.AddOns = mro.flattenAddOns(apiResp.Spec.AddOns)
	state.Status = mro.flattenStatus(apiResp.Status)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (mro *Mk8sResourceOperator) InvokeCreate(req client.Mk8s) (*client.Mk8s, int, error) {
	return mro.Client.CreateMk8s(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (mro *Mk8sResourceOperator) InvokeRead(name string) (*client.Mk8s, int, error) {
	return mro.Client.GetMk8s(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (mro *Mk8sResourceOperator) InvokeUpdate(req client.Mk8s) (*client.Mk8s, int, error) {
	return mro.Client.UpdateMk8s(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (mro *Mk8sResourceOperator) InvokeDelete(name string) error {
	return mro.Client.DeleteMk8s(name)
}

// Builders //

// buildFirewall constructs a []client.Mk8sFirewallRule from the given Terraform state.
func (mro *Mk8sResourceOperator) buildFirewall(state []models.FirewallModel) *[]client.Mk8sFirewallRule {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sFirewallRule{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sFirewallRule{
			SourceCIDR:  BuildString(block.SourceCIDR),
			Description: BuildString(block.Description),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildNetworking constructs a Mk8sNetworkingConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildNetworking(state []models.NetworkingModel) *client.Mk8sNetworkingConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sNetworkingConfig{
		ServiceNetwork: BuildString(block.ServiceNetwork),
		PodNetwork:     BuildString(block.PodNetwork),
		DnsForwarder:   BuildString(block.DnsForwarder),
	}
}

// buildAutoscaler constructs a Mk8sAutoscalerConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAutoscaler(state []models.AutoscalerModel) *client.Mk8sAutoscalerConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAutoscalerConfig{
		Expander:             mro.BuildSetString(block.Expander),
		UnneededTime:         BuildString(block.UnneededTime),
		UnreadyTime:          BuildString(block.UnreadyTime),
		UtilizationThreshold: BuildFloat64(block.UtilizationThreshold),
	}
}

// buildGenericProvider constructs a Mk8sGenericProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildGenericProvider(state []models.GenericProviderModel) *client.Mk8sGenericProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sGenericProvider{
		Location:   BuildString(block.Location),
		Networking: mro.buildNetworking(block.Networking),
		NodePools:  mro.buildGenericProviderNodePools(block.NodePools),
	}
}

// buildGenericProviderNodePools constructs a []client.Mk8sGenericPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildGenericProviderNodePools(state []models.GenericProviderNodePoolModel) *[]client.Mk8sGenericPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sGenericPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sGenericPool{
			Name:   BuildString(block.Name),
			Labels: mro.BuildMapString(block.Labels),
			Taints: mro.buildGenericProviderNodePoolTaints(block.Taints),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildGenericProviderNodePoolTaints constructs a []client.Mk8sTaint from the given Terraform state.
func (mro *Mk8sResourceOperator) buildGenericProviderNodePoolTaints(state []models.GenericProviderNodePoolTaintModel) *[]client.Mk8sTaint {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sTaint{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sTaint{
			Key:    BuildString(block.Key),
			Value:  BuildString(block.Value),
			Effect: BuildString(block.Effect),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildHetznerProvider constructs a Mk8sHetznerProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildHetznerProvider(state []models.HetznerProviderModel) *client.Mk8sHetznerProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sHetznerProvider{
		Region:                   BuildString(block.Region),
		HetznerLabels:            mro.BuildMapString(block.HetznerLabels),
		Networking:               mro.buildNetworking(block.Networking),
		PreInstallScript:         BuildString(block.PreInstallScript),
		TokenSecretLink:          BuildString(block.TokenSecretLink),
		NetworkId:                BuildString(block.NetworkId),
		FirewallId:               BuildString(block.FirewallId),
		NodePools:                mro.buildHetznerProviderNodePools(block.NodePools),
		DedicatedServerNodePools: mro.buildGenericProviderNodePools(block.DedicatedServerNodePools),
		Image:                    BuildString(block.Image),
		SshKey:                   BuildString(block.SshKey),
		Autoscaler:               mro.buildAutoscaler(block.Autoscaler),
		FloatingIpSelector:       mro.BuildMapString(block.FloatingIpSelector),
	}
}

// buildHetznerProviderNodePools constructs a []client.Mk8sHetznerPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildHetznerProviderNodePools(state []models.HetznerProviderNodePoolModel) *[]client.Mk8sHetznerPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sHetznerPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sHetznerPool{
			ServerType:    BuildString(block.ServerType),
			OverrideImage: BuildString(block.OverrideImage),
			MinSize:       BuildInt(block.MinSize),
			MaxSize:       BuildInt(block.MaxSize),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildAwsProvider constructs a Mk8sAwsProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAwsProvider(state []models.AwsProviderModel) *client.Mk8sAwsProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAwsProvider{
		Region:               BuildString(block.Region),
		AwsTags:              mro.BuildMapString(block.AwsTags),
		SkipCreateRoles:      BuildBool(block.SkipCreateRoles),
		Networking:           mro.buildNetworking(block.Networking),
		PreInstallScript:     BuildString(block.PreInstallScript),
		Image:                mro.buildAwsAmi(block.Image),
		DeployRoleArn:        BuildString(block.DeployRoleArn),
		DeployRoleChain:      mro.buildAwsAssumeRoleLink(block.DeployRoleChain),
		VpcId:                BuildString(block.VpcId),
		KeyPair:              BuildString(block.KeyPair),
		DiskEncryptionKeyArn: BuildString(block.DiskEncryptionKeyArn),
		SecurityGroupIds:     mro.BuildSetString(block.SecurityGroupIds),
		ExtraNodePolicies:    mro.BuildSetString(block.ExtraNodePolicies),
		NodePools:            mro.buildAwsProviderNodePools(block.NodePools),
		Autoscaler:           mro.buildAutoscaler(block.Autoscaler),
	}
}

// buildAwsAmi constructs a Mk8sAwsAmi from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAwsAmi(state []models.AwsProviderAmiModel) *client.Mk8sAwsAmi {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAwsAmi{
		Recommended: BuildString(block.Recommended),
		Exact:       BuildString(block.Exact),
	}
}

// buildAwsAssumeRoleLink constructs a []client.Mk8sAwsAssumeRoleLink from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAwsAssumeRoleLink(state []models.AwsProviderAssumeRoleLinkModel) *[]client.Mk8sAwsAssumeRoleLink {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sAwsAssumeRoleLink{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sAwsAssumeRoleLink{
			RoleArn:           BuildString(block.RoleArn),
			ExternalId:        BuildString(block.ExternalId),
			SessionNamePrefix: BuildString(block.SessionNamePrefix),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildAwsProviderNodePools constructs a []client.Mk8sAwsPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAwsProviderNodePools(state []models.AwsProviderNodePoolModel) *[]client.Mk8sAwsPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sAwsPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sAwsPool{
			InstanceTypes:                       mro.BuildSetString(block.InstanceTypes),
			OverrideImage:                       mro.buildAwsAmi(block.OverrideImage),
			BootDiskSize:                        BuildInt(block.BootDiskSize),
			MinSize:                             BuildInt(block.MinSize),
			MaxSize:                             BuildInt(block.MaxSize),
			OnDemandBaseCapacity:                BuildInt(block.OnDemandBaseCapacity),
			OnDemandPercentageAboveBaseCapacity: BuildInt(block.OnDemandPercentageAboveBaseCapacity),
			SpotAllocationStrategy:              BuildString(block.SpotAllocationStrategy),
			SubnetIds:                           mro.BuildSetString(block.SubnetIds),
			ExtraSecurityGroupIds:               mro.BuildSetString(block.ExtraSecurityGroupIds),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildLinodeProvider constructs a Mk8sLinodeProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildLinodeProvider(state []models.LinodeProviderModel) *client.Mk8sLinodeProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sLinodeProvider{
		Region:           BuildString(block.Region),
		TokenSecretLink:  BuildString(block.TokenSecretLink),
		FirewallId:       BuildString(block.FirewallId),
		NodePools:        mro.buildLinodeProviderNodePools(block.NodePools),
		Image:            BuildString(block.Image),
		AuthorizedUsers:  mro.BuildSetString(block.AuthorizedUsers),
		AuthorizedKeys:   mro.BuildSetString(block.AuthorizedKeys),
		VpcId:            BuildString(block.VpcId),
		PreInstallScript: BuildString(block.PreInstallScript),
		Networking:       mro.buildNetworking(block.Networking),
		Autoscaler:       mro.buildAutoscaler(block.Autoscaler),
	}
}

// buildLinodeProviderNodePools constructs a []client.Mk8sLinodePool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildLinodeProviderNodePools(state []models.LinodeProviderNodePoolModel) *[]client.Mk8sLinodePool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sLinodePool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sLinodePool{
			ServerType:    BuildString(block.ServerType),
			OverrideImage: BuildString(block.OverrideImage),
			SubnetId:      BuildString(block.SubnetId),
			MinSize:       BuildInt(block.MinSize),
			MaxSize:       BuildInt(block.MaxSize),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildOblivusProvider constructs a Mk8sOblivusProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildOblivusProvider(state []models.OblivusProviderModel) *client.Mk8sOblivusProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sOblivusProvider{
		Datacenter:         BuildString(block.Datacenter),
		TokenSecretLink:    BuildString(block.TokenSecretLink),
		NodePools:          mro.buildOblivusProviderNodePools(block.NodePools),
		SshKeys:            mro.BuildSetString(block.SshKeys),
		UnmanagedNodePools: mro.buildGenericProviderNodePools(block.UnmanagedNodePool),
		Autoscaler:         mro.buildAutoscaler(block.Autoscaler),
		PreInstallScript:   BuildString(block.PreInstallScript),
	}
}

// buildOblivusProviderNodePools constructs a []client.Mk8sOblivusPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildOblivusProviderNodePools(state []models.OblivusProviderNodePoolModel) *[]client.Mk8sOblivusPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sOblivusPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sOblivusPool{
			MinSize: BuildInt(block.MinSize),
			MaxSize: BuildInt(block.MaxSize),
			Flavor:  BuildString(block.Flavor),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildLambdalabsProvider constructs a Mk8sLambdalabsProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildLambdalabsProvider(state []models.LambdalabsProviderModel) *client.Mk8sLambdalabsProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sLambdalabsProvider{
		Region:             BuildString(block.Region),
		TokenSecretLink:    BuildString(block.TokenSecretLink),
		NodePools:          mro.buildLambdalabsProviderNodePools(block.NodePools),
		SshKey:             BuildString(block.SshKey),
		UnmanagedNodePools: mro.buildGenericProviderNodePools(block.UnmanagedNodePools),
		Autoscaler:         mro.buildAutoscaler(block.Autoscaler),
		FileSystems:        mro.BuildSetString(block.FileSystems),
		PreInstallScript:   BuildString(block.PreInstallScript),
	}
}

// buildLambdalabsProviderNodePools constructs a []client.Mk8sLambdalabsPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildLambdalabsProviderNodePools(state []models.LambdalabsProviderNodePoolModel) *[]client.Mk8sLambdalabsPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sLambdalabsPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sLambdalabsPool{
			MinSize:      BuildInt(block.MinSize),
			MaxSize:      BuildInt(block.MaxSize),
			InstanceType: BuildString(block.InstanceType),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildPaperspaceProvider constructs a Mk8sPaperspaceProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildPaperspaceProvider(state []models.PaperspaceProviderModel) *client.Mk8sPaperspaceProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sPaperspaceProvider{
		Region:             BuildString(block.Region),
		TokenSecretLink:    BuildString(block.TokenSecretLink),
		SharedDrives:       mro.BuildSetString(block.SharedDrives),
		NodePools:          mro.buildPaperspaceProviderNodePools(block.NodePools),
		Autoscaler:         mro.buildAutoscaler(block.Autoscaler),
		UnmanagedNodePools: mro.buildGenericProviderNodePools(block.UnmanagedNodePools),
		PreInstallScript:   BuildString(block.PreInstallScript),
		UserIds:            mro.BuildSetString(block.UserIds),
		NetworkId:          BuildString(block.NetworkId),
	}
}

// buildPaperspaceProviderNodePools constructs a []client.Mk8sPaperspacePool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildPaperspaceProviderNodePools(state []models.PaperspaceProviderNodePoolModel) *[]client.Mk8sPaperspacePool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sPaperspacePool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sPaperspacePool{
			MinSize:      BuildInt(block.MinSize),
			MaxSize:      BuildInt(block.MaxSize),
			PublicIpType: BuildString(block.PublicIpType),
			BootDiskSize: BuildInt(block.BootDiskSize),
			MachineType:  BuildString(block.MachineType),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildEphemeralProvider constructs a Mk8sEphemeralProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildEphemeralProvider(state []models.EphemeralProviderModel) *client.Mk8sEphemeralProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sEphemeralProvider{
		Location:  BuildString(block.Location),
		NodePools: mro.buildEphemeralProviderNodePools(block.NodePools),
	}
}

// buildEphemeralProviderNodePools constructs a []client.Mk8sEphemeralPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildEphemeralProviderNodePools(state []models.EphemeralProviderNodePoolModel) *[]client.Mk8sEphemeralPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sEphemeralPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sEphemeralPool{
			Count:  BuildInt(block.Count),
			Arch:   BuildString(block.Arch),
			Flavor: BuildString(block.Flavor),
			Cpu:    BuildString(block.Cpu),
			Memory: BuildString(block.Memory),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildTritonProvider constructs a Mk8sTritonProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildTritonProvider(state []models.TritonProviderModel) *client.Mk8sTritonProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sTritonProvider{
		Connection:       mro.buildTritonProviderConnection(block.Connection),
		Networking:       mro.buildNetworking(block.Networking),
		PreInstallScript: BuildString(block.PreInstallScript),
		Location:         BuildString(block.Location),
		LoadBalancer:     mro.buildTritonProviderLoadBalancer(block.LoadBalancer),
		PrivateNetworkId: BuildString(block.PrivateNetworkId),
		FirewallEnabled:  BuildBool(block.FirewallEnabled),
		NodePools:        mro.buildTritonProviderNodePools(block.NodePools),
		ImageId:          BuildString(block.ImageId),
		SshKeys:          mro.BuildSetString(block.SshKeys),
		Autoscaler:       mro.buildAutoscaler(block.Autoscaler),
	}
}

// buildTritonProviderConnection constructs a Mk8sTritonConnection from the given Terraform state.
func (mro *Mk8sResourceOperator) buildTritonProviderConnection(state []models.TritonProviderConnectionModel) *client.Mk8sTritonConnection {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sTritonConnection{
		Url:                  BuildString(block.Url),
		Account:              BuildString(block.Account),
		User:                 BuildString(block.User),
		PrivateKeySecretLink: BuildString(block.PrivateKeySecretLink),
	}
}

// buildTritonProviderLoadBalancer constructs a Mk8sTritonLoadBalancer from the given Terraform state.
func (mro *Mk8sResourceOperator) buildTritonProviderLoadBalancer(state []models.TritonProviderLoadBalancerModel) *client.Mk8sTritonLoadBalancer {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sTritonLoadBalancer{
		Manual:  mro.buildTritonProviderLoadBalancerManual(block.Manual),
		Gateway: mro.buildTritonProviderLoadBalancerGateway(block.Gateway),
	}
}

// buildTritonProviderLoadBalancerManual constructs a Mk8sTritonManual from the given Terraform state.
func (mro *Mk8sResourceOperator) buildTritonProviderLoadBalancerManual(state []models.TritonProviderLoadBalancerManualModel) *client.Mk8sTritonManual {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sTritonManual{
		PackageId:         BuildString(block.PackageId),
		ImageId:           BuildString(block.ImageId),
		PublicNetworkId:   BuildString(block.PublicNetworkId),
		PrivateNetworkIds: mro.BuildSetString(block.PrivateNetworkIds),
		Metadata:          mro.BuildMapString(block.Metadata),
		Tags:              mro.BuildMapString(block.Tags),
		Count:             BuildInt(block.Count),
		CnsInternalDomain: BuildString(block.CnsInternalDomain),
		CnsPublicDomain:   BuildString(block.CnsPublicDomain),
	}
}

// buildTritonProviderLoadBalancerGateway constructs a Mk8sTritonGateway from the given Terraform state.
func (mro *Mk8sResourceOperator) buildTritonProviderLoadBalancerGateway(state []models.TritonProviderLoadBalancerGatewayModel) *client.Mk8sTritonGateway {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Construct and return the output
	return &client.Mk8sTritonGateway{}
}

// buildTritonProviderNodePools constructs a []client.Mk8sTritonPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildTritonProviderNodePools(state []models.TritonProviderNodePoolModel) *[]client.Mk8sTritonPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sTritonPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sTritonPool{
			PackageId:         BuildString(block.PackageId),
			OverrideImageId:   BuildString(block.OverrideImageId),
			PublicNetworkId:   BuildString(block.PublicNetworkId),
			PrivateNetworkIds: mro.BuildSetString(block.PrivateNetworkIds),
			TritonTags:        mro.BuildMapString(block.TritonTags),
			MinSize:           BuildInt(block.MinSize),
			MaxSize:           BuildInt(block.MaxSize),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildAzureProvider constructs a Mk8sAzureProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAzureProvider(state []models.AzureProviderModel) *client.Mk8sAzureProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAzureProvider{
		Location:         BuildString(block.Location),
		SubscriptionId:   BuildString(block.SubscriptionId),
		SdkSecretLink:    BuildString(block.SdkSecretLink),
		ResourceGroup:    BuildString(block.ResourceGroup),
		Networking:       mro.buildNetworking(block.Networking),
		PreInstallScript: BuildString(block.PreInstallScript),
		Image:            mro.buildAzureProviderImage(block.Image),
		SshKeys:          mro.BuildSetString(block.SshKeys),
		NetworkId:        BuildString(block.NetworkId),
		Tags:             mro.BuildMapString(block.Tags),
		NodePools:        mro.buildAzureProviderNodePools(block.NodePools),
		Autoscaler:       mro.buildAutoscaler(block.Autoscaler),
	}
}

// buildAzureProviderImage constructs a Mk8sAzureImage from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAzureProviderImage(state []models.AzureProviderImageModel) *client.Mk8sAzureImage {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAzureImage{
		Recommended: BuildString(block.Recommended),
		Reference:   mro.buildAzureProviderImageReference(block.Reference),
	}
}

// buildAzureProviderImageReference constructs a Mk8sAzureImageReference from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAzureProviderImageReference(state []models.AzureProviderImageReferenceModel) *client.Mk8sAzureImageReference {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAzureImageReference{
		Publisher: BuildString(block.Publisher),
		Offer:     BuildString(block.Offer),
		Sku:       BuildString(block.Sku),
		Version:   BuildString(block.Version),
	}
}

// buildAzureProviderNodePools constructs a []client.Mk8sAzurePool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAzureProviderNodePools(state []models.AzureProviderNodePoolModel) *[]client.Mk8sAzurePool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sAzurePool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sAzurePool{
			Size:          BuildString(block.Size),
			SubnetId:      BuildString(block.SubnetId),
			Zones:         mro.BuildSetInt(block.Zones),
			OverrideImage: mro.buildAzureProviderImage(block.OverrideImage),
			BootDiskSize:  BuildInt(block.BootDiskSize),
			MinSize:       BuildInt(block.MinSize),
			MaxSize:       BuildInt(block.MaxSize),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildDigitalOceanProvider constructs a Mk8sDigitalOceanProvider from the given Terraform state.
func (mro *Mk8sResourceOperator) buildDigitalOceanProvider(state []models.DigitalOceanProviderModel) *client.Mk8sDigitalOceanProvider {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sDigitalOceanProvider{
		Region:           BuildString(block.Region),
		DigitalOceanTags: mro.BuildSetString(block.DigitalOceanTags),
		Networking:       mro.buildNetworking(block.Networking),
		PreInstallScript: BuildString(block.PreInstallScript),
		TokenSecretLink:  BuildString(block.TokenSecretLink),
		VpcId:            BuildString(block.VpcId),
		NodePools:        mro.buildDigitalOceanProviderNodePools(block.NodePools),
		Image:            BuildString(block.Image),
		SshKeys:          mro.BuildSetString(block.SshKeys),
		ExtraSshKeys:     mro.BuildSetString(block.ExtraSshKeys),
		Autoscaler:       mro.buildAutoscaler(block.Autoscaler),
		ReservedIps:      mro.BuildSetString(block.ReservedIps),
	}
}

// buildDigitalOceanProviderNodePools constructs a []client.Mk8sDigitalOceanPool from the given Terraform state.
func (mro *Mk8sResourceOperator) buildDigitalOceanProviderNodePools(state []models.DigitalOceanProviderNodePoolModel) *[]client.Mk8sDigitalOceanPool {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Mk8sDigitalOceanPool{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.Mk8sDigitalOceanPool{
			DropletSize:   BuildString(block.DropletSize),
			OverrideImage: BuildString(block.OverrideImage),
			MinSize:       BuildInt(block.MinSize),
			MaxSize:       BuildInt(block.MaxSize),
		}

		// Set embedded attributes
		item.Name = BuildString(block.Name)
		item.Labels = mro.BuildMapString(block.Labels)
		item.Taints = mro.buildGenericProviderNodePoolTaints(block.Taints)

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildAddOns constructs a Mk8sSpecAddOns from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOns(state []models.AddOnsModel) *client.Mk8sSpecAddOns {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sSpecAddOns{
		Dashboard:             mro.buildAddOnConfig(block.Dashboard),
		AzureWorkloadIdentity: mro.buildAddOnAzureWorkloadIdentity(block.AzureWorkloadIdentity),
		AwsWorkloadIdentity:   mro.buildAddOnConfig(block.AwsWorkloadIdentity),
		LocalPathStorage:      mro.buildAddOnConfig(block.LocalPathStorage),
		Metrics:               mro.buildAddOnMetrics(block.Metrics),
		Logs:                  mro.buildAddOnLogs(block.Logs),
		Nvidia:                mro.buildAddOnNvidia(block.Nvidia),
		AwsEFS:                mro.buildAddOnAwsConfig(block.AwsEFS),
		AwsECR:                mro.buildAddOnAwsConfig(block.AwsECR),
		AwsELB:                mro.buildAddOnAwsConfig(block.AwsELB),
		AzureACR:              mro.buildAddOnAzureAcr(block.AzureACR),
		Sysbox:                mro.buildAddOnConfig(block.Sysbox),
	}
}

// buildAddOnConfig builds a non-customizable addon configuration based on the provided state.
func (mro *Mk8sResourceOperator) buildAddOnConfig(state types.Bool) *client.Mk8sNonCustomizableAddonConfig {
	// Convert the Terraform bool value to a Go *bool
	isEnabled := BuildBool(state)

	// If the AddOn flag exists and is true, return a new config to enable the AddOn
	if isEnabled != nil && *isEnabled {
		return &client.Mk8sNonCustomizableAddonConfig{}
	}

	// Return nil to indicate the AddOn is disabled
	return nil
}

// buildAddOnAzureWorkloadIdentity constructs a Mk8sAzureWorkloadIdentityAddOnConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOnAzureWorkloadIdentity(state []models.AddOnAzureWorkloadIdentityModel) *client.Mk8sAzureWorkloadIdentityAddOnConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAzureWorkloadIdentityAddOnConfig{
		TenantId: BuildString(block.TenantId),
	}
}

// buildAddOnMetrics constructs a Mk8sMetricsAddOnConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOnMetrics(state []models.AddOnsMetricsModel) *client.Mk8sMetricsAddOnConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sMetricsAddOnConfig{
		KubeState:       BuildBool(block.KubeState),
		CoreDns:         BuildBool(block.CoreDns),
		Kubelet:         BuildBool(block.Kubelet),
		Apiserver:       BuildBool(block.Apiserver),
		NodeExporter:    BuildBool(block.NodeExporter),
		Cadvisor:        BuildBool(block.Cadvisor),
		ScrapeAnnotated: mro.buildAddOnMetricsScrapeAnnotated(block.ScrapeAnnotated),
	}
}

// buildAddOnMetricsScrapeAnnotated constructs a Mk8sMetricsScrapeAnnotated from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOnMetricsScrapeAnnotated(state []models.AddOnsMetricsScrapeAnnotatedModel) *client.Mk8sMetricsScrapeAnnotated {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sMetricsScrapeAnnotated{
		IntervalSeconds:   BuildInt(block.IntervalSeconds),
		IncludeNamespaces: BuildString(block.IncludeNamespaces),
		ExcludeNamespaces: BuildString(block.ExcludeNamespaces),
		RetainLabels:      BuildString(block.RetainLabels),
	}
}

// buildAddOnLogs constructs a Mk8sLogsAddOnConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOnLogs(state []models.AddOnsLogsModel) *client.Mk8sLogsAddOnConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sLogsAddOnConfig{
		AuditEnabled:      BuildBool(block.AuditEnabled),
		IncludeNamespaces: BuildString(block.IncludeNamespaaces),
		ExcludeNamespaces: BuildString(block.ExcludeNamespaces),
		Docker:            BuildBool(block.Docker),
		Kubelet:           BuildBool(block.Kubelet),
		Kernel:            BuildBool(block.Kernel),
		Events:            BuildBool(block.Events),
	}
}

// buildAddOnNvidia constructs a Mk8sNvidiaAddOnConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOnNvidia(state []models.AddOnsNvidiaModel) *client.Mk8sNvidiaAddOnConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sNvidiaAddOnConfig{
		TaintGPUNodes: BuildBool(block.TaintGpuNodes),
	}
}

// buildAddOnAwsConfig constructs a Mk8sAwsAddOnConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOnAwsConfig(state []models.AddOnsHasRoleArnModel) *client.Mk8sAwsAddOnConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAwsAddOnConfig{
		RoleArn: BuildString(block.RoleArn),
	}
}

// buildAddOnAzureAcr constructs a Mk8sAzureACRAddOnConfig from the given Terraform state.
func (mro *Mk8sResourceOperator) buildAddOnAzureAcr(state []models.AddOnsAzureAcrModel) *client.Mk8sAzureACRAddOnConfig {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.Mk8sAzureACRAddOnConfig{
		ClientId: BuildString(block.ClientId),
	}
}

// Flatteners //

// flattenFirewall transforms *[]client.Mk8sFirewallRule into a []models.FirewallModel.
func (mro *Mk8sResourceOperator) flattenFirewall(input *[]client.Mk8sFirewallRule) []models.FirewallModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.FirewallModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.FirewallModel{
			SourceCIDR:  types.StringPointerValue(item.SourceCIDR),
			Description: types.StringPointerValue(item.Description),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenNetworking transforms *client.Mk8sNetworkingConfig into a []models.NetworkingModel.
func (mro *Mk8sResourceOperator) flattenNetworking(input *client.Mk8sNetworkingConfig) []models.NetworkingModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.NetworkingModel{
		ServiceNetwork: types.StringPointerValue(input.ServiceNetwork),
		PodNetwork:     types.StringPointerValue(input.PodNetwork),
		DnsForwarder:   types.StringPointerValue(input.DnsForwarder),
	}

	// Return a slice containing the single block
	return []models.NetworkingModel{block}
}

// flattenAutoscaler transforms *client.Mk8sAutoscalerConfig into a []models.AutoscalerModel.
func (mro *Mk8sResourceOperator) flattenAutoscaler(input *client.Mk8sAutoscalerConfig) []models.AutoscalerModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AutoscalerModel{
		Expander:             FlattenSetString(input.Expander),
		UnneededTime:         types.StringPointerValue(input.UnneededTime),
		UnreadyTime:          types.StringPointerValue(input.UnreadyTime),
		UtilizationThreshold: FlattenFloat64(input.UtilizationThreshold),
	}

	// Return a slice containing the single block
	return []models.AutoscalerModel{block}
}

// flattenGenericProvider transforms *client.Mk8sGenericProvider into a []models.GenericProviderModel.
func (mro *Mk8sResourceOperator) flattenGenericProvider(input *client.Mk8sGenericProvider) []models.GenericProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.GenericProviderModel{
		Location:   types.StringPointerValue(input.Location),
		Networking: mro.flattenNetworking(input.Networking),
		NodePools:  mro.flattenGenericProviderNodePools(input.NodePools),
	}

	// Return a slice containing the single block
	return []models.GenericProviderModel{block}
}

// flattenGenericProviderNodePools transforms *[]client.Mk8sGenericPool into a []models.GenericProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenGenericProviderNodePools(input *[]client.Mk8sGenericPool) []models.GenericProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.GenericProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.GenericProviderNodePoolModel{
			Name:   types.StringPointerValue(item.Name),
			Labels: FlattenMapString(item.Labels),
			Taints: mro.flattenGenericProviderNodePoolTaints(item.Taints),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenGenericProviderNodePoolTaints transforms *[]client.Mk8sTaint into a []models.GenericProviderNodePoolTaintModel.
func (mro *Mk8sResourceOperator) flattenGenericProviderNodePoolTaints(input *[]client.Mk8sTaint) []models.GenericProviderNodePoolTaintModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.GenericProviderNodePoolTaintModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.GenericProviderNodePoolTaintModel{
			Key:    types.StringPointerValue(item.Key),
			Value:  types.StringPointerValue(item.Value),
			Effect: types.StringPointerValue(item.Effect),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenHetznerProvider transforms *client.Mk8sHetznerProvider into a []models.HetznerProviderModel.
func (mro *Mk8sResourceOperator) flattenHetznerProvider(input *client.Mk8sHetznerProvider) []models.HetznerProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.HetznerProviderModel{
		Region:                   types.StringPointerValue(input.Region),
		HetznerLabels:            FlattenMapString(input.HetznerLabels),
		Networking:               mro.flattenNetworking(input.Networking),
		PreInstallScript:         types.StringPointerValue(input.PreInstallScript),
		TokenSecretLink:          types.StringPointerValue(input.TokenSecretLink),
		NetworkId:                types.StringPointerValue(input.NetworkId),
		FirewallId:               types.StringPointerValue(input.FirewallId),
		NodePools:                mro.flattenHetznerProviderNodePools(input.NodePools),
		DedicatedServerNodePools: mro.flattenGenericProviderNodePools(input.DedicatedServerNodePools),
		Image:                    types.StringPointerValue(input.Image),
		SshKey:                   types.StringPointerValue(input.SshKey),
		Autoscaler:               mro.flattenAutoscaler(input.Autoscaler),
		FloatingIpSelector:       FlattenMapString(input.FloatingIpSelector),
	}

	// Return a slice containing the single block
	return []models.HetznerProviderModel{block}
}

// flattenHetznerProviderNodePools transforms *[]client.Mk8sHetznerPool into a []models.HetznerProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenHetznerProviderNodePools(input *[]client.Mk8sHetznerPool) []models.HetznerProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.HetznerProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.HetznerProviderNodePoolModel{
			ServerType:    types.StringPointerValue(item.ServerType),
			OverrideImage: types.StringPointerValue(item.OverrideImage),
			MinSize:       FlattenInt(item.MinSize),
			MaxSize:       FlattenInt(item.MaxSize),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenAwsProvider transforms *client.Mk8sAwsProvider into a []models.AwsProviderModel.
func (mro *Mk8sResourceOperator) flattenAwsProvider(input *client.Mk8sAwsProvider) []models.AwsProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AwsProviderModel{
		Region:               types.StringPointerValue(input.Region),
		AwsTags:              FlattenMapString(input.AwsTags),
		SkipCreateRoles:      types.BoolPointerValue(input.SkipCreateRoles),
		Networking:           mro.flattenNetworking(input.Networking),
		PreInstallScript:     types.StringPointerValue(input.PreInstallScript),
		Image:                mro.flattenAwsAmi(input.Image),
		DeployRoleArn:        types.StringPointerValue(input.DeployRoleArn),
		DeployRoleChain:      mro.flattenAwsAssumeRoleLink(input.DeployRoleChain),
		VpcId:                types.StringPointerValue(input.VpcId),
		KeyPair:              types.StringPointerValue(input.KeyPair),
		DiskEncryptionKeyArn: types.StringPointerValue(input.DiskEncryptionKeyArn),
		SecurityGroupIds:     FlattenSetString(input.SecurityGroupIds),
		ExtraNodePolicies:    FlattenSetString(input.ExtraNodePolicies),
		NodePools:            mro.flattenAwsProviderNodePools(input.NodePools),
		Autoscaler:           mro.flattenAutoscaler(input.Autoscaler),
	}

	// Return a slice containing the single block
	return []models.AwsProviderModel{block}
}

// flattenAwsAmi transforms *client.Mk8sAwsAmi into a []models.AwsProviderAmiModel.
func (mro *Mk8sResourceOperator) flattenAwsAmi(input *client.Mk8sAwsAmi) []models.AwsProviderAmiModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AwsProviderAmiModel{
		Recommended: types.StringPointerValue(input.Recommended),
		Exact:       types.StringPointerValue(input.Exact),
	}

	// Return a slice containing the single block
	return []models.AwsProviderAmiModel{block}
}

// flattenAwsAssumeRoleLink transforms *[]client.Mk8sAwsAssumeRoleLink into a []models.AwsProviderAssumeRoleLinkModel.
func (mro *Mk8sResourceOperator) flattenAwsAssumeRoleLink(input *[]client.Mk8sAwsAssumeRoleLink) []models.AwsProviderAssumeRoleLinkModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.AwsProviderAssumeRoleLinkModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.AwsProviderAssumeRoleLinkModel{
			RoleArn:           types.StringPointerValue(item.RoleArn),
			ExternalId:        types.StringPointerValue(item.ExternalId),
			SessionNamePrefix: types.StringPointerValue(item.SessionNamePrefix),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenAwsProviderNodePools transforms *[]client.Mk8sAwsPool into a []models.AwsProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenAwsProviderNodePools(input *[]client.Mk8sAwsPool) []models.AwsProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.AwsProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.AwsProviderNodePoolModel{
			InstanceTypes:                       FlattenSetString(item.InstanceTypes),
			OverrideImage:                       mro.flattenAwsAmi(item.OverrideImage),
			BootDiskSize:                        FlattenInt(item.BootDiskSize),
			MinSize:                             FlattenInt(item.MinSize),
			MaxSize:                             FlattenInt(item.MaxSize),
			OnDemandBaseCapacity:                FlattenInt(item.OnDemandBaseCapacity),
			OnDemandPercentageAboveBaseCapacity: FlattenInt(item.OnDemandPercentageAboveBaseCapacity),
			SpotAllocationStrategy:              types.StringPointerValue(item.SpotAllocationStrategy),
			SubnetIds:                           FlattenSetString(item.SubnetIds),
			ExtraSecurityGroupIds:               FlattenSetString(item.ExtraSecurityGroupIds),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenLinodeProvider transforms *client.Mk8sLinodeProvider into a []models.LinodeProviderModel.
func (mro *Mk8sResourceOperator) flattenLinodeProvider(input *client.Mk8sLinodeProvider) []models.LinodeProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.LinodeProviderModel{
		Region:           types.StringPointerValue(input.Region),
		TokenSecretLink:  types.StringPointerValue(input.TokenSecretLink),
		FirewallId:       types.StringPointerValue(input.FirewallId),
		NodePools:        mro.flattenLinodeProviderNodePools(input.NodePools),
		Image:            types.StringPointerValue(input.Image),
		AuthorizedUsers:  FlattenSetString(input.AuthorizedUsers),
		AuthorizedKeys:   FlattenSetString(input.AuthorizedKeys),
		VpcId:            types.StringPointerValue(input.VpcId),
		PreInstallScript: types.StringPointerValue(input.PreInstallScript),
		Networking:       mro.flattenNetworking(input.Networking),
		Autoscaler:       mro.flattenAutoscaler(input.Autoscaler),
	}

	// Return a slice containing the single block
	return []models.LinodeProviderModel{block}
}

// flattenLinodeProviderNodePools transforms *[]client.Mk8sLinodePool into a []models.LinodeProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenLinodeProviderNodePools(input *[]client.Mk8sLinodePool) []models.LinodeProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.LinodeProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.LinodeProviderNodePoolModel{
			ServerType:    types.StringPointerValue(item.ServerType),
			OverrideImage: types.StringPointerValue(item.OverrideImage),
			SubnetId:      types.StringPointerValue(item.SubnetId),
			MinSize:       FlattenInt(item.MinSize),
			MaxSize:       FlattenInt(item.MaxSize),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenOblivusProvider transforms *client.Mk8sOblivusProvider into a []models.OblivusProviderModel.
func (mro *Mk8sResourceOperator) flattenOblivusProvider(input *client.Mk8sOblivusProvider) []models.OblivusProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.OblivusProviderModel{
		Datacenter:        types.StringPointerValue(input.Datacenter),
		TokenSecretLink:   types.StringPointerValue(input.TokenSecretLink),
		NodePools:         mro.flattenOblivusProviderNodePools(input.NodePools),
		SshKeys:           FlattenSetString(input.SshKeys),
		UnmanagedNodePool: mro.flattenGenericProviderNodePools(input.UnmanagedNodePools),
		Autoscaler:        mro.flattenAutoscaler(input.Autoscaler),
		PreInstallScript:  types.StringPointerValue(input.PreInstallScript),
	}

	// Return a slice containing the single block
	return []models.OblivusProviderModel{block}
}

// flattenOblivusProviderNodePools transforms *[]client.Mk8sOblivusPool into a []models.OblivusProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenOblivusProviderNodePools(input *[]client.Mk8sOblivusPool) []models.OblivusProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.OblivusProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.OblivusProviderNodePoolModel{
			MinSize: FlattenInt(item.MinSize),
			MaxSize: FlattenInt(item.MaxSize),
			Flavor:  types.StringPointerValue(item.Flavor),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenLambdalabsProvider transforms *client.Mk8sLambdalabsProvider into a []models.LambdalabsProviderModel.
func (mro *Mk8sResourceOperator) flattenLambdalabsProvider(input *client.Mk8sLambdalabsProvider) []models.LambdalabsProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.LambdalabsProviderModel{
		Region:             types.StringPointerValue(input.Region),
		TokenSecretLink:    types.StringPointerValue(input.TokenSecretLink),
		NodePools:          mro.flattenLambdalabsProviderNodePools(input.NodePools),
		SshKey:             types.StringPointerValue(input.SshKey),
		UnmanagedNodePools: mro.flattenGenericProviderNodePools(input.UnmanagedNodePools),
		Autoscaler:         mro.flattenAutoscaler(input.Autoscaler),
		FileSystems:        FlattenSetString(input.FileSystems),
		PreInstallScript:   types.StringPointerValue(input.PreInstallScript),
	}

	// Return a slice containing the single block
	return []models.LambdalabsProviderModel{block}
}

// flattenLambdalabsProviderNodePools transforms *[]client.Mk8sLambdalabsPool into a []models.LambdalabsProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenLambdalabsProviderNodePools(input *[]client.Mk8sLambdalabsPool) []models.LambdalabsProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.LambdalabsProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.LambdalabsProviderNodePoolModel{
			MinSize:      FlattenInt(item.MinSize),
			MaxSize:      FlattenInt(item.MaxSize),
			InstanceType: types.StringPointerValue(item.InstanceType),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenPaperspaceProvider transforms *client.Mk8sPaperspaceProvider into a []models.PaperspaceProviderModel.
func (mro *Mk8sResourceOperator) flattenPaperspaceProvider(input *client.Mk8sPaperspaceProvider) []models.PaperspaceProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.PaperspaceProviderModel{
		Region:             types.StringPointerValue(input.Region),
		TokenSecretLink:    types.StringPointerValue(input.TokenSecretLink),
		SharedDrives:       FlattenSetString(input.SharedDrives),
		NodePools:          mro.flattenPaperspaceProviderNodePools(input.NodePools),
		Autoscaler:         mro.flattenAutoscaler(input.Autoscaler),
		UnmanagedNodePools: mro.flattenGenericProviderNodePools(input.UnmanagedNodePools),
		PreInstallScript:   types.StringPointerValue(input.PreInstallScript),
		UserIds:            FlattenSetString(input.UserIds),
		NetworkId:          types.StringPointerValue(input.NetworkId),
	}

	// Return a slice containing the single block
	return []models.PaperspaceProviderModel{block}
}

// flattenPaperspaceProviderNodePools transforms *[]client.Mk8sPaperspacePool into a []models.PaperspaceProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenPaperspaceProviderNodePools(input *[]client.Mk8sPaperspacePool) []models.PaperspaceProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.PaperspaceProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.PaperspaceProviderNodePoolModel{
			MinSize:      FlattenInt(item.MinSize),
			MaxSize:      FlattenInt(item.MaxSize),
			PublicIpType: types.StringPointerValue(item.PublicIpType),
			BootDiskSize: FlattenInt(item.BootDiskSize),
			MachineType:  types.StringPointerValue(item.MachineType),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenEphemeralProvider transforms *client.Mk8sEphemeralProvider into a []models.EphemeralProviderModel.
func (mro *Mk8sResourceOperator) flattenEphemeralProvider(input *client.Mk8sEphemeralProvider) []models.EphemeralProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.EphemeralProviderModel{
		Location:  types.StringPointerValue(input.Location),
		NodePools: mro.flattenEphemeralProviderNodePools(input.NodePools),
	}

	// Return a slice containing the single block
	return []models.EphemeralProviderModel{block}
}

// flattenEphemeralProviderNodePools transforms *[]client.Mk8sEphemeralPool into a []models.EphemeralProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenEphemeralProviderNodePools(input *[]client.Mk8sEphemeralPool) []models.EphemeralProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.EphemeralProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.EphemeralProviderNodePoolModel{
			Count:  FlattenInt(item.Count),
			Arch:   types.StringPointerValue(item.Arch),
			Flavor: types.StringPointerValue(item.Flavor),
			Cpu:    types.StringPointerValue(item.Cpu),
			Memory: types.StringPointerValue(item.Memory),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenTritonProvider transforms *client.Mk8sTritonProvider into a []models.TritonProviderModel.
func (mro *Mk8sResourceOperator) flattenTritonProvider(input *client.Mk8sTritonProvider) []models.TritonProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.TritonProviderModel{
		Connection:       mro.flattenTritonProviderConnection(input.Connection),
		Networking:       mro.flattenNetworking(input.Networking),
		PreInstallScript: types.StringPointerValue(input.PreInstallScript),
		Location:         types.StringPointerValue(input.Location),
		LoadBalancer:     mro.flattenTritonProviderLoadBalancer(input.LoadBalancer),
		PrivateNetworkId: types.StringPointerValue(input.PrivateNetworkId),
		FirewallEnabled:  types.BoolPointerValue(input.FirewallEnabled),
		NodePools:        mro.flattenTritonProviderNodePools(input.NodePools),
		ImageId:          types.StringPointerValue(input.ImageId),
		SshKeys:          FlattenSetString(input.SshKeys),
		Autoscaler:       mro.flattenAutoscaler(input.Autoscaler),
	}

	// Return a slice containing the single block
	return []models.TritonProviderModel{block}
}

// flattenTritonProviderConnection transforms *client.Mk8sTritonConnection into a []models.TritonProviderConnectionModel.
func (mro *Mk8sResourceOperator) flattenTritonProviderConnection(input *client.Mk8sTritonConnection) []models.TritonProviderConnectionModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.TritonProviderConnectionModel{
		Url:                  types.StringPointerValue(input.Url),
		Account:              types.StringPointerValue(input.Account),
		User:                 types.StringPointerValue(input.User),
		PrivateKeySecretLink: types.StringPointerValue(input.PrivateKeySecretLink),
	}

	// Return a slice containing the single block
	return []models.TritonProviderConnectionModel{block}
}

// flattenTritonProviderLoadBalancer transforms *client.Mk8sTritonLoadBalancer into a []models.TritonProviderLoadBalancerModel.
func (mro *Mk8sResourceOperator) flattenTritonProviderLoadBalancer(input *client.Mk8sTritonLoadBalancer) []models.TritonProviderLoadBalancerModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.TritonProviderLoadBalancerModel{
		Manual:  mro.flattenTritonProviderLoadBalancerManual(input.Manual),
		Gateway: mro.flattenTritonProviderLoadBalancerGateway(input.Gateway),
	}

	// Return a slice containing the single block
	return []models.TritonProviderLoadBalancerModel{block}
}

// flattenTritonProviderLoadBalancerManual transforms *client.Mk8sTritonManual into a []models.TritonProviderLoadBalancerManualModel.
func (mro *Mk8sResourceOperator) flattenTritonProviderLoadBalancerManual(input *client.Mk8sTritonManual) []models.TritonProviderLoadBalancerManualModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.TritonProviderLoadBalancerManualModel{
		PackageId:         types.StringPointerValue(input.PackageId),
		ImageId:           types.StringPointerValue(input.ImageId),
		PublicNetworkId:   types.StringPointerValue(input.PublicNetworkId),
		PrivateNetworkIds: FlattenSetString(input.PrivateNetworkIds),
		Metadata:          FlattenMapString(input.Metadata),
		Tags:              FlattenMapString(input.Tags),
		Count:             FlattenInt(input.Count),
		CnsInternalDomain: types.StringPointerValue(input.CnsInternalDomain),
		CnsPublicDomain:   types.StringPointerValue(input.CnsPublicDomain),
	}

	// Return a slice containing the single block
	return []models.TritonProviderLoadBalancerManualModel{block}
}

// flattenTritonProviderLoadBalancerGateway transforms *client.Mk8sTritonGateway into a []models.TritonProviderLoadBalancerGatewayModel.
func (mro *Mk8sResourceOperator) flattenTritonProviderLoadBalancerGateway(input *client.Mk8sTritonGateway) []models.TritonProviderLoadBalancerGatewayModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Return a slice containing the single block
	return []models.TritonProviderLoadBalancerGatewayModel{{}}
}

// flattenTritonProviderNodePools transforms *[]client.Mk8sTritonPool into a []models.TritonProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenTritonProviderNodePools(input *[]client.Mk8sTritonPool) []models.TritonProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.TritonProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.TritonProviderNodePoolModel{
			PackageId:         types.StringPointerValue(item.PackageId),
			OverrideImageId:   types.StringPointerValue(item.OverrideImageId),
			PublicNetworkId:   types.StringPointerValue(item.PublicNetworkId),
			PrivateNetworkIds: FlattenSetString(item.PrivateNetworkIds),
			TritonTags:        FlattenMapString(item.TritonTags),
			MinSize:           FlattenInt(item.MinSize),
			MaxSize:           FlattenInt(item.MaxSize),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenAzureProvider transforms *client.Mk8sAzureProvider into a []models.AzureProviderModel.
func (mro *Mk8sResourceOperator) flattenAzureProvider(input *client.Mk8sAzureProvider) []models.AzureProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AzureProviderModel{
		Location:         types.StringPointerValue(input.Location),
		SubscriptionId:   types.StringPointerValue(input.SubscriptionId),
		SdkSecretLink:    types.StringPointerValue(input.SdkSecretLink),
		ResourceGroup:    types.StringPointerValue(input.ResourceGroup),
		Networking:       mro.flattenNetworking(input.Networking),
		PreInstallScript: types.StringPointerValue(input.PreInstallScript),
		Image:            mro.flattenAzureProviderImage(input.Image),
		SshKeys:          FlattenSetString(input.SshKeys),
		NetworkId:        types.StringPointerValue(input.NetworkId),
		Tags:             FlattenMapString(input.Tags),
		NodePools:        mro.flattenAzureProviderNodePools(input.NodePools),
		Autoscaler:       mro.flattenAutoscaler(input.Autoscaler),
	}

	// Return a slice containing the single block
	return []models.AzureProviderModel{block}
}

// flattenAzureProviderImage transforms *client.Mk8sAzureImage into a []models.AzureProviderImageModel.
func (mro *Mk8sResourceOperator) flattenAzureProviderImage(input *client.Mk8sAzureImage) []models.AzureProviderImageModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AzureProviderImageModel{
		Recommended: types.StringPointerValue(input.Recommended),
		Reference:   mro.flattenAzureProviderImageReference(input.Reference),
	}

	// Return a slice containing the single block
	return []models.AzureProviderImageModel{block}
}

// flattenAzureProviderImageReference transforms *client.Mk8sAzureImageReference into a []models.AzureProviderImageReferenceModel.
func (mro *Mk8sResourceOperator) flattenAzureProviderImageReference(input *client.Mk8sAzureImageReference) []models.AzureProviderImageReferenceModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AzureProviderImageReferenceModel{
		Publisher: types.StringPointerValue(input.Publisher),
		Offer:     types.StringPointerValue(input.Offer),
		Sku:       types.StringPointerValue(input.Sku),
		Version:   types.StringPointerValue(input.Version),
	}

	// Return a slice containing the single block
	return []models.AzureProviderImageReferenceModel{block}
}

// flattenAzureProviderNodePools transforms *[]client.Mk8sAzurePool into a []models.AzureProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenAzureProviderNodePools(input *[]client.Mk8sAzurePool) []models.AzureProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.AzureProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.AzureProviderNodePoolModel{
			Size:          types.StringPointerValue(item.Size),
			SubnetId:      types.StringPointerValue(item.SubnetId),
			Zones:         FlattenSetInt(item.Zones),
			OverrideImage: mro.flattenAzureProviderImage(item.OverrideImage),
			BootDiskSize:  FlattenInt(item.BootDiskSize),
			MinSize:       FlattenInt(item.MinSize),
			MaxSize:       FlattenInt(item.MaxSize),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenDigitalOceanProvider transforms *client.Mk8sDigitalOceanProvider into a []models.DigitalOceanProviderModel.
func (mro *Mk8sResourceOperator) flattenDigitalOceanProvider(input *client.Mk8sDigitalOceanProvider) []models.DigitalOceanProviderModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.DigitalOceanProviderModel{
		Region:           types.StringPointerValue(input.Region),
		DigitalOceanTags: FlattenSetString(input.DigitalOceanTags),
		Networking:       mro.flattenNetworking(input.Networking),
		PreInstallScript: types.StringPointerValue(input.PreInstallScript),
		TokenSecretLink:  types.StringPointerValue(input.TokenSecretLink),
		VpcId:            types.StringPointerValue(input.VpcId),
		NodePools:        mro.flattenDigitalOceanProviderNodePools(input.NodePools),
		Image:            types.StringPointerValue(input.Image),
		SshKeys:          FlattenSetString(input.SshKeys),
		ExtraSshKeys:     FlattenSetString(input.ExtraSshKeys),
		Autoscaler:       mro.flattenAutoscaler(input.Autoscaler),
		ReservedIps:      FlattenSetString(input.ReservedIps),
	}

	// Return a slice containing the single block
	return []models.DigitalOceanProviderModel{block}
}

// flattenDigitalOceanProviderNodePools transforms *[]client.Mk8sDigitalOceanPool into a []models.DigitalOceanProviderNodePoolModel.
func (mro *Mk8sResourceOperator) flattenDigitalOceanProviderNodePools(input *[]client.Mk8sDigitalOceanPool) []models.DigitalOceanProviderNodePoolModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.DigitalOceanProviderNodePoolModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.DigitalOceanProviderNodePoolModel{
			DropletSize:   types.StringPointerValue(item.DropletSize),
			OverrideImage: types.StringPointerValue(item.OverrideImage),
			MinSize:       FlattenInt(item.MinSize),
			MaxSize:       FlattenInt(item.MaxSize),
		}

		// Set embedded attributes
		block.Name = types.StringPointerValue(item.Name)
		block.Labels = FlattenMapString(item.Labels)
		block.Taints = mro.flattenGenericProviderNodePoolTaints(item.Taints)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenAddOns transforms *client.Mk8sSpecAddOns into a []models.AddOnsModel.
func (mro *Mk8sResourceOperator) flattenAddOns(input *client.Mk8sSpecAddOns) []models.AddOnsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnsModel{
		Dashboard:             mro.flattenAddOnConfig(input.Dashboard),
		AzureWorkloadIdentity: mro.flattenAddOnAzureWorkloadIdentity(input.AzureWorkloadIdentity),
		AwsWorkloadIdentity:   mro.flattenAddOnConfig(input.AwsWorkloadIdentity),
		LocalPathStorage:      mro.flattenAddOnConfig(input.LocalPathStorage),
		Metrics:               mro.flattenAddOnMetrics(input.Metrics),
		Logs:                  mro.flattenAddOnLogs(input.Logs),
		Nvidia:                mro.flattenAddOnNvidia(input.Nvidia),
		AwsEFS:                mro.flattenAddOnAwsConfig(input.AwsEFS),
		AwsECR:                mro.flattenAddOnAwsConfig(input.AwsECR),
		AwsELB:                mro.flattenAddOnAwsConfig(input.AwsELB),
		AzureACR:              mro.flattenAddOnAzureAcr(input.AzureACR),
		Sysbox:                mro.flattenAddOnConfig(input.Sysbox),
	}

	// Return a slice containing the single block
	return []models.AddOnsModel{block}
}

// flattenAddOnConfig returns a Terraform bool indicating whether the addon config is enabled.
func (mro *Mk8sResourceOperator) flattenAddOnConfig(input *client.Mk8sNonCustomizableAddonConfig) types.Bool {
	// If the input config is nil, the addon is disabled
	if input == nil {
		return types.BoolNull()
	}

	// Otherwise, the addon is enabled
	return types.BoolValue(true)
}

// flattenAddOnAzureWorkloadIdentity transforms *client.Mk8sAzureWorkloadIdentityAddOnConfig into a []models.AddOnAzureWorkloadIdentityModel.
func (mro *Mk8sResourceOperator) flattenAddOnAzureWorkloadIdentity(input *client.Mk8sAzureWorkloadIdentityAddOnConfig) []models.AddOnAzureWorkloadIdentityModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnAzureWorkloadIdentityModel{
		TenantId: types.StringPointerValue(input.TenantId),
	}

	// Return a slice containing the single block
	return []models.AddOnAzureWorkloadIdentityModel{block}
}

// flattenAddOnMetrics transforms *client.Mk8sMetricsAddOnConfig into a []models.AddOnsMetricsModel.
func (mro *Mk8sResourceOperator) flattenAddOnMetrics(input *client.Mk8sMetricsAddOnConfig) []models.AddOnsMetricsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnsMetricsModel{
		KubeState:       types.BoolPointerValue(input.KubeState),
		CoreDns:         types.BoolPointerValue(input.CoreDns),
		Kubelet:         types.BoolPointerValue(input.Kubelet),
		Apiserver:       types.BoolPointerValue(input.Apiserver),
		NodeExporter:    types.BoolPointerValue(input.NodeExporter),
		Cadvisor:        types.BoolPointerValue(input.Cadvisor),
		ScrapeAnnotated: mro.flattenAddOnMetricsScrapeAnnotated(input.ScrapeAnnotated),
	}

	// Return a slice containing the single block
	return []models.AddOnsMetricsModel{block}
}

// flattenAddOnMetricsScrapeAnnotated transforms *client.Mk8sMetricsScrapeAnnotated into a []models.AddOnsMetricsScrapeAnnotatedModel.
func (mro *Mk8sResourceOperator) flattenAddOnMetricsScrapeAnnotated(input *client.Mk8sMetricsScrapeAnnotated) []models.AddOnsMetricsScrapeAnnotatedModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnsMetricsScrapeAnnotatedModel{
		IntervalSeconds:   FlattenInt(input.IntervalSeconds),
		IncludeNamespaces: types.StringPointerValue(input.IncludeNamespaces),
		ExcludeNamespaces: types.StringPointerValue(input.ExcludeNamespaces),
		RetainLabels:      types.StringPointerValue(input.RetainLabels),
	}

	// Return a slice containing the single block
	return []models.AddOnsMetricsScrapeAnnotatedModel{block}
}

// flattenAddOnLogs transforms *client.Mk8sLogsAddOnConfig into a []models.AddOnsLogsModel.
func (mro *Mk8sResourceOperator) flattenAddOnLogs(input *client.Mk8sLogsAddOnConfig) []models.AddOnsLogsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnsLogsModel{
		AuditEnabled:       types.BoolPointerValue(input.AuditEnabled),
		IncludeNamespaaces: types.StringPointerValue(input.IncludeNamespaces),
		ExcludeNamespaces:  types.StringPointerValue(input.ExcludeNamespaces),
		Docker:             types.BoolPointerValue(input.Docker),
		Kubelet:            types.BoolPointerValue(input.Kubelet),
		Kernel:             types.BoolPointerValue(input.Kernel),
		Events:             types.BoolPointerValue(input.Events),
	}

	// Return a slice containing the single block
	return []models.AddOnsLogsModel{block}
}

// flattenAddOnNvidia transforms *client.Mk8sNvidiaAddOnConfig into a []models.AddOnsNvidiaModel.
func (mro *Mk8sResourceOperator) flattenAddOnNvidia(input *client.Mk8sNvidiaAddOnConfig) []models.AddOnsNvidiaModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnsNvidiaModel{
		TaintGpuNodes: types.BoolPointerValue(input.TaintGPUNodes),
	}

	// Return a slice containing the single block
	return []models.AddOnsNvidiaModel{block}
}

// flattenAddOnAwsConfig transforms *client.Mk8sAwsAddOnConfig into a []models.AddOnsHasRoleArnModel.
func (mro *Mk8sResourceOperator) flattenAddOnAwsConfig(input *client.Mk8sAwsAddOnConfig) []models.AddOnsHasRoleArnModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnsHasRoleArnModel{
		RoleArn: types.StringPointerValue(input.RoleArn),
	}

	// Return a slice containing the single block
	return []models.AddOnsHasRoleArnModel{block}
}

// flattenAddOnAzureAcr transforms *client.Mk8sAzureACRAddOnConfig into a []models.AddOnsAzureAcrModel.
func (mro *Mk8sResourceOperator) flattenAddOnAzureAcr(input *client.Mk8sAzureACRAddOnConfig) []models.AddOnsAzureAcrModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AddOnsAzureAcrModel{
		ClientId: types.StringPointerValue(input.ClientId),
	}

	// Return a slice containing the single block
	return []models.AddOnsAzureAcrModel{block}
}

// flattenStatus transforms *client.Mk8sStatus into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatus(input *client.Mk8sStatus) types.List {
	// Get attribute types
	elementType := models.StatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusModel{
		OidcProviderUrl: types.StringPointerValue(input.OidcProviderUrl),
		ServerUrl:       types.StringPointerValue(input.ServerUrl),
		HomeLocation:    types.StringPointerValue(input.HomeLocation),
		AddOns:          mro.flattenStatusAddOn(input.AddOns),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusModel{block})
}

// flattenStatusAddOn transforms *client.Mk8sStatusAddOns into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatusAddOn(input *client.Mk8sStatusAddOns) types.List {
	// Get attribute types
	elementType := models.StatusAddOnsModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusAddOnsModel{
		Dashboard:           mro.flattenStatusAddOnDashboard(input.Dashboard),
		AwsWorkloadIdentity: mro.flattenStatusAddOnAwsWorkloadIdentity(input.AwsWorkloadIdentity),
		Metrics:             mro.flattenStatusAddOnMetrics(input.Metrics),
		Logs:                mro.flattenStatusAddOnLogs(input.Logs),
		AwsECR:              mro.flattenStatusAddOnAwsConfig(input.AwsECR),
		AwsEFS:              mro.flattenStatusAddOnAwsConfig(input.AwsEFS),
		AwsELB:              mro.flattenStatusAddOnAwsConfig(input.AwsELB),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusAddOnsModel{block})
}

// flattenStatusAddOnDashboard transforms *client.Mk8sDashboardAddOnStatus into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatusAddOnDashboard(input *client.Mk8sDashboardAddOnStatus) types.List {
	// Get attribute types
	elementType := models.StatusAddOnsDashboardModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusAddOnsDashboardModel{
		Url: types.StringPointerValue(input.Url),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusAddOnsDashboardModel{block})
}

// flattenStatusAddOnAwsWorkloadIdentity transforms *client.Mk8sAwsWorkloadIdentityAddOnStatus into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatusAddOnAwsWorkloadIdentity(input *client.Mk8sAwsWorkloadIdentityAddOnStatus) types.List {
	// Get attribute types
	elementType := models.StatusAddOnsAwsWorkloadIdentityModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusAddOnsAwsWorkloadIdentityModel{
		OidcProviderConfig: mro.flattenStatusAddOnAwsWorkloadIdentityOidcProviderConfig(input.OidcProviderConfig),
		TrustPolicy:        mro.flattenObjectUnknown(input.TrustPolicy),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusAddOnsAwsWorkloadIdentityModel{block})
}

// flattenStatusAddOnAwsWorkloadIdentityOidcProviderConfig transforms *client.Mk8sOidcProviderConfig into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatusAddOnAwsWorkloadIdentityOidcProviderConfig(input *client.Mk8sOidcProviderConfig) types.List {
	// Get attribute types
	elementType := models.StatusAddOnsAwsWorkloadIdentityOidcProviderConfigModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusAddOnsAwsWorkloadIdentityOidcProviderConfigModel{
		ProviderUrl: types.StringPointerValue(input.ProviderUrl),
		Audience:    types.StringPointerValue(input.Audience),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusAddOnsAwsWorkloadIdentityOidcProviderConfigModel{block})
}

// flattenObjectUnknown flattens an unknown map into a Terraform string containing its JSON representation, or returns a null string if unknown is nil.
func (mro *Mk8sResourceOperator) flattenObjectUnknown(unknown *map[string]interface{}) types.String {
	// If the unknown map is nil, return a null string
	if unknown == nil {
		return types.StringNull()
	}

	// Marshal the map into JSON bytes
	jsonData, _ := json.Marshal(*unknown)

	// Return the JSON bytes as a Terraform string value
	return types.StringValue(string(jsonData))
}

// flattenStatusAddOnMetrics transforms *client.Mk8sMetricsAddOnStatus into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatusAddOnMetrics(input *client.Mk8sMetricsAddOnStatus) types.List {
	// Get attribute types
	elementType := models.StatusAddOnsMetricsModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusAddOnsMetricsModel{
		PrometheusEndpoint: types.StringPointerValue(input.PrometheusEndpoint),
		RemoteWriteConfig:  mro.flattenObjectUnknown(input.RemoteWriteConfig),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusAddOnsMetricsModel{block})
}

// flattenStatusAddOnLogs transforms *client.Mk8sLogsAddOnStatus into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatusAddOnLogs(input *client.Mk8sLogsAddOnStatus) types.List {
	// Get attribute types
	elementType := models.StatusAddOnsLogsModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusAddOnsLogsModel{
		LokiAddress: types.StringPointerValue(input.LokiAddress),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusAddOnsLogsModel{block})
}

// flattenStatusAddOnAwsConfig transforms *client.Mk8sAwsAddOnStatus into a Terraform types.List.
func (mro *Mk8sResourceOperator) flattenStatusAddOnAwsConfig(input *client.Mk8sAwsAddOnStatus) types.List {
	// Get attribute types
	elementType := models.StatusAddOnsAwsStatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusAddOnsAwsStatusModel{
		TrustPolicy: mro.flattenObjectUnknown(input.TrustPolicy),
	}

	// Return the successfully created types.List
	return FlattenList(mro.Ctx, mro.Diags, []models.StatusAddOnsAwsStatusModel{block})
}
