package cpln

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneLocation_basic performs an acceptance test for the resource.
func TestAccControlPlaneLocation_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewLocationResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "LOCATION") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// LocationResourceTest defines the necessary functionality to test the resource.
type LocationResourceTest struct {
	Steps []resource.TestStep
}

// NewLocationResourceTest creates a LocationResourceTest with initialized test cases.
func NewLocationResourceTest() LocationResourceTest {
	// Create a resource test instance
	resourceTest := LocationResourceTest{}

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
func (lrt *LocationResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_location resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_location" {
			continue
		}

		// Retrieve the name for the current resource
		locationName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of location with name: %s", locationName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		location, code, err := TestProvider.client.GetLocation(locationName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if location %s exists: %w", locationName, err)
		}

		// Ensure the Location has a Spec field
		if location.Spec == nil {
			return fmt.Errorf("Location does not have Spec. Name: %s", *location.Name)
		}

		// Ensure the Spec Enabled field is present
		if location.Spec.Enabled == nil {
			return fmt.Errorf("Location enabled is nil. Name: %s", *location.Name)
		}

		// Verify the Enabled flag has reverted to true
		if !*location.Spec.Enabled {
			return fmt.Errorf("Location enabled wasn't reverted to the default value: true. Name: %s", *location.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_location resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case with initial and updated configurations.
func (lrt *LocationResourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	name := "aws-eu-central-1"
	resourceName := "new"

	// Build test steps
	initialConfig, initialStep := lrt.BuildInitialTestStep(resourceName, name)
	caseUpdate1 := lrt.BuildUpdate1TestStep(initialConfig.ProviderTestCase)

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
func (lrt *LocationResourceTest) BuildInitialTestStep(resourceName string, name string) (LocationResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := LocationResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "location",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_location.%s", resourceName),
			Name:            name,
			Description:     "AWS, Europe (Frankfurt)",
		},
		Enabled: strconv.FormatBool(true),
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: lrt.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "enabled", c.Enabled),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (lrt *LocationResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := LocationResourceTestCase{
		ProviderTestCase: initialCase,
		Enabled:          strconv.FormatBool(false),
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: lrt.UpdateWithMinimalOptionals(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "enabled", c.Enabled),
		),
	}
}

// Configs //

// RequiredOnly returns a minimal HCL block for a resource using only required fields.
func (lrt *LocationResourceTest) RequiredOnly(c LocationResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_location" "%s" {
  name    = "%s"
	enabled = %s
}
`, c.ResourceName, c.Name, c.Enabled)
}

// UpdateWithMinimalOptionals
func (lrt *LocationResourceTest) UpdateWithMinimalOptionals(c LocationResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_location" "%s" {
  name = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

	enabled = %s
}
`, c.ResourceName, c.Name, c.Enabled)
}

/*** Resource Test Case ***/

// LocationResourceTestCase defines a specific resource test case.
type LocationResourceTestCase struct {
	ProviderTestCase
	Enabled string
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (lrtc *LocationResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of location: %s. Total resources: %d", lrtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[lrtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", lrtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != lrtc.Name {
			return fmt.Errorf("resource ID %s does not match expected location name %s", rs.Primary.ID, lrtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteLocation, _, err := TestProvider.client.GetLocation(lrtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving location from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteLocation.Name != lrtc.Name {
			return fmt.Errorf("mismatch in location name: expected %s, got %s", lrtc.Name, *remoteLocation.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("location %s verified successfully in both state and external system.", lrtc.Name))
		return nil
	}
}
