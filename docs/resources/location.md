---
page_title: "cpln_location Resource - terraform-provider-cpln"
subcategory: "Location"
description: |-
---

# cpln_location (Resource)

Manages an org's [Location](https://docs.controlplane.com/reference/location).

## Declaration

### Required

- **name** (String) Name of the Location.
- **enabled** (Boolean) Indication if location is enabled.

~> **Note** You need to associate the same tags that are defined in a location; otherwise, the Terraform plan will not be empty. It is common practice to reference the tags from a location data source.

## Outputs

- **cpln_id** (String) The ID, in GUID format, of the location.
- **description** (String) Description of the location.
- **tags** (Map of String) Key-value map of resource tags.
- **cloud_provider** (String) Cloud Provider of the location.
- **region** (String) Region of the location.
- **geo** (Block List, Max: 1) ([see below](#nestedblock--geo))
- **ip_ranges** (List of String) A list of IP ranges of the location.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

<a id="nestedblock--geo"></a>

### `geo`

Location geographical details

- **lat** (Number) Latitude.
- **lon** (Number) Longitude.
- **country** (String) Country.
- **state** (String) State.
- **city** (String) City.
- **continent** (String) Continent.

## Example Usage

```terraform
resource "cpln_location" "example" {
    name    = "aws-eu-central-1"
    enabled = true
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing location resource, execute the following import command:

```terraform
terraform import cpln_location.RESOURCE_NAME LOCATION_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute LOCATION_NAME with the corresponding location name defined in the resource.
