package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneMk8s_basic(t *testing.T) {

	var mk8s client.Mk8s

	name := "mk8s-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	description := "Mk8s description created using Terraform"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "MK8S") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneMk8sCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneMk8sGenericProvider(name+"-generic", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.generic", name+"-generic", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "generic"),
					resource.TestCheckResourceAttr("cpln_mk8s.generic", "name", name+"-generic"),
					resource.TestCheckResourceAttr("cpln_mk8s.generic", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sHetznerProvider(name+"-hetzner", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.hetzner", name+"-hetzner", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "hetzner"),
					resource.TestCheckResourceAttr("cpln_mk8s.hetzner", "name", name+"-hetzner"),
					resource.TestCheckResourceAttr("cpln_mk8s.hetzner", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sAwsProvider(name+"-aws", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.aws", name+"-aws", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "aws"),
					resource.TestCheckResourceAttr("cpln_mk8s.aws", "name", name+"-aws"),
					resource.TestCheckResourceAttr("cpln_mk8s.aws", "description", description),
				),
			},
		},
	})
}

func testAccCheckControlPlaneMk8sExists(resourceName string, mk8sName string, mk8s *client.Mk8s) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneMk8sExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != mk8sName {
			return fmt.Errorf("Mk8s name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)
		_mk8s, _, err := client.GetMk8s(mk8sName)

		if err != nil {
			return err
		}

		if *_mk8s.Name != mk8sName {
			return fmt.Errorf("Mk8s name does not match")
		}

		*mk8s = *_mk8s

		return nil
	}
}

func testAccCheckControlPlaneMk8sAttributes(mk8s *client.Mk8s, providerName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tags := *mk8s.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("Mk8s Tags - `terraform_generated` attribute does not match")
		}

		if tags["acceptance_test"] != "true" {
			return fmt.Errorf("Mk8s Tags - `acceptance_test` attribute does not match")
		}

		// Firewall
		expectedFirewall, _, _ := generateTestMk8sFirewall()

		if diff := deep.Equal(mk8s.Spec.Firewall, expectedFirewall); diff != nil {
			return fmt.Errorf("Mk8s Firewall does not match. Diff: %s", diff)
		}

		// Provider
		expectedProvider := generateTestMk8sProvider(providerName)

		if diff := deep.Equal(mk8s.Spec.Provider, expectedProvider); diff != nil {
			return fmt.Errorf("Mk8s Provider %s does not match. Diff: %s", providerName, diff)
		}

		// Add Ons

		expectedAddOns, _, _ := generateTestMk8sAddOns(providerName)

		if diff := deep.Equal(mk8s.Spec.AddOns, expectedAddOns); diff != nil {
			return fmt.Errorf("Mk8s Add Ons for provider %s does not match. Diff: %s", providerName, diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneMk8sCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_mk8s" {
			continue
		}

		mk8sName := rs.Primary.ID

		mk8s, _, _ := c.GetMk8s(mk8sName)
		if mk8s != nil {
			return fmt.Errorf("Mk8s still exists. Name: %s", *mk8s.Name)
		}
	}

	return nil
}

// Acceptance Tests //

func testAccControlPlaneMk8sGenericProvider(name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_mk8s" "generic" {
		
		name        = "%s"
		description = "%s"

		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
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
				pod_network 	= "10.42.0.0/16"
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
					include_namespaces = "^\\d+$"
					exclude_namespaces  = "^[a-z]$"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^\\d+$"
				exclude_namespaces  = "^[a-z]$"
			}

			nvidia {
				taint_gpu_nodes = true
			}

			azure_acr {
				client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
			}

			sysbox = true
		}
	}
	`, name, description)
}

func testAccControlPlaneMk8sHetznerProvider(name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_mk8s" "hetzner" {
		
		name        = "%s"
		description = "%s"

		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
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
				pod_network 	= "10.42.0.0/16"
			}

			pre_install_script = "#! echo hello world"
			token_secret_link  = "/org/terraform-test-org/secret/hetzner"
			network_id 		   = "2808575"

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

				server_type    = "cx11"
				override_image = "debian-11"
				min_size 	   = 0
				max_size 	   = 0
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

			image 	= "centos-7"
			ssh_key = "10925607"

			autoscaler {
				expander 	  		  = ["most-pods"]
				unneeded_time         = "10m"
				unready_time  		  = "20m"
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
					include_namespaces = "^\\d+$"
					exclude_namespaces  = "^[a-z]$"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^\\d+$"
				exclude_namespaces  = "^[a-z]$"
			}

			nvidia {
				taint_gpu_nodes = true
			}

			azure_acr {
				client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
			}

			sysbox = true
		}
	}
	`, name, description)
}

func testAccControlPlaneMk8sAwsProvider(name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_mk8s" "aws" {
		
		name        = "%s"
		description = "%s"

		tags = {
			terraform_generated = "true"
			acceptance_test     = "true"
		}
	
		version = "1.28.4"
	
		firewall {
			source_cidr = "192.168.1.255"
			description = "hello world"
		}
	
		aws_provider {
	
			region            = "us-west-2"
			skip_create_roles = false
	
			networking {
				service_network = "10.43.0.0/16"
				pod_network 	= "10.42.0.0/16"
			}
	
			pre_install_script = "#! echo hello world"
	
			image {
				recommended = "amazon/al2023"
			}
	
			deploy_role_arn         = "arn:aws:iam::989132402664:role/cpln-mk8s-terraform-test-org"
			vpc_id                  = "vpc-087b3e0f680a7e91e"
			key_pair                = "debug-eks"
			disk_encryption_key_arn = "arn:aws:kms:us-west-2:989132402664:key/2e9f25ea-efb4-49bf-ae39-007be298726d"
	
			security_group_ids = ["sg-0f659b1b0711edce1"]
	
			node_pool {
				name = "my-aws-node-pool"
	
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
					exact = "ami-123"
				}
	
				boot_disk_size                           = 20
				min_size                                 = 0
				max_size                                 = 0
				on_demand_base_capacity                  = 0
				on_demand_percentage_above_base_capacity = 0
				spot_allocation_strategy                 = "lowest-price"
	
				subnet_ids               = ["subnet-077fe72ab6259d9a2"]
				extra_security_group_ids = []
			}
	
			autoscaler {
				expander 	  		  = ["most-pods"]
				unneeded_time         = "10m"
				unready_time  		  = "20m"
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
					include_namespaces = "^\\d+$"
					exclude_namespaces  = "^[a-z]$"
					retain_labels      = "^\\w+$"
				}
			}
	
			logs {
				audit_enabled      = true
				include_namespaces = "^\\d+$"
				exclude_namespaces  = "^[a-z]$"
			}
	
			nvidia {
				taint_gpu_nodes = true
			}
	
			aws_efs {
				role_arn = "arn:aws:iam::123456789012:role/my-custom-role"
			}
	
			aws_ecr {
				role_arn = "arn:aws:iam::123456789012:role/my-custom-role"
			}
	
			aws_elb {
				role_arn = "arn:aws:iam::123456789012:role/my-custom-role"
			}
	
			azure_acr {
				client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
			}

			sysbox = true
		}
	}
	`, name, description)
}

/*** Unit Tests ***/

// Build //

func TestControlPlane_BuildMk8sFirewall(t *testing.T) {

	firewall, expectedFirewall, _ := generateTestMk8sFirewall()

	if diff := deep.Equal(firewall, expectedFirewall); diff != nil {
		t.Errorf("Mk8s Firewall was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAddOns(t *testing.T) {

	addOns, expectedAddOns, _ := generateTestMk8sAddOns("aws")

	if diff := deep.Equal(addOns, expectedAddOns); diff != nil {
		t.Errorf("Mk8s AddOns was not built correctly, Diff: %s", diff)
	}
}

// Providers

func TestControlPlane_BuildMk8sGenericProvider(t *testing.T) {

	generic, expectedGeneric, _ := generateTestMk8sGenericProvider()

	if diff := deep.Equal(generic, expectedGeneric); diff != nil {
		t.Errorf("Mk8s Generic Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sHetznerProvider(t *testing.T) {

	hetzner, expectedHetzner, _ := generateTestMk8sHetznerProvider()

	if diff := deep.Equal(hetzner, expectedHetzner); diff != nil {
		t.Errorf("Mk8s Hetzner Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAwsProvider(t *testing.T) {

	aws, expectedAws, _ := generateTestMk8sAwsProvider()

	if diff := deep.Equal(aws, expectedAws); diff != nil {
		t.Errorf("Mk8s AWS Provider was not built correctly, Diff: %s", diff)
	}
}

// Node Pools

func TestControlPlane_BuildMk8sGenericNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sGenericNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Generic Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sHetznerNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sHetznerNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Hetzner Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAwsNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sAwsNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s AWS Node Pools was not built correctly, Diff: %s", diff)
	}
}

// AWS

func TestControlPlane_BuildMk8sAwsAmi_Recommended(t *testing.T) {

	ami, expectedAmi, _ := generateTestMk8sAwsAmi("recommended")

	if diff := deep.Equal(ami, expectedAmi); diff != nil {
		t.Errorf("Mk8s AWS Ami Recommended was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAwsAmi_Exact(t *testing.T) {

	ami, expectedAmi, _ := generateTestMk8sAwsAmi("exact")

	if diff := deep.Equal(ami, expectedAmi); diff != nil {
		t.Errorf("Mk8s AWS Ami Exact was not built correctly, Diff: %s", diff)
	}
}

// Common

func TestControlPlane_BuildMk8sNetworking(t *testing.T) {

	networking, expectedNetworking, _ := generateTestMk8sNetworking()

	if diff := deep.Equal(networking, expectedNetworking); diff != nil {
		t.Errorf("Mk8s Networking was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sTaints(t *testing.T) {

	taints, expectedTaints, _ := generateTestMk8sTaints()

	if diff := deep.Equal(taints, expectedTaints); diff != nil {
		t.Errorf("Mk8s Taints was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAutoscaler(t *testing.T) {

	autoscaler, expectedAutoscaler, _ := generateTestMk8sAutoscaler()

	if diff := deep.Equal(autoscaler, expectedAutoscaler); diff != nil {
		t.Errorf("Mk8s Autoscaler was not built correctly, Diff: %s", diff)
	}
}

// Add Ons

func TestControlPlane_BuildMk8sAzureWorkloadIdentityAddOn(t *testing.T) {

	azureWorkloadIdentityAddOn, expectedAzureWorkloadIdentityAddOn, _ := generateTestMk8sAzureWorkloadIdentityAddOn()

	if diff := deep.Equal(azureWorkloadIdentityAddOn, expectedAzureWorkloadIdentityAddOn); diff != nil {
		t.Errorf("Mk8s Azure Add On was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sMetricsAddOn(t *testing.T) {

	metricsAddOn, expectedMetricsAddOn, _ := generateTestMk8sMetricsAddOn()

	if diff := deep.Equal(metricsAddOn, expectedMetricsAddOn); diff != nil {
		t.Errorf("Mk8s Metrics Add On was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sMetricsScrapeAnnotated(t *testing.T) {

	scrapeAnnotated, expectedScrapeAnnotated, _ := generateTestMk8sMetricsScrapeAnnotated()

	if diff := deep.Equal(scrapeAnnotated, expectedScrapeAnnotated); diff != nil {
		t.Errorf("Mk8s Metrics Scrape Annotated was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sLogsAddOn(t *testing.T) {

	logsAddOn, expectedLogsAddOn, _ := generateTestMk8sLogsAddOn()

	if diff := deep.Equal(logsAddOn, expectedLogsAddOn); diff != nil {
		t.Errorf("Mk8s Logs Add On was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sNvidiaAddOn(t *testing.T) {

	nvidiaAddOn, expectedNvidiaAddOn, _ := generateTestMk8sNvidiaAddOn()

	if diff := deep.Equal(nvidiaAddOn, expectedNvidiaAddOn); diff != nil {
		t.Errorf("Mk8s Nvidia Add On was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAwsAddOn(t *testing.T) {

	awsAddOn, expectedAwsAddOn, _ := generateTestMk8sAwsAddOn()

	if diff := deep.Equal(awsAddOn, expectedAwsAddOn); diff != nil {
		t.Errorf("Mk8s AWS Add On was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAzureAcrAddOn(t *testing.T) {

	azureAcrAddOn, expectedAzureAcrAddOn, _ := generateTestMk8sAzureAcrAddOn()

	if diff := deep.Equal(azureAcrAddOn, expectedAzureAcrAddOn); diff != nil {
		t.Errorf("Mk8s Azure ACR Add On was not built correctly, Diff: %s", diff)
	}
}

// Flatten //

func TestControlPlane_FlattenMk8sFirewall(t *testing.T) {

	_, expectedFirewall, expectedFlatten := generateTestMk8sFirewall()
	flattenedFirewall := flattenMk8sFirewall(expectedFirewall)

	if diff := deep.Equal(expectedFlatten, flattenedFirewall); diff != nil {
		t.Errorf("Mk8s Firewall was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAddOns(t *testing.T) {

	_, expectedAddOns, expectedFlatten := generateTestMk8sAddOns("aws")
	flattenedAddOns := flattenMk8sAddOns(expectedAddOns)

	if diff := deep.Equal(expectedFlatten, flattenedAddOns); diff != nil {
		t.Errorf("Mk8s Add Ons was not flattened correctly. Diff: %s", diff)
	}
}

// Providers

func TestControlPlane_FlattenMk8sGenericProvider(t *testing.T) {

	_, expectedGeneric, expectedFlatten := generateTestMk8sGenericProvider()
	flattenedGeneric := flattenMk8sGenericProvider(expectedGeneric)

	if diff := deep.Equal(expectedFlatten, flattenedGeneric); diff != nil {
		t.Errorf("Mk8s Generic Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sHetznerProvider(t *testing.T) {

	_, expectedHetzner, expectedFlatten := generateTestMk8sHetznerProvider()
	flattenedHetzner := flattenMk8sHetznerProvider(expectedHetzner)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedHetzner); diff != nil {
		t.Errorf("Mk8s Hetzner Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAwsProvider(t *testing.T) {

	_, expectedAws, expectedFlatten := generateTestMk8sAwsProvider()
	flattenedAws := flattenMk8sAwsProvider(expectedAws)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})
	expectedFlattenItem["security_group_ids"] = expectedFlattenItem["security_group_ids"].(*schema.Set).List()

	// Node Pool
	nodePool := expectedFlattenItem["node_pool"].([]interface{})[0].(map[string]interface{})
	nodePool["instance_types"] = nodePool["instance_types"].(*schema.Set).List()
	nodePool["subnet_ids"] = nodePool["subnet_ids"].(*schema.Set).List()
	nodePool["extra_security_group_ids"] = nodePool["extra_security_group_ids"].(*schema.Set).List()

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedAws); diff != nil {
		t.Errorf("Mk8s AWS Provider was not flattened correctly. Diff: %s", diff)
	}
}

// Node Pools

func TestControlPlane_FlattenMk8sGenericNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sGenericNodePools()
	flattenedNodePools := flattenMk8sGenericNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Generic Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sHetznerNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sHetznerNodePools()
	flattenedNodePools := flattenMk8sHetznerNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Hetzner Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAwsNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sAwsNodePools()
	flattenedNodePools := flattenMk8sAwsNodePools(expectedNodePools)

	// Extract the interface slice from *schema.Set
	nodePool := expectedFlatten[0].(map[string]interface{})
	nodePool["instance_types"] = nodePool["instance_types"].(*schema.Set).List()
	nodePool["subnet_ids"] = nodePool["subnet_ids"].(*schema.Set).List()
	nodePool["extra_security_group_ids"] = nodePool["extra_security_group_ids"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s AWS Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

// AWS

func TestControlPlane_FlattenMk8sAwsAmi_Recommended(t *testing.T) {

	_, expectedAmi, expectedFlatten := generateTestMk8sAwsAmi("recommended")
	flattenedAmi := flattenMk8sAwsAmi(expectedAmi)

	if diff := deep.Equal(expectedFlatten, flattenedAmi); diff != nil {
		t.Errorf("Mk8s AWS Ami Recommended was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAwsAmi_Exact(t *testing.T) {

	_, expectedAmi, expectedFlatten := generateTestMk8sAwsAmi("exact")
	flattenedAmi := flattenMk8sAwsAmi(expectedAmi)

	if diff := deep.Equal(expectedFlatten, flattenedAmi); diff != nil {
		t.Errorf("Mk8s AWS Ami Exact was not flattened correctly. Diff: %s", diff)
	}
}

// Common

func TestControlPlane_FlattenMk8sNetworking(t *testing.T) {

	_, expectedNetworking, expectedFlatten := generateTestMk8sNetworking()
	flattenedNetworking := flattenMk8sNetworking(expectedNetworking)

	if diff := deep.Equal(expectedFlatten, flattenedNetworking); diff != nil {
		t.Errorf("Mk8s Networking was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sTaints(t *testing.T) {

	_, expectedTaints, expectedFlatten := generateTestMk8sTaints()
	flattenedTaints := flattenMk8sTaints(expectedTaints)

	if diff := deep.Equal(expectedFlatten, flattenedTaints); diff != nil {
		t.Errorf("Mk8s Taints was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAutoscaler(t *testing.T) {

	_, expectedAutoscaler, expectedFlatten := generateTestMk8sAutoscaler()
	flattenedAutoscaler := flattenMk8sAutoscaler(expectedAutoscaler)

	// Extract the interface slice from *schema.Set
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})
	expectedFlattenItem["expander"] = expectedFlattenItem["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedAutoscaler); diff != nil {
		t.Errorf("Mk8s Autoscaler was not flattened correctly. Diff: %s", diff)
	}
}

// Add Ons

func TestControlPlane_FlattenMk8sAzureWorkloadIdentityAddOn(t *testing.T) {

	_, expectedAddOn, expectedFlatten := generateTestMk8sAzureWorkloadIdentityAddOn()
	flattenedAddOn := flattenMk8sAzureWorkloadIdentityAddOn(expectedAddOn)

	if diff := deep.Equal(expectedFlatten, flattenedAddOn); diff != nil {
		t.Errorf("Mk8s Azure Add On was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sMetricsAddOn(t *testing.T) {

	_, expectedAddOn, expectedFlatten := generateTestMk8sMetricsAddOn()
	flattenedAddOn := flattenMk8sMetricsAddOn(expectedAddOn)

	if diff := deep.Equal(expectedFlatten, flattenedAddOn); diff != nil {
		t.Errorf("Mk8s Metrics Add On was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sMetricsScrapeAnnotated(t *testing.T) {

	_, expectedScrapeAnnotated, expectedFlatten := generateTestMk8sMetricsScrapeAnnotated()
	flattenedScrapeAnnotated := flattenMk8sMetricsScrapeAnnotated(expectedScrapeAnnotated)

	if diff := deep.Equal(expectedFlatten, flattenedScrapeAnnotated); diff != nil {
		t.Errorf("Mk8s Metrics Scrape Annotated was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sLogsAddOn(t *testing.T) {

	_, expectedAddOn, expectedFlatten := generateTestMk8sLogsAddOn()
	flattenedAddOn := flattenMk8sLogsAddOn(expectedAddOn)

	if diff := deep.Equal(expectedFlatten, flattenedAddOn); diff != nil {
		t.Errorf("Mk8s Logs Add On was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sNvidiaAddOn(t *testing.T) {

	_, expectedAddOn, expectedFlatten := generateTestMk8sNvidiaAddOn()
	flattenedAddOn := flattenMk8sNvidiaAddOn(expectedAddOn)

	if diff := deep.Equal(expectedFlatten, flattenedAddOn); diff != nil {
		t.Errorf("Mk8s Nvidia Add On was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAwsAddOn(t *testing.T) {

	_, expectedAddOn, expectedFlatten := generateTestMk8sAwsAddOn()
	flattenedAddOn := flattenMk8sAwsAddOn(expectedAddOn)

	if diff := deep.Equal(expectedFlatten, flattenedAddOn); diff != nil {
		t.Errorf("Mk8s AWS Add On was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAzureAcrAddOn(t *testing.T) {

	_, expectedAddOn, expectedFlatten := generateTestMk8sAzureAcrAddOn()
	flattenedAddOn := flattenMk8sAzureAcrAddOn(expectedAddOn)

	if diff := deep.Equal(expectedFlatten, flattenedAddOn); diff != nil {
		t.Errorf("Mk8s Azure ACR Add On was not flattened correctly. Diff: %s", diff)
	}
}

/*** Generate ***/

// Build //

func generateTestMk8sFirewall() (*[]client.Mk8sFirewallRule, *[]client.Mk8sFirewallRule, []interface{}) {

	sourceCidr := "192.168.1.255"
	description := "hello world"

	flattened := generateFlatTestMk8sFirewall(sourceCidr, description)
	firewall := buildMk8sFirewall(flattened)
	expectedFirewall := []client.Mk8sFirewallRule{
		{
			SourceCIDR:  &sourceCidr,
			Description: &description,
		},
	}

	return firewall, &expectedFirewall, flattened
}

func generateTestMk8sProvider(provider string) *client.Mk8sProvider {

	output := client.Mk8sProvider{}

	switch provider {
	case "generic":
		generated, _, _ := generateTestMk8sGenericProvider()
		output.Generic = generated
	case "hetzner":
		generated, _, _ := generateTestMk8sHetznerProvider()
		output.Hetzner = generated
	case "aws":
		generated, _, _ := generateTestMk8sAwsProvider()
		output.Aws = generated
	}

	return &output
}

func generateTestMk8sAddOns(providerName string) (*client.Mk8sSpecAddOns, *client.Mk8sSpecAddOns, []interface{}) {

	dashboard := true
	azureWorkloadIdentity, _, flattenedAzureWorkloadIdentity := generateTestMk8sAzureWorkloadIdentityAddOn()
	awsWorkloadIdentity := true
	localPathStorage := true
	metrics, _, flattenedMetrics := generateTestMk8sMetricsAddOn()
	logs, _, flattenedLogs := generateTestMk8sLogsAddOn()
	nvidia, _, flattenedNvidia := generateTestMk8sNvidiaAddOn()
	azureAcr, _, flattenedAzureAcr := generateTestMk8sAzureAcrAddOn()
	sysbox := true

	var awsEfs *client.Mk8sAwsAddOnConfig
	var flattenedAwsEfs []interface{}

	var awsEcr *client.Mk8sAwsAddOnConfig
	var flattenedAwsEcr []interface{}

	var awsElb *client.Mk8sAwsAddOnConfig
	var flattenedAwsElb []interface{}

	if providerName == "aws" {
		awsEfs, _, flattenedAwsEfs = generateTestMk8sAwsAddOn()
		awsEcr, _, flattenedAwsEcr = generateTestMk8sAwsAddOn()
		awsElb, _, flattenedAwsElb = generateTestMk8sAwsAddOn()
	}

	flattened := generateFlatTestMk8sAddOns(dashboard, flattenedAzureWorkloadIdentity, awsWorkloadIdentity, localPathStorage, flattenedMetrics, flattenedLogs, flattenedNvidia, flattenedAwsEfs, flattenedAwsEcr, flattenedAwsElb, flattenedAzureAcr, sysbox)
	addOns := buildMk8sAddOns(flattened)
	expectedAddOns := client.Mk8sSpecAddOns{
		Dashboard:             &client.Mk8sNonCustomizableAddonConfig{},
		AzureWorkloadIdentity: azureWorkloadIdentity,
		AwsWorkloadIdentity:   &client.Mk8sNonCustomizableAddonConfig{},
		LocalPathStorage:      &client.Mk8sNonCustomizableAddonConfig{},
		Metrics:               metrics,
		Logs:                  logs,
		Nvidia:                nvidia,
		AwsEFS:                awsEfs,
		AwsECR:                awsEcr,
		AwsELB:                awsElb,
		AzureACR:              azureAcr,
		Sysbox:                &client.Mk8sNonCustomizableAddonConfig{},
	}

	return addOns, &expectedAddOns, flattened
}

// Providers

func generateTestMk8sGenericProvider() (*client.Mk8sGenericProvider, *client.Mk8sGenericProvider, []interface{}) {

	location := "aws-eu-central-1"
	networking, _, flattenedNetworking := generateTestMk8sNetworking()
	nodePools, _, flattenedNodePools := generateTestMk8sGenericNodePools()

	flattened := generateFlatTestMk8sGenericProvider(location, flattenedNetworking, flattenedNodePools)
	generic := buildMk8sGenericProvider(flattened)
	expectedGeneric := client.Mk8sGenericProvider{
		Location:   &location,
		Networking: networking,
		NodePools:  nodePools,
	}

	return generic, &expectedGeneric, flattened
}

func generateTestMk8sHetznerProvider() (*client.Mk8sHetznerProvider, *client.Mk8sHetznerProvider, []interface{}) {

	region := "fsn1"
	hetznerLabels := map[string]interface{}{
		"hello": "world",
	}
	networking, _, flattenedNetworking := generateTestMk8sNetworking()
	preInstallScript := "#! echo hello world"
	tokenSecretLink := "/org/terraform-test-org/secret/hetzner"
	networkId := "2808575"
	nodePools, _, flattenedNodePools := generateTestMk8sHetznerNodePools()
	dedicatedServerNodePools, _, flattenedDedicatedServerNodePools := generateTestMk8sGenericNodePools()
	image := "centos-7"
	sshKey := "10925607"
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sHetznerProvider(region, hetznerLabels, flattenedNetworking, preInstallScript, tokenSecretLink, networkId, flattenedNodePools, flattenedDedicatedServerNodePools, image, sshKey, flattenedAutoscaler)
	hetzner := buildMk8sHetznerProvider(flattened)
	expectedHetzner := client.Mk8sHetznerProvider{
		Region:                   &region,
		HetznerLabels:            &hetznerLabels,
		Networking:               networking,
		PreInstallScript:         &preInstallScript,
		TokenSecretLink:          &tokenSecretLink,
		NetworkId:                &networkId,
		NodePools:                nodePools,
		DedicatedServerNodePools: dedicatedServerNodePools,
		Image:                    &image,
		SshKey:                   &sshKey,
		Autoscaler:               autoscaler,
	}

	return hetzner, &expectedHetzner, flattened
}

func generateTestMk8sAwsProvider() (*client.Mk8sAwsProvider, *client.Mk8sAwsProvider, []interface{}) {

	region := "us-west-2"
	skipCreateRoles := false
	networking, _, flattenedNetworking := generateTestMk8sNetworking()
	preInstallScript := "#! echo hello world"
	image, _, flattenedImage := generateTestMk8sAwsAmi("recommended")
	deployRoleArn := "arn:aws:iam::989132402664:role/cpln-mk8s-terraform-test-org"
	vpcId := "vpc-087b3e0f680a7e91e"
	keyPair := "debug-eks"
	diskEncryptionKeyArn := "arn:aws:kms:us-west-2:989132402664:key/2e9f25ea-efb4-49bf-ae39-007be298726d"
	securityGroupIds := []string{"sg-0f659b1b0711edce1"}
	nodePools, _, flattenedNodePools := generateTestMk8sAwsNodePools()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sAwsProvider(region, skipCreateRoles, flattenedNetworking, preInstallScript, flattenedImage, deployRoleArn, vpcId, keyPair, diskEncryptionKeyArn, securityGroupIds, flattenedNodePools, flattenedAutoscaler)
	aws := buildMk8sAwsProvider(flattened)
	expectedAws := client.Mk8sAwsProvider{
		Region:               &region,
		SkipCreateRoles:      &skipCreateRoles,
		Networking:           networking,
		PreInstallScript:     &preInstallScript,
		Image:                image,
		DeployRoleArn:        &deployRoleArn,
		VpcId:                &vpcId,
		KeyPair:              &keyPair,
		DiskEncryptionKeyArn: &diskEncryptionKeyArn,
		SecurityGroupIds:     &securityGroupIds,
		NodePools:            nodePools,
		Autoscaler:           autoscaler,
	}

	return aws, &expectedAws, flattened
}

// Node Pools

func generateTestMk8sGenericNodePools() (*[]client.Mk8sGenericPool, *[]client.Mk8sGenericPool, []interface{}) {

	name := "my-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()

	flattened := generateFlatTestMk8sGenericNodePools(name, labels, flattenedTaints)
	nodePools := buildMk8sGenericNodePools(flattened)
	expectedNodePools := []client.Mk8sGenericPool{
		{
			Name:   &name,
			Labels: &labels,
			Taints: taints,
		},
	}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sHetznerNodePools() (*[]client.Mk8sHetznerPool, *[]client.Mk8sHetznerPool, []interface{}) {

	name := "my-hetzner-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()
	serverType := "cx11"
	overrideImage := "debian-11"
	minSize := 0
	maxSize := 0

	flattened := generateFlatTestMk8sHetznerNodePools(name, labels, flattenedTaints, serverType, overrideImage, minSize, maxSize)
	nodePools := buildMk8sHetznerNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sHetznerPool{
		ServerType:    &serverType,
		OverrideImage: &overrideImage,
		MinSize:       &minSize,
		MaxSize:       &maxSize,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sHetznerPool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sAwsNodePools() (*[]client.Mk8sAwsPool, *[]client.Mk8sAwsPool, []interface{}) {

	name := "my-aws-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()
	instanceTypes := []string{"t4g.nano"}
	overrideImage, _, flattenedOverrideImage := generateTestMk8sAwsAmi("exact")
	bootDiskSize := 20
	minSize := 0
	maxSize := 0
	onDemandBaseCapacity := 0
	onDemandPercentageAboveBaseCapacity := 0
	spotAllocationStrategy := "lowest-price"
	subnetIds := []string{"subnet-077fe72ab6259d9a2"}
	extraSecurityGroupIds := []string{}

	flattened := generateFlatTestMk8sAwsNodePools(name, labels, flattenedTaints, instanceTypes, flattenedOverrideImage, bootDiskSize, minSize, maxSize, onDemandBaseCapacity, onDemandPercentageAboveBaseCapacity, spotAllocationStrategy, subnetIds, extraSecurityGroupIds)
	nodePools := buildMk8sAwsNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sAwsPool{
		InstanceTypes:                       &instanceTypes,
		OverrideImage:                       overrideImage,
		BootDiskSize:                        &bootDiskSize,
		MinSize:                             &minSize,
		MaxSize:                             &maxSize,
		OnDemandBaseCapacity:                &onDemandBaseCapacity,
		OnDemandPercentageAboveBaseCapacity: &onDemandPercentageAboveBaseCapacity,
		SpotAllocationStrategy:              &spotAllocationStrategy,
		SubnetIds:                           &subnetIds,
		ExtraSecurityGroupIds:               &extraSecurityGroupIds,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sAwsPool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

// AWS

func generateTestMk8sAwsAmi(choice string) (*client.Mk8sAwsAmi, *client.Mk8sAwsAmi, []interface{}) {

	var recommended *string
	var exact *string

	if choice == "recommended" {
		recommended = GetString("amazon/al2023")
	} else if choice == "exact" {
		exact = GetString("ami-123")
	}

	flattened := generateFlatTestMk8sAwsAmi(recommended, exact)
	ami := buildMk8sAwsAmi(flattened)
	expectedAmi := client.Mk8sAwsAmi{
		Recommended: recommended,
		Exact:       exact,
	}

	return ami, &expectedAmi, flattened
}

// Common

func generateTestMk8sNetworking() (*client.Mk8sNetworkingConfig, *client.Mk8sNetworkingConfig, []interface{}) {

	serviceNetwork := "10.43.0.0/16"
	podNetwork := "10.42.0.0/16"

	flattened := generateFlatTestMk8sNetworking(serviceNetwork, podNetwork)
	networking := buildMk8sNetworking(flattened)
	expectedNetworking := client.Mk8sNetworkingConfig{
		ServiceNetwork: &serviceNetwork,
		PodNetwork:     &podNetwork,
	}

	return networking, &expectedNetworking, flattened
}

func generateTestMk8sTaints() (*[]client.Mk8sTaint, *[]client.Mk8sTaint, []interface{}) {

	key := "hello"
	value := "world"
	effect := "NoSchedule"

	flattened := generateFlatTestMk8sTaints(key, value, effect)
	taints := buildMk8sTaints(flattened)
	expectedTaints := []client.Mk8sTaint{
		{
			Key:    &key,
			Value:  &value,
			Effect: &effect,
		},
	}

	return taints, &expectedTaints, flattened
}

func generateTestMk8sAutoscaler() (*client.Mk8sAutoscalerConfig, *client.Mk8sAutoscalerConfig, []interface{}) {

	expander := []string{"most-pods"}
	unneededTime := "10m"
	unreadyTime := "20m"
	utilizationThreshold := 0.7

	flattened := generateFlatTestMk8sAutoscaler(expander, unneededTime, unreadyTime, utilizationThreshold)
	autoscaler := buildMk8sAutoscaler(flattened)
	expectedAutoscaler := client.Mk8sAutoscalerConfig{
		Expander:             &expander,
		UnneededTime:         &unneededTime,
		UnreadyTime:          &unreadyTime,
		UtilizationThreshold: &utilizationThreshold,
	}

	return autoscaler, &expectedAutoscaler, flattened
}

// Add Ons

func generateTestMk8sAzureWorkloadIdentityAddOn() (*client.Mk8sAzureWorkloadIdentityAddOnConfig, *client.Mk8sAzureWorkloadIdentityAddOnConfig, []interface{}) {

	tenantId := "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"

	flattened := generateFlatTestMk8sAzureWorkloadIdentityAddOn(tenantId)
	azureWorkloadIdentityAddOn := buildMk8sAzureWorkloadIdentityAddOn(flattened)
	expectedAzureWorkloadIdentityAddOn := client.Mk8sAzureWorkloadIdentityAddOnConfig{
		TenantId: &tenantId,
	}

	return azureWorkloadIdentityAddOn, &expectedAzureWorkloadIdentityAddOn, flattened
}

func generateTestMk8sMetricsAddOn() (*client.Mk8sMetricsAddOnConfig, *client.Mk8sMetricsAddOnConfig, []interface{}) {

	kubeState := true
	coreDns := true
	kubelet := true
	apiServer := true
	nodeExporter := true
	cadvisor := true
	scrapeAnnotated, _, flattenedScrapeAnnotated := generateTestMk8sMetricsScrapeAnnotated()

	flattened := generateFlatTestMk8sMetricsAddOn(kubeState, coreDns, kubelet, apiServer, nodeExporter, cadvisor, flattenedScrapeAnnotated)
	metrics := buildMk8sMetricsAddOn(flattened)
	expectedMetrics := client.Mk8sMetricsAddOnConfig{
		KubeState:       &kubeState,
		CoreDns:         &coreDns,
		Kubelet:         &kubelet,
		Apiserver:       &apiServer,
		NodeExporter:    &nodeExporter,
		Cadvisor:        &cadvisor,
		ScrapeAnnotated: scrapeAnnotated,
	}

	return metrics, &expectedMetrics, flattened
}

func generateTestMk8sMetricsScrapeAnnotated() (*client.Mk8sMetricsScrapeAnnotated, *client.Mk8sMetricsScrapeAnnotated, []interface{}) {

	intervalSeconds := 30
	includeNamespaces := "^\\d+$"
	excludeNamespaces := "^[a-z]$"
	retainLabels := "^\\w+$"

	flattened := generateFlatTestMk8sMetricsScrapeAnnotated(intervalSeconds, includeNamespaces, excludeNamespaces, retainLabels)
	scrapeAnnotated := buildMk8sMetricsScrapeAnnotated(flattened)
	expectedScrapeAnnotated := client.Mk8sMetricsScrapeAnnotated{
		IntervalSeconds:   &intervalSeconds,
		IncludeNamespaces: &includeNamespaces,
		ExcludeNamespaces: &excludeNamespaces,
		RetainLabels:      &retainLabels,
	}

	return scrapeAnnotated, &expectedScrapeAnnotated, flattened
}

func generateTestMk8sLogsAddOn() (*client.Mk8sLogsAddOnConfig, *client.Mk8sLogsAddOnConfig, []interface{}) {

	auditEnabled := true
	includeNamespaces := "^\\d+$"
	excludeNamespaces := "^[a-z]$"

	flattened := generateFlatTestMk8sLogsAddOn(auditEnabled, includeNamespaces, excludeNamespaces)
	logs := buildMk8sLogsAddOn(flattened)
	expectedLogs := client.Mk8sLogsAddOnConfig{
		AuditEnabled:      &auditEnabled,
		IncludeNamespaces: &includeNamespaces,
		ExcludeNamespaces: &excludeNamespaces,
	}

	return logs, &expectedLogs, flattened
}

func generateTestMk8sNvidiaAddOn() (*client.Mk8sNvidiaAddOnConfig, *client.Mk8sNvidiaAddOnConfig, []interface{}) {

	taintGpuNodes := true

	flattened := generateFlatTestMk8sNvidiaAddOn(taintGpuNodes)
	nvidia := buildMk8sNvidiaAddOn(flattened)
	expectedNvidia := client.Mk8sNvidiaAddOnConfig{
		TaintGPUNodes: &taintGpuNodes,
	}

	return nvidia, &expectedNvidia, flattened
}

func generateTestMk8sAwsAddOn() (*client.Mk8sAwsAddOnConfig, *client.Mk8sAwsAddOnConfig, []interface{}) {

	roleArn := "arn:aws:iam::123456789012:role/my-custom-role"

	flattened := generateFlatTestMk8sAwsAddOn(roleArn)
	aws := buildMk8sAwsAddOn(flattened)
	expectedAws := client.Mk8sAwsAddOnConfig{
		RoleArn: &roleArn,
	}

	return aws, &expectedAws, flattened
}

func generateTestMk8sAzureAcrAddOn() (*client.Mk8sAzureACRAddOnConfig, *client.Mk8sAzureACRAddOnConfig, []interface{}) {

	clientId := "4e25b134-160b-4a9d-b392-13b381ced5ef"

	flattened := generateFlatTestMk8sAzureAcrAddOn(clientId)
	azureAcr := buildMk8sAzureAcrAddOn(flattened)
	expectedAzureAcr := client.Mk8sAzureACRAddOnConfig{
		ClientId: &clientId,
	}

	return azureAcr, &expectedAzureAcr, flattened
}

// Flatten //

func generateFlatTestMk8sFirewall(sourceCidr string, description string) []interface{} {

	spec := map[string]interface{}{
		"source_cidr": sourceCidr,
		"description": description,
	}

	return []interface{}{
		spec,
	}
}

// Providers

func generateFlatTestMk8sGenericProvider(location string, networking []interface{}, nodePools []interface{}) []interface{} {

	spec := map[string]interface{}{
		"location":   location,
		"networking": networking,
		"node_pool":  nodePools,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sHetznerProvider(region string, hetznerLabels map[string]interface{}, networking []interface{}, preInstallScript string, tokenSecretLink string, networkId string, nodePools []interface{}, dedicatedServerNodePools []interface{}, image string, sshKey string, autoscaler []interface{}) []interface{} {

	spec := map[string]interface{}{
		"region":                     region,
		"hetzner_labels":             hetznerLabels,
		"networking":                 networking,
		"pre_install_script":         preInstallScript,
		"token_secret_link":          tokenSecretLink,
		"network_id":                 networkId,
		"node_pool":                  nodePools,
		"dedicated_server_node_pool": dedicatedServerNodePools,
		"image":                      image,
		"ssh_key":                    sshKey,
		"autoscaler":                 autoscaler,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAwsProvider(region string, skipCreateRoles bool, networking []interface{}, preInstallScript string, image []interface{}, deployRoleArn string, vpcId string, keyPair string, diskEncryptionKeyArn string, securityGroupIds []string, nodePools []interface{}, autoscaler []interface{}) []interface{} {

	spec := map[string]interface{}{
		"region":                  region,
		"skip_create_roles":       skipCreateRoles,
		"networking":              networking,
		"pre_install_script":      preInstallScript,
		"image":                   image,
		"deploy_role_arn":         deployRoleArn,
		"vpc_id":                  vpcId,
		"key_pair":                keyPair,
		"disk_encryption_key_arn": diskEncryptionKeyArn,
		"security_group_ids":      ConvertStringSliceToSet(securityGroupIds),
		"node_pool":               nodePools,
		"autoscaler":              autoscaler,
	}

	return []interface{}{
		spec,
	}
}

// Node Pools

func generateFlatTestMk8sGenericNodePools(name string, labels map[string]interface{}, taints []interface{}) []interface{} {

	spec := map[string]interface{}{
		"name":   name,
		"labels": labels,
		"taint":  taints,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sHetznerNodePools(name string, labels map[string]interface{}, taints []interface{}, serverType string, overrideImage string, minSize int, maxSize int) []interface{} {

	spec := map[string]interface{}{
		"name":           name,
		"labels":         labels,
		"taint":          taints,
		"server_type":    serverType,
		"override_image": overrideImage,
		"min_size":       minSize,
		"max_size":       maxSize,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAwsNodePools(name string, labels map[string]interface{}, taints []interface{}, instanceTypes []string, overrideImage []interface{}, bootDiskSize int, minSize int, maxSize int, onDemandBaseCapacity int, onDemandPercentageAboveBaseCapacity int, spotAllocationStrategy string, subnetIds []string, extraSecurityGroupIds []string) []interface{} {

	spec := map[string]interface{}{
		"name":                    name,
		"labels":                  labels,
		"taint":                   taints,
		"instance_types":          ConvertStringSliceToSet(instanceTypes),
		"override_image":          overrideImage,
		"boot_disk_size":          bootDiskSize,
		"min_size":                minSize,
		"max_size":                maxSize,
		"on_demand_base_capacity": onDemandBaseCapacity,
		"on_demand_percentage_above_base_capacity": onDemandPercentageAboveBaseCapacity,
		"spot_allocation_strategy":                 spotAllocationStrategy,
		"subnet_ids":                               ConvertStringSliceToSet(subnetIds),
		"extra_security_group_ids":                 ConvertStringSliceToSet(extraSecurityGroupIds),
	}

	return []interface{}{
		spec,
	}
}

// AWS

func generateFlatTestMk8sAwsAmi(recommended *string, exact *string) []interface{} {

	spec := make(map[string]interface{})

	if recommended != nil {
		spec["recommended"] = *recommended
	}

	if exact != nil {
		spec["exact"] = *exact
	}

	return []interface{}{
		spec,
	}
}

// Common

func generateFlatTestMk8sNetworking(serviceNetwork string, podNetwork string) []interface{} {

	spec := map[string]interface{}{
		"service_network": serviceNetwork,
		"pod_network":     podNetwork,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sTaints(key string, value string, effect string) []interface{} {

	spec := map[string]interface{}{
		"key":    key,
		"value":  value,
		"effect": effect,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAutoscaler(expander []string, unneededTime string, unreadyTime string, utilizationThreshold float64) []interface{} {

	spec := map[string]interface{}{
		"expander":              ConvertStringSliceToSet(expander),
		"unneeded_time":         unneededTime,
		"unready_time":          unreadyTime,
		"utilization_threshold": utilizationThreshold,
	}

	return []interface{}{
		spec,
	}
}

// Add Ons

func generateFlatTestMk8sAddOns(dashboard bool, azureWorkloadIdentity []interface{}, awsWorkloadIdentity bool, localPathStorage bool, metrics []interface{}, logs []interface{}, nvidia []interface{}, awsEfs []interface{}, awsEcr []interface{}, awsElb []interface{}, azureAcr []interface{}, sysbox bool) []interface{} {

	spec := map[string]interface{}{
		"dashboard":               dashboard,
		"azure_workload_identity": azureWorkloadIdentity,
		"aws_workload_identity":   awsWorkloadIdentity,
		"local_path_storage":      localPathStorage,
		"metrics":                 metrics,
		"logs":                    logs,
		"nvidia":                  nvidia,
		"aws_efs":                 awsEfs,
		"aws_ecr":                 awsEcr,
		"aws_elb":                 awsElb,
		"azure_acr":               azureAcr,
		"sysbox":                  sysbox,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAzureWorkloadIdentityAddOn(tenantId string) []interface{} {

	spec := map[string]interface{}{
		"tenant_id": tenantId,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sMetricsAddOn(kubeState bool, coreDns bool, kubelet bool, apiServer bool, nodeExporter bool, cadvisor bool, scrapeAnnotated []interface{}) []interface{} {

	spec := map[string]interface{}{
		"kube_state":       kubeState,
		"core_dns":         coreDns,
		"kubelet":          kubelet,
		"api_server":       apiServer,
		"node_exporter":    nodeExporter,
		"cadvisor":         cadvisor,
		"scrape_annotated": scrapeAnnotated,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sMetricsScrapeAnnotated(intervalSeconds int, includeNamespaces string, excludeNamespaces string, retainLabels string) []interface{} {

	spec := map[string]interface{}{
		"interval_seconds":   intervalSeconds,
		"include_namespaces": includeNamespaces,
		"exclude_namespaces": excludeNamespaces,
		"retain_labels":      retainLabels,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sLogsAddOn(auditEnabled bool, includeNamespaces string, excludeNamespaces string) []interface{} {

	spec := map[string]interface{}{
		"audit_enabled":      auditEnabled,
		"include_namespaces": includeNamespaces,
		"exclude_namespaces": excludeNamespaces,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sNvidiaAddOn(taintGpuNodes bool) []interface{} {

	spec := map[string]interface{}{
		"taint_gpu_nodes": taintGpuNodes,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAwsAddOn(roleArn string) []interface{} {

	spec := map[string]interface{}{
		"role_arn": roleArn,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAzureAcrAddOn(clientId string) []interface{} {

	spec := map[string]interface{}{
		"client_id": clientId,
	}

	return []interface{}{
		spec,
	}
}
