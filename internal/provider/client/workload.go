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
	Spec   *WorkloadSpec   `json:"spec,omitempty"`
	Status *WorkloadStatus `json:"status,omitempty"`
}

// WorkloadSpec - Workload Specifications
type WorkloadSpec struct {
	Type           *string          `json:"type,omitempty"`
	IdentityLink   *string          `json:"identityLink,omitempty"`
	Containers     *[]ContainerSpec `json:"containers,omitempty"`
	FirewallConfig *FirewallSpec    `json:"firewallConfig,omitempty"`
	DefaultOptions *Options         `json:"defaultOptions,omitempty"`
	LocalOptions   *[]Options       `json:"localOptions,omitempty"`
	Update         bool             `json:"-"`
}

// WorkloadSpecUpdate - Workload Specifications
type WorkloadSpecUpdate struct {
	Type           *string          `json:"type,omitempty"`
	IdentityLink   *string          `json:"identityLink"`
	Containers     *[]ContainerSpec `json:"containers,omitempty"`
	FirewallConfig *FirewallSpec    `json:"firewallConfig,omitempty"`
	DefaultOptions *Options         `json:"defaultOptions,omitempty"`
	LocalOptions   *[]Options       `json:"localOptions"`
}

func (p WorkloadSpec) MarshalJSON() ([]byte, error) {

	type localWorkload WorkloadSpec

	if p.Update && (p.IdentityLink == nil || *p.IdentityLink == "") {
		return json.Marshal(WorkloadSpecUpdate{
			Type:           p.Type,
			IdentityLink:   p.IdentityLink,
			Containers:     p.Containers,
			FirewallConfig: p.FirewallConfig,
			DefaultOptions: p.DefaultOptions,
			LocalOptions:   p.LocalOptions,
		})
	}
	return json.Marshal(localWorkload(p))
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
	Env              *[]NameValue     `json:"env,omitempty"`
	Args             *[]string        `json:"args,omitempty"`
	Volumes          *[]VolumeSpec    `json:"volumes,omitempty"`
	Metrics          *Metrics         `json:"metrics,omitempty"`
	Command          *string          `json:"command,omitempty"`
	WorkingDirectory *string          `json:"workingDir,omitempty"`
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
	InboundAllowCIDR      *[]string `json:"inboundAllowCIDR,omitempty"`
	OutboundAllowCIDR     *[]string `json:"outboundAllowCIDR,omitempty"`
	OutboundAllowHostname *[]string `json:"outboundAllowHostname,omitempty"`
	Update                bool      `json:"-"`
}

// FirewallSpecExternalUpdate - Firewall Spec External
type FirewallSpecExternalUpdate struct {
	InboundAllowCIDR      *[]string `json:"inboundAllowCIDR"`
	OutboundAllowCIDR     *[]string `json:"outboundAllowCIDR"`
	OutboundAllowHostname *[]string `json:"outboundAllowHostname"`
}

func (p FirewallSpecExternal) MarshalJSON() ([]byte, error) {

	type localFirewallSpecExternal FirewallSpecExternal

	if p.Update {
		return json.Marshal(FirewallSpecExternalUpdate{
			InboundAllowCIDR:      p.InboundAllowCIDR,
			OutboundAllowCIDR:     p.OutboundAllowCIDR,
			OutboundAllowHostname: p.OutboundAllowHostname,
		})
	}
	return json.Marshal(localFirewallSpecExternal(p))
}

// FirewallSpecInternal - Firewall Spec Internal
type FirewallSpecInternal struct {
	InboundAllowType     *string   `json:"inboundAllowType,omitempty"`
	InboundAllowWorkload *[]string `json:"inboundAllowWorkload,omitempty"`
	Update               bool      `json:"-"`
}

// FirewallSpecInternaUpdate - Firewall Spec Internal
type FirewallSpecInternaUpdate struct {
	InboundAllowType     *string   `json:"inboundAllowType"`
	InboundAllowWorkload *[]string `json:"inboundAllowWorkload"`
}

func (p FirewallSpecInternal) MarshalJSON() ([]byte, error) {

	type localFirewallSpecInternal FirewallSpecInternal

	if p.Update {
		return json.Marshal(FirewallSpecInternaUpdate{
			InboundAllowType:     p.InboundAllowType,
			InboundAllowWorkload: p.InboundAllowWorkload,
		})
	}
	return json.Marshal(localFirewallSpecInternal(p))
}

// WorkloadStatus - Workload Status
type WorkloadStatus struct {
	ParentID            *string            `json:"parentId,omitempty"`
	CanonicalEndpoint   *string            `json:"canonicalEndpoint,omitempty"`
	Endpoint            *string            `json:"endpoint,omitempty"`
	InternalName        *string            `json:"internalName,omitempty"`
	CurrentReplicaCount *int               `json:"currentReplicaCount,omitempty"`
	HealthCheck         *HealthCheckStatus `json:"healthCheck,omitempty"`
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

// HealthCheckSpec - Health Check Spec (used my readiness and liveness probes)
type HealthCheckSpec struct {
	Exec                *Exec      `json:"exec,omitempty"`
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
	Uri  *string `json:"uri,omitempty"`
	Path *string `json:"path,omitempty"`
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

func (w Workload) RemoveEmptySlices() {

	for _, c := range *w.Spec.Containers {
		if c.Args == nil || len(*c.Args) < 1 {
			c.Args = nil
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

			if w.Spec.FirewallConfig.External.InboundAllowCIDR != nil && len(*w.Spec.FirewallConfig.External.OutboundAllowHostname) < 1 {
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

	workload.(*Workload).RemoveEmptySlices()

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
