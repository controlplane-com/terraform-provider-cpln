package cpln

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlanePolicy_basic performs an acceptance test for the resource.
func TestAccControlPlanePolicy_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewPolicyResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "POLICY") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// PolicyResourceTest defines the necessary functionality to test the resource.
type PolicyResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
}

// NewPolicyResourceTest creates a PolicyResourceTest with initialized test cases.
func NewPolicyResourceTest() PolicyResourceTest {
	// Create a resource test instance
	resourceTest := PolicyResourceTest{
		RandomName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
	}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewTargetSecretScenario()...)
	steps = append(steps, resourceTest.NewTargetWorkloadScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (prt *PolicyResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_policy resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_policy" {
			continue
		}

		// Retrieve the name for the current resource
		policyName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of policy with name: %s", policyName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		policy, code, err := TestProvider.client.GetPolicy(policyName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if policy %s exists: %w", policyName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if policy != nil {
			return fmt.Errorf("CheckDestroy failed: policy %s still exists in the system", *policy.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_policy resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewTargetSecretScenario defines a policy test scenario targeting secrets with create, import, and update steps.
func (prt *PolicyResourceTest) NewTargetSecretScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "new"
	name := fmt.Sprintf("tf-policy-secret-%s", prt.RandomName)

	// Build test steps
	initialConfig, initialStep := prt.BuildTargetSecretTestStep(resourceName, name)
	caseUpdate1 := prt.BuildTargetSecretUpdate1TestStep(initialConfig.ProviderTestCase)

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
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewTargetWorkloadScenario defines a policy test scenario targeting workloads with create, import, and update steps.
func (prt *PolicyResourceTest) NewTargetWorkloadScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "new"
	name := fmt.Sprintf("tf-policy-workload-%s", prt.RandomName)

	// Build test steps
	initialConfig, initialStep := prt.BuildTargetWorkloadTestStep(resourceName, name)
	caseUpdate1 := prt.BuildTargetWorkloadUpdate1TestStep(initialConfig.ProviderTestCase)

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
		// Revert the resource to its initial state
		initialStep,
	}
}

// Test Cases //

// BuildTargetSecretTestStep constructs the initial test step and case for a secret-targeting policy.
func (prt *PolicyResourceTest) BuildTargetSecretTestStep(resourceName string, name string) (PolicyResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := PolicyResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "policy",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_policy.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "policy secret new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: prt.TargetSecretRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckResourceAttr("target_kind", "secret"),
			c.TestCheckResourceAttr("target", "all"),
		),
	}
}

// BuildTargetSecretUpdate1TestStep constructs the update test step for a secret-targeting policy including tags and bindings.
func (prt *PolicyResourceTest) BuildTargetSecretUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := PolicyResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: prt.TargetSecretUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckResourceAttr("target_kind", "secret"),
			c.TestCheckSetAttr("target_links", []string{"/org/terraform-test-org/secret/secret-01", "/org/terraform-test-org/secret/secret-02"}),
			c.TestCheckNestedBlocks("target_query", []map[string]interface{}{
				{
					"fetch": "items",
					"spec": []map[string]interface{}{
						{
							"match": "all",
							"terms": []map[string]interface{}{
								{
									"op":    "=",
									"tag":   "terraform_generated",
									"value": "true",
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("binding", []map[string]interface{}{
				{
					"permissions":     []string{"manage", "edit"},
					"principal_links": []string{"user/support@controlplane.com", "group/viewers", "serviceaccount/service-account-01", "gvc/gvc-01/identity/identity-01"},
				},
			}),
		),
	}
}

// BuildTargetWorkloadTestStep constructs the initial test step and case for a workload-targeting policy.
func (prt *PolicyResourceTest) BuildTargetWorkloadTestStep(resourceName string, name string) (PolicyResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := PolicyResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "policy",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_policy.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "policy workload new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: prt.TargetWorkloadRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckResourceAttr("gvc", "gvc-01"),
			c.TestCheckResourceAttr("target_kind", "workload"),
			c.TestCheckResourceAttr("target", "all"),
		),
	}
}

// BuildTargetWorkloadUpdate1TestStep constructs the update test step for a workload-targeting policy including bindings.
func (prt *PolicyResourceTest) BuildTargetWorkloadUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := PolicyResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: prt.TargetWorkloadUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckResourceAttr("gvc", "gvc-01"),
			c.TestCheckResourceAttr("target_kind", "workload"),
			c.TestCheckSetAttr("target_links", []string{"workload-01"}),
			c.TestCheckNestedBlocks("binding", []map[string]interface{}{
				{
					"permissions":     []string{"manage", "edit"},
					"principal_links": []string{"user/support@controlplane.com", "group/viewers", "serviceaccount/service-account-01", "gvc/gvc-01/identity/identity-01"},
				},
				{
					"permissions":     []string{"manage", "edit", "delete"},
					"principal_links": []string{"/org/terraform-test-org/group/superusers", "serviceaccount/service-account-01"},
				},
			}),
		),
	}
}

// Configs //

// TargetSecretRequiredOnlyHcl returns a minimal HCL configuration for a policy targeting all secrets.
func (prt *PolicyResourceTest) TargetSecretRequiredOnlyHcl(c PolicyResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_policy" "%s" {
  name        = "%s"
  target_kind = "secret"
  target      = "all"
}
`, c.ResourceName, c.Name)
}

// TargetSecretUpdate1Hcl returns an HCL configuration for updating a secret-targeting policy with description, tags, queries, and bindings.
func (prt *PolicyResourceTest) TargetSecretUpdate1Hcl(c PolicyResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_policy" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  target_kind  = "secret"
  target_links = ["/org/terraform-test-org/secret/secret-01", "/org/terraform-test-org/secret/secret-02"]

  target_query {
    spec {
      # match is either "all", "any", or "none"
      match = "all"

      terms {
        op    = "="
        tag   = "terraform_generated"
        value = "true"
      }
    }
  }

  binding {
    permissions     = ["manage", "edit"]
    principal_links = ["user/support@controlplane.com", "group/viewers", "serviceaccount/service-account-01","gvc/gvc-01/identity/identity-01"]
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// TargetWorkloadRequiredOnlyHcl returns a minimal HCL configuration for a policy targeting all workloads.
func (prt *PolicyResourceTest) TargetWorkloadRequiredOnlyHcl(c PolicyResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_policy" "%s" {
  name        = "%s"
	gvc         = "gvc-01"
  target_kind = "workload"
  target      = "all"
}
`, c.ResourceName, c.Name)
}

// TargetWorkloadUpdate1Hcl constructs an HCL block to update a workload-targeting policy with metadata, tags, and binding rules.
func (prt *PolicyResourceTest) TargetWorkloadUpdate1Hcl(c PolicyResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_policy" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  gvc          = "gvc-01"
  target_kind  = "workload"
  target_links = ["workload-01"]

  binding {
    permissions     = ["manage", "edit"]
    principal_links = ["user/support@controlplane.com", "group/viewers", "serviceaccount/service-account-01","gvc/gvc-01/identity/identity-01"]
  }

  binding {
    permissions     = ["manage", "edit", "delete"]
    principal_links = ["/org/terraform-test-org/group/superusers", "serviceaccount/service-account-01"]
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

/*** Resource Test Case ***/

// PolicyResourceTestCase defines a specific resource test case.
type PolicyResourceTestCase struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (prtc *PolicyResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of policy: %s. Total resources: %d", prtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[prtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", prtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != prtc.Name {
			return fmt.Errorf("resource ID %s does not match expected policy name %s", rs.Primary.ID, prtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remotePolicy, _, err := TestProvider.client.GetPolicy(prtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving policy from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remotePolicy.Name != prtc.Name {
			return fmt.Errorf("mismatch in policy name: expected %s, got %s", prtc.Name, *remotePolicy.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("policy %s verified successfully in both state and external system.", prtc.Name))
		return nil
	}
}
