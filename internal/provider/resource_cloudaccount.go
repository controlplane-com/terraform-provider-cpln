package cpln

import (
	"context"
	"fmt"
	"regexp"
	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountCreate,
		ReadContext:   resourceCloudAccountRead,
		UpdateContext: resourceCloudAccountUpdate,
		DeleteContext: resourceCloudAccountDelete,
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
			"aws": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"azure", "gcp"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								v := val.(string)

								re := regexp.MustCompile(`^arn:(aws|aws-us-gov|aws-cn):iam::[0-9]+:role\/.+`)

								if !re.MatchString(v) {
									errs = append(errs, fmt.Errorf("%q is invalid, got: '%s'", key, v))
								}

								return
							},
						},
					},
				},
			},
			"azure": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"aws", "gcp"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_link": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: LinkValidator,
						},
					},
				},
			},
			"gcp": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"aws", "azure"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								v := val.(string)

								re := regexp.MustCompile(`[a-z]([a-z]|-|[0-9])+`)

								if !re.MatchString(v) {
									errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, v))
									return
								}

								if len(v) < 6 || len(v) > 30 {
									errs = append(errs, fmt.Errorf("%q length must be between 6 and 30 characters, got: %d", key, len(v)))
								}

								return
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func getProvider(d *schema.ResourceData) *string {

	secrets := []string{
		"aws", "azure", "gcp",
	}

	for _, s := range secrets {
		if _, v := d.GetOk(s); v {
			return &s
		}
	}

	return nil
}

func resourceCloudAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	ca := client.CloudAccount{}
	ca.Name = GetString(d.Get("name"))
	ca.Description = GetDescriptionString(d.Get("description"), *ca.Name)
	ca.Tags = GetStringMap(d.Get("tags"))

	if ca.Provider = getProvider(d); ca.Provider == nil {
		return diag.FromErr(fmt.Errorf("missing or unable to extract provider. must be 'aws', 'azure', or 'gcp'"))
	}

	aws := d.Get("aws").([]interface{})
	azure := d.Get("azure").([]interface{})
	gcp := d.Get("gcp").([]interface{})

	if aws != nil || azure != nil || gcp != nil {
		ca.Data = &client.CloudAccountConfig{}
	}

	if len(aws) > 0 {
		aws_role_arn := aws[0].(map[string]interface{})
		ca.Data.RoleArn = GetString(aws_role_arn["role_arn"])
	}

	if len(azure) > 0 {
		azure_secret_link := azure[0].(map[string]interface{})
		ca.Data.SecretLink = GetString(azure_secret_link["secret_link"])
	}

	if len(gcp) > 0 {
		gcp_project_id := gcp[0].(map[string]interface{})
		ca.Data.ProjectId = GetString(gcp_project_id["project_id"])
	}

	c := m.(*client.Client)
	newCa, code, err := c.CreateCloudAccount(ca)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setCloudAccount(d, newCa)
}

func resourceCloudAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	ca, code, err := c.GetCloudAccount(d.Id())

	if code == 404 {
		return setGvc(d, nil, c.Org)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setCloudAccount(d, ca)
}

func setCloudAccount(d *schema.ResourceData, ca *client.CloudAccount) diag.Diagnostics {

	if ca == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*ca.Name)

	if err := SetBase(d, ca.Base); err != nil {
		return diag.FromErr(err)
	}

	if ca.Data != nil && ca.Data.RoleArn != nil {

		aws_role_arn := make(map[string]interface{})
		aws_role_arn["role_arn"] = ca.Data.RoleArn

		aws := []interface{}{
			aws_role_arn,
		}

		if err := d.Set("aws", aws); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("aws", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	if ca.Data != nil && ca.Data.SecretLink != nil {

		azure_secret_link := make(map[string]interface{})
		azure_secret_link["secret_link"] = ca.Data.SecretLink

		azure := []interface{}{
			azure_secret_link,
		}

		if err := d.Set("azure", azure); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("azure", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	if ca.Data != nil && ca.Data.ProjectId != nil {

		gcp_project_id := make(map[string]interface{})
		gcp_project_id["project_id"] = ca.Data.ProjectId

		gcp := []interface{}{
			gcp_project_id,
		}

		if err := d.Set("gcp", gcp); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("gcp", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := SetSelfLink(ca.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCloudAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags") {

		caToUpdate := client.CloudAccount{}
		caToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			caToUpdate.Description = GetDescriptionString(d.Get("description"), *caToUpdate.Name)
		}

		if d.HasChange("tags") {
			caToUpdate.Tags = GetTagChanges(d)
		}

		c := m.(*client.Client)
		updatedCa, _, err := c.UpdateCloudAccount(caToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setCloudAccount(d, updatedCa)
	}

	return nil
}

func resourceCloudAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	err := c.DeleteCloudAccount(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
