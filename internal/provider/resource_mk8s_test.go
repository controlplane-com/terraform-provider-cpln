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

// TestAccControlPlaneMk8s_basic performs an acceptance test for the resource.
func TestAccControlPlaneMk8s_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewMk8sResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "MK8S") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.Steps,
	})
}

/*** Resource Test ***/

// Mk8sResourceTest defines the necessary functionality to test the resource.
type Mk8sResourceTest struct {
	Steps      []resource.TestStep
	RandomName string
}

// NewMk8sResourceTest creates a Mk8sResourceTest with initialized test cases.
func NewMk8sResourceTest() Mk8sResourceTest {
	// Create a resource test instance
	resourceTest := Mk8sResourceTest{
		RandomName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
	}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	// steps = append(steps, resourceTest.NewMk8sGenericProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sHetznerProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sAwsProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sLinodeProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sOblivusProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sLambdalabsProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sPaperspaceProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sEphemeralProviderScenario()...)
	// steps = append(steps, resourceTest.NewMk8sTritonProviderScenario()...)
	steps = append(steps, resourceTest.NewMk8sGcpProviderScenario()...)

	// Set the cases for the resource test
	resourceTest.Steps = steps

	// Return the resource test
	return resourceTest
}

// CheckDestroy verifies that all resources have been destroyed.
func (mrt *Mk8sResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_mk8s resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_mk8s" {
			continue
		}

		// Retrieve the name for the current resource
		mk8sName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of mk8s with name: %s", mk8sName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		mk8s, code, err := TestProvider.client.GetMk8s(mk8sName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if mk8s %s exists: %w", mk8sName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if mk8s != nil {
			return fmt.Errorf("CheckDestroy failed: mk8s %s still exists in the system", *mk8s.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_mk8s resources have been successfully destroyed")
	return nil
}

// Test Scenarios //

// NewMk8sGenericProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sGenericProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "generic"
	name := fmt.Sprintf("tf-mk8s-generic-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildGenericProviderTestStep(resourceName, name)
	caseUpdate1 := mrt.BuildGenericProviderUpdate1TestStep(initialConfig.ProviderTestCase)
	caseUpdate2 := mrt.BuildGenericProviderUpdate2TestStep(initialConfig.ProviderTestCase)
	caseUpdate3 := mrt.BuildGenericProviderUpdate3TestStep(initialConfig.ProviderTestCase)

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
		caseUpdate2,
		caseUpdate1,
		// Revert the resource to its initial state
		initialStep,
	}
}

// NewMk8sHetznerProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sHetznerProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "hetzner"
	name := fmt.Sprintf("tf-mk8s-hetzner-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildHetznerProviderTestStep(resourceName, name)

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

// NewMk8sAwsProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sAwsProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "aws"
	name := fmt.Sprintf("tf-mk8s-aws-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildAwsProviderTestStep(resourceName, name)

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

// NewMk8sLinodeProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sLinodeProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "linode"
	name := fmt.Sprintf("tf-mk8s-linode-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildLinodeProviderTestStep(resourceName, name)

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

// NewMk8sOblivusProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sOblivusProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "oblivus"
	name := fmt.Sprintf("tf-mk8s-oblivus-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildOblivusProviderTestStep(resourceName, name)

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

// NewMk8sLambdalabsProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sLambdalabsProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "lambdalabs"
	name := fmt.Sprintf("tf-mk8s-lambdalabs-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildLambdalabsProviderTestStep(resourceName, name)

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

// NewMk8sPaperspaceProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sPaperspaceProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "paperspace"
	name := fmt.Sprintf("tf-mk8s-paperspace-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildPaperspaceProviderTestStep(resourceName, name)

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

// NewMk8sEphemeralProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sEphemeralProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "ephemeral"
	name := fmt.Sprintf("tf-mk8s-ephemeral-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildEphemeralProviderTestStep(resourceName, name)

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

// NewMk8sTritonProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sTritonProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "triton"
	name := fmt.Sprintf("tf-mk8s-triton-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildTritonProviderTestStep(resourceName, name)
	caseUpdate1 := mrt.BuildTritonProviderUpdate1TestStep(initialConfig.ProviderTestCase)

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

// NewMk8sGcpProviderScenario creates a test case with initial and updated configurations.
func (mrt *Mk8sResourceTest) NewMk8sGcpProviderScenario() []resource.TestStep {
	// Define necessary variables
	resourceName := "gcp"
	name := fmt.Sprintf("tf-mk8s-gcp-%s", mrt.RandomName)

	// Build test steps
	initialConfig, initialStep := mrt.BuildGcpProviderTestStep(resourceName, name)

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

// BuildGenericProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildGenericProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s generic new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.GenericProviderRequiredOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.Description, "0"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "0.0.0.0/0",
				},
			}),
			c.TestCheckNestedBlocks("generic_provider", []map[string]interface{}{
				{
					"location": "aws-eu-central-1",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
				},
			}),
		),
	}
}

// BuildGenericProviderUpdate1TestStep returns a test step for the update.
func (mrt *Mk8sResourceTest) BuildGenericProviderUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: mrt.GenericProviderUpdate1Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("generic_provider", []map[string]interface{}{
				{
					"location": "aws-eu-central-1",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"node_pool": []map[string]interface{}{
						{
							"name": "my-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{{}}),
		),
	}
}

// BuildGenericProviderUpdate1TestStep returns a test step for the update.
func (mrt *Mk8sResourceTest) BuildGenericProviderUpdate2TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: mrt.GenericProviderUpdate2Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("generic_provider", []map[string]interface{}{
				{
					"location": "aws-eu-central-1",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"node_pool": []map[string]interface{}{
						{
							"name": "my-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard":               "false",
					"aws_workload_identity":   "false",
					"local_path_storage":      "false",
					"sysbox":                  "false",
					"azure_workload_identity": []map[string]interface{}{{}},
					"metrics":                 []map[string]interface{}{{}},
					"logs":                    []map[string]interface{}{{}},
					"registry_mirror":         []map[string]interface{}{{}},
					"nvidia":                  []map[string]interface{}{{}},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"byok": c.ExpectedByokFull(true),
				},
			}),
		),
	}
}

// BuildGenericProviderUpdate2TestStep returns a test step for the update.
func (mrt *Mk8sResourceTest) BuildGenericProviderUpdate3TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: mrt.GenericProviderUpdate3Hcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("generic_provider", []map[string]interface{}{
				{
					"location": "aws-eu-central-1",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"node_pool": []map[string]interface{}{
						{
							"name": "my-node-pool-01",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
						{
							"name": "my-node-pool-02",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					"byok":   c.ExpectedByokFull(false),
				},
			}),
		),
	}
}

// BuildHetznerProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildHetznerProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s hetzner new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.HetznerProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("hetzner_provider", []map[string]interface{}{
				{
					"region": "fsn1",
					"hetzner_labels": map[string]interface{}{
						"hello": "world",
					},
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"pre_install_script": "#! echo hello world",
					"token_secret_link":  "/org/terraform-test-org/secret/hetzner",
					"network_id":         "2808575",
					"node_pool": []map[string]interface{}{
						{
							"name": "my-hetzner-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
							"server_type":    "cpx11",
							"override_image": "debian-11",
							"min_size":       "0",
							"max_size":       "0",
						},
					},
					"dedicated_server_node_pool": []map[string]interface{}{
						{
							"name": "my-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
					"image":   "centos-7",
					"ssh_key": "10925607",
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
					"floating_ip_selector": map[string]interface{}{
						"floating_ip_1": "123.45.67.89",
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildAwsProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildAwsProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s aws new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.AwsProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("aws_provider", []map[string]interface{}{
				{
					"region": "eu-central-1",
					"aws_tags": map[string]interface{}{
						"hello": "world",
					},
					"skip_create_roles": "false",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"pre_install_script": "#! echo hello world",
					"image": []map[string]interface{}{
						{
							"recommended": "amazon/al2023",
						},
					},
					"deploy_role_arn":         "arn:aws:iam::483676437512:role/cpln-mk8s-terraform-test-org",
					"vpc_id":                  "vpc-03105bd4dc058d3a8",
					"key_pair":                "debug-terraform",
					"disk_encryption_key_arn": "arn:aws:kms:eu-central-1:989132402664:key/2e9f25ea-efb4-49bf-ae39-007be298726d",
					"security_group_ids":      []string{"sg-031480aa7a1e6e38b"},
					"extra_node_policies":     []string{"arn:aws:iam::aws:policy/IAMFullAccess"},
					"deploy_role_chain": []map[string]interface{}{
						{
							"role_arn":            "arn:aws:iam::483676437512:role/mk8s-chain-1",
							"external_id":         "chain-1",
							"session_name_prefix": "foo-",
						},
					},
					"node_pool": []map[string]interface{}{
						{
							"name": "my-hetzner-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
							"instance_types": []string{"t4g.nano"},
							"override_image": []map[string]interface{}{
								{
									"exact": "ami-0c5ee33c81cf67a7f",
								},
							},
							"boot_disk_size":          "20",
							"min_size":                "0",
							"max_size":                "0",
							"on_demand_base_capacity": "0",
							"on_demand_percentage_above_base_capacity": "0",
							"spot_allocation_strategy":                 "lowest-price",
							"subnet_ids":                               []string{"subnet-0e564a042e2a45009"},
							"extra_security_group_ids":                 []string{},
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildLinodeProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildLinodeProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s linode new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.LinodeProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("linode_provider", []map[string]interface{}{
				{
					"region":            "fr-par",
					"token_secret_link": "/org/terraform-test-org/secret/linode",
					"firewall_id":       "168425",
					"image":             "linode/ubuntu24.04",
					"authorized_users":  []string{"juliancpln"},
					// "authorized_keys":    []string{"cplnkey"},
					"vpc_id":             "93666",
					"pre_install_script": "#! echo hello world",
					"node_pool": []map[string]interface{}{
						{
							"name": "my-linode-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
							"server_type":    "g6-nanode-1",
							"override_image": "linode/debian11",
							"subnet_id":      "90623",
							"min_size":       0,
							"max_size":       0,
						},
					},
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildOblivusProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildOblivusProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s oblivus new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.OblivusProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("oblivus_provider", []map[string]interface{}{
				{
					"datacenter":         "OSL1",
					"token_secret_link":  "/org/terraform-test-org/secret/oblivus",
					"pre_install_script": "#! echo hello world",
					"node_pool": []map[string]interface{}{
						{
							"name": "my-oblivus-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
							"min_size": 0,
							"max_size": 0,
							"flavor":   "INTEL_XEON_V3_x4",
						},
					},
					"unmanaged_node_pool": []map[string]interface{}{
						{
							"name": "my-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildLambdalabsProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildLambdalabsProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s lambdalabs new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.LambdalabsProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("lambdalabs_provider", []map[string]interface{}{
				{
					"region":             "europe-central-1",
					"token_secret_link":  "/org/terraform-test-org/secret/lambdalabs",
					"ssh_key":            "julian-test",
					"pre_install_script": "#! echo hello world",
					"node_pool": []map[string]interface{}{
						{
							"name":          "my-lambdalabs-node-pool",
							"instance_type": "cpu_4x_general",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
					"unmanaged_node_pool": []map[string]interface{}{
						{
							"name": "my-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildPaperspaceProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildPaperspaceProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s paperspace new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.PaperspaceProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("paperspace_provider", []map[string]interface{}{
				{
					"region":             "CA1",
					"token_secret_link":  "/org/terraform-test-org/secret/paperspace",
					"pre_install_script": "#! echo hello world",
					"shared_drives":      []string{"california"},
					"network_id":         "nla0jotp",
					"node_pool": []map[string]interface{}{
						{
							"name":           "my-paperspace-node-pool",
							"min_size":       "0",
							"max_size":       "0",
							"public_ip_type": "dynamic",
							"boot_disk_size": "50",
							"machine_type":   "GPU+",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
					"unmanaged_node_pool": []map[string]interface{}{
						{
							"name": "my-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildEphemeralProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildEphemeralProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s ephemeral new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.EphemeralProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("ephemeral_provider", []map[string]interface{}{
				{
					"location": "aws-eu-central-1",
					"node_pool": []map[string]interface{}{
						{
							"name":   "my-ephemeral-node-pool",
							"count":  "1",
							"arch":   "arm64",
							"flavor": "debian",
							"cpu":    "50m",
							"memory": "128Mi",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
				},
			}),
		),
	}
}

// BuildTritonProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildTritonProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s triton new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.TritonProviderWithGatewayLoadBalancerHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("triton_provider", []map[string]interface{}{
				{
					"pre_install_script": "#! echo hello world",
					"location":           "aws-eu-central-1",
					"private_network_id": "6704dae9-00f4-48b5-8bbf-1be538f20587",
					"firewall_enabled":   "false",
					"image_id":           "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"connection": []map[string]interface{}{
						{
							"url":                     "https://us-central-1.api.mnx.io",
							"account":                 "eric_controlplane.com",
							"user":                    "julian_controlplane.com",
							"private_key_secret_link": "/org/terraform-test-org/secret/triton",
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"gateway": []map[string]interface{}{{}},
						},
					},
					"node_pool": []map[string]interface{}{
						{
							"name":                "my-triton-node-pool",
							"package_id":          "da311341-b42b-45a8-9386-78ede625d0a4",
							"override_image_id":   "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e",
							"public_network_id":   "5ff1fe03-075b-4e4c-b85b-73de0c452f77",
							"min_size":            0,
							"max_size":            0,
							"private_network_ids": []string{"6704dae9-00f4-48b5-8bbf-1be538f20587"},
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
							"triton_tags": map[string]interface{}{
								"drink": "water",
							},
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildTritonProviderUpdate1TestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildTritonProviderUpdate1TestStep(initialCase ProviderTestCase) resource.TestStep {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: initialCase,
	}

	// Initialize and return the inital test step
	return resource.TestStep{
		Config: mrt.TritonProviderWithManualLoadBalancerHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("triton_provider", []map[string]interface{}{
				{
					"pre_install_script": "#! echo hello world",
					"location":           "aws-eu-central-1",
					"private_network_id": "6704dae9-00f4-48b5-8bbf-1be538f20587",
					"firewall_enabled":   "false",
					"image_id":           "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"connection": []map[string]interface{}{
						{
							"url":                     "https://us-central-1.api.mnx.io",
							"account":                 "eric_controlplane.com",
							"user":                    "julian_controlplane.com",
							"private_key_secret_link": "/org/terraform-test-org/secret/triton",
						},
					},
					"load_balancer": []map[string]interface{}{
						{
							"manual": []map[string]interface{}{
								{
									"package_id":          "df26ba1d-1261-6fc1-b35c-f1b390bc06ff",
									"image_id":            "8605a524-0655-43b9-adf1-7d572fe797eb",
									"public_network_id":   "5ff1fe03-075b-4e4c-b85b-73de0c452f77",
									"private_network_ids": []string{"6704dae9-00f4-48b5-8bbf-1be538f20587"},
									"count":               "1",
									"cns_internal_domain": "example.com",
									"cns_public_domain":   "example.com",
									"metadata": map[string]interface{}{
										"key1": "value1",
										"key2": "value2",
									},
									"tags": map[string]interface{}{
										"tag1": "value1",
										"tag2": "value2",
									},
									"logging": []map[string]interface{}{
										{
											"node_port":       "32000",
											"external_syslog": "syslog.example.com:514",
										},
									},
								},
							},
						},
					},
					"node_pool": []map[string]interface{}{
						{
							"name":                "my-triton-node-pool",
							"package_id":          "da311341-b42b-45a8-9386-78ede625d0a4",
							"override_image_id":   "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e",
							"public_network_id":   "5ff1fe03-075b-4e4c-b85b-73de0c452f77",
							"min_size":            0,
							"max_size":            0,
							"private_network_ids": []string{"6704dae9-00f4-48b5-8bbf-1be538f20587"},
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
							"triton_tags": map[string]interface{}{
								"drink": "water",
							},
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// BuildGcpProviderTestStep returns a default initial test step and its associated test case for the resource.
func (mrt *Mk8sResourceTest) BuildGcpProviderTestStep(resourceName string, name string) (Mk8sResourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := Mk8sResourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:              "mk8s",
			ResourceName:      resourceName,
			ResourceAddress:   fmt.Sprintf("cpln_mk8s.%s", resourceName),
			Name:              name,
			Description:       name,
			DescriptionUpdate: "mk8s gcp new description",
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: mrt.GcpProviderHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.GetDefaultChecks(c.DescriptionUpdate, "3"),
			c.TestCheckResourceAttr("version", "1.28.4"),
			c.TestCheckNestedBlocks("firewall", []map[string]interface{}{
				{
					"source_cidr": "192.168.1.255",
					"description": "hello world",
				},
			}),
			c.TestCheckNestedBlocks("gcp_provider", []map[string]interface{}{
				{
					"project_id": "coke-267310",
					"region":     "us-west1",
					"gcp_labels": map[string]interface{}{
						"hello": "world",
					},
					"network":     "mk8s",
					"sa_key_link": "/org/terraform-test-org/secret/gcp",
					"networking": []map[string]interface{}{
						{
							"service_network": "10.43.0.0/16",
							"pod_network":     "10.42.0.0/16",
						},
					},
					"pre_install_script": "#! echo hello world",
					"image": []map[string]interface{}{
						{
							"recommended": "ubuntu/jammy-22.04",
						},
					},
					"node_pool": []map[string]interface{}{
						{
							"name": "my-gcp-node-pool",
							"labels": map[string]interface{}{
								"hello": "world",
							},
							"taint": []map[string]interface{}{
								{
									"key":    "hello",
									"value":  "world",
									"effect": "NoSchedule",
								},
							},
							"machine_type": "n1-standard-2",
							"zone":         "us-west1-a",
							"override_image": []map[string]interface{}{
								{
									"recommended": "ubuntu/noble-24.04",
								},
							},
							"boot_disk_size": "30",
							"min_size":       "0",
							"max_size":       "0",
							"subnet":         "mk8s",
						},
					},
					"autoscaler": []map[string]interface{}{
						{
							"expander":              []string{"most-pods"},
							"unneeded_time":         "10m",
							"unready_time":          "20m",
							"utilization_threshold": "0.7",
						},
					},
				},
			}),
			c.TestCheckNestedBlocks("add_ons", []map[string]interface{}{
				{
					"dashboard": "true",
					"azure_workload_identity": []map[string]interface{}{
						{
							"tenant_id": "7f43458a-a34e-4bfa-9e56-e2289e49c4ec",
						},
					},
					"aws_workload_identity": "true",
					"local_path_storage":    "true",
					"metrics": []map[string]interface{}{
						{
							"kube_state":    "true",
							"core_dns":      "true",
							"kubelet":       "true",
							"api_server":    "true",
							"node_exporter": "true",
							"cadvisor":      "true",
							"scrape_annotated": []map[string]interface{}{
								{
									"interval_seconds":   "30",
									"include_namespaces": "^elastic",
									"exclude_namespaces": "^elastic",
									"retain_labels":      "^\\w+$",
								},
							},
						},
					},
					"logs": []map[string]interface{}{
						{
							"audit_enabled":      "true",
							"include_namespaces": "^elastic",
							"exclude_namespaces": "^elastic",
						},
					},
					"registry_mirror": []map[string]interface{}{
						{
							"mirror": []map[string]interface{}{
								{
									"registry": "registry.mycompany.com",
									"mirrors":  []string{"https://mirror1.mycompany.com"},
								},
								{
									"registry": "docker.io",
									"mirrors":  []string{"https://us-mirror.gcr.io"},
								},
							},
						},
					},
					"nvidia": []map[string]interface{}{
						{
							"taint_gpu_nodes": "true",
						},
					},
					"azure_acr": []map[string]interface{}{
						{
							"client_id": "4e25b134-160b-4a9d-b392-13b381ced5ef",
						},
					},
					"sysbox": "true",
					// TODO: Add byok test here
				},
			}),
		),
	}
}

// Configs //

// GenericProviderRequiredOnlyHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) GenericProviderRequiredOnlyHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name    = "%s"
  version = "1.28.4"

  firewall {
    source_cidr = "0.0.0.0/0"
  }

  generic_provider {
    location = "aws-eu-central-1"
		networking {}
  }
}
`, c.ResourceName, c.Name)
}

// GenericProviderUpdate1Hcl returns a test step for the update.
func (mrt *Mk8sResourceTest) GenericProviderUpdate1Hcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  generic_provider {
    location = "aws-eu-central-1"

    networking {
      service_network = "10.43.0.0/16"
      pod_network 	  = "10.42.0.0/16"
    }

    node_pool {
      name = "my-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }
  }

  add_ons {}
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// GenericProviderUpdate1Hcl returns a test step for the update.
func (mrt *Mk8sResourceTest) GenericProviderUpdate2Hcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  generic_provider {
    location = "aws-eu-central-1"

    networking {
      service_network = "10.43.0.0/16"
      pod_network 	  = "10.42.0.0/16"
    }

    node_pool {
      name = "my-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }
  }

  add_ons {
    dashboard             = false
    aws_workload_identity = false
    local_path_storage    = false
    sysbox                = false

    azure_workload_identity {}
    metrics {}
    logs {}
    registry_mirror {}
    nvidia {}

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    byok = {
      location = "/org/terraform-test-org/location/test-byok"
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// GenericProviderUpdate1Hcl returns a test step for the update.
func (mrt *Mk8sResourceTest) GenericProviderUpdate3Hcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  generic_provider {
    location = "aws-eu-central-1"

    networking {
      service_network = "10.43.0.0/16"
      pod_network 	  = "10.42.0.0/16"
    }

    node_pool {
      name = "my-node-pool-01"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

		node_pool {
      name = "my-node-pool-02"
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
      docker             = true
      kubelet            = true
      kernel             = true
      events             = true
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// HetznerProviderHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) HetznerProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  hetzner_provider {

    region = "fsn1"

    hetzner_labels = {
      hello = "world"
    }

    networking {
      service_network = "10.43.0.0/16"
      pod_network     = "10.42.0.0/16"
    }

    pre_install_script = "#! echo hello world"
    token_secret_link  = "/org/terraform-test-org/secret/hetzner"
    network_id         = "2808575"

    node_pool {
      name = "my-hetzner-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }

      server_type    = "cpx11"
      override_image = "debian-11"
      min_size       = 0
      max_size       = 0
    }

    dedicated_server_node_pool {
      name = "my-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

    image   = "centos-7"
    ssh_key = "10925607"

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }

    floating_ip_selector = {
      floating_ip_1 = "123.45.67.89"
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// AwsProviderHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) AwsProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  aws_provider {

    region = "eu-central-1"

    aws_tags = {
      hello = "world"
    }

    skip_create_roles  = false

    networking {}

    pre_install_script = "#! echo hello world"

    image {
      recommended = "amazon/al2023"
    }

    deploy_role_arn         = "arn:aws:iam::483676437512:role/cpln-mk8s-terraform-test-org"
    vpc_id                  = "vpc-03105bd4dc058d3a8"
    key_pair                = "debug-terraform"
    disk_encryption_key_arn = "arn:aws:kms:eu-central-1:989132402664:key/2e9f25ea-efb4-49bf-ae39-007be298726d"

    security_group_ids  = ["sg-031480aa7a1e6e38b"]
    extra_node_policies = ["arn:aws:iam::aws:policy/IAMFullAccess"]

    deploy_role_chain {
      role_arn            = "arn:aws:iam::483676437512:role/mk8s-chain-1"
      external_id         = "chain-1"
      session_name_prefix = "foo-"
    }

    node_pool {
      name = "my-hetzner-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }

      instance_types = ["t4g.nano"]

      override_image {
        exact = "ami-0c5ee33c81cf67a7f"
      }

      boot_disk_size                           = 20
      min_size                                 = 0
      max_size                                 = 0
      on_demand_base_capacity                  = 0
      on_demand_percentage_above_base_capacity = 0
      spot_allocation_strategy                 = "lowest-price"

      subnet_ids               = ["subnet-0e564a042e2a45009"]
      extra_security_group_ids = []
    }

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// LinodeProviderHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) LinodeProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  linode_provider {
    region             = "fr-par"
    token_secret_link  = "/org/terraform-test-org/secret/linode"
    firewall_id        = "168425"
    image              = "linode/ubuntu24.04"
    authorized_users   = ["juliancpln"]
		// authorized_keys    = ["cplnkey"]
    vpc_id             = "93666"
    pre_install_script = "#! echo hello world"

    node_pool {
      name = "my-linode-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }

      server_type    = "g6-nanode-1"
      override_image = "linode/debian11"
      subnet_id      = "90623"
      min_size 	   = 0
      max_size 	   = 0
    }

    networking {
      service_network = "10.43.0.0/16"
      pod_network 	= "10.42.0.0/16"
    }

    autoscaler {
      expander 	  		      = ["most-pods"]
      unneeded_time         = "10m"
      unready_time  		    = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// OblivusProviderHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) OblivusProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  oblivus_provider {
    datacenter         = "OSL1"
    token_secret_link  = "/org/terraform-test-org/secret/oblivus"
    pre_install_script = "#! echo hello world"

    node_pool {
      name           = "my-oblivus-node-pool"
      min_size       = 0
      max_size       = 0
      flavor         = "INTEL_XEON_V3_x4"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

    unmanaged_node_pool {
      name = "my-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// LambdalabsProviderHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) LambdalabsProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  lambdalabs_provider {
    region             = "europe-central-1"
    token_secret_link  = "/org/terraform-test-org/secret/lambdalabs"
    ssh_key            = "julian-test"
    pre_install_script = "#! echo hello world"

    node_pool {
      name           = "my-lambdalabs-node-pool"
      instance_type  = "cpu_4x_general"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

    unmanaged_node_pool {
      name = "my-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// PaperspaceProviderHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) PaperspaceProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  paperspace_provider {
    region             = "CA1"
    token_secret_link  = "/org/terraform-test-org/secret/paperspace"
    shared_drives      = ["california"]
    pre_install_script = "#! echo hello world"
    network_id         = "nla0jotp"

    node_pool {
      name           = "my-paperspace-node-pool"
      min_size       = 0
      max_size       = 0
      public_ip_type = "dynamic"
      boot_disk_size = 50
      machine_type   = "GPU+"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

    unmanaged_node_pool {
      name = "my-node-pool"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// EphemeralProviderHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) EphemeralProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  ephemeral_provider {
    location = "aws-eu-central-1"

    node_pool {
      name   = "my-ephemeral-node-pool"
      count  = 1
      arch   = "arm64"
      flavor = "debian"
      cpu    = "50m"
      memory = "128Mi"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// TritonProviderWithGatewayLoadBalancerHcl returns a test step for the update.
func (mrt *Mk8sResourceTest) TritonProviderWithGatewayLoadBalancerHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  triton_provider {
    pre_install_script = "#! echo hello world"
    location           = "aws-eu-central-1"
    private_network_id = "6704dae9-00f4-48b5-8bbf-1be538f20587"
    firewall_enabled   = false
    image_id           = "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2"

    networking {}

    connection {
      url                     = "https://us-central-1.api.mnx.io"
      account                 = "eric_controlplane.com"
      user                    = "julian_controlplane.com"
      private_key_secret_link = "/org/terraform-test-org/secret/triton"
    }

    load_balancer {
      gateway {}
    }

    node_pool {
      name              = "my-triton-node-pool"
      package_id        = "da311341-b42b-45a8-9386-78ede625d0a4"
      override_image_id = "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
      public_network_id = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
      min_size          = 0
      max_size          = 0

      private_network_ids = ["6704dae9-00f4-48b5-8bbf-1be538f20587"]

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }

      triton_tags = {
        drink = "water"
      }
    }

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// TritonProviderUpdate1Hcl returns a test step for the update.
func (mrt *Mk8sResourceTest) TritonProviderWithManualLoadBalancerHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  triton_provider {
    pre_install_script = "#! echo hello world"
    location           = "aws-eu-central-1"
    private_network_id = "6704dae9-00f4-48b5-8bbf-1be538f20587"
    firewall_enabled   = false
    image_id           = "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2"

    networking {}

    connection {
      url                     = "https://us-central-1.api.mnx.io"
      account                 = "eric_controlplane.com"
      user                    = "julian_controlplane.com"
      private_key_secret_link = "/org/terraform-test-org/secret/triton"
    }

    load_balancer {
      manual {
        package_id          = "df26ba1d-1261-6fc1-b35c-f1b390bc06ff"
        image_id            = "8605a524-0655-43b9-adf1-7d572fe797eb"
        public_network_id   = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
        private_network_ids = ["6704dae9-00f4-48b5-8bbf-1be538f20587"]
        count               = 1
        cns_internal_domain = "example.com"
        cns_public_domain   = "example.com"

        metadata = {
          key1 = "value1"
          key2 = "value2"
        }

        tags = {
          tag1 = "value1"
          tag2 = "value2"
        }

				logging {
				  node_port       = 32000
				  external_syslog = "syslog.example.com:514"
				}
      }
    }

    node_pool {
      name              = "my-triton-node-pool"
      package_id        = "da311341-b42b-45a8-9386-78ede625d0a4"
      override_image_id = "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
      public_network_id = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
      min_size          = 0
      max_size          = 0

      private_network_ids = ["6704dae9-00f4-48b5-8bbf-1be538f20587"]

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }

      triton_tags = {
        drink = "water"
      }
    }

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

// GcpProviderHcl returns a test step.
func (mrt *Mk8sResourceTest) GcpProviderHcl(c Mk8sResourceTestCase) string {
	return fmt.Sprintf(`
resource "cpln_mk8s" "%s" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    "cpln/ignore"       = "true"
  }

  version = "1.28.4"

  firewall {
    source_cidr = "192.168.1.255"
    description = "hello world"
  }

  gcp_provider {
    project_id         = "coke-267310"
    region             = "us-west1"
    network            = "mk8s"
    sa_key_link        = "/org/terraform-test-org/secret/gcp"
    pre_install_script = "#! echo hello world"

    gcp_labels = {
      hello = "world"
    }

    networking {}

    image {
      recommended = "ubuntu/jammy-22.04"
    }

    node_pool {
      name           = "my-gcp-node-pool"
      machine_type   = "n1-standard-2"
      zone           = "us-west1-a"
      boot_disk_size = 30
      min_size       = 0
      max_size       = 0
      subnet         = "mk8s"

      labels = {
        hello = "world"
      }

      taint {
        key    = "hello"
        value  = "world"
        effect = "NoSchedule"
      }

      override_image {
        recommended = "ubuntu/noble-24.04"
      }
    }

    autoscaler {
      expander              = ["most-pods"]
      unneeded_time         = "10m"
      unready_time          = "20m"
      utilization_threshold = 0.7
    }
  }

  add_ons {
    dashboard = true

    azure_workload_identity {
      tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
    }

    aws_workload_identity = true
    local_path_storage    = true

    metrics {
      kube_state    = true
      core_dns      = true
      kubelet       = true
      api_server    = true
      node_exporter = true
      cadvisor      = true

      scrape_annotated {
        interval_seconds   = 30
        include_namespaces = "^elastic"
        exclude_namespaces = "^elastic"
        retain_labels      = "^\\w+$"
      }
    }

    logs {
      audit_enabled      = true
      include_namespaces = "^elastic"
      exclude_namespaces = "^elastic"
    }

    registry_mirror {
      mirror {
        registry = "registry.mycompany.com"
        mirrors = ["https://mirror1.mycompany.com"]
      }

      mirror {
        registry = "docker.io"
        mirrors = ["https://us-mirror.gcr.io"]
      }
    }

    nvidia {
      taint_gpu_nodes = true
    }

    azure_acr {
      client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
    }

    sysbox = true

    byok = {
      ignore_updates = false
      location       = "/org/terraform-test-org/location/test-byok"

      config = {
        actuator = {
          min_cpu    = "50m"
          max_cpu    = "8001m"
          min_memory = "200Mi"
          max_memory = "8000Mi"
          log_level  = "info"
          env = {
            CACHE_PERIOD_DATA_SERVICE = "600"
            LABEL_NODES               = "false"
          }
        }

        middlebox = {
          enabled              = false
          bandwidth_alert_mbps = 650
        }

        common = {
          deployment_replicas = 1

          // pbd {
          //   max_unavailable = 1
          // }
        }

        longhorn = {
          replicas = 2
        }

        ingress = {
          cpu            = "50m"
          memory         = "200Mi"
          target_percent = 6000
        }

        istio = {
          istiod = {
            replicas   = 2
            min_cpu    = "50m"
            max_cpu    = "8001m"
            min_memory = "100Mi"
            max_memory = "8000Mi"
            // pbd        = 10
          }

          ingress_gateway = {
            replicas   = 2
            max_cpu    = "1"
            max_memory = "1Gi"
          }

          sidecar = {
            min_cpu    = "0m"
            min_memory = "200Mi"
          }
        }

        log_splitter = {
          min_cpu         = "1m"
          max_cpu         = "1000m"
          min_memory      = "10Mi"
          max_memory      = "2000Mi"
          mem_buffer_size = "128M"
          per_pod_rate    = 10000
        }

        monitoring = {
          min_memory = "100Mi"
          max_memory = "20Gi"

          kube_state_metrics = {
            min_memory = "25Mi"
          }

          prometheus = {
            main = {
              storage = "10Gi"
            }
          }
        }

        redis = {
          min_cpu    = "10m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = "8Gi"
        }

        redis_ha = {
          min_cpu    = "50m"
          max_cpu    = "2001m"
          min_memory = "100Mi"
          max_memory = "1000Mi"
          storage    = 0
        }

        redis_sentinel = {
          min_cpu    = "10m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
          storage    = 0
        }

        tempo_agent = {
          min_cpu    = "0m"
          min_memory = "10Mi"
        }

        internal_dns = {
          min_cpu    = "0m"
          max_cpu    = "500m"
          min_memory = "10Mi"
          max_memory = "400Mi"
        }
      }
    }
  }
}
`, c.ResourceName, c.Name, c.DescriptionUpdate)
}

/*** Resource Test Case ***/

// Mk8sResourceTestCase defines a specific resource test case.
type Mk8sResourceTestCase struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (mrtc *Mk8sResourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of mk8s: %s. Total resources: %d", mrtc.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[mrtc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", mrtc.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != mrtc.Name {
			return fmt.Errorf("resource ID %s does not match expected mk8s name %s", rs.Primary.ID, mrtc.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteMk8s, _, err := TestProvider.client.GetMk8s(mrtc.Name)
		if err != nil {
			return fmt.Errorf("error retrieving mk8s from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteMk8s.Name != mrtc.Name {
			return fmt.Errorf("mismatch in mk8s name: expected %s, got %s", mrtc.Name, *remoteMk8s.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("mk8s %s verified successfully in both state and external system.", mrtc.Name))
		return nil
	}
}

// ExpectedByokFull builds a focused but deep expectation for BYOK "config" used in Update3 tests.
func (mrtc *Mk8sResourceTestCase) ExpectedByokFull(isDefault bool) map[string]interface{} {
	output := map[string]interface{}{
		"location": "/org/terraform-test-org/location/test-byok",
		"config": map[string]interface{}{
			"actuator": map[string]interface{}{
				"min_cpu":    "50m",
				"max_cpu":    "8001m",
				"min_memory": "200Mi",
				"max_memory": "8000Mi",
				"log_level":  "info",
				"env": map[string]interface{}{
					"CACHE_PERIOD_DATA_SERVICE": "600",
					"LABEL_NODES":               "false",
				},
			},
			"middlebox": map[string]interface{}{
				"enabled":              false,
				"bandwidth_alert_mbps": 650,
			},
			"common": map[string]interface{}{
				"deployment_replicas": 1,
				// "pbd" is optional in your HCL and defaulted in provider; we don't assert it here.
			},
			"longhorn": map[string]interface{}{
				"replicas": 2,
			},
			"ingress": map[string]interface{}{
				"cpu":            "50m",
				"memory":         "200Mi",
				"target_percent": 6000,
			},
			"istio": map[string]interface{}{
				"istiod": map[string]interface{}{
					"replicas":   2,
					"min_cpu":    "50m",
					"max_cpu":    "8001m",
					"min_memory": "100Mi",
					"max_memory": "8000Mi",
					// "pbd" may be defaulted; we skip it to keep test resilient.
				},
				"ingress_gateway": map[string]interface{}{
					"replicas":   2,
					"max_cpu":    "1",
					"max_memory": "1Gi",
				},
				"sidecar": map[string]interface{}{
					"min_cpu":    "0m",
					"min_memory": "200Mi",
				},
			},
			"log_splitter": map[string]interface{}{
				"min_cpu":         "1m",
				"max_cpu":         "1000m",
				"min_memory":      "10Mi",
				"max_memory":      "2000Mi",
				"mem_buffer_size": "128M",
				"per_pod_rate":    10000,
			},
			"monitoring": map[string]interface{}{
				"min_memory": "100Mi",
				"max_memory": "20Gi",
				"kube_state_metrics": map[string]interface{}{
					"min_memory": "25Mi",
				},
				"prometheus": map[string]interface{}{
					"main": map[string]interface{}{
						"storage": "10Gi",
					},
				},
			},
			"redis": map[string]interface{}{
				"min_cpu":    "10m",
				"max_cpu":    "2001m",
				"min_memory": "100Mi",
				"max_memory": "1000Mi",
				"storage":    "8Gi",
			},
			"redis_ha": map[string]interface{}{
				"min_cpu":    "50m",
				"max_cpu":    "2001m",
				"min_memory": "100Mi",
				"max_memory": "1000Mi",
				"storage":    0,
			},
			"redis_sentinel": map[string]interface{}{
				"min_cpu":    "10m",
				"max_cpu":    "500m",
				"min_memory": "10Mi",
				"max_memory": "400Mi",
				"storage":    0,
			},
			"tempo_agent": map[string]interface{}{
				"min_cpu":    "0m",
				"min_memory": "10Mi",
			},
		},
	}

	if isDefault {
		return output
	}

	output["ignore_updates"] = false
	output["config"].(map[string]interface{})["internal_dns"] = map[string]interface{}{
		"min_cpu":    "0m",
		"max_cpu":    "500m",
		"min_memory": "10Mi",
		"max_memory": "400Mi",
	}

	return output
}
