package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// NameValidator ensures the string matches the required pattern and length constraints.
type NameValidator struct{}

// Description provides a plain text description of the validator's behavior.
func (v NameValidator) Description(ctx context.Context) string {
	return "Ensures the string matches the required pattern and length constraints."
}

// MarkdownDescription provides a markdown description of the validator's behavior.
func (v NameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation logic.
func (v NameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	value := req.ConfigValue.ValueString()

	// In case the value is empty, skip, probably the config is not ready yet
	if value == "" {
		return
	}

	// Define the regular expression pattern
	re := regexp.MustCompile(`^[a-z][-a-z0-9]([-a-z0-9])*[a-z0-9]$`)

	// Check if the value matches the pattern
	if !re.MatchString(value) {
		resp.Diagnostics.AddError(
			"Invalid Name",
			fmt.Sprintf("The value '%s' does not match the required pattern: %s", value, re.String()),
		)
	}

	// Check the length constraint
	if len(value) > 63 {
		resp.Diagnostics.AddError(
			"Name Too Long",
			fmt.Sprintf("The value '%s' exceeds the maximum length of 63 characters.", value),
		)
	}
}
