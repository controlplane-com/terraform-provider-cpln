package image

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Image //

type ImageModel struct {
	CplnId     types.String `tfsdk:"cpln_id"`
	Name       types.String `tfsdk:"name"`
	Tags       types.Map    `tfsdk:"tags"`
	SelfLink   types.String `tfsdk:"self_link"`
	Tag        types.String `tfsdk:"tag"`
	Repository types.String `tfsdk:"repository"`
	Digest     types.String `tfsdk:"digest"`
	Manifest   types.List   `tfsdk:"manifest"`
}

func (i ImageModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cpln_id":    types.StringType,
			"name":       types.StringType,
			"tags":       types.MapType{ElemType: types.StringType},
			"self_link":  types.StringType,
			"tag":        types.StringType,
			"repository": types.StringType,
			"digest":     types.StringType,
			"manifest":   types.ListType{ElemType: ManifestModel{}.AttributeTypes()},
		},
	}
}

// Manifest //

type ManifestModel struct {
	Config        types.List   `tfsdk:"config"`
	Layers        types.List   `tfsdk:"layers"`
	MediaType     types.String `tfsdk:"media_type"`
	SchemaVersion types.Int32  `tfsdk:"schema_version"`
}

func (m ManifestModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"config":         types.ListType{ElemType: ManifestConfigModel{}.AttributeTypes()},
			"layers":         types.ListType{ElemType: ManifestConfigModel{}.AttributeTypes()},
			"media_type":     types.StringType,
			"schema_version": types.Int32Type,
		},
	}
}

// Manifest -> Config //

type ManifestConfigModel struct {
	Size      types.Int32  `tfsdk:"size"`
	Digest    types.String `tfsdk:"digest"`
	MediaType types.String `tfsdk:"media_type"`
}

func (m ManifestConfigModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"size":       types.Int32Type,
			"digest":     types.StringType,
			"media_type": types.StringType,
		},
	}
}
