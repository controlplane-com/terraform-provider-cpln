package cpln

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneIdentity_basic performs an acceptance test for the resource.
func TestAccControlPlaneIdentity_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewIdentityResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "IDENTITY") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// IdentityResourceTest defines the necessary functionality to test the resource.
type IdentityResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
}

// NewIdentityResourceTest creates a IdentityResourceTest with initialized test cases.
func NewIdentityResourceTest() IdentityResourceTest {
	// Create a resource test instance
	resourceTest := IdentityResourceTest{
		RandomName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
	}

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
func (irt *IdentityResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_identity resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_gvc" {
			continue
		}

		// Retrieve the name for the current resource
		gvcName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of GVC with name: %s", gvcName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		gvc, code, err := TestProvider.client.GetGvc(gvcName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if GVC %s exists: %w", gvcName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if gvc != nil {
			return fmt.Errorf("CheckDestroy failed: GVC %s still exists in the system", *gvc.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_identity resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewDefaultScenario creates a test case for a group using JMESPATH with initial and updated configurations.
func (irt *IdentityResourceTest) NewDefaultScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("identity-default-%s", random)
	gvcName := fmt.Sprintf("gvc-%s", random)
	agentSelfLink := GetSelfLink(OrgName, "agent", "default-agent")

	// Create the gvc case
	gvcCase := GvcResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "gvc",
			ResourceName:      "new",
			ResourceAddress:   "cpln_gvc.new",
			Name:              gvcName,
			Description:       gvcName,
			DescriptionUpdate: "gvc default description updated",
		},
	}

	// Create a gvc resource test instance
	gvcResourceTest := GvcResourceTest{}

	// Initialize the gvc config
	gvcConfig := gvcResourceTest.GvcRequiredOnly(gvcCase)

	// Build test steps
	initialConfig, initialStep := irt.BuildInitialTestStep(name, gvcConfig, gvcCase)
	caseUpdate1 := irt.BuildUpdate1TestStep(initialConfig.ProviderTestCase, gvcConfig, gvcCase, agentSelfLink)
	caseUpdate2 := irt.BuildUpdate2TestStep(initialConfig.ProviderTestCase, gvcConfig, gvcCase, agentSelfLink)
	caseUpdate3 := irt.BuildUpdate3TestStep(initialConfig.ProviderTestCase, gvcConfig, gvcCase, agentSelfLink)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName:  initialConfig.ResourceAddress,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:%s", gvcName, name),
		},
		// Update & Read
		caseUpdate1,
		caseUpdate2,
		caseUpdate3,
		// Revert the resource to its initial state
		// initialStep, // TODO: Uncomment this once the issue with network resources removal is fixed.
	}
}

// Test Cases //

// BuildInitialTestStep returns a default initial test step and its associated test case for the resource.
func (irt *IdentityResourceTest) BuildInitialTestStep(name string, gvcConfig string, gvcCase GvcResourceTestCase) (IdentityResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := IdentityResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "identity",
			ResourceName:      "new",
			ResourceAddress:   "cpln_identity.new",
			Name:              name,
			GvcName:           gvcCase.Name,
			Description:       name,
			DescriptionUpdate: "identity default description updated",
		},
		GvcConfig: gvcConfig,
		GvcCase:   gvcCase,
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: irt.RequiredOnly(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
		),
	}
}

// BuildUpdate1TestStep returns a test step for the update.
func (irt *IdentityResourceTest) BuildUpdate1TestStep(initialCase ProviderTestCase, gvcConfig string, gvcCase GvcResourceTestCase, agentSelfLink string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := IdentityResourceTestCase{
		ProviderTestCase: initialCase,
		GvcConfig:        gvcConfig,
		GvcCase:          gvcCase,
		NetworkResources: []client.IdentityNetworkResource{
			{
				Name:      StringPointer("test-network-resource-fqdn-01"),
				FQDN:      StringPointer("test-network-resource-fqdn-01.com"),
				AgentLink: &agentSelfLink,
				Ports:     &[]int{1234, 54321},
			},
			{
				Name:      StringPointer("test-network-resource-fqdn-02"),
				FQDN:      StringPointer("test-network-resource-fqdn-02.com"),
				AgentLink: &agentSelfLink,
				Ports:     &[]int{3099, 7890},
			},
		},
		NativeNetworkResources: []client.IdentityNativeNetworkResource{
			{
				Name:  StringPointer("test-native-network-resource-fqdn-01"),
				FQDN:  StringPointer("test-native-network-resource-fqdn-01.com"),
				Ports: &[]int{80, 443, 8080},
				AWSPrivateLink: &client.IdentityAwsPrivateLink{
					EndpointServiceName: StringPointer("com.amazonaws.vpce.us-west-2.vpce-svc-01af6c4c9260ac550"),
				},
			},
			{
				Name:  StringPointer("test-native-network-resource-fqdn-02"),
				FQDN:  StringPointer("test-native-network-resource-fqdn-02.com"),
				Ports: &[]int{8080},
				AWSPrivateLink: &client.IdentityAwsPrivateLink{
					EndpointServiceName: StringPointer("com.amazonaws.vpce.us-east-2.vpce-svc-01af6c4c9260ac550"),
				},
			},
		},
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: irt.WithNetworkResources(c, agentSelfLink),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", gvcCase.Name),
			c.NetworkResourcesRequiredOnlyTestCheck(),
			c.NativeNetworkResourcesWithAwsTestCheck(),
		),
	}
}

// BuildUpdate2TestStep returns a test step for the update.
func (irt *IdentityResourceTest) BuildUpdate2TestStep(initialCase ProviderTestCase, gvcConfig string, gvcCase GvcResourceTestCase, agentSelfLink string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := IdentityResourceTestCase{
		ProviderTestCase: initialCase,
		GvcConfig:        gvcConfig,
		GvcCase:          gvcCase,
		NetworkResources: []client.IdentityNetworkResource{
			{
				Name:      StringPointer("test-network-resource-fqdn-01"),
				FQDN:      StringPointer("test-network-resource-fqdn-01.com"),
				AgentLink: &agentSelfLink,
				Ports:     &[]int{1234, 54321},
			},
			{
				Name:      StringPointer("test-network-resource-fqdn-02"),
				FQDN:      StringPointer("test-network-resource-fqdn-02.com"),
				AgentLink: &agentSelfLink,
				Ports:     &[]int{3099, 7890},
			},
		},
		NativeNetworkResources: []client.IdentityNativeNetworkResource{
			{
				Name:  StringPointer("test-native-network-resource-fqdn-01"),
				FQDN:  StringPointer("test-native-network-resource-fqdn-01.com"),
				Ports: &[]int{80, 443, 8080},
				AWSPrivateLink: &client.IdentityAwsPrivateLink{
					EndpointServiceName: StringPointer("com.amazonaws.vpce.us-west-2.vpce-svc-01af6c4c9260ac550"),
				},
			},
			{
				Name:  StringPointer("test-native-network-resource-fqdn-02"),
				FQDN:  StringPointer("test-native-network-resource-fqdn-02.com"),
				Ports: &[]int{8080},
				AWSPrivateLink: &client.IdentityAwsPrivateLink{
					EndpointServiceName: StringPointer("com.amazonaws.vpce.us-east-2.vpce-svc-01af6c4c9260ac550"),
				},
			},
		},
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: irt.WithMinimalOptionals(c, agentSelfLink),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", gvcCase.Name),
			c.NetworkResourcesRequiredOnlyTestCheck(),
			c.NativeNetworkResourcesWithAwsTestCheck(),
			c.TestCheckNestedBlocks("aws_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-aws-%s", irt.RandomName)),
					"policy_refs":        []string{"aws::/job-function/SupportUser", "aws::AWSSupportAccess"},
				},
			}),
			resource.TestCheckResourceAttr(c.ResourceAddress, "aws_access_policy.0.trust_policy.0.version", "2012-10-17"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "aws_access_policy.0.trust_policy.0.statement.0.Effect", "Allow"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "aws_access_policy.0.trust_policy.0.statement.0.Action", "sts:AssumeRole"),
			c.TestCheckNestedBlocks("gcp_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-gcp-%s", irt.RandomName)),
					"binding": []map[string]interface{}{
						{
							"resource": "//cloudresourcemanager.googleapis.com/projects/cpln-test",
							"roles":    []string{"roles/appengine.appViewer", "roles/actions.Viewer"},
						},
						{
							"resource": "//iam.googleapis.com/projects/cpln-test/serviceAccounts/cpln-tf@cpln-test.iam.gserviceaccount.com",
							"roles":    []string{"roles/editor", "roles/iam.serviceAccountUser"},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("azure_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-azure-%s", irt.RandomName)),
				},
			}),
			c.TestCheckNestedBlocks("ngs_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-ngs-%s", irt.RandomName)),
				},
			}),
		),
	}
}

// BuildUpdate3TestStep returns a test step for the update.
func (irt *IdentityResourceTest) BuildUpdate3TestStep(initialCase ProviderTestCase, gvcConfig string, gvcCase GvcResourceTestCase, agentSelfLink string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := IdentityResourceTestCase{
		ProviderTestCase: initialCase,
		GvcConfig:        gvcConfig,
		GvcCase:          gvcCase,
		NetworkResources: []client.IdentityNetworkResource{
			{
				Name:      StringPointer("test-network-resource-fqdn-01"),
				FQDN:      StringPointer("test-network-resource-fqdn-01.com"),
				AgentLink: &agentSelfLink,
				Ports:     &[]int{1234, 54321},
			},
			{
				Name:      StringPointer("test-network-resource-fqdn-02"),
				FQDN:      StringPointer("test-network-resource-fqdn-02.com"),
				AgentLink: &agentSelfLink,
				Ports:     &[]int{3099, 7890},
			},
		},
		NativeNetworkResources: []client.IdentityNativeNetworkResource{
			{
				Name:  StringPointer("test-native-network-resource-fqdn-01"),
				FQDN:  StringPointer("test-native-network-resource-fqdn-01.com"),
				Ports: &[]int{80, 443, 8080},
				AWSPrivateLink: &client.IdentityAwsPrivateLink{
					EndpointServiceName: StringPointer("com.amazonaws.vpce.us-west-2.vpce-svc-01af6c4c9260ac550"),
				},
			},
			{
				Name:  StringPointer("test-native-network-resource-fqdn-02"),
				FQDN:  StringPointer("test-native-network-resource-fqdn-02.com"),
				Ports: &[]int{8080},
				AWSPrivateLink: &client.IdentityAwsPrivateLink{
					EndpointServiceName: StringPointer("com.amazonaws.vpce.us-east-2.vpce-svc-01af6c4c9260ac550"),
				},
			},
		},
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: irt.WithAllAttributes(c, agentSelfLink),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", gvcCase.Name),
			c.NetworkResourcesRequiredOnlyTestCheck(),
			c.NativeNetworkResourcesWithAwsTestCheck(),
			c.TestCheckNestedBlocks("aws_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-aws-%s", irt.RandomName)),
					"role_name":          "rds-monitoring-role",
				},
			}),
			c.TestCheckNestedBlocks("gcp_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-gcp-%s", irt.RandomName)),
					"scopes":             []string{"https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/cloud-suite"},
					"service_account":    "cpln-tf@cpln-test.iam.gserviceaccount.com",
				},
			}),
			c.TestCheckNestedBlocks("azure_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-azure-%s", irt.RandomName)),
					"role_assignment": []map[string]interface{}{
						{
							"scope": "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group",
							"roles": []string{"AcrPull", "AcrPush"},
						},
						{
							"scope": "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group/providers/Microsoft.Storage/storageAccounts/cplntest",
							"roles": []string{"Support Request Contributor"},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("ngs_access_policy", []map[string]interface{}{
				{
					"cloud_account_link": GetSelfLink(OrgName, "cloudaccount", fmt.Sprintf("tf-ca-ngs-%s", irt.RandomName)),
					"pub": []map[string]interface{}{
						{
							"allow": []string{"pa1", "pa2"},
							"deny":  []string{"pd1", "pd2"},
						},
					},
					"sub": []map[string]interface{}{
						{
							"allow": []string{"sa1", "sa2"},
							"deny":  []string{"sd1", "sd2"},
						},
					},
					"resp": []map[string]interface{}{
						{
							"max": 1,
							"ttl": "5m",
						},
					},
					"subs":    1,
					"data":    2,
					"payload": 3,
				},
			}),
		),
	}
}

// Configs //

// RequiredOnly returns a minimal HCL block for a resource using only required fields.
func (irt *IdentityResourceTest) RequiredOnly(c IdentityResourceTestCase) string {
	return fmt.Sprintf(`
# GVC Resource
%s

resource "cpln_identity" "%s" {
  depends_on = [%s]

  name = "%s"
	gvc  = "%s"
}
`, c.GvcConfig, c.ResourceName, c.GvcCase.ResourceAddress, c.Name, c.GvcCase.Name)
}

// WithNetworkResources returns a minimal HCL block for a resource with network resources.
func (irt *IdentityResourceTest) WithNetworkResources(c IdentityResourceTestCase, agentSelfLink string) string {
	return fmt.Sprintf(`
# GVC Resource
%s

resource "cpln_identity" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"
  gvc         = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  # Network Resources
%s

  # Native Network Resources
%s
}
`, c.GvcConfig, c.ResourceName, c.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate, c.GvcCase.Name,
		irt.NetworkResourcesRequiredOnlyHcl(c.NetworkResources, agentSelfLink), irt.NativeNetworkResourceWithAws(c.NativeNetworkResources),
	)
}

// WithNetworkResources returns a minimal HCL block for a resource with minimal optionals.
func (irt *IdentityResourceTest) WithMinimalOptionals(c IdentityResourceTestCase, agentSelfLink string) string {
	return fmt.Sprintf(`
variable org_name {
  type    = string
  default = "%s"
}

variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_cloud_account" "aws" {
  name = "tf-ca-aws-${var.random_name}"

  aws {
    role_arn = "arn:aws:iam::1234:role/test_role"
  }
}

resource "cpln_cloud_account" "azure" {
  name = "tf-ca-azure-${var.random_name}"

  azure {
    secret_link = "/org/${var.org_name}/secret/tf_secret_azure"
  }
}

resource "cpln_cloud_account" "gcp" {
  name = "tf-ca-gcp-${var.random_name}"

  gcp {
    project_id = "cpln_gcp_project_1234"
  }
}

resource "cpln_cloud_account" "ngs" {
  name = "tf-ca-ngs-${var.random_name}"

  ngs {
    secret_link = "/org/${var.org_name}/secret/tf_secret_nats_account"
  }
}

# GVC Resource
%s

resource "cpln_identity" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"
  gvc         = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  # Network Resources
%s

  # Native Network Resources
%s

  aws_access_policy {
    cloud_account_link = cpln_cloud_account.aws.self_link
    policy_refs        = ["aws::/job-function/SupportUser", "aws::AWSSupportAccess"]

    trust_policy {
      version   = "2012-10-17"
      statement = [{
          Effect    = "Allow"
          Principal = jsonencode({ Service = "ec2.amazonaws.com" })
          Action    = "sts:AssumeRole"
      }]
    }
  }

  gcp_access_policy {
    cloud_account_link = cpln_cloud_account.gcp.self_link

    binding {
      resource = "//cloudresourcemanager.googleapis.com/projects/cpln-test"
      roles    = ["roles/appengine.appViewer", "roles/actions.Viewer"]
    }

    binding {
      resource = "//iam.googleapis.com/projects/cpln-test/serviceAccounts/cpln-tf@cpln-test.iam.gserviceaccount.com"
      roles = ["roles/editor", "roles/iam.serviceAccountUser"]
    }
  }

  azure_access_policy {
    cloud_account_link = cpln_cloud_account.azure.self_link
  }

  ngs_access_policy {
    cloud_account_link = cpln_cloud_account.ngs.self_link
  }
}
`, OrgName, irt.RandomName, c.GvcConfig, c.ResourceName, c.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate, c.GvcCase.Name,
		irt.NetworkResourcesRequiredOnlyHcl(c.NetworkResources, agentSelfLink), irt.NativeNetworkResourceWithAws(c.NativeNetworkResources),
	)
}

// WithAllAttributes returns an HCL block for a resource with all attributes.
func (irt *IdentityResourceTest) WithAllAttributes(c IdentityResourceTestCase, agentSelfLink string) string {
	return fmt.Sprintf(`
variable org_name {
  type    = string
  default = "%s"
}

variable random_name {
  type    = string
  default = "%s"
}

resource "cpln_cloud_account" "aws" {
  name = "tf-ca-aws-${var.random_name}"

  aws {
    role_arn = "arn:aws:iam::1234:role/test_role"
  }
}

resource "cpln_cloud_account" "azure" {
  name = "tf-ca-azure-${var.random_name}"

  azure {
    secret_link = "/org/${var.org_name}/secret/tf_secret_azure"
  }
}

resource "cpln_cloud_account" "gcp" {
  name = "tf-ca-gcp-${var.random_name}"

  gcp {
    project_id = "cpln_gcp_project_1234"
  }
}

resource "cpln_cloud_account" "ngs" {
  name = "tf-ca-ngs-${var.random_name}"

  ngs {
    secret_link = "/org/${var.org_name}/secret/tf_secret_nats_account"
  }
}

# GVC Resource
%s

resource "cpln_identity" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"
  gvc         = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  # Network Resources
%s

  # Native Network Resources
%s

  aws_access_policy {
    cloud_account_link = cpln_cloud_account.aws.self_link
		role_name          = "rds-monitoring-role"
  }

  gcp_access_policy {
    cloud_account_link = cpln_cloud_account.gcp.self_link
		scopes             = ["https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/cloud-suite"]
    service_account    = "cpln-tf@cpln-test.iam.gserviceaccount.com"
  }

  azure_access_policy {
    cloud_account_link = cpln_cloud_account.azure.self_link

    role_assignment {
      scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group"
      roles = ["AcrPull",	"AcrPush"]
    }

    role_assignment {
      scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group/providers/Microsoft.Storage/storageAccounts/cplntest"
      roles = ["Support Request Contributor"]
    }
  }

  ngs_access_policy {
    cloud_account_link = cpln_cloud_account.ngs.self_link

    pub {
      allow = ["pa1", "pa2"]
      deny  = ["pd1", "pd2"]
    }

    sub {
      allow = ["sa1", "sa2"]
      deny  = ["sd1", "sd2"]
    }

    resp {
      max = 1
      ttl = "5m"
    }

    subs = 1
    data = 2
    payload = 3
  }
}
`, OrgName, irt.RandomName, c.GvcConfig, c.ResourceName, c.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate, c.GvcCase.Name,
		irt.NetworkResourcesRequiredOnlyHcl(c.NetworkResources, agentSelfLink), irt.NativeNetworkResourceWithAws(c.NativeNetworkResources),
	)
}

// Network Resources //

// NetworkResourcesRequiredOnlyHcl builds HCL blocks for required-only network resources
func (irt *IdentityResourceTest) NetworkResourcesRequiredOnlyHcl(list []client.IdentityNetworkResource, agentSelfLink string) string {
	// Initialize a slice to collect HCL snippets
	output := []string{}

	// Iterate over each network resource item
	for _, item := range list {
		// Format HCL block with required fields for this resource
		hcl := fmt.Sprintf(`
  network_resource {
    name        = "%s"
    agent_link  = "%s"
    fqdn        = "%s"
    ports       = %s
  }
		`, *item.Name, agentSelfLink, *item.FQDN, IntSliceToString(*item.Ports))

		// Append formatted HCL to output
		output = append(output, hcl)
	}

	// Join all HCL blocks with spacing and return
	return strings.Join(output, "\n\n")
}

// NetworkResourceWithAllAttributesHcl builds HCL blocks for network resources with all attributes
func (irt *IdentityResourceTest) NetworkResourceWithAllAttributesHcl(list []client.IdentityNetworkResource, agentSelfLink string) string {
	// Initialize a slice to collect HCL snippets
	output := []string{}

	// Iterate over each network resource item
	for _, item := range list {
		// Format HCL block with full attribute set for this resource
		hcl := fmt.Sprintf(`
  network_resource {
    name        = "%s"
    agent_link  = "%s"
    ips         = %s
    fqdn        = "%s"
    resolver_ip = "%s"
    ports       = %s
  }
	`, *item.Name, agentSelfLink, StringSliceToString(*item.IPs), *item.FQDN, *item.ResolverIP, IntSliceToString(*item.Ports))

		// Append formatted HCL to output
		output = append(output, hcl)
	}

	// Join all HCL blocks with spacing and return
	return strings.Join(output, "\n\n")
}

// Native Network Resources //

// NativeNetworkResourceWithAws builds HCL block for AWS native network resources
func (irt *IdentityResourceTest) NativeNetworkResourceWithAws(list []client.IdentityNativeNetworkResource) string {
	// Initialize a slice to collect HCL snippets
	output := []string{}

	// Iterate over each native AWS network resource item
	for _, item := range list {
		// Format HCL block with AWS private link configuration
		hcl := fmt.Sprintf(`
  native_network_resource {
    name        = "%s"
    fqdn        = "%s"
    ports       = %s

    aws_private_link {
      endpoint_service_name = "%s"
    }
  }
	`, *item.Name, *item.FQDN, IntSliceToString(*item.Ports), *item.AWSPrivateLink.EndpointServiceName)

		// Append formatted HCL to output
		output = append(output, hcl)
	}

	// Join all HCL blocks with spacing and return
	return strings.Join(output, "\n\n")
}

// NativeNetworkResourceWithGcp builds HCL block for GCP native network resources
func (irt *IdentityResourceTest) NativeNetworkResourceWithGcp(list []client.IdentityNativeNetworkResource) string {
	// Initialize a slice to collect HCL snippets
	output := []string{}

	// Iterate over each native GCP network resource item
	for _, item := range list {
		// Format HCL block with GCP service connect configuration
		hcl := fmt.Sprintf(`
  native_network_resource {
    name        = "%s"
    fqdn        = "%s"
    ports       = %s

    gcp_service_connect {
      target_service = "%s"
    }
  }
	`, *item.Name, *item.FQDN, IntSliceToString(*item.Ports), *item.GCPServiceConnect.TargetService)

		// Append formatted HCL to output
		output = append(output, hcl)
	}

	// Join all HCL blocks with spacing and return
	return strings.Join(output, "\n\n")
}

/*** Resource Test Case ***/

// IdentityResourceTestCase defines a specific resource test case.
type IdentityResourceTestCase struct {
	ProviderTestCase
	GvcConfig              string
	GvcCase                GvcResourceTestCase
	NetworkResources       []client.IdentityNetworkResource
	NativeNetworkResources []client.IdentityNativeNetworkResource
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (irtc *IdentityResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of identity: %s. Total resources: %d", irtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[irtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", irtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != irtc.Name {
			return fmt.Errorf("resource ID %s does not match expected identity name %s", rs.Primary.ID, irtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteIdentity, _, err := TestProvider.client.GetIdentity(irtc.Name, irtc.GvcCase.Name)
		if err != nil {
			return fmt.Errorf("error retrieving identity from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteIdentity.Name != irtc.Name {
			return fmt.Errorf("mismatch in identity name: expected %s, got %s", irtc.Name, *remoteIdentity.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("identity %s verified successfully in both state and external system.", irtc.Name))
		return nil
	}
}

// Network Resources //

// NetworkResourcesRequiredOnlyTestCheck returns a TestCheckFunc validating required network_resource fields
func (irtc *IdentityResourceTestCase) NetworkResourcesRequiredOnlyTestCheck() resource.TestCheckFunc {
	// Initialize a slice to collect network resource maps
	v := []map[string]interface{}{}
	// Iterate over each defined network resource
	for _, item := range irtc.NetworkResources {
		// Append a map of required attributes for this network resource
		v = append(v, map[string]interface{}{
			"name":       *item.Name,
			"agent_link": *item.AgentLink,
			"fqdn":       *item.FQDN,
			"ports":      *item.Ports,
		})
	}
	// Return a nested block test check using the collected network resource maps
	return irtc.TestCheckNestedBlocks("network_resource", v)
}

// NetworkResourcesWithAllAttributesTestCheck returns a TestCheckFunc validating all network_resource fields including IPs and resolver details
func (irtc *IdentityResourceTestCase) NetworkResourcesWithAllAttributesTestCheck(agentSelfLink string) resource.TestCheckFunc {
	// Initialize a slice to collect detailed network resource maps
	v := []map[string]interface{}{}
	// Iterate over each defined network resource
	for _, item := range irtc.NetworkResources {
		// Append a map of all attributes for this network resource
		v = append(v, map[string]interface{}{
			"name":        *item.Name,
			"agent_link":  agentSelfLink,
			"ips":         *item.IPs,
			"fqdn":        *item.FQDN,
			"resolver_ip": *item.ResolverIP,
			"ports":       *item.Ports,
		})
	}
	// Return a nested block test check using the detailed network resource maps
	return irtc.TestCheckNestedBlocks("network_resource", v)
}

// Native Network Resources //

// NativeNetworkResourcesWithAwsTestCheck returns a TestCheckFunc validating AWS native_network_resource blocks
func (irtc *IdentityResourceTestCase) NativeNetworkResourcesWithAwsTestCheck() resource.TestCheckFunc {
	// Initialize a slice to collect AWS native network resource maps
	v := []map[string]interface{}{}

	// Iterate over each defined native network resource for AWS
	for _, item := range irtc.NativeNetworkResources {
		// Append a map of AWS-specific attributes for this native network resource
		v = append(v, map[string]interface{}{
			"name":  *item.Name,
			"fqdn":  *item.FQDN,
			"ports": *item.Ports,
			"aws_private_link": []map[string]interface{}{
				{"endpoint_service_name": *item.AWSPrivateLink.EndpointServiceName},
			},
		})
	}

	// Return a nested block test check using the AWS native network resource maps
	return irtc.TestCheckNestedBlocks("native_network_resource", v)
}

// NativeNetworkResourcesWithGcpTestCheck returns a TestCheckFunc validating GCP native_network_resource blocks
func (irtc *IdentityResourceTestCase) NativeNetworkResourcesWithGcpTestCheck() resource.TestCheckFunc {
	// Initialize a slice to collect GCP native network resource maps
	v := []map[string]interface{}{}

	// Iterate over each defined native network resource for GCP
	for _, item := range irtc.NativeNetworkResources {
		// Append a map of GCP-specific attributes for this native network resource
		v = append(v, map[string]interface{}{
			"name":  *item.Name,
			"fqdn":  *item.FQDN,
			"ports": *item.Ports,
			"gcp_service_connect": []map[string]interface{}{
				{"target_service": *item.GCPServiceConnect.TargetService},
			},
		})
	}

	// Return a nested block test check using the GCP native network resource maps
	return irtc.TestCheckNestedBlocks("native_network_resource", v)
}
