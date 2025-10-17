package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceWorkload_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceWorkload_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewWorkloadDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_WORKLOAD") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// WorkloadDataSourceTest defines the necessary functionality to test the data source.
type WorkloadDataSourceTest struct {
	Steps []resource.TestStep
}

// NewWorkloadDataSourceTest creates a WorkloadDataSourceTest with initialized test cases.
func NewWorkloadDataSourceTest() WorkloadDataSourceTest {
	// Create a data source test instance
	dataSourceTest := WorkloadDataSourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, dataSourceTest.NewDefaultScenario()...)

	// Set the cases for the data source test
	dataSourceTest.Steps = steps

	// Return the data source test
	return dataSourceTest
}

// Test Scenarios //

// NewDefaultScenario creates a test case with initial and updated configurations.
func (gdst *WorkloadDataSourceTest) NewDefaultScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "new"
	gvcName := "default-gvc"
	name := "httpbin-example"

	// Build test steps
	_, initialStep := gdst.BuildDefaultTestStep(dataSourceName, gvcName, name)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildDefaultTestStep returns a default initial test step and its associated test case for the data source.
func (gdst *WorkloadDataSourceTest) BuildDefaultTestStep(dataSourceName string, gvcName string, name string) (WorkloadDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "workload",
			ResourceName:    dataSourceName,
			Name:            name,
			GvcName:         gvcName,
			Description:     name,
			ResourceAddress: fmt.Sprintf("data.cpln_workload.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: gdst.DefaultHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "1"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", c.GvcName),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "standard"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", c.GvcName, name)),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "false"),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "httpbin",
					"image":             "kennethreitz/httpbin",
					"working_directory": "/usr",
					"memory":            "128Mi",
					"cpu":               "50m",
					"ports": []map[string]interface{}{
						{
							"number":   "80",
							"protocol": "http",
						},
					},
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
						"env-name-04": "null",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path":         "/metrics",
							"port":         "8181",
							"drop_metrics": []string{"go_.*", "process_.*", ".*_bucket|.*_sum|.*_count"},
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"grpc": []map[string]interface{}{
								{
									"port": "3000",
								},
							},
							"initial_delay_seconds": "1",
							"period_seconds":        "11",
							"timeout_seconds":       "2",
							"success_threshold":     "2",
							"failure_threshold":     "4",
						},
					},
					"liveness_probe": []map[string]interface{}{
						{
							"http_get": []map[string]interface{}{
								{
									"path":   "/path",
									"port":   "8282",
									"scheme": "HTTPS",
									"http_headers": map[string]interface{}{
										"header-name-01": "header-value-01",
										"header-name-02": "header-value-02",
									},
								},
							},
							"initial_delay_seconds": "2",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"volume": []map[string]interface{}{
						{
							"uri":             "s3://bucket",
							"recovery_policy": "retain",
							"path":            "/testpath01",
						},
						{
							"uri":             "azureblob://storageAccount/container",
							"recovery_policy": "retain",
							"path":            "/testpath02",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("firewall_spec", []map[string]interface{}{
				{
					"external": []map[string]interface{}{
						{
							"inbound_allow_cidr":      []string{"0.0.0.0/0"},
							"inbound_blocked_cidr":    []string{},
							"outbound_allow_hostname": []string{},
							"outbound_allow_cidr":     []string{"0.0.0.0/0"},
							"outbound_blocked_cidr":   []string{},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "none",
							"inbound_allow_workload": []string{},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("options", []map[string]interface{}{
				{
					"timeout_seconds": "30",
					"capacity_ai":     "false",
					"debug":           "true",
					"suspend":         "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "cpu",
							"target":              "100",
							"max_scale":           "1",
							"min_scale":           "1",
							"max_concurrency":     "0",
							"scale_to_zero_delay": "300",
							"keda": []map[string]interface{}{
								{
									"polling_interval":        "1",
									"cooldown_period":         "1",
									"initial_cooldown_period": "1",
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "10",
					"scaling_policy":                   "Parallel",
					"termination_grace_period_seconds": "10",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
					"direct": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("request_retry_policy", []map[string]interface{}{
				{
					"attempts": "3",
					"retry_on": []string{"connect-failure", "refused-stream", "unavailable"},
				},
			}),
		),
	}
}

// Configs //

// DefaultHcl returns a data source HCL.
func (gdst *WorkloadDataSourceTest) DefaultHcl(c WorkloadDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_workload" "%s" {
  name = "%s"
	gvc  = "%s"
}
`, c.ResourceName, c.Name, c.GvcName)
}

/*** Data Source Test Case ***/

// WorkloadDataSourceTestCase defines a specific data source test case.
type WorkloadDataSourceTestCase struct {
	ProviderTestCase
}
