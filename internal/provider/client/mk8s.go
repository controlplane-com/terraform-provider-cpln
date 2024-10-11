package cpln

import "fmt"

type Mk8s struct {
	Base
	Alias       *string     `json:"alias,omitempty"`
	Spec        *Mk8sSpec   `json:"spec,omitempty"`
	SpecReplace *Mk8sSpec   `json:"$replace/spec,omitempty"`
	Status      *Mk8sStatus `json:"status,omitempty"`
}

type Mk8sSpec struct {
	Version  *string             `json:"version,omitempty"`
	Firewall *[]Mk8sFirewallRule `json:"firewall,omitempty"`
	Provider *Mk8sProvider       `json:"provider,omitempty"`
	AddOns   *Mk8sSpecAddOns     `json:"addOns,omitempty"`
}

type Mk8sStatus struct {
	OidcProviderUrl *string           `json:"oidcProviderUrl,omitempty"`
	ServerUrl       *string           `json:"serverUrl,omitempty"`
	HomeLocation    *string           `json:"homeLocation,omitempty"`
	AddOns          *Mk8sStatusAddOns `json:"addOns,omitempty"`
}

/*** Spec ***/

type Mk8sFirewallRule struct {
	SourceCIDR  *string `json:"sourceCIDR,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Mk8sProvider struct {
	Generic    *Mk8sGenericProvider    `json:"generic,omitempty"`
	Hetzner    *Mk8sHetznerProvider    `json:"hetzner,omitempty"`
	Aws        *Mk8sAwsProvider        `json:"aws,omitempty"`
	Linode     *Mk8sLinodeProvider     `json:"linode,omitempty"`
	Oblivus    *Mk8sOblivusProvider    `json:"oblivus,omitempty"`
	Lambdalabs *Mk8sLambdalabsProvider `json:"lambdalabs,omitempty"`
	Paperspace *Mk8sPaperspaceProvider `json:"paperspace,omitempty"`
	Ephemeral  *Mk8sEphemeralProvder   `json:"ephemeral,omitempty"`
}

type Mk8sSpecAddOns struct {
	Dashboard             *Mk8sNonCustomizableAddonConfig       `json:"dashboard,omitempty"`
	AzureWorkloadIdentity *Mk8sAzureWorkloadIdentityAddOnConfig `json:"azureWorkloadIdentity,omitempty"`
	AwsWorkloadIdentity   *Mk8sNonCustomizableAddonConfig       `json:"awsWorkloadIdentity,omitempty"`
	LocalPathStorage      *Mk8sNonCustomizableAddonConfig       `json:"localPathStorage,omitempty"`
	Metrics               *Mk8sMetricsAddOnConfig               `json:"metrics,omitempty"`
	Logs                  *Mk8sLogsAddOnConfig                  `json:"logs,omitempty"`
	Nvidia                *Mk8sNvidiaAddOnConfig                `json:"nvidia,omitempty"`
	AwsEFS                *Mk8sAwsAddOnConfig                   `json:"awsEFS,omitempty"`
	AwsECR                *Mk8sAwsAddOnConfig                   `json:"awsECR,omitempty"`
	AwsELB                *Mk8sAwsAddOnConfig                   `json:"awsELB,omitempty"`
	AzureACR              *Mk8sAzureACRAddOnConfig              `json:"azureACR,omitempty"`
	Sysbox                *Mk8sNonCustomizableAddonConfig       `json:"sysbox,omitempty"`
}

// Providers //

type Mk8sGenericProvider struct {
	Location   *string               `json:"location,omitempty"`
	Networking *Mk8sNetworkingConfig `json:"networking,omitempty"`
	NodePools  *[]Mk8sGenericPool    `json:"nodePools,omitempty"`
}

type Mk8sHetznerProvider struct {
	Region                   *string                 `json:"region,omitempty"`
	HetznerLabels            *map[string]interface{} `json:"hetznerLabels,omitempty"`
	Networking               *Mk8sNetworkingConfig   `json:"networking,omitempty"`
	PreInstallScript         *string                 `json:"preInstallScript,omitempty"`
	TokenSecretLink          *string                 `json:"tokenSecretLink,omitempty"`
	NetworkId                *string                 `json:"networkId,omitempty"`
	FirewallId               *string                 `json:"firewallId,omitempty"`
	NodePools                *[]Mk8sHetznerPool      `json:"nodePools,omitempty"`
	DedicatedServerNodePools *[]Mk8sGenericPool      `json:"dedicatedServerNodePools,omitempty"`
	Image                    *string                 `json:"image,omitempty"`
	SshKey                   *string                 `json:"sshKey,omitempty"`
	Autoscaler               *Mk8sAutoscalerConfig   `json:"autoscaler,omitempty"`
	FloatingIpSelector       *map[string]interface{} `json:"floatingIPSelector,omitempty"`
}

type Mk8sAwsProvider struct {
	Region               *string                 `json:"region,omitempty"`
	AwsTags              *map[string]interface{} `json:"awsTags,omitempty"`
	SkipCreateRoles      *bool                   `json:"skipCreateRoles,omitempty"`
	Networking           *Mk8sNetworkingConfig   `json:"networking,omitempty"`
	PreInstallScript     *string                 `json:"preInstallScript,omitempty"`
	Image                *Mk8sAwsAmi             `json:"image,omitempty"`
	DeployRoleArn        *string                 `json:"deployRoleArn,omitempty"`
	VpcId                *string                 `json:"vpcId,omitempty"`
	KeyPair              *string                 `json:"keyPair,omitempty"`
	DiskEncryptionKeyArn *string                 `json:"diskEncryptionKeyArn,omitempty"`
	SecurityGroupIds     *[]string               `json:"securityGroupIds,omitempty"`
	NodePools            *[]Mk8sAwsPool          `json:"nodePools,omitempty"`
	Autoscaler           *Mk8sAutoscalerConfig   `json:"autoscaler,omitempty"`
}

type Mk8sLinodeProvider struct {
	Region           *string               `json:"region,omitempty"`
	TokenSecretLink  *string               `json:"tokenSecretLink,omitempty"`
	FirewallId       *string               `json:"firewallId,omitempty"`
	NodePools        *[]Mk8sLinodePool     `json:"nodePools,omitempty"`
	Image            *string               `json:"image,omitempty"`
	AuthorizedUsers  *[]string             `json:"authorizedUsers,omitempty"`
	AuthorizedKeys   *[]string             `json:"authorizedKeys,omitempty"`
	VpcId            *string               `json:"vpcId,omitempty"`
	PreInstallScript *string               `json:"preInstallScript,omitempty"`
	Networking       *Mk8sNetworkingConfig `json:"networking,omitempty"`
	Autoscaler       *Mk8sAutoscalerConfig `json:"autoscaler,omitempty"`
}

type Mk8sOblivusProvider struct {
	Datacenter         *string               `json:"datacenter,omitempty"`
	TokenSecretLink    *string               `json:"tokenSecretLink,omitempty"`
	NodePools          *[]Mk8sOblivusPool    `json:"nodePools,omitempty"`
	SshKeys            *[]string             `json:"sshKeys,omitempty"`
	UnmanagedNodePools *[]Mk8sGenericPool    `json:"unmanagedNodePools,omitempty"`
	Autoscaler         *Mk8sAutoscalerConfig `json:"autoscaler,omitempty"`
	PreInstallScript   *string               `json:"preInstallScript,omitempty"`
}

type Mk8sLambdalabsProvider struct {
	Region             *string               `json:"region,omitempty"`
	TokenSecretLink    *string               `json:"tokenSecretLink,omitempty"`
	NodePools          *[]Mk8sLambdalabsPool `json:"nodePools,omitempty"`
	SshKey             *string               `json:"sshKey,omitempty"`
	UnmanagedNodePools *[]Mk8sGenericPool    `json:"unmanagedNodePools,omitempty"`
	Autoscaler         *Mk8sAutoscalerConfig `json:"autoscaler,omitempty"`
	PreInstallScript   *string               `json:"preInstallScript,omitempty"`
}

type Mk8sPaperspaceProvider struct {
	Region             *string               `json:"region,omitempty"`
	TokenSecretLink    *string               `json:"tokenSecretLink,omitempty"`
	SharedDrives       *[]string             `json:"sharedDrives,omitempty"`
	NodePools          *[]Mk8sPaperspacePool `json:"nodePools,omitempty"`
	Autoscaler         *Mk8sAutoscalerConfig `json:"autoscaler,omitempty"`
	UnmanagedNodePools *[]Mk8sGenericPool    `json:"unmanagedNodePools,omitempty"`
	PreInstallScript   *string               `json:"preInstallScript,omitempty"`
	UserIds            *[]string             `json:"userIds,omitempty"`
	NetworkId          *string               `json:"networkId,omitempty"`
}

type Mk8sEphemeralProvder struct {
	Location  *string              `json:"location,omitempty"`
	NodePools *[]Mk8sEphemeralPool `json:"nodePools,omitempty"`
}

// Node Pools //

type Mk8sGenericPool struct {
	Name   *string                 `json:"name,omitempty"`
	Labels *map[string]interface{} `json:"labels,omitempty"`
	Taints *[]Mk8sTaint            `json:"taints,omitempty"`
}

type Mk8sHetznerPool struct {
	Mk8sGenericPool
	ServerType    *string `json:"serverType,omitempty"`
	OverrideImage *string `json:"overrideImage,omitempty"`
	MinSize       *int    `json:"minSize,omitempty"`
	MaxSize       *int    `json:"maxSize,omitempty"`
}

type Mk8sAwsPool struct {
	Mk8sGenericPool
	InstanceTypes                       *[]string   `json:"instanceTypes,omitempty"`
	OverrideImage                       *Mk8sAwsAmi `json:"overrideImage,omitempty"`
	BootDiskSize                        *int        `json:"bootDiskSize,omitempty"`
	MinSize                             *int        `json:"minSize,omitempty"`
	MaxSize                             *int        `json:"maxSize,omitempty"`
	OnDemandBaseCapacity                *int        `json:"onDemandBaseCapacity,omitempty"`
	OnDemandPercentageAboveBaseCapacity *int        `json:"onDemandPercentageAboveBaseCapacity,omitempty"`
	SpotAllocationStrategy              *string     `json:"spotAllocationStrategy,omitempty"`
	SubnetIds                           *[]string   `json:"subnetIds,omitempty"`
	ExtraSecurityGroupIds               *[]string   `json:"extraSecurityGroupIds,omitempty"`
}

type Mk8sLinodePool struct {
	Mk8sGenericPool
	ServerType    *string `json:"serverType,omitempty"`
	OverrideImage *string `json:"overrideImage,omitempty"`
	SubnetId      *string `json:"subnetId,omitempty"`
	MinSize       *int    `json:"minSize,omitempty"`
	MaxSize       *int    `json:"maxSize,omitempty"`
}

type Mk8sOblivusPool struct {
	Mk8sGenericPool
	MinSize *int    `json:"minSize,omitempty"`
	MaxSize *int    `json:"maxSize,omitempty"`
	Flavor  *string `json:"flavor,omitempty"`
}

type Mk8sLambdalabsPool struct {
	Mk8sGenericPool
	MinSize      *int    `json:"minSize,omitempty"`
	MaxSize      *int    `json:"maxSize,omitempty"`
	InstanceType *string `json:"instanceType,omitempty"`
}

type Mk8sPaperspacePool struct {
	Mk8sGenericPool
	MinSize      *int    `json:"minSize,omitempty"`
	MaxSize      *int    `json:"maxSize,omitempty"`
	PublicIpType *string `json:"publicIpType,omitempty"`
	BootDiskSize *int    `json:"bootDiskSize,omitempty"`
	MachineType  *string `json:"machineType,omitempty"`
}

type Mk8sEphemeralPool struct {
	Name   *string                 `json:"name,omitempty"`
	Labels *map[string]interface{} `json:"labels,omitempty"`
	Taints *[]Mk8sTaint            `json:"taints,omitempty"`
	Count  *int                    `json:"count,omitempty"`
	Arch   *string                 `json:"arch,omitempty"`
	Flavor *string                 `json:"flavor,omitempty"`
	Cpu    *string                 `json:"cpu,omitempty"`
	Memory *string                 `json:"memory,omitempty"`
}

// Provider Common //

type Mk8sNetworkingConfig struct {
	ServiceNetwork *string `json:"serviceNetwork,omitempty"`
	PodNetwork     *string `json:"podNetwork,omitempty"`
}

type Mk8sTaint struct {
	Key    *string `json:"key,omitempty"`
	Value  *string `json:"value,omitempty"`
	Effect *string `json:"effect,omitempty"`
}

type Mk8sAutoscalerConfig struct {
	Expander             *[]string `json:"expander,omitempty"`
	UnneededTime         *string   `json:"unneededTime,omitempty"`
	UnreadyTime          *string   `json:"unreadyTime,omitempty"`
	UtilizationThreshold *float64  `json:"utilizationThreshold,omitempty"`
}

// AWS Provider //

type Mk8sAwsAmi struct {
	Recommended *string `json:"recommended,omitempty"`
	Exact       *string `json:"exact,omitempty"`
}

// Add Ons //

type Mk8sAzureWorkloadIdentityAddOnConfig struct {
	TenantId *string `json:"tenantId,omitempty"`
}

type Mk8sMetricsAddOnConfig struct {
	KubeState       *bool                       `json:"kubeState,omitempty"`
	CoreDns         *bool                       `json:"coreDns,omitempty"`
	Kubelet         *bool                       `json:"kubelet,omitempty"`
	Apiserver       *bool                       `json:"apiserver,omitempty"`
	NodeExporter    *bool                       `json:"nodeExporter,omitempty"`
	Cadvisor        *bool                       `json:"cadvisor,omitempty"`
	ScrapeAnnotated *Mk8sMetricsScrapeAnnotated `json:"scrapeAnnotated,omitempty"`
}

type Mk8sMetricsScrapeAnnotated struct {
	IntervalSeconds   *int    `json:"intervalSeconds,omitempty"`
	IncludeNamespaces *string `json:"includeNamespaces,omitempty"`
	ExcludeNamespaces *string `json:"excludeNamespaces,omitempty"`
	RetainLabels      *string `json:"retainLabels,omitempty"`
}

type Mk8sLogsAddOnConfig struct {
	AuditEnabled      *bool   `json:"auditEnabled,omitempty"`
	IncludeNamespaces *string `json:"includeNamespaces,omitempty"`
	ExcludeNamespaces *string `json:"excludeNamespaces,omitempty"`
}

type Mk8sNvidiaAddOnConfig struct {
	TaintGPUNodes *bool `json:"taintGPUNodes,omitempty"`
}

type Mk8sAwsAddOnConfig struct {
	RoleArn *string `json:"roleArn,omitempty"`
}

type Mk8sAzureACRAddOnConfig struct {
	ClientId *string `json:"clientId,omitempty"`
}

/*** Status ***/

type Mk8sStatusAddOns struct {
	Dashboard           *Mk8sDashboardAddOnStatus           `json:"dashboard,omitempty"`
	AwsWorkloadIdentity *Mk8sAwsWorkloadIdentityAddOnStatus `json:"awsWorkloadIdentity,omitempty"`
	Metrics             *Mk8sMetricsAddOnStatus             `json:"metrics,omitempty"`
	Logs                *Mk8sLogsAddOnStatus                `json:"logs,omitempty"`
	AwsECR              *Mk8sAwsAddOnStatus                 `json:"awsECR,omitempty"`
	AwsEFS              *Mk8sAwsAddOnStatus                 `json:"awsEFS,omitempty"`
	AwsELB              *Mk8sAwsAddOnStatus                 `json:"awsELB,omitempty"`
}

// Add Ons //

type Mk8sDashboardAddOnStatus struct {
	Url *string `json:"url,omitempty"`
}

type Mk8sAwsWorkloadIdentityAddOnStatus struct {
	OidcProviderConfig *Mk8sOidcProviderConfig `json:"oidcProviderConfig,omitempty"`
	TrustPolicy        *map[string]interface{} `json:"trustPolicy,omitempty"`
}

type Mk8sOidcProviderConfig struct {
	ProviderUrl *string `json:"providerUrl,omitempty"`
	Audience    *string `json:"audience,omitempty"`
}

type Mk8sMetricsAddOnStatus struct {
	PrometheusEndpoint *string                 `json:"prometheusEndpoint,omitempty"`
	RemoteWriteConfig  *map[string]interface{} `json:"remoteWriteConfig,omitempty"`
}

type Mk8sLogsAddOnStatus struct {
	LokiAddress *string `json:"lokiAddress,omitempty"`
}

type Mk8sAwsAddOnStatus struct {
	TrustPolicy *map[string]interface{} `json:"trustPolicy,omitempty"`
}

type Mk8sNonCustomizableAddonConfig struct{}

/*** Client Functions ***/

func (c *Client) CreateMk8s(mk8s Mk8s) (*Mk8s, int, error) {

	code, err := c.CreateResource("mk8s", *mk8s.Name, mk8s)

	if err != nil {
		return nil, code, err
	}

	return c.GetMk8s(*mk8s.Name)
}

func (c *Client) GetMk8s(name string) (*Mk8s, int, error) {

	mk8s, code, err := c.GetResource(fmt.Sprintf("mk8s/%s", name), new(Mk8s))

	if err != nil {
		return nil, code, err
	}

	return mk8s.(*Mk8s), code, err
}

func (c *Client) UpdateMk8s(mk8s Mk8s) (*Mk8s, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("mk8s/%s", *mk8s.Name), mk8s)

	if err != nil {
		return nil, code, err
	}

	return c.GetMk8s(*mk8s.Name)
}

func (c *Client) DeleteMk8s(name string) error {
	return c.DeleteResource(fmt.Sprintf("mk8s/%s", name))
}
