package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceSecret_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceSecret_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewSecretDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_SECRET") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// SecretDataSourceTest defines the necessary functionality to test the data source.
type SecretDataSourceTest struct {
	Steps []resource.TestStep
}

// NewSecretDataSourceTest creates a SecretDataSourceTest with initialized test cases.
func NewSecretDataSourceTest() SecretDataSourceTest {
	// Create a data source test instance
	dataSourceTest := SecretDataSourceTest{}

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
func (sdst *SecretDataSourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "new"
	name := "test-secret-opaque"

	// Build test steps
	_, initialStep := sdst.BuildDefaultTestStep(dataSourceName, name)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildDefaultTestStep returns a default initial test step and its associated test case for the data source.
func (sdst *SecretDataSourceTest) BuildDefaultTestStep(dataSourceName string, name string) (SecretDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := SecretDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "secret",
			ResourceName:    dataSourceName,
			Name:            name,
			Description:     name,
			ResourceAddress: fmt.Sprintf("data.cpln_secret.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: sdst.DefaultHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "1"),
			c.TestCheckNestedBlocks("opaque", []map[string]interface{}{
				{
					"encoding": "base64",
					"payload":  "VGhpcyBpcyBhbiBvcGFxdWUgaW4gYmFzZTY0",
				},
			}),
		),
	}
}

// Configs //

// DefaultHcl returns a data source HCL.
func (sdst *SecretDataSourceTest) DefaultHcl(c SecretDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_secret" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

/*** Data Source Test Case ***/

// SecretDataSourceTestCase defines a specific data source test case.
type SecretDataSourceTestCase struct {
	ProviderTestCase
}
