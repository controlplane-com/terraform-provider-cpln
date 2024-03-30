package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Structs ***/
type Images struct {
	Kind     string  `json:"kind,omitempty"`
	ItemKind string  `json:"itemKind,omitempty"`
	Items    []Image `json:"items,omitempty"`
	Links    []Link  `json:"links,omitempty"`
}

type Image struct {
	Base
	Tag        *string        `json:"tag,omitempty"`
	Repository *string        `json:"repository,omitempty"`
	Digest     *string        `json:"digest,omitempty"`
	Manifest   *ImageManifest `json:"manifest,omitempty"`
}

type ImageManifest struct {
	Config        *ImageManifestConfig   `json:"config,omitempty"`
	Layers        *[]ImageManifestConfig `json:"layers,omitempty"`
	MediaType     *string                `json:"mediaType,omitempty"`
	SchemaVersion *int                   `json:"schemaVersion,omitempty"`
}

type ImageManifestConfig struct {
	Size      *int    `json:"size,omitempty"`
	Digest    *string `json:"digest,omitempty"`
	MediaType *string `json:"mediaType,omitempty"`
}

/*** Schema ***/
func ImageSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"cpln_id": {
			Type:        schema.TypeString,
			Description: "The ID, in GUID format, of the Image.",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the Image.",
			Required:    true,
		},
		"tags": {
			Type:        schema.TypeMap,
			Description: "Key-value map of resource tags.",
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"self_link": {
			Type:        schema.TypeString,
			Description: "Full link to this resource. Can be referenced by other resources.",
			Computed:    true,
		},
		"tag": {
			Type:        schema.TypeString,
			Description: "Tag of the image.",
			Computed:    true,
		},
		"repository": {
			Type:        schema.TypeString,
			Description: "Respository name of the image.",
			Computed:    true,
		},
		"digest": {
			Type:        schema.TypeString,
			Description: "// TODO: Add description",
			Computed:    true,
		},
		"manifest": {
			Type:        schema.TypeList,
			Description: "// TODO: Add description",
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"config": {
						Type:        schema.TypeList,
						Description: "// TODO: Add description",
						Computed:    true,
						Elem:        ImageManifestConfigSchemaResource(),
					},
					"layers": {
						Type:        schema.TypeList,
						Description: "// TODO: Add description",
						Computed:    true,
						Elem:        ImageManifestConfigSchemaResource(),
					},
					"media_type": {
						Type:        schema.TypeString,
						Description: "// TODO: Add description",
						Computed:    true,
					},
					"schema_version": {
						Type:        schema.TypeInt,
						Description: "// TODO: Add description",
						Computed:    true,
					},
				},
			},
		},
	}
}

func ImagesSchema() map[string]*schema.Schema {
	schema := ImageSchema()
	(*schema["name"]).Required = false
	(*schema["name"]).Computed = true

	return schema
}

func ImageManifestConfigSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"size": {
				Type:        schema.TypeInt,
				Description: "// TODO: Add description",
				Computed:    true,
			},
			"digest": {
				Type:        schema.TypeString,
				Description: "// TODO: Add description",
				Computed:    true,
			},
			"media_type": {
				Type:        schema.TypeString,
				Description: "// TODO: Add description",
				Computed:    true,
			},
		},
	}
}

/*** Functions ***/
func (c *Client) GetImage(name string) (*Image, int, error) {

	image, code, err := c.GetResource(fmt.Sprintf("image/%s", name), new(Image))

	if err != nil {
		return nil, code, err
	}

	return image.(*Image), code, err
}

func (c *Client) GetImages() (*Images, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s/image", c.HostURL, c.Org), nil)

	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, "")

	if err != nil {
		return nil, err
	}

	images := Images{}
	err = json.Unmarshal(body, &images)

	if err != nil {
		return nil, err
	}

	return &images, nil
}
