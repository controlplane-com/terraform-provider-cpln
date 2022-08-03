package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneOrgLogging_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "ORG") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneOrgCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneOrgS3(),
			},
			{
				Config: testAccControlPlaneOrgCoralogix(),
			},
			{
				Config: testAccControlPlaneOrgDatadog(),
			},
			{
				Config: testAccControlPlaneOrgLogzio(),
			},
			{
				Config: testAccControlPlaneOrgLogzio1(),
			},
		},
	})
}

func testAccControlPlaneOrgS3() string {

	TestLogger.Printf("Inside testAccControlPlaneOrg")

	return `

	resource "cpln_secret" "aws" {
		name = "aws-random-tbd"
		description = "aws description aws-random-tbd" 
				
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "aws"
		} 
		
		aws {
			secret_key = "AKIAIOSFODNN7EXAMPLE"
			access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
			role_arn = "arn:awskey" 
		}
	}

	resource "cpln_org_logging" "new" {

		s3_logging {

			bucket = "test-bucket"
			region = "us-east1"
			prefix = "/"

			// AWS Secret Only
			credentials = cpln_secret.aws.self_link
		}	
	}
	`
}

func testAccControlPlaneOrgCoralogix() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgCoralogix")

	return `

	resource "cpln_secret" "opaque" {

		name = "opaque-random-coralogix-tbd"
		description = "opaque description opaque-random-tbd" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		}

		opaque {
			payload = "opaque_secret_payload"
			encoding = "plain"
		}
	}

	resource "cpln_org_logging" "new" {

		coralogix_logging {

			// Valid clusters
			// coralogix.com, coralogix.us, app.coralogix.in, app.eu2.coralogix.com, app.coralogixsg.com,
			cluster = "coralogix.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque.self_link
			
			// Supported variables for App and Subsystem are:
			// {org}, {gvc}, {workload}, {location}
			app = "{workload}"
			subsystem = "{org}"
		}
	}	  	
	`
}

func testAccControlPlaneOrgDatadog() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgDatadog")

	return `

	resource "cpln_secret" "opaque" {

		name = "opaque-random-datadog-tbd"
		description = "opaque description" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		}

		opaque {
			payload = "opaque_secret_payload"
			encoding = "plain"
		}
	}

	resource "cpln_org_logging" "new" {

		datadog_logging {

			// Valid Host
			// http-intake.logs.datadoghq.com, http-intake.logs.us3.datadoghq.com, 
			// http-intake.logs.us5.datadoghq.com, http-intake.logs.datadoghq.eu
			host = "http-intake.logs.datadoghq.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque.self_link	
		}
	}	  	
	`
}

func testAccControlPlaneOrgLogzio() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgLogzio")

	return `

	resource "cpln_secret" "opaque" {

		name = "opaque-random-datadog-tbd"
		description = "opaque description" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		}

		opaque {
			payload = "opaque_secret_payload"
			encoding = "plain"
		}
	}

	resource "cpln_org_logging" "new" {

		logzio_logging {

			// Valid Host
			// listener.logz.io, listener-nl.logz.io 
			listener_host = "listener.logz.io"

			// Opaque Secret Only
			credentials = cpln_secret.opaque.self_link	
		}
	}	  	
	`
}

func testAccControlPlaneOrgLogzio1() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgLogzio1")

	return `

	resource "cpln_secret" "opaque" {

		name = "opaque-random-datadog-tbd"
		description = "opaque description" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		}

		opaque {
			payload = "opaque_secret_payload"
			encoding = "plain"
		}
	}

	resource "cpln_org_logging" "new" {

		logzio_logging {

			// Valid Host
			// listener.logz.io, listener-nl.logz.io 
			listener_host = "listener-nl.logz.io"

			// Opaque Secret Only
			credentials = cpln_secret.opaque.self_link	
		}
	}	  	
	`
}

func testAccCheckControlPlaneOrgCheckDestroy(s *terraform.State) error {

	// TestLogger.Printf("Inside testAccCheckControlPlaneOrgCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneOrgCheckDestroy: rs.Type: %s", rs.Type)

		if rs.Type != "cpln_org_logging" {
			continue
		}

		orgName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneOrgCheckDestroy: Org name: %s", orgName)

		org, _, _ := c.GetOrg()

		if org.Spec.Logging != nil {
			return fmt.Errorf("Org Spec Logging still exists. Org Name: %s", *org.Name)
		}
	}

	return nil
}
