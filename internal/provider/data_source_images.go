package cpln

import (
	"context"
	"strconv"
	"time"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImages() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceImagesRead,
		Schema: map[string]*schema.Schema{
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: client.ImagesSchema(),
				},
			},
		},
	}
}

func dataSourceImagesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	images, err := c.GetImages()

	if err != nil {
		return diag.FromErr(err)
	}

	imageItems := flattenImageItems(&images.Items)

	if err := d.Set("images", imageItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

/*** Flatten Functions ***/
func flattenImageItems(imageItems *[]client.Image) []interface{} {
	if imageItems == nil {
		return []interface{}{}
	}

	specs := make([]interface{}, len(*imageItems))

	for i, item := range *imageItems {

		spec := make(map[string]interface{})

		spec["cpln_id"] = *item.ID
		spec["name"] = *item.Name
		spec["tags"] = GetTags(item.Tags)
		spec["self_link"] = GetSelfLink(item.Links)
		spec["tag"] = *item.Tag
		spec["repository"] = *item.Repository
		spec["digest"] = *item.Digest

		if item.Manifest != nil {
			spec["manifest"] = flattenImageManifest(item.Manifest)
		}

		specs[i] = spec
	}

	return specs
}
