package cpln

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

const MAX_ATTEMPTS = 10

// Domain - Org Defined Domain Name
type Domain struct {
	Base
	Links       *[]Link       `json:"links,omitempty"`
	Spec        *DomainSpec   `json:"spec,omitempty"`
	SpecReplace *DomainSpec   `json:"$replace/spec,omitempty"`
	Status      *DomainStatus `json:"status,omitempty"`
}

type DomainSpec struct {
	DnsMode             *string           `json:"dnsMode,omitempty"` // Enum: "cname", "ns"
	GvcLink             *string           `json:"gvcLink,omitempty"`
	CertChallengeType   *string           `json:"certChallengeType,omitempty"` // Enum: "http01", "dns01"
	WorkloadLink        *string           `json:"workloadLink,omitempty"`
	AcceptAllHosts      *bool             `json:"acceptAllHosts,omitempty"`
	AcceptAllSubdomains *bool             `json:"acceptAllSubdomains,omitempty"`
	Ports               *[]DomainSpecPort `json:"ports,omitempty"`
}

type DomainStatus struct {
	Endpoints   *[]DomainStatusEndpoint        `json:"endPoints,omitempty"`
	Status      *string                        `json:"status,omitempty"`
	Warning     *string                        `json:"warning,omitempty"`
	Locations   *[]DomainStatusLocation        `json:"locations,omitempty"`
	Fingerprint *string                        `json:"fingerprint,omitempty"`
	DnsConfig   *[]DomainStatusDnsConfigRecord `json:"dnsConfig,omitempty"`
}

/*** Spec Related ***/
type DomainSpecPort struct {
	Number   *int           `json:"number,omitempty"`
	Protocol *string        `json:"protocol,omitempty"` // Enum: "http", "http2", "tcp"
	Routes   *[]DomainRoute `json:"routes,omitempty"`
	Cors     *DomainCors    `json:"cors,omitempty"`
	TLS      *DomainTLS     `json:"tls,omitempty"`
}

type DomainRoute struct {
	Prefix        *string              `json:"prefix,omitempty"`
	ReplacePrefix *string              `json:"replacePrefix,omitempty"`
	Regex         *string              `json:"regex,omitempty"`
	WorkloadLink  *string              `json:"workloadLink,omitempty"`
	Port          *int                 `json:"port,omitempty"`
	HostPrefix    *string              `json:"hostPrefix,omitempty"`
	HostRegex     *string              `json:"hostRegex,omitempty"`
	Headers       *DomainRouteHeaders  `json:"headers,omitempty"`
	Replica       *int                 `json:"replica,omitempty"`
	Mirror        *[]DomainRouteMirror `json:"mirror,omitempty"`
	Canaries      *[]DomainRouteCanary `json:"canaries,omitempty"`
}

type DomainRouteMirror struct {
	WorkloadLink *string  `json:"workloadLink,omitempty"`
	Port         *int     `json:"port,omitempty"`
	Percent      *float64 `json:"percent,omitempty"`
}

type DomainRouteCanary struct {
	WorkloadLink *string `json:"workloadLink,omitempty"`
	Port         *int    `json:"port,omitempty"`
	Weight       *int    `json:"weight,omitempty"`
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
	MinProtocolVersion *string            `json:"minProtocolVersion,omitempty"` // Enum: "TLSV1_3", "TLSV1_2", "TLSV1_1", "TLSV1_0"
	CipherSuites       *[]string          `json:"cipherSuites,omitempty"`       // Enum: "ECDHE-ECDSA-AES256-GCM-SHA384", "ECDHE-ECDSA-CHACHA20-POLY1305", "ECDHE-ECDSA-AES128-GCM-SHA256", "ECDHE-RSA-AES256-GCM-SHA384", "ECDHE-RSA-CHACHA20-POLY1305", "ECDHE-RSA-AES128-GCM-SHA256", "AES256-GCM-SHA384", "AES128-GCM-SHA256"
	ClientCertificate  *DomainCertificate `json:"clientCertificate,omitempty"`
	ServerCertificate  *DomainCertificate `json:"serverCertificate,omitempty"`
}

type DomainAllowOrigin struct {
	Exact *string `json:"exact,omitempty"`
	Regex *string `json:"regex,omitempty"`
}

type DomainCertificate struct {
	SecretLink *string `json:"secretLink,omitempty"`
}

type DomainRouteHeaders struct {
	Request *DomainHeaderOperation `json:"request,omitempty"`
}

type DomainHeaderOperation struct {
	Set *map[string]interface{} `json:"set,omitempty"`
}

/*** Status Related ***/
type DomainStatusEndpoint struct {
	URL          *string `json:"url,omitempty"`
	WorkloadLink *string `json:"workloadLink,omitempty"`
}

type DomainStatusLocation struct {
	Name              *string `json:"name,omitempty"`
	CertificateStatus *string `json:"certificateStatus,omitempty"`
}

type DomainStatusDnsConfigRecord struct {
	Type  *string `json:"type,omitempty"`
	TTL   *int    `json:"ttl,omitempty"`
	Host  *string `json:"host,omitempty"`
	Value *string `json:"value,omitempty"`
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

// ErrDomainRouteNotFound - Reported when a domain port carries no route matching the requested identity
var ErrDomainRouteNotFound = errors.New("domain route not found")

// ErrDomainRouteConflict - Reported when an update would move a route onto an identity another route on the port already holds
var ErrDomainRouteConflict = errors.New("domain route identity already taken")

// IdentityKey - Build the key the API uses to identify a route within a port, its path paired with its host
func (r DomainRoute) IdentityKey() string {
	// Identify the route by its prefix or regex
	var path string

	switch {
	case r.Prefix != nil:
		path = "prefix:" + *r.Prefix
	case r.Regex != nil:
		path = "regex:" + *r.Regex
	default:
		return ""
	}

	// The API allows the same path on different hosts, so the host is part of the identity
	var host string

	switch {
	case r.HostPrefix != nil:
		host = "hostPrefix:" + *r.HostPrefix
	case r.HostRegex != nil:
		host = "hostRegex:" + *r.HostRegex
	}

	return path + ";" + host
}

// Identifier - Describe a route identity for error messages
func (r DomainRoute) Identifier() string {
	// Describe the route by its prefix or regex
	var path string

	switch {
	case r.Prefix != nil:
		path = fmt.Sprintf("with prefix '%s'", *r.Prefix)
	case r.Regex != nil:
		path = fmt.Sprintf("with regex '%s'", *r.Regex)
	}

	// Append the host matcher when the route narrows its path to one host
	switch {
	case r.HostPrefix != nil:
		return fmt.Sprintf("%s and host prefix '%s'", path, *r.HostPrefix)
	case r.HostRegex != nil:
		return fmt.Sprintf("%s and host regex '%s'", path, *r.HostRegex)
	}

	return path
}

// indexOfDomainRoute - Locate the route matching the given identity, returning -1 when the port carries no such route
func indexOfDomainRoute(routes []DomainRoute, identity DomainRoute) int {
	// An identity without a prefix or a regex addresses no route
	key := identity.IdentityKey()

	if key == "" {
		return -1
	}

	// Return the position of the first route sharing the identity key
	for index, route := range routes {
		if route.IdentityKey() == key {
			return index
		}
	}

	return -1
}

// AddDomainRoute - Append a route to a domain port
func (c *Client) AddDomainRoute(domainName string, domainPort int, route DomainRoute) (*DomainRoute, int, error) {

	const maxRetries = 5
	backoff := 2 * time.Second
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {

		domain, _, err := c.GetDomain(domainName)

		if err != nil {
			return nil, 0, err
		}

		if domain.Spec == nil || domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return nil, 0, fmt.Errorf("domain is not configured correctly, ports are not set")
		}

		shouldRetry := false

		for index, value := range *domain.Spec.Ports {

			if *value.Number == domainPort {

				// Append a new route
				if (*domain.Spec.Ports)[index].Routes == nil {
					(*domain.Spec.Ports)[index].Routes = &[]DomainRoute{}
				}

				if route.Port != nil && *route.Port == 0 {
					route.Port = nil
				}

				*(*domain.Spec.Ports)[index].Routes = append(*(*domain.Spec.Ports)[index].Routes, route)

				domain.SpecReplace = DeepCopy(domain.Spec).(*DomainSpec)
				domain.Spec = nil
				domain.Status = nil

				// Update resource
				code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

				if err != nil {
					if code == http.StatusConflict && attempt < maxRetries {
						lastErr = err
						time.Sleep(backoff)
						backoff *= 2
						shouldRetry = true
						break
					}
					return nil, 0, err
				}

				// If we got here then route has been added successfully
				return c.GetDomainRoute(domainName, domainPort, route)
			}
		}

		if shouldRetry {
			continue
		}

		// Port not found, return an error
		return nil, 0, fmt.Errorf("unable to add route %s for a domain named '%s'. Port '%d' is not set", route.Identifier(), domainName, domainPort)
	}

	return nil, 0, fmt.Errorf("add domain route failed after %d attempts due to HTTP 409: %w", maxRetries, lastErr)
}

// GetDomainRoute - Get the route matching the given identity at a domain port
func (c *Client) GetDomainRoute(domainName string, domainPort int, identity DomainRoute) (*DomainRoute, int, error) {
	domain, code, err := c.GetDomain(domainName)

	if err != nil {
		return nil, code, err
	}

	if domain.Spec == nil || domain.Spec.Ports == nil {
		return nil, code, err
	}

	for _, value := range *domain.Spec.Ports {
		if *value.Number == domainPort && (value.Routes != nil && len(*value.Routes) > 0) {
			// Return a copy of the route the identity addresses
			if index := indexOfDomainRoute(*value.Routes, identity); index != -1 {
				route := (*value.Routes)[index]
				return &route, code, nil
			}
		}
	}

	return nil, code, err
}

// UpdateDomainRoute - Rewrite the route matching the given identity, which may itself change the route's host
func (c *Client) UpdateDomainRoute(domainName string, domainPort int, identity DomainRoute, route *DomainRoute) (*DomainRoute, int, error) {

	const maxRetries = 5
	backoff := 2 * time.Second
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {

		domain, _, err := c.GetDomain(domainName)

		if err != nil {
			return nil, 0, err
		}

		if domain.Spec == nil || domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return nil, 0, fmt.Errorf("Domain is not configured correctly, ports are not set")
		}

		shouldRetry := false

		for pIndex, value := range *domain.Spec.Ports {

			if *value.Number == domainPort && (value.Routes != nil && len(*value.Routes) > 0) {
				// Locate the route the identity addresses, which holds the values the route had before this update
				rIndex := indexOfDomainRoute(*value.Routes, identity)

				if rIndex != -1 {
					// Refuse to move the route onto a path and host another route on the port already holds
					if conflict := indexOfDomainRoute(*value.Routes, *route); conflict != -1 && conflict != rIndex {
						return nil, 0, fmt.Errorf(
							"domain '%s' already carries a route %s at port %d, so the route %s cannot be moved onto it: %w",
							domainName, route.Identifier(), domainPort, identity.Identifier(), ErrDomainRouteConflict,
						)
					}

					// Modify existing route
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].ReplacePrefix = route.ReplacePrefix
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].WorkloadLink = route.WorkloadLink

					if route.Port == nil || *route.Port == 0 {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Port = nil
					} else {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Port = route.Port
					}

					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].HostPrefix = route.HostPrefix
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].HostRegex = route.HostRegex
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Headers = route.Headers
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Replica = route.Replica
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Mirror = route.Mirror
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Canaries = route.Canaries

					// Update resource
					domain.SpecReplace = DeepCopy(domain.Spec).(*DomainSpec)
					domain.Spec = nil
					domain.Status = nil

					code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

					if err != nil {
						if code == http.StatusConflict && attempt < maxRetries {
							lastErr = err
							time.Sleep(backoff)
							backoff *= 2
							shouldRetry = true
							break
						}
						return nil, 0, err
					}

					// The update may have moved the route to another host, so read it back at its new identity
					return c.GetDomainRoute(domainName, domainPort, *route)
				}
			}
		}

		if shouldRetry {
			continue
		}

		// Route not found, return an error
		return nil, 0, fmt.Errorf("unable to update route %s for a domain named '%s'. Route not found at port %d", identity.Identifier(), domainName, domainPort)
	}

	return nil, 0, fmt.Errorf("update domain route failed after %d attempts due to HTTP 409: %w", maxRetries, lastErr)
}

// RemoveDomainRoute - Delete the route matching the given identity from a domain port
func (c *Client) RemoveDomainRoute(domainName string, domainPort int, identity DomainRoute) error {

	const maxRetries = 5
	backoff := 2 * time.Second
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {

		domain, _, err := c.GetDomain(domainName)

		if err != nil {
			return err
		}

		if domain.Spec == nil || domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return fmt.Errorf("domain is not configured correctly, ports are not set")
		}

		shouldRetry := false
		routeIndex := -1

		for pIndex, value := range *domain.Spec.Ports {

			if *value.Number == domainPort && (value.Routes != nil && len(*value.Routes) > 0) {

				// Locate the route the identity addresses
				routeIndex = indexOfDomainRoute(*value.Routes, identity)

				if routeIndex != -1 {

					// Remove route at index routeIndex
					*(*domain.Spec.Ports)[pIndex].Routes = append((*(*domain.Spec.Ports)[pIndex].Routes)[:routeIndex], (*(*domain.Spec.Ports)[pIndex].Routes)[routeIndex+1:]...)

					// Update resource
					domain.SpecReplace = DeepCopy(domain.Spec).(*DomainSpec)
					domain.Spec = nil
					domain.Status = nil

					code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

					if err != nil {
						if code == http.StatusConflict && attempt < maxRetries {
							lastErr = err
							time.Sleep(backoff)
							backoff *= 2
							shouldRetry = true
							break
						}
						return err
					}

					return nil
				}
			}
		}

		if shouldRetry {
			continue
		}

		// Route not found, return an error
		return fmt.Errorf("unable to delete route %s for a domain named '%s' at port %d: %w", identity.Identifier(), domainName, domainPort, ErrDomainRouteNotFound)
	}

	return fmt.Errorf("remove domain route failed after %d attempts due to HTTP 409: %w", maxRetries, lastErr)
}

func DeepCopy(source interface{}) interface{} {

	sourceValue := reflect.ValueOf(source)

	if sourceValue.Kind() != reflect.Ptr || sourceValue.IsNil() {
		return nil
	}

	sourceType := reflect.TypeOf(source).Elem()
	dest := reflect.New(sourceType).Elem()

	for i := 0; i < sourceValue.Elem().NumField(); i++ {
		sourceFieldValue := sourceValue.Elem().Field(i)
		dest.Field(i).Set(sourceFieldValue)
	}

	return dest.Addr().Interface()
}
