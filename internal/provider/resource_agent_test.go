package cpln

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

/*** Acceptance Test ***/

// TestAccControlPlaneAgent_basic performs an acceptance test for the resource.
func TestAccControlPlaneAgent_basic(t *testing.T) {
	// Initialize a variable to store the API resource retrieved during the test steps
	var testAgent client.Agent

	// Define unique values for the API resource to be used during the test lifecycle
	name := "agent-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	description := "Agent created using terraform for acceptance tests"
	updateDescription := "Agent updated using terraform for acceptance tests"
	resourceName := "cpln_agent.new"

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "AGENT") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             testAccCheckControlPlaneAgentCheckDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccControlPlaneAgentCreateRequiredOnly(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrWith(resourceName, "user_data", testAccControlPlaneAgentValidateUserData),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "agent", name)),
					testAccCheckControlPlaneAgentExists(resourceName, name, &testAgent),
				),
			},
			// ImportState testing
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: testAccControlPlaneAgentUpdateWithOptionals(name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "agent", name)),
					resource.TestCheckResourceAttrWith(resourceName, "user_data", testAccControlPlaneAgentValidateUserData),
				),
			},
			{
				Config: testAccControlPlaneAgentUpdateAddTag(name, updateDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "tags.new_tag", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "agent", name)),
					resource.TestCheckResourceAttrWith(resourceName, "user_data", testAccControlPlaneAgentValidateUserData),
				),
			},
			{
				Config: testAccControlPlaneAgentUpdateRemoveTag(name, updateDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "agent", name)),
					resource.TestCheckResourceAttrWith(resourceName, "user_data", testAccControlPlaneAgentValidateUserData),
				),
			},
			{
				Config: testAccControlPlaneAgentCreateRequiredOnly(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "self_link", GetSelfLink(OrgName, "agent", name)),
					resource.TestCheckResourceAttrWith(resourceName, "user_data", testAccControlPlaneAgentValidateUserData),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccCheckControlPlaneAgentCheckDestroy verifies that all resources have been destroyed.
func testAccCheckControlPlaneAgentCheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_agent resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_agent" {
			continue
		}

		// Retrieve the API resource name for the current resource
		agentName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of agent with name: %s", agentName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		agent, code, err := TestProvider.client.GetAgent(agentName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if agent %s exists: %w", agentName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if agent != nil {
			return fmt.Errorf("CheckDestroy failed: agent %s still exists in the system", *agent.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_agent resources have been successfully destroyed")
	return nil
}

// testAccCheckControlPlaneAgentExists verifies that a specified resource exist within the Terraform state and in the data service.
func testAccCheckControlPlaneAgentExists(resourceName string, agentName string, agent *client.Agent) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of agent: %s. Total resources: %d", agentName, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != agentName {
			return fmt.Errorf("resource ID %s does not match expected agent name %s", rs.Primary.ID, agentName)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteAgent, _, err := TestProvider.client.GetAgent(agentName)
		if err != nil {
			return fmt.Errorf("error retrieving agent from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected API resource name
		if *remoteAgent.Name != agentName {
			return fmt.Errorf("mismatch in agent name: expected %s, got %s", agentName, *remoteAgent.Name)
		}

		// Copy the retrieved API resource data to the pointer provided, for further use in tests
		*agent = *remoteAgent

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("Agent %s verified successfully in both state and external system.", agentName))
		return nil
	}
}

// testAccControlPlaneAgentValidateUserData checks that the attribute value contains a non and valid user data.
func testAccControlPlaneAgentValidateUserData(v string) error {
	// Ensure both required substrings are present.
	if !strings.Contains(v, "registrationToken") || !strings.Contains(v, "agentLink") {
		return fmt.Errorf(
			"attribute value must contain both 'registrationToken' and 'agentLink'; got %q",
			v,
		)
	}

	// If we got here, then the user data is valid
	return nil
}

/*** Configs ***/

// testAccControlPlaneAgentCreateRequiredOnly constructs HCL for creating an agent resource with only the required name field.
func testAccControlPlaneAgentCreateRequiredOnly(name string) string {
	return fmt.Sprintf(`
resource "cpln_agent" "new" {
  name = "%s"
}
`, name)
}

// testAccControlPlaneAgentUpdateWithOptionals constructs HCL to update an agent resource including description and tags.
func testAccControlPlaneAgentUpdateWithOptionals(name string, description string) string {
	return fmt.Sprintf(`
resource "cpln_agent" "new" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}
`, name, description)
}

// testAccControlPlaneAgentUpdateAddTag constructs HCL to update an agent resource by adding a new tag.
func testAccControlPlaneAgentUpdateAddTag(name string, description string) string {
	return fmt.Sprintf(`
resource "cpln_agent" "new" {
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

// testAccControlPlaneAgentUpdateRemoveTag constructs HCL to update an agent resource by removing custom tags.
func testAccControlPlaneAgentUpdateRemoveTag(name string, description string) string {
	return fmt.Sprintf(`
resource "cpln_agent" "new" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }
}
`, name, description)
}
