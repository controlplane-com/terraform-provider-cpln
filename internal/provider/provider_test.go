package cpln

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var TestLogger *log.Logger

func init() {

	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cpln": testAccProvider,
	}

	infoLog := false

	if infoLog {
		TestLogger = log.New(os.Stdout, "TEST LOGGER: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		TestLogger = log.New(io.Discard, "TEST LOGGER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T, testAccName string) {

	if org := os.Getenv("CPLN_ORG"); org == "" {
		t.Fatal("CPLN_ORG must be set for acceptance tests")
	}

	if endpoint := os.Getenv("CPLN_ENDPOINT"); endpoint == "" {
		t.Fatal("CPLN_ENDPOINT must be set for acceptance tests")
	}

	profile := os.Getenv("CPLN_PROFILE")
	token := os.Getenv("CPLN_TOKEN")

	if profile == "" && token == "" {
		t.Fatal("CPLN_PROFILE or CPLN_TOKEN must be set for acceptance tests")
	}

	TestLogger.Print("*********************************************************************")
	TestLogger.Printf("   TERRAFORM PROVIDER - CONTROL PLANE - %s ACCEPTANCE TESTS", testAccName)
	TestLogger.Print("*********************************************************************")
}

// func testAccPreCheckGoogle(t *testing.T, testAccName string) {

// 	if validateDomains := os.Getenv("VALIDATE_DOMAINS"); validateDomains != "false" {

// 		if endpoint := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); endpoint == "" {
// 			t.Fatal("GOOGLE_APPLICATION_CREDENTIALS must be set for acceptance tests")
// 		}

// 		if endpoint := os.Getenv("GOOGLE_PROJECT"); endpoint == "" {
// 			t.Fatal("GOOGLE_PROJECT must be set for acceptance tests")
// 		}
// 	}

// 	testAccPreCheck(t, testAccName)
// }
