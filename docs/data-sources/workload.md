---
page_title: "cpln_workload Data Source - terraform-provider-cpln"
subcategory: "Workload"
description: |-
  
---
# cpln_workload (Data Source)

Use this data source to access information about an existing [Workload](https://docs.controlplane.com/reference/workload) within Control Plane. 

## Required

- **name** (String) Name of the workload.
- **gvc** (String) Name of the GVC that the specified workload belongs to.

## Outputs

The following attributes are exported:

- **id** (String) The unique identifier for this workload.
- **cpln_id** (String) The ID, in GUID format, of the workload.
- **name** (String) Name of the workload.
- **gvc** (String) Name of the associated GVC.
- **type** (String) Workload type. Either `serverless`, `standard`, `stateful`, or `cron`.
- **description** (String) Description of the workload.
- **tags** (Map of String) Key-value map of resource tags.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **identity_link** (String) Full link to the identity used as the access scope for 3rd party cloud resources.
- **support_dynamic_tags** (Boolean) Indicates if Control Plane automatically redeploys when referenced container images are updated in the registry.
- **extras** (String) Extra Kubernetes modifications. Only used for BYOK.
- **container** (Block List) ([see below](#nestedblock--container)).
- **firewall_spec** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec)).
- **options** (Block List, Max: 1) ([see below](#nestedblock--options)).
- **local_options** (Block List) ([see below](#nestedblock--local_options)).
- **job** (Block List, Max: 1) ([see below](#nestedblock--job)).
- **sidecar** (Block List, Max: 1) ([see below](#nestedblock--sidecar)).
- **rollout_options** (Block List, Max: 1) ([see below](#nestedblock--rollout_options)).
- **security_options** (Block List, Max: 1) ([see below](#nestedblock--security_options)).
- **load_balancer** (Block List, Max: 1) ([see below](#nestedblock--load_balancer)).
- **request_retry_policy** (Block List, Max: 1) ([see below](#nestedblock--request_retry_policy)).
- **status** (Block List) ([see below](#nestedblock--status)).

<a id="nestedblock--container"></a>

### `container`

~> **Note** A workload always exposes at least one container definition.

Read-Only:

- **name** (String) Name of the container. Cannot be `istio-proxy`, `queue-proxy`, or `istio-validation`, and cannot start with `cpln_`.
- **image** (String) The full image and tag path.
- **working_directory** (String) Override for the container working directory. Must be an absolute path.
- **port** (Number) The port the container exposes. Only one container can specify a port. Min: `80`. Max: `65535`. Used by the `serverless` workload type. **Deprecated – use `ports`.**
- **memory** (String) Reserved memory when Capacity AI is disabled, or maximum memory when it is enabled. Default: `128Mi`.
- **cpu** (String) Reserved CPU when Capacity AI is disabled, or maximum CPU when it is enabled. Default: `50m`.
- **min_cpu** (String) Minimum CPU when Capacity AI is enabled.
- **min_memory** (String) Minimum memory when Capacity AI is enabled.
- **env** (Map of String) Environment variables exposed to the container.
- **inherit_env** (Boolean) Indicates whether GVC environment variables are inherited. A variable in `env` overrides the same key from the GVC.
- **command** (String) Override for the container entry point.
- **args** (List of String) Command-line arguments passed to the container in order.
- **metrics** (Block List, Max: 1) ([see below](#nestedblock--container--metrics)).
- **ports** (Block List) ([see below](#nestedblock--container--ports)).
- **readiness_probe** (Block List, Max: 1) ([see below](#nestedblock--container--readiness_probe)).
- **liveness_probe** (Block List, Max: 1) ([see below](#nestedblock--container--liveness_probe)).
- **gpu_nvidia** (Block List, Max: 1) ([see below](#nestedblock--container--gpu_nvidia)).
- **gpu_custom** (Block List, Max: 1) ([see below](#nestedblock--container--gpu_custom)).
- **lifecycle** (Block List, Max: 1) ([see below](#nestedblock--container--lifecycle)).
- **volume** (Block List) ([see below](#nestedblock--container--volume)).

~> **Note** The following container ports are reserved and cannot be used: `8012`, `8022`, `9090`, `9091`, `15000`, `15001`, `15006`, `15020`, `15021`, `15090`, `41000`.

<a id="nestedblock--container--metrics"></a>

### `container.metrics`

Read-Only:

- **port** (Number) Port that exposes custom metrics.
- **path** (String) Path where custom metrics are available.
- **drop_metrics** (List of String) Patterns describing metrics to discard.

<a id="nestedblock--container--ports"></a>

### `container.ports`

Read-Only:

- **protocol** (String) Protocol. One of `http`, `http2`, `tcp`, or `grpc`.
- **number** (Number) Port that the container exposes.

<a id="nestedblock--container--readiness_probe"></a>

### `container.readiness_probe`

Read-Only:

- **initial_delay_seconds** (Number) Initial delay before the probe runs. Default: `10`. Min: `0`. Max: `600`.
- **period_seconds** (Number) Interval between probes. Default: `10`. Min: `1`. Max: `600`.
- **timeout_seconds** (Number) Probe timeout. Default: `1`. Min: `1`. Max: `600`.
- **success_threshold** (Number) Minimum consecutive successes to be considered ready. Default: `1`. Min: `1`. Max: `20`.
- **failure_threshold** (Number) Consecutive failures before marking the container unhealthy. Default: `3`. Min: `1`. Max: `20`.
- **exec** (Block List, Max: 1) ([see below](#nestedblock--container--probe--exec)).
- **grpc** (Block List, Max: 1) ([see below](#nestedblock--container--probe--grpc)).
- **tcp_socket** (Block List, Max: 1) ([see below](#nestedblock--container--probe--tcp_socket)).
- **http_get** (Block List, Max: 1) ([see below](#nestedblock--container--probe--http_get)).

<a id="nestedblock--container--liveness_probe"></a>

### `container.liveness_probe`

Read-Only:

- **initial_delay_seconds** (Number) Initial delay before the probe runs. Default: `10`. Min: `0`. Max: `600`.
- **period_seconds** (Number) Interval between probes. Default: `10`. Min: `1`. Max: `600`.
- **timeout_seconds** (Number) Probe timeout. Default: `1`. Min: `1`. Max: `600`.
- **success_threshold** (Number) Minimum consecutive successes to be considered healthy. Default: `1`. Min: `1`. Max: `20`.
- **failure_threshold** (Number) Consecutive failures before restarting the container. Default: `3`. Min: `1`. Max: `20`.
- **exec** (Block List, Max: 1) ([see below](#nestedblock--container--probe--exec)).
- **grpc** (Block List, Max: 1) ([see below](#nestedblock--container--probe--grpc)).
- **tcp_socket** (Block List, Max: 1) ([see below](#nestedblock--container--probe--tcp_socket)).
- **http_get** (Block List, Max: 1) ([see below](#nestedblock--container--probe--http_get)).

<a id="nestedblock--container--probe--exec"></a>

### `container.*_probe.exec`

Read-Only:

- **command** (List of String) Command executed inside the container when the probe runs.

<a id="nestedblock--container--probe--grpc"></a>

### `container.*_probe.grpc`

Read-Only:

- **port** (Number) gRPC port used for the probe.

<a id="nestedblock--container--probe--tcp_socket"></a>

### `container.*_probe.tcp_socket`

Read-Only:

- **port** (Number) TCP port used for the probe.

<a id="nestedblock--container--probe--http_get"></a>

### `container.*_probe.http_get`

Read-Only:

- **path** (String) HTTP path to query. Default: `/`.
- **port** (Number) Port for the HTTP GET. Min: `80`. Max: `65535`.
- **http_headers** (Map of String) Headers included in the probe request.
- **scheme** (String) HTTP scheme. Either `HTTP` or `HTTPS`. Default: `HTTP`.

<a id="nestedblock--container--gpu_nvidia"></a>

### `container.gpu_nvidia`

Read-Only:

- **model** (String) GPU model (for example, `t4`).
- **quantity** (Number) Number of NVIDIA GPUs attached to the container.

<a id="nestedblock--container--gpu_custom"></a>

### `container.gpu_custom`

Read-Only:

- **resource** (String) Name of the custom GPU resource.
- **quantity** (Number) Number of GPUs requested.
- **runtime_class** (String) Runtime class that must be used with the custom GPU.

<a id="nestedblock--container--lifecycle"></a>

### `container.lifecycle`

Read-Only:

- **post_start** (Block List, Max: 1) ([see below](#nestedblock--container--lifecycle--spec)).
- **pre_stop** (Block List, Max: 1) ([see below](#nestedblock--container--lifecycle--spec)).

<a id="nestedblock--container--lifecycle--spec"></a>

### `container.lifecycle.*`

Read-Only:

- **exec** (Block List, Max: 1) ([see below](#nestedblock--container--lifecycle--exec)).

<a id="nestedblock--container--lifecycle--exec"></a>

### `container.lifecycle.*.exec`

Read-Only:

- **command** (List of String) Command executed during the lifecycle hook.

<a id="nestedblock--container--volume"></a>

### `container.volume`

~> **Note** The following paths are reserved and cannot be used: `/dev`, `/dev/log`, `/tmp`, `/var`, `/var/log`.

~> **Note** Valid URI prefixes include `s3://bucket`, `gs://bucket`, `azureblob://storageAccount/container`, `azurefs://storageAccount/share`, `cpln://secret`, `cpln://volumeset`, and `scratch://`.

Read-Only:

- **uri** (String) URI of a volume hosted in Control Plane (Volume Set) or a supported cloud provider.
- **recovery_policy** (String) Recovery policy for persistent volumes. Either `retain` or `recycle`. **Deprecated – no longer used.**
- **path** (String) File-system path where the volume is mounted inside the container.

<a id="nestedblock--firewall_spec"></a>

### `firewall_spec`

Controls inbound and outbound access for external (public) and internal (service-to-service) traffic. Access is restricted by default.

Read-Only:

- **external** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec--external)).
- **internal** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec--internal)).

<a id="nestedblock--firewall_spec--external"></a>

### `firewall_spec.external`

Read-Only:

- **inbound_allow_cidr** (List of String) IPv4/IPv6 addresses or CIDR blocks allowed to reach the workload. `0.0.0.0/0` opens access to the public internet.
- **inbound_blocked_cidr** (List of String) IPv4/IPv6 addresses or CIDR blocks explicitly denied.
- **outbound_allow_hostname** (List of String) Public hostnames the workload can reach. Wildcards are allowed only as a prefix (for example, `*.amazonaws.com`).
- **outbound_allow_cidr** (List of String) IPv4/IPv6 addresses or CIDR blocks the workload can reach. `0.0.0.0/0` enables outbound access to the public internet.
- **outbound_blocked_cidr** (List of String) IPv4/IPv6 addresses or CIDR blocks that are denied even if allow lists include them.
- **outbound_allow_port** (Block List) ([see below](#nestedblock--firewall_spec--external--outbound_allow_port)).
- **http** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec--external--http)).

<a id="nestedblock--firewall_spec--external--outbound_allow_port"></a>

### `firewall_spec.external.outbound_allow_port`

Allows outbound access to specific ports and protocols.

Read-Only:

- **protocol** (String) Either `http`, `https`, or `tcp`. Default: `tcp`.
- **number** (Number) Port number. Max: `65000`.

<a id="nestedblock--firewall_spec--external--http"></a>

### `firewall_spec.external.http`

Firewall options for HTTP workloads.

Read-Only:

- **inbound_header_filter** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec--external--http--inbound_header_filter)).

<a id="nestedblock--firewall_spec--external--http--inbound_header_filter"></a>

### `firewall_spec.external.http.inbound_header_filter`

Configures header-based allow/deny logic.

Read-Only:

- **key** (String) Header name to inspect.
- **allowed_values** (List of String) Regular expressions describing allowed header values. Requests that fail to match any allowed value are filtered.
- **blocked_values** (List of String) Regular expressions describing blocked header values. Requests that match any blocked value are filtered.

<a id="nestedblock--firewall_spec--internal"></a>

### `firewall_spec.internal`

Controls access between workloads.

Read-Only:

- **inbound_allow_type** (String) Internal firewall mode. One of `none`, `same-gvc`, `same-org`, or `workload-list`.
- **inbound_allow_workload** (List of String) Specific workloads allowed when `inbound_allow_type` is `workload-list`.

<a id="nestedblock--options"></a>

### `options`

Exports the workload-level options applied by Control Plane.

Read-Only:

- **timeout_seconds** (Number) Request timeout in seconds. Default: `5`.
- **capacity_ai** (Boolean) Whether Capacity AI is enabled. Default: `true`.
- **capacity_ai_update_minutes** (Number) Minimum interval (in minutes) between Capacity AI reservation updates.
- **debug** (Boolean) Indicates if debug mode is enabled. Default: `false`.
- **suspend** (Boolean) Indicates if the workload is suspended. Default: `false`.
- **autoscaling** (Block List, Max: 1) ([see below](#nestedblock--options--autoscaling)).
- **multi_zone** (Block List, Max: 1) ([see below](#nestedblock--options--multi_zone)).

<a id="nestedblock--options--autoscaling"></a>

### `options.autoscaling`

Read-Only:

- **metric** (String) Scaling metric. One of `concurrency`, `cpu`, `memory`, `rps`, `latency`, `keda`, or `disabled`.
- **metric_percentile** (String) Percentile to target when the metric represents a distribution (for example, latency).
- **target** (Number) Target value for the selected metric. Min: `1`. Max: `20000`. Default: `95`.
- **min_scale** (Number) Minimum replicas allowed. Min: `0`. Max: `max_scale`. Default: `1`.
- **max_scale** (Number) Maximum replicas allowed. Min: `0`. Default: `5`.
- **scale_to_zero_delay** (Number) Seconds without requests before scaling to zero. Min: `30`. Max: `3600`. Default: `300`.
- **max_concurrency** (Number) Maximum concurrent requests per replica. Min: `0`. Max: `1000`. Default: `0`.
- **multi** (Block List, Max: 1) ([see below](#nestedblock--options--autoscaling--multi)).
- **keda** (Block List, Max: 1) ([see below](#nestedblock--options--autoscaling--keda)).

<a id="nestedblock--options--autoscaling--multi"></a>

### `options.autoscaling.multi`

Read-Only:

- **metric** (String) Either `cpu` or `memory`.
- **target** (Number) Target value for the metric. Min: `1`. Max: `20000`.

<a id="nestedblock--options--autoscaling--keda"></a>

### `options.autoscaling.keda`

KEDA (Kubernetes-based Event Driven Autoscaling) configuration.

Read-Only:

- **polling_interval** (Number) Seconds between KEDA polling cycles.
- **cooldown_period** (Number) Cooldown seconds after scaling to zero before scaling up again.
- **initial_cooldown_period** (Number) Initial cooldown after scaling to zero.
- **trigger** (Block List) ([see below](#nestedblock--options--autoscaling--keda--trigger)).
- **advanced** (Block List) ([see below](#nestedblock--options--autoscaling--keda--advanced)).

<a id="nestedblock--options--autoscaling--keda--trigger"></a>

### `options.autoscaling.keda.trigger`

Defines event-driven scaling triggers.

Read-Only:

- **type** (String) Trigger type (for example, `prometheus`, `aws-sqs`).
- **metadata** (Map of String) Configuration parameters required by the trigger.
- **name** (String) Optional trigger name.
- **use_cached_metrics** (Boolean) Indicates whether metrics caching is enabled during the polling interval.
- **metric_type** (String) Metric type used for scaling.
- **authentication_ref** (Block List, Max: 1) ([see below](#nestedblock--options--autoscaling--keda--trigger--authentication_ref)).

<a id="nestedblock--options--autoscaling--keda--trigger--authentication_ref"></a>

### `options.autoscaling.keda.trigger.authentication_ref`

Read-Only:

- **name** (String) Name of the secret listed in `spec.keda.secrets` on the GVC.

<a id="nestedblock--options--autoscaling--keda--advanced"></a>

### `options.autoscaling.keda.advanced`

Advanced KEDA modifiers.

Read-Only:

- **scaling_modifiers** (Block List) ([see below](#nestedblock--options--autoscaling--keda--advanced--scaling_modifiers)).

<a id="nestedblock--options--autoscaling--keda--advanced--scaling_modifiers"></a>

### `options.autoscaling.keda.advanced.scaling_modifiers`

Read-Only:

- **target** (String) New target value for the composed metric.
- **activation_target** (String) Activation target for the composed metric.
- **metric_type** (String) Metric type used for the composed metric.
- **formula** (String) Expression that combines or transforms metrics.

<a id="nestedblock--options--multi_zone"></a>

### `options.multi_zone`

Read-Only:

- **enabled** (Boolean) Indicates if multi-zone execution is enabled.

<a id="nestedblock--local_options"></a>

### `local_options`

Overrides default options for specific Control Plane locations.

Read-Only:

- **location** (String) Location name whose options are overridden.
- All attributes from [`options`](#nestedblock--options) are repeated here with location-specific values.

<a id="nestedblock--job"></a>

### `job`

Exports cron workload settings.

Read-Only:

- **schedule** (String) Cron schedule expression determining job execution times.
- **concurrency_policy** (String) Either `Forbid` or `Replace`. Determines how overlapping jobs are handled.
- **history_limit** (Number) Maximum completed job instances retained. Integer between `1` and `10`. Default: `5`.
- **restart_policy** (String) Either `OnFailure` or `Never`. Default: `Never`.
- **active_deadline_seconds** (Number) Maximum seconds a job can run before it is forcibly stopped.

<a id="nestedblock--sidecar"></a>

### `sidecar`

Read-Only:

- **envoy** (String) Name of the Envoy sidecar configuration attached to the workload.

<a id="nestedblock--rollout_options"></a>

### `rollout_options`

Controls rolling-update behavior.

Read-Only:

- **min_ready_seconds** (Number) Minimum seconds a container must run without crashing to be considered available.
- **max_unavailable_replicas** (String) Maximum replicas that can be unavailable during an update.
- **max_surge_replicas** (String) Maximum replicas above the desired count during an update.
- **scaling_policy** (String) Update strategy. Either `OrderedReady` or `Parallel`. Default: `OrderedReady`.
- **termination_grace_period_seconds** (Number) Seconds allowed for graceful termination, including `preStop` hooks.

~> **Note** `max_surge_replicas` and `max_unavailable_replicas` accept absolute numbers (for example, `2`) or percentages (for example, `50%`), and they cannot both be zero.

<a id="nestedblock--security_options"></a>

### `security_options`

Read-Only:

- **file_system_group_id** (Number) Group ID applied to mounted volumes.

<a id="nestedblock--load_balancer"></a>

### `load_balancer`

Read-Only:

- **replica_direct** (Boolean) When `true`, individual replicas can be reached directly using `replica-<index>` subdomains. Only valid for `stateful` workloads.
- **direct** (Block List, Max: 1) ([see below](#nestedblock--load_balancer--direct)).
- **geo_location** (Block List, Max: 1) ([see below](#nestedblock--load_balancer--geo_location)).

<a id="nestedblock--load_balancer--direct"></a>

### `load_balancer.direct`

Direct load balancers are created in each workload location and expose the workload's standard endpoints. Customers must configure certificates if TLS is required.

Read-Only:

- **enabled** (Boolean) Indicates if the direct load balancer is active.
- **ipset** (String) Name of the IP set associated with the load balancer, if any.
- **port** (Block List) ([see below](#nestedblock--load_balancer--direct--port)).

<a id="nestedblock--load_balancer--direct--port"></a>

### `load_balancer.direct.port`

Read-Only:

- **external_port** (Number) Public-facing port.
- **protocol** (String) Protocol exposed publicly.
- **scheme** (String) Overrides the default `https` URL scheme in generated links.
- **container_port** (Number) Container port receiving the traffic.

<a id="nestedblock--load_balancer--geo_location"></a>

### `load_balancer.geo_location`

Read-Only:

- **enabled** (Boolean) When enabled, geo-location headers are injected into inbound HTTP requests.
- **headers** (Block List, Max: 1) ([see below](#nestedblock--load_balancer--geo_location--headers)).

<a id="nestedblock--load_balancer--geo_location--headers"></a>

### `load_balancer.geo_location.headers`

Read-Only:

- **asn** (String) ASN header value injected into requests.
- **city** (String) City header value.
- **country** (String) Country header value.
- **region** (String) Region header value.

<a id="nestedblock--request_retry_policy"></a>

### `request_retry_policy`

Read-Only:

- **attempts** (Number) Number of retry attempts. Default: `2`.
- **retry_on** (List of String) Retry conditions that trigger another attempt.

<a id="nestedblock--status"></a>

### `status`

Current state of the workload.

Read-Only:

- **parent_id** (String) ID of the parent object.
- **canonical_endpoint** (String) Canonical endpoint for the workload.
- **endpoint** (String) Public endpoint for the workload.
- **internal_name** (String) Internal hostname used for service-to-service communication.
- **replica_internal_names** (List of String)
- **health_check** (Block List) ([see below](#nestedblock--status--health_check)).
- **current_replica_count** (Number) Current number of replicas deployed.
- **resolved_images** (Block List) ([see below](#nestedblock--status--resolved_images)).
- **load_balancer** (Block List) ([see below](#nestedblock--status--load_balancer)).

<a id="nestedblock--status--health_check"></a>

### `status.health_check`

Details about the most recent health checks.

Read-Only:

- **active** (Boolean) Indicates if the health check is active.
- **success** (Boolean) Indicates if the workload is considered healthy.
- **code** (Number) Status code returned by the check.
- **message** (String) Health check message.
- **failures** (Number) Number of recent failures.
- **successes** (Number) Number of recent successes.
- **last_checked** (String) Timestamp (UTC) of the last health check.

<a id="nestedblock--status--resolved_images"></a>

### `status.resolved_images`

Resolved container images when dynamic tags are enabled.

Read-Only:

- **resolved_for_version** (Number) Workload version for which the images were resolved.
- **resolved_at** (String) UTC timestamp when resolution happened.
- **error_messages** (List of String) Errors encountered while resolving images.
- **next_retry_at** (String)
- **images** (Block List) ([see below](#nestedblock--status--resolved_images--images)).

<a id="nestedblock--status--resolved_images--images"></a>

### `status.resolved_images.images`

Read-Only:

- **digest** (String) SHA256 digest uniquely identifying the image content.
- **manifests** (Block List) ([see below](#nestedblock--status--resolved_images--images--manifests)).

<a id="nestedblock--status--resolved_images--images--manifests"></a>

### `status.resolved_images.images.manifests`

Read-Only:

- **image** (String) Name and tag of the resolved image.
- **media_type** (String) MIME type describing the manifest format.
- **digest** (String) SHA256 digest identifying the manifest.
- **platform** (Map of String) Key-value pairs describing the target OS and architecture.

<a id="nestedblock--status--load_balancer"></a>

### `status.load_balancer`

Read-Only:

- **origin** (String) Origin identifier associated with the load balancer.
- **url** (String) Load-balancer endpoint URL exposed by Control Plane.

## Example Usage

```terraform
data "cpln_workload" "workload" {
  name = "workload-example"
  gvc  = "gvc-example"
}

output "workload_id" {
  value = data.cpln_workload.workload.id
}
```
