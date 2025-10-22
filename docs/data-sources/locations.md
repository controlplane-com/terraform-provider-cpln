---
page_title: "cpln_locations Data Source - terraform-provider-cpln"
subcategory: "Location"
description: |-
---

# cpln_locations (Data Source)

Use this data source to access information about all [Locations](https://docs.controlplane.com/reference/location) within Control Plane.

## Outputs

The following attributes are exported:

- **locations** (Block List) ([see below](#nestedblock--locations)).

<a id="nestedblock--locations"></a>

### `locations`

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
data "cpln_locations" "locations" { }

output "locations" {
  value = data.cpln_locations.locations.locations
}
```
