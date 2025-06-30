package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceCloudAccount_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceCloudAccount_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewCloudAccountDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_CLOUD_ACCOUNT") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// CloudAccountDataSourceTest defines the necessary functionality to test the data source.
type CloudAccountDataSourceTest struct {
	Steps []resource.TestStep
}

// NewCloudAccountDataSourceTest creates a CloudAccountDataSourceTest with initialized test cases.
func NewCloudAccountDataSourceTest() CloudAccountDataSourceTest {
	// Create a data source test instance
	dataSourceTest := CloudAccountDataSourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, dataSourceTest.NewDefaultScenario()...)

	// Set the cases for the data source test
	dataSourceTest.Steps = steps

	// Return the data source test
	return dataSourceTest
}

// Test Scenarios //

// NewDefaultScenario creates a test case with initial and updated configurations.
func (cadst *CloudAccountDataSourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "new"

	// Build test steps
	_, initialStep := cadst.BuildDefaultTestStep(dataSourceName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildDefaultTestStep returns a default initial test step and its associated test case for the data source.
func (cadst *CloudAccountDataSourceTest) BuildDefaultTestStep(dataSourceName string) (CloudAccountDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := CloudAccountDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "cloudaccount",
			ResourceName:    dataSourceName,
			ResourceAddress: fmt.Sprintf("data.cpln_cloud_account.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: cadst.DefaultHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", "static-cloud-account"),
			c.TestCheckSetAttr("aws_identifiers", CloudAccountIdentifiers),
		),
	}
}

// Configs //

// DefaultHcl returns a data source HCL.
func (cadst *CloudAccountDataSourceTest) DefaultHcl(c CloudAccountDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_cloud_account" "%s" {}
`, c.ResourceName)
}

/*** Data Source Test Case ***/

// CloudAccountDataSourceTestCase defines a specific data source test case.
type CloudAccountDataSourceTestCase struct {
	ProviderTestCase
}
