package cpln

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// LightstepTracing - LightstepTracing
type LightstepTracing struct {
	Endpoint    *string `json:"endpoint,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
}

type OtelTelemetry struct {
	Endpoint *string `json:"endpoint,omitempty"`
}

type ControlPlaneTracing struct{}

type Provider struct {
	Otel         *OtelTelemetry       `json:"otel,omitempty"`
	Lightstep    *LightstepTracing    `json:"lightstep,omitempty"`
	ControlPlane *ControlPlaneTracing `json:"controlplane,omitempty"`
}

type CustomTag struct {
	Literal *CustomTagValue `json:"literal,omitempty"`
}

type CustomTagValue struct {
	Value *string `json:"value,omitempty"`
}

// Tracing - Tracing
type Tracing struct {
	Sampling   *int                  `json:"sampling,omitempty"`
	Provider   *Provider             `json:"provider,omitempty"`
	CustomTags *map[string]CustomTag `json:"customTags,omitempty"`
}

var tracingOptions = []string{"lightstep_tracing", "otel_tracing", "controlplane_tracing"}

func CustomTagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Key-value map of custom tags.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func LightstepSchema(isExactlyOneOf bool) *schema.Schema {

	schema := schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sampling": {
					Type:         schema.TypeInt,
					Description:  "Determines what percentage of requests should be traced.",
					Required:     true,
					ValidateFunc: validateSamplingFunc,
				},
				"endpoint": {
					Type:        schema.TypeString,
					Description: "Tracing Endpoint Workload. Either the canonical endpoint or internal endpoint.",
					Required:    true,
				},
				"credentials": {
					Type:        schema.TypeString,
					Description: "Full link to referenced Opaque Secret.",
					Optional:    true,
				},
				"custom_tags": CustomTagsSchema(),
			},
		},
	}

	if isExactlyOneOf {
		schema.ExactlyOneOf = tracingOptions
	}

	return &schema
}

func OtelSchema(isExactlyOneOf bool) *schema.Schema {
	schema := schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sampling": {
					Type:         schema.TypeInt,
					Description:  "Determines what percentage of requests should be traced.",
					Required:     true,
					ValidateFunc: validateSamplingFunc,
				},
				"endpoint": {
					Type:        schema.TypeString,
					Description: "Tracing Endpoint Workload. Either the canonical endpoint or internal endpoint.",
					Required:    true,
				},
				"custom_tags": CustomTagsSchema(),
			},
		},
	}

	if isExactlyOneOf {
		schema.ExactlyOneOf = tracingOptions
	}

	return &schema
}

func ControlPlaneTracingSchema(isExactlyOneOf bool) *schema.Schema {
	schema := schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sampling": {
					Type:         schema.TypeInt,
					Description:  "Determines what percentage of requests should be traced.",
					Required:     true,
					ValidateFunc: validateSamplingFunc,
				},
				"custom_tags": CustomTagsSchema(),
			},
		},
	}

	if isExactlyOneOf {
		schema.ExactlyOneOf = tracingOptions
	}

	return &schema
}

func validateSamplingFunc(val interface{}, key string) (warns []string, errs []error) {
	v := val.(int)
	if v < 0 || v > 100 {
		errs = append(errs, fmt.Errorf("%q must be between 0 and 100 inclusive, got: %d", key, v))
		return
	}

	return
}
