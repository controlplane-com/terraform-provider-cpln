package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneLocation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "LOCATION") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneLocationCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneLocation("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_location.new", "enabled", "false"),
					resource.TestCheckResourceAttr("cpln_location.new", "tags.hello", "world"),
				),
			},
			{
				Config: testAccControlPlaneLocation_NoTags("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_location.new", "enabled", "true"),
				),
			},
		},
	})
}

func testAccControlPlaneLocation(enabled string) string {

	TestLogger.Printf("Inside testAccControlPlaneLocation")

	return fmt.Sprintf(`
	resource "cpln_location" "new" {
		name    = "aws-eu-central-1"
		enabled = %s

		tags = {
			hello = "world"
		}
	}
	`, enabled)
}

func testAccControlPlaneLocation_NoTags(enabled string) string {

	TestLogger.Printf("Inside testAccControlPlaneLocation")

	return fmt.Sprintf(`
	resource "cpln_location" "new" {
		name    = "aws-eu-central-1"
		enabled = %s
	}
	`, enabled)
}

func testAccCheckControlPlaneLocationCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneLocationCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneLocationCheckDestroy: rs.Type: %s", rs.Type)

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
