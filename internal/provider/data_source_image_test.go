package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCplnImage_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DATA_SOURCE_IMAGE") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCplnImage("cpln_doc_demo:7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnImageExists("data.cpln_image.specific-image"),
					resource.TestCheckResourceAttr("data.cpln_image.specific-image", "name", "cpln_doc_demo:7"),
				),
			},
			{
				Config: testAccDataSourceCplnLatestImage("call-internal-service-3000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnImageExists("data.cpln_image.latest-image"),
					resource.TestCheckResourceAttr("data.cpln_image.latest-image", "name", "call-internal-service-3000:6"),
				),
			},
		},
	})
}

func testAccDataSourceCplnImage(name string) string {
	return fmt.Sprintf(`data "cpln_image" "specific-image" {
		name = "%s"
	}`, name)
}

func testAccDataSourceCplnLatestImage(name string) string {
	return fmt.Sprintf(`data "cpln_image" "latest-image" {
		name = "%s"
	}`, name)
}

func testAccCheckCplnImageExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Can't find image data source: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}
