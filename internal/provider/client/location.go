package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Locations
type Locations struct {
	Kind     string     `json:"kind,omitempty"`
	ItemKind string     `json:"itemKind,omitempty"`
	Items    []Location `json:"items,omitempty"`
	Links    []Link     `json:"links,omitempty"`
}

// Location
type Location struct {
	Base
	Provider *string         `json:"provider,omitempty"`
	Region   *string         `json:"region,omitempty"`
	Spec     *LocationSpec   `json:"spec,omitempty"`
	Status   *LocationStatus `json:"status,omitempty"`
}

// LocationSpec
type LocationSpec struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// LocationStatus
type LocationStatus struct {
	Geo      *LocationGeo `json:"geo,omitempty"`
	IpRanges *[]string    `json:"ipRanges,omitempty"`
}

// LocationStatus
type LocationGeo struct {
	Lat       *float32 `json:"lat,omitempty"`
	Lon       *float32 `json:"lon,omitempty"`
	State     *string  `json:"state,omitempty"`
	Country   *string  `json:"country,omitempty"`
	Continent *string  `json:"continent,omitempty"`
}

func LocationSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"cpln_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cloud_provider": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ip_ranges": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"self_link": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func LocationsSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"cpln_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cloud_provider": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ip_ranges": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"self_link": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// GetLocation
func (c *Client) GetLocation(name string) (*Location, int, error) {

	location, code, err := c.GetResource(fmt.Sprintf("location/%s", name), new(Location))

	if err != nil {
		return nil, code, err
	}

	return location.(*Location), code, err
}

// GetLocations
func (c *Client) GetLocations() (*Locations, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s/location", c.HostURL, c.Org), nil)

	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, "")

	if err != nil {
		return nil, err
	}

	locations := Locations{}
	err = json.Unmarshal(body, &locations)

	if err != nil {
		return nil, err
	}

	return &locations, nil
}
