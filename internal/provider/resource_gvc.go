package cpln

import (
	"context"
	"fmt"
	"strings"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGvc() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceGvcCreate,
		ReadContext:   resourceGvcRead,
		UpdateContext: resourceGvcUpdate,
		DeleteContext: resourceGvcDelete,
		Schema: map[string]*schema.Schema{
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
				Deprecated: "Selecting a domain on a GVC will be deprecated in the future. Use cpln_domain instead.",
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
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lightstep_tracing": client.LightstepSchema(),
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceGvcCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceGvcCreate")

	gvc := client.Gvc{}

	gvc.Name = GetString(d.Get("name"))
	gvc.Description = GetString(d.Get("description"))
	gvc.Tags = GetStringMap(d.Get("tags"))

	if d.Get("domain") != nil {
		gvc.Spec = &client.GvcSpec{}
		gvc.Spec.Domain = GetString(d.Get("domain"))
	}

	c := m.(*client.Client)

	buildLocations(c.Org, d.Get("locations"), &gvc)
	buildPullSecrets(c.Org, d.Get("pull_secrets"), &gvc)

	traceArray := d.Get("lightstep_tracing").([]interface{})
	if len(traceArray) == 1 {

		if gvc.Spec == nil {
			gvc.Spec = &client.GvcSpec{}
		}

		gvc.Spec.Tracing = buildLightStepTracing(traceArray)
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

func buildLocations(org string, locations interface{}, gvc *client.Gvc) {

	l := []string{}

	if locations != nil {
		for _, location := range locations.(*schema.Set).List() {
			l = append(l, fmt.Sprintf("/org/%s/location/%s", org, location))
		}
	}

	if gvc.Spec == nil {
		gvc.Spec = &client.GvcSpec{}
	}

	if gvc.Spec.StaticPlacement == nil {
		gvc.Spec.StaticPlacement = &client.StaticPlacement{}
	}

	gvc.Spec.StaticPlacement.LocationLinks = &l
}

func buildPullSecrets(org string, pullSecrets interface{}, gvc *client.Gvc) {

	l := []string{}

	if pullSecrets != nil {
		for _, secret := range pullSecrets.(*schema.Set).List() {
			l = append(l, fmt.Sprintf("/org/%s/secret/%s", org, secret))
		}
	}

	if gvc.Spec == nil {
		gvc.Spec = &client.GvcSpec{}
	}

	gvc.Spec.PullSecretLinks = &l
}

func resourceGvcRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceGvcRead")

	c := m.(*client.Client)
	gvc, code, err := c.GetGvc(d.Id())

	if code == 404 {
		return setGvc(d, nil, c.Org)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setGvc(d, gvc, c.Org)
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

	if err := d.Set("alias", gvc.Alias); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(gvc.Links, d); err != nil {
		return diag.FromErr(err)
	}

	if gvc.Spec != nil && gvc.Spec.Tracing != nil {
		if gvc.Spec.Tracing.Lightstep != nil {
			if err := d.Set("lightstep_tracing", flattenLightstepTracing(gvc.Spec.Tracing)); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		if err := d.Set("lightstep_tracing", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

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

func resourceGvcUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceGvcUpdate")

	if d.HasChanges("description", "locations", "tags", "domain", "pull_secrets", "lightstep_tracing") {

		c := m.(*client.Client)

		gvcToUpdate := client.Gvc{}
		gvcToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			gvcToUpdate.Description = GetDescriptionString(d.Get("description"), *gvcToUpdate.Name)
		}

		if d.HasChange("domain") {
			gvcToUpdate.Spec = &client.GvcSpec{}
			gvcToUpdate.Spec.Update = true
			gvcToUpdate.Spec.Domain = GetString(d.Get("domain"))
		}

		if d.HasChange("locations") {
			buildLocations(c.Org, d.Get("locations"), &gvcToUpdate)
		}

		if d.HasChange("pull_secrets") {
			buildPullSecrets(c.Org, d.Get("pull_secrets"), &gvcToUpdate)
		}

		if d.HasChange("tags") {
			gvcToUpdate.Tags = GetTagChanges(d)
		}

		if d.HasChange("lightstep_tracing") {
			traceArray := d.Get("lightstep_tracing").([]interface{})

			if len(traceArray) == 1 {

				if gvcToUpdate.Spec == nil {
					gvcToUpdate.Spec = &client.GvcSpec{}
				}

				gvcToUpdate.Spec.Tracing = buildLightStepTracing(traceArray)
			}
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
