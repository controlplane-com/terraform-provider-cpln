package cpln

import (
	"fmt"
)

// Identity - Identity
type Identity struct {
	Base
	Aws              *AwsIdentity       `json:"aws,omitempty"`
	AwsReplace       *AwsIdentity       `json:"$replace/aws,omitempty"`
	Gcp              *GcpIdentity       `json:"gcp,omitempty"`
	GcpReplace       *GcpIdentity       `json:"$replace/gcp,omitempty"`
	Azure            *AzureIdentity     `json:"azure,omitempty"`
	AzureReplace     *AzureIdentity     `json:"$replace/azure,omitempty"`
	NetworkResources *[]NetworkResource `json:"networkResources,omitempty"`
	Status           *IdentityStatus    `json:"status,omitempty"`
	Drop             *[]string          `json:"$drop,omitempty"`
}

type IdentityStatus struct {
	ObjectName *string `json:"objectName,omitempty"`
}

// NetworkResource - NetworkResource
type NetworkResource struct {
	Name       *string   `json:"name,omitempty"`
	AgentLink  *string   `json:"agentLink,omitempty"`
	IPs        *[]string `json:"IPs,omitempty"`
	FQDN       *string   `json:"FQDN,omitempty"`
	ResolverIP *string   `json:"resolverIP,omitempty"`
	Ports      *[]int    `json:"ports,omitempty"`
}

// AwsIdentity - AwsIdentity
type AwsIdentity struct {
	CloudAccountLink *string   `json:"cloudAccountLink,omitempty"`
	PolicyRefs       *[]string `json:"policyRefs,omitempty"`
	// TrustPolicy      *AwsPolicyDocument `json:"trustPolicy,omitempty"`
	RoleName *string `json:"roleName,omitempty"`
}

// // AwsPolicyDocument - AwsPolicyDocument
// type AwsPolicyDocument struct {
// 	Version   *string   `json:"version,omitempty"`
// 	Statement *[]string `json:"statement,omitempty"`
// }

type GcpBinding struct {
	Resource *string   `json:"resource,omitempty"`
	Roles    *[]string `json:"roles,omitempty"`
}

type GcpIdentity struct {
	CloudAccountLink *string       `json:"cloudAccountLink,omitempty"`
	Scopes           *[]string     `json:"scopes,omitempty"`
	ServiceAccount   *string       `json:"serviceAccount,omitempty"`
	Bindings         *[]GcpBinding `json:"bindings,omitempty"`
}

type AzureRoleAssignment struct {
	Scope *string   `json:"scope,omitempty"`
	Roles *[]string `json:"roles,omitempty"`
}

type AzureIdentity struct {
	CloudAccountLink *string                `json:"cloudAccountLink,omitempty"`
	RoleAssignments  *[]AzureRoleAssignment `json:"roleAssignments,omitempty"`
}

// GetIdentity - Get Identity by name
func (c *Client) GetIdentity(name, gvcName string) (*Identity, int, error) {

	identity, code, err := c.GetResource(fmt.Sprintf("gvc/%s/identity/%s", gvcName, name), new(Identity))
	if err != nil {
		return nil, 0, err
	}

	return identity.(*Identity), code, err
}

// CreateIdentity - Create an Identity
func (c *Client) CreateIdentity(identity Identity, gvcName string) (*Identity, int, error) {

	code, err := c.CreateResource(fmt.Sprintf("gvc/%s/identity", gvcName), *identity.Name, identity)
	if err != nil {
		return nil, code, err
	}

	return c.GetIdentity(*identity.Name, gvcName)
}

// UpdateIdentity - Update an Identity
func (c *Client) UpdateIdentity(identity Identity, gvcName string) (*Identity, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("gvc/%s/identity/%s", gvcName, *identity.Name), identity)
	if err != nil {
		return nil, code, err
	}

	return c.GetIdentity(*identity.Name, gvcName)
}

// DeleteIdentity - Delete Identity by name
func (c *Client) DeleteIdentity(name, gvcName string) error {
	return c.DeleteResource(fmt.Sprintf("gvc/%s/identity/%s", gvcName, name))
}
