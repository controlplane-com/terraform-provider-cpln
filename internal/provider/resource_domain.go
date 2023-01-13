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
										Elem: map[string]*schema.Schema{
											"allow_origins": {
												Type:     schema.TypeList,
												Optional: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"exact": {
															Type: schema.TypeString,
														},
													},
												},
											},
											"allow_methods": {
												Type:     schema.TypeSet,
												Optional: true,
												Elem: &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
											},
											"allow_headers": {
												Type:     schema.TypeSet,
												Optional: true,
												Elem: &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_points": {
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
					},
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
		Status:      buildDomainStatus(d.Get("status").([]interface{})),
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

	if d.HasChanges("description", "tags", "spec", "status") {

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

		if d.HasChange("status") {
			if domainToUpdate.Status == nil {
				domainToUpdate.Status = &client.DomainStatus{}
			}
			domainToUpdate.Status = buildDomainStatus(d.Get("status").([]interface{}))
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
				Type: schema.TypeString,
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
	return &client.DomainSpec{
		DnsMode:        GetString(spec["dns_mode"].(string)),
		GvcLink:        GetString(spec["gvc_link"].(string)),
		AcceptAllHosts: GetBool(spec["accept_all_hosts"].(bool)),
		Ports:          buildSpecPorts(spec["ports"].([]interface{})),
	}
}

func buildSpecPorts(specs []interface{}) *[]client.DomainSpecPort {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []client.DomainSpecPort{}
	for _, item := range specs {
		port := item.(map[string]interface{})
		collection = append(collection, client.DomainSpecPort{
			Number:   GetInt(port["number"].(int)),
			Protocol: GetString(port["protocol"].(string)),
			Routes:   buildRoutes(port["routes"].([]interface{})),
			Cors:     buildCors(port["cors"].([]interface{})),
			TLS:      buildTLS(port["tls"].([]interface{})),
		})
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
		collection = append(collection, client.DomainRoute{
			Prefix:        GetString(route["prefix"].(string)),
			ReplacePrefix: GetString(route["replace_prefix"].(string)),
			WorkloadLink:  GetString(route["workload_link"].(string)),
			Port:          GetInt(route["port"].(int)),
		})
	}

	return &collection
}

func buildCors(specs []interface{}) *client.DomainCors {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	return &client.DomainCors{
		AllowOrigins:     buildAllowOrigins(spec["allow_origins"].([]interface{})),
		AllowMethods:     buildStringArray(spec["allow_methods"].([]interface{})),
		AllowHeaders:     buildStringArray(spec["allow_headers"].([]interface{})),
		MaxAge:           GetString(spec["max_age"].(string)),
		AllowCredentials: GetBool(spec["allow_credentials"].(bool)),
	}
}

func buildTLS(specs []interface{}) *client.DomainTLS {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	return &client.DomainTLS{
		MinProtocolVersion: GetString(spec["min_protocol_version"].(string)),
		CipherSuites:       buildStringArray(spec["cipher_suites"].([]interface{})),
		ClientCertificate:  buildCertificate(spec["client_certificate"].([]interface{})),
		ServerCertificate:  buildCertificate(spec["server_certificate"].([]interface{})),
	}
}

func buildAllowOrigins(specs []interface{}) *[]client.DomainAllowOrigin {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []client.DomainAllowOrigin{}
	for _, item := range specs {
		allowOrigin := item.(map[string]interface{})
		collection = append(collection, client.DomainAllowOrigin{
			Exact: GetString(allowOrigin["exact"].(string)),
		})
	}

	return &collection
}

func buildCertificate(specs []interface{}) *client.DomainCertificate {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	return &client.DomainCertificate{
		SecretLink: GetString(spec["secret_link"].(string)),
	}
}

// Status Related //
func buildDomainStatus(specs []interface{}) *client.DomainStatus {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	return &client.DomainStatus{
		EndPoints:   buildEndPoints(spec["end_points"].([]interface{})),
		Status:      GetString(spec["status"].(string)),
		Warning:     GetString(spec["warning"].(string)),
		Locations:   buildStatusLocations(spec["locations"].([]interface{})),
		Fingerprint: GetString(spec["fingerprint"].(string)),
	}
}

func buildEndPoints(specs []interface{}) *[]client.DomainEndPoint {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []client.DomainEndPoint{}
	for _, item := range specs {
		endPoint := item.(map[string]interface{})
		collection = append(collection, client.DomainEndPoint{
			URL:          GetString(endPoint["url"].(string)),
			WorkloadLink: GetString(endPoint["workload_link"].(string)),
		})
	}

	return &collection
}

func buildStatusLocations(specs []interface{}) *[]client.DomainStatusLocation {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	collection := []client.DomainStatusLocation{}
	for _, item := range specs {
		location := item.(map[string]interface{})
		collection = append(collection, client.DomainStatusLocation{
			Name:              GetString(location["name"].(string)),
			CertificateStatus: GetString(location["certificate_status"].(string)),
		})
	}

	return &collection
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

// Status Related //
func flattenDomainStatus(domainStatus *client.DomainStatus) []interface{} {
	if domainStatus == nil {
		return nil
	}

	status := make(map[string]interface{})
	status["end_points"] = flattenEndPoints(domainStatus.EndPoints)

	if domainStatus.Status != nil {
		status["status"] = *domainStatus.Status
	}

	if domainStatus.Warning != nil {
		status["warning"] = *domainStatus.Warning
	}

	status["locations"] = flattenStatusLocations(domainStatus.Locations)

	if domainStatus.Fingerprint != nil {
		status["fingerprint"] = *domainStatus.Fingerprint
	}

	return []interface{}{
		status,
	}
}

func flattenEndPoints(endPoints *[]client.DomainEndPoint) []interface{} {
	if endPoints == nil || len(*endPoints) == 0 {
		return nil
	}

	collection := make([]interface{}, len(*endPoints))
	for i, item := range *endPoints {

		endPoint := make(map[string]interface{})
		if item.URL != nil {
			endPoint["url"] = *item.URL
		}

		if item.WorkloadLink != nil {
			endPoint["workload_link"] = *item.WorkloadLink
		}

		collection[i] = endPoint
	}

	return collection
}

func flattenStatusLocations(locations *[]client.DomainStatusLocation) []interface{} {
	if locations == nil || len(*locations) == 0 {
		return nil
	}

	collection := make([]interface{}, len(*locations))
	for i, item := range *locations {

		location := make(map[string]interface{})
		if item.Name != nil {
			location["name"] = *item.Name
		}

		if item.CertificateStatus != nil {
			location["certificate_status"] = *item.CertificateStatus
		}

		collection[i] = location
	}

	return collection
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
	if strings != nil || len(*strings) == 0 {
		return nil
	}

	collection := make([]interface{}, len(*strings))
	for i, item := range *strings {
		collection[i] = item
	}

	return collection
}
