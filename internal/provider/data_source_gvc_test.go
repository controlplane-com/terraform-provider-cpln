package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceGvc_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceGvc_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewGvcDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_GVC") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// GvcDataSourceTest defines the necessary functionality to test the data source.
type GvcDataSourceTest struct {
	Steps []resource.TestStep
}

// NewGvcDataSourceTest creates a GvcDataSourceTest with initialized test cases.
func NewGvcDataSourceTest() GvcDataSourceTest {
	// Create a data source test instance
	dataSourceTest := GvcDataSourceTest{}

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
func (gdst *GvcDataSourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "new"
	name := "default-gvc"

	// Build test steps
	_, initialStep := gdst.BuildDefaultTestStep(dataSourceName, name)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildDefaultTestStep returns a default initial test step and its associated test case for the data source.
func (gdst *GvcDataSourceTest) BuildDefaultTestStep(dataSourceName string, name string) (GvcDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := GvcDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "gvc",
			ResourceName:    dataSourceName,
			Name:            name,
			Description:     name,
			ResourceAddress: fmt.Sprintf("data.cpln_gvc.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: gdst.DefaultHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "1"),
		),
	}
}

// Configs //

// DefaultHcl returns a data source HCL.
func (gdst *GvcDataSourceTest) DefaultHcl(c GvcDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_gvc" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

/*** Data Source Test Case ***/

// GvcDataSourceTestCase defines a specific data source test case.
type GvcDataSourceTestCase struct {
	ProviderTestCase
}
