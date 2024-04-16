package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/go-test/deep"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCplnImages_basic(t *testing.T) {

	var images client.ImagesQueryResult

	allImagesQuery := client.Query{
		Kind: GetString("image"),
		Spec: &client.Spec{
			Match: GetString("all"),
		},
	}

	specificRepositoryQuery := client.Query{
		Kind: GetString("image"),
		Spec: &client.Spec{
			Match: GetString("all"),
			Terms: &[]client.Term{
				{
					Op:       GetString("="),
					Property: GetString("repository"),
					Value:    GetString("call-internal-service-3000"),
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DATA_SOURCE_IMAGES") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCplnAllImages(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnImagesExists("data.cpln_images.all-images", &images, allImagesQuery),
					testAccCheckControlPlaneImagesAttributes("data.cpln_images.all-images", &images, 19),
					resource.TestCheckResourceAttrSet("data.cpln_images.all-images", "images.#"),
				),
			},
			{
				Config: testAccDataSourceCplnSpecificImages("call-internal-service-3000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnImagesExists("data.cpln_images.specific-images", &images, specificRepositoryQuery),
					testAccCheckControlPlaneImagesAttributes("data.cpln_images.specific-images", &images, 5),
					resource.TestCheckResourceAttrSet("data.cpln_images.specific-images", "images.#"),
				),
			},
		},
	})
}

func testAccDataSourceCplnAllImages() string {
	return `data "cpln_images" "all-images" {}`
}

func testAccDataSourceCplnSpecificImages(repository string) string {
	return fmt.Sprintf(`data "cpln_images" "specific-images" {
		query {
			fetch = "items"
			spec {
				match = "all"
				terms {
					op 	     = "="
					property = "repository"
					value	 = "%s"
				}
			}
		}
	}`, repository)
}

func testAccCheckCplnImagesExists(resourceName string, images *client.ImagesQueryResult, query client.Query) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Can't find images data source: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Images data source ID not set")
		}

		c := testAccProvider.Meta().(*client.Client)

		_images, err := c.GetImagesQuery(query)

		if err != nil {
			return err
		}

		*images = *_images

		return nil
	}
}

func testAccCheckControlPlaneImagesAttributes(resourceName string, images *client.ImagesQueryResult, expectedAmount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		amount := len(images.Items)

		if amount == 0 {
			return fmt.Errorf("%s has no images", resourceName)
		}

		if diff := deep.Equal(amount, expectedAmount); diff != nil {
			return fmt.Errorf("%s images amount does not match. Diff: %s", resourceName, diff)
		}

		return nil
	}
}
