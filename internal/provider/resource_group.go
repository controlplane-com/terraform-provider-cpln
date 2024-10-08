package cpln

import (
	"context"
	"fmt"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the Group.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the Group.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of Group.",
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
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
			"user_ids_and_emails": {
				Type:        schema.TypeSet,
				Description: "List of either the user ID or email address for a user that exists within the configured org. Group membership will fail if the user ID / email does not exist within the org.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service_accounts": {
				Type:        schema.TypeSet,
				Description: "List of service accounts that exists within the configured org. Group membership will fail if the service account does not exits within the org.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: "Fully qualified link to the this group.",
				Computed:    true,
			},
			"member_query": {
				Type:        schema.TypeList,
				Description: "A predefined set of criteria or conditions used to query and retrieve members within the group.",
				Optional:    true,
				MaxItems:    1,
				Elem:        QuerySchemaResource(),
			},
			"identity_matcher": {
				Type:        schema.TypeList,
				Description: "Executes the expression against the users' claims to decide whether a user belongs to this group. This method is useful for managing the grouping of users logged-in with SAML providers.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:        schema.TypeString,
							Description: "Executes the expression against the users' claims to decide whether a user belongs to this group. This method is useful for managing the grouping of users logged in with SAML providers.",
							Required:    true,
						},
						"language": {
							Type:        schema.TypeString,
							Description: "Language of the expression. Either `jmespath` or `javascript`. Default: `jmespath`.",
							Optional:    true,
							Default:     "jmespath",
						},
					},
				},
			},
			"origin": {
				Type:        schema.TypeString,
				Description: "Origin of the service account. Either `builtin` or `default`.",
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	group := client.Group{}
	group.Name = GetString(d.Get("name"))
	group.Description = GetString(d.Get("description"))
	group.Tags = GetStringMap(d.Get("tags"))

	c := m.(*client.Client)

	buildMemberLinks(c.Org, d.Get("user_ids_and_emails"), d.Get("service_accounts"), &group)
	group.MemberQuery = BuildQueryHelper("user", d.Get("member_query"))
	group.IdentityMatcher = buildIdentityMatcher(d.Get("identity_matcher").([]interface{}))

	newGroup, code, err := c.CreateGroup(group)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setGroup(d, c.Org, newGroup)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	group, code, err := c.GetGroup(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setGroup(d, c.Org, group)
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "user_ids_and_emails", "service_accounts", "member_query", "identity_matcher") {

		c := m.(*client.Client)

		groupToUpdate := client.Group{}
		groupToUpdate.Name = GetString(d.Get("name"))
		groupToUpdate.Description = GetDescriptionString(d.Get("description"), *groupToUpdate.Name)
		groupToUpdate.Tags = GetTagChanges(d)

		if d.HasChange("user_ids_and_emails") || d.HasChange("service_accounts") || d.HasChange("member_query") {

			userMembers := d.Get("user_ids_and_emails")
			serviceAccountMembers := d.Get("service_accounts")
			queryMembers := d.Get("member_query")

			buildMemberLinks(c.Org, userMembers, serviceAccountMembers, &groupToUpdate)
			groupToUpdate.MemberQuery = BuildQueryHelper("user", queryMembers)
		}

		if d.HasChange("identity_matcher") {
			groupToUpdate.IdentityMatcher = buildIdentityMatcher(d.Get("identity_matcher").([]interface{}))
		}

		updatedGroup, _, err := c.UpdateGroup(groupToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setGroup(d, c.Org, updatedGroup)
	}

	return nil
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	err := c.DeleteGroup(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

/*** Build Functions ***/
func buildMemberLinks(org string, users, serviceAccounts interface{}, group *client.Group) {

	memberLinks := []string{}

	if users != nil {
		for _, u := range users.(*schema.Set).List() {
			memberLinks = append(memberLinks, fmt.Sprintf("/org/%s/user/%s", org, u.(string)))
		}
	}

	if serviceAccounts != nil {
		for _, s := range serviceAccounts.(*schema.Set).List() {
			memberLinks = append(memberLinks, fmt.Sprintf("/org/%s/serviceaccount/%s", org, s.(string)))
		}
	}

	group.MemberLinks = &memberLinks
}

func buildIdentityMatcher(specs []interface{}) *client.IdentityMatcher {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	result := client.IdentityMatcher{}

	if spec["expression"] != nil {
		result.Expression = GetString(spec["expression"].(string))
	}

	if spec["language"] != nil {
		result.Language = GetString(spec["language"].(string))
	}

	return &result
}

/*** Flatten Functions ***/
func flattenMemberLinks(org string, memberLinks *[]string) ([]interface{}, []interface{}, error) {

	if org == "" || memberLinks == nil {
		return nil, nil, fmt.Errorf("org is empty or member links is nil")
	}

	linksPrefix := fmt.Sprintf("/org/%s/", org)

	userIDs := []interface{}{}
	userIDPrefix := linksPrefix + "user/"

	serviceAccounts := []interface{}{}
	serviceAccountPrefix := linksPrefix + "serviceaccount/"

	for _, m := range *memberLinks {

		if strings.HasPrefix(m, userIDPrefix) {
			userIDs = append(userIDs, strings.TrimPrefix(m, userIDPrefix))
		} else if strings.HasPrefix(m, serviceAccountPrefix) {
			serviceAccounts = append(serviceAccounts, strings.TrimPrefix(m, serviceAccountPrefix))
		}
	}

	return userIDs, serviceAccounts, nil
}

func flattenIdentityMatcher(identityMatcher *client.IdentityMatcher) []interface{} {
	if identityMatcher == nil {
		return nil
	}

	result := make(map[string]interface{})

	if identityMatcher.Expression != nil {
		result["expression"] = *identityMatcher.Expression
	}

	if identityMatcher.Language != nil {
		result["language"] = *identityMatcher.Language
	}

	return []interface{}{
		result,
	}
}

/*** Helper Functions ***/
func setGroup(d *schema.ResourceData, org string, group *client.Group) diag.Diagnostics {

	if group == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*group.Name)

	if err := SetBase(d, group.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(group.Links, d); err != nil {
		return diag.FromErr(err)
	}

	userIDs, serviceAccounts, err := flattenMemberLinks(org, group.MemberLinks)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_ids_and_emails", userIDs); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("service_accounts", serviceAccounts); err != nil {
		return diag.FromErr(err)
	}

	mqList, err := FlattenQueryHelper(group.MemberQuery)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("member_query", mqList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("identity_matcher", flattenIdentityMatcher(group.IdentityMatcher)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
