package whitespacestring

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringTypable = WhitespaceNormalizedStringType{}

// WhitespaceNormalizedStringType is a custom string type that treats values
// as semantically equal when they differ only by trailing whitespace.
// This prevents plan diffs caused by cosmetic differences in YAML formatting.
type WhitespaceNormalizedStringType struct {
	basetypes.StringType
}

// Equal returns true if the given type is equivalent
func (t WhitespaceNormalizedStringType) Equal(o attr.Type) bool {
	other, ok := o.(WhitespaceNormalizedStringType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns a human-readable string representation of the type
func (t WhitespaceNormalizedStringType) String() string {
	return "WhitespaceNormalizedStringType"
}

// ValueFromString returns a StringValuable type given a StringValue
func (t WhitespaceNormalizedStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := WhitespaceNormalizedStringValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value
func (t WhitespaceNormalizedStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// ValueType returns the Value type
func (t WhitespaceNormalizedStringType) ValueType(ctx context.Context) attr.Value {
	return WhitespaceNormalizedStringValue{}
}

// WhitespaceNormalizedStringValue is a custom string value that implements
// semantic equality by ignoring trailing whitespace differences.
type WhitespaceNormalizedStringValue struct {
	basetypes.StringValue
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ basetypes.StringValuable                      = WhitespaceNormalizedStringValue{}
	_ basetypes.StringValuableWithSemanticEquals    = WhitespaceNormalizedStringValue{}
)

// Type returns the custom type
func (v WhitespaceNormalizedStringValue) Type(ctx context.Context) attr.Type {
	return WhitespaceNormalizedStringType{}
}

// Equal returns true if the given value is equivalent
func (v WhitespaceNormalizedStringValue) Equal(o attr.Value) bool {
	other, ok := o.(WhitespaceNormalizedStringValue)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals returns true if the given value is semantically equal
// by comparing the values after trimming trailing whitespace.
// This allows YAML values with trailing newlines to be treated as equal to
// the same content without trailing newlines, preventing spurious diffs.
func (v WhitespaceNormalizedStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert the new value to our custom type
	newValue, ok := newValuable.(WhitespaceNormalizedStringValue)
	if !ok {
		return false, diags
	}

	// Get the old and new string values
	oldStr := v.StringValue.ValueString()
	newStr := newValue.StringValue.ValueString()

	// Normalize by trimming trailing whitespace (newlines, carriage returns, tabs, spaces)
	normalizedOld := strings.TrimRight(oldStr, "\n\r\t ")
	normalizedNew := strings.TrimRight(newStr, "\n\r\t ")

	// Compare the normalized values
	return normalizedOld == normalizedNew, diags
}
