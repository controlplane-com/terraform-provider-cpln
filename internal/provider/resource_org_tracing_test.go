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

func TestAccControlPlaneOrgTracing_basic(t *testing.T) {

	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "ORG_TRACING") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneOrgTracingCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneOrgTracingLightstep(random, "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 50),
				),
			},
			{
				Config: testAccControlPlaneOrgTracingLightstep(random, "75"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 75),
				),
			},
		},
	})
}

func testAccControlPlaneOrgTracingLightstep(random, sampling string) string {

	TestLogger.Printf("Inside testAccControlPlaneOrgTracingLightstep")

	return fmt.Sprintf(`

	resource "cpln_secret" "opaque" {
		name = "opaque-random-tbd-%s"
		description = "description opaque-random-tbd" 
				
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "opaque"
		} 
		
		opaque {
			payload = "opaque_secret_payload"
			encoding = "plain"
		}
	}

	resource "cpln_org_tracing" "new" {

		lightstep_tracing {

			sampling = %s
			endpoint = "test.cpln.local:8080"

			// Opaque Secret Only
			credentials = cpln_secret.opaque.self_link
		}	
	}
	`, random, sampling)
}

func testAccCheckControlPlaneOrgTracingExists(resourceName string, sampling int) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		client := testAccProvider.Meta().(*client.Client)
		org, code, err := client.GetOrg()

		if code == 404 {
			return fmt.Errorf("Org not found")
		}

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		lightstepTracing, _ := generateLightstepTracing(sampling, *org.Spec.Tracing.Provider.Lightstep.Credentials)
		if diff := deep.Equal(lightstepTracing, org.Spec.Tracing); diff != nil {
			return fmt.Errorf("Org Tracing mismatch. Diff: %s", diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneOrgTracingCheckDestroy(s *terraform.State) error {

	// TestLogger.Printf("Inside testAccCheckControlPlaneOrgTracingCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneOrgTracingCheckDestroy: rs.Type: %s", rs.Type)

		if rs.Type != "cpln_org_tracing" {
			continue
		}

		orgName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneOrgTracingCheckDestroy: Org name: %s", orgName)

		org, _, _ := c.GetOrg()

		if org.Spec.Logging != nil {
			return fmt.Errorf("Org Spec Tracing still exists. Org Name: %s", *org.Name)
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build //
func TestControlPlane_BuildLightstepTracing(t *testing.T) {
	lightstepTracing, expectedLightstepTracing := generateLightstepTracing(50, "/org/terraform-test-org/secret/some-secret")
	if diff := deep.Equal(lightstepTracing, expectedLightstepTracing); diff != nil {
		t.Errorf("Lightstep tracing was not built correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
func generateLightstepTracing(sampling int, credentials string) (*client.Tracing, *client.Tracing) {
	endpoint := "test.cpln.local:8080"

	flattened := generateFlatTestLightstepTracing(sampling, endpoint, credentials)
	lightstepTracing := buildLightStepTracing(flattened)
	expectedLightstepTracing := &client.Tracing{
		Sampling: &sampling,
		Provider: &client.Provider{
			Lightstep: &client.LightstepTracing{
				Endpoint:    &endpoint,
				Credentials: &credentials,
			},
		},
	}

	return lightstepTracing, expectedLightstepTracing
}

func generateFlatTestLightstepTracing(sampling int, endpoint string, credentials string) []interface{} {
	spec := map[string]interface{}{
		"sampling":    sampling,
		"endpoint":    endpoint,
		"credentials": credentials,
	}

	return []interface{}{
		spec,
	}
}
