package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneAuditContext_basic(t *testing.T) {

	var testAuditContext client.AuditContext
	randomName := "audit-context-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "AUDIT-CONTEXT") },
		Providers:    testAccProviders,
		CheckDestroy: testAccControlPlaneAuditContextCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneAuditContext(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneAuditContextExists("cpln_audit_context.tf-audit-context", randomName, &testAuditContext),
					testAccCheckControlPlaneAuditContextAttributes(&testAuditContext),
					resource.TestCheckResourceAttr("cpln_audit_context.tf-audit-context", "name", randomName),
					resource.TestCheckResourceAttr("cpln_audit_context.tf-audit-context", "description", "audit context description "+randomName),
				),
			},
		},
	})
}

func testAccControlPlaneAuditContextCheckDestroy(s *terraform.State) error {

	return nil
}

func testAccControlPlaneAuditContext(name string) string {

	return fmt.Sprintf(`
	
	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_audit_context" "tf-audit-context" {
		name = var.random-name
		description = "audit context description ${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}
	}
	`, name)
}

func testAccCheckControlPlaneAuditContextExists(resourceName string, auditCtxName string, auditCtx *client.AuditContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneAuditContextExists. Resources Length: %d", len(s.RootModule().Resources))
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != auditCtxName {
			return fmt.Errorf("Audit Context name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)
		ac, _, err := client.GetAuditContext(auditCtxName)

		if err != nil {
			return err
		}

		if *ac.Name != auditCtxName {
			return fmt.Errorf("Audit Context name does not match")
		}

		*auditCtx = *ac

		return nil
	}
}

func testAccCheckControlPlaneAuditContextAttributes(auditCtx *client.AuditContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tags := *auditCtx.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("Tags - audit context terraform_generated attribute does not match")
		}

		return nil
	}
}
