package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/image"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &ImagesDataSource{}
	_ datasource.DataSourceWithConfigure = &ImagesDataSource{}
)

/*** Data Source Model ***/

// ImagesDataSourceModel holds the Terraform state for the data source.
type ImagesDataSourceModel struct {
	Images types.List `tfsdk:"images"`
	Query  types.List `tfsdk:"query"`
}

/*** Data Source Configuration ***/

// ImagesDataSource is the data source implementation.
type ImagesDataSource struct {
	EntityBase
}

// NewImagesDataSource returns a new instance of the data source implementation.
func NewImagesDataSource() datasource.DataSource {
	return &ImagesDataSource{}
}

// Metadata provides the data source type name.
func (d *ImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_images"
}

// Configure configures the data source before use.
func (d *ImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the data source.
func (d *ImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Description: "List of all images of the org.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cpln_id": schema.StringAttribute{
							Description: "The ID, in GUID format, of the image.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the image.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Key-value map of resource tags.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"self_link": schema.StringAttribute{
							Description: "Full link to this resource. Can be referenced by other resources.",
							Computed:    true,
						},
						"tag": schema.StringAttribute{
							Description: "Tag of the image.",
							Computed:    true,
						},
						"repository": schema.StringAttribute{
							Description: "Respository name of the image.",
							Computed:    true,
						},
						"digest": schema.StringAttribute{
							Description: "A unique SHA256 hash used to identify a specific image version within the image registry.",
							Computed:    true,
						},
						"manifest": schema.ListNestedAttribute{
							Description: "The manifest provides configuration and layers information about the image. It plays a crucial role in the Docker image distribution system, enabling image creation, verification, and replication in a consistent and secure manner.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"config": d.ConfigSchema("The config is a JSON blob that contains the image configuration data which includes environment variables, default command to run, and other settings necessary to run the container based on this image."),
									"layers": d.ConfigSchema("Layers lists the digests of the image's layers. These layers are filesystem changes or additions made in each step of the Docker image's creation process. The layers are stored separately and pulled as needed, which allows for efficient storage and transfer of images. Each layer is represented by a SHA256 digest, ensuring the integrity and authenticity of the image."),
									"media_type": schema.StringAttribute{
										Description: "Specifies the type of the content represented in the manifest, allowing Docker clients and registries to understand how to handle the document correctly.",
										Computed:    true,
									},
									"schema_version": schema.Int32Attribute{
										Description: "The version of the Docker Image Manifest format.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"query": schema.ListNestedBlock{
				Description: "A predefined set of criteria or conditions used to query and retrieve images within the org.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"fetch": schema.StringAttribute{
							Description: "Type of fetch. Specify either: `links` or `items`. Default: `items`.",
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("items", "links"),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"spec": schema.ListNestedBlock{
							Description: "The specification of the query.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"match": schema.StringAttribute{
										Description: "Type of match. Available values: `all`, `any`, `none`. Default: `all`.",
										Optional:    true,
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("all", "any", "none"),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"terms": schema.ListNestedBlock{
										Description: "Terms can only contain one of the following attributes: `property`, `rel`, `tag`.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"op": schema.StringAttribute{
													Description: "Type of query operation. Available values: `=`, `>`, `>=`, `<`, `<=`, `!=`, `exists`, `!exists`. Default: `=`.",
													Optional:    true,
													Computed:    true,
													Validators: []validator.String{
														stringvalidator.OneOf("=", ">", ">=", "<", "<=", "!=", "~", "exists", "!exists"),
													},
												},
												"property": schema.StringAttribute{
													Description: "Property to use for query evaluation.",
													Optional:    true,
												},
												"rel": schema.StringAttribute{
													Description: "Relation to use for query evaluation.",
													Optional:    true,
												},
												"tag": schema.StringAttribute{
													Description: "Tag key to use for query evaluation.",
													Optional:    true,
												},
												"value": schema.StringAttribute{
													Description: "Testing value for query evaluation.",
													Optional:    true,
												},
											},
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},
			},
		},
	}
}

// Read fetches the current state of the resource.
func (d *ImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state ImagesDataSourceModel

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Create a new operator instance
	operator := ImagesDataSourceOperator{
		Ctx:    ctx,
		Diags:  &resp.Diagnostics,
		Client: d.client,
		Plan:   state,
	}

	// Invoke API to read resource details
	apiResp, err := operator.InvokeRead()

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Build new state from API response
	newState := operator.MapResponseToState(apiResp)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

/*** Schemas ***/

// ConfigSchema returns a ListNestedAttribute schema for the ImagesDataSource using the provided description and nested computed attributes.
func (d *ImagesDataSource) ConfigSchema(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: description,
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"size": schema.Int32Attribute{
					Description: "The size of the image or layer in bytes. This helps in estimating the space required and the download time.",
					Computed:    true,
				},
				"digest": schema.StringAttribute{
					Description: "A unique SHA256 hash used to identify a specific image version within the image registry.",
					Computed:    true,
				},
				"media_type": schema.StringAttribute{
					Description: "Specifies the type of the content represented in the manifest, allowing Docker clients and registries to understand how to handle the document correctly.",
					Computed:    true,
				},
			},
		},
	}
}

/*** Data Source Operator ***/

// ImagesDataSourceOperator is the operator for managing the state.
type ImagesDataSourceOperator struct {
	Ctx    context.Context
	Diags  *diag.Diagnostics
	Client *client.Client
	Plan   ImagesDataSourceModel
}

// MapResponseToState creates a state model from response payload.
func (ieo *ImagesDataSourceOperator) MapResponseToState(images *client.ImagesQueryResult) ImagesDataSourceModel {
	// Initialize a new request payload
	state := ImagesDataSourceModel{}

	// Set specific attributes
	state.Images = ieo.flattenImages(&images.Items)
	state.Query = FlattenQuery(ieo.Ctx, ieo.Diags, &images.Query)

	// Return completed state model
	return state
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (ieo *ImagesDataSourceOperator) InvokeRead() (*client.ImagesQueryResult, error) {
	// Initialize default query to fetch all image resources
	query := client.Query{
		Kind: StringPointer("image"),
		Spec: &client.QuerySpec{
			Match: StringPointer("all"),
		},
	}

	// Check if the plan contains a non-null, known custom query
	if !ieo.Plan.Query.IsNull() && !ieo.Plan.Query.IsUnknown() {
		// Build a custom query from the plan's query expression
		plannedQuery := BuildQuery(ieo.Ctx, ieo.Diags, ieo.Plan.Query)

		// If the custom query is valid, use it instead of the default
		if plannedQuery != nil {
			query = *plannedQuery
		}
	}

	// Execute the image query and return the result with any error encountered
	return ieo.Client.GetImagesQuery(query)
}

// Flatteners //

// flattenImages transforms *[]client.Image into a Terraform types.List.
func (ieo *ImagesDataSourceOperator) flattenImages(input *[]client.Image) types.List {
	// Get attribute types
	elementType := models.ImageModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.ImageModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.ImageModel{
			CplnId:     types.StringPointerValue(item.ID),
			Name:       types.StringPointerValue(item.Name),
			Tags:       FlattenTags(item.Tags),
			SelfLink:   FlattenSelfLink(item.Links),
			Tag:        types.StringPointerValue(item.Tag),
			Repository: types.StringPointerValue(item.Repository),
			Digest:     types.StringPointerValue(item.Digest),
			Manifest:   ieo.flattenManifest(item.Manifest),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(ieo.Ctx, ieo.Diags, blocks)
}

// flattenManifest transforms *client.ImageManifest into a Terraform types.List.
func (ieo *ImagesDataSourceOperator) flattenManifest(input *client.ImageManifest) types.List {
	// Get attribute types
	elementType := models.ManifestModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.ManifestModel{
		Config:        ieo.flattenManifestConfigSingle(input.Config),
		Layers:        ieo.flattenManifestConfigMulti(input.Layers),
		MediaType:     types.StringPointerValue(input.MediaType),
		SchemaVersion: FlattenInt(input.SchemaVersion),
	}

	// Return the successfully created types.List
	return FlattenList(ieo.Ctx, ieo.Diags, []models.ManifestModel{block})
}

// flattenManifestConfigSingle transforms *client.ImageManifestConfig into a Terraform types.List.
func (ieo *ImagesDataSourceOperator) flattenManifestConfigSingle(input *client.ImageManifestConfig) types.List {
	// Get attribute types
	elementType := models.ManifestConfigModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.ManifestConfigModel{
		Size:      FlattenInt(input.Size),
		Digest:    types.StringPointerValue(input.Digest),
		MediaType: types.StringPointerValue(input.MediaType),
	}

	// Return the successfully created types.List
	return FlattenList(ieo.Ctx, ieo.Diags, []models.ManifestConfigModel{block})
}

// flattenManifestConfigMulti transforms *[]client.ImageManifestConfig into a Terraform types.List.
func (ieo *ImagesDataSourceOperator) flattenManifestConfigMulti(input *[]client.ImageManifestConfig) types.List {
	// Get attribute types
	elementType := models.ManifestConfigModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.ManifestConfigModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.ManifestConfigModel{
			Size:      FlattenInt(item.Size),
			Digest:    types.StringPointerValue(item.Digest),
			MediaType: types.StringPointerValue(item.MediaType),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(ieo.Ctx, ieo.Diags, blocks)
}
