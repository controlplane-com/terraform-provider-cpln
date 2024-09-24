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
	Type               *string          `json:"type,omitempty"`
	IdentityLink       *string          `json:"identityLink,omitempty"`
	Containers         *[]ContainerSpec `json:"containers,omitempty"`
	FirewallConfig     *FirewallSpec    `json:"firewallConfig,omitempty"`
	DefaultOptions     *Options         `json:"defaultOptions,omitempty"`
	LocalOptions       *[]Options       `json:"localOptions,omitempty"`
	RolloutOptions     *RolloutOptions  `json:"rolloutOptions,omitempty"`
	Job                *JobSpec         `json:"job,omitempty"`
	SecurityOptions    *SecurityOptions `json:"securityOptions,omitempty"`
	SupportDynamicTags *bool            `json:"supportDynamicTags,omitempty"`
	Sidecar            *WorkloadSidecar `json:"sidecar,omitempty"`
}

// ContainerSpec - Workload Container Definition
type ContainerSpec struct {
	Name             *string          `json:"name,omitempty"`
	Image            *string          `json:"image,omitempty"`
	Port             *int             `json:"port,omitempty"`
	Ports            *[]PortSpec      `json:"ports,omitempty"`
	Memory           *string          `json:"memory,omitempty"`
	ReadinessProbe   *HealthCheckSpec `json:"readinessProbe,omitempty"`
	LivenessProbe    *HealthCheckSpec `json:"livenessProbe,omitempty"`
	CPU              *string          `json:"cpu,omitempty"`
	GPU              *GpuResource     `json:"gpu,omitempty"`
	MinCPU           *string          `json:"minCpu,omitempty"`
	MinMemory        *string          `json:"minMemory,omitempty"`
	Env              *[]NameValue     `json:"env,omitempty"`
	Args             *[]string        `json:"args,omitempty"`
	Volumes          *[]VolumeSpec    `json:"volumes,omitempty"`
	Metrics          *Metrics         `json:"metrics,omitempty"`
	Command          *string          `json:"command,omitempty"`
	InheritEnv       *bool            `json:"inheritEnv,omitempty"`
	WorkingDirectory *string          `json:"workingDir,omitempty"`
	LifeCycle        *LifeCycleSpec   `json:"lifecycle,omitempty"`
}

// GPU - GPU Settings
type GpuResource struct {
	Nvidia *Nvidia `json:"nvidia,omitempty"`
}

type Nvidia struct {
	Model    *string `json:"model,omitempty"`
	Quantity *int    `json:"quantity,omitempty"`
}

// NameValue - Name/Value Struct
type NameValue struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

// PortSpec - Ports
type PortSpec struct {
	Protocol *string `json:"protocol,omitempty"`
	Number   *int    `json:"number,omitempty"`
}

// Options - Options
type Options struct {
	AutoScaling    *AutoScaling `json:"autoscaling,omitempty"`
	CapacityAI     *bool        `json:"capacityAI,omitempty"`
	TimeoutSeconds *int         `json:"timeoutSeconds,omitempty"`
	Debug          *bool        `json:"debug,omitempty"`
	Suspend        *bool        `json:"suspend,omitempty"`
	Location       *string      `json:"location,omitempty"`
}

// AutoScaling - Auto Scaling Options
type AutoScaling struct {
	Metric           *string `json:"metric,omitempty"`
	MetricPercentile *string `json:"metricPercentile,omitempty"`
	Target           *int    `json:"target,omitempty"`
	MaxScale         *int    `json:"maxScale,omitempty"`
	MinScale         *int    `json:"minScale,omitempty"`
	MaxConcurrency   *int    `json:"maxConcurrency,omitempty"`
	ScaleToZeroDelay *int    `json:"scaleToZeroDelay,omitempty"`
}

// FirewallSpec - Firewall Config
type FirewallSpec struct {
	External *FirewallSpecExternal `json:"external,omitempty"`
	Internal *FirewallSpecInternal `json:"internal,omitempty"`
}

// FirewallSpecExternal - Firewall Spec External
type FirewallSpecExternal struct {
	InboundAllowCIDR      *[]string                    `json:"inboundAllowCIDR,omitempty"`
	OutboundAllowCIDR     *[]string                    `json:"outboundAllowCIDR,omitempty"`
	OutboundAllowHostname *[]string                    `json:"outboundAllowHostname,omitempty"`
	OutboundAllowPort     *[]FirewallOutboundAllowPort `json:"outboundAllowPort,omitempty"`
}

type FirewallOutboundAllowPort struct {
	Protocol *string `json:"protocol,omitempty"`
	Number   *int    `json:"number,omitempty"`
}

// FirewallSpecInternal - Firewall Spec Internal
type FirewallSpecInternal struct {
	InboundAllowType     *string   `json:"inboundAllowType,omitempty"`
	InboundAllowWorkload *[]string `json:"inboundAllowWorkload,omitempty"`
}

// WorkloadStatus - Workload Status
type WorkloadStatus struct {
	ParentID            *string            `json:"parentId,omitempty"`
	CanonicalEndpoint   *string            `json:"canonicalEndpoint,omitempty"`
	Endpoint            *string            `json:"endpoint,omitempty"`
	InternalName        *string            `json:"internalName,omitempty"`
	CurrentReplicaCount *int               `json:"currentReplicaCount,omitempty"`
	HealthCheck         *HealthCheckStatus `json:"healthCheck,omitempty"`
	ResolvedImages      *ResolvedImages    `json:"resolvedImages,omitempty"`
}

// HealthCheckStatus - Health Check Status
type HealthCheckStatus struct {
	Active      *bool   `json:"active,omitempty"`
	Success     *bool   `json:"success,omitempty"`
	Code        *int    `json:"code,omitempty"`
	Message     *string `json:"message,omitempty"`
	Failures    *int    `failures:"parentId,omitempty"`
	Successes   *int    `successes:"parentId,omitempty"`
	LastChecked *string `json:"lastChecked,omitempty"`
}

type ResolvedImages struct {
	ResolvedForVersion *int             `json:"resolvedForVersion,omitempty"`
	ResolvedAt         *string          `json:"resolvedAt,omitempty"`
	Images             *[]ResolvedImage `json:"images,omitempty"`
}

type ResolvedImage struct {
	Digest    *string                  `json:"digest,omitempty"`
	Manifests *[]ResolvedImageManifest `json:"manifests,omitempty"`
}

type ResolvedImageManifest struct {
	Image     *string             `json:"image,omitempty"`
	MediaType *string             `json:"mediaType,omitempty"`
	Digest    *string             `json:"digest,omitempty"`
	Platform  *map[string]*string `json:"platform,omitempty"`
}

// HealthCheckSpec - Health Check Spec (used my readiness and liveness probes)
type HealthCheckSpec struct {
	Exec                *Exec      `json:"exec,omitempty"`
	GRPC                *GRPC      `json:"grpc,omitempty"`
	TCPSocket           *TCPSocket `json:"tcpSocket,omitempty"`
	HTTPGet             *HTTPGet   `json:"httpGet,omitempty"`
	InitialDelaySeconds *int       `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int       `json:"periodSeconds,omitempty"`
	TimeoutSeconds      *int       `json:"timeoutSeconds,omitempty"`
	SuccessThreshold    *int       `json:"successThreshold,omitempty"`
	FailureThreshold    *int       `json:"failureThreshold,omitempty"`
}

// VolumeSpec - Volume Spec
type VolumeSpec struct {
	Uri            *string `json:"uri,omitempty"`
	RecoveryPolicy *string `json:"recoveryPolicy,omitempty"`
	Path           *string `json:"path,omitempty"`
}

// Metrics - Metrics
type Metrics struct {
	Path *string `json:"path,omitempty"`
	Port *int    `json:"port,omitempty"`
}

// Exec - Exec
type Exec struct {
	Command *[]string `json:"command,omitempty"`
}

type GRPC struct {
	Port *int `json:"port,omitempty"`
}

// TCPSocket - TCPSocket
type TCPSocket struct {
	Port *int `json:"port,omitempty"`
}

// HTTPGet - HTTPGet
type HTTPGet struct {
	Path        *string      `json:"path,omitempty"`
	Port        *int         `json:"port,omitempty"`
	HTTPHeaders *[]NameValue `json:"httpHeaders,omitempty"`
	Scheme      *string      `json:"scheme,omitempty"`
}

// LifeCycle
type LifeCycleSpec struct {
	PostStart *LifeCycleInner `json:"postStart,omitempty"`
	PreStop   *LifeCycleInner `json:"preStop,omitempty"`
}

// LifeCycle - Inner
type LifeCycleInner struct {
	Exec *Exec `json:"exec,omitempty"`
}

// JobSpec - Cronjob
type JobSpec struct {
	Schedule              *string `json:"schedule,omitempty"`
	ConcurrencyPolicy     *string `json:"concurrencyPolicy,omitempty"` // Enum: [ Forbid, Replace ]
	HistoryLimit          *int    `json:"historyLimit,omitempty"`
	RestartPolicy         *string `json:"restartPolicy,omitempty"` // Enum: [ OnFailure, Never ]
	ActiveDeadlineSeconds *int    `json:"activeDeadlineSeconds,omitempty"`
}

// Rollout Options
type RolloutOptions struct {
	MinReadySeconds        *int    `json:"minReadySeconds,omitempty"`
	MaxUnavailableReplicas *string `json:"maxUnavailableReplicas,omitempty"`
	MaxSurgeReplicas       *string `json:"maxSurgeReplicas,omitempty"`
	ScalingPolicy          *string `json:"scalingPolicy,omitempty"`
}

// Security Options
type SecurityOptions struct {
	FileSystemGroupID *int         `json:"filesystemGroupId,omitempty"`
	GeoLocation       *GeoLocation `json:"geoLocation,omitempty"`
}

type GeoLocation struct {
	Enabled *bool               `json:"enabled,omitempty"`
	Headers *GeoLocationHeaders `json:"headers,omitempty"`
}

type GeoLocationHeaders struct {
	Asn     *string `json:"asn,omitempty"`
	City    *string `json:"city,omitempty"`
	Country *string `json:"country,omitempty"`
	Region  *string `json:"region,omitempty"`
}

// WorkloadSidecar - Workload Sidecar
type WorkloadSidecar struct {
	Envoy *any `json:"envoy,omitempty"`
}

func (w Workload) RemoveEmptySlices() {

	if w.Spec.Containers != nil {
		for _, c := range *w.Spec.Containers {
			if c.Args == nil || len(*c.Args) < 1 {
				c.Args = nil
			}
		}
	}

	if w.Spec.FirewallConfig != nil {

		if w.Spec.FirewallConfig.External != nil {
			if w.Spec.FirewallConfig.External.InboundAllowCIDR != nil && len(*w.Spec.FirewallConfig.External.InboundAllowCIDR) < 1 {
				w.Spec.FirewallConfig.External.InboundAllowCIDR = nil
			}

			if w.Spec.FirewallConfig.External.OutboundAllowCIDR != nil && len(*w.Spec.FirewallConfig.External.OutboundAllowCIDR) < 1 {
				w.Spec.FirewallConfig.External.OutboundAllowCIDR = nil
			}

			if w.Spec.FirewallConfig.External.OutboundAllowHostname != nil && len(*w.Spec.FirewallConfig.External.OutboundAllowHostname) < 1 {
				w.Spec.FirewallConfig.External.OutboundAllowHostname = nil
			}
		}

		if w.Spec.FirewallConfig.Internal != nil && w.Spec.FirewallConfig.Internal.InboundAllowWorkload != nil && len(*w.Spec.FirewallConfig.Internal.InboundAllowWorkload) < 1 {
			w.Spec.FirewallConfig.Internal.InboundAllowWorkload = nil
		}
	}
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
