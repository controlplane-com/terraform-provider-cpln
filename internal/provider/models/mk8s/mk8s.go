package mk8s

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Main Models ***/

// Firewall //

type FirewallModel struct {
	SourceCIDR  types.String `tfsdk:"source_cidr"`
	Description types.String `tfsdk:"description"`
}

// Generic Provider //

type GenericProviderModel struct {
	Location   types.String                   `tfsdk:"location"`
	Networking []NetworkingModel              `tfsdk:"networking"`
	NodePools  []GenericProviderNodePoolModel `tfsdk:"node_pool"`
}

// Generic Provider -> Node Pool //

type GenericProviderNodePoolModel struct {
	Name   types.String                        `tfsdk:"name"`
	Labels types.Map                           `tfsdk:"labels"`
	Taints []GenericProviderNodePoolTaintModel `tfsdk:"taint"`
}

// Generic Provider -> Node Pool -> Taint //

type GenericProviderNodePoolTaintModel struct {
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
	Effect types.String `tfsdk:"effect"`
}

// Hetzner Provider //

type HetznerProviderModel struct {
	Region                   types.String                   `tfsdk:"region"`
	HetznerLabels            types.Map                      `tfsdk:"hetzner_labels"`
	Networking               []NetworkingModel              `tfsdk:"networking"`
	PreInstallScript         types.String                   `tfsdk:"pre_install_script"`
	TokenSecretLink          types.String                   `tfsdk:"token_secret_link"`
	NetworkId                types.String                   `tfsdk:"network_id"`
	FirewallId               types.String                   `tfsdk:"firewall_id"`
	NodePools                []HetznerProviderNodePoolModel `tfsdk:"node_pool"`
	DedicatedServerNodePools []GenericProviderNodePoolModel `tfsdk:"dedicated_server_node_pool"`
	Image                    types.String                   `tfsdk:"image"`
	SshKey                   types.String                   `tfsdk:"ssh_key"`
	Autoscaler               []AutoscalerModel              `tfsdk:"autoscaler"`
	FloatingIpSelector       types.Map                      `tfsdk:"floating_ip_selector"`
}

// Hetzner Provider -> Node Pool //

type HetznerProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	ServerType    types.String `tfsdk:"server_type"`
	OverrideImage types.String `tfsdk:"override_image"`
	MinSize       types.Int32  `tfsdk:"min_size"`
	MaxSize       types.Int32  `tfsdk:"max_size"`
}

// AWS Provider //

type AwsProviderModel struct {
	Region               types.String                     `tfsdk:"region"`
	AwsTags              types.Map                        `tfsdk:"aws_tags"`
	SkipCreateRoles      types.Bool                       `tfsdk:"skip_create_roles"`
	Networking           []NetworkingModel                `tfsdk:"networking"`
	PreInstallScript     types.String                     `tfsdk:"pre_install_script"`
	Image                []AwsProviderAmiModel            `tfsdk:"image"`
	DeployRoleArn        types.String                     `tfsdk:"deploy_role_arn"`
	DeployRoleChain      []AwsProviderAssumeRoleLinkModel `tfsdk:"deploy_role_chain"`
	VpcId                types.String                     `tfsdk:"vpc_id"`
	KeyPair              types.String                     `tfsdk:"key_pair"`
	DiskEncryptionKeyArn types.String                     `tfsdk:"disk_encryption_key_arn"`
	SecurityGroupIds     types.Set                        `tfsdk:"security_group_ids"`
	ExtraNodePolicies    types.Set                        `tfsdk:"extra_node_policies"`
	NodePools            []AwsProviderNodePoolModel       `tfsdk:"node_pool"`
	Autoscaler           []AutoscalerModel                `tfsdk:"autoscaler"`
}

// AWS Provider -> AMI //

type AwsProviderAmiModel struct {
	Recommended types.String `tfsdk:"recommended"`
	Exact       types.String `tfsdk:"exact"`
}

// AWS Provider -> Deploy Role Chain //

type AwsProviderAssumeRoleLinkModel struct {
	RoleArn           types.String `tfsdk:"role_arn"`
	ExternalId        types.String `tfsdk:"external_id"`
	SessionNamePrefix types.String `tfsdk:"session_name_prefix"`
}

// AWS Provider -> Node Pool //

type AwsProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	InstanceTypes                       types.Set             `tfsdk:"instance_types"`
	OverrideImage                       []AwsProviderAmiModel `tfsdk:"override_image"`
	BootDiskSize                        types.Int32           `tfsdk:"boot_disk_size"`
	MinSize                             types.Int32           `tfsdk:"min_size"`
	MaxSize                             types.Int32           `tfsdk:"max_size"`
	OnDemandBaseCapacity                types.Int32           `tfsdk:"on_demand_base_capacity"`
	OnDemandPercentageAboveBaseCapacity types.Int32           `tfsdk:"on_demand_percentage_above_base_capacity"`
	SpotAllocationStrategy              types.String          `tfsdk:"spot_allocation_strategy"`
	SubnetIds                           types.Set             `tfsdk:"subnet_ids"`
	ExtraSecurityGroupIds               types.Set             `tfsdk:"extra_security_group_ids"`
}

// Linode Provider //

type LinodeProviderModel struct {
	Region           types.String                  `tfsdk:"region"`
	TokenSecretLink  types.String                  `tfsdk:"token_secret_link"`
	FirewallId       types.String                  `tfsdk:"firewall_id"`
	NodePools        []LinodeProviderNodePoolModel `tfsdk:"node_pool"`
	Image            types.String                  `tfsdk:"image"`
	AuthorizedUsers  types.Set                     `tfsdk:"authorized_users"`
	AuthorizedKeys   types.Set                     `tfsdk:"authorized_keys"`
	VpcId            types.String                  `tfsdk:"vpc_id"`
	PreInstallScript types.String                  `tfsdk:"pre_install_script"`
	Networking       []NetworkingModel             `tfsdk:"networking"`
	Autoscaler       []AutoscalerModel             `tfsdk:"autoscaler"`
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

// Oblivus Provider //

type OblivusProviderModel struct {
	Datacenter        types.String                   `tfsdk:"datacenter"`
	TokenSecretLink   types.String                   `tfsdk:"token_secret_link"`
	NodePools         []OblivusProviderNodePoolModel `tfsdk:"node_pool"`
	SshKeys           types.Set                      `tfsdk:"ssh_keys"`
	UnmanagedNodePool []GenericProviderNodePoolModel `tfsdk:"unmanaged_node_pool"`
	Autoscaler        []AutoscalerModel              `tfsdk:"autoscaler"`
	PreInstallScript  types.String                   `tfsdk:"pre_install_script"`
}

// Oblivus Provider -> Node Pool //

type OblivusProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	MinSize types.Int32  `tfsdk:"min_size"`
	MaxSize types.Int32  `tfsdk:"max_size"`
	Flavor  types.String `tfsdk:"flavor"`
}

// Lambdalabs Provider //

type LambdalabsProviderModel struct {
	Region             types.String                      `tfsdk:"region"`
	TokenSecretLink    types.String                      `tfsdk:"token_secret_link"`
	NodePools          []LambdalabsProviderNodePoolModel `tfsdk:"node_pool"`
	SshKey             types.String                      `tfsdk:"ssh_key"`
	UnmanagedNodePools []GenericProviderNodePoolModel    `tfsdk:"unmanaged_node_pool"`
	Autoscaler         []AutoscalerModel                 `tfsdk:"autoscaler"`
	FileSystems        types.Set                         `tfsdk:"file_systems"`
	PreInstallScript   types.String                      `tfsdk:"pre_install_script"`
}

// Lambdalabs Provider -> Node Pool //

type LambdalabsProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	MinSize      types.Int32  `tfsdk:"min_size"`
	MaxSize      types.Int32  `tfsdk:"max_size"`
	InstanceType types.String `tfsdk:"instance_type"`
}

// Paperspace Provider //

type PaperspaceProviderModel struct {
	Region             types.String                      `tfsdk:"region"`
	TokenSecretLink    types.String                      `tfsdk:"token_secret_link"`
	SharedDrives       types.Set                         `tfsdk:"shared_drives"`
	NodePools          []PaperspaceProviderNodePoolModel `tfsdk:"node_pool"`
	Autoscaler         []AutoscalerModel                 `tfsdk:"autoscaler"`
	UnmanagedNodePools []GenericProviderNodePoolModel    `tfsdk:"unmanaged_node_pool"`
	PreInstallScript   types.String                      `tfsdk:"pre_install_script"`
	UserIds            types.Set                         `tfsdk:"user_ids"`
	NetworkId          types.String                      `tfsdk:"network_id"`
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

// Ephemeral Provider //

type EphemeralProviderModel struct {
	Location  types.String                     `tfsdk:"location"`
	NodePools []EphemeralProviderNodePoolModel `tfsdk:"node_pool"`
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

// Triton Provider //

type TritonProviderModel struct {
	Connection       []TritonProviderConnectionModel   `tfsdk:"connection"`
	Networking       []NetworkingModel                 `tfsdk:"networking"`
	PreInstallScript types.String                      `tfsdk:"pre_install_script"`
	Location         types.String                      `tfsdk:"location"`
	LoadBalancer     []TritonProviderLoadBalancerModel `tfsdk:"load_balancer"`
	PrivateNetworkId types.String                      `tfsdk:"private_network_id"`
	FirewallEnabled  types.Bool                        `tfsdk:"firewall_enabled"`
	NodePools        []TritonProviderNodePoolModel     `tfsdk:"node_pool"`
	ImageId          types.String                      `tfsdk:"image_id"`
	SshKeys          types.Set                         `tfsdk:"ssh_keys"`
	Autoscaler       []AutoscalerModel                 `tfsdk:"autoscaler"`
}

// Triton Provider -> Connection //

type TritonProviderConnectionModel struct {
	Url                  types.String `tfsdk:"url"`
	Account              types.String `tfsdk:"account"`
	User                 types.String `tfsdk:"user"`
	PrivateKeySecretLink types.String `tfsdk:"private_key_secret_link"`
}

// Triton Provider -> Load Balancer //

type TritonProviderLoadBalancerModel struct {
	Manual  []TritonProviderLoadBalancerManualModel  `tfsdk:"manual"`
	Gateway []TritonProviderLoadBalancerGatewayModel `tfsdk:"gateway"`
}

// Triton Provider -> Load Balancer -> Manual //

type TritonProviderLoadBalancerManualModel struct {
	PackageId         types.String                                   `tfsdk:"package_id"`
	ImageId           types.String                                   `tfsdk:"image_id"`
	PublicNetworkId   types.String                                   `tfsdk:"public_network_id"`
	PrivateNetworkIds types.Set                                      `tfsdk:"private_network_ids"`
	Metadata          types.Map                                      `tfsdk:"metadata"`
	Tags              types.Map                                      `tfsdk:"tags"`
	Logging           []TritonProviderLoadBalancerManualLoggingModel `tfsdk:"logging"`
	Count             types.Int32                                    `tfsdk:"count"`
	CnsInternalDomain types.String                                   `tfsdk:"cns_internal_domain"`
	CnsPublicDomain   types.String                                   `tfsdk:"cns_public_domain"`
}

// Triton Provider -> Load Balancer -> Manual -> Logging //

type TritonProviderLoadBalancerManualLoggingModel struct {
	NodePort       types.Int32  `tfsdk:"node_port"`
	ExternalSyslog types.String `tfsdk:"external_syslog"`
}

// Triton Provider -> Load Balancer -> Gateway //

type TritonProviderLoadBalancerGatewayModel struct{}

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

// Azure Provider //

type AzureProviderModel struct {
	Location         types.String                 `tfsdk:"location"`
	SubscriptionId   types.String                 `tfsdk:"subscription_id"`
	SdkSecretLink    types.String                 `tfsdk:"sdk_secret_link"`
	ResourceGroup    types.String                 `tfsdk:"resource_group"`
	Networking       []NetworkingModel            `tfsdk:"networking"`
	PreInstallScript types.String                 `tfsdk:"pre_install_script"`
	Image            []AzureProviderImageModel    `tfsdk:"image"`
	SshKeys          types.Set                    `tfsdk:"ssh_keys"`
	NetworkId        types.String                 `tfsdk:"network_id"`
	Tags             types.Map                    `tfsdk:"tags"`
	NodePools        []AzureProviderNodePoolModel `tfsdk:"node_pool"`
	Autoscaler       []AutoscalerModel            `tfsdk:"autoscaler"`
}

// Azure Provider -> Image //

type AzureProviderImageModel struct {
	Recommended types.String                       `tfsdk:"recommended"`
	Reference   []AzureProviderImageReferenceModel `tfsdk:"reference"`
}

// Azure Provider -> Image -> Reference //

type AzureProviderImageReferenceModel struct {
	Publisher types.String `tfsdk:"publisher"`
	Offer     types.String `tfsdk:"offer"`
	Sku       types.String `tfsdk:"sku"`
	Version   types.String `tfsdk:"version"`
}

// Azure Provider -> Node Pool //

type AzureProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	Size          types.String              `tfsdk:"size"`
	SubnetId      types.String              `tfsdk:"subnet_id"`
	Zones         types.Set                 `tfsdk:"zones"`
	OverrideImage []AzureProviderImageModel `tfsdk:"override_image"`
	BootDiskSize  types.Int32               `tfsdk:"boot_disk_size"`
	MinSize       types.Int32               `tfsdk:"min_size"`
	MaxSize       types.Int32               `tfsdk:"max_size"`
}

// Digital Ocean Provider //

type DigitalOceanProviderModel struct {
	Region           types.String                        `tfsdk:"region"`
	DigitalOceanTags types.Set                           `tfsdk:"digital_ocean_tags"`
	Networking       []NetworkingModel                   `tfsdk:"networking"`
	PreInstallScript types.String                        `tfsdk:"pre_install_script"`
	TokenSecretLink  types.String                        `tfsdk:"token_secret_link"`
	VpcId            types.String                        `tfsdk:"vpc_id"`
	NodePools        []DigitalOceanProviderNodePoolModel `tfsdk:"node_pool"`
	Image            types.String                        `tfsdk:"image"`
	SshKeys          types.Set                           `tfsdk:"ssh_keys"`
	ExtraSshKeys     types.Set                           `tfsdk:"extra_ssh_keys"`
	Autoscaler       []AutoscalerModel                   `tfsdk:"autoscaler"`
	ReservedIps      types.Set                           `tfsdk:"reserved_ips"`
}

// Digital Ocean Provider -> Node Pool //

type DigitalOceanProviderNodePoolModel struct {
	GenericProviderNodePoolModel
	DropletSize   types.String `tfsdk:"droplet_size"`
	OverrideImage types.String `tfsdk:"override_image"`
	MinSize       types.Int32  `tfsdk:"min_size"`
	MaxSize       types.Int32  `tfsdk:"max_size"`
}

// Add Ons //

type AddOnsModel struct {
	Dashboard             types.Bool                        `tfsdk:"dashboard"`
	AzureWorkloadIdentity []AddOnAzureWorkloadIdentityModel `tfsdk:"azure_workload_identity"`
	AwsWorkloadIdentity   types.Bool                        `tfsdk:"aws_workload_identity"`
	LocalPathStorage      types.Bool                        `tfsdk:"local_path_storage"`
	Metrics               []AddOnsMetricsModel              `tfsdk:"metrics"`
	Logs                  []AddOnsLogsModel                 `tfsdk:"logs"`
	Nvidia                []AddOnsNvidiaModel               `tfsdk:"nvidia"`
	AwsEFS                []AddOnsHasRoleArnModel           `tfsdk:"aws_efs"`
	AwsECR                []AddOnsHasRoleArnModel           `tfsdk:"aws_ecr"`
	AwsELB                []AddOnsHasRoleArnModel           `tfsdk:"aws_elb"`
	AzureACR              []AddOnsAzureAcrModel             `tfsdk:"azure_acr"`
	Sysbox                types.Bool                        `tfsdk:"sysbox"`
}

// Add Ons -> Azure Workload Identity //

type AddOnAzureWorkloadIdentityModel struct {
	TenantId types.String `tfsdk:"tenant_id"`
}

// Add Ons -> Metrics //

type AddOnsMetricsModel struct {
	KubeState       types.Bool                          `tfsdk:"kube_state"`
	CoreDns         types.Bool                          `tfsdk:"core_dns"`
	Kubelet         types.Bool                          `tfsdk:"kubelet"`
	Apiserver       types.Bool                          `tfsdk:"api_server"`
	NodeExporter    types.Bool                          `tfsdk:"node_exporter"`
	Cadvisor        types.Bool                          `tfsdk:"cadvisor"`
	ScrapeAnnotated []AddOnsMetricsScrapeAnnotatedModel `tfsdk:"scrape_annotated"`
}

// Add Ons -> Metrics -> Scrape Annotated //

type AddOnsMetricsScrapeAnnotatedModel struct {
	IntervalSeconds   types.Int32  `tfsdk:"interval_seconds"`
	IncludeNamespaces types.String `tfsdk:"include_namespaces"`
	ExcludeNamespaces types.String `tfsdk:"exclude_namespaces"`
	RetainLabels      types.String `tfsdk:"retain_labels"`
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

// Add Ons -> Nvidia //

type AddOnsNvidiaModel struct {
	TaintGpuNodes types.Bool `tfsdk:"taint_gpu_nodes"`
}

// Add Ons -> Has Role Arn //

type AddOnsHasRoleArnModel struct {
	RoleArn types.String `tfsdk:"role_arn"`
}

// Add Ons -> Azure ACR //

type AddOnsAzureAcrModel struct {
	ClientId types.String `tfsdk:"client_id"`
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

type AutoscalerModel struct {
	Expander             types.Set     `tfsdk:"expander"`
	UnneededTime         types.String  `tfsdk:"unneeded_time"`
	UnreadyTime          types.String  `tfsdk:"unready_time"`
	UtilizationThreshold types.Float64 `tfsdk:"utilization_threshold"`
}
