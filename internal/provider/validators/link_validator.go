package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// LinkValidator ensures the string matches the required pattern and length constraints.
type LinkValidator struct{}

// Description provides a plain text description of the validator's behavior.
func (v LinkValidator) Description(ctx context.Context) string {
	return "Ensures the string matches the required pattern and length constraints."
}

// MarkdownDescription provides a markdown description of the validator's behavior.
func (v LinkValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation logic.
func (v LinkValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	value := req.ConfigValue.ValueString()

	// In case the value is empty, skip, probably the config is not ready yet
	if value == "" {
		return
	}

	// Define the regular expression pattern
	re := regexp.MustCompile(`(\/org\/[^/]+\/.*)|(\/\/.+)`)

	// Check if the value matches the pattern
	if !re.MatchString(value) {
		resp.Diagnostics.AddError(
			"Invalid Link",
			fmt.Sprintf("The value '%s' for '%s' attribute does not match the required pattern: %s", value, req.Path.String(), re.String()),
		)
	}
}
