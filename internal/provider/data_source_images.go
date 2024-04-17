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
			"query": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     QuerySchemaResource(),
			},
		},
	}
}

func dataSourceImagesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	query := client.Query{
		Kind: GetString("image"),
		Spec: &client.Spec{
			Match: GetString("all"),
		},
	}

	if d.Get("query") != nil {
		builtQuery := BuildQueryHelper("image", d.Get("query"))

		if builtQuery != nil {
			query = *builtQuery
		}
	}

	images, err := c.GetImagesQuery(query)

	if err != nil {
		return diag.FromErr(err)
	}

	return setImages(d, images)
}

func setImages(d *schema.ResourceData, images *client.ImagesQueryResult) diag.Diagnostics {

	if err := d.Set("images", flattenImageItems(&images.Items)); err != nil {
		return diag.FromErr(err)
	}

	flattenedQuery, err := FlattenQueryHelper(&images.Query)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("query", flattenedQuery); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

/*** Flatten ***/

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
