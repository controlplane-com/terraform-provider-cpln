package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneDomain_basic(t *testing.T) {

	var testDomain client.Domain

	dName := "cors2.cplntest.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "DOMAIN") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneDomainCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneDomain(dName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneDomainExists(dName, &testDomain),
				),
			},
		},
	})
}

func testAccControlPlaneDomain(domainName string) string {

	TestLogger.Printf("Inside testAccControlPlaneDomain")

	return fmt.Sprintf(`

	resource "cpln_domain" "example_ns_subdomain" {

		name        = "%s"
		description = "Custom domain that can be set on a GVC and used by associated workloads"
	  
		tags = {
		  terraform_generated = "true"
		  example             = "true"
		}
	  
		spec {
		  dns_mode         = "ns"
		  gvc_link         = "/org/myorg/gvc/mygvc"
		  accept_all_hosts = "true"
	  
		  ports {
			number   = 443
			protocol = "http"
	  
			cors {
			  allow_origins {
				exact = "example.com"
			  }
	  
			  allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
			  allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
			  max_age           = "24h"
			  //allow_credentials = "true"
			}
	  
			tls {
			  min_protocol_version = "TLSV1_2"
			  cipher_suites = [
				"ECDHE-ECDSA-AES256-GCM-SHA384",
				"ECDHE-ECDSA-CHACHA20-POLY1305",
				"ECDHE-ECDSA-AES128-GCM-SHA256",
				"ECDHE-RSA-AES256-GCM-SHA384",
				"ECDHE-RSA-CHACHA20-POLY1305",
				"ECDHE-RSA-AES128-GCM-SHA256",
				"AES256-GCM-SHA384",
				"AES128-GCM-SHA256",
			  ]
	  
			  client_certificate {}
			}
		  }
		}
	  }
	`, domainName)
	/*return fmt.Sprintf(`

	resource "cpln_domain" "example" {

		name        = "%s"
		description = "Test hakan"

		tags = {
			terraform_generated = "true"
			example             = "true"
		}

		spec {
			dns_mode         = "ns"
			accept_all_hosts = "true"

			ports {
				number   = 443
				protocol = "http"

				routes {
					prefix        = "/log"
					workload_link = "/org/efe/gvc/kadir/workload/a-log"
					port          = 8080
				}

				routes {
					prefix        = "/canary"
					workload_link = "/org/efe/gvc/kadir/workload/canary"
					port          = 8080
				}
			}
		}
	}
	`, domainName)*/
}

func testAccCheckControlPlaneDomainExists(domainName string, domain *client.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneWorkloadExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[domainName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != domainName {
			return fmt.Errorf("Workload name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)

		d, _, err := client.GetDomain(domainName)

		if err != nil {
			return err
		}

		if *d.Name != domainName {
			return fmt.Errorf("Workload name does not match")
		}

		*domain = *d

		return nil
	}
}

func testAccCheckControlPlaneDomainCheckDestroy(s *terraform.State) error {

	TestLogger.Printf("Inside testAccCheckControlPlaneDomainCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	return nil
}

/*** Unit Tests ***/
// Build Domain Spec //
func TestControlPlane_BuildDomainSpec(t *testing.T) {
	dnsMode := "ns"
	gvcLink := "/org/myorg/gvc/mygvc"
	acceptAllHosts := true
	_, expectedPorts, flattenedPorts := generatePorts()

	domainSpec := buildDomainSpec(generateFlatTestDomainSpec(dnsMode, gvcLink, acceptAllHosts, flattenedPorts))
	expectedDomainSpec := client.DomainSpec{
		DnsMode:        &dnsMode,
		GvcLink:        &gvcLink,
		AcceptAllHosts: &acceptAllHosts,
		Ports:          &expectedPorts,
	}

	if diff := deep.Equal(domainSpec, &expectedDomainSpec); diff != nil {
		t.Errorf("DomainSpec was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildDomainSpec_NoPorts(t *testing.T) {
	dnsMode := "ns"
	gvcLink := "/org/myorg/gvc/mygvc"
	acceptAllHosts := false

	domainSpec := buildDomainSpec(generateFlatTestDomainSpec(dnsMode, gvcLink, acceptAllHosts, nil))
	expectedDomainSpec := client.DomainSpec{
		DnsMode:        &dnsMode,
		GvcLink:        &gvcLink,
		AcceptAllHosts: &acceptAllHosts,
		Ports:          nil,
	}

	if diff := deep.Equal(domainSpec, &expectedDomainSpec); diff != nil {
		t.Errorf("DomainSpec was not built correctly. Diff: %s", diff)
	}
}

// Build Ports //
func TestControlPlane_BuildPorts(t *testing.T) {
	ports, expectedPorts, _ := generatePorts()
	if diff := deep.Equal(ports, &expectedPorts); diff != nil {
		t.Errorf("Ports were not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildPorts_Empty(t *testing.T) {
	ports := buildSpecPorts(generateEmptyInterfaceArray())
	expectedPorts := []client.DomainSpecPort{{}}

	if diff := deep.Equal(ports, &expectedPorts); diff != nil {
		t.Errorf("Ports were not built correctly. Diff: %s", diff)
	}
}

// Build Routes //
func TestControlPlane_BuildRoutes(t *testing.T) {
	routes, expectedRoutes, _ := generateRoutes()
	if diff := deep.Equal(routes, &expectedRoutes); diff != nil {
		t.Errorf("Routes were not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildRoutes_Empty(t *testing.T) {
	routes := buildRoutes(generateEmptyInterfaceArray())
	expectedRoutes := []client.DomainRoute{{}}

	if diff := deep.Equal(routes, &expectedRoutes); diff != nil {
		t.Errorf("Routes were not built correctly. Diff: %s", diff)
	}
}

// Build Cors //
func TestControlPlane_BuildCors(t *testing.T) {
	cors, expectedCors, _ := generateCors()
	if diff := deep.Equal(cors, &expectedCors); diff != nil {
		t.Errorf("Cors was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildCors_Empty(t *testing.T) {
	cors := buildCors(generateEmptyInterfaceArray())
	expectedCors := client.DomainCors{}

	if diff := deep.Equal(cors, &expectedCors); diff != nil {
		t.Errorf("Cors was not built correctly. Diff: %s", diff)
	}
}

// Build TLS Unit Test //
func TestControlPlane_BuildTLS(t *testing.T) {
	tls, expectedTLS, _ := generateTLS()
	if diff := deep.Equal(tls, &expectedTLS); diff != nil {
		t.Errorf("TLS was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildTLS_Empty(t *testing.T) {
	tls := buildTLS(generateEmptyInterfaceArray())
	expectedTLS := client.DomainTLS{}

	if diff := deep.Equal(tls, &expectedTLS); diff != nil {
		t.Errorf("TLS was not built correctly. Diff: %s", diff)
	}
}

// Build Allow Origins Unit Test //
func TestControlPlane_BuildAllowOrigins(t *testing.T) {
	collection, expectedCollection, _ := generateAllowOrigins()
	if diff := deep.Equal(collection, &expectedCollection); diff != nil {
		t.Errorf("Allow Origins was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildAllowOrigins_WithoutExact(t *testing.T) {
	collection := buildAllowOrigins(generateEmptyInterfaceArray())
	expectedCollection := []client.DomainAllowOrigin{{}}

	if diff := deep.Equal(collection, expectedCollection); diff != nil {
		t.Errorf("Allow Origins was not built correctly. Diff: %s", diff)
	}
}

// Build Certificate Unit Test //
func TestControlPlane_BuildCertificate(t *testing.T) {
	secret := "/org/myorg/secret/mysecret"

	cert := buildCertificate(generateFlatTestCertificate(secret))
	expectedCert := client.DomainCertificate{SecretLink: &secret}

	// TODO move expectedCert to a function, can be array of items too for different cases
	if diff := deep.Equal(cert, &expectedCert); diff != nil {
		t.Errorf("Domain Certificate was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildCertificate_WithoutSecret(t *testing.T) {
	cert := buildCertificate(generateEmptyInterfaceArray())
	certTest := client.DomainCertificate{}

	if diff := deep.Equal(cert, &certTest); diff != nil {
		t.Errorf("Domain Certificate was not built correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
func generatePorts() (*[]client.DomainSpecPort, []client.DomainSpecPort, []interface{}) {
	number := 443
	protocol := "http"

	_, expectedRoutes, flattenedRoutes := generateRoutes()
	_, expectedCors, flattenedCors := generateCors()
	_, expectedTLS, flattenedTLS := generateTLS()

	flattenGeneration := generateInterfaceArrayFromMapArray([]map[string]interface{}{
		generateFlatTestPort(number, protocol, flattenedRoutes, flattenedCors, flattenedTLS),
	})

	ports := buildSpecPorts(flattenGeneration)
	expectedPorts := []client.DomainSpecPort{
		{
			Number:   &number,
			Protocol: &protocol,
			Routes:   &expectedRoutes,
			Cors:     &expectedCors,
			TLS:      &expectedTLS,
		},
	}

	return ports, expectedPorts, flattenGeneration
}

func generateRoutes() (*[]client.DomainRoute, []client.DomainRoute, []interface{}) {
	prefix := "/"
	replacePrefix := "/replace"
	workload_link := "/org/myorg/gvc/mygvc/workload_two"
	port := 8080

	flattened := generateInterfaceArrayFromMapArray([]map[string]interface{}{
		generateFlatTestRoute(prefix, replacePrefix, workload_link, port),
	})

	routes := buildRoutes(flattened)
	expectedRoutes := []client.DomainRoute{{
		Prefix:        &prefix,
		ReplacePrefix: &replacePrefix,
		WorkloadLink:  &workload_link,
		Port:          &port,
	}}

	return routes, expectedRoutes, flattened
}

func generateCors() (*client.DomainCors, client.DomainCors, []interface{}) {
	allowMethods := []string{"GET", "POST"}
	allowHeaders := []string{"Content-Type", "Accept", "Authorization"}
	maxAge := "24h"
	allowCredentials := true

	_, expectedAllowOrigins, flattenedAllowOrigins := generateAllowOrigins()
	flattened := generateFlatTestCors(flattenedAllowOrigins, flattenStringsArray(&allowMethods), flattenStringsArray(&allowHeaders), maxAge, allowCredentials)

	cors := buildCors(flattened)
	expectedCors := client.DomainCors{
		AllowOrigins:     &expectedAllowOrigins,
		AllowMethods:     &allowMethods,
		AllowHeaders:     &allowHeaders,
		MaxAge:           &maxAge,
		AllowCredentials: &allowCredentials,
	}

	return cors, expectedCors, flattened
}

func generateTLS() (*client.DomainTLS, client.DomainTLS, []interface{}) {
	minProtocolVersion := "TLSv1_1"
	cipherSuites := []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA"}
	clientSecret := "/org/myorg/secret/mysecret_client"
	serverSecret := "/org/myorg/secret/mysecret_server"

	cipherSuitesFlattened := flattenStringsArray(&cipherSuites)
	clientCertificate := generateFlatTestCertificate(clientSecret)
	serverCertificate := generateFlatTestCertificate(serverSecret)

	flattened := generateFlatTestTLS(minProtocolVersion, cipherSuitesFlattened, clientCertificate, serverCertificate)

	tls := buildTLS(flattened)
	expectedTLS := client.DomainTLS{
		MinProtocolVersion: &minProtocolVersion,
		CipherSuites:       &cipherSuites,
		ClientCertificate:  &client.DomainCertificate{SecretLink: &clientSecret},
		ServerCertificate:  &client.DomainCertificate{SecretLink: &serverSecret},
	}

	return tls, expectedTLS, flattened
}

func generateAllowOrigins() (*[]client.DomainAllowOrigin, []client.DomainAllowOrigin, []interface{}) {
	exact := "example.com"
	flattened := generateFlatTestAllowOrigins(exact)

	collection := buildAllowOrigins(flattened)
	expectedCollection := []client.DomainAllowOrigin{{Exact: &exact}}

	return collection, expectedCollection, flattened
}

/*** Flatten ***/
func generateFlatTestDomainSpec(dnsMode string, gvcLink string, acceptAllHosts bool, ports []interface{}) []interface{} {
	spec := map[string]interface{}{
		"dns_mode":         dnsMode,
		"gvc_link":         gvcLink,
		"accept_all_hosts": acceptAllHosts,
		"ports":            ports,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestPort(number int, protocol string, routes []interface{}, cors []interface{}, tls []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"number":   number,
		"protocol": protocol,
		"routes":   routes,
		"cors":     cors,
		"tls":      tls,
	}
}

func generateFlatTestRoute(prefix string, replacePrefix string, workloadLink string, port int) map[string]interface{} {
	return map[string]interface{}{
		"prefix":         prefix,
		"replace_prefix": replacePrefix,
		"workload_link":  workloadLink,
		"port":           port,
	}
}

func generateFlatTestCors(allowOrigins []interface{}, allowMethods []interface{}, allowHeaders []interface{}, maxAge string, allowCredentials bool) []interface{} {
	spec := map[string]interface{}{
		"allow_origins":     allowOrigins,
		"allow_methods":     allowMethods,
		"allow_headers":     allowHeaders,
		"max_age":           maxAge,
		"allow_credentials": allowCredentials,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestTLS(minProtocolVersion string, cipherSuites []interface{}, clientCertificate []interface{}, serverCertificate []interface{}) []interface{} {
	spec := map[string]interface{}{
		"min_protocol_version": minProtocolVersion,
		"cipher_suites":        cipherSuites,
		"client_certificate":   clientCertificate,
		"server_certificate":   serverCertificate,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestAllowOrigins(exact string) []interface{} {
	spec := map[string]interface{}{
		"exact": exact,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestCertificate(secretLink string) []interface{} {
	spec := map[string]interface{}{
		"secret_link": secretLink,
	}

	return []interface{}{
		spec,
	}
}

func generateInterfaceArrayFromMapArray(specs []map[string]interface{}) []interface{} {
	collection := make([]interface{}, len(specs))
	for i, spec := range specs {
		collection[i] = spec
	}

	return collection
}

func generateEmptyInterfaceArray() []interface{} {
	return []interface{}{
		map[string]interface{}{},
	}
}
