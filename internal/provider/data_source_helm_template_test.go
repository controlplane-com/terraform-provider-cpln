package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneHelmTemplate_basic performs an acceptance test for the data source.
func TestAccControlPlaneHelmTemplate_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewHelmTemplateDataSourceTest()

	// Run the acceptance test case for the data source, covering read functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "HELM") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// HelmTemplateDataSourceTest defines the necessary functionality to test the data source.
type HelmTemplateDataSourceTest struct {
	Steps []resource.TestStep
}

// NewHelmTemplateDataSourceTest creates a HelmTemplateDataSourceTest with initialized test cases.
func NewHelmTemplateDataSourceTest() HelmTemplateDataSourceTest {
	// Create a data source test instance
	dataSourceTest := HelmTemplateDataSourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, dataSourceTest.NewRequiredOnlyScenario()...)
	steps = append(steps, dataSourceTest.NewWithSetValuesScenario()...)

	// Set the cases for the data source test
	dataSourceTest.Steps = steps

	// Return the data source test
	return dataSourceTest
}

// Test Scenarios //

// NewRequiredOnlyScenario creates a test case for helm template rendering using only required fields.
func (htt *HelmTemplateDataSourceTest) NewRequiredOnlyScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("helm-tmpl-req-%s", random)
	gvcName := fmt.Sprintf("helm-tmpl-req-gvc-%s", random)
	dataSourceName := "required-only"

	// Build test steps
	_, readStep := htt.BuildRequiredOnlyReadTestStep(dataSourceName, name, gvcName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read & Verify
		readStep,
	}
}

// NewWithSetValuesScenario creates a test case for helm template rendering with set, set_string, and multiple values.
func (htt *HelmTemplateDataSourceTest) NewWithSetValuesScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("helm-tmpl-set-%s", random)
	gvcName := fmt.Sprintf("helm-tmpl-set-gvc-%s", random)
	dataSourceName := "with-set-values"

	// Build test steps
	_, readStep := htt.BuildWithSetValuesReadTestStep(dataSourceName, name, gvcName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read & Verify
		readStep,
	}
}

// Test Cases //

// BuildRequiredOnlyReadTestStep returns a read test step using only required fields.
func (htt *HelmTemplateDataSourceTest) BuildRequiredOnlyReadTestStep(dataSourceName string, name string, gvcName string) (HelmTemplateDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := HelmTemplateDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "helm_template",
			ResourceName:    dataSourceName,
			ResourceAddress: fmt.Sprintf("data.cpln_helm_template.%s", dataSourceName),
			Name:            name,
		},
		GvcName:         gvcName,
		GvcResourceName: "helm_tmpl_gvc_required_only",
		Chart:           "../../testdata/helm/sample-chart",
		ValuesFile:      TestdataAbsPath("../../testdata/helm/values/initial.yaml"),
	}

	// Initialize and return the read test step
	return c, resource.TestStep{
		Config: htt.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("gvc", c.GvcName),
			c.TestCheckResourceAttr("chart", c.Chart),
			c.TestCheckResourceAttr("values.#", "1"),
			c.TestCheckResourceAttrSet("manifest"),
		),
	}
}

// BuildWithSetValuesReadTestStep returns a read test step with set, set_string, and multiple values files.
func (htt *HelmTemplateDataSourceTest) BuildWithSetValuesReadTestStep(dataSourceName string, name string, gvcName string) (HelmTemplateDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	overrideValuesFile := TestdataAbsPath("../../testdata/helm/values/override.yaml")
	c := HelmTemplateDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "helm_template",
			ResourceName:    dataSourceName,
			ResourceAddress: fmt.Sprintf("data.cpln_helm_template.%s", dataSourceName),
			Name:            name,
		},
		GvcName:         gvcName,
		GvcResourceName: "helm_tmpl_gvc_with_set_values",
		Chart:           "../../testdata/helm/sample-chart",
		ValuesFile:      TestdataAbsPath("../../testdata/helm/values/initial.yaml"),
	}

	// Initialize and return the read test step
	return c, resource.TestStep{
		Config: htt.WithSetValuesAndMultipleValues(c, overrideValuesFile),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("gvc", c.GvcName),
			c.TestCheckResourceAttr("chart", c.Chart),
			c.TestCheckResourceAttr("values.#", "2"),
			c.TestCheckResourceAttr("set.secret.data", "override-value"),
			c.TestCheckResourceAttr("set_string.secret.name", fmt.Sprintf("%s-secret", name)),
			c.TestCheckResourceAttrSet("manifest"),
		),
	}
}

// Configs //

// RequiredOnly returns a minimal HCL block for the data source using only required fields.
func (htt *HelmTemplateDataSourceTest) RequiredOnly(c HelmTemplateDataSourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name = "%s"
}

data "cpln_helm_template" "%s" {
  name  = "%s"
  gvc   = cpln_gvc.%s.name
  chart = "%s"

  values = [file("%s")]

  set_string = {
    "secret.name" = "%s-secret"
  }
}
`, c.GvcResourceName, c.GvcName, c.ResourceName, c.Name, c.GvcResourceName, c.Chart, c.ValuesFile, c.Name)
}

// WithSetValuesAndMultipleValues returns an HCL block with set, set_string, and two values files.
func (htt *HelmTemplateDataSourceTest) WithSetValuesAndMultipleValues(c HelmTemplateDataSourceTestCase, overrideValuesFile string) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name = "%s"
}

data "cpln_helm_template" "%s" {
  name  = "%s"
  gvc   = cpln_gvc.%s.name
  chart = "%s"

  values = [file("%s"), file("%s")]

  set = {
    "secret.data" = "override-value"
  }

  set_string = {
    "secret.name" = "%s-secret"
  }
}
`, c.GvcResourceName, c.GvcName, c.ResourceName, c.Name, c.GvcResourceName, c.Chart, c.ValuesFile, overrideValuesFile, c.Name)
}

/*** Data Source Test Case ***/

// HelmTemplateDataSourceTestCase defines a specific data source test case.
type HelmTemplateDataSourceTestCase struct {
	ProviderTestCase
	GvcName         string
	GvcResourceName string
	Chart           string
	ValuesFile      string
}
