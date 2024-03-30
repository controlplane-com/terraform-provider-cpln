package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCplnImages_basic(t *testing.T) {

	var images client.Images
	resourceName := "data.cpln_images.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DATA_SOURCE_IMAGES") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCplnImages(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnImagesExists(resourceName, &images),
					testAccCheckControlPlaneImagesAttributes(&images),
					resource.TestCheckResourceAttrSet(resourceName, "images.#"),
				),
			},
		},
	})
}

func testAccDataSourceCplnImages() string {
	return `data "cpln_images" "test" {}`
}

func testAccCheckCplnImagesExists(resourceName string, images *client.Images) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Can't find images data source: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Images data source ID not set")
		}

		client := testAccProvider.Meta().(*client.Client)
		_images, err := client.GetImages()

		if err != nil {
			return err
		}

		*images = *_images

		return nil
	}
}

func testAccCheckControlPlaneImagesAttributes(images *client.Images) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if len(images.Items) == 0 {
			return fmt.Errorf("Images data source has no images")
		}

		return nil
	}
}
