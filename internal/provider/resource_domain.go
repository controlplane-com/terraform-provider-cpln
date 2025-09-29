package cpln

import (
	"context"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/domain"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &DomainResource{}
	_ resource.ResourceWithImportState = &DomainResource{}
	_ resource.ResourceWithModifyPlan  = &DomainResource{}
)

/*** Resource Model ***/

// DomainResourceModel holds the Terraform state for the resource.
type DomainResourceModel struct {
	EntityBaseModel
	Spec   types.List `tfsdk:"spec"`
	Status types.List `tfsdk:"status"`
}

/*** Resource Configuration ***/

// DomainResource is the resource implementation.
type DomainResource struct {
	EntityBase
	Operations EntityOperations[DomainResourceModel, client.Domain]
}

// NewDomainResource returns a new instance of the resource implementation.
func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// Configure configures the resource before use.
func (dr *DomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	dr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	dr.Operations = NewEntityOperations(dr.client, &DomainResourceOperator{})
}

// ModifyPlan handles plan modifications.
func (r *DomainResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If no existing state or plan provided, skip further processing
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Declare domain resource models for state and plan
	var st, pl DomainResourceModel

	// Populate models from stored state
	resp.Diagnostics.Append(req.State.Get(ctx, &st)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &pl)...)

	// Exit if retrieving state or plan resulted in error
	if resp.Diagnostics.HasError() {
		return
	}

	// Skip if either current or planned spec is null or unknown
	if st.Spec.IsNull() || st.Spec.IsUnknown() || pl.Spec.IsNull() || pl.Spec.IsUnknown() {
		return
	}

	// Skip warning if spec unchanged between state and plan
	if st.Spec.Equal(pl.Spec) {
		return
	}

	// Warn about temporary outage caused by domain changes
	resp.Diagnostics.AddWarning(
		"Updating domain will cause a temporary outage",
		"Changing the domain triggers DNS/TLS (and possibly cert) updates. Expect brief downtime until propagation completes.",
	)
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (dr *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (dr *DomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_domain"
}

// Schema defines the schema for the resource.
func (dr *DomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(dr.EntityBaseAttributes("Domain"), map[string]schema.Attribute{
			"status": schema.ListNestedAttribute{
				Description: "Domain status.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"status": schema.StringAttribute{
							Description: "Status of Domain. Possible values: `initializing`, `ready`, `pendingDnsConfig`, `pendingCertificate`, `usedByGvc`.",
							Computed:    true,
						},
						"warning": schema.StringAttribute{
							Description: "Warning message.",
							Computed:    true,
						},
						"fingerprint": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"endpoints": schema.ListNestedAttribute{
							Description: "List of configured domain endpoints.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"url": schema.StringAttribute{
										Description: "URL of endpoint.",
										Computed:    true,
									},
									"workload_link": schema.StringAttribute{
										Description: "Full link to associated workload.",
										Computed:    true,
									},
								},
							},
						},
						"locations": schema.ListNestedAttribute{
							Description: "Contains the cloud provider name, region, and certificate status.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the location.",
										Computed:    true,
									},
									"certificate_status": schema.StringAttribute{
										Description: "The current validity or status of the SSL/TLS certificate.",
										Computed:    true,
									},
								},
							},
						},
						"dns_config": schema.ListNestedAttribute{
							Description: "List of required DNS record entries.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Description: "The DNS record type specifies the type of data the DNS record contains. Valid values: `CNAME`, `NS`, `TXT`.",
										Computed:    true,
									},
									"ttl": schema.Int32Attribute{
										Description: "Time to live (TTL) is a value that signifies how long (in seconds) a DNS record should be cached by a resolver or a browser before a new request should be sent to refresh the data. Lower TTL values mean records are updated more frequently, which is beneficial for dynamic DNS configurations or during DNS migrations. Higher TTL values reduce the load on DNS servers and improve the speed of name resolution for end users by relying on cached data.",
										Computed:    true,
									},
									"host": schema.StringAttribute{
										Description: "The host in DNS terminology refers to the domain or subdomain that the DNS record is associated with. It's essentially the name that is being queried or managed. For example, in a DNS record for `www.example.com`, `www` is a host in the domain `example.com`.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "The value of a DNS record contains the data the record is meant to convey, based on the type of the record.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"spec": schema.ListNestedBlock{
				Description: "Domain specification.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"dns_mode": schema.StringAttribute{
							Description: "In `cname` dnsMode, Control Plane will configure workloads to accept traffic for the domain but will not manage DNS records for the domain. End users must configure CNAME records in their own DNS pointed to the canonical workload endpoint. Currently `cname` dnsMode requires that a TLS server certificate be configured when subdomain based routing is used. In `ns` dnsMode, Control Plane will manage the subdomains and create all necessary DNS records. End users configure NS records to forward DNS requests to the Control Plane managed DNS servers. Valid values: `cname`, `ns`. Default: `cname`.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("cname"),
						},
						"gvc_link": schema.StringAttribute{
							Description: "This value is set to a target GVC (using a full link) for use by subdomain based routing. Each workload in the GVC will receive a subdomain in the form ${workload.name}.${domain.name}. **Do not include if path based routing is used.**",
							Optional:    true,
						},
						"cert_challenge_type": schema.StringAttribute{
							Description: "Defines the method used to prove domain ownership for certificate issuance.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("http01", "dns01"),
							},
						},
						"workload_link": schema.StringAttribute{
							Description: "Creates a unique subdomain for each replica of a stateful workload, enabling direct access to individual instances.",
							Optional:    true,
						},
						"accept_all_hosts": schema.BoolAttribute{
							Description: "Allows domain to accept wildcards. The associated GVC must have dedicated load balancing enabled.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"accept_all_subdomains": schema.BoolAttribute{
							Description: "Accept all subdomains will accept any host that is a sub domain of the domain so *.$DOMAIN",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
					Blocks: map[string]schema.Block{
						"ports": schema.SetNestedBlock{
							Description: "Domain port specifications.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"number": schema.Int32Attribute{
										Description: "Sets or overrides headers to all http requests for this route.",
										Optional:    true,
										Computed:    true,
										Default:     int32default.StaticInt32(443),
									},
									"protocol": schema.StringAttribute{
										Description: "Allowed protocol. Valid values: `http`, `http2`, `tcp`. Default: `http2`.",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("http2"),
									},
								},
								Blocks: map[string]schema.Block{
									"cors": schema.ListNestedBlock{
										Description: "A security feature implemented by web browsers to allow resources on a web page to be requested from another domain outside the domain from which the resource originated.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"allow_methods": schema.SetAttribute{
													Description: "Specifies the HTTP methods (such as `GET`, `POST`, `PUT`, `DELETE`, etc.) that are allowed for a cross-origin request to a specific resource.",
													ElementType: types.StringType,
													Optional:    true,
												},
												"allow_headers": schema.SetAttribute{
													Description: "Specifies the custom HTTP headers that are allowed in a cross-origin request to a specific resource.",
													ElementType: types.StringType,
													Optional:    true,
												},
												"expose_headers": schema.SetAttribute{
													Description: "The HTTP headers that a server allows to be exposed to the client in response to a cross-origin request. These headers provide additional information about the server's capabilities or requirements, aiding in proper handling of the request by the client's browser or application.",
													ElementType: types.StringType,
													Optional:    true,
												},
												"max_age": schema.StringAttribute{
													Description: "Maximum amount of time that a preflight request result can be cached by the client browser. Input is expected as a duration string (i.e, 24h, 20m, etc.).",
													Optional:    true,
													Computed:    true,
													Default:     stringdefault.StaticString("24h"),
												},
												"allow_credentials": schema.BoolAttribute{
													Description: "Determines whether the client-side code (typically running in a web browser) is allowed to include credentials (such as cookies, HTTP authentication, or client-side SSL certificates) in cross-origin requests.",
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
												},
											},
											Blocks: map[string]schema.Block{
												"allow_origins": schema.SetNestedBlock{
													Description: "Determines which origins are allowed to access a particular resource on a server from a web browser.",
													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"exact": schema.StringAttribute{
																Description: "Value of allowed origin.",
																Optional:    true,
																Validators: []validator.String{
																	stringvalidator.ExactlyOneOf(
																		path.MatchRelative().AtParent().AtName("exact"),
																		path.MatchRelative().AtParent().AtName("regex"),
																	),
																},
															},
															"regex": schema.StringAttribute{
																Description: "",
																Optional:    true,
																Validators: []validator.String{
																	stringvalidator.ExactlyOneOf(
																		path.MatchRelative().AtParent().AtName("exact"),
																		path.MatchRelative().AtParent().AtName("regex"),
																	),
																},
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
									"tls": schema.ListNestedBlock{
										Description: "Used for TLS connections for this Domain. End users are responsible for certificate updates.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"min_protocol_version": schema.StringAttribute{
													Description: "Minimum TLS version to accept. Minimum is `1.0`. Default: `1.2`.",
													Optional:    true,
													Computed:    true,
													Default:     stringdefault.StaticString("TLSV1_2"),
												},
												"cipher_suites": schema.SetAttribute{
													Description: "Allowed cipher suites. Refer to the [Domain Reference](https://docs.controlplane.com/reference/domain#cipher-suites) for details.",
													ElementType: types.StringType,
													Optional:    true,
													Computed:    true,
													Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
														types.StringValue("AES128-GCM-SHA256"),
														types.StringValue("AES256-GCM-SHA384"),
														types.StringValue("ECDHE-ECDSA-AES128-GCM-SHA256"),
														types.StringValue("ECDHE-ECDSA-AES256-GCM-SHA384"),
														types.StringValue("ECDHE-ECDSA-CHACHA20-POLY1305"),
														types.StringValue("ECDHE-RSA-AES128-GCM-SHA256"),
														types.StringValue("ECDHE-RSA-AES256-GCM-SHA384"),
														types.StringValue("ECDHE-RSA-CHACHA20-POLY1305"),
													})),
												},
											},
											Blocks: map[string]schema.Block{
												"client_certificate": schema.ListNestedBlock{
													Description:  "The certificate authority PEM, stored as a TLS Secret, used to verify the authority of the client certificate. The only verification performed checks that the CN of the PEM matches the Domain (i.e., CN=*.DOMAIN).",
													NestedObject: dr.CertificateSchema("The secret will include a client certificate authority cert in PEM format used to verify requests which include client certificates. The key subject must match the domain and the key usage properties must be configured for client certificate authorization. The secret type must be keypair."),
													Validators: []validator.List{
														listvalidator.SizeAtMost(1),
													},
												},
												"server_certificate": schema.ListNestedBlock{
													Description:  "Configure an optional custom server certificate for the domain. When the port number is 443 and this is not supplied, a certificate is provisioned automatically.",
													NestedObject: dr.CertificateSchema("When provided, this is used as the server certificate authority. The secret type must be keypair and the content must be PEM encoded."),
													Validators: []validator.List{
														listvalidator.SizeAtMost(1),
													},
												},
											},
										},
										Validators: []validator.List{
											listvalidator.IsRequired(),
											listvalidator.SizeAtMost(1),
										},
									},
								},
							},
							Validators: []validator.Set{
								setvalidator.IsRequired(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

// Create creates the resource.
func (dr *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, dr.Operations)
}

// Read fetches the current state of the resource.
func (dr *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, dr.Operations)
}

// Update modifies the resource.
func (dr *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, dr.Operations)
}

// Delete removes the resource.
func (dr *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, dr.Operations)
}

/*** Schemas ***/

// CertificateSchema creates a nested block schema using the provided description.
func (dr *DomainResource) CertificateSchema(description string) schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"secret_link": schema.StringAttribute{
				Description: description,
				Optional:    true,
			},
		},
	}
}

/*** Resource Operator ***/

// DomainResourceOperator is the operator for managing the state.
type DomainResourceOperator struct {
	EntityOperator[DomainResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (dro *DomainResourceOperator) NewAPIRequest(isUpdate bool) client.Domain {
	// Initialize a new request payload
	requestPayload := client.Domain{}

	// Populate Base fields from state
	dro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Build the domain spec struct
	var spec *client.DomainSpec = dro.buildSpec(dro.Plan.Spec)

	// Set specific attributes
	if isUpdate {
		requestPayload.SpecReplace = spec

		// Fetch the domain to inspect its routes
		domain, _, err := dro.InvokeRead(dro.Plan.Name.ValueString())

		// Handle error
		if err != nil {
			dro.Diags.AddError("Unable to Fetch Domain", fmt.Sprintf("Unable to fetch domain during update, details: %v", err))
			return requestPayload
		}

		// Inspect the routes of the domain ports only if ports is set
		if spec != nil && spec.Ports != nil && domain.Spec != nil && domain.Spec.Ports != nil {
			// Initialize a map to hold copies of routes keyed by port number
			routeMap := make(map[int][]client.DomainRoute, len(*domain.Spec.Ports))

			// Iterate over each port in the domain spec
			for _, port := range *domain.Spec.Ports {
				// Skip ports without any routes or without port number
				if port.Number == nil || port.Routes == nil || len(*port.Routes) == 0 {
					continue
				}

				// Dereference the Routes pointer to get original slice
				source := *port.Routes

				// Create a new slice to hold deep-copied routes
				destination := make([]client.DomainRoute, len(source))

				// Copy routes to avoid aliasing original slice
				copy(destination, source)

				// Store the copied routes in the map under the port number
				routeMap[*port.Number] = destination
			}

			// Iterate over ports in the new spec
			for i := range *spec.Ports {
				// Get pointer to the current port in the new spec
				up := &(*spec.Ports)[i]

				// Set routes to this port if routes exist in the map for this port number
				if routes, ok := routeMap[*up.Number]; ok {
					// Assign the copied routes back to the updated port
					up.Routes = &routes
				}
			}
		}
	} else {
		requestPayload.Spec = spec
	}

	// Return the request payload object
	return requestPayload
}

// MapResponseToState creates a state model from response payload.
func (dro *DomainResourceOperator) MapResponseToState(domain *client.Domain, isCreate bool) DomainResourceModel {
	// Initialize empty state model
	state := DomainResourceModel{}

	// Populate common fields from base resource data
	state.From(domain.Base)

	// In case the self link is empty, construct one
	state.SelfLink = types.StringValue(GetSelfLink(dro.Client.Org, "domain", *domain.Name))

	// Set specific attributes
	state.Spec = dro.flattenSpec(domain.Spec)
	state.Status = dro.flattenStatus(domain.Status)

	// Return the built state
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (dro *DomainResourceOperator) InvokeCreate(req client.Domain) (*client.Domain, int, error) {
	return dro.Client.CreateDomain(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (dro *DomainResourceOperator) InvokeRead(name string) (*client.Domain, int, error) {
	return dro.Client.GetDomain(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (dro *DomainResourceOperator) InvokeUpdate(req client.Domain) (*client.Domain, int, error) {
	return dro.Client.UpdateDomain(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (dro *DomainResourceOperator) InvokeDelete(name string) error {
	return dro.Client.DeleteDomain(name)
}

// Builders //

// buildSpec constructs a DomainSpec struct from the given Terraform state.
func (dro *DomainResourceOperator) buildSpec(state types.List) *client.DomainSpec {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SpecModel](dro.Ctx, dro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.DomainSpec{
		DnsMode:             BuildString(block.DnsMode),
		GvcLink:             BuildString(block.GvcLink),
		CertChallengeType:   BuildString(block.CertChallengeType),
		WorkloadLink:        BuildString(block.WorkloadLink),
		AcceptAllHosts:      BuildBool(block.AcceptAllHosts),
		AcceptAllSubdomains: BuildBool(block.AcceptAllSubdomains),
		Ports:               dro.buildSpecPorts(block.Ports),
	}
}

// buildSpecPorts constructs a []client.DomainSpecPort slice from the given Terraform state.
func (dro *DomainResourceOperator) buildSpecPorts(state types.Set) *[]client.DomainSpecPort {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.SpecPortsModel](dro.Ctx, dro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Declare the result slice
	result := []client.DomainSpecPort{}

	// Iterate over each block and construct a result item
	for _, block := range blocks {
		// Construct the item
		item := client.DomainSpecPort{
			Number:   BuildInt(block.Number),
			Protocol: BuildString(block.Protocol),
			Cors:     dro.buildSpecPortCors(block.Cors),
			TLS:      dro.buildSpecPortTls(block.TLS),
		}

		// Add the item to the result slice
		result = append(result, item)
	}

	// Return the result
	return &result
}

// buildSpecPortCors constructs a DomainCors struct from the given Terraform state.
func (dro *DomainResourceOperator) buildSpecPortCors(state types.List) *client.DomainCors {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SpecPortsCorsModel](dro.Ctx, dro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.DomainCors{
		AllowOrigins:     dro.buildSpecPortCorsAllowOrigins(block.AllowOrigins),
		AllowMethods:     BuildSetString(dro.Ctx, dro.Diags, block.AllowMethods),
		AllowHeaders:     BuildSetString(dro.Ctx, dro.Diags, block.AllowHeaders),
		ExposeHeaders:    BuildSetString(dro.Ctx, dro.Diags, block.ExposeHeaders),
		MaxAge:           BuildString(block.MaxAge),
		AllowCredentials: BuildBool(block.AllowCredentials),
	}
}

// buildSpecPortCorsAllowOrigins constructs a []client.DomainAllowOrigin slice from the given Terraform state.
func (dro *DomainResourceOperator) buildSpecPortCorsAllowOrigins(state types.Set) *[]client.DomainAllowOrigin {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.SpecPortsCorsAllowOriginsModel](dro.Ctx, dro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Declare the result slice
	result := []client.DomainAllowOrigin{}

	// Iterate over each block and construct a result item
	for _, block := range blocks {
		// Construct the item
		item := client.DomainAllowOrigin{
			Exact: BuildString(block.Exact),
			Regex: BuildString(block.Regex),
		}

		// Add the item to the result slice
		result = append(result, item)
	}

	// Return the result
	return &result
}

// buildSpecPortTls constructs a DomainTLS struct from the given Terraform state.
func (dro *DomainResourceOperator) buildSpecPortTls(state types.List) *client.DomainTLS {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SpecPortsTlsModel](dro.Ctx, dro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.DomainTLS{
		MinProtocolVersion: BuildString(block.MinProtocolVersion),
		CipherSuites:       BuildSetString(dro.Ctx, dro.Diags, block.CipherSuites),
		ClientCertificate:  dro.buildSpecPortTlsCertificate(block.ClientCertificate),
		ServerCertificate:  dro.buildSpecPortTlsCertificate(block.ServerCertificate),
	}
}

// buildSpecPortTlsCertificate constructs a DomainCertificate struct from the given Terraform state.
func (dro *DomainResourceOperator) buildSpecPortTlsCertificate(state types.List) *client.DomainCertificate {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SpecPortsTlsCertificateModel](dro.Ctx, dro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.DomainCertificate{
		SecretLink: BuildString(block.SecretLink),
	}
}

// Flatteners //

// flattenSpec transforms client.DomainSpec into a Terraform types.List.
func (dro *DomainResourceOperator) flattenSpec(input *client.DomainSpec) types.List {
	// Get attribute types
	elementType := models.SpecModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SpecModel{
		DnsMode:             types.StringPointerValue(input.DnsMode),
		GvcLink:             types.StringPointerValue(input.GvcLink),
		CertChallengeType:   types.StringPointerValue(input.CertChallengeType),
		WorkloadLink:        types.StringPointerValue(input.WorkloadLink),
		AcceptAllHosts:      types.BoolPointerValue(input.AcceptAllHosts),
		AcceptAllSubdomains: types.BoolPointerValue(input.AcceptAllSubdomains),
		Ports:               dro.flattenSpecPorts(input.Ports),
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, []models.SpecModel{block})
}

// flattenSpecPorts transforms []client.DomainSpecPort into a Terraform types.Set.
func (dro *DomainResourceOperator) flattenSpecPorts(input *[]client.DomainSpecPort) types.Set {
	// Get attribute types
	elementType := models.SpecPortsModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	blocks := []models.SpecPortsModel{}

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.SpecPortsModel{
			Number:   FlattenInt(item.Number),
			Protocol: types.StringPointerValue(item.Protocol),
			Cors:     dro.flattenSpecPortCors(item.Cors),
			TLS:      dro.flattenSpecPortTls(item.TLS),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(dro.Ctx, dro.Diags, blocks)
}

// flattenSpecPortCors transforms client.DomainCors into a Terraform types.List.
func (dro *DomainResourceOperator) flattenSpecPortCors(input *client.DomainCors) types.List {
	// Get attribute types
	elementType := models.SpecPortsCorsModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SpecPortsCorsModel{
		AllowOrigins:     dro.flattenSpecPortCorsAllowOrigins(input.AllowOrigins),
		AllowMethods:     FlattenSetString(input.AllowMethods),
		AllowHeaders:     FlattenSetString(input.AllowHeaders),
		ExposeHeaders:    FlattenSetString(input.ExposeHeaders),
		MaxAge:           types.StringPointerValue(input.MaxAge),
		AllowCredentials: types.BoolPointerValue(input.AllowCredentials),
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, []models.SpecPortsCorsModel{block})
}

// flattenSpecPortCorsAllowOrigins transforms []client.DomainAllowOrigin into a Terraform types.Set.
func (dro *DomainResourceOperator) flattenSpecPortCorsAllowOrigins(input *[]client.DomainAllowOrigin) types.Set {
	// Get attribute types
	elementType := models.SpecPortsCorsAllowOriginsModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	blocks := []models.SpecPortsCorsAllowOriginsModel{}

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.SpecPortsCorsAllowOriginsModel{
			Exact: types.StringPointerValue(item.Exact),
			Regex: types.StringPointerValue(item.Regex),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(dro.Ctx, dro.Diags, blocks)
}

// flattenSpecPortTls transforms client.DomainTLS into a Terraform types.List.
func (dro *DomainResourceOperator) flattenSpecPortTls(input *client.DomainTLS) types.List {
	// Get attribute types
	elementType := models.SpecPortsTlsModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SpecPortsTlsModel{
		MinProtocolVersion: types.StringPointerValue(input.MinProtocolVersion),
		CipherSuites:       FlattenSetString(input.CipherSuites),
		ClientCertificate:  dro.flattenSpecPortTlsCertificate(input.ClientCertificate),
		ServerCertificate:  dro.flattenSpecPortTlsCertificate(input.ServerCertificate),
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, []models.SpecPortsTlsModel{block})
}

// flattenSpecPortTlsCertificate transforms client.DomainCertificate into a Terraform types.List.
func (dro *DomainResourceOperator) flattenSpecPortTlsCertificate(input *client.DomainCertificate) types.List {
	// Get attribute types
	elementType := models.SpecPortsTlsCertificateModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SpecPortsTlsCertificateModel{
		SecretLink: types.StringPointerValue(input.SecretLink),
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, []models.SpecPortsTlsCertificateModel{block})
}

// flattenStatus transforms client.DomainStatus into a Terraform types.List.
func (dro *DomainResourceOperator) flattenStatus(input *client.DomainStatus) types.List {
	// Get attribute types
	elementType := models.StatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusModel{
		Endpoints:   dro.flattenStatusEndpoints(input.Endpoints),
		Status:      types.StringPointerValue(input.Status),
		Warning:     types.StringPointerValue(input.Warning),
		Locations:   dro.flattenStatusLocations(input.Locations),
		Fingerprint: types.StringPointerValue(input.Fingerprint),
		DnsConfig:   dro.flattenStatusDnsConfig(input.DnsConfig),
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, []models.StatusModel{block})
}

// flattenStatusEndpoints transforms []client.DomainEndpoint into a Terraform types.List.
func (dro *DomainResourceOperator) flattenStatusEndpoints(input *[]client.DomainStatusEndpoint) types.List {
	// Get attribute types
	elementType := models.StatusEndpointModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	blocks := []models.StatusEndpointModel{}

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.StatusEndpointModel{
			URL:          types.StringPointerValue(item.URL),
			WorkloadLink: types.StringPointerValue(item.WorkloadLink),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, blocks)
}

// flattenStatusLocations transforms []client.DomainStatusLocation into a Terraform types.List.
func (dro *DomainResourceOperator) flattenStatusLocations(input *[]client.DomainStatusLocation) types.List {
	// Get attribute types
	elementType := models.StatusLocationModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	blocks := []models.StatusLocationModel{}

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.StatusLocationModel{
			Name:              types.StringPointerValue(item.Name),
			CertificateStatus: types.StringPointerValue(item.CertificateStatus),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, blocks)
}

// flattenStatusDnsConfig transforms []client.DnsConfigRecord into a Terraform types.List.
func (dro *DomainResourceOperator) flattenStatusDnsConfig(input *[]client.DomainStatusDnsConfigRecord) types.List {
	// Get attribute types
	elementType := models.StatusDnsConfigModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	blocks := []models.StatusDnsConfigModel{}

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.StatusDnsConfigModel{
			Type:  types.StringPointerValue(item.Type),
			TTL:   FlattenInt(item.TTL),
			Host:  types.StringPointerValue(item.Host),
			Value: types.StringPointerValue(item.Value),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(dro.Ctx, dro.Diags, blocks)
}
