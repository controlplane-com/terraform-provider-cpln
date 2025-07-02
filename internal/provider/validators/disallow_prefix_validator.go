package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DisallowPrefixValidator holds the prefix that must not be used
type DisallowPrefixValidator struct {
	// Prefix is the forbidden starting substring
	Prefix string
}

// Description returns a human-friendly description of this validator.
func (v DisallowPrefixValidator) Description(ctx context.Context) string {
	// Build the description message with the forbidden prefix
	return fmt.Sprintf("Must not start with %q", v.Prefix)
}

// MarkdownDescription returns the markdown-formatted description.
func (v DisallowPrefixValidator) MarkdownDescription(ctx context.Context) string {
	// Reuse the plain Description for markdown output
	return v.Description(ctx)
}

// ValidateString checks whether the given value starts with the forbidden prefix.
func (v DisallowPrefixValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the configuration value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	// Extract the actual string value from the request
	val := req.ConfigValue.ValueString()

	// If the value starts with the forbidden prefix, record an error
	if strings.HasPrefix(val, v.Prefix) {
		// Report an error diagnostic with a summary and detail
		resp.Diagnostics.AddError(
			// Error summary title
			"Invalid Name",
			// Detailed error message explaining the violation
			fmt.Sprintf("cannot start with %q, got: %s", v.Prefix, val),
		)
	}
}
