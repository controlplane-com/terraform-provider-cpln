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

// TestAccControlPlaneIpSet_basic performs an acceptance test for the resource.
func TestAccControlPlaneIpSet_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewIpSetResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "IPSET") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// IpSetResourceTest defines the necessary functionality to test the resource.
type IpSetResourceTest struct {
	Steps []resource.TestStep
}

// NewIpSetResourceTest creates a IpSetResourceTest with initialized test cases.
func NewIpSetResourceTest() IpSetResourceTest {
	// Create a resource test instance
	resourceTest := IpSetResourceTest{}

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
func (isrt *IpSetResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_ipset resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_ipset" {
			continue
		}

		// Retrieve the name for the current resource
		ipsetName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of ipset with name: %s", ipsetName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		ipset, code, err := TestProvider.client.GetIpSet(ipsetName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if ipset %s exists: %w", ipsetName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if ipset != nil {
			return fmt.Errorf("CheckDestroy failed: ipset %s still exists in the system", *ipset.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_ipset resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case for a group using JMESPATH with initial and updated configurations.
func (isrt *IpSetResourceTest) NewDefaultScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("ipset-default-%s", random)
	resourceName := "new"

	// Build test steps
	initialConfig, initialStep := isrt.BuildInitialTestStep(resourceName, name)
	caseUpdate1 := isrt.BuildUpdate1TestStep(initialConfig.ProviderTestCase, resourceName)

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

// BuildInitialTestStep returns a default initial test step and its associated test case for the resource.
func (isrt *IpSetResourceTest) BuildInitialTestStep(resourceName string, name string) (IpSetResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := IpSetResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "ipset",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_ipset.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "ipset default description updated",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: isrt.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (isrt *IpSetResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase, resourceName string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := IpSetResourceTestCase{
		ProviderTestCase: initialCase,
		Link:             fmt.Sprintf("/org/%s/gvc/default-gvc", OrgName),
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: isrt.UpdateWithMinimalOptionals(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "link", c.Link),
			c.TestCheckNestedBlocks("location", []map[string]interface{}{
				{
					"name":             fmt.Sprintf("/org/%s/location/aws-eu-central-1", OrgName),
					"retention_policy": "keep",
				},
			}),
		),
	}
}

// Configs //

// RequiredOnly returns a minimal HCL block for a resource using only required fields.
func (isrt *IpSetResourceTest) RequiredOnly(c IpSetResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_ipset" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

// UpdateWithMinimalOptionals constructs an HCL configuration for an IP set resource including minimal optionals like tags, link, and location settings
func (isrt *IpSetResourceTest) UpdateWithMinimalOptionals(c IpSetResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_ipset" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

	link = "%s"

	location {
	  name             = "/org/%s/location/aws-eu-central-1"
		retention_policy = "keep"
	}
}
`, c.ResourceName, c.Name, c.DescriptionUpdate, c.Link, OrgName)
}

/*** Resource Test Case ***/

// IpSetResourceTestCase defines a specific resource test case.
type IpSetResourceTestCase struct {
	ProviderTestCase
	Link string
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (isrtc *IpSetResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of ipset: %s. Total resources: %d", isrtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[isrtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", isrtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != isrtc.Name {
			return fmt.Errorf("resource ID %s does not match expected ipset name %s", rs.Primary.ID, isrtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteIpSet, _, err := TestProvider.client.GetIpSet(isrtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving ipset from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteIpSet.Name != isrtc.Name {
			return fmt.Errorf("mismatch in ipset name: expected %s, got %s", isrtc.Name, *remoteIpSet.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("ipset %s verified successfully in both state and external system.", isrtc.Name))
		return nil
	}
}
