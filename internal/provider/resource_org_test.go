package cpln

import (
	"testing"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	observability := buildObservability(generateFlatTestOrgObservability(*expectedObservability.LogsRetentionDays, *expectedObservability.MetricsRetentionDays, *expectedObservability.TracesRetentionDays))

	if diff := deep.Equal(observability, expectedObservability); diff != nil {
		t.Errorf("Org Observability was not built correctly, Diff: %s", diff)
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
	expectedFlatten := generateFlatTestOrgObservability(*expectedObservability.LogsRetentionDays, *expectedObservability.MetricsRetentionDays, *expectedObservability.TracesRetentionDays)
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

	expectedObservability := client.Observability{
		LogsRetentionDays:    &logsRetentionDays,
		MetricsRetentionDays: &metricsRetentionDays,
		TracesRetentionDays:  &tracesRetentionDays,
	}

	return &expectedObservability
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

func generateFlatTestOrgObservability(logsRetentionDays int, metricsRetentionDays int, tracesRetentionDays int) []interface{} {

	spec := map[string]interface{}{
		"logs_retention_days":    logsRetentionDays,
		"metrics_retention_days": metricsRetentionDays,
		"traces_retention_days":  tracesRetentionDays,
	}

	return []interface{}{
		spec,
	}
}
