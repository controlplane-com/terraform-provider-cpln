---
page_title: "cpln_volume_set Resource - terraform-provider-cpln"
subcategory: "Volume Set"
description: |-
---

# cpln_volume_set (Resource)

A [volume set](https://docs.controlplane.com/reference/volumeset) is a collection of storage volumes. Each volume set can be used by at most one [stateful workload](https://docs.controlplane.com/reference/workload#stateful). Volumes are not deleted until the volume set is deleted.

Refer to the [Volume Set Reference Page](https://docs.controlplane.com/reference/volumeset) for additional details.

## Declaration

### Required

- **name** (String) Name of the Volume Set.
- **gvc** (String) Name of the associated GVC.
- **initial_capacity** (Integer) The initial size in GB of volumes in this set. Minimum value: `10`.
- **performance_class** (String) Each volume set has a single, immutable, performance class. Valid classes: `general-purpose-ssd` or `high-throughput-ssd`
- **file_system_type** (String) Each volume set has a single, immutable file system. Valid types: `xfs` or `ext4`

### Optional

- **description** (String) Description of the Volume Set.
- **tags** (Map of String) Key-value map of resource tags.
- **snapshots** (Block List, Max: 1) ([see below](#nestedblock--snapshots)).
- **autoscaling** (Block List, Max: 1) ([see below](#nestedblock--autoscaling)).

<a id="nestedblock--snapshots"></a>

### `snapshots`

Optional:

- **create_final_snapshot** (Boolean) If true, a volume snapshot will be created immediately before deletion of any volume in this set. Default: `true`
- **retention_duration** (String) The default retention period for volume snapshots. This string should contain a floating point number followed by either d, h, or m. For example, "10d" would retain snapshots for 10 days.
- **schedule** (String) A standard cron schedule expression used to determine when your job should execute.

<a id="nestedblock--autoscaling"></a>

### `autoscaling`

Required:

- **max_capacity** (Integer) The maximum size in GB for a volume in this set. A volume cannot grow to be bigger than this value. Minimum value: `10`.
- **min_free_percentage** (Integer) The guaranteed free space on the volume as a percentage of the volume's total size. Control Plane will try to maintain at least that many percent free by scaling up the total size. Minimum percentage: `1`. Maximum Percentage: `100`.
- **scaling_factor** (Float64) When scaling is necessary, then `new_capacity = current_capacity * storageScalingFactor`. Minimum value: `1.1`.

## Outputs

- **cpln_id** (String) ID, in GUID format, of the Volume Set.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (Block List, Max: 1) ([see below](#nestedblock--status)).

<a id="nestedblock--status"></a>

### `status`

- **parent_id** (String) The GVC ID.
- **used_by_workload** (String) The url of the workload currently using this volume set (if any).
- **locations** (List of String) Contains a list of actual volumes grouped by location.

## Example Usage

```terraform
resource "cpln_gvc" "new" {
    name        = "gvc-for-volume-set"
    description = "This is a GVC description"

    locations = ["aws-eu-central-1", "aws-us-west-2"]

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }
}

resource "cpln_volume_set" "new" {

    name 		= "volume-set-example"
    description = "This is a Volume Set description"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    gvc 			  = cpln_gvc.new.name
    initial_capacity  = 1000
    performance_class = "high-throughput-ssd"
    file_system_type  = "xfs"

    snapshots {
        create_final_snapshot = false
        retention_duration    = "2d"
        schedule              = "* * 1 * 1"
    }

    autoscaling {
        max_capacity        = 2048
        min_free_percentage = 2
        scaling_factor      = 2.2
    }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing volume set resource, execute the following import command:

```terraform
terraform import cpln_volume_set.RESOURCE_NAME GVC_NAME:VOLUME_SET_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute GVC_NAME and VOLUME_SET_NAME with the corresponding GVC and volume set name defined in the resource.
