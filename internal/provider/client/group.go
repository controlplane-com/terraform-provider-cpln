package cpln

import (
	"fmt"
)

// Group - Control Plane Group
type Group struct {
	Base
	MemberLinks     *[]string        `json:"memberLinks,omitempty"`
	MemberQuery     *Query           `json:"memberQuery,omitempty"`
	IdentityMatcher *IdentityMatcher `json:"identityMatcher,omitempty"`
	Origin          *string          `json:"origin,omitempty"`
}

type IdentityMatcher struct {
	Expression *string `json:"expression,omitempty"`
	Language   *string `json:"language,omitempty"` // Enum: [ jmespath, javascript ]
}

// GetGroup - Get Group by name
func (c *Client) GetGroup(name string) (*Group, int, error) {

	group, code, err := c.GetResource(fmt.Sprintf("group/%s", name), new(Group))

	if err != nil {
		return nil, code, err
	}

	return group.(*Group), code, err
}

// CreateGroup - Create a new Group
func (c *Client) CreateGroup(group Group) (*Group, int, error) {

	code, err := c.CreateResource("group", *group.Name, group)
	if err != nil {
		return nil, code, err
	}

	return c.GetGroup(*group.Name)
}

// UpdateGroup - Update an existing Group
func (c *Client) UpdateGroup(group Group) (*Group, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("group/%s", *group.Name), group)
	if err != nil {
		return nil, code, err
	}

	return c.GetGroup(*group.Name)
}

// DeleteGroup - Delete Group by name
func (c *Client) DeleteGroup(name string) error {
	return c.DeleteResource(fmt.Sprintf("group/%s", name))
}
