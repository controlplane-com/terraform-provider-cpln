package cpln

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
	WorkloadLink string // HCL expression (e.g., `"//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"`)
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
	steps = append(steps, resourceTest.NewOptionalTlsScenario()...)
	steps = append(steps, resourceTest.NewCanaryLifecycleScenario()...)
	steps = append(steps, resourceTest.NewHostRouteLifecycleScenario()...)
	steps = append(steps, resourceTest.NewRouteOrderLifecycleScenario()...)
	steps = append(steps, resourceTest.NewHttpsPortWithoutTlsScenario()...)

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

	// Build a revert step that uses hclBase() to keep GVC/workload alive for the coexistence scenario,
	// while still reverting the domain to its minimal config (hclBase includes the same minimal domain).
	revertStep := resource.TestStep{
		Config: drt.hclBase(),
		Check:  initialStep.Check,
	}

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
		// Domain Route Import (using domain name instead of full link)
		// Verifies that bare domain names are normalized to full self-links,
		// preventing a RequiresReplace diff on the next plan.
		{
			ResourceName:     "cpln_domain_route.first-route",
			ImportState:      true,
			ImportStateId:    fmt.Sprintf("domain-acctest-%s.%s:443:/first", drt.RandomName, name),
			ImportStateCheck: domainRouteLinkCheck(subDomainSelfLink),
		},
		{
			ResourceName:     "cpln_domain_route.second-route",
			ImportState:      true,
			ImportStateId:    fmt.Sprintf("domain-acctest-%s.%s:80:/second", drt.RandomName, name),
			ImportStateCheck: domainRouteLinkCheck(subDomainSelfLink),
		},
		{
			ResourceName:     "cpln_domain_route.third-route",
			ImportState:      true,
			ImportStateId:    fmt.Sprintf("domain-acctest-%s.%s:80:/third", drt.RandomName, name),
			ImportStateCheck: domainRouteLinkCheck(subDomainSelfLink),
		},
		{
			ResourceName:     "cpln_domain_route.fourth-route",
			ImportState:      true,
			ImportStateId:    fmt.Sprintf("domain-acctest-%s.%s:443:/user/.*/profile", drt.RandomName, name),
			ImportStateCheck: domainRouteLinkCheck(subDomainSelfLink),
		},
		// Inline Routes: Create & Read
		caseUpdate4,
		// Inline Routes: Update & Read
		caseUpdate5,
		// Revert the domain to its initial state (keep GVC/workload for coexistence scenario)
		revertStep,
	}
}

// NewCoexistenceScenario creates a test scenario covering all unique route coexistence transitions.
func (drt *DomainResourceTest) NewCoexistenceScenario() []resource.TestStep {
	// Define necessary variables
	subDomainName := fmt.Sprintf("route-coexist-%s.%s", drt.RandomName, drt.ApexDomain)

	// Build test steps
	caseCreate := drt.BuildCoexistenceCreateTestStep(subDomainName)
	caseStripRoutes := drt.BuildNoRoutesTestStep(subDomainName, "no routes")
	caseExternalOnly := drt.BuildExternalRoutesOnlyTestStep(subDomainName, "external routes only")
	caseCoexistence := drt.BuildCoexistenceTestStep(subDomainName)
	caseInlineOnly := drt.BuildInlineOnlyTestStep(subDomainName)
	caseCleanup := drt.BuildNoRoutesTestStep(subDomainName, "final cleanup")

	// Return the complete test steps
	return []resource.TestStep{
		// Create domain + inline routes + external routes in a single apply
		caseCreate,
		// Remove inline route blocks and external route resources — inline routes
		// persist on the API (mergeRoutes preserves all API routes when no inline
		// blocks are defined) but are not stored in Terraform state
		caseStripRoutes,
		// Add external routes only (domain has no inline route blocks)
		caseExternalOnly,
		// Add inline routes alongside external (coexistence)
		caseCoexistence,
		// Remove external, keep inline only
		caseInlineOnly,
		// Import State — verify routes are stored in state (backward compat)
		{
			ResourceName:     "cpln_domain.subdomain",
			ImportState:      true,
			ImportStateId:    subDomainName,
			ImportStateCheck: domainImportWithRoutesCheck(),
		},
		// Cleanup: remove inline route blocks and external route resources
		caseCleanup,
	}
}

// NewOptionalTlsScenario creates a test scenario covering a TCP subdomain that omits the optional tls block.
func (drt *DomainResourceTest) NewOptionalTlsScenario() []resource.TestStep {
	// Define necessary variables
	subDomainName := fmt.Sprintf("tcp-no-tls-%s.%s", drt.RandomName, drt.ApexDomain)

	// Build test steps
	caseCreate := drt.BuildOptionalTlsCreateTestStep(subDomainName)
	casePlanStable := drt.BuildOptionalTlsPlanStableTestStep(subDomainName)

	// Return the complete test steps
	return []resource.TestStep{
		// Create a TCP subdomain that omits the optional tls block
		caseCreate,
		// Re-plan the same configuration to confirm the API response does not
		// reintroduce a tls block in state and the plan stays empty
		casePlanStable,
	}
}

// NewCanaryLifecycleScenario walks a single route's canary block through every cardinality transition on a standalone cpln_domain_route.
func (drt *DomainResourceTest) NewCanaryLifecycleScenario() []resource.TestStep {
	// Define the subdomain that hosts the canary route across the lifecycle
	subDomainName := fmt.Sprintf("canary-life-%s.%s", drt.RandomName, drt.ApexDomain)

	// Build the per-stage test steps
	absentStep := drt.BuildCanaryAbsentTestStep(subDomainName)
	requiredOnlyStep := drt.BuildCanaryRequiredOnlyTestStep(subDomainName)
	allMultiStep := drt.BuildCanaryAllMultiTestStep(subDomainName)
	expandedStep := drt.BuildCanaryExpandedTestStep(subDomainName)

	// Walk the canary block: absent -> required-only -> all-attrs+multiple(incl weight 0) -> expand -> shrink -> remove
	return []resource.TestStep{
		// Route exists with no canary blocks
		absentStep,
		// One canary with only its required attributes (workload_link + weight)
		requiredOnlyStep,
		// Two canaries with every attribute set, including a weight 0 (disabled-but-retained) toggle
		allMultiStep,
		// Grow to three canaries so the build/flatten loop runs at a non-trivial count
		expandedStep,
		// Shrink back to two canaries (re-use the all-attrs step)
		allMultiStep,
		// Remove every canary while the route persists (re-use the absent step)
		absentStep,
	}
}

// NewHostRouteLifecycleScenario walks standalone cpln_domain_route resources that share one path across several hosts.
func (drt *DomainResourceTest) NewHostRouteLifecycleScenario() []resource.TestStep {
	// Define the subdomain that hosts the routes across the lifecycle
	subDomainName := fmt.Sprintf("host-route-%s.%s", drt.RandomName, drt.ApexDomain)

	// Define the self link the import identifiers address
	subDomainSelfLink := GetSelfLink(OrgName, "domain", subDomainName)

	// Build the per-stage test steps
	catchAllStep := drt.BuildHostRouteCatchAllTestStep(subDomainName)
	multiHostStep := drt.BuildHostRouteMultiHostTestStep(subDomainName)
	swappedStep := drt.BuildHostRouteSwappedTestStep(subDomainName)
	renamedHostStep := drt.BuildHostRouteRenamedHostTestStep(subDomainName)
	expandedStep := drt.BuildHostRouteExpandedTestStep(subDomainName)

	// Walk routes that share a path: every host -> three hosts -> retarget -> rename a host -> expand -> shrink -> remove
	return []resource.TestStep{
		// One route owns the path and matches every host
		catchAllStep,
		// Two more routes take the same path on their own hosts, which the API stores by descending host length
		multiHostStep,
		// Import the route on one host, which the path alone cannot address
		{
			ResourceName:            "cpln_domain_route.store",
			ImportState:             true,
			ImportStateId:           fmt.Sprintf("%s:443:/:store.", subDomainSelfLink),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"workload_link"},
		},
		// Import the route that matches every host, which a trailing empty host segment addresses
		{
			ResourceName:            "cpln_domain_route.catch-all",
			ImportState:             true,
			ImportStateId:           fmt.Sprintf("%s:443:/:", subDomainSelfLink),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"workload_link"},
		},
		// Import by domain name rather than self link, with the host segment still resolving the route
		{
			ResourceName:     "cpln_domain_route.blog",
			ImportState:      true,
			ImportStateId:    fmt.Sprintf("%s:443:/:blog.", subDomainName),
			ImportStateCheck: domainRouteLinkCheck(subDomainSelfLink),
		},
		// Refuse to import when the path alone addresses three routes
		{
			ResourceName:  "cpln_domain_route.store",
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:443:/", subDomainSelfLink),
			ExpectError:   regexp.MustCompile(`Append\s+the\s+host\s+to\s+the\s+import\s+identifier`),
		},
		// Swap the workloads the two host routes point at, which only lands correctly when each update finds its own route
		swappedStep,
		// Move a route to another host in place, without destroying and recreating it
		renamedHostStep,
		// A fourth route takes the same path on a host regex
		expandedStep,
		// Shrink back by dropping the host regex route and returning the renamed route to its host (re-use the swapped step)
		swappedStep,
		// Remove every host route while the route that matches every host persists (re-use the first step)
		catchAllStep,
	}
}

// NewRouteOrderLifecycleScenario walks a port's inline route blocks through the cardinality and ordering transitions the API's route sort disturbs.
func (drt *DomainResourceTest) NewRouteOrderLifecycleScenario() []resource.TestStep {
	// Define the subdomain that hosts the inline routes across the lifecycle
	subDomainName := fmt.Sprintf("route-order-%s.%s", drt.RandomName, drt.ApexDomain)

	// Build the per-stage test steps
	absentStep := drt.BuildRouteOrderAbsentTestStep(subDomainName)
	requiredOnlyStep := drt.BuildRouteOrderRequiredOnlyTestStep(subDomainName)
	unsortedMultiStep := drt.BuildRouteOrderUnsortedMultiTestStep(subDomainName)
	withExternalStep := drt.BuildRouteOrderWithExternalTestStep(subDomainName)
	expandedStep := drt.BuildRouteOrderExpandedTestStep(subDomainName)
	regexStep := drt.BuildRouteOrderRegexTestStep(subDomainName)

	// Walk the inline route blocks: absent -> required-only -> unsorted multiple -> external alongside -> expand -> shrink -> regex -> remove
	return []resource.TestStep{
		// Port carries no inline route blocks
		absentStep,
		// One route with only its required attributes (prefix + workload_link)
		requiredOnlyStep,
		// Five routes declared shortest prefix first, the exact order the API rewrites, including two routes that share a prefix on different hosts
		unsortedMultiStep,
		// Keep the same five routes while an externally owned cpln_domain_route sorts into the middle of them
		withExternalStep,
		// Grow to six routes declared in another order the API rewrites
		expandedStep,
		// Shrink back to five routes (re-use the unsorted multiple step)
		unsortedMultiStep,
		// A regex route on the port switches the API's route sort off, so the declared order must survive untouched
		regexStep,
		// Remove every inline route while the port persists (re-use the absent step)
		absentStep,
	}
}

// NewHttpsPortWithoutTlsScenario covers an https port that omits the optional tls block the API fills in server side.
func (drt *DomainResourceTest) NewHttpsPortWithoutTlsScenario() []resource.TestStep {
	// Define the subdomain that hosts the https port across the scenario
	subDomainName := fmt.Sprintf("https-no-tls-%s.%s", drt.RandomName, drt.ApexDomain)

	// Build the per-stage test steps
	withoutTlsStep := drt.BuildHttpsPortWithoutTlsTestStep(subDomainName)
	planStableStep := drt.BuildHttpsPortWithoutTlsPlanStableTestStep(subDomainName)
	mixedTlsStep := drt.BuildHttpsPortMixedTlsTestStep(subDomainName)
	explicitTlsStep := drt.BuildHttpsPortExplicitTlsTestStep(subDomainName)

	// Walk the tls block: absent -> stable plan -> absent beside a port that sets one -> set -> absent again
	return []resource.TestStep{
		// Port 443 omits tls while the API fills in its default tls configuration
		withoutTlsStep,
		// Re-plan the same configuration to confirm the API default never reaches state
		planStableStep,
		// Port 443 omits tls while port 80 sets one, so the suppression is resolved per port
		mixedTlsStep,
		// Port 443 sets tls explicitly and the configured values are stored
		explicitTlsStep,
		// Remove the tls block again (re-use the first step)
		withoutTlsStep,
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

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
			domainRoute2.TestCheckNestedBlocks("mirror", []map[string]interface{}{
				{
					"workload_link": workloadSelfLink,
					"port":          "8080",
					"percent":       "50",
				},
				{
					"workload_link": workloadSelfLink,
					"percent":       "25.5",
				},
			}),
			domainRoute2.TestCheckNestedBlocks("canary", []map[string]interface{}{
				{
					"workload_link": workloadSelfLink,
					"port":          "8080",
					"weight":        "50",
				},
				{
					"workload_link": workloadSelfLink,
					"weight":        "25",
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

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
			domainRoute2.TestCheckNestedBlocks("mirror", []map[string]interface{}{
				{
					"workload_link": workloadSelfLink,
					"port":          "8080",
					"percent":       "50",
				},
				{
					"workload_link": workloadSelfLink,
					"percent":       "25.5",
				},
			}),
			domainRoute2.TestCheckNestedBlocks("canary", []map[string]interface{}{
				{
					"workload_link": workloadSelfLink,
					"port":          "8080",
					"weight":        "50",
				},
				{
					"workload_link": workloadSelfLink,
					"weight":        "25",
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

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
									"mirror": []map[string]interface{}{
										{
											"workload_link": workloadSelfLink,
											"port":          "8080",
											"percent":       "50",
										},
										{
											"workload_link": workloadSelfLink,
											"percent":       "25.5",
										},
									},
									"canary": []map[string]interface{}{
										{
											"workload_link": workloadSelfLink,
											"port":          "8080",
											"weight":        "50",
										},
										{
											"workload_link": workloadSelfLink,
											"weight":        "25",
										},
									},
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

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
									"mirror": []map[string]interface{}{
										{
											"workload_link": workloadSelfLink,
											"port":          "8080",
											"percent":       "75",
										},
									},
									"canary": []map[string]interface{}{
										{
											"workload_link": workloadSelfLink,
											"port":          "8080",
											"weight":        "75",
										},
									},
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
									"mirror": []map[string]interface{}{
										{
											"workload_link": workloadSelfLink,
											"percent":       "30",
										},
									},
									"canary": []map[string]interface{}{
										{
											"workload_link": workloadSelfLink,
											"weight":        "30",
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)
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
			{Prefix: "/inline-a", WorkloadLink: `"//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"`, Port: 8080},
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)
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
			{Prefix: "/inline-a", WorkloadLink: `"//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"`, Port: 8080},
			{Prefix: "/inline-b", WorkloadLink: `"//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"`, Port: 8080},
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
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

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
			{Prefix: "/inline-a", WorkloadLink: `"//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"`, Port: 8080},
			{Prefix: "/inline-b", WorkloadLink: `"//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"`, Port: 8080},
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

// BuildOptionalTlsCreateTestStep returns a test step that creates a TCP subdomain without a tls block.
func (drt *DomainResourceTest) BuildOptionalTlsCreateTestStep(subDomainName string) resource.TestStep {
	// Define necessary variables
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

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
		Config: drt.hclBase() + drt.hclSubDomainTcpNoTls(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			// Sub domain spec should expose a single TCP port routed to the workload with no tls block
			subDomain.TestCheckNestedBlocks("spec", []map[string]interface{}{
				{
					"dns_mode": "cname",
					"ports": []map[string]interface{}{
						{
							"number":   "5432",
							"protocol": "tcp",
							"route": []map[string]interface{}{
								{
									"prefix":        "/",
									"workload_link": workloadSelfLink,
									"port":          "8080",
								},
							},
						},
					},
				},
			}),

			// Strictly assert that no tls block was persisted to state
			subDomain.TestCheckResourceAttr("spec.0.ports.0.tls.#", "0"),
		),
	}
}

// BuildOptionalTlsPlanStableTestStep returns a plan-only test step that re-runs the no-tls config to verify no plan drift.
func (drt *DomainResourceTest) BuildOptionalTlsPlanStableTestStep(subDomainName string) resource.TestStep {
	// Initialize and return the test step
	return resource.TestStep{
		Config:   drt.hclBase() + drt.hclSubDomainTcpNoTls(subDomainName),
		PlanOnly: true,
	}
}

// newCanaryLifecycleCases builds the subdomain and route test cases shared across the canary lifecycle stages.
func (drt *DomainResourceTest) newCanaryLifecycleCases(subDomainName string) (DomainResourceTestCase, DomainRouteResourceTestCase, string) {
	// Build the subdomain case used to resolve the route's domain_link self link
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Build the standalone route case that owns the canary blocks
	route := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "canary-route",
			ResourceAddress: "cpln_domain_route.canary-route",
		},
	}

	// Construct the workload self link in the short form the config uses
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

	// Return the shared cases and link
	return subDomain, route, workloadSelfLink
}

// BuildCanaryAbsentTestStep returns a step where the route exists with no canary blocks.
func (drt *DomainResourceTest) BuildCanaryAbsentTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload link
	subDomain, route, workloadSelfLink := drt.newCanaryLifecycleCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.CanaryAbsentHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			route.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			route.TestCheckResourceAttr("domain_port", "443"),
			route.TestCheckResourceAttr("prefix", "/canary"),
			route.TestCheckResourceAttr("workload_link", workloadSelfLink),
			// No canary blocks are present on the route
			route.TestCheckResourceAttr("canary.#", "0"),
		),
	}
}

// BuildCanaryRequiredOnlyTestStep returns a step with a single canary that sets only its required attributes.
func (drt *DomainResourceTest) BuildCanaryRequiredOnlyTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload link
	subDomain, route, workloadSelfLink := drt.newCanaryLifecycleCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.CanaryRequiredOnlyHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			route.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			route.TestCheckResourceAttr("prefix", "/canary"),
			route.TestCheckNestedBlocks("canary", []map[string]interface{}{
				{
					"workload_link": workloadSelfLink,
					"weight":        "50",
				},
			}),
		),
	}
}

// BuildCanaryAllMultiTestStep returns a step with two canaries covering every attribute and a weight 0 toggle.
func (drt *DomainResourceTest) BuildCanaryAllMultiTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload link
	subDomain, route, workloadSelfLink := drt.newCanaryLifecycleCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.CanaryAllMultiHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			route.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			route.TestCheckResourceAttr("prefix", "/canary"),
			route.TestCheckNestedBlocks("canary", []map[string]interface{}{
				{
					"workload_link": workloadSelfLink,
					"port":          "8080",
					"weight":        "60",
				},
				{
					"workload_link": workloadSelfLink,
					"weight":        "0",
				},
			}),
		),
	}
}

// BuildCanaryExpandedTestStep returns a step with three canaries so the build/flatten loop runs at a non-trivial count.
func (drt *DomainResourceTest) BuildCanaryExpandedTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload link
	subDomain, route, workloadSelfLink := drt.newCanaryLifecycleCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.CanaryExpandedHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			route.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			route.TestCheckResourceAttr("prefix", "/canary"),
			route.TestCheckNestedBlocks("canary", []map[string]interface{}{
				{
					"workload_link": workloadSelfLink,
					"port":          "8080",
					"weight":        "30",
				},
				{
					"workload_link": workloadSelfLink,
					"weight":        "20",
				},
				{
					"workload_link": workloadSelfLink,
					"port":          "8080",
					"weight":        "0",
				},
			}),
		),
	}
}

// newHostRouteCases builds the subdomain test case, the route test cases, and the workload links shared across the host route lifecycle stages.
func (drt *DomainResourceTest) newHostRouteCases(subDomainName string) (DomainResourceTestCase, DomainRouteResourceTestCase, DomainRouteResourceTestCase, DomainRouteResourceTestCase, string, string) {
	// Build the subdomain case used to resolve each route's domain_link self link
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Build the route case that matches every host
	catchAll := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "catch-all",
			ResourceAddress: "cpln_domain_route.catch-all",
		},
	}

	// Build the route case that narrows the same path to one host
	store := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "store",
			ResourceAddress: "cpln_domain_route.store",
		},
	}

	// Build the route case that narrows the same path to another host
	blog := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "blog",
			ResourceAddress: "cpln_domain_route.blog",
		},
	}

	// Construct both workload self links in the short form the config uses
	primaryLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)
	alternateLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-alt-%s", drt.RandomName, drt.RandomName)

	// Return the shared cases and links
	return subDomain, catchAll, store, blog, primaryLink, alternateLink
}

// BuildHostRouteCatchAllTestStep returns a step where a single route owns the path and matches every host.
func (drt *DomainResourceTest) BuildHostRouteCatchAllTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload links
	subDomain, catchAll, _, _, primaryLink, _ := drt.newHostRouteCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HostRouteCatchAllHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			catchAll.TestCheckResourceAttr("id", fmt.Sprintf("%s_443_/", subDomain.GetSelfLink())),
			catchAll.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			catchAll.TestCheckResourceAttr("domain_port", "443"),
			catchAll.TestCheckResourceAttr("prefix", "/"),
			catchAll.TestCheckResourceAttr("workload_link", primaryLink),
			catchAll.TestCheckResourceAttr("port", "8080"),
			// No host narrows the route, so its identifier stays the one earlier provider versions issued
			resource.TestCheckNoResourceAttr(catchAll.ResourceAddress, "host_prefix"),
			// The port carries the one route and nothing else
			domainRouteHostMappingCheck(subDomainName, 443, map[string]string{
				"": fmt.Sprintf("workload-%s", drt.RandomName),
			}),
		),
	}
}

// BuildHostRouteMultiHostTestStep returns a step where three routes share the path, two of them narrowed to a host.
func (drt *DomainResourceTest) BuildHostRouteMultiHostTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload links
	subDomain, catchAll, store, blog, primaryLink, alternateLink := drt.newHostRouteCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HostRouteMultiHostHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			// The route that matches every host keeps its own values while the host routes exist beside it
			catchAll.TestCheckResourceAttr("id", fmt.Sprintf("%s_443_/", subDomain.GetSelfLink())),
			catchAll.TestCheckResourceAttr("prefix", "/"),
			catchAll.TestCheckResourceAttr("workload_link", primaryLink),
			// The route on the store host carries the host in its identifier
			store.TestCheckResourceAttr("id", fmt.Sprintf("%s_443_/_store.", subDomain.GetSelfLink())),
			store.TestCheckResourceAttr("prefix", "/"),
			store.TestCheckResourceAttr("host_prefix", "store."),
			store.TestCheckResourceAttr("workload_link", alternateLink),
			store.TestCheckResourceAttr("port", "8080"),
			// The route on the blog host carries its own host in its identifier
			blog.TestCheckResourceAttr("id", fmt.Sprintf("%s_443_/_blog.", subDomain.GetSelfLink())),
			blog.TestCheckResourceAttr("prefix", "/"),
			blog.TestCheckResourceAttr("host_prefix", "blog."),
			blog.TestCheckResourceAttr("workload_link", primaryLink),
			// Each host on the port reaches the workload its own route declares
			domainRouteHostMappingCheck(subDomainName, 443, map[string]string{
				"":       fmt.Sprintf("workload-%s", drt.RandomName),
				"store.": fmt.Sprintf("workload-alt-%s", drt.RandomName),
				"blog.":  fmt.Sprintf("workload-%s", drt.RandomName),
			}),
		),
	}
}

// BuildHostRouteSwappedTestStep returns a step where the two host routes exchange the workloads they point at.
func (drt *DomainResourceTest) BuildHostRouteSwappedTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload links
	subDomain, catchAll, store, blog, primaryLink, alternateLink := drt.newHostRouteCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HostRouteSwappedHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			// The route that matches every host is untouched by the updates beside it
			catchAll.TestCheckResourceAttr("id", fmt.Sprintf("%s_443_/", subDomain.GetSelfLink())),
			catchAll.TestCheckResourceAttr("workload_link", primaryLink),
			// The store route now reaches the primary workload and rewrites the path it forwards
			store.TestCheckResourceAttr("host_prefix", "store."),
			store.TestCheckResourceAttr("workload_link", primaryLink),
			store.TestCheckResourceAttr("replace_prefix", "/shop"),
			// The blog route now reaches the alternate workload
			blog.TestCheckResourceAttr("host_prefix", "blog."),
			blog.TestCheckResourceAttr("workload_link", alternateLink),
			// Each update landed on its own route rather than on the first route sharing the path
			domainRouteHostMappingCheck(subDomainName, 443, map[string]string{
				"":       fmt.Sprintf("workload-%s", drt.RandomName),
				"store.": fmt.Sprintf("workload-%s", drt.RandomName),
				"blog.":  fmt.Sprintf("workload-alt-%s", drt.RandomName),
			}),
		),
	}
}

// BuildHostRouteRenamedHostTestStep returns a step where a route moves to another host without being destroyed and recreated.
func (drt *DomainResourceTest) BuildHostRouteRenamedHostTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload links
	subDomain, _, store, blog, primaryLink, alternateLink := drt.newHostRouteCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HostRouteRenamedHostHcl(subDomainName),
		ConfigPlanChecks: resource.ConfigPlanChecks{
			PreApply: []plancheck.PlanCheck{
				// Moving a route to another host updates it in place, it never destroys and recreates it
				plancheck.ExpectResourceAction(blog.ResourceAddress, plancheck.ResourceActionUpdate),
			},
		},
		Check: resource.ComposeAggregateTestCheckFunc(
			// The renamed route carries its new host in its identifier
			blog.TestCheckResourceAttr("id", fmt.Sprintf("%s_443_/_news.", subDomain.GetSelfLink())),
			blog.TestCheckResourceAttr("host_prefix", "news."),
			blog.TestCheckResourceAttr("workload_link", alternateLink),
			// The route on the host that was left alone keeps its own values
			store.TestCheckResourceAttr("host_prefix", "store."),
			store.TestCheckResourceAttr("workload_link", primaryLink),
			// The port now reaches the renamed host and no longer reaches the old one
			domainRouteHostMappingCheck(subDomainName, 443, map[string]string{
				"":       fmt.Sprintf("workload-%s", drt.RandomName),
				"store.": fmt.Sprintf("workload-%s", drt.RandomName),
				"news.":  fmt.Sprintf("workload-alt-%s", drt.RandomName),
			}),
		),
	}
}

// BuildHostRouteExpandedTestStep returns a step where a fourth route takes the same path on a host regex.
func (drt *DomainResourceTest) BuildHostRouteExpandedTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared cases and workload links
	subDomain, catchAll, _, _, primaryLink, alternateLink := drt.newHostRouteCases(subDomainName)

	// Build the route case that narrows the same path to a host regex
	api := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "api",
			ResourceAddress: "cpln_domain_route.api",
		},
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HostRouteExpandedHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			// The host regex route carries its regex in its identifier
			api.TestCheckResourceAttr("id", fmt.Sprintf("%s_443_/_^api.*$", subDomain.GetSelfLink())),
			api.TestCheckResourceAttr("prefix", "/"),
			api.TestCheckResourceAttr("host_regex", "^api.*$"),
			api.TestCheckResourceAttr("workload_link", alternateLink),
			// Four routes now share the path, each reaching the workload its own route declares
			domainRouteHostMappingCheck(subDomainName, 443, map[string]string{
				"":        fmt.Sprintf("workload-%s", drt.RandomName),
				"store.":  fmt.Sprintf("workload-%s", drt.RandomName),
				"news.":   fmt.Sprintf("workload-alt-%s", drt.RandomName),
				"^api.*$": fmt.Sprintf("workload-alt-%s", drt.RandomName),
			}),
			// The route that matches every host is still reachable under its own identifier
			catchAll.TestCheckResourceAttr("workload_link", primaryLink),
		),
	}
}

// newRouteOrderCases builds the subdomain test case and workload link shared across the route order lifecycle stages.
func (drt *DomainResourceTest) newRouteOrderCases(subDomainName string) (DomainResourceTestCase, string) {
	// Build the subdomain case that owns the inline route blocks
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Construct the workload self link in the short form the config uses
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

	// Return the shared case and link
	return subDomain, workloadSelfLink
}

// BuildRouteOrderAbsentTestStep returns a step where the https port carries no inline route blocks.
func (drt *DomainResourceTest) BuildRouteOrderAbsentTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, _ := drt.newRouteOrderCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.RouteOrderAbsentHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("name", subDomainName),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.number", "443"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.protocol", "http"),
			// No inline route blocks are present on the port
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.#", "0"),
		),
	}
}

// BuildRouteOrderRequiredOnlyTestStep returns a step with a single route that sets only its required attributes.
func (drt *DomainResourceTest) BuildRouteOrderRequiredOnlyTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, workloadSelfLink := drt.newRouteOrderCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.RouteOrderRequiredOnlyHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.#", "1"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.prefix", "/"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.workload_link", workloadSelfLink),
		),
	}
}

// BuildRouteOrderUnsortedMultiTestStep returns a step with five routes declared shortest prefix first so the API's sort rewrites the order.
func (drt *DomainResourceTest) BuildRouteOrderUnsortedMultiTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, workloadSelfLink := drt.newRouteOrderCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.RouteOrderUnsortedMultiHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.#", "5"),

			// Index assertions are deliberate: the API returns these routes sorted by descending prefix length, state must keep the declared order
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.prefix", "/"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.prefix", "/api"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.prefix", "/api/users"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.port", "8080"),

			// These two routes share a prefix and differ only by host, which is what the API's own route uniqueness rule allows
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.prefix", "/a"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.host_prefix", "www"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.prefix", "/a"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.host_prefix", "api-staging"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.port", "8080"),

			// Prove the API really did rewrite this order, so the assertions above keep testing the reordering
			domainApiRouteOrderDiffersCheck(subDomainName, []string{"/", "/api", "/api/users", "/a|www", "/a|api-staging"}),
		),
	}
}

// BuildRouteOrderWithExternalTestStep returns a step where an external cpln_domain_route lands in the middle of the API's sorted order.
func (drt *DomainResourceTest) BuildRouteOrderWithExternalTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, workloadSelfLink := drt.newRouteOrderCases(subDomainName)

	// Build the external route case that owns a route outside the domain's inline blocks
	externalRoute := DomainRouteResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "external-route",
			ResourceAddress: "cpln_domain_route.external-route",
		},
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.RouteOrderWithExternalHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			// The inline blocks keep their declared order and exclude the externally owned route
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.#", "5"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.prefix", "/"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.prefix", "/api"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.prefix", "/api/users"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.prefix", "/a"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.host_prefix", "www"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.prefix", "/a"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.host_prefix", "api-staging"),

			// The externally owned route survives the domain update that reorders everything around it
			externalRoute.TestCheckResourceAttr("domain_link", subDomain.GetSelfLink()),
			externalRoute.TestCheckResourceAttr("domain_port", "443"),
			externalRoute.TestCheckResourceAttr("prefix", "/external"),
			externalRoute.TestCheckResourceAttr("workload_link", workloadSelfLink),

			// The API holds all six routes and sorted the external one into the middle of the inline ones
			domainApiRouteOrderDiffersCheck(subDomainName, []string{"/", "/api", "/api/users", "/a|www", "/a|api-staging", "/external"}),
		),
	}
}

// BuildRouteOrderExpandedTestStep returns a step with six routes declared in another order the API's sort rewrites.
func (drt *DomainResourceTest) BuildRouteOrderExpandedTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, workloadSelfLink := drt.newRouteOrderCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.RouteOrderExpandedHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.#", "6"),

			// The declared order differs from both the previous step and the API's sorted order
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.prefix", "/a"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.host_prefix", "www"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.prefix", "/api"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.prefix", "/"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.prefix", "/api/users"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.prefix", "/a"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.host_prefix", "api-staging"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.4.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.5.prefix", "/api/v1/orders"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.5.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.5.port", "8080"),

			// Prove the API really did rewrite this order, so the assertions above keep testing the reordering
			domainApiRouteOrderDiffersCheck(subDomainName, []string{"/a|www", "/api", "/", "/api/users", "/a|api-staging", "/api/v1/orders"}),
		),
	}
}

// BuildRouteOrderRegexTestStep returns a step where a regex route switches the API's route sort off and the declared order must survive untouched.
func (drt *DomainResourceTest) BuildRouteOrderRegexTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, workloadSelfLink := drt.newRouteOrderCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.RouteOrderRegexHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.#", "4"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.prefix", "/"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.prefix", "/api"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.1.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.prefix", "/api/users"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.2.port", "8080"),

			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.regex", "/health.*"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.workload_link", workloadSelfLink),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.3.port", "8080"),
		),
	}
}

// newHttpsPortWithoutTlsCases builds the subdomain test case and workload link shared across the https port tls stages.
func (drt *DomainResourceTest) newHttpsPortWithoutTlsCases(subDomainName string) (DomainResourceTestCase, string) {
	// Build the subdomain case that owns the https port
	subDomain := DomainResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "domain",
			ResourceName:    "subdomain",
			ResourceAddress: "cpln_domain.subdomain",
			Name:            subDomainName,
		},
	}

	// Construct the workload self link in the short form the config uses
	workloadSelfLink := fmt.Sprintf("//gvc/gvc-%s/workload/workload-%s", drt.RandomName, drt.RandomName)

	// Return the shared case and link
	return subDomain, workloadSelfLink
}

// BuildHttpsPortWithoutTlsTestStep returns a step where the https port omits the tls block the API fills in server side.
func (drt *DomainResourceTest) BuildHttpsPortWithoutTlsTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, workloadSelfLink := drt.newHttpsPortWithoutTlsCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HttpsPortWithoutTlsHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("name", subDomainName),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.number", "443"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.protocol", "http"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.prefix", "/"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.route.0.workload_link", workloadSelfLink),
			// The API answers an https port with a default tls configuration which must stay out of state
			subDomain.TestCheckResourceAttr("spec.0.ports.0.tls.#", "0"),
		),
	}
}

// BuildHttpsPortWithoutTlsPlanStableTestStep returns a plan-only step that re-runs the no-tls config to verify no plan drift.
func (drt *DomainResourceTest) BuildHttpsPortWithoutTlsPlanStableTestStep(subDomainName string) resource.TestStep {
	// Initialize and return the test step
	return resource.TestStep{
		Config:   drt.HttpsPortWithoutTlsHcl(subDomainName),
		PlanOnly: true,
	}
}

// BuildHttpsPortMixedTlsTestStep returns a step where an http2 port 443 omits tls while port 80 sets one.
func (drt *DomainResourceTest) BuildHttpsPortMixedTlsTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, _ := drt.newHttpsPortWithoutTlsCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HttpsPortMixedTlsHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("spec.0.ports.#", "2"),

			// Port 443 omits tls on http2 as well as http and keeps the API default out of state
			subDomain.TestCheckResourceAttr("spec.0.ports.0.number", "443"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.protocol", "http2"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.tls.#", "0"),

			// Port 80 sets tls and keeps exactly what it declared
			subDomain.TestCheckResourceAttr("spec.0.ports.1.number", "80"),
			subDomain.TestCheckResourceAttr("spec.0.ports.1.tls.0.min_protocol_version", "TLSV1_2"),
			subDomain.TestCheckSetAttr("spec.0.ports.1.tls.0.cipher_suites", []string{
				"ECDHE-ECDSA-AES256-GCM-SHA384",
				"ECDHE-RSA-AES256-GCM-SHA384",
			}),
		),
	}
}

// BuildHttpsPortExplicitTlsTestStep returns a step where the https port declares a tls block with non-default values.
func (drt *DomainResourceTest) BuildHttpsPortExplicitTlsTestStep(subDomainName string) resource.TestStep {
	// Resolve the shared case and workload link
	subDomain, _ := drt.newHttpsPortWithoutTlsCases(subDomainName)

	// Initialize and return the test step
	return resource.TestStep{
		Config: drt.HttpsPortExplicitTlsHcl(subDomainName),
		Check: resource.ComposeAggregateTestCheckFunc(
			subDomain.TestCheckResourceAttr("spec.0.ports.#", "1"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.number", "443"),
			subDomain.TestCheckResourceAttr("spec.0.ports.0.tls.0.min_protocol_version", "TLSV1_1"),
			subDomain.TestCheckSetAttr("spec.0.ports.0.tls.0.cipher_suites", []string{
				"AES256-GCM-SHA384",
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
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  replica       = 1
}

resource "cpln_domain_route" "second-route" {
  domain_link = cpln_domain.subdomain.self_link
  domain_port = 80

  prefix         = "/second"
  replace_prefix = "/"
  workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
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

  mirror {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    port          = 8080
    percent       = 50
  }

  mirror {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    percent       = 25.5
  }

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    port          = 8080
    weight        = 50
  }

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    weight        = 25
  }
}

resource "cpln_domain_route" "third-route" {
  domain_link = cpln_domain.subdomain.self_link
  domain_port = 80

  prefix         = "/third"
  replace_prefix = "/"
  workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
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
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
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
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
}

resource "cpln_domain_route" "second-route" {
  domain_link = cpln_domain.subdomain.self_link
  domain_port = 80

  prefix         = "/second"
  replace_prefix = "/"
  workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
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

  mirror {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    port          = 8080
    percent       = 50
  }

  mirror {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    percent       = 25.5
  }

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    port          = 8080
    weight        = 50
  }

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    weight        = 25
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
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix         = "/app"
        replace_prefix = "/"
        workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port           = 8080

        mirror {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          port          = 8080
          percent       = 50
        }

        mirror {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          percent       = 25.5
        }

        canary {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          port          = 8080
          weight        = 50
        }

        canary {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          weight        = 25
        }
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
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080

        mirror {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          port          = 8080
          percent       = 75
        }

        canary {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          port          = 8080
          weight        = 75
        }
      }

      route {
        prefix         = "/app"
        replace_prefix = "/"
        workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port           = 8080

        headers {
          request {
            set = {
              "X-Forwarded-Proto" = "https"
            }
          }
        }

        mirror {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          percent       = 30
        }

        canary {
          workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
          weight        = 30
        }
      }

      route {
        regex         = "/user/.*/profile"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
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

// hclSubDomainTcpNoTls returns HCL for a CNAME subdomain that exposes a TCP port and intentionally omits the tls block.
func (drt *DomainResourceTest) hclSubDomainTcpNoTls(name string) string {
	return fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  name = "%s"

  tags = {
    terraform_generated = "true"
  }

  spec {
    dns_mode = "cname"

    ports {
      number   = 5432
      protocol = "tcp"

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }
    }
  }
}
`, name)
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
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = %d
}
`, resourceName, routeKey, port)
}

// CanaryAbsentHcl returns HCL for a route on the subdomain with no canary blocks.
func (drt *DomainResourceTest) CanaryAbsentHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "canary lifecycle") + `
resource "cpln_domain_route" "canary-route" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/canary"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
}
`
}

// CanaryRequiredOnlyHcl returns HCL for a route with a single canary that sets only its required attributes.
func (drt *DomainResourceTest) CanaryRequiredOnlyHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "canary lifecycle") + `
resource "cpln_domain_route" "canary-route" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/canary"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    weight        = 50
  }
}
`
}

// CanaryAllMultiHcl returns HCL for a route with two canaries covering every attribute and a weight 0 toggle.
func (drt *DomainResourceTest) CanaryAllMultiHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "canary lifecycle") + `
resource "cpln_domain_route" "canary-route" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/canary"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    port          = 8080
    weight        = 60
  }

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    weight        = 0
  }
}
`
}

// CanaryExpandedHcl returns HCL for a route with three canaries so the build/flatten loop runs at a non-trivial count.
func (drt *DomainResourceTest) CanaryExpandedHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "canary lifecycle") + `
resource "cpln_domain_route" "canary-route" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/canary"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    port          = 8080
    weight        = 30
  }

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    weight        = 20
  }

  canary {
    workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
    port          = 8080
    weight        = 0
  }
}
`
}

// hclHostRouteWorkload returns HCL for the second workload the host routes point at.
func (drt *DomainResourceTest) hclHostRouteWorkload() string {
	return `
resource "cpln_workload" "alternate" {
  gvc         = cpln_gvc.new.name
  name        = "workload-alt-${var.random_name}"
  description = "Host route test workload"
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
`
}

// HostRouteCatchAllHcl returns HCL for a single route that owns the path and matches every host.
func (drt *DomainResourceTest) HostRouteCatchAllHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "host route lifecycle") + drt.hclHostRouteWorkload() + `
resource "cpln_domain_route" "catch-all" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = 8080
}
`
}

// HostRouteMultiHostHcl returns HCL for three routes that share the path, two of them narrowed to a host.
func (drt *DomainResourceTest) HostRouteMultiHostHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "host route lifecycle") + drt.hclHostRouteWorkload() + `
resource "cpln_domain_route" "catch-all" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = 8080
}

resource "cpln_domain_route" "store" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  host_prefix   = "store."
  workload_link = "//gvc/${cpln_workload.alternate.gvc}/workload/${cpln_workload.alternate.name}"
  port          = 8080
}

resource "cpln_domain_route" "blog" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  host_prefix   = "blog."
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = 8080
}
`
}

// HostRouteSwappedHcl returns HCL where the two host routes exchange the workloads they point at.
func (drt *DomainResourceTest) HostRouteSwappedHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "host route lifecycle") + drt.hclHostRouteWorkload() + `
resource "cpln_domain_route" "catch-all" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = 8080
}

resource "cpln_domain_route" "store" {
  domain_link    = cpln_domain.subdomain.self_link
  domain_port    = 443
  prefix         = "/"
  replace_prefix = "/shop"
  host_prefix    = "store."
  workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port           = 8080
}

resource "cpln_domain_route" "blog" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  host_prefix   = "blog."
  workload_link = "//gvc/${cpln_workload.alternate.gvc}/workload/${cpln_workload.alternate.name}"
  port          = 8080
}
`
}

// HostRouteRenamedHostHcl returns HCL where one route moves from its host to another one.
func (drt *DomainResourceTest) HostRouteRenamedHostHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "host route lifecycle") + drt.hclHostRouteWorkload() + `
resource "cpln_domain_route" "catch-all" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = 8080
}

resource "cpln_domain_route" "store" {
  domain_link    = cpln_domain.subdomain.self_link
  domain_port    = 443
  prefix         = "/"
  replace_prefix = "/shop"
  host_prefix    = "store."
  workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port           = 8080
}

resource "cpln_domain_route" "blog" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  host_prefix   = "news."
  workload_link = "//gvc/${cpln_workload.alternate.gvc}/workload/${cpln_workload.alternate.name}"
  port          = 8080
}
`
}

// HostRouteExpandedHcl returns HCL where a fourth route takes the same path on a host regex.
func (drt *DomainResourceTest) HostRouteExpandedHcl(subDomainName string) string {
	return drt.hclBase() + drt.hclSubDomain(subDomainName, "host route lifecycle") + drt.hclHostRouteWorkload() + `
resource "cpln_domain_route" "catch-all" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = 8080
}

resource "cpln_domain_route" "store" {
  domain_link    = cpln_domain.subdomain.self_link
  domain_port    = 443
  prefix         = "/"
  replace_prefix = "/shop"
  host_prefix    = "store."
  workload_link  = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port           = 8080
}

resource "cpln_domain_route" "blog" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  host_prefix   = "news."
  workload_link = "//gvc/${cpln_workload.alternate.gvc}/workload/${cpln_workload.alternate.name}"
  port          = 8080
}

resource "cpln_domain_route" "api" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/"
  host_regex    = "^api.*$"
  workload_link = "//gvc/${cpln_workload.alternate.gvc}/workload/${cpln_workload.alternate.name}"
  port          = 8080
}
`
}

// RouteOrderAbsentHcl returns HCL for a subdomain whose https port declares no inline route blocks.
func (drt *DomainResourceTest) RouteOrderAbsentHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

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
`, subDomainName)
}

// RouteOrderRequiredOnlyHcl returns HCL for a subdomain with a single inline route that sets only its required attributes.
func (drt *DomainResourceTest) RouteOrderRequiredOnlyHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

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

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
      }
    }
  }
}
`, subDomainName)
}

// RouteOrderUnsortedMultiHcl returns HCL for five inline routes declared shortest prefix first, the order the API's route sort rewrites.
func (drt *DomainResourceTest) RouteOrderUnsortedMultiHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

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

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/api"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/api/users"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/a"
        host_prefix   = "www"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/a"
        host_prefix   = "api-staging"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }
    }
  }
}
`, subDomainName)
}

// RouteOrderWithExternalHcl returns the unsorted multi route config plus an externally owned route on the same port.
func (drt *DomainResourceTest) RouteOrderWithExternalHcl(subDomainName string) string {
	return drt.RouteOrderUnsortedMultiHcl(subDomainName) + `
resource "cpln_domain_route" "external-route" {
  domain_link   = cpln_domain.subdomain.self_link
  domain_port   = 443
  prefix        = "/external"
  workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
  port          = 8080
}
`
}

// RouteOrderExpandedHcl returns HCL for six inline routes declared in another order the API's route sort rewrites.
func (drt *DomainResourceTest) RouteOrderExpandedHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

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

      route {
        prefix        = "/a"
        host_prefix   = "www"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/api"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/api/users"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/a"
        host_prefix   = "api-staging"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/api/v1/orders"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }
    }
  }
}
`, subDomainName)
}

// RouteOrderRegexHcl returns HCL for four inline routes where a regex route switches the API's route sort off.
func (drt *DomainResourceTest) RouteOrderRegexHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

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

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/api"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        prefix        = "/api/users"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }

      route {
        regex         = "/health.*"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }
    }
  }
}
`, subDomainName)
}

// HttpsPortWithoutTlsHcl returns HCL for a subdomain whose https port omits the optional tls block.
func (drt *DomainResourceTest) HttpsPortWithoutTlsHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

  tags = {
    terraform_generated = "true"
  }

  spec {
    dns_mode         = "ns"
    accept_all_hosts = true

    ports {
      number   = 443
      protocol = "http"

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }
    }
  }
}
`, subDomainName)
}

// HttpsPortMixedTlsHcl returns HCL for a subdomain where an http2 port 443 omits tls and a second port declares one.
func (drt *DomainResourceTest) HttpsPortMixedTlsHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

  tags = {
    terraform_generated = "true"
  }

  spec {
    dns_mode         = "ns"
    accept_all_hosts = true

    ports {
      number   = 443
      protocol = "http2"

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }
    }

    ports {
      number   = 80
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
`, subDomainName)
}

// HttpsPortExplicitTlsHcl returns HCL for a subdomain whose https port declares a tls block with non-default values.
func (drt *DomainResourceTest) HttpsPortExplicitTlsHcl(subDomainName string) string {
	return drt.hclBase() + fmt.Sprintf(`
resource "cpln_domain" "subdomain" {
  depends_on = [cpln_domain.new]
  name       = "%s"

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
        min_protocol_version = "TLSV1_1"
        cipher_suites = [
          "AES256-GCM-SHA384",
        ]
      }

      route {
        prefix        = "/"
        workload_link = "//gvc/${cpln_workload.new.gvc}/workload/${cpln_workload.new.name}"
        port          = 8080
      }
    }
  }
}
`, subDomainName)
}

// domainApiRouteOrderDiffersCheck verifies the API stored the first port's routes in an order other than the declared one.
func domainApiRouteOrderDiffersCheck(domainName string, declaredOrder []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// Read the domain straight from the API to inspect the order it stored
		domain, _, err := TestProvider.client.GetDomain(domainName)

		if err != nil {
			return fmt.Errorf("error occurred while fetching domain %s: %w", domainName, err)
		}

		// Reject a domain that carries no ports to inspect
		if domain.Spec == nil || domain.Spec.Ports == nil || len(*domain.Spec.Ports) == 0 {
			return fmt.Errorf("expected domain %s to expose at least one port", domainName)
		}

		// Collect the route identities the API stored for the first port
		apiOrder := []string{}
		port := (*domain.Spec.Ports)[0]

		if port.Routes != nil {
			for _, route := range *port.Routes {
				// Identify the route the same way the declared order spells it
				identity := ""

				if route.Prefix != nil {
					identity = *route.Prefix
				} else if route.Regex != nil {
					identity = *route.Regex
				}

				if route.HostPrefix != nil {
					identity += "|" + *route.HostPrefix
				}

				apiOrder = append(apiOrder, identity)
			}
		}

		// The declared routes must all be present, otherwise the comparison below is meaningless
		if len(apiOrder) != len(declaredOrder) {
			return fmt.Errorf(
				"expected domain %s to store %d routes on its first port, got %d: %v",
				domainName, len(declaredOrder), len(apiOrder), apiOrder,
			)
		}

		// This step only exercises the reordering while the API keeps rewriting the declared order
		if strings.Join(apiOrder, ",") == strings.Join(declaredOrder, ",") {
			return fmt.Errorf(
				"expected the API to store domain %s routes in an order other than the declared one, got %v for both",
				domainName, apiOrder,
			)
		}

		return nil
	}
}

// domainRouteHostMappingCheck returns a TestCheckFunc that verifies which workload every host on a domain port routes to.
func domainRouteHostMappingCheck(domainName string, domainPort int, expected map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// Read the domain straight from the API to inspect the routes it stored
		domain, _, err := TestProvider.client.GetDomain(domainName)

		if err != nil {
			return fmt.Errorf("error occurred while fetching domain %s: %w", domainName, err)
		}

		// Reject a domain that carries no ports to inspect
		if domain.Spec == nil || domain.Spec.Ports == nil {
			return fmt.Errorf("expected domain %s to expose at least one port", domainName)
		}

		// Collect the workload every host on the port routes to
		actual := map[string]string{}

		for _, port := range *domain.Spec.Ports {
			// Skip ports that carry another number or no routes at all
			if port.Number == nil || *port.Number != domainPort || port.Routes == nil {
				continue
			}

			for _, route := range *port.Routes {
				// Reject a route the API stored without a workload to reach
				if route.WorkloadLink == nil {
					return fmt.Errorf("expected every route on domain %s port %d to carry a workload link", domainName, domainPort)
				}

				// Key the route by the host it matches on, which is empty when it matches every host
				host := ""

				if route.HostPrefix != nil {
					host = *route.HostPrefix
				} else if route.HostRegex != nil {
					host = *route.HostRegex
				}

				actual[host] = GetNameFromSelfLink(*route.WorkloadLink)
			}
		}

		// Compare the mapping the API stored against the one the configuration declares
		if !reflect.DeepEqual(actual, expected) {
			return fmt.Errorf("expected domain %s port %d to route %v, got %v", domainName, domainPort, expected, actual)
		}

		return nil
	}
}

// domainImportWithRoutesCheck returns an ImportStateCheckFunc that verifies routes WERE imported into state.
func domainImportWithRoutesCheck() resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		// Validate exactly one resource was imported
		if len(states) != 1 {
			return fmt.Errorf("expected 1 imported state, got %d", len(states))
		}

		// Check that routes are stored in the spec (backward compat for standard import)
		routeCount, exists := states[0].Attributes["spec.0.ports.0.route.#"]
		if !exists || routeCount == "0" {
			return fmt.Errorf(
				"expected routes in imported state (spec.0.ports.0.route.# should be > 0), got exists=%v count=%q",
				exists, routeCount,
			)
		}

		return nil
	}
}

// domainRouteLinkCheck returns an ImportStateCheckFunc that verifies the imported domain_link matches the expected self-link.
func domainRouteLinkCheck(expectedLink string) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		// Validate exactly one resource was imported
		if len(states) != 1 {
			return fmt.Errorf("expected 1 imported state, got %d", len(states))
		}

		// Verify domain_link is the full self-link, not a bare domain name
		domainLink := states[0].Attributes["domain_link"]
		if domainLink != expectedLink {
			return fmt.Errorf(
				"expected domain_link to be %q, got %q",
				expectedLink, domainLink,
			)
		}

		return nil
	}
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
