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
				Config: testAccControlPlaneCustomLocation("byok-loc", "desc-1", "byok", "true", "bar", "qux"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_custom_location.new", "name", "byok-loc"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "description", "desc-1"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "cloud_provider", "byok"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "enabled", "true"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "tags.foo", "bar"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "tags.baz", "qux"),
				),
			},
			{
				Config: testAccControlPlaneCustomLocation("byok-loc", "desc-2", "byok", "false", "bar-2", "qux-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_custom_location.new", "description", "desc-2"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "tags.foo", "bar-2"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "tags.baz", "qux-2"),
					resource.TestCheckResourceAttr("cpln_custom_location.new", "enabled", "false"),
				),
			},
		},
	})
}

func testAccControlPlaneCustomLocation(name string, description string, provider string, enabled string, val1 string, val2 string) string {

	TestLogger.Printf("Inside testAccControlPlaneCustomLocation")

	return fmt.Sprintf(`
	resource "cpln_custom_location" "new" {
		name         	= "%s"
		description 	= "%s"
		cloud_provider  = "%s"
		enabled 	 	= "%s"

		tags = {
			"foo"   = "%s"
			"baz"	= "%s"
		}
	}
    `, name, description, provider, enabled, val1, val2)
}

func testAccCheckControlPlaneCustomLocationCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneCustomLocationCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneCustomLocationCheckDestroy: rs.Type: %s", rs.Type)

		locationName := rs.Primary.ID

		err := c.DeleteCustomLocation(locationName)

		if err != nil {
			return fmt.Errorf("Error deleting Custom Location. Name: %s", locationName)
		}
	}

	return nil
}
