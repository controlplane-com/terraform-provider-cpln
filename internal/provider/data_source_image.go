package cpln

import (
	"context"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/image"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &ImageDataSource{}
	_ datasource.DataSourceWithConfigure = &ImageDataSource{}
)

/*** Data Source Model ***/

// ImageDataSourceModel holds the Terraform state for the data source.
type ImageDataSourceModel struct {
	EntityBaseModel
	Tag        types.String `tfsdk:"tag"`
	Repository types.String `tfsdk:"repository"`
	Digest     types.String `tfsdk:"digest"`
	Manifest   types.List   `tfsdk:"manifest"`
}

/*** Data Source Configuration ***/

// ImageDataSource is the data source implementation.
type ImageDataSource struct {
	EntityBase
	Operations EntityOperations[ImageDataSourceModel, client.Image]
}

// NewImageDataSource returns a new instance of the data source implementation.
func NewImageDataSource() datasource.DataSource {
	return &ImageDataSource{}
}

// Metadata provides the data source type name.
func (d *ImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_image"
}

// Configure configures the data source before use.
func (d *ImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	d.Operations = NewEntityOperations(d.client, &ImageDataSourceOperator{})
}

// Schema defines the schema for the data source.
func (d *ImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this image.",
				Computed:    true,
			},
			"cpln_id": schema.StringAttribute{
				Description: "The ID, in GUID format, of the image.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the image.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the image.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value map of resource tags.",
				ElementType: types.StringType,
				Optional:    true,
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
	}
}

// Read fetches the current state of the resource.
func (d *ImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state ImageDataSourceModel

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Create a new operator instance
	operator := d.Operations.NewOperator(ctx, &resp.Diagnostics, state)

	// Invoke API to read resource details
	apiResp, code, err := operator.InvokeRead(state.Name.ValueString())

	// Remove resource from state if not found
	if code == 404 {
		// Drop resource from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Build new state from API response
	newState := operator.MapResponseToState(apiResp, true)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

/*** Schemas ***/

// ConfigSchema returns a ListNestedAttribute schema for the ImagesDataSource using the provided description and nested computed attributes.
func (d *ImageDataSource) ConfigSchema(description string) schema.ListNestedAttribute {
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

// ImageDataSourceOperator is the operator for managing the state.
type ImageDataSourceOperator struct {
	EntityOperator[ImageDataSourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (ieo *ImageDataSourceOperator) NewAPIRequest(isUpdate bool) client.Image {
	return client.Image{}
}

// MapResponseToState creates a state model from response payload.
func (ieo *ImageDataSourceOperator) MapResponseToState(image *client.Image, isCreate bool) ImageDataSourceModel {
	// Initialize a new request payload
	state := ImageDataSourceModel{}

	// Populate common fields from base resource data
	state.From(image.Base)

	// Set specific attributes
	state.Tag = types.StringPointerValue(image.Tag)
	state.Repository = types.StringPointerValue(image.Repository)
	state.Digest = types.StringPointerValue(image.Digest)
	state.Manifest = ieo.flattenManifest(image.Manifest)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (ieo *ImageDataSourceOperator) InvokeCreate(req client.Image) (*client.Image, int, error) {
	return nil, 0, nil
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (ieo *ImageDataSourceOperator) InvokeRead(name string) (*client.Image, int, error) {
	// Check if the provided name includes a tag separator (:)
	hasColon := len(strings.SplitN(name, ":", 2)) == 2

	// Initialize variables to capture the API response
	var image *client.Image
	var code int
	var err error

	// Use GetImage when a specific tag is provided, otherwise fetch the latest image
	if hasColon {
		image, code, err = ieo.Client.GetImage(name)
	} else {
		image, code, err = ieo.Client.GetLatestImage(name)
	}

	// Return the obtained image, status code, and error (if any)
	return image, code, err
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (ieo *ImageDataSourceOperator) InvokeUpdate(req client.Image) (*client.Image, int, error) {
	return nil, 0, nil
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (ieo *ImageDataSourceOperator) InvokeDelete(name string) error {
	return nil
}

// Flatteners //

// flattenManifest transforms *client.ImageManifest into a Terraform types.List.
func (ieo *ImageDataSourceOperator) flattenManifest(input *client.ImageManifest) types.List {
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
func (ieo *ImageDataSourceOperator) flattenManifestConfigSingle(input *client.ImageManifestConfig) types.List {
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
func (ieo *ImageDataSourceOperator) flattenManifestConfigMulti(input *[]client.ImageManifestConfig) types.List {
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
