package cpln

import (
	"context"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMk8sKubeconfig() *schema.Resource {

	var profileOrServiceAccount []string = []string{"profile", "service_account"}

	return &schema.Resource{
		CreateContext: resourceMk8sKubeconfigCreate,
		ReadContext:   resourceMk8sKubeconfigRead,
		DeleteContext: resourceMk8sKubeconfigDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the MK8s to create the Kubeconfig for.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"profile": {
				Type:         schema.TypeString,
				Description:  "Profile name to extract the token from.",
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: profileOrServiceAccount,
			},
			"service_account": {
				Type:         schema.TypeString,
				Description:  "A service account to add a key to.",
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: profileOrServiceAccount,
				ValidateFunc: NameValidator,
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Description: "The Kubeconfig of your MK8s cluster in YAML format.",
				Sensitive:   true,
				Computed:    true,
			},
		},
	}
}

func resourceMk8sKubeconfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	// Get values from attributes
	mk8sName := GetString(d.Get("name"))
	profile := GetString(d.Get("profile"))
	serviceAccountName := GetString(d.Get("service_account"))

	// Create the Kubeconfig
	kubeconfig, err := c.CreateMk8sKubeconfig(*mk8sName, profile, serviceAccountName)
	if err != nil {
		return diag.FromErr(err)
	}

	return setMk8sKubeconfig(d, *mk8sName, profile, serviceAccountName, *kubeconfig)
}

func resourceMk8sKubeconfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceMk8sKubeconfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func setMk8sKubeconfig(d *schema.ResourceData, mk8sName string, profile *string, serviceAccount *string, kubeconfig string) diag.Diagnostics {
	// Set unique id to identify this resource
	d.SetId(getMk8sKubeconfigUnqiueId(mk8sName, profile, serviceAccount))

	// Set kubeconfig
	d.Set("kubeconfig", kubeconfig)

	return nil
}

// Helpers //

func getMk8sKubeconfigUnqiueId(mk8sName string, profile *string, serviceAccount *string) string {
	// Create an identity based on the profile name
	if profile != nil && len(*profile) != 0 {
		return fmt.Sprintf("%s:profile:%s", mk8sName, *profile)
	}

	// If profile is empty, then service account shouldn't, let's create an identity based on it
	return fmt.Sprintf("%s:service_account:%s", mk8sName, *serviceAccount)
}
