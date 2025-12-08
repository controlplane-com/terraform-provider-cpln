package cpln

import (
	"errors"
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneCatalogTemplate_basic performs an acceptance test for the resource.
func Skip_TestAccControlPlaneCatalogTemplate_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewCatalogTemplateResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "CATALOG_TEMPLATE") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// CatalogTemplateResourceTest defines the necessary functionality to test the resource.
type CatalogTemplateResourceTest struct {
	Steps []resource.TestStep
}

// NewCatalogTemplateResourceTest creates a CatalogTemplateResourceTest with initialized test cases.
func NewCatalogTemplateResourceTest() CatalogTemplateResourceTest {
	// Create a resource test instance
	resourceTest := CatalogTemplateResourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewWithGvcScenario()...)
	steps = append(steps, resourceTest.NewWithoutGvcScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (ctrt *CatalogTemplateResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_catalog_template resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_catalog_template" {
			continue
		}

		// Retrieve the name for the current resource
		catalogTemplateName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of catalog template with name: %s", catalogTemplateName))

		// Build the query to find helm release secrets for this release
		query := client.Query{
			Kind: StringPointer("secret"),
			Spec: &client.QuerySpec{
				Match: StringPointer("all"),
				Terms: &[]client.QueryTerm{
					{
						Op:       StringPointer("~"),
						Property: StringPointer("name"),
						Value:    StringPointer("cpln-helm-release-"),
					},
					{
						Op:    StringPointer("="),
						Tag:   StringPointer("name"),
						Value: StringPointer(catalogTemplateName),
					},
					{
						Op:  StringPointer("exists"),
						Tag: StringPointer("cpln/marketplace"),
					},
					{
						Op:  StringPointer("exists"),
						Tag: StringPointer("cpln/marketplace-template"),
					},
					{
						Op:  StringPointer("exists"),
						Tag: StringPointer("cpln/marketplace-template-version"),
					},
				},
			},
		}

		// Use the TestProvider client to check if the API resource still exists in the data service
		catalogTemplate, code, err := TestProvider.client.GetMarketplaceRelease(catalogTemplateName, query)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if catalog template %s exists: %w", catalogTemplateName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if catalogTemplate != nil && catalogTemplate.Name == catalogTemplateName {
			return fmt.Errorf("CheckDestroy failed: catalog template %s still exists in the system", catalogTemplate.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_catalog_template resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewWithGvcScenario creates a test case for a catalog template deployment with a GVC provided.
func (ctrt *CatalogTemplateResourceTest) NewWithGvcScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("catalog-template-%s", random)
	gvcName := fmt.Sprintf("catalog-gvc-%s", random)
	resourceName := "with-gvc"

	// Define template details
	templateName := "redis"
	version := "3.0.1"
	updateVersion := "3.0.2"

	// Build test steps
	initialConfig, initialStep := ctrt.BuildInitialTestStepWithGvc(resourceName, name, gvcName, templateName, version)
	caseUpdate1 := ctrt.BuildUpdate1TestStepWithGvc(initialConfig.ProviderTestCase, gvcName, templateName, version)
	caseUpdate2 := ctrt.BuildUpdate2TestStepWithGvc(initialConfig.ProviderTestCase, gvcName, templateName, updateVersion)

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
		initialStep,
	}
}

// NewWithoutGvcScenario creates a test case for a catalog template deployment without a GVC (template creates its own GVC).
func (ctrt *CatalogTemplateResourceTest) NewWithoutGvcScenario() []resource.TestStep {
	// Generate a unique name for the resources
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("catalog-template-no-gvc-%s", random)
	resourceName := "without-gvc"

	// Define template details for a template that creates its own GVC
	templateName := "cockroach"
	version := "1.0.0"

	// Build test steps
	initialConfig, initialStep := ctrt.BuildInitialTestStepWithoutGvc(resourceName, name, templateName, version)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName: initialConfig.ResourceAddress,
			ImportState:  true,
		},
	}
}

// Test Cases //

// BuildInitialTestStepWithGvc returns a default initial test step and its associated test case for the resource with a GVC.
func (ctrt *CatalogTemplateResourceTest) BuildInitialTestStepWithGvc(resourceName string, name string, gvcName string, templateName string, version string) (CatalogTemplateResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := CatalogTemplateResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "catalog_template",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_catalog_template.%s", resourceName),
			Name:            name,
		},
		TemplateName: templateName,
		Version:      version,
		GvcName:      gvcName,
		Values:       getValuesRedisInitial(),
	}

	// Initialize and return the initial test step
	return c, resource.TestStep{
		Config: ctrt.RequiredOnlyWithGvc(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			resource.TestCheckResourceAttr(c.ResourceAddress, "id", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "name", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "template", c.TemplateName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "version", c.Version),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", c.GvcName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "values", c.Values),
			resource.TestCheckResourceAttrSet(c.ResourceAddress, "resources.#"),
		),
	}
}

// BuildUpdate1TestStepWithGvc returns a test step for the first update (update values).
func (ctrt *CatalogTemplateResourceTest) BuildUpdate1TestStepWithGvc(initialCase ProviderTestCase, gvcName string, templateName string, version string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CatalogTemplateResourceTestCase{
		ProviderTestCase: initialCase,
		TemplateName:     templateName,
		Version:          version,
		GvcName:          gvcName,
		Values:           getValuesRedisInitial(),
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: ctrt.Update1WithGvc(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			resource.TestCheckResourceAttr(c.ResourceAddress, "id", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "name", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "template", c.TemplateName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "version", c.Version),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", c.GvcName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "values", c.Values),
			resource.TestCheckResourceAttrSet(c.ResourceAddress, "resources.#"),
		),
	}
}

// BuildUpdate2TestStepWithGvc returns a test step for the second update (change version).
func (ctrt *CatalogTemplateResourceTest) BuildUpdate2TestStepWithGvc(initialCase ProviderTestCase, gvcName string, templateName string, version string) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := CatalogTemplateResourceTestCase{
		ProviderTestCase: initialCase,
		TemplateName:     templateName,
		Version:          version,
		GvcName:          gvcName,
		Values:           getValuesRedisUpdate2(),
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: ctrt.Update2WithGvc(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			resource.TestCheckResourceAttr(c.ResourceAddress, "id", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "name", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "template", c.TemplateName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "version", c.Version),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", c.GvcName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "values", c.Values),
			resource.TestCheckResourceAttrSet(c.ResourceAddress, "resources.#"),
		),
	}
}

// BuildInitialTestStepWithoutGvc returns a default initial test step and its associated test case for the resource without a GVC.
func (ctrt *CatalogTemplateResourceTest) BuildInitialTestStepWithoutGvc(resourceName string, name string, templateName string, version string) (CatalogTemplateResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := CatalogTemplateResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "catalog_template",
			ResourceName:    resourceName,
			ResourceAddress: fmt.Sprintf("cpln_catalog_template.%s", resourceName),
			Name:            name,
		},
		TemplateName: templateName,
		Version:      version,
		GvcName:      name,
		Values:       getValuesCockroachTemplate(name),
	}

	// Initialize and return the initial test step
	return c, resource.TestStep{
		Config: ctrt.RequiredOnlyWithoutGvc(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			resource.TestCheckResourceAttr(c.ResourceAddress, "id", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "name", c.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "template", c.TemplateName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "version", c.Version),
			resource.TestCheckResourceAttr(c.ResourceAddress, "values", c.Values),
			resource.TestCheckResourceAttrSet(c.ResourceAddress, "resources.#"),
		),
	}
}

// Configs //

// RequiredOnlyWithGvc returns a minimal HCL block for a resource using only required fields with a GVC.
func (ctrt *CatalogTemplateResourceTest) RequiredOnlyWithGvc(c CatalogTemplateResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "catalog_test_gvc" {
  name = "%s"
}

resource "cpln_catalog_template" "%s" {
  name     = "%s"
  template = "%s"
  version  = "%s"
  gvc      = cpln_gvc.catalog_test_gvc.name
  values   = <<-EOT
%sEOT
}
`, c.GvcName, c.ResourceName, c.Name, c.TemplateName, c.Version, c.Values)
}

// Update1WithGvc returns an HCL block for the first update with different values.
func (ctrt *CatalogTemplateResourceTest) Update1WithGvc(c CatalogTemplateResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "catalog_test_gvc" {
  name = "%s"
}

resource "cpln_catalog_template" "%s" {
  name     = "%s"
  template = "%s"
  version  = "%s"
  gvc      = cpln_gvc.catalog_test_gvc.name
  values   = <<-EOT
%sEOT
}
`, c.GvcName, c.ResourceName, c.Name, c.TemplateName, c.Version, c.Values)
}

// Update2WithGvc returns an HCL block for the second update with version change.
func (ctrt *CatalogTemplateResourceTest) Update2WithGvc(c CatalogTemplateResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_gvc" "catalog_test_gvc" {
  name = "%s"
}

resource "cpln_catalog_template" "%s" {
  name     = "%s"
  template = "%s"
  version  = "%s"
  gvc      = cpln_gvc.catalog_test_gvc.name
  values   = <<-EOT
%sEOT
}
`, c.GvcName, c.ResourceName, c.Name, c.TemplateName, c.Version, c.Values)
}

// RequiredOnlyWithoutGvc returns a minimal HCL block for a resource using only required fields without a GVC.
func (ctrt *CatalogTemplateResourceTest) RequiredOnlyWithoutGvc(c CatalogTemplateResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_catalog_template" "%s" {
  name     = "%s"
  template = "%s"
  version  = "%s"
  values   = <<-EOT
%sEOT
}
`, c.ResourceName, c.Name, c.TemplateName, c.Version, c.Values)
}

/*** Resource Test Case ***/

// CatalogTemplateResourceTestCase defines a specific resource test case.
type CatalogTemplateResourceTestCase struct {
	ProviderTestCase
	TemplateName string
	Version      string
	GvcName      string
	Values       string
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (ctrtc *CatalogTemplateResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of catalog template: %s. Total resources: %d", ctrtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[ctrtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", ctrtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != ctrtc.Name {
			return fmt.Errorf("resource ID %s does not match expected catalog template name %s", rs.Primary.ID, ctrtc.Name)
		}

		// Build the query to find helm release secrets for this release
		query := client.Query{
			Kind: StringPointer("secret"),
			Spec: &client.QuerySpec{
				Match: StringPointer("all"),
				Terms: &[]client.QueryTerm{
					{
						Op:       StringPointer("~"),
						Property: StringPointer("name"),
						Value:    StringPointer("cpln-helm-release-"),
					},
					{
						Op:    StringPointer("="),
						Tag:   StringPointer("name"),
						Value: StringPointer(ctrtc.Name),
					},
					{
						Op:  StringPointer("exists"),
						Tag: StringPointer("cpln/marketplace"),
					},
					{
						Op:  StringPointer("exists"),
						Tag: StringPointer("cpln/marketplace-template"),
					},
					{
						Op:  StringPointer("exists"),
						Tag: StringPointer("cpln/marketplace-template-version"),
					},
				},
			},
		}

		// Retrieve the API resource from the external system using the provider client
		remoteCatalogTemplate, _, err := TestProvider.client.GetMarketplaceRelease(ctrtc.Name, query)
		if err != nil {
			return fmt.Errorf("error retrieving catalog template from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if remoteCatalogTemplate.Name != ctrtc.Name {
			return fmt.Errorf("mismatch in catalog template name: expected %s, got %s", ctrtc.Name, remoteCatalogTemplate.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("Catalog Template %s verified successfully in both state and external system.", ctrtc.Name))
		return nil
	}
}

/*** Test Values Helpers ***/

// getValuesRedisInitial returns the initial Redis configuration values used in multiple test cases
func getValuesRedisInitial() string {
	return `redis:
  image: redis/redis-stack:7.4.0-v3
  resources:
    cpu: 200m
    memory: 256Mi
    minCpu: 80m
    minMemory: 128Mi
  replicas: 3
  timeoutSeconds: 15
  auth:
    fromSecret:
      enabled: false
      name: example-redis-auth-password
      passwordKey: password
    password:
      enabled: false
      value: your-password
  serverCommand: redis-stack-server  # Can be overridden based on the version of redis image
  publicAccess:
    enabled: false
    address: redis-test.example-cpln.com
  firewall:
    internal_inboundAllowType: "same-gvc" # Options: same-org / same-gvc(Recommended) / workload-list
    external_inboundAllowCIDR: 0.0.0.0/0 # Provide a comma-separated list
    # # You can specify additional workloads with either same-gvc or workload-list:
    # inboundAllowWorkload:
    #   - //gvc/main-redis/workload/main-redis-sentinel
    #   - //gvc/client-gvc/workload/client
    external_outboundAllowCIDR: "0.0.0.0/0" # Provide a comma-separated list
  env: []
  dataDir: /data
  persistence:
    enabled: false
    volumes:
      data:
        initialCapacity: 10 # In GB
        performanceClass: general-purpose-ssd # general-purpose-ssd / high-throughput-ssd (Min 1000GB)
        fileSystemType: ext4 # ext4 / xfs
        snapshots:
          retentionDuration: 7d
          schedule: 0 0 * * * # UTC
        autoscaling:
          maxCapacity: 100 # In GB
          minFreePercentage: 20
          scalingFactor: 1.2

sentinel:
  image: redis/redis-stack:7.4.0-v3
  resources:
    cpu: 200m
    memory: 256Mi
    minCpu: 80m
    minMemory: 128Mi
  replicas: 3
  timeoutSeconds: 10
  quorumAutoCalculation: true  # Set to false if you want to override quorum. Quorum is (replicas/2)+1
  quorumOverride: null  # Only used if quorumAutoCalculation is false
  auth:
    fromSecret:
      enabled: false
      name: example-redis-auth-password
      passwordKey: password
    password:
      enabled: false
      value: your-password
  publicAccess:
    enabled: false
    address: redis-sentinel-test.example-cpln.com
  firewall:
    internal_inboundAllowType: "same-gvc" # Options: same-org / same-gvc(Recommended)
    external_inboundAllowCIDR: 0.0.0.0/0 # Provide a comma-separated list
    # # You can specify additional workloads with either same-gvc or workload-list:
    # inboundAllowWorkload:
    #   - //gvc/main-redis/workload/main-redis-sentinel
    #   - //gvc/client-gvc/workload/client
    external_outboundAllowCIDR: "0.0.0.0/0" # Provide a comma-separated list
  env: []
  persistence:
    enabled: false
    volumes:
      data:
        initialCapacity: 10 # In GB
        performanceClass: general-purpose-ssd # general-purpose-ssd / high-throughput-ssd (Min 1000GB)
        fileSystemType: ext4 # ext4 / xfs
        snapshots:
          retentionDuration: 7d
          schedule: 0 0 * * * # UTC
        autoscaling:
          maxCapacity: 50 # In GB
          minFreePercentage: 20
          scalingFactor: 1.2
`
}

// getValuesRedisUpdate2 returns the updated Redis configuration values for version update test
func getValuesRedisUpdate2() string {
	return `redis:
  image: redis:7.2
  resources:
    cpu: 200m
    memory: 256Mi
    minCpu: 80m
    minMemory: 128Mi
  replicas: 2
  timeoutSeconds: 15
  auth:
    fromSecret:
      enabled: false
      name: example-redis-auth-password
      passwordKey: password
    password:
      enabled: false
      value: your-password
  serverCommand: redis-server  # Can be overridden based on the version of redis image
  # extraArgs: "--maxclients 20000"
  publicAccess:
    enabled: false
    address: redis-test.example-cpln.com
  firewall:
    internal_inboundAllowType: "same-gvc" # Options: same-org / same-gvc(Recommended) / workload-list
    # external_inboundAllowCIDR: 0.0.0.0/0 # Provide a comma-separated list
    # # You can specify additional workloads with either same-gvc or workload-list:
    # inboundAllowWorkload:
    #   - //gvc/main-redis/workload/main-redis-sentinel
    #   - //gvc/client-gvc/workload/client
    # external_outboundAllowCIDR: "0.0.0.0/0" # Provide a comma-separated list
  env: []
  tags: {}
  dataDir: /data
  persistence:
    enabled: false
    volumes:
      data:
        initialCapacity: 10 # In GB
        performanceClass: general-purpose-ssd # general-purpose-ssd / high-throughput-ssd (Min 1000GB)
        fileSystemType: ext4 # ext4 / xfs
        snapshots:
          retentionDuration: 7d
          schedule: 0 0 * * * # UTC
        autoscaling:
          maxCapacity: 100 # In GB
          minFreePercentage: 20
          scalingFactor: 1.2

sentinel:
  image: redis:7.2
  resources:
    cpu: 200m
    memory: 256Mi
    minCpu: 80m
    minMemory: 128Mi
  replicas: 3
  timeoutSeconds: 10
  quorumAutoCalculation: true  # Set to false if you want to override quorum. Quorum is (replicas/2)+1
  quorumOverride: null  # Only used if quorumAutoCalculation is false
  auth:
    fromSecret:
      enabled: false
      name: example-redis-auth-password
      passwordKey: password
    password:
      enabled: false
      value: your-password
  publicAccess:
    enabled: false
    address: redis-sentinel-test.example-cpln.com
  firewall:
    internal_inboundAllowType: "same-gvc" # Options: same-org / same-gvc(Recommended)
    # external_inboundAllowCIDR: 0.0.0.0/0 # Provide a comma-separated list
    # # You can specify additional workloads with either same-gvc or workload-list:
    # inboundAllowWorkload:
    #   - //gvc/main-redis/workload/main-redis-sentinel
    #   - //gvc/client-gvc/workload/client
    # external_outboundAllowCIDR: "0.0.0.0/0" # Provide a comma-separated list
  env: []
  tags: {}
  persistence:
    enabled: false
    volumes:
      data:
        initialCapacity: 10 # In GB
        performanceClass: general-purpose-ssd # general-purpose-ssd / high-throughput-ssd (Min 1000GB)
        fileSystemType: ext4 # ext4 / xfs
        snapshots:
          retentionDuration: 7d
          schedule: 0 0 * * * # UTC
        autoscaling:
          maxCapacity: 50 # In GB
          minFreePercentage: 20
          scalingFactor: 1.2
`
}

// getValuesCockroachTemplate returns the Cockroach template values with GVC name interpolated
func getValuesCockroachTemplate(gvcName string) string {
	return fmt.Sprintf(`# Global Virtual Cluster (gvc) settings
gvc:
  name: %s
  locations:
    - name: aws-eu-central-1
      replicas: 3

resources:
  cpu: 2000m
  memory: 4096Mi

database:
  name: mydb
  user: myuser

cockroach_defaults:
  workload_name: cockroach
  sql_port: 26257
  http_port: 8080

internal_access:
  type: same-gvc # options: same-gvc, same-org, workload-list
  workloads:  # Note: can only be used if type is same-gvc or workload-list
    #- //gvc/GVC_NAME/workload/WORKLOAD_NAME
    #- //gvc/GVC_NAME/workload/WORKLOAD_NAME
`, gvcName)
}
