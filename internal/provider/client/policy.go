package cpln

import (
	"encoding/json"
	"fmt"
)

// Policy - Policy
type Policy struct {
	Base
	TargetKind  *string    `json:"targetKind,omitempty"`
	TargetLinks *[]string  `json:"targetLinks"`
	TargetQuery *Query     `json:"targetQuery,omitempty"`
	Target      *string    `json:"target,omitempty"`
	Origin      *string    `json:"origin,omitempty"`
	Bindings    *[]Binding `json:"bindings"`
	Update      bool       `json:"-"`
}

type PolicyUpdate struct {
	Base
	TargetKind  *string    `json:"targetKind,omitempty"`
	TargetLinks *[]string  `json:"targetLinks"`
	TargetQuery *Query     `json:"targetQuery"`
	Target      *string    `json:"target"`
	Origin      *string    `json:"origin,omitempty"`
	Bindings    *[]Binding `json:"bindings"`
}

func (p Policy) MarshalJSON() ([]byte, error) {

	type localPolicy Policy

	if p.Update && (p.Target == nil || *p.Target == "") {
		return json.Marshal(PolicyUpdate{
			Base:        p.Base,
			TargetKind:  p.TargetKind,
			TargetLinks: p.TargetLinks,
			TargetQuery: p.TargetQuery,
			Target:      p.Target,
			Origin:      p.Origin,
			Bindings:    p.Bindings,
		})
	}
	return json.Marshal(localPolicy(p))
}

// Binding - Binding
type Binding struct {
	Permissions    *[]string `json:"permissions,omitempty"`
	PrincipalLinks *[]string `json:"principalLinks,omitempty"`
}

// GetPolicy - Get Policy by name
func (c *Client) GetPolicy(name string) (*Policy, int, error) {

	policy, code, err := c.GetResource(fmt.Sprintf("policy/%s", name), new(Policy))

	if err != nil {
		return nil, 0, err
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
func (c *Client) UpdatePolicy(policy Policy) (*Policy, int, error) {

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
