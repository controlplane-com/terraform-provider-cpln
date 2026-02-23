package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
)

// Ensure resource implements required interfaces at compile time
var (
	_ resource.Resource                = &HelmReleaseResource{}
	_ resource.ResourceWithImportState = &HelmReleaseResource{}
)

/*** Resource Model ***/

// HelmReleaseResourceModel holds the Terraform state for the resource.
type HelmReleaseResourceModel struct {
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
	Wait                  types.Bool   `tfsdk:"wait"`
	Timeout               types.Int32  `tfsdk:"timeout"`
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
	MaxHistory            types.Int32  `tfsdk:"max_history"`
	Status                types.String `tfsdk:"status"`
	Revision              types.Int32  `tfsdk:"revision"`
	Manifest              types.String `tfsdk:"manifest"`
	Resources             types.Map    `tfsdk:"resources"`
}

// GetID returns the ID field from the helm release resource model.
func (m HelmReleaseResourceModel) GetID() types.String {
	return m.ID
}

/*** Resource Configuration ***/

// HelmReleaseResource is the resource implementation.
type HelmReleaseResource struct {
	EntityBase
	Operations EntityOperations[HelmReleaseResourceModel, client.HelmReleaseState]
}

// NewHelmReleaseResource returns a new instance of the resource implementation.
func NewHelmReleaseResource() resource.Resource {
	return &HelmReleaseResource{}
}

// Configure configures the resource before use.
func (r *HelmReleaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	r.Operations = NewEntityOperations(r.client, &HelmReleaseResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (r *HelmReleaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (r *HelmReleaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_helm_release"
}

// Schema defines the schema for the resource.
func (r *HelmReleaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Helm chart deployments on Control Plane using the `cpln helm` command. This resource allows you to install, upgrade, and uninstall Helm charts that deploy Control Plane resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this helm release (same as name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The release name for this helm deployment.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gvc": schema.StringAttribute{
				Description: "The GVC (Global Virtual Cloud) to use for the helm deployment. Required only if the chart deploys GVC-scoped resources and the GVC is not defined within the chart manifests.",
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
			"wait": schema.BoolAttribute{
				Description: "If set to true, will wait until all Workloads are in a ready state before marking the release as successful.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"timeout": schema.Int32Attribute{
				Description: "The amount of seconds to wait for workloads to be ready before timing out. Only used when wait is true. Default is 300 seconds.",
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(300),
			},
			"dependency_update": schema.BoolAttribute{
				Description: "Update dependencies if they are missing before installing the chart.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: "Add a custom description for the release.",
				Optional:    true,
			},
			"verify": schema.BoolAttribute{
				Description: "Verify the package before using it.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"render_subchart_notes": schema.BoolAttribute{
				Description: "If set, render subchart notes along with the parent on install/upgrade.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"postrender": schema.SingleNestedAttribute{
				Description: "Post-renderer configuration. Specifies a binary to run after helm renders the manifests.",
				Optional:    true,
				Attributes:  postrenderSchemaAttributes(),
			},
			"max_history": schema.Int32Attribute{
				Description: "Maximum number of revisions saved per release. Use 0 for no limit. Default is 10. Only used on upgrade.",
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(10),
			},
			"status": schema.StringAttribute{
				Description: "The current status of the helm release.",
				Computed:    true,
			},
			"revision": schema.Int32Attribute{
				Description: "The current revision number of the helm release.",
				Computed:    true,
			},
			"manifest": schema.StringAttribute{
				Description: "The rendered manifest of the helm release.",
				Computed:    true,
			},
			"resources": schema.MapAttribute{
				Description: "Rendered manifests keyed by kind/gvc/name (e.g., workload/my-gvc/my-workload). GVC may be empty for GVC-level resources.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Create creates the resource.
func (r *HelmReleaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, r.Operations)
}

// Read fetches the current state of the resource.
func (r *HelmReleaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, r.Operations)
}

// Update modifies the resource.
func (r *HelmReleaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, r.Operations)
}

// Delete removes the resource.
func (r *HelmReleaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, r.Operations)
}

/*** Schemas ***/

// postrenderSchemaAttributes returns the schema attributes for the postrender nested block.
func postrenderSchemaAttributes() map[string]schema.Attribute {
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

/*** Resource Operator ***/

// HelmReleaseResourceOperator is the operator for managing the helm release state.
type HelmReleaseResourceOperator struct {
	EntityOperator[HelmReleaseResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (op *HelmReleaseResourceOperator) NewAPIRequest(isUpdate bool) client.HelmReleaseState {
	return client.HelmReleaseState{
		Name: op.Plan.Name.ValueString(),
	}
}

// MapResponseToState creates a state model from response payload.
func (op *HelmReleaseResourceOperator) MapResponseToState(resp *client.HelmReleaseState, isCreate bool) HelmReleaseResourceModel {
	// Start from the current plan/state to preserve config fields
	state := op.Plan

	// Set computed fields from response
	state.Name = types.StringValue(resp.Name)
	state.ID = state.Name
	state.Status = types.StringValue(resp.Status)
	state.Revision = types.Int32Value(int32(resp.Revision))
	state.Manifest = types.StringValue(resp.Manifest)

	// Build the resources map
	state.Resources = op.flattenResources(resp.Resources)

	return state
}

// InvokeCreate invokes helm install to create a new release.
// If install fails with a 409 conflict (release already exists from a partial install), it falls back to helm upgrade.
func (op *HelmReleaseResourceOperator) InvokeCreate(req client.HelmReleaseState) (*client.HelmReleaseState, int, error) {
	// Build the command arguments
	args := []string{"helm", "install", op.Plan.Name.ValueString(), op.Plan.Chart.ValueString()}

	// Add common arguments
	commonArgs, tempFiles, err := op.Client.BuildHelmArgs(args, op.buildCommonConfig())
	defer client.RemoveTempFiles(tempFiles)
	if err != nil {
		return nil, 0, fmt.Errorf("could not install release %s: %w", req.Name, err)
	}

	// Execute helm install
	if _, err := ExecuteCplnCommand(commonArgs); err != nil {
		// On 409 conflict, the release already exists (e.g. a previous install partially succeeded), fall back to upgrade
		if strings.Contains(err.Error(), "409") {
			return op.InvokeUpdate(req)
		}

		return nil, 0, fmt.Errorf("could not install release %s: %w", req.Name, err)
	}

	// Fetch release info
	return op.getRelease(req.Name)
}

// InvokeRead fetches the current state of an existing release.
func (op *HelmReleaseResourceOperator) InvokeRead(name string) (*client.HelmReleaseState, int, error) {
	return op.getRelease(name)
}

// InvokeUpdate invokes helm upgrade to update an existing release.
func (op *HelmReleaseResourceOperator) InvokeUpdate(req client.HelmReleaseState) (*client.HelmReleaseState, int, error) {
	// Build the command arguments
	args := []string{"helm", "upgrade", op.Plan.Name.ValueString(), op.Plan.Chart.ValueString()}

	// Add common arguments including max_history for upgrade
	cfg := op.buildCommonConfig()
	cfg.MaxHistory = op.Plan.MaxHistory

	// Build and validate the full argument list
	commonArgs, tempFiles, err := op.Client.BuildHelmArgs(args, cfg)
	defer client.RemoveTempFiles(tempFiles)
	if err != nil {
		return nil, 0, fmt.Errorf("could not upgrade release %s: %w", req.Name, err)
	}

	// Execute helm upgrade with retry on 409 conflict
	if err := executeHelmWithRetry(commonArgs); err != nil {
		return nil, 0, fmt.Errorf("could not upgrade release %s: %w", req.Name, err)
	}

	// Fetch release info
	return op.getRelease(req.Name)
}

// InvokeDelete invokes helm uninstall to delete a release.
func (op *HelmReleaseResourceOperator) InvokeDelete(name string) error {
	// Build the command arguments
	args := []string{"helm", "uninstall", name}
	args = op.Client.AppendCplnContextArgs(args, op.Plan.Gvc.ValueString())

	// Execute helm uninstall
	_, err := ExecuteCplnCommand(args)
	if err != nil {
		// If release is already gone, that's fine
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "release: not found") {
			return nil
		}

		return fmt.Errorf("could not uninstall release %s: %w", name, err)
	}

	return nil
}

// Flatteners //

// flattenResources converts a map of rendered manifests to a Terraform map type.
func (op *HelmReleaseResourceOperator) flattenResources(resources map[string]string) types.Map {
	if len(resources) == 0 {
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value, len(resources))
	for key, value := range resources {
		elements[key] = types.StringValue(value)
	}

	return types.MapValueMust(types.StringType, elements)
}

// Helpers //

// getRelease fetches all release info using cpln helm get all <name> -o json.
func (op *HelmReleaseResourceOperator) getRelease(releaseName string) (*client.HelmReleaseState, int, error) {
	// Build the command arguments
	args := []string{"helm", "get", "all", releaseName, "-o", "json"}
	args = op.Client.AppendCplnAuthArgs(args)

	// Execute helm get all
	output, err := ExecuteCplnCommand(args)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, 404, fmt.Errorf("release %s not found", releaseName)
		}

		return nil, 0, fmt.Errorf("could not get release %s: %w", releaseName, err)
	}

	// Parse the JSON response
	var resp client.HelmGetAllResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse helm get all output: %w", err)
	}

	// Parse manifest YAML documents into resources map
	resources := op.parseManifestResources(resp.Manifest)

	return &client.HelmReleaseState{
		Name:      resp.Name,
		Status:    resp.Info.Status,
		Revision:  resp.Version,
		Manifest:  strings.TrimSpace(resp.Manifest),
		Resources: resources,
	}, 0, nil
}

// parseManifestResources splits a multi-document YAML manifest into a map keyed by resource identity.
func (op *HelmReleaseResourceOperator) parseManifestResources(manifest string) map[string]string {
	resources := make(map[string]string)

	for doc := range strings.SplitSeq(manifest, "---") {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		var res client.HelmManifestResource
		if err := yaml.Unmarshal([]byte(doc), &res); err != nil {
			continue
		}

		if res.Kind == "" || res.Name == "" {
			continue
		}

		kind := strings.ToLower(res.Kind)

		// GVC-scoped resources include the GVC in the key for uniqueness
		if IsGvcScopedResource(kind) {
			resources[fmt.Sprintf("%s/%s/%s", kind, res.Gvc, res.Name)] = doc
		} else {
			resources[fmt.Sprintf("%s/%s", kind, res.Name)] = doc
		}
	}

	return resources
}

// executeHelmWithRetry executes a cpln helm command with retry logic on 409 conflict errors.
func executeHelmWithRetry(args []string) error {
	const maxRetries = 5
	backoff := 2 * time.Second
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err := ExecuteCplnCommand(args)
		if err == nil {
			return nil
		}

		lastErr = err

		// Retry on 409 conflict
		if strings.Contains(err.Error(), "409") {
			if attempt == maxRetries {
				break
			}

			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		// Non-retryable error
		return err
	}

	return lastErr
}

// buildCommonConfig builds a HelmCommonConfig from the current plan.
func (op *HelmReleaseResourceOperator) buildCommonConfig() client.HelmCommonConfig {
	return client.HelmCommonConfig{
		Gvc:                   op.Plan.Gvc,
		Repository:            op.Plan.Repository,
		Version:               op.Plan.Version,
		Values:                op.Plan.Values,
		Set:                   op.Plan.Set,
		SetString:             op.Plan.SetString,
		SetFile:               op.Plan.SetFile,
		Wait:                  op.Plan.Wait,
		Timeout:               op.Plan.Timeout,
		Description:           op.Plan.Description,
		Verify:                op.Plan.Verify,
		RepositoryUsername:    op.Plan.RepositoryUsername,
		RepositoryPassword:    op.Plan.RepositoryPassword,
		RepositoryCaFile:      op.Plan.RepositoryCaFile,
		RepositoryCertFile:    op.Plan.RepositoryCertFile,
		RepositoryKeyFile:     op.Plan.RepositoryKeyFile,
		InsecureSkipTLSVerify: op.Plan.InsecureSkipTLSVerify,
		RenderSubchartNotes:   op.Plan.RenderSubchartNotes,
		Postrender:            op.Plan.Postrender,
		DependencyUpdate:      op.Plan.DependencyUpdate,
	}
}

