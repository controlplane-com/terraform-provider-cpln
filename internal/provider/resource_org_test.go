package cpln

import (
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccControlPlaneOrg_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "ORG") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneOrgCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneOrg(),
			},
		},
	})
}

func testAccControlPlaneOrg() string {

	TestLogger.Printf("Inside testAccControlPlaneOrg")

	return `
	  resource "cpln_org" "example" {
  
		session_timeout_seconds = 50000
		description = "testing"

		tags = {
			terraform_generated = "true"
			example             = "true"
    	}
	  
		observability {
		  logs_retention_days    = 55
		  metrics_retention_days = 65
		  traces_retention_days  = 75
			default_alert_emails   = ["bob@example.com", "rob@example.com", "abby@example.com"]
		}

		security {
			threat_detection {
				enabled 		 = true
				minimum_severity = "warning"
				syslog {
					transport = "tcp"
					host 	  = "example.com"
					port  	  = 8080
				}
			}
		}
	  }
    `
}

/*** Unit Tests ***/
// Build //
func TestControlPlane_BuildOrgAuthConfig(t *testing.T) {

	expectedAuthConfig := generateTestOrgAuthConfig()
	authConfig := buildAuthConfig(generateFlatTestOrgAuthConfig(true, *expectedAuthConfig.DomainAutoMembers, *expectedAuthConfig.SamlOnly))

	if diff := deep.Equal(authConfig, expectedAuthConfig); diff != nil {
		t.Errorf("Org Auth Config was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildOrgObservability(t *testing.T) {

	expectedObservability := generateTestOrgObservability()
	observability := buildObservability(generateFlatTestOrgObservability(true, *expectedObservability.LogsRetentionDays, *expectedObservability.MetricsRetentionDays, *expectedObservability.TracesRetentionDays, *expectedObservability.DefaultAlertEmails))

	if diff := deep.Equal(observability, expectedObservability); diff != nil {
		t.Errorf("Org Observability was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildOrgSecurity(t *testing.T) {

	expectedSecurity := generateTestOrgSecurity()
	flattenedSyslog := generateFlatTestOrgThreatDetectionSyslog(*expectedSecurity.ThreatDetection.Syslog.Transport, *expectedSecurity.ThreatDetection.Syslog.Host, *expectedSecurity.ThreatDetection.Syslog.Port)
	security := buildOrgSecurity(generateFlatTestOrgSecurity(generateFlatTestOrgThreatDetection(*expectedSecurity.ThreatDetection.Enabled, *expectedSecurity.ThreatDetection.MinimumSeverity, flattenedSyslog)))

	if diff := deep.Equal(security, expectedSecurity); diff != nil {
		t.Errorf("Org Security was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildOrgThreatDetection(t *testing.T) {

	expectedThreatDetection := generateTestOrgThreatDetection()
	flattenedSyslog := generateFlatTestOrgThreatDetectionSyslog(*expectedThreatDetection.Syslog.Transport, *expectedThreatDetection.Syslog.Host, *expectedThreatDetection.Syslog.Port)
	threatDetection := buildOrgThreatDetection(generateFlatTestOrgThreatDetection(*expectedThreatDetection.Enabled, *expectedThreatDetection.MinimumSeverity, flattenedSyslog))

	if diff := deep.Equal(threatDetection, expectedThreatDetection); diff != nil {
		t.Errorf("Org Threat Detection was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildOrgThreatDetectionSyslog(t *testing.T) {

	expectedSyslog := generateTestOrgThreatDetectionSyslog()
	syslog := buildOrgThreatDetectionSyslog(generateFlatTestOrgThreatDetectionSyslog(*expectedSyslog.Transport, *expectedSyslog.Host, *expectedSyslog.Port))

	if diff := deep.Equal(syslog, expectedSyslog); diff != nil {
		t.Errorf("Org Threat Detection Syslog was not built correctly, Diff: %s", diff)
	}
}

// Flatten //
func TestControlPlane_FlattenOrgAuthConfig(t *testing.T) {

	authConfig := generateTestOrgAuthConfig()
	expectedFlatten := generateFlatTestOrgAuthConfig(false, *authConfig.DomainAutoMembers, *authConfig.SamlOnly)
	flattenedAuthConfig := flattenAuthConfig(authConfig)

	if diff := deep.Equal(expectedFlatten, flattenedAuthConfig); diff != nil {
		t.Errorf("Org Auth Config was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenOrgObservability(t *testing.T) {

	expectedObservability := generateTestOrgObservability()
	expectedFlatten := generateFlatTestOrgObservability(false, *expectedObservability.LogsRetentionDays, *expectedObservability.MetricsRetentionDays, *expectedObservability.TracesRetentionDays, *expectedObservability.DefaultAlertEmails)
	flattenedObservability := flattenObservability(expectedObservability)

	if diff := deep.Equal(expectedFlatten, flattenedObservability); diff != nil {
		t.Errorf("Org Observability was not flattened correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
// Build //
func generateTestOrgAuthConfig() *client.AuthConfig {

	domainAutoMembers := []string{"example.com"}
	samlOnly := false

	expectedAuthConfig := client.AuthConfig{
		DomainAutoMembers: &domainAutoMembers,
		SamlOnly:          &samlOnly,
	}

	return &expectedAuthConfig
}

func generateTestOrgObservability() *client.Observability {

	logsRetentionDays := 60
	metricsRetentionDays := 50
	tracesRetentionDays := 40
	defaultAlertEmails := []string{"bob@example.com", "rob@example.com", "abby@example.com"}

	expectedObservability := client.Observability{
		LogsRetentionDays:    &logsRetentionDays,
		MetricsRetentionDays: &metricsRetentionDays,
		TracesRetentionDays:  &tracesRetentionDays,
		DefaultAlertEmails:   &defaultAlertEmails,
	}

	return &expectedObservability
}

func generateTestOrgSecurity() *client.OrgSecurity {

	threatDetection := generateTestOrgThreatDetection()

	expectedSecurity := client.OrgSecurity{
		ThreatDetection: threatDetection,
	}

	return &expectedSecurity
}

func generateTestOrgThreatDetection() *client.OrgThreatDetection {

	enabled := true
	minimumSeverity := "warning"
	syslog := generateTestOrgThreatDetectionSyslog()

	expectedThreatDetection := client.OrgThreatDetection{
		Enabled:         &enabled,
		MinimumSeverity: &minimumSeverity,
		Syslog:          syslog,
	}

	return &expectedThreatDetection
}

func generateTestOrgThreatDetectionSyslog() *client.OrgThreatDetectionSyslog {

	transport := "tcp"
	host := "example.com"
	port := 8080

	expectedSyslog := client.OrgThreatDetectionSyslog{
		Transport: &transport,
		Host:      &host,
		Port:      &port,
	}

	return &expectedSyslog
}

// Flatten //
func generateFlatTestOrgAuthConfig(useSet bool, domainAutoMembers []string, samlOnly bool) []interface{} {

	stringFunc := schema.HashSchema(StringSchema())
	interfaceSlice := make([]interface{}, len(domainAutoMembers))

	for i, v := range domainAutoMembers {
		interfaceSlice[i] = v
	}

	spec := map[string]interface{}{
		"saml_only": samlOnly,
	}

	if useSet {
		spec["domain_auto_members"] = schema.NewSet(stringFunc, interfaceSlice)
	} else {
		spec["domain_auto_members"] = interfaceSlice
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestOrgObservability(useSet bool, logsRetentionDays int, metricsRetentionDays int, tracesRetentionDays int, defaultAlertEmails []string) []interface{} {

	defaultAlertEmailsInterfaceSlice := make([]interface{}, len(defaultAlertEmails))

	for i, v := range defaultAlertEmails {
		defaultAlertEmailsInterfaceSlice[i] = v
	}

	spec := map[string]interface{}{
		"logs_retention_days":    logsRetentionDays,
		"metrics_retention_days": metricsRetentionDays,
		"traces_retention_days":  tracesRetentionDays,
	}

	if useSet {
		stringFunc := schema.HashSchema(StringSchema())
		spec["default_alert_emails"] = schema.NewSet(stringFunc, defaultAlertEmailsInterfaceSlice)
	} else {
		spec["default_alert_emails"] = defaultAlertEmailsInterfaceSlice
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestOrgSecurity(threatDetection []interface{}) []interface{} {

	spec := map[string]interface{}{
		"threat_detection": threatDetection,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestOrgThreatDetection(enabled bool, minimumSeverity string, syslog []interface{}) []interface{} {

	spec := map[string]interface{}{
		"enabled":          enabled,
		"minimum_severity": minimumSeverity,
		"syslog":           syslog,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestOrgThreatDetectionSyslog(transport string, host string, port int) []interface{} {

	spec := map[string]interface{}{
		"transport": transport,
		"host":      host,
		"port":      port,
	}

	return []interface{}{
		spec,
	}
}
