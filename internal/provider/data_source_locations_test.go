package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceLocations_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceLocations_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewLocationsDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_LOCATIONS") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// LocationsDataSourceTest defines the necessary functionality to test the data source.
type LocationsDataSourceTest struct {
	Steps []resource.TestStep
}

// NewLocationsDataSourceTest creates a LocationsDataSourceTest with initialized test cases.
func NewLocationsDataSourceTest() LocationsDataSourceTest {
	// Create a data source test instance
	dataSourceTest := LocationsDataSourceTest{}

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
func (ldst *LocationsDataSourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "new"

	// Build test steps
	_, initialStep := ldst.BuildDefaultTestStep(dataSourceName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildDefaultTestStep returns a test case for the data source.
func (ldst *LocationsDataSourceTest) BuildDefaultTestStep(dataSourceName string) (LocationsDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := LocationsDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "location",
			ResourceName:    dataSourceName,
			ResourceAddress: fmt.Sprintf("data.cpln_locations.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: ldst.DefaultHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			resource.TestCheckResourceAttrSet(c.ResourceAddress, "locations.#"),
		),
	}
}

// Configs //

// DefaultHcl returns a data source HCL.
func (ldst *LocationsDataSourceTest) DefaultHcl(c LocationsDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_locations" "%s" {}
`, c.ResourceName)
}

/*** Data Source Test Case ***/

// LocationsDataSourceTestCase defines a specific data source test case.
type LocationsDataSourceTestCase struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (ldstc *LocationsDataSourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[ldstc.ResourceAddress]

		// Return an error if the resource is not found in state
		if !ok {
			return fmt.Errorf("Can't find locations data source: %s", ldstc.ResourceAddress)
		}

		// Ensure the Terraform state has set the resource ID
		if rs.Primary.ID == "" {
			return fmt.Errorf("locations data source ID not set")
		}

		// Indicate successful existence check
		return nil
	}
}
