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

// TestAccControlPlaneWorkload_basic performs an acceptance test for the resource.
func TestAccControlPlaneWorkload_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewWorkloadResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "WORKLOAD") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// WorkloadResourceTest defines the necessary functionality to test the resource.
type WorkloadResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
	GvcConfig  string
	GvcCase    GvcResourceTestCase
}

// NewWorkloadResourceTest creates a WorkloadResourceTest with initialized test cases.
func NewWorkloadResourceTest() WorkloadResourceTest {
	// Generate a unique random identifier
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	// Generate the GVC name
	gvcName := fmt.Sprintf("gvc-%s", random)

	// Create resource test instances
	gvcResourceTest := GvcResourceTest{}

	// Create a GVC case
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

	// Initialize the gvc config
	gvcConfig := gvcResourceTest.GvcRequiredOnly(gvcCase)

	// Create a resource test instance
	resourceTest := WorkloadResourceTest{
		RandomName: random,
		GvcConfig:  gvcConfig,
		GvcCase:    gvcCase,
	}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, resourceTest.NewK8sVolumeUriScenario()...)
	steps = append(steps, resourceTest.NewServerlessScenario()...)
	steps = append(steps, resourceTest.NewStandardScenario()...)
	steps = append(steps, resourceTest.NewCronScenario()...)
	steps = append(steps, resourceTest.NewStatefulScenario()...)
	steps = append(steps, resourceTest.NewVmScenario()...)
	steps = append(steps, resourceTest.NewVmDefaultsScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (wrt *WorkloadResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_workload resources. Total resources: %d", len(s.RootModule().Resources)))

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
	tflog.Info(TestLoggerContext, "All cpln_workload resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewK8sVolumeUriScenario defines a scenario verifying a k8s://secret volume uri is accepted and persisted to state.
func (wrt *WorkloadResourceTest) NewK8sVolumeUriScenario() []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("workload-k8s-volume-%s", wrt.RandomName)

	// Build the test step
	_, initialStep := wrt.BuildK8sVolumeUriTestStep(name)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
	}
}

// NewServerlessScenario defines a full serverless workload lifecycle test case including creation, updates, import, and state restoration.
func (wrt *WorkloadResourceTest) NewServerlessScenario() []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("workload-serverless-%s", wrt.RandomName)

	// Build test steps
	initialConfig, initialStep := wrt.BuildServerlessTestStep(name)
	caseUpdate1 := wrt.BuildServerlessUpdate1TestStep(initialConfig.ProviderTestCase)
	caseUpdate2 := wrt.BuildServerlessUpdate2TestStep(initialConfig.ProviderTestCase)
	caseUpdate3 := wrt.BuildServerlessUpdate3TestStep(initialConfig.ProviderTestCase)
	caseUpdate4 := wrt.BuildServerlessUpdate4TestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName:  initialConfig.ResourceAddress,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:%s", wrt.GvcCase.Name, name),
		},
		// Update & Read
		caseUpdate1,
		caseUpdate2,
		caseUpdate3,
		caseUpdate4,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewStandardScenario defines a standard workload test case with creation and import verification only.
func (wrt *WorkloadResourceTest) NewStandardScenario() []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("workload-standard-%s", wrt.RandomName)

	// Build test steps
	initialConfig, initialStep := wrt.BuildStandardTestStep(name)
	caseUpdate1 := wrt.BuildStandardUpdate1TestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName:  initialConfig.ResourceAddress,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:%s", wrt.GvcCase.Name, name),
		},
		// Update & Read
		caseUpdate1,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewCronScenario defines a cron workload test case including creation, an update, import, and state restoration.
func (wrt *WorkloadResourceTest) NewCronScenario() []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("workload-cron-%s", wrt.RandomName)

	// Build test steps
	initialConfig, initialStep := wrt.BuildCronTestStep(name)
	caseUpdate1 := wrt.BuildCronUpdate1TestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName:  initialConfig.ResourceAddress,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:%s", wrt.GvcCase.Name, name),
		},
		// Update & Read
		caseUpdate1,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewStatefulScenario defines a stateful workload test case with creation and import validation.
func (wrt *WorkloadResourceTest) NewStatefulScenario() []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("workload-stateful-%s", wrt.RandomName)

	// Build test steps
	initialConfig, initialStep := wrt.BuildStatefulTestStep(name)
	caseUpdate1 := wrt.BuildStatefulUpdate1TestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps
	return []resource.TestStep{
		// Create & Read
		initialStep,
		// Import State
		{
			ResourceName:  initialConfig.ResourceAddress,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:%s", wrt.GvcCase.Name, name),
		},
		// Update & Read
		caseUpdate1,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewVmScenario defines a vm workload test case including creation, an update, import, and state restoration.
func (wrt *WorkloadResourceTest) NewVmScenario() []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("workload-vm-%s", wrt.RandomName)

	// Build test steps
	initialConfig, initialStep := wrt.BuildVmTestStep(name)
	caseUpdate1 := wrt.BuildVmUpdate1TestStep(initialConfig.ProviderTestCase)
	caseUpdate2 := wrt.BuildVmUpdate2TestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps. The slice walks the VM through its full lifecycle:
	// absent optional blocks -> add cloud-init/access-credential -> grow collections and flip
	// the boot-source/cloud-init xor branches -> shrink back -> remove the optional blocks again.
	// Re-using initialStep and caseUpdate1 on the way down exercises the build/flatten round-trip
	// at every cardinality without drift.
	return []resource.TestStep{
		// Create & Read (oci boot source, no cloud-init/access-credential)
		initialStep,
		// Import State
		{
			ResourceName:  initialConfig.ResourceAddress,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("%s:%s", wrt.GvcCase.Name, name),
		},
		// Update & Read (add cloud-init + a single access-credential)
		caseUpdate1,
		// Update & Read (flip to http boot source + base64 cloud-init, grow collections to two entries)
		caseUpdate2,
		// Update & Read (shrink back to the single-entry oci configuration)
		caseUpdate1,
		// Revert the resource to its initial state (remove the optional blocks entirely)
		initialStep,
	}
}

// NewVmDefaultsScenario verifies the vm firmware, network, and clock object/list defaults materialize without drift when omitted.
func (wrt *WorkloadResourceTest) NewVmDefaultsScenario() []resource.TestStep {
	// Generate a unique name for the resources
	name := fmt.Sprintf("workload-vm-defaults-%s", wrt.RandomName)

	// Build test steps
	initialConfig, omittedStep := wrt.BuildVmDefaultsOmittedTestStep(name)
	setStep := wrt.BuildVmDefaultsSetTestStep(initialConfig.ProviderTestCase)

	// Return the complete test steps. Walk the optional defaulted blocks through omit -> set -> omit to
	// prove the defaults apply on create, an explicit value overrides them, and removing the blocks
	// re-applies the defaults (the Computed+Default revert path).
	return []resource.TestStep{
		// Omit firmware/network/clock -> defaults materialize
		omittedStep,
		// Set them explicitly -> overrides
		setStep,
		// Remove them again -> defaults re-apply
		omittedStep,
	}
}

// Test Cases //

// BuildK8sVolumeUriTestStep constructs a workload test step that mounts a k8s://secret volume and asserts its uri and path persist to state.
func (wrt *WorkloadResourceTest) BuildK8sVolumeUriTestStep(name string) (WorkloadResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "workload",
			ResourceName:    "new",
			ResourceAddress: "cpln_workload.new",
			Name:            name,
			GvcName:         wrt.GvcCase.Name,
			Description:     name,
		},
	}

	// Initialize and return the test step. A k8s://secret volume only mounts on byok clusters, but
	// the spec is accepted at create just like an s3:// volume pointing at a missing bucket, so
	// applying it and asserting the persisted uri/path is safe (create never waits on the deployment).
	return c, resource.TestStep{
		Config: wrt.K8sVolumeUriHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "serverless"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "false"),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":        "container-01",
					"image":       "gcr.io/knative-samples/helloworld-go",
					"cpu":         "50m",
					"memory":      "128Mi",
					"inherit_env": "false",
					"ports": []map[string]interface{}{
						{
							"number":   "8080",
							"protocol": "http",
						},
					},
					"volume": []map[string]interface{}{
						{
							"uri":             "k8s://secret/example-secret",
							"recovery_policy": "retain",
							"path":            "/k8s-secret",
						},
					},
				},
			}),
		),
	}
}

// BuildServerlessTestStep constructs the initial serverless workload test configuration and test step.
func (wrt *WorkloadResourceTest) BuildServerlessTestStep(name string) (WorkloadResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "workload",
			ResourceName:      "new",
			ResourceAddress:   "cpln_workload.new",
			Name:              name,
			GvcName:           wrt.GvcCase.Name,
			Description:       name,
			DescriptionUpdate: "workload default description updated",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: wrt.ServerlessRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "0"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "serverless"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "false"),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":        "container-01",
					"image":       "gcr.io/knative-samples/helloworld-go",
					"cpu":         "50m",
					"memory":      "128Mi",
					"inherit_env": "false",
					"ports": []map[string]interface{}{
						{
							"number":   "8080",
							"protocol": "http",
						},
					},
				},
			}),
		),
	}
}

// BuildServerlessUpdate1TestStep returns the first update test step for a serverless workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildServerlessUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
		Envoy:            `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras:           `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: wrt.ServerlessUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "serverless"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"port":              "8080",
					"memory":            "128Mi",
					"min_cpu":           "25m",
					"min_memory":        "32Mi",
					"cpu":               "50m",
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
						"env-name-04": "",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket":            []map[string]interface{}{{}},
							"initial_delay_seconds": "10",
							"period_seconds":        "10",
							"timeout_seconds":       "1",
							"success_threshold":     "1",
							"failure_threshold":     "3",
						},
					},
					"liveness_probe": []map[string]interface{}{
						{
							"http_get": []map[string]interface{}{
								{
									"path":   "/path",
									"port":   "8080",
									"scheme": "HTTPS",
								},
							},
							"initial_delay_seconds": "60",
							"period_seconds":        "10",
							"timeout_seconds":       "1",
							"success_threshold":     "1",
							"failure_threshold":     "3",
						},
					},
					"lifecycle": []map[string]interface{}{{}},
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
			c.TestCheckNestedBlocks("firewall_spec", []map[string]interface{}{{}}),
			c.TestCheckNestedBlocks("options", []map[string]interface{}{
				{
					"timeout_seconds": "5",
					"capacity_ai":     "true",
					"debug":           "false",
					"suspend":         "false",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "95",
							"max_scale":           "5",
							"min_scale":           "1",
							"max_concurrency":     "0",
							"scale_to_zero_delay": "300",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("local_options", []map[string]interface{}{
				{
					"location":        "aws-eu-central-1",
					"timeout_seconds": "5",
					"capacity_ai":     "true",
					"debug":           "false",
					"suspend":         "false",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "95",
							"max_scale":           "5",
							"min_scale":           "1",
							"max_concurrency":     "0",
							"scale_to_zero_delay": "300",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "0",
					"scaling_policy":                   "OrderedReady",
					"termination_grace_period_seconds": "90",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{{}}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
				},
			}),
			c.TestCheckNestedBlocks("request_retry_policy", []map[string]interface{}{
				{
					"attempts": "2",
					"retry_on": []string{"connect-failure", "refused-stream", "unavailable", "cancelled", "resource-exhausted", "retriable-status-codes"},
				},
			}),
		),
	}
}

// BuildServerlessUpdate2TestStep returns the second update test step for a serverless workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildServerlessUpdate2TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
		Envoy:            `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras:           `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: wrt.ServerlessUpdate2Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "serverless"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"port":              "8080",
					"memory":            "7Gi",
					"cpu":               "2",
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
					"gpu_nvidia": []map[string]interface{}{
						{
							"model":    "t4",
							"quantity": "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{{}},
							"pre_stop":   []map[string]interface{}{{}},
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
							"inbound_allow_cidr":      []string{},
							"inbound_blocked_cidr":    []string{},
							"outbound_allow_hostname": []string{},
							"outbound_allow_cidr":     []string{},
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
							"metric":              "concurrency",
							"target":              "100",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("local_options", []map[string]interface{}{
				{
					"location":        "aws-eu-central-1",
					"timeout_seconds": "30",
					"capacity_ai":     "false",
					"debug":           "true",
					"suspend":         "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "100",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
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
					"run_as_user":          "1000",
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
					"geo_location": []map[string]interface{}{
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

// BuildServerlessUpdate3TestStep returns the third update test step for a serverless workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildServerlessUpdate3TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
		Envoy:            `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras:           `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: wrt.ServerlessUpdate3Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "serverless"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"port":              "8080",
					"memory":            "128Mi",
					"cpu":               "50m",
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
						"env-name-04": "",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path":         "/metrics",
							"port":         "8181",
							"drop_metrics": []string{"go_.*", ".*_bucket|.*_sum|.*_count"},
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket": []map[string]interface{}{
								{
									"port": "3200",
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
							"exec": []map[string]interface{}{
								{
									"command": []string{"command-01", "command-02"},
								},
							},
							"initial_delay_seconds": "2",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource": "amd.com/gpu",
							"quantity": "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{{}},
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
					"timeout_seconds":            "30",
					"capacity_ai":                "false",
					"capacity_ai_update_minutes": "5",
					"debug":                      "true",
					"suspend":                    "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "100",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("local_options", []map[string]interface{}{
				{
					"location":        "aws-eu-central-1",
					"timeout_seconds": "30",
					"capacity_ai":     "false",
					"debug":           "true",
					"suspend":         "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "100",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
				{
					"location":        "aws-us-west-2",
					"timeout_seconds": "30",
					"capacity_ai":     "false",
					"debug":           "true",
					"suspend":         "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "90",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "20",
					"max_unavailable_replicas":         "10",
					"scaling_policy":                   "OrderedReady",
					"termination_grace_period_seconds": "20",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
					"direct": []map[string]interface{}{
						{
							"enabled": "false",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildServerlessUpdate4TestStep returns the fourth update test step for a serverless workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildServerlessUpdate4TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
		Envoy:            `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras:           `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: wrt.ServerlessUpdate4Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "serverless"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", fmt.Sprintf("//gvc/%s/identity/identity-%s", wrt.GvcCase.Name, wrt.RandomName)),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"port":              "8080",
					"memory":            "128Mi",
					"cpu":               "50m",
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
						"env-name-04": "env-value-04",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket": []map[string]interface{}{
								{
									"port": "3200",
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
							"exec": []map[string]interface{}{
								{
									"command": []string{"command-01", "command-02"},
								},
							},
							"initial_delay_seconds": "2",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource":      "amd.com/gpu",
							"runtime_class": "amd",
							"quantity":      "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "Allow-Header",
											"allowed_values": []string{"reg", "req2"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "workload-list",
							"inbound_allow_workload": []string{"//gvc/new/workload/non-existing-workload-01", "/org/terraform-test-org/gvc/new/workload/non-existing-workload-02"},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("options", []map[string]interface{}{
				{
					"timeout_seconds":            "30",
					"capacity_ai":                "false",
					"capacity_ai_update_minutes": "10",
					"debug":                      "true",
					"suspend":                    "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "100",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("local_options", []map[string]interface{}{
				{
					"location":        "aws-eu-central-1",
					"timeout_seconds": "30",
					"capacity_ai":     "false",
					"debug":           "true",
					"suspend":         "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "100",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
				{
					"location":        "aws-us-west-2",
					"timeout_seconds": "30",
					"capacity_ai":     "false",
					"debug":           "true",
					"suspend":         "true",
					"autoscaling": []map[string]interface{}{
						{
							"metric":              "concurrency",
							"target":              "90",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "30",
					"max_unavailable_replicas":         "20",
					"max_surge_replicas":               "30",
					"scaling_policy":                   "Parallel",
					"termination_grace_period_seconds": "40",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
					"direct": []map[string]interface{}{
						{
							"enabled": "false",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildStandardTestStep constructs the initial test configuration and test step for a standard workload.
func (wrt *WorkloadResourceTest) BuildStandardTestStep(name string) (WorkloadResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "workload",
			ResourceName:      "new",
			ResourceAddress:   "cpln_workload.new",
			Name:              name,
			GvcName:           wrt.GvcCase.Name,
			Description:       name,
			DescriptionUpdate: "workload default description updated",
		},
		Envoy:  `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras: `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: wrt.StandardHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "standard"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"memory":            "128Mi",
					"cpu":               "50m",
					"ports": []map[string]interface{}{
						{
							"protocol": "http2",
							"number":   "8080",
						},
						{
							"protocol": "grpc",
							"number":   "3000",
						},
						{
							"protocol": "tcp",
							"number":   "3001",
						},
					},
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket": []map[string]interface{}{
								{
									"port": "3200",
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
							"exec": []map[string]interface{}{
								{
									"command": []string{"command-01", "command-02"},
								},
							},
							"initial_delay_seconds": "2",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource":      "amd.com/gpu",
							"runtime_class": "amd",
							"quantity":      "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "Allow-Header",
											"blocked_values": []string{"blocked", "blocked2"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "same-org",
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
							"metric_percentile":   "p50",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
							"multi": []map[string]interface{}{
								{
									"metric": "cpu",
									"target": 50,
								},
								{
									"metric": "memory",
									"target": 50,
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "0",
					"scaling_policy":                   "OrderedReady",
					"termination_grace_period_seconds": "90",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
					"direct": []map[string]interface{}{
						{
							"enabled": "false",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildStandardUpdate1TestStep returns the first update test step for a standard workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildStandardUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
		Envoy:            `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras:           `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: wrt.StandardUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "standard"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"memory":            "128Mi",
					"cpu":               "50m",
					"ports": []map[string]interface{}{
						{
							"protocol": "http2",
							"number":   "8080",
						},
						{
							"protocol": "grpc",
							"number":   "3000",
						},
						{
							"protocol": "tcp",
							"number":   "3001",
						},
					},
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket": []map[string]interface{}{
								{
									"port": "3200",
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
							"exec": []map[string]interface{}{
								{
									"command": []string{"command-01", "command-02"},
								},
							},
							"initial_delay_seconds": "2",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource":      "amd.com/gpu",
							"runtime_class": "amd",
							"quantity":      "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "Allow-Header",
											"blocked_values": []string{"blocked", "blocked2"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "same-org",
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
							"metric":              "keda",
							"metric_percentile":   "p50",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
							"keda": []map[string]interface{}{
								{
									"polling_interval":        "30",
									"cooldown_period":         "60",
									"initial_cooldown_period": "10",
									"trigger": []map[string]interface{}{
										{
											"type":               "cpu",
											"name":               "cpu-trigger-01",
											"use_cached_metrics": "true",
											"metric_type":        "Utilization",
											"metadata": map[string]interface{}{
												"type":  "Utilization",
												"value": "50",
											},
										},
										{
											"type":               "rabbitmq",
											"name":               "rabbitmq-trigger",
											"use_cached_metrics": "false",
											"metric_type":        "AverageValue",
											"metadata": map[string]interface{}{
												"host":        "amqp://user:pass@rabbitmq:5672/",
												"queueName":   "jobs",
												"queueLength": "30",
											},
											"authentication_ref": []map[string]interface{}{
												{
													"name": "rabbitmq-auth",
												},
											},
										},
									},
									"advanced": []map[string]interface{}{
										{
											"scaling_modifiers": []map[string]interface{}{
												{
													"target":            "5",
													"activation_target": "1",
													"metric_type":       "Value",
													"formula":           "m * 2",
												},
											},
										},
									},
									"fallback": []map[string]interface{}{
										{
											"failure_threshold": "3",
											"replicas":          "1",
											"behavior":          "static",
										},
									},
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "2",
					"max_unavailable_replicas":         "10",
					"max_surge_replicas":               "20",
					"scaling_policy":                   "Parallel",
					"termination_grace_period_seconds": "10",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
					"direct": []map[string]interface{}{
						{
							"enabled": "false",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildCronTestStep constructs the initial test configuration and test step for a cron workload.
func (wrt *WorkloadResourceTest) BuildCronTestStep(name string) (WorkloadResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "workload",
			ResourceName:      "new",
			ResourceAddress:   "cpln_workload.new",
			Name:              name,
			GvcName:           wrt.GvcCase.Name,
			Description:       name,
			DescriptionUpdate: "workload default description updated",
		},
		Envoy:  `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras: `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: wrt.CronHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "cron"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"memory":            "128Mi",
					"cpu":               "50m",
					"ports": []map[string]interface{}{
						{
							"protocol": "http2",
							"number":   "8080",
						},
						{
							"protocol": "grpc",
							"number":   "3000",
						},
						{
							"protocol": "tcp",
							"number":   "3001",
						},
					},
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource":      "amd.com/gpu",
							"runtime_class": "amd",
							"quantity":      "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "Allow-Header",
											"blocked_values": []string{"blocked", "blocked2"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "same-org",
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
							"metric":              "concurrency",
							"target":              "95",
							"max_scale":           "1",
							"min_scale":           "1",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("job", []map[string]interface{}{
				{
					"schedule":           "* * * * *",
					"concurrency_policy": "Forbid",
					"history_limit":      "5",
					"restart_policy":     "Never",
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "0",
					"scaling_policy":                   "OrderedReady",
					"termination_grace_period_seconds": "90",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
					"direct": []map[string]interface{}{
						{
							"enabled": "false",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildCronUpdate1TestStep returns the first update test step for a cron workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildCronUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
		Envoy:            `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras:           `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: wrt.CronUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "cron"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"memory":            "128Mi",
					"cpu":               "50m",
					"ports": []map[string]interface{}{
						{
							"protocol": "http2",
							"number":   "8080",
						},
						{
							"protocol": "grpc",
							"number":   "3000",
						},
						{
							"protocol": "tcp",
							"number":   "3001",
						},
					},
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource":      "amd.com/gpu",
							"runtime_class": "amd",
							"quantity":      "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "Allow-Header",
											"blocked_values": []string{"blocked", "blocked2"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "same-org",
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
							"metric":              "concurrency",
							"target":              "95",
							"max_scale":           "1",
							"min_scale":           "1",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("job", []map[string]interface{}{
				{
					"schedule":                "* * * * *",
					"concurrency_policy":      "Forbid",
					"history_limit":           "5",
					"restart_policy":          "Never",
					"active_deadline_seconds": "1200",
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "2",
					"max_unavailable_replicas":         "10",
					"max_surge_replicas":               "20",
					"scaling_policy":                   "Parallel",
					"termination_grace_period_seconds": "10",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "false",
					"direct": []map[string]interface{}{
						{
							"enabled": "false",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildStatefulTestStep constructs the initial test configuration and test step for a stateful workload.
func (wrt *WorkloadResourceTest) BuildStatefulTestStep(name string) (WorkloadResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "workload",
			ResourceName:      "new",
			ResourceAddress:   "cpln_workload.new",
			Name:              name,
			GvcName:           wrt.GvcCase.Name,
			Description:       name,
			DescriptionUpdate: "workload default description updated",
		},
		Envoy:  `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras: `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: wrt.StatefulHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "stateful"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"memory":            "128Mi",
					"cpu":               "50m",
					"ports": []map[string]interface{}{
						{
							"protocol": "http2",
							"number":   "8080",
						},
						{
							"protocol": "grpc",
							"number":   "3000",
						},
						{
							"protocol": "tcp",
							"number":   "3001",
						},
					},
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket": []map[string]interface{}{
								{
									"port": "3200",
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
							"exec": []map[string]interface{}{
								{
									"command": []string{"command-01", "command-02"},
								},
							},
							"initial_delay_seconds": "2",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource":      "amd.com/gpu",
							"runtime_class": "amd",
							"quantity":      "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "Allow-Header",
											"blocked_values": []string{"blocked", "blocked2"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "same-org",
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
							"metric_percentile":   "p50",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
							"multi": []map[string]interface{}{
								{
									"metric": "cpu",
									"target": 50,
								},
								{
									"metric": "memory",
									"target": 50,
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "0",
					"scaling_policy":                   "OrderedReady",
					"termination_grace_period_seconds": "90",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "true",
					"direct": []map[string]interface{}{
						{
							"enabled": "true",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildStatefulUpdate1TestStep returns the first update test step for a stateful workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildStatefulUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
		Envoy:            `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`,
		Extras:           `{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"cpln.io/nodeType","operator":"In","values":["tasks"]}]}]}}}}`,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: wrt.StatefulUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "stateful"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(c.Extras)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":              "container-01",
					"image":             "gcr.io/knative-samples/helloworld-go",
					"working_directory": "/usr",
					"memory":            "128Mi",
					"cpu":               "50m",
					"ports": []map[string]interface{}{
						{
							"protocol": "http2",
							"number":   "8080",
						},
						{
							"protocol": "grpc",
							"number":   "3000",
						},
						{
							"protocol": "tcp",
							"number":   "3001",
						},
					},
					"env": map[string]interface{}{
						"env-name-01": "env-value-01",
						"env-name-02": "env-value-02",
						"env-name-03": "env-value-03",
					},
					"inherit_env": "true",
					"command":     "override-command",
					"args":        []string{"arg-01", "arg-02", "arg-03"},
					"metrics": []map[string]interface{}{
						{
							"path": "/metrics",
							"port": "8181",
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket": []map[string]interface{}{
								{
									"port": "3200",
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
							"exec": []map[string]interface{}{
								{
									"command": []string{"command-01", "command-02"},
								},
							},
							"initial_delay_seconds": "2",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"gpu_custom": []map[string]interface{}{
						{
							"resource":      "amd.com/gpu",
							"runtime_class": "amd",
							"quantity":      "1",
						},
					},
					"lifecycle": []map[string]interface{}{
						{
							"post_start": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"exec": []map[string]interface{}{
										{
											"command": []string{"command-01", "command-02", "command-03"},
										},
									},
								},
							},
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
							"inbound_blocked_cidr":    []string{"127.0.0.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com", "*.cpln.io"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"::1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "http",
									"number":   "80",
								},
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "Allow-Header",
											"blocked_values": []string{"blocked", "blocked2"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "same-org",
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
							"metric":              "keda",
							"metric_percentile":   "p50",
							"max_scale":           "3",
							"min_scale":           "2",
							"max_concurrency":     "500",
							"scale_to_zero_delay": "400",
							"keda": []map[string]interface{}{
								{
									"trigger": []map[string]interface{}{
										{
											"type":        "cpu",
											"name":        "cpu-trigger-01",
											"metric_type": "Utilization",
											"metadata": map[string]interface{}{
												"type":  "Utilization",
												"value": "50",
											},
										},
									},
									"advanced": []map[string]interface{}{
										{
											"scaling_modifiers": []map[string]interface{}{{}},
										},
									},
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("sidecar", []map[string]interface{}{
				{
					"envoy": CanonicalizeEnvoyJSON(c.Envoy),
				},
			}),
			c.TestCheckNestedBlocks("rollout_options", []map[string]interface{}{
				{
					"min_ready_seconds":                "2",
					"max_surge_replicas":               "20",
					"scaling_policy":                   "Parallel",
					"termination_grace_period_seconds": "10",
				},
			}),
			c.TestCheckNestedBlocks("security_options", []map[string]interface{}{
				{
					"file_system_group_id": "1",
					"run_as_user":          "1000",
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"replica_direct": "true",
					"direct": []map[string]interface{}{
						{
							"enabled": "true",
							"ipset":   "my-ipset-01",
							"port": []map[string]interface{}{
								{
									"external_port":  "22",
									"protocol":       "TCP",
									"scheme":         "http",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "198.51.100.0/24",
									"city":    "Los Angeles",
									"country": "USA",
									"region":  "North America",
								},
							},
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

// BuildVmTestStep constructs the initial test configuration and test step for a vm workload.
func (wrt *WorkloadResourceTest) BuildVmTestStep(name string) (WorkloadResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "workload",
			ResourceName:      "new",
			ResourceAddress:   "cpln_workload.new",
			Name:              name,
			GvcName:           wrt.GvcCase.Name,
			Description:       name,
			DescriptionUpdate: "workload default description updated",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: wrt.VmHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "vm"),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":   "vm-container",
					"cpu":    "2000m",
					"memory": "2Gi",
					"volume": []map[string]interface{}{
						{
							"uri":        fmt.Sprintf("cpln://volumeset/vmdata-%s", wrt.RandomName),
							"name":       "data-disk",
							"bus":        "virtio",
							"boot_order": "2",
						},
					},
				},
			}),
			c.TestCheckObjectAttr("vm", map[string]interface{}{
				"guest_os":     "linux",
				"run_strategy": "Always",
				"hostname":     "vm-host",
				"boot_disk": map[string]interface{}{
					"bus":        "virtio",
					"boot_order": "1",
					"source": map[string]interface{}{
						"oci": map[string]interface{}{
							"image": "quay.io/containerdisks/ubuntu:22.04",
						},
					},
					"persist": map[string]interface{}{
						"volume_set": fmt.Sprintf("cpln://volumeset/vmboot-%s", wrt.RandomName),
					},
				},
				"cpu": map[string]interface{}{
					"sockets": "2",
					"threads": "1",
				},
				"firmware": map[string]interface{}{
					"bootloader":  "efi",
					"secure_boot": "false",
					"uuid":        "5d8e7a3c-1f2b-4c6d-8e9f-0a1b2c3d4e5f",
					"smbios": map[string]interface{}{
						"manufacturer": "ControlPlane",
						"product":      "cpln-vm",
					},
				},
				"network": []map[string]interface{}{
					{
						"name": "default",
					},
				},
				"clock": map[string]interface{}{
					"timezone": "UTC",
				},
			}),
		),
	}
}

// BuildVmUpdate1TestStep returns the first update test step for a vm workload based on the initial test case.
func (wrt *WorkloadResourceTest) BuildVmUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the update test step
	return resource.TestStep{
		Config: wrt.VmUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "vm"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":   "vm-container",
					"cpu":    "4000m",
					"memory": "4Gi",
					"volume": []map[string]interface{}{
						{
							"uri":        fmt.Sprintf("cpln://volumeset/vmdata-%s", wrt.RandomName),
							"name":       "data-disk",
							"bus":        "scsi",
							"boot_order": "2",
						},
					},
				},
			}),
			c.TestCheckObjectAttr("vm", map[string]interface{}{
				"guest_os":     "linux",
				"run_strategy": "Manual",
				"hostname":     "vm-host-updated",
				"subdomain":    "vms",
				"boot_disk": map[string]interface{}{
					"bus":        "virtio",
					"boot_order": "1",
					"source": map[string]interface{}{
						"oci": map[string]interface{}{
							"image": "quay.io/containerdisks/ubuntu:24.04",
						},
					},
					"persist": map[string]interface{}{
						"volume_set": fmt.Sprintf("cpln://volumeset/vmboot-%s", wrt.RandomName),
					},
				},
				"cpu": map[string]interface{}{
					"sockets": "4",
					"threads": "2",
				},
				"firmware": map[string]interface{}{
					"bootloader":  "efi",
					"secure_boot": "false",
					"uuid":        "5d8e7a3c-1f2b-4c6d-8e9f-0a1b2c3d4e5f",
					"serial":      "vm-serial-01",
					"smbios": map[string]interface{}{
						"manufacturer": "ControlPlane",
						"product":      "cpln-vm",
						"version":      "2.0",
						"sku":          "sku-01",
						"family":       "cpln",
					},
				},
				"network": []map[string]interface{}{
					{
						"name": "default",
					},
				},
				"cloud_init": map[string]interface{}{
					"user_data":              "#cloud-config\nruncmd:\n  - echo hello\n",
					"ssh_public_key_secrets": []string{GetSelfLink(OrgName, "secret", fmt.Sprintf("vm-ssh-%s", wrt.RandomName))},
				},
				"access_credential": []map[string]interface{}{
					{
						"ssh_public_key_secret": GetSelfLink(OrgName, "secret", fmt.Sprintf("vm-ssh-%s", wrt.RandomName)),
						"users":                 []string{"root", "ubuntu"},
						"delivery_method":       "qemuGuestAgent",
					},
				},
				"clock": map[string]interface{}{
					"timezone": "America/New_York",
				},
			}),
		),
	}
}

// BuildVmUpdate2TestStep returns the second update test step for a vm workload, flipping the boot source to
// http and cloud-init to base64, growing the access-credential and ssh-key collections to two entries, and
// exercising the remaining enum values (sata bus, bios bootloader, configDrive delivery, RerunOnFailure).
func (wrt *WorkloadResourceTest) BuildVmUpdate2TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the update test step
	return resource.TestStep{
		Config: wrt.VmUpdate2Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.DescriptionUpdate, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "vm"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "identity_link", GetSelfLinkWithGvc(OrgName, "identity", wrt.GvcCase.Name, fmt.Sprintf("identity-%s", wrt.RandomName))),
			resource.TestCheckResourceAttr(c.ResourceAddress, "support_dynamic_tags", "true"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "extras", CanonicalizeEnvoyJSON(`{"tolerations":[{"key":"cpln.io/nodeType","operator":"Equal","value":"vm","effect":"NoSchedule"}]}`)),
			c.TestCheckNestedBlocks("container", []map[string]interface{}{
				{
					"name":        "vm-container",
					"cpu":         "4000m",
					"memory":      "4Gi",
					"inherit_env": "true",
					"ports": []map[string]interface{}{
						{
							"protocol": "tcp",
							"number":   "8080",
						},
						{
							"protocol": "http",
							"number":   "80",
						},
					},
					"env": map[string]interface{}{
						"ENV_KEY": "env-value",
						"APP_ENV": "production",
					},
					"metrics": []map[string]interface{}{
						{
							"port":         "8181",
							"path":         "/metrics",
							"drop_metrics": []string{"envoy_.*", "go_gc_.*"},
						},
					},
					"readiness_probe": []map[string]interface{}{
						{
							"tcp_socket": []map[string]interface{}{
								{
									"port": "8080",
								},
							},
							"initial_delay_seconds": "5",
							"period_seconds":        "10",
							"timeout_seconds":       "2",
							"success_threshold":     "1",
							"failure_threshold":     "3",
						},
					},
					"liveness_probe": []map[string]interface{}{
						{
							"http_get": []map[string]interface{}{
								{
									"path":   "/healthz",
									"port":   "80",
									"scheme": "HTTP",
									"http_headers": map[string]interface{}{
										"X-Custom-Header": "custom-value",
									},
								},
							},
							"initial_delay_seconds": "10",
							"period_seconds":        "15",
							"timeout_seconds":       "3",
							"success_threshold":     "1",
							"failure_threshold":     "5",
						},
					},
					"volume": []map[string]interface{}{
						{
							"uri":        fmt.Sprintf("cpln://volumeset/vmdata-%s", wrt.RandomName),
							"name":       "data-disk",
							"bus":        "virtio",
							"boot_order": "2",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("firewall_spec", []map[string]interface{}{
				{
					"external": []map[string]interface{}{
						{
							"inbound_allow_cidr":      []string{"0.0.0.0/0"},
							"inbound_blocked_cidr":    []string{"192.0.2.1"},
							"outbound_allow_hostname": []string{"*.controlplane.com"},
							"outbound_allow_cidr":     []string{},
							"outbound_blocked_cidr":   []string{"198.51.100.1"},
							"outbound_allow_port": []map[string]interface{}{
								{
									"protocol": "https",
									"number":   "443",
								},
							},
							"http": []map[string]interface{}{
								{
									"inbound_header_filter": []map[string]interface{}{
										{
											"key":            "X-Allowed",
											"allowed_values": []string{"^v1$", "^v2$"},
										},
									},
								},
							},
						},
					},
					"internal": []map[string]interface{}{
						{
							"inbound_allow_type":     "same-gvc",
							"inbound_allow_workload": []string{},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("load_balancer", []map[string]interface{}{
				{
					"direct": []map[string]interface{}{
						{
							"enabled": "true",
							"ipset":   "vm-ipset",
							"port": []map[string]interface{}{
								{
									"external_port":  "8080",
									"protocol":       "TCP",
									"scheme":         "tcp",
									"container_port": "80",
								},
							},
						},
					},
					"geo_location": []map[string]interface{}{
						{
							"enabled": "true",
							"headers": []map[string]interface{}{
								{
									"asn":     "x-geo-asn",
									"city":    "x-geo-city",
									"country": "x-geo-country",
									"region":  "x-geo-region",
								},
							},
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
			c.TestCheckObjectAttr("vm", map[string]interface{}{
				"guest_os":     "linux",
				"run_strategy": "RerunOnFailure",
				"hostname":     "vm-host-2",
				"subdomain":    "vms2",
				"boot_disk": map[string]interface{}{
					"bus":        "sata",
					"boot_order": "1",
					"source": map[string]interface{}{
						"http": map[string]interface{}{
							"url":      "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img",
							"checksum": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						},
					},
					"persist": map[string]interface{}{
						"volume_set": fmt.Sprintf("cpln://volumeset/vmboot-%s", wrt.RandomName),
					},
				},
				"cpu": map[string]interface{}{
					"sockets": "4",
					"threads": "2",
				},
				"firmware": map[string]interface{}{
					"bootloader":  "bios",
					"secure_boot": "false",
					"uuid":        "5d8e7a3c-1f2b-4c6d-8e9f-0a1b2c3d4e5f",
					"serial":      "vm-serial-02",
					"smbios": map[string]interface{}{
						"manufacturer": "ControlPlane",
						"product":      "cpln-vm",
						"version":      "3.0",
						"sku":          "sku-02",
						"family":       "cpln",
					},
				},
				"network": []map[string]interface{}{
					{
						"name": "default",
					},
				},
				"cloud_init": map[string]interface{}{
					"user_data_base64": "I2Nsb3VkLWNvbmZpZwpwYWNrYWdlczoKICAtIGh0b3AK",
					"ssh_public_key_secrets": []string{
						GetSelfLink(OrgName, "secret", fmt.Sprintf("vm-ssh-%s", wrt.RandomName)),
						GetSelfLink(OrgName, "secret", fmt.Sprintf("vm-ssh2-%s", wrt.RandomName)),
					},
				},
				"access_credential": []map[string]interface{}{
					{
						"ssh_public_key_secret": GetSelfLink(OrgName, "secret", fmt.Sprintf("vm-ssh-%s", wrt.RandomName)),
						"users":                 []string{"root", "ubuntu", "admin"},
						"delivery_method":       "configDrive",
					},
					{
						"ssh_public_key_secret": GetSelfLink(OrgName, "secret", fmt.Sprintf("vm-ssh2-%s", wrt.RandomName)),
						"users":                 []string{"deploy"},
						"delivery_method":       "qemuGuestAgent",
					},
				},
				"clock": map[string]interface{}{
					"timezone": "Europe/London",
				},
			}),
		),
	}
}

// BuildVmDefaultsOmittedTestStep constructs a vm workload test step that omits firmware/network/clock and asserts the materialized defaults.
func (wrt *WorkloadResourceTest) BuildVmDefaultsOmittedTestStep(name string) (WorkloadResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "workload",
			ResourceName:    "new",
			ResourceAddress: "cpln_workload.new",
			Name:            name,
			GvcName:         wrt.GvcCase.Name,
			Description:     name,
		},
	}

	// Initialize and return the test step
	return c, resource.TestStep{
		Config: wrt.VmDefaultsOmittedHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "vm"),
			c.TestCheckObjectAttr("vm", map[string]interface{}{
				"firmware": map[string]interface{}{
					"bootloader":  "efi",
					"secure_boot": "false",
				},
				"network": []map[string]interface{}{
					{
						"name": "default",
					},
				},
				"clock": map[string]interface{}{
					"timezone": "UTC",
				},
			}),
		),
	}
}

// BuildVmDefaultsSetTestStep constructs a vm workload test step that sets firmware/network/clock explicitly to override the defaults.
func (wrt *WorkloadResourceTest) BuildVmDefaultsSetTestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := WorkloadResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the test step
	return resource.TestStep{
		Config: wrt.VmDefaultsSetHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.GetDefaultChecks(c.Description, "2"),
			resource.TestCheckResourceAttr(c.ResourceAddress, "gvc", wrt.GvcCase.Name),
			resource.TestCheckResourceAttr(c.ResourceAddress, "type", "vm"),
			c.TestCheckObjectAttr("vm", map[string]interface{}{
				"firmware": map[string]interface{}{
					"bootloader":  "bios",
					"secure_boot": "false",
				},
				"network": []map[string]interface{}{
					{
						"name": "default",
					},
				},
				"clock": map[string]interface{}{
					"timezone": "America/New_York",
				},
			}),
		),
	}
}

// Configs //

// K8sVolumeUriHcl returns a serverless workload configuration that mounts a k8s://secret volume to exercise the uri scheme validator.
func (wrt *WorkloadResourceTest) K8sVolumeUriHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name = "%s"
  gvc  = %s
  type = "serverless"

  container {
    name  = "container-01"
    image = "gcr.io/knative-samples/helloworld-go"

    ports {
      protocol = "http"
      number   = "8080"
    }

    volume {
      uri  = "k8s://secret/example-secret"
      path = "/k8s-secret"
    }
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name,
		wrt.GvcCase.GetResourceNameAttr(),
	)
}

// ServerlessRequiredOnlyHcl returns a minimal serverless workload configuration with only required fields set.
func (wrt *WorkloadResourceTest) ServerlessRequiredOnlyHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name = "%s"
  gvc  = %s
  type = "serverless"

  container {
    name  = "container-01"
    image = "gcr.io/knative-samples/helloworld-go"

    ports {
      protocol = "http"
      number   = "8080"
    }
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name,
		wrt.GvcCase.GetResourceNameAttr(),
	)
}

// ServerlessUpdate1Hcl returns an extended serverless workload configuration with additional fields like tags, probes, and autoscaling.
func (wrt *WorkloadResourceTest) ServerlessUpdate1Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "serverless"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    port              = 8080
    memory            = "128Mi"
    min_cpu           = "25m"
    min_memory        = "32Mi"
    cpu               = "50m"

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
      env-name-04 = ""
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    readiness_probe {

      tcp_socket {}
    }

    liveness_probe {

      http_get {
        path   = "/path"
        scheme = "HTTPS"
      }
    }

    lifecycle {}

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri  = "azureblob://storageAccount/container"
      path = "/testpath02"
    }
  }

  firewall_spec {}

  options {
    autoscaling {}
  }

  local_options {
    location = "aws-eu-central-1"
    autoscaling {}
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {}
  security_options {}
  load_balancer {}
  request_retry_policy {}
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// ServerlessUpdate2Hcl returns a serverless workload configuration with detailed autoscaling, GPU, lifecycle, and firewall settings.
func (wrt *WorkloadResourceTest) ServerlessUpdate2Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "serverless"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    port              = 8080
    memory            = "7Gi"
    cpu               = "2"

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
      env-name-04 = null
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path         = "/metrics"
      port         = 8181
			drop_metrics = ["go_.*", "process_.*", ".*_bucket|.*_sum|.*_count"]
    }

    readiness_probe {

      grpc {
        port = 3000
      }

      initial_delay_seconds = 1
      period_seconds        = 11
      timeout_seconds       = 2
      success_threshold     = 2
      failure_threshold     = 4
    }

    liveness_probe {

      http_get {
        path = "/path"
        port = 8282
        scheme = "HTTPS"
        http_headers = {
          header-name-01 = "header-value-01"
          header-name-02 = "header-value-02"
        }
      }

      initial_delay_seconds = 2
      period_seconds        = 10
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    gpu_nvidia {
      model    = "t4"
      quantity = 1
    }

    lifecycle {
      post_start {}
      pre_stop {}
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {}
    internal {}
  }

  options {
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 100
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  local_options {
    location = "aws-eu-central-1"
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 100
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {
    min_ready_seconds                = 10
    scaling_policy                   = "Parallel"
    termination_grace_period_seconds = 10
  }

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = false

    direct {
      enabled = false
    }

    geo_location {}
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// ServerlessUpdate3Hcl returns a serverless workload configuration with multi-location local_options and load balancer details.
func (wrt *WorkloadResourceTest) ServerlessUpdate3Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "serverless"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    port              = 8080
    memory            = "128Mi"
    cpu               = "50m"

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
      env-name-04 = ""
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
			drop_metrics = ["go_.*", ".*_bucket|.*_sum|.*_count"]
    }

    readiness_probe {

      tcp_socket {
        port = 3200
      }

      initial_delay_seconds = 1
      period_seconds        = 11
      timeout_seconds       = 2
      success_threshold     = 2
      failure_threshold     = 4
    }

    liveness_probe {

      exec {
        command = ["command-01", "command-02"]
      }

      initial_delay_seconds = 2
      period_seconds        = 10
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    gpu_custom {
      resource = "amd.com/gpu"
      quantity = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {}
    }

    internal {
      inbound_allow_type     = "none"
      inbound_allow_workload = []
    }
  }

  options {
    timeout_seconds            = 30
    capacity_ai                = false
		capacity_ai_update_minutes = 5
    debug                      = true
    suspend                    = true

    autoscaling {
      metric              = "concurrency"
      target              = 100
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  local_options {
    location        = "aws-eu-central-1"
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 100
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  local_options {
    location        = "aws-us-west-2"
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 90
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {
    min_ready_seconds                = 20
    max_unavailable_replicas         = "10"
    scaling_policy                   = "OrderedReady"
    termination_grace_period_seconds = 20
  }

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = false

    direct {
      enabled = false
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// ServerlessUpdate4Hcl returns a serverless workload configuration with extended firewall HTTP filter and workload allowlist.
func (wrt *WorkloadResourceTest) ServerlessUpdate4Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "serverless"
  identity_link        = "//gvc/${cpln_identity.new.gvc}/identity/${cpln_identity.new.name}"
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    port              = 8080
    memory            = "128Mi"
    cpu               = "50m"

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
      env-name-04 = "env-value-04"
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    readiness_probe {

      tcp_socket {
        port = 3200
      }

      initial_delay_seconds = 1
      period_seconds        = 11
      timeout_seconds       = 2
      success_threshold     = 2
      failure_threshold     = 4
    }

    liveness_probe {

      exec {
        command = ["command-01", "command-02"]
      }

      initial_delay_seconds = 2
      period_seconds        = 10
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    gpu_custom {
      resource      = "amd.com/gpu"
      runtime_class = "amd"
      quantity      = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "Allow-Header"
          allowed_values = ["reg", "req2"]
        }
      }
    }

    internal {
      inbound_allow_type     = "workload-list"
      inbound_allow_workload = ["//gvc/new/workload/non-existing-workload-01", "/org/terraform-test-org/gvc/new/workload/non-existing-workload-02"]
    }
  }

  options {
    timeout_seconds            = 30
    capacity_ai                = false
		capacity_ai_update_minutes = 10
    debug                      = true
    suspend                    = true

    autoscaling {
      metric              = "concurrency"
      target              = 100
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  local_options {
    location        = "aws-eu-central-1"
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 100
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  local_options {
    location        = "aws-us-west-2"
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 90
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {
    min_ready_seconds                = 30
    max_unavailable_replicas         = "20"
    max_surge_replicas               = "30"
    scaling_policy                   = "Parallel"
    termination_grace_period_seconds = 40
  }

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = false

    direct {
      enabled = false
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// StandardHcl returns a standard workload configuration with all supported features including multi autoscaling and GPU.
func (wrt *WorkloadResourceTest) StandardHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "standard"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    memory            = "128Mi"
    cpu               = "50m"

    ports {
      protocol = "http2"
      number   = "8080" 
    }

    ports {
      protocol = "grpc"
      number   = "3000" 
    }

    ports {
      protocol = "tcp"
      number   = "3001" 
    }

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    readiness_probe {

      tcp_socket {
        port = 3200
      }

      initial_delay_seconds = 1
      period_seconds        = 11
      timeout_seconds       = 2
      success_threshold     = 2
      failure_threshold     = 4
    }

    liveness_probe {

      exec {
        command = ["command-01", "command-02"]
      }

      initial_delay_seconds = 2
      period_seconds        = 10
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    gpu_custom {
      resource      = "amd.com/gpu"
      runtime_class = "amd"
      quantity      = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "Allow-Header"
          blocked_values = ["blocked", "blocked2"]
        }
      }
    }

    internal {
      inbound_allow_type     = "same-org"
      inbound_allow_workload = []
    }
  }

  options {
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric_percentile   = "p50"
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400

      multi {
        metric = "cpu"
        target = 50
      }

      multi {
        metric = "memory"
        target = 50
      }
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {}

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = false

    direct {
      enabled = false
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// StandardUpdate1Hcl returns a standard workload configuration with all supported features including keda autoscaling and GPU.
func (wrt *WorkloadResourceTest) StandardUpdate1Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "standard"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    memory            = "128Mi"
    cpu               = "50m"

    ports {
      protocol = "http2"
      number   = "8080" 
    }

    ports {
      protocol = "grpc"
      number   = "3000" 
    }

    ports {
      protocol = "tcp"
      number   = "3001" 
    }

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    readiness_probe {

      tcp_socket {
        port = 3200
      }

      initial_delay_seconds = 1
      period_seconds        = 11
      timeout_seconds       = 2
      success_threshold     = 2
      failure_threshold     = 4
    }

    liveness_probe {

      exec {
        command = ["command-01", "command-02"]
      }

      initial_delay_seconds = 2
      period_seconds        = 10
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    gpu_custom {
      resource      = "amd.com/gpu"
      runtime_class = "amd"
      quantity      = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "Allow-Header"
          blocked_values = ["blocked", "blocked2"]
        }
      }
    }

    internal {
      inbound_allow_type     = "same-org"
      inbound_allow_workload = []
    }
  }

  options {
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "keda"
      metric_percentile   = "p50"
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400

      keda {
        polling_interval        = 30
        cooldown_period         = 60
        initial_cooldown_period = 10

        trigger {
          type               = "cpu"
          name               = "cpu-trigger-01"
          use_cached_metrics = true
          metric_type        = "Utilization"

          metadata = {
            type  = "Utilization"
            value = "50"
          }
        }

        trigger {
          type               = "rabbitmq"
          name               = "rabbitmq-trigger"
          use_cached_metrics = false
          metric_type        = "AverageValue"

          metadata = {
            host        = "amqp://user:pass@rabbitmq:5672/"
            queueName   = "jobs"
            queueLength = "30"
          }

          authentication_ref {
            name = "rabbitmq-auth"
          }
        }

        advanced {
          scaling_modifiers {
            target            = "5"
            activation_target = "1"
            metric_type       = "Value"
            formula           = "m * 2"
          }
        }

        fallback {
          failure_threshold = 3
          replicas          = 1
          behavior          = "static"
        }
      }
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {
    min_ready_seconds                = 2
    max_unavailable_replicas         = "10"
    max_surge_replicas               = "20"
    scaling_policy                   = "Parallel"
    termination_grace_period_seconds = 10
  }

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = false

    direct {
      enabled = false
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// CronHcl returns a cron workload configuration with job scheduling and standard runtime configuration.
func (wrt *WorkloadResourceTest) CronHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "cron"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    memory            = "128Mi"
    cpu               = "50m"

    ports {
      protocol = "http2"
      number   = "8080" 
    }

    ports {
      protocol = "grpc"
      number   = "3000" 
    }

    ports {
      protocol = "tcp"
      number   = "3001" 
    }

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    gpu_custom {
      resource      = "amd.com/gpu"
      runtime_class = "amd"
      quantity      = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "Allow-Header"
          blocked_values = ["blocked", "blocked2"]
        }
      }
    }

    internal {
      inbound_allow_type     = "same-org"
      inbound_allow_workload = []
    }
  }

  options {
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 95
      max_scale           = 1
      min_scale           = 1
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  job {
    schedule = "* * * * *"
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {}

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = false

    direct {
      enabled = false
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// CronUpdate1Hcl returns an updated cron workload configuration with advanced job options like concurrency and deadline.
func (wrt *WorkloadResourceTest) CronUpdate1Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "cron"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    memory            = "128Mi"
    cpu               = "50m"

    ports {
      protocol = "http2"
      number   = "8080" 
    }

    ports {
      protocol = "grpc"
      number   = "3000" 
    }

    ports {
      protocol = "tcp"
      number   = "3001" 
    }

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    gpu_custom {
      resource      = "amd.com/gpu"
      runtime_class = "amd"
      quantity      = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "Allow-Header"
          blocked_values = ["blocked", "blocked2"]
        }
      }
    }

    internal {
      inbound_allow_type     = "same-org"
      inbound_allow_workload = []
    }
  }

  options {
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "concurrency"
      target              = 95
      max_scale           = 1
      min_scale           = 1
      max_concurrency     = 500
      scale_to_zero_delay = 400
    }
  }

  job {
    schedule                = "* * * * *"
    concurrency_policy      = "Forbid"
    history_limit           = 5
    restart_policy          = "Never"
    active_deadline_seconds = 1200
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {
    min_ready_seconds                = 2
    max_unavailable_replicas         = "10"
    max_surge_replicas               = "20"
    scaling_policy                   = "Parallel"
    termination_grace_period_seconds = 10
  }

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = false

    direct {
      enabled = false
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// StatefulHcl returns a stateful workload configuration with persistent networking, GPU, lifecycle, and autoscaling features.
func (wrt *WorkloadResourceTest) StatefulHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "stateful"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    memory            = "128Mi"
    cpu               = "50m"

    ports {
      protocol = "http2"
      number   = "8080" 
    }

    ports {
      protocol = "grpc"
      number   = "3000" 
    }

    ports {
      protocol = "tcp"
      number   = "3001" 
    }

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    readiness_probe {

      tcp_socket {
        port = 3200
      }

      initial_delay_seconds = 1
      period_seconds        = 11
      timeout_seconds       = 2
      success_threshold     = 2
      failure_threshold     = 4
    }

    liveness_probe {

      exec {
        command = ["command-01", "command-02"]
      }

      initial_delay_seconds = 2
      period_seconds        = 10
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    gpu_custom {
      resource      = "amd.com/gpu"
      runtime_class = "amd"
      quantity      = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "Allow-Header"
          blocked_values = ["blocked", "blocked2"]
        }
      }
    }

    internal {
      inbound_allow_type     = "same-org"
      inbound_allow_workload = []
    }
  }

  options {
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric_percentile   = "p50"
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400

      multi {
        metric = "cpu"
        target = 50
      }

      multi {
        metric = "memory"
        target = 50
      }
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {}

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = true

    direct {
      enabled = true
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// StatefulUpdate1Hcl returns a stateful workload configuration with persistent networking, GPU, lifecycle, and keda autoscaling features.
func (wrt *WorkloadResourceTest) StatefulUpdate1Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

# Identity Resource
resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_workload" "%s" {
  depends_on = [%s]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "stateful"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode(%s)

  container {
    name              = "container-01"
    image             = "gcr.io/knative-samples/helloworld-go"
    working_directory = "/usr"
    memory            = "128Mi"
    cpu               = "50m"

    ports {
      protocol = "http2"
      number   = "8080" 
    }

    ports {
      protocol = "grpc"
      number   = "3000" 
    }

    ports {
      protocol = "tcp"
      number   = "3001" 
    }

    env = {
      env-name-01 = "env-value-01"
      env-name-02 = "env-value-02"
      env-name-03 = "env-value-03"
    }

    inherit_env = true
    command     = "override-command"
    args        = ["arg-01", "arg-02", "arg-03"]

    metrics {
      path = "/metrics"
      port = 8181
    }

    readiness_probe {

      tcp_socket {
        port = 3200
      }

      initial_delay_seconds = 1
      period_seconds        = 11
      timeout_seconds       = 2
      success_threshold     = 2
      failure_threshold     = 4
    }

    liveness_probe {

      exec {
        command = ["command-01", "command-02"]
      }

      initial_delay_seconds = 2
      period_seconds        = 10
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    gpu_custom {
      resource      = "amd.com/gpu"
      runtime_class = "amd"
      quantity      = 1
    }

    lifecycle {
      post_start {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
      pre_stop {
        exec {
          command = ["command-01", "command-02", "command-03"]
        }
      }
    }

    volume {
      uri             = "s3://bucket"
      recovery_policy = "retain"
      path            = "/testpath01"
    }

    volume {
      uri             = "azureblob://storageAccount/container"
      path            = "/testpath02"
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["127.0.0.1"]
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["::1"]

      outbound_allow_port {
        protocol = "http"
        number   = 80
      }

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "Allow-Header"
          blocked_values = ["blocked", "blocked2"]
        }
      }
    }

    internal {
      inbound_allow_type     = "same-org"
      inbound_allow_workload = []
    }
  }

  options {
    timeout_seconds = 30
    capacity_ai     = false
    debug           = true
    suspend         = true

    autoscaling {
      metric              = "keda"
      metric_percentile   = "p50"
      max_scale           = 3
      min_scale           = 2
      max_concurrency     = 500
      scale_to_zero_delay = 400

      keda {
        trigger {
          type        = "cpu"
          name        = "cpu-trigger-01"
          metric_type = "Utilization"

          metadata = {
            type  = "Utilization"
            value = "50"
          }
        }

        advanced {
          scaling_modifiers {}
        }
      }
    }
  }

  sidecar {
    envoy = jsonencode(%s)
  }

  rollout_options {
    min_ready_seconds                = 2
    max_surge_replicas               = "20"
    scaling_policy                   = "Parallel"
    termination_grace_period_seconds = 10
  }

  security_options {
    file_system_group_id = 1
    run_as_user          = 1000
  }

  load_balancer {
    replica_direct = true

    direct {
      enabled = true
      ipset   = "my-ipset-01"

      port {
        external_port  = 22
        protocol       = "TCP"
        scheme         = "http"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "198.51.100.0/24"
        city    = "Los Angeles"
        country = "USA"
        region  = "North America"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, wrt.GvcCase.ResourceAddress, c.Name, c.DescriptionUpdate,
		wrt.GvcCase.GetResourceNameAttr(), c.Extras, c.Envoy,
	)
}

// VmHcl returns a vm workload configuration with an OCI boot disk, persistent boot volume, cpu topology, firmware, and a data disk.
func (wrt *WorkloadResourceTest) VmHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_volume_set" "vm_boot" {
  depends_on = [%s]

  name              = "vmboot-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_volume_set" "vm_data" {
  depends_on = [%s]

  name              = "vmdata-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_workload" "%s" {
  depends_on = [cpln_volume_set.vm_boot, cpln_volume_set.vm_data]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc  = %s
  type = "vm"

  container {
    name   = "vm-container"
    cpu    = "2000m"
    memory = "2Gi"

    volume {
      uri        = "cpln://volumeset/vmdata-${var.random_name}"
      name       = "data-disk"
      bus        = "virtio"
      boot_order = 2
    }
  }

  vm = {
    boot_disk = {
      source = {
        oci = {
          image = "quay.io/containerdisks/ubuntu:22.04"
        }
      }

      persist = {
        volume_set = "cpln://volumeset/vmboot-${var.random_name}"
      }

      bus        = "virtio"
      boot_order = 1
    }

    cpu = {
      sockets = 2
      threads = 1
    }

    firmware = {
      bootloader  = "efi"
      secure_boot = false
      uuid        = "5d8e7a3c-1f2b-4c6d-8e9f-0a1b2c3d4e5f"

      smbios = {
        manufacturer = "ControlPlane"
        product      = "cpln-vm"
      }
    }

    guest_os = "linux"

    network = [
      {
        name = "default"
      }
    ]

    run_strategy = "Always"

    clock = {
      timezone = "UTC"
    }

    hostname = "vm-host"
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), wrt.GvcCase.ResourceAddress, wrt.GvcCase.GetResourceNameAttr(), wrt.GvcCase.ResourceAddress,
		wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, c.Name, c.Description, wrt.GvcCase.GetResourceNameAttr(),
	)
}

// VmUpdate1Hcl returns an updated vm workload configuration with cloud-init, access credentials, full firmware, and modified topology.
func (wrt *WorkloadResourceTest) VmUpdate1Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_volume_set" "vm_boot" {
  depends_on = [%s]

  name              = "vmboot-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_volume_set" "vm_data" {
  depends_on = [%s]

  name              = "vmdata-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_secret" "vm_ssh" {
  name = "vm-ssh-${var.random_name}"

  opaque {
    payload  = "c3NoLXJzYSBBQUFBQjNOemFDMTljMkVBQUFBREFRQUJBQUFB"
    encoding = "base64"
  }
}

resource "cpln_policy" "vm_secrets" {
  name        = "vm-secrets-${var.random_name}"
  target_kind = "secret"
  target      = "all"

  binding {
    permissions     = ["reveal"]
    principal_links = [cpln_identity.new.self_link]
  }
}

resource "cpln_workload" "%s" {
  depends_on = [cpln_identity.new, cpln_policy.vm_secrets, cpln_volume_set.vm_boot, cpln_volume_set.vm_data, cpln_secret.vm_ssh]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc           = %s
  type          = "vm"
  identity_link = cpln_identity.new.self_link

  container {
    name   = "vm-container"
    cpu    = "4000m"
    memory = "4Gi"

    volume {
      uri        = "cpln://volumeset/vmdata-${var.random_name}"
      name       = "data-disk"
      bus        = "scsi"
      boot_order = 2
    }
  }

  vm = {
    boot_disk = {
      source = {
        oci = {
          image = "quay.io/containerdisks/ubuntu:24.04"
        }
      }

      persist = {
        volume_set = "cpln://volumeset/vmboot-${var.random_name}"
      }

      bus        = "virtio"
      boot_order = 1
    }

    cpu = {
      sockets = 4
      threads = 2
    }

    firmware = {
      bootloader  = "efi"
      secure_boot = false
      uuid        = "5d8e7a3c-1f2b-4c6d-8e9f-0a1b2c3d4e5f"
      serial      = "vm-serial-01"

      smbios = {
        manufacturer = "ControlPlane"
        product      = "cpln-vm"
        version      = "2.0"
        sku          = "sku-01"
        family       = "cpln"
      }
    }

    guest_os = "linux"

    network = [
      {
        name = "default"
      }
    ]

    cloud_init = {
      user_data              = "#cloud-config\nruncmd:\n  - echo hello\n"
      ssh_public_key_secrets = [cpln_secret.vm_ssh.self_link]
    }

    access_credential = [
      {
        ssh_public_key_secret = cpln_secret.vm_ssh.self_link
        users                 = ["root", "ubuntu"]
        delivery_method       = "qemuGuestAgent"
      }
    ]

    run_strategy = "Manual"

    clock = {
      timezone = "America/New_York"
    }

    hostname  = "vm-host-updated"
    subdomain = "vms"
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), wrt.GvcCase.ResourceAddress, wrt.GvcCase.GetResourceNameAttr(), wrt.GvcCase.ResourceAddress,
		wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, c.Name, c.DescriptionUpdate, wrt.GvcCase.GetResourceNameAttr(),
	)
}

// VmUpdate2Hcl returns a fully-featured vm workload. On top of the vm block (http boot source, base64
// cloud-init, two access-credentials/ssh-keys, and the remaining enum values), it layers on every
// workload-level attribute that round-trips cleanly with type=vm: identity_link, support_dynamic_tags, an
// enriched container (ports, env, inherit_env, metrics, tcp_socket/http_get probes), firewall_spec,
// load_balancer (direct + geo_location), and request_retry_policy. This surfaces any conflict between vm
// and other workload features.
//
// Intentionally omitted because the API REJECTS them for type=vm: security_options, job, container
// image/port/command/args/lifecycle/working_directory, exec & grpc probes, load_balancer.replica_direct.
//
// Intentionally omitted because the API STRIPS Computed/defaulted fields for type=vm, which collides with
// the provider's client-side defaults and produces an "inconsistent result after apply": the whole
// `options` block (capacity_ai defaults true but is stripped; autoscaling only allows metric=disabled) and
// `rollout_options` (scaling_policy defaults OrderedReady but is stripped, as is max_surge_replicas). The
// existing cron scenario avoids `options` for the same reason. Surfacing those cleanly needs the provider
// to skip those defaults for type=vm; until then they can't be asserted without perpetual drift.
//
// Also omitted: sidecar/extras (envoy filters and BYOK k8s modifications are not meaningful on a cloud GVC).
func (wrt *WorkloadResourceTest) VmUpdate2Hcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

resource "cpln_identity" "new" {
  name = "identity-${var.random_name}"
  gvc  = %s
}

resource "cpln_volume_set" "vm_boot" {
  depends_on = [%s]

  name              = "vmboot-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_volume_set" "vm_data" {
  depends_on = [%s]

  name              = "vmdata-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_secret" "vm_ssh" {
  name = "vm-ssh-${var.random_name}"

  opaque {
    payload  = "c3NoLXJzYSBBQUFBQjNOemFDMTljMkVBQUFBREFRQUJBQUFB"
    encoding = "base64"
  }
}

resource "cpln_secret" "vm_ssh2" {
  name = "vm-ssh2-${var.random_name}"

  opaque {
    payload  = "c3NoLXJzYSBBQUFBQjNOemFDMTljMkVCQkJCQ0RRQUJCQkJC"
    encoding = "base64"
  }
}

resource "cpln_policy" "vm_secrets" {
  name        = "vm-secrets-${var.random_name}"
  target_kind = "secret"
  target      = "all"

  binding {
    permissions     = ["reveal"]
    principal_links = [cpln_identity.new.self_link]
  }
}

resource "cpln_workload" "%s" {
  depends_on = [cpln_identity.new, cpln_policy.vm_secrets, cpln_volume_set.vm_boot, cpln_volume_set.vm_data, cpln_secret.vm_ssh, cpln_secret.vm_ssh2]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc                  = %s
  type                 = "vm"
  identity_link        = cpln_identity.new.self_link
  support_dynamic_tags = true
  extras               = jsonencode({ tolerations = [{ key = "cpln.io/nodeType", operator = "Equal", value = "vm", effect = "NoSchedule" }] })

  container {
    name        = "vm-container"
    cpu         = "4000m"
    memory      = "4Gi"
    inherit_env = true

    ports {
      protocol = "tcp"
      number   = 8080
    }

    ports {
      protocol = "http"
      number   = 80
    }

    env = {
      ENV_KEY = "env-value"
      APP_ENV = "production"
    }

    metrics {
      port         = 8181
      path         = "/metrics"
      drop_metrics = ["envoy_.*", "go_gc_.*"]
    }

    readiness_probe {
      tcp_socket {
        port = 8080
      }

      initial_delay_seconds = 5
      period_seconds        = 10
      timeout_seconds       = 2
      success_threshold     = 1
      failure_threshold     = 3
    }

    liveness_probe {
      http_get {
        path   = "/healthz"
        port   = 80
        scheme = "HTTP"

        http_headers = {
          "X-Custom-Header" = "custom-value"
        }
      }

      initial_delay_seconds = 10
      period_seconds        = 15
      timeout_seconds       = 3
      success_threshold     = 1
      failure_threshold     = 5
    }

    volume {
      uri        = "cpln://volumeset/vmdata-${var.random_name}"
      name       = "data-disk"
      bus        = "virtio"
      boot_order = 2
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      inbound_blocked_cidr    = ["192.0.2.1"]
      outbound_allow_hostname = ["*.controlplane.com"]
      outbound_allow_cidr     = []
      outbound_blocked_cidr   = ["198.51.100.1"]

      outbound_allow_port {
        protocol = "https"
        number   = 443
      }

      http {
        inbound_header_filter {
          key            = "X-Allowed"
          allowed_values = ["^v1$", "^v2$"]
        }
      }
    }

    internal {
      inbound_allow_type     = "same-gvc"
      inbound_allow_workload = []
    }
  }

  load_balancer {
    direct {
      enabled = true
      ipset   = "vm-ipset"

      port {
        external_port  = 8080
        protocol       = "TCP"
        scheme         = "tcp"
        container_port = 80
      }
    }

    geo_location {
      enabled = true

      headers {
        asn     = "x-geo-asn"
        city    = "x-geo-city"
        country = "x-geo-country"
        region  = "x-geo-region"
      }
    }
  }

  request_retry_policy {
    attempts = 3
    retry_on = ["connect-failure", "refused-stream", "unavailable"]
  }

  vm = {
    boot_disk = {
      source = {
        http = {
          url      = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
          checksum = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
        }
      }

      persist = {
        volume_set = "cpln://volumeset/vmboot-${var.random_name}"
      }

      bus        = "sata"
      boot_order = 1
    }

    cpu = {
      sockets = 4
      threads = 2
    }

    firmware = {
      bootloader  = "bios"
      secure_boot = false
      uuid        = "5d8e7a3c-1f2b-4c6d-8e9f-0a1b2c3d4e5f"
      serial      = "vm-serial-02"

      smbios = {
        manufacturer = "ControlPlane"
        product      = "cpln-vm"
        version      = "3.0"
        sku          = "sku-02"
        family       = "cpln"
      }
    }

    guest_os = "linux"

    network = [
      {
        name = "default"
      }
    ]

    cloud_init = {
      user_data_base64       = "I2Nsb3VkLWNvbmZpZwpwYWNrYWdlczoKICAtIGh0b3AK"
      ssh_public_key_secrets = [cpln_secret.vm_ssh.self_link, cpln_secret.vm_ssh2.self_link]
    }

    access_credential = [
      {
        ssh_public_key_secret = cpln_secret.vm_ssh.self_link
        users                 = ["root", "ubuntu", "admin"]
        delivery_method       = "configDrive"
      },
      {
        ssh_public_key_secret = cpln_secret.vm_ssh2.self_link
        users                 = ["deploy"]
        delivery_method       = "qemuGuestAgent"
      }
    ]

    run_strategy = "RerunOnFailure"

    clock = {
      timezone = "Europe/London"
    }

    hostname  = "vm-host-2"
    subdomain = "vms2"
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.GetResourceNameAttr(), wrt.GvcCase.ResourceAddress, wrt.GvcCase.GetResourceNameAttr(), wrt.GvcCase.ResourceAddress,
		wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, c.Name, c.DescriptionUpdate, wrt.GvcCase.GetResourceNameAttr(),
	)
}

// VmDefaultsOmittedHcl returns a vm workload that omits firmware, network, and clock so their defaults are materialized.
func (wrt *WorkloadResourceTest) VmDefaultsOmittedHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

resource "cpln_volume_set" "vm_boot" {
  depends_on = [%s]

  name              = "vmboot-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_workload" "%s" {
  depends_on = [cpln_volume_set.vm_boot]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc  = %s
  type = "vm"

  container {
    name   = "vm-container"
    cpu    = "1000m"
    memory = "1Gi"
  }

  vm = {
    boot_disk = {
      source = {
        oci = {
          image = "quay.io/containerdisks/ubuntu:22.04"
        }
      }

      persist = {
        volume_set = "cpln://volumeset/vmboot-${var.random_name}"
      }
    }
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.ResourceAddress, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, c.Name, c.Description, wrt.GvcCase.GetResourceNameAttr(),
	)
}

// VmDefaultsSetHcl returns a vm workload that sets firmware, network, and clock explicitly to override their defaults.
func (wrt *WorkloadResourceTest) VmDefaultsSetHcl(c WorkloadResourceTestCase) string {
	return fmt.Sprintf(`
variable "random_name" {
  type    = string
  default = "%s"
}

# GVC Resource
%s

resource "cpln_volume_set" "vm_boot" {
  depends_on = [%s]

  name              = "vmboot-${var.random_name}"
  gvc               = %s
  initial_capacity  = 10
  performance_class = "general-purpose-ssd"
  file_system_type  = "ext4"
}

resource "cpln_workload" "%s" {
  depends_on = [cpln_volume_set.vm_boot]

  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  gvc  = %s
  type = "vm"

  container {
    name   = "vm-container"
    cpu    = "1000m"
    memory = "1Gi"
  }

  vm = {
    boot_disk = {
      source = {
        oci = {
          image = "quay.io/containerdisks/ubuntu:22.04"
        }
      }

      persist = {
        volume_set = "cpln://volumeset/vmboot-${var.random_name}"
      }
    }

    firmware = {
      bootloader  = "bios"
      secure_boot = false
    }

    network = [
      {
        name = "default"
      }
    ]

    clock = {
      timezone = "America/New_York"
    }
  }
}
`, wrt.RandomName, wrt.GvcConfig, wrt.GvcCase.ResourceAddress, wrt.GvcCase.GetResourceNameAttr(), c.ResourceName, c.Name, c.Description, wrt.GvcCase.GetResourceNameAttr(),
	)
}

/*** Resource Test Case ***/

// WorkloadResourceTestCase defines a specific resource test case.
type WorkloadResourceTestCase struct {
	ProviderTestCase
	Envoy  string
	Extras string
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (wrtc *WorkloadResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of workload: %s. Total resources: %d", wrtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[wrtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", wrtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != wrtc.Name {
			return fmt.Errorf("resource ID %s does not match expected workload name %s", rs.Primary.ID, wrtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteWorkload, _, err := TestProvider.client.GetWorkload(wrtc.Name, wrtc.GvcName)
		if err != nil {
			return fmt.Errorf("error retrieving workload from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteWorkload.Name != wrtc.Name {
			return fmt.Errorf("mismatch in workload name: expected %s, got %s", wrtc.Name, *remoteWorkload.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("workload %s verified successfully in both state and external system.", wrtc.Name))
		return nil
	}
}
