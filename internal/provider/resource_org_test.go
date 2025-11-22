package cpln

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneOrg_basic performs an acceptance test for the resource.
func TestAccControlPlaneOrg_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewOrgResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "ORG") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// OrgResourceTest defines the necessary functionality to test the resource.
type OrgResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
}

// NewOrgResourceTest creates a OrgResourceTest with initialized test cases.
func NewOrgResourceTest() OrgResourceTest {
	// Create a resource test instance
	resourceTest := OrgResourceTest{
		RandomName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
	}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewDefaultScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (ort *OrgResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_org resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_org" {
			continue
		}

		// Retrieve the name for the current resource
		orgName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of org with name: %s", orgName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		org, code, err := TestProvider.client.GetOrg()

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if org %s exists: %w", orgName, err)
		}

		// Make sure that logging is reset after resource deletion
		if org.Spec.Logging != nil || (org.Spec.ExtraLogging != nil && len(*org.Spec.ExtraLogging) != 0) {
			return fmt.Errorf("Org Spec Logging still exists. Org Name: %s", *org.Name)
		}

		// Make sure that tracing is reset after resource deletion
		if org.Spec.Tracing != nil {
			return fmt.Errorf("Org Spec Tracing still exists. Org Name: %s", *org.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_org resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case with initial and updated configurations.
func (ort *OrgResourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "new"

	cLogging := OrgLoggingResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "org",
			ResourceName:    "tf-logging",
			ResourceAddress: "cpln_org_logging.tf-logging",
			Name:            OrgName,
			Description:     OrgName,
		},
	}

	// Build test steps
	initialConfig, initialStep := ort.BuildInitialTestStep(resourceName)
	caseUpdate1 := ort.BuildUpdate1TestStep(initialConfig.ProviderTestCase)
	caseUpdate2 := ort.BuildUpdate2TestStep(initialConfig.ProviderTestCase, cLogging)
	caseUpdate3 := ort.BuildUpdate3TestStep(initialConfig.ProviderTestCase, cLogging)

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
		// Revert the resource to its initial state
		initialStep,
	}
}

// Test Cases //

// BuildInitialTestStep returns a default initial test step and its associated test case for the resource.
func (ort *OrgResourceTest) BuildInitialTestStep(resourceName string) (OrgResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := OrgResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "org",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_org.%s", resourceName),
			Name:              OrgName,
			Description:       OrgName,
			DescriptionUpdate: "testing",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: ort.HclRequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("observability", []map[string]interface{}{
				{
					"logs_retention_days":    "30",
					"metrics_retention_days": "30",
					"traces_retention_days":  "30",
				},
			}),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (ort *OrgResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: ort.HclUpdateWithMinimalOptionals(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckNestedBlocks("observability", []map[string]interface{}{
				{
					"logs_retention_days":    "55",
					"metrics_retention_days": "65",
					"traces_retention_days":  "75",
					"default_alert_emails":   []string{"bob@example.com", "rob@example.com"},
				},
			}),
			c.TestCheckNestedBlocks("security", []map[string]interface{}{
				{},
			}),
		),
	}
}

// BuildUpdate2TestStep returns a test step for the update.
func (ort *OrgResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase, cLogging OrgLoggingResourceTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: ort.HclUpdateWithLogging(c, cLogging),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckNestedBlocks("observability", []map[string]interface{}{
				{
					"logs_retention_days":    "55",
					"metrics_retention_days": "65",
					"traces_retention_days":  "75",
					"default_alert_emails":   []string{"bob@example.com", "rob@example.com", "abby@example.com"},
				},
			}),
			c.TestCheckNestedBlocks("security", []map[string]interface{}{
				{},
			}),
			cLogging.TestCheckNestedBlocks("coralogix_logging", []map[string]interface{}{
				{
					"cluster":     "coralogix.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-coralogix-%s", ort.RandomName)),
					"app":         "{workload}",
					"subsystem":   "{org}",
				},
			}),
			cLogging.TestCheckNestedBlocks("datadog_logging", []map[string]interface{}{
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-00-%s", ort.RandomName)),
				},
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-01-%s", ort.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate3TestStep returns a test step for the update.
func (ort *OrgResourceTest) BuildUpdate3TestStep(initialCase ProviderTestCase, cLogging OrgLoggingResourceTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: ort.HclUpdateWithAllAttributes(c, cLogging),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckNestedBlocks("observability", []map[string]interface{}{
				{
					"logs_retention_days":    "55",
					"metrics_retention_days": "65",
					"traces_retention_days":  "75",
					"default_alert_emails":   []string{"bob@example.com", "rob@example.com", "abby@example.com", "david@example.com"},
				},
			}),
			c.TestCheckNestedBlocks("security", []map[string]interface{}{
				{
					"threat_detection": []map[string]interface{}{
						{
							"enabled":          "true",
							"minimum_severity": "warning",
							"syslog": []map[string]interface{}{
								{
									"transport": "tcp",
									"host":      "example.com",
									"port":      "8080",
								},
							},
						},
					},
				},
			}),
			cLogging.TestCheckNestedBlocks("coralogix_logging", []map[string]interface{}{
				{
					"cluster":     "coralogix.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-coralogix-%s", ort.RandomName)),
					"app":         "{workload}",
					"subsystem":   "{org}",
				},
			}),
			cLogging.TestCheckNestedBlocks("datadog_logging", []map[string]interface{}{
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-00-%s", ort.RandomName)),
				},
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-01-%s", ort.RandomName)),
				},
			}),
		),
	}
}

// Configs //

// HclRequiredOnly returns a minimal HCL block for an Org resource with default observability settings.
func (ort *OrgResourceTest) HclRequiredOnly(c OrgResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org" "%s" {
  observability {
    logs_retention_days    = 30
    metrics_retention_days = 30
    traces_retention_days  = 30
  }
}
`, c.ResourceName)
}

// HclUpdateWithMinimalOptionals returns a minimal update HCL block for an Org resource including description, tags, and observability tweaks.
func (ort *OrgResourceTest) HclUpdateWithMinimalOptionals(c OrgResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org" "%s" {
  description             = "%s"
  session_timeout_seconds = 50000

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  observability {
    logs_retention_days    = 55
    metrics_retention_days = 65
    traces_retention_days  = 75
    default_alert_emails   = ["bob@example.com", "rob@example.com"]
  }

  security {}
}
`, c.ResourceName, c.DescriptionUpdate)
}

// HclUpdateWithLogging returns a minimal update HCL block for an Org resource including description, tags, and observability tweaks and logging.
func (ort *OrgResourceTest) HclUpdateWithLogging(c OrgResourceTestCase, cLogging OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_org" "%s" {
  description             = "%s"
  session_timeout_seconds = 50000

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  observability {
    logs_retention_days    = 55
    metrics_retention_days = 65
    traces_retention_days  = 75
    default_alert_emails   = ["bob@example.com", "rob@example.com", "abby@example.com"]
  }

  security {}
}

resource "cpln_secret" "opaque-coralogix" {
  name        = "tf-opaque-random-coralogix-${var.random_name}"
  description = "opaque description opaque-random-tbd"

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

resource "cpln_secret" "opaque-datadog" {
  name        = "tf-opaque-random-datadog-00-${var.random_name}"
	description = "opaque description"

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

resource "cpln_secret" "opaque-datadog-1" {
  name        = "tf-opaque-random-datadog-01-${var.random_name}"
  description = "opaque description"

  tags = {
    terraform_generated = "true"
    acceptance_test = "true"
    secret_type = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "%s" {

  coralogix_logging {
    cluster = "coralogix.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-coralogix.self_link

    app       = "{workload}"
    subsystem = "{org}"
  }

  datadog_logging {
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-datadog.self_link
  }

  datadog_logging {
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-datadog-1.self_link
  }
}
`, ort.RandomName, c.ResourceName, c.DescriptionUpdate, cLogging.ResourceName)
}

// HclUpdateWithAllAttributes returns a full update HCL block for an Org resource with threat detection and observability customizations.
func (ort *OrgResourceTest) HclUpdateWithAllAttributes(c OrgResourceTestCase, cLogging OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_org" "%s" {
  description             = "%s"
  session_timeout_seconds = 50000

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  observability {
    logs_retention_days    = 55
    metrics_retention_days = 65
    traces_retention_days  = 75
    default_alert_emails   = ["bob@example.com", "rob@example.com", "abby@example.com", "david@example.com"]
  }

  security {
    threat_detection {
      enabled          = true
      minimum_severity = "warning"

      syslog {
        transport = "tcp"
        host      = "example.com"
        port      = 8080
      }
    }
  }
}

resource "cpln_secret" "opaque-coralogix" {
  name        = "tf-opaque-random-coralogix-${var.random_name}"
  description = "opaque description opaque-random-tbd"

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

resource "cpln_secret" "opaque-datadog" {
  name        = "tf-opaque-random-datadog-00-${var.random_name}"
	description = "opaque description"

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

resource "cpln_secret" "opaque-datadog-1" {
  name        = "tf-opaque-random-datadog-01-${var.random_name}"
  description = "opaque description"

  tags = {
    terraform_generated = "true"
    acceptance_test = "true"
    secret_type = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "%s" {

  coralogix_logging {
    cluster = "coralogix.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-coralogix.self_link

    app       = "{workload}"
    subsystem = "{org}"
  }

  datadog_logging {
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-datadog.self_link
  }

  datadog_logging {
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-datadog-1.self_link
  }
}
`, ort.RandomName, c.ResourceName, c.DescriptionUpdate, cLogging.ResourceName)
}

/*** Resource Test Case ***/

// OrgResourceTestCase defines a specific resource test case.
type OrgResourceTestCase struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (ortc *OrgResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of org: %s. Total resources: %d", ortc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[ortc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", ortc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != ortc.Name {
			return fmt.Errorf("resource ID %s does not match expected org name %s", rs.Primary.ID, ortc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteOrg, _, err := TestProvider.client.GetOrg()
		if err != nil {
			return fmt.Errorf("error retrieving org from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteOrg.Name != ortc.Name {
			return fmt.Errorf("mismatch in org name: expected %s, got %s", ortc.Name, *remoteOrg.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("org %s verified successfully in both state and external system.", ortc.Name))
		return nil
	}
}
