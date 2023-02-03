package cpln

import (
	"context"

	client "terraform-provider-cpln/internal/provider/client"

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
			"domain_name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"domain_port": {
				Type:    schema.TypeInt,
				Default: 443,
			},
			"prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replace_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workload_link": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceDomainRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("domain_name").(string)
	domainPort := d.Get("domain_port").(int)
	route := client.DomainRoute{
		Prefix:        GetString(d.Get("prefix")),
		ReplacePrefix: GetString(d.Get("replace_prefix")),
		WorkloadLink:  GetString(d.Get("workload_link")),
		Port:          GetInt(d.Get("port")),
	}

	c := m.(*client.Client)

	newRoute, err := c.AddDomainRoute(domainName, domainPort, route)
	if err != nil {
		return diag.FromErr(err)
	}

	return setDomainRoute(d, domainName, domainPort, newRoute)
}

func resourceDomainRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Id()
	c := m.(*client.Client)
	domain, code, err := c.GetDomain(domainName)

	// TODO fix this logic for all resources, we don't need to set gvc to nil when a domain is not found, unrelated
	if code == 404 {
		return setGvc(d, nil, c.Org)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	domainPort := d.Get("domain_port").(int)
	prefix := d.Get("prefix").(string)

	if domain.Spec.Ports != nil && len(*domain.Spec.Ports) > 0 {
		for _, port := range *domain.Spec.Ports {
			if port.Number != &domainPort {
				continue
			}

			if port.Routes == nil || len(*port.Routes) == 0 {
				continue
			}

			for _, route := range *port.Routes {
				if route.Prefix != &prefix {
					continue
				}

				return setDomainRoute(d, domainName, domainPort, &route)
			}
		}
	}

	return nil
}

func resourceDomainRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("domain_name").(string)
	domainPort := d.Get("domain_port").(int)
	route := client.DomainRoute{
		Prefix:        GetString(d.Get("prefix")),
		ReplacePrefix: GetString(d.Get("replace_prefix")),
		WorkloadLink:  GetString(d.Get("workload_link")),
		Port:          GetInt(d.Get("port")),
	}

	c := m.(*client.Client)

	newRoute, err := c.UpdateDomainRoute(domainName, domainPort, route)
	if err != nil {
		return diag.FromErr(err)
	}

	return setDomainRoute(d, domainName, domainPort, newRoute)
}

func resourceDomainRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("domain_name").(string)
	domainPort := d.Get("domain_port").(int)
	prefix := d.Get("prefix").(string)

	c := m.(*client.Client)

	_, err := c.RemoveDomainRoute(domainName, domainPort, prefix)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setDomainRoute(d *schema.ResourceData, domainName string, domainPort int, route *client.DomainRoute) diag.Diagnostics {
	if route == nil {
		d.SetId("")
		return nil
	}

	d.SetId(domainName)

	if err := d.Set("domain_name", domainName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("domain_port", domainPort); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("prefix", route.Prefix); err != nil {
		return diag.FromErr(err)
	}

	if route.ReplacePrefix != nil && *route.ReplacePrefix != "" {
		if err := d.Set("replace_prefix", route.ReplacePrefix); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("workload_link", route.WorkloadLink); err != nil {
		return diag.FromErr(err)
	}

	if route.Port != nil {
		if err := d.Set("port", route.Port); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
