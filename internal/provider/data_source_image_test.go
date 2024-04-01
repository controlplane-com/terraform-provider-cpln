package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCplnImage_basic(t *testing.T) {

	imageName := "cpln_doc_demo:7"
	resourceName := "data.cpln_image.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DATA_SOURCE_IMAGE") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCplnImage(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", imageName),
				),
			},
		},
	})
}

func testAccDataSourceCplnImage(name string) string {
	return fmt.Sprintf(`data "cpln_image" "test" {
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
