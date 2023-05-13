package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ServiceAccount - Service Account
type ServiceAccount struct {
	Base
	Keys   *[]ServiceAccountKey `json:"keys,omitempty"`
	Origin *string              `json:"origin,omitempty"`
}

// ServiceAccountKey - Service Account Key
type ServiceAccountKey struct {
	Name        string  `json:"name,omitempty"`
	Created     *string `json:"created,omitempty"`
	Key         string  `json:"key,omitempty"`
	Description *string `json:"description,omitempty"`
}

// GetServiceAccount - Get Service Account by name
func (c *Client) GetServiceAccount(name string) (*ServiceAccount, int, error) {

	serviceAccount, code, err := c.GetResource(fmt.Sprintf("serviceaccount/%s", name), new(ServiceAccount))

	if err != nil {
		return nil, code, err
	}

	return serviceAccount.(*ServiceAccount), code, err
}

// CreateServiceAccount - Create a new Service Account
func (c *Client) CreateServiceAccount(serviceaccount ServiceAccount) (*ServiceAccount, int, error) {

	code, err := c.CreateResource("serviceaccount", *serviceaccount.Name, serviceaccount)
	if err != nil {
		return nil, code, err
	}

	return c.GetServiceAccount(*serviceaccount.Name)
}

// AddServiceAccountKey - Add Service Account Key
func (c *Client) AddServiceAccountKey(serviceAccountName, description string) (*ServiceAccountKey, error) {

	key := make(map[string]string)

	key["description"] = description

	s, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/org/%s/serviceaccount/%s/-addKey", c.HostURL, c.Org, serviceAccountName), strings.NewReader(string(s)))
	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, "application/json")
	if err != nil {
		return nil, err
	}

	saKey := ServiceAccountKey{}

	err = json.Unmarshal(body, &saKey)
	if err != nil {
		return nil, err
	}

	saKey.Name = saKey.Key[0:strings.Index(saKey.Key, ".")]

	return &saKey, nil
}

// RemoveServiceAccountKey = Remove Service Account Key
func (c *Client) RemoveServiceAccountKey(serviceAccountName, keyName string) error {

	removeKey := make(map[string][]string)

	removeKey["$drop/keys"] = []string{
		keyName,
	}

	s, err := json.Marshal(removeKey)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/org/%s/serviceaccount/%s", c.HostURL, c.Org, serviceAccountName), strings.NewReader(string(s)))
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, "application/json")
	if err != nil {
		return err
	}

	return nil
}

// UpdateServiceAccount - Update an existing ServiceAccount
func (c *Client) UpdateServiceAccount(serviceaccount ServiceAccount) (*ServiceAccount, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("serviceaccount/%s", *serviceaccount.Name), serviceaccount)
	if err != nil {
		return nil, code, err
	}

	return c.GetServiceAccount(*serviceaccount.Name)
}

// DeleteServiceAccount - Delete ServiceAccount by name
func (c *Client) DeleteServiceAccount(name string) error {
	return c.DeleteResource(fmt.Sprintf("serviceaccount/%s", name))
}
