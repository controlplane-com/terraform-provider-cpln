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
- **enabled** (Boolean) Indication if custom location is enabled.
- **cloud_provider** (String) Provider of the custom location, must be `byok`.

### Optional

- **tags** (Map of String) Key-value map of resource tags.
- **description** - (String) Description of Custom Location.

## Outputs

- **cpln_id** (String) The ID, in GUID format, of the custom location.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform
resource "cpln_custom_location" "example" {
    name            = "custom-location-1"
    description 	= "custom location description"
    cloud_provider  = "byok"
    enabled 	 	= "true"

    tags = {
        "foo"   = "bar"
        "baz"	= "qux"
    }
}

```
