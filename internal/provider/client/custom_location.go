package cpln

import (
	"fmt"
)

func (c *Client) GetCustomLocation(name string) (*Location, int, error) {

	location, code, err := c.GetResource(fmt.Sprintf("location/%s", name), new(Location))

	if err != nil {
		return nil, code, err
	}

	return location.(*Location), code, err
}

func (c *Client) CreateCustomLocation(location Location) (*Location, int, error) {

	code, err := c.CreateResource("location", *location.Name, location)
	if err != nil {
		return nil, code, err
	}

	return c.GetLocation(*location.Name)
}

func (c *Client) UpdateCustomLocation(location Location) (*Location, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("location/%s", *location.Name), location)
	if err != nil {
		return nil, code, err
	}

	return c.GetLocation(*location.Name)
}

func (c *Client) DeleteCustomLocation(name string) error {
	return c.DeleteResource(fmt.Sprintf("location/%s", name))
}
