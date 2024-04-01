package cpln

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretDataObjectsNames = []string{
	"aws", "azure_connector", "azure_sdk", "docker", "dictionary", "ecr",
	"gcp", "keypair", "opaque", "tls", "userpass", "nats_account",
}

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the Secret.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the secret.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the Secret.",
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
				Type:        schema.TypeString,
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"dictionary": {
				Type:         schema.TypeMap,
				Description:  "List of unique key-value pairs. [Reference Page](https://docs.controlplane.com/reference/secret#dictionary).",
				Optional:     true,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dictionary_as_envs": {
				Type:        schema.TypeMap,
				Description: "If a dictionary secret is defined, this output will be a key-value map in the following format: `key = cpln://secret/SECRET_NAME.key`.",
				Computed:    true,
			},
			"opaque": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#opaque).",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payload": {
							Type:        schema.TypeString,
							Description: "Plain text or base64 encoded string. Use `encoding` attribute to specify encoding.",
							Required:    true,
							Sensitive:   true,
						},
						"encoding": {
							Type:         schema.TypeString,
							Description:  "Available encodings: `plain`, `base64`. Default: `plain`.",
							Optional:     true,
							Default:      "plain",
							ValidateFunc: EncodingValidator,
						},
					},
				},
			},
			"tls": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#tls).",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Description: "Private Certificate.",
							Required:    true,
							Sensitive:   true,
						},
						"cert": {
							Type:        schema.TypeString,
							Description: "Public Certificate.",
							Required:    true,
						},
						"chain": {
							Type:        schema.TypeString,
							Description: "Chain Certificate.",
							Optional:    true,
							Default:     "",
						},
					},
				},
			},
			"gcp": {
				Type:             schema.TypeString,
				Description:      "JSON string containing the GCP secret. [Reference Page](https://docs.controlplane.com/reference/secret#gcp)",
				Optional:         true,
				Sensitive:        true,
				ExactlyOneOf:     []string{"aws", "azure_connector", "azure_sdk", "docker", "dictionary", "ecr", "gcp", "keypair", "nats_account", "opaque", "tls", "userpass"},
				DiffSuppressFunc: diffSuppressJSON,
			},
			"aws": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#aws).",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_key": {
							Type:        schema.TypeString,
							Description: "Secret Key provided by AWS.",
							Required:    true,
							Sensitive:   true,
						},
						"access_key": {
							Type:         schema.TypeString,
							Description:  "Access Key provided by AWS.",
							Required:     true,
							Sensitive:    true,
							ValidateFunc: AwsAccessKeyValidator,
						},
						"role_arn": {
							Type:         schema.TypeString,
							Description:  "Role ARN provided by AWS.",
							Optional:     true,
							Default:      "",
							ValidateFunc: AwsRoleArnValidator,
						},
						"external_id": {
							Type:        schema.TypeString,
							Description: "AWS IAM Role External ID.",
							Optional:    true,
						},
					},
				},
			},
			"ecr": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#ecr)",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_key": {
							Type:        schema.TypeString,
							Description: "Secret Key provided by AWS.",
							Required:    true,
							Sensitive:   true,
						},
						"access_key": {
							Type:         schema.TypeString,
							Description:  "Access Key provided by AWS.",
							Required:     true,
							ValidateFunc: AwsAccessKeyValidator,
						},
						"role_arn": {
							Type:         schema.TypeString,
							Description:  "Role ARN provided by AWS.",
							Optional:     true,
							Default:      "",
							ValidateFunc: AwsRoleArnValidator,
						},
						"external_id": {
							Type:        schema.TypeString,
							Description: "AWS IAM Role External ID. Used when setting up cross-account access to your ECR repositories.",
							Optional:    true,
						},
						"repos": {
							Type:        schema.TypeSet,
							Description: "List of ECR repositories.",
							Required:    true,
							Elem: &schema.Schema{
								MinItems: 1,
								MaxItems: 20,
								Type:     schema.TypeString,
							},
						},
					},
				},
			},
			"docker": {
				Type:             schema.TypeString,
				Description:      "JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#docker).",
				Optional:         true,
				Sensitive:        true,
				ExactlyOneOf:     []string{"aws", "azure_connector", "azure_sdk", "docker", "dictionary", "ecr", "gcp", "keypair", "nats_account", "opaque", "tls", "userpass"},
				DiffSuppressFunc: diffSuppressJSON,
			},
			"userpass": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#username).",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:         schema.TypeString,
							Description:  "Username.",
							Required:     true,
							ValidateFunc: EmptyValidator,
						},
						"password": {
							Type:         schema.TypeString,
							Description:  "Password.",
							Required:     true,
							Sensitive:    true,
							ValidateFunc: EmptyValidator,
						},
						"encoding": {
							Type:         schema.TypeString,
							Description:  "Available encodings: `plain`, `base64`. Default: `plain`.",
							Optional:     true,
							Default:      "plain",
							ValidateFunc: EncodingValidator,
						},
					},
				},
			},
			"keypair": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#keypair).",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_key": {
							Type:        schema.TypeString,
							Description: "Secret/Private Key.",
							Required:    true,
							Sensitive:   true,
						},
						"public_key": {
							Type:        schema.TypeString,
							Description: "Public Key.",
							Optional:    true,
							Default:     "",
						},
						"passphrase": {
							Type:        schema.TypeString,
							Description: "Passphrase for private key.",
							Optional:    true,
							Sensitive:   true,
							Default:     "",
						},
					},
				},
			},
			"azure_sdk": {
				Type:             schema.TypeString,
				Description:      "JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#azure).",
				Optional:         true,
				Sensitive:        true,
				ExactlyOneOf:     []string{"aws", "azure_connector", "azure_sdk", "docker", "dictionary", "ecr", "gcp", "keypair", "nats_account", "opaque", "tls", "userpass"},
				DiffSuppressFunc: diffSuppressJSON,
			},
			"azure_connector": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#azure-connector).",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Description: "Deployment URL.",
							Required:    true,
							Sensitive:   true,
						},
						"code": {
							Type:        schema.TypeString,
							Description: "Code/Key to authenticate to deployment URL.",
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"nats_account": {
				Type:         schema.TypeList,
				Description:  "[Reference Page](https://docs.controlplane.com/reference/secret#nats-account).",
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: secretDataObjectsNames,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Description: "Account ID.",
							Required:    true,
						},
						"private_key": {
							Type:        schema.TypeString,
							Description: "Private Key.",
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func diffSuppressJSON(k, old, new string, d *schema.ResourceData) bool {

	if old != "" && new != "" {

		bo, _ := json.Marshal(json.RawMessage(old))
		bn, _ := json.Marshal(json.RawMessage(new))

		return bytes.Equal(bo, bn)
	}

	return old == new
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	secret := client.Secret{}
	secret.Name = GetString(d.Get("name"))
	secret.Description = DescriptionHelper(*secret.Name, d.Get("description").(string))
	secret.Tags = GetStringMap(d.Get("tags"))

	if secret.Type = getSecretType(d); secret.Type == nil {
		return diag.FromErr(fmt.Errorf("unable to extract secret type"))
	}

	sType := *secret.Type
	data := d.Get(sType)

	if err := buildData(*secret.Type, data, &secret, false); err != nil {
		return diag.FromErr(err)
	}

	if *secret.Type == "azure_sdk" {
		*secret.Type = "azure-sdk"
	}

	if *secret.Type == "azure_connector" {
		*secret.Type = "azure-connector"
	}

	if *secret.Type == "nats_account" {
		*secret.Type = "nats-account"
	}

	c := m.(*client.Client)
	newSecret, code, err := c.CreateSecret(secret)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setSecret(d, newSecret)
}

func getSecretType(d *schema.ResourceData) *string {

	for _, s := range secretDataObjectsNames {
		if _, v := d.GetOk(s); v {
			return &s
		}
	}

	return nil
}

func buildData(secretType string, data interface{}, secret *client.Secret, update bool) error {

	var dataToSet *interface{}
	secret.Type = &secretType

	dataType := "unknown"
	dType := data

	switch data.(type) {
	case string:
		dataType = "string"
	case []interface{}:
		dataType = "interface"
	case map[string]interface{}:
		{
			dataType = "interface"

			dType = []interface{}{
				data,
			}
		}
	}

	if dataType == "string" && (secretType == "gcp" || secretType == "docker" || secretType == "azure_sdk") {

		dataString := dType.(string)
		dataToSet = GetInterface(dataString)

	} else if dataType == "interface" {

		dataArray := dType.([]interface{})

		if len(dataArray) == 1 {

			if secretType == "aws" || secretType == "ecr" || secretType == "keypair" || secretType == "tls" || secretType == "nats_account" {

				secretData := dataArray[0].(map[string]interface{})
				dataMap := make(map[string]interface{})

				if secretType == "aws" || secretType == "ecr" {

					dataMap["secretKey"] = secretData["secret_key"]
					dataMap["accessKey"] = secretData["access_key"]

					if secretData["role_arn"] != nil && secretData["role_arn"] != "" {
						dataMap["roleArn"] = secretData["role_arn"]
					} else {
						if update {
							dataMap["roleArn"] = nil
						}
					}

					if secretData["external_id"] != nil && secretData["external_id"] != "" {
						dataMap["externalId"] = secretData["external_id"]
					} else if update {
						dataMap["externalId"] = nil
					}
				}

				if secretType == "ecr" {

					repos := []string{}

					for _, value := range secretData["repos"].(*schema.Set).List() {
						repos = append(repos, value.(string))
					}

					if len(repos) > 0 {
						dataMap["repos"] = repos
					}
				}

				if secretType == "keypair" {

					dataMap["secretKey"] = secretData["secret_key"]

					if secretData["public_key"] != nil && secretData["public_key"] != "" {
						dataMap["publicKey"] = secretData["public_key"]
					} else {
						if update {
							dataMap["publicKey"] = nil
						}
					}

					if secretData["passphrase"] != nil && secretData["passphrase"] != "" {
						dataMap["passphrase"] = secretData["passphrase"]
					} else {
						if update {
							dataMap["passphrase"] = nil
						}
					}
				}

				if secretType == "tls" {

					dataMap["key"] = secretData["key"]

					if secretData["cert"] != nil && secretData["cert"] != "" {
						dataMap["cert"] = secretData["cert"]
					} else {
						if update {
							dataMap["cert"] = nil
						}
					}

					if secretData["chain"] != nil && secretData["chain"] != "" {
						dataMap["chain"] = secretData["chain"]
					} else {
						if update {
							dataMap["chain"] = nil
						}
					}
				}

				if secretType == "nats_account" {
					dataMap["accountId"] = secretData["account_id"]
					dataMap["privateKey"] = secretData["private_key"]
				}

				output := []interface{}{
					dataMap,
				}

				dataToSet = &output[0]

			} else {

				sData := make(map[string]interface{})

				for k, v := range dataArray[0].(map[string]interface{}) {
					if v != "" {
						sData[k] = v
					}
				}

				dataToSet = GetInterface(sData)
			}
		}
	}

	if dataToSet == nil {
		return fmt.Errorf("invalid secret input or data type. Secret type: %s", secretType)
	}

	if update {
		secret.DataReplace = dataToSet
	} else {
		secret.Data = dataToSet
	}

	return nil
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	secret, code, err := c.GetSecret(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setSecret(d, secret)
}

func resourceSecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "opaque", "tls", "gcp", "aws", "docker", "userpass", "keypair", "azure_sdk", "dictionary", "ecr", "azure_connector", "nats_account") {

		secretToUpdate := client.Secret{}
		secretToUpdate.Name = GetString(d.Get("name"))

		changedSecret := []string{}

		if d.HasChange("aws") {
			changedSecret = append(changedSecret, "aws")
		}

		if d.HasChange("azure_connector") {
			changedSecret = append(changedSecret, "azure_connector")
		}

		if d.HasChange("azure_sdk") {
			changedSecret = append(changedSecret, "azure_sdk")
		}

		if d.HasChange("docker") {
			changedSecret = append(changedSecret, "docker")
		}

		if d.HasChange("dictionary") {
			changedSecret = append(changedSecret, "dictionary")
		}

		if d.HasChange("ecr") {
			changedSecret = append(changedSecret, "ecr")
		}

		if d.HasChange("gcp") {
			changedSecret = append(changedSecret, "gcp")
		}

		if d.HasChange("keypair") {
			changedSecret = append(changedSecret, "keypair")
		}

		if d.HasChange("opaque") {
			changedSecret = append(changedSecret, "opaque")
		}

		if d.HasChange("tls") {
			changedSecret = append(changedSecret, "tls")
		}

		if d.HasChange("userpass") {
			changedSecret = append(changedSecret, "userpass")
		}

		if d.HasChange("nats_account") {
			changedSecret = append(changedSecret, "nats_account")
		}

		if d.HasChange("description") {
			secretToUpdate.Description = DescriptionHelper(*secretToUpdate.Name, d.Get("description").(string))
		}

		if d.HasChange("tags") {
			secretToUpdate.Tags = GetTagChanges(d)
		}

		if len(changedSecret) == 1 {

			s := changedSecret[0]

			data := d.Get(s)

			secretToUpdate.Type = GetString(s)

			if err := buildData(*secretToUpdate.Type, data, &secretToUpdate, true); err != nil {
				return diag.FromErr(err)
			}
		}

		c := m.(*client.Client)
		updatedSecret, _, err := c.UpdateSecret(secretToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setSecret(d, updatedSecret)
	}

	return nil
}

func setSecret(d *schema.ResourceData, secret *client.Secret) diag.Diagnostics {

	if secret == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*secret.Name)

	if err := SetBase(d, secret.Base); err != nil {
		return diag.FromErr(err)
	}

	if *secret.Type == "azure-sdk" {
		*secret.Type = "azure_sdk"
	}

	if *secret.Type == "azure-connector" {
		*secret.Type = "azure_connector"
	}

	if *secret.Type == "nats-account" {
		*secret.Type = "nats_account"
	}

	if err := d.Set("dictionary_as_envs", nil); err != nil {
		return diag.FromErr(err)
	}

	if secret.Data != nil {

		data := *secret.Data

		if *secret.Type == "gcp" || *secret.Type == "docker" || *secret.Type == "azure_sdk" {

			if err := d.Set(*secret.Type, data.(string)); err != nil {
				return diag.FromErr(err)
			}
		} else if *secret.Type == "dictionary" {

			secretData := data.(map[string]interface{})

			if err := d.Set(*secret.Type, secretData); err != nil {
				return diag.FromErr(err)
			}

			dict_as_envs := make(map[string]string)
			for key := range secretData {
				dict_as_envs[key] = fmt.Sprintf("cpln://secret/%s.%s", *secret.Name, key)
			}

			if err := d.Set("dictionary_as_envs", dict_as_envs); err != nil {
				return diag.FromErr(err)
			}
		} else {

			setData := make([]interface{}, 1)
			bData := make(map[string]interface{})

			if *secret.Type == "aws" || *secret.Type == "ecr" || *secret.Type == "keypair" || *secret.Type == "nats_account" {

				secretData := data.(map[string]interface{})

				if *secret.Type == "aws" || *secret.Type == "ecr" {

					bData["secret_key"] = secretData["secretKey"]
					bData["access_key"] = secretData["accessKey"]
					bData["role_arn"] = secretData["roleArn"]
					bData["external_id"] = secretData["externalId"]

					if *secret.Type == "ecr" {
						bData["repos"] = secretData["repos"]
					}
				}

				if *secret.Type == "keypair" {
					bData["secret_key"] = secretData["secretKey"]
					bData["public_key"] = secretData["publicKey"]
					bData["passphrase"] = secretData["passphrase"]
				}

				if *secret.Type == "nats_account" {
					bData["account_id"] = secretData["accountId"]
					bData["private_key"] = secretData["privateKey"]
				}

				setData[0] = bData

			} else {
				setData[0] = secret.Data
			}

			if err := d.Set(*secret.Type, setData); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err := SetSelfLink(secret.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	err := c.DeleteSecret(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
