---
page_title: "cpln_gvc Data Source - terraform-provider-cpln"
subcategory: "Global Virtual Cloud"
description: |-
  
---
# cpln_gvc (Data Source)

Use this data source to access information about an existing [Global Virtual Cloud (GVC)](https://docs.controlplane.com/reference/gvc) within Control Plane. 

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the GVC.
- **name** (String) Name of the GVC.
- **locations** (List of String) A list of [locations](https://docs.controlplane.com/reference/location#current) making up the Global Virtual Cloud.
- **description** (String) Description of the GVC.
- **domain** (String) Custom domain name used by associated workloads.
- **pull_secrets** (List of String) A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.
- **tags** (Map of String) Key-value map of resource tags.
- **lightstep_tracing** (Block List, Max: 1) ([see below](#nestedblock--lightstep_tracing)).
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

<a id="nestedblock--lightstep_tracing"></a>
### `lightstep_tracing`

- **sampling** (Int) Sampling percentage.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internal endpoint.
- **credentials** (String) Full link to referenced Opaque Secret.

## Example Usage

```terraform
data "cpln_gvc" "gvc" {}

output "gvc_id" {
  value = data.cpln_gvc.id
}

output "gvc_locations" {
  value = data.cpln_gvc.locations
}
```