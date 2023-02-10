package cpln

import (
	"context"
	"fmt"
	"strings"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicy() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
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
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				Elem:         StringSchema(),
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_kind": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: KindValidator,
			},
			"target_links": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 200,
				Elem:     StringSchema(),
			},
			"target_query": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     QuerySchemaResource(),
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "" && v != "all" {
						errs = append(errs, fmt.Errorf("%q must be set to 'all', got: %s", key, v))
					}

					return
				},
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"binding": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem:     BindingResource(),
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func BindingResource() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"permissions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     StringSchema(),
			},
			"principal_links": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 200,
				Elem:     StringSchema(),
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	policy := client.Policy{}
	policy.Name = GetString(d.Get("name"))
	policy.Description = GetString(d.Get("description"))
	policy.Tags = GetStringMap(d.Get("tags"))
	policy.TargetKind = GetString(d.Get("target_kind"))
	policy.Target = GetString(d.Get("target"))

	c := m.(*client.Client)

	buildTargetLinks(c.Org, *policy.TargetKind, d.Get("target_links"), &policy)
	policy.TargetQuery = BuildQueryHelper(*policy.TargetKind, d.Get("target_query"))
	buildBindings(c.Org, d.Get("binding"), &policy)

	newPolicy, code, err := c.CreatePolicy(policy)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setPolicy(c.Org, d, newPolicy)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	policy, code, err := c.GetPolicy(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setPolicy(c.Org, d, policy)
}

func setPolicy(org string, d *schema.ResourceData, policy *client.Policy) diag.Diagnostics {

	if policy == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*policy.Name)

	if err := SetBase(d, policy.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("target_kind", policy.TargetKind); err != nil {
		return diag.FromErr(err)
	}

	targetLinks := flattenTargetLinks(policy.TargetLinks)

	if err := d.Set("target_links", targetLinks); err != nil {
		return diag.FromErr(err)
	}

	targetQuery, err := FlattenQueryHelper(policy.TargetQuery)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("target_query", targetQuery); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("target", policy.Target); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("origin", policy.Origin); err != nil {
		return diag.FromErr(err)
	}

	bindings := flattenBindings(org, policy.Bindings)

	if err := d.Set("binding", bindings); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(policy.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "target_links", "target_query", "target", "binding") {

		c := m.(*client.Client)

		policyToUpdate := client.Policy{}
		policyToUpdate.Update = true
		policyToUpdate.Name = GetString(d.Get("name"))

		policyToUpdate.Target = GetString(d.Get("target"))
		policyToUpdate.TargetQuery = BuildQueryHelper("user", d.Get("target_query"))
		buildTargetLinks(c.Org, d.Get("target_kind").(string), d.Get("target_links"), &policyToUpdate)
		buildBindings(c.Org, d.Get("binding"), &policyToUpdate)

		if d.HasChange("description") {
			policyToUpdate.Description = GetDescriptionString(d.Get("description"), *policyToUpdate.Name)
		}

		if d.HasChange("tags") {
			policyToUpdate.Tags = GetTagChanges(d)
		}

		updatedPolicy, _, err := c.UpdatePolicy(policyToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setPolicy(c.Org, d, updatedPolicy)
	}

	return nil
}

func buildTargetLinks(org, kind string, targets interface{}, policy *client.Policy) {

	targetLinks := []string{}

	if targets != nil {

		targetArray := targets.(*schema.Set)

		for _, t := range targetArray.List() {
			targetLinks = append(targetLinks, fmt.Sprintf("/org/%s/%s/%s", org, kind, t))
		}
	}

	policy.TargetLinks = &targetLinks
}

func buildBindings(org string, bindings interface{}, policy *client.Policy) {

	bindingsArray := []client.Binding{}

	if bindings != nil {

		bSet := bindings.(*schema.Set)

		bArray := bSet.List()

		for _, binding := range bArray {

			permissions := []string{}
			principalLinks := []string{}

			b := binding.(map[string]interface{})

			pArray := b["permissions"].(*schema.Set)
			plArray := b["principal_links"].(*schema.Set)

			for _, p := range pArray.List() {
				permissions = append(permissions, p.(string))
			}

			for _, b := range plArray.List() {

				principal := fmt.Sprintf(`/org/%s/%s`, org, b.(string))
				principalLinks = append(principalLinks, principal)
			}

			if len(permissions) > 0 || len(principalLinks) > 0 {

				localBinding := client.Binding{}

				if len(permissions) > 0 {
					localBinding.Permissions = &permissions
				}

				if len(principalLinks) > 0 {
					localBinding.PrincipalLinks = &principalLinks
				}

				bindingsArray = append(bindingsArray, localBinding)
			}
		}
	}

	policy.Bindings = &bindingsArray
}

func flattenTargetLinks(targetLinks *[]string) []interface{} {

	if targetLinks == nil || len(*targetLinks) < 1 {
		return nil
	}

	output := []interface{}{}

	for _, m := range *targetLinks {
		output = append(output, m[strings.LastIndexAny(m, "/")+1:])
	}

	return output
}

func flattenBindings(org string, bindings *[]client.Binding) []interface{} {

	if bindings == nil || len(*bindings) < 1 {
		return nil
	}

	flatBindings := []interface{}{}

	for _, binding := range *bindings {

		b := make(map[string]interface{})

		permissions := []interface{}{}

		for _, p := range *binding.Permissions {
			permissions = append(permissions, p)
		}

		if len(permissions) > 0 {
			b["permissions"] = permissions
		}

		principalLinks := []interface{}{}

		for _, p := range *binding.PrincipalLinks {

			principal := strings.TrimPrefix(p, fmt.Sprintf(`/org/%s/`, org))
			principalLinks = append(principalLinks, principal)
		}

		if len(principalLinks) > 0 {
			b["principal_links"] = principalLinks
		}

		if len(permissions) > 0 || len(principalLinks) > 0 {
			flatBindings = append(flatBindings, b)
		}
	}

	return flatBindings
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourcePolicyDelete")

	c := m.(*client.Client)
	err := c.DeletePolicy(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
