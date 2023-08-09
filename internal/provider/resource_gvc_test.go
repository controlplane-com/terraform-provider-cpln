package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO:
// Add to TestAcc: Add test for locations and tags

/*** Acc Tests ***/
func TestAccControlPlaneGvc_basic(t *testing.T) {

	var testGvc client.Gvc

	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rName := "gvc-" + random

	ep := resource.ExternalProvider{
		Source:            "time",
		VersionConstraint: "0.7.2",
	}

	eps := map[string]resource.ExternalProvider{
		"time": ep,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, "GVC") },
		Providers:         testAccProviders,
		ExternalProviders: eps,
		CheckDestroy:      testAccCheckControlPlaneGvcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneGvc(random, random, rName, "GVC created using terraform for acceptance tests", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneGvcExists("cpln_gvc.new", rName, &testGvc),
					testAccCheckControlPlaneGvcAttributes(50, &testGvc),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "GVC created using terraform for acceptance tests"),
				),
			},
			{
				Config: testAccControlPlaneGvc(random, random, rName, "GVC created using terraform for acceptance tests", "75"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneGvcExists("cpln_gvc.new", rName, &testGvc),
					testAccCheckControlPlaneGvcAttributes(75, &testGvc),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "GVC created using terraform for acceptance tests"),
				),
			},
			{
				Config: testAccControlPlaneGvc(random, random, rName+"renamed", "Renamed GVC created using terraform for acceptance tests", "75"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneGvcExists("cpln_gvc.new", rName+"renamed", &testGvc),
					testAccCheckControlPlaneGvcAttributes(75, &testGvc),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "Renamed GVC created using terraform for acceptance tests"),
				),
			},
		},
	})
}

func testAccControlPlaneGvc(random, random2, name, description, sampling string) string {

	return fmt.Sprintf(`

	resource "cpln_secret" "docker" {
		name = "docker-secret-%s"
		description = "docker secret" 
					
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
			secret_type = "docker"
		} 
			
		docker = "{\"auths\":{\"your-registry-server\":{\"username\":\"your-name\",\"password\":\"your-pword\",\"email\":\"your-email\",\"auth\":\"<Secret>\"}}}"
	}

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

	resource "time_sleep" "wait_30_seconds" {
		depends_on = [cpln_secret.docker]
	  	destroy_duration = "30s"
	}

	resource "cpln_gvc" "new" {

		depends_on = [time_sleep.wait_30_seconds]
		
		name        = "%s"	
		description = "%s"

		locations = ["aws-eu-central-1", "aws-us-west-2"]

		pull_secrets = [cpln_secret.docker.name]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}

		env = {
			env-name-01 = "env-value-01",
			env-name-02 = "env-value-02",
		}

		lightstep_tracing {

			sampling = %s
			endpoint = "test.cpln.local:8080"

			// Opaque Secret Only
			credentials = cpln_secret.opaque.self_link
		}	

		load_balancer {
			dedicated = true
		}

	  }`, random, random2, name, description, sampling)
}

func testAccCheckControlPlaneGvcExists(resourceName, gvcName string, gvc *client.Gvc) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != gvcName {
			return fmt.Errorf("GVC name does not match")
		}

		// Validate the data
		client := testAccProvider.Meta().(*client.Client)
		newGvc, code, err := client.GetGvc(gvcName)

		if code == 404 {
			return fmt.Errorf("GVC not found")
		}

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		if *newGvc.Name != gvcName {
			return fmt.Errorf("Gvc name does not match")
		}

		*gvc = *newGvc

		return nil
	}
}

func testAccCheckControlPlaneGvcAttributes(sampling int, gvc *client.Gvc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tags := *gvc.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("Tags - GVC terraform_generated attribute does not match")
		}

		if tags["acceptance_test"] != "true" {
			return fmt.Errorf("Tags - GVC acceptance_test attribute does not match")
		}

		lightstepTracing, _ := generateLightstepTracing(sampling, *gvc.Spec.Tracing.Provider.Lightstep.Credentials)
		if diff := deep.Equal(lightstepTracing, gvc.Spec.Tracing); diff != nil {
			return fmt.Errorf("GVC Tracing mismatch. Diff: %s", diff)
		}

		expectedLoadBalancer, _, _ := generateTestLoadBalancer()
		if diff := deep.Equal(expectedLoadBalancer, gvc.Spec.LoadBalancer); diff != nil {
			return fmt.Errorf("LoadBalancer attributes do not match. Diff: %s", diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneGvcDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneGvcDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_gvc" {
			continue
		}

		gvcName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneGvcDestroy: gvcName: %s", gvcName)

		gvc, _, _ := c.GetGvc(gvcName)
		if gvc != nil {
			return fmt.Errorf("GVC still exists. Name: %s", *gvc.Name)
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build //
func TestControlPlane_BuildLocations(t *testing.T) {

	org := "unit-test-org"

	locations := []interface{}{
		"us-east-2",
		"us-west-1",
	}

	stringFunc := schema.HashSchema(StringSchema())
	unitTestGvc := client.Gvc{}
	buildLocations(org, schema.NewSet(stringFunc, locations), &unitTestGvc)

	testLocation := []string{}

	for _, location := range locations {
		testLocation = append(testLocation, fmt.Sprintf("/org/%s/location/%s", org, location))
	}

	if diff := deep.Equal(unitTestGvc.Spec.StaticPlacement.LocationLinks, &testLocation); diff != nil {
		t.Errorf("LocationLinks did not built the location links correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildPullSecrets(t *testing.T) {

	org := "unit-test-org"

	pullSecrets := []interface{}{
		"gcr-secret",
		"docker-secret",
	}

	stringFunc := schema.HashSchema(StringSchema())
	unitTestGvc := client.Gvc{}
	buildPullSecrets(org, schema.NewSet(stringFunc, pullSecrets), &unitTestGvc)

	testPullSecrets := []string{}

	for _, pullSecret := range pullSecrets {
		testPullSecrets = append(testPullSecrets, fmt.Sprintf("/org/%s/secret/%s", org, pullSecret))
	}

	if diff := deep.Equal(unitTestGvc.Spec.PullSecretLinks, &testPullSecrets); diff != nil {
		t.Errorf("PullSecretLinks did not built the pull secret links correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildLoadBalancer(t *testing.T) {
	loadBalancer, expectedLoadBalancer, _ := generateTestLoadBalancer()
	if diff := deep.Equal(loadBalancer, expectedLoadBalancer); diff != nil {
		t.Errorf("LoadBalancer was not built correctly, Diff: %s", diff)
	}
}

// Flatten //
func TestControlPlane_FlattenLocations(t *testing.T) {

	org := "unit-test-org"

	locations := []string{
		"/org/unit-test-org/location/us-east-2",
		"/org/unit-test-org/location/us-west-1",
	}

	flatLocations := []string{
		"us-east-2",
		"us-west-1",
	}

	gvcSpec := client.GvcSpec{}
	gvcSpec.StaticPlacement = &client.StaticPlacement{
		LocationLinks: &locations,
	}

	flattenedLocations := flattenLocations(&gvcSpec, org)

	for i, location := range flatLocations {
		if flattenedLocations[i].(string) != location {
			t.Errorf("FlattenLocations did not flatten the locations correctly. Result: %s. Wanted: %s", flattenedLocations[i].(string), location)
		}
	}
}

func TestControlPlane_FlattenPullSecrets(t *testing.T) {

	org := "unit-test-org"

	pullSecrets := []string{
		"/org/unit-test-org/secret/gcp-secret",
		"/org/unit-test-org/secret/docker-secret",
	}

	flatPullSecrets := []string{
		"gcp-secret",
		"docker-secret",
	}

	gvcSpec := client.GvcSpec{
		PullSecretLinks: &pullSecrets,
	}

	flattenedPullSecrets := flattenPullSecrets(&gvcSpec, org)

	for i, pullSecret := range flatPullSecrets {
		if flattenedPullSecrets[i].(string) != pullSecret {
			t.Errorf("FlattenPullSecrets did not flatten the pull secrets correctly. Result: %s. Wanted: %s", flattenedPullSecrets[i].(string), pullSecret)
		}
	}
}

func TestControlPlane_FlattenLoadBalancer(t *testing.T) {
	_, expectedLoadBalancer, expectedFlatten := generateTestLoadBalancer()
	flattenLoadBalancer := flattenLoadBalancer(expectedLoadBalancer)

	if diff := deep.Equal(expectedFlatten, flattenLoadBalancer); diff != nil {
		t.Errorf("LoadBalancer was not flattened correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
func generateTestLoadBalancer() (*client.LoadBalancer, *client.LoadBalancer, []interface{}) {
	dedicated := true

	flatten := generateFlatTestLoadBalancer(dedicated)
	loadBalancer := buildLoadBalancer(flatten)
	expectedLoadBalancer := &client.LoadBalancer{
		Dedicated: &dedicated,
	}

	return loadBalancer, expectedLoadBalancer, flatten
}

// Flatten //
func generateFlatTestLoadBalancer(dedicated bool) []interface{} {
	spec := map[string]interface{}{
		"dedicated": dedicated,
	}

	return []interface{}{
		spec,
	}
}
