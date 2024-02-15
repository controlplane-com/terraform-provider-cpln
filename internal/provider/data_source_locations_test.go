package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCplnLocations_basic(t *testing.T) {

	resourceName := "data.cpln_locations.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DATA_SOURCE_LOCATIONS") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCplnLocationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnLocationsExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "locations.#"),
					// Add more detailed checks here, e.g., for specific attributes within locations
				),
			},
		},
	})
}

func testAccDataSourceCplnLocationsConfig() string {
	return `data "cpln_locations" "test" {}`
}

func testAccCheckCplnLocationsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find locations data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Locations data source ID not set")
		}

		return nil
	}
}
