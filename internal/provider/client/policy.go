package cpln

import (
	"fmt"
)

// Policy - Policy
type Policy struct {
	Base
	TargetKind  *string    `json:"targetKind,omitempty"`
	TargetLinks *[]string  `json:"targetLinks,omitempty"`
	TargetQuery *Query     `json:"targetQuery,omitempty"`
	Target      *string    `json:"target,omitempty"`
	Origin      *string    `json:"origin,omitempty"`
	Bindings    *[]Binding `json:"bindings"`
	Update      bool       `json:"-"`
}

type PolicyUpdate struct {
	Base
	TargetKind  *string    `json:"targetKind,omitempty"`
	TargetLinks *[]string  `json:"targetLinks,omitempty"`
	TargetQuery *Query     `json:"targetQuery"`
	Target      *string    `json:"target"`
	Origin      *string    `json:"origin,omitempty"`
	Bindings    *[]Binding `json:"bindings"`
}

type Binding struct {
	Permissions    *[]string `json:"permissions,omitempty"`
	PrincipalLinks *[]string `json:"principalLinks,omitempty"`
}

// GetPolicy - Get Policy by name
func (c *Client) GetPolicy(name string) (*Policy, int, error) {

	policy, code, err := c.GetResource(fmt.Sprintf("policy/%s", name), new(Policy))

	if err != nil {
		return nil, code, err
	}

	return policy.(*Policy), code, err
}

// CreatePolicy - Create an Policy
func (c *Client) CreatePolicy(policy Policy) (*Policy, int, error) {

	code, err := c.CreateResource("policy", *policy.Name, policy)
	if err != nil {
		return nil, code, err
	}

	return c.GetPolicy(*policy.Name)
}

// UpdatePolicy - Update an Policy
func (c *Client) UpdatePolicy(policy PolicyUpdate) (*Policy, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("policy/%s", *policy.Name), policy)
	if err != nil {
		return nil, code, err
	}

	return c.GetPolicy(*policy.Name)
}

// DeletePolicy - Delete Policy by name
func (c *Client) DeletePolicy(name string) error {
	return c.DeleteResource(fmt.Sprintf("policy/%s", name))
}
