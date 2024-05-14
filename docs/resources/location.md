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
- **tags** (Map of String) Key-value map of resource tags.
- **enabled** (Boolean) Indication if location is enabled.

~> **Note** You need to associate the same tags that are defined in a location; otherwise, the Terraform plan will not be empty. It is common practice to reference the tags from a location data source.

## Outputs

- **cpln_id** (String) The ID, in GUID format, of the location.
- **description** (String) Description of the location.
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

### Reference Tags from Data Source

```terraform
data "cpln_location" "main-location" {
    name = "aws-eu-central-1"
}

resource "cpln_location" "reference-tags-example" {
    name    = "aws-eu-central-1"
    enabled = true

    tags = data.cpln_location.main-location.tags
}
```

### Hard Code Location Tags

```terraform
resource "cpln_location" "example" {
    name    = "aws-eu-central-1"
    enabled = true

    tags = {
        "cpln/city"      = "Frankfurt"
        "cpln/continent" = "Europe"
        "cpln/country"   = "Germany"
    }
}
```