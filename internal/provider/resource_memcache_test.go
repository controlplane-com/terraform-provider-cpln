package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneMemcache_basic(t *testing.T) {

	var memcache client.Memcache
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	memcacheName := "memcache-" + randomName
	requiredOnlyMemcacheName := "required-only-memcache-" + randomName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "MEMCACHE") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneMemcacheCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneMemcache(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMemcacheExists("cpln_memcache.new", memcacheName, &memcache),
					testAccCheckControlPlaneMemcacheAttributes(&memcache, "default"),
					resource.TestCheckResourceAttr("cpln_memcache.new", "name", memcacheName),
					resource.TestCheckResourceAttr("cpln_memcache.new", "description", "Memcache description for "+memcacheName),
				),
			},
			{
				Config: testAccControlPlaneMemcache_Update(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMemcacheExists("cpln_memcache.new", memcacheName, &memcache),
					testAccCheckControlPlaneMemcacheAttributes(&memcache, "update"),
					resource.TestCheckResourceAttr("cpln_memcache.new", "name", memcacheName),
					resource.TestCheckResourceAttr("cpln_memcache.new", "description", "Updated Memcache description for "+memcacheName),
				),
			},
			{
				Config: testAccControlPlaneMemcache_RequiredOnly(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneMemcacheExists("cpln_memcache.new", requiredOnlyMemcacheName, &memcache),
					testAccCheckControlPlaneMemcacheAttributes(&memcache, "required-only"),
					resource.TestCheckResourceAttr("cpln_memcache.new", "name", requiredOnlyMemcacheName),
					resource.TestCheckResourceAttr("cpln_memcache.new", "description", "Memcache description for "+requiredOnlyMemcacheName),
				),
			},
		},
	})
}

func testAccControlPlaneMemcache(randomName string) string {
	return fmt.Sprintf(`
	variable "random-name" {
		type 	= string
		default = "%s"
	}

	resource "cpln_memcache" "new" {

		name 		= "memcache-${var.random-name}"
		description = "Memcache description for memcache-${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test 	= "true"
		}

		node_count = 1
		node_size  = 0.5
		version    = "1.5.22"

		options {
			eviction_disabled 	 = true
			idle_timeout_seconds = 600
			max_item_size 		 = 1024
			max_connections      = 1024
		}

		locations  = ["/org/terraform-test-org/location/aws-us-west-2"]
	}
	`, randomName)
}

func testAccControlPlaneMemcache_Update(randomName string) string {
	return fmt.Sprintf(`
	variable "random-name" {
		type 	= string
		default = "%s"
	}

	resource "cpln_memcache" "new" {

		name 		= "memcache-${var.random-name}"
		description = "Updated Memcache description for memcache-${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test 	= "false"
		}

		node_count = 2
		node_size  = 1
		version    = "1.6.17"

		options {
			eviction_disabled 	 = false
			idle_timeout_seconds = 650
			max_item_size 		 = 512
			max_connections      = 512
		}

		locations  = ["/org/terraform-test-org/location/aws-us-west-1"]
	}
	`, randomName)
}

func testAccControlPlaneMemcache_RequiredOnly(randomName string) string {
	return fmt.Sprintf(`
	variable "random-name" {
		type 	= string
		default = "%s"
	}

	resource "cpln_memcache" "new" {

		name 		= "required-only-memcache-${var.random-name}"
		description = "Memcache description for required-only-memcache-${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test 	= "true"
		}

		node_count = 1
		node_size  = 0.25

		options {
			eviction_disabled 	 = true
			idle_timeout_seconds = 600
			max_item_size 		 = 1024
			max_connections      = 1024
		}

		locations  = ["/org/terraform-test-org/location/aws-us-west-2"]
	}
	`, randomName)
}

func testAccCheckControlPlaneMemcacheExists(resourceName string, memcacheName string, memcache *client.Memcache) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneMemcacheExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != memcacheName {
			return fmt.Errorf("Memcache name does not match. %s != %s", rs.Primary.ID, memcacheName)
		}

		client := testAccProvider.Meta().(*client.Client)
		tempMemcache, _, err := client.GetMemcache(memcacheName)

		if err != nil {
			return err
		}

		if *tempMemcache.Name != memcacheName {
			return fmt.Errorf("Memcache name does not match")
		}

		*memcache = *tempMemcache

		return nil
	}
}

func testAccCheckControlPlaneMemcacheAttributes(memcache *client.Memcache, memcacheType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tags := *memcache.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("Tags - Memcache terraform_generated attribute does not match")
		}

		var expectedMemcacheOptions *client.MemcacheOptions

		switch memcacheType {
		case "update":
			expectedMemcacheOptions, _, _ = generateTestMemcacheOptions(false, 650, 512, 512)
		default:
			expectedMemcacheOptions, _, _ = generateTestMemcacheOptions(true, 600, 1024, 1024)
		}

		if diff := deep.Equal(memcache.Spec.Options, expectedMemcacheOptions); diff != nil {
			return fmt.Errorf("MemcacheOptions was not built correctly. Diff: %s", diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneMemcacheCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy For Memcache. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_memcache" {
			continue
		}

		memcache, _, _ := c.GetMemcache(rs.Primary.ID)
		if memcache != nil {
			return fmt.Errorf("Memcache still exists. Name: %s", *memcache.Name)
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build //
func TestControlPlane_BuildMemcacheOptions(t *testing.T) {
	memcacheOptions, expectedMemcacheOptions, _ := generateTestMemcacheOptions(true, 600, 1024, 1024)
	if diff := deep.Equal(memcacheOptions, expectedMemcacheOptions); diff != nil {
		t.Errorf("MemcacheOptions was not built correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
func generateTestMemcacheOptions(evictionDisabled bool, idleTimeoutSeconds int, maxItemSize int, maxConnections int) (*client.MemcacheOptions, *client.MemcacheOptions, []interface{}) {
	flattened := generateFlatTestMemcacheOptions(evictionDisabled, idleTimeoutSeconds, maxItemSize, maxConnections)
	memcacheOptions := buildMemcacheOptions(flattened)
	expectedMemcacheOptions := &client.MemcacheOptions{
		EvictionsDisabled:  &evictionDisabled,
		IdleTimeoutSeconds: &idleTimeoutSeconds,
		MaxItemSizeKiB:     &maxItemSize,
		MaxConnections:     &maxConnections,
	}

	return memcacheOptions, expectedMemcacheOptions, flattened
}

// Flatten //
func generateFlatTestMemcacheOptions(evictionDisabled bool, idleTimeoutSeconds int, maxItemSize int, maxConnections int) []interface{} {
	spec := map[string]interface{}{
		"eviction_disabled":    evictionDisabled,
		"idle_timeout_seconds": idleTimeoutSeconds,
		"max_item_size":        maxItemSize,
		"max_connections":      maxConnections,
	}

	return []interface{}{
		spec,
	}
}
