package cpln

import (
	"fmt"
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
	DnsMode        *string           `json:"dnsMode,omitempty"` // Enum: "cname", "ns"
	GvcLink        *string           `json:"gvcLink,omitempty"`
	AcceptAllHosts *bool             `json:"acceptAllHosts,omitempty"`
	Ports          *[]DomainSpecPort `json:"ports,omitempty"`
}

type DomainStatus struct {
	Endpoints   *[]DomainEndpoint       `json:"endPoints,omitempty"`
	Status      *string                 `json:"status,omitempty"`
	Warning     *string                 `json:"warning,omitempty"`
	Locations   *[]DomainStatusLocation `json:"locations,omitempty"`
	Fingerprint *string                 `json:"fingerprint,omitempty"`
	DnsConfig   *[]DnsConfigRecord      `json:"dnsConfig,omitempty"`
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
	Prefix        *string             `json:"prefix,omitempty"`
	ReplacePrefix *string             `json:"replacePrefix,omitempty"`
	Regex         *string             `json:"regex,omitempty"`
	WorkloadLink  *string             `json:"workloadLink,omitempty"`
	Port          *int                `json:"port,omitempty"`
	HostPrefix    *string             `json:"hostPrefix,omitempty"`
	Headers       *DomainRouteHeaders `json:"headers,omitempty"`
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

type DomainRouteHeaders struct {
	Request *DomainHeaderOperation `json:"request,omitempty"`
}

type DomainHeaderOperation struct {
	Set *map[string]interface{} `json:"set,omitempty"`
}

/*** Status Related ***/
type DomainEndpoint struct {
	URL          *string `json:"url,omitempty"`
	WorkloadLink *string `json:"workloadLink,omitempty"`
}

type DomainStatusLocation struct {
	Name              *string `json:"name,omitempty"`
	CertificateStatus *string `json:"certificateStatus,omitempty"`
}

type DnsConfigRecord struct {
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

	time.Sleep(15 * time.Second)

	return c.GetDomain(*domain.Name)
}

// UpdateDomain - Update an existing domain
func (c *Client) UpdateDomain(domain Domain) (*Domain, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)
	if err != nil {
		return nil, code, err
	}

	time.Sleep(15 * time.Second)

	return c.GetDomain(*domain.Name)
}

// DeleteDomain - Delete domain by name
func (c *Client) DeleteDomain(name string) error {
	return c.DeleteResource(fmt.Sprintf("domain/%s", name))
}

/*** Domain Route ***/
func (c *Client) AddDomainRoute(domainName string, domainPort int, route DomainRoute) error {

	domain, _, err := c.GetDomain(domainName)

	if err != nil {
		return err
	}

	if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
		return fmt.Errorf("domain is not configured correctly, ports are not set")
	}

	for index, value := range *domain.Spec.Ports {

		if *value.Number == domainPort {

			// Append a new route
			if (*domain.Spec.Ports)[index].Routes == nil {
				(*domain.Spec.Ports)[index].Routes = &[]DomainRoute{}
			}

			if *route.Port == 0 {
				route.Port = nil
			}

			*(*domain.Spec.Ports)[index].Routes = append(*(*domain.Spec.Ports)[index].Routes, route)

			domain.SpecReplace = DeepCopy(domain.Spec).(*DomainSpec)
			domain.Spec = nil
			domain.Status = nil

			// Update resource
			_, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

			if err != nil {
				return err
			}

			return nil
		}
	}

	// Port not found, return an error
	routeIdentifier := ""

	if route.Prefix != nil {
		routeIdentifier = fmt.Sprintf("with prefix '%s'", *route.Prefix)
	}

	if route.Regex != nil {
		routeIdentifier = fmt.Sprintf("with regex '%s'", *route.Regex)
	}

	return fmt.Errorf("unable to add route %s for a domain named '%s'. Port '%d' is not set", routeIdentifier, domainName, domainPort)
}

func (c *Client) UpdateDomainRoute(domainName string, domainPort int, route *DomainRoute) error {

	domain, _, err := c.GetDomain(domainName)

	if err != nil {
		return err
	}

	if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
		return fmt.Errorf("Domain is not configured correctly, ports are not set")
	}

	for pIndex, value := range *domain.Spec.Ports {

		if *value.Number == domainPort && (value.Routes != nil && len(*value.Routes) > 0) {

			for rIndex, _route := range *value.Routes {

				if (_route.Prefix != nil && route.Prefix != nil && *_route.Prefix == *route.Prefix) ||
					(_route.Regex != nil && route.Regex != nil && *_route.Regex == *route.Regex) {

					// Modify existing route
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].ReplacePrefix = route.ReplacePrefix
					(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].WorkloadLink = route.WorkloadLink

					if route.Port == nil || *route.Port == 0 {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Port = nil
					} else {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Port = route.Port
					}

					if route.HostPrefix == nil || *route.HostPrefix == "" {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].HostPrefix = nil
					} else {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].HostPrefix = route.HostPrefix
					}

					if route.Headers == nil || route.Headers.Request == nil || route.Headers.Request.Set == nil {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Headers = nil
					} else {
						(*(*domain.Spec.Ports)[pIndex].Routes)[rIndex].Headers = route.Headers
					}

					// Update resource
					domain.SpecReplace = DeepCopy(domain.Spec).(*DomainSpec)
					domain.Spec = nil
					domain.Status = nil

					_, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

					if err != nil {
						return err
					}

					return nil
				}
			}
		}
	}

	// Port not found, return an error
	routeIdentifier := ""

	if route.Prefix != nil {
		routeIdentifier = fmt.Sprintf("with prefix '%s'", *route.Prefix)
	}

	if route.Regex != nil {
		routeIdentifier = fmt.Sprintf("with regex '%s'", *route.Regex)
	}

	return fmt.Errorf("unable to update route %s for a domain named '%s'. Port '%d' is not set", routeIdentifier, domainName, domainPort)
}

func (c *Client) RemoveDomainRoute(domainName string, domainPort int, prefix *string, regex *string) error {

	domain, _, err := c.GetDomain(domainName)

	if err != nil {
		return err
	}

	if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
		return fmt.Errorf("domain is not configured correctly, ports are not set")
	}

	routeIndex := -1

	for pIndex, value := range *domain.Spec.Ports {

		if *value.Number == domainPort && (value.Routes != nil && len(*value.Routes) > 0) {

			for _index, _route := range *value.Routes {

				if (prefix != nil && _route.Prefix != nil && *_route.Prefix == *prefix) ||
					(regex != nil && _route.Regex != nil && *_route.Regex == *regex) {
					routeIndex = _index
					break
				}
			}

			if routeIndex != -1 {

				// Remove route at index routeIndex
				*(*domain.Spec.Ports)[pIndex].Routes = append((*(*domain.Spec.Ports)[pIndex].Routes)[:routeIndex], (*(*domain.Spec.Ports)[pIndex].Routes)[routeIndex+1:]...)

				// Update resource
				domain.SpecReplace = DeepCopy(domain.Spec).(*DomainSpec)
				domain.Spec = nil
				domain.Status = nil

				_, err := c.UpdateResource(fmt.Sprintf("domain/%s", *domain.Name), domain)

				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	// Route not found, return an error
	routeIdentifier := ""

	if prefix != nil {
		routeIdentifier = fmt.Sprintf("with prefix '%s'", *prefix)
	}

	if regex != nil {
		routeIdentifier = fmt.Sprintf("with regex '%s'", *regex)
	}

	return fmt.Errorf("unable to delete route %s for a domain named '%s'. Route not found at port %d", routeIdentifier, domainName, domainPort)
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
