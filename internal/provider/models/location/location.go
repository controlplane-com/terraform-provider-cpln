package location

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Location //

type LocationModel struct {
	CplnId        types.String `tfsdk:"cpln_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Tags          types.Map    `tfsdk:"tags"`
	SelfLink      types.String `tfsdk:"self_link"`
	Origin        types.String `tfsdk:"origin"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Region        types.String `tfsdk:"region"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Geo           types.List   `tfsdk:"geo"`
	IpRanges      types.Set    `tfsdk:"ip_ranges"`
}

func (l LocationModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cpln_id":        types.StringType,
			"name":           types.StringType,
			"description":    types.StringType,
			"tags":           types.MapType{ElemType: types.StringType},
			"self_link":      types.StringType,
			"origin":         types.StringType,
			"cloud_provider": types.StringType,
			"region":         types.StringType,
			"enabled":        types.BoolType,
			"geo":            types.ListType{ElemType: GeoModel{}.AttributeTypes()},
			"ip_ranges":      types.SetType{ElemType: types.StringType},
		},
	}
}

// GeoModel //

type GeoModel struct {
	Lat       types.Float32 `tfsdk:"lat"`
	Lon       types.Float32 `tfsdk:"lon"`
	Country   types.String  `tfsdk:"country"`
	State     types.String  `tfsdk:"state"`
	City      types.String  `tfsdk:"city"`
	Continent types.String  `tfsdk:"continent"`
}

func (g GeoModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"lat":       types.Float64Type,
			"lon":       types.Float64Type,
			"country":   types.StringType,
			"state":     types.StringType,
			"city":      types.StringType,
			"continent": types.StringType,
		},
	}
}
