package cpln

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var routeRegexOrPrefixAttribute = []string{"prefix", "regex"}

func resourceDomainRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainRouteCreate,
		ReadContext:   resourceDomainRouteRead,
		UpdateContext: resourceDomainRouteUpdate,
		DeleteContext: resourceDomainRouteDelete,
		Schema: map[string]*schema.Schema{
			"domain_link": {
				Type:        schema.TypeString,
				Description: "The self link of the domain to add the route to.",
				ForceNew:    true,
				Required:    true,
			},
			"domain_port": {
				Type:        schema.TypeInt,
				Description: "The port the route corresponds to. Default: 443",
				ForceNew:    true,
				Optional:    true,
				Default:     443,
			},
			"prefix": {
				Type:         schema.TypeString,
				Description:  "The path will match any unmatched path prefixes for the subdomain.",
				ForceNew:     true,
				Optional:     true,
				ExactlyOneOf: routeRegexOrPrefixAttribute,
			},
			"replace_prefix": {
				Type:        schema.TypeString,
				Description: "A path prefix can be configured to be replaced when forwarding the request to the Workload.",
				Optional:    true,
			},
			"regex": {
				Type:         schema.TypeString,
				Description:  "Used to match URI paths. Uses the google re2 regex syntax.",
				ForceNew:     true,
				Optional:     true,
				ExactlyOneOf: routeRegexOrPrefixAttribute,
			},
			"workload_link": {
				Type:        schema.TypeString,
				Description: "The link of the workload to map the prefix to.",
				Required:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "For the linked workload, the port to route traffic to.",
				Optional:    true,
			},
			"host_prefix": {
				Type:          schema.TypeString,
				Description:   "This option allows forwarding traffic for different host headers to different workloads. This will only be used when the target GVC has dedicated load balancing enabled and the Domain is configured for wildcard support. Please contact us on Slack or at support@controlplane.com for additional details.",
				Optional:      true,
				ConflictsWith: []string{"host_regex"},
			},
			"host_regex": {
				Type:          schema.TypeString,
				Description:   "A regex to match the host header. This will only be used when the target GVC has dedicated load balancing enabled and the Domain is configure for wildcard support. Contact your account manager for details.",
				Optional:      true,
				ConflictsWith: []string{"host_prefix"},
			},
			"headers": {
				Type:        schema.TypeList,
				Description: "Modify the headers for all http requests for this route.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request": {
							Type:        schema.TypeList,
							Description: "Manipulates HTTP headers.",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"set": {
										Type:        schema.TypeMap,
										Description: "Sets or overrides headers to all http requests for this route.",
										Optional:    true,
										Elem:        StringSchema(),
									},
									"placeholder_attribute": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
						"placeholder_attribute": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateDomainRoute,
		},
	}
}

func importStateDomainRoute(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	// Extract required values from ResourceData
	parts := strings.SplitN(d.Id(), ":", 3)

	var domainLink string
	var domainPortStr string
	var prefixOrRegex string

	if len(parts) == 3 {
		domainLink = parts[0]
		domainPortStr = parts[1]
		prefixOrRegex = parts[2]
	}

	if domainLink == "" || domainPortStr == "" || prefixOrRegex == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected ID syntax: 'domain_link:domain_port:[prefix|regex]'. Example: 'terraform import cpln_domain_route.RESOURCE_NAME DOMAIN_LINK:DOMAIN_PORT:[PREFIX|REGEX]'", d.Id())
	}

	// Convert domainPortStr to integer
	domainPort, err := strconv.Atoi(domainPortStr)

	if err != nil {
		return nil, fmt.Errorf("unexpected format of ID (%s), domain port is invalid, must be a integer. value provided: %s. error: %s", d.Id(), parts[1], err.Error())
	}

	// Figure out if it was prefix OR regex that was provided
	var prefix *string
	var regex *string

	c := m.(*client.Client)

	domain, code, err := c.GetDomain(GetNameFromSelfLink(domainLink))

	if code == 404 {
		return nil, fmt.Errorf("domain '%s' not found", domainLink)
	}

	if err != nil {
		return nil, err
	}

	if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
		return nil, fmt.Errorf("domain '%s' does not have any port configured", domainLink)
	}

	for _, port := range *domain.Spec.Ports {

		if *port.Number != domainPort || port.Routes == nil || len(*port.Routes) == 0 {
			continue
		}

		found := false

		for _, route := range *port.Routes {

			if route.Prefix != nil && *route.Prefix == prefixOrRegex {
				prefix = route.Prefix
				found = true
				break
			}

			if route.Regex != nil && *route.Regex == prefixOrRegex {
				regex = route.Regex
				found = true
				break
			}
		}

		if found {
			break
		}
	}

	// Set values and Id
	if err := d.Set("domain_link", domainLink); err != nil {
		return nil, err
	}

	if err := d.Set("domain_port", domainPort); err != nil {
		return nil, err
	}

	var routeIdentifier string

	if prefix != nil {
		routeIdentifier = *prefix

		if err := d.Set("prefix", *prefix); err != nil {
			return nil, err
		}
	}

	if regex != nil {
		routeIdentifier = *regex

		if err := d.Set("regex", *regex); err != nil {
			return nil, err
		}
	}

	d.SetId(fmt.Sprintf("%s_%d_%s", domainLink, domainPort, routeIdentifier))

	return []*schema.ResourceData{d}, nil
}

func resourceDomainRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domainLink := d.Get("domain_link").(string)
	domainPort := d.Get("domain_port").(int)

	route := client.DomainRoute{
		ReplacePrefix: GetString(d.Get("replace_prefix")),
		WorkloadLink:  GetString(d.Get("workload_link")),
		Port:          GetInt(d.Get("port")),
		HostPrefix:    GetString(d.Get("host_prefix")),
		HostRegex:     GetString(d.Get("host_regex")),
	}

	if d.Get("prefix") != nil {
		route.Prefix = GetString(d.Get("prefix"))
	}

	if d.Get("regex") != nil {
		route.Regex = GetString(d.Get("regex"))
	}

	if d.Get("headers") != nil {
		route.Headers = buildDomainRouteHeaders(d.Get("headers").([]interface{}))
	}

	c := m.(*client.Client)
	err := c.AddDomainRoute(GetNameFromSelfLink(domainLink), domainPort, route)

	if err != nil {
		return diag.FromErr(err)
	}

	return setDomainRoute(d, domainLink, domainPort, &route)
}

func resourceDomainRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domainLink := d.Get("domain_link").(string)
	domainPort := d.Get("domain_port").(int)

	var prefix *string
	var regex *string

	if d.Get("prefix") != nil {
		prefix = GetString(d.Get("prefix").(string))
	}

	if d.Get("regex") != nil {
		regex = GetString(d.Get("regex").(string))
	}

	c := m.(*client.Client)
	domain, code, err := c.GetDomain(GetNameFromSelfLink(domainLink))

	if code == 404 {
		return setDomainRoute(d, domainLink, domainPort, nil)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	for _, value := range *domain.Spec.Ports {

		if *value.Number == domainPort && (value.Routes != nil && len(*value.Routes) > 0) {

			for _, route := range *value.Routes {

				if (prefix != nil && route.Prefix != nil && *route.Prefix == *prefix) ||
					(regex != nil && route.Regex != nil && *route.Regex == *regex) {
					return setDomainRoute(d, domainLink, domainPort, &route)
				}
			}
		}
	}

	return diag.Errorf("route not found in port %d", domainPort)
}

func resourceDomainRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("prefix", "replace_prefix", "workload_link", "port", "host_prefix", "host_regex", "headers") {

		domainLink := d.Get("domain_link").(string)
		domainPort := d.Get("domain_port").(int)

		route := &client.DomainRoute{
			ReplacePrefix: GetString(d.Get("replace_prefix")),
			WorkloadLink:  GetString(d.Get("workload_link")),
			Port:          GetInt(d.Get("port")),
			HostPrefix:    GetString(d.Get("host_prefix")),
			HostRegex:     GetString(d.Get("host_regex")),
		}

		if d.Get("prefix") != nil {
			route.Prefix = GetString(d.Get("prefix"))
		}

		if d.Get("regex") != nil {
			route.Regex = GetString(d.Get("regex"))
		}

		if d.Get("headers") != nil {
			route.Headers = buildDomainRouteHeaders(d.Get("headers").([]interface{}))
		}

		c := m.(*client.Client)

		err := c.UpdateDomainRoute(GetNameFromSelfLink(domainLink), domainPort, route)

		if err != nil {
			return diag.FromErr(err)
		}

		return setDomainRoute(d, domainLink, domainPort, route)
	}

	return nil
}

func resourceDomainRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domainLink := d.Get("domain_link").(string)
	domainPort := d.Get("domain_port").(int)

	var prefix *string
	var regex *string

	if d.Get("prefix") != nil {
		prefix = GetString(d.Get("prefix"))
	}

	if d.Get("regex") != nil {
		regex = GetString(d.Get("regex"))
	}

	c := m.(*client.Client)

	err := c.RemoveDomainRoute(GetNameFromSelfLink(domainLink), domainPort, prefix, regex)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setDomainRoute(d *schema.ResourceData, domainLink string, domainPort int, route *client.DomainRoute) diag.Diagnostics {

	if route == nil {
		d.SetId("")
		return nil
	}

	var prefixOrRegex string

	if route.Prefix != nil {
		prefixOrRegex = *route.Prefix
	}

	if route.Regex != nil {
		prefixOrRegex = *route.Regex
	}

	d.SetId(fmt.Sprintf("%s_%d_%s", domainLink, domainPort, prefixOrRegex))

	if err := d.Set("domain_link", domainLink); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("domain_port", domainPort); err != nil {
		return diag.FromErr(err)
	}

	if route.Prefix != nil {
		if err := d.Set("prefix", route.Prefix); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("replace_prefix", route.ReplacePrefix); err != nil {
		return diag.FromErr(err)
	}

	if route.Regex != nil {
		if err := d.Set("regex", route.Regex); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("workload_link", route.WorkloadLink); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("port", route.Port); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("host_prefix", route.HostPrefix); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("host_regex", route.HostRegex); err != nil {
		return diag.FromErr(err)
	}

	if route.Headers != nil {
		if err := d.Set("headers", flattenDomainRouteHeaders(route.Headers)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

/*** Build ***/

func buildDomainRouteHeaders(specs []interface{}) *client.DomainRouteHeaders {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.DomainRouteHeaders{}

	if spec["request"] != nil {
		output.Request = buildDomainHeaderOperation(spec["request"].([]interface{}))
	}

	return &output
}

func buildDomainHeaderOperation(specs []interface{}) *client.DomainHeaderOperation {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.DomainHeaderOperation{}

	if spec["set"] != nil {
		output.Set = GetStringMap(spec["set"])
	}

	return &output
}

/*** Flatten ***/

func flattenDomainRouteHeaders(headers *client.DomainRouteHeaders) []interface{} {

	if headers == nil {
		return nil
	}

	spec := map[string]interface{}{
		"placeholder_attribute": true,
	}

	if headers.Request != nil {
		spec["request"] = flattenDomainHeaderOperation(headers.Request)
	}

	return []interface{}{
		spec,
	}
}

func flattenDomainHeaderOperation(request *client.DomainHeaderOperation) []interface{} {

	if request == nil {
		return nil
	}

	spec := map[string]interface{}{
		"placeholder_attribute": true,
	}

	if request.Set != nil {
		spec["set"] = *request.Set
	}

	return []interface{}{
		spec,
	}
}
