---
page_title: "cpln_custom_location Resource - terraform-provider-cpln"
subcategory: "Custom Location"
description: |-
---

# cpln_custom_location (Resource)

Manages an org's [Custom Location](https://docs.controlplane.com/reference/location#byok-locations).

## Declaration

### Required

- **name** (String) Name of the Custom Location.
- **cloud_provider** (String) Provider of the custom location, Available providers: [`byok`].
- **enabled** (Boolean) Indication if custom location is enabled.

### Optional

- **description** (String) Description of Custom Location.
- **tags** (Map of String) Key-value map of resource tags.

## Outputs

- **cpln_id** (String) The ID, in GUID format, of the custom location.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **region** (String) Region of the location.

## Example Usage

```terraform
resource "cpln_custom_location" "example" {
  name           = "custom-location-1"
  description    = "custom location description"
  cloud_provider = "byok"
  enabled        = "true"

  tags = {
    "foo" = "bar"
    "baz" = "qux"
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing custom location resource, execute the following import command:

```terraform
terraform import cpln_custom_location.RESOURCE_NAME CUSTOM_LOCATION_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute CUSTOM_LOCATION_NAME with the corresponding custom location name defined in the resource.
