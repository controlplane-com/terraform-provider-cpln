package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneServiceAccount_basic(t *testing.T) {

	randomName := "service-account-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "SERVICE-ACCOUNT") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneServiceAccountCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneServiceAccount(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_service_account.tf-sa", "name", randomName),
					resource.TestCheckResourceAttr("cpln_service_account_key.tf_sa_key_01", "description", "key-01-"+randomName),
					resource.TestCheckResourceAttr("cpln_service_account_key.tf_sa_key_02", "description", "key-02-"+randomName),
				),
			},
		},
	})
}

func testAccControlPlaneServiceAccount(name string) string {

	return fmt.Sprintf(`
	
		variable "random-name" {
			type = string
			default = "%s"
		}

		resource "cpln_service_account" "tf-sa" {

			name = var.random-name
			description = "service account description ${var.random-name}" 
			
			tags = {
				terraform_generated = "true"
				acceptance_test = "true"
			}
		}

		resource "cpln_service_account_key" "tf_sa_key_01" {
			service_account_name = cpln_service_account.tf-sa.name
			description = "key-01-${var.random-name}"
		}

		resource "cpln_service_account_key" "tf_sa_key_02" {

			// remove below to test parallel adding of keys
			depends_on = [cpln_service_account_key.tf_sa_key_01]

			service_account_name = cpln_service_account.tf-sa.name
			description = "key-02-${var.random-name}"
		}
	`, name)
}

func testAccCheckControlPlaneServiceAccountCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_service_account" {
			continue
		}

		saName := rs.Primary.ID

		sa, _, _ := c.GetServiceAccount(saName)
		if sa != nil {
			return fmt.Errorf("Service Account still exists. Name: %s", *sa.Name)
		}
	}

	return nil
}
