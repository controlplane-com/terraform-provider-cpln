package cpln

import (
	"context"
	"fmt"
	"regexp"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var cloudProvidersNames = []string{
	"aws", "azure", "gcp", "ngs",
}

func CloudAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpln_id": {
			Type:        schema.TypeString,
			Description: "The ID, in GUID format, of the Cloud Account.",
			Computed:    true,
		},
		"name": {
			Type:         schema.TypeString,
			Description:  "Name of the Cloud Account.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: NameValidator,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the Cloud Account.",
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
		"self_link": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aws": {
			Type:         schema.TypeList,
			Description: "A leading cloud computing platform provided by 
				Amazon, offering a wide range of services for computing 
				power, storage, and other functionalities on a pay-as-you-go 
				basis.",
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: cloudProvidersNames,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role_arn": {
						Type:        schema.TypeString,
						Description: "Amazon Resource Name (ARN) Role.",
						Required:    true,
						ForceNew:    true,
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
			Type:         schema.TypeList,
			Description: "Microsoft's cloud computing platform, providing
				various services such as virtual machines, databases, and AI tools, 
				facilitating scalable and flexible solutions for businesses.",
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: cloudProvidersNames,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"secret_link": {
						Type:         schema.TypeString,
						Description:  "Full link to an Azure secret. (e.g., /org/ORG_NAME/secret/AZURE_SECRET).",
						Required:     true,
						ForceNew:     true,
						ValidateFunc: LinkValidator,
					},
				},
			},
		},
		"gcp": {
			Type:         schema.TypeList,
			Description: "Google's comprehensive suite of cloud computing 
				services, offering infrastructure, data storage, machine learning, 
				and application development tools for building and scaling applications
				and services."
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: cloudProvidersNames,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"project_id": {
						Type:        schema.TypeString,
						Description: "GCP project ID. Obtained from the GCP cloud console.",
						Required:    true,
						ForceNew:    true,
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
		"ngs": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: cloudProvidersNames,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"secret_link": {
						Type:         schema.TypeString,
						Description:  "Full link to a NATS Account Secret secret. (e.g., /org/ORG_NAME/secret/NATS_ACCOUNT_SECRET).",
						Required:     true,
						ForceNew:     true,
						ValidateFunc: LinkValidator,
					},
				},
			},
		},
		"gcp_service_account_name": {
			Type:        schema.TypeString,
			Description: "GCP service account name used during the configuration of the cloud account at GCP.",
			Computed:    true,
		},
		"gcp_roles": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type:        schema.TypeString,
				Description: "GCP roles used during the configuration of the cloud account at GCP.",
			},
		},
	}
}

func resourceCloudAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountCreate,
		ReadContext:   resourceCloudAccountRead,
		UpdateContext: resourceCloudAccountUpdate,
		DeleteContext: resourceCloudAccountDelete,
		Schema:        CloudAccountSchema(),
		Importer:      &schema.ResourceImporter{},
	}
}

func getProvider(d *schema.ResourceData) *string {

	for _, s := range cloudProvidersNames {
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
		return diag.FromErr(fmt.Errorf("missing or unable to extract provider. must be 'aws', 'azure', 'gcp' or 'ngs'"))
	}

	aws := d.Get("aws").([]interface{})
	azure := d.Get("azure").([]interface{})
	gcp := d.Get("gcp").([]interface{})
	ngs := d.Get("ngs").([]interface{})

	if aws != nil || azure != nil || gcp != nil || ngs != nil {
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

	if len(ngs) > 0 {
		ngs_secret_link := ngs[0].(map[string]interface{})
		ca.Data.SecretLink = GetString(ngs_secret_link["secret_link"])
	}

	c := m.(*client.Client)
	newCa, code, err := c.CreateCloudAccount(ca)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setCloudAccount(d, newCa, c.Org)
}

func resourceCloudAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	ca, code, err := c.GetCloudAccount(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setCloudAccount(d, ca, c.Org)
}

func setCloudAccount(d *schema.ResourceData, ca *client.CloudAccount, orgName string) diag.Diagnostics {

	if ca == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*ca.Name)

	if err := SetBase(d, ca.Base); err != nil {
		return diag.FromErr(err)
	}

	if ca.Data != nil {

		if ca.Data.RoleArn != nil {

			aws_role_arn := make(map[string]interface{})
			aws_role_arn["role_arn"] = ca.Data.RoleArn

			aws := []interface{}{
				aws_role_arn,
			}

			if err := d.Set("aws", aws); err != nil {
				return diag.FromErr(err)
			}
		}

		if ca.Data.SecretLink != nil {

			provider_secret_link := make(map[string]interface{})
			provider_secret_link["secret_link"] = ca.Data.SecretLink

			provider := []interface{}{
				provider_secret_link,
			}

			if err := d.Set(*ca.Provider, provider); err != nil {
				return diag.FromErr(err)
			}
		}

		if ca.Data.ProjectId != nil {

			gcp_project_id := make(map[string]interface{})
			gcp_project_id["project_id"] = ca.Data.ProjectId

			gcp := []interface{}{
				gcp_project_id,
			}

			serviceAccountName := "cpln-" + orgName + "@cpln-prod01.iam.gserviceaccount.com"
			if err := d.Set("gcp_service_account_name", serviceAccountName); err != nil {
				return diag.FromErr(err)
			}

			roles := []interface{}{"roles/viewer", "roles/iam.serviceAccountAdmin", "roles/iam.serviceAccountTokenCreator"}
			if err := d.Set("gcp_roles", roles); err != nil {
				return diag.FromErr(err)
			}

			if err := d.Set("gcp", gcp); err != nil {
				return diag.FromErr(err)
			}
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

		return setCloudAccount(d, updatedCa, c.Org)
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
