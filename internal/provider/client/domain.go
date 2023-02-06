package cpln

import (
	"fmt"
)

// Domain - Org Defined Domain Name
type Domain struct {
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
func (c *Client) AddDomainRoute(domainName string, domainPort int, route DomainRoute) (*DomainRoute, error) {
	domain, _, err := c.GetDomain(domainName)

	if err != nil {
		return nil, err
	}

	// append new route to the domain's correct port
	found := false

	if domain.Spec.Ports != nil && len(*domain.Spec.Ports) > 0 {

		for i, port := range *domain.Spec.Ports {
			if *port.Number != domainPort {
				continue
			}

			found = true

			if (*domain.Spec.Ports)[i].Routes == nil {
				(*domain.Spec.Ports)[i].Routes = &[]DomainRoute{}
			}

			*(*domain.Spec.Ports)[i].Routes = append(*(*domain.Spec.Ports)[i].Routes, route)
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("Domain port was not found")
	}

	_, err = c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

func (c *Client) UpdateDomainRoute(domainName string, domainPort int, route *DomainRoute) (*DomainRoute, error) {
	domain, _, err := c.GetDomain(domainName)

	if err != nil {
		return nil, err
	}

	// update the values of the route in place
	found := false

	if domain.Spec.Ports != nil && len(*domain.Spec.Ports) > 0 {
		for i, port := range *domain.Spec.Ports {
			if *port.Number != domainPort {
				continue
			}

			for j, _route := range *port.Routes {
				if *_route.Prefix != *route.Prefix {
					continue
				}

				*(*(*domain.Spec.Ports)[i].Routes)[j].Prefix = *route.Prefix

				if route.ReplacePrefix != nil {
					*(*(*domain.Spec.Ports)[i].Routes)[j].ReplacePrefix = *route.ReplacePrefix
				}

				*(*(*domain.Spec.Ports)[i].Routes)[j].WorkloadLink = *route.WorkloadLink

				if route.Port != nil {
					*(*(*domain.Spec.Ports)[i].Routes)[j].Port = *route.Port
				}

				found = true
				break
			}

			if found {
				break
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("Domain route to update was not found")
	}

	_, err = c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)
	if err != nil {
		return nil, err
	}

	return route, nil
}

func (c *Client) RemoveDomainRoute(domainName string, domainPort int, prefix string) (bool, error) {
	domain, _, err := c.GetDomain(domainName)

	if err != nil {
		return false, err
	}

	// update the values of the route in place
	if domain.Spec.Ports != nil && len(*domain.Spec.Ports) > 0 {
		for i, port := range *domain.Spec.Ports {
			if *port.Number != domainPort {
				continue
			}

			routeIndex := -1

			for j, _route := range *port.Routes {
				if *_route.Prefix != prefix {
					continue
				}

				routeIndex = j
				break
			}

			if routeIndex == -1 {
				continue
			}

			*(*domain.Spec.Ports)[i].Routes = append((*(*domain.Spec.Ports)[i].Routes)[:routeIndex], (*(*domain.Spec.Ports)[i].Routes)[routeIndex+1:]...)
			break
		}
	}

	_, err = c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)
	if err != nil {
		return false, err
	}

	return true, nil
}
