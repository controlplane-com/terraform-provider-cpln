package cpln

import (
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneIpSet_basic(t *testing.T) {

	var ipSet client.IpSet

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := "ipset-" + randomName
	description := "IpSet description created using Terraform"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "IPSET") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneIpSetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneIpSet(name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneIpSetExists("cpln_ipset.new", name, &ipSet),
					testAccCheckControlPlaneIpSetAttributes(&ipSet, randomName, ""),
					resource.TestCheckResourceAttr("cpln_ipset.new", "name", name),
					resource.TestCheckResourceAttr("cpln_ipset.new", "description", description),
				),
			},
			{
				Config: testAccControlPlaneIpSetWithLinkOnly(randomName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneIpSetExists("cpln_ipset.new", name, &ipSet),
					testAccCheckControlPlaneIpSetAttributes(&ipSet, randomName, "link_only"),
					resource.TestCheckResourceAttr("cpln_ipset.new", "name", name),
					resource.TestCheckResourceAttr("cpln_ipset.new", "description", description),
				),
			},
			{
				Config: testAccControlPlaneIpSetWithAllArguments(randomName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneIpSetExists("cpln_ipset.new", name, &ipSet),
					testAccCheckControlPlaneIpSetAttributes(&ipSet, randomName, "all_attributes"),
					resource.TestCheckResourceAttr("cpln_ipset.new", "name", name),
					resource.TestCheckResourceAttr("cpln_ipset.new", "description", description),
				),
			},
		},
	})
}

func testAccCheckControlPlaneIpSetExists(resourceName string, name string, ipSet *client.IpSet) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneIpSetExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != name {
			return fmt.Errorf("IpSet name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)
		_ipSet, _, err := client.GetIpSet(name)

		if err != nil {
			return err
		}

		if *_ipSet.Name != name {
			return fmt.Errorf("IpSet name does not match")
		}

		*ipSet = *_ipSet

		return nil
	}
}

func testAccCheckControlPlaneIpSetAttributes(ipSet *client.IpSet, randomName string, update string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tags := *ipSet.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("IpSet Tags - `terraform_generated` attribute does not match")
		}

		if tags["acceptance_test"] != "true" {
			return fmt.Errorf("IpSet Tags - `acceptance_test` attribute does not match")
		}

		expectedLink := fmt.Sprintf("/org/terraform-test-org/gvc/ipset-gvc-%s/workload/httpbin-example-%s", randomName, randomName)

		switch update {
		case "link_only":
			if diff := deep.Equal(ipSet.Spec.Link, &expectedLink); diff != nil {
				return fmt.Errorf("IpSet Link does not match. Diff: %s", diff)
			}

		case "all_attributes":
			if diff := deep.Equal(ipSet.Spec.Link, &expectedLink); diff != nil {
				return fmt.Errorf("IpSet Link does not match. Diff: %s", diff)
			}

			// Locations
			expectedLocations, _, _ := generateTestIpSetLocations()

			if diff := deep.Equal(ipSet.Spec.Locations, expectedLocations); diff != nil {
				return fmt.Errorf("IpSet Locations does not match. Diff: %s", diff)
			}
		}

		return nil
	}
}

func testAccCheckControlPlaneIpSetCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_ipset" {
			continue
		}

		ipSet, _, _ := c.GetIpSet(rs.Primary.ID)

		if ipSet != nil {
			return fmt.Errorf("IpSet still exists. Name: %s", *ipSet.Name)
		}
	}

	return nil
}

// SECTION Acceptance Tests

// SECTION Create

func testAccControlPlaneIpSet(name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_ipset" "new" {
		
		name        = "%s"
		description = "%s"

		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}
	}
	`, name, description)
}

// !SECTION

// SECTION Update

func testAccControlPlaneIpSetWithLinkOnly(randomName string, name string, description string) string {

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "new" {
		name        = "ipset-gvc-${var.random-name}"	
		description = "ipset-gvc-${var.random-name}"

		locations = ["aws-eu-central-1"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}

	resource "cpln_workload" "new" {

		gvc = cpln_gvc.new.name

		name = "httpbin-example-${var.random-name}"
		description = "httpbin-example-${var.random-name}"
		type = "serverless"
		support_dynamic_tags = false

		container {

			name = "httpbin"
			image = "kennethreitz/httpbin"
			cpu = "50m"
			memory = "128Mi"

			ports {
				number = 80
				protocol = "http"
			}
		}

		options {

			timeout_seconds = 5
			capacity_ai = true
			debug = false
			suspend = true

			autoscaling {
				metric = "concurrency"
				target = 100
				min_scale = 1
				max_scale = 1
				scale_to_zero_delay = 300
				max_concurrency = 1000
			}
		}

		firewall_spec {
		
			external {
				inbound_allow_cidr = ["0.0.0.0/0"]
				outbound_allow_cidr = ["0.0.0.0/0"]
			}

			internal {
				inbound_allow_type = "none"
			}
		}
	}

	resource "cpln_ipset" "new" {
		
		name        = "%s"
		description = "%s"

		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}

		link = cpln_workload.new.self_link
	}
	`, randomName, name, description)
}

func testAccControlPlaneIpSetWithAllArguments(randomName string, name string, description string) string {

	return fmt.Sprintf(`

	variable "random-name" {
		type = string
		default = "%s"
	}

	data "cpln_location" "aws-eu-central-1" {
		name = "aws-eu-central-1"
	}

	resource "cpln_gvc" "new" {
		name        = "ipset-gvc-${var.random-name}"	
		description = "ipset-gvc-${var.random-name}"

		locations = ["aws-eu-central-1"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}

	resource "cpln_workload" "new" {

		gvc = cpln_gvc.new.name

		name = "httpbin-example-${var.random-name}"
		description = "httpbin-example-${var.random-name}"
		type = "serverless"
		support_dynamic_tags = false

		container {

			name = "httpbin"
			image = "kennethreitz/httpbin"
			cpu = "50m"
			memory = "128Mi"

			ports {
				number = 80
				protocol = "http"
			}
		}

		options {

			timeout_seconds = 5
			capacity_ai = true
			debug = false
			suspend = true

			autoscaling {
				metric = "concurrency"
				target = 100
				min_scale = 1
				max_scale = 1
				scale_to_zero_delay = 300
				max_concurrency = 1000
			}
		}

		firewall_spec {
		
			external {
				inbound_allow_cidr = ["0.0.0.0/0"]
				outbound_allow_cidr = ["0.0.0.0/0"]
			}

			internal {
				inbound_allow_type = "none"
			}
		}
	}

	resource "cpln_ipset" "new" {
		
		name        = "%s"
		description = "%s"

		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}

		link = cpln_workload.new.self_link
		
		location {
			name             = data.cpln_location.aws-eu-central-1.self_link
			retention_policy = "keep"
		}
	}
	`, randomName, name, description)
}

// !SECTION
// !SECTION

// SECTION Unit Tests

// SECTION Build

func TestControlPlane_BuildIpSetLocations(t *testing.T) {

	locations, expectedLocations, _ := generateTestIpSetLocations()

	if diff := deep.Equal(locations, expectedLocations); diff != nil {
		t.Errorf("IpSet Locations was not built correctly, Diff: %s", diff)
	}
}

// !SECTION

// SECTION Flatten

func TestControlPlane_FlattenIpSetLocations(t *testing.T) {

	_, expectedLocations, expectedFlatten := generateTestIpSetLocations()
	flattenedLocations := flattenIpSetLocations(expectedLocations)

	if diff := deep.Equal(expectedFlatten, flattenedLocations); diff != nil {
		t.Errorf("IpSet Locations was not flattened correctly. Diff: %s", diff)
	}
}

// !SECTION
// !SECTION

// SECTION Generate

// SECTION Build

func generateTestIpSetLocations() (*[]client.IpSetLocation, *[]client.IpSetLocation, []interface{}) {

	name := "/org/terraform-test-org/location/aws-eu-central-1"
	retentionPolicy := "keep"

	flattened := generateFlatTestIpSetLocations(name, retentionPolicy)
	locations := buildIpSetLocations(flattened)
	expectedLocations := &[]client.IpSetLocation{
		{
			Name:            &name,
			RetentionPolicy: &retentionPolicy,
		},
	}

	return locations, expectedLocations, flattened
}

// !SECTION

// SECTION Flatten

func generateFlatTestIpSetLocations(name string, retentionPolicy string) []interface{} {

	spec := map[string]interface{}{
		"name":             name,
		"retention_policy": retentionPolicy,
	}

	return []interface{}{
		spec,
	}
}

// !SECTION
// !SECTION
