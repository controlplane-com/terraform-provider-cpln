---
page_title: "cpln_gvc Resource - terraform-provider-cpln"
subcategory: "Global Virtual Cloud"
description: |-
---

# cpln_gvc (Resource)

Manages an org's [Global Virtual Cloud (GVC)](https://docs.controlplane.com/reference/gvc).

## Declaration

### Required

- **name** (String) Name of the GVC.

### Optional

- **description** (String) Description of the GVC.
- **tags** (Map of String) Key-value map of resource tags.
- **locations** (List of String) A list of [locations](https://docs.controlplane.com/reference/location#current) making up the Global Virtual Cloud.
- **pull_secrets** (List of String) A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.
- **domain** (String) Custom domain name used by associated workloads.
- **endpoint_naming_format** (String) Customizes the subdomain format for the canonical workload endpoint. `legacy` leaves it as '${workloadName}-${gvcName}.cpln.app'. `org` follows the scheme '${workloadName}-${gvcName}.${orgEndpointPrefix}.cpln.app'.
- **alias_workload_link** (String) A link to a workload in this GVC whose canonical endpoint backs the GVC alias DNS record. When set, the GVC alias is published as a CNAME to the workload's canonical endpoint, inheriting its HTTP health probes and per-location geo failover. When unset, the alias resolves directly to cluster ingress endpoints with no application-level health awareness. Has no effect while the referenced workload is globally suspended.
- **env** (Array of Name-Value Pair) Key-value array of resource env variables.
- **lightstep_tracing** (Block List, Max: 1) ([see below](#nestedblock--lightstep_tracing)).
- **otel_tracing** (Block List, Max: 1) ([see below](#nestedblock--otel_tracing)).
- **controlplane_tracing** (Block List, Max: 1) ([see below](#nestedblock--controlplane_tracing)).
- **load_balancer** (Block List, Max: 1) ([see below](#nestedblock--load_balancer)).
- **keda** (Block List, Max: 1) ([see below](#nestedblock--keda)).
- **location_query** (Block List, Max: 1) ([see below](#nestedblock--location_query)).
- **location_options** (Block List) ([see below](#nestedblock--location_options)).

~> **Note** Only one of the tracing blocks can be defined.

<a id="nestedblock--lightstep_tracing"></a>

### `lightstep_tracing`

Required:

- **sampling** (Int) Determines what percentage of requests should be traced.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or internal endpoint.

Optional:

- **credentials** (String) Full link to referenced Opaque Secret.
- **custom_tags** (Map of String) Key-value map of custom tags.

<a id="nestedblock--otel_tracing"></a>

### `otel_tracing`

Required:

- **sampling** (Int) Determines what percentage of requests should be traced.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or internal endpoint.

Optional:

- **custom_tags** (Map of String) Key-value map of custom tags.

<a id="nestedblock--controlplane_tracing"></a>

### `controlplane_tracing`

Required:

- **sampling** (Int) Determines what percentage of requests should be traced.

Optional:

- **custom_tags** (Map of String) Key-value map of custom tags.

<a id="nestedblock--load_balancer"></a>

### `load_balancer`

Optional:

- **dedicated** (Boolean) Creates a dedicated load balancer in each location and enables additional Domain features: custom ports, protocols and wildcard hostnames. Charges apply for each location.
- **multi_zone** (Block List, Max: 1) ([see below](#nestedblock--load_balancer--multi_zone)).
- **trusted_proxies** (Int) Controls the address used for request logging and for setting the X-Envoy-External-Address header. If set to 1, then the last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If set to 2, then the second to last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If the XFF header does not have at least two addresses or does not exist then the source client IP address will be used instead.
- **redirect** (Block List, Max: 1) ([see below](#nestedblock--load_balancer--redirect)).
- **ipset** (String) The link or the name of the IP Set that will be used for this load balancer.

<a id="nestedblock--load_balancer--multi_zone"></a>

### `load_balancer.multi_zone`

- **enabled** (Boolean) Default: `false`.

<a id="nestedblock--load_balancer--redirect"></a>

### `load_balancer.redirect`

Specify the url to be redirected to for different http status codes.

Optional:

- **class** (Block List, Max: 1) ([see below](#nestedblock--load_balancer--redirect--class)).

<a id="nestedblock--load_balancer--redirect--class"></a>

### `load_balancer.redirect.class`

Specify the redirect url for all status codes in a class.

Optional:

- **status_5xx** (String) Specify the redirect url for any 500 level status code.
- **status_401** (String) An optional url redirect for 401 responses. Supports envoy format strings to include request information. E.g. https://your-oauth-server/oauth2/authorize?return_to=%REQ(:path)%&client_id=your-client-id

<a id="nestedblock--keda"></a>

### `keda`

- **enabled** (Boolean) Enable KEDA for this GVC. KEDA is a Kubernetes-based event-driven autoscaler that allows you to scale workloads based on external events. When enabled, a keda operator will be deployed in the GVC and workloads in the GVC can use KEDA to scale based on external metrics.
- **identity_link** (String) A link to an Identity resource that will be used for KEDA. This will allow the keda operator to access cloud and network resources.
- **secrets** (List of String) A list of secrets to be used as TriggerAuthentication objects. The TriggerAuthentication object will be named after the secret and can be used by triggers on workloads in this GVC.

<a id="nestedblock--location_query"></a>

### `location_query`

A query that dynamically selects the locations making up the Global Virtual Cloud.

Optional:

- **fetch** (String) Type of fetch. Specify either: `links` or `items`. Default: `items`.
- **spec** (Block List, Max: 1) ([see below](#nestedblock--location_query--spec)).

<a id="nestedblock--location_query--spec"></a>

### `location_query.spec`

Optional:

- **match** (String) Type of match. Available values: `all`, `any`, `none`. Default: `all`.
- **terms** (Block List) ([see below](#nestedblock--location_query--spec--terms)).

<a id="nestedblock--location_query--spec--terms"></a>

### `location_query.spec.terms`

Terms can only contain one of the following attributes: `property`, `rel`, `tag`.

Optional:

- **op** (String) Type of query operation. Available values: `=`, `>`, `>=`, `<`, `<=`, `!=`, `~`, `=~`, `exists`, `!exists`, `contains`. Default: `=`.
- **property** (String) Property to use for query evaluation.
- **rel** (String) Relation to use for query evaluation.
- **tag** (String) Tag key to use for query evaluation.
- **value** (String) Testing value for query evaluation.

<a id="nestedblock--location_options"></a>

### `location_options`

Per-location routing options for DNS geo routing. Allows configuring priority-based failover and latency adjustments per location. Each entry references a location listed in `locations`.

Required:

- **name** (String) Name of the location these options apply to.

Optional:

- **routing_tier** (Number) Routing tier for DNS geo routing. Lower value = higher priority. Locations with the same `routing_tier` form a group; within a group, lowest latency wins. If all locations in the highest-priority group are unavailable, the next group is used.
- **latency_offset_ms** (Number) Artificial latency offset in milliseconds added to measured latency. Positive values push traffic away from this location, negative values attract traffic. Default: `0`.
- **latency_tolerance_ms** (Number) Maximum acceptable latency in milliseconds. If measured latency exceeds this value, the location is treated as unavailable for DNS geo routing.

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the GVC.
- **alias** (String) The alias name of the GVC.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform

resource "cpln_secret" "docker" {
  name        = "docker-secret"
  description = "docker secret"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "docker"
  }

  docker = "{\"auths\":{\"your-registry-server\":{\"username\":\"your-name\",\"password\":\"your-pword\",\"email\":\"your-email\",\"auth\":\"<Secret>\"}}}"
}

resource "cpln_secret" "opaque" {

  name        = "opaque-random-tbd"
  description = "description opaque-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_gvc" "example" {

  name        = "gvc-example"
  description = "Example GVC"

  # Org endpoint naming format gives us: ${workloadName}-${gvcName}.${orgEndpointPrefix}.cpln.app
  endpoint_naming_format = "org"

  # Publishes the GVC alias as a CNAME to the workload's canonical endpoint for health-aware DNS.
  # alias_workload_link = "/org/terraform-test-org/gvc/gvc-example/workload/my-workload"

  # Example Locations: `aws-eu-central-1`, `aws-us-west-2`, `azure-east2`, `gcp-us-east1`
  locations = ["aws-eu-central-1", "aws-us-west-2"]

  # Per-location DNS geo routing options.
  # Locations in the highest-priority tier (lowest routing_tier) are tried first;
  # within a tier, the location with the lowest measured latency wins.
  location_options {
    name                 = "aws-eu-central-1"
    routing_tier         = 1
    latency_tolerance_ms = 150
  }

  location_options {
    name         = "aws-us-west-2"
    routing_tier = 2
  }

  # As an alternative to the explicit `locations` list above, a `location_query`
  # can dynamically select the locations making up the GVC.
  # location_query {
  #   fetch = "links"
  #
  #   spec {
  #     match = "all"
  #
  #     terms {
  #       op       = "="
  #       property = "provider"
  #       value    = "aws"
  #     }
  #   }
  # }

  # domain = "app.example.com"
  pull_secrets = [cpln_secret.docker.name]

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  env = {
    env_var_key          = "env_var_value"
    workload_can_inherit = "true"
  }

  lightstep_tracing {

    sampling = 50
    endpoint = "test.cpln.local:8080"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link
  }

  load_balancer {
    dedicated       = true
    trusted_proxies = 1
    ipset           = "my-ipset"

    multi_zone {
      enabled = false
    }

    redirect {
      class {
        status_5xx = "https://example.com/error/5xx"
        status_401 = "https://your-oauth-server/oauth2/authorize?return_to=%REQ(:path)%&client_id=your-client-id"
      }
    }
  }

  keda {
    enabled       = true
    identity_link = "/org/terraform-test-org/gvc/gvc-example/identity/non-existant-identity"
    secrets       = ["/org/terraform-test-org/secret/my-secret-01", "/org/terraform-test-org/secret/my-secret-02", "/org/terraform-test-org/secret/my-secret-03"]
  }
}

```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing GVC resource, execute the following import command:

```terraform
terraform import cpln_gvc.RESOURCE_NAME GVC_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute GVC_NAME with the corresponding GVC defined in the resource.
