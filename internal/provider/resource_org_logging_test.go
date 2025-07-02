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

// TestAccControlPlaneOrgLogging_basic performs an acceptance test for the resource.
func TestAccControlPlaneOrgLogging_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewOrgLoggingResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "ORG_LOGGING") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// OrgLoggingResourceTest defines the necessary functionality to test the resource.
type OrgLoggingResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
}

// NewOrgLoggingResourceTest creates a OrgLoggingResourceTest with initialized test cases.
func NewOrgLoggingResourceTest() OrgLoggingResourceTest {
	// Create a resource test instance
	resourceTest := OrgLoggingResourceTest{
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
func (olrt *OrgLoggingResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_org_logging resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_org_logging" {
			continue
		}

		// Retrieve the name for the current resource
		orgName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of org with name: %s", orgName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		org, _, _ := TestProvider.client.GetOrg()

		// Make sure the org has no logging spec at all
		if org.Spec.Logging != nil || (org.Spec.ExtraLogging != nil && len(*org.Spec.ExtraLogging) != 0) {
			return fmt.Errorf("Org Spec Logging still exists. Org Name: %s", *org.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_org_logging resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case with initial and updated configurations.
func (olrt *OrgLoggingResourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "tf-logging"

	// Build test steps
	initialConfig, initialStep := olrt.BuildInitialTestStep(resourceName)
	caseUpdate1 := olrt.BuildUpdate1TestStep(initialConfig.ProviderTestCase)
	caseUpdate2 := olrt.BuildUpdate2TestStep(initialConfig.ProviderTestCase)
	caseUpdate3 := olrt.BuildUpdate3TestStep(initialConfig.ProviderTestCase)
	caseUpdate4 := olrt.BuildUpdate4TestStep(initialConfig.ProviderTestCase)
	caseUpdate5 := olrt.BuildUpdate5TestStep(initialConfig.ProviderTestCase)
	caseUpdate6 := olrt.BuildUpdate6TestStep(initialConfig.ProviderTestCase)
	caseUpdate7 := olrt.BuildUpdate7TestStep(initialConfig.ProviderTestCase)
	caseUpdate8 := olrt.BuildUpdate8TestStep(initialConfig.ProviderTestCase)
	caseUpdate9 := olrt.BuildUpdate9TestStep(initialConfig.ProviderTestCase)
	caseUpdate10 := olrt.BuildUpdate10TestStep(initialConfig.ProviderTestCase)
	caseUpdate11 := olrt.BuildUpdate11TestStep(initialConfig.ProviderTestCase)
	caseUpdate12 := olrt.BuildUpdate12TestStep(initialConfig.ProviderTestCase)
	caseUpdate13 := olrt.BuildUpdate13TestStep(initialConfig.ProviderTestCase)

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
		caseUpdate5,
		caseUpdate6,
		caseUpdate7,
		caseUpdate8,
		caseUpdate9,
		caseUpdate10,
		caseUpdate11,
		caseUpdate12,
		caseUpdate13,
		// Revert the resource to its initial state
		initialStep,
	}
}

// Test Cases //

// BuildInitialTestStep returns a default initial test step and its associated test case for the resource.
func (olrt *OrgLoggingResourceTest) BuildInitialTestStep(resourceName string) (OrgLoggingResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "org",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_org_logging.%s", resourceName),
			Name:            OrgName,
			Description:     OrgName,
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: olrt.HclLoggingS3(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("s3_logging", []map[string]interface{}{
				{
					"bucket":      "test-bucket",
					"region":      "us-east1",
					"prefix":      "/",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-aws-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingCoralogix(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("coralogix_logging", []map[string]interface{}{
				{
					"cluster":     "coralogix.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-%s", olrt.RandomName)),
					"app":         "{workload}",
					"subsystem":   "{org}",
				},
			}),
		),
	}
}

// BuildUpdate2TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingDatadog(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("datadog_logging", []map[string]interface{}{
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-00-%s", olrt.RandomName)),
				},
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-01-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate3TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate3TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase:   initialCase,
		LogzioListenerHost: "listener.logz.io",
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingLogzio(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("logzio_logging", []map[string]interface{}{
				{
					"listener_host": c.LogzioListenerHost,
					"credentials":   GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate4TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate4TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase:   initialCase,
		LogzioListenerHost: "listener-nl.logz.io",
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingLogzio(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("logzio_logging", []map[string]interface{}{
				{
					"listener_host": c.LogzioListenerHost,
					"credentials":   GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate5TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate5TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingElasticAws(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("elastic_logging", []map[string]interface{}{
				{
					"aws": []map[string]interface{}{
						{
							"host":        "es.amazonaws.com",
							"port":        "8080",
							"index":       "my-index",
							"type":        "my-type",
							"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-aws-%s", olrt.RandomName)),
							"region":      "us-east-1",
						},
					},
				},
			}),
		),
	}
}

// BuildUpdate6TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate6TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingElasticCloud(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("elastic_logging", []map[string]interface{}{
				{
					"elastic_cloud": []map[string]interface{}{
						{
							"index":       "my-index",
							"type":        "my-type",
							"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-userpass-elastic-cloud-%s", olrt.RandomName)),
							"cloud_id":    "my-cloud-id",
						},
					},
				},
			}),
		),
	}
}

// BuildUpdate7TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate7TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingElasticGeneric(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("elastic_logging", []map[string]interface{}{
				{
					"generic": []map[string]interface{}{
						{
							"host":        "example.com",
							"port":        "9200",
							"path":        "/var/log/elasticsearch/",
							"index":       "my-index",
							"type":        "my-type",
							"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-userpass-elastic-generic-%s", olrt.RandomName)),
						},
					},
				},
			}),
		),
	}
}

// BuildUpdate8TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate8TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingCloudWatch(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("cloud_watch_logging", []map[string]interface{}{
				{
					"region":         "us-east-1",
					"retention_days": "1",
					"group_name":     "demo-group-name",
					"stream_name":    "demo-stream-name",
					"extract_fields": map[string]interface{}{
						"log_level": "$.level",
					},
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate9TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate9TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingFluentd(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("fluentd_logging", []map[string]interface{}{
				{
					"host": "example.com",
					"port": "24224",
				},
			}),
		),
	}
}

// BuildUpdate10TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate10TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingStackdriver(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("stackdriver_logging", []map[string]interface{}{
				{
					"location":    "us-east4",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate11TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate11TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclLoggingSyslog(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("syslog_logging", []map[string]interface{}{
				{
					"host":     "syslog.example.com",
					"port":     "443",
					"mode":     "tcp",
					"format":   "rfc5424",
					"severity": "6",
				},
			}),
		),
	}
}

// BuildUpdate12TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate12TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclThreeUniqueLoggings(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("s3_logging", []map[string]interface{}{
				{
					"bucket":      "test-bucket",
					"region":      "us-east1",
					"prefix":      "/",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-aws-%s", olrt.RandomName)),
				},
			}),
			c.TestCheckNestedBlocks("coralogix_logging", []map[string]interface{}{
				{
					"cluster":     "coralogix.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-coralogix-%s", olrt.RandomName)),
					"app":         "{workload}",
					"subsystem":   "{org}",
				},
			}),
			c.TestCheckNestedBlocks("datadog_logging", []map[string]interface{}{
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-00-%s", olrt.RandomName)),
				},
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-01-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate13TestStep returns a test step for the update.
func (olrt *OrgLoggingResourceTest) BuildUpdate13TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := OrgLoggingResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: olrt.HclTwoUniqueLoggings(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("coralogix_logging", []map[string]interface{}{
				{
					"cluster":     "coralogix.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-coralogix-%s", olrt.RandomName)),
					"app":         "{workload}",
					"subsystem":   "{org}",
				},
			}),
			c.TestCheckNestedBlocks("datadog_logging", []map[string]interface{}{
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-00-%s", olrt.RandomName)),
				},
				{
					"host":        "http-intake.logs.datadoghq.com",
					"credentials": GetSelfLink(OrgName, "secret", fmt.Sprintf("tf-opaque-random-datadog-01-%s", olrt.RandomName)),
				},
			}),
		),
	}
}

// Configs //

// HclLoggingS3 returns a minimal HCL block for a resource using only required fields.
func (olrt *OrgLoggingResourceTest) HclLoggingS3(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "aws" {
  name        = "tf-aws-${var.random_name}"
  description = "aws description aws-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "aws"
  }

  aws {
    secret_key = "AKIAIOSFODNN7EXAMPLE"
    access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    role_arn   = "arn:awskey" 
  }
}

resource "cpln_org_logging" "%s" {
  s3_logging {
    bucket = "test-bucket"
    region = "us-east1"
    prefix = "/"

    // AWS Secret Only
    credentials = cpln_secret.aws.self_link
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingCoralogix returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingCoralogix(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "opaque" {
  name        = "tf-opaque-${var.random_name}"
  description = "opaque description opaque-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "%s" {

  coralogix_logging {

    // Valid clusters
    // coralogix.com, coralogix.us, app.coralogix.in, app.eu2.coralogix.com, app.coralogixsg.com
    cluster = "coralogix.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link

    // Supported variables for App and Subsystem are:
    // {org}, {gvc}, {workload}, {location}
    app       = "{workload}"
    subsystem = "{org}"
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingDatadog returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingDatadog(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "opaque-00" {
  name        = "tf-opaque-00-${var.random_name}"
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

resource "cpln_secret" "opaque-01" {
  name        = "tf-opaque-01-${var.random_name}"
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

resource "cpln_org_logging" "%s" {

  datadog_logging {

    // Valid Host
    // http-intake.logs.datadoghq.com, http-intake.logs.us3.datadoghq.com, 
    // http-intake.logs.us5.datadoghq.com, http-intake.logs.datadoghq.eu
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-00.self_link  
  }

  datadog_logging {
    host = "http-intake.logs.datadoghq.com"

    // Opaque Secret Only
    credentials = cpln_secret.opaque-01.self_link  
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingLogzio returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingLogzio(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "opaque" {
  name        = "tf-opaque-${var.random_name}"
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

resource "cpln_org_logging" "%s" {

  logzio_logging {

    // Valid Listener Hosts
    // listener.logz.io, listener-nl.logz.io 
    listener_host = "%s"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link  
  }
}
`, olrt.RandomName, c.ResourceName, c.LogzioListenerHost)
}

// HclLoggingElasticAws returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingElasticAws(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "aws" {
  name        = "tf-aws-${var.random_name}"
  description = "aws description"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "aws"
  }

  aws {
    secret_key = "AKIAIOSFODNN7EXAMPLEUPDATE"
    access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEYUPDATE"
    role_arn = "arn:awskeyupdate"
  }
}

resource "cpln_org_logging" "%s" {

  elastic_logging {
    aws {
      host        = "es.amazonaws.com"
      port        = 8080
      index       = "my-index"
      type        = "my-type"
      credentials = cpln_secret.aws.self_link
      region      = "us-east-1"
    }
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingElasticCloud returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingElasticCloud(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "userpass-elastic-cloud" {
  name        = "tf-userpass-elastic-cloud-${var.random_name}"
  description = "userpass-elastic-cloud description"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "userpass"
  }

  userpass {
    username = "cpln_username"
    password = "cpln_password"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "%s" {

  elastic_logging {
    elastic_cloud {
      index       = "my-index"
      type        = "my-type"
      credentials = cpln_secret.userpass-elastic-cloud.self_link
      cloud_id    = "my-cloud-id"
    }
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingElasticGeneric returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingElasticGeneric(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "userpass-elastic-generic" {
  name        = "tf-userpass-elastic-generic-${var.random_name}"
  description = "userpass-elastic-generic description"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "userpass"
  }

  userpass {
    username = "cpln_username"
    password = "cpln_password"
    encoding = "plain"
  }
}

resource "cpln_org_logging" "%s" {

  elastic_logging {
    generic {
      host  = "example.com"
      port  = 9200
      path  = "/var/log/elasticsearch/"
      index = "my-index"
      type  = "my-type"
      credentials = cpln_secret.userpass-elastic-generic.self_link
    }
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingCloudWatch returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingCloudWatch(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "opaque" {
  name        = "tf-opaque-${var.random_name}"
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

resource "cpln_org_logging" "%s" {

  cloud_watch_logging {
    region         = "us-east-1"
    retention_days = 1
    group_name     = "demo-group-name"
    stream_name    = "demo-stream-name"

    extract_fields  = {
      log_level = "$.level"
    }

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingFluentd returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingFluentd(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org_logging" "%s" {

  fluentd_logging {
    host = "example.com"
    port = 24224
  }
}
`, c.ResourceName)
}

// HclLoggingStackdriver returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingStackdriver(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "opaque" {
  name        = "tf-opaque-${var.random_name}"
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

resource "cpln_org_logging" "%s" {

  stackdriver_logging {
    location = "us-east4"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link
  }
}
`, olrt.RandomName, c.ResourceName)
}

// HclLoggingSyslog returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclLoggingSyslog(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_org_logging" "%s" {

  syslog_logging {
    host     = "syslog.example.com"
    port     = 443
    mode     = "tcp"
    format   = "rfc5424"
    severity = 6
  }
}
`, c.ResourceName)
}

// HclThreeUniqueLoggings returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclThreeUniqueLoggings(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_secret" "aws" {
  name        = "tf-aws-${var.random_name}"
  description = "aws description aws-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "aws"
  }

  aws {
    secret_key = "AKIAIOSFODNN7EXAMPLE"
    access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    role_arn   = "arn:awskey"
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

  s3_logging {
    bucket = "test-bucket"
    region = "us-east1"
    prefix = "/"

    // AWS Secret Only
    credentials = cpln_secret.aws.self_link
  }

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
`, olrt.RandomName, c.ResourceName)
}

// HclTwoUniqueLoggings returns a HCL block for a resource.
func (olrt *OrgLoggingResourceTest) HclTwoUniqueLoggings(c OrgLoggingResourceTestCase) string {
	return fmt.Sprintf(`
variable random_name {
  type    = string
  default = "%s"
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
`, olrt.RandomName, c.ResourceName)
}

/*** Resource Test Case ***/

// OrgLoggingResourceTestCase defines a specific resource test case.
type OrgLoggingResourceTestCase struct {
	ProviderTestCase
	LogzioListenerHost string
}
