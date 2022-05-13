---
page_title: "cpln_domain Resource - terraform-provider-cpln"
subcategory: "Domain"
description: |-
  
---

# cpln_domain (Resource)

Manages an org's custom [Domain](https://docs.controlplane.com/reference/domain).

The required DNS entries must exist before using Terraform to manage a `Domain`.

Refer to the [Configure a Domain](https://docs.controlplane.com/guides/configure-domain#dns-entries)
page for additional details. 

During the creation of a domain, Control Plane will verify that the DNS entries exists. If they do 
not exist, the Terraform script will fail.

## Declaration

### Required

- **name** (String) Domain name. Must be a valid domain name with at least three segments (e.g., test.example.com). Control Plane will validate the existence of the domain with DNS. Create and Update will fail if the required DNS entries cannot be validated.

### Optional

- **description** (String) Description for the domain name.
- **tags** (Map of String) Key-value map of resource tags.

## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources. 

## Example Usage

```terraform
resource "cpln_domain" "example" {

  name        = "app.example.com"
  description = "Custom domain that can be set on a GVC and used by associated workloads"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}
```

