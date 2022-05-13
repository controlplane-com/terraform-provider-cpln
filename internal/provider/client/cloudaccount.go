package cpln

import (
	"fmt"
)

// CloudAccount - CloudAccount
type CloudAccount struct {
	Base
	Provider *string             `json:"provider,omitempty"`
	Data     *CloudAccountConfig `json:"data,omitempty"`
	Status   *CloudAccountStatus `json:"status,omitempty"`
}

// CloudAccountConfig - CloudAccountConfig
type CloudAccountConfig struct {
	RoleArn    *string `json:"roleArn,omitempty"`
	ProjectId  *string `json:"projectId,omitempty"`
	SecretLink *string `json:"secretLink,omitempty"`
}

// CloudAccountStatus - CloudAccountStatus
type CloudAccountStatus struct {
	Usable      *bool   `json:"usable,omitempty"`
	LastChecked *string `json:"lastChecked,omitempty"`
	LastError   *string `json:"lastError,omitempty"`
}

// GetCloudAccount - Get CloudAccount by name
func (c *Client) GetCloudAccount(name string) (*CloudAccount, int, error) {

	cloudAccount, code, err := c.GetResource(fmt.Sprintf("cloudaccount/%s", name), new(CloudAccount))

	if err != nil {
		return nil, 0, err
	}

	return cloudAccount.(*CloudAccount), code, err
}

// CreateCloudAccount - Create an CloudAccount
func (c *Client) CreateCloudAccount(cloudaccount CloudAccount) (*CloudAccount, int, error) {

	code, err := c.CreateResource("cloudaccount", *cloudaccount.Name, cloudaccount)
	if err != nil {
		return nil, code, err
	}

	return c.GetCloudAccount(*cloudaccount.Name)
}

// UpdateCloudAccount - Update an CloudAccount
func (c *Client) UpdateCloudAccount(cloudaccount CloudAccount) (*CloudAccount, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("cloudaccount/%s", *cloudaccount.Name), cloudaccount)
	if err != nil {
		return nil, code, err
	}

	return c.GetCloudAccount(*cloudaccount.Name)
}

// DeleteCloudAccount - Delete CloudAccount by name
func (c *Client) DeleteCloudAccount(name string) error {
	return c.DeleteResource(fmt.Sprintf("cloudaccount/%s", name))
}
