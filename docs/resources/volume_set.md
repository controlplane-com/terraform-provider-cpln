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
- **initial_capacity** (Integer) The initial volume size in this set, specified in GB. The minimum size for the performance class `general-purpose-ssd` is `10 GB`, while `high-throughput-ssd` requires at least `200 GB`.
- **performance_class** (String) Each volume set has a single, immutable, performance class. Valid classes: `general-purpose-ssd` or `high-throughput-ssd`

### Optional

- **description** (String) Description of the Volume Set.
- **tags** (Map of String) Key-value map of resource tags.
- **storage_class_suffix** (String) For self-hosted locations only. The storage class used for volumes in this set will be {performanceClass}-{fileSystemType}-{storageClassSuffix} if it exists, otherwise it will be {performanceClass}-{fileSystemType}
- **file_system_type** (String) Each volume set has a single, immutable file system. Valid types: `xfs` or `ext4`. Default: `ext4`.
- **snapshots** (Block List, Max: 1) ([see below](#nestedblock--snapshots)).
- **autoscaling** (Block List, Max: 1) ([see below](#nestedblock--autoscaling)).
- **mount_options** (Block List, Max: 1) ([see below](#nestedblock--mount_options))

<a id="nestedblock--snapshots"></a>

### `snapshots`

Optional:

- **create_final_snapshot** (Boolean) If true, a volume snapshot will be created immediately before deletion of any volume in this set. Default: `true`
- **retention_duration** (String) The default retention period for volume snapshots. This string should contain a floating point number followed by either d, h, or m. For example, "10d" would retain snapshots for 10 days.
- **schedule** (String) A standard cron schedule expression used to determine when a snapshot will be taken. (i.e., `0 * * * *` Every hour). Note: snapshots cannot be scheduled more often than once per hour.

~> Use a tool, such as [Crontab Guru](https://crontab.guru/), to easily generate a cron schedule expression.

<a id="nestedblock--autoscaling"></a>

### `autoscaling`

Required:

- **max_capacity** (Integer) The maximum size in GB for a volume in this set. A volume cannot grow to be bigger than this value. Minimum value: `10`.
- **min_free_percentage** (Integer) The guaranteed free space on the volume as a percentage of the volume's total size. Control Plane will try to maintain at least that many percent free by scaling up the total size. Minimum percentage: `1`. Maximum Percentage: `100`.
- **scaling_factor** (Float64) When scaling is necessary, then `new_capacity = current_capacity * storageScalingFactor`. Minimum value: `1.1`.

<a id="nestedblock--mount_options"></a>

### `mount_options`

Optionals:

- **max_cpu** (String) Default: 2000m
- **min_cpu** (String) Default: 500m
- **min_memory** (String) Default: 1Gi
- **max_memory** (String) Default: 2Gi

## Outputs

- **cpln_id** (String) ID, in GUID format, of the Volume Set.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (Block List, Max: 1) ([see below](#nestedblock--status)).
- **volumeset_link** (String) Output used when linking a volume set to a workload, in the format: `cpln://volumeset/VOLUME_SET_NAME`.

<a id="nestedblock--status"></a>

### `status`

- **parent_id** (String) The GVC ID.
- **used_by_workload** (String) The url of the workload currently using this volume set (if any).
- **binding_id** (String) Uniquely identifies the connection between the volume set and its workload. Every time a new connection is made, a new id is generated (e.g., If a workload is updated to remove the volume set, then updated again to reattach it, the volume set will have a new binding id).
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

  gvc 			     = cpln_gvc.new.name
  initial_capacity     = 1000
  performance_class    = "high-throughput-ssd"
  file_system_type     = "xfs"
  storage_class_suffix = "demo-class"

  snapshots {
    create_final_snapshot = false
    retention_duration    = "2d"
    schedule              = "0 * * * *"
  }

  autoscaling {
    max_capacity        = 2048
    min_free_percentage = 2
    scaling_factor      = 2.2
  }

  mount_options {
    resources {}
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
