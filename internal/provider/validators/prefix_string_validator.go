package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure implementation satisfies the validator.String interface
var _ validator.String = PrefixStringValidator{}

// PrefixStringValidator checks if a string starts with a given prefix or matches a regex pattern.
type PrefixStringValidator struct {
	// Prefix is the expected beginning of the string
	Prefix string

	// Regex is the compiled pattern used to validate the string prefix
	Regex *regexp.Regexp

	// Label is used in error messages to describe what kind of value is expected
	Label string
}

// NewPrefixStringValidator returns a new validator with the given prefix or regex.
func NewPrefixStringValidator(prefix, label string) PrefixStringValidator {
	// Create a new PrefixStringValidator with a compiled regex to match the exact prefix
	return PrefixStringValidator{
		Prefix: prefix,
		Regex:  regexp.MustCompile(fmt.Sprintf("^%s.*", regexp.QuoteMeta(prefix))),
		Label:  label,
	}
}

// Description returns a plain text description of the validator's behavior.
func (v PrefixStringValidator) Description(_ context.Context) string {
	// Provide a human-readable explanation for what this validator checks
	return fmt.Sprintf("Value must start with %q", v.Prefix)
}

// MarkdownDescription returns a markdown description of the validator's behavior.
func (v PrefixStringValidator) MarkdownDescription(_ context.Context) string {
	// Provide a markdown-formatted explanation for what this validator checks
	return fmt.Sprintf("The value must start with the prefix **%q**, typically used for %s.", v.Prefix, v.Label)
}

// ValidateString performs validation on the input string.
func (v PrefixStringValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Extract the actual string value
	val := req.ConfigValue.ValueString()

	// Check if the value does not match the expected prefix pattern
	if !v.Regex.MatchString(val) {
		// Add a diagnostic error indicating the value is invalid
		resp.Diagnostics.AddAttributeError(
			req.Path,
			fmt.Sprintf("Invalid %s", v.Label),
			fmt.Sprintf("The value must start with %q, got: %q", v.Prefix, val),
		)
	}
}
