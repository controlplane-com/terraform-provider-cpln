package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneAgent_basic(t *testing.T) {

	var testAgent client.Agent

	aName := "agent-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "AGENT") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneAgentCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneAgent(aName, "Agent created using terraform for acceptance tests"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneAgentExists("cpln_agent.new", aName, &testAgent),
					resource.TestCheckResourceAttr("cpln_agent.new", "description", "Agent created using terraform for acceptance tests"),
				),
			},
		},
	})
}

func testAccCheckControlPlaneAgentExists(resourceName, agentName string, agent *client.Agent) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneAgentExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != agentName {
			return fmt.Errorf("Agent name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)

		wl, _, err := client.GetAgent(agentName)

		if err != nil {
			return err
		}

		if *wl.Name != agentName {
			return fmt.Errorf("Agent name does not match")
		}

		*agent = *wl

		return nil
	}
}

func testAccControlPlaneAgent(agentName, agentDescription string) string {

	TestLogger.Printf("Inside testAccControlPlaneAgent")

	return fmt.Sprintf(`

	resource "cpln_agent" "new" {
		name        = "%s"	
		description = "%s"
  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}
	  	
	`, agentName, agentDescription)
}

func testAccCheckControlPlaneAgentCheckDestroy(s *terraform.State) error {

	// TestLogger.Printf("Inside testAccCheckControlPlaneAgentCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneAgentCheckDestroy: rs.Type: %s", rs.Type)

		if rs.Type != "cpln_agent" {
			continue
		}

		agentName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneAgentCheckDestroy: agent name: %s", agentName)

		agent, _, _ := c.GetAgent(agentName)
		if agent != nil {
			return fmt.Errorf("Agent still exists. Name: %s", *agent.Name)
		}
	}

	return nil
}
