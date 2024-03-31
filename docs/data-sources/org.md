---
page_title: "cpln_org Data Source - terraform-provider-cpln"
subcategory: "Org"
description: |-
  
---
# cpln_org (Data Source)

Output the ID and name of the current [org](https://docs.controlplane.com/reference/org). 

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the org.
- **name** (String) The name of org.

## Example Usage

```terraform
data "cpln_org" "org" {}

output "org_id" {
  value = data.cpln_org.org.id
}

output "org_name" {
  value = data.cpln_org.org.name
}
```