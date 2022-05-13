---
page_title: "cpln_service_account Resource - terraform-provider-cpln"
subcategory: "Service Account"
description: |-
  
---

# cpln_service_account (Resource)

Manages an org's [Service Accounts](https://docs.controlplane.com/reference/serviceaccount).


## Declaration

### Required

- **name** (String) Name of the Service Account.

### Optional

- **description** (String) Description of the Service Account.
- **tags** (Map of String) Key-value map of resource tags.

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Secret.
- **origin** (String) Origin of the Policy. Either `builtin` or `default`.
- **self_link** (String) Full link to this resource. Can be referenced by other resources. 

## Example Usage

```terraform
resource "cpln_service_account" "example" {

  name        = "service-account-example"
  description = "Example Service Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_service_account_key" "example" {

  service_account_name = cpln_service_account.example.name
  description          = "Service Account Key"
}


resource "cpln_service_account_key" "example_02" {

  // When adding another key, use `depends_on` to add the keys synchronously 
  depends_on = [cpln_service_account_key.example]

  service_account_name = cpln_service_account.example.name
  description          = "Service Account Key #2"
}

output "key_01" {
  value = cpln_service_account_key.example.key
}

output "key_02" {
  value = cpln_service_account_key.example_02.key
}
```
