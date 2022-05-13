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

func TestAccControlPlaneIdentity_basic(t *testing.T) {

	var testIdentity client.Identity

	gName := "gvc-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	iName := "identity-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	aName := "agent-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	orgName := os.Getenv("CPLN_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "IDENTITY") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneIdentityCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneIdentity(orgName, randomName, gName, "GVC created using terraform for Identity acceptance tests", aName, iName, "Identity created using terraform for acceptance tests"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneIdentityExists("cpln_identity.test_identity", iName, gName, &testIdentity),
					// testAccCheckControlPlaneWorkloadAttributes(&testIdentity),
					resource.TestCheckResourceAttr("cpln_gvc.test_gvc", "description", "GVC created using terraform for Identity acceptance tests"),
					resource.TestCheckResourceAttr("cpln_identity.test_identity", "description", "Identity created using terraform for acceptance tests"),
				),
			},
			// {
			// 	Config: testAccControlPlaneIdentity(gName, "GVC created using terraform for acceptance tests", wName+"renamed", "Renamed Workload created using terraform for acceptance tests"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckControlPlaneWorkloadExists("cpln_workload.new", wName+"renamed", gName, &testWorkload),
			// 		testAccCheckControlPlaneWorkloadAttributes(&testWorkload),
			// 		resource.TestCheckResourceAttr("cpln_workload.new", "description", "Renamed Workload created using terraform for acceptance tests"),
			// 	),
			// },
		},
	})
}

func testAccControlPlaneIdentity(orgName, randomName, gvcName, gvcDescription, agentName, identityName, identityDescription string) string {

	TestLogger.Printf("Inside testAccControlPlaneIdentity")

	return fmt.Sprintf(`

	variable org_name {
		type = string
		default = "%s"
	}

	variable random_name {
		type = string
		default = "%s"
	}
	
	resource "cpln_gvc" "test_gvc" {
		name        = "%s"	
		description = "%s"

		locations = ["aws-us-west-2"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}

	resource "cpln_agent" "test_agent" {
		name        = "%s"
		description = "Test Agent created using Terraform"
	}

	resource "cpln_cloud_account" "test_aws_cloud_account" {

		name = "tf-ca-aws-${var.random_name}"
		description = "cloud account description tf-ca-aws" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		aws {
			role_arn = "arn:aws:iam::1234:role/test_role"
		}
	}

	resource "cpln_cloud_account" "test_azure_cloud_account" {

		name = "tf-ca-azure-${var.random_name}"
		description = "cloud account description tf-ca-azure" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		azure {
			// Use the full link for now
			// secret_link = "//secret/tf_secret_azure"
			secret_link = "/org/${var.org_name}/secret/tf_secret_azure"
		}
	}

	resource "cpln_cloud_account" "test_gcp_cloud_account" {

		name = "tf-ca-gcp-${var.random_name}"
		description = "cloud account description tf-ca-gcp" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		gcp {
			project_id = "cpln_gcp_project_1234"
		}
		
	}
	  
	resource "cpln_identity" "test_identity" {

		// depends_on = [cpln_cloud_account.test_aws_cloud_account]

  		gvc = cpln_gvc.test_gvc.name

		name        = "%s"	
		description = "%s"
 
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}

		network_resource {
			name = "test-network-resource-fqdn"
			agent_link = cpln_agent.test_agent.self_link
			fqdn = "domain.example.com"
			ports = [1234, 5432]
		}

		network_resource {
			name = "test-network-resource-fqdn-rip"
			agent_link = cpln_agent.test_agent.self_link
			fqdn = "domain2.example.com"
			resolver_ip = "192.168.1.1"
			ports = [12345, 54321]
		}

		network_resource {
			name = "test-network-resource-ip"
			agent_link = cpln_agent.test_agent.self_link
			ips = ["192.168.1.1", "192.168.1.250"]
			ports = [3099, 7890]
		}

		aws_access_policy {
			cloud_account_link = cpln_cloud_account.test_aws_cloud_account.self_link
			
			// role_name = "rds-monitoring-role"

			policy_refs = ["aws::/job-function/SupportUser", "aws::AWSSupportAccess"]
			
			// trust_policy {
			// 	version = ""
			// 	statement = ""
			// }
		}

		azure_access_policy {
			cloud_account_link = cpln_cloud_account.test_azure_cloud_account.self_link

			role_assignment {
				scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group"
				roles = ["AcrPull",	"AcrPush"]
			}

			role_assignment {
				scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group/providers/Microsoft.Storage/storageAccounts/cplntest"
				roles = ["Support Request Contributor"]
			}
		}

		gcp_access_policy {
			
			cloud_account_link = cpln_cloud_account.test_gcp_cloud_account.self_link
			scopes = "https://www.googleapis.com/auth/cloud-platform"
		
			// service_account = "cpln-tf@cpln-test.iam.gserviceaccount.com"
			
			binding {
				resource = "//cloudresourcemanager.googleapis.com/projects/cpln-test"
				roles = ["roles/appengine.appViewer", "roles/actions.Viewer"]
			}

			binding {
				resource = "//iam.googleapis.com/projects/cpln-test/serviceAccounts/cpln-tf@cpln-test.iam.gserviceaccount.com"
				roles = ["roles/editor", "roles/iam.serviceAccountUser"]
			}
		}
	}
	
	`, orgName, randomName, gvcName, gvcDescription, agentName, identityName, identityDescription)
}

func testAccCheckControlPlaneIdentityExists(resourceName, identityName, gvcName string, identity *client.Identity) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneIdentityExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != identityName {
			return fmt.Errorf("Identity name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)

		wl, _, err := client.GetIdentity(identityName, gvcName)

		if err != nil {
			return err
		}

		if *wl.Name != identityName {
			return fmt.Errorf("Identity name does not match")
		}

		*identity = *wl

		return nil
	}
}

func testAccCheckControlPlaneIdentityCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneIdentityCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneIdentityCheckDestroy: rs.Type: %s", rs.Type)

		if rs.Type != "cpln_gvc" {
			continue
		}

		gvcName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneIdentityCheckDestroy: gvcName: %s", gvcName)

		gvc, _, _ := c.GetGvc(gvcName)
		if gvc != nil {
			return fmt.Errorf("GVC still exists. Name: %s. Associated Identities might still exist", *gvc.Name)
		}
	}

	return nil
}

// func TestControlPlane_FlattenIdentityStatus(t *testing.T) {

// 	status := &client.IdentityStatus{
// 		ObjectName: "cpln-terraform-test-o-qwx0zftz",
// 	}

// 	flatStatus := map[string]interface{}{
// 		"objectName": "cpln-terraform-test-o-qwx0zftz",
// 	}

// 	flattenedStatus := "" // flattenIdentityStatus(status)

// 	if diff := deep.Equal(flattenedStatus, flatStatus); diff != nil {
// 		t.Errorf("Workload Status was not flattened correctly. Diff: %s", diff)
// 	}
// }
