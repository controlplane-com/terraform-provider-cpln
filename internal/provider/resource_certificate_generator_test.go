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

// TestAccControlPlaneCertificateGenerator_basic performs an acceptance test for the resource.
func TestAccControlPlaneCertificateGenerator_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewCertificateGeneratorResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "CERTIFICATE_GENERATOR") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// CertificateGeneratorResourceTest defines the necessary functionality to test the resource.
type CertificateGeneratorResourceTest struct {
	Steps []resource.TestStep
}

// NewCertificateGeneratorResourceTest creates a CertificateGeneratorResourceTest with initialized test cases.
func NewCertificateGeneratorResourceTest() CertificateGeneratorResourceTest {
	// Create a resource test instance
	resourceTest := CertificateGeneratorResourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewAcmeLifecycleScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (cgrt *CertificateGeneratorResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_certificate_generator resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_certificate_generator" {
			continue
		}

		// Retrieve the name for the current resource
		certificateGeneratorName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of certificate generator with name: %s", certificateGeneratorName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		certificateGenerator, code, err := TestProvider.client.GetCertificateGenerator(certificateGeneratorName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if certificate generator %s exists: %w", certificateGeneratorName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if certificateGenerator != nil {
			return fmt.Errorf("CheckDestroy failed: certificate generator %s still exists in the system", *certificateGenerator.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_certificate_generator resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewAcmeLifecycleScenario walks the acme challenge through http01, every dns01 provider, and back.
func (cgrt *CertificateGeneratorResourceTest) NewAcmeLifecycleScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-certificate-generator-%s", random)
	resourceName := "new"

	// Build test steps
	initialConfig, http01RequiredOnlyStep := cgrt.BuildHttp01RequiredOnlyTestStep(resourceName, name)
	http01AllStep := cgrt.BuildHttp01AllTestStep(initialConfig.ProviderTestCase)
	dns01CplnStep := cgrt.BuildDns01CplnTestStep(initialConfig.ProviderTestCase)
	dns01Route53Step := cgrt.BuildDns01Route53TestStep(initialConfig.ProviderTestCase)
	dns01CloudDnsStep := cgrt.BuildDns01CloudDnsTestStep(initialConfig.ProviderTestCase)
	dns01AzureDnsStep := cgrt.BuildDns01AzureDnsTestStep(initialConfig.ProviderTestCase)
	dns01CloudflareStep := cgrt.BuildDns01CloudflareTestStep(initialConfig.ProviderTestCase)
	dns01Route53RequiredOnlyStep := cgrt.BuildDns01Route53RequiredOnlyTestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		http01RequiredOnlyStep,
		// Import State
		{
			ResourceName: initialConfig.ResourceAddress,
			ImportState:  true,
		},
		// Update & Read
		http01AllStep,
		// Switch the challenge type, which recreates the resource
		dns01CplnStep,
		dns01Route53Step,
		dns01CloudDnsStep,
		dns01AzureDnsStep,
		dns01CloudflareStep,
		dns01Route53RequiredOnlyStep,
		// Revert the resource to its initial state, which recreates it again
		http01RequiredOnlyStep,
	}
}

// Test Cases //

// BuildHttp01RequiredOnlyTestStep returns a default initial test step and its associated test case for the resource.
func (cgrt *CertificateGeneratorResourceTest) BuildHttp01RequiredOnlyTestStep(resourceName string, name string) (CertificateGeneratorResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "certificategenerator",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_certificate_generator.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "certificate generator description updated",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: cgrt.Http01RequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server":                   "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":                    nil,
				"external_account_binding": nil,
				"challenge": map[string]interface{}{
					"http01": map[string]interface{}{},
					"dns01":  nil,
				},
			}),
			c.TestCheckResourceAttr("status.#", "1"),
		),
	}
}

// BuildHttp01AllTestStep returns a test step that sets every optional field on an http01 certificate generator.
func (cgrt *CertificateGeneratorResourceTest) BuildHttp01AllTestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: cgrt.Http01AllHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server": "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":  "ops@example.com",
				"external_account_binding": map[string]interface{}{
					"key_id":      "terraform-test-key-id",
					"secret_link": fmt.Sprintf("/org/%s/secret/%s-eab", OrgName, c.Name),
				},
				"challenge": map[string]interface{}{
					"http01": map[string]interface{}{},
					"dns01":  nil,
				},
			}),
		),
	}
}

// BuildDns01CplnTestStep returns a test step that switches to a dns01 challenge delegated to Control Plane.
func (cgrt *CertificateGeneratorResourceTest) BuildDns01CplnTestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: cgrt.Dns01CplnHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server": "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":  "ops@example.com",
				"external_account_binding": map[string]interface{}{
					"key_id":      "terraform-test-key-id",
					"secret_link": fmt.Sprintf("/org/%s/secret/%s-eab", OrgName, c.Name),
				},
				"challenge": map[string]interface{}{
					"http01": nil,
					"dns01": map[string]interface{}{
						"provider":   "cpln",
						"route53":    nil,
						"cloud_dns":  nil,
						"azure_dns":  nil,
						"cloudflare": nil,
					},
				},
			}),
		),
	}
}

// BuildDns01Route53TestStep returns a test step that solves dns01 challenges in Route53 with every field set.
func (cgrt *CertificateGeneratorResourceTest) BuildDns01Route53TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: cgrt.Dns01Route53Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server": "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":  "ops@example.com",
				"external_account_binding": map[string]interface{}{
					"key_id":      "terraform-test-key-id",
					"secret_link": fmt.Sprintf("/org/%s/secret/%s-eab", OrgName, c.Name),
				},
				"challenge": map[string]interface{}{
					"http01": nil,
					"dns01": map[string]interface{}{
						"provider": "route53",
						"route53": map[string]interface{}{
							"region":         "us-east-1",
							"hosted_zone_id": "Z0123456789ABCDEFGHIJ",
							"secret_link":    fmt.Sprintf("/org/%s/secret/%s-aws", OrgName, c.Name),
						},
						"cloud_dns":  nil,
						"azure_dns":  nil,
						"cloudflare": nil,
					},
				},
			}),
		),
	}
}

// BuildDns01CloudDnsTestStep returns a test step that solves dns01 challenges in Google Cloud DNS.
func (cgrt *CertificateGeneratorResourceTest) BuildDns01CloudDnsTestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: cgrt.Dns01CloudDnsHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server": "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":  "ops@example.com",
				"external_account_binding": map[string]interface{}{
					"key_id":      "terraform-test-key-id",
					"secret_link": fmt.Sprintf("/org/%s/secret/%s-eab", OrgName, c.Name),
				},
				"challenge": map[string]interface{}{
					"http01": nil,
					"dns01": map[string]interface{}{
						"provider": "cloudDns",
						"route53":  nil,
						"cloud_dns": map[string]interface{}{
							"project":     "cpln12345",
							"secret_link": fmt.Sprintf("/org/%s/secret/%s-gcp", OrgName, c.Name),
						},
						"azure_dns":  nil,
						"cloudflare": nil,
					},
				},
			}),
		),
	}
}

// BuildDns01AzureDnsTestStep returns a test step that solves dns01 challenges in Azure DNS with every field set.
func (cgrt *CertificateGeneratorResourceTest) BuildDns01AzureDnsTestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: cgrt.Dns01AzureDnsHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server": "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":  "ops@example.com",
				"external_account_binding": map[string]interface{}{
					"key_id":      "terraform-test-key-id",
					"secret_link": fmt.Sprintf("/org/%s/secret/%s-eab", OrgName, c.Name),
				},
				"challenge": map[string]interface{}{
					"http01": nil,
					"dns01": map[string]interface{}{
						"provider":  "azureDns",
						"route53":   nil,
						"cloud_dns": nil,
						"azure_dns": map[string]interface{}{
							"resource_group_name": "dns-resource-group",
							"hosted_zone_name":    "example.com",
							"secret_link":         fmt.Sprintf("/org/%s/secret/%s-azure-sdk", OrgName, c.Name),
						},
						"cloudflare": nil,
					},
				},
			}),
		),
	}
}

// BuildDns01CloudflareTestStep returns a test step that solves dns01 challenges in Cloudflare using a short secret link.
func (cgrt *CertificateGeneratorResourceTest) BuildDns01CloudflareTestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: cgrt.Dns01CloudflareHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server": "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":  "ops@example.com",
				"external_account_binding": map[string]interface{}{
					"key_id":      "terraform-test-key-id",
					"secret_link": fmt.Sprintf("/org/%s/secret/%s-eab", OrgName, c.Name),
				},
				"challenge": map[string]interface{}{
					"http01": nil,
					"dns01": map[string]interface{}{
						"provider":  "cloudflare",
						"route53":   nil,
						"cloud_dns": nil,
						"azure_dns": nil,
						"cloudflare": map[string]interface{}{
							"secret_link": fmt.Sprintf("//secret/%s-cloudflare", c.Name),
						},
					},
				},
			}),
		),
	}
}

// BuildDns01Route53RequiredOnlyTestStep returns a test step that solves dns01 challenges in Route53 with only the required fields.
func (cgrt *CertificateGeneratorResourceTest) BuildDns01Route53RequiredOnlyTestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CertificateGeneratorResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: cgrt.Dns01Route53RequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckObjectAttr("acme", map[string]interface{}{
				"server": "https://acme-staging-v02.api.letsencrypt.org/directory",
				"email":  "ops@example.com",
				"external_account_binding": map[string]interface{}{
					"key_id":      "terraform-test-key-id",
					"secret_link": fmt.Sprintf("/org/%s/secret/%s-eab", OrgName, c.Name),
				},
				"challenge": map[string]interface{}{
					"http01": nil,
					"dns01": map[string]interface{}{
						"provider": "route53",
						"route53": map[string]interface{}{
							"region":         "us-east-1",
							"hosted_zone_id": nil,
							"secret_link":    fmt.Sprintf("/org/%s/secret/%s-aws", OrgName, c.Name),
						},
						"cloud_dns":  nil,
						"azure_dns":  nil,
						"cloudflare": nil,
					},
				},
			}),
		),
	}
}

// Configs //

// Http01RequiredOnlyHcl returns a certificate generator HCL with only the required fields and an http01 challenge.
func (cgrt *CertificateGeneratorResourceTest) Http01RequiredOnlyHcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name = "%s"

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"

    challenge = {
      http01 = {}
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name)
}

// Http01AllHcl returns a certificate generator HCL with every optional field and an http01 challenge.
func (cgrt *CertificateGeneratorResourceTest) Http01AllHcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "terraform-test-key-id"
      secret_link = cpln_secret.eab.self_link
    }

    challenge = {
      http01 = {}
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// Dns01CplnHcl returns a certificate generator HCL with a dns01 challenge delegated to Control Plane.
func (cgrt *CertificateGeneratorResourceTest) Dns01CplnHcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "terraform-test-key-id"
      secret_link = cpln_secret.eab.self_link
    }

    challenge = {
      dns01 = {}
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// Dns01Route53Hcl returns a certificate generator HCL with a dns01 challenge solved in Route53 with every field set.
func (cgrt *CertificateGeneratorResourceTest) Dns01Route53Hcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "terraform-test-key-id"
      secret_link = cpln_secret.eab.self_link
    }

    challenge = {
      dns01 = {
        provider = "route53"

        route53 = {
          region         = "us-east-1"
          hosted_zone_id = "Z0123456789ABCDEFGHIJ"
          secret_link    = cpln_secret.aws.self_link
        }
      }
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// Dns01CloudDnsHcl returns a certificate generator HCL with a dns01 challenge solved in Google Cloud DNS.
func (cgrt *CertificateGeneratorResourceTest) Dns01CloudDnsHcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "terraform-test-key-id"
      secret_link = cpln_secret.eab.self_link
    }

    challenge = {
      dns01 = {
        provider = "cloudDns"

        cloud_dns = {
          project     = "cpln12345"
          secret_link = cpln_secret.gcp.self_link
        }
      }
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// Dns01AzureDnsHcl returns a certificate generator HCL with a dns01 challenge solved in Azure DNS with every field set.
func (cgrt *CertificateGeneratorResourceTest) Dns01AzureDnsHcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "terraform-test-key-id"
      secret_link = cpln_secret.eab.self_link
    }

    challenge = {
      dns01 = {
        provider = "azureDns"

        azure_dns = {
          resource_group_name = "dns-resource-group"
          hosted_zone_name    = "example.com"
          secret_link         = cpln_secret.azure_sdk.self_link
        }
      }
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// Dns01CloudflareHcl returns a certificate generator HCL with a dns01 challenge solved in Cloudflare, using a short secret link.
func (cgrt *CertificateGeneratorResourceTest) Dns01CloudflareHcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "terraform-test-key-id"
      secret_link = cpln_secret.eab.self_link
    }

    challenge = {
      dns01 = {
        provider = "cloudflare"

        cloudflare = {
          secret_link = "//secret/${cpln_secret.cloudflare.name}"
        }
      }
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// Dns01Route53RequiredOnlyHcl returns a certificate generator HCL with a dns01 challenge solved in Route53 with only the required fields.
func (cgrt *CertificateGeneratorResourceTest) Dns01Route53RequiredOnlyHcl(c CertificateGeneratorResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s-aws"

  aws {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

resource "cpln_secret" "gcp" {
  name = "%s-gcp"
  gcp  = jsonencode({
    type                        = "service_account"
    project_id                  = "cpln12345"
    private_key_id              = "pvt_key"
    private_key                 = "key"
    client_email                = "support@cpln.io"
    client_id                   = "12744"
    auth_uri                    = "cloud.google.com"
    token_uri                   = "token.cloud.google.com"
    auth_provider_x509_cert_url = "cert.google.com"
    client_x509_cert_url        = "cert.google.com"
  })
}

resource "cpln_secret" "azure_sdk" {
  name      = "%s-azure-sdk"
  azure_sdk = jsonencode({
    subscriptionId = "2cd8674e-4f89-4a1f-b420-7a1361b46ef7"
    tenantId       = "292f5674-c8b0-488b-9ff8-6d30d77f38d9"
    clientId       = "649846ce-d862-49d5-a5eb-7d5aad90f54e"
    clientSecret   = "cpln"
  })
}

resource "cpln_secret" "cloudflare" {
  name = "%s-cloudflare"

  opaque {
    payload = "cloudflare-api-token"
  }
}

resource "cpln_secret" "eab" {
  name = "%s-eab"

  opaque {
    payload  = "ZmFrZS1obWFjLWtleQ"
    encoding = "plain"
  }
}

resource "cpln_certificate_generator" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  acme = {
    server = "https://acme-staging-v02.api.letsencrypt.org/directory"
    email  = "ops@example.com"

    external_account_binding = {
      key_id      = "terraform-test-key-id"
      secret_link = cpln_secret.eab.self_link
    }

    challenge = {
      dns01 = {
        provider = "route53"

        route53 = {
          region      = "us-east-1"
          secret_link = cpln_secret.aws.self_link
        }
      }
    }
  }
}
`, c.Name, c.Name, c.Name, c.Name, c.Name, c.ResourceName, c.Name, c.DescriptionUpdate)
}

/*** Resource Test Case ***/

// CertificateGeneratorResourceTestCase defines a specific resource test case.
type CertificateGeneratorResourceTestCase struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (cgrtc *CertificateGeneratorResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of certificate generator: %s. Total resources: %d", cgrtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[cgrtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", cgrtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != cgrtc.Name {
			return fmt.Errorf("resource ID %s does not match expected certificate generator name %s", rs.Primary.ID, cgrtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteCertificateGenerator, _, err := TestProvider.client.GetCertificateGenerator(cgrtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving certificate generator from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteCertificateGenerator.Name != cgrtc.Name {
			return fmt.Errorf("mismatch in certificate generator name: expected %s, got %s", cgrtc.Name, *remoteCertificateGenerator.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("certificate generator %s verified successfully in both state and external system.", cgrtc.Name))
		return nil
	}
}
