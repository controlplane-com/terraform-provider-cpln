package helm

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PostrenderModel represents the postrender nested attribute.
type PostrenderModel struct {
	BinaryPath types.String `tfsdk:"binary_path"`
	Args       types.List   `tfsdk:"args"`
}
