package workload

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Container //

type ContainerModel struct {
	Name             types.String `tfsdk:"name"`
	Image            types.String `tfsdk:"image"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
	Metrics          types.List   `tfsdk:"metrics"`
	Port             types.Int32  `tfsdk:"port"`
	Ports            types.List   `tfsdk:"ports"`
	Memory           types.String `tfsdk:"memory"`
	ReadinessProbe   types.List   `tfsdk:"readiness_probe"`
	LivenessProbe    types.List   `tfsdk:"liveness_probe"`
	Cpu              types.String `tfsdk:"cpu"`
	MinCpu           types.String `tfsdk:"min_cpu"`
	MinMemory        types.String `tfsdk:"min_memory"`
	Env              types.Map    `tfsdk:"env"`
	GpuNvidia        types.List   `tfsdk:"gpu_nvidia"`
	GpuCustom        types.List   `tfsdk:"gpu_custom"`
	InheritEnv       types.Bool   `tfsdk:"inherit_env"`
	Command          types.String `tfsdk:"command"`
	Args             types.Set    `tfsdk:"args"`
	Lifecycle        types.List   `tfsdk:"lifecycle"`
	Volumes          types.List   `tfsdk:"volume"`
}

func (c ContainerModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":              types.StringType,
			"image":             types.StringType,
			"working_directory": types.StringType,
			"metrics":           types.ListType{ElemType: ContainerMetricsModel{}.AttributeTypes()},
			"port":              types.Int32Type,
			"ports":             types.ListType{ElemType: ContainerPortModel{}.AttributeTypes()},
			"memory":            types.StringType,
			"readiness_probe":   types.ListType{ElemType: ContainerHealthCheckModel{}.AttributeTypes()},
			"liveness_probe":    types.ListType{ElemType: ContainerHealthCheckModel{}.AttributeTypes()},
			"cpu":               types.StringType,
			"min_cpu":           types.StringType,
			"min_memory":        types.StringType,
			"env":               types.MapType{ElemType: types.StringType},
			"gpu_nvidia":        types.ListType{ElemType: ContainerGpuNvidiaModel{}.AttributeTypes()},
			"gpu_custom":        types.ListType{ElemType: ContainerGpuCustomModel{}.AttributeTypes()},
			"inherit_env":       types.BoolType,
			"command":           types.StringType,
			"args":              types.SetType{ElemType: types.StringType},
			"lifecycle":         types.ListType{ElemType: ContainerLifecycleModel{}.AttributeTypes()},
			"volume":            types.ListType{ElemType: ContainerVolumeModel{}.AttributeTypes()},
		},
	}
}

// Container -> Metrics //

type ContainerMetricsModel struct {
	Port types.Int32  `tfsdk:"port"`
	Path types.String `tfsdk:"path"`
}

func (c ContainerMetricsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"port": types.Int32Type,
			"path": types.StringType,
		},
	}
}

// Container -> Port //

type ContainerPortModel struct {
	Protocol types.String `tfsdk:"protocol"`
	Number   types.Int32  `tfsdk:"number"`
}

func (c ContainerPortModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"protocol": types.StringType,
			"number":   types.Int32Type,
		},
	}
}

// Container -> Health Check //

type ContainerHealthCheckModel struct {
	Exec                types.List  `tfsdk:"exec"`
	Grpc                types.List  `tfsdk:"grpc"`
	TcpSocket           types.List  `tfsdk:"tcp_socket"`
	HttpGet             types.List  `tfsdk:"http_get"`
	InitialDelaySeconds types.Int32 `tfsdk:"initial_delay_seconds"`
	PeriodSeconds       types.Int32 `tfsdk:"period_seconds"`
	TimeoutSeconds      types.Int32 `tfsdk:"timeout_seconds"`
	SuccessThreshold    types.Int32 `tfsdk:"success_threshold"`
	FailureThreshold    types.Int32 `tfsdk:"failure_threshold"`
}

func (c ContainerHealthCheckModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"exec":                  types.ListType{ElemType: ContainerExecModel{}.AttributeTypes()},
			"grpc":                  types.ListType{ElemType: ContainerHealthCheckGrpcModel{}.AttributeTypes()},
			"tcp_socket":            types.ListType{ElemType: ContainerHealthCheckTcpSocketModel{}.AttributeTypes()},
			"http_get":              types.ListType{ElemType: ContainerHealthCheckHttpGetModel{}.AttributeTypes()},
			"initial_delay_seconds": types.Int32Type,
			"period_seconds":        types.Int32Type,
			"timeout_seconds":       types.Int32Type,
			"success_threshold":     types.Int32Type,
			"failure_threshold":     types.Int32Type,
		},
	}
}

// Container -> Health Check -> GRPC //

type ContainerHealthCheckGrpcModel struct {
	Port types.Int32 `tfsdk:"port"`
}

func (c ContainerHealthCheckGrpcModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"port": types.Int32Type,
		},
	}
}

// Container -> Health Check -> TCP Socket //

type ContainerHealthCheckTcpSocketModel struct {
	Port types.Int32 `tfsdk:"port"`
}

func (c ContainerHealthCheckTcpSocketModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"port": types.Int32Type,
		},
	}
}

// Container -> Health Check -> HTTP Get //

type ContainerHealthCheckHttpGetModel struct {
	Path        types.String `tfsdk:"path"`
	Port        types.Int32  `tfsdk:"port"`
	HttpHeaders types.Map    `tfsdk:"http_headers"`
	Scheme      types.String `tfsdk:"scheme"`
}

func (c ContainerHealthCheckHttpGetModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"path":         types.StringType,
			"port":         types.Int32Type,
			"http_headers": types.MapType{ElemType: types.StringType},
			"scheme":       types.StringType,
		},
	}
}

// Container -> GPU Nvidia //

type ContainerGpuNvidiaModel struct {
	Model    types.String `tfsdk:"model"`
	Quantity types.Int32  `tfsdk:"quantity"`
}

func (c ContainerGpuNvidiaModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"model":    types.StringType,
			"quantity": types.Int32Type,
		},
	}
}

// Container -> GPU Custom //

type ContainerGpuCustomModel struct {
	Resource     types.String `tfsdk:"resource"`
	RuntimeClass types.String `tfsdk:"runtime_class"`
	Quantity     types.Int32  `tfsdk:"quantity"`
}

func (c ContainerGpuCustomModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource":      types.StringType,
			"runtime_class": types.StringType,
			"quantity":      types.Int32Type,
		},
	}
}

// Container -> Lifecycle //

type ContainerLifecycleModel struct {
	PostStart types.List `tfsdk:"post_start"`
	PreStop   types.List `tfsdk:"pre_stop"`
}

func (c ContainerLifecycleModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"post_start": types.ListType{ElemType: ContainerLifecycleSpecModel{}.AttributeTypes()},
			"pre_stop":   types.ListType{ElemType: ContainerLifecycleSpecModel{}.AttributeTypes()},
		},
	}
}

// Container -> Lifecycle -> Spec //

type ContainerLifecycleSpecModel struct {
	Exec types.List `tfsdk:"exec"`
}

func (c ContainerLifecycleSpecModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"exec": types.ListType{ElemType: ContainerExecModel{}.AttributeTypes()},
		},
	}
}

// Container -> Volume //

type ContainerVolumeModel struct {
	Uri            types.String `tfsdk:"uri"`
	RecoveryPolicy types.String `tfsdk:"recovery_policy"`
	Path           types.String `tfsdk:"path"`
}

func (c ContainerVolumeModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"uri":             types.StringType,
			"recovery_policy": types.StringType,
			"path":            types.StringType,
		},
	}
}

// Container -> Exec //

type ContainerExecModel struct {
	Command types.Set `tfsdk:"command"`
}

func (c ContainerExecModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"command": types.SetType{ElemType: types.StringType},
		},
	}
}

// Firewall //

type FirewallModel struct {
	External types.List `tfsdk:"external"`
	Internal types.List `tfsdk:"internal"`
}

func (f FirewallModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"external": types.ListType{ElemType: FirewallExternalModel{}.AttributeTypes()},
			"internal": types.ListType{ElemType: FirewallInternalModel{}.AttributeTypes()},
		},
	}
}

// Firewall -> External //

type FirewallExternalModel struct {
	InboundAllowCidr      types.Set  `tfsdk:"inbound_allow_cidr"`
	InboundBlockedCidr    types.Set  `tfsdk:"inbound_blocked_cidr"`
	OutboundAllowHostname types.Set  `tfsdk:"outbound_allow_hostname"`
	OutboundAllowPort     types.List `tfsdk:"outbound_allow_port"`
	OutboundAllowCidr     types.Set  `tfsdk:"outbound_allow_cidr"`
	OutboundBlockedCidr   types.Set  `tfsdk:"outbound_blocked_cidr"`
	Http                  types.List `tfsdk:"http"`
}

func (f FirewallExternalModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"inbound_allow_cidr":      types.SetType{ElemType: types.StringType},
			"inbound_blocked_cidr":    types.SetType{ElemType: types.StringType},
			"outbound_allow_hostname": types.SetType{ElemType: types.StringType},
			"outbound_allow_port":     types.ListType{ElemType: FirewallExternalOutboundAllowPortModel{}.AttributeTypes()},
			"outbound_allow_cidr":     types.SetType{ElemType: types.StringType},
			"outbound_blocked_cidr":   types.SetType{ElemType: types.StringType},
			"http":                    types.ListType{ElemType: FirewallExternalHttpModel{}.AttributeTypes()},
		},
	}
}

// Firewall -> External -> Outbound Allow Port //

type FirewallExternalOutboundAllowPortModel struct {
	Protocol types.String `tfsdk:"protocol"`
	Number   types.Int32  `tfsdk:"number"`
}

func (f FirewallExternalOutboundAllowPortModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"protocol": types.StringType,
			"number":   types.Int32Type,
		},
	}
}

// Firewall -> External -> HTTP //

type FirewallExternalHttpModel struct {
	InboundHeaderFilter types.List `tfsdk:"inbound_header_filter"`
}

func (f FirewallExternalHttpModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"inbound_header_filter": types.ListType{ElemType: FirewallExternalHttpHeaderFilterModel{}.AttributeTypes()},
		},
	}
}

// Firewall -> External -> HTTP -> Header Filter //

type FirewallExternalHttpHeaderFilterModel struct {
	Key           types.String `tfsdk:"key"`
	AllowedValues types.Set    `tfsdk:"allowed_values"`
	BlockedValues types.Set    `tfsdk:"blocked_values"`
}

func (f FirewallExternalHttpHeaderFilterModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":            types.StringType,
			"allowed_values": types.SetType{ElemType: types.StringType},
			"blocked_values": types.SetType{ElemType: types.StringType},
		},
	}
}

// Firewall -> Internal //

type FirewallInternalModel struct {
	InboundAllowType     types.String `tfsdk:"inbound_allow_type"`
	InboundAllowWorkload types.Set    `tfsdk:"inbound_allow_workload"`
}

func (f FirewallInternalModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"inbound_allow_type":     types.StringType,
			"inbound_allow_workload": types.SetType{ElemType: types.StringType},
		},
	}
}

// Options //

type OptionsModel struct {
	Autoscaling    types.List  `tfsdk:"autoscaling"`
	TimeoutSeconds types.Int32 `tfsdk:"timeout_seconds"`
	CapacityAI     types.Bool  `tfsdk:"capacity_ai"`
	Debug          types.Bool  `tfsdk:"debug"`
	Suspend        types.Bool  `tfsdk:"suspend"`
	MultiZone      types.List  `tfsdk:"multi_zone"`
}

func (o OptionsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"autoscaling":     types.ListType{ElemType: OptionsAutoscalingModel{}.AttributeTypes()},
			"timeout_seconds": types.Int32Type,
			"capacity_ai":     types.BoolType,
			"debug":           types.BoolType,
			"suspend":         types.BoolType,
			"multi_zone":      types.ListType{ElemType: OptionsMultiZoneModel{}.AttributeTypes()},
		},
	}
}

// Options -> Autoscaling //

type OptionsAutoscalingModel struct {
	Metric           types.String `tfsdk:"metric"`
	Multi            types.List   `tfsdk:"multi"`
	MetricPercentile types.String `tfsdk:"metric_percentile"`
	Target           types.Int32  `tfsdk:"target"`
	MinScale         types.Int32  `tfsdk:"min_scale"`
	MaxScale         types.Int32  `tfsdk:"max_scale"`
	ScaleToZeroDelay types.Int32  `tfsdk:"scale_to_zero_delay"`
	MaxConcurrency   types.Int32  `tfsdk:"max_concurrency"`
	Keda             types.List   `tfsdk:"keda"`
}

func (o OptionsAutoscalingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"metric":              types.StringType,
			"multi":               types.ListType{ElemType: OptionsAutoscalingMultiModel{}.AttributeTypes()},
			"metric_percentile":   types.StringType,
			"target":              types.Int32Type,
			"min_scale":           types.Int32Type,
			"max_scale":           types.Int32Type,
			"scale_to_zero_delay": types.Int32Type,
			"max_concurrency":     types.Int32Type,
			"keda":                types.ListType{ElemType: OptionsAutoscalingKedaModel{}.AttributeTypes()},
		},
	}
}

// Options -> Autoscaling -> Multi //

type OptionsAutoscalingMultiModel struct {
	Metric types.String `tfsdk:"metric"`
	Target types.Int32  `tfsdk:"target"`
}

func (o OptionsAutoscalingMultiModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"metric": types.StringType,
			"target": types.Int32Type,
		},
	}
}

// Options -> Autoscaling -> Keda //

type OptionsAutoscalingKedaModel struct {
	Triggers types.List `tfsdk:"trigger"`
	Advanced types.List `tfsdk:"advanced"`
}

func (o OptionsAutoscalingKedaModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"trigger":  types.ListType{ElemType: OptionsAutoscalingKedaTriggerModel{}.AttributeTypes()},
			"advanced": types.ListType{ElemType: OptionsAutoscalingKedaAdvancedModel{}.AttributeTypes()},
		},
	}
}

// Options -> Autoscaling -> Keda -> Trigger //

type OptionsAutoscalingKedaTriggerModel struct {
	Type             types.String `tfsdk:"type"`
	Metadata         types.Map    `tfsdk:"metadata"`
	Name             types.String `tfsdk:"name"`
	UseCachedMetrics types.Bool   `tfsdk:"use_cached_metrics"`
	MetricType       types.String `tfsdk:"metric_type"`
}

func (o OptionsAutoscalingKedaTriggerModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":               types.StringType,
			"metadata":           types.MapType{ElemType: types.StringType},
			"name":               types.StringType,
			"use_cached_metrics": types.BoolType,
			"metric_type":        types.StringType,
		},
	}
}

// Options -> Autoscaling -> Keda -> Advanced //

type OptionsAutoscalingKedaAdvancedModel struct {
	ScalingModifiers types.List `tfsdk:"scaling_modifiers"`
}

func (o OptionsAutoscalingKedaAdvancedModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"scaling_modifiers": types.ListType{ElemType: OptionsAutoscalingKedaAdvancedScalingModifiersModel{}.AttributeTypes()},
		},
	}
}

// Options -> Autoscaling -> Keda -> Advanced -> Scaling Modifier //

type OptionsAutoscalingKedaAdvancedScalingModifiersModel struct {
	Target           types.String `tfsdk:"target"`
	ActivationTarget types.String `tfsdk:"activation_target"`
	MetricType       types.String `tfsdk:"metric_type"`
	Formula          types.String `tfsdk:"formula"`
}

func (o OptionsAutoscalingKedaAdvancedScalingModifiersModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"target":            types.StringType,
			"activation_target": types.StringType,
			"metric_type":       types.StringType,
			"formula":           types.StringType,
		},
	}
}

// Options -> Muli Zone //

type OptionsMultiZoneModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

func (o OptionsMultiZoneModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
		},
	}
}

// Local Options //

type LocalOptionsModel struct {
	OptionsModel
	Location types.String `tfsdk:"location"`
}

func (o LocalOptionsModel) AttributeTypes() attr.Type {
	// Get the attribute types from OptionsModel
	base := OptionsModel{}.AttributeTypes().(types.ObjectType)

	// Copy the map to avoid mutating the original
	merged := map[string]attr.Type{}
	maps.Copy(merged, base.AttrTypes)

	// Add or override attributes
	merged["location"] = types.StringType

	// Return merged object type
	return types.ObjectType{
		AttrTypes: merged,
	}
}

// Job //

type JobModel struct {
	Schedule              types.String `tfsdk:"schedule"`
	ConcurrencyPolicy     types.String `tfsdk:"concurrency_policy"`
	HistoryLimit          types.Int32  `tfsdk:"history_limit"`
	RestartPolicy         types.String `tfsdk:"restart_policy"`
	ActiveDeadlineSeconds types.Int32  `tfsdk:"active_deadline_seconds"`
}

func (j JobModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"schedule":                types.StringType,
			"concurrency_policy":      types.StringType,
			"history_limit":           types.Int32Type,
			"restart_policy":          types.StringType,
			"active_deadline_seconds": types.Int32Type,
		},
	}
}

// Sidecar //

type SidecarModel struct {
	Envoy types.String `tfsdk:"envoy"`
}

func (s SidecarModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"envoy": types.StringType,
		},
	}
}

// Rollout Options //

type RolloutOptionsModel struct {
	MinReadySeconds               types.Int32  `tfsdk:"min_ready_seconds"`
	MaxUnavailableReplicas        types.String `tfsdk:"max_unavailable_replicas"`
	MaxSurgeReplicas              types.String `tfsdk:"max_surge_replicas"`
	ScalingPolicy                 types.String `tfsdk:"scaling_policy"`
	TerminationGracePeriodSeconds types.Int32  `tfsdk:"termination_grace_period_seconds"`
}

func (r RolloutOptionsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"min_ready_seconds":                types.Int32Type,
			"max_unavailable_replicas":         types.StringType,
			"max_surge_replicas":               types.StringType,
			"scaling_policy":                   types.StringType,
			"termination_grace_period_seconds": types.Int32Type,
		},
	}
}

// Security Options //

type SecurityOptionsModel struct {
	FileSystemGroupId types.Int32 `tfsdk:"file_system_group_id"`
}

func (s SecurityOptionsModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"file_system_group_id": types.Int32Type,
		},
	}
}

// Load Balancer //

type LoadBalancerModel struct {
	Direct        types.List `tfsdk:"direct"`
	GeoLocation   types.List `tfsdk:"geo_location"`
	ReplicaDirect types.Bool `tfsdk:"replica_direct"`
}

func (l LoadBalancerModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"direct":         types.ListType{ElemType: LoadBalancerDirectModel{}.AttributeTypes()},
			"geo_location":   types.ListType{ElemType: LoadBalancerGeoLocationModel{}.AttributeTypes()},
			"replica_direct": types.BoolType,
		},
	}
}

// Load Balancer -> Direct //

type LoadBalancerDirectModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Ports   types.List   `tfsdk:"port"`
	IpSet   types.String `tfsdk:"ipset"`
}

func (l LoadBalancerDirectModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
			"port":    types.ListType{ElemType: LoadBalancerDirectPortModel{}.AttributeTypes()},
			"ipset":   types.StringType,
		},
	}
}

// Load Balancer -> Direct -> Port //

type LoadBalancerDirectPortModel struct {
	ExternalPort  types.Int32  `tfsdk:"external_port"`
	Protocol      types.String `tfsdk:"protocol"`
	Scheme        types.String `tfsdk:"scheme"`
	ContainerPort types.Int32  `tfsdk:"container_port"`
}

func (l LoadBalancerDirectPortModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"external_port":  types.Int32Type,
			"protocol":       types.StringType,
			"scheme":         types.StringType,
			"container_port": types.Int32Type,
		},
	}
}

// Load Balancer -> Geo Location //

type LoadBalancerGeoLocationModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
	Headers types.List `tfsdk:"headers"`
}

func (l LoadBalancerGeoLocationModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
			"headers": types.ListType{ElemType: LoadBalancerGeoLocationHeadersModel{}.AttributeTypes()},
		},
	}
}

// Load Balancer -> Geo Location -> Headers //

type LoadBalancerGeoLocationHeadersModel struct {
	Asn     types.String `tfsdk:"asn"`
	City    types.String `tfsdk:"city"`
	Country types.String `tfsdk:"country"`
	Region  types.String `tfsdk:"region"`
}

func (l LoadBalancerGeoLocationHeadersModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"asn":     types.StringType,
			"city":    types.StringType,
			"country": types.StringType,
			"region":  types.StringType,
		},
	}
}

// Request Retry On //

type RequestRetryPolicyModel struct {
	Attempts types.Int32 `tfsdk:"attempts"`
	RetryOn  types.Set   `tfsdk:"retry_on"`
}

func (r RequestRetryPolicyModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attempts": types.Int32Type,
			"retry_on": types.SetType{ElemType: types.StringType},
		},
	}
}

// Status //

type StatusModel struct {
	ParentId            types.String `tfsdk:"parent_id"`
	CanonicalEndpoint   types.String `tfsdk:"canonical_endpoint"`
	Endpoint            types.String `tfsdk:"endpoint"`
	InternalName        types.String `tfsdk:"internal_name"`
	HealthCheck         types.List   `tfsdk:"health_check"`
	CurrentReplicaCount types.Int32  `tfsdk:"current_replica_count"`
	ResolvedImages      types.List   `tfsdk:"resolved_images"`
	LoadBalancer        types.List   `tfsdk:"load_balancer"`
}

func (s StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"parent_id":             types.StringType,
			"canonical_endpoint":    types.StringType,
			"endpoint":              types.StringType,
			"internal_name":         types.StringType,
			"health_check":          types.ListType{ElemType: StatusHealthCheckModel{}.AttributeTypes()},
			"current_replica_count": types.Int32Type,
			"resolved_images":       types.ListType{ElemType: StatusResolvedImagesModel{}.AttributeTypes()},
			"load_balancer":         types.ListType{ElemType: StatusLoadBalancerModel{}.AttributeTypes()},
		},
	}
}

// Status -> Health Check //

type StatusHealthCheckModel struct {
	Active      types.Bool   `tfsdk:"active"`
	Success     types.Bool   `tfsdk:"success"`
	Code        types.Int32  `tfsdk:"code"`
	Message     types.String `tfsdk:"message"`
	Failures    types.Int32  `tfsdk:"failures"`
	Successes   types.Int32  `tfsdk:"successes"`
	LastChecked types.String `tfsdk:"last_checked"`
}

func (s StatusHealthCheckModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"active":       types.BoolType,
			"success":      types.BoolType,
			"code":         types.Int32Type,
			"message":      types.StringType,
			"failures":     types.Int32Type,
			"successes":    types.Int32Type,
			"last_checked": types.StringType,
		},
	}
}

// Status -> Resolved Images //

type StatusResolvedImagesModel struct {
	ResolvedForVersion types.Int32  `tfsdk:"resolved_for_version"`
	ResolvedAt         types.String `tfsdk:"resolved_at"`
	ErrorMessages      types.Set    `tfsdk:"error_messages"`
	Images             types.List   `tfsdk:"images"`
}

func (s StatusResolvedImagesModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resolved_for_version": types.Int32Type,
			"resolved_at":          types.StringType,
			"error_messages":       types.SetType{ElemType: types.StringType},
			"images":               types.ListType{ElemType: StatusResolvedImageModel{}.AttributeTypes()},
		},
	}
}

// Status -> Resolved Image //

type StatusResolvedImageModel struct {
	Digest    types.String `tfsdk:"digest"`
	Manifests types.List   `tfsdk:"manifests"`
}

func (s StatusResolvedImageModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"digest":    types.StringType,
			"manifests": types.ListType{ElemType: StatusResolvedImageManifestModel{}.AttributeTypes()},
		},
	}
}

// Status -> Resolved Image -> Manifest //

type StatusResolvedImageManifestModel struct {
	Image     types.String `tfsdk:"image"`
	MediaType types.String `tfsdk:"media_type"`
	Digest    types.String `tfsdk:"digest"`
	Platform  types.Map    `tfsdk:"platform"`
}

func (s StatusResolvedImageManifestModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"image":      types.StringType,
			"media_type": types.StringType,
			"digest":     types.StringType,
			"platform":   types.MapType{ElemType: types.StringType},
		},
	}
}

// Status -> Load Balancer //

type StatusLoadBalancerModel struct {
	Origin types.String `tfsdk:"origin"`
	Url    types.String `tfsdk:"url"`
}

func (s StatusLoadBalancerModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"origin": types.StringType,
			"url":    types.StringType,
		},
	}
}
