package cpln

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/secret"
	modifiers "github.com/controlplane-com/terraform-provider-cpln/internal/provider/modifiers"
	validators "github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// secretAttributeNames conatins all the data attribute names defined by the secret resource.
var secretAttributeNames = []string{
	"aws", "azure_connector", "azure_sdk", "docker", "dictionary", "ecr",
	"gcp", "keypair", "opaque", "tls", "userpass", "nats_account",
}

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &SecretResource{}
	_ resource.ResourceWithImportState = &SecretResource{}
)

/*** Resource Model ***/

// SecretResourceModel holds the Terraform state for the resource.
type SecretResourceModel struct {
	EntityBaseModel
	Opaque           []models.OpaqueModel         `tfsdk:"opaque"`
	TLS              []models.TlsModel            `tfsdk:"tls"`
	GCP              types.String                 `tfsdk:"gcp"`
	AWS              []models.AwsModel            `tfsdk:"aws"`
	ECR              []models.EcrModel            `tfsdk:"ecr"`
	Docker           types.String                 `tfsdk:"docker"`
	UsernamePassword []models.UserpassModel       `tfsdk:"userpass"`
	KeyPair          []models.KeyPairModel        `tfsdk:"keypair"`
	Dictionary       types.Map                    `tfsdk:"dictionary"`
	DictionaryAsEnvs types.Map                    `tfsdk:"dictionary_as_envs"`
	AzureSdk         types.String                 `tfsdk:"azure_sdk"`
	AzureConnector   []models.AzureConnectorModel `tfsdk:"azure_connector"`
	NatsAccount      []models.NatsAccountModel    `tfsdk:"nats_account"`
	SecretLink       types.String                 `tfsdk:"secret_link"`
}

/*** Resource Configuration ***/

// SecretResource is the resource implementation.
type SecretResource struct {
	EntityBase
	Operations EntityOperations[SecretResourceModel, client.Secret]
}

// NewSecretResource returns a new instance of the resource implementation.
func NewSecretResource() resource.Resource {
	return &SecretResource{}
}

// Configure configures the resource before use.
func (sr *SecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	sr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	sr.Operations = NewEntityOperations(sr.client, NewSecretResourceOperator())
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (sr *SecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (sr *SecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_secret"
}

// Schema defines the schema for the resource.
func (sr *SecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(sr.EntityBaseAttributes("secret"), map[string]schema.Attribute{
			"gcp": schema.StringAttribute{
				MarkdownDescription: "JSON string containing the GCP secret. [Reference Page](https://docs.controlplane.com/reference/secret#gcp)",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					modifiers.SuppressDiffOnEqualJSON{},
				},
			},
			"docker": schema.StringAttribute{
				MarkdownDescription: "JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#docker).",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					modifiers.SuppressDiffOnEqualJSON{},
				},
			},
			"dictionary": schema.MapAttribute{
				MarkdownDescription: "List of unique key-value pairs. [Reference Page](https://docs.controlplane.com/reference/secret#dictionary).",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"dictionary_as_envs": schema.MapAttribute{
				MarkdownDescription: "If a dictionary secret is defined, this output will be a key-value map in the following format: `key = cpln://secret/SECRET_NAME.key`.",
				ElementType:         types.StringType,
				Computed:            true,
				PlanModifiers: []planmodifier.Map{
					modifiers.DictionaryAsEnvsPlanModifier{},
				},
			},
			"azure_sdk": schema.StringAttribute{
				MarkdownDescription: "JSON string containing the Docker secret. [Reference Page](https://docs.controlplane.com/reference/secret#azure).",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					modifiers.SuppressDiffOnEqualJSON{},
				},
			},
			"secret_link": schema.StringAttribute{
				Description: "Output used when linking a secret to an environment variable or volume.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"opaque": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#opaque).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"payload": schema.StringAttribute{
							Description: "Plain text or base64 encoded string. Use `encoding` attribute to specify encoding.",
							Required:    true,
							Sensitive:   true,
						},
						"encoding": schema.StringAttribute{
							Description: "Available encodings: `plain`, `base64`. Default: `plain`.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("plain"),
							Validators: []validator.String{
								stringvalidator.OneOf("plain", "base64"),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"tls": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#tls).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "Private Certificate.",
							Required:    true,
							Sensitive:   true,
						},
						"cert": schema.StringAttribute{
							Description: "Public Certificate.",
							Required:    true,
						},
						"chain": schema.StringAttribute{
							Description: "Chain Certificate.",
							Optional:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"aws": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#aws).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"access_key": schema.StringAttribute{
							Description: "Access Key provided by AWS.",
							Required:    true,
							Sensitive:   true,
							Validators: []validator.String{
								validators.NewPrefixStringValidator("AKIA", "AWS Access Key"),
							},
						},
						"secret_key": schema.StringAttribute{
							Description: "Secret Key provided by AWS.",
							Required:    true,
							Sensitive:   true,
						},
						"role_arn": schema.StringAttribute{
							Description: "Role ARN provided by AWS.",
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								validators.NewPrefixStringValidator("arn:", "AWS Role ARN"),
							},
						},
						"external_id": schema.StringAttribute{
							Description: "AWS IAM Role External ID.",
							Optional:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"ecr": schema.ListNestedBlock{
				MarkdownDescription: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"access_key": schema.StringAttribute{
							Description: "Access Key provided by AWS.",
							Required:    true,
							Sensitive:   true,
							Validators: []validator.String{
								validators.NewPrefixStringValidator("AKIA", "AWS Access Key"),
							},
						},
						"secret_key": schema.StringAttribute{
							Description: "Secret Key provided by AWS.",
							Required:    true,
							Sensitive:   true,
						},
						"role_arn": schema.StringAttribute{
							Description: "Role ARN provided by AWS.",
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								validators.NewPrefixStringValidator("arn:", "AWS Role ARN"),
							},
						},
						"external_id": schema.StringAttribute{
							Description: "AWS IAM Role External ID. Used when setting up cross-account access to your ECR repositories.",
							Optional:    true,
						},
						"repos": schema.SetAttribute{
							Description: "List of ECR repositories.",
							ElementType: types.StringType,
							Required:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.SizeAtMost(20),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"userpass": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#username).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							Description: "Username.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"password": schema.StringAttribute{
							Description: "Password.",
							Required:    true,
							Sensitive:   true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"encoding": schema.StringAttribute{
							Description: "Available encodings: `plain`, `base64`. Default: `plain`.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("plain"),
							Validators: []validator.String{
								stringvalidator.OneOf("plain", "base64"),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"keypair": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#keypair).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"secret_key": schema.StringAttribute{
							Description: "Secret/Private Key.",
							Required:    true,
							Sensitive:   true,
						},
						"public_key": schema.StringAttribute{
							Description: "Public Key.",
							Optional:    true,
						},
						"passphrase": schema.StringAttribute{
							Description: "Passphrase for private key.",
							Optional:    true,
							Sensitive:   true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"azure_connector": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#azure-connector).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: "Deployment URL.",
							Required:    true,
							Sensitive:   true,
						},
						"code": schema.StringAttribute{
							Description: "Code/Key to authenticate to deployment URL.",
							Required:    true,
							Sensitive:   true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"nats_account": schema.ListNestedBlock{
				MarkdownDescription: "[Reference Page](https://docs.controlplane.com/reference/secret#nats-account).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Description: "Account ID.",
							Required:    true,
						},
						"private_key": schema.StringAttribute{
							Description: "Private Key.",
							Required:    true,
							Sensitive:   true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

// ConfigValidators enforces mutual exclusivity between attributes.
func (mr *SecretResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	expressions := []path.Expression{}

	// Iterate over each secret attribute and append to the expressions slice
	for _, attributeName := range secretAttributeNames {
		expressions = append(expressions, path.MatchRoot(attributeName))
	}

	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(expressions...),
	}
}

// Create creates the resource.
func (sr *SecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, sr.Operations)
}

// Read fetches the current state of the resource.
func (sr *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, sr.Operations)
}

// Update modifies the resource.
func (sr *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, sr.Operations)
}

// Delete removes the resource.
func (sr *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, sr.Operations)
}

/*** Resource Operator ***/

// SecretResourceOperator is the operator for managing the state.
type SecretResourceOperator struct {
	EntityOperator[SecretResourceModel]
	DataBuilders map[string]func(SecretResourceModel) *interface{}
}

// NewSecretResourceOperator initializes the operator with per-type data builders.
func NewSecretResourceOperator() *SecretResourceOperator {
	// Instantiate operator instance
	op := &SecretResourceOperator{}

	// Assign data builders map with functions for each secret type
	op.DataBuilders = map[string]func(SecretResourceModel) *interface{}{
		// Handler for opaque secret type
		"opaque": func(s SecretResourceModel) *interface{} {
			// Build opaque secret data payload
			if m := op.buildOpaque(s.Opaque); m != nil {
				// Wrap opaque data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when Opaque data is not defined
			return nil
		},
		// Handler for tls secret type
		"tls": func(s SecretResourceModel) *interface{} {
			// Build TLS secret data payload
			if m := op.buildTls(s.TLS); m != nil {
				// Wrap TLS data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when TLS data is not defined
			return nil
		},
		// Handler for gcp secret type
		"gcp": func(s SecretResourceModel) *interface{} {
			// Build string representation of GCP field
			if data := GetInterface(BuildString(s.GCP)); data != nil {
				// Return GCP data when valid
				return data
			}

			// Return nil when GCP data is not defined
			return nil
		},
		// Handler for aws secret type
		"aws": func(s SecretResourceModel) *interface{} {
			// Build AWS secret data payload
			if m := op.buildAws(s.AWS); m != nil {
				// Wrap AWS data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when AWS data is not defined
			return nil
		},
		// Handler for ecr secret type
		"ecr": func(s SecretResourceModel) *interface{} {
			// Build ECR secret data payload
			if m := op.buildEcr(s.ECR); m != nil {
				// Wrap ECR data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when ECR data is not defined
			return nil
		},
		// Handler for docker secret type
		"docker": func(s SecretResourceModel) *interface{} {
			// Build string representation of Docker field
			if data := GetInterface(BuildString(s.Docker)); data != nil {
				// Return Docker data when valid
				return data
			}

			// Return nil when Docker data is not defined
			return nil
		},
		// Handler for userpass secret type
		"userpass": func(s SecretResourceModel) *interface{} {
			// Build username/password secret data payload
			if m := op.buildUserpass(s.UsernamePassword); m != nil {
				// Wrap username/password data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when username/password data is not defined
			return nil
		},
		// Handler for keypair secret type
		"keypair": func(s SecretResourceModel) *interface{} {
			// Build keypair secret data payload
			if m := op.buildKeypair(s.KeyPair); m != nil {
				// Wrap keypair data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when keypair data is not defined
			return nil
		},
		// Handler for dictionary secret type
		"dictionary": func(s SecretResourceModel) *interface{} {
			// Build map string from Dictionary field
			if m := op.BuildMapString(s.Dictionary); m != nil {
				// Wrap map string in interface pointer
				return GetInterface(*m)
			}

			// Return nil when Dictionary data is not defined
			return nil
		},
		// Handler for azure-sdk secret type
		"azure-sdk": func(s SecretResourceModel) *interface{} {
			// Build string representation of AzureSdk field
			if data := GetInterface(BuildString(s.AzureSdk)); data != nil {
				// Return AzureSdk data when valid
				return data
			}

			// Return nil when AzureSdk data is not defined
			return nil
		},
		// Handler for azure-connector secret type
		"azure-connector": func(s SecretResourceModel) *interface{} {
			// Build azure-connector secret data payload
			if m := op.buildAzureConnector(s.AzureConnector); m != nil {
				// Wrap azure-connector data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when azure-connector data is not defined
			return nil
		},
		// Handler for nats-account secret type
		"nats-account": func(s SecretResourceModel) *interface{} {
			// Build nats-account secret data payload
			if m := op.buildNatsAccount(s.NatsAccount); m != nil {
				// Wrap nats-account data in interface pointer
				return GetInterface(*m)
			}

			// Return nil when nats-account data is not defined
			return nil
		},
	}

	// Return initialized operator instance
	return op
}

// NewAPIRequest creates a request payload from a state model.
func (sro *SecretResourceOperator) NewAPIRequest(isUpdate bool) client.Secret {
	// Initialize a new request payload
	requestPayload := client.Secret{}

	// Populate Base fields from state
	sro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Get secret type
	secretType := sro.buildType(sro.Plan)

	// Set specific attributes
	requestPayload.Type = secretType

	// Build the secret data
	data := sro.buildData(sro.Plan, *secretType)

	// Determine whether we should replace data or not
	if isUpdate {
		requestPayload.DataReplace = data
	} else {
		requestPayload.Data = data
	}

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (sro *SecretResourceOperator) MapResponseToState(apiResp *client.Secret, isCreate bool) SecretResourceModel {
	// Initialize empty state model
	state := SecretResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Set specific attributes
	state.SecretLink = types.StringValue(fmt.Sprintf("cpln://secret/%s", *apiResp.Name))

	// Initialize the attributes within the state
	state.Opaque = nil
	state.TLS = nil
	state.GCP = types.StringNull()
	state.AWS = nil
	state.ECR = nil
	state.Docker = types.StringNull()
	state.UsernamePassword = nil
	state.KeyPair = nil
	state.Dictionary = types.MapNull(types.StringType)
	state.DictionaryAsEnvs = types.MapNull(types.StringType)
	state.AzureSdk = types.StringNull()
	state.AzureConnector = nil
	state.NatsAccount = nil

	// Handle secret data if present
	if apiResp.Data != nil {
		data := *apiResp.Data
		switch *apiResp.Type {
		case "opaque":
			// Flatten opaque secret payload
			state.Opaque = sro.flattenOpaque(data.(map[string]interface{}))

		case "tls":
			// Flatten TLS secret payload
			state.TLS = sro.flattenTls(data.(map[string]interface{}))

		case "gcp":
			// Flatten GCP secret payload
			state.GCP = PreserveJSONFormatting(data, sro.Plan.GCP)

		case "aws":
			// Flatten AWS secret payload
			state.AWS = sro.flattenAws(data.(map[string]interface{}))

		case "ecr":
			// Flatten ECR secret payload
			state.ECR = sro.flattenEcr(data.(map[string]interface{}))

		case "docker":
			// Flatten Docker secret payload
			state.Docker = PreserveJSONFormatting(data, sro.Plan.Docker)

		case "userpass":
			// Flatten username/password secret payload
			state.UsernamePassword = sro.flattenUserpass(data.(map[string]interface{}))

		case "keypair":
			// Flatten keypair secret payload
			state.KeyPair = sro.flattenKeyPair(data.(map[string]interface{}))

		case "dictionary":
			// Flatten map string for dictionary payload
			dataMap := data.(map[string]interface{})
			state.Dictionary = FlattenMapString(&dataMap)

			// Build environment variables for dictionary entries
			dictAsEnvs := make(map[string]interface{})
			for key := range dataMap {
				dictAsEnvs[key] = fmt.Sprintf("cpln://secret/%s.%s", *apiResp.Name, key)
			}
			state.DictionaryAsEnvs = FlattenMapString(&dictAsEnvs)

		case "azure-sdk":
			// Set Azure SDK secret value
			state.AzureSdk = PreserveJSONFormatting(data, sro.Plan.AzureSdk)

		case "azure-connector":
			// Flatten Azure connector secret payload
			state.AzureConnector = sro.flattenAzureConnector(data.(map[string]interface{}))

		case "nats-account":
			// Flatten NATS account secret payload
			state.NatsAccount = sro.flattenNatsAccount(data.(map[string]interface{}))
		}
	}

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (sro *SecretResourceOperator) InvokeCreate(req client.Secret) (*client.Secret, int, error) {
	return sro.Client.CreateSecret(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (sro *SecretResourceOperator) InvokeRead(name string) (*client.Secret, int, error) {
	return sro.Client.GetSecret(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (sro *SecretResourceOperator) InvokeUpdate(req client.Secret) (*client.Secret, int, error) {
	return sro.Client.UpdateSecret(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (sro *SecretResourceOperator) InvokeDelete(name string) error {
	return sro.Client.DeleteSecret(name)
}

// Builders //

// buildType determines the secret type based on which attribute is set in the state.
func (op *SecretResourceOperator) buildType(state SecretResourceModel) *string {
	// Use reflection to examine the fields of the state model
	v := reflect.ValueOf(state)
	t := v.Type()

	// Iterate over the list of known secret attribute names
	for _, attrName := range secretAttributeNames {
		// Iterate over the struct fields of the state model
		for i := range t.NumField() {
			field := t.Field(i)

			// Skip if the field's tfsdk tag doesn't match the current attribute name
			if field.Tag.Get("tfsdk") != attrName {
				continue
			}

			// Retrieve the field's value
			fv := v.Field(i)
			set := false

			// Check if the field is a slice and is non-empty
			if fv.Kind() == reflect.Slice {
				if fv.Len() > 0 {
					set = true
				}
			} else {
				// For other types that implement attr.Value, check if the value is defined
				if val, ok := fv.Interface().(attr.Value); ok && !val.IsNull() && !val.IsUnknown() {
					set = true
				}
			}

			// If the attribute is set, convert underscores to hyphens and return the type name
			if set {
				typeName := strings.ReplaceAll(attrName, "_", "-")
				return &typeName
			}

			// Exit the inner loop once the relevant field is found
			break
		}
	}

	// Return nil if no matching attribute is set
	return nil
}

// BuildData constructs the secret data payload or reports an error if invalid.
func (op *SecretResourceOperator) buildData(state SecretResourceModel, secretType string) *interface{} {
	// Check if a data builder is registered for the secret type
	if builder, ok := op.DataBuilders[secretType]; ok {
		// Invoke the builder function to assemble the payload
		if data := builder(state); data != nil {
			// Return the built data when successful
			return data
		}
	}

	// Report an error when no valid data could be produced
	op.Diags.AddError("Invalid Secret Data", fmt.Sprintf("invalid input for secret type %q", secretType))

	// Indicate failure by returning nil
	return nil
}

// buildOpaque constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildOpaque(state []models.OpaqueModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &map[string]interface{}{
		"payload":  BuildString(block.Payload),
		"encoding": BuildString(block.Encoding),
	}
}

// buildTls constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildTls(state []models.TlsModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct the output
	output := map[string]interface{}{
		"key":  BuildString(block.Key),
		"cert": BuildString(block.Cert),
	}

	// Set chain if specified
	if chain := BuildString(block.Chain); chain != nil {
		output["chain"] = chain
	}

	// Return the output
	return &output
}

// buildAws constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildAws(state []models.AwsModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct the output
	output := map[string]interface{}{
		"accessKey": BuildString(block.AccessKey),
		"secretKey": BuildString(block.SecretKey),
	}

	// Set roleArn if specified
	if roleArn := BuildString(block.RoleArn); roleArn != nil {
		output["roleArn"] = roleArn
	}

	// Set externalId if specified
	if externalId := BuildString(block.ExternalId); externalId != nil {
		output["externalId"] = externalId
	}

	// Return the output
	return &output
}

// buildEcr constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildEcr(state []models.EcrModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct the output
	output := map[string]interface{}{
		"accessKey": BuildString(block.AccessKey),
		"secretKey": BuildString(block.SecretKey),
		"repos":     sro.BuildSetString(block.Repos),
	}

	// Set roleArn if specified
	if roleArn := BuildString(block.RoleArn); roleArn != nil {
		output["roleArn"] = roleArn
	}

	// Set externalId if specified
	if externalId := BuildString(block.ExternalId); externalId != nil {
		output["externalId"] = externalId
	}

	// Return the output
	return &output
}

// buildUserpass constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildUserpass(state []models.UserpassModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &map[string]interface{}{
		"username": BuildString(block.Username),
		"password": BuildString(block.Password),
		"encoding": BuildString(block.Encoding),
	}
}

// buildKeypair constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildKeypair(state []models.KeyPairModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct the output
	output := map[string]interface{}{
		"secretKey": BuildString(block.SecretKey),
	}

	// Set the publicKey if specified
	if publicKey := BuildString(block.PublicKey); publicKey != nil {
		output["publicKey"] = publicKey
	}

	// Set the passphrase if specified
	if passphrase := BuildString(block.Passphrase); passphrase != nil {
		output["passphrase"] = passphrase
	}

	// Return the output
	return &output
}

// buildAzureConnector constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildAzureConnector(state []models.AzureConnectorModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &map[string]interface{}{
		"url":  BuildString(block.Url),
		"code": BuildString(block.Code),
	}
}

// buildNatsAccount constructs a map from the given Terraform state.
func (sro *SecretResourceOperator) buildNatsAccount(state []models.NatsAccountModel) *map[string]interface{} {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &map[string]interface{}{
		"accountId":  BuildString(block.AccountId),
		"privateKey": BuildString(block.PrivateKey),
	}
}

// Flatteners //

// flattenOpaque transforms *interface{} into a []models.OpaqueModel.
func (sro *SecretResourceOperator) flattenOpaque(input map[string]interface{}) []models.OpaqueModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.OpaqueModel{
		Payload:  types.StringValue(input["payload"].(string)),
		Encoding: types.StringValue(input["encoding"].(string)),
	}

	// Return a slice containing the single block
	return []models.OpaqueModel{block}
}

// flattenTls transforms *interface{} into a []models.TlsModel.
func (sro *SecretResourceOperator) flattenTls(input map[string]interface{}) []models.TlsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.TlsModel{
		Key:  types.StringValue(input["key"].(string)),
		Cert: types.StringValue(input["cert"].(string)),
	}

	// Set chain if specified
	if chain, ok := input["chain"]; ok {
		block.Chain = types.StringValue(chain.(string))
	}

	// Return a slice containing the single block
	return []models.TlsModel{block}
}

// flattenAws transforms *interface{} into a []models.AwsModel.
func (sro *SecretResourceOperator) flattenAws(input map[string]interface{}) []models.AwsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AwsModel{
		AccessKey: types.StringValue(input["accessKey"].(string)),
		SecretKey: types.StringValue(input["secretKey"].(string)),
	}

	// Set roleArn if specified
	if roleArn, ok := input["roleArn"]; ok {
		block.RoleArn = types.StringValue(roleArn.(string))
	}

	// Set externalId if specified
	if externalId, ok := input["externalId"]; ok {
		block.ExternalId = types.StringValue(externalId.(string))
	}

	// Return a slice containing the single block
	return []models.AwsModel{block}
}

// flattenEcr transforms *interface{} into a []models.EcrModel.
func (sro *SecretResourceOperator) flattenEcr(input map[string]interface{}) []models.EcrModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	repos := ToStringSlice(input["repos"].([]interface{}))

	// Build a single block
	block := models.EcrModel{
		AccessKey: types.StringValue(input["accessKey"].(string)),
		SecretKey: types.StringValue(input["secretKey"].(string)),
		Repos:     FlattenSetString(&repos),
	}

	// Set roleArn if specified
	if roleArn, ok := input["roleArn"]; ok {
		block.RoleArn = types.StringValue(roleArn.(string))
	}

	// Set externalId if specified
	if externalId, ok := input["externalId"]; ok {
		block.ExternalId = types.StringValue(externalId.(string))
	}

	// Return a slice containing the single block
	return []models.EcrModel{block}
}

// flattenUserpass transforms *interface{} into a []models.UserpassModel.
func (sro *SecretResourceOperator) flattenUserpass(input map[string]interface{}) []models.UserpassModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.UserpassModel{
		Username: types.StringValue(input["username"].(string)),
		Password: types.StringValue(input["password"].(string)),
		Encoding: types.StringValue(input["encoding"].(string)),
	}

	// Return a slice containing the single block
	return []models.UserpassModel{block}
}

// flattenKeyPair transforms *interface{} into a []models.KeyPairModel.
func (sro *SecretResourceOperator) flattenKeyPair(input map[string]interface{}) []models.KeyPairModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.KeyPairModel{
		SecretKey: types.StringValue(input["secretKey"].(string)),
	}

	// Set publicKey if specified
	if publicKey, ok := input["publicKey"]; ok {
		block.PublicKey = types.StringValue(publicKey.(string))
	}

	// Set passphrase if specified
	if passphrase, ok := input["passphrase"]; ok {
		block.Passphrase = types.StringValue(passphrase.(string))
	}

	// Return a slice containing the single block
	return []models.KeyPairModel{block}
}

// flattenAzureConnector transforms *interface{} into a []models.AzureConnectorModel.
func (sro *SecretResourceOperator) flattenAzureConnector(input map[string]interface{}) []models.AzureConnectorModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AzureConnectorModel{
		Url:  types.StringValue(input["url"].(string)),
		Code: types.StringValue(input["code"].(string)),
	}

	// Return a slice containing the single block
	return []models.AzureConnectorModel{block}
}

// flattenNatsAccount transforms *interface{} into a []models.NatsAccountModel.
func (sro *SecretResourceOperator) flattenNatsAccount(input map[string]interface{}) []models.NatsAccountModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.NatsAccountModel{
		AccountId:  types.StringValue(input["accountId"].(string)),
		PrivateKey: types.StringValue(input["privateKey"].(string)),
	}

	// Return a slice containing the single block
	return []models.NatsAccountModel{block}
}
