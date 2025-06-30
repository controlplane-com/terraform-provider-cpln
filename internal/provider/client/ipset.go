package cpln

import "fmt"

type IpSet struct {
	Base
	Spec        *IpSetSpec   `json:"spec,omitempty"`
	SpecReplace *IpSetSpec   `json:"$replace/spec,omitempty"`
	Status      *IpSetStatus `json:"status,omitempty"`
}

type IpSetSpec struct {
	Link      *string          `json:"link,omitempty"`
	Locations *[]IpSetLocation `json:"locations,omitempty"`
}

type IpSetLocation struct {
	Name            *string `json:"name,omitempty"`
	RetentionPolicy *string `json:"retentionPolicy,omitempty"`
}

type IpSetStatus struct {
	IpAddresses *[]IpSetIpAddress `json:"ipAddresses,omitempty"`
	Error       *string           `json:"error,omitempty"`
	Warning     *string           `json:"warning,omitempty"`
}

type IpSetIpAddress struct {
	Name    *string `json:"name,omitempty"`
	IP      *string `json:"ip,omitempty"`
	ID      *string `json:"id,omitempty"`
	State   *string `json:"state,omitempty"`
	Created *string `json:"created,omitempty"`
}

// GetIpSet - Get IP Set by name
func (c *Client) GetIpSet(name string) (*IpSet, int, error) {

	ipSet, code, err := c.GetResource(fmt.Sprintf("ipset/%s", name), new(IpSet))

	if err != nil {
		return nil, code, err
	}

	return ipSet.(*IpSet), code, err
}

// CreateIpSet - Create a new IP Set
func (c *Client) CreateIpSet(ipSet IpSet) (*IpSet, int, error) {

	code, err := c.CreateResource("ipset", *ipSet.Name, ipSet)

	if err != nil {
		return nil, code, err
	}

	return c.GetIpSet(*ipSet.Name)
}

// UpdateIpSet - Update an existing IP Set
func (c *Client) UpdateIpSet(ipSet IpSet) (*IpSet, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("ipset/%s", *ipSet.Name), ipSet)

	if err != nil {
		return nil, code, err
	}

	return c.GetIpSet(*ipSet.Name)
}

// DeleteIpSet - Delete IP Set by name
func (c *Client) DeleteIpSet(name string) error {
	return c.DeleteResource(fmt.Sprintf("ipset/%s", name))
}
