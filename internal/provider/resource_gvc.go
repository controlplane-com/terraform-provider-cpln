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
			Type:        schema.TypeString,
			Description: "The ID, in GUID format, of the GVC.",
			Computed:    true,
		},
		"name": {
			Type:         schema.TypeString,
			Description:  "Name of the GVC.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: NameValidator,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the GVC.",
			Optional:         true,
			ValidateFunc:     DescriptionValidator,
			DiffSuppressFunc: DiffSuppressDescription,
		},
		"domain": {
			Type:        schema.TypeString,
			Description: "Custom domain name used by associated workloads.",
			Optional:    true,
			Deprecated:  "Selecting a domain on a GVC will be deprecated in the future. Use the 'cpln_domain resource' instead.",
		},
		"alias": {
			Type:        schema.TypeString,
			Description: "The alias name of the GVC.",
			Computed:    true,
		},
		"pull_secrets": {
			Type:        schema.TypeSet,
			Description: "A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"locations": {
			Type:        schema.TypeSet,
			Description: "A list of [locations](https://docs.controlplane.com/reference/location#current) making up the Global Virtual Cloud.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:        schema.TypeMap,
			Description: "Key-value map of resource tags.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ValidateFunc: TagValidator,
		},
		"env": {
			Type:        schema.TypeMap,
			Description: "Key-value array of resource env variables.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"self_link": {
			Type:        schema.TypeString,
			Description: "Full link to this resource. Can be referenced by other resources.",
			Computed:    true,
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
			Type:        schema.TypeList,
			Description: "Dedicated load balancer configuration.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"dedicated": {
						Type:        schema.TypeBool,
						Description: "Creates a dedicated load balancer in each location and enables additional Domain features: custom ports, protocols and wildcard hostnames. Charges apply for each location.",
						Optional:    true,
					},
					"trusted_proxies": {
						Type:        schema.TypeInt,
						Description: "Controls the address used for request logging and for setting the X-Envoy-External-Address header. If set to 1, then the last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If set to 2, then the second to last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If the XFF header does not have at least two addresses or does not exist then the source client IP address will be used instead.",
						Optional:    true,
						Default:     0,
					},
					"redirect": {
						Type:        schema.TypeList,
						Description: "Specify the url to be redirected to for different http status codes.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"class": {
									Type:        schema.TypeList,
									Description: "Specify the redirect url for all status codes in a class.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"status_5xx": {
												Type:        schema.TypeString,
												Description: "Specify the redirect url for any 500 level status code.",
												Optional:    true,
											},
											"_sentinel": {
												Type:     schema.TypeBool,
												Optional: true,
												Default:  true,
											},
										},
									},
								},
								"_sentinel": {
									Type:     schema.TypeBool,
									Optional: true,
									Default:  true,
								},
							},
						},
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
	output := client.LoadBalancer{}

	if spec["dedicated"] != nil {
		output.Dedicated = GetBool(spec["dedicated"].(bool))
	}

	if spec["trusted_proxies"] != nil {
		output.TrustedProxies = GetInt(spec["trusted_proxies"].(int))
	}

	if spec["redirect"] != nil {
		output.Redirect = buildRedirect(spec["redirect"].([]interface{}))
	}

	return &output
}

func buildRedirect(specs []interface{}) *client.Redirect {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Redirect{}

	if spec["class"] != nil {
		output.Class = buildRedirectClass(spec["class"].([]interface{}))
	}

	return &output
}

func buildRedirectClass(specs []interface{}) *client.RedirectClass {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.RedirectClass{}

	if spec["status_5xx"] != nil {
		output.Status5XX = GetString(spec["status_5xx"])
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

	loadBalancer := map[string]interface{}{}

	if gvcSpec.Dedicated != nil {
		loadBalancer["dedicated"] = *gvcSpec.Dedicated
	}

	if gvcSpec.TrustedProxies != nil {
		loadBalancer["trusted_proxies"] = *gvcSpec.TrustedProxies
	}

	if gvcSpec.Redirect != nil {
		loadBalancer["redirect"] = flattenRedirect(gvcSpec.Redirect)
	}

	return []interface{}{
		loadBalancer,
	}
}

func flattenRedirect(spec *client.Redirect) []interface{} {

	if spec == nil {
		return nil
	}

	redirect := map[string]interface{}{
		"_sentinel": true,
	}

	if spec.Class != nil {
		redirect["class"] = flattenRedirectClass(spec.Class)
	}

	return []interface{}{
		redirect,
	}
}

func flattenRedirectClass(spec *client.RedirectClass) []interface{} {

	if spec == nil {
		return nil
	}

	class := map[string]interface{}{
		"_sentinel": true,
	}

	if spec.Status5XX != nil {
		class["status_5xx"] = *spec.Status5XX
	}

	return []interface{}{
		class,
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
