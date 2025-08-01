package workload

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Container //

type ContainerModel struct {
	Name             types.String                `tfsdk:"name"`
	Image            types.String                `tfsdk:"image"`
	WorkingDirectory types.String                `tfsdk:"working_directory"`
	Metrics          []ContainerMetricsModel     `tfsdk:"metrics"`
	Port             types.Int32                 `tfsdk:"port"`
	Ports            []ContainerPortModel        `tfsdk:"ports"`
	Memory           types.String                `tfsdk:"memory"`
	ReadinessProbe   []ContainerHealthCheckModel `tfsdk:"readiness_probe"`
	LivenessProbe    []ContainerHealthCheckModel `tfsdk:"liveness_probe"`
	Cpu              types.String                `tfsdk:"cpu"`
	MinCpu           types.String                `tfsdk:"min_cpu"`
	MinMemory        types.String                `tfsdk:"min_memory"`
	Env              types.Map                   `tfsdk:"env"`
	GpuNvidia        []ContainerGpuNvidiaModel   `tfsdk:"gpu_nvidia"`
	GpuCustom        []ContainerGpuCustomModel   `tfsdk:"gpu_custom"`
	InheritEnv       types.Bool                  `tfsdk:"inherit_env"`
	Command          types.String                `tfsdk:"command"`
	Args             types.Set                   `tfsdk:"args"`
	Lifecycle        []ContainerLifecycleModel   `tfsdk:"lifecycle"`
	Volumes          []ContainerVolumeModel      `tfsdk:"volume"`
}

// Container -> Metrics //

type ContainerMetricsModel struct {
	Port types.Int32  `tfsdk:"port"`
	Path types.String `tfsdk:"path"`
}

// Container -> Port //

type ContainerPortModel struct {
	Protocol types.String `tfsdk:"protocol"`
	Number   types.Int32  `tfsdk:"number"`
}

// Container -> Health Check //

type ContainerHealthCheckModel struct {
	Exec                []ContainerExecModel                 `tfsdk:"exec"`
	Grpc                []ContainerHealthCheckGrpcModel      `tfsdk:"grpc"`
	TcpSocket           []ContainerHealthCheckTcpSocketModel `tfsdk:"tcp_socket"`
	HttpGet             []ContainerHealthCheckHttpGetModel   `tfsdk:"http_get"`
	InitialDelaySeconds types.Int32                          `tfsdk:"initial_delay_seconds"`
	PeriodSeconds       types.Int32                          `tfsdk:"period_seconds"`
	TimeoutSeconds      types.Int32                          `tfsdk:"timeout_seconds"`
	SuccessThreshold    types.Int32                          `tfsdk:"success_threshold"`
	FailureThreshold    types.Int32                          `tfsdk:"failure_threshold"`
}

// Container -> Health Check -> GRPC //

type ContainerHealthCheckGrpcModel struct {
	Port types.Int32 `tfsdk:"port"`
}

// Container -> Health Check -> TCP Socket //

type ContainerHealthCheckTcpSocketModel struct {
	Port types.Int32 `tfsdk:"port"`
}

// Container -> Health Check -> HTTP Get //

type ContainerHealthCheckHttpGetModel struct {
	Path        types.String `tfsdk:"path"`
	Port        types.Int32  `tfsdk:"port"`
	HttpHeaders types.Map    `tfsdk:"http_headers"`
	Scheme      types.String `tfsdk:"scheme"`
}

// Container -> GPU Nvidia //

type ContainerGpuNvidiaModel struct {
	Model    types.String `tfsdk:"model"`
	Quantity types.Int32  `tfsdk:"quantity"`
}

// Container -> GPU Custom //

type ContainerGpuCustomModel struct {
	Resource     types.String `tfsdk:"resource"`
	RuntimeClass types.String `tfsdk:"runtime_class"`
	Quantity     types.Int32  `tfsdk:"quantity"`
}

// Container -> Lifecycle //

type ContainerLifecycleModel struct {
	PostStart []ContainerLifecycleSpecModel `tfsdk:"post_start"`
	PreStop   []ContainerLifecycleSpecModel `tfsdk:"pre_stop"`
}

// Container -> Lifecycle -> Spec //

type ContainerLifecycleSpecModel struct {
	Exec []ContainerExecModel `tfsdk:"exec"`
}

// Container -> Volume //

type ContainerVolumeModel struct {
	Uri            types.String `tfsdk:"uri"`
	RecoveryPolicy types.String `tfsdk:"recovery_policy"`
	Path           types.String `tfsdk:"path"`
}

// Container -> Exec //

type ContainerExecModel struct {
	Command types.Set `tfsdk:"command"`
}

// Firewall //

type FirewallModel struct {
	External []FirewallExternalModel `tfsdk:"external"`
	Internal []FirewallInternalModel `tfsdk:"internal"`
}

// Firewall -> External //

type FirewallExternalModel struct {
	InboundAllowCidr      types.Set                                `tfsdk:"inbound_allow_cidr"`
	InboundBlockedCidr    types.Set                                `tfsdk:"inbound_blocked_cidr"`
	OutboundAllowHostname types.Set                                `tfsdk:"outbound_allow_hostname"`
	OutboundAllowPort     []FirewallExternalOutboundAllowPortModel `tfsdk:"outbound_allow_port"`
	OutboundAllowCidr     types.Set                                `tfsdk:"outbound_allow_cidr"`
	OutboundBlockedCidr   types.Set                                `tfsdk:"outbound_blocked_cidr"`
	Http                  []FirewallExternalHttpModel              `tfsdk:"http"`
}

// Firewall -> External -> Outbound Allow Port //

type FirewallExternalOutboundAllowPortModel struct {
	Protocol types.String `tfsdk:"protocol"`
	Number   types.Int32  `tfsdk:"number"`
}

// Firewall -> External -> HTTP //

type FirewallExternalHttpModel struct {
	InboundHeaderFilter []FirewallExternalHttpHeaderFilterModel `tfsdk:"inbound_header_filter"`
}

// Firewall -> External -> HTTP -> Header Filter //

type FirewallExternalHttpHeaderFilterModel struct {
	Key           types.String `tfsdk:"key"`
	AllowedValues types.Set    `tfsdk:"allowed_values"`
	BlockedValues types.Set    `tfsdk:"blocked_values"`
}

// Firewall -> Internal //

type FirewallInternalModel struct {
	InboundAllowType     types.String `tfsdk:"inbound_allow_type"`
	InboundAllowWorkload types.Set    `tfsdk:"inbound_allow_workload"`
}

// Options //

type OptionsModel struct {
	Autoscaling    []OptionsAutoscalingModel `tfsdk:"autoscaling"`
	TimeoutSeconds types.Int32               `tfsdk:"timeout_seconds"`
	CapacityAI     types.Bool                `tfsdk:"capacity_ai"`
	Debug          types.Bool                `tfsdk:"debug"`
	Suspend        types.Bool                `tfsdk:"suspend"`
	MultiZone      []OptionsMultiZoneModel   `tfsdk:"multi_zone"`
}

// Options -> Autoscaling //

type OptionsAutoscalingModel struct {
	Metric           types.String                   `tfsdk:"metric"`
	Multi            []OptionsAutoscalingMultiModel `tfsdk:"multi"`
	MetricPercentile types.String                   `tfsdk:"metric_percentile"`
	Target           types.Int32                    `tfsdk:"target"`
	MinScale         types.Int32                    `tfsdk:"min_scale"`
	MaxScale         types.Int32                    `tfsdk:"max_scale"`
	ScaleToZeroDelay types.Int32                    `tfsdk:"scale_to_zero_delay"`
	MaxConcurrency   types.Int32                    `tfsdk:"max_concurrency"`
	Keda             []OptionsAutoscalingKedaModel  `tfsdk:"keda"`
}

// Options -> Autoscaling -> Multi //

type OptionsAutoscalingMultiModel struct {
	Metric types.String `tfsdk:"metric"`
	Target types.Int32  `tfsdk:"target"`
}

// Options -> Autoscaling -> Keda //

type OptionsAutoscalingKedaModel struct {
	Triggers []OptionsAutoscalingKedaTriggerModel  `tfsdk:"trigger"`
	Advanced []OptionsAutoscalingKedaAdvancedModel `tfsdk:"advanced"`
}

// Options -> Autoscaling -> Keda -> Trigger //

type OptionsAutoscalingKedaTriggerModel struct {
	Type             types.String `tfsdk:"type"`
	Metadata         types.Map    `tfsdk:"metadata"`
	Name             types.String `tfsdk:"name"`
	UseCachedMetrics types.Bool   `tfsdk:"use_cached_metrics"`
	MetricType       types.String `tfsdk:"metric_type"`
}

// Options -> Autoscaling -> Keda -> Advanced //

type OptionsAutoscalingKedaAdvancedModel struct {
	ScalingModifiers []OptionsAutoscalingKedaAdvancedScalingModifiersModel `tfsdk:"scaling_modifiers"`
}

// Options -> Autoscaling -> Keda -> Advanced -> Scaling Modifier //

type OptionsAutoscalingKedaAdvancedScalingModifiersModel struct {
	Target           types.String `tfsdk:"target"`
	ActivationTarget types.String `tfsdk:"activation_target"`
	MetricType       types.String `tfsdk:"metric_type"`
	Formula          types.String `tfsdk:"formula"`
}

// Options -> Muli Zone //

type OptionsMultiZoneModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

// Local Options //

type LocalOptionsModel struct {
	OptionsModel
	Location types.String `tfsdk:"location"`
}

// Job //

type JobModel struct {
	Schedule              types.String `tfsdk:"schedule"`
	ConcurrencyPolicy     types.String `tfsdk:"concurrency_policy"`
	HistoryLimit          types.Int32  `tfsdk:"history_limit"`
	RestartPolicy         types.String `tfsdk:"restart_policy"`
	ActiveDeadlineSeconds types.Int32  `tfsdk:"active_deadline_seconds"`
}

// Sidecar //

type SidecarModel struct {
	Envoy types.String `tfsdk:"envoy"`
}

// Rollout Options //

type RolloutOptionsModel struct {
	MinReadySeconds               types.Int32  `tfsdk:"min_ready_seconds"`
	MaxUnavailableReplicas        types.String `tfsdk:"max_unavailable_replicas"`
	MaxSurgeReplicas              types.String `tfsdk:"max_surge_replicas"`
	ScalingPolicy                 types.String `tfsdk:"scaling_policy"`
	TerminationGracePeriodSeconds types.Int32  `tfsdk:"termination_grace_period_seconds"`
}

// Security Options //

type SecurityOptionsModel struct {
	FileSystemGroupId types.Int32 `tfsdk:"file_system_group_id"`
}

// Load Balancer //

type LoadBalancerModel struct {
	Direct        []LoadBalancerDirectModel      `tfsdk:"direct"`
	GeoLocation   []LoadBalancerGeoLocationModel `tfsdk:"geo_location"`
	ReplicaDirect types.Bool                     `tfsdk:"replica_direct"`
}

// Load Balancer -> Direct //

type LoadBalancerDirectModel struct {
	Enabled types.Bool                    `tfsdk:"enabled"`
	Ports   []LoadBalancerDirectPortModel `tfsdk:"port"`
	IpSet   types.String                  `tfsdk:"ipset"`
}

// Load Balancer -> Direct -> Port //

type LoadBalancerDirectPortModel struct {
	ExternalPort  types.Int32  `tfsdk:"external_port"`
	Protocol      types.String `tfsdk:"protocol"`
	Scheme        types.String `tfsdk:"scheme"`
	ContainerPort types.Int32  `tfsdk:"container_port"`
}

// Load Balancer -> Geo Location //

type LoadBalancerGeoLocationModel struct {
	Enabled types.Bool                            `tfsdk:"enabled"`
	Headers []LoadBalancerGeoLocationHeadersModel `tfsdk:"headers"`
}

// Load Balancer -> Geo Location -> Headers //

type LoadBalancerGeoLocationHeadersModel struct {
	Asn     types.String `tfsdk:"asn"`
	City    types.String `tfsdk:"city"`
	Country types.String `tfsdk:"country"`
	Region  types.String `tfsdk:"region"`
}

// Request Retry On //

type RequestRetryPolicyModel struct {
	Attempts types.Int32 `tfsdk:"attempts"`
	RetryOn  types.Set   `tfsdk:"retry_on"`
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
