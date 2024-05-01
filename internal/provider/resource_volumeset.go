package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Main ***/
func resourceVolumeSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeSetCreate,
		ReadContext:   resourceVolumeSetRead,
		UpdateContext: resourceVolumeSetUpdate,
		DeleteContext: resourceVolumeSetDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "ID, in GUID format, of the Volume Set.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the Volume Set.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the Volume Set.",
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:        schema.TypeMap,
				Description: "Key-value map of resource tags.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeList,
				Description: "Status of the Volume Set.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_id": {
							Type:        schema.TypeString,
							Description: "The GVC ID.",
							Computed:    true,
						},
						"used_by_workload": {
							Type:        schema.TypeString,
							Description: "The url of the workload currently using this volume set (if any).",
							Computed:    true,
						},
						"binding_id": {
							Type:        schema.TypeString,
							Description: "Uniquely identifies the connection between the volume set and its workload. Every time a new connection is made, a new id is generated (e.g., If a workload is updated to remove the volume set, then updated again to reattach it, the volume set will have a new binding id).",
							Computed:    true,
						},
						"locations": {
							Type:        schema.TypeSet,
							Description: "Contains a list of actual volumes grouped by location.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
					},
				},
			},
			"gvc": {
				Type:         schema.TypeString,
				Description:  "Name of the associated GVC.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"initial_capacity": {
				Type:        schema.TypeInt,
				Description: "The initial size in GB of volumes in this set. Minimum value: `10`.",
				Required:    true,
			},
			"performance_class": {
				Type:        schema.TypeString,
				Description: "Each volume set has a single, immutable, performance class. Valid classes: `general-purpose-ssd` or `high-throughput-ssd`",
				Required:    true,
				ForceNew:    true,
			},
			"storage_class_suffix": {
				Type:        schema.TypeString,
				Description: "For self-hosted locations only. The storage class used for volumes in this set will be {performanceClass}-{fileSystemType}-{storageClassSuffix} if it exists, otherwise it will be {performanceClass}-{fileSystemType}",
				Optional:    true,
			},
			"file_system_type": {
				Type:        schema.TypeString,
				Description: "Each volume set has a single, immutable file system. Valid types: `xfs` or `ext4`",
				Optional:    true,
				ForceNew:    true,
				Default:     "ext4",
			},
			"snapshots": {
				Type:     schema.TypeList,
				Description: "Point-in-time copies of data stored within the volume set, capturing the state of the data at a specific moment.",
				Optional: true,
				Default:  nil,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_final_snapshot": {
							Type:        schema.TypeBool,
							Description: "If true, a volume snapshot will be created immediately before deletion of any volume in this set. Default: `true`",
							Optional:    true,
							Default:     true,
						},
						"retention_duration": {
							Type:        schema.TypeString,
							Description: "The default retention period for volume snapshots. This string should contain a floating point number followed by either d, h, or m. For example, \"10d\" would retain snapshots for 10 days.",
							Optional:    true,
						},
						"schedule": {
							Type:        schema.TypeString,
							Description: "A standard cron schedule expression used to determine when a snapshot will be taken. (i.e., `0 * * * *` Every hour). Note: snapshots cannot be scheduled more often than once per hour.",
							Optional:    true,
						},
					},
				},
			},
			"autoscaling": {
				Type:     schema.TypeList,
				Optional: true,
				Default:  nil,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_capacity": {
							Type:        schema.TypeInt,
							Description: "The maximum size in GB for a volume in this set. A volume cannot grow to be bigger than this value. Minimum value: `10`.",
							Required:    true,
						},
						"min_free_percentage": {
							Type:        schema.TypeInt,
							Description: "The guaranteed free space on the volume as a percentage of the volume's total size. Control Plane will try to maintain at least that many percent free by scaling up the total size. Minimum percentage: `1`. Maximum Percentage: `100`.",
							Required:    true,
						},
						"scaling_factor": {
							Type:        schema.TypeFloat,
							Description: "When scaling is necessary, then `new_capacity = current_capacity * storageScalingFactor`. Minimum value: `1.1`.",
							Required:    true,
						},
					},
				},
			},
			"volumeset_link": {
				Type:        schema.TypeString,
				Description: "Output used when linking a volume set to a workload.",
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateVolumeSet,
		},
	}
}

func importStateVolumeSet(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.SplitN(d.Id(), ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected ID syntax: 'gvc:volume_set'. Example: 'terraform import cpln_volume_set.RESOURCE_NAME GVC_NAME:VOLUME_SET_NAME'", d.Id())
	}

	d.Set("gvc", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func setVolumeSet(d *schema.ResourceData, volumeSet *client.VolumeSet) diag.Diagnostics {

	if volumeSet == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*volumeSet.Name)

	if err := SetBase(d, volumeSet.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(volumeSet.Links, d); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("volumeset_link", "cpln://volumeset/"+*volumeSet.Name); err != nil {
		return diag.FromErr(err)
	}

	// Set VolumeSet Status
	if err := d.Set("status", flattenVolumeSetStatus(volumeSet.Status)); err != nil {
		return diag.FromErr(err)
	}

	// Set VolumeSet Spec
	if err := d.Set("initial_capacity", volumeSet.Spec.InitialCapacity); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("performance_class", volumeSet.Spec.PerformanceClass); err != nil {
		return diag.FromErr(err)
	}

	if volumeSet.Spec.StorageClassSuffix != nil {

		if err := d.Set("storage_class_suffix", volumeSet.Spec.StorageClassSuffix); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("file_system_type", volumeSet.Spec.FileSystemType); err != nil {
		return diag.FromErr(err)
	}

	if volumeSet.Spec.Snapshots != nil {

		if err := d.Set("snapshots", flattenVolumeSetSnapshots(volumeSet.Spec.Snapshots)); err != nil {
			return diag.FromErr(err)
		}
	}

	if volumeSet.Spec.AutoScaling != nil {

		if err := d.Set("autoscaling", flattenVolumeSetAutoscaling(volumeSet.Spec.AutoScaling)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceVolumeSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	volumeSet := client.VolumeSet{}

	volumeSet.Name = GetString(d.Get("name"))
	volumeSet.Description = GetDescriptionString(d.Get("description"), *volumeSet.Name)
	volumeSet.Tags = GetStringMap(d.Get("tags"))
	volumeSet.Spec = &client.VolumeSetSpec{
		InitialCapacity:  GetInt(d.Get("initial_capacity")),
		PerformanceClass: GetString(d.Get("performance_class")),
		FileSystemType:   GetString(d.Get("file_system_type")),
	}

	if d.Get("storage_class_suffix") != nil {
		volumeSet.Spec.StorageClassSuffix = GetString(d.Get("storage_class_suffix"))
	}

	if d.Get("snapshots") != nil {
		volumeSet.Spec.Snapshots = buildVolumeSetSnapshots(d.Get("snapshots").([]interface{}))
	}

	if d.Get("autoscaling") != nil {
		volumeSet.Spec.AutoScaling = buildVolumeSetAutoscaling(d.Get("autoscaling").([]interface{}))
	}

	c := m.(*client.Client)
	newVolumeSet, code, err := c.CreateVolumeSet(volumeSet, d.Get("gvc").(string))

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setVolumeSet(d, newVolumeSet)
}

func resourceVolumeSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	volumeSet, code, err := c.GetVolumeSet(d.Id(), d.Get("gvc").(string))

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setVolumeSet(d, volumeSet)
}

func resourceVolumeSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "initial_capacity", "performance_class", "storage_class_suffix", "file_system_type", "snapshots", "autoscaling") {

		volumeSetToUpdate := client.VolumeSet{
			SpecReplace: &client.VolumeSetSpec{
				InitialCapacity:  GetInt(d.Get("initial_capacity")),
				PerformanceClass: GetString(d.Get("performance_class")),
				FileSystemType:   GetString(d.Get("file_system_type")),
			},
		}
		volumeSetToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			volumeSetToUpdate.Description = GetDescriptionString(d.Get("description"), *volumeSetToUpdate.Name)
		}

		if d.HasChange("tags") {
			volumeSetToUpdate.Tags = GetTagChanges(d)
		}

		if d.HasChange("storage_class_suffix") {
			volumeSetToUpdate.SpecReplace.StorageClassSuffix = GetString(d.Get("storage_class_suffix"))
		}

		if d.HasChange("snapshots") {
			volumeSetToUpdate.SpecReplace.Snapshots = buildVolumeSetSnapshots(d.Get("snapshots").([]interface{}))
		}

		if d.HasChange("autoscaling") {
			volumeSetToUpdate.SpecReplace.AutoScaling = buildVolumeSetAutoscaling(d.Get("autoscaling").([]interface{}))
		}

		// Perform update
		c := m.(*client.Client)
		updatedVolumeSet, _, err := c.UpdateVolumeSet(volumeSetToUpdate, d.Get("gvc").(string))

		if err != nil {
			return diag.FromErr(err)
		}

		return setVolumeSet(d, updatedVolumeSet)
	}

	return nil
}

func resourceVolumeSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	err := c.DeleteVolumeSet(d.Id(), d.Get("gvc").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

/*** Build ***/
func buildVolumeSetSnapshots(specs []interface{}) *client.VolumeSetSnapshots {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.VolumeSetSnapshots{}

	if spec["create_final_snapshot"] != nil {
		output.CreateFinalSnapshot = GetBool(spec["create_final_snapshot"])
	}

	if spec["retention_duration"] != nil {
		output.RetentionDuration = GetString(spec["retention_duration"])
	}

	if spec["schedule"] != nil {
		output.Schedule = GetString(spec["schedule"])
	}

	return &output
}

func buildVolumeSetAutoscaling(specs []interface{}) *client.VolumeSetScaling {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.VolumeSetScaling{}

	if spec["max_capacity"] != nil {
		output.MaxCapacity = GetInt(spec["max_capacity"])
	}

	if spec["min_free_percentage"] != nil {
		output.MinFreePercentage = GetInt(spec["min_free_percentage"])
	}

	if spec["scaling_factor"] != nil {
		output.ScalingFactor = GetFloat64(spec["scaling_factor"])
	}

	return &output
}

/*** Flatten ***/
func flattenVolumeSetStatus(status *client.VolumeSetStatus) []interface{} {

	spec := map[string]interface{}{}

	if status.ParentID != nil {
		spec["parent_id"] = *status.ParentID
	}

	if status.UsedByWorkload != nil {
		spec["used_by_workload"] = *status.UsedByWorkload
	}

	if status.BindingID != nil {
		spec["binding_id"] = *status.BindingID
	}

	if status.Locations != nil {
		spec["locations"] = flattenVolumeSetStatusLocations(status.Locations)
	}

	return []interface{}{
		spec,
	}
}

func flattenVolumeSetStatusLocations(locations *[]interface{}) []interface{} {

	result := make([]interface{}, len(*locations))

	for i, location := range *locations {

		jsonData, err := json.Marshal(location)

		if err != nil {
			result[i] = fmt.Sprintf("Error serializing to JSON: %s", err)

		} else {
			result[i] = string(jsonData)
		}
	}

	return result
}

func flattenVolumeSetSnapshots(snapshots *client.VolumeSetSnapshots) []interface{} {

	if snapshots == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if snapshots.CreateFinalSnapshot != nil {
		spec["create_final_snapshot"] = *snapshots.CreateFinalSnapshot
	}

	if snapshots.RetentionDuration != nil {
		spec["retention_duration"] = *snapshots.RetentionDuration
	}

	if snapshots.Schedule != nil {
		spec["schedule"] = *snapshots.Schedule
	}

	return []interface{}{
		spec,
	}
}

func flattenVolumeSetAutoscaling(autoscaling *client.VolumeSetScaling) []interface{} {

	if autoscaling == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if autoscaling.MaxCapacity != nil {
		spec["max_capacity"] = *autoscaling.MaxCapacity
	}

	if autoscaling.MinFreePercentage != nil {
		spec["min_free_percentage"] = *autoscaling.MinFreePercentage
	}

	if autoscaling.ScalingFactor != nil {
		spec["scaling_factor"] = *autoscaling.ScalingFactor
	}

	return []interface{}{
		spec,
	}
}
