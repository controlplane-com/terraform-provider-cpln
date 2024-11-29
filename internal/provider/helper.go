package cpln

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Exported Functions ***/

// MergeAttributes combines multiple maps of schema.Attribute into a single map.
// Later maps in the arguments override keys from earlier ones.
func MergeAttributes(maps ...map[string]schema.Attribute) map[string]schema.Attribute {
	merged := make(map[string]schema.Attribute)

	for _, m := range maps {
		for key, value := range m {
			merged[key] = value
		}
	}

	return merged
}

// BUILDS

// BuildString converts a types.String to a *string.
// If the types.String value is null or unknown, the function returns nil.
// Otherwise, it returns a pointer to the underlying string value.
func BuildString(input types.String) *string {
	// Check if the types.String input is null or unknown
	// Return nil to represent the absence of a input
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	return input.ValueStringPointer()
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

// SCHEMAS

// SchemaSelfLink retrieves the "self" link from a list of client.Link objects and returns it as a types.String.
// If no "self" link is found, it returns an empty types.String.
func SchemaSelfLink(links *[]client.Link) types.String {
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

// SchemaTags converts a *map[string]interface{} to a types.Map.
// If the input map is nil, it returns an empty types.Map with StringType elements.
// Each interface{} value is converted to types.String, with non-string values being
// formatted as strings, and nil values set as types.StringNull.
func SchemaTags(input *map[string]interface{}) types.Map {
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

/*** Private Functions ***/

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
// This function ensures consistent string formatting across various types,
// including numeric, boolean, and string values.
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
// Handles booleans, numbers, and other types, converting them into strings.
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
