package cpln

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var OrgName string = os.Getenv("CPLN_ORG")
var TestProvider *CplnProvider
var TestLoggerContext = context.Background()

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
func GetProviderServer() map[string]func() (tfprotov6.ProviderServer, error) {
	// Initialize a new instance of the CplnProvider using the "test" version
	p := New("test")()

	// Type assert the newly created provider instance to *CplnProvider
	TestProvider = p.(*CplnProvider)

	// Return a map of provider factories for Terraform testing framework
	return map[string]func() (tfprotov6.ProviderServer, error){
		"cpln": providerserver.NewProtocol6WithError(p),
	}
}

// testAccPreCheck verifies that all required environment variables are set before running
// an acceptance test. This function ensures that the test has sufficient configuration to
// interact with the API and outputs a header message to indicate the start of a specific test.
func testAccPreCheck(t *testing.T, testAccName string) {
	// Check for required organization name environment variable
	if OrgName == "" {
		t.Fatal("CPLN_ORG must be set for acceptance tests")
	}

	// Check for required API endpoint environment variable
	if endpoint := os.Getenv("CPLN_ENDPOINT"); endpoint == "" {
		t.Fatal("CPLN_ENDPOINT must be set for acceptance tests")
	}

	// Retrieve optional authentication parameters (profile or token)
	profile := os.Getenv("CPLN_PROFILE")
	token := os.Getenv("CPLN_TOKEN")

	// Ensure that either CPLN_PROFILE or CPLN_TOKEN is set for authentication
	if profile == "" && token == "" {
		t.Fatal("CPLN_PROFILE or CPLN_TOKEN must be set for acceptance tests")
	}

	// Log a header message indicating the start of the specified test
	tflog.Info(TestLoggerContext, "*********************************************************************")
	tflog.Info(TestLoggerContext, fmt.Sprintf("   TERRAFORM PROVIDER - CONTROL PLANE - %s ACCEPTANCE TESTS", testAccName))
	tflog.Info(TestLoggerContext, "*********************************************************************")
}
