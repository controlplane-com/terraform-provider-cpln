package cpln

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDomain_basic performs an acceptance test for the resource.
func TestAccControlPlaneDomain_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewDomainResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DOMAIN") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// DomainResourceTest defines the necessary functionality to test the resource.
type DomainResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
	ApexDomain string
}

// DomainInlineRouteConfig represents an inline route in HCL config.
type DomainInlineRouteConfig struct {
	Prefix       string
	Regex        string
	WorkloadLink string // HCL expression (e.g., "cpln_workload.new.self_link")
	Port         int
}

// NewDomainResourceTest creates a DomainResourceTest with initialized test cases.
func NewDomainResourceTest() DomainResourceTest {
	// Create a resource test instance
	resourceTest := DomainResourceTest{
		RandomName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		ApexDomain: "erickotler.com",
	}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewDefaultScenario()...)
	steps = append(steps, resourceTest.NewCoexistenceScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (drt *DomainResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_domain resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_domain" {
			continue
		}

		// Retrieve the name for the current resource
		domainName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of domain with name: %s", domainName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		domain, code, err := TestProvider.client.GetDomain(domainName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if domain %s exists: %w", domainName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if domain != nil {
			return fmt.Errorf("CheckDestroy failed: domain %s still exists in the system", *domain.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_domain resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case with initial and updated configurations.
func (drt *DomainResourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "new"
	name := drt.ApexDomain
	subDomainSelfLink := GetSelfLink(OrgName, "domain", fmt.Sprintf("domain-acctest-%s.%s", drt.RandomName, name))

	// Build test steps
	initialConfig, initialStep := drt.BuildDefaultTestStep(resourceName, name)
	caseUpdate1 := drt.BuildUpdate1TestStep(initialConfig.ProviderTestCase)
	caseUpdate2 := drt.BuildUpdate2TestStep(initialConfig.ProviderTestCase)
	caseUpdate3 := drt.BuildUpdate3TestStep(initialConfig.ProviderTestCase)
	caseUpdate4 := drt.BuildUpdate4TestStep(initialConfig.ProviderTestCase)
	caseUpdate5 := drt.BuildUpdate5TestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName: initialConfig.ResourceAddress,
			ImportState:  true,
		},
		// Update & Read
		caseUpdate1,
		caseUpdate2,
		caseUpdate3,
		// Domain Route Import
		{
			ResourceName:  "cpln_domain_route.first-route",
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:443:/first", subDomainSelfLink),
		},
		{
			ResourceName:  "cpln_domain_route.second-route",
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:80:/second", subDomainSelfLink),
		},
		{
			ResourceName:  "cpln_domain_route.third-route",
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:80:/third", subDomainSelfLink),
		},
		{
			ResourceName:  "cpln_domain_route.fourth-route",
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:443:/user/.*/profile", subDomainSelfLink),
		},
		// Inline Routes: Create & Read
		caseUpdate4,
		// Inline Routes: Update & Read
		caseUpdate5,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewCoexistenceScenario creates a test scenario covering all unique route coexistence transitions.
func (drt *DomainResourceTest) NewCoexistenceScenario() []resource.TestStep {
	// Define necessary variables
	subDomainName := fmt.Sprintf("route-coexist-%s.%s", drt.RandomName, drt.ApexDomain)

	// Build test steps
	return []resource.TestStep{
		// Create domain + inline routes + external routes in a single apply
		drt.BuildCoexistenceCreateTestStep(subDomainName),
		// Strip all routes
		drt.BuildNoRoutesTestStep(subDomainName, "no routes"),
		// Add external routes only
		drt.BuildExternalRoutesOnlyTestStep(subDomainName, "external routes only"),
		// Add inline routes alongside external (coexistence)
		drt.BuildCoexistenceTestStep(subDomainName),
		// Remove external, keep inline only
		drt.BuildInlineOnlyTestStep(subDomainName),
		// Import State
		{
			ResourceName:  "cpln_domain.subdomain",
			ImportState:   true,
			ImportStateId: subDomainName,
		},
		// Cleanup: remove all routes
		drt.BuildNoRoutesTestStep(subDomainName, "final cleanup"),
	}
}

// Test Cases //

// BuildDefaultTestStep returns a default initial test step and its associated test case for the resource.
func (drt *DomainResourceTest) BuildDefaultTestStep(resourceName string, name string) (DomainResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "domain",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_domain.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "domain new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: drt.RequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "cname",
					"accept_all_hosts": "false",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http2",
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"AES128-GCM-SHA256",
										"AES256-GCM-SHA384",
										"ECDHE-ECDSA-AES128-GCM-SHA256",
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-ECDSA-CHACHA20-POLY1305",
										"ECDHE-RSA-AES128-GCM-SHA256",
										"ECDHE-RSA-AES256-GCM-SHA384",
										"ECDHE-RSA-CHACHA20-POLY1305",
									},
								},
							},
						},
					},
				},
			}),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (drt *DomainResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := DomainResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: drt.Update1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":              "cname",
					"gvc_link":              "/org/terraform-test-org/gvc/gvc-01",
					"cert_challenge_type":   "dns01",
					"accept_all_hosts":      "false",
					"accept_all_subdomains": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http2",
							"cors": []map[string]interface{}{
								{
									"allow_origins": []map[string]interface{}{
										{
											"exact": "*",
										},
										{
											"exact": "*.erickotler.com",
										},
										{
											"regex": `^https://example\.com$`,
										},
									},
									"allow_methods":     []string{"GET", "OPTIONS", "POST"},
									"allow_headers":     []string{"authorization", "host"},
									"expose_headers":    []string{"accept/type"},
									"max_age":           "12h",
									"allow_credentials": "true",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_1",
									"cipher_suites":        []string{"AES256-GCM-SHA384"},
									"client_certificate": []map[string]interface{}{
										{
											"secret_link": "/org/terraform-test-org/secret/aa-tbd-2",
										},
									},
									"server_certificate": []map[string]interface{}{
										{
											"secret_link": "/org/terraform-test-org/secret/aa-tbd-2",
										},
									},
								},
							},
						},
					},
				},
			}),
		),
	}
}

// BuildUpdate2TestStep returns a test step for the update.
func (drt *DomainResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := DomainResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Create the sub-domain test case
	subDomainName := fmt.Sprintf("domain-acctest-%s.%s", drt.RandomName, initialCase.Name)
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "domain",
			ResourceName:      "subdomain",
			ResourceAddress:   "cpln_domain.subdomain",
			Name:              subDomainName,
			Description:       subDomainName,
			DescriptionUpdate: "domain new description",
		},
	}

	// Create the domain route test cases
	domainRoute1 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "first-route",
			ResourceAddress: "cpln_domain_route.first-route",
		},
	}

	domainRoute2 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "second-route",
			ResourceAddress: "cpln_domain_route.second-route",
		},
	}

	domainRoute3 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "third-route",
			ResourceAddress: "cpln_domain_route.third-route",
		},
	}

	domainRoute4 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "fourth-route",
			ResourceAddress: "cpln_domain_route.fourth-route",
		},
	}

	// Construct the workload self link
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: drt.Update2Hcl(c, subDomain),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Apex Domain
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":              "cname",
					"gvc_link":              "/org/terraform-test-org/gvc/gvc-01",
					"cert_challenge_type":   "http01",
					"accept_all_hosts":      "false",
					"accept_all_subdomains": "false",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http2",
							"cors": []map[string]interface{}{
								{
									"allow_origins": []map[string]interface{}{
										{
											"exact": "*",
										},
										{
											"exact": "*.erickotler.com",
										},
										{
											"regex": `^https://example\.com$`,
										},
									},
									"allow_methods":     []string{"GET", "OPTIONS", "POST"},
									"allow_headers":     []string{"authorization", "host"},
									"expose_headers":    []string{"accept/type"},
									"max_age":           "12h",
									"allow_credentials": "true",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_1",
									"cipher_suites":        []string{"AES256-GCM-SHA384"},
									"client_certificate": []map[string]interface{}{
										{
											"secret_link": "/org/terraform-test-org/secret/aa-tbd-2",
										},
									},
									"server_certificate": []map[string]interface{}{
										{
											"secret_link": "/org/terraform-test-org/secret/aa-tbd-2",
										},
									},
								},
							},
						},
					},
				},
			}),

			// Sub Domain
			subDomain.GetDefaultChecks(subDomain.DescriptionUpdate, "1"),
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"cors": []map[string]interface{}{
								{
									"allow_origins": []map[string]interface{}{
										{
											"exact": "example.com",
										},
										{
											"exact": "*",
										},
									},
									"allow_methods":     []string{"allow_method_1", "allow_method_2", "allow_method_3"},
									"allow_headers":     []string{"allow_header_1", "allow_header_2", "allow_header_3"},
									"expose_headers":    []string{"expose_header_1", "expose_header_2", "expose_header_3"},
									"max_age":           "24h",
									"allow_credentials": "true",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-ECDSA-CHACHA20-POLY1305",
										"ECDHE-ECDSA-AES128-GCM-SHA256",
										"ECDHE-RSA-AES256-GCM-SHA384",
										"ECDHE-RSA-CHACHA20-POLY1305",
										"ECDHE-RSA-AES128-GCM-SHA256",
										"AES256-GCM-SHA384",
										"AES128-GCM-SHA256",
									},
									"client_certificate": []map[string]interface{}{{}},
								},
							},
						},
						{
							"number":   "80",
							"protocol": "http",
							"cors": []map[string]interface{}{
								{
									"allow_origins": []map[string]interface{}{
										{
											"exact": "example.com",
										},
										{
											"exact": "*",
										},
									},
									"allow_methods":     []string{"allow_method"},
									"allow_headers":     []string{"allow_header"},
									"expose_headers":    []string{"expose_header"},
									"max_age":           "24h",
									"allow_credentials": "true",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
									},
								},
							},
						},
					},
				},
			}),

			// First Route
			domainRoute1.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			domainRoute1.TestCheckResourceAttr("domain_port", "443"),
			domainRoute1.TestCheckResourceAttr("prefix", "/first"),
			domainRoute1.TestCheckResourceAttr("replica", "1"),
			domainRoute1.TestCheckResourceAttr("workload_link", workloadSelfLink),

			// Second Route
			domainRoute2.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			domainRoute2.TestCheckResourceAttr("domain_port", "80"),
			domainRoute2.TestCheckResourceAttr("prefix", "/second"),
			domainRoute2.TestCheckResourceAttr("replace_prefix", "/"),
			domainRoute2.TestCheckResourceAttr("workload_link", workloadSelfLink),
			domainRoute2.TestCheckResourceAttr("port", "443"),
			domainRoute2.TestCheckResourceAttr("host_prefix", "my.thing."),
			domainRoute2.TestCheckResourceAttr("replica", "0"),
			domainRoute2.TestCheckNestedBlocks("headers", []map[string]interface{}{
				{
					"request": []map[string]interface{}{
						{
							"set": map[string]interface{}{
								"Host":         "example.com",
								"Content-Type": "application/json",
							},
						},
					},
				},
			}),

			// Third Route
			domainRoute3.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			domainRoute3.TestCheckResourceAttr("domain_port", "80"),
			domainRoute3.TestCheckResourceAttr("prefix", "/third"),
			domainRoute3.TestCheckResourceAttr("replace_prefix", "/"),
			domainRoute3.TestCheckResourceAttr("workload_link", workloadSelfLink),
			domainRoute3.TestCheckResourceAttr("port", "443"),
			domainRoute3.TestCheckResourceAttr("host_regex", "reg"),
			domainRoute3.TestCheckNestedBlocks("headers", []map[string]interface{}{
				{
					"request": []map[string]interface{}{
						{
							"set": map[string]interface{}{
								"Host":         "example.com",
								"Content-Type": "application/json",
							},
						},
					},
				},
			}),

			// Fourth Route
			domainRoute4.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			domainRoute4.TestCheckResourceAttr("domain_port", "443"),
			domainRoute4.TestCheckResourceAttr("regex", "/user/.*/profile"),
			domainRoute4.TestCheckResourceAttr("workload_link", workloadSelfLink),
			domainRoute4.TestCheckResourceAttr("port", "80"),
		),
	}
}

// BuildUpdate3TestStep returns a test step for the update.
func (drt *DomainResourceTest) BuildUpdate3TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := DomainResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Create the sub-domain test case
	subDomainName := fmt.Sprintf("domain-acctest-%s.%s", drt.RandomName, initialCase.Name)
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "domain",
			ResourceName:      "subdomain",
			ResourceAddress:   "cpln_domain.subdomain",
			Name:              subDomainName,
			Description:       subDomainName,
			DescriptionUpdate: "domain new description",
		},
	}

	// Create the domain route test cases
	domainRoute1 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "first-route",
			ResourceAddress: "cpln_domain_route.first-route",
		},
	}

	domainRoute2 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "second-route",
			ResourceAddress: "cpln_domain_route.second-route",
		},
	}

	// Construct the workload self link
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: drt.Update2Hcl(c, subDomain),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Apex Domain
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "cname",
					"gvc_link":         "/org/terraform-test-org/gvc/gvc-01",
					"accept_all_hosts": "false",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http2",
							"cors": []map[string]interface{}{
								{
									"allow_origins": []map[string]interface{}{
										{
											"exact": "*",
										},
										{
											"exact": "*.erickotler.com",
										},
										{
											"regex": `^https://example\.com$`,
										},
									},
									"allow_methods":     []string{"GET", "OPTIONS", "POST"},
									"allow_headers":     []string{"authorization", "host"},
									"expose_headers":    []string{"accept/type"},
									"max_age":           "12h",
									"allow_credentials": "true",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_1",
									"cipher_suites":        []string{"AES256-GCM-SHA384"},
									"client_certificate": []map[string]interface{}{
										{
											"secret_link": "/org/terraform-test-org/secret/aa-tbd-2",
										},
									},
									"server_certificate": []map[string]interface{}{
										{
											"secret_link": "/org/terraform-test-org/secret/aa-tbd-2",
										},
									},
								},
							},
						},
					},
				},
			}),

			// Sub Domain
			subDomain.GetDefaultChecks(subDomain.DescriptionUpdate, "1"),
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"cors": []map[string]interface{}{
								{
									"allow_origins": []map[string]interface{}{
										{
											"exact": "example.com",
										},
										{
											"exact": "*",
										},
									},
									"allow_methods":     []string{"allow_method_1", "allow_method_2", "allow_method_3"},
									"allow_headers":     []string{"allow_header_1", "allow_header_2", "allow_header_3"},
									"expose_headers":    []string{"expose_header_1", "expose_header_2", "expose_header_3"},
									"max_age":           "24h",
									"allow_credentials": "true",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-ECDSA-CHACHA20-POLY1305",
										"ECDHE-ECDSA-AES128-GCM-SHA256",
										"ECDHE-RSA-AES256-GCM-SHA384",
										"ECDHE-RSA-CHACHA20-POLY1305",
										"ECDHE-RSA-AES128-GCM-SHA256",
										"AES256-GCM-SHA384",
										"AES128-GCM-SHA256",
									},
									"client_certificate": []map[string]interface{}{{}},
								},
							},
						},
						{
							"number":   "80",
							"protocol": "http",
							"cors": []map[string]interface{}{
								{
									"allow_origins": []map[string]interface{}{
										{
											"exact": "example.com",
										},
										{
											"exact": "*",
										},
									},
									"allow_methods":     []string{"allow_method"},
									"allow_headers":     []string{"allow_header"},
									"expose_headers":    []string{"expose_header"},
									"max_age":           "24h",
									"allow_credentials": "true",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
									},
								},
							},
						},
					},
				},
			}),

			// First Route
			domainRoute1.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			domainRoute1.TestCheckResourceAttr("domain_port", "443"),
			domainRoute1.TestCheckResourceAttr("prefix", "/first"),
			domainRoute1.TestCheckResourceAttr("workload_link", workloadSelfLink),

			// Second Route
			domainRoute2.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			domainRoute2.TestCheckResourceAttr("domain_port", "80"),
			domainRoute2.TestCheckResourceAttr("prefix", "/second"),
			domainRoute2.TestCheckResourceAttr("replace_prefix", "/"),
			domainRoute2.TestCheckResourceAttr("workload_link", workloadSelfLink),
			domainRoute2.TestCheckResourceAttr("port", "443"),
			domainRoute2.TestCheckResourceAttr("host_prefix", "my.thing."),
			domainRoute2.TestCheckNestedBlocks("headers", []map[string]interface{}{
				{
					"request": []map[string]interface{}{
						{
							"set": map[string]interface{}{
								"Host":         "example.com",
								"Content-Type": "application/json",
							},
						},
					},
				},
			}),
		),
	}
}

// BuildUpdate4TestStep returns a test step for inline routes creation on the subdomain.
func (drt *DomainResourceTest) BuildUpdate4TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := DomainResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Create the sub-domain test case with inline routes
	subDomainName := fmt.Sprintf("domain-acctest-%s.%s", drt.RandomName, initialCase.Name)
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "domain",
			ResourceName:      "subdomain",
			ResourceAddress:   "cpln_domain.subdomain",
			Name:              subDomainName,
			Description:       subDomainName,
			DescriptionUpdate: "domain with inline routes",
		},
	}

	// Construct the workload self link
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.Update4Hcl(c, subDomain),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Apex Domain
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),

			// Sub Domain with inline routes
			subDomain.GetDefaultChecks(subDomain.DescriptionUpdate, "1"),
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"route": []map[string]interface{}{
								{
									"prefix":        "/api",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
								{
									"prefix":         "/app",
									"replace_prefix": "/",
									"workload_link":  workloadSelfLink,
									"port":           "8080",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-ECDSA-CHACHA20-POLY1305",
										"ECDHE-ECDSA-AES128-GCM-SHA256",
										"ECDHE-RSA-AES256-GCM-SHA384",
										"ECDHE-RSA-CHACHA20-POLY1305",
										"ECDHE-RSA-AES128-GCM-SHA256",
										"AES256-GCM-SHA384",
										"AES128-GCM-SHA256",
									},
								},
							},
						},
					},
				},
			}),
		),
	}
}

// BuildUpdate5TestStep returns a test step for updating inline routes on the subdomain.
func (drt *DomainResourceTest) BuildUpdate5TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := DomainResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Create the sub-domain test case with updated inline routes
	subDomainName := fmt.Sprintf("domain-acctest-%s.%s", drt.RandomName, initialCase.Name)
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "domain",
			ResourceName:      "subdomain",
			ResourceAddress:   "cpln_domain.subdomain",
			Name:              subDomainName,
			Description:       subDomainName,
			DescriptionUpdate: "domain with inline routes updated",
		},
	}

	// Construct the workload self link
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.Update5Hcl(c, subDomain),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Apex Domain
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),

			// Sub Domain with updated inline routes
			subDomain.GetDefaultChecks(subDomain.DescriptionUpdate, "1"),
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"route": []map[string]interface{}{
								{
									"prefix":        "/api",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
								{
									"prefix":         "/app",
									"replace_prefix": "/",
									"workload_link":  workloadSelfLink,
									"port":           "8080",
									"headers": []map[string]interface{}{
										{
											"request": []map[string]interface{}{
												{
													"set": map[string]interface{}{
														"X-Forwarded-Proto": "https",
													},
												},
											},
										},
									},
								},
								{
									"regex":         "/user/.*/profile",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-ECDSA-CHACHA20-POLY1305",
										"ECDHE-ECDSA-AES128-GCM-SHA256",
										"ECDHE-RSA-AES256-GCM-SHA384",
										"ECDHE-RSA-CHACHA20-POLY1305",
										"ECDHE-RSA-AES128-GCM-SHA256",
										"AES256-GCM-SHA384",
										"AES128-GCM-SHA256",
									},
								},
							},
						},
					},
				},
			}),
		),
	}
}

// BuildCoexistenceCreateTestStep returns a test step that creates a subdomain with both inline and external routes.
func (drt *DomainResourceTest) BuildCoexistenceCreateTestStep(subDomainName string) resource.TestStep {
	// Define necessary variables
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))
	subDomainSelfLink := GetSelfLink(OrgName, "domain", subDomainName)

	// Create the sub-domain test case
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Create the domain route test cases
	route1 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			ResourceName:    "route-a",
			ResourceAddress: "cpln_domain_route.route-a",
		},
	}

	route2 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			ResourceName:    "route-b",
			ResourceAddress: "cpln_domain_route.route-b",
		},
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.hclBase() + drt.hclSubDomainWithInlineRoutes(subDomainName, "create with coexistence", []DomainInlineRouteConfig{
			{Prefix: "/inline-a", WorkloadLink: "cpln_workload.new.self_link", Port: 8080},
		}) +
			drt.hclDomainRoute("route-a", "/ext-a", "", 8080) +
			drt.hclDomainRoute("route-b", "/ext-b", "", 8080),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Sub domain spec should show inline routes
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"route": []map[string]interface{}{
								{
									"prefix":        "/inline-a",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-RSA-AES256-GCM-SHA384",
									},
								},
							},
						},
					},
				},
			}),

			// External routes targeting the same subdomain
			route1.TestCheckResourceAttr("domain_link", subDomainSelfLink),
			route1.TestCheckResourceAttr("prefix", "/ext-a"),
			route1.TestCheckResourceAttr("workload_link", workloadSelfLink),
			route1.TestCheckResourceAttr("port", "8080"),
			route2.TestCheckResourceAttr("domain_link", subDomainSelfLink),
			route2.TestCheckResourceAttr("prefix", "/ext-b"),
			route2.TestCheckResourceAttr("workload_link", workloadSelfLink),
			route2.TestCheckResourceAttr("port", "8080"),
		),
	}
}

// BuildNoRoutesTestStep returns a test step for subdomain with no routes.
func (drt *DomainResourceTest) BuildNoRoutesTestStep(subDomainName string, description string) resource.TestStep {
	// Create the sub-domain test case
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.hclBase() + drt.hclSubDomain(subDomainName, description),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("name", subDomainName),
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-RSA-AES256-GCM-SHA384",
									},
								},
							},
						},
					},
				},
			}),
		),
	}
}

// BuildExternalRoutesOnlyTestStep returns a test step with cpln_domain_route resources only (no inline routes).
func (drt *DomainResourceTest) BuildExternalRoutesOnlyTestStep(subDomainName string, description string) resource.TestStep {
	// Define necessary variables
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))
	subDomainSelfLink := GetSelfLink(OrgName, "domain", subDomainName)

	// Create the sub-domain test case
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Create the domain route test cases
	route1 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			ResourceName:    "route-a",
			ResourceAddress: "cpln_domain_route.route-a",
		},
	}

	route2 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			ResourceName:    "route-b",
			ResourceAddress: "cpln_domain_route.route-b",
		},
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.hclBase() + drt.hclSubDomain(subDomainName, description) +
			drt.hclDomainRoute("route-a", "/ext-a", "", 8080) +
			drt.hclDomainRoute("route-b", "/ext-b", "", 8080),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Sub domain spec with no inline routes
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-RSA-AES256-GCM-SHA384",
									},
								},
							},
						},
					},
				},
			}),

			// External routes targeting the same subdomain
			route1.TestCheckResourceAttr("domain_link", subDomainSelfLink),
			route1.TestCheckResourceAttr("prefix", "/ext-a"),
			route1.TestCheckResourceAttr("workload_link", workloadSelfLink),
			route1.TestCheckResourceAttr("port", "8080"),
			route2.TestCheckResourceAttr("domain_link", subDomainSelfLink),
			route2.TestCheckResourceAttr("prefix", "/ext-b"),
			route2.TestCheckResourceAttr("workload_link", workloadSelfLink),
			route2.TestCheckResourceAttr("port", "8080"),
		),
	}
}

// BuildCoexistenceTestStep returns a test step with both inline routes and cpln_domain_route resources.
func (drt *DomainResourceTest) BuildCoexistenceTestStep(subDomainName string) resource.TestStep {
	// Define necessary variables
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))
	subDomainSelfLink := GetSelfLink(OrgName, "domain", subDomainName)

	// Create the sub-domain test case
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Create the domain route test cases
	route1 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			ResourceName:    "route-a",
			ResourceAddress: "cpln_domain_route.route-a",
		},
	}

	route2 := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			ResourceName:    "route-b",
			ResourceAddress: "cpln_domain_route.route-b",
		},
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.hclBase() + drt.hclSubDomainWithInlineRoutes(subDomainName, "coexistence", []DomainInlineRouteConfig{
			{Prefix: "/inline-a", WorkloadLink: "cpln_workload.new.self_link", Port: 8080},
			{Prefix: "/inline-b", WorkloadLink: "cpln_workload.new.self_link", Port: 8080},
		}) +
			drt.hclDomainRoute("route-a", "/ext-a", "", 8080) +
			drt.hclDomainRoute("route-b", "/ext-b", "", 8080),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Sub domain spec should only show inline routes (not external)
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"route": []map[string]interface{}{
								{
									"prefix":        "/inline-a",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
								{
									"prefix":        "/inline-b",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-RSA-AES256-GCM-SHA384",
									},
								},
							},
						},
					},
				},
			}),

			// External routes targeting the same subdomain
			route1.TestCheckResourceAttr("domain_link", subDomainSelfLink),
			route1.TestCheckResourceAttr("prefix", "/ext-a"),
			route1.TestCheckResourceAttr("workload_link", workloadSelfLink),
			route1.TestCheckResourceAttr("port", "8080"),
			route2.TestCheckResourceAttr("domain_link", subDomainSelfLink),
			route2.TestCheckResourceAttr("prefix", "/ext-b"),
			route2.TestCheckResourceAttr("workload_link", workloadSelfLink),
			route2.TestCheckResourceAttr("port", "8080"),
		),
	}
}

// BuildInlineOnlyTestStep returns a test step with only inline routes (no cpln_domain_route).
func (drt *DomainResourceTest) BuildInlineOnlyTestStep(subDomainName string) resource.TestStep {
	// Define necessary variables
	workloadSelfLink := GetSelfLinkWithGvc(OrgName, "workload", fmt.Sprintf("gvc-%s", drt.RandomName), fmt.Sprintf("workload-%s", drt.RandomName))

	// Create the sub-domain test case
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.hclBase() + drt.hclSubDomainWithInlineRoutes(subDomainName, "inline only", []DomainInlineRouteConfig{
			{Prefix: "/inline-a", WorkloadLink: "cpln_workload.new.self_link", Port: 8080},
			{Prefix: "/inline-b", WorkloadLink: "cpln_workload.new.self_link", Port: 8080},
		}),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode":         "ns",
					"accept_all_hosts": "true",
					"ports": []map[string]interface{}{
						{
							"number":   "443",
							"protocol": "http",
							"route": []map[string]interface{}{
								{
									"prefix":        "/inline-a",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
								{
									"prefix":        "/inline-b",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
							},
							"tls": []map[string]interface{}{
								{
									"min_protocol_version": "TLSV1_2",
									"cipher_suites": []string{
										"ECDHE-ECDSA-AES256-GCM-SHA384",
										"ECDHE-RSA-AES256-GCM-SHA384",
									},
								},
							},
						},
					},
				},
			}),
		),
	}
}

// Configs //

// RequiredOnlyHcl returns a minimal HCL block for a resource using only required fields.
func (drt *DomainResourceTest) RequiredOnlyHcl(c DomainResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_domain" "%s" {
  name        = "%s"

  spec {
    ports {
      tls {}
    }
  }
}
`, c.ResourceName, c.Name)
}

// Update1Hcl returns a minimal HCL block for a resource using only required fields.
func (drt *DomainResourceTest) Update1Hcl(c DomainResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_domain" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode              = "cname"
    gvc_link              = "/org/terraform-test-org/gvc/gvc-01"
    cert_challenge_type   = "dns01"
    accept_all_hosts      = false
    accept_all_subdomains = true

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

        allow_origins {						
          regex = "^https://example\\.com$"
        }

        allow_methods     = ["GET", "OPTIONS", "POST"]
        allow_headers     = ["authorization", "host"]
        expose_headers    = ["accept/type"]
        max_age           = "12h"
        allow_credentials = true
      }

      tls {
        min_protocol_version = "TLSV1_1"
        cipher_suites        = ["AES256-GCM-SHA384"]

        client_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }

        server_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }
			}
		}
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// Update2Hcl returns a minimal HCL block for a resource using only required fields.
func (drt *DomainResourceTest) Update2Hcl(c DomainResourceTestCase, subDomain DomainResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

resource "cpln_gvc" "new" {

  name        = "gvc-${var.random_name}"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
  }
}

resource "cpln_workload" "new" {

  gvc = cpln_gvc.new.name

  name        = "workload-${var.random_name}"
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

resource "cpln_domain" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode              = "cname"
    gvc_link              = "/org/terraform-test-org/gvc/gvc-01"
    cert_challenge_type   = "http01"
    accept_all_hosts      = false
    accept_all_subdomains = false

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

        allow_origins {						
          regex = "^https://example\\.com$"
        }

        allow_methods     = ["GET", "OPTIONS", "POST"]
        allow_headers     = ["authorization", "host"]
        expose_headers    = ["accept/type"]
        max_age           = "12h"
        allow_credentials = true
      }

      tls {
        min_protocol_version = "TLSV1_1"
        cipher_suites        = ["AES256-GCM-SHA384"]

        client_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }

        server_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }
			}
		}
  }
}

resource "cpln_domain" "%s" {

  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    terraform_generated = "true"
  }

  spec {
    dns_mode = "ns"
    accept_all_hosts = true

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
  domain_link   = %s
  prefix        = "/first"
  workload_link = cpln_workload.new.self_link
  replica       = 1
}

resource "cpln_domain_route" "second-route" {
  domain_link = cpln_domain.subdomain.self_link
  domain_port = 80

  prefix         = "/second"
  replace_prefix = "/"
  workload_link  = cpln_workload.new.self_link
  port 		       = 443
  host_prefix    = "my.thing."
  replica        = 0

  headers {
    request {
      set = {
        Host           = "example.com"
        "Content-Type" = "application/json"
      }
    }
  }
}

resource "cpln_domain_route" "third-route" {
  domain_link = cpln_domain.subdomain.self_link
  domain_port = 80

  prefix         = "/third"
  replace_prefix = "/"
  workload_link  = cpln_workload.new.self_link
  port 		       = 443
  host_regex     = "reg"

  headers {
    request {
      set = {
        Host           = "example.com"
        "Content-Type" = "application/json"
      }
    }
  }
}

resource "cpln_domain_route" "fourth-route" {
  domain_link   = cpln_domain.subdomain.self_link
  regex         = "/user/.*/profile"
  workload_link = cpln_workload.new.self_link
  port          = 80
}
`, drt.RandomName, c.ResourceName, c.Name, c.DescriptionUpdate, subDomain.ResourceName, c.ResourceAddress, subDomain.Name, subDomain.DescriptionUpdate,
		subDomain.GetSelfLinkAttr(),
	)
}

// Update2Hcl returns a minimal HCL block for a resource using only required fields.
func (drt *DomainResourceTest) Update3Hcl(c DomainResourceTestCase, subDomain DomainResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

resource "cpln_gvc" "new" {

  name        = "gvc-${var.random_name}"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
  }
}

resource "cpln_workload" "new" {

  gvc = cpln_gvc.new.name

  name        = "workload-${var.random_name}"
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

resource "cpln_domain" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode         = "cname"
    gvc_link         = "/org/terraform-test-org/gvc/gvc-01"
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

        allow_origins {						
          regex = "^https://example\\.com$"
        }

        allow_methods     = ["GET", "OPTIONS", "POST"]
        allow_headers     = ["authorization", "host"]
        expose_headers    = ["accept/type"]
        max_age           = "12h"
        allow_credentials = true
      }

      tls {
        min_protocol_version = "TLSV1_1"
        cipher_suites        = ["AES256-GCM-SHA384"]

        client_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }

        server_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }
			}
		}
  }
}

resource "cpln_domain" "%s" {

  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    terraform_generated = "true"
  }

  spec {
    dns_mode = "ns"
    accept_all_hosts = true

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
  domain_link   = %s
  prefix        = "/first"
  workload_link = cpln_workload.new.self_link
}

resource "cpln_domain_route" "second-route" {
  domain_link = cpln_domain.subdomain.self_link
  domain_port = 80

  prefix         = "/second"
  replace_prefix = "/"
  workload_link  = cpln_workload.new.self_link
  port 		       = 443
  host_prefix    = "my.thing."

  headers {
    request {
      set = {
        Host           = "example.com"
        "Content-Type" = "application/json"
      }
    }
  }
}
`, drt.RandomName, c.ResourceName, c.Name, c.DescriptionUpdate, subDomain.ResourceName, c.ResourceAddress, subDomain.Name, subDomain.DescriptionUpdate,
		subDomain.GetSelfLinkAttr(),
	)
}

// Update4Hcl returns an HCL block for a subdomain with inline routes (replaces external cpln_domain_route resources).
func (drt *DomainResourceTest) Update4Hcl(c DomainResourceTestCase, subDomain DomainResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

resource "cpln_gvc" "new" {

  name        = "gvc-${var.random_name}"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
  }
}

resource "cpln_workload" "new" {

  gvc = cpln_gvc.new.name

  name        = "workload-${var.random_name}"
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

resource "cpln_domain" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode         = "cname"
    gvc_link         = "/org/terraform-test-org/gvc/gvc-01"
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

        allow_origins {
          regex = "^https://example\\.com$"
        }

        allow_methods     = ["GET", "OPTIONS", "POST"]
        allow_headers     = ["authorization", "host"]
        expose_headers    = ["accept/type"]
        max_age           = "12h"
        allow_credentials = true
      }

      tls {
        min_protocol_version = "TLSV1_1"
        cipher_suites        = ["AES256-GCM-SHA384"]

        client_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }

        server_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }
			}
		}
  }
}

resource "cpln_domain" "%s" {

  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
  }

  spec {
    dns_mode = "ns"
    accept_all_hosts = true

    ports {
      number   = 443
      protocol = "http"

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

      route {
        prefix        = "/api"
        workload_link = cpln_workload.new.self_link
        port          = 8080
      }

      route {
        prefix         = "/app"
        replace_prefix = "/"
        workload_link  = cpln_workload.new.self_link
        port           = 8080
      }
    }
  }
}
`, drt.RandomName, c.ResourceName, c.Name, c.DescriptionUpdate, subDomain.ResourceName, c.ResourceAddress, subDomain.Name, subDomain.DescriptionUpdate)
}

// Update5Hcl returns an HCL block for updating inline routes (add a third route, add headers to second route).
func (drt *DomainResourceTest) Update5Hcl(c DomainResourceTestCase, subDomain DomainResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

resource "cpln_gvc" "new" {

  name        = "gvc-${var.random_name}"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
  }
}

resource "cpln_workload" "new" {

  gvc = cpln_gvc.new.name

  name        = "workload-${var.random_name}"
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

resource "cpln_domain" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode         = "cname"
    gvc_link         = "/org/terraform-test-org/gvc/gvc-01"
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

        allow_origins {
          regex = "^https://example\\.com$"
        }

        allow_methods     = ["GET", "OPTIONS", "POST"]
        allow_headers     = ["authorization", "host"]
        expose_headers    = ["accept/type"]
        max_age           = "12h"
        allow_credentials = true
      }

      tls {
        min_protocol_version = "TLSV1_1"
        cipher_suites        = ["AES256-GCM-SHA384"]

        client_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }

        server_certificate {
          secret_link = "/org/terraform-test-org/secret/aa-tbd-2"
        }
			}
		}
  }
}

resource "cpln_domain" "%s" {

  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
  }

  spec {
    dns_mode = "ns"
    accept_all_hosts = true

    ports {
      number   = 443
      protocol = "http"

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

      route {
        prefix        = "/api"
        workload_link = cpln_workload.new.self_link
        port          = 8080
      }

      route {
        prefix         = "/app"
        replace_prefix = "/"
        workload_link  = cpln_workload.new.self_link
        port           = 8080

        headers {
          request {
            set = {
              "X-Forwarded-Proto" = "https"
            }
          }
        }
      }

      route {
        regex         = "/user/.*/profile"
        workload_link = cpln_workload.new.self_link
        port          = 8080
      }
    }
  }
}
`, drt.RandomName, c.ResourceName, c.Name, c.DescriptionUpdate, subDomain.ResourceName, c.ResourceAddress, subDomain.Name, subDomain.DescriptionUpdate)
}

// hclBase returns the shared infrastructure HCL (GVC, workload, apex domain).
func (drt *DomainResourceTest) hclBase() string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

resource "cpln_gvc" "new" {
  name        = "gvc-${var.random_name}"
  description = "Route coexistence test GVC"
  locations   = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
  }
}

resource "cpln_workload" "new" {
  gvc         = cpln_gvc.new.name
  name        = "workload-${var.random_name}"
  description = "Route coexistence test workload"
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

resource "cpln_domain" "new" {
  name = "%s"

  spec {
    ports {
      tls {}
    }
  }
}
`, drt.RandomName, drt.ApexDomain)
}

// hclSubDomain returns HCL for the subdomain without inline routes.
func (drt *DomainResourceTest) hclSubDomain(name string, description string) string {
	return fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on  = [cpln_domain.new]
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
  }

  spec {
    dns_mode         = "ns"
    accept_all_hosts = true

    ports {
      number   = 443
      protocol = "http"

      tls {
        min_protocol_version = "TLSV1_2"
        cipher_suites = [
          "ECDHE-ECDSA-AES256-GCM-SHA384",
          "ECDHE-RSA-AES256-GCM-SHA384",
        ]
      }
    }
  }
}
`, name, description)
}

// hclSubDomainWithInlineRoutes returns HCL for the subdomain with inline routes.
func (drt *DomainResourceTest) hclSubDomainWithInlineRoutes(name string, description string, routes []DomainInlineRouteConfig) string {
	var routeBlocks strings.Builder
	for _, r := range routes {
		if r.Prefix != "" {
			routeBlocks.WriteString(fmt.Sprintf(`
      route {
        prefix        = "%s"
        workload_link = %s
        port          = %d
      }
`, r.Prefix, r.WorkloadLink, r.Port))
		} else {
			routeBlocks.WriteString(fmt.Sprintf(`
      route {
        regex         = "%s"
        workload_link = %s
        port          = %d
      }
`, r.Regex, r.WorkloadLink, r.Port))
		}
	}

	return fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on  = [cpln_domain.new]
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
  }

  spec {
    dns_mode         = "ns"
    accept_all_hosts = true

    ports {
      number   = 443
      protocol = "http"

      tls {
        min_protocol_version = "TLSV1_2"
        cipher_suites = [
          "ECDHE-ECDSA-AES256-GCM-SHA384",
          "ECDHE-RSA-AES256-GCM-SHA384",
        ]
      }
%s
    }
  }
}
`, name, description, routeBlocks.String())
}

// hclDomainRoute returns HCL for a cpln_domain_route resource.
func (drt *DomainResourceTest) hclDomainRoute(resourceName string, prefix string, regex string, port int) string {
	var routeKey string
	if prefix != "" {
		routeKey = fmt.Sprintf(`prefix        = "%s"`, prefix)
	} else {
		routeKey = fmt.Sprintf(`regex         = "%s"`, regex)
	}

	return fmt.Sprintf(`
resource "cpln_domain_route" "%s" {
  domain_link   = cpln_domain.subdomain.self_link
  %s
  workload_link = cpln_workload.new.self_link
  port          = %d
}
`, resourceName, routeKey, port)
}

/*** Resource Test Cases ***/

// DomainResourceTestCase defines a specific resource test case.
type DomainResourceTestCase struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (drtc *DomainResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of domain: %s. Total resources: %d", drtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[drtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", drtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != drtc.Name {
			return fmt.Errorf("resource ID %s does not match expected domain name %s", rs.Primary.ID, drtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteDomain, _, err := TestProvider.client.GetDomain(drtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving domain from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteDomain.Name != drtc.Name {
			return fmt.Errorf("mismatch in domain name: expected %s, got %s", drtc.Name, *remoteDomain.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("Domain %s verified successfully in both state and external system.", drtc.Name))
		return nil
	}
}

// DomainRouteResourceTestCase defines a specific resource test case.
type DomainRouteResourceTestCase struct {
	ProviderTestCase
}
