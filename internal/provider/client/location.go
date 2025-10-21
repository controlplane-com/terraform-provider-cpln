package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Origin   *string         `json:"origin,omitempty"`
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

func (c *Client) CreateCustomLocation(location Location) (*Location, int, error) {

	code, err := c.CreateResource("location", *location.Name, location)
	if err != nil {
		return nil, code, err
	}

	return c.GetLocation(*location.Name)
}

func (c *Client) UpdateLocation(location Location) (*Location, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("location/%s", *location.Name), location)
	if err != nil {
		return nil, code, err
	}

	return c.GetLocation(*location.Name)
}

// UpdateLocationToDefault patches the specified location to its default state.
func (c *Client) UpdateLocationToDefault(location Location) (*Location, int, error) {
	// Remove the Terraform-managed tag before sending the update
	c.RemoveManagedByTerraformTag(&location.Base)

	// Marshal the location struct into JSON for the HTTP request body
	payload, err := json.Marshal(location)

	// Return immediately if marshaling fails
	if err != nil {
		return nil, 0, err
	}

	// Build the PATCH request targeting the location endpoint
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/org/%s/location/%s", c.HostURL, c.Org, *location.Name),
		strings.NewReader(string(payload)),
	)

	// Return if request creation fails
	if err != nil {
		return nil, 0, err
	}

	// Execute the HTTP request and capture its status code
	_, code, err := c.doRequest(req, "application/json")

	// Propagate any errors from the request
	if err != nil {
		return nil, code, err
	}
	// Retrieve and return the updated location from the API
	return c.GetLocation(*location.Name)
}

func (c *Client) DeleteCustomLocation(name string) error {
	return c.DeleteResource(fmt.Sprintf("location/%s", name))
}
