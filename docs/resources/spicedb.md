---
page_title: "cpln_spicedb Resource - terraform-provider-cpln"
subcategory: "SpiceDB Cluster"
description: |-
---

# cpln_spicedb (Resource)

## Declaration

## Required

- **name** (String) Name of the SpiceDB Cluster.
- **version** (String) //TODO: Add description
- **locations** (List of String) A list of [locations](https://docs.controlplane.com/reference/location#current).


## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **external_link** (String) //TODO: Add description

```terraform
resource "cpln_spicedb" "new" {
    name 		= "new-spicedb"
    description = "This is a SpiceDB Cluster description"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version   = "1.14.1"
    locations = ["aws-eu-central-1", "aws-us-west-2"]
}
```