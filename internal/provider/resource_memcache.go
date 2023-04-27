package cpln

import (
	"context"
	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMemcache() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceMemcacheCreate,
		ReadContext:   resourceMemcacheRead,
		UpdateContext: resourceMemcacheUpdate,
		DeleteContext: resourceMemcacheDelete,
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
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_size": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eviction_disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"idle_timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_item_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"locations": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceMemcacheCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Build Memcache Cluster
	memcache := client.Memcache{
		Spec: &client.MemcacheClusterSpec{
			NodeCount:   GetInt(d.Get("node_count")),
			NodeSizeGiB: GetFloat64(d.Get("node_size")),
		},
	}

	memcache.Name = GetString(d.Get("name"))
	memcache.Description = GetDescriptionString(d.Get("description"), *memcache.Name)
	memcache.Tags = GetStringMap(d.Get("tags"))

	if d.Get("version") != nil {
		memcache.Spec.Version = GetString(d.Get("version"))
	}

	memcache.Spec.Options = buildMemcacheOptions(d.Get("options").([]interface{}))
	memcache.Spec.Locations = buildStringsArrayFromSet(d.Get("locations"))

	// Post Memcache
	c := m.(*client.Client)
	newMemcache, code, err := c.CreateMemcache(memcache)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setMemcache(d, newMemcache)
}

func resourceMemcacheRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	memcache, code, err := c.GetMemcache(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setMemcache(d, memcache)
}

func resourceMemcacheUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "node_count", "node_size", "version", "options", "locations") {

		memcacheToUpdate := client.Memcache{
			Spec: &client.MemcacheClusterSpec{},
		}
		memcacheToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			memcacheToUpdate.Description = GetDescriptionString(d.Get("description"), *memcacheToUpdate.Name)
		}

		if d.HasChange("tags") {
			memcacheToUpdate.Tags = GetTagChanges(d)
		}

		if d.HasChange("node_count") {
			memcacheToUpdate.Spec.NodeCount = GetInt(d.Get("node_count"))
		}

		if d.HasChange("node_size") {
			memcacheToUpdate.Spec.NodeSizeGiB = GetFloat64(d.Get("node_size"))
		}

		if d.HasChange("version") {
			memcacheToUpdate.Spec.Version = GetString(d.Get("version"))
		}

		if d.HasChange("options") {
			memcacheToUpdate.Spec.Options = buildMemcacheOptions(d.Get("options").([]interface{}))
		}

		if d.HasChange("locations") {
			memcacheToUpdate.Spec.Locations = buildStringsArrayFromSet(d.Get("locations"))
		}

		c := m.(*client.Client)
		updatedMemcache, _, err := c.UpdateMemcache(memcacheToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setMemcache(d, updatedMemcache)
	}

	return nil
}

func resourceMemcacheDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	err := c.DeleteMemcache(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setMemcache(d *schema.ResourceData, memcache *client.Memcache) diag.Diagnostics {

	if memcache == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*memcache.Name)

	if err := SetBase(d, memcache.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(memcache.Links, d); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("node_count", memcache.Spec.NodeCount); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("node_size", memcache.Spec.NodeSizeGiB); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("version", memcache.Spec.Version); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("options", flattenMemcacheOptions(memcache.Spec.Options)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", flattenReferencedStringsArray(memcache.Spec.Locations)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

/*** Build Functions ***/
func buildMemcacheOptions(specs []interface{}) *client.MemcacheOptions {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	result := &client.MemcacheOptions{}

	if spec["eviction_disabled"] != nil {
		result.EvictionsDisabled = GetBool(spec["eviction_disabled"])
	}

	if spec["idle_timeout_seconds"] != nil {
		result.IdleTimeoutSeconds = GetInt(spec["idle_timeout_seconds"])
	}

	if spec["max_item_size"] != nil {
		result.MaxItemSizeKiB = GetInt(spec["max_item_size"])
	}

	if spec["max_connections"] != nil {
		result.MaxConnections = GetInt(spec["max_connections"])
	}

	return result
}

/*** Flatten Functions ***/
func flattenMemcacheOptions(options *client.MemcacheOptions) []interface{} {

	if options == nil {
		return nil
	}

	result := make(map[string]interface{})

	if options.EvictionsDisabled != nil {
		result["eviction_disabled"] = *options.EvictionsDisabled
	}

	if options.IdleTimeoutSeconds != nil {
		result["idle_timeout_seconds"] = *options.IdleTimeoutSeconds
	}

	if options.MaxItemSizeKiB != nil {
		result["max_item_size"] = *options.MaxItemSizeKiB
	}

	if options.MaxConnections != nil {
		result["max_connections"] = *options.MaxConnections
	}

	return []interface{}{
		result,
	}
}
