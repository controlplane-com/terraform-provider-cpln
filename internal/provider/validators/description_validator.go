package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DescriptionValidator ensures the string does not have leading/trailing whitespace and does not exceed 250 characters.
type DescriptionValidator struct{}

// Description provides a plain text description of the validator's behavior.
func (v DescriptionValidator) Description(ctx context.Context) string {
	return "Ensures the string does not contain leading or trailing whitespace and does not exceed 250 characters."
}

// MarkdownDescription provides a markdown description of the validator's behavior.
func (v DescriptionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation logic.
func (v DescriptionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	value := req.ConfigValue.ValueString()

	// Trim whitespace and compare to check for leading/trailing spaces
	trimmedValue := strings.TrimSpace(value)
	if value != trimmedValue {
		resp.Diagnostics.AddError(
			"Invalid Description",
			fmt.Sprintf("The description '%s' contains leading or trailing whitespace.", value),
		)
	}

	// Check length constraint
	if len(value) > 250 {
		resp.Diagnostics.AddError(
			"Description Too Long",
			fmt.Sprintf("The description '%s' exceeds the maximum length of 250 characters; got length: %d.", value, len(value)),
		)
	}
}
