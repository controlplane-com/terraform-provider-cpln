package cpln

import (
	"context"
	"fmt"
	"sync"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/org"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceLock = &sync.Mutex{}
var loggingTypes = []string{
	"s3_logging", "coralogix_logging", "datadog_logging", "logzio_logging",
	"elastic_logging", "cloud_watch_logging", "fluentd_logging", "stackdriver_logging",
	"syslog_logging",
}

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &OrgLoggingResource{}
	_ resource.ResourceWithImportState = &OrgLoggingResource{}
)

/*** Resource Model ***/

// OrgLoggingResourceModel holds the Terraform state for the resource.
type OrgLoggingResourceModel struct {
	ID                 types.String                     `tfsdk:"id"`
	CplnID             types.String                     `tfsdk:"cpln_id"`
	Name               types.String                     `tfsdk:"name"`
	Description        types.String                     `tfsdk:"description"`
	Tags               types.Map                        `tfsdk:"tags"`
	SelfLink           types.String                     `tfsdk:"self_link"`
	S3Logging          []models.S3LoggingModel          `tfsdk:"s3_logging"`
	CoralogixLogging   []models.CoralogixLoggingModel   `tfsdk:"coralogix_logging"`
	DatadogLogging     []models.DatadogLoggingModel     `tfsdk:"datadog_logging"`
	LogzioLogging      []models.LogzioLoggingModel      `tfsdk:"logzio_logging"`
	ElasticLogging     []models.ElasticLoggingModel     `tfsdk:"elastic_logging"`
	CloudWatchLogging  []models.CloudWatchModel         `tfsdk:"cloud_watch_logging"`
	FluentdLogging     []models.FluentdLoggingModel     `tfsdk:"fluentd_logging"`
	StackdriverLogging []models.StackdriverLoggingModel `tfsdk:"stackdriver_logging"`
	SyslogLogging      []models.SyslogLoggingModel      `tfsdk:"syslog_logging"`
}

/*** Resource Configuration ***/

// OrgLoggingResource is the resource implementation.
type OrgLoggingResource struct {
	EntityBase
}

// NewOrgLoggingResource returns a new instance of the resource implementation.
func NewOrgLoggingResource() resource.Resource {
	return &OrgLoggingResource{}
}

// Configure configures the resource before use.
func (olr *OrgLoggingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	olr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (olr *OrgLoggingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (olr *OrgLoggingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_org_logging"
}

// Schema defines the schema for the resource.
func (olr *OrgLoggingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this Org Logging.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cpln_id": schema.StringAttribute{
				Description: "The ID, in GUID format, of the Org.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the Org.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the Org.",
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"s3_logging": schema.ListNestedBlock{
				MarkdownDescription: "[Documentation Reference](https://docs.controlplane.com/external-logging/s3)",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"bucket": schema.StringAttribute{
							Description: "Name of S3 bucket.",
							Required:    true,
						},
						"region": schema.StringAttribute{
							Description: "AWS region where bucket is located.",
							Required:    true,
						},
						"prefix": schema.StringAttribute{
							Description: "Bucket path prefix. Default: \"/\".",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("/"),
						},
						"credentials": schema.StringAttribute{
							Description: "Full link to referenced AWS Secret.",
							Required:    true,
						},
					},
				},
			},
			"coralogix_logging": schema.ListNestedBlock{
				MarkdownDescription: "[Documentation Reference](https://docs.controlplane.com/external-logging/coralogix)",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cluster": schema.StringAttribute{
							Description: "Coralogix cluster URI.",
							Required:    true,
						},
						"credentials": schema.StringAttribute{
							Description: "Full link to referenced Opaque Secret.",
							Required:    true,
						},
						"app": schema.StringAttribute{
							Description: "App name to be displayed in Coralogix dashboard.",
							Optional:    true,
						},
						"subsystem": schema.StringAttribute{
							Description: "Subsystem name to be displayed in Coralogix dashboard.",
							Optional:    true,
						},
					},
				},
			},
			"datadog_logging": schema.ListNestedBlock{
				MarkdownDescription: "[Documentation Reference](https://docs.controlplane.com/external-logging/datadog)",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Description: "Datadog host URI.",
							Required:    true,
						},
						"credentials": schema.StringAttribute{
							Description: "Full link to referenced Opaque Secret.",
							Required:    true,
						},
					},
				},
			},
			"logzio_logging": schema.ListNestedBlock{
				MarkdownDescription: "[Documentation Reference](https://docs.controlplane.com/external-logging/logz-io)",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"listener_host": schema.StringAttribute{
							Description: "Logzio listener host URI.",
							Required:    true,
						},
						"credentials": schema.StringAttribute{
							Description: "Full link to referenced Opaque Secret.",
							Required:    true,
						},
					},
				},
			},
			"elastic_logging": schema.ListNestedBlock{
				Description: "For logging and analyzing data within an org using Elastic Logging.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"aws": schema.ListNestedBlock{
							Description: "For targeting Amazon Web Services (AWS) ElasticSearch.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										Description: "A valid AWS ElasticSearch hostname (must end with es.amazonaws.com).",
										Required:    true,
									},
									"port": schema.Int32Attribute{
										Description: "Port. Default: 443",
										Required:    true,
									},
									"index": schema.StringAttribute{
										Description: "Logging Index.",
										Required:    true,
									},
									"type": schema.StringAttribute{
										Description: "Logging Type.",
										Required:    true,
									},
									"credentials": schema.StringAttribute{
										Description: "Full Link to a secret of type `aws`.",
										Required:    true,
									},
									"region": schema.StringAttribute{
										Description: "Valid AWS region.",
										Required:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"elastic_cloud": schema.ListNestedBlock{
							Description: "For targeting Elastic Cloud.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"index": schema.StringAttribute{
										Description: "Logging Index.",
										Required:    true,
									},
									"type": schema.StringAttribute{
										Description: "Logging Type.",
										Required:    true,
									},
									"credentials": schema.StringAttribute{
										Description: "Full Link to a secret of type `userpass`.",
										Required:    true,
									},
									"cloud_id": schema.StringAttribute{
										MarkdownDescription: "[Cloud ID](https://www.elastic.co/guide/en/cloud/current/ec-cloud-id.html)",
										Required:            true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"generic": schema.ListNestedBlock{
							Description: "For targeting generic Elastic Search providers.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										Description: "A valid Elastic Search provider hostname.",
										Required:    true,
									},
									"port": schema.Int32Attribute{
										Description: "Port. Default: 443",
										Required:    true,
									},
									"path": schema.StringAttribute{
										Description: "Logging path.",
										Required:    true,
									},
									"index": schema.StringAttribute{
										Description: "Logging Index.",
										Required:    true,
									},
									"type": schema.StringAttribute{
										Description: "Logging Type.",
										Required:    true,
									},
									"credentials": schema.StringAttribute{
										Description: "Full Link to a secret of type `userpass`.",
										Required:    true,
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
			"cloud_watch_logging": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "Valid AWS region.",
							Required:    true,
						},
						"credentials": schema.StringAttribute{
							Description: "Full Link to a secret of type `opaque`.",
							Required:    true,
						},
						"retention_days": schema.Int32Attribute{
							Description: "Length, in days, for how log data is kept before it is automatically deleted.",
							Optional:    true,
						},
						"group_name": schema.StringAttribute{
							Description: "A container for log streams with common settings like retention. Used to categorize logs by application or service type.",
							Required:    true,
						},
						"stream_name": schema.StringAttribute{
							Description: "A sequence of log events from the same source within a log group. Typically represents individual instances of services or applications.",
							Required:    true,
						},
						"extract_fields": schema.MapAttribute{
							Description: "Enable custom data extraction from log entries for enhanced querying and analysis.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"fluentd_logging": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Description: "The hostname or IP address of a remote log storage system.",
							Required:    true,
						},
						"port": schema.Int32Attribute{
							Description: "Port. Default: 24224",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(24224),
						},
					},
				},
			},
			"stackdriver_logging": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"credentials": schema.StringAttribute{
							Description: "Full Link to a secret of type `opaque`.",
							Required:    true,
						},
						"location": schema.StringAttribute{
							Description: "A Google Cloud Provider region.",
							Required:    true,
						},
					},
				},
			},
			"syslog_logging": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Description: "Hostname of Syslog Endpoint.",
							Required:    true,
						},
						"port": schema.Int32Attribute{
							Description: "Port of Syslog Endpoint.",
							Required:    true,
						},
						"mode": schema.StringAttribute{
							Description: "Log Mode. Valid values: TCP, TLS, or UDP.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("tcp"),
						},
						"format": schema.StringAttribute{
							Description: "Log Format. Valid values: RFC3164 or RFC5424.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("rfc5424"),
						},
						"severity": schema.Int32Attribute{
							Description: "Severity Level. See documentation for details. Valid values: 0 to 7.",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(6),
						},
					},
				},
			},
		},
	}
}

// Create creates the resource.
func (olr *OrgLoggingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Acquire lock to ensure only one operation modifies the resource at a time
	resourceLock.Lock()
	defer resourceLock.Unlock()

	// Declare variable to hold the planned state from Terraform configuration
	var plannedState OrgLoggingResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the planned state
	loggings := olr.buildRequest(ctx, &resp.Diagnostics, plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the create request to the API client
	responsePayload, _, err := olr.client.UpdateOrgLogging(loggings)

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating org logging: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := olr.buildState(responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (olr *OrgLoggingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState OrgLoggingResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the org
	responsePayload, code, err := olr.client.GetOrg()

	// Handle the case where the org is not found (HTTP 404),
	// indicating it has been deleted outside of Terraform. Remove it from state
	if code == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle any other errors that occur during the API call
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading org logging: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := olr.buildState(responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Update modifies the resource.
func (olr *OrgLoggingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedState OrgLoggingResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the planned state
	loggings := olr.buildRequest(ctx, &resp.Diagnostics, plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the update request to the API with the modified data
	responsePayload, _, err := olr.client.UpdateOrgLogging(loggings)

	// Handle errors from the API update request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating org logging: %s", err))
		return
	}

	// Map the API response to the Terraform finalState
	finalState := olr.buildState(responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Delete removes the resource.
func (olr *OrgLoggingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrgLoggingResourceModel

	// Retrieve the state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Send a delete request to the API using the name from the state
	_, _, err := olr.client.UpdateOrgLogging(nil)

	// Handle errors from the API delete request
	if err != nil {
		// If an error occurs during the delete request, add an error to diagnostics
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting org logging: %s", err))
		return
	}

	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}

/*** Operations ***/

// buildRequest creates a request payload from a state model.
func (olr *OrgLoggingResource) buildRequest(ctx context.Context, diags *diag.Diagnostics, state OrgLoggingResourceModel) *[]client.Logging {
	// Initialize slice to hold logging configurations
	loggings := []client.Logging{}

	// Iterate over each logging type defined and collect loggings
	for _, loggingType := range loggingTypes {
		// Placeholder for logging entries to add
		var loggingToAdd *[]client.Logging

		switch loggingType {
		case "s3_logging":
			loggingToAdd = olr.buildS3Logging(state.S3Logging)
		case "coralogix_logging":
			loggingToAdd = olr.buildCoralogixLogging(state.CoralogixLogging)
		case "datadog_logging":
			loggingToAdd = olr.buildDatadogLogging(state.DatadogLogging)
		case "logzio_logging":
			loggingToAdd = olr.buildLogzioLogging(state.LogzioLogging)
		case "elastic_logging":
			loggingToAdd = olr.buildElasticLogging(state.ElasticLogging)
		case "cloud_watch_logging":
			loggingToAdd = olr.buildCloudWatchLogging(ctx, diags, state.CloudWatchLogging)
		case "fluentd_logging":
			loggingToAdd = olr.buildFluentdLogging(state.FluentdLogging)
		case "stackdriver_logging":
			loggingToAdd = olr.buildStackdriverLogging(state.StackdriverLogging)
		case "syslog_logging":
			loggingToAdd = olr.buildSyslogLogging(state.SyslogLogging)
		default:
			continue
		}

		// If logging configuration exists, append to payload
		if loggingToAdd != nil {
			loggings = append(loggings, *loggingToAdd...)
		}
	}

	// Validate loggings
	olr.validateLoggings(diags, loggings)

	// Return constructed request payload
	return &loggings
}

// buildState creates a state model from response payload.
func (olr *OrgLoggingResource) buildState(apiResp *client.Org) OrgLoggingResourceModel {
	// Initialize empty state model
	state := OrgLoggingResourceModel{}

	// Set specific attributes
	state.ID = types.StringPointerValue(apiResp.Name)
	state.CplnID = types.StringPointerValue(apiResp.ID)
	state.Name = types.StringPointerValue(apiResp.Name)
	state.Description = types.StringPointerValue(apiResp.Description)
	state.Tags = FlattenTags(apiResp.Tags)
	state.SelfLink = FlattenSelfLink(apiResp.Links)

	// Only process logging if Spec is non-nil
	if apiResp.Spec != nil {
		// Initialize slice to collect all logging entries
		loggings := []client.Logging{}

		// Append primary logging entries if defined
		if apiResp.Spec.Logging != nil {
			loggings = append(loggings, *apiResp.Spec.Logging)
		}

		// Append extra logging entries if present
		if apiResp.Spec.ExtraLogging != nil && len(*apiResp.Spec.ExtraLogging) > 0 {
			loggings = append(loggings, *apiResp.Spec.ExtraLogging...)
		}

		// Declare slices for each logging type
		var s3Array []client.S3Logging
		var coralogixArray []client.CoralogixLogging
		var dataDogArray []client.DatadogLogging
		var logzioArray []client.LogzioLogging
		var elasticArray []client.ElasticLogging
		var cloudWatchArray []client.CloudWatchLogging
		var fluentdArray []client.FluentdLogging
		var stackdriverArray []client.StackdriverLogging
		var syslogArray []client.SyslogLogging

		// Iterate over each logging entry to categorize by type
		for _, logging := range loggings {
			// Collect S3 logging entries
			if logging.S3 != nil {
				s3Array = append(s3Array, *logging.S3)
			}

			// Collect Coralogix logging entries
			if logging.Coralogix != nil {
				coralogixArray = append(coralogixArray, *logging.Coralogix)
			}

			// Collect Datadog logging entries
			if logging.Datadog != nil {
				dataDogArray = append(dataDogArray, *logging.Datadog)
			}

			// Collect Logz.io logging entries
			if logging.Logzio != nil {
				logzioArray = append(logzioArray, *logging.Logzio)
			}

			// Collect Elastic logging entries
			if logging.Elastic != nil {
				elasticArray = append(elasticArray, *logging.Elastic)
			}

			// Collect CloudWatch logging entries
			if logging.CloudWatch != nil {
				cloudWatchArray = append(cloudWatchArray, *logging.CloudWatch)
			}

			// Collect Fluentd logging entries
			if logging.Fluentd != nil {
				fluentdArray = append(fluentdArray, *logging.Fluentd)
			}

			// Collect Stackdriver logging entries
			if logging.Stackdriver != nil {
				stackdriverArray = append(stackdriverArray, *logging.Stackdriver)
			}

			// Collect Syslog logging entries
			if logging.Syslog != nil {
				syslogArray = append(syslogArray, *logging.Syslog)
			}
		}

		// Flatten loggings
		state.S3Logging = olr.flattenS3Logging(&s3Array)
		state.CoralogixLogging = olr.flattenCoralogixLogging(&coralogixArray)
		state.DatadogLogging = olr.flattenDatadogLogging(&dataDogArray)
		state.LogzioLogging = olr.flattenLogzioLogging(&logzioArray)
		state.ElasticLogging = olr.flattenElasticLogging(&elasticArray)
		state.CloudWatchLogging = olr.flattenCloudWatchLogging(&cloudWatchArray)
		state.FluentdLogging = olr.flattenFluentdLogging(&fluentdArray)
		state.StackdriverLogging = olr.flattenStackdriverLogging(&stackdriverArray)
		state.SyslogLogging = olr.flattenSyslogLogging(&syslogArray)
	}

	// Return completed state model
	return state
}

// Builders //

// buildS3Logging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildS3Logging(state []models.S3LoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.S3Logging{
			Bucket:      BuildString(block.Bucket),
			Region:      BuildString(block.Region),
			Prefix:      BuildString(block.Prefix),
			Credentials: BuildString(block.Credentials),
		}

		// Construct logging
		logging := client.Logging{
			S3: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildCoralogixLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildCoralogixLogging(state []models.CoralogixLoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.CoralogixLogging{
			Cluster:     BuildString(block.Cluster),
			Credentials: BuildString(block.Credentials),
			App:         BuildString(block.App),
			Subsystem:   BuildString(block.Subsystem),
		}

		// Construct logging
		logging := client.Logging{
			Coralogix: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildDatadogLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildDatadogLogging(state []models.DatadogLoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.DatadogLogging{
			Host:        BuildString(block.Host),
			Credentials: BuildString(block.Credentials),
		}

		// Construct logging
		logging := client.Logging{
			Datadog: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildLogzioLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildLogzioLogging(state []models.LogzioLoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.LogzioLogging{
			ListenerHost: BuildString(block.ListenerHost),
			Credentials:  BuildString(block.Credentials),
		}

		// Construct logging
		logging := client.Logging{
			Logzio: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildElasticLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildElasticLogging(state []models.ElasticLoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.ElasticLogging{
			AWS:          olr.buildElasticLoggingAws(block.AWS),
			ElasticCloud: olr.buildElasticLoggingElasticCloud(block.ElasticCloud),
			Generic:      olr.buildElasticLoggingGeneric(block.Generic),
		}

		// Construct logging
		logging := client.Logging{
			Elastic: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildElasticLoggingAws constructs a AWSLogging from the given Terraform state.
func (olr *OrgLoggingResource) buildElasticLoggingAws(state []models.ElasticLoggingAwsModel) *client.AWSLogging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.AWSLogging{
		Host:        BuildString(block.Host),
		Port:        BuildInt(block.Port),
		Index:       BuildString(block.Index),
		Type:        BuildString(block.Type),
		Credentials: BuildString(block.Credentials),
		Region:      BuildString(block.Region),
	}
}

// buildElasticLoggingElasticCloud constructs a ElasticCloudLogging from the given Terraform state.
func (olr *OrgLoggingResource) buildElasticLoggingElasticCloud(state []models.ElasticLoggingElasticCloudModel) *client.ElasticCloudLogging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.ElasticCloudLogging{
		Index:       BuildString(block.Index),
		Type:        BuildString(block.Type),
		Credentials: BuildString(block.Credentials),
		CloudID:     BuildString(block.CloudID),
	}
}

// buildElasticLoggingGeneric constructs a GenericLogging from the given Terraform state.
func (olr *OrgLoggingResource) buildElasticLoggingGeneric(state []models.ElasticLoggingGenericModel) *client.GenericLogging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.GenericLogging{
		Host:        BuildString(block.Host),
		Port:        BuildInt(block.Port),
		Path:        BuildString(block.Path),
		Index:       BuildString(block.Index),
		Type:        BuildString(block.Type),
		Credentials: BuildString(block.Credentials),
	}
}

// buildCloudWatchLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildCloudWatchLogging(ctx context.Context, diags *diag.Diagnostics, state []models.CloudWatchModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.CloudWatchLogging{
			Region:        BuildString(block.Region),
			Credentials:   BuildString(block.Credentials),
			RetentionDays: BuildInt(block.RetentionDays),
			GroupName:     BuildString(block.GroupName),
			StreamName:    BuildString(block.StreamName),
			ExtractFields: BuildMapString(ctx, diags, block.ExtractFields),
		}

		// Construct logging
		logging := client.Logging{
			CloudWatch: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildFluentdLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildFluentdLogging(state []models.FluentdLoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.FluentdLogging{
			Host: BuildString(block.Host),
			Port: BuildInt(block.Port),
		}

		// Construct logging
		logging := client.Logging{
			Fluentd: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildStackdriverLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildStackdriverLogging(state []models.StackdriverLoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.StackdriverLogging{
			Credentials: BuildString(block.Credentials),
			Location:    BuildString(block.Location),
		}

		// Construct logging
		logging := client.Logging{
			Stackdriver: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// buildSyslogLogging constructs a []client.Logging from the given Terraform state.
func (olr *OrgLoggingResource) buildSyslogLogging(state []models.SyslogLoggingModel) *[]client.Logging {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Prepare the output slice
	output := []client.Logging{}

	// Iterate over each block and construct an output item
	for _, block := range state {
		// Construct the item
		item := client.SyslogLogging{
			Host:     BuildString(block.Host),
			Port:     BuildInt(block.Port),
			Mode:     BuildString(block.Mode),
			Format:   BuildString(block.Format),
			Severity: BuildInt(block.Severity),
		}

		// Construct logging
		logging := client.Logging{
			Syslog: &item,
		}

		// Add the item to the output slice
		output = append(output, logging)
	}

	// Return a pointer to the output
	return &output
}

// Flatteners //

// flattenS3Logging transforms *[]client.S3Logging into a []models.S3LoggingModel.
func (olr *OrgLoggingResource) flattenS3Logging(input *[]client.S3Logging) []models.S3LoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.S3LoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.S3LoggingModel{
			Bucket:      types.StringPointerValue(item.Bucket),
			Region:      types.StringPointerValue(item.Region),
			Prefix:      types.StringPointerValue(item.Prefix),
			Credentials: types.StringPointerValue(item.Credentials),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenCoralogixLogging transforms *[]client.CoralogixLogging into a []models.CoralogixLoggingModel.
func (olr *OrgLoggingResource) flattenCoralogixLogging(input *[]client.CoralogixLogging) []models.CoralogixLoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.CoralogixLoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.CoralogixLoggingModel{
			Cluster:     types.StringPointerValue(item.Cluster),
			Credentials: types.StringPointerValue(item.Credentials),
			App:         types.StringPointerValue(item.App),
			Subsystem:   types.StringPointerValue(item.Subsystem),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenDatadogLogging transforms *[]client.DatadogLogging into a []models.DatadogLoggingModel.
func (olr *OrgLoggingResource) flattenDatadogLogging(input *[]client.DatadogLogging) []models.DatadogLoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.DatadogLoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.DatadogLoggingModel{
			Host:        types.StringPointerValue(item.Host),
			Credentials: types.StringPointerValue(item.Credentials),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenLogzioLogging transforms *[]client.LogzioLogging into a []models.LogzioLoggingModel.
func (olr *OrgLoggingResource) flattenLogzioLogging(input *[]client.LogzioLogging) []models.LogzioLoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.LogzioLoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.LogzioLoggingModel{
			ListenerHost: types.StringPointerValue(item.ListenerHost),
			Credentials:  types.StringPointerValue(item.Credentials),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenElasticLogging transforms *[]client.ElasticLogging into a []models.ElasticLoggingModel.
func (olr *OrgLoggingResource) flattenElasticLogging(input *[]client.ElasticLogging) []models.ElasticLoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.ElasticLoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.ElasticLoggingModel{
			AWS:          olr.flattenElasticLoggingAws(item.AWS),
			ElasticCloud: olr.flattenElasticLoggingElasticCloud(item.ElasticCloud),
			Generic:      olr.flattenElasticLoggingGeneric(item.Generic),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenElasticLoggingAws transforms *client.AWSLogging into a []models.ElasticLoggingAwsModel.
func (olr *OrgLoggingResource) flattenElasticLoggingAws(input *client.AWSLogging) []models.ElasticLoggingAwsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ElasticLoggingAwsModel{
		Host:        types.StringPointerValue(input.Host),
		Port:        FlattenInt(input.Port),
		Index:       types.StringPointerValue(input.Index),
		Type:        types.StringPointerValue(input.Type),
		Credentials: types.StringPointerValue(input.Credentials),
		Region:      types.StringPointerValue(input.Region),
	}

	// Return a slice containing the single block
	return []models.ElasticLoggingAwsModel{block}
}

// flattenElasticLoggingElasticCloud transforms *client.ElasticCloudLogging into a []models.ElasticLoggingElasticCloudModel.
func (olr *OrgLoggingResource) flattenElasticLoggingElasticCloud(input *client.ElasticCloudLogging) []models.ElasticLoggingElasticCloudModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ElasticLoggingElasticCloudModel{
		Index:       types.StringPointerValue(input.Index),
		Type:        types.StringPointerValue(input.Type),
		Credentials: types.StringPointerValue(input.Credentials),
		CloudID:     types.StringPointerValue(input.CloudID),
	}

	// Return a slice containing the single block
	return []models.ElasticLoggingElasticCloudModel{block}
}

// flattenElasticLoggingGeneric transforms *client.GenericLogging into a []models.ElasticLoggingGenericModel.
func (olr *OrgLoggingResource) flattenElasticLoggingGeneric(input *client.GenericLogging) []models.ElasticLoggingGenericModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.ElasticLoggingGenericModel{
		Host:        types.StringPointerValue(input.Host),
		Port:        FlattenInt(input.Port),
		Path:        types.StringPointerValue(input.Path),
		Index:       types.StringPointerValue(input.Index),
		Type:        types.StringPointerValue(input.Type),
		Credentials: types.StringPointerValue(input.Credentials),
	}

	// Return a slice containing the single block
	return []models.ElasticLoggingGenericModel{block}
}

// flattenCloudWatchLogging transforms *[]client.CloudWatchLogging into a []models.CloudWatchModel.
func (olr *OrgLoggingResource) flattenCloudWatchLogging(input *[]client.CloudWatchLogging) []models.CloudWatchModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.CloudWatchModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.CloudWatchModel{
			Region:        types.StringPointerValue(item.Region),
			Credentials:   types.StringPointerValue(item.Credentials),
			RetentionDays: FlattenInt(item.RetentionDays),
			GroupName:     types.StringPointerValue(item.GroupName),
			StreamName:    types.StringPointerValue(item.StreamName),
			ExtractFields: FlattenMapString(item.ExtractFields),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenFluentdLogging transforms *[]client.FluentdLogging into a []models.FluentdLoggingModel.
func (olr *OrgLoggingResource) flattenFluentdLogging(input *[]client.FluentdLogging) []models.FluentdLoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.FluentdLoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.FluentdLoggingModel{
			Host: types.StringPointerValue(item.Host),
			Port: FlattenInt(item.Port),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenStackdriverLogging transforms *[]client.StackdriverLogging into a []models.StackdriverLoggingModel.
func (olr *OrgLoggingResource) flattenStackdriverLogging(input *[]client.StackdriverLogging) []models.StackdriverLoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.StackdriverLoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.StackdriverLoggingModel{
			Credentials: types.StringPointerValue(item.Credentials),
			Location:    types.StringPointerValue(item.Location),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

// flattenSyslogLogging transforms *[]client.SyslogLogging into a []models.SyslogLoggingModel.
func (olr *OrgLoggingResource) flattenSyslogLogging(input *[]client.SyslogLogging) []models.SyslogLoggingModel {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return nil
	}

	// Define the blocks slice
	var blocks []models.SyslogLoggingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.SyslogLoggingModel{
			Host:     types.StringPointerValue(item.Host),
			Port:     FlattenInt(item.Port),
			Mode:     types.StringPointerValue(item.Mode),
			Format:   types.StringPointerValue(item.Format),
			Severity: FlattenInt(item.Severity),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully accumulated blocks
	return blocks
}

/*** Helpers ***/

// ValidateLogging ensures the logging configurations meet size requirements.
func (olr *OrgLoggingResource) validateLoggings(diags *diag.Diagnostics, loggings []client.Logging) {
	// Ensure no more than 4 providers
	if len(loggings) > 4 {
		diags.AddError("Invalid Logging Size", "max of 4 external logging providers allowed")
		return
	}

	// Error if no providers defined
	if len(loggings) == 0 {
		diags.AddError("Empty Logging", "at least one external logging providers must be defined")
		return
	}
}
