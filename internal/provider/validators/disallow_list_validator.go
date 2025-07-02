package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DisallowListValidator ensures a value is not one of a forbidden set
type DisallowListValidator struct {
	// Forbidden is the list of disallowed string values
	Forbidden []string
}

// Description returns a human-friendly description of this validator.
func (v DisallowListValidator) Description(ctx context.Context) string {
	// Build the description listing all forbidden values
	return fmt.Sprintf("Must not be any of %v", v.Forbidden)
}

// MarkdownDescription returns the markdown-formatted description.
func (v DisallowListValidator) MarkdownDescription(ctx context.Context) string {
	// Reuse the plain Description for markdown output
	return v.Description(ctx)
}

// ValidateString checks whether the given value matches any forbidden entry.
func (v DisallowListValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the configuration value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Extract the actual string value from the request
	val := req.ConfigValue.ValueString()

	// Iterate through each forbidden value
	for _, f := range v.Forbidden {
		// If the value matches a forbidden entry, record an error and stop
		if val == f {
			resp.Diagnostics.AddError(
				// Error summary title
				"Invalid Name",
				// Detailed error message naming the disallowed value
				fmt.Sprintf("cannot be %q", f),
			)
			return
		}
	}
}
