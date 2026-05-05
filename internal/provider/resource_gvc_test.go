package cpln

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

/*** Acceptance Test ***/

// TestAccControlPlaneGvc_basic performs an acceptance test for the resource.
func TestAccControlPlaneGvc_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewGvcResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "GVC") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// GvcResourceTest defines the necessary functionality to test the resource.
type GvcResourceTest struct {
	Steps []resource.TestStep
}

// NewGvcResourceTest creates a GvcResourceTest with initialized test cases.
func NewGvcResourceTest() GvcResourceTest {
	// Create a resource test instance
	resourceTest := GvcResourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewOrgNamingScenario()...)
	steps = append(steps, resourceTest.NewDefaultNamingScenario()...)
	steps = append(steps, resourceTest.NewLocationOptionsLifecycleScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (grt *GvcResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_gvc resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_gvc" {
			continue
		}

		// Retrieve the name for the current resource
		gvcName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of GVC with name: %s", gvcName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		gvc, code, err := TestProvider.client.GetGvc(gvcName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if GVC %s exists: %w", gvcName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if gvc != nil {
			return fmt.Errorf("CheckDestroy failed: GVC %s still exists in the system", *gvc.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_gvc resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewOrgNamingScenario creates a test case for a GVC with endpoint naming format set to "org" with initial and updated configurations.
func (grt *GvcResourceTest) NewOrgNamingScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("gvc-default-%s", random)
	opaqueName := fmt.Sprintf("opaque-%s", random)
	dockerName := "test-gvc-docker-pull-secret"
	resourceName := "with-org-endpoint-naming-format"

	// Create the opaque secret case
	opaqueSecretCase := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.opaque",
			Name:              opaqueName,
			Description:       opaqueName,
			DescriptionUpdate: "secret description updated",
		},
	}

	// Get secret config
	opaqueSecretConfig := opaqueSecretCase.OpaqueRequiredOnly("opaque_secret_payload")

	// Declare the endpoint naming format for this test case
	endpointNamingFormat := "org"

	// Build test steps
	initialConfig, initialStep := grt.BuildInitialTestStep(resourceName, name)
	caseUpdate1 := grt.BuildUpdate1TestStep(initialConfig.ProviderTestCase, endpointNamingFormat, dockerName, opaqueSecretConfig)
	caseUpdate2 := grt.BuildUpdate2TestStep(initialConfig.ProviderTestCase, endpointNamingFormat, dockerName, opaqueSecretCase, opaqueSecretConfig)
	caseUpdate3 := grt.BuildUpdate3TestStep(initialConfig.ProviderTestCase, endpointNamingFormat, dockerName, opaqueSecretCase, opaqueSecretConfig)
	caseUpdate4 := grt.BuildUpdate4TestStep(initialConfig.ProviderTestCase, endpointNamingFormat, dockerName, opaqueSecretCase, opaqueSecretConfig)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName: initialConfig.ResourceAddress,
			ImportState:  true,
		},
		// Update & Read
		caseUpdate1,
		caseUpdate2,
		caseUpdate3,
		caseUpdate4,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewDefaultNamingScenario creates a test case for a GVC with endpoint naming format set to "default" with initial and updated configurations.
func (grt *GvcResourceTest) NewDefaultNamingScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("gvc-default-%s", random)
	opaqueName := fmt.Sprintf("opaque-%s", random)
	dockerName := "test-gvc-docker-pull-secret"
	resourceName := "with-default-endpoint-naming-format"

	// Create the opaque secret case
	opaqueSecretCase := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.opaque",
			Name:              opaqueName,
			Description:       opaqueName,
			DescriptionUpdate: "secret description updated",
		},
	}

	// Get secret config
	opaqueSecretConfig := opaqueSecretCase.OpaqueRequiredOnly("opaque_secret_payload")

	// Build test steps
	initialConfig, initialStep := grt.BuildInitialTestStepWithEndpointNamingFormat(resourceName, name, "default")
	caseUpdate1 := grt.BuildUpdate1TestStep(initialConfig.ProviderTestCase, "default", dockerName, opaqueSecretConfig)
	caseUpdate2 := grt.BuildUpdate2TestStep(initialConfig.ProviderTestCase, "legacy", dockerName, opaqueSecretCase, opaqueSecretConfig)
	caseUpdate3 := grt.BuildUpdate3TestStep(initialConfig.ProviderTestCase, "legacy", dockerName, opaqueSecretCase, opaqueSecretConfig)
	caseUpdate4 := grt.BuildUpdate4TestStep(initialConfig.ProviderTestCase, "legacy", dockerName, opaqueSecretCase, opaqueSecretConfig)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName: initialConfig.ResourceAddress,
			ImportState:  true,
		},
		// Update & Read
		caseUpdate1,
		caseUpdate2,
		caseUpdate3,
		caseUpdate4,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewLocationOptionsLifecycleScenario walks the location_options SetNestedBlock through its full state lifecycle to catch state-drift bugs at empty/null and cardinality transitions.
func (grt *GvcResourceTest) NewLocationOptionsLifecycleScenario() []resource.TestStep {
	// Generate a unique name for the resource
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("gvc-loc-opts-%s", random)
	resourceName := "with-location-options-lifecycle"

	// Create a shared test case for all lifecycle steps
	c := GvcResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "gvc",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_gvc.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: name,
		},
	}

	// Build the lifecycle test steps
	absentStep := grt.BuildLocationOptionsAbsentStep(c)
	requiredOnlyStep := grt.BuildLocationOptionsRequiredOnlyStep(c)
	multiAllStep := grt.BuildLocationOptionsMultiAllStep(c)
	shrunkStep := grt.BuildLocationOptionsShrunkStep(c)

	// Walk the lifecycle: absent → 1 min → 2 full → 1 full (shrink) → 2 full (regrow) → absent
	return []resource.TestStep{
		absentStep,
		requiredOnlyStep,
		multiAllStep,
		shrunkStep,
		multiAllStep,
		absentStep,
	}
}

// Test Cases //

// BuildInitialTestStep returns a default initial test step and its associated test case for the GVC resource.
func (grt *GvcResourceTest) BuildInitialTestStep(resourceName string, name string) (GvcResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := GvcResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "gvc",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_gvc.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "gvc default description updated",
		},
		EndpointNamingFormat: "org",
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: grt.GvcRequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "endpoint_naming_format", c.EndpointNamingFormat),
		),
	}
}

// BuildInitialTestStepWithEndpointNamingFormat returns a default initial test step and its associated test case for the GVC resource.
func (grt *GvcResourceTest) BuildInitialTestStepWithEndpointNamingFormat(resourceName string, name string, endpointNamingFormat string) (GvcResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := GvcResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "gvc",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_gvc.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "gvc default description updated",
		},
		EndpointNamingFormat: endpointNamingFormat,
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: grt.GvcRequiredOnlyWithEndpointNamingFormat(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "endpoint_naming_format", c.EndpointNamingFormat),
		),
	}
}

// BuildUpdate1TestStep constructs the first update test step with optional tracing, load balancer, and Envoy settings for the GVC resource.
func (grt *GvcResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase, endpointNamingFormat string, dockerName string, opaqueSecretConfig string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := GvcResourceTestCase{
		ProviderTestCase:     initialCase,
		EndpointNamingFormat: endpointNamingFormat,
		Locations:            []string{"aws-eu-central-1"},
		PullSecrets:          []string{dockerName},
		Env: map[string]interface{}{
			"env-name-01": "env-value-01",
			"env-name-02": "env-value-02",
		},
		Tracing: client.Tracing{
			Sampling: Float64Pointer(55.55),
			Provider: &client.TracingProvider{
				Lightstep: &client.TracingProviderLightstep{
					Endpoint: StringPointer("test.cpln.local:8080"),
				},
			},
		},
		LoadBalancer: client.GvcLoadBalancer{
			TrustedProxies: IntPointer(0),
		},
		Envoy: `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
	}

	// Initialize the tracing block
	lightstepTracingRequiredOnlyBlock := grt.LightstepTracingRequiredOnly(c)

	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.UpdateWithMinimalOptionals(c, opaqueSecretConfig, lightstepTracingRequiredOnlyBlock),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "endpoint_naming_format", c.EndpointNamingFormat),
			c.TestCheckSetAttr("locations", c.Locations),
			c.TestCheckSetAttr("pull_secrets", c.PullSecrets),
			c.TestCheckMapAttr("env", ConvertMapToStringMap(c.Env)),
			c.TestCheckNestedBlocks("lightstep_tracing", []map[string]interface{}{
				{
					"sampling": strconv.FormatFloat(*c.Tracing.Sampling, 'f', 2, 64),
					"endpoint": *c.Tracing.Provider.Lightstep.Endpoint,
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"trusted_proxies": strconv.Itoa(*c.LoadBalancer.TrustedProxies),
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("keda", []map[string]interface{}{
				{
					"enabled": "false",
				},
			}),
		),
	}
}

// BuildUpdate2TestStep builds the second update test step including advanced load balancer, custom tracing tags, and nested redirect settings.
func (grt *GvcResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase, endpointNamingFormat string, dockerName string, opaqueSecretCase SecretResourceTestScenario, opaqueSecretConfig string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := GvcResourceTestCase{
		ProviderTestCase:     initialCase,
		EndpointNamingFormat: endpointNamingFormat,
		Locations:            []string{"aws-eu-central-1", "aws-us-west-2"},
		PullSecrets:          []string{dockerName},
		Env: map[string]interface{}{
			"env-name-01": "env-value-01",
			"env-name-02": "env-value-02",
			"env-name-03": "env-value-03",
		},
		Tracing: client.Tracing{
			Sampling: Float64Pointer(50),
			Provider: &client.TracingProvider{
				Lightstep: &client.TracingProviderLightstep{
					Endpoint:    StringPointer("test.cpln.local:80"),
					Credentials: StringPointer(opaqueSecretCase.GetSelfLink()),
				},
			},
			CustomTags: &map[string]client.TracingCustomTag{
				"key": {
					Literal: &client.TracingCustomTagValue{
						Value: StringPointer("value"),
					},
				},
			},
		},
		LoadBalancer: client.GvcLoadBalancer{
			Dedicated:      BoolPointer(false),
			TrustedProxies: IntPointer(2),
			IpSet:          StringPointer("my-ipset-01"),
			MultiZone: &client.GvcLoadBalancerMultiZone{
				Enabled: BoolPointer(true),
			},
			Redirect: &client.GvcLoadBalancerRedirect{
				Class: &client.GvcLoadBalancerRedirectClass{
					Status5XX: StringPointer("https://example.org/error/5xx"),
					Status401: StringPointer("https://your-oauth-server/oauth2/authorize?return_to=%%REQ(:path)%%&client_id=your-client-id-01"),
				},
			},
		},
		Envoy: `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"15s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
	}

	// Convert tracing custom tags to map[string]interface{}
	customTags := grt.ConvertCustomTagsToMap(*c.Tracing.CustomTags)

	// Initialize the tracing block
	lightstepTracingWithOptionalsBlock := grt.LightstepTracingWithOptionals(c, opaqueSecretCase.GetSelfLinkAttr(), customTags)

	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.UpdateWithAllOptionals(c, opaqueSecretConfig, lightstepTracingWithOptionalsBlock),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "endpoint_naming_format", c.EndpointNamingFormat),
			c.TestCheckSetAttr("locations", c.Locations),
			c.TestCheckSetAttr("pull_secrets", c.PullSecrets),
			c.TestCheckMapAttr("env", ConvertMapToStringMap(c.Env)),
			c.TestCheckNestedBlocks("lightstep_tracing", []map[string]interface{}{
				{
					"sampling":    fmt.Sprintf("%.0f", *c.Tracing.Sampling),
					"endpoint":    *c.Tracing.Provider.Lightstep.Endpoint,
					"credentials": opaqueSecretCase.GetSelfLink(),
					"custom_tags": customTags,
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"dedicated":       strconv.FormatBool(*c.LoadBalancer.Dedicated),
					"trusted_proxies": strconv.Itoa(*c.LoadBalancer.TrustedProxies),
					"ipset":           *c.LoadBalancer.IpSet,
					"multi_zone": []map[string]interface{}{
						{
							"enabled": strconv.FormatBool(*c.LoadBalancer.MultiZone.Enabled),
						},
					},
					"redirect": []map[string]interface{}{
						{
							"class": []map[string]interface{}{
								{
									"status_5xx": *c.LoadBalancer.Redirect.Class.Status5XX,
									"status_401": *c.LoadBalancer.Redirect.Class.Status401,
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("keda", []map[string]interface{}{
				{
					"enabled":       "true",
					"identity_link": fmt.Sprintf("/org/%s/gvc/%s/identity/non-existant-identity", OrgName, c.Name),
					"secrets": []string{
						fmt.Sprintf("/org/%s/secret/non-existant-secret-01", OrgName),
						fmt.Sprintf("/org/%s/secret/non-existant-secret-02", OrgName),
						fmt.Sprintf("/org/%s/secret/non-existant-secret-03", OrgName),
					},
				},
			}),
		),
	}
}

// BuildUpdate3TestStep builds the third update test step including advanced load balancer, custom tracing tags, and nested redirect settings.
func (grt *GvcResourceTest) BuildUpdate3TestStep(initialCase ProviderTestCase, endpointNamingFormat string, dockerName string, opaqueSecretCase SecretResourceTestScenario, opaqueSecretConfig string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := GvcResourceTestCase{
		ProviderTestCase:     initialCase,
		EndpointNamingFormat: endpointNamingFormat,
		Locations:            []string{"aws-eu-central-1", "aws-us-west-2"},
		PullSecrets:          []string{dockerName},
		Env: map[string]interface{}{
			"env-name-01": "env-value-01",
			"env-name-02": "env-value-02",
			"env-name-03": "env-value-03",
		},
		Tracing: client.Tracing{
			Sampling: Float64Pointer(50),
			Provider: &client.TracingProvider{
				Otel: &client.TracingProviderOtel{
					Endpoint: StringPointer("test.cpln.local:80"),
				},
			},
			CustomTags: &map[string]client.TracingCustomTag{
				"key": {
					Literal: &client.TracingCustomTagValue{
						Value: StringPointer("value"),
					},
				},
			},
		},
		LoadBalancer: client.GvcLoadBalancer{
			Dedicated:      BoolPointer(false),
			TrustedProxies: IntPointer(2),
			IpSet:          StringPointer("my-ipset-01"),
			MultiZone: &client.GvcLoadBalancerMultiZone{
				Enabled: BoolPointer(true),
			},
			Redirect: &client.GvcLoadBalancerRedirect{
				Class: &client.GvcLoadBalancerRedirectClass{
					Status5XX: StringPointer("https://example.org/error/5xx"),
					Status401: StringPointer("https://your-oauth-server/oauth2/authorize?return_to=%%REQ(:path)%%&client_id=your-client-id-01"),
				},
			},
		},
		Envoy: `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"15s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
	}

	// Convert tracing custom tags to map[string]interface{}
	customTags := grt.ConvertCustomTagsToMap(*c.Tracing.CustomTags)

	// Initialize the tracing block
	otelTracingWithOptionalsBlock := grt.OtelTracingHcl(c, customTags)

	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.UpdateWithAllOptionals(c, opaqueSecretConfig, otelTracingWithOptionalsBlock),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "endpoint_naming_format", c.EndpointNamingFormat),
			c.TestCheckSetAttr("locations", c.Locations),
			c.TestCheckSetAttr("pull_secrets", c.PullSecrets),
			c.TestCheckMapAttr("env", ConvertMapToStringMap(c.Env)),
			c.TestCheckNestedBlocks("otel_tracing", []map[string]interface{}{
				{
					"sampling":    fmt.Sprintf("%.0f", *c.Tracing.Sampling),
					"endpoint":    *c.Tracing.Provider.Otel.Endpoint,
					"custom_tags": customTags,
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"dedicated":       strconv.FormatBool(*c.LoadBalancer.Dedicated),
					"trusted_proxies": strconv.Itoa(*c.LoadBalancer.TrustedProxies),
					"ipset":           *c.LoadBalancer.IpSet,
					"multi_zone": []map[string]interface{}{
						{
							"enabled": strconv.FormatBool(*c.LoadBalancer.MultiZone.Enabled),
						},
					},
					"redirect": []map[string]interface{}{
						{
							"class": []map[string]interface{}{
								{
									"status_5xx": *c.LoadBalancer.Redirect.Class.Status5XX,
									"status_401": *c.LoadBalancer.Redirect.Class.Status401,
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("keda", []map[string]interface{}{
				{
					"enabled":       "true",
					"identity_link": fmt.Sprintf("/org/%s/gvc/%s/identity/non-existant-identity", OrgName, c.Name),
				},
			}),
		),
	}
}

// BuildUpdate4TestStep builds the fourth update test step including advanced load balancer, custom tracing tags, and nested redirect settings.
func (grt *GvcResourceTest) BuildUpdate4TestStep(initialCase ProviderTestCase, endpointNamingFormat string, dockerName string, opaqueSecretCase SecretResourceTestScenario, opaqueSecretConfig string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := GvcResourceTestCase{
		ProviderTestCase:     initialCase,
		EndpointNamingFormat: endpointNamingFormat,
		Locations:            []string{"aws-eu-central-1", "aws-us-west-2"},
		PullSecrets:          []string{dockerName},
		Env: map[string]interface{}{
			"env-name-01": "env-value-01",
			"env-name-02": "env-value-02",
			"env-name-03": "env-value-03",
		},
		Tracing: client.Tracing{
			Sampling: Float64Pointer(50),
			Provider: &client.TracingProvider{
				ControlPlane: &client.TracingProviderControlPlane{},
			},
			CustomTags: &map[string]client.TracingCustomTag{
				"key": {
					Literal: &client.TracingCustomTagValue{
						Value: StringPointer("value"),
					},
				},
			},
		},
		LoadBalancer: client.GvcLoadBalancer{
			Dedicated:      BoolPointer(false),
			TrustedProxies: IntPointer(2),
			IpSet:          StringPointer("my-ipset-01"),
			MultiZone: &client.GvcLoadBalancerMultiZone{
				Enabled: BoolPointer(true),
			},
			Redirect: &client.GvcLoadBalancerRedirect{
				Class: &client.GvcLoadBalancerRedirectClass{
					Status5XX: StringPointer("https://example.org/error/5xx"),
					Status401: StringPointer("https://your-oauth-server/oauth2/authorize?return_to=%%REQ(:path)%%&client_id=your-client-id-01"),
				},
			},
		},
		Envoy: `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"15s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
	}

	// Convert tracing custom tags to map[string]interface{}
	customTags := grt.ConvertCustomTagsToMap(*c.Tracing.CustomTags)

	// Initialize the tracing block
	cplnTracingWithOptionalsBlock := grt.ControlPlaneTracingHcl(c, customTags)

	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.UpdateWithAllOptionals(c, opaqueSecretConfig, cplnTracingWithOptionalsBlock),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "endpoint_naming_format", c.EndpointNamingFormat),
			c.TestCheckSetAttr("locations", c.Locations),
			c.TestCheckNestedBlocks("location_options", []map[string]interface{}{
				{
					"name":                 "aws-eu-central-1",
					"routing_tier":         "1",
					"latency_tolerance_ms": "150",
				},
				{
					"name":              "aws-us-west-2",
					"routing_tier":      "2",
					"latency_offset_ms": "0",
				},
			}),
			c.TestCheckSetAttr("pull_secrets", c.PullSecrets),
			c.TestCheckMapAttr("env", ConvertMapToStringMap(c.Env)),
			c.TestCheckNestedBlocks("controlplane_tracing", []map[string]interface{}{
				{
					"sampling":    fmt.Sprintf("%.0f", *c.Tracing.Sampling),
					"custom_tags": customTags,
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"dedicated":       strconv.FormatBool(*c.LoadBalancer.Dedicated),
					"trusted_proxies": strconv.Itoa(*c.LoadBalancer.TrustedProxies),
					"ipset":           *c.LoadBalancer.IpSet,
					"multi_zone": []map[string]interface{}{
						{
							"enabled": strconv.FormatBool(*c.LoadBalancer.MultiZone.Enabled),
						},
					},
					"redirect": []map[string]interface{}{
						{
							"class": []map[string]interface{}{
								{
									"status_5xx": *c.LoadBalancer.Redirect.Class.Status5XX,
									"status_401": *c.LoadBalancer.Redirect.Class.Status401,
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("keda", []map[string]interface{}{
				{
					"enabled":       "true",
					"identity_link": fmt.Sprintf("/org/%s/gvc/%s/identity/non-existant-identity", OrgName, c.Name),
				},
			}),
		),
	}
}

// BuildLocationOptionsAbsentStep returns a test step where the GVC has no location_options block configured.
func (grt *GvcResourceTest) BuildLocationOptionsAbsentStep(c GvcResourceTestCase) resource.TestStep {
	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.LocationOptionsAbsentHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			resource.TestCheckResourceAttr(c.ResourceAddress, "location_options.#", "0"),
		),
	}
}

// BuildLocationOptionsRequiredOnlyStep returns a test step with a single location_options block carrying only the required `name` attribute.
func (grt *GvcResourceTest) BuildLocationOptionsRequiredOnlyStep(c GvcResourceTestCase) resource.TestStep {
	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.LocationOptionsRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.TestCheckNestedBlocks("location_options", []map[string]interface{}{
				{
					"name": "aws-eu-central-1",
				},
			}),
		),
	}
}

// BuildLocationOptionsMultiAllStep returns a test step with two location_options blocks, both with all attributes populated.
func (grt *GvcResourceTest) BuildLocationOptionsMultiAllStep(c GvcResourceTestCase) resource.TestStep {
	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.LocationOptionsMultiAllHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.TestCheckNestedBlocks("location_options", []map[string]interface{}{
				{
					"name":                 "aws-eu-central-1",
					"routing_tier":         "1",
					"latency_offset_ms":    "0",
					"latency_tolerance_ms": "150",
				},
				{
					"name":                 "aws-us-west-2",
					"routing_tier":         "2",
					"latency_offset_ms":    "5",
					"latency_tolerance_ms": "200",
				},
			}),
		),
	}
}

// BuildLocationOptionsShrunkStep returns a test step with a single fully-populated location_options block, exercising removal of one entry while retaining the other with full attribute coverage.
func (grt *GvcResourceTest) BuildLocationOptionsShrunkStep(c GvcResourceTestCase) resource.TestStep {
	// Initialize and return the test step
	return resource.TestStep{
		Config: grt.LocationOptionsShrunkHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.TestCheckNestedBlocks("location_options", []map[string]interface{}{
				{
					"name":                 "aws-eu-central-1",
					"routing_tier":         "1",
					"latency_offset_ms":    "0",
					"latency_tolerance_ms": "150",
				},
			}),
		),
	}
}

// Configs //

// LocationOptionsAbsentHcl returns the HCL for a GVC with no location_options block configured.
func (grt *GvcResourceTest) LocationOptionsAbsentHcl(c GvcResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name        = "%s"
  description = "%s"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}
`, c.ResourceName, c.Name, c.Description)
}

// LocationOptionsRequiredOnlyHcl returns the HCL for a GVC with a single location_options block carrying only the required `name` attribute.
func (grt *GvcResourceTest) LocationOptionsRequiredOnlyHcl(c GvcResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name        = "%s"
  description = "%s"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  location_options {
    name = "aws-eu-central-1"
  }
}
`, c.ResourceName, c.Name, c.Description)
}

// LocationOptionsMultiAllHcl returns the HCL for a GVC with two location_options blocks, both with all attributes populated.
func (grt *GvcResourceTest) LocationOptionsMultiAllHcl(c GvcResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name        = "%s"
  description = "%s"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  location_options {
    name                 = "aws-eu-central-1"
    routing_tier         = 1
    latency_offset_ms    = 0
    latency_tolerance_ms = 150
  }

  location_options {
    name                 = "aws-us-west-2"
    routing_tier         = 2
    latency_offset_ms    = 5
    latency_tolerance_ms = 200
  }
}
`, c.ResourceName, c.Name, c.Description)
}

// LocationOptionsShrunkHcl returns the HCL for a GVC with a single location_options block carrying all attributes.
func (grt *GvcResourceTest) LocationOptionsShrunkHcl(c GvcResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name        = "%s"
  description = "%s"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  location_options {
    name                 = "aws-eu-central-1"
    routing_tier         = 1
    latency_offset_ms    = 0
    latency_tolerance_ms = 150
  }
}
`, c.ResourceName, c.Name, c.Description)
}

// GvcRequiredOnly returns a minimal HCL block for a GVC using only required fields.
func (grt *GvcResourceTest) GvcRequiredOnly(c GvcResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

// GvcRequiredOnlyWithEndpointNamingFormat returns a minimal HCL block for a GVC using only required fields.
func (grt *GvcResourceTest) GvcRequiredOnlyWithEndpointNamingFormat(c GvcResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name                   = "%s"
  endpoint_naming_format = "%s"
}
`, c.ResourceName, c.Name, c.EndpointNamingFormat)
}

// UpdateWithMinimalOptionals returns a HCL block for a GVC using minimal optional attributes.
func (grt *GvcResourceTest) UpdateWithMinimalOptionals(c GvcResourceTestCase, opaqueSecretResource string, tracingBlock string) string {
	return fmt.Sprintf(`
# Opaque Secret Resource
%s

resource "cpln_gvc" "%s" {
  name        = "%s"
  description = "%s"

	endpoint_naming_format = "%s"
  locations              = %s
  pull_secrets           = %s

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  # Env Block
  env = %s

  # Tracing Block
  %s

  load_balancer {}

  sidecar {
    envoy = jsonencode(%s)
  }

  keda {}
}
`, opaqueSecretResource, c.ResourceName, c.Name, c.DescriptionUpdate, c.EndpointNamingFormat, StringSliceToString(c.Locations), StringSliceToString(c.PullSecrets),
		MapToHCL(c.Env, 2), tracingBlock, c.Envoy,
	)
}

// UpdateWithAllOptionals returns a HCL block for a GVC using all attributes.
func (grt *GvcResourceTest) UpdateWithAllOptionals(c GvcResourceTestCase, opaqueSecretResource string, tracingBlock string) string {
	return fmt.Sprintf(`
# Opaque Secret Resource
%s

resource "cpln_gvc" "%s" {
  name        = "%s"
  description = "%s"

	endpoint_naming_format = "%s"
  locations              = %s
  pull_secrets           = %s

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  location_options {
    name                 = "aws-eu-central-1"
    routing_tier         = 1
    latency_tolerance_ms = 150
  }

  location_options {
    name              = "aws-us-west-2"
    routing_tier      = 2
    latency_offset_ms = 0
  }

  # Env Block
  env = %s

  # Tracing Block
  %s

  load_balancer {
    dedicated       = %s
    trusted_proxies = %d
    ipset           = "%s"

    multi_zone {
      enabled = %s
    }

    redirect {
      class {
        status_5xx = "%s"
        status_401 = "%s"
      }
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  keda {
    enabled       = true
    identity_link = "/org/%s/gvc/%s/identity/non-existant-identity"
    secrets       = ["/org/terraform-test-org/secret/non-existant-secret-01", "/org/terraform-test-org/secret/non-existant-secret-02", "/org/terraform-test-org/secret/non-existant-secret-03"]
  }
}
`, opaqueSecretResource, c.ResourceName, c.Name, c.DescriptionUpdate, c.EndpointNamingFormat, StringSliceToString(c.Locations), StringSliceToString(c.PullSecrets),
		MapToHCL(c.Env, 2), tracingBlock, strconv.FormatBool(*c.LoadBalancer.Dedicated), *c.LoadBalancer.TrustedProxies, *c.LoadBalancer.IpSet,
		strconv.FormatBool(*c.LoadBalancer.MultiZone.Enabled), *c.LoadBalancer.Redirect.Class.Status5XX, *c.LoadBalancer.Redirect.Class.Status401, c.Envoy, OrgName, c.Name,
	)
}

// Tracing Config //

// LightstepTracingRequiredOnly defines the HCL for the lightstep tracing with minimal attributes.
func (grt *GvcResourceTest) LightstepTracingRequiredOnly(c GvcResourceTestCase) string {
	return fmt.Sprintf(`
  lightstep_tracing {
    sampling = "%f"
    endpoint = "%s"
  }
`, *c.Tracing.Sampling, *c.Tracing.Provider.Lightstep.Endpoint)
}

// LightstepTracingWithOptionals defines the HCL for the lightstep tracing with all attributes.
func (grt *GvcResourceTest) LightstepTracingWithOptionals(c GvcResourceTestCase, credentials string, customTags map[string]interface{}) string {
	return fmt.Sprintf(`
  lightstep_tracing {
    sampling = "%f"
    endpoint = "%s"

    # Opaque Secret Only
    credentials = %s

    # Custom Tags
    custom_tags = %s
  }
`, *c.Tracing.Sampling, *c.Tracing.Provider.Lightstep.Endpoint, credentials, MapToHCL(customTags, 2))
}

// OtelTracingHcl defines the HCL for the lightstep tracing.
func (grt *GvcResourceTest) OtelTracingHcl(c GvcResourceTestCase, customTags map[string]interface{}) string {
	return fmt.Sprintf(`
  otel_tracing {
    sampling = "%f"
    endpoint = "%s"

    # Custom Tags
    custom_tags = %s
  }
`, *c.Tracing.Sampling, *c.Tracing.Provider.Otel.Endpoint, MapToHCL(customTags, 2))
}

// ControlPlaneTracingHcl defines the HCL for the lightstep tracing.
func (grt *GvcResourceTest) ControlPlaneTracingHcl(c GvcResourceTestCase, customTags map[string]interface{}) string {
	return fmt.Sprintf(`
  controlplane_tracing {
    sampling = "%f"

    # Custom Tags
    custom_tags = %s
  }
`, *c.Tracing.Sampling, MapToHCL(customTags, 2))
}

// Helpers //

// ConvertCustomTagsToMap converts map[string]client.TracingCustomTag instances to a plain map for HCL generation and test comparisons.
func (grt *GvcResourceTest) ConvertCustomTagsToMap(tags map[string]client.TracingCustomTag) map[string]interface{} {
	// Initialize output map with capacity matching the number of tags
	out := make(map[string]interface{}, len(tags))

	// Populate the map with literal values from tags
	for key, tag := range tags {
		// If the tag has a literal value, use it
		if tag.Literal != nil {
			out[key] = *tag.Literal.Value
		}
	}

	// Return the resulting map
	return out
}

/*** Resource Test Case ***/

// GvcResourceTestCase defines a specific resource test case.
type GvcResourceTestCase struct {
	ProviderTestCase
	EndpointNamingFormat string
	Locations            []string
	PullSecrets          []string
	Env                  map[string]interface{}
	Tracing              client.Tracing
	LoadBalancer         client.GvcLoadBalancer
	Envoy                string
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (grtc *GvcResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of GVC: %s. Total resources: %d", grtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[grtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", grtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != grtc.Name {
			return fmt.Errorf("resource ID %s does not match expected GVC name %s", rs.Primary.ID, grtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteGvc, _, err := TestProvider.client.GetGvc(grtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving GVC from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteGvc.Name != grtc.Name {
			return fmt.Errorf("mismatch in GVC name: expected %s, got %s", grtc.Name, *remoteGvc.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("GVC %s verified successfully in both state and external system.", grtc.Name))
		return nil
	}
}
