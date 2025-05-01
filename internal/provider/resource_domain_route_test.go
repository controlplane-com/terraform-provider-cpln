package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneDomainRoute_basic(t *testing.T) {

	domainName := "erickotler.com"
	orgName := "terraform-test-org"

	var domain client.Domain
	var org client.Org

	// Formatted
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	subDomainName := fmt.Sprintf("domain-route-%s.%s", random, domainName)
	gvcName := fmt.Sprintf("gvc-%s", random)
	workloadName := fmt.Sprintf("workload-%s", random)

	// Links
	subDomainLink := fmt.Sprintf("/org/%s/domain/%s", orgName, subDomainName)
	workloadLink := fmt.Sprintf("/org/%s/gvc/%s/workload/%s", orgName, gvcName, workloadName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, "DOMAIN_ROUTE") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneDomainRoute_Prefix(random, domainName, subDomainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.apex", domainName, &domain, &org),
					testAccCheckControlPlaneDomainExists("cpln_domain.subdomain", subDomainName, &domain, &org),
					testAccCheckControlPlaneDomainRouteExists("cpln_domain_route.first-route", fmt.Sprintf("%s_443_/first", subDomainLink)),
					testAccCheckControlPlaneDomainRouteExists("cpln_domain_route.second-route", fmt.Sprintf("%s_80_/second", subDomainLink)),
					testAccCheckControlPlaneDomainRouteExists("cpln_domain_route.third-route", fmt.Sprintf("%s_80_/third", subDomainLink)),
					testAccCheckControlPlaneDomainRouteAttributes("domain-with-prefix", &domain, workloadLink),
					resource.TestCheckResourceAttr("cpln_domain.subdomain", "name", subDomainName),
					resource.TestCheckResourceAttr("cpln_domain.subdomain", "description", "NS - Path Based"),
				),
			},
			{
				Config: testAccControlPlaneDomainRoute_Regex(random, domainName, subDomainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists("cpln_domain.apex", domainName, &domain, &org),
					testAccCheckControlPlaneDomainExists("cpln_domain.subdomain", subDomainName, &domain, &org),
					testAccCheckControlPlaneDomainRouteExists("cpln_domain_route.first-route", fmt.Sprintf("%s_443_/user/.*/profile", subDomainLink)),
					testAccCheckControlPlaneDomainRouteAttributes("domain-with-regex", &domain, workloadLink),
					resource.TestCheckResourceAttr("cpln_domain.subdomain", "name", subDomainName),
					resource.TestCheckResourceAttr("cpln_domain.subdomain", "description", "NS - Path Based"),
				),
			},
		},
	})
}

func testAccControlPlaneDomainRoute_Prefix(random string, domainName string, subDomainName string) string {
	return fmt.Sprintf(`

	variable "random-name" {
		type    = string
		default = "%s"
	}

	resource "cpln_gvc" "new" {

		name        = "gvc-${var.random-name}"
		description = "Example GVC"

		locations = ["aws-eu-central-1", "aws-us-west-2"]

		tags = {
			terraform_generated = "true"
		}
	}

	resource "cpln_workload" "new" {

		gvc = cpln_gvc.new.name

		name        = "workload-${var.random-name}"
		description = "Example Workload"
		type        = "serverless"

		tags = {
			terraform_generated = "true"
		}

		container {
			name   = "container-01"
			image  = "gcr.io/knative-samples/helloworld-go"
			cpu    = "50m"
			memory = "128Mi"
			port   = 8080
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

	resource "cpln_domain" "apex" {
		name        = "%s"
		description = "Apex Domain Description"

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
		
		depends_on = [cpln_domain.apex]

		name        = "%s"
		description = "NS - Path Based"

		tags = {
		  terraform_generated = "true"
		}

		spec {

			dns_mode = "ns"

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

					allow_methods     = ["allow_method"]
					allow_headers     = ["allow_header"]
					expose_headers    = ["expose_header"]
					max_age           = "24h"
					allow_credentials = "true"
				}

				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
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

					allow_methods     = ["allow_method"]
					allow_headers     = ["allow_header"]
					expose_headers    = ["expose_header"]
					max_age           = "24h"
					allow_credentials = "true"
				}

				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
					]
				}
			}
		}
	}

	resource "cpln_domain_route" "first-route" {

		domain_link = cpln_domain.subdomain.self_link

		prefix 		  = "/first"
		workload_link = cpln_workload.new.self_link
	}

	resource "cpln_domain_route" "second-route" {

		depends_on = [cpln_domain_route.first-route]

		domain_link = cpln_domain.subdomain.self_link
		domain_port = 80

		prefix 		     = "/second"
		replace_prefix = "/"
		workload_link  = cpln_workload.new.self_link
		port 		       = 443
		host_prefix    = "my.thing."

		headers {
			request {
				set = {
					Host = "example.com"
					"Content-Type" = "application/json"
				}
			}
		}
	}

	resource "cpln_domain_route" "third-route" {

		depends_on = [cpln_domain_route.second-route]

		domain_link = cpln_domain.subdomain.self_link
		domain_port = 80

		prefix 		     = "/third"
		replace_prefix = "/"
		workload_link  = cpln_workload.new.self_link
		port 		       = 443
		host_regex     = "reg"

		headers {
			request {
				set = {
					Host = "example.com"
					"Content-Type" = "application/json"
				}
			}
		}
	}
	`, random, domainName, subDomainName)
}

func testAccControlPlaneDomainRoute_Regex(random string, domainName string, subDomainName string) string {
	return fmt.Sprintf(`

	variable "random-name" {
		type    = string
		default = "%s"
	}

	resource "cpln_gvc" "new" {

		name        = "gvc-${var.random-name}"
		description = "Example GVC"

		locations = ["aws-eu-central-1", "aws-us-west-2"]

		tags = {
			terraform_generated = "true"
		}
	}

	resource "cpln_workload" "new" {

		gvc = cpln_gvc.new.name

		name        = "workload-${var.random-name}"
		description = "Example Workload"
		type        = "serverless"

		tags = {
			terraform_generated = "true"
		}

		container {
			name   = "container-01"
			image  = "gcr.io/knative-samples/helloworld-go"
			cpu    = "50m"
			memory = "128Mi"
			port   = 8080
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

	resource "cpln_domain" "apex" {
		name        = "%s"
		description = "Apex Domain Description"

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
		
		depends_on = [cpln_domain.apex]

		name        = "%s"
		description = "NS - Path Based"

		tags = {
		  terraform_generated = "true"
		}

		spec {

			dns_mode = "ns"

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

					allow_methods     = ["allow_method"]
					allow_headers     = ["allow_header"]
					expose_headers    = ["expose_header"]
					max_age           = "24h"
					allow_credentials = "true"
				}

				tls {
					min_protocol_version = "TLSV1_2"
					cipher_suites = [
						"ECDHE-ECDSA-AES256-GCM-SHA384",
					]
				}
			}
		}
	}

	resource "cpln_domain_route" "first-route" {

		domain_link = cpln_domain.subdomain.self_link

		regex 		  = "/user/.*/profile"
		workload_link = cpln_workload.new.self_link
		port 		  = 8080
	}
	`, random, domainName, subDomainName)
}

func testAccCheckControlPlaneDomainRouteExists(resourceName string, routeIdentifier string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneDomainRouteExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != routeIdentifier {
			return fmt.Errorf("Domain Route ID '%s' does not match '%s'", rs.Primary.ID, routeIdentifier)
		}

		return nil
	}
}

func testAccCheckControlPlaneDomainRouteAttributes(state string, domain *client.Domain, workloadLink string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if domain.Spec == nil {
			return fmt.Errorf("Domain Spec is nil")
		}

		if domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return fmt.Errorf("Domain Ports are nil or empty")
		}

		var expectedDomainPorts *[]client.DomainSpecPort

		switch state {
		case "domain-with-prefix":
			expectedDomainPorts = generateTestDomainRoutePorts_Prefix(workloadLink)
		case "domain-with-regex":
			expectedDomainPorts = generateTestDomainRoutePorts_Regex(workloadLink)
		}

		if diff := deep.Equal(domain.Spec.Ports, expectedDomainPorts); diff != nil {
			return fmt.Errorf("Domain Ports do not match. Diff: %s", diff)
		}

		return nil
	}
}

/*** Unit Tests ***/

// Build //

func TestControlPlane_BuildDomainRouteHeaders(t *testing.T) {

	headers, expectedHeaders, _ := generateTestDomainRouteHeaders()

	if diff := deep.Equal(headers, expectedHeaders); diff != nil {
		t.Errorf("Domain Route Headers were not built correctly, Diff: %s", diff)
	}
}

// Flatten //

func TestControlPlane_FlattenDomainRouteHeaders(t *testing.T) {

	_, expectedHeaders, expectedFlatten := generateTestDomainRouteHeaders()
	flattenedHeaders := flattenDomainRouteHeaders(expectedHeaders)

	if diff := deep.Equal(expectedFlatten, flattenedHeaders); diff != nil {
		t.Errorf("Domain Route Headers were not flattened correctly. Diff: %s", diff)
	}
}

/*** Generate ***/

// Build //

func generateTestDomainRoutePorts_Prefix(workloadLink string) *[]client.DomainSpecPort {
	headers, _, _ := generateTestDomainRouteHeaders()

	return &[]client.DomainSpecPort{
		{
			Number:   GetInt(443),
			Protocol: GetString("http"),
			Routes: &[]client.DomainRoute{
				{
					Prefix:       GetString("/first"),
					WorkloadLink: GetString(workloadLink),
				},
			},
			Cors: &client.DomainCors{
				AllowOrigins: &[]client.DomainAllowOrigin{
					{
						Exact: GetString("example.com"),
					},
					{
						Exact: GetString("*"),
					},
				},
				AllowMethods:     &[]string{"allow_method"},
				AllowHeaders:     &[]string{"allow_header"},
				ExposeHeaders:    &[]string{"expose_header"},
				MaxAge:           GetString("24h"),
				AllowCredentials: GetBool(true),
			},
			TLS: &client.DomainTLS{
				MinProtocolVersion: GetString("TLSV1_2"),
				CipherSuites: &[]string{
					"ECDHE-ECDSA-AES256-GCM-SHA384",
				},
			},
		},
		{
			Number:   GetInt(80),
			Protocol: GetString("http"),
			Routes: &[]client.DomainRoute{
				{
					Prefix:        GetString("/second"),
					ReplacePrefix: GetString("/"),
					WorkloadLink:  GetString(workloadLink),
					Port:          GetInt(443),
					HostPrefix:    GetString("my.thing."),
					Headers:       headers,
				},
				{
					Prefix:        GetString("/third"),
					ReplacePrefix: GetString("/"),
					WorkloadLink:  GetString(workloadLink),
					Port:          GetInt(443),
					HostRegex:     GetString("reg"),
					Headers:       headers,
				},
			},
			Cors: &client.DomainCors{
				AllowOrigins: &[]client.DomainAllowOrigin{
					{
						Exact: GetString("example.com"),
					},
					{
						Exact: GetString("*"),
					},
				},
				AllowMethods:     &[]string{"allow_method"},
				AllowHeaders:     &[]string{"allow_header"},
				ExposeHeaders:    &[]string{"expose_header"},
				MaxAge:           GetString("24h"),
				AllowCredentials: GetBool(true),
			},
			TLS: &client.DomainTLS{
				MinProtocolVersion: GetString("TLSV1_2"),
				CipherSuites: &[]string{
					"ECDHE-ECDSA-AES256-GCM-SHA384",
				},
			},
		},
	}
}

func generateTestDomainRoutePorts_Regex(workloadLink string) *[]client.DomainSpecPort {

	return &[]client.DomainSpecPort{
		{
			Number:   GetInt(443),
			Protocol: GetString("http"),
			Routes: &[]client.DomainRoute{
				{
					Regex:        GetString("/user/.*/profile"),
					WorkloadLink: GetString(workloadLink),
					Port:         GetInt(8080),
				},
			},
			Cors: &client.DomainCors{
				AllowOrigins: &[]client.DomainAllowOrigin{
					{
						Exact: GetString("example.com"),
					},
					{
						Exact: GetString("*"),
					},
				},
				AllowMethods:     &[]string{"allow_method"},
				AllowHeaders:     &[]string{"allow_header"},
				ExposeHeaders:    &[]string{"expose_header"},
				MaxAge:           GetString("24h"),
				AllowCredentials: GetBool(true),
			},
			TLS: &client.DomainTLS{
				MinProtocolVersion: GetString("TLSV1_2"),
				CipherSuites: &[]string{
					"ECDHE-ECDSA-AES256-GCM-SHA384",
				},
			},
		},
	}
}

func generateTestDomainRouteHeaders() (*client.DomainRouteHeaders, *client.DomainRouteHeaders, []interface{}) {

	request, _, flattenedRequest := generateTestDomainHeaderOperation()

	flattened := generateFlatTestDomainRouteHeaders(flattenedRequest)
	headers := buildDomainRouteHeaders(flattened)
	expectedHeaders := client.DomainRouteHeaders{
		Request: request,
	}

	return headers, &expectedHeaders, flattened
}

func generateTestDomainHeaderOperation() (*client.DomainHeaderOperation, *client.DomainHeaderOperation, []interface{}) {

	set := map[string]interface{}{
		"Host":         "example.com",
		"Content-Type": "application/json",
	}

	flattened := generateFlatTestDomainHeaderOperation(set)
	request := buildDomainHeaderOperation(flattened)
	expectedRequest := client.DomainHeaderOperation{
		Set: &set,
	}

	return request, &expectedRequest, flattened
}

// Flatten //

func generateFlatTestDomainRouteHeaders(request []interface{}) []interface{} {

	spec := map[string]interface{}{
		"request":               request,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestDomainHeaderOperation(set map[string]interface{}) []interface{} {

	spec := map[string]interface{}{
		"set":                   set,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}
