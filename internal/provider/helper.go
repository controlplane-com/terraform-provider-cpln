package cpln

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"slices"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	commonmodel "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/common"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

/*** Exported Functions ***/

// MergeAttributes combines multiple maps of schema.Attribute into a single map.
func MergeAttributes(maps ...map[string]schema.Attribute) map[string]schema.Attribute {
	merged := make(map[string]schema.Attribute)

	for _, m := range maps {
		for key, value := range m {
			merged[key] = value
		}
	}

	return merged
}

// GetNameFromSelfLink extracts the resource name from a selfLink string.
func GetNameFromSelfLink(selfLink string) string {
	// Split the selfLink by "/" separator
	parts := strings.Split(selfLink, "/")

	// Return the last element of the split parts as the name
	return parts[len(parts)-1]
}

// GetSelfLink construct the self link of the specified resource.
func GetSelfLink(orgName string, kind string, name string) string {
	return fmt.Sprintf("/org/%s/%s/%s", orgName, kind, name)
}

// GetSelfLinkWithGvc construct the self link of the specified resource.
func GetSelfLinkWithGvc(orgName string, kind string, gvc string, name string) string {
	return fmt.Sprintf("/org/%s/gvc/%s/%s/%s", orgName, gvc, kind, name)
}

// GetDomainLock returns a per-domain mutex for serializing route operations.
func GetDomainLock(domainName string) *sync.Mutex {
	mu, _ := domainOperationLocks.LoadOrStore(domainName, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

// DomainRouteKey returns a unique key for a DomainRoute based on its prefix or regex.
func DomainRouteKey(route client.DomainRoute) string {
	if route.Prefix != nil {
		return "prefix:" + *route.Prefix
	}

	if route.Regex != nil {
		return "regex:" + *route.Regex
	}

	return ""
}

// StringPointerFromInterface converts an interface{} to a *string.
func StringPointerFromInterface(input interface{}) *string {
	// Return nil if the input is nil
	if input == nil {
		return nil
	}

	// Assert that the input is a string
	strValue := input.(string)

	// Return pointer to the validated string value
	return &strValue
}

// BoolPointer returns a pointer to the input bool.
func BoolPointer(input bool) *bool {
	return &input
}

// StringPointer returns a pointer to the input string.
func StringPointer(input string) *string {
	return &input
}

// IntPointer returns a pointer to the input int.
func IntPointer(input int) *int {
	return &input
}

// Float64Pointer returns a pointer to the input float64.
func Float64Pointer(input float64) *float64 {
	return &input
}

// Float32Pointer returns a pointer to the input float32.
func Float32Pointer(input float32) *float32 {
	return &input
}

// IsGvcScopedResource returns true if the provided kind is scoped to GVC.
func IsGvcScopedResource(kind string) bool {
	return slices.Contains(GvcScopedKinds, kind)
}

// GetInterface returns a pointer to the provided interface value, or nil if the input is nil.
func GetInterface(s interface{}) *interface{} {
	// Check if the input is nil
	if s == nil {
		// Return nil for a nil input
		return nil
	}

	// Return pointer to the input interface value
	return &s
}

// ParseValueAndUnit extracts the numeric value and unit from a resource quantity string.
func ParseValueAndUnit(value string) (int, string) {
	// Compile regex for numeric characters
	numberRegex := regexp.MustCompile("[0-9]+")

	// Compile regex for alphabetic characters
	charactersRegex := regexp.MustCompile("[A-Za-z]+")

	// Find numeric substring and convert to integer
	number, _ := strconv.Atoi(numberRegex.FindString(value))

	// Find unit substring
	characters := charactersRegex.FindString(value)

	// Return the parsed number and characters
	return number, characters
}

// ToStringSlice converts a slice of interface{} to a slice of strings, returning an error if any element is not a string.
func ToStringSlice(ifaces []interface{}) []string {
	// Initialize a slice of strings with capacity matching the number of interface elements
	strs := make([]string, 0, len(ifaces))

	// Loop through each element in the input slice
	for _, v := range ifaces {
		// Add the asserted string to the result slice
		strs = append(strs, v.(string))
	}

	// Return the resulting slice and a nil error
	return strs
}

// PreserveJSONFormatting returns the plan value if both raw API string and plan string parse to semantically identical JSON, otherwise returns the API string value.
func PreserveJSONFormatting(raw interface{}, plan types.String) types.String {
	// Convert raw interface to string value
	rawAPI, ok := raw.(string)

	// Check if raw value is not a string
	if !ok {
		// Return original plan value
		return plan
	}

	// Extract raw plan JSON string
	rawPlan := plan.ValueString()

	// Declare variables to hold parsed JSON structures
	var apiObj, planObj interface{}

	// Attempt to unmarshal raw API JSON into apiObj
	if err := json.Unmarshal([]byte(rawAPI), &apiObj); err == nil {
		// Attempt to unmarshal raw plan JSON into planObj
		if err := json.Unmarshal([]byte(rawPlan), &planObj); err == nil {
			// Re-marshal apiObj to canonical JSON
			apiCanon, _ := json.Marshal(apiObj)

			// Re-marshal planObj to canonical JSON
			planCanon, _ := json.Marshal(planObj)

			// Compare canonical JSON for semantic equality
			if bytes.Equal(apiCanon, planCanon) {
				// Use plan value if JSON matches
				return plan
			}
		}
	}

	// Return API value preserving formatting
	return types.StringValue(rawAPI)
}

// StringSliceToString returns a string representation of the string slice.
func StringSliceToString(items []string) string {
	// Format the string slice into a quoted space-separated list using fmt
	q := fmt.Sprintf("%q", items)

	// Replace spaces between quoted items with comma and space to mimic JSON-like list
	return strings.ReplaceAll(q, `" "`, `", "`)
}

// IntSliceToString converts a slice of ints to a formatted string representation.
func IntSliceToString(nums []int) string {
	// Return empty slice representation if input is empty
	if len(nums) == 0 {
		return "[]"
	}

	// Create a string builder for efficient concatenation
	var sb strings.Builder

	// Write opening bracket
	sb.WriteString("[")

	// Iterate through numbers and append to builder
	for i, n := range nums {
		// Append number as string
		sb.WriteString(fmt.Sprintf("%d", n))
		// Append comma and space if not the last element
		if i < len(nums)-1 {
			sb.WriteString(", ")
		}
	}

	// Write closing bracket
	sb.WriteString("]")

	// Return the complete string
	return sb.String()
}

// IntSliceToStringSlice converts a slice of ints to a slice of strings.
func IntSliceToStringSlice(nums []int) []string {
	// Pre-allocate a string slice with the same length as the input slice
	result := make([]string, len(nums))

	// Iterate over each integer and convert it to a string
	for i, n := range nums {
		// Convert integer to string using strconv.Itoa
		result[i] = strconv.Itoa(n)
	}

	// Return the resulting slice of strings
	return result
}

// MapToHCL converts a map of key-value pairs to a formatted HCL block as a string.
func MapToHCL(dict map[string]interface{}, indentLevel int) string {
	// Define the default unit of indentation as two spaces
	unitIndent := "  "

	// Create the full base indentation string by repeating unitIndent
	indent := strings.Repeat(unitIndent, indentLevel)

	// Create a string builder to efficiently build the HCL string
	var b strings.Builder

	// Write the opening curly brace with indentation
	b.WriteString(indent + "{\n")

	// Initialize a slice to hold map keys
	keys := make([]string, 0, len(dict))

	// Iterate over the map to collect all keys
	for k := range dict {
		// Append each key to the slice
		keys = append(keys, k)
	}

	// Sort the keys alphabetically for consistent output
	sort.Strings(keys)

	// Define the indentation level for the map entries (one level deeper)
	entryIndent := indent + strings.Repeat(unitIndent, 1)

	// Iterate over the sorted keys to generate HCL entries
	for _, k := range keys {
		// Write each key-value pair in HCL format with proper indentation
		b.WriteString(fmt.Sprintf("%s%s = \"%v\"\n", entryIndent, k, dict[k]))
	}

	// Write the closing curly brace with base indentation
	b.WriteString(indent + "}")

	// Return the constructed HCL string
	return b.String()
}

// ConvertMapToStringMap converts a map with interface{} values to a map with string values.
func ConvertMapToStringMap(input map[string]interface{}) map[string]string {
	// Create a new map to hold stringified key-value pairs
	result := make(map[string]string)

	// Iterate over each key-value pair in the input map
	for k, v := range input {
		// Convert the value to a string using fmt.Sprint
		result[k] = fmt.Sprint(v)
	}

	// Return the resulting map with string values
	return result
}

// CanonicalizeEnvoyJSON parses the given Envoy JSON string into a generic Go data structure and then marshals it back to JSON.
func CanonicalizeEnvoyJSON(envoyStr string) string {
	// Declare a variable to hold the unmarshaled JSON structure
	var envoy interface{}

	// Unmarshal the input JSON string into the generic interface
	json.Unmarshal([]byte(envoyStr), &envoy)

	// Marshal the generic interface back into JSON bytes
	jsonOut, _ := json.Marshal(envoy)

	// Return the canonical JSON output
	return string(jsonOut)
}

// StringifyStringValue converts a types.String into a readable string representation
func StringifyStringValue(v types.String) string {
	// Return placeholder when the value is unknown
	if v.IsUnknown() {
		return "<unknown>"
	}

	// Return placeholder when the value is null
	if v.IsNull() {
		return "<null>"
	}

	// Return the actual string value when known and non-null
	return v.ValueString()
}

// ExecuteCplnCommand executes a cpln CLI command and returns the standard output.
func ExecuteCplnCommand(args []string) (string, error) {
	cmd := exec.Command("cpln", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errorMessage := ""
		stdoutString := stdout.String()
		stderrString := stderr.String()

		if len(strings.TrimSpace(stdoutString)) != 0 {
			errorMessage = fmt.Sprintf("Stdout: %s", stdoutString)
		}

		if len(strings.TrimSpace(stderrString)) != 0 {
			if len(errorMessage) != 0 {
				errorMessage = fmt.Sprintf("%s ", errorMessage)
			}
			errorMessage = fmt.Sprintf("%sStderr: %s", errorMessage, stderrString)
		}

		return "", fmt.Errorf("cpln command failed: %s. %s", err, errorMessage)
	}

	return stdout.String(), nil
}

// TestdataAbsPath returns the absolute path to a file under testdata.
func TestdataAbsPath(relativePath string) string {
	absPath, _ := filepath.Abs(relativePath)
	return absPath
}

// Builders //

// BuildString converts a types.String to a pointer to string.
func BuildString(input types.String) *string {
	// Determine if input has no valid value
	if input.IsNull() || input.IsUnknown() {
		// Indicate absence of value by returning nil
		return nil
	}

	// Return pointer to the underlying string value
	return input.ValueStringPointer()
}

// BuildInt converts a types.Int32 to a pointer to int.
func BuildInt(input types.Int32) *int {
	// Determine if input has no valid value
	if input.IsNull() || input.IsUnknown() {
		// Indicate absence of value by returning nil
		return nil
	}

	// Extract the Int32 value and cast to native int
	result := int(input.ValueInt32())

	// Provide pointer to the cast integer for further use
	return &result
}

// BuildFloat32 converts a types.Float64 to a pointer to float64.
func BuildFloat32(input types.Float32) *float32 {
	// Determine if input has no valid value
	if input.IsNull() || input.IsUnknown() {
		// Indicate absence of value by returning nil
		return nil
	}

	// Extract the Float32 value
	output := input.ValueFloat32()

	// Provide pointer to the float32 for further use
	return &output
}

// BuildFloat64 converts a types.Float64 to a pointer to float64.
func BuildFloat64(input types.Float64) *float64 {
	// Determine if input has no valid value
	if input.IsNull() || input.IsUnknown() {
		// Indicate absence of value by returning nil
		return nil
	}

	// Extract the Float64 value
	output := input.ValueFloat64()

	// Provide pointer to the float64 for further use
	return &output
}

// BuildBool converts a types.Bool to a pointer to bool.
func BuildBool(input types.Bool) *bool {
	// Determine if input has no valid value
	if input.IsNull() || input.IsUnknown() {
		// Indicate absence of value by returning nil
		return nil
	}

	// Return pointer to the underlying bool value
	return input.ValueBoolPointer()
}

// BuildTags filters and converts a types.Map to *map[string]interface{}.
func BuildTags(input types.Map) *map[string]interface{} {
	// Initialize the output map
	output := make(map[string]interface{})

	// Check if the input map is null or unknown
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	// Iterate over each key-value pair in the converted map
	for key, value := range input.Elements() {
		// Skip tags with prefixes that indicate server-generated metadata
		if shouldIgnoreTag(key) {
			continue
		}

		// Check if the types.String value is null or unknown; if so, set the output value as nil
		if value.IsNull() || value.IsUnknown() {
			output[key] = nil
		} else {
			// Otherwise, convert the types.String to a regular string and add it to the output map
			output[key] = formatFrameworkTypesToString(value)
		}
	}

	return &output
}

// BuildMapString converts a Terraform types.Map with tfsdk.StringType into a Go map[string]interface{} with nil entries for null or unknown values.
func BuildMapString(ctx context.Context, diags *diag.Diagnostics, input types.Map) *map[string]interface{} {
	// Exit early if map is null or unknown
	if input.IsNull() || input.IsUnknown() {
		// No data to process
		return nil
	}

	// Declare intermediate map to unmarshal Terraform values
	var intermediate map[string]types.String

	// Decode Terraform map into intermediate representation
	diags.Append(input.ElementsAs(ctx, &intermediate, false)...)

	// Abort on diagnostic errors
	if diags.HasError() {
		return nil
	}

	// Create output map with capacity for intermediate entries
	output := make(map[string]interface{}, len(intermediate))

	// Iterate over intermediate values
	for key, value := range intermediate {
		// Assign nil for null or unknown values
		if value.IsNull() || value.IsUnknown() {
			output[key] = nil
		} else {
			// Extract string value for known entries
			output[key] = formatFrameworkTypesToString(value)
		}
	}

	// Return populated map without error
	return &output
}

// BuildSetString converts a Terraform types.Set with tfsdk.StringType into a Go *[]string.
func BuildSetString(ctx context.Context, diags *diag.Diagnostics, input types.Set) *[]string {
	// Exit early if set itself is null or unknown
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	// Prepare an intermediate slice to unmarshal Terraform values
	var intermediate []types.String

	// Decode Terraform set elements into the intermediate slice
	diags.Append(input.ElementsAs(ctx, &intermediate, false)...)

	// Abort if any diagnostic errors occurred during decoding
	if diags.HasError() {
		return nil
	}

	// Build the output slice, preallocating for efficiency
	output := make([]string, 0, len(intermediate))

	// Iterate and extract each known string value
	for _, elem := range intermediate {
		// Skip null or unknown entries
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}

		// Add the element to the output slice
		output = append(output, elem.ValueString())
	}

	// Return a pointer to the populated slice
	return &output
}

// BuildListString converts a Terraform types.List with tfsdk.StringType into a Go *[]string.
func BuildListString(ctx context.Context, diags *diag.Diagnostics, input types.List) *[]string {
	// Exit early if list itself is null or unknown
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	// Prepare an intermediate slice to unmarshal Terraform values
	var intermediate []types.String

	// Decode Terraform list elements into the intermediate slice
	diags.Append(input.ElementsAs(ctx, &intermediate, false)...)

	// Abort if any diagnostic errors occurred during decoding
	if diags.HasError() {
		return nil
	}

	// Build the output slice, preallocating for efficiency
	output := make([]string, 0, len(intermediate))

	// Iterate and extract each known string value
	for _, elem := range intermediate {
		// Skip null or unknown entries
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}

		// Add the element to the output slice
		output = append(output, elem.ValueString())
	}

	// Return a pointer to the populated slice
	return &output
}

// BuildSetInt converts a Terraform types.Set with tfsdk.Int32Type into a Go *[]int.
func BuildSetInt(ctx context.Context, diags *diag.Diagnostics, input types.Set) *[]int {
	// Exit early if set itself is null or unknown
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	// Prepare an intermediate slice to unmarshal Terraform values
	var intermediate []types.Int32

	// Decode Terraform set elements into the intermediate slice
	diags.Append(input.ElementsAs(ctx, &intermediate, false)...)

	// Abort if any diagnostic errors occurred during decoding
	if diags.HasError() {
		return nil
	}

	// Build the output slice, preallocating for efficiency
	output := make([]int, 0, len(intermediate))

	// Iterate and extract each known string value
	for _, elem := range intermediate {
		// Skip null or unknown entries
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}

		// Add the element to the output slice
		output = append(output, int(elem.ValueInt32()))
	}

	// Return a pointer to the populated slice
	return &output
}

// BuildList extracts a slice of blocks of type T from a Terraform types.List.
func BuildList[T any](ctx context.Context, diags *diag.Diagnostics, l types.List) ([]T, bool) {
	// Return nil, false if list is null or unknown
	if l.IsNull() || l.IsUnknown() {
		return nil, false
	}

	// Prepare a slice to hold decoded blocks
	var blocks []T

	// Decode list elements into blocks and append any diagnostics
	diags.Append(l.ElementsAs(ctx, &blocks, false)...)

	// If decoding produced errors, abort and return false
	if diags.HasError() {
		return nil, false
	}

	// If there were no blocks, return nil
	if len(blocks) == 0 {
		return nil, false
	}

	// Return decoded blocks and success indicator
	return blocks, true
}

// BuildSet extracts a slice of blocks of type T from a Terraform types.Set.
func BuildSet[T any](ctx context.Context, diags *diag.Diagnostics, s types.Set) ([]T, bool) {
	// Return nil, false if set is null or unknown
	if s.IsNull() || s.IsUnknown() {
		return nil, false
	}

	// Prepare a slice to hold decoded blocks
	var blocks []T

	// Decode set elements into blocks and append any diagnostics
	diags.Append(s.ElementsAs(ctx, &blocks, false)...)

	// If decoding produced errors, abort and return false
	if diags.HasError() {
		return nil, false
	}

	// If there were no blocks, return nil
	if len(blocks) == 0 {
		return nil, false
	}

	// Return decoded blocks and success indicator
	return blocks, true
}

// BuildObject extracts a block of type T from a Terraform types.Object.
func BuildObject[T any](ctx context.Context, diags *diag.Diagnostics, o types.Object) (*T, bool) {
	// Return nil, false if object is null or unknown
	if o.IsNull() || o.IsUnknown() {
		return nil, false
	}

	// Prepare the destination value
	var block T

	// Decode the object into the destination using framework helper
	diags.Append(o.As(ctx, &block, basetypes.ObjectAsOptions{})...)

	// Abort on diagnostics errors
	if diags.HasError() {
		return nil, false
	}

	return &block, true
}

// Flatteners //

// FlattenInt converts an *int into a Terraform types.Int32.
func FlattenInt(input *int) types.Int32 {
	// Exit early on nil pointer to represent an absent value
	if input == nil {
		return types.Int32Null()
	}

	// Cast the Go int to int32 and return a concrete types.Int32 value
	return types.Int32Value(int32(*input))
}

// FlattenFloat64 converts a *float64 into a Terraform types.Float64.
func FlattenFloat64(input *float64) types.Float64 {
	// Exit early on nil pointer to represent an absent value
	if input == nil {
		return types.Float64Null()
	}

	// Return the constructed types.Float64
	return types.Float64Value(*input)
}

// FlattenSelfLink retrieves the "self" link from a list of client.Link objects and returns it as a types.String.
func FlattenSelfLink(links *[]client.Link) types.String {
	// Initialize selfLink as an empty string
	var selfLink string

	// Check if links is non-nil and contains elements
	if links != nil && len(*links) > 0 {
		// Iterate through each link in the slice to find the "self" link
		for _, ls := range *links {
			if ls.Rel == "self" {
				selfLink = ls.Href
				break // Stop searching once the "self" link is found
			}
		}
	}

	// Return the selfLink as a types.String, either with the found URL or as an empty types.String
	return types.StringValue(selfLink)
}

// FlattenTags converts a *map[string]interface{} to a types.Map.
func FlattenTags(input *map[string]interface{}) types.Map {
	// If the input map is nil, return an empty types.Map immediately
	if input == nil {
		return types.MapNull(types.StringType)
	}

	// Prepare elements for the types.Map
	elements := make(map[string]attr.Value)

	// Iterate over each key-value pair in the input map
	for key, value := range *input {
		// Skip tags that should be ignored
		if shouldIgnoreTag(key) {
			continue
		}

		elements[key] = formatToTypeString(value)
	}

	// Return the constructed types.Map
	return types.MapValueMust(types.StringType, elements)
}

// FlattenMapString converts a pointer to a Go map[string]interface{} into a Terraform types.Map.
func FlattenMapString(input *map[string]interface{}) types.Map {
	// Check if the input pointer is nil
	if input == nil {
		// Represent nil input as a null map
		return types.MapNull(types.StringType)
	}

	// Prepare elements for the types.Map
	elements := make(map[string]attr.Value)

	// Iterate over each key-value pair in the input map
	for key, value := range *input {
		elements[key] = formatToTypeString(value)
	}

	// Return the constructed types.Map
	return types.MapValueMust(types.StringType, elements)
}

// FlattenSetString converts a Go *[]string into a Terraform types.Set with tfsdk.StringType.
func FlattenSetString(input *[]string) types.Set {
	// No input slice means no data → return a null set
	if input == nil {
		return types.SetNull(types.StringType)
	}

	// Convert each Go string into an attr.Value (types.String)
	values := make([]attr.Value, len(*input))
	for i, s := range *input {
		values[i] = types.StringValue(s)
	}

	// Build and return the Set, panicking on any internal error
	return types.SetValueMust(types.StringType, values)
}

// FlattenListString converts a Go *[]string into a Terraform types.List with tfsdk.StringType.
func FlattenListString(input *[]string) types.List {
	// No input slice means no data → return a null list
	if input == nil {
		return types.ListNull(types.StringType)
	}

	// Convert each Go string into an attr.Value (types.String)
	values := make([]attr.Value, len(*input))
	for i, s := range *input {
		values[i] = types.StringValue(s)
	}

	// Build and return the List, panicking on any internal error
	return types.ListValueMust(types.StringType, values)
}

// FlattenSetInt converts a Go *[]int into a Terraform types.Set with tfsdk.Int32Type.
func FlattenSetInt(input *[]int) types.Set {
	// No input slice means no data → return a null set
	if input == nil {
		return types.SetNull(types.Int32Type)
	}

	// Convert each Go string into an attr.Value (types.String)
	values := make([]attr.Value, len(*input))
	for i, s := range *input {
		values[i] = types.Int32Value(int32(s))
	}

	// Build and return the Set, panicking on any internal error
	return types.SetValueMust(types.Int32Type, values)
}

// FlattenList creates a Terraform types.List from a slice of generic Model blocks.
func FlattenList[T commonmodel.Model](ctx context.Context, diags *diag.Diagnostics, blocks []T) types.List {
	// Declare a zero value to access attribute types
	var zero T

	// Obtain the element attribute types for the list
	elemType := zero.AttributeTypes()

	// Guard clause for existing diagnostics errors or empty input
	if diags.HasError() || len(blocks) == 0 {
		return types.ListNull(elemType)
	}

	// Convert the slice of blocks into a Terraform list while collecting diagnostics
	l, d := types.ListValueFrom(ctx, elemType, blocks)

	// Merge any diagnostics from the conversion into the main diagnostics
	diags.Append(d...)

	// If the conversion produced errors, return a null list
	if d.HasError() {
		return types.ListNull(elemType)
	}

	// Return the successfully built list
	return l
}

// FlattenSet creates a Terraform types.Set from a slice of generic Model blocks.
func FlattenSet[T commonmodel.Model](ctx context.Context, diags *diag.Diagnostics, blocks []T) types.Set {
	// Declare a zero value to access attribute types
	var zero T

	// Obtain the element attribute types for the set
	elemType := zero.AttributeTypes()

	// Guard clause for existing diagnostics errors or empty input
	if diags.HasError() || len(blocks) == 0 {
		return types.SetNull(elemType)
	}

	// Convert the slice of blocks into a Terraform set while collecting diagnostics
	l, d := types.SetValueFrom(ctx, elemType, blocks)

	// Merge any diagnostics from the conversion into the main diagnostics
	diags.Append(d...)

	// If the conversion produced errors, return a null set
	if d.HasError() {
		return types.SetNull(elemType)
	}

	// Return the successfully built set
	return l
}

// FlattenObject creates a Terraform types.Object from a block implementing Model.
func FlattenObject[T commonmodel.Model](ctx context.Context, diags *diag.Diagnostics, block *T) types.Object {
	// Access attribute types via a zero value of T
	var zero T

	// Attempt to interpret the attribute types as an ObjectType
	objectType := zero.AttributeTypes().(types.ObjectType)

	// Return a null object when the input block is absent
	if block == nil {
		return types.ObjectNull(objectType.AttrTypes)
	}

	// Build an Object value from the provided block while collecting diagnostics
	objectValue, diag := types.ObjectValueFrom(ctx, objectType.AttrTypes, *block)

	// Merge any diagnostics from the conversion into the main diagnostics
	diags.Append(diag...)

	// Return a null object if the conversion produced errors
	if diag.HasError() {
		return types.ObjectNull(objectType.AttrTypes)
	}

	// Return the successfully constructed Object value
	return objectValue
}

/*** Local Functions ***/

// shouldIgnoreTag checks if a tag key starts with any of the prefixes in IgnoredTagPrefixes.
func shouldIgnoreTag(key string) bool {
	// Iterate through each prefix in IgnoredTagPrefixes to check if the key starts with the prefix
	for _, prefix := range IgnoredTagPrefixes {
		// If the key starts with the current prefix, return true to ignore the tag
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}

	// If no prefixes match, return false to indicate the tag should not be ignored
	return false
}

// FormatTypeToString converts an interface{} value to its string representation.
func formatToTypeString(v interface{}) types.String {
	switch v := v.(type) {
	case nil:
		// Set the types.String to null if the value is nil
		return types.StringNull()
	case float64, float32:
		// Format floating-point numbers without decimal places
		return types.StringValue(fmt.Sprintf("%.0f", v))
	case int, int8, int16, int32, int64:
		// Format signed integers directly
		return types.StringValue(fmt.Sprintf("%d", v))
	case uint, uint8, uint16, uint32, uint64:
		// Format unsigned integers directly
		return types.StringValue(fmt.Sprintf("%d", v))
	case bool:
		// Format boolean values as "true" or "false"
		return types.StringValue(fmt.Sprintf("%t", v))
	case string:
		// Return string values as-is
		return types.StringValue(v)
	default:
		// Fallback for unsupported types
		return types.StringValue(fmt.Sprintf("%v", v))
	}
}

// formatFrameworkTypesToString converts attr.Value to its string representation.
func formatFrameworkTypesToString(value attr.Value) string {
	switch v := value.(type) {
	case types.Float32, types.Float64:
		// Format floating-point numbers without decimal places
		return fmt.Sprintf("%.0f", v)
	case types.Int32, types.Int64:
		// Format signed integers directly
		return fmt.Sprintf("%d", v)
	case types.Bool:
		// Format boolean values as "true" or "false"
		return fmt.Sprintf("%t", v.ValueBool())
	case types.String:
		return v.ValueString()
	default:
		// Fallback for unsupported types
		return fmt.Sprintf("%v", v)
	}
}
