package cpln

import (
	"fmt"
	"os"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneCloudAccount_basic(t *testing.T) {

	var testCloudAccountAws, testCloudAccountAzure, testCloudAccountGcp, testCloudAccountNgs client.CloudAccount

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	updateRole := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	orgName := os.Getenv("CPLN_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "CLOUD-ACCOUNT") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneCloudAccountCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneCloudAccount(orgName, randomName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-aws", "tf-ca-aws-"+randomName, &testCloudAccountAws),
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-azure", "tf-ca-azure-"+randomName, &testCloudAccountAzure),
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-gcp", "tf-ca-gcp-"+randomName, &testCloudAccountGcp),
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-ngs", "tf-ca-ngs-"+randomName, &testCloudAccountNgs),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-aws", "name", "tf-ca-aws-"+randomName),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-azure", "name", "tf-ca-azure-"+randomName),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-gcp", "name", "tf-ca-gcp-"+randomName),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-ngs", "name", "tf-ca-ngs-"+randomName),
				),
			},
			{
				Config: testAccControlPlaneCloudAccount(orgName, randomName, updateRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-aws", "tf-ca-aws-"+randomName, &testCloudAccountAws),
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-azure", "tf-ca-azure-"+randomName, &testCloudAccountAzure),
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-gcp", "tf-ca-gcp-"+randomName, &testCloudAccountGcp),
					testAccCheckControlPlaneCloudAccountExists("cpln_cloud_account.tf-ca-ngs", "tf-ca-ngs-"+randomName, &testCloudAccountNgs),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-aws", "name", "tf-ca-aws-"+randomName),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-azure", "name", "tf-ca-azure-"+randomName),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-gcp", "name", "tf-ca-gcp-"+randomName),
					resource.TestCheckResourceAttr("cpln_cloud_account.tf-ca-ngs", "name", "tf-ca-ngs-"+randomName),
				),
			},
		},
	})
}

func testAccControlPlaneCloudAccount(orgName, name, update string) string {

	return fmt.Sprintf(`

	variable org_name {
		type = string
		default = "%s"
	}

	variable random_name {
		type = string
		default = "%s"
	}

	variable update_name {
		type = string
		default = "%s"
	}

	resource "cpln_cloud_account" "tf-ca-aws" {

		name = "tf-ca-aws-${var.random_name}"
		description = "cloud account description tf-ca-aws-${var.update_name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		aws {
			role_arn = "arn:aws:iam::1234:role/test_role${var.update_name}"
		}
	}

	resource "cpln_cloud_account" "tf-ca-azure" {

		name = "tf-ca-azure-${var.random_name}"
		description = "cloud account description tf-ca-azure-${var.update_name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		azure {
			// Use the full link for now
			// secret_link = "//secret/tf_secret_azure"
			secret_link = "/org/${var.org_name}/secret/tf_secret_azure${var.update_name}"
		}
	}

	resource "cpln_cloud_account" "tf-ca-gcp" {

		name = "tf-ca-gcp-${var.random_name}"
		description = "cloud account description tf-ca-gcp-${var.update_name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		gcp {
			project_id = "cpln_gcp_project_${var.update_name}"
		}
	}

	resource "cpln_cloud_account" "tf-ca-ngs" {

		name = "tf-ca-ngs-${var.random_name}"
		description = "cloud account description tf-ca-ngs-${var.update_name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		ngs {
			// secret_link = cpln_secret.nats_account.self_link
			secret_link = "/org/${var.org_name}/secret/tf_secret_nats${var.update_name}"
		}
	}
	
	
	`, orgName, name, update)
}

func testAccCheckControlPlaneCloudAccountExists(resourceName, cloudAccountName string, cloudAccount *client.CloudAccount) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != cloudAccountName {
			return fmt.Errorf("Cloud Account name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)

		ca, _, err := client.GetCloudAccount(cloudAccountName)

		if err != nil {
			return err
		}

		if *ca.Name != cloudAccountName {
			return fmt.Errorf("Cloud Account name does not match")
		}

		*cloudAccount = *ca

		return nil
	}
}

func testAccCheckControlPlaneCloudAccountCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy For Cloud Account. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_cloud_account" {
			continue
		}

		saName := rs.Primary.ID

		sa, _, _ := c.GetServiceAccount(saName)
		if sa != nil {
			return fmt.Errorf("Cloud Account still exists. Name: %s", *sa.Name)
		}
	}

	return nil
}
