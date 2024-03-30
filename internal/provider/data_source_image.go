package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,
		Schema:      client.ImageSchema(),
	}
}

func dataSourceImageRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	imageName := d.Get("name").(string)

	image, _, err := c.GetImage(imageName)

	if err != nil {
		return diag.FromErr(err)
	}

	return setImage(d, image)
}

/*** Flatten Functions ***/
func flattenImageManifest(manifest *client.ImageManifest) []interface{} {
	if manifest == nil {
		return nil
	}

	spec := make(map[string]interface{})

	if manifest.Config != nil {
		spec["config"] = flattenImageManifestConfig(manifest.Config)
	}

	if manifest.Layers != nil {
		spec["layers"] = flattenImageManifestLayers(manifest.Layers)
	}

	if manifest.MediaType != nil {
		spec["media_type"] = *manifest.MediaType
	}

	if manifest.SchemaVersion != nil {
		spec["schema_version"] = *manifest.SchemaVersion
	}

	return []interface{}{
		spec,
	}
}

func flattenImageManifestConfig(config *client.ImageManifestConfig) []interface{} {
	if config == nil {
		return nil
	}

	spec := make(map[string]interface{})

	if config.Size != nil {
		spec["size"] = *config.Size
	}

	if config.Digest != nil {
		spec["digest"] = *config.Digest
	}

	if config.MediaType != nil {
		spec["media_type"] = *config.MediaType
	}

	return []interface{}{
		spec,
	}
}

func flattenImageManifestLayers(layers *[]client.ImageManifestConfig) []interface{} {
	if len(*layers) == 0 {
		return nil
	}

	specs := []interface{}{}

	for _, layer := range *layers {
		flattenedLayer := flattenImageManifestConfig(&layer)

		if flattenedLayer == nil {
			continue
		}

		specs = append(specs, flattenedLayer...)
	}

	return specs
}

/*** Helper Functions ***/
func setImage(d *schema.ResourceData, image *client.Image) diag.Diagnostics {

	if image == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*image.Name)

	if err := SetBase(d, image.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tag", image.Tag); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("repository", image.Repository); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("digest", image.Digest); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("manifest", flattenImageManifest(image.Manifest)); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(image.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
