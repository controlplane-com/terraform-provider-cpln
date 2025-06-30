package cpln

import (
	"errors"
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

/*** Acceptance Test ***/

// TestAccControlPlaneGroup_basic performs an acceptance test for the resource.
func TestAccControlPlaneGroup_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewGroupResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "GROUP") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// GroupResourceTest defines the necessary functionality to test the resource.
type GroupResourceTest struct {
	Steps []resource.TestStep
}

// NewGroupResourceTest creates a GroupResourceTest with initialized test cases.
func NewGroupResourceTest() GroupResourceTest {
	// Create a resource test instance
	resourceTest := GroupResourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewDefaultScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (grt *GroupResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_group resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_group" {
			continue
		}

		// Retrieve the name for the current resource
		groupName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of group with name: %s", groupName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		group, code, err := TestProvider.client.GetGroup(groupName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if group %s exists: %w", groupName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if group != nil {
			return fmt.Errorf("CheckDestroy failed: group %s still exists in the system", *group.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_group resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case for a group using JMESPATH with initial and updated configurations.
func (grt *GroupResourceTest) NewDefaultScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("sa-default-%s", random)
	resourceName := "new"
	serviceAccountName := fmt.Sprintf("sa-%s", random)

	// Create the service account case
	serviceAccountCase := ServiceAccountResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "serviceaccount",
			ResourceName:      "new",
			ResourceAddress:   "cpln_service_account.new",
			Name:              serviceAccountName,
			Description:       serviceAccountName,
			DescriptionUpdate: "service account default description updated",
		},
	}

	// Create a service account resource test instance
	serviceAccountResourceTest := ServiceAccountResourceTest{}

	// Initialize the service account config
	serviceAccountConfig := serviceAccountResourceTest.RequiredOnly(serviceAccountCase)

	// Build test steps
	initialConfig, initialStep := grt.BuildInitialTestStep(resourceName, name)
	caseUpdate1 := grt.BuildUpdate1TestStep(initialConfig.ProviderTestCase, resourceName, serviceAccountConfig, serviceAccountCase)
	caseUpdate2 := grt.BuildUpdate2TestStep(initialConfig.ProviderTestCase, resourceName, serviceAccountConfig, serviceAccountCase)

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
		// Revert the resource to its initial state
		// initialStep, // TODO: Uncomment this once the issue with memberQuery and identityMatcher removal is fixed.
	}
}

// Test Cases //

// BuildInitialTestStep returns a default initial test step and its associated test case for the resource.
func (grt *GroupResourceTest) BuildInitialTestStep(resourceName string, name string) (GroupResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := GroupResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "group",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_group.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "group default description updated",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: grt.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (grt *GroupResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase, resourceName string, serviceAccountConfig string, serviceAccountCase ServiceAccountResourceTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := GroupResourceTestCase{
		ProviderTestCase: initialCase,
		UserIdsAndEmails: []string{"unittest@controlplane.com"},
		MemberQuery: client.Query{
			Fetch: StringPointer("items"),
			Spec: &client.QuerySpec{
				Match: StringPointer("all"),
				Terms: &[]client.QueryTerm{
					{
						Op:    StringPointer("="),
						Tag:   StringPointer("firebase/sign_in_provider"),
						Value: StringPointer("microsoft.com"),
					},
				},
			},
		},
		IdentityMatcher: client.GroupIdentityMatcher{
			Expression: StringPointer("groups"),
			Language:   StringPointer("jmespath"),
		},
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: grt.UpdateWithMinimalOptionals(c, serviceAccountConfig, serviceAccountCase),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckSetAttr("user_ids_and_emails", c.UserIdsAndEmails),
			c.TestCheckSetAttr("service_accounts", []string{serviceAccountCase.Name}),
			c.TestCheckNestedBlocks("member_query", []map[string]interface{}{
				{
					"fetch": *c.MemberQuery.Fetch,
					"spec": []map[string]interface{}{
						{
							"match": *c.MemberQuery.Spec.Match,
							"terms": []map[string]interface{}{
								{
									"op":    *(*c.MemberQuery.Spec.Terms)[0].Op,
									"tag":   *(*c.MemberQuery.Spec.Terms)[0].Tag,
									"value": *(*c.MemberQuery.Spec.Terms)[0].Value,
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("identity_matcher", []map[string]interface{}{
				{
					"expression": *c.IdentityMatcher.Expression,
					"language":   *c.IdentityMatcher.Language,
				},
			}),
		),
	}
}

// BuildUpdate2TestStep returns a test step for the update.
func (grt *GroupResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase, resourceName string, serviceAccountConfig string, serviceAccountCase ServiceAccountResourceTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := GroupResourceTestCase{
		ProviderTestCase: initialCase,
		UserIdsAndEmails: []string{"unittest@controlplane.com", "unittest2@controlplane.com"},
		MemberQuery: client.Query{
			Fetch: StringPointer("items"),
			Spec: &client.QuerySpec{
				Match: StringPointer("all"),
				Terms: &[]client.QueryTerm{
					{
						Op:    StringPointer("="),
						Tag:   StringPointer("firebase/sign_in_provider"),
						Value: StringPointer("microsoft.com"),
					},
				},
			},
		},
		IdentityMatcher: client.GroupIdentityMatcher{
			Expression: StringPointer("if ($.includes('groups')) { const y = $.groups; }"),
			Language:   StringPointer("javascript"),
		},
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: grt.UpdateWithAllOptionals(c, serviceAccountConfig, serviceAccountCase),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			c.TestCheckSetAttr("user_ids_and_emails", c.UserIdsAndEmails),
			c.TestCheckSetAttr("service_accounts", []string{serviceAccountCase.Name}),
			c.TestCheckNestedBlocks("member_query", []map[string]interface{}{
				{
					"fetch": *c.MemberQuery.Fetch,
					"spec": []map[string]interface{}{
						{
							"match": *c.MemberQuery.Spec.Match,
							"terms": []map[string]interface{}{
								{
									"op":    *(*c.MemberQuery.Spec.Terms)[0].Op,
									"tag":   *(*c.MemberQuery.Spec.Terms)[0].Tag,
									"value": *(*c.MemberQuery.Spec.Terms)[0].Value,
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("identity_matcher", []map[string]interface{}{
				{
					"expression": *c.IdentityMatcher.Expression,
					"language":   *c.IdentityMatcher.Language,
				},
			}),
		),
	}
}

// Configs //

// RequiredOnly returns a minimal HCL block for a resource using only required fields.
func (grt *GroupResourceTest) RequiredOnly(c GroupResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_group" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

// UpdateWithMinimalOptionals constructs an HCL configuration for a group resource using minimal optional fields including member query and identity matcher
func (grt *GroupResourceTest) UpdateWithMinimalOptionals(c GroupResourceTestCase, serviceAccountConfig string, serviceAccountCase ServiceAccountResourceTestCase) string {
	return fmt.Sprintf(`
# Service Account Resource
%s

resource "cpln_group" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  user_ids_and_emails = %s
  service_accounts    = [%s]

  member_query {
    spec {
      terms {
        tag   = "%s"
        value = "%s"
      }
    }
  }

  identity_matcher {
    expression = "%s"
    # language default value is 'jmespath'
	}
}
`, serviceAccountConfig, c.ResourceName, serviceAccountCase.ResourceAddress, c.Name, c.DescriptionUpdate, StringSliceToString(c.UserIdsAndEmails),
		serviceAccountCase.GetResourceNameAttr(), *(*c.MemberQuery.Spec.Terms)[0].Tag, *(*c.MemberQuery.Spec.Terms)[0].Value, *c.IdentityMatcher.Expression,
	)
}

// UpdateWithAllOptionals constructs an HCL configuration for a group resource using all optional fields including fetch, match, and full identity matcher settings.
func (grt *GroupResourceTest) UpdateWithAllOptionals(c GroupResourceTestCase, serviceAccountConfig string, serviceAccountCase ServiceAccountResourceTestCase) string {
	return fmt.Sprintf(`
# Service Account Resource
%s

resource "cpln_group" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  user_ids_and_emails = %s
  service_accounts    = [%s]

  member_query {
    fetch = "%s"

    spec {
      match = "%s"

      terms {
        op    = "%s"
        tag   = "%s"
        value = "%s"
      }
    }
  }

  identity_matcher {
    expression = "%s"
    language   = "%s"
  }
}
`, serviceAccountConfig, c.ResourceName, serviceAccountCase.ResourceAddress, c.Name, c.DescriptionUpdate, StringSliceToString(c.UserIdsAndEmails), serviceAccountCase.GetResourceNameAttr(),
		*c.MemberQuery.Fetch, *c.MemberQuery.Spec.Match, *(*c.MemberQuery.Spec.Terms)[0].Op, *(*c.MemberQuery.Spec.Terms)[0].Tag, *(*c.MemberQuery.Spec.Terms)[0].Value,
		*c.IdentityMatcher.Expression, *c.IdentityMatcher.Language,
	)
}

/*** Resource Test Case ***/

// GroupResourceTestCase defines a specific resource test case.
type GroupResourceTestCase struct {
	ProviderTestCase
	UserIdsAndEmails []string
	MemberQuery      client.Query
	IdentityMatcher  client.GroupIdentityMatcher
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (grtc *GroupResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of group: %s. Total resources: %d", grtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[grtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", grtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != grtc.Name {
			return fmt.Errorf("resource ID %s does not match expected group name %s", rs.Primary.ID, grtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteGroup, _, err := TestProvider.client.GetGroup(grtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving group from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteGroup.Name != grtc.Name {
			return fmt.Errorf("mismatch in group name: expected %s, got %s", grtc.Name, *remoteGroup.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("Group %s verified successfully in both state and external system.", grtc.Name))
		return nil
	}
}
