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
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 50, "lightstep", false),
				),
			},
			{
				Config: testAccControlPlaneOrgTracingLightstep(random, "75"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 75, "lightstep", false),
				),
			},
			{
				Config: testAccControlPlaneOrgTracingOtel("50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 50, "otel", false),
				),
			},
			{
				Config: testAccControlPlaneOrgTracingOtel("75"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 75, "otel", false),
				),
			},
			{
				Config: testAccControlPlaneOrgTracingControlPlane("50", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 50, "controlplane", false),
				),
			},
			{
				Config: testAccControlPlaneOrgTracingControlPlane("75", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneOrgTracingExists("cpln_org_tracing.new", 75, "controlplane", true),
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

func testAccControlPlaneOrgTracingOtel(sampling string) string {

	TestLogger.Printf("Inside testAccControlPlaneOrgTracingOtel")

	return fmt.Sprintf(`

	resource "cpln_org_tracing" "new" {

		otel_tracing {

			sampling = %s
			endpoint = "test.cpln.local:8080"
		}	
	}
	`, sampling)
}

func testAccControlPlaneOrgTracingControlPlane(sampling string, withCustomTags bool) string {

	TestLogger.Printf("Inside testAccControlPlaneOrgTracingControlPlane")

	customTags := ""

	if withCustomTags {
		customTags = `custom_tags = {
			hello = "world"
		}`
	}

	return fmt.Sprintf(`

	resource "cpln_org_tracing" "new" {

		controlplane_tracing {

			sampling = %s
			%s
		}
	}
	`, sampling, customTags)
}

func testAccCheckControlPlaneOrgTracingExists(resourceName string, sampling float64, tracingType string, withCustomTags bool) resource.TestCheckFunc {
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

		switch tracingType {
		case "lightstep":
			lightstepTracing, _ := generateLightstepTracing(sampling, *org.Spec.Tracing.Provider.Lightstep.Credentials)
			if diff := deep.Equal(lightstepTracing, org.Spec.Tracing); diff != nil {
				return fmt.Errorf("Org Tracing mismatch. Diff: %s", diff)
			}
		case "otel":
			otelTracing, _ := generateOtelTracing(sampling, "test.cpln.local:8080")
			if diff := deep.Equal(otelTracing, org.Spec.Tracing); diff != nil {
				return fmt.Errorf("Org Tracing mismatch. Diff: %s", diff)
			}
		case "controlplane":
			controlplaneTracing, _ := generateControlPlaneTracing(sampling, withCustomTags)
			if diff := deep.Equal(controlplaneTracing, org.Spec.Tracing); diff != nil {
				return fmt.Errorf("Org Tracing mismatch. Diff: %s", diff)
			}
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

func TestControlPlane_BuildOtelTracing(t *testing.T) {
	otelTracing, expectedOtelTracing := generateOtelTracing(50, "test.cpln.local:8080")
	if diff := deep.Equal(otelTracing, expectedOtelTracing); diff != nil {
		t.Errorf("Otel tracing was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildControlPlaneTracing_WithoutCustomTags(t *testing.T) {
	controlPlaneTracing, expectedControlPlaneTracing := generateControlPlaneTracing(50, false)
	if diff := deep.Equal(controlPlaneTracing, expectedControlPlaneTracing); diff != nil {
		t.Errorf("Control Plane tracing was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildControlPlaneTracing_WithCustomTags(t *testing.T) {
	controlPlaneTracing, expectedControlPlaneTracing := generateControlPlaneTracing(50, true)
	if diff := deep.Equal(controlPlaneTracing, expectedControlPlaneTracing); diff != nil {
		t.Errorf("Control Plane tracing was not built correctly. Diff: %s", diff)
	}
}

/*** Generate ***/

func generateLightstepTracing(sampling float64, credentials string) (*client.Tracing, *client.Tracing) {
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

func generateOtelTracing(sampling float64, endpoint string) (*client.Tracing, *client.Tracing) {
	flattened := generateFlatTestOtelTracing(sampling, endpoint)
	otelTracing := buildOtelTracing(flattened)
	expectedOtelTracing := &client.Tracing{
		Sampling: &sampling,
		Provider: &client.Provider{
			Otel: &client.OtelTelemetry{
				Endpoint: &endpoint,
			},
		},
	}

	return otelTracing, expectedOtelTracing
}

func generateControlPlaneTracing(sampling float64, withCustomTags bool) (*client.Tracing, *client.Tracing) {
	var customTags *map[string]interface{}

	if withCustomTags {
		customTags = &map[string]interface{}{
			"hello": "world",
		}
	}

	flattened := generateFlatTestControlPlaneTracing(sampling, customTags)
	controlPlaneTracing := buildControlPlaneTracing(flattened)
	expectedControlPlaneTracing := &client.Tracing{
		Sampling:   &sampling,
		CustomTags: buildCustomTags(customTags),
		Provider: &client.Provider{
			ControlPlane: &client.ControlPlaneTracing{},
		},
	}

	return controlPlaneTracing, expectedControlPlaneTracing
}

// Flatten //

func generateFlatTestLightstepTracing(sampling float64, endpoint string, credentials string) []interface{} {
	spec := map[string]interface{}{
		"sampling":    sampling,
		"endpoint":    endpoint,
		"credentials": credentials,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestOtelTracing(sampling float64, endpoint string) []interface{} {
	spec := map[string]interface{}{
		"sampling": sampling,
		"endpoint": endpoint,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestControlPlaneTracing(sampling float64, customTags *map[string]interface{}) []interface{} {
	spec := map[string]interface{}{
		"sampling": sampling,
	}

	if customTags != nil {
		spec["custom_tags"] = *customTags
	}

	return []interface{}{
		spec,
	}
}
