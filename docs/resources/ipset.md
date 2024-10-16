---
page_title: "cpln_ipset Resource - terraform-provider-cpln"
subcategory: "IpSet"
description: |-
---

# cpln_ipset (Resource)

Manages an org's IpSet.

## Declaration

### Required

- **name** (String) Name of the IpSet.

### Optional

- **description** - (String) Description of the IpSet.
- **tags** (Map of String) Key-value map of resource tags.
- **link** (String) The self link of a workload.
- **location** (Block List) ([see below](#nestedblock--location)).

<a id="nestedblock--location"></a>

### `location`

Required:

- **name** (String) The self link of a location.
- **retention_policy** (String) Exactly one of: `keep` and `free`.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the IpSet.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (Block List, Max: 1) ([see below](#nestedblock--status)).

<a id="nestedblock--status"></a>

### `status`

Status of the IpSet.

Read-Only:

- **ip_address** (Block List) ([see below](#nestedblock--status-ip_address))
- **error** (String)

<a id="nestedblock--status--ip_address"></a>

### `status.ip_address`

- **name** (String)
- **ip** (String)
- **id** (String)
- **state** (String)
- **created** (String)

## Example Usage

```terrafrom
resource "cpln_ipset" "new" {
		
  name        = "example"
  description = "example"

  tags = {
    terraform_generated = "true"
  }

  link = "SELF_LINK_TO_WORKLOAD"
  
  location {
    name             = "SELF_LINK_TO_LOCATION"
    retention_policy = "keep"
  }
}
```