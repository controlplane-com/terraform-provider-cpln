package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Discovery struct {
	Firebase  map[string]string `json:"firebase,omitempty"`
	Endpoints map[string]string `json:"endpoints,omitempty"`
}

func (c *Client) GetDiscovery() (*Discovery, int, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/discovery", c.HostURL), nil)

	if err != nil {
		return nil, 0, err
	}

	body, code, err := c.doRequest(req, "")
	if err != nil {
		return nil, code, err
	}

	discovery := Discovery{}
	err = json.Unmarshal(body, &discovery)
	if err != nil {
		return nil, code, err
	}

	return &discovery, code, err
}

func (c *Client) GetBillingNgEndpoint() (string, int, error) {
	discovery, code, err := c.GetDiscovery()

	if err != nil {
		return "", code, err
	}

	return discovery.Endpoints["billing-ng"], code, err
}
