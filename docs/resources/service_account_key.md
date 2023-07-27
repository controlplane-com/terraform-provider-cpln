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

## Import Syntax

To update a statefile with an existing service account key resource, execute the following import command:

```terraform
terraform import cpln_service_account_key.RESOURCE_NAME SERVICE_ACCOUNT_NAME:KEY_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute SERVICE_ACCOUNT_NAME with the name of the service account and KEY_NAME with the corresponding key name defined in the resource. (key name can be obtained from the console UI or CLI).

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
