package cpln

import (
	"fmt"
)

// Identity - Identity
type Identity struct {
	Base
	Aws                    *IdentityAws                     `json:"aws,omitempty"`
	AwsReplace             *IdentityAws                     `json:"$replace/aws,omitempty"`
	Gcp                    *IdentityGcp                     `json:"gcp,omitempty"`
	GcpReplace             *IdentityGcp                     `json:"$replace/gcp,omitempty"`
	Azure                  *IdentityAzure                   `json:"azure,omitempty"`
	AzureReplace           *IdentityAzure                   `json:"$replace/azure,omitempty"`
	Ngs                    *IdentityNgs                     `json:"ngs,omitempty"`
	NgsReplace             *IdentityNgs                     `json:"$replace/ngs,omitempty"`
	NetworkResources       *[]IdentityNetworkResource       `json:"networkResources,omitempty"`
	NativeNetworkResources *[]IdentityNativeNetworkResource `json:"nativeNetworkResources,omitempty"`
	Status                 *IdentityStatus                  `json:"status,omitempty"`
	Drop                   *[]string                        `json:"$drop,omitempty"`
}

type IdentityAws struct {
	CloudAccountLink *string                 `json:"cloudAccountLink,omitempty"`
	PolicyRefs       *[]string               `json:"policyRefs,omitempty"`
	RoleName         *string                 `json:"roleName,omitempty"`
	TrustPolicy      *IdentityAwsTrustPolicy `json:"trustPolicy,omitempty"`
}

type IdentityAwsTrustPolicy struct {
	Version   *string                   `json:"Version,omitempty"`
	Statement *[]map[string]interface{} `json:"Statement,omitempty"`
}

type IdentityGcp struct {
	CloudAccountLink *string               `json:"cloudAccountLink,omitempty"`
	Scopes           *[]string             `json:"scopes,omitempty"`
	ServiceAccount   *string               `json:"serviceAccount,omitempty"`
	Bindings         *[]IdentityGcpBinding `json:"bindings,omitempty"`
}

type IdentityGcpBinding struct {
	Resource *string   `json:"resource,omitempty"`
	Roles    *[]string `json:"roles,omitempty"`
}

type IdentityAzure struct {
	CloudAccountLink *string                        `json:"cloudAccountLink,omitempty"`
	RoleAssignments  *[]IdentityAzureRoleAssignment `json:"roleAssignments,omitempty"`
}

type IdentityAzureRoleAssignment struct {
	Scope *string   `json:"scope,omitempty"`
	Roles *[]string `json:"roles,omitempty"`
}

type IdentityNgs struct {
	CloudAccountLink *string          `json:"cloudAccountLink,omitempty"`
	Pub              *IdentityNgsPerm `json:"pub,omitempty"`
	Sub              *IdentityNgsPerm `json:"sub,omitempty"`
	Resp             *IdentityNgsResp `json:"resp,omitempty"`
	Subs             *int             `json:"subs,omitempty"`
	Data             *int             `json:"data,omitempty"`
	Payload          *int             `json:"payload,omitempty"`
}

type IdentityNgsPerm struct {
	Allow *[]string `json:"allow,omitempty"`
	Deny  *[]string `json:"deny,omitempty"`
}

type IdentityNgsResp struct {
	Max *int    `json:"max,omitempty"`
	TTL *string `json:"ttl,omitempty"`
}

type IdentityNetworkResource struct {
	Name       *string   `json:"name,omitempty"`
	AgentLink  *string   `json:"agentLink,omitempty"`
	IPs        *[]string `json:"IPs,omitempty"`
	FQDN       *string   `json:"FQDN,omitempty"`
	ResolverIP *string   `json:"resolverIP,omitempty"`
	Ports      *[]int    `json:"ports,omitempty"`
}

type IdentityNativeNetworkResource struct {
	Name              *string                    `json:"name,omitempty"`
	FQDN              *string                    `json:"FQDN,omitempty"`
	Ports             *[]int                     `json:"ports,omitempty"`
	AWSPrivateLink    *IdentityAwsPrivateLink    `json:"awsPrivateLink,omitempty"`
	GCPServiceConnect *IdentityGcpServiceConnect `json:"gcpServiceConnect,omitempty"`
}

type IdentityAwsPrivateLink struct {
	EndpointServiceName *string `json:"endpointServiceName,omitempty"`
}

type IdentityGcpServiceConnect struct {
	TargetService *string `json:"targetService,omitempty"`
}

type IdentityStatus struct {
	ObjectName *string                 `json:"objectName,omitempty"`
	Aws        *IdentityStatusProvider `json:"aws,omitempty"`
	Gcp        *IdentityStatusProvider `json:"gcp,omitempty"`
	Azure      *IdentityStatusProvider `json:"azure,omitempty"`
}

type IdentityStatusProvider struct {
	LastError *string `json:"lastError,omitempty"`
	Usable    *bool   `json:"usable,omitempty"`
}

// GetIdentity - Get Identity by name
func (c *Client) GetIdentity(name, gvcName string) (*Identity, int, error) {

	identity, code, err := c.GetResource(fmt.Sprintf("gvc/%s/identity/%s", gvcName, name), new(Identity))
	if err != nil {
		return nil, code, err
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
