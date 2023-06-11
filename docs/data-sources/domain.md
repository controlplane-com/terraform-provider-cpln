---
page_title: "cpln_domain Data Source - terraform-provider-cpln"
subcategory: "Domain"
description: |-
  
---
# cpln_domain (Data Source)

Use this data source to access information about a [Domain](https://docs.controlplane.com/reference/domain) within Control Plane.

## Required

- **name** (String) Name of the domain.

## Outputs

- **cpln_id** (String) The ID, in GUID format, of the domain.
- **name** (String) Name of the domain.
- **description** (String) Description of the domain.
- **tags** (Map of String) Key-value map of resource tags.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform
data "cpln_domain" "domain" {}

output "domain" {
  value = data.cpln_domain.domain
}
```