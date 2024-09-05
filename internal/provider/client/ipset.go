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
	IpAddresses *[]IpAddress `json:"ipAddresses,omitempty"`
	Error       *string      `json:"error,omitempty"`
}

type IpAddress struct {
	Name    *string `json:"name,omitempty"`
	Ip      *string `json:"ip,omitempty"`
	Id      *string `json:"id,omitempty"`
	State   *string `json:"state,omitempty"`
	Created *string `json:"created,omitempty"`
}

func (c *Client) GetIpSet(name string) (*IpSet, int, error) {

	ipSet, code, err := c.GetResource(fmt.Sprintf("ipset/%s", name), new(IpSet))

	if err != nil {
		return nil, code, err
	}

	return ipSet.(*IpSet), code, err
}

func (c *Client) CreateIpSet(ipSet IpSet) (*IpSet, int, error) {

	code, err := c.CreateResource("ipset", *ipSet.Name, ipSet)

	if err != nil {
		return nil, code, err
	}

	return c.GetIpSet(*ipSet.Name)
}

func (c *Client) UpdateIpSet(ipSet IpSet) (*IpSet, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("ipset/%s", *ipSet.Name), ipSet)

	if err != nil {
		return nil, code, err
	}

	return c.GetIpSet(*ipSet.Name)
}

func (c *Client) DeleteIpSet(name string) error {
	return c.DeleteResource(fmt.Sprintf("ipset/%s", name))
}
