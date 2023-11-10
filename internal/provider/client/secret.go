package cpln

import (
	"fmt"
)

// Secret - Secret
type Secret struct {
	Base
	Type 		*string      `json:"type,omitempty"`
	Data 	    *interface{} `json:"data,omitempty"`
	DataReplace *interface{} `json:"$replace/data,omitempty"`
}

// GetSecret - Get secret by name
func (c *Client) GetSecret(name string) (*Secret, int, error) {

	secret, code, err := c.GetResource(fmt.Sprintf("secret/%s/-reveal", name), new(Secret))

	if err != nil {
		return nil, code, err
	}

	return secret.(*Secret), code, err
}

// CreateSecret - Create a new Secret
func (c *Client) CreateSecret(secret Secret) (*Secret, int, error) {

	code, err := c.CreateResource("secret", *secret.Name, secret)
	if err != nil {
		return nil, code, err
	}

	return c.GetSecret(*secret.Name)
}

// UpdateSecret - Update an existing secret
func (c *Client) UpdateSecret(secret Secret) (*Secret, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("secret/%s", *secret.Name), secret)
	if err != nil {
		return nil, code, err
	}

	return c.GetSecret(*secret.Name)
}

// DeleteSecret - Delete secret by name
func (c *Client) DeleteSecret(name string) error {
	return c.DeleteResource(fmt.Sprintf("secret/%s", name))
}
