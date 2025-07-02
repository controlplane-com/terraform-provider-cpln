package common

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model //

type Model interface {
	AttributeTypes() attr.Type
}

// Query //

type QueryModel struct {
	Fetch types.String `tfsdk:"fetch"`
	Spec  types.List   `tfsdk:"spec"`
}

func (q QueryModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"fetch": types.StringType,
			"spec":  types.ListType{ElemType: QuerySpecModel{}.AttributeTypes()},
		},
	}
}

// Query -> Spec //

type QuerySpecModel struct {
	Match types.String `tfsdk:"match"`
	Terms types.List   `tfsdk:"terms"`
}

func (q QuerySpecModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"match": types.StringType,
			"terms": types.ListType{ElemType: QuerySpecTermModel{}.AttributeTypes()},
		},
	}
}

// Query -> Spec -> Terms //

type QuerySpecTermModel struct {
	Op       types.String `tfsdk:"op"`
	Property types.String `tfsdk:"property"`
	Rel      types.String `tfsdk:"rel"`
	Tag      types.String `tfsdk:"tag"`
	Value    types.String `tfsdk:"value"`
}

func (q QuerySpecTermModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"op":       types.StringType,
			"property": types.StringType,
			"rel":      types.StringType,
			"tag":      types.StringType,
			"value":    types.StringType,
		},
	}
}

// Lightstep Tracing //

type LightstepTracingModel struct {
	Sampling    types.Float64 `tfsdk:"sampling"`
	Endpoint    types.String  `tfsdk:"endpoint"`
	Credentials types.String  `tfsdk:"credentials"`
	CustomTags  types.Map     `tfsdk:"custom_tags"`
}

func (l LightstepTracingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sampling":    types.Float64Type,
			"endpoint":    types.StringType,
			"credentials": types.StringType,
			"custom_tags": types.MapType{ElemType: types.StringType},
		},
	}
}

// Otel Tracing //

type OtelTracingModel struct {
	Sampling   types.Float64 `tfsdk:"sampling"`
	Endpoint   types.String  `tfsdk:"endpoint"`
	CustomTags types.Map     `tfsdk:"custom_tags"`
}

func (o OtelTracingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sampling":    types.Float64Type,
			"endpoint":    types.StringType,
			"custom_tags": types.MapType{ElemType: types.StringType},
		},
	}
}

// Control Plane Tracing //

type ControlPlaneTracingModel struct {
	Sampling   types.Float64 `tfsdk:"sampling"`
	CustomTags types.Map     `tfsdk:"custom_tags"`
}

func (c ControlPlaneTracingModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sampling":    types.Float64Type,
			"custom_tags": types.MapType{ElemType: types.StringType},
		},
	}
}
