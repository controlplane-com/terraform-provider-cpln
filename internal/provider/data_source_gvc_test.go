package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCplnGvc_basic(t *testing.T) {

	resourceName := "data.cpln_gvc.test"
	gvcName := "default-gvc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DATA_SOURCE_GVC") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCplnGvcConfig(gvcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnGvcExists(resourceName),
				),
			},
		},
	})
}

func testAccDataSourceCplnGvcConfig(gvcName string) string {
	return fmt.Sprintf(`data "cpln_gvc" "test" {
		name = "%s"
	}`, gvcName)
}

func testAccCheckCplnGvcExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Can't find gvc data source: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("GVC data source ID not set")
		}

		return nil
	}
}
