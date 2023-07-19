package cpln

import (
	"context"
	client "terraform-provider-cpln/internal/provider/client"

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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_by_workload": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locations": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
			"gvc": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"initial_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"performance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_system_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Optional: true,
				Default:  nil,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_final_snapshot": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"retention_duration": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeInt,
							Optional: true,
						},
						"min_free_percentage": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scaling_factor": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
		},
	}
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

	if d.HasChanges("description", "tags", "initial_capacity", "performance_class", "file_system_type", "snapshots", "autoscaling") {

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

	if status.Locations != nil {
		spec["locations"] = flattenVolumeSetStatusLocations(status.Locations)
	}

	return []interface{}{
		spec,
	}
}

func flattenVolumeSetStatusLocations(locations *[]string) []interface{} {

	result := make([]interface{}, len(*locations))

	for i, location := range *locations {
		result[i] = location
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
