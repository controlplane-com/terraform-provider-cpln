---
page_title: "cpln_workload Resource - terraform-provider-cpln"
subcategory: "Workload"
description: |-
---

# cpln_workload (Resource)

Manages a GVC's [Workload](https://docs.controlplane.com/reference/workload).

## Declaration

### Required

- **name** (String) Name of the Workload.
- **gvc** (String) Name of the associated GVC.
- **type** (String) Workload Type. Either `serverless`, `standard`, `stateful`, or `cron`.
- **container** (Block List) ([see below](#nestedblock--container)).
- **options** (Block List, Max: 1) ([see below](#nestedblock--options)).

### Optional

- **description** (String) Description of the Workload.
- **firewall_spec** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec)).
- **identity_link** (String) Full link to an Identity.
- **local_options** (Block List) ([see below](#nestedblock--options)).
- **tags** (Map of String) Key-value map of resource tags.
- **job** (Block List, Max: 1) ([see below](#nestedblock--job)) [Cron Job Reference Page](https://docs.controlplane.com/reference/workload#cron).
- **rollout_options** (Block List, Max: 1) ([see below](#nestedblock--rollout_options))
- **security_options** (Block List, Max: 1) ([see below](#nestedblock--security_options))
- **support_dynamic_tags** (Boolean) Workload will automatically redeploy when one of the container images is updated in the container registry. Default: false.

<a id="nestedblock--container"></a>

### `container`

~> **Note** A Workload must contain at least one container.

Required:

- **name** (String) Name of the container.

  - The following rules apply to the name of a container:
    - Cannot be: 'istio-proxy', 'queue-proxy', or 'istio-validation'.
    - Cannot start with: `cpln_`.

- **image** (String) The full image and tag path.

Optional:

- **port** (Number) The port the container exposes. Only one container is allowed to specify a port. Min: `80`. Max: `65535`. Used by `serverless` Workload type. **DEPRECATED - Use `ports`.**

- **ports** (Block List) ([see below](#nestedblock--container--ports)).

~> **Note** The ports listed below are blocked and are not allowed to be used.
Containers which attempt to use these ports will not be able to bind:
8012, 8022, 9090, 9091, 15000, 15001, 15006, 15020, 15021, 15090, 41000.

- **args** (List of String) Command line arguments passed to the container at runtime.
- **env** (Map of String) Name-Value list of environment variables.
- **command** (String) Override the entry point.
- **inherit_env** (Boolean) Enables inheritance of GVC environment variables. A variable in spec.env will override a GVC variable with the same name.
- **cpu** (String) Reserved CPU of the workload when capacityAI is disabled. Maximum CPU when CapacityAI is enabled. Default: "50m".
- **memory** (String) Reserved memory of the workload when capacityAI is disabled. Maximum memory when CapacityAI is enabled. Default: "128Mi".
- **gpu_nvidia** (Block List, Max: 1) ([see below](#nestedblock--container--gpu_nvidia))
- **liveness_probe** (Block List, Max: 1) Liveness Probe ([see below](#nestedblock--container--liveness_probe)).
- **readiness_probe** (Block List, Max: 1) Readiness Probe ([see below](#nestedblock--container--readiness_probe)).

- **metrics** (Block List, Max: 1) ([see below](#nestedblock--container--metrics)) [Reference Page](https://docs.controlplane.com/reference/workload#metrics).
- **volume** (Block List) ([see below](#nestedblock--container--volume)) [Reference Page](https://docs.controlplane.com/reference/workload#volumes).
- **working_directory** (String) Override the working directory. Must be an absolute path.
- **lifecycle** (Block List, Max: 1) LifeCycle ([see below](#nestedblock--container--lifecycle)) [Reference Page](https://docs.controlplane.com/reference/workload#lifecycle).

<a id="nestedblock--container--ports"></a>

### `container.ports`

Required:

- **protocol** (String) Protocol. Choice of: `http`, `http2`, or `grpc`.
- **number** (String) Port to expose.

<a id="nestedblock--container--gpu_nvidia"></a>

### `container.gpu_nvidia`

Required:

- **model** (String) GPU Model (i.e.: t4)
- **quantity** (Int) Number of GPUs.

<a id="nestedblock--container--liveness_probe"></a>

### `container.liveness_probe`

Optional:

- **failure_threshold** (Number) Failure Threshold. Default: 3. Min: 1. Max: 20.
- **initial_delay_seconds** (Number) Initial Delay in seconds. Default: 0. Min: 0. Max: 120.
- **period_seconds** (Number) Period Seconds. Default: 10. Min: 1. Max: 60.
- **success_threshold** (Number) Success Threshold. Default: 1. Min: 1. Max: 20.
- **timeout_seconds** (Number) Timeout in seconds. Default: 1. Min: 1. Max: 60.

- **exec** (Block List, Max: 1) ([see below](#nestedblock--container--liveness_probe--exec)).
- **http_get** (Block List, Max: 1) ([see below](#nestedblock--container--liveness_probe--http_get)).
- **tcp_socket** (Block List, Max: 1) ([see below](#nestedblock--container--liveness_probe--tcp_socket)).

<a id="nestedblock--container--liveness_probe--exec"></a>

### `container.liveness_probe.exec`

Required:

- **command** (List of Strings, Min: 1) List of commands to execute.

<a id="nestedblock--container--liveness_probe--http_get"></a>

### `container.liveness_probe.http_get`

Optional:

- **path** (String) Path. Default: "/".
- **port** (Number) Port. Min: `80`. Max: `65535`.
- **http_headers** (Map of String) Name-Value list of HTTP Headers to send to container.
- **scheme** (String) HTTP Scheme. Valid values: "HTTP", "HTTPS". Default: "HTTP".

<a id="nestedblock--container--liveness_probe--tcp_socket"></a>

### `container.liveness_probe.tcp_socket`

Optional:

- **port** (Number) TCP Socket Port.

<a id="nestedblock--container--readiness_probe"></a>

### `container.readiness_probe`

Optional:

- **failure_threshold** (Number) Failure Threshold. Default: 3. Min: 1. Max: 20.
- **initial_delay_seconds** (Number) Initial Delay in seconds. Default: 0. Min: 0. Max: 120.
- **period_seconds** (Number) Period Seconds. Default: 10. Min: 1. Max: 60.
- **success_threshold** (Number) Success Threshold. Default: 1. Min: 1. Max: 20.
- **timeout_seconds** (Number) Timeout in seconds. Default: 1. Min: 1. Max: 60.

- **exec** (Block List, Max: 1) ([see below](#nestedblock--container--readiness_probe--exec)).
- **http_get** (Block List, Max: 1) ([see below](#nestedblock--container--readiness_probe--http_get)).
- **tcp_socket** (Block List, Max: 1) ([see below](#nestedblock--container--readiness_probe--tcp_socket)).

<a id="nestedblock--container--readiness_probe--exec"></a>

### `container.readiness_probe.exec`

Required:

- **command** (List of Strings, Min: 1) List of commands to execute.

<a id="nestedblock--container--readiness_probe--http_get"></a>

### `container.readiness_probe.http_get`

Optional:

- **http_headers** (Map of String) Name-Value list of HTTP Headers to send to container.
- **path** (String) Path. Default: "/".
- **port** (Number) Port. Min: `80`. Max: `65535`.
- **scheme** (String) HTTP Scheme. Valid values: "HTTP", "HTTPS". Default: "HTTP".

<a id="nestedblock--container--readiness_probe--tcp_socket"></a>

### `container.readiness_probe.tcp_socket`

Optional:

- **port** (Number) TCP Socket Port.

<a id="nestedblock--container--volume"></a>

### `container.volume`

Required:

- **uri** (String) URI of volume at cloud provider.
- **path** (String) File path added to workload pointing to the volume.

~> **Note** The following list of paths are reserved and cannot be used: `/dev`, `/dev/log`, `/tmp`, `/var`, `/var/log`.

<a id="nestedblock--container--metrics"></a>

### `container.metrics`

Required:

- **path** (String) Path from container emitting custom metrics
- **port** (Number) Port from container emitting custom metrics

<a id="nestedblock--container--lifecycle"></a>

### `container.lifecycle`

Optional:

- **post_start** (Block List, Max: 1) ([see below](#nestedblock--container--lifecycle--spec)).
- **pre_stop** (Block List, Max: 1) ([see below](#nestedblock--container--lifecycle--spec)).

<a id="nestedblock--container--lifecycle--spec"></a>

### `container.lifecycle.spec`

Optional:

- **exec** (Block List, Max: 1) ([see below](#nestedblock--container--lifecycle--spec--exec)).

<a id="nestedblock--container--lifecycle--spec--exec"></a>

### `container.lifecycle.spec.exec`

Required:

- **command** (List of Strings, Min: 1) List of commands to execute.

<a id="nestedblock--firewall_spec"></a>

### `firewall_spec`

Control of inbound and outbound access to the workload for external (public) and internal (service to service) traffic. Access is restricted by default.

Optional:

- **external** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec--external)).
- **internal** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec--internal)).

<a id="nestedblock--firewall_spec--external"></a>

### `firewall_spec.external`

Optional:

- **inbound_allow_cidr** (List of String) he list of ipv4/ipv6 addresses or cidr blocks that are allowed to access this workload. No external access is allowed by default. Specify '0.0.0.0/0' to allow access to the public internet.
- **outbound_allow_hostname** (List of String) The list of public hostnames that this workload is allowed to reach. No outbound access is allowed by default. A wildcard `*` is allowed on the prefix of the hostname only, ex: `*.amazonaws.com`. Use `outboundAllowCIDR` to allow access to all external websites.
- **outbound_allow_cidr** (List of String) The list of ipv4/ipv6 addresses or cidr blocks that this workload is allowed reach. No outbound access is allowed by default. Specify '0.0.0.0/0' to allow outbound access to the public internet.
- **outbound_allow_port** (Block List) ([see below](#nestedblock--firewall_spec--external--outbound_allow_port)).

<a id="nestedblock--firewall_spec--external--outbound_allow_port"></a>

### `firewall_spec.external.outbound_allow_port`

Allow outbound access to specific ports and protocols. When not specified, communication to address ranges in outboundAllowCIDR is allowed on all ports and communication to names in outboundAllowHostname is allowed on ports 80/443.

Required:

- **protocol** (String) Either `http`, `https` or `tcp`. Default: `tcp`.
- **number** (Number) Port number. Max: 65000

<a id="nestedblock--firewall_spec--internal"></a>

### `firewall_spec.internal`

The internal firewall is used to control access between workloads.

Optional:

- **inbound_allow_type** (String) Used to control the internal firewall configuration and mutual tls. Allowed Values: "none", "same-gvc", "same-org", "workload-list".

  - 'none': no access is allowed between this workload and other workloads on Control Plane.
  - 'same-gvc': workloads running on the same Global Virtual Cloud are allowed to access this workload internally.
  - 'same-org': workloads running on the same Control Plane Organization are allowed to access this workload internally.
  - 'workload-list': specific workloads provided in the 'inboundAllowWorkload' array are allowed to access this workload internally.

- **inbound_allow_workload** (List of String) A list of specific workloads which are allowed to access this workload internally. This list is only used if the 'inboundAllowType' is set to 'workload-list'.

<a id="nestedblock--options"></a>

### `options`

Optional:

- **autoscaling** (Block List, Max: 1) ([see below](#nestedblock--options--autoscaling)).
- **timeout_seconds** (Number) Timeout in seconds. Default: `5`.
- **capacity_ai** (Boolean) Capacity AI. Default: `true`.
- **debug** (Boolean) Debug mode. Default: `false`
- **suspend** (Boolean) Workload suspend. Default: `false`

- **location** (String) Valid only for `local_options`. Override options for a specific location.

<a id="nestedblock--options--autoscaling"></a>

### `options.autoscaling`

Optional:

- **metric** (String) Valid values: `concurrency`, `cpu`, `latency`, or `rps`.
- **target** (Number) Control Plane will scale the number of replicas for this deployment up/down in order to be as close as possible to the target metric across all replicas of a deployment. Min: `0`. Max: `20000`. Default: `95`.
- **max_scale** (Number) The maximum allowed number of replicas. Min: `0`. Default `5`.
- **min_scale** (Number) The minimum allowed number of replicas. Control Plane can scale the workload down to 0 when there is no traffic and scale up immediately to fulfill new requests. Min: `0`. Max: `max_scale`. Default `1`.
- **scale_to_zero_delay** (Number) The amount of time (in seconds) with no requests received before a workload is scaled to 0. Min: `30`. Max: `3600`. Default: `300`.
- **max_concurrency** (Number) A hard maximum for the number of concurrent requests allowed to a replica. If no replicas are available to fulfill the request then it will be queued until a replica with capacity is available and delivered as soon as one is available again. Capacity can be available from requests completing or when a new replica is available from scale out.Min: `0`. Max: `1000`. Default `0`.

<a id="nestedblock--job"></a>

### `job`

~> **Note** A CRON workload must contain a `job`.<br/><br/>Capacity AI must be false and min/max scale must equal 1.

Required:

- **schedule** (String) A standard cron [schedule expression](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#schedule-syntax) used to determine when your job should execute.

Optional:

- **concurrency_policy** (String) Either 'Forbid' or 'Replace'. This determines what Control Plane will do when the schedule requires a job to start, while a prior instance of the job is still running. Enum: [ Forbid, Replace ] Default: `Forbid`
- **history_limit** (Number) The maximum number of completed job instances to display. This should be an integer between 1 and 10. Default: `5`
- **restart_policy** (String) Either 'OnFailure' or 'Never'. This determines what Control Plane will do when a job instance fails. Enum: [ OnFailure, Never ] Default: `Never`
- **active_deadline_seconds** (Number) The maximum number of seconds Control Plane will wait for the job to complete. If a job does not succeed or fail in the allotted time, Control Plane will stop the job, moving it into the Removed status.

<a id="nestedblock--rollout_options"></a>

### `rollout_options`

Optional:

- **min_ready_seconds** (Number) The minimum number of seconds a container must run without crashing to be considered available.
- **max_unavailable_replicas** (String) The number of replicas that can be unavailable during the update process.
- **max_surge_replicas** (String) The number of replicas that can be created above the desired amount of replicas during an update.

~> **Note** Both max_surge_replicas and max_unavailable_replicas can be specified as either an integer (e.g. 2) or a percentage (e.g. 50%), and they cannot both be zero.

<a id="nestedblock--security_options"></a>

### `security_options`

Required:

- **file_system_group_id** (Number) The group id assigned to any mounted volume

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Workload.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (List of Object) ([see below](#nestedatt--status)).

<a id="nestedatt--status"></a>

### `status`

Status of the workload.

Read-Only:

- **canonical_endpoint** (String) Canonical endpoint for the workload.
- **endpoint** (String) Endpoint for the workload.
- **health_check** (List of Object) ([see below](#nestedobjatt--status--health_check)).
- **parent_id** (String) ID of the parent object.

<a id="nestedobjatt--status--health_check"></a>

### `status.health_check`

Current health status.

Read-Only:

- **active** (Boolean) Active boolean for the associated workload.
- **code** (Number) Current output code for the associated workload.
- **failures** (Number) Failure integer for the associated workload.
- **last_checked** (String) Timestamp in UTC of the last health check.
- **message** (String) Current health status for the associated workload.
- **success** (Boolean) Success boolean for the associated workload.
- **successes** (Number) Success integer for the associated workload.

## Example Usage - Serverless

```terraform
resource "cpln_gvc" "example" {
  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_identity" "example" {

  gvc = cpln_gvc.example.name

  name        = "identity-example"
  description = "Example Identity"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_workload" "new" {

  gvc = cpln_gvc.example.name

  type = "serverless"

  name        = "workload-example"
  description = "Example Workload"

  support_dynamic_tags = false

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_link = cpln_identity.example.self_link

  container {
    name   = "container-01"
    image  = "gcr.io/knative-samples/helloworld-go"

    memory = "128Mi"
    cpu    = "50m"

    ports {
			protocol = "http"
			number   = "8080"
		}

    command = "override-command"
    working_directory = "/usr"

    inherit_env = false

    env = {
      env-name-01 = "env-value-01",
      env-name-02 = "env-value-02",
    }

    args = ["arg-01", "arg-02"]

    readiness_probe {

      tcp_socket {
        port = 8181
      }

      period_seconds        = 11
      timeout_seconds       = 2
      failure_threshold     = 4
      success_threshold     = 2
      initial_delay_seconds = 1
    }

    liveness_probe {

      http_get {
        path   = "/path"
        port   = 8282
        scheme = "HTTPS"
        http_headers = {
          header-name-01 = "header-value-01"
          header-name-02 = "header-value-02"
        }
      }

      period_seconds        = 10
      timeout_seconds       = 3
      failure_threshold     = 5
      success_threshold     = 1
      initial_delay_seconds = 2
    }

    lifecycle {

      post_start {
        exec {
          command = ["command_post", "arg_1", "arg_2"]
        }
      }

      pre_stop {
        exec {
          command = ["command_pre", "arg_1", "arg_2"]
        }
      }
    }

    volume {
      uri  = "s3://bucket"
      path = "/s3"
    }
  }

  options {
    capacity_ai     = false
    timeout_seconds = 30
    suspend         = false

    autoscaling {
      metric          = "concurrency"
      target          = 100
      max_scale       = 3
      min_scale       = 2
      max_concurrency = 500
    }
  }

  local_options {

    location        = "aws-us-west-2"
    capacity_ai     = false
    timeout_seconds = 30
    suspend         = false

    autoscaling {
      metric          = "concurrency"
      target          = 100
      max_scale       = 3
      min_scale       = 2
      max_concurrency = 500
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      outbound_allow_cidr     = []
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]

      outbound_allow_port {
				protocol = "http"
				number   = 80
			}

			outbound_allow_port {
				protocol = "https"
				number   = 443
			}
    }
    internal {
      # Allowed Types: "none", "same-gvc", "same-org", "workload-list"
      inbound_allow_type     = "none"
      inbound_allow_workload = []
    }
  }

  security_options {
    file_system_group_id = 1
  }

}
```

## Example Usage - Standard

```terraform

resource "cpln_gvc" "example" {
  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_identity" "example" {

  gvc = cpln_gvc.example.name

  name        = "identity-example"
  description = "Example Identity"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_workload" "new" {

  gvc = cpln_gvc.example.name

  type = "standard"

  name        = "workload-example"
  description = "Example Workload"

  support_dynamic_tags = false

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_link = cpln_identity.example.self_link

  container {
    name   = "container-01"
    image  = "gcr.io/knative-samples/helloworld-go"
    memory = "128Mi"
    cpu    = "50m"

		ports {
		  protocol = "http"
			number   = "80"
		}

		ports {
			protocol = "http2"
			number   = "8080"
	  }

    env = {
      env-name-01 = "env-value-01",
      env-name-02 = "env-value-02",
    }

    lifecycle {

      post_start {
        exec {
          command = ["command_post", "arg_1", "arg_2"]
        }
      }

      pre_stop {
        exec {
          command = ["command_pre", "arg_1", "arg_2"]
        }
      }
    }
  }

  options {
    capacity_ai     = false
    timeout_seconds = 30
    suspend         = false

    autoscaling {
      metric          = "cpu"
      target          = 60
      max_scale       = 3
      min_scale       = 2
      max_concurrency = 500
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      outbound_allow_cidr     = []
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]

      outbound_allow_port {
				protocol = "http"
				number   = 80
			}

			outbound_allow_port {
				protocol = "https"
				number   = 443
			}
    }
    internal {
      # Allowed Types: "none", "same-gvc", "same-org", "workload-list"
      inbound_allow_type     = "none"
      inbound_allow_workload = []
    }
  }

  rollout_options {
    min_ready_seconds = 2
    max_unavailable_replicas = "10"
    max_surge_replicas = "20"
  }

  security_options {
    file_system_group_id = 1
  }

}

```

## Example Usage - Cron

```terraform
resource "cpln_gvc" "example" {

  name        = "gvc-example"
  description = "Example GVC"
  locations   = ["aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_identity" "example" {

  gvc         = cpln_gvc.example.name
  name        = "identity-example"
  description = "Example Identity"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_workload" "new" {

  gvc         = cpln_gvc.example.name

  type = "cron"

  name        = "workload-example"
  description = "Example Workload"

  support_dynamic_tags = false

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_link = cpln_identity.example.self_link

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    memory            = "128Mi"
    cpu               = "50m"
    command           = "override-command"
    working_directory = "/usr"

    env = {
      env-name-01 = "env-value-01",
      env-name-02 = "env-value-02",
    }

    args = ["arg-01", "arg-02"]

    volume {
      uri  = "s3://bucket"
      path = "/testpath01"
    }

    volume {
      uri  = "azureblob://storageAccount/container"
      path = "/testpath02"
    }

    metrics {
      path = "/metrics"
      port = 8181
    }
  }

  options {
    suspend     = false
    capacity_ai = false

    autoscaling {
    	max_scale = 1
	    min_scale = 1
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]

      outbound_allow_port {
				protocol = "http"
				number   = 80
			}

			outbound_allow_port {
				protocol = "https"
				number   = 443
			}
    }
  }

  security_options {
    file_system_group_id = 1
  }

  job {
    schedule                = "* * * * *"
    concurrency_policy      = "Forbid"
    history_limit           = 5
    restart_policy          = "Never"
    active_deadline_seconds = 1200
  }
}

```

## Example Usage - Serverless Workload with a GPU resource

```terraform

resource "cpln_gvc" "new" {
  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-us-west-2", "gcp-us-east1"]

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}

resource "cpln_identity" "new" {

  gvc = cpln_gvc.new.name

  name        = "identity-example"
  description = "Identity created using terraform"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}


resource "cpln_workload" "new" {

  gvc = cpln_gvc.new.name

  name        = "workload-example"
  description = "Example Workload"
  type        = "serverless"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true

  container {
    name   = "container-01"
    image  = "gcr.io/knative-samples/helloworld-go"

    memory = "7Gi"
    cpu    = "2"

    ports {
			protocol = "http"
			number   = "8080"
		}

    gpu_nvidia {
      model    = "t4"
      quantity = 1
    }

    command           = "override-command"
    working_directory = "/usr"

    env = {
      env-name-01 = "env-value-01",
      env-name-02 = "env-value-02",
    }

    args = ["arg-01", "arg-02"]

    volume {
      uri  = "s3://bucket"
      path = "/testpath01"
    }

    volume {
      uri  = "azureblob://storageAccount/container"
      path = "/testpath02"
    }

    metrics {
      path = "/metrics"
      port = 8181
    }

    readiness_probe {

      tcp_socket {
        port = 8181
      }

      period_seconds        = 11
      timeout_seconds       = 2
      failure_threshold     = 4
      success_threshold     = 2
      initial_delay_seconds = 1
    }

    liveness_probe {

      http_get {
        path   = "/path"
        port   = 8282
        scheme = "HTTPS"
        http_headers = {
          header-name-01 = "header-value-01"
          header-name-02 = "header-value-02"
        }
      }

      period_seconds        = 10
      timeout_seconds       = 3
      failure_threshold     = 5
      success_threshold     = 1
      initial_delay_seconds = 2
    }

    lifecycle {

      post_start {
        exec {
          command = ["command_post", "arg_1", "arg_2"]
        }
      }

      pre_stop {
        exec {
          command = ["command_pre", "arg_1", "arg_2"]
        }
      }
    }
  }

  options {
    capacity_ai     = false
    timeout_seconds = 30
    suspend         = false

    autoscaling {
      metric              = "concurrency"
      target              = 100
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      outbound_allow_cidr     = []
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
    }

    internal {
      # Allowed Types: "none", "same-gvc", "same-org", "workload-list"
      inbound_allow_type     = "none"
      inbound_allow_workload = []
    }
  }

  security_options {
    file_system_group_id = 1
  }


}

```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing workload resource, execute the following import command:

```terraform
terraform import cpln_workload.RESOURCE_NAME GVC_NAME:WORKLOAD_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute GVC_NAME and WORKLOAD_NAME with the corresponding GVC and workload name defined in the resource.
