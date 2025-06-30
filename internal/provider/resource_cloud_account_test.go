package cpln

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type CloudAccountResourceTest struct {
	resourceName string
	name         string
	data         string
	updateData   string
}

/*** Acceptance Test ***/

// TestAccControlPlaneCloudAccount_basic performs an acceptance test for the resource.
func TestAccControlPlaneCloudAccount_basic(t *testing.T) {
	// Define unique values for the API resource to be used during the test lifecycle
	description := "Cloud Account created using terraform for acceptance tests"
	updateDescription := "Cloud Account updated using terraform for acceptance tests"

	// Initialize the data
	var aws = CloudAccountResourceTest{
		resourceName: "cpln_cloud_account.aws",
		name:         "cloud-account-aws" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		data:         "arn:aws:iam::1234:role/test_role",
		updateData:   "arn:aws:iam::1234:role/test_role_updated",
	}

	var azure = CloudAccountResourceTest{
		resourceName: "cpln_cloud_account.azure",
		name:         "cloud-account-azure" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		data:         fmt.Sprintf("/org/%s/secret/tf_secret_azure", OrgName),
		updateData:   fmt.Sprintf("/org/%s/secret/tf_secret_azure_updated", OrgName),
	}

	var gcp = CloudAccountResourceTest{
		resourceName: "cpln_cloud_account.gcp",
		name:         "cloud-account-gcp" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		data:         "cpln_gcp_project",
		updateData:   "cpln_gcp_project_updated",
	}

	var ngs = CloudAccountResourceTest{
		resourceName: "cpln_cloud_account.ngs",
		name:         "cloud-account-ngs" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		data:         fmt.Sprintf("/org/%s/secret/tf_secret_nats", OrgName),
		updateData:   fmt.Sprintf("/org/%s/secret/tf_secret_nats_updated", OrgName),
	}

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "CLOUD_ACCOUNT") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             testAccCheckControlPlaneCloudAccountCheckDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccControlPlaneCloudAccountCreateRequiredOnly(aws, azure, gcp, ngs),
				Check: resource.ComposeAggregateTestCheckFunc(
					// AWS Resource Check
					resource.TestCheckResourceAttr(aws.resourceName, "id", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "name", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "description", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(aws.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", aws.name)),
					resource.TestCheckResourceAttr(aws.resourceName, "aws.0.role_arn", aws.data),
					testAccCheckControlPlaneCloudAccountExists(aws.resourceName, aws.name),

					// Azure Resource Check
					resource.TestCheckResourceAttr(azure.resourceName, "id", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "name", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "description", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(azure.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", azure.name)),
					resource.TestCheckResourceAttr(azure.resourceName, "azure.0.secret_link", azure.data),
					testAccCheckControlPlaneCloudAccountExists(azure.resourceName, azure.name),

					// GCP Resource Check
					resource.TestCheckResourceAttr(gcp.resourceName, "id", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "name", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "description", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(gcp.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", gcp.name)),
					resource.TestCheckResourceAttr(gcp.resourceName, "gcp.0.project_id", gcp.data),
					testAccCheckControlPlaneCloudAccountExists(gcp.resourceName, gcp.name),

					// NGS Resource Check
					resource.TestCheckResourceAttr(ngs.resourceName, "id", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "name", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "description", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(ngs.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", ngs.name)),
					resource.TestCheckResourceAttr(ngs.resourceName, "ngs.0.secret_link", ngs.data),
					testAccCheckControlPlaneCloudAccountExists(ngs.resourceName, ngs.name),
				),
			},
			// ImportState testing
			{
				ResourceName: aws.resourceName,
				ImportState:  true,
			},
			{
				ResourceName: azure.resourceName,
				ImportState:  true,
			},
			{
				ResourceName: gcp.resourceName,
				ImportState:  true,
			},
			{
				ResourceName: ngs.resourceName,
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: testAccControlPlaneCloudAccountUpdateWithOptionals(description, aws, azure, gcp, ngs),
				Check: resource.ComposeAggregateTestCheckFunc(
					// AWS Resource Check
					resource.TestCheckResourceAttr(aws.resourceName, "id", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "name", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "description", description),
					resource.TestCheckResourceAttr(aws.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(aws.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", aws.name)),
					resource.TestCheckResourceAttr(aws.resourceName, "aws.0.role_arn", aws.data),
					testAccCheckControlPlaneCloudAccountExists(aws.resourceName, aws.name),

					// Azure Resource Check
					resource.TestCheckResourceAttr(azure.resourceName, "id", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "name", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "description", description),
					resource.TestCheckResourceAttr(azure.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(azure.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", azure.name)),
					resource.TestCheckResourceAttr(azure.resourceName, "azure.0.secret_link", azure.data),
					testAccCheckControlPlaneCloudAccountExists(azure.resourceName, azure.name),

					// GCP Resource Check
					resource.TestCheckResourceAttr(gcp.resourceName, "id", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "name", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "description", description),
					resource.TestCheckResourceAttr(gcp.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(gcp.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", gcp.name)),
					resource.TestCheckResourceAttr(gcp.resourceName, "gcp.0.project_id", gcp.data),
					testAccCheckControlPlaneCloudAccountExists(gcp.resourceName, gcp.name),

					// NGS Resource Check
					resource.TestCheckResourceAttr(ngs.resourceName, "id", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "name", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "description", description),
					resource.TestCheckResourceAttr(ngs.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(ngs.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", ngs.name)),
					resource.TestCheckResourceAttr(ngs.resourceName, "ngs.0.secret_link", ngs.data),
					testAccCheckControlPlaneCloudAccountExists(ngs.resourceName, ngs.name),
				),
			},
			{
				Config: testAccControlPlaneCloudAccountUpdateAddTag(updateDescription, aws, azure, gcp, ngs),
				Check: resource.ComposeAggregateTestCheckFunc(
					// AWS Resource Check
					resource.TestCheckResourceAttr(aws.resourceName, "id", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "name", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(aws.resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(aws.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", aws.name)),
					resource.TestCheckResourceAttr(aws.resourceName, "aws.0.role_arn", aws.data),
					testAccCheckControlPlaneCloudAccountExists(aws.resourceName, aws.name),

					// Azure Resource Check
					resource.TestCheckResourceAttr(azure.resourceName, "id", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "name", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(azure.resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(azure.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", azure.name)),
					resource.TestCheckResourceAttr(azure.resourceName, "azure.0.secret_link", azure.data),
					testAccCheckControlPlaneCloudAccountExists(azure.resourceName, azure.name),

					// GCP Resource Check
					resource.TestCheckResourceAttr(gcp.resourceName, "id", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "name", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(gcp.resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(gcp.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", gcp.name)),
					resource.TestCheckResourceAttr(gcp.resourceName, "gcp.0.project_id", gcp.data),
					testAccCheckControlPlaneCloudAccountExists(gcp.resourceName, gcp.name),

					// NGS Resource Check
					resource.TestCheckResourceAttr(ngs.resourceName, "id", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "name", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(ngs.resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(ngs.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", ngs.name)),
					resource.TestCheckResourceAttr(ngs.resourceName, "ngs.0.secret_link", ngs.data),
					testAccCheckControlPlaneCloudAccountExists(ngs.resourceName, ngs.name),
				),
			},
			{
				Config: testAccControlPlaneCloudAccountUpdateWithOptionals(updateDescription, aws, azure, gcp, ngs),
				Check: resource.ComposeAggregateTestCheckFunc(
					// AWS Resource Check
					resource.TestCheckResourceAttr(aws.resourceName, "id", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "name", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(aws.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(aws.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", aws.name)),
					resource.TestCheckResourceAttr(aws.resourceName, "aws.0.role_arn", aws.data),
					testAccCheckControlPlaneCloudAccountExists(aws.resourceName, aws.name),

					// Azure Resource Check
					resource.TestCheckResourceAttr(azure.resourceName, "id", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "name", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(azure.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(azure.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", azure.name)),
					resource.TestCheckResourceAttr(azure.resourceName, "azure.0.secret_link", azure.data),
					testAccCheckControlPlaneCloudAccountExists(azure.resourceName, azure.name),

					// GCP Resource Check
					resource.TestCheckResourceAttr(gcp.resourceName, "id", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "name", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(gcp.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(gcp.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", gcp.name)),
					resource.TestCheckResourceAttr(gcp.resourceName, "gcp.0.project_id", gcp.data),
					testAccCheckControlPlaneCloudAccountExists(gcp.resourceName, gcp.name),

					// NGS Resource Check
					resource.TestCheckResourceAttr(ngs.resourceName, "id", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "name", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(ngs.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(ngs.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", ngs.name)),
					resource.TestCheckResourceAttr(ngs.resourceName, "ngs.0.secret_link", ngs.data),
					testAccCheckControlPlaneCloudAccountExists(ngs.resourceName, ngs.name),
				),
			},
			// Update Replace
			{
				Config: testAccControlPlaneCloudAccountUpdateReplace(updateDescription, aws, azure, gcp, ngs),
				Check: resource.ComposeAggregateTestCheckFunc(
					// AWS Resource Check
					resource.TestCheckResourceAttr(aws.resourceName, "id", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "name", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(aws.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(aws.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", aws.name)),
					resource.TestCheckResourceAttr(aws.resourceName, "aws.0.role_arn", aws.updateData),
					testAccCheckControlPlaneCloudAccountExists(aws.resourceName, aws.name),

					// Azure Resource Check
					resource.TestCheckResourceAttr(azure.resourceName, "id", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "name", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(azure.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(azure.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", azure.name)),
					resource.TestCheckResourceAttr(azure.resourceName, "azure.0.secret_link", azure.updateData),
					testAccCheckControlPlaneCloudAccountExists(azure.resourceName, azure.name),

					// GCP Resource Check
					resource.TestCheckResourceAttr(gcp.resourceName, "id", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "name", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(gcp.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(gcp.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", gcp.name)),
					resource.TestCheckResourceAttr(gcp.resourceName, "gcp.0.project_id", gcp.updateData),
					testAccCheckControlPlaneCloudAccountExists(gcp.resourceName, gcp.name),

					// NGS Resource Check
					resource.TestCheckResourceAttr(ngs.resourceName, "id", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "name", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(ngs.resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(ngs.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", ngs.name)),
					resource.TestCheckResourceAttr(ngs.resourceName, "ngs.0.secret_link", ngs.updateData),
					testAccCheckControlPlaneCloudAccountExists(ngs.resourceName, ngs.name),
				),
			},
			// Update to Required Only
			{
				Config: testAccControlPlaneCloudAccountCreateRequiredOnly(aws, azure, gcp, ngs),
				Check: resource.ComposeAggregateTestCheckFunc(
					// AWS Resource Check
					resource.TestCheckResourceAttr(aws.resourceName, "id", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "name", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "description", aws.name),
					resource.TestCheckResourceAttr(aws.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(aws.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", aws.name)),
					resource.TestCheckResourceAttr(aws.resourceName, "aws.0.role_arn", aws.data),
					testAccCheckControlPlaneCloudAccountExists(aws.resourceName, aws.name),

					// Azure Resource Check
					resource.TestCheckResourceAttr(azure.resourceName, "id", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "name", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "description", azure.name),
					resource.TestCheckResourceAttr(azure.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(azure.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", azure.name)),
					resource.TestCheckResourceAttr(azure.resourceName, "azure.0.secret_link", azure.data),
					testAccCheckControlPlaneCloudAccountExists(azure.resourceName, azure.name),

					// GCP Resource Check
					resource.TestCheckResourceAttr(gcp.resourceName, "id", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "name", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "description", gcp.name),
					resource.TestCheckResourceAttr(gcp.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(gcp.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", gcp.name)),
					resource.TestCheckResourceAttr(gcp.resourceName, "gcp.0.project_id", gcp.data),
					testAccCheckControlPlaneCloudAccountExists(gcp.resourceName, gcp.name),

					// NGS Resource Check
					resource.TestCheckResourceAttr(ngs.resourceName, "id", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "name", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "description", ngs.name),
					resource.TestCheckResourceAttr(ngs.resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(ngs.resourceName, "self_link", GetSelfLink(OrgName, "cloudaccount", ngs.name)),
					resource.TestCheckResourceAttr(ngs.resourceName, "ngs.0.secret_link", ngs.data),
					testAccCheckControlPlaneCloudAccountExists(ngs.resourceName, ngs.name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccCheckControlPlaneCloudAccountCheckDestroy verifies that all resources have been destroyed.
func testAccCheckControlPlaneCloudAccountCheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_cloud_account resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_cloud_account" {
			continue
		}

		// Retrieve the name for the current resource
		cloudAccountName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of cloud account with name: %s", cloudAccountName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		cloudAccount, code, err := TestProvider.client.GetCloudAccount(cloudAccountName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if cloud account %s exists: %w", cloudAccountName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if cloudAccount != nil {
			return fmt.Errorf("CheckDestroy failed: cloud account %s still exists in the system", *cloudAccount.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_cloud_account resources have been successfully destroyed")
	return nil
}

// testAccCheckControlPlaneCloudAccountExists verifies that a specified resource exist within the Terraform state and in the data service.
func testAccCheckControlPlaneCloudAccountExists(resourceName string, cloudAccountName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of cloud account: %s. Total resources: %d", cloudAccountName, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != cloudAccountName {
			return fmt.Errorf("resource ID %s does not match expected cloud account name %s", rs.Primary.ID, cloudAccountName)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteCloudAccount, _, err := TestProvider.client.GetCloudAccount(cloudAccountName)
		if err != nil {
			return fmt.Errorf("error retrieving cloud account from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteCloudAccount.Name != cloudAccountName {
			return fmt.Errorf("mismatch in cloud account name: expected %s, got %s", cloudAccountName, *remoteCloudAccount.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("Cloud Account %s verified successfully in both state and external system.", cloudAccountName))
		return nil
	}
}

/*** Configs ***/

// testAccControlPlaneCloudAccountCreateRequiredOnly constructs HCL for creating cloud accounts with only required fields for AWS, Azure, GCP, and NGS.
func testAccControlPlaneCloudAccountCreateRequiredOnly(aws CloudAccountResourceTest, azure CloudAccountResourceTest, gcp CloudAccountResourceTest, ngs CloudAccountResourceTest) string {
	return fmt.Sprintf(`
resource "cpln_cloud_account" "aws" {
  name = "%s"

  aws {
    role_arn = "%s"
  }
}

resource "cpln_cloud_account" "azure" {
  name = "%s"

  azure {
    secret_link = "%s"
  }
}

resource "cpln_cloud_account" "gcp" {
  name = "%s"

  gcp {
    project_id = "%s"
  }
}

resource "cpln_cloud_account" "ngs" {
  name = "%s"

  ngs {
    secret_link = "%s"
  }
}
`, aws.name, aws.data, azure.name, azure.data, gcp.name, gcp.data, ngs.name, ngs.data)
}

// testAccControlPlaneCloudAccountUpdateWithOptionals constructs HCL for updating cloud accounts including description and tags for AWS, Azure, GCP, and NGS.
func testAccControlPlaneCloudAccountUpdateWithOptionals(description string, aws CloudAccountResourceTest, azure CloudAccountResourceTest, gcp CloudAccountResourceTest, ngs CloudAccountResourceTest) string {
	return fmt.Sprintf(`
resource "cpln_cloud_account" "aws" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  aws {
    role_arn = "%s"
  }
}

resource "cpln_cloud_account" "azure" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  azure {
    secret_link = "%s"
  }
}

resource "cpln_cloud_account" "gcp" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gcp {
    project_id = "%s"
  }
}

resource "cpln_cloud_account" "ngs" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  ngs {
    secret_link = "%s"
  }
}
`, aws.name, description, aws.data, azure.name, description, azure.data, gcp.name, description, gcp.data, ngs.name, description, ngs.data)
}

// testAccControlPlaneCloudAccountUpdateAddTag constructs HCL to update cloud accounts by adding a new tag alongside existing tags for all providers.
func testAccControlPlaneCloudAccountUpdateAddTag(description string, aws CloudAccountResourceTest, azure CloudAccountResourceTest, gcp CloudAccountResourceTest, ngs CloudAccountResourceTest) string {
	return fmt.Sprintf(`
resource "cpln_cloud_account" "aws" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    new_tag             = "true"
  }

  aws {
    role_arn = "%s"
  }
}

resource "cpln_cloud_account" "azure" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    new_tag             = "true"
  }

  azure {
    secret_link = "%s"
  }
}

resource "cpln_cloud_account" "gcp" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    new_tag             = "true"
  }

  gcp {
    project_id = "%s"
  }
}

resource "cpln_cloud_account" "ngs" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    new_tag             = "true"
  }

  ngs {
    secret_link = "%s"
  }
}
`, aws.name, description, aws.data, azure.name, description, azure.data, gcp.name, description, gcp.data, ngs.name, description, ngs.data)
}

// testAccControlPlaneCloudAccountUpdateReplace constructs HCL to update cloud accounts by replacing provider credentials for all providers.
func testAccControlPlaneCloudAccountUpdateReplace(description string, aws CloudAccountResourceTest, azure CloudAccountResourceTest, gcp CloudAccountResourceTest, ngs CloudAccountResourceTest) string {
	return fmt.Sprintf(`
resource "cpln_cloud_account" "aws" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  aws {
    role_arn = "%s"
  }
}

resource "cpln_cloud_account" "azure" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  azure {
    secret_link = "%s"
  }
}

resource "cpln_cloud_account" "gcp" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gcp {
    project_id = "%s"
  }
}

resource "cpln_cloud_account" "ngs" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  ngs {
    secret_link = "%s"
  }
}
`, aws.name, description, aws.updateData, azure.name, description, azure.updateData, gcp.name, description, gcp.updateData, ngs.name, description, ngs.updateData)
}
