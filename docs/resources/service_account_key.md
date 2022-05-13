---
page_title: "cpln_service_account_key Resource - terraform-provider-cpln"
subcategory: "Service Account"
description: |-
  
---

# cpln_service_account_key (Resource)

Manages an org's [Service Account Keys](https://docs.controlplane.com/reference/serviceaccount#keys).

Used in conjunction with a Service Account.

**A key can only be created and deleted. Updates will fail.**

## Declaration

### Required

- **description** (String) Description of the Service Account Key.
- **service_account_name** (String) The name of an existing Service Account this key will belong to.

## Outputs

The following attributes are exported:

- **created** (String) The timestamp, in UTC, when the key was created.
- **key** (String, Sensitive) The generated key.
- **name** (String) The generated name of the key.

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