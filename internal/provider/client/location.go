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
	Country   *string  `json:"country,omitempty"`
	State     *string  `json:"state,omitempty"`
	City      *string  `json:"city,omitempty"`
	Continent *string  `json:"continent,omitempty"`
}

func LocationSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"cpln_id": {
			Type:        schema.TypeString,
			Description: "The ID, in GUID format, of the location.",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the location.",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description of the location.",
			Computed:    true,
		},
		"tags": {
			Type:        schema.TypeMap,
			Description: "Key-value map of resource tags.",
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cloud_provider": {
			Type:        schema.TypeString,
			Description: "Cloud Provider of the location.",
			Computed:    true,
		},
		"region": {
			Type:        schema.TypeString,
			Description: "Region of the location.",
			Computed:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indication if location is enabled.",
			Computed:    true,
		},
		"geo": GeoSchema(),
		"ip_ranges": {
			Type:        schema.TypeSet,
			Description: "A list of IP ranges of the location.",
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
	}
}

func LocationsSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"cpln_id": {
			Type:        schema.TypeString,
			Description: "The ID, in GUID format, of the location.",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the location.",
			Computed:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description of the location.",
			Computed:    true,
		},
		"tags": {
			Type:        schema.TypeMap,
			Description: "Key-value map of resource tags.",
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cloud_provider": {
			Type:        schema.TypeString,
			Description: "Cloud Provider of the location.",
			Computed:    true,
		},
		"region": {
			Type:        schema.TypeString,
			Description: "Region of the location.",
			Computed:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indication if location is enabled.",
			Computed:    true,
		},
		"geo": GeoSchema(),
		"ip_ranges": {
			Type:        schema.TypeSet,
			Description: "A list of IP ranges of the location.",
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
	}
}

func GeoSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"lat": {
					Type:        schema.TypeFloat,
					Description: "Latitude.",
					Optional:    true,
				},
				"lon": {
					Type:        schema.TypeFloat,
					Description: "Longitude.",
					Optional:    true,
				},
				"country": {
					Type:        schema.TypeString,
					Description: "Country.",
					Optional:    true,
				},
				"state": {
					Type:        schema.TypeString,
					Description: "State.",
					Optional:    true,
				},
				"city": {
					Type:        schema.TypeString,
					Description: "City.",
					Optional:    true,
				},
				"continent": {
					Type:        schema.TypeString,
					Description: "Continent.",
					Optional:    true,
				},
			},
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

func (c *Client) UpdateLocation(location Location) (*Location, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("location/%s", *location.Name), location)
	if err != nil {
		return nil, code, err
	}

	return c.GetLocation(*location.Name)
}
