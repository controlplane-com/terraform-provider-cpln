package cpln

import (
	"fmt"
)

// Domain - Org Defined Domain Name
type Domain struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Tags        *map[string]interface{} `json:"tags,omitempty"`
	Links       *[]Link                 `json:"links,omitempty"`
}

// GetDomain - Get Domain by name
func (c *Client) GetDomain(name string) (*Domain, int, error) {

	domain, code, err := c.GetResource(fmt.Sprintf("domain/%s", name), new(Domain))

	if err != nil {
		return nil, code, err
	}

	return domain.(*Domain), code, err
}

// CreateDomain - Create a new Domain
func (c *Client) CreateDomain(domain Domain) (*Domain, int, error) {

	code, err := c.CreateResource("domain", *domain.Name, domain)
	if err != nil {
		return nil, code, err
	}

	return c.GetDomain(*domain.Name)
}

// UpdateDomain - Update an existing domain
func (c *Client) UpdateDomain(domain Domain) (*Domain, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)
	if err != nil {
		return nil, code, err
	}

	return c.GetDomain(*domain.Name)
}

// DeleteDomain - Delete domain by name
func (c *Client) DeleteDomain(name string) error {
	return c.DeleteResource(fmt.Sprintf("domain/%s", name))
}
