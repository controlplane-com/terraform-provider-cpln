package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneDomain_basic(t *testing.T) {

	var testDomain client.Domain

	dName := "cors2.cplntest.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "DOMAIN") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneDomainCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneDomain(dName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists(dName, &testDomain),
				),
			},
		},
	})
}

func testAccControlPlaneDomain(domainName string) string {

	TestLogger.Printf("Inside testAccControlPlaneDomain")

	return fmt.Sprintf(`

	resource "cpln_domain" "example" {

		name        = "%s"
		description = "Test hakan"
	
		tags = {
			terraform_generated = "true"
			example             = "true"
		}
	
		spec {
			dns_mode         = "ns"
			accept_all_hosts = "true"
	
			ports {
				number   = 443
				protocol = "http"
	
				routes {
					prefix        = "/log"
					workload_link = "/org/efe/gvc/kadir/workload/a-log"
					port          = 8080
				}
	
				routes {
					prefix        = "/canary"
					workload_link = "/org/efe/gvc/kadir/workload/canary"
					port          = 8080
				}
			}
		}
	}	
	`, domainName)
}

func testAccCheckControlPlaneDomainExists(domainName string, domain *client.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[domainName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != domainName {
			return fmt.Errorf("Workload name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)

		d, _, err := client.GetDomain(domainName)

		if err != nil {
			return err
		}

		if *d.Name != domainName {
			return fmt.Errorf("Workload name does not match")
		}

		*domain = *d

		return nil
	}
}

func testAccCheckControlPlaneDomainCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneDomainCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	return nil
}

// Unit Tests

// Certificate With Secret
func TestControlPlane_BuildCertificate(t *testing.T) {
	secret := "/org/myorg/secret/mysecret"
	cert := buildCertificate(generateFlatTestCertificate(secret))

	expectedCert := client.DomainCertificate{SecretLink: &secret}

	// TODO move expectedCert to a function, can be array of items too for different cases
	if diff := deep.Equal(cert, &expectedCert); diff != nil {
		t.Errorf("Domain Certificate was not built correctly. Diff: %s", diff)
	}
}

func generateFlatTestCertificate(secretLink string) []interface{} {
	spec := map[string]interface{}{
		"secret_link": secretLink,
	}

	return []interface{}{
		spec,
	}
}

// Certificate Without Secret
func TestControlPlane_BuildCertificate_WithoutSecret(t *testing.T) {
	cert := buildCertificate(generateFlatTestCertificateWithoutSecret())
	certTest := client.DomainCertificate{}

	if diff := deep.Equal(cert, &certTest); diff != nil {
		t.Errorf("Domain Certificate was not built correctly. Diff: %s", diff)
	}
}

func generateFlatTestCertificateWithoutSecret() []interface{} {
	spec := map[string]interface{}{}

	return []interface{}{
		spec,
	}
}
