package cpln

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
func (ptc *ProviderTestCase) TestCheckSetAttr(key string, expectedValues []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// Look up the resource from the Terraform state by its address
		rs, ok := state.RootModule().Resources[ptc.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource %q not in state", ptc.ResourceAddress)
		}

		// Get the raw attribute map of the resource
		attrs := rs.Primary.Attributes

		// Read the actual number of elements in the set from the state
		actualCount := ptc.ReadCollectionElementCount(attrs, key)

		// If both expected and actual sets are empty, succeed early
		if actualCount == 0 && len(expectedValues) == 0 {
			return nil
		}

		// Fail if the element count does not match the expected count
		if actualCount != len(expectedValues) {
			return fmt.Errorf("%s.# = %d; expected %d", key, actualCount, len(expectedValues))
		}

		// Prepare a slice to hold the actual values retrieved from the state
		actual := make([]string, 0, actualCount)

		// Collect actual values from the state using indexed keys
		for i := range actualCount {
			// Build the key path for the i-th element in the set
			ik := fmt.Sprintf("%s.%d", key, i)

			// Retrieve the value from the state attributes
			v, ok := attrs[ik]
			if !ok {
				return fmt.Errorf("missing %s in state", ik)
			}

			// Append the value to the actual slice
			actual = append(actual, v)
		}

		// Copy expected values into a new slice to avoid mutating the input
		expected := make([]string, len(expectedValues))
		copy(expected, expectedValues)

		// Sort both actual and expected slices for order-independent comparison
		sort.Strings(actual)
		sort.Strings(expected)

		// Compare the sorted slices and fail if they differ
		if !reflect.DeepEqual(actual, expected) {
			return fmt.Errorf("%s values mismatch\n  got:  %v\n  expected: %v", key, actual, expected)
		}

		// Return success if all checks passed
		return nil
	}
}

// TestCheckSetAttr generates a TestCheckFunc to verify the count and members of a map attribute for the resource.
func (ptc *ProviderTestCase) TestCheckMapAttr(key string, expected map[string]string) resource.TestCheckFunc {
	// Initialize count of non-null entries
	nonNullCount := 0

	// Count non-null entries in the expected map (nulls are treated as "must be absent")
	for _, v := range expected {
		// Only count non-null entries
		if v != "null" {
			nonNullCount++
		}
	}

	// Initialize slice of TestCheckFunc with a count check
	checks := []resource.TestCheckFunc{
		// Verify that the resource map has the expected number of elements
		resource.TestCheckResourceAttr(ptc.ResourceAddress, fmt.Sprintf("%s.%%", key), fmt.Sprint(nonNullCount)),
	}

	// Append a check for each element key-value in the map attribute
	for k, v := range expected {
		// Skip if value is null
		if v == "null" {
			continue
		}

		// Add TestCheckResourceAttr for the current item
		checks = append(checks, resource.TestCheckResourceAttr(ptc.ResourceAddress, fmt.Sprintf("%s.%s", key, k), v))
	}

	// Compose all checks into a single aggregate TestCheckFunc
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// TestCheckNestedBlocks verifies that a nested block attribute contains the expected values
func (ptc *ProviderTestCase) TestCheckNestedBlocks(key string, expected []map[string]interface{}) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// Retrieve the resource from the Terraform state using its address
		rs, ok := state.RootModule().Resources[ptc.ResourceAddress]

		// Return error if resource is not found in the state
		if !ok {
			return fmt.Errorf("resource %q not in state", ptc.ResourceAddress)
		}

		// Get the attributes of the resource from the state
		attrs := rs.Primary.Attributes

		// Read the actual number of nested blocks from the state
		actualCount := ptc.ReadCollectionElementCount(attrs, key)

		// Check if the actual count matches the expected count
		if actualCount != len(expected) {
			return fmt.Errorf("%s.# = %d; expected %d", key, actualCount, len(expected))
		}

		// Track which actual elements have been matched to expected elements
		used := make([]bool, actualCount)

		// Iterate over each expected nested block element
		for expectedIndex, expectedElem := range expected {
			// Initialize variable to store the index of a matching actual element
			matchIndex := -1

			// Iterate over actual elements to find a match for the expected element
			for i := range actualCount {
				// Skip actual elements that have already been matched
				if used[i] {
					continue
				}

				// Check if the actual element matches the expected element
				if ptc.MatchesElementAtIndex(attrs, key, i, expectedElem) {
					// Mark the actual element as used
					used[i] = true

					// Store the index of the matching actual element
					matchIndex = i

					// Build the base key for the matched element
					base := fmt.Sprintf("%s.%d", key, i)

					// Recursively check nested children for the matched element
					if err := ptc.checkNestedChildrenRecursively(state, base, expectedElem); err != nil {
						return fmt.Errorf("nested mismatch at %s (expected element #%d): %w", base, expectedIndex, err)
					}

					// Break out of the loop after finding a match
					break
				}
			}

			// If no match was found for the expected element, provide diagnostics
			if matchIndex == -1 {
				// Initialize a string builder for diagnostic output
				var diag strings.Builder

				// Iterate over actual elements to gather diagnostic information
				for i := range actualCount {
					// Build the base key for the candidate actual element
					idxBase := fmt.Sprintf("%s.%d", key, i)

					// Build token bags for state and expected values
					stateTokens := ptc.BuildStateTokenBag(attrs, idxBase)
					expectTokens := ptc.BuildExpectedTokenBag(fmt.Sprintf("%s.[*]", key), expectedElem)

					// Identify missing expected tokens in the state
					missing := ptc.MissingExpectedTokens(stateTokens, expectTokens)

					// Write candidate information to the diagnostic output
					fmt.Fprintf(&diag, "\n--- Candidate %s ---\n", idxBase)
					fmt.Fprintf(&diag, "State subtree under %s:\n%s", idxBase, ptc.DumpStateSubtree(attrs, idxBase))

					// Indicate whether missing tokens were found
					if len(missing) == 0 {
						fmt.Fprintf(&diag, "No missing tokens (mismatch may be in primitive fields)\n")
					} else {
						fmt.Fprintf(&diag, "Missing expected tokens (%d):\n", len(missing))
						for _, m := range missing {
							fmt.Fprintf(&diag, "  %s\n", m)
						}
					}
				}

				// Return error indicating no match was found, along with diagnostics
				return fmt.Errorf("no match found in %s for expected element #%d: %+v\nDetails:%s",
					key, expectedIndex, expectedElem, diag.String())
			}
		}

		// Return nil if all expected elements were matched successfully
		return nil
	}
}

// ReadCollectionElementCount reads "<base>.#" and returns the list/set length.
func (ptc *ProviderTestCase) ReadCollectionElementCount(attrs map[string]string, base string) int {
	// Build the key for the collection count attribute
	key := fmt.Sprintf("%s.#", base)

	// Retrieve the raw count value from the attributes map
	raw, ok := attrs[key]

	// Return zero if the key is not present in the map
	if !ok {
		return 0
	}

	// Convert the raw string value to an integer
	n, _ := strconv.Atoi(raw)

	// Return the parsed integer count
	return n
}

// MatchesElementAtIndex determines if a candidate state element matches the expected element by comparing token bags.
func (ptc *ProviderTestCase) MatchesElementAtIndex(attrs map[string]string, base string, idx int, elem map[string]interface{}) bool {
	// Build the concrete state prefix for the candidate index
	statePrefix := fmt.Sprintf("%s.%d", base, idx)

	// Build the expected wildcard prefix to normalize indices
	expectPrefix := fmt.Sprintf("%s.[*]", base)
	expectPrefix = ptc.NormalizeIndexWildcards(expectPrefix)

	// Collect normalized tokens from the state under the candidate prefix
	stateTokens := ptc.BuildStateTokenBag(attrs, statePrefix)

	// Collect normalized tokens from the expected element
	expectTokens := ptc.BuildExpectedTokenBag(expectPrefix, elem)

	// Return true only if state covers all expected tokens (ignoring extra defaults)
	return ptc.StateTokenBagCoversExpected(stateTokens, expectTokens)
}

// CheckNestedChildrenRecursively walks all fields of an expected element under a resolved index and validates them against state.
func (ptc *ProviderTestCase) checkNestedChildrenRecursively(state *terraform.State, base string, elem map[string]interface{}) error {
	// Get the resource attributes from Terraform state
	rs := state.RootModule().Resources[ptc.ResourceAddress]
	attrs := rs.Primary.Attributes

	// Iterate through each expected field
	for fieldName, fieldVal := range elem {
		// Build the full key path for the field
		childKey := fmt.Sprintf("%s.%s", base, fieldName)

		// Dispatch by type of the expected value
		switch v := fieldVal.(type) {

		// Validate primitive values by direct equality
		case string, bool, int, int32, int64, float64:
			expected := ptc.PrimToString(v)
			got := attrs[childKey]
			if got != expected {
				return fmt.Errorf("%s = %q; expected %q", childKey, got, expected)
			}

		// Validate a slice of strings as a set (order-independent)
		case []string:
			if err := ptc.TestCheckSetAttr(childKey, v)(state); err != nil {
				return err
			}

		// Validate a slice of ints as a set (converted to strings)
		case []int:
			intAsStr := IntSliceToStringSlice(v)
			if err := ptc.TestCheckSetAttr(childKey, intAsStr)(state); err != nil {
				return err
			}

		// Validate nested block collections recursively
		case []map[string]interface{}:
			if err := ptc.TestCheckNestedBlocks(childKey, v)(state); err != nil {
				return err
			}

		// Validate map attributes with null-sentinel semantics
		case map[string]interface{}:
			// Convert expected map values to string form
			expectedMap := ConvertMapToStringMap(v)

			// Read actual map entries from state
			actualMap := ptc.ReadMapEntries(attrs, childKey)

			// Ensure no unexpected keys are present
			for actualKey := range actualMap {
				if _, ok := expectedMap[actualKey]; !ok {
					return fmt.Errorf("%s has unexpected key %q (value %q)", childKey, actualKey, actualMap[actualKey])
				}
			}

			// Validate all expected keys
			for expKey, expVal := range expectedMap {
				fullKey := fmt.Sprintf("%s.%s", childKey, expKey)
				gotVal, present := actualMap[expKey]

				// Null sentinel means key should be absent or empty
				if expVal == "null" {
					if present && gotVal != "" {
						return fmt.Errorf("%s should be absent or empty (expected null), but found %q", fullKey, gotVal)
					}
					continue
				}

				// Non-null entries must be present and match exactly
				if !present {
					return fmt.Errorf("missing required map entry %s", fullKey)
				}
				if gotVal != expVal {
					return fmt.Errorf("%s = %q; expected %q", fullKey, gotVal, expVal)
				}
			}

		// Reject unsupported nested types
		default:
			return fmt.Errorf("unsupported nested type at %s: %T", childKey, v)
		}
	}

	// Return success if all fields validated
	return nil
}

// DumpStateSubtree prints all attributes under a given prefix, excluding meta counters.
func (ptc *ProviderTestCase) DumpStateSubtree(attrs map[string]string, prefix string) string {
	// Ensure prefix ends with a dot
	p := prefix
	if !strings.HasSuffix(p, ".") {
		p += "."
	}

	// Collect keys that belong to the prefix and are not meta counters
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		if strings.HasPrefix(k, p) && !strings.HasSuffix(k, ".#") && !strings.HasSuffix(k, ".%") {
			keys = append(keys, k)
		}
	}

	// Sort keys for deterministic output
	sort.Strings(keys)

	// Build the dump string
	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s = %q\n", k, attrs[k])
	}
	return b.String()
}

// MissingExpectedTokens returns all expected tokens absent or undercounted in state.
func (ptc *ProviderTestCase) MissingExpectedTokens(state map[string]int, expect map[string]int) []string {
	// Track missing tokens
	missing := []string{}
	for k, need := range expect {
		have := state[k]
		if have < need {
			missing = append(missing, fmt.Sprintf("%s  (need %d, have %d)", k, need, have))
		}
	}

	// Sort for deterministic error output
	sort.Strings(missing)
	return missing
}

// StateTokenBagCoversExpected returns true if state contains all expected tokens with sufficient counts.
func (ptc *ProviderTestCase) StateTokenBagCoversExpected(state map[string]int, expect map[string]int) bool {
	for k, ve := range expect {
		if vs, ok := state[k]; !ok || vs < ve {
			return false
		}
	}

	return true
}

// BuildStateTokenBag flattens state attributes under a concrete prefix into normalized "path=value" tokens.
func (ptc *ProviderTestCase) BuildStateTokenBag(attrs map[string]string, concretePrefix string) map[string]int {
	// Initialize output map
	out := make(map[string]int)

	// Regex to normalize numeric indices in paths
	reIdx := regexp.MustCompile(`\.\d+`)

	// Prefix to match all relevant keys
	prefix := concretePrefix + "."

	// Walk through all attributes
	for key, val := range attrs {
		// Skip attributes outside the prefix
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// Skip meta counters
		if strings.HasSuffix(key, ".#") || strings.HasSuffix(key, ".%") {
			continue
		}

		// Normalize numeric indices in the path
		normPath := reIdx.ReplaceAllString(key, ".[*]")

		// Normalize numeric values to stable string form
		normVal := ptc.NormalizeNumberString(val)

		// Build token as "path=value"
		token := normPath + "=" + normVal

		// Increment token count
		out[token]++
	}
	return out
}

// BuildExpectedTokenBag flattens the expected element into normalized "path=value" tokens.
func (ptc *ProviderTestCase) BuildExpectedTokenBag(prefix string, elem map[string]interface{}) map[string]int {
	// Normalize prefix to use wildcards for indices
	base := ptc.NormalizeIndexWildcards(prefix)

	// Initialize output map
	out := make(map[string]int)

	// Recursively flatten the expected element
	ptc.buildExpectedTokenBagRec(out, base, elem)
	return out
}

// BuildExpectedTokenBagRec recursively flattens expected structures into token form.
func (ptc *ProviderTestCase) buildExpectedTokenBagRec(out map[string]int, base string, val interface{}) {
	switch v := val.(type) {

	// Map attributes or object fields
	case map[string]interface{}:
		for key, mv := range v {
			// Skip keys whose value is explicitly "null"
			if s, ok := mv.(string); ok && s == "null" {
				continue
			}
			childPath := ptc.NormalizeIndexWildcards(fmt.Sprintf("%s.%s", base, key))
			ptc.buildExpectedTokenBagRec(out, childPath, mv)
		}

	// List/set of nested blocks
	case []map[string]interface{}:
		childBase := ptc.NormalizeIndexWildcards(fmt.Sprintf("%s.[*]", base))
		for _, child := range v {
			ptc.buildExpectedTokenBagRec(out, childBase, child)
		}

	// List/set of strings
	case []string:
		path := ptc.NormalizeIndexWildcards(fmt.Sprintf("%s.[*]", base))
		for _, s := range v {
			token := path + "=" + ptc.NormalizeNumberString(s)
			out[token]++
		}

	// List/set of ints
	case []int:
		path := ptc.NormalizeIndexWildcards(fmt.Sprintf("%s.[*]", base))
		for _, n := range v {
			token := path + "=" + strconv.Itoa(n)
			out[token]++
		}

	// Primitive values
	case string, bool, int, int32, int64, float64:
		path := ptc.NormalizeIndexWildcards(base)
		token := path + "=" + ptc.PrimToString(v)
		out[token]++

	// Unsupported types
	default:
		panic(fmt.Sprintf("unsupported expected type at %s: %T", base, v))
	}
}

// NormalizeIndexWildcards replaces numeric indices in a path with wildcards.
func (ptc *ProviderTestCase) NormalizeIndexWildcards(path string) string {
	re := regexp.MustCompile(`\.\d+`)
	return re.ReplaceAllString(path, ".[*]")
}

// PrimToString converts a primitive value into a stable string representation.
func (ptc *ProviderTestCase) PrimToString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return ptc.NormalizeNumberString(x)
	case bool:
		return strconv.FormatBool(x)
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	default:
		panic(fmt.Sprintf("unsupported primitive type %T in PrimToString", v))
	}
}

// NormalizeNumberString returns canonical string for numeric-looking strings; otherwise returns input unchanged.
func (ptc *ProviderTestCase) NormalizeNumberString(s string) string {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return strconv.FormatInt(i, 10)
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
	return s
}

// ReadMapEntries extracts direct map entries from state under the given base path.
func (ptc *ProviderTestCase) ReadMapEntries(attrs map[string]string, base string) map[string]string {
	// Build prefix for map keys
	prefix := base + "."

	// Initialize output map
	out := map[string]string{}

	// Iterate all state attributes
	for k, v := range attrs {
		// Skip keys outside the prefix
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		// Skip meta counters
		if strings.HasSuffix(k, ".%") || strings.HasSuffix(k, ".#") {
			continue
		}

		// Extract remainder after the base prefix
		rest := strings.TrimPrefix(k, prefix)

		// Skip if the remainder indicates deeper nesting
		if strings.Contains(rest, ".") {
			continue
		}

		// Add direct entry
		out[rest] = v
	}
	return out
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
