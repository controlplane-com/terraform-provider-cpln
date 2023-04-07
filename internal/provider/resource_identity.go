package cpln

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentity() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceIdentityCreate,
		ReadContext:   resourceIdentityRead,
		UpdateContext: resourceIdentityUpdate,
		DeleteContext: resourceIdentityDelete,
		Schema: map[string]*schema.Schema{
			"gvc": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
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
			"status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_resource": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"agent_link": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: LinkValidator,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resolver_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ips": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ports": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"aws_access_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_account_link": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: LinkValidator,
						},
						"policy_refs": {
							Type:     schema.TypeSet,
							Optional: true,
							// ConflictsWith: []string{"aws_access_policy.role_name"},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"role_name": {
							Type:     schema.TypeString,
							Optional: true,
							// ConflictsWith: []string{"aws_access_policy.policy_refs"},
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								v := val.(string)

								re := regexp.MustCompile(`^([a-zA-Z0-9/+=,.@_-])+$`)

								if !re.MatchString(v) {
									errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, v))
									return
								}

								if len(v) > 65 {
									errs = append(errs, fmt.Errorf("%q length must at most 65 characters, got: %d", key, len(v)))
								}

								return
							},
						},
						// "trust_policy": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	MaxItems: 1,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"version": {
						// 				Type:     schema.TypeString,
						// 				Optional: true,
						// 			},
						// 			"statement": {
						// 				Type:     schema.TypeString,
						// 				Optional: true,
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},
			"gcp_access_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_account_link": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: LinkValidator,
						},
						"scopes": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "https://www.googleapis.com/auth/cloud-platform",
						},
						"service_account": {
							Type:     schema.TypeString,
							Optional: true,
							// ConflictsWith: []string{"gcp_access_policy.binding"},
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								v := val.(string)

								re := regexp.MustCompile(`^.*\.gserviceaccount\.com$`)

								if !re.MatchString(v) {
									errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, v))
									return
								}

								return
							},
						},
						"binding": {
							Type:     schema.TypeList,
							Optional: true,
							// ConflictsWith: []string{"gcp_access_policy.service_account"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"roles": {
										Type:     schema.TypeSet,
										Optional: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"azure_access_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_account_link": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: LinkValidator,
						},
						"role_assignment": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"roles": {
										Type:     schema.TypeSet,
										Optional: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"ngs_access_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_account_link": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: LinkValidator,
						},
						"pub": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     permResource(),
						},
						"sub": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     permResource(),
						},
						"resp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
									"ttl": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"subs": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"data": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"payload": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateIdentity,
		},
	}
}

func permResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deny": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func importStateIdentity(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.SplitN(d.Id(), ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected gvc:identity", d.Id())
	}

	d.Set("gvc", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceIdentityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceIdentityCreate")

	gvcName := d.Get("gvc").(string)

	identity := client.Identity{}
	identity.Name = GetString(d.Get("name"))
	identity.Description = GetString(d.Get("description"))
	identity.Tags = GetStringMap(d.Get("tags"))

	buildNetworkResources(d.Get("network_resource").([]interface{}), &identity)
	buildAwsIdentity(d.Get("aws_access_policy").([]interface{}), &identity, false)
	buildAzureIdentity(d.Get("azure_access_policy").([]interface{}), &identity, false)
	buildGcpIdentity(d.Get("gcp_access_policy").([]interface{}), &identity, false)
	buildNgsIdentity(d.Get("ngs_access_policy").([]interface{}), &identity, false)

	c := m.(*client.Client)
	newIdentity, code, err := c.CreateIdentity(identity, gvcName)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setIdentity(d, newIdentity, gvcName)
}

func buildNetworkResources(networkResources []interface{}, identity *client.Identity) {

	newNetworkResources := []client.NetworkResource{}

	for _, networkresource := range networkResources {

		n := networkresource.(map[string]interface{})

		newNetworkResource := client.NetworkResource{
			Name:      GetString(n["name"]),
			AgentLink: GetString(n["agent_link"]),
		}

		newNetworkResource.FQDN = GetString(n["fqdn"])
		newNetworkResource.ResolverIP = GetString(n["resolver_ip"])

		if n["ips"] != nil {
			ips := []string{}

			for _, value := range n["ips"].(*schema.Set).List() {
				ips = append(ips, value.(string))
			}

			if len(ips) > 0 {
				newNetworkResource.IPs = &ips
			}
		}

		if n["ports"] != nil {
			ports := []int{}

			for _, value := range n["ports"].(*schema.Set).List() {
				ports = append(ports, value.(int))
			}

			if len(ports) > 0 {
				newNetworkResource.Ports = &ports
			}
		}

		newNetworkResources = append(newNetworkResources, newNetworkResource)
	}

	identity.NetworkResources = &newNetworkResources
}

func buildAwsIdentity(awsIdentities []interface{}, identity *client.Identity, update bool) {

	if len(awsIdentities) == 1 {

		a := awsIdentities[0].(map[string]interface{})

		newAwsIdentity := &client.AwsIdentity{
			CloudAccountLink: GetString(a["cloud_account_link"]),
		}

		if a["role_name"] != nil {
			r := strings.TrimSpace(a["role_name"].(string))

			if r != "" {
				newAwsIdentity.RoleName = GetString(r)
			}
		}

		if a["policy_refs"] != nil {

			prs := a["policy_refs"].(*schema.Set).List()

			if len(prs) > 0 {
				pr := []string{}

				for _, p := range prs {
					pr = append(pr, p.(string))
				}

				newAwsIdentity.PolicyRefs = &pr
			}
		}

		if update {
			identity.AwsReplace = newAwsIdentity
		} else {
			identity.Aws = newAwsIdentity
		}
	} else {

		if update {
			if identity.Drop == nil {
				identity.Drop = &[]string{}
			}

			list := *identity.Drop
			newList := append(list, "aws")
			identity.Drop = &newList
		}
	}
}

func buildAzureIdentity(azureIdentities []interface{}, identity *client.Identity, update bool) {

	if len(azureIdentities) == 1 {

		a := azureIdentities[0].(map[string]interface{})

		newAzureIdentity := &client.AzureIdentity{
			CloudAccountLink: GetString(a["cloud_account_link"]),
		}

		if a["role_assignment"] != nil {

			ra := a["role_assignment"].([]interface{})

			if len(ra) > 0 {

				localRoleAssignments := []client.AzureRoleAssignment{}

				for _, r := range ra {

					if r != nil {

						localRa := client.AzureRoleAssignment{}

						rm := r.(map[string]interface{})

						localRa.Scope = GetString(rm["scope"].(string))

						if rm["roles"] != nil && len(rm["roles"].(*schema.Set).List()) > 0 {

							localRoles := []string{}

							for _, sr := range rm["roles"].(*schema.Set).List() {
								localRoles = append(localRoles, sr.(string))
							}

							localRa.Roles = &localRoles
						}

						localRoleAssignments = append(localRoleAssignments, localRa)
					}
				}

				newAzureIdentity.RoleAssignments = &localRoleAssignments
			}
		}

		if update {
			identity.AzureReplace = newAzureIdentity
		} else {
			identity.Azure = newAzureIdentity
		}
	} else {

		if update {
			if identity.Drop == nil {
				identity.Drop = &[]string{}
			}

			list := *identity.Drop
			newList := append(list, "azure")
			identity.Drop = &newList
		}
	}
}

func buildGcpIdentity(gcpIdentities []interface{}, identity *client.Identity, update bool) {

	if len(gcpIdentities) == 1 {

		a := gcpIdentities[0].(map[string]interface{})

		newGcpIdentity := &client.GcpIdentity{
			CloudAccountLink: GetString(a["cloud_account_link"]),
		}

		if a["scopes"] != nil {

			s := a["scopes"].(string)

			splitScope := strings.Split(s, ",")

			if len(splitScope) > 0 {
				newGcpIdentity.Scopes = &splitScope
			}
		}

		if a["service_account"] != nil {
			r := strings.TrimSpace(a["service_account"].(string))

			if r != "" {
				newGcpIdentity.ServiceAccount = GetString(r)
			}
		}

		if a["binding"] != nil {

			bs := a["binding"].([]interface{})

			if len(bs) > 0 {

				localBindings := []client.GcpBinding{}

				for _, b := range bs {

					if b != nil {

						localBs := client.GcpBinding{}

						rm := b.(map[string]interface{})

						localBs.Resource = GetString(rm["resource"].(string))

						if rm["roles"] != nil && len(rm["roles"].(*schema.Set).List()) > 0 {

							localRoles := []string{}

							for _, sr := range rm["roles"].(*schema.Set).List() {
								localRoles = append(localRoles, sr.(string))
							}

							localBs.Roles = &localRoles
						}

						localBindings = append(localBindings, localBs)
					}
				}

				newGcpIdentity.Bindings = &localBindings
			}
		}

		if update {
			identity.GcpReplace = newGcpIdentity
		} else {
			identity.Gcp = newGcpIdentity
		}
	} else {

		if update {
			if identity.Drop == nil {
				identity.Drop = &[]string{}
			}

			list := *identity.Drop
			newList := append(list, "gcp")
			identity.Drop = &newList
		}
	}
}

func buildPerm(perm []interface{}) *client.NgsPerm {

	if len(perm) > 0 {

		localPub := client.NgsPerm{}

		for _, p := range perm {

			if p != nil {

				rm := p.(map[string]interface{})

				if rm["allow"] != nil && len(rm["allow"].(*schema.Set).List()) > 0 {

					localAllow := []string{}

					for _, sr := range rm["allow"].(*schema.Set).List() {
						localAllow = append(localAllow, sr.(string))
					}

					localPub.Allow = &localAllow
				}

				if rm["deny"] != nil && len(rm["deny"].(*schema.Set).List()) > 0 {

					localDeny := []string{}

					for _, sr := range rm["deny"].(*schema.Set).List() {
						localDeny = append(localDeny, sr.(string))
					}

					localPub.Deny = &localDeny
				}
			}
		}

		return &localPub
	}

	return nil

}

func buildNgsIdentity(ngsIdentities []interface{}, identity *client.Identity, update bool) {

	if len(ngsIdentities) == 1 {

		a := ngsIdentities[0].(map[string]interface{})

		newNgsIdentity := &client.NgsIdentity{
			CloudAccountLink: GetString(a["cloud_account_link"]),
		}

		if a["pub"] != nil {
			newNgsIdentity.Pub = buildPerm(a["pub"].([]interface{}))
		}

		if a["sub"] != nil {
			newNgsIdentity.Sub = buildPerm(a["sub"].([]interface{}))
		}

		if a["resp"] != nil {

			ps := a["resp"].([]interface{})

			if len(ps) > 0 {

				localResp := client.NgsResp{}

				for _, p := range ps {

					rm := p.(map[string]interface{})

					localResp.Max = GetInt(rm["max"].(int))
					localResp.TTL = GetString(rm["ttl"].(string))
				}

				newNgsIdentity.Resp = &localResp
			}
		}

		if a["subs"] != nil {
			subs := a["subs"].(int)
			newNgsIdentity.Subs = GetInt(subs)
		}

		if a["data"] != nil {
			data := a["data"].(int)
			newNgsIdentity.Data = GetInt(data)
		}

		if a["payload"] != nil {
			payload := a["payload"].(int)
			newNgsIdentity.Payload = GetInt(payload)
		}

		if update {
			identity.NgsReplace = newNgsIdentity
		} else {
			identity.Ngs = newNgsIdentity
		}
	} else {

		if update {
			if identity.Drop == nil {
				identity.Drop = &[]string{}
			}

			list := *identity.Drop
			newList := append(list, "ngs")
			identity.Drop = &newList
		}
	}
}

func flattenNetworkResources(networkResources *[]client.NetworkResource) []interface{} {

	if networkResources != nil && len(*networkResources) > 0 {

		nrs := make([]interface{}, len(*networkResources))

		for i, networkResource := range *networkResources {

			nr := make(map[string]interface{})

			nr["name"] = networkResource.Name
			nr["agent_link"] = networkResource.AgentLink

			if networkResource.FQDN != nil {
				nr["fqdn"] = networkResource.FQDN
			}

			if networkResource.ResolverIP != nil {
				nr["resolver_ip"] = networkResource.ResolverIP
			}

			if networkResource.IPs != nil && len(*networkResource.IPs) > 0 {
				nr["ips"] = []interface{}{}

				for _, ip := range *networkResource.IPs {
					nr["ips"] = append(nr["ips"].([]interface{}), ip)
				}
			}

			if networkResource.Ports != nil && len(*networkResource.Ports) > 0 {
				nr["ports"] = []interface{}{}

				for _, port := range *networkResource.Ports {
					nr["ports"] = append(nr["ports"].([]interface{}), port)
				}
			}

			nrs[i] = nr
		}

		return nrs
	}

	return make([]interface{}, 0)
}

func flattenAwsIdentity(awsIdentity *client.AwsIdentity) []interface{} {

	if awsIdentity != nil {

		output := make(map[string]interface{})

		output["cloud_account_link"] = *awsIdentity.CloudAccountLink

		if awsIdentity.PolicyRefs != nil && len(*awsIdentity.PolicyRefs) > 0 {

			pr := []interface{}{}

			for _, p := range *awsIdentity.PolicyRefs {
				pr = append(pr, p)
			}
			output["policy_refs"] = pr
		}

		if awsIdentity.RoleName != nil && strings.TrimSpace(*awsIdentity.RoleName) != "" {
			output["role_name"] = *awsIdentity.RoleName
		}

		return []interface{}{
			output,
		}
	}

	return nil
}

func flattenAzureIdentity(azureIdentity *client.AzureIdentity) []interface{} {

	if azureIdentity != nil {

		output := make(map[string]interface{})

		output["cloud_account_link"] = *azureIdentity.CloudAccountLink

		if azureIdentity.RoleAssignments != nil && len(*azureIdentity.RoleAssignments) > 0 {

			roleAssignment := []interface{}{}

			for _, r := range *azureIdentity.RoleAssignments {

				ra := make(map[string]interface{})

				if r.Scope != nil {
					ra["scope"] = *r.Scope
				}

				if r.Roles != nil && len(*r.Roles) > 0 {
					roles := []interface{}{}

					for _, rr := range *r.Roles {
						roles = append(roles, rr)
					}

					ra["roles"] = roles
				}

				roleAssignment = append(roleAssignment, ra)
			}

			output["role_assignment"] = roleAssignment
		}

		return []interface{}{
			output,
		}
	}

	return nil
}

func flattenGcpIdentity(gcpIdentity *client.GcpIdentity) []interface{} {

	if gcpIdentity != nil {

		output := make(map[string]interface{})

		output["cloud_account_link"] = *gcpIdentity.CloudAccountLink

		if gcpIdentity.Scopes != nil && len(*gcpIdentity.Scopes) > 0 {
			joinScopes := strings.Join(*gcpIdentity.Scopes, ",")
			output["scopes"] = joinScopes
		}

		if gcpIdentity.ServiceAccount != nil && strings.TrimSpace(*gcpIdentity.ServiceAccount) != "" {
			output["service_account"] = *gcpIdentity.ServiceAccount
		}

		if gcpIdentity.Bindings != nil && len(*gcpIdentity.Bindings) > 0 {

			bindings := []interface{}{}

			for _, b := range *gcpIdentity.Bindings {

				bs := make(map[string]interface{})

				if b.Resource != nil {
					bs["resource"] = *b.Resource
				}

				if b.Roles != nil && len(*b.Roles) > 0 {
					roles := []interface{}{}

					for _, rr := range *b.Roles {
						roles = append(roles, rr)
					}

					bs["roles"] = roles
				}

				bindings = append(bindings, bs)
			}

			output["binding"] = bindings
		}

		return []interface{}{
			output,
		}
	}

	return nil
}

func flattenPerm(perm *client.NgsPerm) []interface{} {

	if perm != nil {

		ps := []interface{}{}
		bs := make(map[string]interface{})

		if perm.Allow != nil && len(*perm.Allow) > 0 {
			allowDeny := []interface{}{}

			for _, ad := range *perm.Allow {
				allowDeny = append(allowDeny, ad)
			}

			bs["allow"] = allowDeny
		}

		if perm.Deny != nil && len(*perm.Deny) > 0 {
			allowDeny := []interface{}{}

			for _, ad := range *perm.Deny {
				allowDeny = append(allowDeny, ad)
			}

			bs["deny"] = allowDeny
		}

		ps = append(ps, bs)

		return ps
	}

	return nil
}

func flattenNgsIdentity(ngsIdentity *client.NgsIdentity) []interface{} {

	if ngsIdentity != nil {

		output := make(map[string]interface{})

		output["cloud_account_link"] = *ngsIdentity.CloudAccountLink
		output["pub"] = flattenPerm(ngsIdentity.Pub)
		output["sub"] = flattenPerm(ngsIdentity.Sub)

		if ngsIdentity.Resp != nil {

			rs := make(map[string]interface{})

			if ngsIdentity.Resp.Max != nil {
				rs["max"] = *ngsIdentity.Resp.Max
			}

			if ngsIdentity.Resp.TTL != nil {
				rs["ttl"] = *ngsIdentity.Resp.TTL
			}

			resps := []interface{}{}
			resps = append(resps, rs)

			output["resp"] = resps
		}

		if ngsIdentity.Subs != nil {
			output["subs"] = *ngsIdentity.Subs
		}

		if ngsIdentity.Data != nil {
			output["data"] = *ngsIdentity.Data
		}

		if ngsIdentity.Payload != nil {
			output["payload"] = *ngsIdentity.Payload
		}

		return []interface{}{
			output,
		}
	}

	return nil
}

func resourceIdentityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceIdentityRead")

	gvcName := d.Get("gvc").(string)

	c := m.(*client.Client)
	identity, code, err := c.GetIdentity(d.Id(), gvcName)

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setIdentity(d, identity, gvcName)
}

func setIdentity(d *schema.ResourceData, identity *client.Identity, gvcName string) diag.Diagnostics {

	if identity == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*identity.Name)

	if err := SetBase(d, identity.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", flattenIdentityStatus(identity.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("network_resource", flattenNetworkResources(identity.NetworkResources)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("aws_access_policy", flattenAwsIdentity(identity.Aws)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("azure_access_policy", flattenAzureIdentity(identity.Azure)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("gcp_access_policy", flattenGcpIdentity(identity.Gcp)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ngs_access_policy", flattenNgsIdentity(identity.Ngs)); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(identity.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenIdentityStatus(status *client.IdentityStatus) interface{} {

	if status != nil && status.ObjectName != nil && strings.TrimSpace(*status.ObjectName) != "" {
		fs := make(map[string]interface{})
		fs["objectName"] = status.ObjectName
		return fs
	}

	return nil
}

func resourceIdentityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceIdentityUpdate")

	if d.HasChanges("description", "tags", "network_resource", "aws_access_policy", "azure_access_policy", "gcp_access_policy", "ngs_access_policy") {

		gvcName := d.Get("gvc").(string)

		identityToUpdate := client.Identity{}
		identityToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			identityToUpdate.Description = GetDescriptionString(d.Get("description"), *identityToUpdate.Name)
		}

		if d.HasChange("tags") {
			identityToUpdate.Tags = GetTagChanges(d)
		}

		if d.HasChange("network_resource") {
			buildNetworkResources(d.Get("network_resource").([]interface{}), &identityToUpdate)
		}

		if d.HasChange("aws_access_policy") {
			buildAwsIdentity(d.Get("aws_access_policy").([]interface{}), &identityToUpdate, true)
		}

		if d.HasChange("azure_access_policy") {
			buildAzureIdentity(d.Get("azure_access_policy").([]interface{}), &identityToUpdate, true)
		}

		if d.HasChange("gcp_access_policy") {
			buildGcpIdentity(d.Get("gcp_access_policy").([]interface{}), &identityToUpdate, true)
		}

		if d.HasChange("ngs_access_policy") {
			buildNgsIdentity(d.Get("ngs_access_policy").([]interface{}), &identityToUpdate, true)
		}

		c := m.(*client.Client)
		updatedIdentity, _, err := c.UpdateIdentity(identityToUpdate, gvcName)
		if err != nil {
			return diag.FromErr(err)
		}

		return setIdentity(d, updatedIdentity, gvcName)
	}

	return nil
}

func resourceIdentityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// log.Printf("[INFO] Method: resourceIdentityDelete")

	c := m.(*client.Client)
	err := c.DeleteIdentity(d.Id(), d.Get("gvc").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
