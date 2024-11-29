package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// TagValidator ensures the map of tags does not exceed 50 entries.
type TagValidator struct{}

// Description provides a plain text description of the validator's behavior.
func (v TagValidator) Description(ctx context.Context) string {
	return "Ensures that the map does not contain more than 50 entries."
}

// MarkdownDescription provides a markdown description of the validator's behavior.
func (v TagValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap checks if the tag map contains more than 50 entries.
func (v TagValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	// Check the number of entries in the tag map
	if len(req.ConfigValue.Elements()) > 50 {
		resp.Diagnostics.AddError(
			"Too Many Tags",
			fmt.Sprintf("The %q map cannot contain more than 50 tags; got length: %d", req.Path, len(req.ConfigValue.Elements())),
		)
	}
}
