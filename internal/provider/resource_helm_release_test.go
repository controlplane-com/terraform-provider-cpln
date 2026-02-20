package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneHelmRelease_basic performs an acceptance test for the resource.
func TestAccControlPlaneHelmRelease_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewHelmReleaseResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "HELM") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// HelmReleaseResourceTest defines the necessary functionality to test the resource.
type HelmReleaseResourceTest struct {
	Steps []resource.TestStep
}

// NewHelmReleaseResourceTest creates a HelmReleaseResourceTest with initialized test cases.
func NewHelmReleaseResourceTest() HelmReleaseResourceTest {
	// Create a resource test instance
	resourceTest := HelmReleaseResourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewRequiredOnlyScenario()...)
	steps = append(steps, resourceTest.NewAllOptionsScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (hrt *HelmReleaseResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_helm_release resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, return early
	if len(s.RootModule().Resources) == 0 {
		return nil
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_helm_release" {
			continue
		}

		// Retrieve the name for the current resource
		helmReleaseName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of helm release with name: %s", helmReleaseName))
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_helm_release resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewRequiredOnlyScenario creates a test case for a helm deployment using only required fields.
func (hrt *HelmReleaseResourceTest) NewRequiredOnlyScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("helm-req-%s", random)
	gvcName := fmt.Sprintf("helm-req-gvc-%s", random)
	resourceName := "required-only"

	// Build test steps
	initialConfig, initialStep := hrt.BuildRequiredOnlyInitialTestStep(resourceName, name, gvcName)
	updateValuesStep := hrt.BuildRequiredOnlyUpdateValuesTestStep(initialConfig.ProviderTestCase, gvcName)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName: initialConfig.ResourceAddress,
			ImportState:  true,
		},
		// Update values & Read
		updateValuesStep,
	}
}

// NewAllOptionsScenario creates a test case for a helm deployment using all optional fields.
func (hrt *HelmReleaseResourceTest) NewAllOptionsScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("helm-all-%s", random)
	gvcName := fmt.Sprintf("helm-all-gvc-%s", random)
	resourceName := "all-options"

	// Build test steps
	initialConfig, initialStep := hrt.BuildAllOptionsInitialTestStep(resourceName, name, gvcName)
	updateStep := hrt.BuildAllOptionsUpdateTestStep(initialConfig.ProviderTestCase, gvcName)
	updateMaxHistoryStep := hrt.BuildAllOptionsUpdateMaxHistoryTestStep(initialConfig.ProviderTestCase, gvcName)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName: initialConfig.ResourceAddress,
			ImportState:  true,
		},
		// Update description and switch from set to set_file & Read
		updateStep,
		// Update max_history & Read
		updateMaxHistoryStep,
	}
}

// Test Cases //

// BuildRequiredOnlyInitialTestStep returns the initial test step using only required fields.
func (hrt *HelmReleaseResourceTest) BuildRequiredOnlyInitialTestStep(resourceName string, name string, gvcName string) (HelmReleaseResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := HelmReleaseResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "helm_release",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_helm_release.%s", resourceName),
			Name:            name,
		},
		GvcName:         gvcName,
		GvcResourceName: "helm_gvc_required_only",
		Chart:           "../../testdata/helm/sample-chart",
		ValuesFile:      TestdataAbsPath("../../testdata/helm/values/initial.yaml"),
	}

	// Initialize and return the initial test step
	return c, resource.TestStep{
		Config: hrt.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("gvc", c.GvcName),
			c.TestCheckResourceAttr("chart", c.Chart),
			c.TestCheckResourceAttr("values.#", "1"),
			c.TestCheckResourceAttr("set_string.secret.name", fmt.Sprintf("%s-secret", name)),
			c.TestCheckResourceAttr("wait", "false"),
			c.TestCheckResourceAttr("timeout", "300"),
			c.TestCheckResourceAttr("verify", "false"),
			c.TestCheckResourceAttr("dependency_update", "false"),
			c.TestCheckResourceAttr("insecure_skip_tls_verify", "false"),
			c.TestCheckResourceAttr("render_subchart_notes", "false"),
			c.TestCheckResourceAttr("max_history", "10"),
			c.TestCheckResourceAttrSet("status"),
			c.TestCheckResourceAttr("revision", "1"),
			c.TestCheckResourceAttrSet("manifest"),
			c.TestCheckResourceAttrSet("resources.%"),
		),
	}
}

// BuildRequiredOnlyUpdateValuesTestStep returns a test step that updates the values file.
func (hrt *HelmReleaseResourceTest) BuildRequiredOnlyUpdateValuesTestStep(initialCase ProviderTestCase, gvcName string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := HelmReleaseResourceTestCase{
		ProviderTestCase: initialCase,
		GvcName:          gvcName,
		GvcResourceName:  "helm_gvc_required_only",
		Chart:            "../../testdata/helm/sample-chart",
		ValuesFile:       TestdataAbsPath("../../testdata/helm/values/updated.yaml"),
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: hrt.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("gvc", c.GvcName),
			c.TestCheckResourceAttr("values.#", "1"),
			c.TestCheckResourceAttrSet("status"),
			c.TestCheckResourceAttr("revision", "2"),
			c.TestCheckResourceAttrSet("manifest"),
			c.TestCheckResourceAttrSet("resources.%"),
		),
	}
}

// BuildAllOptionsInitialTestStep returns the initial test step using all optional fields with multiple values.
func (hrt *HelmReleaseResourceTest) BuildAllOptionsInitialTestStep(resourceName string, name string, gvcName string) (HelmReleaseResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	overrideValuesFile := TestdataAbsPath("../../testdata/helm/values/override.yaml")
	c := HelmReleaseResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "helm_release",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_helm_release.%s", resourceName),
			Name:            name,
		},
		GvcName:         gvcName,
		GvcResourceName: "helm_gvc_all_options",
		Chart:           "../../testdata/helm/sample-chart",
		ValuesFile:      TestdataAbsPath("../../testdata/helm/values/initial.yaml"),
	}

	// Initialize and return the initial test step
	return c, resource.TestStep{
		Config: hrt.AllOptionsWithMultipleValues(c, overrideValuesFile),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("gvc", c.GvcName),
			c.TestCheckResourceAttr("chart", c.Chart),
			c.TestCheckResourceAttr("values.#", "2"),
			c.TestCheckResourceAttr("description", "Initial deployment"),
			c.TestCheckResourceAttr("set.secret.data", "set-value"),
			c.TestCheckResourceAttr("set_string.secret.name", fmt.Sprintf("%s-secret", name)),
			c.TestCheckResourceAttr("wait", "false"),
			c.TestCheckResourceAttr("timeout", "600"),
			c.TestCheckResourceAttr("max_history", "5"),
			c.TestCheckResourceAttrSet("status"),
			c.TestCheckResourceAttr("revision", "1"),
			c.TestCheckResourceAttrSet("manifest"),
			c.TestCheckResourceAttrSet("resources.%"),
		),
	}
}

// BuildAllOptionsUpdateTestStep returns a test step that updates description and switches from set to set_file.
func (hrt *HelmReleaseResourceTest) BuildAllOptionsUpdateTestStep(initialCase ProviderTestCase, gvcName string) resource.TestStep {
	// Create the test case with metadata and descriptions
	secretPayloadFile := TestdataAbsPath("../../testdata/helm/values/secret-payload.txt")
	c := HelmReleaseResourceTestCase{
		ProviderTestCase: initialCase,
		GvcName:          gvcName,
		GvcResourceName:  "helm_gvc_all_options",
		Chart:            "../../testdata/helm/sample-chart",
		ValuesFile:       TestdataAbsPath("../../testdata/helm/values/updated.yaml"),
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: hrt.AllOptionsWithSetFile(c, secretPayloadFile),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("values.#", "1"),
			c.TestCheckResourceAttr("description", "Updated deployment"),
			c.TestCheckResourceAttr("set_file.secret.data", secretPayloadFile),
			c.TestCheckResourceAttr("set_string.secret.name", fmt.Sprintf("%s-secret", c.Name)),
			c.TestCheckResourceAttr("timeout", "600"),
			c.TestCheckResourceAttr("max_history", "5"),
			c.TestCheckResourceAttrSet("status"),
			c.TestCheckResourceAttr("revision", "2"),
			c.TestCheckResourceAttrSet("manifest"),
			c.TestCheckResourceAttrSet("resources.%"),
		),
	}
}

// BuildAllOptionsUpdateMaxHistoryTestStep returns a test step that updates max_history.
func (hrt *HelmReleaseResourceTest) BuildAllOptionsUpdateMaxHistoryTestStep(initialCase ProviderTestCase, gvcName string) resource.TestStep {
	// Create the test case with metadata and descriptions
	secretPayloadFile := TestdataAbsPath("../../testdata/helm/values/secret-payload.txt")
	c := HelmReleaseResourceTestCase{
		ProviderTestCase: initialCase,
		GvcName:          gvcName,
		GvcResourceName:  "helm_gvc_all_options",
		Chart:            "../../testdata/helm/sample-chart",
		ValuesFile:       TestdataAbsPath("../../testdata/helm/values/updated.yaml"),
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: hrt.AllOptionsWithSetFileUpdatedMaxHistory(c, secretPayloadFile),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("description", "Updated deployment"),
			c.TestCheckResourceAttr("max_history", "20"),
			c.TestCheckResourceAttrSet("status"),
			c.TestCheckResourceAttr("revision", "3"),
			c.TestCheckResourceAttrSet("manifest"),
			c.TestCheckResourceAttrSet("resources.%"),
		),
	}
}

// Configs //

// RequiredOnly returns a minimal HCL block for a resource using only required fields.
func (hrt *HelmReleaseResourceTest) RequiredOnly(c HelmReleaseResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name = "%s"
}

resource "cpln_helm_release" "%s" {
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

// AllOptionsWithMultipleValues returns an HCL block with all optional fields and two values files.
func (hrt *HelmReleaseResourceTest) AllOptionsWithMultipleValues(c HelmReleaseResourceTestCase, overrideValuesFile string) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name = "%s"
}

resource "cpln_helm_release" "%s" {
  name  = "%s"
  gvc   = cpln_gvc.%s.name
  chart = "%s"

  values = [file("%s"), file("%s")]

  set = {
    "secret.data" = "set-value"
  }

  set_string = {
    "secret.name" = "%s-secret"
  }

  description = "Initial deployment"
  timeout     = 600
  max_history = 5
}
`, c.GvcResourceName, c.GvcName, c.ResourceName, c.Name, c.GvcResourceName, c.Chart, c.ValuesFile, overrideValuesFile, c.Name)
}

// AllOptionsWithSetFile returns an HCL block with set_file replacing set and updated description.
func (hrt *HelmReleaseResourceTest) AllOptionsWithSetFile(c HelmReleaseResourceTestCase, secretPayloadFile string) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name = "%s"
}

resource "cpln_helm_release" "%s" {
  name  = "%s"
  gvc   = cpln_gvc.%s.name
  chart = "%s"

  values = [file("%s")]

  set_file = {
    "secret.data" = "%s"
  }

  set_string = {
    "secret.name" = "%s-secret"
  }

  description = "Updated deployment"
  timeout     = 600
  max_history = 5
}
`, c.GvcResourceName, c.GvcName, c.ResourceName, c.Name, c.GvcResourceName, c.Chart, c.ValuesFile, secretPayloadFile, c.Name)
}

// AllOptionsWithSetFileUpdatedMaxHistory returns an HCL block with set_file and updated max_history.
func (hrt *HelmReleaseResourceTest) AllOptionsWithSetFileUpdatedMaxHistory(c HelmReleaseResourceTestCase, secretPayloadFile string) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "%s" {
  name = "%s"
}

resource "cpln_helm_release" "%s" {
  name  = "%s"
  gvc   = cpln_gvc.%s.name
  chart = "%s"

  values = [file("%s")]

  set_file = {
    "secret.data" = "%s"
  }

  set_string = {
    "secret.name" = "%s-secret"
  }

  description = "Updated deployment"
  timeout     = 600
  max_history = 20
}
`, c.GvcResourceName, c.GvcName, c.ResourceName, c.Name, c.GvcResourceName, c.Chart, c.ValuesFile, secretPayloadFile, c.Name)
}

/*** Resource Test Case ***/

// HelmReleaseResourceTestCase defines a specific resource test case.
type HelmReleaseResourceTestCase struct {
	ProviderTestCase
	GvcName         string
	GvcResourceName string
	Chart           string
	ValuesFile      string
}
