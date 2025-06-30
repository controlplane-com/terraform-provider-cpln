package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Discovery holds Firebase and service endpoint mappings
type Discovery struct {
	Firebase  map[string]string `json:"firebase,omitempty"`
	Endpoints map[string]string `json:"endpoints,omitempty"`
}

// GetDiscovery fetches discovery information and returns the parsed object, HTTP status code, and any error.
func (c *Client) GetDiscovery() (*Discovery, int, error) {
	// Build HTTP GET request for the discovery endpoint
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/discovery", c.HostURL), nil)

	// Return early if request creation fails
	if err != nil {
		return nil, 0, err
	}

	// Perform the HTTP request and capture response body, status code, and any error
	body, code, err := c.doRequest(req, "")

	// propagate errors from the request execution
	if err != nil {
		return nil, code, err
	}

	// Initialize a Discovery struct to hold the unmarshaled JSON
	discovery := Discovery{}

	// Unmarshal the response body into the discovery struct
	err = json.Unmarshal(body, &discovery)

	// Propagate JSON parsing errors
	if err != nil {
		return nil, code, err
	}

	// Return the discovery data, HTTP status code, and nil error
	return &discovery, code, err
}

// GetBillingNgEndpoint retrieves the billing-ng service endpoint from discovery.
func (c *Client) GetBillingNgEndpoint() (string, int, error) {
	// Call GetDiscovery to obtain the discovery information
	discovery, code, err := c.GetDiscovery()

	// Propagate errors from discovery retrieval
	if err != nil {
		return "", code, err
	}

	// Return the billing-ng endpoint URL along with status code and error
	return discovery.Endpoints["billing-ng"], code, err
}
