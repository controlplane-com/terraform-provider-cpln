package cpln

import (
	"context"
	"fmt"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &HelmTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &HelmTemplateDataSource{}
)

// HelmTemplateDataSourceModel holds the Terraform state for the data source.
type HelmTemplateDataSourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Gvc                   types.String `tfsdk:"gvc"`
	Chart                 types.String `tfsdk:"chart"`
	Repository            types.String `tfsdk:"repository"`
	Version               types.String `tfsdk:"version"`
	Values                types.List   `tfsdk:"values"`
	Set                   types.Map    `tfsdk:"set"`
	SetString             types.Map    `tfsdk:"set_string"`
	SetFile               types.Map    `tfsdk:"set_file"`
	DependencyUpdate      types.Bool   `tfsdk:"dependency_update"`
	Description           types.String `tfsdk:"description"`
	Verify                types.Bool   `tfsdk:"verify"`
	RepositoryUsername    types.String `tfsdk:"repository_username"`
	RepositoryPassword    types.String `tfsdk:"repository_password"`
	RepositoryCaFile      types.String `tfsdk:"repository_ca_file"`
	RepositoryCertFile    types.String `tfsdk:"repository_cert_file"`
	RepositoryKeyFile     types.String `tfsdk:"repository_key_file"`
	InsecureSkipTLSVerify types.Bool   `tfsdk:"insecure_skip_tls_verify"`
	RenderSubchartNotes   types.Bool   `tfsdk:"render_subchart_notes"`
	Postrender            types.Object `tfsdk:"postrender"`
	Manifest              types.String `tfsdk:"manifest"`
}

// HelmTemplateDataSource is the data source implementation.
type HelmTemplateDataSource struct {
	client *client.Client
}

// NewHelmTemplateDataSource returns a new instance of the data source implementation.
func NewHelmTemplateDataSource() datasource.DataSource {
	return &HelmTemplateDataSource{}
}

// Metadata provides the data source type name.
func (d *HelmTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_helm_template"
}

// Configure configures the data source before use.
func (d *HelmTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

// Schema defines the schema for the data source.
func (d *HelmTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Renders Helm chart templates using the `cpln helm template` command without installing. Useful for previewing rendered manifests or feeding them into other resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this data source (same as name).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The release name to use for rendering the templates.",
				Required:    true,
			},
			"gvc": schema.StringAttribute{
				Description: "The GVC (Global Virtual Cloud) context for rendering the helm chart templates. Required only if the chart contains GVC-scoped resources and the GVC is not defined within the chart manifests.",
				Optional:    true,
			},
			"chart": schema.StringAttribute{
				Description: "Path to the chart. This can be a local path to a chart directory or packaged chart, or a URL/path when used with --repo.",
				Required:    true,
			},
			"repository": schema.StringAttribute{
				Description: "Chart repository URL where to locate the requested chart. Can be a Helm repository URL or an OCI registry URL.",
				Optional:    true,
			},
			"version": schema.StringAttribute{
				Description: "Specify a version constraint for the chart version to use. This can be a specific tag (e.g., 1.1.1) or a valid range (e.g., ^2.0.0). If not specified, the latest version is used.",
				Optional:    true,
			},
			"values": schema.ListAttribute{
				Description: "List of values in raw YAML to pass to helm.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"set": schema.MapAttribute{
				Description: "Set values on the command line. Map of key-value pairs. Equivalent to using --set flag.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"set_string": schema.MapAttribute{
				Description: "Set STRING values on the command line. Map of key-value pairs. Equivalent to using --set-string flag.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"set_file": schema.MapAttribute{
				Description: "Set values from files specified via the command line. Map of key to file path. Equivalent to using --set-file flag.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"dependency_update": schema.BoolAttribute{
				Description: "Update dependencies if they are missing before rendering the chart.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Add a custom description.",
				Optional:    true,
			},
			"verify": schema.BoolAttribute{
				Description: "Verify the package before using it.",
				Optional:    true,
			},
			"repository_username": schema.StringAttribute{
				Description: "Chart repository username where to locate the requested chart.",
				Optional:    true,
			},
			"repository_password": schema.StringAttribute{
				Description: "Chart repository password where to locate the requested chart.",
				Optional:    true,
				Sensitive:   true,
			},
			"repository_ca_file": schema.StringAttribute{
				Description: "Verify certificates of HTTPS-enabled servers using this CA bundle.",
				Optional:    true,
			},
			"repository_cert_file": schema.StringAttribute{
				Description: "Identify HTTPS client using this SSL certificate file.",
				Optional:    true,
			},
			"repository_key_file": schema.StringAttribute{
				Description: "Identify HTTPS client using this SSL key file.",
				Optional:    true,
			},
			"insecure_skip_tls_verify": schema.BoolAttribute{
				Description: "Skip TLS certificate checks for the chart download.",
				Optional:    true,
			},
			"render_subchart_notes": schema.BoolAttribute{
				Description: "If set, render subchart notes along with the parent.",
				Optional:    true,
			},
			"postrender": schema.SingleNestedAttribute{
				Description: "Post-renderer configuration. Specifies a binary to run after helm renders the manifests.",
				Optional:    true,
				Attributes:  postrenderDataSourceSchemaAttributes(),
			},
			"manifest": schema.StringAttribute{
				Description: "The rendered manifest output from helm template.",
				Computed:    true,
			},
		},
	}
}

// postrenderDataSourceSchemaAttributes returns the schema attributes for the postrender nested block in the data source.
func postrenderDataSourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"binary_path": schema.StringAttribute{
			Description: "The path to an executable to be used for post rendering.",
			Required:    true,
		},
		"args": schema.ListAttribute{
			Description: "Arguments to the post-renderer.",
			ElementType: types.StringType,
			Optional:    true,
		},
	}
}

// Read executes cpln helm template and stores the rendered manifest.
func (d *HelmTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read the config from the request
	var config HelmTemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the command arguments
	args := []string{"helm", "template", config.Name.ValueString(), config.Chart.ValueString()}

	// Build the common config from the data source model
	cfg := client.HelmCommonConfig{
		Gvc:                   config.Gvc,
		Repository:            config.Repository,
		Version:               config.Version,
		Values:                config.Values,
		Set:                   config.Set,
		SetString:             config.SetString,
		SetFile:               config.SetFile,
		Description:           config.Description,
		Verify:                config.Verify,
		RepositoryUsername:    config.RepositoryUsername,
		RepositoryPassword:    config.RepositoryPassword,
		RepositoryCaFile:      config.RepositoryCaFile,
		RepositoryCertFile:    config.RepositoryCertFile,
		RepositoryKeyFile:     config.RepositoryKeyFile,
		InsecureSkipTLSVerify: config.InsecureSkipTLSVerify,
		RenderSubchartNotes:   config.RenderSubchartNotes,
		Postrender:            config.Postrender,
		DependencyUpdate:      config.DependencyUpdate,
	}

	// Build and validate the full argument list
	commonArgs, tempFiles, err := d.client.BuildHelmArgs(args, cfg)
	defer client.RemoveTempFiles(tempFiles)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build helm template arguments", err.Error())
		return
	}

	// Execute helm template
	output, err := ExecuteCplnCommand(commonArgs)
	if err != nil {
		resp.Diagnostics.AddError("Helm template failed", err.Error())
		return
	}

	// Set computed fields
	config.ID = config.Name
	config.Manifest = types.StringValue(strings.TrimSpace(output))

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
