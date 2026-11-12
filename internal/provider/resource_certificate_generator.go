package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/certificate_generator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &CertificateGeneratorResource{}
	_ resource.ResourceWithImportState = &CertificateGeneratorResource{}
)

/*** Resource Model ***/

// CertificateGeneratorResourceModel holds the Terraform state for the resource.
type CertificateGeneratorResourceModel struct {
	EntityBaseModel
	Acme   types.Object `tfsdk:"acme"`
	Status types.List   `tfsdk:"status"`
}

/*** Resource Configuration ***/

// CertificateGeneratorResource is the resource implementation.
type CertificateGeneratorResource struct {
	EntityBase
	Operations EntityOperations[CertificateGeneratorResourceModel, client.CertificateGenerator]
}

// NewCertificateGeneratorResource returns a new instance of the resource implementation.
func NewCertificateGeneratorResource() resource.Resource {
	return &CertificateGeneratorResource{}
}

// Configure configures the resource before use.
func (cgr *CertificateGeneratorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	cgr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	cgr.Operations = NewEntityOperations(cgr.client, &CertificateGeneratorResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (cgr *CertificateGeneratorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (cgr *CertificateGeneratorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_certificate_generator"
}

// Schema defines the schema for the resource.
func (cgr *CertificateGeneratorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(cgr.EntityBaseAttributes("Certificate Generator"), map[string]schema.Attribute{
			"acme": schema.SingleNestedAttribute{
				Description: "Issues domain certificates from an ACME certificate authority.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"server": schema.StringAttribute{
						Description: "ACME directory URL of the CA, e.g. `https://acme-v02.api.letsencrypt.org/directory`. Must use `https`. Changing it recreates the certificate generator.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"email": schema.StringAttribute{
						Description: "Contact email registered with the ACME account.",
						Optional:    true,
					},
					"external_account_binding": schema.SingleNestedAttribute{
						Description: "Binds the ACME account to an existing account at the CA. Required by most commercial CAs.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"key_id": schema.StringAttribute{
								Description: "Key identifier issued by the CA.",
								Required:    true,
							},
							"secret_link": schema.StringAttribute{
								Description: "Link to a secret of type `opaque` whose payload is the base64url encoded HMAC key exactly as issued by the CA (encoding `plain`).",
								Required:    true,
							},
						},
					},
					"challenge": schema.SingleNestedAttribute{
						Description: "How ACME challenges are solved. Exactly one of `http01` or `dns01` must be set. Changing the challenge type recreates the certificate generator.",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"http01": schema.SingleNestedAttribute{
								Description: "Solve ACME challenges over HTTP on the domain itself. Only domain-name certificates can be issued; wildcards require `dns01`.",
								Optional:    true,
								Attributes:  map[string]schema.Attribute{},
								Validators: []validator.Object{
									objectvalidator.ExactlyOneOf(
										path.MatchRelative().AtParent().AtName("dns01"),
									),
								},
								PlanModifiers: []planmodifier.Object{
									cgr.RequiresReplaceOnChallengeTypeChange(),
								},
							},
							"dns01": schema.SingleNestedAttribute{
								Description: "Solve ACME challenges with TXT records, either delegated to Control Plane (`cpln`) or written into DNS zones the org controls. Required for wildcard certificates.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"provider": schema.StringAttribute{
										Description: "DNS provider that hosts the zones of the domains using this generator. `cpln`: domains CNAME `_acme-challenge` to Control Plane, no credentials needed. Valid values: `cpln`, `route53`, `cloudDns`, `azureDns`, `cloudflare`. Default: `cpln`.",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("cpln"),
										Validators: []validator.String{
											stringvalidator.OneOf("cpln", "route53", "cloudDns", "azureDns", "cloudflare"),
										},
									},
									"route53": schema.SingleNestedAttribute{
										Description: "Route53 settings. Required when `provider` is `route53`.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"region": schema.StringAttribute{
												Description: "AWS region of the Route53 API to use, e.g. `us-east-1`.",
												Required:    true,
											},
											"hosted_zone_id": schema.StringAttribute{
												Description: "Restrict record updates to this hosted zone; discovered from the domain when omitted.",
												Optional:    true,
											},
											"secret_link": schema.StringAttribute{
												Description: "Link to a secret of type `aws`. The access key must be allowed to change TXT records in the hosted zone.",
												Required:    true,
											},
										},
									},
									"cloud_dns": schema.SingleNestedAttribute{
										Description: "Google Cloud DNS settings. Required when `provider` is `cloudDns`.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"project": schema.StringAttribute{
												Description: "GCP project that owns the Cloud DNS zone.",
												Required:    true,
											},
											"secret_link": schema.StringAttribute{
												Description: "Link to a secret of type `gcp` holding a service account key with the DNS Administrator role.",
												Required:    true,
											},
										},
									},
									"azure_dns": schema.SingleNestedAttribute{
										Description: "Azure DNS settings. Required when `provider` is `azureDns`.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"resource_group_name": schema.StringAttribute{
												Description: "Resource group that contains the DNS zone.",
												Required:    true,
											},
											"hosted_zone_name": schema.StringAttribute{
												Description: "Restrict record updates to this DNS zone; discovered from the domain when omitted.",
												Optional:    true,
											},
											"secret_link": schema.StringAttribute{
												Description: "Link to a secret of type `azure-sdk`. Its service principal must be a DNS Zone Contributor on the zone.",
												Required:    true,
											},
										},
									},
									"cloudflare": schema.SingleNestedAttribute{
										Description: "Cloudflare settings. Required when `provider` is `cloudflare`.",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"secret_link": schema.StringAttribute{
												Description: "Link to a secret of type `opaque` whose payload is a Cloudflare API token with Zone:DNS:Edit permission.",
												Required:    true,
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									cgr.RequiresReplaceOnChallengeTypeChange(),
								},
							},
						},
					},
				},
			},
			"status": schema.ListNestedAttribute{
				Description: "The status of the Certificate Generator.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ready": schema.BoolAttribute{
							Description: "True once the generator is registered with the CA and can issue certificates.",
							Computed:    true,
						},
						"message": schema.StringAttribute{
							Description: "Reason the generator is not ready.",
							Computed:    true,
						},
						"last_updated": schema.StringAttribute{
							Description: "Timestamp of the last status update.",
							Computed:    true,
						},
					},
				},
			},
		}),
	}
}

// Create creates the resource.
func (cgr *CertificateGeneratorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, cgr.Operations)
}

// Read fetches the current state of the resource.
func (cgr *CertificateGeneratorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, cgr.Operations)
}

// Update modifies the resource.
func (cgr *CertificateGeneratorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, cgr.Operations)
}

// Delete removes the resource.
func (cgr *CertificateGeneratorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, cgr.Operations)
}

/*** Plan Modifiers ***/

// RequiresReplaceOnChallengeTypeChange returns a plan modifier that forces a replace when a challenge type is added or removed, since the API fixes it at creation.
func (cgr *CertificateGeneratorResource) RequiresReplaceOnChallengeTypeChange() planmodifier.Object {
	return objectplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
			// Skip while the planned value is still unknown
			if req.PlanValue.IsUnknown() {
				return
			}

			// Replace when the challenge type is set in only one of the prior state and the plan
			resp.RequiresReplace = req.StateValue.IsNull() != req.PlanValue.IsNull()
		},
		"Recreate when the challenge type changes",
		"The API does not allow switching between http01 and dns01 after creation, so the certificate generator is recreated.",
	)
}

/*** Resource Operator ***/

// CertificateGeneratorResourceOperator is the operator for managing the state.
type CertificateGeneratorResourceOperator struct {
	EntityOperator[CertificateGeneratorResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (cgro *CertificateGeneratorResourceOperator) NewAPIRequest(isUpdate bool) client.CertificateGenerator {
	// Initialize a new request payload
	requestPayload := client.CertificateGenerator{}

	// Initialize the spec struct
	var spec *client.CertificateGeneratorSpec = &client.CertificateGeneratorSpec{}

	// Populate Base fields from state
	cgro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Map planned state attributes to the API struct
	if isUpdate {
		requestPayload.SpecReplace = spec
	} else {
		requestPayload.Spec = spec
	}

	// Set specific attributes
	spec.Acme = cgro.buildAcme(cgro.Plan.Acme)

	// Derive the provider from the acme attribute, the only provider the API supports
	if spec.Acme != nil {
		spec.Provider = StringPointer("acme")
	}

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (cgro *CertificateGeneratorResourceOperator) MapResponseToState(apiResp *client.CertificateGenerator, isCreate bool) CertificateGeneratorResourceModel {
	// Initialize empty state model
	state := CertificateGeneratorResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Just in case the spec is nil
	if apiResp.Spec != nil {
		state.Acme = cgro.flattenAcme(apiResp.Spec.Acme, cgro.Plan.Acme)
	} else {
		state.Acme = types.ObjectNull(models.AcmeModel{}.AttributeTypes().(types.ObjectType).AttrTypes)
	}

	// Set specific attributes
	state.Status = cgro.flattenStatus(apiResp.Status)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (cgro *CertificateGeneratorResourceOperator) InvokeCreate(req client.CertificateGenerator) (*client.CertificateGenerator, int, error) {
	return cgro.Client.CreateCertificateGenerator(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (cgro *CertificateGeneratorResourceOperator) InvokeRead(name string) (*client.CertificateGenerator, int, error) {
	return cgro.Client.GetCertificateGenerator(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (cgro *CertificateGeneratorResourceOperator) InvokeUpdate(req client.CertificateGenerator) (*client.CertificateGenerator, int, error) {
	return cgro.Client.UpdateCertificateGenerator(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (cgro *CertificateGeneratorResourceOperator) InvokeDelete(name string) error {
	return cgro.Client.DeleteCertificateGenerator(name)
}

// Builders //

// buildAcme constructs a CertificateGeneratorAcme from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcme(state types.Object) *client.CertificateGeneratorAcme {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.AcmeModel](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorAcme{
		Server:                 BuildString(block.Server),
		Email:                  BuildString(block.Email),
		ExternalAccountBinding: cgro.buildAcmeExternalAccountBinding(block.ExternalAccountBinding),
		Challenge:              cgro.buildAcmeChallenge(block.Challenge),
	}
}

// buildAcmeExternalAccountBinding constructs a CertificateGeneratorExternalAccountBinding from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeExternalAccountBinding(state types.Object) *client.CertificateGeneratorExternalAccountBinding {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.ExternalAccountBindingModel](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorExternalAccountBinding{
		KeyId:      BuildString(block.KeyId),
		SecretLink: BuildString(block.SecretLink),
	}
}

// buildAcmeChallenge constructs a CertificateGeneratorAcmeChallenge from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeChallenge(state types.Object) *client.CertificateGeneratorAcmeChallenge {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.ChallengeModel](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorAcmeChallenge{
		Http01: cgro.buildAcmeChallengeHttp01(block.Http01),
		Dns01:  cgro.buildAcmeChallengeDns01(block.Dns01),
	}
}

// buildAcmeChallengeHttp01 constructs a CertificateGeneratorHttp01Challenge from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeChallengeHttp01(state types.Object) *client.CertificateGeneratorHttp01Challenge {
	// Return nil when the challenge is not configured
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Return an empty challenge, which the API expects as {}
	return &client.CertificateGeneratorHttp01Challenge{}
}

// buildAcmeChallengeDns01 constructs a CertificateGeneratorDns01Challenge from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeChallengeDns01(state types.Object) *client.CertificateGeneratorDns01Challenge {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.Dns01Model](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorDns01Challenge{
		Provider:   BuildString(block.Provider),
		Route53:    cgro.buildAcmeChallengeDns01Route53(block.Route53),
		CloudDns:   cgro.buildAcmeChallengeDns01CloudDns(block.CloudDns),
		AzureDns:   cgro.buildAcmeChallengeDns01AzureDns(block.AzureDns),
		Cloudflare: cgro.buildAcmeChallengeDns01Cloudflare(block.Cloudflare),
	}
}

// buildAcmeChallengeDns01Route53 constructs a CertificateGeneratorRoute53Dns01 from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeChallengeDns01Route53(state types.Object) *client.CertificateGeneratorRoute53Dns01 {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.Route53Model](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorRoute53Dns01{
		Region:       BuildString(block.Region),
		HostedZoneId: BuildString(block.HostedZoneId),
		SecretLink:   BuildString(block.SecretLink),
	}
}

// buildAcmeChallengeDns01CloudDns constructs a CertificateGeneratorCloudDnsDns01 from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeChallengeDns01CloudDns(state types.Object) *client.CertificateGeneratorCloudDnsDns01 {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.CloudDnsModel](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorCloudDnsDns01{
		Project:    BuildString(block.Project),
		SecretLink: BuildString(block.SecretLink),
	}
}

// buildAcmeChallengeDns01AzureDns constructs a CertificateGeneratorAzureDnsDns01 from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeChallengeDns01AzureDns(state types.Object) *client.CertificateGeneratorAzureDnsDns01 {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.AzureDnsModel](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorAzureDnsDns01{
		ResourceGroupName: BuildString(block.ResourceGroupName),
		HostedZoneName:    BuildString(block.HostedZoneName),
		SecretLink:        BuildString(block.SecretLink),
	}
}

// buildAcmeChallengeDns01Cloudflare constructs a CertificateGeneratorCloudflareDns01 from the given Terraform state.
func (cgro *CertificateGeneratorResourceOperator) buildAcmeChallengeDns01Cloudflare(state types.Object) *client.CertificateGeneratorCloudflareDns01 {
	// Convert Terraform object into model block using generic helper
	block, ok := BuildObject[models.CloudflareModel](cgro.Ctx, cgro.Diags, state)

	// Return nil if conversion failed or object was null
	if !ok {
		return nil
	}

	// Construct and return the output
	return &client.CertificateGeneratorCloudflareDns01{
		SecretLink: BuildString(block.SecretLink),
	}
}

// Flatteners //

// flattenAcme transforms *client.CertificateGeneratorAcme into a Terraform types.Object, keeping the link forms of the prior state.
func (cgro *CertificateGeneratorResourceOperator) flattenAcme(input *client.CertificateGeneratorAcme, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.AcmeModel{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior block so nested links keep the form the user wrote
	prior := models.AcmeModel{}
	if block, ok := BuildObject[models.AcmeModel](cgro.Ctx, cgro.Diags, priorState); ok {
		prior = *block
	}

	// Build a single block
	block := models.AcmeModel{
		Server:                 types.StringPointerValue(input.Server),
		Email:                  types.StringPointerValue(input.Email),
		ExternalAccountBinding: cgro.flattenAcmeExternalAccountBinding(input.ExternalAccountBinding, prior.ExternalAccountBinding),
		Challenge:              cgro.flattenAcmeChallenge(input.Challenge, prior.Challenge),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenAcmeExternalAccountBinding transforms *client.CertificateGeneratorExternalAccountBinding into a Terraform types.Object, keeping the prior secret link form.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeExternalAccountBinding(input *client.CertificateGeneratorExternalAccountBinding, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.ExternalAccountBindingModel{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior secret link so its short or full form is preserved
	var priorSecretLink types.String
	if prior, ok := BuildObject[models.ExternalAccountBindingModel](cgro.Ctx, cgro.Diags, priorState); ok {
		priorSecretLink = prior.SecretLink
	}

	// Build a single block
	block := models.ExternalAccountBindingModel{
		KeyId:      types.StringPointerValue(input.KeyId),
		SecretLink: cgro.FlattenLinkString(priorSecretLink, input.SecretLink, cgro.Client.Org),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenAcmeChallenge transforms *client.CertificateGeneratorAcmeChallenge into a Terraform types.Object, keeping the link forms of the prior state.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeChallenge(input *client.CertificateGeneratorAcmeChallenge, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.ChallengeModel{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior block so nested links keep the form the user wrote
	prior := models.ChallengeModel{}
	if block, ok := BuildObject[models.ChallengeModel](cgro.Ctx, cgro.Diags, priorState); ok {
		prior = *block
	}

	// Build a single block
	block := models.ChallengeModel{
		Http01: cgro.flattenAcmeChallengeHttp01(input.Http01),
		Dns01:  cgro.flattenAcmeChallengeDns01(input.Dns01, prior.Dns01),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenAcmeChallengeHttp01 transforms *client.CertificateGeneratorHttp01Challenge into a Terraform types.Object.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeChallengeHttp01(input *client.CertificateGeneratorHttp01Challenge) types.Object {
	// Get attribute types
	elementType := models.Http01Model{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &models.Http01Model{})
}

// flattenAcmeChallengeDns01 transforms *client.CertificateGeneratorDns01Challenge into a Terraform types.Object, keeping the link forms of the prior state.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeChallengeDns01(input *client.CertificateGeneratorDns01Challenge, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.Dns01Model{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior block so nested links keep the form the user wrote
	prior := models.Dns01Model{}
	if block, ok := BuildObject[models.Dns01Model](cgro.Ctx, cgro.Diags, priorState); ok {
		prior = *block
	}

	// Build a single block
	block := models.Dns01Model{
		Provider:   types.StringPointerValue(input.Provider),
		Route53:    cgro.flattenAcmeChallengeDns01Route53(input.Route53, prior.Route53),
		CloudDns:   cgro.flattenAcmeChallengeDns01CloudDns(input.CloudDns, prior.CloudDns),
		AzureDns:   cgro.flattenAcmeChallengeDns01AzureDns(input.AzureDns, prior.AzureDns),
		Cloudflare: cgro.flattenAcmeChallengeDns01Cloudflare(input.Cloudflare, prior.Cloudflare),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenAcmeChallengeDns01Route53 transforms *client.CertificateGeneratorRoute53Dns01 into a Terraform types.Object, keeping the prior secret link form.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeChallengeDns01Route53(input *client.CertificateGeneratorRoute53Dns01, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.Route53Model{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior secret link so its short or full form is preserved
	var priorSecretLink types.String
	if prior, ok := BuildObject[models.Route53Model](cgro.Ctx, cgro.Diags, priorState); ok {
		priorSecretLink = prior.SecretLink
	}

	// Build a single block
	block := models.Route53Model{
		Region:       types.StringPointerValue(input.Region),
		HostedZoneId: types.StringPointerValue(input.HostedZoneId),
		SecretLink:   cgro.FlattenLinkString(priorSecretLink, input.SecretLink, cgro.Client.Org),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenAcmeChallengeDns01CloudDns transforms *client.CertificateGeneratorCloudDnsDns01 into a Terraform types.Object, keeping the prior secret link form.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeChallengeDns01CloudDns(input *client.CertificateGeneratorCloudDnsDns01, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.CloudDnsModel{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior secret link so its short or full form is preserved
	var priorSecretLink types.String
	if prior, ok := BuildObject[models.CloudDnsModel](cgro.Ctx, cgro.Diags, priorState); ok {
		priorSecretLink = prior.SecretLink
	}

	// Build a single block
	block := models.CloudDnsModel{
		Project:    types.StringPointerValue(input.Project),
		SecretLink: cgro.FlattenLinkString(priorSecretLink, input.SecretLink, cgro.Client.Org),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenAcmeChallengeDns01AzureDns transforms *client.CertificateGeneratorAzureDnsDns01 into a Terraform types.Object, keeping the prior secret link form.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeChallengeDns01AzureDns(input *client.CertificateGeneratorAzureDnsDns01, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.AzureDnsModel{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior secret link so its short or full form is preserved
	var priorSecretLink types.String
	if prior, ok := BuildObject[models.AzureDnsModel](cgro.Ctx, cgro.Diags, priorState); ok {
		priorSecretLink = prior.SecretLink
	}

	// Build a single block
	block := models.AzureDnsModel{
		ResourceGroupName: types.StringPointerValue(input.ResourceGroupName),
		HostedZoneName:    types.StringPointerValue(input.HostedZoneName),
		SecretLink:        cgro.FlattenLinkString(priorSecretLink, input.SecretLink, cgro.Client.Org),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenAcmeChallengeDns01Cloudflare transforms *client.CertificateGeneratorCloudflareDns01 into a Terraform types.Object, keeping the prior secret link form.
func (cgro *CertificateGeneratorResourceOperator) flattenAcmeChallengeDns01Cloudflare(input *client.CertificateGeneratorCloudflareDns01, priorState types.Object) types.Object {
	// Get attribute types
	elementType := models.CloudflareModel{}.AttributeTypes().(types.ObjectType)

	// Check if the input is nil
	if input == nil {
		// Return a null object
		return types.ObjectNull(elementType.AttrTypes)
	}

	// Read the prior secret link so its short or full form is preserved
	var priorSecretLink types.String
	if prior, ok := BuildObject[models.CloudflareModel](cgro.Ctx, cgro.Diags, priorState); ok {
		priorSecretLink = prior.SecretLink
	}

	// Build a single block
	block := models.CloudflareModel{
		SecretLink: cgro.FlattenLinkString(priorSecretLink, input.SecretLink, cgro.Client.Org),
	}

	// Return the successfully created types.Object
	return FlattenObject(cgro.Ctx, cgro.Diags, &block)
}

// flattenStatus transforms *client.CertificateGeneratorStatus into a Terraform types.List.
func (cgro *CertificateGeneratorResourceOperator) flattenStatus(input *client.CertificateGeneratorStatus) types.List {
	// Get attribute types
	elementType := models.StatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusModel{
		Ready:       types.BoolPointerValue(input.Ready),
		Message:     types.StringPointerValue(input.Message),
		LastUpdated: types.StringPointerValue(input.LastUpdated),
	}

	// Return the successfully created types.List
	return FlattenList(cgro.Ctx, cgro.Diags, []models.StatusModel{block})
}
