package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Main ***/
func resourceGvc() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceGvcCreate,
		ReadContext:   resourceGvcRead,
		UpdateContext: resourceGvcUpdate,
		DeleteContext: resourceGvcDelete,
		Schema:        GvcSchema(),
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			// Check if both attributes are set
			if len(diff.Get("lightstep_tracing").([]interface{})) > 0 && len(diff.Get("otel_tracing").([]interface{})) > 0 && len(diff.Get("controlplane_tracing").([]interface{})) > 0 {
				return fmt.Errorf("only one of lightstep_tracing, otel_tracing or controlplane_tracing can be specified")
			}
			return nil
		},
	}
}

func GvcSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpln_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: NameValidator,
		},
		"description": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateFunc:     DescriptionValidator,
			DiffSuppressFunc: DiffSuppressDescription,
		},
		"domain": {
			Type:       schema.TypeString,
			Optional:   true,
			Deprecated: "Selecting a domain on a GVC will be deprecated in the future. Use the 'cpln_domain resource' instead.",
		},
		"alias": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pull_secrets": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"locations": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ValidateFunc: TagValidator,
		},
		"env": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"self_link": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"lightstep_tracing":    client.LightstepSchema(false),
		"otel_tracing":         client.OtelSchema(false),
		"controlplane_tracing": client.ControlPlaneTracingSchema(false),
		"sidecar": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"envoy": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"load_balancer": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"dedicated": {
						Type:     schema.TypeBool,
						Required: true,
					},
					"trusted_proxies": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  0,
					},
				},
			},
		},
	}
}

func resourceGvcCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceGvcCreate")

	gvc := client.Gvc{}

	gvc.Name = GetString(d.Get("name"))
	gvc.Description = GetString(d.Get("description"))
	gvc.Tags = GetStringMap(d.Get("tags"))

	gvc.Spec = &client.GvcSpec{}

	gvc.Spec.Domain = GetString(d.Get("domain"))

	gvcEnv := []client.NameValue{}
	keys, envMap := MapSortHelper(d.Get("env"))

	for _, k := range keys {
		envName := k
		envValue := envMap[envName].(string)

		localEnv := client.NameValue{
			Name:  &envName,
			Value: &envValue,
		}

		gvcEnv = append(gvcEnv, localEnv)
	}

	if len(keys) > 0 {
		gvc.Spec.Env = &gvcEnv
	}

	c := m.(*client.Client)

	buildLocations(c.Org, d.Get("locations"), gvc.Spec)
	buildPullSecrets(c.Org, d.Get("pull_secrets"), gvc.Spec)
	gvc.Spec.LoadBalancer = buildLoadBalancer(d.Get("load_balancer").([]interface{}))

	gvc.Spec.Tracing = buildLightStepTracing(d.Get("lightstep_tracing").([]interface{}))

	if gvc.Spec.Tracing == nil {
		gvc.Spec.Tracing = buildOtelTracing(d.Get("otel_tracing").([]interface{}))
	}

	if gvc.Spec.Tracing == nil {
		gvc.Spec.Tracing = buildControlPlaneTracing(d.Get("controlplane_tracing").([]interface{}))
	}

	if d.Get("sidecar") != nil {
		gvc.Spec.Sidecar = buildGvcSidecar(d.Get("sidecar").([]interface{}))
	}

	newGvc, code, err := c.CreateGvc(gvc)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setGvc(d, newGvc, c.Org)
}

func resourceGvcRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceGvcRead")

	c := m.(*client.Client)
	gvc, code, err := c.GetGvc(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setGvc(d, gvc, c.Org)
}

func resourceGvcUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceGvcUpdate")

	if d.HasChanges("description", "locations", "env", "tags", "domain", "pull_secrets", "lightstep_tracing", "otel_tracing", "controlplane_tracing", "load_balancer", "sidecar") {

		c := m.(*client.Client)

		gvcToUpdate := client.Gvc{}
		gvcToUpdate.Name = GetString(d.Get("name"))
		gvcToUpdate.Description = GetDescriptionString(d.Get("description"), *gvcToUpdate.Name)
		gvcToUpdate.Tags = GetTagChanges(d)

		gvcToUpdate.SpecReplace = &client.GvcSpec{}
		gvcToUpdate.SpecReplace.Domain = GetString(d.Get("domain"))
		buildLocations(c.Org, d.Get("locations"), gvcToUpdate.SpecReplace)
		buildPullSecrets(c.Org, d.Get("pull_secrets"), gvcToUpdate.SpecReplace)
		gvcToUpdate.SpecReplace.Env = GetGVCEnvChanges(d)
		gvcToUpdate.SpecReplace.LoadBalancer = buildLoadBalancer(d.Get("load_balancer").([]interface{}))
		gvcToUpdate.SpecReplace.Sidecar = buildGvcSidecar(d.Get("sidecar").([]interface{}))

		gvcToUpdate.SpecReplace.Tracing = buildLightStepTracing(d.Get("lightstep_tracing").([]interface{}))

		if gvcToUpdate.SpecReplace.Tracing == nil {
			gvcToUpdate.SpecReplace.Tracing = buildOtelTracing(d.Get("otel_tracing").([]interface{}))
		}

		if gvcToUpdate.SpecReplace.Tracing == nil {
			gvcToUpdate.SpecReplace.Tracing = buildControlPlaneTracing(d.Get("controlplane_tracing").([]interface{}))
		}

		updatedGvc, _, err := c.UpdateGvc(gvcToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setGvc(d, updatedGvc, c.Org)
	}

	return nil
}

func resourceGvcDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceGvcDelete")

	c := m.(*client.Client)
	err := c.DeleteGvc(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setGvc(d *schema.ResourceData, gvc *client.Gvc, org string) diag.Diagnostics {

	if gvc == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*gvc.Name)

	if err := SetBase(d, gvc.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("domain", flattenDomain(gvc.Spec)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", flattenLocations(gvc.Spec, org)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pull_secrets", flattenPullSecrets(gvc.Spec, org)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("load_balancer", flattenLoadBalancer(gvc.Spec.LoadBalancer)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("sidecar", flattenGvcSidecar(gvc.Spec.Sidecar)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("alias", gvc.Alias); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(gvc.Links, d); err != nil {
		return diag.FromErr(err)
	}

	if gvc.Spec != nil && gvc.Spec.Env != nil {
		if len(*gvc.Spec.Env) > 0 {

			envMap := make(map[string]interface{}, len(*gvc.Spec.Env))

			for _, envObj := range *gvc.Spec.Env {
				key := envObj.Name
				value := envObj.Value
				envMap[*key] = value
			}

			if err := d.Set("env", envMap); err != nil {
				return diag.FromErr(err)
			}

		} else {

			emptyEnvMap := make(map[string]interface{}, 0)

			if err := d.Set("env", emptyEnvMap); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {

		emptyEnvMap := make(map[string]interface{}, 0)

		if err := d.Set("env", emptyEnvMap); err != nil {
			return diag.FromErr(err)
		}
	}

	if gvc.Spec != nil && gvc.Spec.Tracing != nil && gvc.Spec.Tracing.Provider != nil && gvc.Spec.Tracing.Provider.Lightstep != nil {
		if err := d.Set("lightstep_tracing", flattenLightstepTracing(gvc.Spec.Tracing)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("lightstep_tracing", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	if gvc.Spec != nil && gvc.Spec.Tracing != nil && gvc.Spec.Tracing.Provider != nil && gvc.Spec.Tracing.Provider.Otel != nil {
		if err := d.Set("otel_tracing", flattenOtelTracing(gvc.Spec.Tracing)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("otel_tracing", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	if gvc.Spec != nil && gvc.Spec.Tracing != nil && gvc.Spec.Tracing.Provider != nil && gvc.Spec.Tracing.Provider.ControlPlane != nil {
		if err := d.Set("controlplane_tracing", flattenControlPlaneTracing(gvc.Spec.Tracing)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("controlplane_tracing", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

/*** Build ***/
func buildLocations(org string, locations interface{}, gvcSpec *client.GvcSpec) {

	l := []string{}

	if locations != nil {
		for _, location := range locations.(*schema.Set).List() {
			l = append(l, fmt.Sprintf("/org/%s/location/%s", org, location))
		}
	}

	if gvcSpec.StaticPlacement == nil {
		gvcSpec.StaticPlacement = &client.StaticPlacement{}
	}

	gvcSpec.StaticPlacement.LocationLinks = &l
}

func buildPullSecrets(org string, pullSecrets interface{}, gvcSpec *client.GvcSpec) {

	l := []string{}

	if pullSecrets != nil {
		for _, secret := range pullSecrets.(*schema.Set).List() {
			l = append(l, fmt.Sprintf("/org/%s/secret/%s", org, secret))
		}
	}

	gvcSpec.PullSecretLinks = &l
}

func buildLoadBalancer(specs []interface{}) *client.LoadBalancer {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.LoadBalancer{
		Dedicated: GetBool(spec["dedicated"].(bool)),
	}

	if spec["trusted_proxies"] != nil {
		output.TrustedProxies = GetInt(spec["trusted_proxies"].(int))
	}

	return &output
}

func buildGvcSidecar(specs []interface{}) *client.GvcSidecar {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.GvcSidecar{}

	// Attempt to unmarshal `envoy`
	var envoy interface{}
	err := json.Unmarshal([]byte(spec["envoy"].(string)), &envoy)
	if err != nil {
		log.Fatalf("Error occurred during unmarshaling 'envoy' value. Error: %s", err.Error())
	}

	// Set envoy
	output.Envoy = &envoy

	return &output
}

/*** Flatten ***/
func flattenDomain(gvcSpec *client.GvcSpec) *string {

	if gvcSpec != nil && gvcSpec.Domain != nil {
		return gvcSpec.Domain
	}

	return nil
}

func flattenLocations(gvcSpec *client.GvcSpec, org string) []interface{} {

	if gvcSpec != nil && gvcSpec.StaticPlacement != nil && gvcSpec.StaticPlacement.LocationLinks != nil && len(*gvcSpec.StaticPlacement.LocationLinks) > 0 {

		l := make([]interface{}, len(*gvcSpec.StaticPlacement.LocationLinks))

		for i, location := range *gvcSpec.StaticPlacement.LocationLinks {
			location = strings.TrimPrefix(location, fmt.Sprintf("/org/%s/location/", org))
			l[i] = location
		}

		return l
	}

	return make([]interface{}, 0)
}

func flattenPullSecrets(gvcSpec *client.GvcSpec, org string) []interface{} {

	if gvcSpec != nil && gvcSpec.PullSecretLinks != nil && len(*gvcSpec.PullSecretLinks) > 0 {

		l := make([]interface{}, len(*gvcSpec.PullSecretLinks))

		for i, secret := range *gvcSpec.PullSecretLinks {
			secret = strings.TrimPrefix(secret, fmt.Sprintf("/org/%s/secret/", org))
			l[i] = secret
		}

		return l
	}

	return make([]interface{}, 0)
}

func flattenLoadBalancer(gvcSpec *client.LoadBalancer) []interface{} {
	if gvcSpec == nil {
		return nil
	}

	loadBalancer := map[string]interface{}{
		"dedicated": *gvcSpec.Dedicated,
	}

	if gvcSpec.TrustedProxies != nil {
		loadBalancer["trusted_proxies"] = *gvcSpec.TrustedProxies
	}

	return []interface{}{
		loadBalancer,
	}
}

func flattenGvcSidecar(gvcSpec *client.GvcSidecar) []interface{} {
	if gvcSpec == nil {
		return nil
	}

	// Attempt to marshal `envoy`
	jsonOut, err := json.Marshal(*gvcSpec.Envoy)
	if err != nil {
		log.Fatalf("Error occurred during marshaling 'envoy' value. Error: %s", err.Error())
	}

	sidecar := map[string]interface{}{
		"envoy": string(jsonOut),
	}

	return []interface{}{
		sidecar,
	}
}
