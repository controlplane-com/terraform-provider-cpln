package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Gvcs - GVC's
type Gvcs struct {
	Kind  string `json:"kind,omitempty"`
	Items []Gvc  `json:"items,omitempty"`
	Links []Link `json:"links,omitempty"`
	Query Query  `json:"query,omitempty"`
}

// Gvc - Global Virtual Cloud
type Gvc struct {
	Base
	Spec        *GvcSpec `json:"spec,omitempty"`
	SpecReplace *GvcSpec `json:"$replace/spec,omitempty"`
	Alias       *string  `json:"alias,omitempty"`
}

// GvcSpec - GVC Spec
type GvcSpec struct {
	StaticPlacement *StaticPlacement `json:"staticPlacement,omitempty"`
	PullSecretLinks *[]string        `json:"pullSecretLinks,omitempty"`
	Domain          *string          `json:"domain,omitempty"`
	Tracing         *Tracing         `json:"tracing,omitempty"`
	Env             *[]NameValue     `json:"env,omitempty"`
	LoadBalancer    *LoadBalancer    `json:"loadBalancer,omitempty"`
}

// StaticPlacement - Static Placement
type StaticPlacement struct {
	LocationLinks *[]string `json:"locationLinks,omitempty"`
	LocationQuery *Query    `json:"locationQuery,omitempty"`
}

// LoadBalancer - Load Balancer
type LoadBalancer struct {
	Dedicated *bool `json:"dedicated,omitempty"`
}

// GetGvcs - Get All Gvcs
func (c *Client) GetGvcs() (*Gvcs, error) {

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/org/%s/gvc/-query", c.HostURL, c.Org), nil)

	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, "")

	if err != nil {
		return nil, err
	}

	gvcs := Gvcs{}
	err = json.Unmarshal(body, &gvcs)

	if err != nil {
		return nil, err
	}

	return &gvcs, nil
}

// GetGvc - Get GVC by name
func (c *Client) GetGvc(name string) (*Gvc, int, error) {

	gvc, code, err := c.GetResource(fmt.Sprintf("gvc/%s", name), new(Gvc))

	if err != nil {
		return nil, code, err
	}

	return gvc.(*Gvc), code, err
}

// CreateGvc - Create a new GVC
func (c *Client) CreateGvc(gvc Gvc) (*Gvc, int, error) {

	code, err := c.CreateResource("gvc", *gvc.Name, gvc)
	if err != nil {
		return nil, code, err
	}

	return c.GetGvc(*gvc.Name)
}

// UpdateGvc - Update an existing GVC
func (c *Client) UpdateGvc(gvc Gvc) (*Gvc, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("gvc/%s", *gvc.Name), gvc)
	if err != nil {
		return nil, code, err
	}

	return c.GetGvc(*gvc.Name)
}

// DeleteGvc - Delete GVC by name
func (c *Client) DeleteGvc(name string) error {
	return c.DeleteResource(fmt.Sprintf("gvc/%s", name))
}
