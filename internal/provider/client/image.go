package cpln

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

type ImagesQueryResult struct {
	Kind  string  `json:"kind,omitempty"`
	Items []Image `json:"items,omitempty"`
	Links []Link  `json:"links,omitempty"`
	Query Query   `json:"query,omitempty"`
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

	if len(images.Items) == 0 {

		// Convert query to JSON
		queryJson, _ := json.Marshal(query)

		// Format error
		err = fmt.Errorf("no images were found")

		if query.Spec != nil && query.Spec.Terms != nil {

			// Format if query was modified
			if len(*query.Spec.Terms) > 0 || *query.Spec.Match != "all" {
				err = fmt.Errorf("no images were found following the query: %s", string(queryJson))
			}

			// Format if query was looking for a specific repository
			if len(*query.Spec.Terms) == 1 && (*query.Spec.Terms)[0].Property != nil && *(*query.Spec.Terms)[0].Property == "repository" && (*query.Spec.Terms)[0].Value != nil {
				err = fmt.Errorf("image base name %s does not exist, make sure you have spelled the base name correctly", *(*query.Spec.Terms)[0].Value)
			}
		}

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
