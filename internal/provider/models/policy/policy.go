package policy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Binding //

type BindingModel struct {
	Permissions    types.Set `tfsdk:"permissions"`
	PrincipalLinks types.Set `tfsdk:"principal_links"`
}
