package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceLocation_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceLocation_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewLocationDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_LOCATION") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// LocationDataSourceTest defines the necessary functionality to test the data source.
type LocationDataSourceTest struct {
	Steps []resource.TestStep
}

// NewLocationDataSourceTest creates a LocationDataSourceTest with initialized test cases.
func NewLocationDataSourceTest() LocationDataSourceTest {
	// Create a data source test instance
	dataSourceTest := LocationDataSourceTest{}

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
func (ldst *LocationDataSourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "new"
	name := "aws-eu-central-1"

	// Build test steps
	_, initialStep := ldst.BuildDefaultTestStep(dataSourceName, name)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildDefaultTestStep returns a default initial test step and its associated test case for the data source.
func (ldst *LocationDataSourceTest) BuildDefaultTestStep(dataSourceName string, name string) (LocationDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := LocationDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "location",
			ResourceName:    dataSourceName,
			Name:            name,
			Description:     "AWS, Europe (Frankfurt)",
			ResourceAddress: fmt.Sprintf("data.cpln_location.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: ldst.DefaultHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckResourceAttr("cloud_provider", "aws"),
			c.TestCheckResourceAttr("region", "eu-central-1"),
			c.TestCheckResourceAttr("enabled", "true"),
		),
	}
}

// Configs //

// DefaultHcl returns a data source HCL.
func (ldst *LocationDataSourceTest) DefaultHcl(c LocationDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_location" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

/*** Data Source Test Case ***/

// LocationDataSourceTestCase defines a specific data source test case.
type LocationDataSourceTestCase struct {
	ProviderTestCase
}
