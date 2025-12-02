package cpln

import (
	"context"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/catalog_template"
	whitespacestring "github.com/controlplane-com/terraform-provider-cpln/internal/provider/types/whitespacestring"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces at compile time
var (
	_ resource.Resource                = &CatalogTemplateResource{}
	_ resource.ResourceWithImportState = &CatalogTemplateResource{}
)

/*** Resource Model ***/

// CatalogTemplateResourceModel holds the Terraform state for the resource.
type CatalogTemplateResourceModel struct {
	ID        types.String                                     `tfsdk:"id"`
	Name      types.String                                     `tfsdk:"name"`
	Template  types.String                                     `tfsdk:"template"`
	Version   types.String                                     `tfsdk:"version"`
	Gvc       types.String                                     `tfsdk:"gvc"`
	Values    whitespacestring.WhitespaceNormalizedStringValue `tfsdk:"values"`
	Resources types.List                                       `tfsdk:"resources"`
}

// GetID returns the ID field from the catalog template resource model.
func (m CatalogTemplateResourceModel) GetID() types.String {
	// Return the stored ID value
	return m.ID
}

/*** Resource Configuration ***/

// CatalogTemplateResource is the resource implementation.
type CatalogTemplateResource struct {
	EntityBase
	Operations EntityOperations[CatalogTemplateResourceModel, client.MarketplaceRelease]
}

// NewCatalogTemplateResource returns a new instance of the resource implementation.
func NewCatalogTemplateResource() resource.Resource {
	return &CatalogTemplateResource{}
}

// Configure configures the resource before use.
func (r *CatalogTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	r.Operations = NewEntityOperations(r.client, &CatalogTemplateResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (r *CatalogTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (r *CatalogTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_catalog_template"
}

// Schema defines the schema for the resource.
func (r *CatalogTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Control Plane Catalog Template deployment. This resource allows you to install, update, and uninstall applications from the Control Plane marketplace catalog.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this catalog template deployment (same as name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The release name for this catalog release.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template": schema.StringAttribute{
				Description: "The name of the catalog template to deploy (e.g., 'postgres', 'nginx', 'redis').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Description: "The version of the catalog template to deploy.",
				Required:    true,
			},
			"gvc": schema.StringAttribute{
				Description: "The GVC where the template will be deployed. Leave empty if the template creates its own GVC (check template's createsGvc field).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"values": schema.StringAttribute{
				Description: "The values file content (YAML format) for customizing the template deployment.",
				Required:    true,
				CustomType:  whitespacestring.WhitespaceNormalizedStringType{},
			},
			"resources": schema.ListNestedAttribute{
				Description: "List of Control Plane resources created by this release. Populated from the helm release secret.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"kind": schema.StringAttribute{
							Description: "The kind of resource (e.g., 'workload', 'secret', 'gvc').",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the resource.",
							Computed:    true,
						},
						"link": schema.StringAttribute{
							Description: "The full Control Plane link to the resource.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Create creates the resource.
func (r *CatalogTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, r.Operations)
}

// Read fetches the current state of the resource.
func (r *CatalogTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, r.Operations)
}

// Update modifies the resource.
func (r *CatalogTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, r.Operations)
}

// Delete removes the resource.
func (r *CatalogTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, r.Operations)
}

/*** Resource Operator ***/

// CatalogTemplateResourceOperator is the operator for managing the state.
type CatalogTemplateResourceOperator struct {
	EntityOperator[CatalogTemplateResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (ctro *CatalogTemplateResourceOperator) NewAPIRequest(isUpdate bool) client.MarketplaceRelease {
	return client.MarketplaceRelease{
		Name:     ctro.Plan.Name.ValueString(),
		Template: ctro.Plan.Template.ValueString(),
		Version:  ctro.Plan.Version.ValueString(),
		Gvc:      ctro.Plan.Gvc.ValueStringPointer(),
		Values:   ctro.Plan.Values.ValueString(),
	}
}

// MapResponseToState creates a state model from response payload.
func (ctro *CatalogTemplateResourceOperator) MapResponseToState(release *client.MarketplaceRelease, isCreate bool) CatalogTemplateResourceModel {
	// Initialize state model with current plan values
	state := CatalogTemplateResourceModel{}

	// Set the ID to the release name
	state.ID = types.StringValue(release.Name)

	// Map basic release metadata to state
	state.Name = types.StringValue(release.Name)
	state.Template = types.StringValue(release.Template)
	state.Version = types.StringValue(release.Version)
	state.Gvc = types.StringPointerValue(release.Gvc)
	state.Values = whitespacestring.WhitespaceNormalizedStringValue{
		StringValue: types.StringValue(release.Values),
	}
	state.Resources = ctro.flattenResources(release.Resources)

	// Return the populated state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (ctro *CatalogTemplateResourceOperator) InvokeCreate(req client.MarketplaceRelease) (*client.MarketplaceRelease, int, error) {
	// Get template details to check if it creates its own GVC
	template, err := ctro.Client.GetMarketplaceTemplate(req.Template)
	if err != nil {
		// Return error if template fetch fails
		return nil, 0, fmt.Errorf("could not get template %s: %w", req.Template, err)
	}

	// Find the specified version to check createsGvc flag
	var createsGvc bool
	if template.Versions != nil {
		// Look up the version in the versions map
		if versionData, ok := (*template.Versions)[req.Version]; ok {
			// Extract the createsGvc flag
			if versionData.CreatesGvc != nil {
				createsGvc = *versionData.CreatesGvc
			}
		} else {
			// Version not found in template
			return nil, 0, fmt.Errorf("version %s not found in template %s", req.Version, req.Template)
		}
	}

	// Determine GVC value, send empty string when template creates GVC
	gvcValue := req.Gvc
	if createsGvc {
		// Template creates its own GVC, send empty string like console
		gvcValue = StringPointer("")
	} else if gvcValue == nil || *gvcValue == "" {
		// GVC is required but not provided
		return nil, 0, fmt.Errorf("the template %s version %s does not create its own GVC. You must specify a gvc", req.Template, req.Version)
	}

	// This is an actual install operation
	// Build the install request with action="install"
	installReq := client.MarketplaceInstallRequest{
		Org:      &ctro.Client.Org,
		Gvc:      gvcValue,
		Name:     &req.Name,
		Template: &req.Template,
		Version:  &req.Version,
		Values:   &req.Values,
		Action:   StringPointer("install"),
	}

	// Call the install endpoint to create the release
	_, err = ctro.Client.InstallMarketplaceRelease(installReq)
	if err != nil {
		return nil, 0, fmt.Errorf("could not install release %s: %s", req.Name, err.Error())
	}

	// Query the backend to get the release information including resources
	releaseInfo, code, err := ctro.queryRelease(req.Name)
	if err != nil {
		// Return error if we can't fetch release info
		return nil, code, fmt.Errorf("release was created but could not fetch details: %w", err)
	}

	// Return the release info
	return releaseInfo, code, nil
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (ctro *CatalogTemplateResourceOperator) InvokeRead(name string) (*client.MarketplaceRelease, int, error) {
	// Query the backend to get the release information
	releaseInfo, code, err := ctro.queryRelease(name)
	if err != nil {
		// Return error if query fails
		return nil, code, err
	}

	// If release doesn't exist, return 404
	if code == 404 {
		return nil, code, fmt.Errorf("release %s not found", name)
	}

	// Return the release info
	return releaseInfo, code, nil
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (ctro *CatalogTemplateResourceOperator) InvokeUpdate(req client.MarketplaceRelease) (*client.MarketplaceRelease, int, error) {
	// This is an actual upgrade operation
	// Build the upgrade request with action="upgrade"
	upgradeReq := client.MarketplaceInstallRequest{
		Org:      &ctro.Client.Org,
		Gvc:      req.Gvc,
		Name:     &req.Name,
		Template: &req.Template,
		Version:  &req.Version,
		Values:   &req.Values,
		Action:   StringPointer("upgrade"),
	}

	// Perform the release upgrade
	_, err := ctro.Client.InstallMarketplaceRelease(upgradeReq)
	if err != nil {
		return nil, 0, fmt.Errorf("could not upgrade release %s: %s", req.Name, err.Error())
	}

	// Query the backend to get updated release information
	releaseInfo, code, err := ctro.queryRelease(req.Name)
	if err != nil {
		// Return error if we can't fetch release info
		return nil, code, fmt.Errorf("release was upgraded but could not fetch details: %w", err)
	}

	// Return the release info
	return releaseInfo, code, nil
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (ctro *CatalogTemplateResourceOperator) InvokeDelete(name string) error {
	// Build the uninstall request
	uninstallReq := client.MarketplaceUninstallRequest{
		Org:  &ctro.Client.Org,
		Name: StringPointer(name),
	}

	// Uninstall the marketplace release
	_, err := ctro.Client.UninstallMarketplaceRelease(uninstallReq)
	if err != nil {
		return fmt.Errorf("could not uninstall release %s: %s", name, err.Error())
	}

	// Return success
	return nil
}

// Flatteners //

// flattenResources maps the resources array from API response to Terraform list type.
func (ctro *CatalogTemplateResourceOperator) flattenResources(input *[]client.HelmReleaseResource) types.List {
	// Get attribute types
	elementType := models.HelmReleaseResourceModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.HelmReleaseResourceModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.HelmReleaseResourceModel{
			Kind: types.StringValue(item.Kind),
			Name: types.StringPointerValue(item.Template.Name),
			Link: types.StringValue(item.Link),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(ctro.Ctx, ctro.Diags, blocks)
}

// Helpers //

// queryRelease queries the Control Plane API for helm release secrets.
func (ctro *CatalogTemplateResourceOperator) queryRelease(releaseName string) (*client.MarketplaceRelease, int, error) {
	// Build query to find helm release secrets
	query := client.Query{
		Kind: StringPointer("secret"),
		Spec: &client.QuerySpec{
			Match: StringPointer("all"),
			Terms: &[]client.QueryTerm{
				{
					Op:       StringPointer("~"),
					Property: StringPointer("name"),
					Value:    StringPointer("cpln-helm-release-"),
				},
				{
					Op:    StringPointer("="),
					Tag:   StringPointer("name"),
					Value: StringPointer(releaseName),
				},
				{
					Op:  StringPointer("exists"),
					Tag: StringPointer("cpln/marketplace"),
				},
				{
					Op:  StringPointer("exists"),
					Tag: StringPointer("cpln/marketplace-template"),
				},
				{
					Op:  StringPointer("exists"),
					Tag: StringPointer("cpln/marketplace-template-version"),
				},
			},
		},
	}

	// Query Control Plane API for helm release secrets
	releaseInfo, code, err := ctro.Client.GetMarketplaceRelease(releaseName, query)
	if err != nil {
		// Return error if query fails
		return nil, code, fmt.Errorf("could not query release %s: %w", releaseName, err)
	}

	// Return the release info
	return releaseInfo, code, nil
}
