package cpln

import (
	"context"
	"time"

	client "terraform-provider-cpln/internal/provider/client"

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//TODO: ValidateFunc: NameValidator,
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
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"gvc_link": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"accept_all_hosts": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"number": {
										Type:     schema.TypeInt, // Float instead?
										Optional: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"routes": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prefix": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"replace_prefix": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"workload_link": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
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
														// Optional: true,
													},
												},
												"allow_headers": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														// Optional: true,
													},
												},
												"max_age": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"allow_credentials": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"tls": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min_protocol_version": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"cipher_suites": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
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
				Type: schema.TypeList,
				// MaxItems: 1,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domain := client.Domain{
		Name:        GetString(d.Get("name").(string)),
		Description: GetString(d.Get("description").(string)),
		Tags:        GetStringMap(d.Get("tags")),
		Spec:        buildDomainSpec(d.Get("spec").([]interface{})),
	}

	c := m.(*client.Client)
	count := 0

	for {

		newDomain, code, err := c.CreateDomain(domain)

		if code == 409 {
			return ResourceExistsHelper()
		}

		if code != 400 {

			if err != nil {
				return diag.FromErr(err)
			}

			return setDomain(d, newDomain)
		}

		if count++; count > 16 {
			// Exit loop after timeout

			var diags diag.Diagnostics

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to verify domain ownership",
				Detail:   "Please review and run terraform apply again",
			})

			return diags
		}

		time.Sleep(15 * time.Second)
	}
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	domain, code, err := c.GetDomain(d.Id())

	if code == 404 {
		return setGvc(d, nil, c.Org)
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

		// Check for changes
		if d.HasChange("description") {
			domainToUpdate.Description = GetDescriptionString(d.Get("description"), *domainToUpdate.Name)
		}

		if d.HasChange("tags") {
			domainToUpdate.Tags = GetTagChanges(d)
		}

		if d.HasChange("spec") {
			if domainToUpdate.Spec == nil {
				domainToUpdate.Spec = &client.DomainSpec{}
			}
			domainToUpdate.Spec = buildDomainSpec(d.Get("spec").([]interface{}))
		}

		// Apply update
		c := m.(*client.Client)
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
	err := c.DeleteDomain(d.Id())
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
func buildDomainSpec(specs []interface{}) *client.DomainSpec {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	result := &client.DomainSpec{}

	if spec["dns_mode"] != nil {
		result.DnsMode = GetString(spec["dns_mode"].(string))
	}

	if spec["gvc_link"] != nil {
		result.GvcLink = GetString(spec["gvc_link"].(string))
	}

	if spec["accept_all_hosts"] != nil {
		result.AcceptAllHosts = GetBool(spec["accept_all_hosts"].(bool))
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
			newPort.Number = GetInt(port["number"].(int))
		}

		if port["protocol"] != nil {
			newPort.Protocol = GetString(port["protocol"].(string))
		}

		if port["routes"] != nil {
			newPort.Routes = buildRoutes(port["routes"].([]interface{}))
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

func buildRoutes(specs []interface{}) *[]client.DomainRoute {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []client.DomainRoute{}
	for _, item := range specs {
		route := item.(map[string]interface{})
		newRoute := client.DomainRoute{}

		if route["prefix"] != nil {
			newRoute.Prefix = GetString(route["prefix"].(string))
		}

		if route["replace_prefix"] != nil {
			newRoute.ReplacePrefix = GetString(route["replace_prefix"].(string))
		}

		if route["workload_link"] != nil {
			newRoute.WorkloadLink = GetString(route["workload_link"].(string))
		}

		if route["port"] != nil {
			newRoute.Port = GetInt(route["port"].(int))
		}

		collection = append(collection, newRoute)
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
		result.AllowMethods = buildStringArray(spec["allow_methods"].([]interface{}))
	}

	if spec["allow_headers"] != nil {
		result.AllowHeaders = buildStringArray(spec["allow_headers"].([]interface{}))
	}

	if spec["max_age"] != nil {
		result.MaxAge = GetString(spec["max_age"].(string))
	}

	if spec["allow_credentials"] != nil {
		result.AllowCredentials = GetBool(spec["allow_credentials"].(bool))
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
		result.MinProtocolVersion = GetString(spec["min_protocol_version"].(string))
	}

	if spec["cipher_suites"] != nil {
		result.CipherSuites = buildStringArray(spec["cipher_suites"].([]interface{}))
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
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	result := client.DomainCertificate{}

	if spec["secret_link"] != nil {
		result.SecretLink = GetString(spec["secret_link"].(string))
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

		port["routes"] = flattenRoutes(item.Routes)
		port["cors"] = flattenCors(item.Cors)
		port["tls"] = flattenTLS(item.TLS)

		collection[i] = port
	}

	return collection
}

func flattenRoutes(routes *[]client.DomainRoute) []interface{} {
	if routes == nil || len(*routes) == 0 {
		return nil
	}

	collection := make([]interface{}, len(*routes))
	for i, item := range *routes {

		route := make(map[string]interface{})
		if item.Prefix != nil {
			route["prefix"] = *item.Prefix
		}

		if item.ReplacePrefix != nil {
			route["replace_prefix"] = *item.ReplacePrefix
		}

		if item.WorkloadLink != nil {
			route["workload_link"] = *item.WorkloadLink
		}

		if item.Port != nil {
			route["port"] = *item.Port
		}

		collection[i] = route
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
