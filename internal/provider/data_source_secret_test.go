package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCplnSecret_basic(t *testing.T) {

	resourceName := "data.cpln_secret.test"
	secretName := "test-secret-opaque"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DATA_SOURCE_SECRET") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCplnSecret(secretName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCplnSecretExists(resourceName),
				),
			},
		},
	})
}

func testAccDataSourceCplnSecret(secretName string) string {
	return fmt.Sprintf(`data "cpln_secret" "test" {
		name = "%s"
	}`, secretName)
}

func testAccCheckCplnSecretExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Can't find secret data source: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Secret data source ID not set")
		}

		return nil
	}
}
