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


### Optional

- **description** (String) Description of the Workload.

- **type** (String) Workload Type. Either `serverless` or `standard`. Default: `serverless`. 

- **container** (Block List) ([see below](#nestedblock--container)).
- **firewall_spec** (Block List, Max: 1) ([see below](#nestedblock--firewall_spec)).
- **identity_link** (String) Full link to an Identity.
- **options** (Block List, Max: 1) ([see below](#nestedblock--options)).
- **local_options** (Block List, Max: 1) ([see below](#nestedblock--options)).
- **tags** (Map of String) Key-value map of resource tags.


<a id="nestedblock--container"></a>
 ### `container`

~> **Note** A Workload must contain at least one container.

Required:

- **name** (String) Name of the container.
  - The following rules apply to the name of a container:
    - Cannot be: 'istio-proxy', 'queue-proxy', 'istio-validation'.
    - Cannot start with: `cpln_`.

- **image** (String) The full image and tag path.


Optional:

- **port** (Number) The port the container exposes. Only one container is allowed to specify a port. Min: `80`. Max: `65535`. Used by `serverless` Workload type.

- **ports** (Block List) ([see below](#nestedblock--container--ports)). 

~> **Note**  The ports listed below are blocked and are not allowed to be used.
Containers which attempt to use these ports will not be able to bind:
8012, 8022, 9090, 9091, 15000, 15001, 15006, 15020, 15021, 15090, 41000.


- **args** (List of String) Command line arguments passed to the container at runtime.
- **env** (Map of String) Name-Value list of environment variables.
- **command** (String) Override the entry point. 
- **cpu** (String) Reserved CPU of the workload when capacityAI is disabled. Maximum CPU when CapacityAI is enabled. Default: "50m".
- **memory** (String) Reserved memory of the workload when capacityAI is disabled. Maximum memory when CapacityAI is enabled. Default: "128Mi".
  
- **liveness_probe** (Block List, Max: 1) Liveness Probe  ([see below](#nestedblock--container--liveness_probe)).
- **readiness_probe** (Block List, Max: 1) Readiness Probe ([see below](#nestedblock--container--readiness_probe)).

- **metrics** (Block List, Max: 1) ([see below](#nestedblock--container--metrics)) [Reference Page](https://docs.controlplane.com/reference/workload#metrics).
  
- **volume** (Block List) ([see below](#nestedblock--container--volume)) [Reference Page](https://docs.controlplane.com/reference/workload#volumes).
- **working_directory** (String) Override the working directory. Must be an absolute path.

<a id="nestedblock--container--ports"></a>
 ### `container.ports`

Required:

- **protocol** (String) Protocol. Choice of: `http`, `http2`, or `grpc`.
- **number** (String) Port to expose.

<a id="nestedblock--container--ports"></a>


<a id="nestedblock--container--liveness_probe"></a>
 ### `container.liveness_probe`

Optional:

- **failure_threshold** (Number) Failure Threshold.  Default: 3. Min: 1. Max: 20.
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

- **http_headers** (Map of String) Name-Value list of HTTP Headers to send to container.
- **path** (String) Path. Default: "/".
- **port** (Number) Port. Min: `80`. Max: `65535`.
- **scheme** (String) HTTP Scheme. Valid values: "HTTP", "HTTPS". Default: "HTTP".


<a id="nestedblock--container--liveness_probe--tcp_socket"></a>
 ### `container.liveness_probe.tcp_socket`

Optional:

- **port** (Number) TCP Socket Port.



<a id="nestedblock--container--readiness_probe"></a>
 ### `container.readiness_probe`

Optional:

- **failure_threshold** (Number) Failure Threshold.  Default: 3. Min: 1. Max: 20.
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

<a id="nestedblock--container--volume"></a>
 ### `container.metrics`

Required:

- **path** (String) Path from container emitting custom metrics
- **port** (Number) Port from container emitting custom metrics



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
- **outbound_allow_cidr** (List of String) The list of ipv4/ipv6 addresses or cidr blocks that this workload is allowed reach. No outbound access is allowed by default. Specify '0.0.0.0/0' to allow outbound access to the public internet.
- **outbound_allow_hostname** (List of String) The list of public hostnames that this workload is allowed to reach. No outbound access is allowed by default. A wildcard `*` is allowed on the prefix of the hostname only, ex: `*.amazonaws.com`. Use `outboundAllowCIDR` to allow access to all external websites.

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
- **capacity_ai** (Boolean) Capacity AI. Default: `true`.
- **debug** (Boolean) Debug mode. Default: `false`
- **timeout_seconds** (Number) Timeout in seconds. Default: `5`.

- **location** (String) Valid only for `local_options`. Local options override for a specific location.

<a id="nestedblock--options--autoscaling"></a>
 ### `options.autoscaling`

Optional:

- **metric** (String) Valid values: `concurrency`, `cpu`, `rps`. Default: `concurrency`.
  
- **max_concurrency** (Number) A hard maximum for the number of concurrent requests allowed to a replica. If no replicas are available to fulfill the request then it will be queued until a replica with capacity is available and delivered as soon as one is available again. Capacity can be available from requests completing or when a new replica is available from scale out.Min: `0`. Max: `1000`. Default `0`.
- **max_scale** (Number) The maximum allowed number of replicas. Min: `0`. Default `5`.
- **min_scale** (Number) The minimum allowed number of replicas. Control Plane can scale the workload down to 0 when there is no traffic and scale up immediately to fulfill new requests. Min: `0`. Max: `max_scale`. Default `1`.
- **scale_to_zero_delay** (Number) The amount of time (in seconds) with no requests received before a workload is scaled to 0. Min: `30`. Max: `3600`. Default: `300`.
- **target** (Number) Control Plane will scale the number of replicas for this deployment up/down in order to be as close as possible to the target metric across all replicas of a deployment. Min: `0`. Max: `20000`. Default: `100`.


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

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Workload.
- **self_link** (String) Full link to this resource. Can be referenced by other resources. 
- **status** (List of Object) ([see below](#nestedatt--status)).


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

  name        = "workload-example"
  description = "Example Workload"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_link = cpln_identity.example.self_link

  type = "serverless" 

  container {
    name   = "container-01"
    image  = "gcr.io/knative-samples/helloworld-go"
    port   = 8080
    memory = "128Mi"
    cpu    = "50m"

    command = "override-command"
		working_directory = "/usr"

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

    volume {
      uri  = "s3://bucket"
      path = "/s3"
    }
  }
 
  options {
    capacity_ai     = false
    timeout_seconds = 30

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
    }
    internal {
      # Allowed Types: "none", "same-gvc", "same-org", "workload-list"
      inbound_allow_type     = "none"
      inbound_allow_workload = []
    }
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

  name        = "workload-example"
  description = "Example Workload"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_link = cpln_identity.example.self_link

  type = "standard" 

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
  }
 
  options {
    capacity_ai     = false
    timeout_seconds = 30

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
    }
    internal {
      # Allowed Types: "none", "same-gvc", "same-org", "workload-list"
      inbound_allow_type     = "none"
      inbound_allow_workload = []
    }
  }
}

```