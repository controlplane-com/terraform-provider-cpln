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
			"domain_link": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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
	domainLink := d.Get("domain_link").(string)
	route := client.DomainRoute{
		Prefix:        GetString(d.Get("prefix")),
		ReplacePrefix: GetString(d.Get("replace_prefix")),
		WorkloadLink:  GetString(d.Get("workload_link")),
		Port:          GetInt(d.Get("port")),
	}

	c := m.(*client.Client)
	newRoute, err := c.AddDomainRoute(GetNameFromSelfLink(domainLink), route)

	if err != nil {
		return diag.FromErr(err)
	}

	return setDomainRoute(d, domainLink, newRoute)
}

func resourceDomainRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainLink := d.Id()
	prefix := d.Get("prefix").(string)

	c := m.(*client.Client)
	domain, code, err := c.GetDomain(GetNameFromSelfLink(domainLink))

	if code == 404 {
		return setDomainRoute(d, domainLink, nil)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	for _, route := range *(*domain.Spec.Ports)[0].Routes {
		if route.Prefix != &prefix {
			continue
		}

		return setDomainRoute(d, domainLink, &route)
	}

	return nil
}

func resourceDomainRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainLink := d.Get("domain_link").(string)
	route := &client.DomainRoute{
		Prefix:        GetString(d.Get("prefix")),
		ReplacePrefix: GetString(d.Get("replace_prefix")),
		WorkloadLink:  GetString(d.Get("workload_link")),
		Port:          GetInt(d.Get("port")),
	}

	c := m.(*client.Client)

	newRoute, err := c.UpdateDomainRoute(GetNameFromSelfLink(domainLink), route)
	if err != nil {
		return diag.FromErr(err)
	}

	return setDomainRoute(d, domainLink, newRoute)
}

func resourceDomainRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainLink := d.Get("domain_link").(string)
	prefix := d.Get("prefix").(string)

	c := m.(*client.Client)

	_, err := c.RemoveDomainRoute(GetNameFromSelfLink(domainLink), prefix)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setDomainRoute(d *schema.ResourceData, domainLink string, route *client.DomainRoute) diag.Diagnostics {
	if route == nil {
		d.SetId("")
		return nil
	}

	d.SetId(domainLink)

	if err := d.Set("domain_link", domainLink); err != nil {
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
