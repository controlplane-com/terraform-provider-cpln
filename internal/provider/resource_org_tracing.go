package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

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
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the organization.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the organization.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of Org.",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeMap,
				Description: "Key-value map of the Org's tags.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lightstep_tracing":    client.LightstepSchema(true),
			"otel_tracing":         client.OtelSchema(true),
			"controlplane_tracing": client.ControlPlaneTracingSchema(true),
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceOrgTracingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgTracingCreate")

	c := m.(*client.Client)

	var traceCreate *client.Tracing

	traceCreate = buildLightStepTracing(d.Get("lightstep_tracing").([]interface{}))

	if traceCreate == nil {
		traceCreate = buildOtelTracing(d.Get("otel_tracing").([]interface{}))
	}

	if traceCreate == nil {
		traceCreate = buildControlPlaneTracing(d.Get("controlplane_tracing").([]interface{}))
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

	if org.Spec != nil && org.Spec.Tracing != nil && org.Spec.Tracing.Provider != nil && org.Spec.Tracing.Provider.Otel != nil {
		if err := d.Set("otel_tracing", flattenOtelTracing(org.Spec.Tracing)); err != nil {
			return diag.FromErr(err)
		}
	}

	if org.Spec != nil && org.Spec.Tracing != nil && org.Spec.Tracing.Provider != nil && org.Spec.Tracing.Provider.ControlPlane != nil {
		if err := d.Set("controlplane_tracing", flattenControlPlaneTracing(org.Spec.Tracing)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func flattenCustomTags(customTags *map[string]client.CustomTag) map[string]interface{} {
	if customTags == nil {
		return nil
	}

	stringMap := map[string]interface{}{}

	for key, value := range *customTags {
		stringMap[key] = value.Literal.Value
	}

	return stringMap
}

func buildCustomTags(stringMap *map[string]interface{}) *map[string]client.CustomTag {
	if stringMap == nil || len(*stringMap) == 0 {
		return nil
	}

	customTags := map[string]client.CustomTag{}

	for key, value := range *stringMap {
		customTags[key] = client.CustomTag{
			Literal: &client.CustomTagValue{
				Value: GetString(value),
			},
		}
	}

	return &customTags
}

func flattenLightstepTracing(trace *client.Tracing) []interface{} {

	if trace != nil {

		outputMap := make(map[string]interface{})

		outputMap["sampling"] = *trace.Sampling
		outputMap["endpoint"] = *trace.Provider.Lightstep.Endpoint
		outputMap["credentials"] = *trace.Provider.Lightstep.Credentials

		if trace.CustomTags != nil {
			outputMap["custom_tags"] = flattenCustomTags(trace.CustomTags)
		}

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
			Sampling:   GetInt(trace["sampling"]),
			CustomTags: buildCustomTags(GetStringMap(trace["custom_tags"])),
			Provider: &client.Provider{
				Lightstep: iTrace,
			},
		}
	}

	return nil
}

func flattenOtelTracing(trace *client.Tracing) []interface{} {

	if trace != nil {

		outputMap := make(map[string]interface{})

		outputMap["sampling"] = *trace.Sampling
		outputMap["endpoint"] = *trace.Provider.Otel.Endpoint

		if trace.CustomTags != nil {
			outputMap["custom_tags"] = flattenCustomTags(trace.CustomTags)
		}

		output := make([]interface{}, 1)
		output[0] = outputMap

		return output
	}

	return nil
}

func buildOtelTracing(tracing []interface{}) *client.Tracing {

	if len(tracing) == 1 {

		trace := tracing[0].(map[string]interface{})

		iTrace := &client.OtelTelemetry{}
		iTrace.Endpoint = GetString(trace["endpoint"])

		return &client.Tracing{
			Sampling:   GetInt(trace["sampling"]),
			CustomTags: buildCustomTags(GetStringMap(trace["custom_tags"])),
			Provider: &client.Provider{
				Otel: iTrace,
			},
		}
	}

	return nil
}

func flattenControlPlaneTracing(trace *client.Tracing) []interface{} {
	if trace == nil {
		return nil
	}

	outputMap := make(map[string]interface{})

	outputMap["sampling"] = *trace.Sampling

	if trace.CustomTags != nil {
		outputMap["custom_tags"] = flattenCustomTags(trace.CustomTags)
	}

	return []interface{}{
		outputMap,
	}
}

func buildControlPlaneTracing(tracing []interface{}) *client.Tracing {
	if len(tracing) == 0 {
		return nil
	}

	trace := tracing[0].(map[string]interface{})
	iTrace := &client.ControlPlaneTracing{}

	return &client.Tracing{
		Sampling:   GetInt(trace["sampling"]),
		CustomTags: buildCustomTags(GetStringMap(trace["custom_tags"])),
		Provider: &client.Provider{
			ControlPlane: iTrace,
		},
	}
}

func resourceOrgTracingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceOrgTracingUpdate")

	if d.HasChanges("lightstep_tracing", "otel_tracing", "controlplane_tracing") {

		c := m.(*client.Client)

		var traceUpdate *client.Tracing

		traceUpdate = buildLightStepTracing(d.Get("lightstep_tracing").([]interface{}))

		if traceUpdate == nil {
			traceUpdate = buildOtelTracing(d.Get("otel_tracing").([]interface{}))
		}

		if traceUpdate == nil {
			traceUpdate = buildControlPlaneTracing(d.Get("controlplane_tracing").([]interface{}))
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
