package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

/*** Acceptance Test ***/

// TestAccControlPlaneAuditContext_basic performs an acceptance test for the resource.
func TestAccControlPlaneAuditContext_basic(t *testing.T) {
	// Initialize a variable to store the API resource retrieved during the test steps
	var testAuditCtx client.AuditContext

	// Define unique values for the API resource to be used during the test lifecycle
	name := "audit-context-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	description := "Audit Context created using terraform for acceptance tests"
	updateDescription := "Audit Context updated using terraform for acceptance tests"
	resourceName := "cpln_audit_context.new"

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "AUDIT_CONTEXT") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccControlPlaneAuditContextCreateRequiredOnly(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "auditctx", name)),
					testAccCheckControlPlaneAuditContextExists(resourceName, name, &testAuditCtx),
				),
			},
			// ImportState testing
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: testAccControlPlaneAuditContextUpdateWithOptionals(name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "auditctx", name)),
				),
			},
			{
				Config: testAccControlPlaneAuditContextUpdateAddTag(name, updateDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "tags.new_tag", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "auditctx", name)),
				),
			},
			{
				Config: testAccControlPlaneAuditContextUpdateRemoveTag(name, updateDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "auditctx", name)),
				),
			},
			{
				Config: testAccControlPlaneAuditContextCreateRequiredOnly(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "auditctx", name)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccCheckControlPlaneAuditContextExists verifies that a specified resource exist within the Terraform state and in the data service.
func testAccCheckControlPlaneAuditContextExists(resourceName string, auditCtxName string, auditCtx *client.AuditContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of audit context: %s. Total resources: %d", auditCtxName, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != auditCtxName {
			return fmt.Errorf("resource ID %s does not match expected audit context name %s", rs.Primary.ID, auditCtxName)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteAuditCtx, _, err := TestProvider.client.GetAuditContext(auditCtxName)
		if err != nil {
			return fmt.Errorf("error retrieving audit context from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected API resource name
		if *remoteAuditCtx.Name != auditCtxName {
			return fmt.Errorf("mismatch in audit context name: expected %s, got %s", auditCtxName, *remoteAuditCtx.Name)
		}

		// Copy the retrieved API resource data to the pointer provided, for further use in tests
		*auditCtx = *remoteAuditCtx

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("Audit Context %s verified successfully in both state and external system.", auditCtxName))
		return nil
	}
}

/*** Configs ***/

// testAccControlPlaneAuditContextCreateRequiredOnly constructs HCL for creating an audit context with only the required name field.
func testAccControlPlaneAuditContextCreateRequiredOnly(name string) string {
	return fmt.Sprintf(`
resource "cpln_audit_context" "new" {
  name = "%s"
}
`, name)
}

// testAccControlPlaneAuditContextUpdateWithOptionals constructs HCL to update an audit context including description and tags.
func testAccControlPlaneAuditContextUpdateWithOptionals(name string, description string) string {
	return fmt.Sprintf(`
resource "cpln_audit_context" "new" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}
`, name, description)
}

// testAccControlPlaneAuditContextUpdateAddTag constructs HCL to add a new tag to an audit context.
func testAccControlPlaneAuditContextUpdateAddTag(name string, description string) string {
	return fmt.Sprintf(`
resource "cpln_audit_context" "new" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
	  new_tag             = "true"
  }
}
`, name, description)
}

// testAccControlPlaneAuditContextUpdateRemoveTag constructs HCL to update an audit context by removing custom tags.
func testAccControlPlaneAuditContextUpdateRemoveTag(name string, description string) string {
	return fmt.Sprintf(`
resource "cpln_audit_context" "new" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}
`, name, description)
}
