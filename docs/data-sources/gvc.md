---
page_title: "cpln_gvc Data Source - terraform-provider-cpln"
subcategory: "Global Virtual Cloud"
description: |-
  
---
# cpln_gvc (Data Source)

Use this data source to access information about an existing [Global Virtual Cloud (GVC)](https://docs.controlplane.com/reference/gvc) within Control Plane. 

## Required

- **name** (String) Name of the GVC.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the GVC.
- **name** (String) Name of the GVC.
- **alias** (String) The alias name of the GVC.
- **description** (String) Description of the GVC.
- **tags** (Map of String) Key-value map of resource tags.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **domain** (String) Custom domain name used by associated workloads.
- **locations** (List of String) A list of [locations](https://docs.controlplane.com/reference/location#current) making up the Global Virtual Cloud.
- **pull_secrets** (List of String) A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.
- **lightstep_tracing** (Block List, Max: 1) ([see below](#nestedblock--lightstep_tracing)).
- **otel_tracing** (Block List, Max: 1) ([see below](#nestedblock--otel_tracing)).
- **controlplane_tracing** (Block List, Max: 1) ([see below](#nestedblock--controlplane_tracing)).
- **load_balancer** (Block List, Max: 1) ([see below](#nestedblock--load_balancer)).

<a id="nestedblock--lightstep_tracing"></a>

### `lightstep_tracing`

- **sampling** (Int) Sampling percentage.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internal endpoint.
- **credentials** (String) Full link to referenced Opaque Secret.

<a id="nestedblock--otel_tracing"></a>

### `otel_tracing`

- **sampling** (Int) Determines what percentage of requests should be traced.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or internal endpoint.
- **custom_tags** (Map of String) Key-value map of custom tags.

<a id="nestedblock--controlplane_tracing"></a>

### `controlplane_tracing`

- **sampling** (Int) Determines what percentage of requests should be traced.
- **custom_tags** (Map of String) Key-value map of custom tags.

<a id="nestedblock--load_balancer"></a>

### `load_balancer`

- **dedicated** (Boolean) Creates a dedicated load balancer in each location and enables additional Domain features: custom ports, protocols and wildcard hostnames. Charges apply for each location.

- **trusted_proxies** (Int) Controls the address used for request logging and for setting the X-Envoy-External-Address header. If set to 1, then the last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If set to 2, then the second to last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If the XFF header does not have at least two addresses or does not exist then the source client IP address will be used instead.

## Example Usage

```terraform
data "cpln_gvc" "gvc" {}

output "gvc_id" {
  value = data.cpln_gvc.gvc.id
}

output "gvc_locations" {
  value = data.cpln_gvc.gvc.locations
}
```