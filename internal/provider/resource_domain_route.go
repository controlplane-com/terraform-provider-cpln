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
				Type:    schema.TypeString,
				Default: "/",
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
	route := client.DomainRoute{
		Prefix:        GetString(d.Get("prefix")),
		ReplacePrefix: GetString(d.Get("replace_prefix")),
		WorkloadLink:  GetString(d.Get("workload_link")),
		Port:          GetInt(d.Get("port")),
	}

	c := m.(*client.Client)
	newRoute, err := c.AddDomainRoute(domainName, route)

	if err != nil {
		return diag.FromErr(err)
	}

	return setDomainRoute(d, domainName, newRoute)
}

func resourceDomainRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Id()
	c := m.(*client.Client)
	domain, code, err := c.GetDomain(domainName)

	if code == 404 {
		return setGvc(d, nil, c.Org)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if domain.Spec.Ports != nil && len(*domain.Spec.Ports) > 0 {
		for _, port := range *domain.Spec.Ports {
			if port.Routes == nil || len(*port.Routes) == 0 {
				continue
			}

			for _, route := range *port.Routes {
				if err := setDomainRoute(d, domainName, &route); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func resourceDomainRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	err := c.RemoveDomainRoute()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setDomainRoute(d *schema.ResourceData, domainName string, route *client.DomainRoute) diag.Diagnostics {
	if route == nil {
		d.SetId("")
		return nil
	}

	d.SetId(domainName)

	if err := d.Set("domain_name", domainName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("prefix", route.Prefix); err != nil {
		return diag.FromErr(err)
	}

	if route.ReplacePrefix != nil && *route.ReplacePrefix != "" {
		if err := d.Set("replace_prefix", route.Prefix); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("workload_link", route.Prefix); err != nil {
		return diag.FromErr(err)
	}

	if route.Port != nil {
		if err := d.Set("port", route.Port); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
