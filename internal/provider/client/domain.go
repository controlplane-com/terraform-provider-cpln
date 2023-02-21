package cpln

import (
	"fmt"
)

const MAX_ATTEMPTS = 10

// Domain - Org Defined Domain Name
type Domain struct {
	Base
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Tags        *map[string]interface{} `json:"tags,omitempty"`
	Links       *[]Link                 `json:"links,omitempty"`
	Spec        *DomainSpec             `json:"spec,omitempty"`
	Status      *DomainStatus           `json:"status,omitempty"`
}

type DomainSpec struct {
	DnsMode        *string           `json:"dnsMode,omitempty"` // Enum: "cname", "ns"
	GvcLink        *string           `json:"gvcLink,omitempty"`
	AcceptAllHosts *bool             `json:"acceptAllHosts,omitempty"`
	Ports          *[]DomainSpecPort `json:"ports,omitempty"`
}

type DomainStatus struct {
	EndPoints   *[]DomainEndPoint       `json:"endPoints,omitempty"`
	Status      *string                 `json:"status,omitempty"` // Enum: "initializing", "ready", "pendingDnsConfig", "pendingCertificate", "usedByGvc"
	Warning     *string                 `json:"warning,omitempty"`
	Locations   *[]DomainStatusLocation `json:"locations,omitempty"`
	Fingerprint *string                 `json:"fingerprint,omitempty"`
}

/*** Spec Related ***/
type DomainSpecPort struct {
	Number   *int           `json:"number,omitempty"`
	Protocol *string        `json:"protocol,omitempty"` // Enum: "http", "http2"
	Routes   *[]DomainRoute `json:"routes,omitempty"`
	Cors     *DomainCors    `json:"cors,omitempty"`
	TLS      *DomainTLS     `json:"tls,omitempty"`
}

type DomainRoute struct {
	Prefix        *string `json:"prefix,omitempty"`
	ReplacePrefix *string `json:"replacePrefix,omitempty"`
	WorkloadLink  *string `json:"workloadLink,omitempty"`
	Port          *int    `json:"port,omitempty"`
}

type DomainCors struct {
	AllowOrigins     *[]DomainAllowOrigin `json:"allowOrigins,omitempty"`
	AllowMethods     *[]string            `json:"allowMethods,omitempty"`
	AllowHeaders     *[]string            `json:"allowHeaders,omitempty"`
	ExposeHeaders    *[]string            `json:"exposeHeaders,omitempty"`
	MaxAge           *string              `json:"maxAge,omitempty"`
	AllowCredentials *bool                `json:"allowCredentials,omitempty"`
}

type DomainTLS struct {
	MinProtocolVersion *string            `json:"minProtocolVersion,omitempty"` // Enum: "TLSV1_2", "TLSV1_1", "TLSV1_0"
	CipherSuites       *[]string          `json:"cipherSuites,omitempty"`       // Enum: "ECDHE-ECDSA-AES256-GCM-SHA384", "ECDHE-ECDSA-CHACHA20-POLY1305", "ECDHE-ECDSA-AES128-GCM-SHA256", "ECDHE-RSA-AES256-GCM-SHA384", "ECDHE-RSA-CHACHA20-POLY1305", "ECDHE-RSA-AES128-GCM-SHA256", "AES256-GCM-SHA384", "AES128-GCM-SHA256"
	ClientCertificate  *DomainCertificate `json:"clientCertificate,omitempty"`
	ServerCertificate  *DomainCertificate `json:"serverCertificate,omitempty"`
}

type DomainAllowOrigin struct {
	Exact *string `json:"exact,omitempty"`
}

type DomainCertificate struct {
	SecretLink *string `json:"secretLink,omitempty"`
}

/*** Status Related ***/
type DomainEndPoint struct {
	URL          *string `json:"url,omitempty"`
	WorkloadLink *string `json:"workloadLink,omitempty"`
}

type DomainStatusLocation struct {
	Name              *string `json:"name,omitempty"`
	CertificateStatus *string `json:"certificateStatus,omitempty"` // Enum: "initializing", "ready", "pendingDnsConfig", "pendingCertificate "
}

// GetDomain - Get Domain by name
func (c *Client) GetDomain(name string) (*Domain, int, error) {

	domain, code, err := c.GetResource(fmt.Sprintf("domain/%s", name), new(Domain))

	if err != nil {
		return nil, code, err
	}

	return domain.(*Domain), code, err
}

// CreateDomain - Create a new Domain
func (c *Client) CreateDomain(domain Domain) (*Domain, int, error) {

	code, err := c.CreateResource("domain", *domain.Name, domain)
	if err != nil {
		return nil, code, err
	}

	return c.GetDomain(*domain.Name)
}

// UpdateDomain - Update an existing domain
func (c *Client) UpdateDomain(domain Domain) (*Domain, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)
	if err != nil {
		return nil, code, err
	}

	return c.GetDomain(*domain.Name)
}

// DeleteDomain - Delete domain by name
func (c *Client) DeleteDomain(name string) error {
	return c.DeleteResource(fmt.Sprintf("domain/%s", name))
}

/*** Domain Route ***/
func (c *Client) AddDomainRoute(domainName string, route DomainRoute) (*DomainRoute, error) {

	for {
		domain, _, err := c.GetDomain(domainName)

		if err != nil {
			return nil, err
		}

		if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return nil, fmt.Errorf("Domain is not configured correctly, ports are not set")
		}

		// Append a new route
		if (*domain.Spec.Ports)[0].Routes == nil {
			(*domain.Spec.Ports)[0].Routes = &[]DomainRoute{}
		}

		*(*domain.Spec.Ports)[0].Routes = append(*(*domain.Spec.Ports)[0].Routes, route)

		// Update resource
		code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

		if code == 409 {
			continue
		}

		if err != nil {
			return nil, err
		}

		break
	}

	return &route, nil
}

func (c *Client) UpdateDomainRoute(domainName string, route *DomainRoute) (*DomainRoute, error) {

	for {
		domain, _, err := c.GetDomain(domainName)

		if err != nil {
			return nil, err
		}

		if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return nil, fmt.Errorf("Domain is not configured correctly, ports are not set")
		}

		found := false

		// Update given route
		for j, _route := range *(*domain.Spec.Ports)[0].Routes {
			if *_route.Prefix != *route.Prefix {
				continue
			}

			*(*(*domain.Spec.Ports)[0].Routes)[j].Prefix = *route.Prefix
			*(*(*domain.Spec.Ports)[0].Routes)[j].WorkloadLink = *route.WorkloadLink

			if route.ReplacePrefix != nil {
				*(*(*domain.Spec.Ports)[0].Routes)[j].ReplacePrefix = *route.ReplacePrefix
			}

			if route.Port != nil {
				*(*(*domain.Spec.Ports)[0].Routes)[j].Port = *route.Port
			}

			found = true
			break
		}

		if !found {
			return nil, fmt.Errorf("Domain route specified was not found")
		}

		// Update resource
		code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

		if code == 409 {
			continue
		}

		if err != nil {
			return nil, err
		}

		break
	}

	return route, nil
}

func (c *Client) RemoveDomainRoute(domainName string, prefix string) (bool, error) {

	for {
		domain, _, err := c.GetDomain(domainName)

		if err != nil {
			return false, err
		}

		if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return false, fmt.Errorf("Domain is not configured correctly, ports are not set")
		}

		// Remove route
		routeIndex := -1

		for j, _route := range *(*domain.Spec.Ports)[0].Routes {
			if *_route.Prefix != prefix {
				continue
			}

			routeIndex = j
			break
		}

		if routeIndex == -1 {
			continue
		}

		*(*domain.Spec.Ports)[0].Routes = append((*(*domain.Spec.Ports)[0].Routes)[:routeIndex], (*(*domain.Spec.Ports)[0].Routes)[routeIndex+1:]...)

		// Update resource
		code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

		if code == 409 {
			continue
		}

		if err != nil {
			return false, err
		}

		break
	}

	return true, nil
}
