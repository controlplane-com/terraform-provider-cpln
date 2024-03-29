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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO validate domain name
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     DescriptionDomainValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "cname",
						},
						"gvc_link": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"accept_all_hosts": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ports": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"number": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  443,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "http2",
									},
									"cors": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allow_origins": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"exact": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
												"allow_methods": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														// TODO Disregard uppercase lowercase
													},
												},
												"allow_headers": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														// TODO Disregard uppercase lowercase
													},
												},
												"expose_headers": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														// TODO Disregard uppercase lowercase
													},
												},
												"max_age": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "24h",
												},
												"allow_credentials": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
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
													Type:     schema.TypeString,
													Optional: true,
													Default:  "TLSV1_2",
												},
												"cipher_suites": {
													Type:     schema.TypeSet,
													Optional: true,
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
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem:     certificateResource(),
												},
												"server_certificate": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem:     certificateResource(),
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
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"workload_link": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"warning": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"locations": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"certificate_status": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dns_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"host": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
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
				Type:     schema.TypeString,
				Optional: true,
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
