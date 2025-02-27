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
			// Create
			{
				Config: testAccControlPlaneMk8sGenericProvider(name+"-generic", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.generic", name+"-generic", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "generic", ""),
					resource.TestCheckResourceAttr("cpln_mk8s.generic", "name", name+"-generic"),
					resource.TestCheckResourceAttr("cpln_mk8s.generic", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sHetznerProvider(name+"-hetzner", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.hetzner", name+"-hetzner", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "hetzner", ""),
					resource.TestCheckResourceAttr("cpln_mk8s.hetzner", "name", name+"-hetzner"),
					resource.TestCheckResourceAttr("cpln_mk8s.hetzner", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sAwsProvider(name+"-aws", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.aws", name+"-aws", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "aws", ""),
					resource.TestCheckResourceAttr("cpln_mk8s.aws", "name", name+"-aws"),
					resource.TestCheckResourceAttr("cpln_mk8s.aws", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sLinodeProvider(name+"-linode", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.linode", name+"-linode", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "linode", ""),
					resource.TestCheckResourceAttr("cpln_mk8s.linode", "name", name+"-linode"),
					resource.TestCheckResourceAttr("cpln_mk8s.linode", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sOblivusProvider(name+"-oblivus", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.oblivus", name+"-oblivus", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "oblivus", ""),
					resource.TestCheckResourceAttr("cpln_mk8s.oblivus", "name", name+"-oblivus"),
					resource.TestCheckResourceAttr("cpln_mk8s.oblivus", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sLambdalabsProvider(name+"-lambdalabs", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.lambdalabs", name+"-lambdalabs", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "lambdalabs", ""),
					resource.TestCheckResourceAttr("cpln_mk8s.lambdalabs", "name", name+"-lambdalabs"),
					resource.TestCheckResourceAttr("cpln_mk8s.lambdalabs", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sPaperspaceProvider(name+"-paperspace", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.paperspace", name+"-paperspace", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "paperspace", ""),
					resource.TestCheckResourceAttr("cpln_mk8s.paperspace", "name", name+"-paperspace"),
					resource.TestCheckResourceAttr("cpln_mk8s.paperspace", "description", description),
				),
			},
			{
				Config: testAccControlPlaneEphemeralProvider(name+"-ephemeral", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.ephemeral", name+"-ephemeral", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "ephemeral", "no-sysbox"),
					resource.TestCheckResourceAttr("cpln_mk8s.ephemeral", "name", name+"-ephemeral"),
					resource.TestCheckResourceAttr("cpln_mk8s.ephemeral", "description", description),
				),
			},
			{
				Config: testAccControlPlaneTritonProvider(name+"-triton", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.triton", name+"-triton", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "triton", "gateway"),
					resource.TestCheckResourceAttr("cpln_mk8s.triton", "name", name+"-triton"),
					resource.TestCheckResourceAttr("cpln_mk8s.triton", "description", description),
				),
			},
			// Update
			{
				Config: testAccControlPlaneMk8sHetznerProviderUpdate(name+"-hetzner", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.hetzner", name+"-hetzner", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "hetzner", "case1"),
					resource.TestCheckResourceAttr("cpln_mk8s.hetzner", "name", name+"-hetzner"),
					resource.TestCheckResourceAttr("cpln_mk8s.hetzner", "description", description),
				),
			},
			{
				Config: testAccControlPlaneMk8sLambdalabsProviderUpdate(name+"-lambdalabs", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.lambdalabs", name+"-lambdalabs", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "lambdalabs", "case1"),
					resource.TestCheckResourceAttr("cpln_mk8s.lambdalabs", "name", name+"-lambdalabs"),
					resource.TestCheckResourceAttr("cpln_mk8s.lambdalabs", "description", description),
				),
			},
			{
				Config: testAccControlPlaneTritonProviderUpdate(name+"-triton", description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMk8sExists("cpln_mk8s.triton", name+"-triton", &mk8s),
					testAccCheckControlPlaneMk8sAttributes(&mk8s, "triton", "manual"),
					resource.TestCheckResourceAttr("cpln_mk8s.triton", "name", name+"-triton"),
					resource.TestCheckResourceAttr("cpln_mk8s.triton", "description", description),
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

func testAccCheckControlPlaneMk8sAttributes(mk8s *client.Mk8s, providerName string, update string) resource.TestCheckFunc {
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
		expectedProvider := generateTestMk8sProvider(providerName, update)

		if diff := deep.Equal(mk8s.Spec.Provider, expectedProvider); diff != nil {
			return fmt.Errorf("Mk8s Provider %s does not match. Diff: %s", providerName, diff)
		}

		// Add Ons
		expectedAddOns, _, _ := generateTestMk8sAddOns(providerName, update)

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

// SECTION Acceptance Tests

// SECTION Create

func testAccControlPlaneMk8sGenericProvider(name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_mk8s" "generic" {
		
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
					include_namespaces = "^elastic"
					exclude_namespaces  = "^elastic"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^elastic"
				exclude_namespaces  = "^elastic"
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
					exclude_namespaces  = "^elastic"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^elastic"
				exclude_namespaces  = "^elastic"
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
			"cpln/ignore"       = "true"
		}

		version = "1.28.4"

		firewall {
			source_cidr = "192.168.1.255"
			description = "hello world"
		}

		aws_provider {

			region            = "eu-central-1"

			aws_tags = {
				hello = "world"
			}

			skip_create_roles = false

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

			nvidia {
				taint_gpu_nodes = true
			}

			aws_efs {
				role_arn = "arn:aws:iam::123456789012:role/my-custom-role"
			}

			aws_ecr {
				role_arn = "arn:aws:iam::123456789012:role/my-custom-role"
			}

			aws_elb {}

			azure_acr {
				client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
			}
			
			sysbox = true
		}
	}
	`, name, description)
}

func testAccControlPlaneMk8sLinodeProvider(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "linode" {
		
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
			image              = "linode/ubuntu24.04"
			vpc_id             = "93666"
			firewall_id        = "168425"
			pre_install_script = "#! echo hello world"

			authorized_users = ["juliancpln"]

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

func testAccControlPlaneMk8sOblivusProvider(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "oblivus" {
		
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
				min_size 	     = 0
				max_size 	     = 0
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
					include_namespaces = "^elastic"
					exclude_namespaces  = "^elastic"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^elastic"
				exclude_namespaces  = "^elastic"
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

func testAccControlPlaneMk8sLambdalabsProvider(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "lambdalabs" {
		
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
				name = "my-lambdalabs-node-pool"

				labels = {
					hello = "world"
				}

				taint {
					key    = "hello"
					value  = "world"
					effect = "NoSchedule"
				}

				instance_type = "cpu_4x_general"
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
					include_namespaces = "^elastic"
					exclude_namespaces  = "^elastic"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^elastic"
				exclude_namespaces  = "^elastic"
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

func testAccControlPlaneMk8sPaperspaceProvider(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "paperspace" {
		
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
			pre_install_script = "#! echo hello world"
			shared_drives      = ["california"]
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

func testAccControlPlaneEphemeralProvider(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "ephemeral" {
		
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
				name = "my-node-pool"

				labels = {
					hello = "world"
				}

				taint {
					key    = "hello"
					value  = "world"
					effect = "NoSchedule"
				}

				count  = 1
				arch   = "arm64"
				flavor = "debian"
				cpu    = "50m"
				memory = "128Mi"
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
					exclude_namespaces  = "^elastic"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^elastic"
				exclude_namespaces  = "^elastic"
			}

			nvidia {
				taint_gpu_nodes = true
			}

			azure_acr {
				client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
			}
		}
	}
	`, name, description)
}

func testAccControlPlaneTritonProvider(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "triton" {
		
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
				name                = "my-triton-node-pool"
				package_id          = "da311341-b42b-45a8-9386-78ede625d0a4"
				override_image_id   = "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
				public_network_id   = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
				min_size            = 0
				max_size            = 0

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

// !SECTION

// SECTION Update

func testAccControlPlaneMk8sHetznerProviderUpdate(name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_mk8s" "hetzner" {
		
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

			floating_ip_selector = {
				floating_ip_1 = "123.45.67.89"
				floating_ip_2 = "98.76.54.32"
			}
		}

		add_ons {
			dashboard = false

			azure_workload_identity {
				tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
			}

			aws_workload_identity = false
			local_path_storage    = false

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
					exclude_namespaces  = "^elastic"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^elastic"
				exclude_namespaces  = "^elastic"
			}

			nvidia {
				taint_gpu_nodes = true
			}

			azure_acr {
				client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
			}

			sysbox = false
		}
	}
	`, name, description)
}

func testAccControlPlaneMk8sLambdalabsProviderUpdate(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "lambdalabs" {
		
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
				name = "my-lambdalabs-node-pool"

				labels = {
					hello = "world"
				}

				taint {
					key    = "hello"
					value  = "world"
					effect = "NoSchedule"
				}

				min_size      = 0
				max_size      = 0
				instance_type = "cpu_4x_general"
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
				expander 	  		  = ["most-pods"]
				unneeded_time         = "10m"
				unready_time  		  = "20m"
				utilization_threshold = 0.7
			}
		}

		add_ons {
			dashboard = false

			azure_workload_identity {
				tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
			}

			aws_workload_identity = false
			local_path_storage    = false

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
					exclude_namespaces  = "^elastic"
					retain_labels      = "^\\w+$"
				}
			}

			logs {
				audit_enabled      = true
				include_namespaces = "^elastic"
				exclude_namespaces  = "^elastic"
			}

			nvidia {
				taint_gpu_nodes = true
			}

			azure_acr {
				client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
			}

			sysbox = false
		}
	}
	`, name, description)
}

func testAccControlPlaneTritonProviderUpdate(name string, description string) string {
	return fmt.Sprintf(`

	resource "cpln_mk8s" "triton" {
		
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
				}
			}

			node_pool {
				name                = "my-triton-node-pool"
				package_id          = "da311341-b42b-45a8-9386-78ede625d0a4"
				override_image_id   = "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
				public_network_id   = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
				min_size            = 0
				max_size            = 0

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

// !SECTION
// !SECTION

// SECTION Unit Tests

// SECTION Build

func TestControlPlane_BuildMk8sFirewall(t *testing.T) {

	firewall, expectedFirewall, _ := generateTestMk8sFirewall()

	if diff := deep.Equal(firewall, expectedFirewall); diff != nil {
		t.Errorf("Mk8s Firewall was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sAddOns(t *testing.T) {

	addOns, expectedAddOns, _ := generateTestMk8sAddOns("aws", "")

	if diff := deep.Equal(addOns, expectedAddOns); diff != nil {
		t.Errorf("Mk8s AddOns was not built correctly, Diff: %s", diff)
	}
}

// SECTION Providers

func TestControlPlane_BuildMk8sGenericProvider(t *testing.T) {

	generic, expectedGeneric, _ := generateTestMk8sGenericProvider()

	if diff := deep.Equal(generic, expectedGeneric); diff != nil {
		t.Errorf("Mk8s Generic Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sHetznerProvider(t *testing.T) {

	hetzner, expectedHetzner, _ := generateTestMk8sHetznerProvider("")

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

func TestControlPlane_BuildMk8sLinodeProvider(t *testing.T) {

	linode, expectedLinode, _ := generateTestMk8sLinodeProvider()

	if diff := deep.Equal(linode, expectedLinode); diff != nil {
		t.Errorf("Mk8s Linode Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sOblivusProvider(t *testing.T) {

	oblivus, expectedOblivus, _ := generateTestMk8sOblivusProvider()

	if diff := deep.Equal(oblivus, expectedOblivus); diff != nil {
		t.Errorf("Mk8s Oblivus Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sLambdalabsProvider(t *testing.T) {

	lambdalabs, expectedLambdalabs, _ := generateTestMk8sLambdalabsProvider()

	if diff := deep.Equal(lambdalabs, expectedLambdalabs); diff != nil {
		t.Errorf("Mk8s Lambdalabs Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sPaperspaceProvider(t *testing.T) {

	paperspace, expectedPaperspace, _ := generateTestMk8sPaperspaceProvider()

	if diff := deep.Equal(paperspace, expectedPaperspace); diff != nil {
		t.Errorf("Mk8s Paperspace Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sEphemeralProvider(t *testing.T) {

	ephemeral, expectedEphemeral, _ := generateTestMk8sEphemeralProvider()

	if diff := deep.Equal(ephemeral, expectedEphemeral); diff != nil {
		t.Errorf("Mk8s Ephemeral Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sTritonProvider(t *testing.T) {

	triton, expectedTriton, _ := generateTestMk8sTritonProvider("gateway")

	if diff := deep.Equal(triton, expectedTriton); diff != nil {
		t.Errorf("Mk8s Triton Provider was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sDigitalOceanProvider(t *testing.T) {

	digitalOcean, expectedDigitalOcean, _ := generateTestMk8sDigitalOceanProvider()

	if diff := deep.Equal(digitalOcean, expectedDigitalOcean); diff != nil {
		t.Errorf("Mk8s Digital Ocean Provider was not built correctly, Diff: %s", diff)
	}
}

// !SECTION

// SECTION Node Pools

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

func TestControlPlane_BuildMk8sLinodeNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sLinodeNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Linode Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sOblivusNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sOblivusNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Oblivus Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sLambdalabsNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sLambdalabsNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Lambdalabs Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sPaperspaceNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sPaperspaceNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Paperspace Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sEphemeralNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sEphemeralNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Ephemeral Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sTritonNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sTritonNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Triton Node Pools was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sDigitalOceanNodePools(t *testing.T) {

	nodePools, expectedNodePools, _ := generateTestMk8sDigitalOceanNodePools()

	if diff := deep.Equal(nodePools, expectedNodePools); diff != nil {
		t.Errorf("Mk8s Digital Ocean Node Pools was not built correctly, Diff: %s", diff)
	}
}

// !SECTION

// SECTION AWS

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

func TestControlPlane_BuildMk8sAwsDeployRoleChain(t *testing.T) {

	deployRoleChain, expectedDeployRoleChain, _ := generateTestMk8sAwsDeployRoleChain()

	if diff := deep.Equal(deployRoleChain, expectedDeployRoleChain); diff != nil {
		t.Errorf("Mk8s AWS Deploy Role Chain was not built correctly, Diff: %s", diff)
	}
}

// !SECTION

// SECTION Triton

func TestControlPlane_BuildMk8sTritonConnection(t *testing.T) {

	connection, expectedConnection, _ := generateTestMk8sTritonConnection()

	if diff := deep.Equal(connection, expectedConnection); diff != nil {
		t.Errorf("Mk8s Triton Connection was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sTritonLoadBalancerGateway(t *testing.T) {

	loadBalancer, expectedLoadBalancer, _ := generateTestMk8sTritonLoadBalancer("gateway")

	if diff := deep.Equal(loadBalancer, expectedLoadBalancer); diff != nil {
		t.Errorf("Mk8s Triton Load Balancer Gateway was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildMk8sTritonLoadBalancerManual(t *testing.T) {

	loadBalancer, expectedLoadBalancer, _ := generateTestMk8sTritonLoadBalancer("manual")

	if diff := deep.Equal(loadBalancer, expectedLoadBalancer); diff != nil {
		t.Errorf("Mk8s Triton Load Balancer Manual was not built correctly, Diff: %s", diff)
	}
}

// !SECTION

// SECTION Common

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

// !SECTION

// SECTION Add Ons

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

	awsAddOn, expectedAwsAddOn, _ := generateTestMk8sAwsAddOn(true)

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

// !SECTION
// !SECTION

// SECTION Flatten

func TestControlPlane_FlattenMk8sFirewall(t *testing.T) {

	_, expectedFirewall, expectedFlatten := generateTestMk8sFirewall()
	flattenedFirewall := flattenMk8sFirewall(expectedFirewall)

	if diff := deep.Equal(expectedFlatten, flattenedFirewall); diff != nil {
		t.Errorf("Mk8s Firewall was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sAddOns(t *testing.T) {

	_, expectedAddOns, expectedFlatten := generateTestMk8sAddOns("aws", "")
	flattenedAddOns := flattenMk8sAddOns(expectedAddOns)

	if diff := deep.Equal(expectedFlatten, flattenedAddOns); diff != nil {
		t.Errorf("Mk8s Add Ons was not flattened correctly. Diff: %s", diff)
	}
}

// SECTION Providers

func TestControlPlane_FlattenMk8sGenericProvider(t *testing.T) {

	_, expectedGeneric, expectedFlatten := generateTestMk8sGenericProvider()
	flattenedGeneric := flattenMk8sGenericProvider(expectedGeneric)

	if diff := deep.Equal(expectedFlatten, flattenedGeneric); diff != nil {
		t.Errorf("Mk8s Generic Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sHetznerProvider(t *testing.T) {

	_, expectedHetzner, expectedFlatten := generateTestMk8sHetznerProvider("")
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
	expectedFlattenItem["extra_node_policies"] = expectedFlattenItem["extra_node_policies"].(*schema.Set).List()

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

func TestControlPlane_FlattenMk8sLinodeProvider(t *testing.T) {

	_, expectedLinode, expectedFlatten := generateTestMk8sLinodeProvider()
	flattenedLinode := flattenMk8sLinodeProvider(expectedLinode)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})
	expectedFlattenItem["authorized_users"] = expectedFlattenItem["authorized_users"].(*schema.Set).List()
	expectedFlattenItem["authorized_keys"] = expectedFlattenItem["authorized_keys"].(*schema.Set).List()

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedLinode); diff != nil {
		t.Errorf("Mk8s Linode Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sOblivusProvider(t *testing.T) {

	_, expectedOblivus, expectedFlatten := generateTestMk8sOblivusProvider()
	flattenedOblivus := flattenMk8sOblivusProvider(expectedOblivus)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})
	expectedFlattenItem["ssh_keys"] = expectedFlattenItem["ssh_keys"].(*schema.Set).List()

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedOblivus); diff != nil {
		t.Errorf("Mk8s Oblivus Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sLambdalabsProvider(t *testing.T) {

	_, expectedLambdalabs, expectedFlatten := generateTestMk8sLambdalabsProvider()
	flattenedLambdalabs := flattenMk8sLambdalabsProvider(expectedLambdalabs)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedLambdalabs); diff != nil {
		t.Errorf("Mk8s Lambdalabs Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sPaperspaceProvider(t *testing.T) {

	_, expectedPaperspace, expectedFlatten := generateTestMk8sPaperspaceProvider()
	flattenedPaperspace := flattenMk8sPaperspaceProvider(expectedPaperspace)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})

	// Shared Drives
	expectedFlattenItem["shared_drives"] = expectedFlattenItem["shared_drives"].(*schema.Set).List()

	// User IDs
	expectedFlattenItem["user_ids"] = expectedFlattenItem["user_ids"].(*schema.Set).List()

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedPaperspace); diff != nil {
		t.Errorf("Mk8s Paperspace Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sEphemeralProvider(t *testing.T) {

	_, expectedEphemeral, expectedFlatten := generateTestMk8sEphemeralProvider()
	flattenedEphemeral := flattenMk8sEphemeralProvider(expectedEphemeral)

	if diff := deep.Equal(expectedFlatten, flattenedEphemeral); diff != nil {
		t.Errorf("Mk8s Ephemeral Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sTritonProvider(t *testing.T) {

	_, expectedTriton, expectedFlatten := generateTestMk8sTritonProvider("gateway")
	flattenedTriton := flattenMk8sTritonProvider(expectedTriton)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})

	// SSH Keys
	expectedFlattenItem["ssh_keys"] = expectedFlattenItem["ssh_keys"].(*schema.Set).List()

	// Node Pool
	nodePool := expectedFlattenItem["node_pool"].([]interface{})[0].(map[string]interface{})
	nodePool["private_network_ids"] = nodePool["private_network_ids"].(*schema.Set).List()

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedTriton); diff != nil {
		t.Errorf("Mk8s Triton Provider was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sDigitalOceanProvider(t *testing.T) {

	_, expectedDigitalOcean, expectedFlatten := generateTestMk8sDigitalOceanProvider()
	flattenedDigitalOcean := flattenMk8sDigitalOceanProvider(expectedDigitalOcean)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})

	// Digital Ocean Tags
	expectedFlattenItem["digital_ocean_tags"] = expectedFlattenItem["digital_ocean_tags"].(*schema.Set).List()

	// SSH Keys
	expectedFlattenItem["ssh_keys"] = expectedFlattenItem["ssh_keys"].(*schema.Set).List()

	// Extra SSH Keys
	expectedFlattenItem["extra_ssh_keys"] = expectedFlattenItem["extra_ssh_keys"].(*schema.Set).List()

	// Reserved IPs
	expectedFlattenItem["reserved_ips"] = expectedFlattenItem["reserved_ips"].(*schema.Set).List()

	// Autoscaler
	autoscaler := expectedFlattenItem["autoscaler"].([]interface{})[0].(map[string]interface{})
	autoscaler["expander"] = autoscaler["expander"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedDigitalOcean); diff != nil {
		t.Errorf("Mk8s Digital Ocean Provider was not flattened correctly. Diff: %s", diff)
	}
}

// !SECTION

// SECTION Node Pools

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

func TestControlPlane_FlattenMk8sLinodeNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sLinodeNodePools()
	flattenedNodePools := flattenMk8sLinodeNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Linode Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sOblivusNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sOblivusNodePools()
	flattenedNodePools := flattenMk8sOblivusNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Oblivus Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sLambdalabsNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sLambdalabsNodePools()
	flattenedNodePools := flattenMk8sLambdalabsNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Lambdalabs Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sPaperspaceNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sPaperspaceNodePools()
	flattenedNodePools := flattenMk8sPaperspaceNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Paperspace Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sEphemeralNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sEphemeralNodePools()
	flattenedNodePools := flattenMk8sEphemeralNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Ephemeral Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sTritonNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sTritonNodePools()
	flattenedNodePools := flattenMk8sTritonNodePools(expectedNodePools)

	// Extract the interface slice from *schema.Set
	// Provider
	expectedFlattenItem := expectedFlatten[0].(map[string]interface{})

	// Private Network IDs
	expectedFlattenItem["private_network_ids"] = expectedFlattenItem["private_network_ids"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Triton Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sDigitalOceanNodePools(t *testing.T) {

	_, expectedNodePools, expectedFlatten := generateTestMk8sDigitalOceanNodePools()
	flattenedNodePools := flattenMk8sDigitalOceanNodePools(expectedNodePools)

	if diff := deep.Equal(expectedFlatten, flattenedNodePools); diff != nil {
		t.Errorf("Mk8s Digital Ocean Node Pools was not flattened correctly. Diff: %s", diff)
	}
}

// !SECTION

// SECTION AWS

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

func TestControlPlane_FlattenMk8sAwsDeployRoleChain(t *testing.T) {

	_, expectedDeployRoleChain, expectedFlatten := generateTestMk8sAwsDeployRoleChain()
	flattenedDeployRoleChain := flattenMk8sAwsDeployRoleChain(expectedDeployRoleChain)

	if diff := deep.Equal(expectedFlatten, flattenedDeployRoleChain); diff != nil {
		t.Errorf("Mk8s AWS Deploy Role Chain was not flattened correctly. Diff: %s", diff)
	}
}

// !SECTION

// SECTION Triton

func TestControlPlane_FlattenMk8sTritonConnection(t *testing.T) {

	_, expectedConnection, expectedFlatten := generateTestMk8sTritonConnection()
	flattenedConnection := flattenMk8sTritonConnection(expectedConnection)

	if diff := deep.Equal(expectedFlatten, flattenedConnection); diff != nil {
		t.Errorf("Mk8s Triton Connection was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sTritonLoadBalancerGateway(t *testing.T) {

	_, expectedLoadBalancer, expectedFlatten := generateTestMk8sTritonLoadBalancer("gateway")
	flattenedLoadBalancer := flattenMk8sTritonLoadBalancer(expectedLoadBalancer)

	if diff := deep.Equal(expectedFlatten, flattenedLoadBalancer); diff != nil {
		t.Errorf("Mk8s Triton LoadBalancer Gateway was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenMk8sTritonLoadBalancerManual(t *testing.T) {

	_, expectedLoadBalancer, expectedFlatten := generateTestMk8sTritonLoadBalancer("manual")
	flattenedLoadBalancer := flattenMk8sTritonLoadBalancer(expectedLoadBalancer)

	// Extract the interface slice from *schema.Set
	loadBalancer := expectedFlatten[0].(map[string]interface{})
	manual := loadBalancer["manual"].([]interface{})[0].(map[string]interface{})
	manual["private_network_ids"] = manual["private_network_ids"].(*schema.Set).List()

	if diff := deep.Equal(expectedFlatten, flattenedLoadBalancer); diff != nil {
		t.Errorf("Mk8s Triton LoadBalancer Manual was not flattened correctly. Diff: %s", diff)
	}
}

// !SECTION

// SECTION Common

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

// !SECTION

// SECTION Add Ons

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

	_, expectedAddOn, expectedFlatten := generateTestMk8sAwsAddOn(true)
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

// !SECTION
// !SECTION
// !SECTION

// SECTION Generate

// SECTION Build

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

func generateTestMk8sProvider(provider string, update string) *client.Mk8sProvider {

	output := client.Mk8sProvider{}

	switch provider {
	case "generic":
		generated, _, _ := generateTestMk8sGenericProvider()
		output.Generic = generated
	case "hetzner":
		generated, _, _ := generateTestMk8sHetznerProvider(update)
		output.Hetzner = generated
	case "aws":
		generated, _, _ := generateTestMk8sAwsProvider()
		output.Aws = generated
	case "linode":
		generated, _, _ := generateTestMk8sLinodeProvider()
		output.Linode = generated
	case "oblivus":
		generated, _, _ := generateTestMk8sOblivusProvider()
		output.Oblivus = generated
	case "lambdalabs":
		generated, _, _ := generateTestMk8sLambdalabsProvider()
		output.Lambdalabs = generated
	case "paperspace":
		generated, _, _ := generateTestMk8sPaperspaceProvider()
		output.Paperspace = generated
	case "ephemeral":
		generated, _, _ := generateTestMk8sEphemeralProvider()
		output.Ephemeral = generated
	case "triton":
		generated, _, _ := generateTestMk8sTritonProvider(update)
		output.Triton = generated
	}

	return &output
}

func generateTestMk8sAddOns(providerName string, update string) (*client.Mk8sSpecAddOns, *client.Mk8sSpecAddOns, []interface{}) {

	dashboard := true
	azureWorkloadIdentity, _, flattenedAzureWorkloadIdentity := generateTestMk8sAzureWorkloadIdentityAddOn()
	awsWorkloadIdentity := true
	localPathStorage := true
	metrics, _, flattenedMetrics := generateTestMk8sMetricsAddOn()
	logs, _, flattenedLogs := generateTestMk8sLogsAddOn()
	nvidia, _, flattenedNvidia := generateTestMk8sNvidiaAddOn()
	azureAcr, _, flattenedAzureAcr := generateTestMk8sAzureAcrAddOn()
	var sysbox *bool

	switch update {
	case "case1":
		dashboard = false
		awsWorkloadIdentity = false
		localPathStorage = false
		sysbox = GetBool(false)
	case "no-sysbox":
		sysbox = nil
	default:
		sysbox = GetBool(true)
	}

	var awsEfs *client.Mk8sAwsAddOnConfig
	var flattenedAwsEfs []interface{}

	var awsEcr *client.Mk8sAwsAddOnConfig
	var flattenedAwsEcr []interface{}

	var awsElb *client.Mk8sAwsAddOnConfig
	var flattenedAwsElb []interface{}

	if providerName == "aws" {
		awsEfs, _, flattenedAwsEfs = generateTestMk8sAwsAddOn(true)
		awsEcr, _, flattenedAwsEcr = generateTestMk8sAwsAddOn(true)
		awsElb, _, flattenedAwsElb = generateTestMk8sAwsAddOn(false)
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
	}

	if sysbox != nil {
		expectedAddOns.Sysbox = &client.Mk8sNonCustomizableAddonConfig{}
	}

	return addOns, &expectedAddOns, flattened
}

// SECTION Providers

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

func generateTestMk8sHetznerProvider(update string) (*client.Mk8sHetznerProvider, *client.Mk8sHetznerProvider, []interface{}) {

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
	floatingIpSelector := map[string]interface{}{
		"floating_ip_1": "123.45.67.89",
	}

	// Handle updates
	switch update {
	case "case1":
		floatingIpSelector["floating_ip_2"] = "98.76.54.32"
	}

	flattened := generateFlatTestMk8sHetznerProvider(region, hetznerLabels, flattenedNetworking, preInstallScript, tokenSecretLink, networkId, flattenedNodePools, flattenedDedicatedServerNodePools, image, sshKey, flattenedAutoscaler, floatingIpSelector)
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
		FloatingIpSelector:       &floatingIpSelector,
	}

	return hetzner, &expectedHetzner, flattened
}

func generateTestMk8sAwsProvider() (*client.Mk8sAwsProvider, *client.Mk8sAwsProvider, []interface{}) {

	region := "eu-central-1"
	awsTags := map[string]interface{}{
		"hello": "world",
	}
	skipCreateRoles := false
	networking, _, flattenedNetworking := generateTestMk8sNetworking()
	preInstallScript := "#! echo hello world"
	image, _, flattenedImage := generateTestMk8sAwsAmi("recommended")
	deployRoleArn := "arn:aws:iam::483676437512:role/cpln-mk8s-terraform-test-org"
	deployRoleChain, _, flattenedDeployRoleChain := generateTestMk8sAwsDeployRoleChain()
	vpcId := "vpc-03105bd4dc058d3a8"
	keyPair := "debug-terraform"
	diskEncryptionKeyArn := "arn:aws:kms:eu-central-1:989132402664:key/2e9f25ea-efb4-49bf-ae39-007be298726d"
	securityGroupIds := []string{"sg-031480aa7a1e6e38b"}
	extraNodePolicies := []string{"arn:aws:iam::aws:policy/IAMFullAccess"}
	nodePools, _, flattenedNodePools := generateTestMk8sAwsNodePools()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sAwsProvider(region, awsTags, skipCreateRoles, flattenedNetworking, preInstallScript, flattenedImage, deployRoleArn, flattenedDeployRoleChain, vpcId, keyPair, diskEncryptionKeyArn, securityGroupIds, extraNodePolicies, flattenedNodePools, flattenedAutoscaler)
	aws := buildMk8sAwsProvider(flattened)
	expectedAws := client.Mk8sAwsProvider{
		Region:               &region,
		AwsTags:              &awsTags,
		SkipCreateRoles:      &skipCreateRoles,
		Networking:           networking,
		PreInstallScript:     &preInstallScript,
		Image:                image,
		DeployRoleArn:        &deployRoleArn,
		DeployRoleChain:      deployRoleChain,
		VpcId:                &vpcId,
		KeyPair:              &keyPair,
		DiskEncryptionKeyArn: &diskEncryptionKeyArn,
		SecurityGroupIds:     &securityGroupIds,
		ExtraNodePolicies:    &extraNodePolicies,
		NodePools:            nodePools,
		Autoscaler:           autoscaler,
	}

	return aws, &expectedAws, flattened
}

func generateTestMk8sLinodeProvider() (*client.Mk8sLinodeProvider, *client.Mk8sLinodeProvider, []interface{}) {

	region := "fr-par"
	tokenSecretLink := "/org/terraform-test-org/secret/linode"
	image := "linode/ubuntu24.04"
	vpcId := "93666"
	firewallId := "168425"
	preInstallScript := "#! echo hello world"
	authorizedUsers := []string{"juliancpln"}
	authorizedKeys := []string{}
	nodePools, _, flattenedNodePools := generateTestMk8sLinodeNodePools()
	networking, _, flattenedNetworking := generateTestMk8sNetworking()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sLinodeProvider(region, tokenSecretLink, firewallId, flattenedNodePools, image, authorizedUsers, authorizedKeys, vpcId, preInstallScript, flattenedNetworking, flattenedAutoscaler)
	linode := buildMk8sLinodeProvider(flattened)
	expectedLinode := client.Mk8sLinodeProvider{
		Region:           &region,
		TokenSecretLink:  &tokenSecretLink,
		FirewallId:       &firewallId,
		NodePools:        nodePools,
		Image:            &image,
		AuthorizedUsers:  &authorizedUsers,
		AuthorizedKeys:   &authorizedKeys,
		VpcId:            &vpcId,
		PreInstallScript: &preInstallScript,
		Networking:       networking,
		Autoscaler:       autoscaler,
	}

	return linode, &expectedLinode, flattened
}

func generateTestMk8sOblivusProvider() (*client.Mk8sOblivusProvider, *client.Mk8sOblivusProvider, []interface{}) {

	datacenter := "OSL1"
	tokenSecretLink := "/org/terraform-test-org/secret/oblivus"
	preInstallScript := "#! echo hello world"
	sshKeys := []string{}
	nodePools, _, flattenedNodePools := generateTestMk8sOblivusNodePools()
	unmanagedNodePools, _, flattenedUnmanagedNodePools := generateTestMk8sGenericNodePools()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sOblivusProvider(datacenter, tokenSecretLink, flattenedNodePools, sshKeys, flattenedUnmanagedNodePools, flattenedAutoscaler, preInstallScript)
	oblivus := buildMk8sOblivusProvider(flattened)
	expectedOblivus := client.Mk8sOblivusProvider{
		Datacenter:         &datacenter,
		TokenSecretLink:    &tokenSecretLink,
		NodePools:          nodePools,
		SshKeys:            &sshKeys,
		UnmanagedNodePools: unmanagedNodePools,
		Autoscaler:         autoscaler,
		PreInstallScript:   &preInstallScript,
	}

	return oblivus, &expectedOblivus, flattened
}

func generateTestMk8sLambdalabsProvider() (*client.Mk8sLambdalabsProvider, *client.Mk8sLambdalabsProvider, []interface{}) {

	region := "europe-central-1"
	tokenSecretLink := "/org/terraform-test-org/secret/lambdalabs"
	nodePools, _, flattenedNodePools := generateTestMk8sLambdalabsNodePools()
	sshKey := "julian-test"
	unmanagedNodePools, _, flattenedUnmanagedNodePools := generateTestMk8sGenericNodePools()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()
	preInstallScript := "#! echo hello world"

	flattened := generateFlatTestMk8sLambdalabsProvider(region, tokenSecretLink, flattenedNodePools, sshKey, flattenedUnmanagedNodePools, flattenedAutoscaler, preInstallScript)
	lambdalabs := buildMk8sLambdalabsProvider(flattened)
	expectedLambdalabs := client.Mk8sLambdalabsProvider{
		Region:             &region,
		TokenSecretLink:    &tokenSecretLink,
		NodePools:          nodePools,
		SshKey:             &sshKey,
		UnmanagedNodePools: unmanagedNodePools,
		Autoscaler:         autoscaler,
		PreInstallScript:   &preInstallScript,
	}

	return lambdalabs, &expectedLambdalabs, flattened
}

func generateTestMk8sPaperspaceProvider() (*client.Mk8sPaperspaceProvider, *client.Mk8sPaperspaceProvider, []interface{}) {

	region := "CA1"
	tokenSecretLink := "/org/terraform-test-org/secret/paperspace"
	preInstallScript := "#! echo hello world"
	sharedDrivers := []string{"california"}
	userIds := []string{}
	networkId := "nla0jotp"

	nodePools, _, flattenedNodePools := generateTestMk8sPaperspaceNodePools()
	unmanagedNodePools, _, flattenedUnmanagedNodePools := generateTestMk8sGenericNodePools()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sPaperspaceProvider(region, tokenSecretLink, sharedDrivers, flattenedNodePools, flattenedAutoscaler, flattenedUnmanagedNodePools, preInstallScript, userIds, networkId)
	paperspace := buildMk8sPaperspaceProvider(flattened)
	expectedPaperspace := client.Mk8sPaperspaceProvider{
		Region:             &region,
		TokenSecretLink:    &tokenSecretLink,
		SharedDrives:       &sharedDrivers,
		NodePools:          nodePools,
		Autoscaler:         autoscaler,
		UnmanagedNodePools: unmanagedNodePools,
		PreInstallScript:   &preInstallScript,
		UserIds:            &userIds,
		NetworkId:          &networkId,
	}

	return paperspace, &expectedPaperspace, flattened
}

func generateTestMk8sEphemeralProvider() (*client.Mk8sEphemeralProvider, *client.Mk8sEphemeralProvider, []interface{}) {

	location := "aws-eu-central-1"
	nodePools, _, flattenedNodePools := generateTestMk8sEphemeralNodePools()

	flattened := generateFlatTestMk8sEphemeralProvider(location, flattenedNodePools)
	ephemeral := buildMk8sEphemeralProvider(flattened)
	expectedEphemeral := client.Mk8sEphemeralProvider{
		Location:  &location,
		NodePools: nodePools,
	}

	return ephemeral, &expectedEphemeral, flattened
}

func generateTestMk8sTritonProvider(update string) (*client.Mk8sTritonProvider, *client.Mk8sTritonProvider, []interface{}) {

	preInstallScript := "#! echo hello world"
	location := "aws-eu-central-1"
	privateNetworkId := "6704dae9-00f4-48b5-8bbf-1be538f20587"
	firewallEnabled := false
	imageId := "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2"
	sshKeys := []string{}

	connection, _, flattenedConnection := generateTestMk8sTritonConnection()
	networking, _, flattenedNetworking := generateTestMk8sNetworking()
	loadBalancer, _, flattenedLoadBalancer := generateTestMk8sTritonLoadBalancer(update)
	nodePools, _, flattenedNodePools := generateTestMk8sTritonNodePools()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sTritonProvider(flattenedConnection, flattenedNetworking, preInstallScript, location, flattenedLoadBalancer, privateNetworkId, firewallEnabled, flattenedNodePools, imageId, sshKeys, flattenedAutoscaler)
	triton := buildMk8sTritonProvider(flattened)
	expectedTriton := client.Mk8sTritonProvider{
		Connection:       connection,
		Networking:       networking,
		PreInstallScript: &preInstallScript,
		Location:         &location,
		LoadBalancer:     loadBalancer,
		PrivateNetworkId: &privateNetworkId,
		FirewallEnabled:  &firewallEnabled,
		NodePools:        nodePools,
		ImageId:          &imageId,
		SshKeys:          &sshKeys,
		Autoscaler:       autoscaler,
	}

	return triton, &expectedTriton, flattened
}

func generateTestMk8sDigitalOceanProvider() (*client.Mk8sDigitalOceanProvider, *client.Mk8sDigitalOceanProvider, []interface{}) {

	region := "ams3"
	preInstallScript := "#! echo hello world"
	tokenSecretLink := "/org/terraform-test-org/secret/digitalocean"
	vpcId := "6704dae9-00f4-48b5-8bbf-1be538f20587"
	image := "debian-11"
	digitalOceanTags := []string{}
	sshKeys := []string{}
	extraSshKeys := []string{}
	reservedIps := []string{}

	networking, _, flattenedNetworking := generateTestMk8sNetworking()
	nodePools, _, flattenedNodePools := generateTestMk8sDigitalOceanNodePools()
	autoscaler, _, flattenedAutoscaler := generateTestMk8sAutoscaler()

	flattened := generateFlatTestMk8sDigitalOceanProvider(region, digitalOceanTags, flattenedNetworking, preInstallScript, tokenSecretLink, vpcId, flattenedNodePools, image, sshKeys, extraSshKeys, flattenedAutoscaler, reservedIps)
	digitalOcean := buildMk8sDigitalOceanProvider(flattened)
	expectedDigitalOcean := client.Mk8sDigitalOceanProvider{
		Region:           &region,
		DigitalOceanTags: &digitalOceanTags,
		Networking:       networking,
		PreInstallScript: &preInstallScript,
		TokenSecretLink:  &tokenSecretLink,
		VpcId:            &vpcId,
		NodePools:        nodePools,
		Image:            &image,
		SshKeys:          &sshKeys,
		ExtraSshKeys:     &extraSshKeys,
		Autoscaler:       autoscaler,
		ReservedIps:      &reservedIps,
	}

	return digitalOcean, &expectedDigitalOcean, flattened
}

// !SECTION

// SECTION Node Pools

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
	subnetIds := []string{"subnet-0e564a042e2a45009"}
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

func generateTestMk8sLinodeNodePools() (*[]client.Mk8sLinodePool, *[]client.Mk8sLinodePool, []interface{}) {

	name := "my-linode-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()
	serverType := "g6-nanode-1"
	overrideImage := "linode/debian11"
	subnetId := "90623"
	minSize := 0
	maxSize := 0

	flattened := generateFlatTestMk8sLinodeNodePools(name, labels, flattenedTaints, serverType, overrideImage, subnetId, minSize, maxSize)
	nodePools := buildMk8sLinodeNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sLinodePool{
		ServerType:    &serverType,
		OverrideImage: &overrideImage,
		SubnetId:      &subnetId,
		MinSize:       &minSize,
		MaxSize:       &maxSize,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sLinodePool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sOblivusNodePools() (*[]client.Mk8sOblivusPool, *[]client.Mk8sOblivusPool, []interface{}) {

	name := "my-oblivus-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()
	minSize := 0
	maxSize := 0
	flavor := "INTEL_XEON_V3_x4"

	flattened := generateFlatTestMk8sOblivusNodePools(name, labels, flattenedTaints, minSize, maxSize, flavor)
	nodePools := buildMk8sOblivusNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sOblivusPool{
		MinSize: &minSize,
		MaxSize: &maxSize,
		Flavor:  &flavor,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sOblivusPool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sLambdalabsNodePools() (*[]client.Mk8sLambdalabsPool, *[]client.Mk8sLambdalabsPool, []interface{}) {

	name := "my-lambdalabs-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()
	minSize := 0
	maxSize := 0
	instanceType := "cpu_4x_general"

	flattened := generateFlatTestMk8sLambdalabsNodePools(name, labels, flattenedTaints, minSize, maxSize, instanceType)
	nodePools := buildMk8sLambdalabsNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sLambdalabsPool{
		MinSize:      &minSize,
		MaxSize:      &maxSize,
		InstanceType: &instanceType,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sLambdalabsPool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sPaperspaceNodePools() (*[]client.Mk8sPaperspacePool, *[]client.Mk8sPaperspacePool, []interface{}) {

	name := "my-paperspace-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()
	minSize := 0
	maxSize := 0
	publicIpType := "dynamic"
	bootDiskSize := 50
	machineType := "GPU+"

	flattened := generateFlatTestMk8sPaperspaceNodePools(name, labels, flattenedTaints, minSize, maxSize, publicIpType, bootDiskSize, machineType)
	nodePools := buildMk8sPaperspaceNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sPaperspacePool{
		MinSize:      &minSize,
		MaxSize:      &maxSize,
		PublicIpType: &publicIpType,
		BootDiskSize: &bootDiskSize,
		MachineType:  &machineType,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sPaperspacePool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sEphemeralNodePools() (*[]client.Mk8sEphemeralPool, *[]client.Mk8sEphemeralPool, []interface{}) {

	name := "my-node-pool"
	labels := map[string]interface{}{
		"hello": "world",
	}
	taints, _, flattenedTaints := generateTestMk8sTaints()
	count := 1
	arch := "arm64"
	flavor := "debian"
	cpu := "50m"
	memory := "128Mi"

	flattened := generateFlatTestMk8sEphemeralNodePools(name, labels, flattenedTaints, count, arch, flavor, cpu, memory)
	nodePools := buildMk8sEphemeralNodePools(flattened)
	expectedNodePools := []client.Mk8sEphemeralPool{
		{
			Name:   &name,
			Labels: &labels,
			Taints: taints,
			Count:  &count,
			Arch:   &arch,
			Flavor: &flavor,
			Cpu:    &cpu,
			Memory: &memory,
		},
	}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sTritonNodePools() (*[]client.Mk8sTritonPool, *[]client.Mk8sTritonPool, []interface{}) {

	name := "my-triton-node-pool"
	packageId := "da311341-b42b-45a8-9386-78ede625d0a4"
	overrideImageId := "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
	publicNetworkId := "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
	minSize := 0
	maxSize := 0

	privateNetworkIds := []string{"6704dae9-00f4-48b5-8bbf-1be538f20587"}

	labels := map[string]interface{}{
		"hello": "world",
	}

	tritonTags := map[string]interface{}{
		"drink": "water",
	}

	taints, _, flattenedTaints := generateTestMk8sTaints()

	flattened := generateFlatTestMk8sTritonNodePools(name, labels, flattenedTaints, packageId, overrideImageId, publicNetworkId, privateNetworkIds, tritonTags, minSize, maxSize)
	nodePools := buildMk8sTritonNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sTritonPool{
		PackageId:         &packageId,
		OverrideImageId:   &overrideImageId,
		PublicNetworkId:   &publicNetworkId,
		PrivateNetworkIds: &privateNetworkIds,
		TritonTags:        &tritonTags,
		MinSize:           &minSize,
		MaxSize:           &maxSize,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sTritonPool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

func generateTestMk8sDigitalOceanNodePools() (*[]client.Mk8sDigitalOceanPool, *[]client.Mk8sDigitalOceanPool, []interface{}) {

	name := "my-digital-ocean-node-pool"
	dropletSize := "da311341-b42b-45a8-9386-78ede625d0a4"
	overrideImage := "debian-11"
	minSize := 0
	maxSize := 0

	labels := map[string]interface{}{
		"hello": "world",
	}

	taints, _, flattenedTaints := generateTestMk8sTaints()

	flattened := generateFlatTestMk8sDigitalOceanNodePools(name, labels, flattenedTaints, dropletSize, overrideImage, minSize, maxSize)
	nodePools := buildMk8sDigitalOceanNodePools(flattened)

	// Define expected node pool
	expectedNodePool := client.Mk8sDigitalOceanPool{
		DropletSize:   &dropletSize,
		OverrideImage: &overrideImage,
		MinSize:       &minSize,
		MaxSize:       &maxSize,
	}

	expectedNodePool.Name = &name
	expectedNodePool.Labels = &labels
	expectedNodePool.Taints = taints

	// Define expected node pools
	expectedNodePools := []client.Mk8sDigitalOceanPool{expectedNodePool}

	return nodePools, &expectedNodePools, flattened
}

// !SECTION

// SECTION AWS

func generateTestMk8sAwsAmi(choice string) (*client.Mk8sAwsAmi, *client.Mk8sAwsAmi, []interface{}) {

	var recommended *string
	var exact *string

	if choice == "recommended" {
		recommended = GetString("amazon/al2023")
	} else if choice == "exact" {
		exact = GetString("ami-0c5ee33c81cf67a7f")
	}

	flattened := generateFlatTestMk8sAwsAmi(recommended, exact)
	ami := buildMk8sAwsAmi(flattened)
	expectedAmi := client.Mk8sAwsAmi{
		Recommended: recommended,
		Exact:       exact,
	}

	return ami, &expectedAmi, flattened
}

func generateTestMk8sAwsDeployRoleChain() (*[]client.Mk8sAwsAssumeRoleLink, *[]client.Mk8sAwsAssumeRoleLink, []interface{}) {

	roleArn := "arn:aws:iam::483676437512:role/mk8s-chain-1"
	externalId := "chain-1"
	sessionNamePrefix := "foo-"

	flattened := generateFlatTestMk8sAwsDeployRoleChain(roleArn, externalId, sessionNamePrefix)
	deployRoleChain := buildMk8sAwsDeployRoleChain(flattened)

	expectedDeployRoleChain := []client.Mk8sAwsAssumeRoleLink{
		{
			RoleArn:           &roleArn,
			ExternalId:        &externalId,
			SessionNamePrefix: &sessionNamePrefix,
		},
	}

	return deployRoleChain, &expectedDeployRoleChain, flattened
}

// !SECTION

// SECTION Triton

func generateTestMk8sTritonConnection() (*client.Mk8sTritonConnection, *client.Mk8sTritonConnection, []interface{}) {

	url := "https://us-central-1.api.mnx.io"
	account := "eric_controlplane.com"
	user := "julian_controlplane.com"
	privateKeySecretLink := "/org/terraform-test-org/secret/triton"

	flattened := generateFlatTestMk8sTritonConnection(url, account, user, privateKeySecretLink)
	connection := buildMk8sTritonConnection(flattened)
	expectedConnection := client.Mk8sTritonConnection{
		Url:                  &url,
		Account:              &account,
		User:                 &user,
		PrivateKeySecretLink: &privateKeySecretLink,
	}

	return connection, &expectedConnection, flattened
}

func generateTestMk8sTritonLoadBalancer(option string) (*client.Mk8sTritonLoadBalancer, *client.Mk8sTritonLoadBalancer, []interface{}) {
	var gateway *client.Mk8sTritonGateway
	var manual *client.Mk8sTritonManual
	var gatewayFlattened *[]interface{}
	var manualFlattened *[]interface{}

	switch option {
	case "gateway":
		_flattened := generateFlatTestMk8sTritonGateway()
		gatewayFlattened = &_flattened
		gateway = &client.Mk8sTritonGateway{}
	case "manual":
		packageId := "df26ba1d-1261-6fc1-b35c-f1b390bc06ff"
		imageId := "8605a524-0655-43b9-adf1-7d572fe797eb"
		publicNetworkId := "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
		privateNetworkIds := []string{"6704dae9-00f4-48b5-8bbf-1be538f20587"}
		count := 1
		cnsInternalDomain := "example.com"
		cnsPublicDomain := "example.com"

		metadata := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		tags := map[string]interface{}{
			"tag1": "value1",
			"tag2": "value2",
		}

		_flattened := generateFlatTestMk8sTritonManual(packageId, imageId, publicNetworkId, privateNetworkIds, metadata, tags, count, cnsInternalDomain, cnsPublicDomain)
		manualFlattened = &_flattened
		manual = &client.Mk8sTritonManual{
			PackageId:         &packageId,
			ImageId:           &imageId,
			PublicNetworkId:   &publicNetworkId,
			PrivateNetworkIds: &privateNetworkIds,
			Metadata:          &metadata,
			Tags:              &tags,
			Count:             &count,
			CnsInternalDomain: &cnsInternalDomain,
			CnsPublicDomain:   &cnsPublicDomain,
		}
	}

	flattened := generateFlatTestMk8sTritonLoadBalancer(gatewayFlattened, manualFlattened)
	loadBalancer := buildMk8sTritonLoadBalancer(flattened)
	expectedLoadBalancer := client.Mk8sTritonLoadBalancer{
		Gateway: gateway,
		Manual:  manual,
	}

	return loadBalancer, &expectedLoadBalancer, flattened
}

// !SECTION

// SECTION Common

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

// !SECTION

// SECTION Add Ons

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
	includeNamespaces := "^elastic"
	excludeNamespaces := "^elastic"
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
	includeNamespaces := "^elastic"
	excludeNamespaces := "^elastic"

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

func generateTestMk8sAwsAddOn(withRoleArn bool) (*client.Mk8sAwsAddOnConfig, *client.Mk8sAwsAddOnConfig, []interface{}) {

	var roleArn *string = nil

	if withRoleArn {
		roleArn = GetString("arn:aws:iam::123456789012:role/my-custom-role")
	}

	flattened := generateFlatTestMk8sAwsAddOn(roleArn)
	aws := buildMk8sAwsAddOn(flattened)
	expectedAws := client.Mk8sAwsAddOnConfig{
		RoleArn: roleArn,
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

// !SECTION
// !SECTION

// SECTION Flatten

func generateFlatTestMk8sFirewall(sourceCidr string, description string) []interface{} {

	spec := map[string]interface{}{
		"source_cidr": sourceCidr,
		"description": description,
	}

	return []interface{}{
		spec,
	}
}

// SECTION Providers

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

func generateFlatTestMk8sHetznerProvider(region string, hetznerLabels map[string]interface{}, networking []interface{}, preInstallScript string, tokenSecretLink string, networkId string, nodePools []interface{}, dedicatedServerNodePools []interface{}, image string, sshKey string, autoscaler []interface{}, floatingIpSelector map[string]interface{}) []interface{} {

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
		"floating_ip_selector":       floatingIpSelector,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAwsProvider(region string, awsTags map[string]interface{}, skipCreateRoles bool, networking []interface{}, preInstallScript string, image []interface{}, deployRoleArn string, deployRoleChain []interface{}, vpcId string, keyPair string, diskEncryptionKeyArn string, securityGroupIds []string, extraNodePolicies []string, nodePools []interface{}, autoscaler []interface{}) []interface{} {

	spec := map[string]interface{}{
		"region":                  region,
		"aws_tags":                awsTags,
		"skip_create_roles":       skipCreateRoles,
		"networking":              networking,
		"pre_install_script":      preInstallScript,
		"image":                   image,
		"deploy_role_arn":         deployRoleArn,
		"deploy_role_chain":       deployRoleChain,
		"vpc_id":                  vpcId,
		"key_pair":                keyPair,
		"disk_encryption_key_arn": diskEncryptionKeyArn,
		"security_group_ids":      ConvertStringSliceToSet(securityGroupIds),
		"extra_node_policies":     ConvertStringSliceToSet(extraNodePolicies),
		"node_pool":               nodePools,
		"autoscaler":              autoscaler,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sLinodeProvider(region string, tokenSecretLink string, firewallId string, nodePools []interface{}, image string, authorizedUsers []string, authorizedKeys []string, vpcId string, preInstallScript string, networking []interface{}, autoscaler []interface{}) []interface{} {

	spec := map[string]interface{}{
		"region":             region,
		"token_secret_link":  tokenSecretLink,
		"firewall_id":        firewallId,
		"node_pool":          nodePools,
		"image":              image,
		"authorized_users":   ConvertStringSliceToSet(authorizedUsers),
		"authorized_keys":    ConvertStringSliceToSet(authorizedKeys),
		"vpc_id":             vpcId,
		"pre_install_script": preInstallScript,
		"networking":         networking,
		"autoscaler":         autoscaler,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sOblivusProvider(datacenter string, tokenSecretLink string, nodePools []interface{}, sshKeys []string, unmanagedNodePools []interface{}, autoscaler []interface{}, preInstallScript string) []interface{} {

	spec := map[string]interface{}{
		"datacenter":          datacenter,
		"token_secret_link":   tokenSecretLink,
		"node_pool":           nodePools,
		"ssh_keys":            ConvertStringSliceToSet(sshKeys),
		"unmanaged_node_pool": unmanagedNodePools,
		"autoscaler":          autoscaler,
		"pre_install_script":  preInstallScript,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sLambdalabsProvider(region string, tokenSecretLink string, nodePools []interface{}, sshKey string, unmanagedNodePools []interface{}, autoscaler []interface{}, preInstallScript string) []interface{} {

	spec := map[string]interface{}{
		"region":              region,
		"token_secret_link":   tokenSecretLink,
		"node_pool":           nodePools,
		"ssh_key":             sshKey,
		"unmanaged_node_pool": unmanagedNodePools,
		"autoscaler":          autoscaler,
		"pre_install_script":  preInstallScript,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sPaperspaceProvider(region string, tokenSecretLink string, sharedDrivers []string, nodePools []interface{}, autoscaler []interface{}, unmanagedNodePools []interface{}, preInstallScript string, userIds []string, networkId string) []interface{} {

	spec := map[string]interface{}{
		"region":              region,
		"token_secret_link":   tokenSecretLink,
		"shared_drives":       ConvertStringSliceToSet(sharedDrivers),
		"node_pool":           nodePools,
		"autoscaler":          autoscaler,
		"unmanaged_node_pool": unmanagedNodePools,
		"pre_install_script":  preInstallScript,
		"user_ids":            ConvertStringSliceToSet(userIds),
		"network_id":          networkId,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sEphemeralProvider(location string, nodePools []interface{}) []interface{} {

	spec := map[string]interface{}{
		"location":  location,
		"node_pool": nodePools,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sTritonProvider(connection []interface{}, networking []interface{}, preInstallScript string, location string, loadBalancer []interface{}, privateNetworkId string, firewallEnabled bool, nodePools []interface{}, imageId string, sshKeys []string, autoscaler []interface{}) []interface{} {

	spec := map[string]interface{}{
		"connection":         connection,
		"networking":         networking,
		"pre_install_script": preInstallScript,
		"location":           location,
		"load_balancer":      loadBalancer,
		"private_network_id": privateNetworkId,
		"firewall_enabled":   firewallEnabled,
		"node_pool":          nodePools,
		"image_id":           imageId,
		"ssh_keys":           ConvertStringSliceToSet(sshKeys),
		"autoscaler":         autoscaler,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sDigitalOceanProvider(region string, digitalOceanTags []string, networking []interface{}, preInstallScript string, tokenSecretLink string, vpcId string, nodePools []interface{}, image string, sshKeys []string, extraSshKeys []string, autoscaler []interface{}, reservedIps []string) []interface{} {

	spec := map[string]interface{}{
		"region":             region,
		"digital_ocean_tags": ConvertStringSliceToSet(digitalOceanTags),
		"networking":         networking,
		"pre_install_script": preInstallScript,
		"token_secret_link":  tokenSecretLink,
		"vpc_id":             vpcId,
		"node_pool":          nodePools,
		"image":              image,
		"ssh_keys":           ConvertStringSliceToSet(sshKeys),
		"extra_ssh_keys":     ConvertStringSliceToSet(extraSshKeys),
		"autoscaler":         autoscaler,
		"reserved_ips":       ConvertStringSliceToSet(reservedIps),
	}

	return []interface{}{
		spec,
	}
}

// !SECTION

// SECTION Node Pools

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

func generateFlatTestMk8sLinodeNodePools(name string, labels map[string]interface{}, taints []interface{}, serverType string, overrideImage string, subnetId string, minSize int, maxSize int) []interface{} {

	spec := map[string]interface{}{
		"name":           name,
		"labels":         labels,
		"taint":          taints,
		"server_type":    serverType,
		"override_image": overrideImage,
		"subnet_id":      subnetId,
		"min_size":       minSize,
		"max_size":       maxSize,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sOblivusNodePools(name string, labels map[string]interface{}, taints []interface{}, minSize int, maxSize int, flavor string) []interface{} {

	spec := map[string]interface{}{
		"name":     name,
		"labels":   labels,
		"taint":    taints,
		"min_size": minSize,
		"max_size": maxSize,
		"flavor":   flavor,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sLambdalabsNodePools(name string, labels map[string]interface{}, taints []interface{}, minSize int, maxSize int, instanceType string) []interface{} {

	spec := map[string]interface{}{
		"name":          name,
		"labels":        labels,
		"taint":         taints,
		"min_size":      minSize,
		"max_size":      maxSize,
		"instance_type": instanceType,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sPaperspaceNodePools(name string, labels map[string]interface{}, taints []interface{}, minSize int, maxSize int, publicIpType string, bootDiskSize int, machineType string) []interface{} {

	spec := map[string]interface{}{
		"name":           name,
		"labels":         labels,
		"taint":          taints,
		"min_size":       minSize,
		"max_size":       maxSize,
		"public_ip_type": publicIpType,
		"boot_disk_size": bootDiskSize,
		"machine_type":   machineType,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sEphemeralNodePools(name string, labels map[string]interface{}, taints []interface{}, count int, arch string, flavor string, cpu string, memory string) []interface{} {

	spec := map[string]interface{}{
		"name":   name,
		"labels": labels,
		"taint":  taints,
		"count":  count,
		"arch":   arch,
		"flavor": flavor,
		"cpu":    cpu,
		"memory": memory,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sTritonNodePools(name string, labels map[string]interface{}, taints []interface{}, packageId string, overrideImageId string, publicNetworkId string, privateNetworkIds []string, tritonTags map[string]interface{}, minSize int, maxSize int) []interface{} {

	spec := map[string]interface{}{
		"name":                name,
		"labels":              labels,
		"taint":               taints,
		"package_id":          packageId,
		"override_image_id":   overrideImageId,
		"public_network_id":   publicNetworkId,
		"private_network_ids": ConvertStringSliceToSet(privateNetworkIds),
		"triton_tags":         tritonTags,
		"min_size":            minSize,
		"max_size":            maxSize,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sDigitalOceanNodePools(name string, labels map[string]interface{}, taints []interface{}, dropletSize string, overrideImage string, minSize int, maxSize int) []interface{} {

	spec := map[string]interface{}{
		"name":           name,
		"labels":         labels,
		"taint":          taints,
		"droplet_size":   dropletSize,
		"override_image": overrideImage,
		"min_size":       minSize,
		"max_size":       maxSize,
	}

	return []interface{}{
		spec,
	}
}

// !SECTION

// SECTION AWS

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

func generateFlatTestMk8sAwsDeployRoleChain(roleArn string, externalId string, sessionNamePrefix string) []interface{} {

	spec := map[string]interface{}{
		"role_arn":            roleArn,
		"external_id":         externalId,
		"session_name_prefix": sessionNamePrefix,
	}

	return []interface{}{
		spec,
	}
}

// !SECTION

// SECTION Triton

func generateFlatTestMk8sTritonConnection(url string, account string, user string, privateKeySecretLink string) []interface{} {

	spec := map[string]interface{}{
		"url":                     url,
		"account":                 account,
		"user":                    user,
		"private_key_secret_link": privateKeySecretLink,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sTritonLoadBalancer(gateway *[]interface{}, manual *[]interface{}) []interface{} {
	spec := map[string]interface{}{}

	if gateway != nil {
		spec["gateway"] = *gateway
	}

	if manual != nil {
		spec["manual"] = *manual
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sTritonGateway() []interface{} {
	spec := map[string]interface{}{
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sTritonManual(packageId string, imageId string, publicNetworkId string, privateNetworkIds []string, metadata map[string]interface{}, tags map[string]interface{}, count int, cnsInternalDomain string, cnsPublicDomain string) []interface{} {
	spec := map[string]interface{}{
		"package_id":          packageId,
		"image_id":            imageId,
		"public_network_id":   publicNetworkId,
		"private_network_ids": ConvertStringSliceToSet(privateNetworkIds),
		"metadata":            metadata,
		"tags":                tags,
		"count":               count,
		"cns_internal_domain": cnsInternalDomain,
		"cns_public_domain":   cnsPublicDomain,
	}

	return []interface{}{
		spec,
	}
}

// !SECTION

// SECTION Common

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

// !SECTION

// SECTION Add Ons

func generateFlatTestMk8sAddOns(dashboard bool, azureWorkloadIdentity []interface{}, awsWorkloadIdentity bool, localPathStorage bool, metrics []interface{}, logs []interface{}, nvidia []interface{}, awsEfs []interface{}, awsEcr []interface{}, awsElb []interface{}, azureAcr []interface{}, sysbox *bool) []interface{} {

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
	}

	if sysbox != nil {
		spec["sysbox"] = *sysbox
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAzureWorkloadIdentityAddOn(tenantId string) []interface{} {

	spec := map[string]interface{}{
		"tenant_id":             tenantId,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sMetricsAddOn(kubeState bool, coreDns bool, kubelet bool, apiServer bool, nodeExporter bool, cadvisor bool, scrapeAnnotated []interface{}) []interface{} {

	spec := map[string]interface{}{
		"kube_state":            kubeState,
		"core_dns":              coreDns,
		"kubelet":               kubelet,
		"api_server":            apiServer,
		"node_exporter":         nodeExporter,
		"cadvisor":              cadvisor,
		"scrape_annotated":      scrapeAnnotated,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sMetricsScrapeAnnotated(intervalSeconds int, includeNamespaces string, excludeNamespaces string, retainLabels string) []interface{} {

	spec := map[string]interface{}{
		"interval_seconds":      intervalSeconds,
		"include_namespaces":    includeNamespaces,
		"exclude_namespaces":    excludeNamespaces,
		"retain_labels":         retainLabels,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sLogsAddOn(auditEnabled bool, includeNamespaces string, excludeNamespaces string) []interface{} {

	spec := map[string]interface{}{
		"audit_enabled":         auditEnabled,
		"include_namespaces":    includeNamespaces,
		"exclude_namespaces":    excludeNamespaces,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sNvidiaAddOn(taintGpuNodes bool) []interface{} {

	spec := map[string]interface{}{
		"taint_gpu_nodes":       taintGpuNodes,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestMk8sAwsAddOn(roleArn *string) []interface{} {

	spec := map[string]interface{}{
		"placeholder_attribute": true,
	}

	if roleArn != nil {
		spec["role_arn"] = *roleArn
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

// !SECTION
// !SECTION
// !SECTION
