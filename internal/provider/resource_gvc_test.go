package cpln

import (
	"encoding/json"
	"fmt"
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/go-test/deep"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO: Add to TestAcc: Add test for locations and tags

const gvcEnvoyJson = `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"10s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`
const gvcEnvoyJsonUpdated = `{"clusters":[{"name":"provider_gcp","type":"STRICT_DNS","connect_timeout":"15s","dns_lookup_family":"V4_ONLY","lb_policy":"ROUND_ROBIN","load_assignment":{"cluster_name":"provider_gcp","endpoints":[{"lb_endpoints":[{"endpoint":{"address":{"socket_address":{"address":"www.googleapis.com","port_value":443}}}}]}]},"transport_socket":{"name":"envoy.transport_sockets.tls","typed_config":{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext","sni":"www.googleapis.com"}}}],"http":[{"name":"envoy.filters.http.grpc_web","typed_config":{"@type":"type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb"}}],"volumes":[{"path":"/etc/config","recoveryPolicy":"retain","uri":"scratch://config"}]}`

/*** Acc Tests ***/
func TestAccControlPlaneGvc_basic(t *testing.T) {

	var testGvc client.Gvc

	org := "terraform-test-org"
	random := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rName := "gvc-" + random

	ep := resource.ExternalProvider{
		Source:            "time",
		VersionConstraint: "0.9.2",
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
				Config: testAccControlPlaneGvc(random, random, rName, "GVC created using terraform for acceptance tests", "55.55", gvcEnvoyJson, 1, "my-ipset", "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneGvcExists("cpln_gvc.new", rName, &testGvc),
					testAccCheckControlPlaneGvcAttributes(55.55, gvcEnvoyJson, 1, &testGvc, org, "name"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "GVC created using terraform for acceptance tests"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "endpoint_naming_format", "default"),
				),
			},
			{
				Config: testAccControlPlaneGvc(random, random, rName, "GVC created using terraform for acceptance tests", "75", gvcEnvoyJsonUpdated, 2, "my-ipset", "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneGvcExists("cpln_gvc.new", rName, &testGvc),
					testAccCheckControlPlaneGvcAttributes(75, gvcEnvoyJsonUpdated, 2, &testGvc, org, "name"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "GVC created using terraform for acceptance tests"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "endpoint_naming_format", "default"),
				),
			},
			{
				Config: testAccControlPlaneGvc(random, random, rName+"renamed", "Renamed GVC created using terraform for acceptance tests", "75", gvcEnvoyJsonUpdated, 2, "my-ipset", "org"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneGvcExists("cpln_gvc.new", rName+"renamed", &testGvc),
					testAccCheckControlPlaneGvcAttributes(75, gvcEnvoyJsonUpdated, 2, &testGvc, org, "name"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "Renamed GVC created using terraform for acceptance tests"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "endpoint_naming_format", "org"),
				),
			},

			// GVC With Load Balancer - IP Set Complete Link
			{
				Config: testAccControlPlaneGvc(random, random, rName+"renamed", "Renamed GVC created using terraform for acceptance tests", "75", gvcEnvoyJsonUpdated, 2, fmt.Sprintf("/org/%s/ipset/my-ipset", org), "org"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneGvcExists("cpln_gvc.new", rName+"renamed", &testGvc),
					testAccCheckControlPlaneGvcAttributes(75, gvcEnvoyJsonUpdated, 2, &testGvc, org, "complete-link"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "description", "Renamed GVC created using terraform for acceptance tests"),
					resource.TestCheckResourceAttr("cpln_gvc.new", "endpoint_naming_format", "org"),
				),
			},
		},
	})
}

func testAccControlPlaneGvc(random, random2, name, description, sampling string, envoy string, trustedProxies int, ipset string, endpointNamingFormat string) string {

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

		endpoint_naming_format = "%s"
		locations              = ["aws-eu-central-1", "aws-us-west-2"]
		pull_secrets           = [cpln_secret.docker.name]

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
			trusted_proxies = %d
			ipset = "%s"

			redirect {
				class {
					status_5xx = "https://example.com/error/5xx"
					status_401 = "https://your-oauth-server/oauth2/authorize?return_to=%%REQ(:path)%%&client_id=your-client-id"
				}
			}
		}

		sidecar {
			envoy = jsonencode(%s)
		}

	  }`, random, random2, name, description, endpointNamingFormat, sampling, trustedProxies, ipset, envoy)
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
			return err
		}

		if *newGvc.Name != gvcName {
			return fmt.Errorf("Gvc name does not match")
		}

		*gvc = *newGvc

		return nil
	}
}

func testAccCheckControlPlaneGvcAttributes(sampling float64, envoy string, trustedProxies int, gvc *client.Gvc, org string, ipSetClass string) resource.TestCheckFunc {
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

		expectedLoadBalancer, _, _ := generateTestLoadBalancer(trustedProxies, org, ipSetClass)
		if diff := deep.Equal(expectedLoadBalancer, gvc.Spec.LoadBalancer); diff != nil {
			return fmt.Errorf("Load Balancer attributes do not match. Diff: %s", diff)
		}

		expectedSidecar, _, _ := generateTestGvcSidecar(envoy)
		if diff := deep.Equal(expectedSidecar, gvc.Spec.Sidecar); diff != nil {
			return fmt.Errorf("Sidecar attributes do not match. Diff: %s", diff)
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
	unitTestGvc.Spec = &client.GvcSpec{}
	buildLocations(org, schema.NewSet(stringFunc, locations), unitTestGvc.Spec)

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
	unitTestGvc.Spec = &client.GvcSpec{}
	buildPullSecrets(org, schema.NewSet(stringFunc, pullSecrets), unitTestGvc.Spec)

	testPullSecrets := []string{}

	for _, pullSecret := range pullSecrets {
		testPullSecrets = append(testPullSecrets, fmt.Sprintf("/org/%s/secret/%s", org, pullSecret))
	}

	if diff := deep.Equal(unitTestGvc.Spec.PullSecretLinks, &testPullSecrets); diff != nil {
		t.Errorf("PullSecretLinks did not built the pull secret links correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildLoadBalancer_IpSetNameClass(t *testing.T) {
	loadBalancer, expectedLoadBalancer, _ := generateTestLoadBalancer(1, "terraform-test-org", "name")
	if diff := deep.Equal(loadBalancer, expectedLoadBalancer); diff != nil {
		t.Errorf("Load Balancer - IP Set Name Only Class was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildLoadBalancer_IpSetCompleteLinkClass(t *testing.T) {
	loadBalancer, expectedLoadBalancer, _ := generateTestLoadBalancer(1, "terraform-test-org", "complete-link")
	if diff := deep.Equal(loadBalancer, expectedLoadBalancer); diff != nil {
		t.Errorf("Load Balancer - IP Set Complete Link was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildLoadBalancer_IpSetShortLinkClass(t *testing.T) {
	loadBalancer, expectedLoadBalancer, _ := generateTestLoadBalancer(1, "terraform-test-org", "short-link")
	if diff := deep.Equal(loadBalancer, expectedLoadBalancer); diff != nil {
		t.Errorf("Load Balancer - IP Set Short Link was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildGvcSidecar(t *testing.T) {
	sidecar, expectedSidecar, _ := generateTestGvcSidecar(gvcEnvoyJson)
	if diff := deep.Equal(sidecar, expectedSidecar); diff != nil {
		t.Errorf("GVC Sidecar was not built correctly, Diff: %s", diff)
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
	_, expectedLoadBalancer, expectedFlatten := generateTestLoadBalancer(1, "terraform-test-org", "name")
	flattenLoadBalancer := flattenLoadBalancer(expectedLoadBalancer, "name", "terraform-test-org")

	if diff := deep.Equal(expectedFlatten, flattenLoadBalancer); diff != nil {
		t.Errorf("LoadBalancer was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenGvcSidecar(t *testing.T) {
	_, expectedSidecar, expectedFlatten := generateTestGvcSidecar(gvcEnvoyJson)
	flattenSidecar := flattenGvcSidecar(expectedSidecar)

	if diff := deep.Equal(expectedFlatten, flattenSidecar); diff != nil {
		t.Errorf("Sidecar was not flattened correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
func generateTestLoadBalancer(trustedProxies int, org string, class string) (*client.LoadBalancer, *client.LoadBalancer, []interface{}) {
	dedicated := true
	ipsetName := "my-ipset"
	_, expectedRedirect, redirectFlatten := generateTestRedirect()

	var ipset string
	var expectedIpSet string

	switch class {
	case "complete-link":
		ipset = fmt.Sprintf("/org/%s/ipset/%s", org, ipsetName)
		expectedIpSet = ipset
	case "short-link":
		ipset = fmt.Sprintf("//ipset/%s", ipsetName)
		expectedIpSet = ipset
	default:
		ipset = ipsetName
		expectedIpSet = fmt.Sprintf("/org/%s/ipset/%s", org, ipsetName)
	}

	flatten := generateFlatTestLoadBalancer(dedicated, trustedProxies, redirectFlatten, ipset)
	loadBalancer := buildLoadBalancer(flatten, org)
	expectedLoadBalancer := &client.LoadBalancer{
		Dedicated:      &dedicated,
		TrustedProxies: &trustedProxies,
		Redirect:       expectedRedirect,
		IpSet:          &expectedIpSet,
	}

	return loadBalancer, expectedLoadBalancer, flatten
}

func generateTestRedirect() (*client.Redirect, *client.Redirect, []interface{}) {
	_, expectedClass, classFlatten := generateTestRedirectClass()

	flatten := generateFlatTestRedirect(classFlatten)
	redirect := buildRedirect(flatten)
	expectedRedirect := &client.Redirect{
		Class: expectedClass,
	}

	return redirect, expectedRedirect, flatten
}

func generateTestRedirectClass() (*client.RedirectClass, *client.RedirectClass, []interface{}) {
	status5XX := "https://example.com/error/5xx"
	status401 := "https://your-oauth-server/oauth2/authorize?return_to=%REQ(:path)%&client_id=your-client-id"

	flatten := generateFlatTestRedirectClass(status5XX, status401)
	class := buildRedirectClass(flatten)
	expectedClass := &client.RedirectClass{
		Status5XX: &status5XX,
		Status401: &status401,
	}

	return class, expectedClass, flatten
}

func generateTestGvcSidecar(stringifiedJson string) (*client.GvcSidecar, *client.GvcSidecar, []interface{}) {
	// Attempt to unmarshal `envoy`
	var envoy interface{}

	json.Unmarshal([]byte(stringifiedJson), &envoy)
	jsonOut, _ := json.Marshal(envoy)

	flatten := generateFlatTestGvcSidecar(string(jsonOut))
	sidecar := buildGvcSidecar(flatten)
	expectedSidecar := &client.GvcSidecar{
		Envoy: &envoy,
	}

	return sidecar, expectedSidecar, flatten
}

// Flatten //
func generateFlatTestLoadBalancer(dedicated bool, trustedProxies int, redirect []interface{}, ipset string) []interface{} {
	spec := map[string]interface{}{
		"dedicated":       dedicated,
		"trusted_proxies": trustedProxies,
		"redirect":        redirect,
		"ipset":           ipset,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestRedirect(class []interface{}) []interface{} {
	spec := map[string]interface{}{
		"class":                 class,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestRedirectClass(status5XX string, status401 string) []interface{} {
	spec := map[string]interface{}{
		"status_5xx":            status5XX,
		"status_401":            status401,
		"placeholder_attribute": true,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestGvcSidecar(envoy string) []interface{} {
	spec := map[string]interface{}{
		"envoy": envoy,
	}

	return []interface{}{
		spec,
	}
}
