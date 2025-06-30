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

type CustomLocationResourceTest struct {
	resourceName string
	name         string
	region       string
}

/*** Acceptance Test ***/

// TestAccControlPlaneCustomLocation_basic performs an acceptance test for the resource.
func TestAccControlPlaneCustomLocation_basic(t *testing.T) {
	// Define unique values for the API resource to be used during the test lifecycle
	description := "Custom Location created using terraform for acceptance tests"
	updateDescription := "Custom Location updated using terraform for acceptance tests"

	// Initialize the data
	var byokEnabledName string = "custom-location-byok-enabled" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	var byokDisabledName string = "custom-location-byok-disabled" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	var byokEnabled = CustomLocationResourceTest{
		resourceName: "cpln_custom_location.byok-enabled",
		name:         byokEnabledName,
		region:       byokEnabledName,
	}

	var byokDisabled = CustomLocationResourceTest{
		resourceName: "cpln_custom_location.byok-disabled",
		name:         byokDisabledName,
		region:       byokDisabledName,
	}

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "LOCATION") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             testAccCheckControlPlaneCustomLocationCheckDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccControlPlaneCustomLocationCreate(byokEnabled, byokDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					// BYOK Enabled
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "id", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "name", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "description", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokEnabled.name)),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "region", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "enabled", "true"),
					testAccCheckControlPlaneCustomLocationExists(byokEnabled.resourceName, byokEnabled.name),

					// BYOK Disabled
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "id", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "name", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "description", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokDisabled.name)),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "region", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "enabled", "false"),
					testAccCheckControlPlaneCustomLocationExists(byokDisabled.resourceName, byokDisabled.name),
				),
			},
			// ImportState testing
			{
				ResourceName: byokEnabled.resourceName,
				ImportState:  true,
			},
			{
				ResourceName: byokDisabled.resourceName,
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: testAccControlPlaneCustomLocationUpdate(description, byokEnabled, byokDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					// BYOK Enabled
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "id", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "name", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "description", description),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokEnabled.name)),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "region", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "enabled", "false"),
					testAccCheckControlPlaneCustomLocationExists(byokEnabled.resourceName, byokEnabled.name),

					// BYOK Disabled
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "id", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "name", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "description", description),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokDisabled.name)),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "region", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "enabled", "true"),
					testAccCheckControlPlaneCustomLocationExists(byokDisabled.resourceName, byokDisabled.name),
				),
			},
			{
				Config: testAccControlPlaneCustomLocationUpdateAddTag(updateDescription, byokEnabled, byokDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					// BYOK Enabled
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "id", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "name", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokEnabled.name)),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "region", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "enabled", "false"),
					testAccCheckControlPlaneCustomLocationExists(byokEnabled.resourceName, byokEnabled.name),

					// BYOK Disabled
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "id", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "name", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokDisabled.name)),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "region", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "enabled", "true"),
					testAccCheckControlPlaneCustomLocationExists(byokDisabled.resourceName, byokDisabled.name),
				),
			},
			{
				Config: testAccControlPlaneCustomLocationUpdate(description, byokEnabled, byokDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					// BYOK Enabled
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "id", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "name", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "description", description),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokEnabled.name)),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "region", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "enabled", "false"),
					testAccCheckControlPlaneCustomLocationExists(byokEnabled.resourceName, byokEnabled.name),

					// BYOK Disabled
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "id", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "name", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "description", description),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokDisabled.name)),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "region", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "enabled", "true"),
					testAccCheckControlPlaneCustomLocationExists(byokDisabled.resourceName, byokDisabled.name),
				),
			},
			// Update to Required Only
			{
				Config: testAccControlPlaneCustomLocationCreate(byokEnabled, byokDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					// BYOK Enabled
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "id", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "name", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "description", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokEnabled.name)),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "region", byokEnabled.name),
					resource.TestCheckResourceAttr(byokEnabled.resourceName, "enabled", "true"),
					testAccCheckControlPlaneCustomLocationExists(byokEnabled.resourceName, byokEnabled.name),

					// BYOK Disabled
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "id", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "name", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "description", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "self_link", GetSelfLink(OrgName, "location", byokDisabled.name)),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "cloud_provider", "byok"),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "region", byokDisabled.name),
					resource.TestCheckResourceAttr(byokDisabled.resourceName, "enabled", "false"),
					testAccCheckControlPlaneCustomLocationExists(byokDisabled.resourceName, byokDisabled.name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccCheckControlPlaneCustomLocationCheckDestroy verifies that all resources have been destroyed.
func testAccCheckControlPlaneCustomLocationCheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_custom_location resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_custom_location" {
			continue
		}

		// Retrieve the name for the current resource
		customLocationName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of custom location with name: %s", customLocationName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		customLocation, code, err := TestProvider.client.GetLocation(customLocationName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if custom location %s exists: %w", customLocationName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if customLocation != nil {
			return fmt.Errorf("CheckDestroy failed: custom location %s still exists in the system", *customLocation.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_custom_location resources have been successfully destroyed")
	return nil
}

// testAccCheckControlPlaneCustomLocationExists verifies that a specified resource exist within the Terraform state and in the data service.
func testAccCheckControlPlaneCustomLocationExists(resourceName string, customLocationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of custom location: %s. Total resources: %d", customLocationName, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != customLocationName {
			return fmt.Errorf("resource ID %s does not match expected custom location name %s", rs.Primary.ID, customLocationName)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteCustomLocation, _, err := TestProvider.client.GetLocation(customLocationName)
		if err != nil {
			return fmt.Errorf("error retrieving custom location from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteCustomLocation.Name != customLocationName {
			return fmt.Errorf("mismatch in custom location name: expected %s, got %s", customLocationName, *remoteCustomLocation.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("Custom Location %s verified successfully in both state and external system.", customLocationName))
		return nil
	}
}

/*** Configs ***/

// testAccControlPlaneCustomLocationCreate constructs HCL blocks to create two custom locations: one enabled and one disabled for BYOK.
func testAccControlPlaneCustomLocationCreate(byokEnabled CustomLocationResourceTest, byokDisabled CustomLocationResourceTest) string {
	return fmt.Sprintf(`
resource "cpln_custom_location" "byok-enabled" {
  name           = "%s"
	cloud_provider = "byok"
	enabled        = "true"
}

resource "cpln_custom_location" "byok-disabled" {
  name           = "%s"
	cloud_provider = "byok"
	enabled        = "false"
}
`, byokEnabled.name, byokDisabled.name)
}

// testAccControlPlaneCustomLocationUpdate constructs HCL blocks to update two BYOK custom locations with tags and description.
func testAccControlPlaneCustomLocationUpdate(description string, byokEnabled CustomLocationResourceTest, byokDisabled CustomLocationResourceTest) string {
	return fmt.Sprintf(`
resource "cpln_custom_location" "byok-enabled" {
  name        = "%s"
	description = "%s"

	tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

	cloud_provider = "byok"
	enabled        = "false"
}

resource "cpln_custom_location" "byok-disabled" {
  name        = "%s"
	description = "%s"

	tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

	cloud_provider = "byok"
	enabled        = "true"
}
`, byokEnabled.name, description, byokDisabled.name, description)
}

// testAccControlPlaneCustomLocationUpdateAddTag constructs HCL blocks to add a new tag to BYOK custom locations.
func testAccControlPlaneCustomLocationUpdateAddTag(description string, byokEnabled CustomLocationResourceTest, byokDisabled CustomLocationResourceTest) string {
	return fmt.Sprintf(`
resource "cpln_custom_location" "byok-enabled" {
  name        = "%s"
	description = "%s"

	tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    new_tag             = "true"
  }

	cloud_provider = "byok"
	enabled        = "false"
}

resource "cpln_custom_location" "byok-disabled" {
  name        = "%s"
	description = "%s"

	tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    new_tag             = "true"
  }

	cloud_provider = "byok"
	enabled        = "true"
}
`, byokEnabled.name, description, byokDisabled.name, description)
}
