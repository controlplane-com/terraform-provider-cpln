package cpln

import (
	"context"
	"fmt"
	"strings"
	client "terraform-provider-cpln/internal/provider/client"

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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"user_ids_and_emails": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service_accounts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_query": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     QuerySchemaResource(),
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
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

	newGroup, code, err := c.CreateGroup(group)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setGroup(d, c.Org, newGroup)
}

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

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	group, code, err := c.GetGroup(d.Id())

	if code == 404 {
		return setGvc(d, nil, c.Org)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setGroup(d, c.Org, group)
}

func setGroup(d *schema.ResourceData, org string, group *client.Group) diag.Diagnostics {

	if group == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*group.Name)

	if err := SetBase(d, group.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("origin", group.Origin); err != nil {
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

	if err := SetSelfLink(group.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

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

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "user_ids_and_emails", "service_accounts", "member_query") {

		c := m.(*client.Client)

		groupToUpdate := client.Group{}
		groupToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			groupToUpdate.Description = GetDescriptionString(d.Get("description"), *groupToUpdate.Name)
		}

		if d.HasChange("tags") {
			groupToUpdate.Tags = GetTagChanges(d)
		}

		if d.HasChange("user_ids_and_emails") || d.HasChange("service_accounts") || d.HasChange("member_query") {

			userMembers := d.Get("user_ids_and_emails")
			serviceAccountMembers := d.Get("service_accounts")
			queryMembers := d.Get("member_query")

			buildMemberLinks(c.Org, userMembers, serviceAccountMembers, &groupToUpdate)
			groupToUpdate.MemberQuery = BuildQueryHelper("user", queryMembers)
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
