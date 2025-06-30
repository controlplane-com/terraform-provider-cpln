package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceOrg_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceOrg_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewOrgDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_ORG") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// OrgDataSourceTest defines the necessary functionality to test the data source.
type OrgDataSourceTest struct {
	Steps []resource.TestStep
}

// NewOrgDataSourceTest creates a OrgDataSourceTest with initialized test cases.
func NewOrgDataSourceTest() OrgDataSourceTest {
	// Create a data source test instance
	dataSourceTest := OrgDataSourceTest{}

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
func (odst *OrgDataSourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "new"

	// Build test steps
	_, initialStep := odst.BuildDefaultTestStep(dataSourceName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildDefaultTestStep returns a default initial test step and its associated test case for the data source.
func (odst *OrgDataSourceTest) BuildDefaultTestStep(dataSourceName string) (OrgDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := OrgDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "org",
			ResourceName:    dataSourceName,
			Name:            OrgName,
			Description:     OrgName,
			ResourceAddress: fmt.Sprintf("data.cpln_org.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: odst.DefaultHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("observability", []map[string]interface{}{
				{
					"logs_retention_days":    "30",
					"metrics_retention_days": "30",
					"traces_retention_days":  "30",
					"default_alert_emails":   []string{},
				},
			}),
		),
	}
}

// Configs //

// DefaultHcl returns a data source HCL.
func (odst *OrgDataSourceTest) DefaultHcl(c OrgDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_org" "%s" {}
`, c.ResourceName)
}

/*** Data Source Test Case ***/

// OrgDataSourceTestCase defines a specific data source test case.
type OrgDataSourceTestCase struct {
	ProviderTestCase
}
