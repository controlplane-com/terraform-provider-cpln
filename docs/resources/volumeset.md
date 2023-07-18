---
page_title: "cpln_volume_set Resource - terraform-provider-cpln"
subcategory: "Volume Set"
description: |-
---

# cpln_volume_set (Resource)

A [volume set](https://docs.controlplane.com/reference/volume-sets) is a collection of storage volumes. Each volume set can be used by at most one [stateful workload](https://docs.controlplane.com/reference/workload#stateful). Volumes are not deleted until the volume set is deleted.

Refer to the [Volume Set Reference Page](https://docs.controlplane.com/reference/volume-sets) for additional details.

## Declaration

### Required

- **name** (String) Name of the Volume Set.
- **gvc** (String) Name of the associated GVC.
- **initial_capacity** (Integer) The initial size in GB of volumes in this set.
- **performance_class** (String) Each volume set has a single, immutable, performance class.
- **file_system_type** (String) Each volume set has a single, immutable file system.

### Optional

- **description** (String) Description of the Volume Set.
- **tags** (Map of String) Key-value map of resource tags.
- **snapshots** (Block List, Max: 1) ([see below](#nestedblock--snapshots)).
- **autoscaling** (Block List, Max: 1) ([see below](#nestedblock--autoscaling)).

<a id="nestedblock--snapshots"></a>

### `snapshots`

Optional:

- **create_final_snapshot** (Boolean) If true, a volume snapshot will be created immediately before deletion of any volume in this set.
- **retention_duration** (String) The default retention period for volume snapshots. This string should contain a floating point number followed by either d, h, or m. For example, "10d" would retain snapshots for 10 days.

<a id="nestedblock--autoscaling"></a>

### `autoscaling`

Optional:

- **max_capacity** (Integer) The maximum size in GB for a volume in this set. A volume cannot grow to be bigger than this value.
- **min_free_percentage** (Integer) The guaranteed free space on the volume as a percentage of the volume's total size. ControlPlane will try to maintain at least that many percent free by scaling up the total size.
- **scaling_factor** (Float64) When scaling is necessary, then `new_capacity = current_capacity * storageScalingFactor`.

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
    performance_class = "premium-low-latency-ssd"
    file_system_type  = "xfs"

    snapshots {
        create_final_snapshot = false
        retention_duration    = "2d"
    }

    autoscaling {
        max_capacity        = 2048
        min_free_percentage = 2
        scaling_factor      = 2.2
    }
}
```
