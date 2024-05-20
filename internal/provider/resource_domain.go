package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Main ***/
func resourceDomain() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the Domain.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Domain name. (e.g., `example.com` / `test.example.com`). Control Plane will validate the existence of the domain with DNS. Create and Update will fail if the required DNS entries cannot be validated.",
				Required:    true,
				ForceNew:    true,
				// TODO validate domain name
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the domain name.",
				Optional:         true,
				ValidateFunc:     DescriptionDomainValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:        schema.TypeMap,
				Description: "Key-value map of resource tags.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"spec": {
				Type:        schema.TypeList,
				Description: "Domain specificiation.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_mode": {
							Type:        schema.TypeString,
							Description: "In `cname` dnsMode, Control Plane will configure workloads to accept traffic for the domain but will not manage DNS records for the domain. End users must configure CNAME records in their own DNS pointed to the canonical workload endpoint. Currently `cname` dnsMode requires that a TLS server certificate be configured when subdomain based routing is used. In `ns` dnsMode, Control Plane will manage the subdomains and create all necessary DNS records. End users configure NS records to forward DNS requests to the Control Plane managed DNS servers. Valid values: `cname`, `ns`. Default: `cname`.",
							Optional:    true,
							Default:     "cname",
						},
						"gvc_link": {
							Type:        schema.TypeString,
							Description: "This value is set to a target GVC (using a full link) for use by subdomain based routing. Each workload in the GVC will receive a subdomain in the form ${workload.name}.${domain.name}. **Do not include if path based routing is used.**",
							Optional:    true,
						},
						"accept_all_hosts": {
							Type:        schema.TypeBool,
							Description: "Allows domain to accept wildcards. The associated GVC must have dedicated load balancing enabled.",
							Optional:    true,
							Default:     false,
						},
						"ports": {
							Type:        schema.TypeList,
							Description: "Domain port specifications.",
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"number": {
										Type:        schema.TypeInt,
										Description: "Port to expose externally. Values: `80`, `443`. Default: `443`.",
										Optional:    true,
										Default:     443,
									},
									"protocol": {
										Type:        schema.TypeString,
										Description: "Allowed protocol. Valid values: `http`, `http2`, `tcp`. Default: `http2`.",
										Optional:    true,
										Default:     "http2",
									},
									"cors": {
										Type:        schema.TypeList,
										Description: "A security feature implemented by web browsers to allow resources on a web page to be requested from another domain outside the domain from which the resource originated.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allow_origins": {
													Type:        schema.TypeList,
													Description: "Determines which origins are allowed to access a particular resource on a server from a web browser.",
													Optional:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"exact": {
																Type:        schema.TypeString,
																Description: "Value of allowed origin.",
																Required:    true,
															},
														},
													},
												},
												"allow_methods": {
													Type:        schema.TypeSet,
													Description: "Specifies the HTTP methods (such as `GET`, `POST`, `PUT`, `DELETE`, etc.) that are allowed for a cross-origin request to a specific resource.",
													Optional:    true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														// TODO Disregard uppercase lowercase
													},
												},
												"allow_headers": {
													Type:        schema.TypeSet,
													Description: "Specifies the custom HTTP headers that are allowed in a cross-origin request to a specific resource.",
													Optional:    true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														// TODO Disregard uppercase lowercase
													},
												},
												"expose_headers": {
													Type:        schema.TypeSet,
													Description: "The HTTP headers that a server allows to be exposed to the client in response to a cross-origin request. These headers provide additional information about the server's capabilities or requirements, aiding in proper handling of the request by the client's browser or application.",
													Optional:    true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														// TODO Disregard uppercase lowercase
													},
												},
												"max_age": {
													Type:        schema.TypeString,
													Description: "Maximum amount of time that a preflight request result can be cached by the client browser. Input is expected as a duration string (i.e, 24h, 20m, etc.).",
													Optional:    true,
													Default:     "24h",
												},
												"allow_credentials": {
													Type:        schema.TypeBool,
													Description: "Determines whether the client-side code (typically running in a web browser) is allowed to include credentials (such as cookies, HTTP authentication, or client-side SSL certificates) in cross-origin requests.",
													Optional:    true,
													Default:     false,
												},
											},
										},
									},
									"tls": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min_protocol_version": {
													Type:        schema.TypeString,
													Description: "Minimum TLS version to accept. Minimum is `1.0`. Default: `1.2`.",
													Optional:    true,
													Default:     "TLSV1_2",
												},
												"cipher_suites": {
													Type:        schema.TypeSet,
													Description: "Allowed cipher suites. Refer to the [Domain Reference](https://docs.controlplane.com/reference/domain#cipher-suites) for details.",
													Optional:    true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													DefaultFunc: func() (interface{}, error) {
														return []string{
															"AES128-GCM-SHA256",
															"AES256-GCM-SHA384",
															"ECDHE-ECDSA-AES128-GCM-SHA256",
															"ECDHE-ECDSA-AES256-GCM-SHA384",
															"ECDHE-ECDSA-CHACHA20-POLY1305",
															"ECDHE-RSA-AES128-GCM-SHA256",
															"ECDHE-RSA-AES256-GCM-SHA384",
															"ECDHE-RSA-CHACHA20-POLY1305",
														}, nil
													},
												},
												"client_certificate": {
													Type:        schema.TypeList,
													Description: "The certificate authority PEM, stored as a TLS Secret, used to verify the authority of the client certificate. The only verification performed checks that the CN of the PEM matches the Domain (i.e., CN=*.DOMAIN).",
													Optional:    true,
													MaxItems:    1,
													Elem:        certificateResource(),
												},
												"server_certificate": {
													Type:        schema.TypeList,
													Description: "Custom Server Certificate.",
													Optional:    true,
													MaxItems:    1,
													Elem:        certificateResource(),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoints": {
							Type:        schema.TypeList,
							Description: "List of configured domain endpoints.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Description: "URL of endpoint.",
										Optional:    true,
									},
									"workload_link": {
										Type:        schema.TypeString,
										Description: "Full link to associated workload.",
										Optional:    true,
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of Domain. Possible values: `initializing`, `ready`, `pendingDnsConfig`, `pendingCertificate`, `usedByGvc`.",
							Optional:    true,
						},
						"warning": {
							Type:        schema.TypeString,
							Description: "Warning message.",
							Optional:    true,
						},
						"locations": {
							Type:        schema.TypeList,
							Description: "Contains the cloud provider name, region, and certificate status.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "The name of the location.",
										Optional:    true,
									},
									"certificate_status": {
										Type:        schema.TypeString,
										Description: "The current validity or status of the SSL/TLS certificate.",
										Optional:    true,
									},
								},
							},
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dns_config": {
							Type:        schema.TypeList,
							Description: "List of required DNS record entries.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "The DNS record type specifies the type of data the DNS record contains. Valid values: `CNAME`, `NS`, `TXT`.",
										Optional:    true,
									},
									"ttl": {
										Type:        schema.TypeInt,
										Description: "Time to live (TTL) is a value that signifies how long (in seconds) a DNS record should be cached by a resolver or a browser before a new request should be sent to refresh the data. Lower TTL values mean records are updated more frequently, which is beneficial for dynamic DNS configurations or during DNS migrations. Higher TTL values reduce the load on DNS servers and improve the speed of name resolution for end users by relying on cached data.",
										Optional:    true,
									},
									"host": {
										Type:        schema.TypeString,
										Description: "The host in DNS terminology refers to the domain or subdomain that the DNS record is associated with. It's essentially the name that is being queried or managed. For example, in a DNS record for `www.example.com`, `www` is a host in the domain `example.com`.",
										Optional:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "The value of a DNS record contains the data the record is meant to convey, based on the type of the record.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domain := client.Domain{
		Name:        GetString(d.Get("name")),
		Description: GetString(d.Get("description")),
		Tags:        GetStringMap(d.Get("tags")),
		Spec:        buildDomainSpec(d.Get("spec")),
	}

	c := m.(*client.Client)

	newDomain, code, err := c.CreateDomain(domain)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setDomain(d, newDomain)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	domain, code, err := c.GetDomain(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setDomain(d, domain)
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "spec") {

		domainToUpdate := client.Domain{
			Name: GetString(d.Get("name")),
		}

		if d.HasChange("description") {
			domainToUpdate.Description = GetDescriptionString(d.Get("description"), *domainToUpdate.Name)
		}

		if d.HasChange("tags") {
			domainToUpdate.Tags = GetTagChanges(d)
		}

		c := m.(*client.Client)

		if d.HasChange("spec") {

			domain, _, err := c.GetDomain(*domainToUpdate.Name)

			if err != nil {
				return diag.FromErr(err)
			}

			domainToUpdate.SpecReplace = buildDomainSpec(d.Get("spec"))

			// For loop is to restore routes created using domain_routes
			for uIndex, uValue := range *domainToUpdate.SpecReplace.Ports {

				for dIndex, dValue := range *domain.Spec.Ports {

					if *uValue.Number == *dValue.Number && ((*domain.Spec.Ports)[dIndex].Routes != nil && len(*(*domain.Spec.Ports)[dIndex].Routes) > 0) {

						destination := make([]client.DomainRoute, len(*(*domain.Spec.Ports)[dIndex].Routes))
						copy(destination, (*(*domain.Spec.Ports)[dIndex].Routes))
						(*domainToUpdate.SpecReplace.Ports)[uIndex].Routes = &destination
					}
				}
			}
		}

		// Apply update

		updatedDomain, _, err := c.UpdateDomain(domainToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setDomain(d, updatedDomain)
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	id := d.Id()
	err := c.DeleteDomain(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

/*** Resources ***/
func certificateResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"secret_link": {
				Type:        schema.TypeString,
				Description: "Full link to a TLS secret.",
				Optional:    true,
			},
		},
	}
}

/*** Build Functions ***/
// Spec Related //
func buildDomainSpec(input interface{}) *client.DomainSpec {

	specs := input.([]interface{})

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	result := &client.DomainSpec{}

	if spec["dns_mode"] != nil {
		result.DnsMode = GetString(spec["dns_mode"])
	}

	if spec["gvc_link"] != nil {
		result.GvcLink = GetString(spec["gvc_link"])
	}

	if spec["accept_all_hosts"] != nil {
		result.AcceptAllHosts = GetBool(spec["accept_all_hosts"])
	}

	if spec["ports"] != nil {
		result.Ports = buildSpecPorts(spec["ports"].([]interface{}))
	}

	return result
}

func buildSpecPorts(specs []interface{}) *[]client.DomainSpecPort {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []client.DomainSpecPort{}

	for _, item := range specs {
		port := item.(map[string]interface{})
		newPort := client.DomainSpecPort{}

		if port["number"] != nil {
			newPort.Number = GetInt(port["number"])
		}

		if port["protocol"] != nil {
			newPort.Protocol = GetString(port["protocol"])
		}

		if port["cors"] != nil {
			newPort.Cors = buildCors(port["cors"].([]interface{}))
		}

		if port["tls"] != nil {
			newPort.TLS = buildTLS(port["tls"].([]interface{}))
		}

		collection = append(collection, newPort)
	}

	return &collection
}

func buildCors(specs []interface{}) *client.DomainCors {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	result := &client.DomainCors{}

	if spec["allow_origins"] != nil {
		result.AllowOrigins = buildAllowOrigins(spec["allow_origins"].([]interface{}))
	}

	if spec["allow_methods"] != nil {
		result.AllowMethods = buildStringArray(spec["allow_methods"].(*schema.Set).List())
	}

	if spec["allow_headers"] != nil {
		result.AllowHeaders = buildStringArray(spec["allow_headers"].(*schema.Set).List())
	}

	if spec["expose_headers"] != nil {
		result.ExposeHeaders = buildStringArray(spec["expose_headers"].(*schema.Set).List())
	}

	if spec["max_age"] != nil {
		result.MaxAge = GetString(spec["max_age"])
	}

	if spec["allow_credentials"] != nil {
		result.AllowCredentials = GetBool(spec["allow_credentials"])
	}

	return result
}

func buildTLS(specs []interface{}) *client.DomainTLS {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	result := &client.DomainTLS{}

	if spec["min_protocol_version"] != nil {
		result.MinProtocolVersion = GetString(spec["min_protocol_version"])
	}

	if spec["cipher_suites"] != nil {
		result.CipherSuites = buildStringArray(spec["cipher_suites"].(*schema.Set).List())
	}

	if spec["client_certificate"] != nil {
		result.ClientCertificate = buildCertificate(spec["client_certificate"].([]interface{}))
	}

	if spec["server_certificate"] != nil {
		result.ServerCertificate = buildCertificate(spec["server_certificate"].([]interface{}))
	}

	return result
}

func buildAllowOrigins(specs []interface{}) *[]client.DomainAllowOrigin {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []client.DomainAllowOrigin{}

	for _, item := range specs {
		allowOrigin := item.(map[string]interface{})
		newAllowOrigin := client.DomainAllowOrigin{}
		if allowOrigin["exact"] != nil {
			newAllowOrigin.Exact = GetString(allowOrigin["exact"].(string))
		}
		collection = append(collection, newAllowOrigin)
	}

	return &collection
}

func buildCertificate(specs []interface{}) *client.DomainCertificate {

	if len(specs) == 0 {
		return nil
	}

	spec := specs[0]
	result := client.DomainCertificate{}

	if spec == nil {
		return &result
	}

	specAsMap := spec.(map[string]interface{})

	if specAsMap["secret_link"] != nil {
		result.SecretLink = GetString(specAsMap["secret_link"].(string))
	}

	return &result
}

/*** Flatten Functions ***/
// Spec Related //

func flattenDomainSpec(domainSpec *client.DomainSpec) []interface{} {

	if domainSpec == nil {
		return nil
	}

	spec := make(map[string]interface{})
	if domainSpec.DnsMode != nil {
		spec["dns_mode"] = *domainSpec.DnsMode
	}

	if domainSpec.GvcLink != nil {
		spec["gvc_link"] = *domainSpec.GvcLink
	}

	if domainSpec.AcceptAllHosts != nil {
		spec["accept_all_hosts"] = *domainSpec.AcceptAllHosts
	}
	spec["ports"] = flattenSpecPorts(domainSpec.Ports)

	return []interface{}{
		spec,
	}
}

func flattenSpecPorts(ports *[]client.DomainSpecPort) []interface{} {

	if ports == nil || len(*ports) == 0 {
		return nil
	}

	collection := make([]interface{}, len(*ports))
	for i, item := range *ports {

		port := make(map[string]interface{})
		if item.Number != nil {
			port["number"] = *item.Number
		}

		if item.Protocol != nil {
			port["protocol"] = *item.Protocol
		}

		port["cors"] = flattenCors(item.Cors)
		port["tls"] = flattenTLS(item.TLS)

		collection[i] = port
	}

	return collection
}

func flattenCors(cors *client.DomainCors) []interface{} {

	if cors == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["allow_origins"] = flattenAllowOrigins(cors.AllowOrigins)
	result["allow_methods"] = flattenStringsArray(cors.AllowMethods)
	result["allow_headers"] = flattenStringsArray(cors.AllowHeaders)
	result["expose_headers"] = flattenStringsArray(cors.ExposeHeaders)

	if cors.MaxAge != nil {
		result["max_age"] = *cors.MaxAge
	}

	if cors.AllowCredentials != nil {
		result["allow_credentials"] = *cors.AllowCredentials
	}

	return []interface{}{
		result,
	}
}

func flattenTLS(tls *client.DomainTLS) []interface{} {

	if tls == nil {
		return nil
	}

	result := make(map[string]interface{})
	if tls.MinProtocolVersion != nil {
		result["min_protocol_version"] = *tls.MinProtocolVersion
	}

	result["cipher_suites"] = flattenStringsArray(tls.CipherSuites)
	result["client_certificate"] = flattenCertificate(tls.ClientCertificate)
	result["server_certificate"] = flattenCertificate(tls.ServerCertificate)

	return []interface{}{
		result,
	}
}

func flattenAllowOrigins(allowOrigins *[]client.DomainAllowOrigin) []interface{} {

	if allowOrigins == nil || len(*allowOrigins) == 0 {
		return nil
	}

	collection := make([]interface{}, len(*allowOrigins))

	for i, item := range *allowOrigins {

		allowOrigin := make(map[string]interface{})
		if item.Exact != nil {
			allowOrigin["exact"] = *item.Exact
		}

		collection[i] = allowOrigin
	}

	return collection
}

func flattenCertificate(certificate *client.DomainCertificate) []interface{} {

	if certificate == nil {
		return nil
	}

	result := make(map[string]interface{})
	if certificate.SecretLink != nil {
		result["secret_link"] = *certificate.SecretLink
	}

	return []interface{}{
		result,
	}
}

func flattenDomainStatus(status *client.DomainStatus) []interface{} {
	if status == nil {
		return nil
	}

	spec := make(map[string]interface{})

	endpoints := flattenDomainStatusEndpoints(status.Endpoints)
	if endpoints != nil {
		spec["endpoints"] = endpoints
	}

	if status.Status != nil {
		spec["status"] = *status.Status
	}

	if status.Warning != nil {
		spec["warning"] = *status.Warning
	}

	locations := flattenDomainStatusLocations(status.Locations)
	if locations != nil {
		spec["locations"] = locations
	}

	if status.Fingerprint != nil {
		spec["fingerprint"] = *status.Fingerprint
	}

	dnsConfig := flattenDomainStatusDnsConfig(status.DnsConfig)
	if dnsConfig != nil {
		spec["dns_config"] = dnsConfig
	}

	return []interface{}{
		spec,
	}
}

func flattenDomainStatusEndpoints(endpoints *[]client.DomainEndpoint) []interface{} {
	if endpoints == nil || len(*endpoints) == 0 {
		return nil
	}

	specs := []interface{}{}

	for _, endpoint := range *endpoints {
		if endpoint.URL == nil && endpoint.WorkloadLink == nil {
			continue
		}

		spec := make(map[string]interface{})

		if endpoint.URL != nil {
			spec["url"] = *endpoint.URL
		}

		if endpoint.WorkloadLink != nil {
			spec["workload_link"] = *endpoint.WorkloadLink
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenDomainStatusLocations(locations *[]client.DomainStatusLocation) []interface{} {
	if locations == nil || len(*locations) == 0 {
		return nil
	}

	specs := []interface{}{}

	for _, location := range *locations {
		if location.Name == nil && location.CertificateStatus == nil {
			continue
		}

		spec := make(map[string]interface{})

		if location.Name != nil {
			spec["name"] = *location.Name
		}

		if location.CertificateStatus != nil {
			spec["certificate_status"] = *location.CertificateStatus
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenDomainStatusDnsConfig(configs *[]client.DnsConfigRecord) []interface{} {
	if configs == nil || len(*configs) == 0 {
		return nil
	}

	specs := []interface{}{}

	for _, config := range *configs {
		if config.Type == nil && config.TTL == nil && config.Host == nil && config.Value == nil {
			continue
		}

		spec := make(map[string]interface{})

		if config.Type != nil {
			spec["type"] = *config.Type
		}

		if config.TTL != nil {
			spec["ttl"] = *config.TTL
		}

		if config.Host != nil {
			spec["host"] = *config.Host
		}

		if config.Value != nil {
			spec["value"] = *config.Value
		}

		specs = append(specs, spec)
	}

	return specs
}

/*** Helper Functions ***/
func setDomain(d *schema.ResourceData, domain *client.Domain) diag.Diagnostics {

	if domain == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*domain.Name)

	if err := d.Set("name", domain.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", domain.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", domain.Tags); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(domain.Links, d); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("spec", flattenDomainSpec(domain.Spec)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", flattenDomainStatus(domain.Status)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Build //
func buildStringArray(specs []interface{}) *[]string {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []string{}
	for _, item := range specs {
		collection = append(collection, item.(string))
	}

	return &collection
}

// Flatten //
func flattenStringsArray(strings *[]string) []interface{} {

	if strings == nil || len(*strings) == 0 {
		return nil
	}

	collection := make([]interface{}, len(*strings))
	for i, item := range *strings {
		collection[i] = item
	}

	return collection
}
