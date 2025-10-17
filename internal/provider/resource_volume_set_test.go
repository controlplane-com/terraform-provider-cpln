package cpln

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneVolumeSet_basic performs an acceptance test for the resource.
func TestAccControlPlaneVolumeSet_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewVolumeSetResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "VOLUME_SET") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// VolumeSetResourceTest defines the necessary functionality to test the resource.
type VolumeSetResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
}

// NewVolumeSetResourceTest creates a VolumeSetResourceTest with initialized test cases.
func NewVolumeSetResourceTest() VolumeSetResourceTest {
	// Create a resource test instance
	resourceTest := VolumeSetResourceTest{
		RandomName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
	}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewDefaultScenario("ext4", "general-purpose-ssd")...)
	steps = append(steps, resourceTest.NewDefaultScenario("xfs", "high-throughput-ssd")...)
	steps = append(steps, resourceTest.NewDefaultScenario("shared", "shared")...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (vsrt *VolumeSetResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_volume_set resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_gvc" {
			continue
		}

		// Retrieve the name for the current resource
		gvcName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of GVC with name: %s", gvcName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		gvc, code, err := TestProvider.client.GetGvc(gvcName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if GVC %s exists: %w", gvcName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if gvc != nil {
			return fmt.Errorf("CheckDestroy failed: GVC %s still exists in the system", *gvc.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_volume_set resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case for a volume set resource with required fields and multiple update steps.
func (vsrt *VolumeSetResourceTest) NewDefaultScenario(fileSystemType string, performanceClass string) []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("volume-set-default-%s", vsrt.RandomName)
	gvcName := fmt.Sprintf("gvc-%s", vsrt.RandomName)

	// Create the gvc case
	gvcCase := GvcResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "gvc",
			ResourceName:      "new",
			ResourceAddress:   "cpln_gvc.new",
			Name:              gvcName,
			Description:       gvcName,
			DescriptionUpdate: "gvc default description updated",
		},
	}

	// Create a gvc resource test instance
	gvcResourceTest := GvcResourceTest{}

	// Initialize the gvc config
	gvcConfig := gvcResourceTest.GvcRequiredOnly(gvcCase)

	// Build test steps
	initialConfig, initialStep := vsrt.BuildInitialTestStep(name, gvcConfig, gvcCase, fileSystemType, performanceClass)
	caseUpdate1 := vsrt.BuildUpdate1TestStep(initialConfig.ProviderTestCase, gvcConfig, gvcCase, fileSystemType, performanceClass)
	caseUpdate2 := vsrt.BuildUpdate2TestStep(initialConfig.ProviderTestCase, gvcConfig, gvcCase, fileSystemType, performanceClass)
	// caseUpdate3 := vsrt.BuildUpdate3TestStep(initialConfig.ResourceTestCase, gvcConfig, gvcCase)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName:  initialConfig.ResourceAddress,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:%s", gvcName, name),
		},
		// Update & Read
		caseUpdate1,
		caseUpdate2,
		// caseUpdate3,
		// Revert the resource to its initial state
		initialStep,
	}
}

// Test Cases //

// BuildInitialTestStep returns a test step and case for creating the volume set with required attributes.
func (vsrt *VolumeSetResourceTest) BuildInitialTestStep(name string, gvcConfig string, gvcCase GvcResourceTestCase, fileSystemType string, performanceClass string) (VolumeSetResourceTestCase, resource.TestStep) {
	initialCapacity := "10"

	if performanceClass == "high-throughput-ssd" {
		initialCapacity = "1000"
	}

	// Create the test case with metadata and descriptions
	c := VolumeSetResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "volumeset",
			ResourceName:      "new",
			ResourceAddress:   "cpln_volume_set.new",
			Name:              name,
			GvcName:           gvcCase.Name,
			Description:       name,
			DescriptionUpdate: "volume set default description updated",
		},
		GvcConfig:        gvcConfig,
		GvcCase:          gvcCase,
		InitialCapacity:  initialCapacity,
		PerformanceClass: performanceClass,
		FileSystemType:   fileSystemType,
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: vsrt.HclRequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "initial_capacity", c.InitialCapacity),
			resource.TestCheckResourceAttr(c.ResourceAddress, "performance_class", c.PerformanceClass),
			resource.TestCheckResourceAttr(c.ResourceAddress, "file_system_type", c.FileSystemType),
			resource.TestCheckResourceAttr(c.ResourceAddress, "volumeset_link", fmt.Sprintf("cpln://volumeset/%s", c.Name)),
			resource.TestCheckNoResourceAttr(c.ResourceAddress, "custom_encryption.#"),
		),
	}
}

// BuildUpdate1TestStep returns a test step with minimal optional fields for the volume set update.
func (vsrt *VolumeSetResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase, gvcConfig string, gvcCase GvcResourceTestCase, fileSystemType string, performanceClass string) resource.TestStep {
	initialCapacity := "10"

	if performanceClass == "high-throughput-ssd" {
		initialCapacity = "1000"
	}

	// Create the test case with metadata and descriptions
	c := VolumeSetResourceTestCase{
		ProviderTestCase: initialCase,
		GvcConfig:        gvcConfig,
		GvcCase:          gvcCase,
		InitialCapacity:  initialCapacity,
		PerformanceClass: performanceClass,
		FileSystemType:   fileSystemType,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: vsrt.HclWithMinimalOptionals(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", gvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "initial_capacity", c.InitialCapacity),
			resource.TestCheckResourceAttr(c.ResourceAddress, "performance_class", c.PerformanceClass),
			resource.TestCheckResourceAttr(c.ResourceAddress, "file_system_type", c.FileSystemType),
			resource.TestCheckResourceAttr(c.ResourceAddress, "storage_class_suffix", "demo-class"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "volumeset_link", fmt.Sprintf("cpln://volumeset/%s", c.Name)),
			c.TestCheckNestedBlocks("snapshots", []map[string]interface{}{
				{
					"create_final_snapshot": "true",
				},
			}),
			c.TestCheckNestedBlocks("autoscaling", []map[string]interface{}{
				{},
			}),
			c.TestCheckNestedBlocks("mount_options", []map[string]interface{}{
				{},
			}),
			c.TestCheckMapObjectAttr("custom_encryption.0.regions", map[string]interface{}{
				"aws-us-west-2": map[string]interface{}{
					"key_id": "arn:aws:kms:us-west-2:123456789012:key/minimal",
				},
			}),
		),
	}
}

// BuildUpdate2TestStep returns a test step with extended optional fields including autoscaling and snapshots.
func (vsrt *VolumeSetResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase, gvcConfig string, gvcCase GvcResourceTestCase, fileSystemType string, performanceClass string) resource.TestStep {
	initialCapacity := "20"

	if performanceClass == "high-throughput-ssd" {
		initialCapacity = "2000"
	}

	// Create the test case with metadata and descriptions
	c := VolumeSetResourceTestCase{
		ProviderTestCase: initialCase,
		GvcConfig:        gvcConfig,
		GvcCase:          gvcCase,
		InitialCapacity:  initialCapacity,
		PerformanceClass: performanceClass,
		FileSystemType:   fileSystemType,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: vsrt.HclWithMinimalOptionals2(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", gvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "initial_capacity", c.InitialCapacity),
			resource.TestCheckResourceAttr(c.ResourceAddress, "performance_class", c.PerformanceClass),
			resource.TestCheckResourceAttr(c.ResourceAddress, "file_system_type", c.FileSystemType),
			resource.TestCheckResourceAttr(c.ResourceAddress, "storage_class_suffix", "demo-class"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "volumeset_link", fmt.Sprintf("cpln://volumeset/%s", c.Name)),
			c.TestCheckNestedBlocks("snapshots", []map[string]interface{}{
				{
					"create_final_snapshot": "false",
					"retention_duration":    "2d",
					"schedule":              "0 * * * *",
				},
			}),
			c.TestCheckNestedBlocks("autoscaling", []map[string]interface{}{
				{
					"max_capacity":        "2000",
					"min_free_percentage": "2",
					"scaling_factor":      "2.2",
				},
			}),
			c.TestCheckNestedBlocks("mount_options", []map[string]interface{}{
				{
					"resources": []map[string]interface{}{
						{
							"max_cpu":    "2000m",
							"min_cpu":    "500m",
							"min_memory": "1Gi",
							"max_memory": "2Gi",
						},
					},
				},
			}),
			c.TestCheckMapObjectAttr("custom_encryption.0.regions", map[string]interface{}{
				"aws-us-west-2": map[string]interface{}{
					"key_id": "arn:aws:kms:us-west-2:123456789012:key/minimal",
				},
				"aws-us-east-1": map[string]interface{}{
					"key_id": "arn:aws:kms:us-east-1:123456789012:key/extended",
				},
			}),
		),
	}
}

// Configs //

// HclRequiredOnly returns HCL for volume set with only required fields.
func (vsrt *VolumeSetResourceTest) HclRequiredOnly(c VolumeSetResourceTestCase) string {
	return fmt.Sprintf(`
# GVC Resource
%s

resource "cpln_volume_set" "%s" {
  depends_on = [%s]

  name              = "%s"
  gvc               = "%s"
  initial_capacity  = %s
  performance_class = "%s"
  file_system_type  = "%s"
}
`, c.GvcConfig, c.ResourceName, c.GvcCase.ResourceAddress, c.Name, c.GvcCase.Name, c.InitialCapacity, c.PerformanceClass,
		c.FileSystemType,
	)
}

// HclWithMinimalOptionals returns HCL for volume set with minimal optional fields including tags and mount options.
func (vsrt *VolumeSetResourceTest) HclWithMinimalOptionals(c VolumeSetResourceTestCase) string {
	return fmt.Sprintf(`
# GVC Resource
%s

resource "cpln_volume_set" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = "%s"
  initial_capacity     = %s
  performance_class    = "%s"
  file_system_type     = "%s"

	custom_encryption {
		regions = {
			aws-us-west-2 = {
				key_id = "arn:aws:kms:us-west-2:123456789012:key/minimal"
			}
		}
	}

  storage_class_suffix = "demo-class"

  snapshots {}
  autoscaling {}
  mount_options {}
}
`, c.GvcConfig, c.ResourceName, c.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate, c.GvcCase.Name, c.InitialCapacity,
		c.PerformanceClass, c.FileSystemType,
	)
}

// HclWithMinimalOptionals2 returns HCL for volume set with extended optional fields including autoscaling and snapshot settings.
func (vsrt *VolumeSetResourceTest) HclWithMinimalOptionals2(c VolumeSetResourceTestCase) string {
	return fmt.Sprintf(`
# GVC Resource
%s

resource "cpln_volume_set" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = "%s"
  initial_capacity     = %s
  performance_class    = "%s"
  file_system_type     = "%s"

	custom_encryption {
		regions = {
			aws-us-west-2 = {
				key_id = "arn:aws:kms:us-west-2:123456789012:key/minimal"
			}

			"aws-us-east-1" = {
				key_id = "arn:aws:kms:us-east-1:123456789012:key/extended"
			}
		}
	}

  storage_class_suffix = "demo-class"

  snapshots {
    create_final_snapshot = false
    retention_duration    = "2d"
    schedule              = "0 * * * *"
  }

  autoscaling {
    max_capacity        = 2000
    min_free_percentage = 2
    scaling_factor      = 2.2
  }

  mount_options {
    resources {}
  }
}
`, c.GvcConfig, c.ResourceName, c.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate, c.GvcCase.Name, c.InitialCapacity,
		c.PerformanceClass, c.FileSystemType,
	)
}

/*** Resource Test Case ***/

// VolumeSetResourceTestCase defines a specific resource test case.
type VolumeSetResourceTestCase struct {
	ProviderTestCase
	GvcConfig        string
	GvcCase          GvcResourceTestCase
	InitialCapacity  string
	PerformanceClass string
	FileSystemType   string
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (vsrtc *VolumeSetResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of volume set: %s. Total resources: %d", vsrtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[vsrtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", vsrtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != vsrtc.Name {
			return fmt.Errorf("resource ID %s does not match expected volume set name %s", rs.Primary.ID, vsrtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteVolumeSet, _, err := TestProvider.client.GetVolumeSet(vsrtc.Name, vsrtc.GvcCase.Name)
		if err != nil {
			return fmt.Errorf("error retrieving volume set from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteVolumeSet.Name != vsrtc.Name {
			return fmt.Errorf("mismatch in volume set name: expected %s, got %s", vsrtc.Name, *remoteVolumeSet.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("volume set %s verified successfully in both state and external system.", vsrtc.Name))
		return nil
	}
}
