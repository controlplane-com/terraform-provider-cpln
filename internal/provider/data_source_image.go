package cpln

import (
	"context"
	"strings"

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

	name := d.Get("name").(string)
	hasColon := hasColon(name)

	var image *client.Image
	var err error

	// Get specific image if a colon is specified
	if hasColon {
		image, _, err = c.GetImage(name)
	} else {
		// Fetch latest image
		image, _, err = c.GetLatestImage(name)
	}

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

/*** Helpers ***/

func hasColon(input string) bool {

	parts := strings.SplitN(input, ":", 2)
	return len(parts) == 2
}
