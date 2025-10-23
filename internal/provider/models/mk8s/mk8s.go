package mk8s

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Main Models ***/

// Firewall //

type FirewallModel struct {
	SourceCIDR  types.String `tfsdk:"source_cidr"`
	Description types.String `tfsdk:"description"`
}

func (f FirewallModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"source_cidr": types.StringType,
			"description": types.StringType,
		},
	}
}

// Generic Provider //

type GenericProviderModel struct {
	Location   types.String `tfsdk:"location"`
	Networking types.List   `tfsdk:"networking"`
	NodePools  types.Set    `tfsdk:"node_pool"`
}

func (g GenericProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"location":   types.StringType,
			"networking": types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"node_pool":  types.SetType{ElemType: GenericProviderNodePoolModel{}.AttributeTypes()},
		},
	}
}

// Generic Provider -> Node Pool //

type GenericProviderNodePoolModel struct {
	Name   types.String `tfsdk:"name"`
	Labels types.Map    `tfsdk:"labels"`
	Taints types.Set    `tfsdk:"taint"`
}

func (g GenericProviderNodePoolModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":   types.StringType,
			"labels": types.MapType{ElemType: types.StringType},
			"taint":  types.SetType{ElemType: GenericProviderNodePoolTaintModel{}.AttributeTypes()},
		},
	}
}

// Generic Provider -> Node Pool -> Taint //

type GenericProviderNodePoolTaintModel struct {
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
	Effect types.String `tfsdk:"effect"`
}

func (g GenericProviderNodePoolTaintModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":    types.StringType,
			"value":  types.StringType,
			"effect": types.StringType,
		},
	}
}

// Hetzner Provider //

type HetznerProviderModel struct {
	Region                   types.String `tfsdk:"region"`
	HetznerLabels            types.Map    `tfsdk:"hetzner_labels"`
	Networking               types.List   `tfsdk:"networking"`
	PreInstallScript         types.String `tfsdk:"pre_install_script"`
	TokenSecretLink          types.String `tfsdk:"token_secret_link"`
	NetworkId                types.String `tfsdk:"network_id"`
	FirewallId               types.String `tfsdk:"firewall_id"`
	NodePools                types.Set    `tfsdk:"node_pool"`
	DedicatedServerNodePools types.Set    `tfsdk:"dedicated_server_node_pool"`
	Image                    types.String `tfsdk:"image"`
	SshKey                   types.String `tfsdk:"ssh_key"`
	Autoscaler               types.List   `tfsdk:"autoscaler"`
	FloatingIpSelector       types.Map    `tfsdk:"floating_ip_selector"`
}

func (h HetznerProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":                     types.StringType,
			"hetzner_labels":             types.MapType{ElemType: types.StringType},
			"networking":                 types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"pre_install_script":         types.StringType,
			"token_secret_link":          types.StringType,
			"network_id":                 types.StringType,
			"firewall_id":                types.StringType,
			"node_pool":                  types.SetType{ElemType: HetznerProviderNodePoolModel{}.AttributeTypes()},
			"dedicated_server_node_pool": types.SetType{ElemType: GenericProviderNodePoolModel{}.AttributeTypes()},
			"image":                      types.StringType,
			"ssh_key":                    types.StringType,
			"autoscaler":                 types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
			"floating_ip_selector":       types.MapType{ElemType: types.StringType},
		},
	}
}

// Hetzner Provider -> Node Pool //

type HetznerProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	ServerType    types.String `tfsdk:"server_type"`
	OverrideImage types.String `tfsdk:"override_image"`
	MinSize       types.Int32  `tfsdk:"min_size"`
	MaxSize       types.Int32  `tfsdk:"max_size"`
}

func (h HetznerProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"server_type":    types.StringType,
		"override_image": types.StringType,
		"min_size":       types.Int32Type,
		"max_size":       types.Int32Type,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// AWS Provider //

type AwsProviderModel struct {
	Region               types.String `tfsdk:"region"`
	AwsTags              types.Map    `tfsdk:"aws_tags"`
	SkipCreateRoles      types.Bool   `tfsdk:"skip_create_roles"`
	Networking           types.List   `tfsdk:"networking"`
	PreInstallScript     types.String `tfsdk:"pre_install_script"`
	Image                types.List   `tfsdk:"image"`
	DeployRoleArn        types.String `tfsdk:"deploy_role_arn"`
	DeployRoleChain      types.List   `tfsdk:"deploy_role_chain"`
	VpcId                types.String `tfsdk:"vpc_id"`
	KeyPair              types.String `tfsdk:"key_pair"`
	DiskEncryptionKeyArn types.String `tfsdk:"disk_encryption_key_arn"`
	SecurityGroupIds     types.Set    `tfsdk:"security_group_ids"`
	ExtraNodePolicies    types.Set    `tfsdk:"extra_node_policies"`
	NodePools            types.Set    `tfsdk:"node_pool"`
	Autoscaler           types.List   `tfsdk:"autoscaler"`
}

func (a AwsProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":                  types.StringType,
			"aws_tags":                types.MapType{ElemType: types.StringType},
			"skip_create_roles":       types.BoolType,
			"networking":              types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"pre_install_script":      types.StringType,
			"image":                   types.ListType{ElemType: AwsProviderAmiModel{}.AttributeTypes()},
			"deploy_role_arn":         types.StringType,
			"deploy_role_chain":       types.ListType{ElemType: AwsProviderAssumeRoleLinkModel{}.AttributeTypes()},
			"vpc_id":                  types.StringType,
			"key_pair":                types.StringType,
			"disk_encryption_key_arn": types.StringType,
			"security_group_ids":      types.SetType{ElemType: types.StringType},
			"extra_node_policies":     types.SetType{ElemType: types.StringType},
			"node_pool":               types.SetType{ElemType: AwsProviderNodePoolModel{}.AttributeTypes()},
			"autoscaler":              types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
		},
	}
}

// AWS Provider -> AMI //

type AwsProviderAmiModel struct {
	Recommended types.String `tfsdk:"recommended"`
	Exact       types.String `tfsdk:"exact"`
}

func (a AwsProviderAmiModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"recommended": types.StringType,
			"exact":       types.StringType,
		},
	}
}

// AWS Provider -> Deploy Role Chain //

type AwsProviderAssumeRoleLinkModel struct {
	RoleArn           types.String `tfsdk:"role_arn"`
	ExternalId        types.String `tfsdk:"external_id"`
	SessionNamePrefix types.String `tfsdk:"session_name_prefix"`
}

func (a AwsProviderAssumeRoleLinkModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"role_arn":            types.StringType,
			"external_id":         types.StringType,
			"session_name_prefix": types.StringType,
		},
	}
}

// AWS Provider -> Node Pool //

type AwsProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	InstanceTypes                       types.Set    `tfsdk:"instance_types"`
	OverrideImage                       types.List   `tfsdk:"override_image"`
	BootDiskSize                        types.Int32  `tfsdk:"boot_disk_size"`
	MinSize                             types.Int32  `tfsdk:"min_size"`
	MaxSize                             types.Int32  `tfsdk:"max_size"`
	OnDemandBaseCapacity                types.Int32  `tfsdk:"on_demand_base_capacity"`
	OnDemandPercentageAboveBaseCapacity types.Int32  `tfsdk:"on_demand_percentage_above_base_capacity"`
	SpotAllocationStrategy              types.String `tfsdk:"spot_allocation_strategy"`
	SubnetIds                           types.Set    `tfsdk:"subnet_ids"`
	ExtraSecurityGroupIds               types.Set    `tfsdk:"extra_security_group_ids"`
}

func (a AwsProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"instance_types":          types.SetType{ElemType: types.StringType},
		"override_image":          types.ListType{ElemType: AwsProviderAmiModel{}.AttributeTypes()},
		"boot_disk_size":          types.Int32Type,
		"min_size":                types.Int32Type,
		"max_size":                types.Int32Type,
		"on_demand_base_capacity": types.Int32Type,
		"on_demand_percentage_above_base_capacity": types.Int32Type,
		"spot_allocation_strategy":                 types.StringType,
		"subnet_ids":                               types.SetType{ElemType: types.StringType},
		"extra_security_group_ids":                 types.SetType{ElemType: types.StringType},
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Linode Provider //

type LinodeProviderModel struct {
	Region           types.String `tfsdk:"region"`
	TokenSecretLink  types.String `tfsdk:"token_secret_link"`
	FirewallId       types.String `tfsdk:"firewall_id"`
	NodePools        types.Set    `tfsdk:"node_pool"`
	Image            types.String `tfsdk:"image"`
	AuthorizedUsers  types.Set    `tfsdk:"authorized_users"`
	AuthorizedKeys   types.Set    `tfsdk:"authorized_keys"`
	VpcId            types.String `tfsdk:"vpc_id"`
	PreInstallScript types.String `tfsdk:"pre_install_script"`
	Networking       types.List   `tfsdk:"networking"`
	Autoscaler       types.List   `tfsdk:"autoscaler"`
}

func (l LinodeProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":             types.StringType,
			"token_secret_link":  types.StringType,
			"firewall_id":        types.StringType,
			"node_pool":          types.SetType{ElemType: LinodeProviderNodePoolModel{}.AttributeTypes()},
			"image":              types.StringType,
			"authorized_users":   types.SetType{ElemType: types.StringType},
			"authorized_keys":    types.SetType{ElemType: types.StringType},
			"vpc_id":             types.StringType,
			"pre_install_script": types.StringType,
			"networking":         types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"autoscaler":         types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
		},
	}
}

// Linode Provider -> Node Pool //

type LinodeProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	ServerType    types.String `tfsdk:"server_type"`
	OverrideImage types.String `tfsdk:"override_image"`
	SubnetId      types.String `tfsdk:"subnet_id"`
	MinSize       types.Int32  `tfsdk:"min_size"`
	MaxSize       types.Int32  `tfsdk:"max_size"`
}

func (l LinodeProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"server_type":    types.StringType,
		"override_image": types.StringType,
		"subnet_id":      types.StringType,
		"min_size":       types.Int32Type,
		"max_size":       types.Int32Type,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Oblivus Provider //

type OblivusProviderModel struct {
	Datacenter        types.String `tfsdk:"datacenter"`
	TokenSecretLink   types.String `tfsdk:"token_secret_link"`
	NodePools         types.Set    `tfsdk:"node_pool"`
	SshKeys           types.Set    `tfsdk:"ssh_keys"`
	UnmanagedNodePool types.Set    `tfsdk:"unmanaged_node_pool"`
	Autoscaler        types.List   `tfsdk:"autoscaler"`
	PreInstallScript  types.String `tfsdk:"pre_install_script"`
}

func (o OblivusProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"datacenter":          types.StringType,
			"token_secret_link":   types.StringType,
			"node_pool":           types.SetType{ElemType: OblivusProviderNodePoolModel{}.AttributeTypes()},
			"ssh_keys":            types.SetType{ElemType: types.StringType},
			"unmanaged_node_pool": types.SetType{ElemType: GenericProviderNodePoolModel{}.AttributeTypes()},
			"autoscaler":          types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
			"pre_install_script":  types.StringType,
		},
	}
}

// Oblivus Provider -> Node Pool //

type OblivusProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	MinSize types.Int32  `tfsdk:"min_size"`
	MaxSize types.Int32  `tfsdk:"max_size"`
	Flavor  types.String `tfsdk:"flavor"`
}

func (o OblivusProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"min_size": types.Int32Type,
		"max_size": types.Int32Type,
		"flavor":   types.StringType,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Lambdalabs Provider //

type LambdalabsProviderModel struct {
	Region             types.String `tfsdk:"region"`
	TokenSecretLink    types.String `tfsdk:"token_secret_link"`
	NodePools          types.Set    `tfsdk:"node_pool"`
	SshKey             types.String `tfsdk:"ssh_key"`
	UnmanagedNodePools types.Set    `tfsdk:"unmanaged_node_pool"`
	Autoscaler         types.List   `tfsdk:"autoscaler"`
	FileSystems        types.Set    `tfsdk:"file_systems"`
	PreInstallScript   types.String `tfsdk:"pre_install_script"`
}

func (l LambdalabsProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":              types.StringType,
			"token_secret_link":   types.StringType,
			"node_pool":           types.SetType{ElemType: LambdalabsProviderNodePoolModel{}.AttributeTypes()},
			"ssh_key":             types.StringType,
			"unmanaged_node_pool": types.SetType{ElemType: GenericProviderNodePoolModel{}.AttributeTypes()},
			"autoscaler":          types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
			"file_systems":        types.SetType{ElemType: types.StringType},
			"pre_install_script":  types.StringType,
		},
	}
}

// Lambdalabs Provider -> Node Pool //

type LambdalabsProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	MinSize      types.Int32  `tfsdk:"min_size"`
	MaxSize      types.Int32  `tfsdk:"max_size"`
	InstanceType types.String `tfsdk:"instance_type"`
}

func (l LambdalabsProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"min_size":      types.Int32Type,
		"max_size":      types.Int32Type,
		"instance_type": types.StringType,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Paperspace Provider //

type PaperspaceProviderModel struct {
	Region             types.String `tfsdk:"region"`
	TokenSecretLink    types.String `tfsdk:"token_secret_link"`
	SharedDrives       types.Set    `tfsdk:"shared_drives"`
	NodePools          types.Set    `tfsdk:"node_pool"`
	Autoscaler         types.List   `tfsdk:"autoscaler"`
	UnmanagedNodePools types.Set    `tfsdk:"unmanaged_node_pool"`
	PreInstallScript   types.String `tfsdk:"pre_install_script"`
	UserIds            types.Set    `tfsdk:"user_ids"`
	NetworkId          types.String `tfsdk:"network_id"`
}

func (p PaperspaceProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":              types.StringType,
			"token_secret_link":   types.StringType,
			"shared_drives":       types.SetType{ElemType: types.StringType},
			"node_pool":           types.SetType{ElemType: PaperspaceProviderNodePoolModel{}.AttributeTypes()},
			"autoscaler":          types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
			"unmanaged_node_pool": types.SetType{ElemType: GenericProviderNodePoolModel{}.AttributeTypes()},
			"pre_install_script":  types.StringType,
			"user_ids":            types.SetType{ElemType: types.StringType},
			"network_id":          types.StringType,
		},
	}
}

// Paperspace Provider -> Node Pool //

type PaperspaceProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	MinSize      types.Int32  `tfsdk:"min_size"`
	MaxSize      types.Int32  `tfsdk:"max_size"`
	PublicIpType types.String `tfsdk:"public_ip_type"`
	BootDiskSize types.Int32  `tfsdk:"boot_disk_size"`
	MachineType  types.String `tfsdk:"machine_type"`
}

func (p PaperspaceProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"min_size":       types.Int32Type,
		"max_size":       types.Int32Type,
		"public_ip_type": types.StringType,
		"boot_disk_size": types.Int32Type,
		"machine_type":   types.StringType,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Ephemeral Provider //

type EphemeralProviderModel struct {
	Location  types.String `tfsdk:"location"`
	NodePools types.Set    `tfsdk:"node_pool"`
}

func (e EphemeralProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"location":  types.StringType,
			"node_pool": types.SetType{ElemType: EphemeralProviderNodePoolModel{}.AttributeTypes()},
		},
	}
}

// Ephemeral Provider -> Node Pool //

type EphemeralProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	Count  types.Int32  `tfsdk:"count"`
	Arch   types.String `tfsdk:"arch"`
	Flavor types.String `tfsdk:"flavor"`
	Cpu    types.String `tfsdk:"cpu"`
	Memory types.String `tfsdk:"memory"`
}

func (e EphemeralProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"count":  types.Int32Type,
		"arch":   types.StringType,
		"flavor": types.StringType,
		"cpu":    types.StringType,
		"memory": types.StringType,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Triton Provider //

type TritonProviderModel struct {
	Connection       types.List   `tfsdk:"connection"`
	Networking       types.List   `tfsdk:"networking"`
	PreInstallScript types.String `tfsdk:"pre_install_script"`
	Location         types.String `tfsdk:"location"`
	LoadBalancer     types.List   `tfsdk:"load_balancer"`
	PrivateNetworkId types.String `tfsdk:"private_network_id"`
	FirewallEnabled  types.Bool   `tfsdk:"firewall_enabled"`
	NodePools        types.Set    `tfsdk:"node_pool"`
	ImageId          types.String `tfsdk:"image_id"`
	SshKeys          types.Set    `tfsdk:"ssh_keys"`
	Autoscaler       types.List   `tfsdk:"autoscaler"`
}

func (t TritonProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"connection":         types.ListType{ElemType: TritonProviderConnectionModel{}.AttributeTypes()},
			"networking":         types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"pre_install_script": types.StringType,
			"location":           types.StringType,
			"load_balancer":      types.ListType{ElemType: TritonProviderLoadBalancerModel{}.AttributeTypes()},
			"private_network_id": types.StringType,
			"firewall_enabled":   types.BoolType,
			"node_pool":          types.SetType{ElemType: TritonProviderNodePoolModel{}.AttributeTypes()},
			"image_id":           types.StringType,
			"ssh_keys":           types.SetType{ElemType: types.StringType},
			"autoscaler":         types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
		},
	}
}

// Triton Provider -> Connection //

type TritonProviderConnectionModel struct {
	Url                  types.String `tfsdk:"url"`
	Account              types.String `tfsdk:"account"`
	User                 types.String `tfsdk:"user"`
	PrivateKeySecretLink types.String `tfsdk:"private_key_secret_link"`
}

func (t TritonProviderConnectionModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"url":                     types.StringType,
			"account":                 types.StringType,
			"user":                    types.StringType,
			"private_key_secret_link": types.StringType,
		},
	}
}

// Triton Provider -> Load Balancer //

type TritonProviderLoadBalancerModel struct {
	Manual  types.List `tfsdk:"manual"`
	None    types.List `tfsdk:"none"`
	Gateway types.List `tfsdk:"gateway"`
}

func (t TritonProviderLoadBalancerModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"manual":  types.ListType{ElemType: TritonProviderLoadBalancerManualModel{}.AttributeTypes()},
			"none":    types.ListType{ElemType: TritonProviderLoadBalancerNoneModel{}.AttributeTypes()},
			"gateway": types.ListType{ElemType: TritonProviderLoadBalancerGatewayModel{}.AttributeTypes()},
		},
	}
}

// Triton Provider -> Load Balancer -> Manual //

type TritonProviderLoadBalancerManualModel struct {
	PackageId         types.String `tfsdk:"package_id"`
	ImageId           types.String `tfsdk:"image_id"`
	PublicNetworkId   types.String `tfsdk:"public_network_id"`
	PrivateNetworkIds types.Set    `tfsdk:"private_network_ids"`
	Metadata          types.Map    `tfsdk:"metadata"`
	Tags              types.Map    `tfsdk:"tags"`
	Logging           types.List   `tfsdk:"logging"`
	Count             types.Int32  `tfsdk:"count"`
	CnsInternalDomain types.String `tfsdk:"cns_internal_domain"`
	CnsPublicDomain   types.String `tfsdk:"cns_public_domain"`
}

func (t TritonProviderLoadBalancerManualModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"package_id":          types.StringType,
			"image_id":            types.StringType,
			"public_network_id":   types.StringType,
			"private_network_ids": types.SetType{ElemType: types.StringType},
			"metadata":            types.MapType{ElemType: types.StringType},
			"tags":                types.MapType{ElemType: types.StringType},
			"logging":             types.ListType{ElemType: TritonProviderLoadBalancerManualLoggingModel{}.AttributeTypes()},
			"count":               types.Int32Type,
			"cns_internal_domain": types.StringType,
			"cns_public_domain":   types.StringType,
		},
	}
}

// Triton Provider -> Load Balancer -> Manual -> Logging //

type TritonProviderLoadBalancerManualLoggingModel struct {
	NodePort       types.Int32  `tfsdk:"node_port"`
	ExternalSyslog types.String `tfsdk:"external_syslog"`
}

func (t TritonProviderLoadBalancerManualLoggingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"node_port":       types.Int32Type,
			"external_syslog": types.StringType,
		},
	}
}

// Triton Provider -> Load Balancer -> None //

type TritonProviderLoadBalancerNoneModel struct{}

func (t TritonProviderLoadBalancerNoneModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{},
	}
}

// Triton Provider -> Load Balancer -> Gateway //

type TritonProviderLoadBalancerGatewayModel struct{}

func (t TritonProviderLoadBalancerGatewayModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{},
	}
}

// Triton Provider -> Node Pool //

type TritonProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	PackageId         types.String `tfsdk:"package_id"`
	OverrideImageId   types.String `tfsdk:"override_image_id"`
	PublicNetworkId   types.String `tfsdk:"public_network_id"`
	PrivateNetworkIds types.Set    `tfsdk:"private_network_ids"`
	TritonTags        types.Map    `tfsdk:"triton_tags"`
	MinSize           types.Int32  `tfsdk:"min_size"`
	MaxSize           types.Int32  `tfsdk:"max_size"`
}

func (t TritonProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"package_id":          types.StringType,
		"override_image_id":   types.StringType,
		"public_network_id":   types.StringType,
		"private_network_ids": types.SetType{ElemType: types.StringType},
		"triton_tags":         types.MapType{ElemType: types.StringType},
		"min_size":            types.Int32Type,
		"max_size":            types.Int32Type,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Azure Provider //

type AzureProviderModel struct {
	Location         types.String `tfsdk:"location"`
	SubscriptionId   types.String `tfsdk:"subscription_id"`
	SdkSecretLink    types.String `tfsdk:"sdk_secret_link"`
	ResourceGroup    types.String `tfsdk:"resource_group"`
	Networking       types.List   `tfsdk:"networking"`
	PreInstallScript types.String `tfsdk:"pre_install_script"`
	Image            types.List   `tfsdk:"image"`
	SshKeys          types.Set    `tfsdk:"ssh_keys"`
	NetworkId        types.String `tfsdk:"network_id"`
	Tags             types.Map    `tfsdk:"tags"`
	NodePools        types.Set    `tfsdk:"node_pool"`
	Autoscaler       types.List   `tfsdk:"autoscaler"`
}

func (a AzureProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"location":           types.StringType,
			"subscription_id":    types.StringType,
			"sdk_secret_link":    types.StringType,
			"resource_group":     types.StringType,
			"networking":         types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"pre_install_script": types.StringType,
			"image":              types.ListType{ElemType: AzureProviderImageModel{}.AttributeTypes()},
			"ssh_keys":           types.SetType{ElemType: types.StringType},
			"network_id":         types.StringType,
			"tags":               types.MapType{ElemType: types.StringType},
			"node_pool":          types.SetType{ElemType: AzureProviderNodePoolModel{}.AttributeTypes()},
			"autoscaler":         types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
		},
	}
}

// Azure Provider -> Image //

type AzureProviderImageModel struct {
	Recommended types.String `tfsdk:"recommended"`
	Reference   types.List   `tfsdk:"reference"`
}

func (a AzureProviderImageModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"recommended": types.StringType,
			"reference":   types.ListType{ElemType: AzureProviderImageReferenceModel{}.AttributeTypes()},
		},
	}
}

// Azure Provider -> Image -> Reference //

type AzureProviderImageReferenceModel struct {
	Publisher types.String `tfsdk:"publisher"`
	Offer     types.String `tfsdk:"offer"`
	Sku       types.String `tfsdk:"sku"`
	Version   types.String `tfsdk:"version"`
}

func (a AzureProviderImageReferenceModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"publisher": types.StringType,
			"offer":     types.StringType,
			"sku":       types.StringType,
			"version":   types.StringType,
		},
	}
}

// Azure Provider -> Node Pool //

type AzureProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	Size          types.String `tfsdk:"size"`
	SubnetId      types.String `tfsdk:"subnet_id"`
	Zones         types.Set    `tfsdk:"zones"`
	OverrideImage types.List   `tfsdk:"override_image"`
	BootDiskSize  types.Int32  `tfsdk:"boot_disk_size"`
	MinSize       types.Int32  `tfsdk:"min_size"`
	MaxSize       types.Int32  `tfsdk:"max_size"`
}

func (a AzureProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"size":           types.StringType,
		"subnet_id":      types.StringType,
		"zones":          types.SetType{ElemType: types.Int32Type},
		"override_image": types.ListType{ElemType: AzureProviderImageModel{}.AttributeTypes()},
		"boot_disk_size": types.Int32Type,
		"min_size":       types.Int32Type,
		"max_size":       types.Int32Type,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Digital Ocean Provider //

type DigitalOceanProviderModel struct {
	Region           types.String `tfsdk:"region"`
	DigitalOceanTags types.Set    `tfsdk:"digital_ocean_tags"`
	Networking       types.List   `tfsdk:"networking"`
	PreInstallScript types.String `tfsdk:"pre_install_script"`
	TokenSecretLink  types.String `tfsdk:"token_secret_link"`
	VpcId            types.String `tfsdk:"vpc_id"`
	NodePools        types.Set    `tfsdk:"node_pool"`
	Image            types.String `tfsdk:"image"`
	SshKeys          types.Set    `tfsdk:"ssh_keys"`
	ExtraSshKeys     types.Set    `tfsdk:"extra_ssh_keys"`
	Autoscaler       types.List   `tfsdk:"autoscaler"`
	ReservedIps      types.Set    `tfsdk:"reserved_ips"`
}

func (d DigitalOceanProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":             types.StringType,
			"digital_ocean_tags": types.SetType{ElemType: types.StringType},
			"networking":         types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"pre_install_script": types.StringType,
			"token_secret_link":  types.StringType,
			"vpc_id":             types.StringType,
			"node_pool":          types.SetType{ElemType: DigitalOceanProviderNodePoolModel{}.AttributeTypes()},
			"image":              types.StringType,
			"ssh_keys":           types.SetType{ElemType: types.StringType},
			"extra_ssh_keys":     types.SetType{ElemType: types.StringType},
			"autoscaler":         types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
			"reserved_ips":       types.SetType{ElemType: types.StringType},
		},
	}
}

// Digital Ocean Provider -> Node Pool //

type DigitalOceanProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	DropletSize   types.String `tfsdk:"droplet_size"`
	OverrideImage types.String `tfsdk:"override_image"`
	MinSize       types.Int32  `tfsdk:"min_size"`
	MaxSize       types.Int32  `tfsdk:"max_size"`
}

func (d DigitalOceanProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"droplet_size":   types.StringType,
		"override_image": types.StringType,
		"min_size":       types.Int32Type,
		"max_size":       types.Int32Type,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Gcp Provider //

type GcpProviderModel struct {
	ProjectId        types.String `tfsdk:"project_id"`
	Region           types.String `tfsdk:"region"`
	GcpLabels        types.Map    `tfsdk:"gcp_labels"`
	Network          types.String `tfsdk:"network"`
	SaKeyLink        types.String `tfsdk:"sa_key_link"`
	Networking       types.List   `tfsdk:"networking"`
	PreInstallScript types.String `tfsdk:"pre_install_script"`
	Image            types.List   `tfsdk:"image"`
	NodePools        types.Set    `tfsdk:"node_pool"`
	Autoscaler       types.List   `tfsdk:"autoscaler"`
}

func (g GcpProviderModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"project_id":         types.StringType,
			"region":             types.StringType,
			"gcp_labels":         types.MapType{ElemType: types.StringType},
			"network":            types.StringType,
			"sa_key_link":        types.StringType,
			"networking":         types.ListType{ElemType: NetworkingModel{}.AttributeTypes()},
			"pre_install_script": types.StringType,
			"image":              types.ListType{ElemType: GcpProviderImageModel{}.AttributeTypes()},
			"node_pool":          types.SetType{ElemType: GcpProviderNodePoolModel{}.AttributeTypes()},
			"autoscaler":         types.ListType{ElemType: AutoscalerModel{}.AttributeTypes()},
		},
	}
}

// Gcp Provider -> Image //

type GcpProviderImageModel struct {
	Recommended types.String `tfsdk:"recommended"`
}

func (g GcpProviderImageModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"recommended": types.StringType,
		},
	}
}

// Gcp Provider -> Node Pool //

type GcpProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	MachineType   types.String `tfsdk:"machine_type"`
	Zone          types.String `tfsdk:"zone"`
	OverrideImage types.List   `tfsdk:"override_image"`
	BootDiskSize  types.Int32  `tfsdk:"boot_disk_size"`
	MinSize       types.Int32  `tfsdk:"min_size"`
	MaxSize       types.Int32  `tfsdk:"max_size"`
	Subnet        types.String `tfsdk:"subnet"`
}

func (g GcpProviderNodePoolModel) AttributeTypes() attr.Type {
	// Get the attribute types from GenericProviderNodePoolModel
	base := GenericProviderNodePoolModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{
		"machine_type":   types.StringType,
		"zone":           types.StringType,
		"override_image": types.ListType{ElemType: GcpProviderImageModel{}.AttributeTypes()},
		"boot_disk_size": types.Int32Type,
		"min_size":       types.Int32Type,
		"max_size":       types.Int32Type,
		"subnet":         types.StringType,
	}

	// Add the attributes from base
	maps.Copy(merged, base.AttrTypes)

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Add Ons //

type AddOnsModel struct {
	Dashboard             types.Bool   `tfsdk:"dashboard"`
	AzureWorkloadIdentity types.List   `tfsdk:"azure_workload_identity"`
	AwsWorkloadIdentity   types.Bool   `tfsdk:"aws_workload_identity"`
	LocalPathStorage      types.Bool   `tfsdk:"local_path_storage"`
	Metrics               types.List   `tfsdk:"metrics"`
	Logs                  types.List   `tfsdk:"logs"`
	RegistryMirror        types.List   `tfsdk:"registry_mirror"`
	Nvidia                types.List   `tfsdk:"nvidia"`
	AwsEFS                types.List   `tfsdk:"aws_efs"`
	AwsECR                types.List   `tfsdk:"aws_ecr"`
	AwsELB                types.List   `tfsdk:"aws_elb"`
	AzureACR              types.List   `tfsdk:"azure_acr"`
	Byok                  types.Object `tfsdk:"byok"`
	Sysbox                types.Bool   `tfsdk:"sysbox"`
}

func (a AddOnsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"dashboard":               types.BoolType,
			"azure_workload_identity": types.ListType{ElemType: AddOnAzureWorkloadIdentityModel{}.AttributeTypes()},
			"aws_workload_identity":   types.BoolType,
			"local_path_storage":      types.BoolType,
			"metrics":                 types.ListType{ElemType: AddOnsMetricsModel{}.AttributeTypes()},
			"logs":                    types.ListType{ElemType: AddOnsLogsModel{}.AttributeTypes()},
			"registry_mirror":         types.ListType{ElemType: AddOnsRegistryMirror{}.AttributeTypes()},
			"nvidia":                  types.ListType{ElemType: AddOnsNvidiaModel{}.AttributeTypes()},
			"aws_efs":                 types.ListType{ElemType: AddOnsHasRoleArnModel{}.AttributeTypes()},
			"aws_ecr":                 types.ListType{ElemType: AddOnsHasRoleArnModel{}.AttributeTypes()},
			"aws_elb":                 types.ListType{ElemType: AddOnsHasRoleArnModel{}.AttributeTypes()},
			"azure_acr":               types.ListType{ElemType: AddOnsAzureAcrModel{}.AttributeTypes()},
			"byok":                    AddOnsByokModel{}.AttributeTypes(),
			"sysbox":                  types.BoolType,
		},
	}
}

// Add Ons -> Azure Workload Identity //

type AddOnAzureWorkloadIdentityModel struct {
	TenantId types.String `tfsdk:"tenant_id"`
}

func (a AddOnAzureWorkloadIdentityModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"tenant_id": types.StringType,
		},
	}
}

// Add Ons -> Metrics //

type AddOnsMetricsModel struct {
	KubeState       types.Bool `tfsdk:"kube_state"`
	CoreDns         types.Bool `tfsdk:"core_dns"`
	Kubelet         types.Bool `tfsdk:"kubelet"`
	Apiserver       types.Bool `tfsdk:"api_server"`
	NodeExporter    types.Bool `tfsdk:"node_exporter"`
	Cadvisor        types.Bool `tfsdk:"cadvisor"`
	ScrapeAnnotated types.List `tfsdk:"scrape_annotated"`
}

func (a AddOnsMetricsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"kube_state":       types.BoolType,
			"core_dns":         types.BoolType,
			"kubelet":          types.BoolType,
			"api_server":       types.BoolType,
			"node_exporter":    types.BoolType,
			"cadvisor":         types.BoolType,
			"scrape_annotated": types.ListType{ElemType: AddOnsMetricsScrapeAnnotatedModel{}.AttributeTypes()},
		},
	}
}

// Add Ons -> Metrics -> Scrape Annotated //

type AddOnsMetricsScrapeAnnotatedModel struct {
	IntervalSeconds   types.Int32  `tfsdk:"interval_seconds"`
	IncludeNamespaces types.String `tfsdk:"include_namespaces"`
	ExcludeNamespaces types.String `tfsdk:"exclude_namespaces"`
	RetainLabels      types.String `tfsdk:"retain_labels"`
}

func (a AddOnsMetricsScrapeAnnotatedModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"interval_seconds":   types.Int32Type,
			"include_namespaces": types.StringType,
			"exclude_namespaces": types.StringType,
			"retain_labels":      types.StringType,
		},
	}
}

// Add Ons -> Logs //

type AddOnsLogsModel struct {
	AuditEnabled       types.Bool   `tfsdk:"audit_enabled"`
	IncludeNamespaaces types.String `tfsdk:"include_namespaces"`
	ExcludeNamespaces  types.String `tfsdk:"exclude_namespaces"`
	Docker             types.Bool   `tfsdk:"docker"`
	Kubelet            types.Bool   `tfsdk:"kubelet"`
	Kernel             types.Bool   `tfsdk:"kernel"`
	Events             types.Bool   `tfsdk:"events"`
}

func (a AddOnsLogsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"audit_enabled":      types.BoolType,
			"include_namespaces": types.StringType,
			"exclude_namespaces": types.StringType,
			"docker":             types.BoolType,
			"kubelet":            types.BoolType,
			"kernel":             types.BoolType,
			"events":             types.BoolType,
		},
	}
}

// Add Ons -> Registry Mirror //

type AddOnsRegistryMirror struct {
	Mirrors types.Set `tfsdk:"mirror"`
}

func (a AddOnsRegistryMirror) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"mirror": types.SetType{ElemType: AddOnsRegistryConfig{}.AttributeTypes()},
		},
	}
}

// Add Ons -> Registry Mirror -> Mirrors //

type AddOnsRegistryConfig struct {
	Registry types.String `tfsdk:"registry"`
	Mirrors  types.Set    `tfsdk:"mirrors"`
}

func (a AddOnsRegistryConfig) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"registry": types.StringType,
			"mirrors":  types.SetType{ElemType: types.StringType},
		},
	}
}

// Add Ons -> Nvidia //

type AddOnsNvidiaModel struct {
	TaintGpuNodes types.Bool `tfsdk:"taint_gpu_nodes"`
}

func (a AddOnsNvidiaModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"taint_gpu_nodes": types.BoolType,
		},
	}
}

// Add Ons -> Has Role Arn //

type AddOnsHasRoleArnModel struct {
	RoleArn types.String `tfsdk:"role_arn"`
}

func (a AddOnsHasRoleArnModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"role_arn": types.StringType,
		},
	}
}

// Add Ons -> Azure ACR //

type AddOnsAzureAcrModel struct {
	ClientId types.String `tfsdk:"client_id"`
}

func (a AddOnsAzureAcrModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"client_id": types.StringType,
		},
	}
}

// Add Ons -> Byok //

type AddOnsByokModel struct {
	IgnoreUpdates types.Bool   `tfsdk:"ignore_updates"`
	Location      types.String `tfsdk:"location"`
	Config        types.Object `tfsdk:"config"`
}

func (a AddOnsByokModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ignore_updates": types.BoolType,
			"location":       types.StringType,
			"config":         AddOnsByokConfigModel{}.AttributeTypes(),
		},
	}
}

// Add Ons -> Byok -> Config //

type AddOnsByokConfigModel struct {
	Actuator      types.Object `tfsdk:"actuator"`
	Middlebox     types.Object `tfsdk:"middlebox"`
	Common        types.Object `tfsdk:"common"`
	Longhorn      types.Object `tfsdk:"longhorn"`
	Ingress       types.Object `tfsdk:"ingress"`
	Istio         types.Object `tfsdk:"istio"`
	LogSplitter   types.Object `tfsdk:"log_splitter"`
	Monitoring    types.Object `tfsdk:"monitoring"`
	Redis         types.Object `tfsdk:"redis"`
	RedisHa       types.Object `tfsdk:"redis_ha"`
	RedisSentinel types.Object `tfsdk:"redis_sentinel"`
	TempoAgent    types.Object `tfsdk:"tempo_agent"`
	InternalDns   types.Object `tfsdk:"internal_dns"`
}

func (a AddOnsByokConfigModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"actuator":       AddOnsByokActuatorModel{}.AttributeTypes(),
			"middlebox":      AddOnsByokMiddleboxModel{}.AttributeTypes(),
			"common":         AddOnsByokCommonModel{}.AttributeTypes(),
			"longhorn":       AddOnsByokLonghornModel{}.AttributeTypes(),
			"ingress":        AddOnsByokIngressModel{}.AttributeTypes(),
			"istio":          AddOnsByokIstioModel{}.AttributeTypes(),
			"log_splitter":   AddOnsByokLogSplitterModel{}.AttributeTypes(),
			"monitoring":     AddOnsByokMonitoringModel{}.AttributeTypes(),
			"redis":          AddOnsByokRedisStringModel{}.AttributeTypes(),
			"redis_ha":       AddOnsByokRedisIntModel{}.AttributeTypes(),
			"redis_sentinel": AddOnsByokRedisIntModel{}.AttributeTypes(),
			"tempo_agent":    AddOnsByokTempoAgentModel{}.AttributeTypes(),
			"internal_dns":   AddOnsByokInternalDnsModel{}.AttributeTypes(),
		},
	}
}

// Add Ons -> Byok -> Config -> Actuator //

type AddOnsByokActuatorModel struct {
	MinCpu    types.String `tfsdk:"min_cpu"`
	MaxCpu    types.String `tfsdk:"max_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
	MaxMemory types.String `tfsdk:"max_memory"`
	LogLevel  types.String `tfsdk:"log_level"`
	Env       types.Map    `tfsdk:"env"`
}

func (a AddOnsByokActuatorModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_cpu":    types.StringType,
			"max_cpu":    types.StringType,
			"min_memory": types.StringType,
			"max_memory": types.StringType,
			"log_level":  types.StringType,
			"env":        types.MapType{ElemType: types.StringType},
		},
	}
}

// Add Ons -> Byok -> Config -> Middlebox //

type AddOnsByokMiddleboxModel struct {
	Enabled            types.Bool  `tfsdk:"enabled"`
	BandwidthAlertMbps types.Int32 `tfsdk:"bandwidth_alert_mbps"`
}

func (a AddOnsByokMiddleboxModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled":              types.BoolType,
			"bandwidth_alert_mbps": types.Int32Type,
		},
	}
}

// Add Ons -> Byok -> Config -> Common //

type AddOnsByokCommonModel struct {
	DeploymentReplicas types.Int32  `tfsdk:"deployment_replicas"`
	Pdb                types.Object `tfsdk:"pdb"`
}

func (a AddOnsByokCommonModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deployment_replicas": types.Int32Type,
			"pdb":                 AddOnsByokCommonPdbModel{}.AttributeTypes(),
		},
	}
}

// Add Ons -> Byok -> Config -> Common -> pdb //

type AddOnsByokCommonPdbModel struct {
	MaxUnavailable types.Int32 `tfsdk:"max_unavailable"`
}

func (a AddOnsByokCommonPdbModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"max_unavailable": types.Int32Type,
		},
	}
}

// Add Ons -> Byok -> Config -> Longhorn //

type AddOnsByokLonghornModel struct {
	Replicas types.Int32 `tfsdk:"replicas"`
}

func (a AddOnsByokLonghornModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"replicas": types.Int32Type,
		},
	}
}

// Add Ons -> Byok -> Config -> Ingress //

type AddOnsByokIngressModel struct {
	Cpu           types.String  `tfsdk:"cpu"`
	Memory        types.String  `tfsdk:"memory"`
	TargetPercent types.Float32 `tfsdk:"target_percent"`
}

func (a AddOnsByokIngressModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cpu":            types.StringType,
			"memory":         types.StringType,
			"target_percent": types.Float32Type,
		},
	}
}

// Add Ons -> Byok -> Config -> Istio //

type AddOnsByokIstioModel struct {
	Istiod         types.Object `tfsdk:"istiod"`
	IngressGateway types.Object `tfsdk:"ingress_gateway"`
	Sidecar        types.Object `tfsdk:"sidecar"`
}

func (a AddOnsByokIstioModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"istiod":          AddOnsByokIstioIstiodModel{}.AttributeTypes(),
			"ingress_gateway": AddOnsByokIstioIngressGatewayModel{}.AttributeTypes(),
			"sidecar":         AddOnsByokIstioSidecarModel{}.AttributeTypes(),
		},
	}
}

// Add Ons -> Byok -> Config -> Istio -> Istiod //

type AddOnsByokIstioIstiodModel struct {
	Replicas  types.Int32  `tfsdk:"replicas"`
	MinCpu    types.String `tfsdk:"min_cpu"`
	MaxCpu    types.String `tfsdk:"max_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
	MaxMemory types.String `tfsdk:"max_memory"`
	Pdb       types.Int32  `tfsdk:"pdb"`
}

func (a AddOnsByokIstioIstiodModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"replicas":   types.Int32Type,
			"min_cpu":    types.StringType,
			"max_cpu":    types.StringType,
			"min_memory": types.StringType,
			"max_memory": types.StringType,
			"pdb":        types.Int32Type,
		},
	}
}

// Add Ons -> Byok -> Config -> Istio -> Ingress Gateway //

type AddOnsByokIstioIngressGatewayModel struct {
	Replicas  types.Int32  `tfsdk:"replicas"`
	MaxCpu    types.String `tfsdk:"max_cpu"`
	MaxMemory types.String `tfsdk:"max_memory"`
}

func (a AddOnsByokIstioIngressGatewayModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"replicas":   types.Int32Type,
			"max_cpu":    types.StringType,
			"max_memory": types.StringType,
		},
	}
}

// Add Ons -> Byok -> Config -> Istio -> Sidecar //

type AddOnsByokIstioSidecarModel struct {
	MinCpu    types.String `tfsdk:"min_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
}

func (a AddOnsByokIstioSidecarModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_cpu":    types.StringType,
			"min_memory": types.StringType,
		},
	}
}

// Add Ons -> Byok -> Config -> Log Splitter //

type AddOnsByokLogSplitterModel struct {
	MinCpu        types.String `tfsdk:"min_cpu"`
	MaxCpu        types.String `tfsdk:"max_cpu"`
	MinMemory     types.String `tfsdk:"min_memory"`
	MaxMemory     types.String `tfsdk:"max_memory"`
	MemBufferSize types.String `tfsdk:"mem_buffer_size"`
	PerPodRate    types.Int32  `tfsdk:"per_pod_rate"`
}

func (a AddOnsByokLogSplitterModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_cpu":         types.StringType,
			"max_cpu":         types.StringType,
			"min_memory":      types.StringType,
			"max_memory":      types.StringType,
			"mem_buffer_size": types.StringType,
			"per_pod_rate":    types.Int32Type,
		},
	}
}

// Add Ons -> Byok -> Config -> Monitoring //

type AddOnsByokMonitoringModel struct {
	MinMemory        types.String `tfsdk:"min_memory"`
	MaxMemory        types.String `tfsdk:"max_memory"`
	KubeStateMetrics types.Object `tfsdk:"kube_state_metrics"`
	Prometheus       types.Object `tfsdk:"prometheus"`
}

func (a AddOnsByokMonitoringModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_memory":         types.StringType,
			"max_memory":         types.StringType,
			"kube_state_metrics": AddOnsByokMonitoringKubeStateMetricsModel{}.AttributeTypes(),
			"prometheus":         AddOnsByokMonitoringPrometheusModel{}.AttributeTypes(),
		},
	}
}

// Add Ons -> Byok -> Config -> Monitoring -> Kube State Metrics //

type AddOnsByokMonitoringKubeStateMetricsModel struct {
	MinMemory types.String `tfsdk:"min_memory"`
}

func (a AddOnsByokMonitoringKubeStateMetricsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_memory": types.StringType,
		},
	}
}

// Add Ons -> Byok -> Config -> Monitoring -> Prometheus //

type AddOnsByokMonitoringPrometheusModel struct {
	Main types.Object `tfsdk:"main"`
}

func (a AddOnsByokMonitoringPrometheusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"main": AddOnsByokMonitoringPrometheusMainModel{}.AttributeTypes(),
		},
	}
}

// Add Ons -> Byok -> Config -> Monitoring -> Prometheus -> Main //

type AddOnsByokMonitoringPrometheusMainModel struct {
	Storage types.String `tfsdk:"storage"`
}

func (a AddOnsByokMonitoringPrometheusMainModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"storage": types.StringType,
		},
	}
}

// Add Ons -> Byok -> Config -> Redis //

type AddOnsByokRedisStringModel struct {
	MinCpu    types.String `tfsdk:"min_cpu"`
	MaxCpu    types.String `tfsdk:"max_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
	MaxMemory types.String `tfsdk:"max_memory"`
	Storage   types.String `tfsdk:"storage"`
}

func (a AddOnsByokRedisStringModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_cpu":    types.StringType,
			"max_cpu":    types.StringType,
			"min_memory": types.StringType,
			"max_memory": types.StringType,
			"storage":    types.StringType,
		},
	}
}

type AddOnsByokRedisIntModel struct {
	MinCpu    types.String `tfsdk:"min_cpu"`
	MaxCpu    types.String `tfsdk:"max_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
	MaxMemory types.String `tfsdk:"max_memory"`
	Storage   types.Int32  `tfsdk:"storage"`
}

func (a AddOnsByokRedisIntModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_cpu":    types.StringType,
			"max_cpu":    types.StringType,
			"min_memory": types.StringType,
			"max_memory": types.StringType,
			"storage":    types.Int32Type,
		},
	}
}

// Add Ons -> Byok -> Config -> Tempo Agent //

type AddOnsByokTempoAgentModel struct {
	MinCpu    types.String `tfsdk:"min_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
}

func (a AddOnsByokTempoAgentModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_cpu":    types.StringType,
			"min_memory": types.StringType,
		},
	}
}

// Add Ons -> Byok -> Config -> Internal DNS //

type AddOnsByokInternalDnsModel struct {
	MinCpu    types.String `tfsdk:"min_cpu"`
	MaxCpu    types.String `tfsdk:"max_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
	MaxMemory types.String `tfsdk:"max_memory"`
}

func (a AddOnsByokInternalDnsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_cpu":    types.StringType,
			"max_cpu":    types.StringType,
			"min_memory": types.StringType,
			"max_memory": types.StringType,
		},
	}
}

// Status //

type StatusModel struct {
	OidcProviderUrl types.String `tfsdk:"oidc_provider_url"`
	ServerUrl       types.String `tfsdk:"server_url"`
	HomeLocation    types.String `tfsdk:"home_location"`
	AddOns          types.List   `tfsdk:"add_ons"`
}

func (s StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"oidc_provider_url": types.StringType,
			"server_url":        types.StringType,
			"home_location":     types.StringType,
			"add_ons":           types.ListType{ElemType: StatusAddOnsModel{}.AttributeTypes()},
		},
	}
}

// Status -> Add Ons //

type StatusAddOnsModel struct {
	Dashboard           types.List `tfsdk:"dashboard"`
	AwsWorkloadIdentity types.List `tfsdk:"aws_workload_identity"`
	Metrics             types.List `tfsdk:"metrics"`
	Logs                types.List `tfsdk:"logs"`
	AwsECR              types.List `tfsdk:"aws_ecr"`
	AwsEFS              types.List `tfsdk:"aws_efs"`
	AwsELB              types.List `tfsdk:"aws_elb"`
}

func (s StatusAddOnsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"dashboard":             types.ListType{ElemType: StatusAddOnsDashboardModel{}.AttributeTypes()},
			"aws_workload_identity": types.ListType{ElemType: StatusAddOnsAwsWorkloadIdentityModel{}.AttributeTypes()},
			"metrics":               types.ListType{ElemType: StatusAddOnsMetricsModel{}.AttributeTypes()},
			"logs":                  types.ListType{ElemType: StatusAddOnsLogsModel{}.AttributeTypes()},
			"aws_ecr":               types.ListType{ElemType: StatusAddOnsAwsStatusModel{}.AttributeTypes()},
			"aws_efs":               types.ListType{ElemType: StatusAddOnsAwsStatusModel{}.AttributeTypes()},
			"aws_elb":               types.ListType{ElemType: StatusAddOnsAwsStatusModel{}.AttributeTypes()},
		},
	}
}

// Status -> Add Ons -> Dashboard //

type StatusAddOnsDashboardModel struct {
	Url types.String `tfsdk:"url"`
}

func (s StatusAddOnsDashboardModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"url": types.StringType,
		},
	}
}

// Status -> Add Ons -> AWS Workload Identity //

type StatusAddOnsAwsWorkloadIdentityModel struct {
	OidcProviderConfig types.List   `tfsdk:"oidc_provider_config"`
	TrustPolicy        types.String `tfsdk:"trust_policy"`
}

func (s StatusAddOnsAwsWorkloadIdentityModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"oidc_provider_config": types.ListType{ElemType: StatusAddOnsAwsWorkloadIdentityOidcProviderConfigModel{}.AttributeTypes()},
			"trust_policy":         types.StringType,
		},
	}
}

// Status -> Add Ons -> AWS Workload Identity -> Oidc Provider Config //

type StatusAddOnsAwsWorkloadIdentityOidcProviderConfigModel struct {
	ProviderUrl types.String `tfsdk:"provider_url"`
	Audience    types.String `tfsdk:"audience"`
}

func (s StatusAddOnsAwsWorkloadIdentityOidcProviderConfigModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"provider_url": types.StringType,
			"audience":     types.StringType,
		},
	}
}

// Status -> Add Ons -> Metrics //

type StatusAddOnsMetricsModel struct {
	PrometheusEndpoint types.String `tfsdk:"prometheus_endpoint"`
	RemoteWriteConfig  types.String `tfsdk:"remote_write_config"`
}

func (s StatusAddOnsMetricsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"prometheus_endpoint": types.StringType,
			"remote_write_config": types.StringType,
		},
	}
}

// Status -> Add Ons -> Logs //

type StatusAddOnsLogsModel struct {
	LokiAddress types.String `tfsdk:"loki_address"`
}

func (s StatusAddOnsLogsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"loki_address": types.StringType,
		},
	}
}

// Status -> Add Ons -> AWS Config //

type StatusAddOnsAwsStatusModel struct {
	TrustPolicy types.String `tfsdk:"trust_policy"`
}

func (s StatusAddOnsAwsStatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"trust_policy": types.StringType,
		},
	}
}

/*** Common Models ***/

type NetworkingModel struct {
	ServiceNetwork types.String `tfsdk:"service_network"`
	PodNetwork     types.String `tfsdk:"pod_network"`
	DnsForwarder   types.String `tfsdk:"dns_forwarder"`
}

func (n NetworkingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"service_network": types.StringType,
			"pod_network":     types.StringType,
			"dns_forwarder":   types.StringType,
		},
	}
}

type AutoscalerModel struct {
	Expander             types.Set     `tfsdk:"expander"`
	UnneededTime         types.String  `tfsdk:"unneeded_time"`
	UnreadyTime          types.String  `tfsdk:"unready_time"`
	UtilizationThreshold types.Float64 `tfsdk:"utilization_threshold"`
}

func (a AutoscalerModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"expander":              types.SetType{ElemType: types.StringType},
			"unneeded_time":         types.StringType,
			"unready_time":          types.StringType,
			"utilization_threshold": types.Float64Type,
		},
	}
}
