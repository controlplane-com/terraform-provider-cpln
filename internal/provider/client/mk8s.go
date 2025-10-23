package cpln

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

const (
	CplnCommand    = "cpln"
	EnvCplnProfile = "CPLN_PROFILE"
	EnvCplnOrg     = "CPLN_ORG"

	// K8s Related
	K8sConfigVersion                  = "v1"
	K8sConfigKind                     = "Config"
	K8sExecAuthVersion                = "client.authentication.k8s.io/v1"
	K8sExecInteractiveModeNever       = "Never"
	K8sExecInteractiveModeIfAvailable = "IfAvailable"
	K8sExecInteractiveModeAlways      = "Always"
)

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
	Generic      *Mk8sGenericProvider      `json:"generic,omitempty"`
	Hetzner      *Mk8sHetznerProvider      `json:"hetzner,omitempty"`
	Aws          *Mk8sAwsProvider          `json:"aws,omitempty"`
	Linode       *Mk8sLinodeProvider       `json:"linode,omitempty"`
	Oblivus      *Mk8sOblivusProvider      `json:"oblivus,omitempty"`
	Lambdalabs   *Mk8sLambdalabsProvider   `json:"lambdalabs,omitempty"`
	Paperspace   *Mk8sPaperspaceProvider   `json:"paperspace,omitempty"`
	Ephemeral    *Mk8sEphemeralProvider    `json:"ephemeral,omitempty"`
	Triton       *Mk8sTritonProvider       `json:"triton,omitempty"`
	Azure        *Mk8sAzureProvider        `json:"azure,omitempty"`
	DigitalOcean *Mk8sDigitalOceanProvider `json:"digitalocean,omitempty"`
	Gcp          *Mk8sGcpProvider          `json:"gcp,omitempty"`
}

type Mk8sSpecAddOns struct {
	Dashboard             *Mk8sNonCustomizableAddonConfig `json:"dashboard,omitempty"`
	AzureWorkloadIdentity *Mk8sAzureWorkloadIdentityAddOn `json:"azureWorkloadIdentity,omitempty"`
	AwsWorkloadIdentity   *Mk8sNonCustomizableAddonConfig `json:"awsWorkloadIdentity,omitempty"`
	LocalPathStorage      *Mk8sNonCustomizableAddonConfig `json:"localPathStorage,omitempty"`
	Metrics               *Mk8sMetricsAddOn               `json:"metrics,omitempty"`
	Logs                  *Mk8sLogsAddOn                  `json:"logs,omitempty"`
	RegistryMirror        *Mk8sRegistryMirrorAddOn        `json:"registryMirror,omitempty"`
	Nvidia                *Mk8sNvidiaAddOn                `json:"nvidia,omitempty"`
	AwsEFS                *Mk8sAwsAddOn                   `json:"awsEFS,omitempty"`
	AwsECR                *Mk8sAwsAddOn                   `json:"awsECR,omitempty"`
	AwsELB                *Mk8sAwsAddOn                   `json:"awsELB,omitempty"`
	AzureACR              *Mk8sAzureACRAddOn              `json:"azureACR,omitempty"`
	Sysbox                *Mk8sNonCustomizableAddonConfig `json:"sysbox,omitempty"`
	Byok                  *Mk8sByokAddOn                  `json:"byok,omitempty"`
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
	Region               *string                  `json:"region,omitempty"`
	AwsTags              *map[string]interface{}  `json:"awsTags,omitempty"`
	SkipCreateRoles      *bool                    `json:"skipCreateRoles,omitempty"`
	Networking           *Mk8sNetworkingConfig    `json:"networking,omitempty"`
	PreInstallScript     *string                  `json:"preInstallScript,omitempty"`
	Image                *Mk8sAwsAmi              `json:"image,omitempty"`
	DeployRoleArn        *string                  `json:"deployRoleArn,omitempty"`
	DeployRoleChain      *[]Mk8sAwsAssumeRoleLink `json:"deployRoleChain,omitempty"`
	VpcId                *string                  `json:"vpcId,omitempty"`
	KeyPair              *string                  `json:"keyPair,omitempty"`
	DiskEncryptionKeyArn *string                  `json:"diskEncryptionKeyArn,omitempty"`
	SecurityGroupIds     *[]string                `json:"securityGroupIds,omitempty"`
	ExtraNodePolicies    *[]string                `json:"extraNodePolicies,omitempty"`
	NodePools            *[]Mk8sAwsPool           `json:"nodePools,omitempty"`
	Autoscaler           *Mk8sAutoscalerConfig    `json:"autoscaler,omitempty"`
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
	FileSystems        *[]string             `json:"fileSystems,omitempty"`
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

type Mk8sEphemeralProvider struct {
	Location  *string              `json:"location,omitempty"`
	NodePools *[]Mk8sEphemeralPool `json:"nodePools,omitempty"`
}

type Mk8sTritonProvider struct {
	Connection       *Mk8sTritonConnection   `json:"connection,omitempty"`
	Networking       *Mk8sNetworkingConfig   `json:"networking,omitempty"`
	PreInstallScript *string                 `json:"preInstallScript,omitempty"`
	Location         *string                 `json:"location,omitempty"`
	LoadBalancer     *Mk8sTritonLoadBalancer `json:"loadBalancer,omitempty"`
	PrivateNetworkId *string                 `json:"privateNetworkId,omitempty"`
	FirewallEnabled  *bool                   `json:"firewallEnabled,omitempty"`
	NodePools        *[]Mk8sTritonPool       `json:"nodePools,omitempty"`
	ImageId          *string                 `json:"imageId,omitempty"`
	SshKeys          *[]string               `json:"sshKeys,omitempty"`
	Autoscaler       *Mk8sAutoscalerConfig   `json:"autoscaler,omitempty"`
}

type Mk8sAzureProvider struct {
	Location         *string                 `json:"location,omitempty"`
	SubscriptionId   *string                 `tfsdjsonk:"subscriptionId,omitempty"`
	SdkSecretLink    *string                 `json:"sdkSecretLink,omitempty"`
	ResourceGroup    *string                 `json:"resourceGroup,omitempty"`
	Networking       *Mk8sNetworkingConfig   `json:"networking,omitempty"`
	PreInstallScript *string                 `json:"preInstallScript,omitempty"`
	Image            *Mk8sAzureImage         `json:"image,omitempty"`
	SshKeys          *[]string               `json:"sshKeys,omitempty"`
	NetworkId        *string                 `json:"networkId,omitempty"`
	Tags             *map[string]interface{} `json:"tags,omitempty"`
	NodePools        *[]Mk8sAzurePool        `json:"nodePools,omitempty"`
	Autoscaler       *Mk8sAutoscalerConfig   `json:"autoscaler,omitempty"`
}

type Mk8sDigitalOceanProvider struct {
	Region           *string                 `json:"region,omitempty"`
	DigitalOceanTags *[]string               `json:"digitalOceanTags,omitempty"`
	Networking       *Mk8sNetworkingConfig   `json:"networking,omitempty"`
	PreInstallScript *string                 `json:"preInstallScript,omitempty"`
	TokenSecretLink  *string                 `json:"tokenSecretLink,omitempty"`
	VpcId            *string                 `json:"vpcId,omitempty"`
	NodePools        *[]Mk8sDigitalOceanPool `json:"nodePools,omitempty"`
	Image            *string                 `json:"image,omitempty"`
	SshKeys          *[]string               `json:"sshKeys,omitempty"`
	ExtraSshKeys     *[]string               `json:"extraSshKeys,omitempty"`
	Autoscaler       *Mk8sAutoscalerConfig   `json:"autoscaler,omitempty"`
	ReservedIps      *[]string               `json:"reservedIps,omitempty"`
}

type Mk8sGcpProvider struct {
	ProjectId        *string                 `json:"projectId,omitempty"`
	Region           *string                 `json:"region,omitempty"`
	GcpLabels        *map[string]interface{} `json:"gcpLabels,omitempty"`
	Network          *string                 `json:"network,omitempty"`
	SaKeyLink        *string                 `json:"saKeyLink,omitempty"`
	Networking       *Mk8sNetworkingConfig   `json:"networking,omitempty"`
	PreInstallScript *string                 `json:"preInstallScript,omitempty"`
	Image            *Mk8sGcpImage           `json:"image,omitempty"`
	NodePools        *[]Mk8sGcpPool          `json:"nodePools,omitempty"`
	Autoscaler       *Mk8sAutoscalerConfig   `json:"autoscaler,omitempty"`
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

type Mk8sTritonPool struct {
	Mk8sGenericPool
	PackageId         *string                 `json:"packageId,omitempty"`
	OverrideImageId   *string                 `json:"overrideImageId,omitempty"`
	PublicNetworkId   *string                 `json:"publicNetworkId,omitempty"`
	PrivateNetworkIds *[]string               `json:"privateNetworkIds,omitempty"`
	TritonTags        *map[string]interface{} `json:"tritonTags,omitempty"`
	MinSize           *int                    `json:"minSize,omitempty"`
	MaxSize           *int                    `json:"maxSize,omitempty"`
}

type Mk8sAzurePool struct {
	Mk8sGenericPool
	Size          *string         `json:"size,omitempty"`
	SubnetId      *string         `json:"subnetId,omitempty"`
	Zones         *[]int          `json:"zone,omitempty"`
	OverrideImage *Mk8sAzureImage `json:"overrideImage,omitempty"`
	BootDiskSize  *int            `json:"bootDiskSize,omitempty"`
	MinSize       *int            `json:"minSize,omitempty"`
	MaxSize       *int            `json:"maxSize,omitempty"`
}

type Mk8sDigitalOceanPool struct {
	Mk8sGenericPool
	DropletSize   *string `json:"dropletSize,omitempty"`
	OverrideImage *string `json:"overrideImage,omitempty"`
	MinSize       *int    `json:"minSize,omitempty"`
	MaxSize       *int    `json:"maxSize,omitempty"`
}

type Mk8sGcpPool struct {
	Mk8sGenericPool
	MachineType   *string       `json:"machineType,omitempty"`
	Zone          *string       `json:"zone,omitempty"`
	OverrideImage *Mk8sGcpImage `json:"overrideImage,omitempty"`
	BootDiskSize  *int          `json:"bootDiskSize,omitempty"`
	MinSize       *int          `json:"minSize,omitempty"`
	MaxSize       *int          `json:"maxSize,omitempty"`
	Subnet        *string       `json:"subnet,omitempty"`
}

// Provider Common //

type Mk8sNetworkingConfig struct {
	ServiceNetwork *string `json:"serviceNetwork,omitempty"`
	PodNetwork     *string `json:"podNetwork,omitempty"`
	DnsForwarder   *string `json:"dnsForwarder,omitempty"`
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

type Mk8sAwsAssumeRoleLink struct {
	RoleArn           *string `json:"roleArn,omitempty"`
	ExternalId        *string `json:"externalId,omitempty"`
	SessionNamePrefix *string `json:"sessionNamePrefix,omitempty"`
}

// Triton Provider //

type Mk8sTritonConnection struct {
	Url                  *string `json:"url,omitempty"`
	Account              *string `json:"account,omitempty"`
	User                 *string `json:"user,omitempty"`
	PrivateKeySecretLink *string `json:"privateKeySecretLink,omitempty"`
}

type Mk8sTritonLoadBalancer struct {
	Manual  *Mk8sTritonManual           `json:"manual,omitempty"`
	None    *Mk8sTritonLoadBalancerNone `json:"none,omitempty"`
	Gateway *Mk8sTritonGateway          `json:"gateway,omitempty"`
}

type Mk8sTritonManual struct {
	PackageId         *string                  `json:"packageId,omitempty"`
	ImageId           *string                  `json:"imageId,omitempty"`
	PublicNetworkId   *string                  `json:"publicNetworkId,omitempty"`
	PrivateNetworkIds *[]string                `json:"privateNetworkIds,omitempty"`
	Metadata          *map[string]interface{}  `json:"metadata,omitempty"`
	Tags              *map[string]interface{}  `json:"tags,omitempty"`
	Logging           *Mk8sTritonManualLogging `json:"logging,omitempty"`
	Count             *int                     `json:"count,omitempty"`
	CnsInternalDomain *string                  `json:"cnsInternalDomain,omitempty"`
	CnsPublicDomain   *string                  `json:"cnsPublicDomain,omitempty"`
}

type Mk8sTritonManualLogging struct {
	NodePort       *int    `json:"nodePort,omitempty"`
	ExternalSyslog *string `json:"externalSyslog,omitempty"`
}

type Mk8sTritonGateway struct{}

type Mk8sTritonLoadBalancerNone struct{}

// Azure Provider //

type Mk8sAzureImage struct {
	Recommended *string                  `json:"recommended,omitempty"`
	Reference   *Mk8sAzureImageReference `json:"reference,omitempty"`
}

type Mk8sAzureImageReference struct {
	Publisher *string `json:"publisher,omitempty"`
	Offer     *string `json:"offer,omitempty"`
	Sku       *string `json:"sku,omitempty"`
	Version   *string `json:"version,omitempty"`
}

// Gcp //

type Mk8sGcpImage struct {
	Recommended *string `json:"recommended,omitempty"`
}

// Add Ons //

type Mk8sAzureWorkloadIdentityAddOn struct {
	TenantId *string `json:"tenantId,omitempty"`
}

type Mk8sMetricsAddOn struct {
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

type Mk8sLogsAddOn struct {
	AuditEnabled      *bool   `json:"auditEnabled,omitempty"`
	IncludeNamespaces *string `json:"includeNamespaces,omitempty"`
	ExcludeNamespaces *string `json:"excludeNamespaces,omitempty"`
	Docker            *bool   `json:"docker,omitempty"`
	Kubelet           *bool   `json:"kubelet,omitempty"`
	Kernel            *bool   `json:"kernel,omitempty"`
	Events            *bool   `json:"events,omitempty"`
}

type Mk8sRegistryMirrorAddOn struct {
	Mirrors *[]Mk8sAddOnRegistry `json:"mirrors,omitempty"`
}

type Mk8sAddOnRegistry struct {
	Registry *string   `json:"registry,omitempty"`
	Mirrors  *[]string `json:"mirrors,omitempty"`
}

type Mk8sNvidiaAddOn struct {
	TaintGPUNodes *bool `json:"taintGPUNodes,omitempty"`
}

type Mk8sAwsAddOn struct {
	RoleArn *string `json:"roleArn,omitempty"`
}

type Mk8sAzureACRAddOn struct {
	ClientId *string `json:"clientId,omitempty"`
}

type Mk8sByokAddOn struct {
	IgnoreUpdates *bool                `json:"ignoreUpdates,omitempty"`
	Location      *string              `json:"location,omitempty"`
	Config        *Mk8sByokAddOnConfig `json:"config,omitempty"`
}

type Mk8sByokAddOnConfig struct {
	Actuator      *Mk8sByokAddOnConfigActuator    `json:"actuator,omitempty"`
	Middlebox     *Mk8sByokAddOnConfigMiddlebox   `json:"middlebox,omitempty"`
	Common        *Mk8sByokAddOnConfigCommon      `json:"common,omitempty"`
	Longhorn      *Mk8sByokAddOnConfigLonghorn    `json:"longhorn,omitempty"`
	Ingress       *Mk8sByokAddOnConfigIngress     `json:"ingress,omitempty"`
	Istio         *Mk8sByokAddOnConfigIstio       `json:"istio,omitempty"`
	LogSplitter   *Mk8sByokAddOnConfigLogSplitter `json:"logSplitter,omitempty"`
	Monitoring    *Mk8sByokAddOnConfigMonitoring  `json:"monitoring,omitempty"`
	Redis         *Mk8sByokAddOnConfigRedisString `json:"redis,omitempty"`
	RedisHa       *Mk8sByokAddOnConfigRedisInt    `json:"redisHa,omitempty"`
	RedisSentinel *Mk8sByokAddOnConfigRedisInt    `json:"redisSentinel,omitempty"`
	TempoAgent    *Mk8sByokAddOnConfigTempoAgent  `json:"tempoAgent,omitempty"`
	InternalDns   *Mk8sByokAddOnConfigInternalDns `json:"internalDns,omitempty"`
}

type Mk8sByokAddOnConfigActuator struct {
	MinCpu    *string                 `json:"minCpu,omitempty"`
	MaxCpu    *string                 `json:"maxCpu,omitempty"`
	MinMemory *string                 `json:"minMemory,omitempty"`
	MaxMemory *string                 `json:"maxMemory,omitempty"`
	LogLevel  *string                 `json:"logLevel,omitempty"`
	Env       *map[string]interface{} `json:"env,omitempty"`
}

type Mk8sByokAddOnConfigMiddlebox struct {
	Enabled            *bool `json:"enabled,omitempty"`
	BandwidthAlertMbps *int  `json:"bandwidthAlertMbps,omitempty"`
}

type Mk8sByokAddOnConfigCommon struct {
	DeploymentReplicas *int                          `json:"deploymentReplicas,omitempty"`
	Pdb                *Mk8sByokAddOnConfigCommonPdb `json:"pdb,omitempty"`
}

type Mk8sByokAddOnConfigCommonPdb struct {
	MaxUnavailable *int `json:"maxUnavailable,omitempty"`
}

type Mk8sByokAddOnConfigLonghorn struct {
	Replicas *int `json:"replicas,omitempty"`
}

type Mk8sByokAddOnConfigIngress struct {
	Cpu           *string  `json:"cpu,omitempty"`
	Memory        *string  `json:"memory,omitempty"`
	TargetPercent *float32 `json:"targetPercent,omitempty"`
}

type Mk8sByokAddOnConfigIstio struct {
	Istiod         *Mk8sByokAddOnConfigIstioIstiod         `json:"istiod,omitempty"`
	IngressGateway *Mk8sByokAddOnConfigIstioIngressGateway `json:"ingressgateway,omitempty"`
	Sidecar        *Mk8sByokAddOnConfigIstioSidecar        `json:"sidecar,omitempty"`
}

type Mk8sByokAddOnConfigIstioIstiod struct {
	Replicas  *int    `json:"replicas,omitempty"`
	MinCpu    *string `json:"minCpu,omitempty"`
	MaxCpu    *string `json:"maxCpu,omitempty"`
	MinMemory *string `json:"minMemory,omitempty"`
	MaxMemory *string `json:"maxMemory,omitempty"`
	Pdb       *int    `json:"pdb,omitempty"`
}

type Mk8sByokAddOnConfigIstioIngressGateway struct {
	Replicas  *int    `json:"replicas,omitempty"`
	MaxCpu    *string `json:"maxCpu,omitempty"`
	MaxMemory *string `json:"maxMemory,omitempty"`
}

type Mk8sByokAddOnConfigIstioSidecar struct {
	MinCpu    *string `json:"minCpu,omitempty"`
	MinMemory *string `json:"minMemory,omitempty"`
}

type Mk8sByokAddOnConfigLogSplitter struct {
	MinCpu        *string `json:"minCpu,omitempty"`
	MaxCpu        *string `json:"maxCpu,omitempty"`
	MinMemory     *string `json:"minMemory,omitempty"`
	MaxMemory     *string `json:"maxMemory,omitempty"`
	MemBufferSize *string `json:"memBufferSize,omitempty"`
	PerPodRate    *int    `json:"perPodRate,omitempty"`
}

type Mk8sByokAddOnConfigMonitoring struct {
	MinMemory        *string                                        `json:"minMemory,omitempty"`
	MaxMemory        *string                                        `json:"maxMemory,omitempty"`
	KubeStateMetrics *Mk8sByokAddOnConfigMonitoringKubeStateMetrics `json:"kubeStateMetrics,omitempty"`
	Prometheus       *Mk8sByokAddOnConfigMonitoringPrometheus       `json:"prometheus,omitempty"`
}

type Mk8sByokAddOnConfigMonitoringKubeStateMetrics struct {
	MinMemory *string `json:"minMemory,omitempty"`
}

type Mk8sByokAddOnConfigMonitoringPrometheus struct {
	Main *Mk8sByokAddOnConfigMonitoringPrometheusMain `json:"main,omitempty"`
}

type Mk8sByokAddOnConfigMonitoringPrometheusMain struct {
	Storage *string `json:"storage"`
}

type Mk8sByokAddOnConfigRedisString struct {
	MinCpu    *string `json:"minCpu,omitempty"`
	MaxCpu    *string `json:"maxCpu,omitempty"`
	MinMemory *string `json:"minMemory,omitempty"`
	MaxMemory *string `json:"maxMemory,omitempty"`
	Storage   *string `json:"storage,omitempty"`
}

type Mk8sByokAddOnConfigRedisInt struct {
	MinCpu    *string `json:"minCpu,omitempty"`
	MaxCpu    *string `json:"maxCpu,omitempty"`
	MinMemory *string `json:"minMemory,omitempty"`
	MaxMemory *string `json:"maxMemory,omitempty"`
	Storage   *int    `json:"storage,omitempty"`
}

type Mk8sByokAddOnConfigTempoAgent struct {
	MinCpu    *string `json:"minCpu,omitempty"`
	MinMemory *string `json:"minMemory,omitempty"`
}

type Mk8sByokAddOnConfigInternalDns struct {
	MinCpu    *string `json:"minCpu,omitempty"`
	MaxCpu    *string `json:"maxCpu,omitempty"`
	MinMemory *string `json:"minMemory,omitempty"`
	MaxMemory *string `json:"maxMemory,omitempty"`
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

/*** MK8s Common ***/

type Mk8sCacertsResponse struct {
	Cacerts *string `json:"cacerts,omitempty"`
}

/*** K8s Config Related ***/

type K8sConfig struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Kind           string                 `yaml:"kind"`
	Preferences    map[string]interface{} `yaml:"preferences,omitempty"`
	Clusters       []K8sNamedCluster      `yaml:"clusters"`
	Contexts       []K8sNamedContext      `yaml:"contexts"`
	CurrentContext string                 `yaml:"current-context"`
	Users          []K8sNamedUser         `yaml:"users"`
}

type K8sNamedCluster struct {
	Name    string     `yaml:"name"`
	Cluster K8sCluster `yaml:"cluster"`
}

type K8sCluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"` // base64 encoded
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

type K8sNamedContext struct {
	Name    string     `yaml:"name"`
	Context K8sContext `yaml:"context"`
}

type K8sContext struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

type K8sNamedUser struct {
	Name string  `yaml:"name"`
	User K8sUser `yaml:"user"`
}

type K8sUser struct {
	ClientCertificateData string        `yaml:"client-certificate-data,omitempty"` // base64 encoded
	ClientKeyData         string        `yaml:"client-key-data,omitempty"`         // base64 encoded
	Token                 string        `yaml:"token,omitempty"`
	Username              string        `yaml:"username,omitempty"`
	Password              string        `yaml:"password,omitempty"`
	Exec                  K8sExecConfig `yaml:"exec,omitempty"`
}

type K8sExecConfig struct {
	APIVersion         string          `yaml:"apiVersion,omitempty"`
	Command            string          `yaml:"command"`
	Args               []string        `yaml:"args,omitempty"`
	Env                []K8sExecEnvVar `yaml:"env,omitempty"`
	ProvideClusterInfo bool            `yaml:"provideClusterInfo,omitempty"`
	InteractiveMode    string          `yaml:"interactiveMode,omitempty"`
}

type K8sExecEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

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

// CreateMk8sKubeconfig retrieves MK8s cluster info and cacerts, builds a kubeconfig in YAML format, and returns a pointer to the YAML string along with an error.
func (c *Client) CreateMk8sKubeconfig(mk8sName string, profileName *string, serviceAccountName *string) (*string, error) {
	// Construct the /-cacerts link out of the MK8s link
	cacertsLink := fmt.Sprintf("org/%s/mk8s/%s/-cacerts", c.Org, mk8sName)

	// Get the cluster
	mk8s, _, err := c.GetResource(fmt.Sprintf("mk8s/%s", mk8sName), new(Mk8s))
	if err != nil {
		return nil, err
	}

	// Get the cacerts
	cacerts, _, err := c.Get(cacertsLink, new(Mk8sCacertsResponse))
	if err != nil {
		return nil, err
	}

	// Cacerts cannot be nil
	if cacerts.(*Mk8sCacertsResponse).Cacerts == nil {
		return nil, fmt.Errorf("the MK8s cluster '%s' has empty cacerts, please try again later", mk8sName)
	}

	// Build the kubeconfig
	kubeconfig, err := c.buildKubeconfig(mk8s.(*Mk8s), cacerts.(*Mk8sCacertsResponse), profileName, serviceAccountName)
	if err != nil {
		return nil, err
	}

	// Convert kubeconfig to YAML format
	kubeconfigYamlBytes, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error marshalling to YAML: %v", err)
	}

	// Convert YAML bytes to YAML string
	kubeconfigYaml := string(kubeconfigYamlBytes)

	// Return the kubeconfig in YAML format
	return &kubeconfigYaml, err
}

// buildKubeconfig constructs a K8sConfig object for the given MK8s cluster using cacerts and user data, and returns it or an error.
func (c *Client) buildKubeconfig(mk8s *Mk8s, cacerts *Mk8sCacertsResponse, profileName *string, serviceAccountName *string) (*K8sConfig, error) {
	// Extract the server url
	serverUrl := mk8s.Status.ServerUrl

	// Handle server url does not exist
	if serverUrl == nil {
		return nil, fmt.Errorf("the specified MK8s cluster '%s' has no serverUrl", *mk8s.Name)
	}

	// Construct the cluster name
	clusterName := fmt.Sprintf("%s/%s/%s", c.Org, *mk8s.Name, *mk8s.Alias)

	// Build the Kubernetes user
	user, err := c.buildK8sUser(*mk8s.Name, clusterName, profileName, serviceAccountName)
	if err != nil {
		return nil, err
	}

	// Build and return the kubeconfig
	return &K8sConfig{
		APIVersion:     K8sConfigVersion,
		Kind:           K8sConfigKind,
		CurrentContext: clusterName,
		Users: []K8sNamedUser{
			*user,
		},
		Clusters: []K8sNamedCluster{
			{
				Name: clusterName,
				Cluster: K8sCluster{
					CertificateAuthorityData: base64.StdEncoding.EncodeToString([]byte(*cacerts.Cacerts)),
					Server:                   *serverUrl,
				},
			},
		},
		Contexts: []K8sNamedContext{
			{
				Name: clusterName,
				Context: K8sContext{
					Cluster: clusterName,
					User:    user.Name,
				},
			},
		},
	}, nil
}

// buildK8sUser creates a K8sNamedUser based on either a profile token or service account key (ensuring only one is provided) and returns it or an error.
func (c *Client) buildK8sUser(mk8sName string, clusterName string, profileName *string, serviceAccountName *string) (*K8sNamedUser, error) {
	// Profile and service account cannot be defined together
	if profileName != nil && len(*profileName) != 0 && serviceAccountName != nil && len(*serviceAccountName) != 0 {
		return nil, fmt.Errorf("exactly one of cpln profile or an existing service account can be specified in order to create the kubeconfig")
	}

	// Create a user using profile name
	if profileName != nil && len(*profileName) != 0 {
		// Extract token from the specified profile
		token, err := c.ExtractTokenFromProfile(*profileName)
		if err != nil {
			return nil, err
		}

		// Remove Bearer from the start
		*token = strings.TrimPrefix(*token, "Bearer ")

		// Build K8s username
		username, isServiceAccountToken := buildK8sProfileUsername(c.Org, clusterName, *token)

		// Handle service account token from profile differently
		if isServiceAccountToken {
			return &K8sNamedUser{
				Name: username,
				User: K8sUser{
					Token: *token,
				},
			}, nil
		}

		// Construct and return the user
		return &K8sNamedUser{
			Name: username,
			User: K8sUser{
				Exec: K8sExecConfig{
					APIVersion: K8sExecAuthVersion,
					Command:    CplnCommand,
					Args:       []string{"mk8s", "auth", mk8sName},
					Env: []K8sExecEnvVar{
						{
							Name:  EnvCplnProfile,
							Value: *profileName,
						},
						{
							Name:  EnvCplnOrg,
							Value: c.Org,
						},
					},
					ProvideClusterInfo: false,
					InteractiveMode:    K8sExecInteractiveModeIfAvailable,
				},
			},
		}, nil
	}

	// Create a user using a service account, this will add a new key to the specified service account
	if serviceAccountName != nil && len(*serviceAccountName) != 0 {
		// Create a new key for the kubeconfig
		key, err := c.AddServiceAccountKey(*serviceAccountName, fmt.Sprintf("A Kubeconfig key for cluster '%s'", mk8sName))
		if err != nil {
			return nil, err
		}

		// Declare the username
		username := buildK8sServiceAccountUsername(clusterName, *serviceAccountName)

		// Construct and return the user
		return &K8sNamedUser{
			Name: username,
			User: K8sUser{
				Token: key.Key,
			},
		}, nil
	}

	// If none of the above, then the user must provide either of them
	return nil, fmt.Errorf("at lease one of a cpln profile or an existing service account must be specified in order to create the kubeconfig")
}

// buildK8sProfileUsername parses the token to extract an email for username (or treats it as a service account token)
// and returns the username along with a bool indicating if it is a service account token.
func buildK8sProfileUsername(org string, clusterName string, token string) (string, bool) {
	// Assuming that this is a refresh token, let's attempt to parse the JWT
	// token into a struct and extract the email of the user
	jwtToken, _, err := new(jwt.Parser).ParseUnverified(token, &CplnClaims{})
	if err == nil && jwtToken != nil && jwtToken.Claims != nil {
		// Refer to the claims as CplnClaims
		cplnClaims := jwtToken.Claims.(*CplnClaims)

		// If the parse was successful, let's use the email in the username
		if len(cplnClaims.Email) != 0 {
			return fmt.Sprintf("%s/%s", org, cplnClaims.Email), false
		}
	}

	// Otherwise, let's treat it as a service account token
	serviceAccountToken := ParseServiceAccountToken(token)

	// If the service account was parsed successfully, then we will be able to create a username for it
	if serviceAccountToken != nil {
		return buildK8sServiceAccountUsername(clusterName, serviceAccountToken.Name), true
	}

	// Default to a generic username if none of the above worked
	return fmt.Sprintf("mk8s-user-%d", time.Now().UnixMilli()), false
}

// buildK8sServiceAccountUsername formats and returns a K8s username string for a service account.
func buildK8sServiceAccountUsername(clusterName string, serviceAccountName string) string {
	return fmt.Sprintf("%s/sa:%s", clusterName, serviceAccountName)
}
