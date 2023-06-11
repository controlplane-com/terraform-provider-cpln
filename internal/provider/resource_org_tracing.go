package cpln

import (
	"context"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrgTracing() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceOrgTracingCreate,
		ReadContext:   resourceOrgTracingRead,
		UpdateContext: resourceOrgTracingUpdate,
		DeleteContext: resourceOrgTracingDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lightstep_tracing": client.LightstepSchema(),
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceOrgTracingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgTracingCreate")

	c := m.(*client.Client)

	var traceCreate *client.Tracing

	traceArray := d.Get("lightstep_tracing").([]interface{})
	if len(traceArray) == 1 {
		traceCreate = buildLightStepTracing(traceArray)
	}

	org, _, err := c.UpdateOrgTracing(traceCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	return setOrgTracing(d, org)
}

func resourceOrgTracingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgTracingRead")

	c := m.(*client.Client)
	org, _, err := c.GetOrg()

	if err != nil {
		return diag.FromErr(err)
	}

	return setOrgTracing(d, org)
}

func setOrgTracing(d *schema.ResourceData, org *client.Org) diag.Diagnostics {

	if org == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*org.Name)

	if err := SetBase(d, org.Base); err != nil {
		return diag.FromErr(err)
	}

	if org.Spec != nil && org.Spec.Tracing != nil && org.Spec.Tracing.Provider != nil && org.Spec.Tracing.Provider.Lightstep != nil {
		if err := d.Set("lightstep_tracing", flattenLightstepTracing(org.Spec.Tracing)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func flattenLightstepTracing(trace *client.Tracing) []interface{} {

	if trace != nil {

		outputMap := make(map[string]interface{})

		outputMap["sampling"] = *trace.Sampling
		outputMap["endpoint"] = *trace.Provider.Lightstep.Endpoint
		outputMap["credentials"] = *trace.Provider.Lightstep.Credentials

		output := make([]interface{}, 1)
		output[0] = outputMap

		return output
	}

	return nil
}

func buildLightStepTracing(tracing []interface{}) *client.Tracing {

	if len(tracing) == 1 {

		trace := tracing[0].(map[string]interface{})

		iTrace := &client.LightstepTracing{}
		iTrace.Endpoint = GetString(trace["endpoint"])
		iTrace.Credentials = GetString(trace["credentials"])

		return &client.Tracing{
			Sampling: GetInt(trace["sampling"]),
			Provider: &client.Provider{
				Lightstep: iTrace,
			},
		}
	}

	return nil
}

func resourceOrgTracingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgTracingUpdate")

	if d.HasChanges("lightstep_tracing") {

		c := m.(*client.Client)

		var traceUpdate *client.Tracing

		if d.HasChange("lightstep_tracing") {
			traceArray := d.Get("lightstep_tracing").([]interface{})

			if traceArray != nil {
				traceUpdate = buildLightStepTracing(traceArray)
			}
		}

		org, _, err := c.UpdateOrgTracing(traceUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setOrgTracing(d, org)
	}

	return nil
}

func resourceOrgTracingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgTracingDelete")

	c := m.(*client.Client)

	_, _, err := c.UpdateOrgTracing(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
