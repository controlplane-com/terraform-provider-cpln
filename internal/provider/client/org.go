package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Org - Organization
type Org struct {
	Base
	Status *OrgStatus `json:"status,omitempty"`
}

// OrgStatus - Organization Status
type OrgStatus struct {
	AccountLink *string `json:"accountLink,omitempty"`
}

// GetOrg - Get Organziation By Name
func (c *Client) GetOrg(name string) (*Org, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s", c.HostURL, name), nil)

	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, "")
	if err != nil {
		return nil, err
	}

	org := Org{}
	err = json.Unmarshal(body, &org)
	if err != nil {
		return nil, err
	}

	return &org, nil
}
