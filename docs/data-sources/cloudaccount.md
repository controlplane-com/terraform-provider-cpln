---
page_title: "cpln_cloud_account Data Source - terraform-provider-cpln"
subcategory: "Cloud Account"
description: |-
  
---

# cpln_cloud_account (Data Source)

Use this data source to access information about an existing [Cloud Account](https://docs.controlplane.com/reference/cloudaccount) within Control Plane.

## Outputs

- **aws_identifiers** (String)

## Example Usage

```terraform
data "cpln_cloud_account" "this" {}

output "cloud_account" {
  value = data.cpln_cloud_account.this
}
```

