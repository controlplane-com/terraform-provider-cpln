package cpln

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneServiceAccount_basic performs an acceptance test for the resource.
func TestAccControlPlaneServiceAccount_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewServiceAccountResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "SERVICE-ACCOUNT") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// ServiceAccountResourceTest defines the necessary functionality to test the resource.
type ServiceAccountResourceTest struct {
	Steps []resource.TestStep
}

// NewServiceAccountResourceTest creates a ServiceAccountResourceTest with initialized test cases.
func NewServiceAccountResourceTest() ServiceAccountResourceTest {
	// Create a resource test instance
	resourceTest := ServiceAccountResourceTest{}

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
func (sart *ServiceAccountResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_service_account resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_service_account" {
			continue
		}

		// Retrieve the name for the current resource
		serviceAccountName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of serviceAccount with name: %s", serviceAccountName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		serviceAccount, code, err := TestProvider.client.GetServiceAccount(serviceAccountName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if serviceAccount %s exists: %w", serviceAccountName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if serviceAccount != nil {
			return fmt.Errorf("CheckDestroy failed: serviceAccount %s still exists in the system", *serviceAccount.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_service_account resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case for a Service Account with initial and updated configurations.
func (sart *ServiceAccountResourceTest) NewDefaultScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("sa-default-%s", random)
	resourceName := "new"

	// Build test steps
	initialConfig, initialStep := sart.BuildInitialTestStep(resourceName, name)
	caseUpdate1 := sart.BuildUpdate1TestStep(initialConfig.ProviderTestCase, resourceName, name)
	caseUpdate2 := sart.BuildUpdate2TestStep(initialConfig.ProviderTestCase, resourceName, name)
	caseUpdate3 := sart.BuildUpdate3TestStep(initialConfig.ProviderTestCase, resourceName, name)

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
func (sart *ServiceAccountResourceTest) BuildInitialTestStep(resourceName string, name string) (ServiceAccountResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := ServiceAccountResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "serviceaccount",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_service_account.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "service account default description updated",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: sart.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (sart *ServiceAccountResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase, resourceName string, name string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := ServiceAccountResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: sart.UpdateWithOptionals(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
		),
	}
}

// BuildUpdate2TestStep returns a test step for the update.
func (sart *ServiceAccountResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase, resourceName string, name string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := ServiceAccountResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Create service account key resource cases
	keyCase1ResourceName := fmt.Sprintf("%s-0", resourceName)
	keyCase1 := ServiceAccountKeyResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "serviceaccount",
			ResourceName:    keyCase1ResourceName,
			ResourceAddress: fmt.Sprintf("cpln_service_account_key.%s", keyCase1ResourceName),
		},
		DependsOn:              c.ResourceAddress,
		ServiceAccountNameAttr: c.GetResourceNameAttr(),
		KeyDescription:         "key-01",
	}

	// Construct the service account key resources
	serviceAccountKeyResources := []string{
		sart.ServiceAccountKeyHcl(keyCase1),
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: sart.UpdateWithOptionalsWithKeys(c, strings.Join(serviceAccountKeyResources, "\n\n")),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),

			// Key 01
			resource.TestCheckResourceAttr(keyCase1.ResourceAddress, "description", keyCase1.KeyDescription),
			resource.TestCheckResourceAttr(keyCase1.ResourceAddress, "service_account_name", c.Name),
			resource.TestCheckResourceAttrSet(keyCase1.ResourceAddress, "key"),
		),
	}
}

// BuildUpdate3TestStep returns a test step for the update.
func (sart *ServiceAccountResourceTest) BuildUpdate3TestStep(initialCase ProviderTestCase, resourceName string, name string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := ServiceAccountResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Create service account key resource cases
	keyCase1ResourceName := fmt.Sprintf("%s-0", resourceName)
	keyCase1 := ServiceAccountKeyResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "serviceaccount",
			ResourceName:    keyCase1ResourceName,
			ResourceAddress: fmt.Sprintf("cpln_service_account_key.%s", keyCase1ResourceName),
		},
		DependsOn:              c.ResourceAddress,
		ServiceAccountNameAttr: c.GetResourceNameAttr(),
		KeyDescription:         "key-01",
	}

	keyCase2ResourceName := fmt.Sprintf("%s-1", resourceName)
	keyCase2 := ServiceAccountKeyResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "serviceaccount",
			ResourceName:    keyCase2ResourceName,
			ResourceAddress: fmt.Sprintf("cpln_service_account_key.%s", keyCase2ResourceName),
		},
		DependsOn:              keyCase1.ResourceAddress,
		ServiceAccountNameAttr: c.GetResourceNameAttr(),
		KeyDescription:         "key-02",
	}

	// Construct the service account key resources
	serviceAccountKeyResources := []string{
		sart.ServiceAccountKeyHcl(keyCase1),
		sart.ServiceAccountKeyHcl(keyCase2),
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: sart.UpdateWithOptionalsWithKeys(c, strings.Join(serviceAccountKeyResources, "\n\n")),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),

			// Key 01
			resource.TestCheckResourceAttr(keyCase1.ResourceAddress, "description", keyCase1.KeyDescription),
			resource.TestCheckResourceAttr(keyCase1.ResourceAddress, "service_account_name", c.Name),
			resource.TestCheckResourceAttrSet(keyCase1.ResourceAddress, "key"),

			// Key 02
			resource.TestCheckResourceAttr(keyCase2.ResourceAddress, "description", keyCase2.KeyDescription),
			resource.TestCheckResourceAttr(keyCase2.ResourceAddress, "service_account_name", c.Name),
			resource.TestCheckResourceAttrSet(keyCase2.ResourceAddress, "key"),
		),
	}
}

// Configs //

// RequiredOnly returns a minimal HCL block for the resource using only required fields.
func (sart *ServiceAccountResourceTest) RequiredOnly(c ServiceAccountResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_service_account" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

// UpdateWithOptionals returns a HCL block for the resource using all attributes.
func (sart *ServiceAccountResourceTest) UpdateWithOptionals(c ServiceAccountResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_service_account" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// UpdateWithOptionalsWithKeys returns a HCL block for the resource using all attributes and service account keys.
func (sart *ServiceAccountResourceTest) UpdateWithOptionalsWithKeys(c ServiceAccountResourceTestCase, serviceAccountKeyResources string) string {
	return fmt.Sprintf(`
resource "cpln_service_account" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}

# Service Account Key Resources
%s
`, c.ResourceName, c.Name, c.DescriptionUpdate, serviceAccountKeyResources)
}

// Keys //

// ServiceAccountKeyHcl defines the HCL for the service account key resource.
func (sart *ServiceAccountResourceTest) ServiceAccountKeyHcl(c ServiceAccountKeyResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_service_account_key" "%s" {
  depends_on = [%s]

  service_account_name = %s
  description          = "%s"
}
`, c.ResourceName, c.DependsOn, c.ServiceAccountNameAttr, c.KeyDescription)
}

/*** Resource Test Case ***/

// ServiceAccountResourceTestCase defines a specific resource test case.
type ServiceAccountResourceTestCase struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (sartc *ServiceAccountResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of serviceAccount: %s. Total resources: %d", sartc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[sartc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", sartc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != sartc.Name {
			return fmt.Errorf("resource ID %s does not match expected serviceAccount name %s", rs.Primary.ID, sartc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteServiceAccount, _, err := TestProvider.client.GetServiceAccount(sartc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving serviceAccount from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteServiceAccount.Name != sartc.Name {
			return fmt.Errorf("mismatch in serviceAccount name: expected %s, got %s", sartc.Name, *remoteServiceAccount.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("serviceAccount %s verified successfully in both state and external system.", sartc.Name))
		return nil
	}
}

/*** Service Account Key - Resource Test Case ***/

// ServiceAccountKeyResourceTestCase defines a specific resource test case.
type ServiceAccountKeyResourceTestCase struct {
	ProviderTestCase
	DependsOn              string
	ServiceAccountNameAttr string
	KeyDescription         string
}
