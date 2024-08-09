package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneCustomLocation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "LOCATION") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneCustomLocationCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneCustomLocation("byok-loc-name-02", "byok-loc-desc", "byok", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_custom_location.new", "enabled", "false"),
				),
			},
			// {
			// 	Config: testAccControlPlaneCustomLocation("true"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("cpln_custom_location.new", "enabled", "true"),
			// 	),
			// },
			// {
			// 	Config: testAccControlPlaneCustomLocation_ReferenceTags("false"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("cpln_custom_location.new", "enabled", "false"),
			// 	),
			// },
			// {
			// 	Config: testAccControlPlaneCustomLocation_ReferenceTags("true"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("cpln_custom_location.new", "enabled", "true"),
			// 	),
			// },
		},
	})
}

func testAccControlPlaneCustomLocation(name string, description string, provider string, enabled string) string {

	TestLogger.Printf("Inside testAccControlPlaneCustomLocation")

	return fmt.Sprintf(`
	resource "cpln_custom_location" "new" {
		name         	= "%s"
		description 	= "%s"
		cloud_provider  = "%s"
		enabled 	 	= "%s"

		tags = {
			"cpln/city"      = "Frankfurt"
			"cpln/continent" = "Europe"
			"cpln/country"   = "Germany"
		}
	}
    `, name, description, provider, enabled)
}

// func testAccControlPlaneCustomLocation_ReferenceTags(enabled string) string {

// 	TestLogger.Printf("Inside testAccControlPlaneCustomLocation")

// 	return fmt.Sprintf(`
// 	data "cpln_custom_location" "main-location" {
// 		name = "aws-eu-central-1"
// 	}

// 	resource "cpln_custom_location" "new" {
// 		name    = "aws-eu-central-1"
// 		enabled = %s

// 		tags = data.cpln_custom_location.main-location.tags
// 	}
//     `, enabled)
// }

func testAccCheckControlPlaneCustomLocationCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneCustomLocationCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneCustomLocationCheckDestroy: rs.Type: %s", rs.Type)

		locationName := rs.Primary.ID

		location, _, _ := c.GetLocation(locationName)

		if location == nil {
			return fmt.Errorf("Location does not exists. Name: %s", locationName)
		}

		if location.Spec == nil {
			return fmt.Errorf("Location does not have Spec. Name: %s", *location.Name)
		}

		if location.Spec.Enabled == nil {
			return fmt.Errorf("Location enabled is nil. Name: %s", *location.Name)
		}

		if !*location.Spec.Enabled {
			return fmt.Errorf("Location enabled wasn't reverted to the default value: true. Name: %s", *location.Name)
		}
	}

	return nil
}
