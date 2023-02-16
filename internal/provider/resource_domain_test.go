package cpln

import (
	"fmt"
	"os"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func ToBeUpdated_TestAccControlPlaneDomain_basic(t *testing.T) {

	aName := "globalvirtualcloud.com"
	dName := "domain-testacc.globalvirtualcloud.com"
	zone := "cpln-test"

	ep := resource.ExternalProvider{
		Source:            "google",
		VersionConstraint: "3.72.0",
	}

	eps := map[string]resource.ExternalProvider{
		"google": ep,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckGoogle(t, "DOMAIN") },
		Providers:         testAccProviders,
		ExternalProviders: eps,
		CheckDestroy:      testAccCheckControlPlaneDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneDomain(aName, dName, "Domain created using Terraform", zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_domain.new", "name", dName),
					resource.TestCheckResourceAttr("cpln_domain.new", "description", "Domain created using Terraform"),
				),
			},
			{
				Config: testAccControlPlaneDomain(aName, dName, "Domain created using Terraform - Updated", zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_domain.new", "name", dName),
					resource.TestCheckResourceAttr("cpln_domain.new", "description", "Domain created using Terraform - Updated"),
				),
			},
		},
	})
}

func testAccCheckControlPlaneDomainDestroy(s *terraform.State) error {

	// TODO: Test that DNS records have been removed

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_domain" {
			continue
		}

		domainName := rs.Primary.ID

		domain, _, _ := c.GetDomain(domainName)
		if domain != nil {
			return fmt.Errorf("Domain still exists. Name: %s", *domain.Name)
		}
	}

	return nil
}

func testAccControlPlaneDomain(apex, domain, description, managedZone string) string {

	variables := fmt.Sprintf(`

		variable domain_name_apex {
			type = string
			default = "%s"
		}

		variable domain_name {
			type = string
			default = "%s"
		}

		variable managed_zone {
			type = string
			default = "%s"
		}

		variable description {
			type = string
			default = "%s"
		}
		`, apex, domain, managedZone, description)

	domainSetup := `
	
	data "cpln_org" "org" {}
	 
	resource "google_dns_record_set" "ns" {

		name         = "${var.domain_name}."
		managed_zone = var.managed_zone
		type         = "NS"
		ttl          = 1800
	  
		rrdatas = ["ns1.cpln.io.", "ns2.cpln.io.", "ns3.cpln.io.", "ns4.cpln.io."]
	}
	  
	resource "google_dns_record_set" "txt" {

		name         = "_cpln-${google_dns_record_set.ns.name}"
		managed_zone = var.managed_zone
		type         = "TXT"
		ttl          = 600
	  
		rrdatas = [data.cpln_org.org.id]
	}

	resource "cpln_domain" "new" {

		depends_on = [google_dns_record_set.ns, google_dns_record_set.txt]
		
		name        = var.domain_name	
		description = var.description
	
		tags = {
		terraform_generated = "true"
		acceptance_test = "true"
		}
	}
	`

	cplnDomainSetup := `

	    resource "cpln_domain" "new_apex" {
			
			name = var.domain_name_apex	
		
			tags = {
			  terraform_generated = "true"
			  acceptance_test = "true"
			}
	    }

		resource "cpln_domain" "new" {
			
			depends_on = [cpln_domain.new_apex]

			name        = var.domain_name	
			description = var.description
		
			tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			}
		}
	`

	if validateDomains := os.Getenv("VALIDATE_DOMAINS"); validateDomains == "false" {
		return fmt.Sprintf("%s %s", variables, cplnDomainSetup)
	}

	return fmt.Sprintf("%s %s", variables, domainSetup)
}
