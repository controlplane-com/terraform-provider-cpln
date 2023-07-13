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

var org string = "terraform-test-org"

func TestAccControlPlaneSpicedb_basic(t *testing.T) {

	var spicedb client.Spicedb

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	spicedbName := "spice-db-" + randomName
	description := "SpiceDB description created using Terraform"
	locations := `["aws-eu-central-1"]`

	// Update variables
	descriptionUpdated := "SpiceDB description updated using Terraform"
	locationsUpdated := `["aws-eu-central-1", "aws-us-west-2"]`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "SPICE-DB") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneSpicedbCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneSpicedb(spicedbName, description, locations),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneSpicedbExists("cpln_spicedb.new", spicedbName, &spicedb),
					testAccCheckControlPlaneSpicedbAttributes(&spicedb, "new"),
					resource.TestCheckResourceAttr("cpln_spicedb.new", "name", spicedbName),
					resource.TestCheckResourceAttr("cpln_spicedb.new", "description", description),
				),
			},
			{
				Config: testAccControlPlaneSpicedb(spicedbName, descriptionUpdated, locationsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneSpicedbExists("cpln_spicedb.new", spicedbName, &spicedb),
					testAccCheckControlPlaneSpicedbAttributes(&spicedb, "update"),
					resource.TestCheckResourceAttr("cpln_spicedb.new", "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccControlPlaneSpicedb(spicedbName string, description string, locations string) string {

	return fmt.Sprintf(`

	resource "cpln_spicedb" "new" {
		name 		= "%s"
		description = "%s"

		tags = {
			terraform_generated = "true"
			acceptance_test     = "true"
		}

		version   = "1.14.1"
		locations = %s
	}

	`, spicedbName, description, locations)
}

func testAccCheckControlPlaneSpicedbExists(resourceName string, spicedbName string, spicedb *client.Spicedb) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneSpicedbExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != spicedbName {
			return fmt.Errorf("SpiceDB name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)
		_spicedb, _, err := client.GetSpicedb(spicedbName)

		if err != nil {
			return err
		}

		if *_spicedb.Name != spicedbName {
			return fmt.Errorf("SpiceDB name does not match")
		}

		*spicedb = *_spicedb

		return nil
	}
}

func testAccCheckControlPlaneSpicedbAttributes(spicedb *client.Spicedb, state string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tags := *spicedb.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("SpiceDB Tags - `terraform_generated` attribute does not match")
		}

		if tags["acceptance_test"] != "true" {
			return fmt.Errorf("SpiceDB Tags - `acceptance_test` attribute does not match")
		}

		expectedVersion, expectedLocations := generateTestSpicedbClusterProperties(state)

		// Check version
		if diff := deep.Equal(spicedb.Spec.Version, expectedVersion); diff != nil {
			return fmt.Errorf("SpiceDB Version does not match. Diff: %s", diff)
		}

		// Check Locations
		if diff := deep.Equal(spicedb.Spec.Locations, expectedLocations); diff != nil {
			return fmt.Errorf("SpiceDB Locations do not match. Diff: %s", diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneSpicedbCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy For SpiceDB. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_spicedb" {
			continue
		}

		name := rs.Primary.ID

		spicedb, _, _ := c.GetSpicedb(name)

		if spicedb != nil {
			return fmt.Errorf("SpiceDB still exists. Name: %s", *spicedb.Name)
		}
	}

	return nil
}

/*** Generate ***/
func generateTestSpicedbClusterProperties(state string) (*string, *[]string) {

	version := "1.14.1"
	locations := []string{generateSpicedbClusterLocationLink("aws-eu-central-1")}

	if state == "update" {
		locations = append(locations, generateSpicedbClusterLocationLink("aws-us-west-2"))
	}

	return &version, &locations
}

func generateSpicedbClusterLocationLink(locationName string) string {
	return fmt.Sprintf("/org/%s/location/%s", org, locationName)
}
