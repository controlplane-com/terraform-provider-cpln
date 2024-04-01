package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func getGVC() string {
	gvc := "/org/terraform-test-org/gvc/domain-test-gvc"
	return gvc
}

func getWorkloadOne() string {
	workload := getGVC() + "/workload/wl1"
	return workload
}

func getWorkloadTwo() string {
	workload := getGVC() + "/workload/wl2"
	return workload
}

// func getDomainOne() string {
// 	domain := "domain-test.example.com"
// 	return domain
// }

// func getDomainTwo() string {
// 	domain := "example.example.com"
// 	return domain
// }

// func getDomainThree() string {
// 	domain := "example2.example.com"
// 	return domain
// }

// func getDomainFour() string {
// 	domain := "example3.example.com"
// 	return domain
// }

// TODO: Once the pipline is configured to run the acc tests, it will be set to not validate the apex txt record and we can use '.example.com'
// func getTestApex() string {
// 	return ".erickotler.com"
// }

func getTestApex() string {
	return "erickotler.com"
}

func NeedToFix_TestAccControlPlaneDomain_basic(t *testing.T) {

	var domain client.Domain
	var org client.Org

	randomName := "domain-acctest-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	domainName := randomName + "." + getTestApex()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "DOMAIN") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneDomainCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainApexClean(getTestApex(), getTestApex()+" Description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.domain_apex", getTestApex(), &domain, &org),
				),
			},
			{
				Config: testAccDomainApexClean(getTestApex(), getTestApex()+" Description Updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.domain_apex", getTestApex(), &domain, &org),
				),
			},
			{
				Config: testAccDomainApex(randomName, getTestApex(), getTestApex()+" Description Updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.domain_apex", getTestApex(), &domain, &org),
				),
			},
			{
				Config: testAccDomainApex(randomName, getTestApex(), getTestApex()+" Description Updated Again"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.domain_apex", getTestApex(), &domain, &org),
				),
			},
			{
				Config: testAccControlPlaneDomainSubdomain(randomName, getTestApex(), getTestApex()+" Description", domainName, "ns"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.subdomain", domainName, &domain, &org),
					// testAccCheckControlPlaneDomainNSSubdomain(&domain, &org, "gvc-"+randomName),
				),
			},
			{
				Config: testAccControlPlaneDomainSubdomain(randomName, getTestApex(), getTestApex()+" Description - Updated", domainName, "ns"),
			},
			{
				Config: testAccControlPlaneDomainPathBased(randomName, getTestApex(), getTestApex()+" Description", domainName, "ns"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.subdomain", domainName, &domain, &org),
					// testAccCheckControlPlaneDomainNSPathBased(&domain, &org, randomName),
				),
			},
			{
				Config: testAccControlPlaneDomainPathBased(randomName, getTestApex(), getTestApex()+" Description - Updated", domainName, "ns"),
			},
			{
				Config: testAccControlPlaneDomainPathBased(randomName, getTestApex(), getTestApex()+" Description", domainName, "cname"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.subdomain", domainName, &domain, &org),
					// testAccCheckControlPlaneDomainNSPathBased(&domain, &org, randomName),
				),
			},
			{
				Config: testAccControlPlaneDomainPathBasedUpdateRoutePort(randomName, getTestApex(), getTestApex()+" Description - Updated", domainName, "cname"),
			},
		},
	})
}

func testAccDomainApexClean(domain, description string) string {

	TestLogger.Printf("Inside testAccDomainApex")

	return fmt.Sprintf(`

	resource "cpln_domain" "domain_apex" {
		name        = "%s"
		description = "%s"
	  
		tags = {
		  terraform_generated = "true"
		}

		spec {
			ports {
				tls { }
			 }
		}
		
	}`, domain, description)
}

func testAccDomainApex(random, domain, description string) string {

	TestLogger.Printf("Inside testAccDomainApex")

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "domain_gvc" {

		name        = "gvc-${var.random-name}"
		description = "Example GVC"

		locations = ["aws-eu-central-1", "aws-us-west-2"]

		tags = {
		  terraform_generated = "true"
		}
	}

	resource "cpln_domain" "domain_apex" {
		name        = "%s"
		description = "%s"
	  
		tags = {
		  terraform_generated = "true"
		}

		spec {

			dns_mode         = "cname"
		  	gvc_link         = cpln_gvc.domain_gvc.self_link
			accept_all_hosts = false

			ports {

				number = 443
				protocol = "http2"

				cors {

					allow_origins {						
						exact = "*"
					}

					allow_origins {						
						exact = "*.erickotler.com"
					}

					allow_methods = ["GET", "OPTIONS", "POST"]
					allow_headers = ["authorization", "host"]
					expose_headers = ["accept/type"]
					max_age = "12h"
					allow_credentials = true
				}

				tls {

					min_protocol_version = "TLSV1_1"
					cipher_suites = ["AES256-GCM-SHA384"]

					client_certificate {
						secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
					}

					server_certificate {
						secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
					}
				}
			}
		}
		
	}`, random, domain, description)
}

func testAccControlPlaneDomainSubdomain(random, apex, description, domain, dnsMode string) string {

	TestLogger.Printf("Inside testAccControlPlaneDomain")

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "domain_gvc" {

		name        = "gvc-${var.random-name}"
		description = "Example GVC"

		locations = ["aws-eu-central-1", "aws-us-west-2"]

		tags = {
		  terraform_generated = "true"
		}
	}

	resource "cpln_domain" "domain_apex" {
		name        = "%s"
		description = "%s"
	  
		tags = {
		  terraform_generated = "true"
		}

		spec {
			ports {
				tls { }
			 }
		}	
	}

	resource "cpln_domain" "subdomain" {

		depends_on = [cpln_domain.domain_apex]

		name        = "%s"
		description = "NS - Subdomain Based"
	  
		tags = {
		  terraform_generated = "true"
		}
	  
		spec {
		  dns_mode         = "%s"
		  gvc_link         = cpln_gvc.domain_gvc.self_link
	  
		  ports {
				number   = 443
				protocol = "http"
			
				cors {
					allow_origins {
						exact = "example.com"
					}

					allow_origins {
						exact = "*"
					}
			
					allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
					allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
					expose_headers    = ["expose_header_1", "expose_header_2", "expose_header_3"]
					max_age           = "24h"
					allow_credentials = "true"
				}
	  
				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
						"ECDHE-ECDSA-CHACHA20-POLY1305",
						"ECDHE-ECDSA-AES128-GCM-SHA256",
						"ECDHE-RSA-AES256-GCM-SHA384",
						"ECDHE-RSA-CHACHA20-POLY1305",
						"ECDHE-RSA-AES128-GCM-SHA256",
						"AES256-GCM-SHA384",
						"AES128-GCM-SHA256",
					]
					client_certificate {}
				}
			}
		}
	}`, random, apex, description, domain, dnsMode)
}

func testAccControlPlaneDomainPathBased(random, apex, description, domain, dnsMode string) string {

	TestLogger.Printf("Inside testAccControlPlaneDomain")

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "domain_gvc" {

		name        = "gvc-${var.random-name}"
		description = "Example GVC"

		locations = ["aws-eu-central-1", "aws-us-west-2"]

		tags = {
		  terraform_generated = "true"
		}
	}

	resource "cpln_workload" "new" {

		gvc = cpln_gvc.domain_gvc.name

		name        = "workload-${var.random-name}"
		description = "Example Workload"
		type        = "serverless"

		tags = {
		  terraform_generated = "true"
		}

		container {
		  name   = "container-01"
		  image  = "gcr.io/knative-samples/helloworld-go"
		  port   = 8080
		  memory = "128Mi"
		  cpu    = "50m"
		}

		options {
		  capacity_ai     = false
		  timeout_seconds = 30
		  suspend         = true

		  autoscaling {
			metric          = "concurrency"
			target          = 100
			max_scale       = 0
			min_scale       = 0
			max_concurrency = 500
		  }
		}
	}

	resource "cpln_domain" "domain_apex" {
		name        = "%s"
		description = "%s"
	  
		tags = {
		  terraform_generated = "true"
		}

		spec {
			ports {
				tls { }
			 }
		}	
	}

	resource "cpln_domain" "subdomain" {
		
		depends_on = [cpln_domain.domain_apex]

		name        = "%s"
		description = "NS - Path Based"
	  
		tags = {
		  terraform_generated = "true"
		}
	  
		spec {

			dns_mode = "%s"
			
			ports {
				number   = 443
				protocol = "http"

				cors {
					allow_origins {
						exact = "example.com"
					}

					allow_origins {
						exact = "*"
					}
			
					allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
					allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
					expose_headers     = ["expose_header_1", "expose_header_2", "expose_header_3"]
					max_age           = "24h"
					allow_credentials = "true"
				}
	  
				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
						"ECDHE-ECDSA-CHACHA20-POLY1305",
						"ECDHE-ECDSA-AES128-GCM-SHA256",
						"ECDHE-RSA-AES256-GCM-SHA384",
						"ECDHE-RSA-CHACHA20-POLY1305",
						"ECDHE-RSA-AES128-GCM-SHA256",
						"AES256-GCM-SHA384",
						"AES128-GCM-SHA256",
					]
				}
			}

			ports {
				number   = 80
				protocol = "http"

				cors {
					allow_origins {
						exact = "example.com"
					}

					allow_origins {
						exact = "*"
					}
			
					allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
					allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
					expose_headers     = ["expose_header_1", "expose_header_2", "expose_header_3"]
					max_age           = "24h"
					allow_credentials = "true"
				}
	  
				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
						"ECDHE-ECDSA-CHACHA20-POLY1305",
						"ECDHE-ECDSA-AES128-GCM-SHA256",
						"ECDHE-RSA-AES256-GCM-SHA384",
						"ECDHE-RSA-CHACHA20-POLY1305",
						"ECDHE-RSA-AES128-GCM-SHA256",
						"AES256-GCM-SHA384",
						"AES128-GCM-SHA256",
					]
				}
			}
		}
	}
	
	resource "cpln_domain_route" "route_first" {

		domain_link = cpln_domain.subdomain.self_link
		// domain_port = 443

		prefix = "/first"
		workload_link = cpln_workload.new.self_link
		host_prefix   = "my.thing."
		// port = 80
	}

	resource "cpln_domain_route" "route_second" {

		depends_on = [cpln_domain_route.route_first]
		
		domain_link = cpln_domain.subdomain.self_link
		// domain_port = 443

		prefix = "/second"
		replace_prefix = "/"
		workload_link = cpln_workload.new.self_link
		port = 443
		host_prefix   = "my.thing."
	}

	resource "cpln_domain_route" "route_3" {

		depends_on = [cpln_domain_route.route_second]
		
		domain_link = cpln_domain.subdomain.self_link
		domain_port = 80

		prefix = "/3"
		workload_link = cpln_workload.new.self_link
		port = 443
		host_prefix   = "my.thing."
	}
	
	`, random, apex, description, domain, dnsMode)
}

func testAccControlPlaneDomainPathBasedUpdateRoutePort(random, apex, description, domain, dnsMode string) string {

	TestLogger.Printf("Inside testAccControlPlaneDomain")

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "domain_gvc" {

		name        = "gvc-${var.random-name}"
		description = "Example GVC"

		locations = ["aws-eu-central-1", "aws-us-west-2"]

		tags = {
		  terraform_generated = "true"
		}
	}

	resource "cpln_workload" "new" {

		gvc = cpln_gvc.domain_gvc.name

		name        = "workload-${var.random-name}"
		description = "Example Workload"
		type        = "serverless"

		tags = {
		  terraform_generated = "true"
		}

		container {
		  name   = "container-01"
		  image  = "gcr.io/knative-samples/helloworld-go"
		  port   = 8080
		  memory = "128Mi"
		  cpu    = "50m"
		}

		options {
		  capacity_ai     = false
		  timeout_seconds = 30
		  suspend         = true

		  autoscaling {
			metric          = "concurrency"
			target          = 100
			max_scale       = 0
			min_scale       = 0
			max_concurrency = 500
		  }
		}
	}

	resource "cpln_domain" "domain_apex" {
		name        = "%s"
		description = "%s"
	  
		tags = {
		  terraform_generated = "true"
		}

		spec {
			ports {
				tls { }
			 }
		}	
	}

	resource "cpln_domain" "subdomain" {
		
		depends_on = [cpln_domain.domain_apex]

		name        = "%s"
		description = "NS - Path Based"
	  
		tags = {
		  terraform_generated = "true"
		}
	  
		spec {

			dns_mode = "%s"
			
			ports {
				number   = 443
				protocol = "http"

				cors {
					allow_origins {
						exact = "example.com"
					}

					allow_origins {
						exact = "*"
					}
			
					allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
					allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
					expose_headers     = ["expose_header_1", "expose_header_2", "expose_header_3"]
					max_age           = "24h"
					allow_credentials = "true"
				}
	  
				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
						"ECDHE-ECDSA-CHACHA20-POLY1305",
						"ECDHE-ECDSA-AES128-GCM-SHA256",
						"ECDHE-RSA-AES256-GCM-SHA384",
						"ECDHE-RSA-CHACHA20-POLY1305",
						"ECDHE-RSA-AES128-GCM-SHA256",
						"AES256-GCM-SHA384",
						"AES128-GCM-SHA256",
					]
				}
			}

			ports {
				number   = 80
				protocol = "http"

				cors {
					allow_origins {
						exact = "example.com"
					}

					allow_origins {
						exact = "*"
					}
			
					allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
					allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
					expose_headers     = ["expose_header_1", "expose_header_2", "expose_header_3"]
					max_age           = "24h"
					allow_credentials = "true"
				}
	  
				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
						"ECDHE-ECDSA-CHACHA20-POLY1305",
						"ECDHE-ECDSA-AES128-GCM-SHA256",
						"ECDHE-RSA-AES256-GCM-SHA384",
						"ECDHE-RSA-CHACHA20-POLY1305",
						"ECDHE-RSA-AES128-GCM-SHA256",
						"AES256-GCM-SHA384",
						"AES128-GCM-SHA256",
					]
				}
			}
		}
	}
	
	resource "cpln_domain_route" "route_first" {

		domain_link = cpln_domain.subdomain.self_link
		// domain_port = 443

		prefix = "/first"
		workload_link = cpln_workload.new.self_link
		port = 80
		host_prefix   = "my.thing.update."
	}

	resource "cpln_domain_route" "route_second" {

		depends_on = [cpln_domain_route.route_first]
		
		domain_link = cpln_domain.subdomain.self_link
		// domain_port = 443

		prefix = "/second"
		replace_prefix = "/"
		workload_link = cpln_workload.new.self_link
		port = 443
		host_prefix   = "my.thing.update."
	}

	resource "cpln_domain_route" "route_3" {

		depends_on = [cpln_domain_route.route_second]
		
		domain_link = cpln_domain.subdomain.self_link
		domain_port = 80

		prefix = "/3"
		workload_link = cpln_workload.new.self_link
		port = 443
		host_prefix   = "my.thing.update."
	}
	
	`, random, apex, description, domain, dnsMode)
}

func testAccCheckControlPlaneDomainExists(resourceName string, domainName string, domain *client.Domain, orgName *client.Org) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadExists. Resources Length: %d", len(s.RootModule().Resources))

		resources := s.RootModule().Resources
		rs, ok := resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != domainName {
			return fmt.Errorf("Domain name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)

		d, _, err := client.GetDomain(domainName)

		if err != nil {
			return err
		}

		if *d.Name != domainName {
			return fmt.Errorf("Domain name does not match")
		}

		*domain = *d

		o, _, err := client.GetOrg()

		if err != nil {
			return err
		}

		*orgName = *o

		return nil
	}
}

func testAccCheckControlPlaneDomainNSSubdomain(domain *client.Domain, org *client.Org, gvc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadNsSubdomain. Resources Length: %d", len(s.RootModule().Resources))

		expectedDnsMode := "ns"
		dnsMode := domain.Spec.DnsMode

		if *dnsMode != expectedDnsMode {
			return fmt.Errorf("DnsMode does not match, value: %v, expected: %v", *dnsMode, expectedDnsMode)
		}

		gvcLink := domain.Spec.GvcLink
		gvcName := "/org/" + *org.Name + "/gvc/" + gvc

		if *gvcLink != gvcName {
			return fmt.Errorf("GvcLink does not match, value %v, expected: %v", *gvcLink, gvcName)
		}

		return nil
	}
}

func testAccCheckControlPlaneDomainNSPathBased(domain *client.Domain, org *client.Org, randomName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadNsPathBased. Resources Length: %d", len(s.RootModule().Resources))

		expectedDnsMode := "ns"
		dnsMode := domain.Spec.DnsMode

		if *dnsMode != expectedDnsMode {
			return fmt.Errorf("DnsMode does not match, value: %v, expected: %v", dnsMode, expectedDnsMode)
		}

		port1 := 80
		prefix1 := "/first"
		hostPrefix := "my.thing." // On update this will fail
		wl := "/org/" + *org.Name + "/gvc/gvc-" + randomName + "/workload/workload-" + randomName

		routes := []client.DomainRoute{
			{
				Prefix:       &prefix1,
				WorkloadLink: &wl,
				Port:         &port1,
				HostPrefix:   &hostPrefix,
			},
		}

		if diff := deep.Equal(&routes, (*domain.Spec.Ports)[0].Routes); diff != nil {
			return fmt.Errorf("Domain spec port routes does not match. Diff: %s", diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneDomainCNameSubdomain(domain *client.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadCNameSubdomain. Resources Length: %d", len(s.RootModule().Resources))

		expectedDnsMode := "cname"
		dnsMode := domain.Spec.DnsMode

		if *dnsMode != expectedDnsMode {
			return fmt.Errorf("DnsMode does not match, value: %v, expected: %v", *dnsMode, expectedDnsMode)
		}

		expectedGVCLink := getGVC()
		gvcLink := domain.Spec.GvcLink
		if *gvcLink != expectedGVCLink {
			return fmt.Errorf("GVCLink does not match, value: %v, expected: %v", *gvcLink, expectedGVCLink)
		}

		return nil
	}
}

func testAccCheckControlPlaneDomainCNamePathBased(domain *client.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadCNamePathBased. Resources Length: %d", len(s.RootModule().Resources))

		expectedDnsMode := "cname"
		dnsMode := domain.Spec.DnsMode

		if *dnsMode != expectedDnsMode {
			return fmt.Errorf("DnsMode does not match, value: %v, expected: %v", *dnsMode, expectedDnsMode)
		}

		prefix1 := "/first"
		prefix2 := "/second"
		wl1 := getWorkloadOne()
		wl2 := getWorkloadTwo()
		port1 := 8080
		port2 := 8081
		hostPrefix1 := "my.thing." // On update this will fail
		hostPrefix2 := "my."       // On update this will fail
		routes := []client.DomainRoute{
			{
				Prefix:       &prefix1,
				WorkloadLink: &wl1,
				Port:         &port1,
				HostPrefix:   &hostPrefix1,
			},
			{
				Prefix:       &prefix2,
				WorkloadLink: &wl2,
				Port:         &port2,
				HostPrefix:   &hostPrefix2,
			},
		}

		if diff := deep.Equal(&routes, (*domain.Spec.Ports)[0].Routes); diff != nil {
			return fmt.Errorf("Domain spec port routes does not match. Diff: %s", diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneDomainCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneDomainCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneDomainCheckDestroy: rs.Type: %s", rs.Type)

		if rs.Type != "cpln_domain" {
			continue
		}

		domainName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneDomainCheckDestroy: domainName: %s", domainName)

		domain, _, _ := c.GetDomain(domainName)
		if domain != nil {
			return fmt.Errorf("Domain still exists. Name: %s.", *domain.Name)
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build Domain Spec //
func TestControlPlane_BuildDomainSpec(t *testing.T) {
	dnsMode := "ns"
	gvcLink := getGVC()
	acceptAllHosts := true
	_, expectedPorts, flattenedPorts := generatePorts()

	domainSpec := buildDomainSpec(generateFlatTestDomainSpec(dnsMode, gvcLink, acceptAllHosts, flattenedPorts))
	expectedDomainSpec := client.DomainSpec{
		DnsMode:        &dnsMode,
		GvcLink:        &gvcLink,
		AcceptAllHosts: &acceptAllHosts,
		Ports:          &expectedPorts,
	}

	if diff := deep.Equal(domainSpec, &expectedDomainSpec); diff != nil {
		t.Errorf("DomainSpec was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildDomainSpec_NoPorts(t *testing.T) {
	dnsMode := "ns"
	gvcLink := getGVC()
	acceptAllHosts := false

	domainSpec := buildDomainSpec(generateFlatTestDomainSpec(dnsMode, gvcLink, acceptAllHosts, nil))
	expectedDomainSpec := client.DomainSpec{
		DnsMode:        &dnsMode,
		GvcLink:        &gvcLink,
		AcceptAllHosts: &acceptAllHosts,
		Ports:          nil,
	}

	if diff := deep.Equal(domainSpec, &expectedDomainSpec); diff != nil {
		t.Errorf("DomainSpec was not built correctly. Diff: %s", diff)
	}
}

// Build Ports //
func TestControlPlane_BuildPorts(t *testing.T) {
	ports, expectedPorts, _ := generatePorts()
	if diff := deep.Equal(ports, &expectedPorts); diff != nil {
		t.Errorf("Ports were not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildPorts_Empty(t *testing.T) {
	ports := buildSpecPorts(generateEmptyInterfaceArray())
	expectedPorts := []client.DomainSpecPort{{}}

	if diff := deep.Equal(ports, &expectedPorts); diff != nil {
		t.Errorf("Ports were not built correctly. Diff: %s", diff)
	}
}

// Build Cors //
func TestControlPlane_BuildCors(t *testing.T) {

	allowMethods := []string{"2", "3", "1"}
	allowHeaders := []string{"2"}
	exposeHeaders := []string{"3"}
	maxAge := "24h"
	allowCredentials := true

	stringFunc := schema.HashSchema(StringSchema())

	_, expectedAllowOrigins, flattenedAllowOrigins := generateAllowOrigins()
	flattened := generateFlatTestCors(flattenedAllowOrigins,
		schema.NewSet(stringFunc, flattenStringsArray(&allowMethods)),
		schema.NewSet(stringFunc, flattenStringsArray(&allowHeaders)),
		schema.NewSet(stringFunc, flattenStringsArray(&exposeHeaders)), maxAge, allowCredentials)

	cors := buildCors(flattened)
	expectedCors := client.DomainCors{
		AllowOrigins:     &expectedAllowOrigins,
		AllowMethods:     &allowMethods,
		AllowHeaders:     &allowHeaders,
		ExposeHeaders:    &exposeHeaders,
		MaxAge:           &maxAge,
		AllowCredentials: &allowCredentials,
	}

	if diff := deep.Equal(cors, &expectedCors); diff != nil {
		t.Errorf("Cors was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildCors_Empty(t *testing.T) {
	cors := buildCors(generateEmptyInterfaceArray())
	expectedCors := client.DomainCors{}

	if diff := deep.Equal(cors, &expectedCors); diff != nil {
		t.Errorf("Cors was not built correctly. Diff: %s", diff)
	}
}

// Build TLS Unit Test //
func TestControlPlane_BuildTLS(t *testing.T) {
	tls, expectedTLS, _ := generateTLS()
	if diff := deep.Equal(tls, &expectedTLS); diff != nil {
		t.Errorf("TLS was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildTLS_Empty(t *testing.T) {
	tls := buildTLS(generateEmptyInterfaceArray())
	expectedTLS := client.DomainTLS{}

	if diff := deep.Equal(tls, &expectedTLS); diff != nil {
		t.Errorf("TLS was not built correctly. Diff: %s", diff)
	}
}

// Build Allow Origins Unit Test //
func TestControlPlane_BuildAllowOrigins(t *testing.T) {
	collection, expectedCollection, _ := generateAllowOrigins()
	if diff := deep.Equal(collection, &expectedCollection); diff != nil {
		t.Errorf("Allow Origins was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildAllowOrigins_WithoutExact(t *testing.T) {
	collection := buildAllowOrigins(generateEmptyInterfaceArray())
	expectedCollection := []client.DomainAllowOrigin{{}}

	if diff := deep.Equal(collection, &expectedCollection); diff != nil {
		t.Errorf("Allow Origins was not built correctly. Diff: %s", diff)
	}
}

// Build Certificate Unit Test //
func TestControlPlane_BuildCertificate(t *testing.T) {
	secret := "/org/myorg/secret/mysecret"

	cert := buildCertificate(generateFlatTestCertificate(secret))
	expectedCert := client.DomainCertificate{SecretLink: &secret}

	// TODO move expectedCert to a function, can be array of items too for different cases
	if diff := deep.Equal(cert, &expectedCert); diff != nil {
		t.Errorf("Domain Certificate was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildCertificate_WithoutSecret(t *testing.T) {
	cert := buildCertificate(generateEmptyInterfaceArray())
	certTest := client.DomainCertificate{}

	if diff := deep.Equal(cert, &certTest); diff != nil {
		t.Errorf("Domain Certificate was not built correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
func generatePorts() (*[]client.DomainSpecPort, []client.DomainSpecPort, []interface{}) {
	number := 443
	protocol := "http"

	_, expectedCors, flattenedCors := generateCors()
	_, expectedTLS, flattenedTLS := generateTLS()

	flattenGeneration := generateInterfaceArrayFromMapArray([]map[string]interface{}{
		generateFlatTestPort(number, protocol, flattenedCors, flattenedTLS),
	})

	ports := buildSpecPorts(flattenGeneration)
	expectedPorts := []client.DomainSpecPort{
		{
			Number:   &number,
			Protocol: &protocol,
			Cors:     &expectedCors,
			TLS:      &expectedTLS,
		},
	}

	return ports, expectedPorts, flattenGeneration
}

func generateCors() (*client.DomainCors, client.DomainCors, []interface{}) {
	allowMethods := []string{"1"}
	allowHeaders := []string{"2"}
	exposeHeaders := []string{"3"}
	maxAge := "24h"
	allowCredentials := true

	stringFunc := schema.HashSchema(StringSchema())

	_, expectedAllowOrigins, flattenedAllowOrigins := generateAllowOrigins()
	flattened := generateFlatTestCors(flattenedAllowOrigins,
		schema.NewSet(stringFunc, flattenStringsArray(&allowMethods)),
		schema.NewSet(stringFunc, flattenStringsArray(&allowHeaders)),
		schema.NewSet(stringFunc, flattenStringsArray(&exposeHeaders)), maxAge, allowCredentials)

	cors := buildCors(flattened)
	expectedCors := client.DomainCors{
		AllowOrigins:     &expectedAllowOrigins,
		AllowMethods:     &allowMethods,
		AllowHeaders:     &allowHeaders,
		ExposeHeaders:    &exposeHeaders,
		MaxAge:           &maxAge,
		AllowCredentials: &allowCredentials,
	}

	return cors, expectedCors, flattened
}

func generateTLS() (*client.DomainTLS, client.DomainTLS, []interface{}) {
	minProtocolVersion := "TLSv1_1"
	cipherSuites := []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA"}
	clientSecret := "/org/myorg/secret/mysecret_client"
	serverSecret := "/org/myorg/secret/mysecret_server"

	stringFunc := schema.HashSchema(StringSchema())
	cipherSuitesFlattened := schema.NewSet(stringFunc, flattenStringsArray(&cipherSuites))
	clientCertificate := generateFlatTestCertificate(clientSecret)
	serverCertificate := generateFlatTestCertificate(serverSecret)

	flattened := generateFlatTestTLS(minProtocolVersion, cipherSuitesFlattened, clientCertificate, serverCertificate)

	tls := buildTLS(flattened)
	expectedTLS := client.DomainTLS{
		MinProtocolVersion: &minProtocolVersion,
		CipherSuites:       &cipherSuites,
		ClientCertificate:  &client.DomainCertificate{SecretLink: &clientSecret},
		ServerCertificate:  &client.DomainCertificate{SecretLink: &serverSecret},
	}

	return tls, expectedTLS, flattened
}

func generateAllowOrigins() (*[]client.DomainAllowOrigin, []client.DomainAllowOrigin, []interface{}) {
	exact := "example.com"
	flattened := generateFlatTestAllowOrigins(exact)

	collection := buildAllowOrigins(flattened)
	expectedCollection := []client.DomainAllowOrigin{{Exact: &exact}}

	return collection, expectedCollection, flattened
}

/*** Flatten ***/
func generateFlatTestDomainSpec(dnsMode string, gvcLink string, acceptAllHosts bool, ports []interface{}) []interface{} {
	spec := map[string]interface{}{
		"dns_mode":         dnsMode,
		"gvc_link":         gvcLink,
		"accept_all_hosts": acceptAllHosts,
		"ports":            ports,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestPort(number int, protocol string, cors []interface{}, tls []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"number":   number,
		"protocol": protocol,
		"cors":     cors,
		"tls":      tls,
	}
}

func generateFlatTestRoute(prefix string, replacePrefix string, workloadLink string, port int) map[string]interface{} {
	return map[string]interface{}{
		"prefix":         prefix,
		"replace_prefix": replacePrefix,
		"workload_link":  workloadLink,
		"port":           port,
	}
}

func generateFlatTestCors(allowOrigins []interface{}, allowMethods interface{}, allowHeaders interface{}, exposeHeaders interface{}, maxAge string, allowCredentials bool) []interface{} {

	spec := map[string]interface{}{
		"allow_origins":     allowOrigins,
		"allow_methods":     allowMethods,
		"allow_headers":     allowHeaders,
		"expose_headers":    exposeHeaders,
		"max_age":           maxAge,
		"allow_credentials": allowCredentials,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestTLS(minProtocolVersion string, cipherSuites interface{}, clientCertificate []interface{}, serverCertificate []interface{}) []interface{} {
	spec := map[string]interface{}{
		"min_protocol_version": minProtocolVersion,
		"cipher_suites":        cipherSuites,
		"client_certificate":   clientCertificate,
		"server_certificate":   serverCertificate,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestAllowOrigins(exact string) []interface{} {
	spec := map[string]interface{}{
		"exact": exact,
	}

	return []interface{}{
		spec,
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

func generateInterfaceArrayFromMapArray(specs []map[string]interface{}) []interface{} {
	collection := make([]interface{}, len(specs))
	for i, spec := range specs {
		collection[i] = spec
	}

	return collection
}

func generateEmptyInterfaceArray() []interface{} {
	return []interface{}{
		map[string]interface{}{},
	}
}
