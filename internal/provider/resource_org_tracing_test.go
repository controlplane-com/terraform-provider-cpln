package cpln

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

/*** Acceptance Test ***/

// TestAccControlPlaneOrgTracing_basic performs an acceptance test for the resource.
func TestAccControlPlaneOrgTracing_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewOrgTracingResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "ORG_TRACING") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// OrgTracingResourceTest defines the necessary functionality to test the resource.
type OrgTracingResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
}

// NewOrgTracingResourceTest creates a OrgTracingResourceTest with initialized test cases.
func NewOrgTracingResourceTest() OrgTracingResourceTest {
	// Create a resource test instance
	resourceTest := OrgTracingResourceTest{
		RandomName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
	}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewLightstepTracingScenario()...)
	steps = append(steps, resourceTest.NewOtelTracingScenario()...)
	steps = append(steps, resourceTest.NewCplnTracingScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (otrt *OrgTracingResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_org_tracing resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_org_tracing" {
			continue
		}

		// Retrieve the name for the current resource
		orgName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of org with name: %s", orgName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		org, _, _ := TestProvider.client.GetOrg()

		// Make sure the org has no tracing spec at all
		if org.Spec.Tracing != nil {
			return fmt.Errorf("Org Spec Tracing still exists. Org Name: %s", *org.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_org_tracing resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewLightstepTracingScenario defines an org-level tracing test scenario using Lightstep with create, import, and update steps.
func (otrt *OrgTracingResourceTest) NewLightstepTracingScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "tf-tracing"

	// Build test steps
	initialConfig, initialStep := otrt.BuildLightstepTracingInitialTestStep(resourceName)
	caseUpdate1 := otrt.BuildLightstepTracingUpdate1TestStep(initialConfig.ProviderTestCase)

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
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewOtelTracingScenario defines an org-level tracing test scenario using OpenTelemetry with create, import, and update steps.
func (otrt *OrgTracingResourceTest) NewOtelTracingScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "tf-tracing"

	// Build test steps
	initialConfig, initialStep := otrt.BuildOtelTracingInitialTestStep(resourceName)
	caseUpdate1 := otrt.BuildOtelTracingUpdate1TestStep(initialConfig.ProviderTestCase)

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
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewCplnTracingScenario defines an org-level tracing test scenario using Control Plane implementation with create, import, and update steps.
func (otrt *OrgTracingResourceTest) NewCplnTracingScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "tf-tracing"

	// Build test steps
	initialConfig, initialStep := otrt.BuildCplnTracingInitialTestStep(resourceName)
	caseUpdate1 := otrt.BuildCplnTracingUpdate1TestStep(initialConfig.ProviderTestCase)

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
		// Revert the resource to its initial state
		initialStep,
	}
}

// Test Cases //

// BuildLightstepTracingInitialTestStep constructs the initial test step and case for Lightstep tracing configuration.
func (otrt *OrgTracingResourceTest) BuildLightstepTracingInitialTestStep(resourceName string) (OrgTracingResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := OrgTracingResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "org",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_org_tracing.%s", resourceName),
			Name:            OrgName,
			Description:     OrgName,
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: otrt.LightstepTracingRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("lightstep_tracing", []map[string]interface{}{
				{
					"sampling": "55.55",
					"endpoint": "test.cpln.local:8080",
				},
			}),
		),
	}
}

// BuildLightstepTracingUpdate1TestStep constructs the update test step for Lightstep tracing including credentials and custom tags.
func (otrt *OrgTracingResourceTest) BuildLightstepTracingUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgTracingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: otrt.LightstepTracingWithAllAttributesHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("lightstep_tracing", []map[string]interface{}{
				{
					"sampling":    "55.55",
					"endpoint":    "test.cpln.local:8080",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-%s", otrt.RandomName)),
					"custom_tags": map[string]interface{}{
						"key": "value",
					},
				},
			}),
		),
	}
}

// BuildOtelTracingInitialTestStep constructs the initial test step and case for OpenTelemetry tracing configuration.
func (otrt *OrgTracingResourceTest) BuildOtelTracingInitialTestStep(resourceName string) (OrgTracingResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := OrgTracingResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "org",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_org_tracing.%s", resourceName),
			Name:            OrgName,
			Description:     OrgName,
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: otrt.OtelTracingRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("otel_tracing", []map[string]interface{}{
				{
					"sampling": "70",
					"endpoint": "test.cpln.local:80",
				},
			}),
		),
	}
}

// BuildOtelTracingUpdate1TestStep constructs the update test step for OpenTelemetry tracing including custom tags.
func (otrt *OrgTracingResourceTest) BuildOtelTracingUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgTracingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: otrt.OtelTracingWithAllAttributesHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("otel_tracing", []map[string]interface{}{
				{
					"sampling": "60.2",
					"endpoint": "test.cpln.local:443",
					"custom_tags": map[string]interface{}{
						"key01": "value-01",
						"key02": "value-02",
					},
				},
			}),
		),
	}
}

// BuildCplnTracingInitialTestStep constructs the initial test step and case for Control Plane tracing configuration.
func (otrt *OrgTracingResourceTest) BuildCplnTracingInitialTestStep(resourceName string) (OrgTracingResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := OrgTracingResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "org",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_org_tracing.%s", resourceName),
			Name:            OrgName,
			Description:     OrgName,
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: otrt.ControlPlaneTracingRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("controlplane_tracing", []map[string]interface{}{
				{
					"sampling": "75",
				},
			}),
		),
	}
}

// BuildCplnTracingUpdate1TestStep constructs the update test step for Control Plane tracing including custom tags.
func (otrt *OrgTracingResourceTest) BuildCplnTracingUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgTracingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: otrt.ControlPlaneTracingWithAllAttributesHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("controlplane_tracing", []map[string]interface{}{
				{
					"sampling": "50",
					"custom_tags": map[string]interface{}{
						"key01": "value-01",
						"key02": "value-02",
						"key03": "value-03",
					},
				},
			}),
		),
	}
}

// Configs //

// LightstepTracingRequiredOnlyHcl returns a minimal HCL block for Lightstep tracing configuration.
func (otrt *OrgTracingResourceTest) LightstepTracingRequiredOnlyHcl(c OrgTracingResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org_tracing" "%s" {

  lightstep_tracing {
    sampling = 55.55
    endpoint = "test.cpln.local:8080"
	}
}
`, c.ResourceName)
}

// LightstepTracingHcl returns an HCL block for Lightstep tracing including credentials and custom_tags.
func (otrt *OrgTracingResourceTest) LightstepTracingWithAllAttributesHcl(c OrgTracingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "opaque" {
  name        = "tf-opaque-${var.random_name}"
  description = "description opaque-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
	} 

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_tracing" "%s" {

  lightstep_tracing {
    sampling = 55.55
    endpoint = "test.cpln.local:8080"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link

    custom_tags = {
		  "key" = "value",
		}
	}
}
`, otrt.RandomName, c.ResourceName)
}

// OtelTracingRequiredOnlyHcl returns a minimal HCL block for Otel tracing configuration.
func (otrt *OrgTracingResourceTest) OtelTracingRequiredOnlyHcl(c OrgTracingResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org_tracing" "%s" {

  otel_tracing {
    sampling = 70
    endpoint = "test.cpln.local:80"
	}
}
`, c.ResourceName)
}

// OtelTracingWithAllAttributesHcl returns an HCL block for Otel tracing including custom_tags.
func (otrt *OrgTracingResourceTest) OtelTracingWithAllAttributesHcl(c OrgTracingResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org_tracing" "%s" {

  otel_tracing {
    sampling = 60.2
    endpoint = "test.cpln.local:443"

    custom_tags = {
      key01 = "value-01"
      key02 = "value-02"
    }
  }
}
`, c.ResourceName)
}

// ControlPlaneTracingRequiredOnlyHcl returns a minimal HCL block for Control Plane tracing configuration.
func (otrt *OrgTracingResourceTest) ControlPlaneTracingRequiredOnlyHcl(c OrgTracingResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org_tracing" "%s" {

  controlplane_tracing {
    sampling = 75
	}
}
`, c.ResourceName)
}

// ControlPlaneTracingWithAllAttributesHcl returns an HCL block for Control Plane tracing including custom_tags.
func (otrt *OrgTracingResourceTest) ControlPlaneTracingWithAllAttributesHcl(c OrgTracingResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org_tracing" "%s" {

  controlplane_tracing {
    sampling = 50

    custom_tags = {
      key01 = "value-01"
      key02 = "value-02"
      key03 = "value-03"
    }
  }
}
`, c.ResourceName)
}

/*** Resource Test Case ***/

// OrgTracingResourceTestCase defines a specific resource test case.
type OrgTracingResourceTestCase struct {
	ProviderTestCase
	LogzioListenerHost string
}
