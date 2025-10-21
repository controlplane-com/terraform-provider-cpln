---
page_title: "cpln_location Data Source - terraform-provider-cpln"
subcategory: "Location"
description: |-
---

# cpln_location (Data Source)

Use this data source to access information about a [Location](https://docs.controlplane.com/reference/location) within Control Plane.

## Required

- **name** (String) Name of the location (i.e. `aws-us-west-2`).

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the location.
- **name** (String) Name of the location.
- **description** (String) Description of the location.
- **tags** (Map of String) Key-value map of resource tags.
- **origin** (String)
- **cloud_provider** (String) Cloud Provider of the location.
- **region** (String) Region of the location.
- **enabled** (Boolean) Indication if location is enabled.
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
data "cpln_location" "location" {
    name = "aws-us-west-2"
}

output "location" {
  value = data.cpln_location.location
}

output "location_enabled" {
  value = data.cpln_location.location.enabled
}
```
