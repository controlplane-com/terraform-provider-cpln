package cpln

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Structs ***/

// ProviderTestCase defines a specific resource test scenario.
type ProviderTestCase struct {
	Kind              string
	ResourceName      string
	ResourceAddress   string
	Name              string
	GvcName           string
	Description       string
	DescriptionUpdate string
}

// GetResourceNameAttr construct the resource name attribute of the specified resource.
func (ptc *ProviderTestCase) GetResourceNameAttr() string {
	return fmt.Sprintf("%s.name", ptc.ResourceAddress)
}

// GetSelfLink construct the self link of the specified resource.
func (ptc *ProviderTestCase) GetSelfLink() string {
	if ptc.Kind == "org" {
		return fmt.Sprintf("/org/%s", ptc.Name)
	}

	if ptc.GvcName != "" {
		return GetSelfLinkWithGvc(OrgName, ptc.Kind, ptc.GvcName, ptc.Name)
	}

	return GetSelfLink(OrgName, ptc.Kind, ptc.Name)
}

// GetSelfLinkAttr construct the self_link attribute of the specified resource.
func (ptc *ProviderTestCase) GetSelfLinkAttr() string {
	return fmt.Sprintf("%s.self_link", ptc.ResourceAddress)
}

// GetDefaultChecks returns a composed TestCheckFunc that verifies the default attributes of the resource for the specified kind.
func (ptc *ProviderTestCase) GetDefaultChecks(description string, tagsCount string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(ptc.ResourceAddress, "id", ptc.Name),
		resource.TestCheckResourceAttr(ptc.ResourceAddress, "name", ptc.Name),
		resource.TestCheckResourceAttr(ptc.ResourceAddress, "description", description),
		resource.TestCheckResourceAttr(ptc.ResourceAddress, "tags.%", tagsCount),
		resource.TestCheckResourceAttr(ptc.ResourceAddress, "self_link", ptc.GetSelfLink()),
	)
}

// TestCheckResourceAttr generates a TestCheckFunc to verify the count and members of a set attribute for the resource.
func (ptc *ProviderTestCase) TestCheckResourceAttr(key string, value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(ptc.ResourceAddress, key, value)
}

// TestCheckSetAttr generates a TestCheckFunc to verify the count and members of a set attribute for the resource.
func (ptc *ProviderTestCase) TestCheckSetAttr(key string, value []string) resource.TestCheckFunc {
	// Initialize slice of TestCheckFunc with a count check
	checks := []resource.TestCheckFunc{
		// Verify that the resource set has the expected number of elements
		resource.TestCheckResourceAttr(ptc.ResourceAddress, fmt.Sprintf("%s.#", key), fmt.Sprint(len(value))),
	}

	// Append a check for each element value in the set attribute
	for _, item := range value {
		// Add TestCheckTypeSetElemAttr for the current item
		checks = append(checks, resource.TestCheckTypeSetElemAttr(ptc.ResourceAddress, fmt.Sprintf("%s.*", key), item))
	}

	// Compose all checks into a single aggregate TestCheckFunc
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// TestCheckSetAttr generates a TestCheckFunc to verify the count and members of a map attribute for the resource.
func (ptc *ProviderTestCase) TestCheckMapAttr(key string, value map[string]string) resource.TestCheckFunc {
	// Initialize slice of TestCheckFunc with a count check
	checks := []resource.TestCheckFunc{
		// Verify that the resource map has the expected number of elements
		resource.TestCheckResourceAttr(ptc.ResourceAddress, fmt.Sprintf("%s.%%", key), fmt.Sprint(len(value))),
	}

	// Append a check for each element key-value in the map attribute
	for _key, _value := range value {
		// Add TestCheckResourceAttr for the current item
		checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, fmt.Sprintf("%s.%s", key, _key), _value))
	}

	// Compose all checks into a single aggregate TestCheckFunc
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// TestCheckNestedBlocks verifies that a nested block attribute contains the expected values.
func (ptc *ProviderTestCase) TestCheckNestedBlocks(key string, value []map[string]interface{}) resource.TestCheckFunc {
	// Initialize a slice to collect all the test check functions
	var checks []resource.TestCheckFunc

	// Append a check for the number of top-level elements in the list
	checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, fmt.Sprintf("%s.#", key), strconv.Itoa(len(value))))

	// Loop through each block and validate its fields
	for idx, block := range value {
		for _key, _value := range block {
			// Construct the full path for the field, e.g., "opaque.0.payload"
			path := fmt.Sprintf("%s.%d.%s", key, idx, _key)

			// Check the value depending on its type
			switch v := _value.(type) {
			case string:
				// Check a string value
				checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, path, v))

			case bool:
				// Check a boolean value, converted to string
				checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, path, strconv.FormatBool(v)))

			case int:
				// Check an integer value
				checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, path, strconv.Itoa(v)))

			case int32:
				// Check a 32-bit integer value
				checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, path, strconv.FormatInt(int64(v), 10)))

			case int64:
				// Check a 64-bit integer value
				checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, path, strconv.FormatInt(v, 10)))

			case float64:
				// Check a float64 value with flexible precision
				checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, path, strconv.FormatFloat(v, 'f', -1, 64)))

			case []string:
				// Check a string slice value
				checks = append(checks, ptc.TestCheckSetAttr(path, v))

			case []int:
				// Check an int slice value
				checks = append(checks, ptc.TestCheckSetAttr(path, IntSliceToStringSlice(v)))

			case []map[string]interface{}:
				// Recursively check a nested list of blocks
				checks = append(checks, ptc.TestCheckNestedBlocks(path, v))

			case map[string]interface{}:
				// Treat a nested block as a list of one and recurse
				checks = append(checks, ptc.TestCheckMapAttr(path, ConvertMapToStringMap(v)))

			default:
				// Panic on unsupported types
				panic(fmt.Sprintf("unsupported type %T for nested attribute '%s'", v, path))
			}
		}
	}

	// Aggregate all the checks into a single test check function
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

/*** Variables ***/

// Declare global test variables
var TestProvider *CplnProvider
var TestLoggerContext context.Context = context.Background()
var TestDataDirectoryPath string = "../../testdata"

/*** Functions ***/

// testAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
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

// MustLoadTestData loads the contents of a file from the testdata directory as a string and fails the test if it cannot be read.
func MustLoadTestData(filename string) string {
	// Construct the full file path relative to the testdata directory
	path := filepath.Join(TestDataDirectoryPath, filename)

	// Attempt to read the file content
	data, err := os.ReadFile(path)

	// Fail the test immediately if reading fails
	if err != nil {
		panic(fmt.Sprintf("failed to read %s: %v", path, err))
	}

	// Return the file contents as a string
	return string(data)
}

// testAccPreCheck verifies that all required environment variables are set before running an acceptance test.
func testAccPreCheck(t *testing.T, testAccName string) {
	// Check for required organization name environment variable
	if OrgName == "" {
		t.Fatal("CPLN_ORG must be set for acceptance tests")
	}

	// // Check for required API endpoint environment variable
	// if endpoint := os.Getenv("CPLN_ENDPOINT"); endpoint == "" {
	// 	t.Fatal("CPLN_ENDPOINT must be set for acceptance tests")
	// }

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
