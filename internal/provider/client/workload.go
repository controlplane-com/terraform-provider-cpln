package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Workloads - GVC Workloads
type Workloads struct {
	Kind     string     `json:"kind,omitempty"`
	ItemKind string     `json:"itemKind,omitempty"`
	Links    []Link     `json:"links,omitempty"`
	Items    []Workload `json:"items,omitempty"`
}

// Workload - GVC Workload
type Workload struct {
	Base
	Spec        *WorkloadSpec   `json:"spec,omitempty"`
	SpecReplace *WorkloadSpec   `json:"$replace/spec,omitempty"`
	Status      *WorkloadStatus `json:"status,omitempty"`
}

// WorkloadSpec - Workload Specifications
type WorkloadSpec struct {
	Type               *string                     `json:"type,omitempty"`
	IdentityLink       *string                     `json:"identityLink,omitempty"`
	Containers         *[]WorkloadContainer        `json:"containers,omitempty"`
	FirewallConfig     *WorkloadFirewall           `json:"firewallConfig,omitempty"`
	DefaultOptions     *WorkloadOptions            `json:"defaultOptions,omitempty"`
	LocalOptions       *[]WorkloadOptions          `json:"localOptions,omitempty"`
	Job                *WorkloadJob                `json:"job,omitempty"`
	Sidecar            *WorkloadSidecar            `json:"sidecar,omitempty"`
	SupportDynamicTags *bool                       `json:"supportDynamicTags,omitempty"`
	RolloutOptions     *WorkloadRolloutOptions     `json:"rolloutOptions,omitempty"`
	SecurityOptions    *WorkloadSecurityOptions    `json:"securityOptions,omitempty"`
	LoadBalancer       *WorkloadLoadBalancer       `json:"loadBalancer,omitempty"`
	Extras             *any                        `json:"extras,omitempty"`
	RequestRetryPolicy *WorkloadRequestRetryPolicy `json:"requestRetryPolicy,omitempty"`
}

// WorkloadContainer - Workload Container Definition
type WorkloadContainer struct {
	Name             *string                       `json:"name,omitempty"`
	Image            *string                       `json:"image,omitempty"`
	WorkingDirectory *string                       `json:"workingDir,omitempty"`
	Metrics          *WorkloadContainerMetrics     `json:"metrics,omitempty"`
	Port             *int                          `json:"port,omitempty"`
	Ports            *[]WorkloadContainerPort      `json:"ports,omitempty"`
	Memory           *string                       `json:"memory,omitempty"`
	ReadinessProbe   *WorkloadHealthCheck          `json:"readinessProbe,omitempty"`
	LivenessProbe    *WorkloadHealthCheck          `json:"livenessProbe,omitempty"`
	CPU              *string                       `json:"cpu,omitempty"`
	MinCPU           *string                       `json:"minCpu,omitempty"`
	MinMemory        *string                       `json:"minMemory,omitempty"`
	Env              *[]WorkloadContainerNameValue `json:"env,omitempty"`
	GPU              *WorkloadContainerGpu         `json:"gpu,omitempty"`
	InheritEnv       *bool                         `json:"inheritEnv,omitempty"`
	Command          *string                       `json:"command,omitempty"`
	Args             *[]string                     `json:"args,omitempty"`
	LifeCycle        *WorkloadLifeCycle            `json:"lifecycle,omitempty"`
	Volumes          *[]WorkloadContainerVolume    `json:"volumes,omitempty"`
}

type WorkloadContainerMetrics struct {
	Path *string `json:"path,omitempty"`
	Port *int    `json:"port,omitempty"`
}

type WorkloadContainerPort struct {
	Protocol *string `json:"protocol,omitempty"`
	Number   *int    `json:"number,omitempty"`
}

// WorkloadHealthCheck - Health Check (used my readiness and liveness probes)
type WorkloadHealthCheck struct {
	Exec                *WorkloadExec                 `json:"exec,omitempty"`
	GRPC                *WorkloadHealthCheckGrpc      `json:"grpc,omitempty"`
	TCPSocket           *WorkloadHealthCheckTcpSocket `json:"tcpSocket,omitempty"`
	HTTPGet             *WorkloadHealthCheckHttpGet   `json:"httpGet,omitempty"`
	InitialDelaySeconds *int                          `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int                          `json:"periodSeconds,omitempty"`
	TimeoutSeconds      *int                          `json:"timeoutSeconds,omitempty"`
	SuccessThreshold    *int                          `json:"successThreshold,omitempty"`
	FailureThreshold    *int                          `json:"failureThreshold,omitempty"`
}

// WorkloadExec - WorkloadExec
type WorkloadExec struct {
	Command *[]string `json:"command,omitempty"`
}

type WorkloadHealthCheckGrpc struct {
	Port *int `json:"port,omitempty"`
}

// WorkloadHealthCheckTcpSocket - WorkloadHealthCheckTcpSocket
type WorkloadHealthCheckTcpSocket struct {
	Port *int `json:"port,omitempty"`
}

// WorkloadHealthCheckHttpGet - WorkloadHealthCheckHttpGet
type WorkloadHealthCheckHttpGet struct {
	Path        *string                       `json:"path,omitempty"`
	Port        *int                          `json:"port,omitempty"`
	HttpHeaders *[]WorkloadContainerNameValue `json:"httpHeaders,omitempty"`
	Scheme      *string                       `json:"scheme,omitempty"`
}

// WorkloadContainerNameValue - Name/Value Struct
type WorkloadContainerNameValue struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

type WorkloadContainerGpu struct {
	Nvidia *WorkloadContainerGpuNvidia `json:"nvidia,omitempty"`
	Custom *WorkloadContainerGpuCustom `json:"custom,omitempty"`
}

type WorkloadContainerGpuNvidia struct {
	Model    *string `json:"model,omitempty"`
	Quantity *int    `json:"quantity,omitempty"`
}

type WorkloadContainerGpuCustom struct {
	Resource     *string `json:"resource,omitempty"`
	RuntimeClass *string `json:"runtimeClass,omitempty"`
	Quantity     *int    `json:"quantity,omitempty"`
}

// LifeCycle
type WorkloadLifeCycle struct {
	PostStart *WorkloadLifeCycleSpec `json:"postStart,omitempty"`
	PreStop   *WorkloadLifeCycleSpec `json:"preStop,omitempty"`
}

// LifeCycle - Inner
type WorkloadLifeCycleSpec struct {
	Exec *WorkloadExec `json:"exec,omitempty"`
}

// WorkloadContainerVolume - Volume Spec
type WorkloadContainerVolume struct {
	Uri            *string `json:"uri,omitempty"`
	RecoveryPolicy *string `json:"recoveryPolicy,omitempty"`
	Path           *string `json:"path,omitempty"`
}

// WorkloadFirewall - Firewall Config
type WorkloadFirewall struct {
	External *WorkloadFirewallExternal `json:"external,omitempty"`
	Internal *WorkloadFirewallInternal `json:"internal,omitempty"`
}

// WorkloadFirewallExternal
type WorkloadFirewallExternal struct {
	InboundAllowCidr      *[]string                            `json:"inboundAllowCIDR,omitempty"`
	InboundBlockedCidr    *[]string                            `json:"inboundBlockedCIDR,omitempty"`
	OutboundAllowHostname *[]string                            `json:"outboundAllowHostname,omitempty"`
	OutboundAllowPort     *[]WorkloadFirewallOutboundAllowPort `json:"outboundAllowPort,omitempty"`
	OutboundAllowCidr     *[]string                            `json:"outboundAllowCIDR,omitempty"`
	OutboundBlockedCidr   *[]string                            `json:"outboundBlockedCIDR,omitempty"`
	Http                  *WorkloadFirewallExternalHttp        `json:"http,omitempty"`
}

type WorkloadFirewallOutboundAllowPort struct {
	Protocol *string `json:"protocol,omitempty"`
	Number   *int    `json:"number,omitempty"`
}

type WorkloadFirewallExternalHttp struct {
	InboundHeaderFilter *[]WorkloadFirewallExternalHttpHeaderFilter `json:"inboundHeaderFilter,omitempty"`
}

type WorkloadFirewallExternalHttpHeaderFilter struct {
	Key           *string   `json:"key,omitempty"`
	AllowedValues *[]string `json:"allowedValues,omitempty"`
	BlockedValues *[]string `json:"blockedValues,omitempty"`
}

// WorkloadFirewallInternal - Firewall Internal
type WorkloadFirewallInternal struct {
	InboundAllowType     *string   `json:"inboundAllowType,omitempty"`
	InboundAllowWorkload *[]string `json:"inboundAllowWorkload,omitempty"`
}

// WorkloadOptions - WorkloadOptions
type WorkloadOptions struct {
	AutoScaling             *WorkloadOptionsAutoscaling `json:"autoscaling,omitempty"`
	TimeoutSeconds          *int                        `json:"timeoutSeconds,omitempty"`
	CapacityAI              *bool                       `json:"capacityAI,omitempty"`
	CapacityAIUpdateMinutes *int                        `json:"capacityAIUpdateMinutes,omitempty"`
	Debug                   *bool                       `json:"debug,omitempty"`
	Suspend                 *bool                       `json:"suspend,omitempty"`
	MultiZone               *WorkloadOptionsMultiZone   `json:"multiZone,omitempty"`
	Location                *string                     `json:"location,omitempty"`
}

// WorkloadOptionsAutoscaling - Auto Scaling Options
type WorkloadOptionsAutoscaling struct {
	Metric           *string                            `json:"metric,omitempty"`
	Multi            *[]WorkloadOptionsAutoscalingMulti `json:"multi,omitempty"`
	MetricPercentile *string                            `json:"metricPercentile,omitempty"`
	Target           *int                               `json:"target,omitempty"`
	MinScale         *int                               `json:"minScale,omitempty"`
	MaxScale         *int                               `json:"maxScale,omitempty"`
	ScaleToZeroDelay *int                               `json:"scaleToZeroDelay,omitempty"`
	MaxConcurrency   *int                               `json:"maxConcurrency,omitempty"`
	Keda             *WorkloadOptionsAutoscalingKeda    `json:"keda,omitempty"`
}

// WorkloadOptionsAutoscalingMulti - Multi Metrics
type WorkloadOptionsAutoscalingMulti struct {
	Metric *string `json:"metric,omitempty"`
	Target *int    `json:"target,omitempty"`
}

type WorkloadOptionsAutoscalingKeda struct {
	PollingInterval       *int                                     `json:"pollingInterval,omitempty"`
	CooldownPeriod        *int                                     `json:"cooldownPeriod,omitempty"`
	InitialCooldownPeriod *int                                     `json:"initialCooldownPeriod,omitempty"`
	Triggers              *[]WorkloadOptionsAutoscalingKedaTrigger `json:"triggers,omitempty"`
	Advanced              *WorkloadOptionsAutoscalingKedaAdvanced  `json:"advanced,omitempty"`
}

type WorkloadOptionsAutoscalingKedaTrigger struct {
	Type              *string                                                 `json:"type,omitempty"`
	Metadata          *map[string]interface{}                                 `json:"metadata,omitempty"`
	Name              *string                                                 `json:"name,omitempty"`
	UseCachedMetrics  *bool                                                   `json:"useCachedMetrics,omitempty"`
	MetricType        *string                                                 `json:"metricType,omitempty"`
	AuthenticationRef *WorkloadOptionsAutoscalingKedaTriggerAuthenticationRef `json:"authenticationRef,omitempty"`
}

type WorkloadOptionsAutoscalingKedaTriggerAuthenticationRef struct {
	Name *string `json:"name,omitempty"`
}

type WorkloadOptionsAutoscalingKedaAdvanced struct {
	ScalingModifiers *WorkloadOptionsAutoscalingKedaAdvancedScalingModifiers `json:"scalingModifiers,omitempty"`
}

type WorkloadOptionsAutoscalingKedaAdvancedScalingModifiers struct {
	Target           *string `json:"target,omitempty"`
	ActivationTarget *string `json:"activationTarget,omitempty"`
	MetricType       *string `json:"metricType,omitempty"`
	Formula          *string `json:"formula,omitempty"`
}

type WorkloadOptionsMultiZone struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// WorkloadJob - Cronjob
type WorkloadJob struct {
	Schedule              *string `json:"schedule,omitempty"`
	ConcurrencyPolicy     *string `json:"concurrencyPolicy,omitempty"` // Enum: [ Forbid, Replace ]
	HistoryLimit          *int    `json:"historyLimit,omitempty"`
	RestartPolicy         *string `json:"restartPolicy,omitempty"` // Enum: [ OnFailure, Never ]
	ActiveDeadlineSeconds *int    `json:"activeDeadlineSeconds,omitempty"`
}

// WorkloadSidecar - Workload Sidecar
type WorkloadSidecar struct {
	Envoy *any `json:"envoy,omitempty"`
}

// Rollout Options
type WorkloadRolloutOptions struct {
	MinReadySeconds               *int    `json:"minReadySeconds,omitempty"`
	MaxUnavailableReplicas        *string `json:"maxUnavailableReplicas,omitempty"`
	MaxSurgeReplicas              *string `json:"maxSurgeReplicas,omitempty"`
	ScalingPolicy                 *string `json:"scalingPolicy,omitempty"`
	TerminationGracePeriodSeconds *int    `json:"terminationGracePeriodSeconds,omitempty"`
}

// Security Options
type WorkloadSecurityOptions struct {
	FileSystemGroupId *int `json:"filesystemGroupId,omitempty"`
}

// WorkloadLoadBalancer - Workload Load Balancer
type WorkloadLoadBalancer struct {
	Direct        *WorkloadLoadBalancerDirect      `json:"direct,omitempty"`
	GeoLocation   *WorkloadLoadBalancerGeoLocation `json:"geoLocation,omitempty"`
	ReplicaDirect *bool                            `json:"replicaDirect,omitempty"`
}

// WorkloadLoadBalancerDirect - Workload Load Balancer Direct
type WorkloadLoadBalancerDirect struct {
	Enabled *bool                             `json:"enabled,omitempty"`
	Ports   *[]WorkloadLoadBalancerDirectPort `json:"ports,omitempty"`
	IpSet   *string                           `json:"ipSet,omitempty"`
}

type WorkloadLoadBalancerDirectPort struct {
	ExternalPort  *int    `json:"externalPort,omitempty"`
	Protocol      *string `json:"protocol,omitempty"`
	Scheme        *string `json:"scheme,omitempty"`
	ContainerPort *int    `json:"containerPort,omitempty"`
}

type WorkloadLoadBalancerGeoLocation struct {
	Enabled *bool                                   `json:"enabled,omitempty"`
	Headers *WorkloadLoadBalancerGeoLocationHeaders `json:"headers,omitempty"`
}

type WorkloadLoadBalancerGeoLocationHeaders struct {
	Asn     *string `json:"asn,omitempty"`
	City    *string `json:"city,omitempty"`
	Country *string `json:"country,omitempty"`
	Region  *string `json:"region,omitempty"`
}

type WorkloadRequestRetryPolicy struct {
	Attempts *int      `json:"attempts,omitempty"`
	RetryOn  *[]string `json:"retryOn,omitempty"`
}

// WorkloadStatus - Workload Status
type WorkloadStatus struct {
	ParentId            *string                       `json:"parentId,omitempty"`
	CanonicalEndpoint   *string                       `json:"canonicalEndpoint,omitempty"`
	Endpoint            *string                       `json:"endpoint,omitempty"`
	InternalName        *string                       `json:"internalName,omitempty"`
	HealthCheck         *WorkloadStatusHealthCheck    `json:"healthCheck,omitempty"`
	CurrentReplicaCount *int                          `json:"currentReplicaCount,omitempty"`
	ResolvedImages      *WorkloadStatusResolvedImages `json:"resolvedImages,omitempty"`
	LoadBalancer        *[]WorkloadStatusLoadBalancer `json:"loadBalancer,omitempty"`
}

// WorkloadStatusHealthCheck - Health Check Status
type WorkloadStatusHealthCheck struct {
	Active      *bool   `json:"active,omitempty"`
	Success     *bool   `json:"success,omitempty"`
	Code        *int    `json:"code,omitempty"`
	Message     *string `json:"message,omitempty"`
	Failures    *int    `failures:"parentId,omitempty"`
	Successes   *int    `successes:"parentId,omitempty"`
	LastChecked *string `json:"lastChecked,omitempty"`
}

type WorkloadStatusResolvedImages struct {
	ResolvedForVersion *int                           `json:"resolvedForVersion,omitempty"`
	ResolvedAt         *string                        `json:"resolvedAt,omitempty"`
	ErrorMessages      *[]string                      `json:"errorMessages,omitempty"`
	Images             *[]WorkloadStatusResolvedImage `json:"images,omitempty"`
}

type WorkloadStatusResolvedImage struct {
	Digest    *string                                `json:"digest,omitempty"`
	Manifests *[]WorkloadStatusResolvedImageManifest `json:"manifests,omitempty"`
}

type WorkloadStatusResolvedImageManifest struct {
	Image     *string                 `json:"image,omitempty"`
	MediaType *string                 `json:"mediaType,omitempty"`
	Digest    *string                 `json:"digest,omitempty"`
	Platform  *map[string]interface{} `json:"platform,omitempty"`
}

type WorkloadStatusLoadBalancer struct {
	Origin *string `json:"origin,omitempty"`
	Url    *string `json:"url,omitempty"`
}

// GetWorkloads - Get Workloads by GVC name
func (c *Client) GetWorkloads(gvcName string) (*[]Workload, int, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s/gvc/%s/workload", c.HostURL, c.Org, gvcName), nil)
	if err != nil {
		return nil, 0, err
	}

	body, code, err := c.doRequest(req, "")
	if err != nil {
		return nil, code, err
	}

	workloads := Workloads{}
	err = json.Unmarshal(body, &workloads)
	if err != nil {
		return nil, 0, err
	}

	return &workloads.Items, code, nil
}

// GetWorkload - Get Workload by name
func (c *Client) GetWorkload(name, gvcName string) (*Workload, int, error) {

	workload, code, err := c.GetResource(fmt.Sprintf("gvc/%s/workload/%s", gvcName, name), new(Workload))
	if err != nil {
		return nil, code, err
	}

	// workload.(*Workload).RemoveEmptySlices()

	return workload.(*Workload), code, err
}

// CreateWorkload - Create a new Workload
func (c *Client) CreateWorkload(workload Workload, gvcName string) (*Workload, int, error) {

	// log.Printf("[INFO] About to create Workload with Name: %s", workload.Name)

	code, err := c.CreateResource(fmt.Sprintf("gvc/%s/workload", gvcName), *workload.Name, workload)
	if err != nil {
		return nil, code, err
	}

	// log.Printf("[INFO] Created Workload with Name: %s", workload.Name)

	return c.GetWorkload(*workload.Name, gvcName)
}

// UpdateWorkload - Update an existing workload
func (c *Client) UpdateWorkload(workload Workload, gvcName string) (*Workload, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("gvc/%s/workload/%s", gvcName, *workload.Name), workload)
	if err != nil {
		return nil, code, err
	}

	return c.GetWorkload(*workload.Name, gvcName)
}

// DeleteWorkload - Delete Workload by name
func (c *Client) DeleteWorkload(name, gvcName string) error {
	// log.Printf("[INFO] Deleting Workload with name: %s", name)
	return c.DeleteResource(fmt.Sprintf("gvc/%s/workload/%s", gvcName, name))
}
