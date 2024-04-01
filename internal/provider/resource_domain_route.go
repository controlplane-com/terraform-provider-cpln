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
				Type:        schema.TypeString,
				Description: "The path will match any unmatched path prefixes for the subdomain.",
				ForceNew:    true,
				Required:    true,
			},
			"replace_prefix": {
				Type:        schema.TypeString,
				Description: "A path prefix can be configured to be replaced when forwarding the request to the Workload.",
				Optional:    true,
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
				Type:        schema.TypeString,
				Description: "This option allows forwarding traffic for different host headers to different workloads. This will only be used when the target GVC has dedicated load balancing enabled and the Domain is configured for wildcard support. Please contact us on Slack or at support@controlplane.com for additional details.",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateDomainRoute,
		},
	}
}

func importStateDomainRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.SplitN(d.Id(), ":", 3)

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected ID syntax: 'domain_link:domain_port:prefix'. Example: 'terraform import cpln_domain_route.RESOURCE_NAME DOMAIN_LINK:DOMAIN_PORT:PREFIX'", d.Id())
	}

	domainLink := parts[0]
	domainPort, err := strconv.Atoi(parts[1])
	routePrefix := parts[2]

	if err != nil {
		return nil, fmt.Errorf("unexpected format of ID (%s), domain port is invalid, must be a integer. value provided: %s. error: %s", d.Id(), parts[1], err.Error())
	}

	d.Set("domain_link", domainLink)
	d.Set("domain_port", domainPort)
	d.Set("prefix", routePrefix)
	d.SetId(fmt.Sprintf("%s_%d_%s", domainLink, domainPort, routePrefix))

	return []*schema.ResourceData{d}, nil
}

func resourceDomainRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domainLink := d.Get("domain_link").(string)
	domainPort := d.Get("domain_port").(int)

	route := client.DomainRoute{
		Prefix:        GetString(d.Get("prefix")),
		ReplacePrefix: GetString(d.Get("replace_prefix")),
		WorkloadLink:  GetString(d.Get("workload_link")),
		Port:          GetInt(d.Get("port")),
		HostPrefix:    GetString(d.Get("host_prefix")),
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
	prefix := d.Get("prefix").(string)

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

				if *route.Prefix != prefix {
					continue
				}

				return setDomainRoute(d, domainLink, domainPort, &route)
			}
		}
	}

	return nil
}

func resourceDomainRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("prefix", "replace_prefix", "workload_link", "port", "host_prefix") {

		domainLink := d.Get("domain_link").(string)
		domainPort := d.Get("domain_port").(int)

		route := &client.DomainRoute{
			Prefix:        GetString(d.Get("prefix")),
			ReplacePrefix: GetString(d.Get("replace_prefix")),
			WorkloadLink:  GetString(d.Get("workload_link")),
			Port:          GetInt(d.Get("port")),
			HostPrefix:    GetString(d.Get(("host_prefix"))),
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
	prefix := d.Get("prefix").(string)

	c := m.(*client.Client)

	err := c.RemoveDomainRoute(GetNameFromSelfLink(domainLink), domainPort, prefix)

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

	d.SetId(fmt.Sprintf("%s_%d_%s", domainLink, domainPort, *route.Prefix))

	if err := d.Set("domain_link", domainLink); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("domain_port", domainPort); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("prefix", route.Prefix); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("replace_prefix", route.ReplacePrefix); err != nil {
		return diag.FromErr(err)
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

	return nil
}
