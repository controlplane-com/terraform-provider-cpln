package cpln

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Structs ***/

type Image struct {
	Base
	LastModified *string        `json:"lastModified,omitempty"`
	Tag          *string        `json:"tag,omitempty"`
	Repository   *string        `json:"repository,omitempty"`
	Digest       *string        `json:"digest,omitempty"`
	Manifest     *ImageManifest `json:"manifest,omitempty"`
}

type ImagesQueryResult struct {
	Kind  string  `json:"kind,omitempty"`
	Items []Image `json:"items,omitempty"`
	Links []Link  `json:"links,omitempty"`
	Query Query   `json:"query,omitempty"`
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
			Description: "A unique SHA256 hash used to identify a specific image version within the image registry.",
			Computed:    true,
		},
		"manifest": {
			Type:        schema.TypeList,
			Description: "The manifest provides configuration and layers information about the image. It plays a crucial role in the Docker image distribution system, enabling image creation, verification, and replication in a consistent and secure manner.",
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"config": {
						Type:        schema.TypeList,
						Description: "The config is a JSON blob that contains the image configuration data which includes environment variables, default command to run, and other settings necessary to run the container based on this image.",
						Computed:    true,
						Elem:        ImageManifestConfigSchemaResource(),
					},
					"layers": {
						Type:        schema.TypeList,
						Description: "Layers lists the digests of the image's layers. These layers are filesystem changes or additions made in each step of the Docker image's creation process. The layers are stored separately and pulled as needed, which allows for efficient storage and transfer of images. Each layer is represented by a SHA256 digest, ensuring the integrity and authenticity of the image.",
						Computed:    true,
						Elem:        ImageManifestConfigSchemaResource(),
					},
					"media_type": {
						Type:        schema.TypeString,
						Description: "Specifies the type of the content represented in the manifest, allowing Docker clients and registries to understand how to handle the document correctly.",
						Computed:    true,
					},
					"schema_version": {
						Type:        schema.TypeInt,
						Description: "The version of the Docker Image Manifest format.",
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
				Description: "The size of the image or layer in bytes. This helps in estimating the space required and the download time.",
				Computed:    true,
			},
			"digest": {
				Type:        schema.TypeString,
				Description: "A unique SHA256 hash used to identify a specific image version within the image registry.",
				Computed:    true,
			},
			"media_type": {
				Type:        schema.TypeString,
				Description: "Specifies the type of the content represented in the manifest, allowing Docker clients and registries to understand how to handle the document correctly.",
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

func (c *Client) GetLatestImage(name string) (*Image, int, error) {

	image, code, err := c.GetResource(fmt.Sprintf("image/-latest/%s", name), new(Image))

	if err != nil {
		return nil, code, err
	}

	return image.(*Image), code, err
}

func (c *Client) GetImagesQuery(query Query) (*ImagesQueryResult, error) {

	// Marshal query into a JSON byte slice
	jsonData, jsonError := json.Marshal(query)

	if jsonError != nil {
		return nil, jsonError
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/org/%s/image/-query", c.HostURL, c.Org), bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, "application/json")

	if err != nil {
		return nil, err
	}

	images := ImagesQueryResult{}
	err = json.Unmarshal(body, &images)

	if err != nil {
		return nil, err
	}

	nextLink := GetLinkHref(images.Links, "next")

	for {

		if nextLink == nil {
			break
		}

		nextPage, _, err := c.Get(*nextLink, new(ImagesQueryResult))

		if err != nil {
			return nil, err
		}

		newQueryResult := nextPage.(*ImagesQueryResult)
		nextLink = GetLinkHref(newQueryResult.Links, "next")

		// Append items of next page
		images.Items = append(images.Items, newQueryResult.Items...)

		// Make links valid
		images.Links = newQueryResult.Links
	}

	return &images, nil
}
