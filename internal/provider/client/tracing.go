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

// Tracing - Tracing
type Tracing struct {
	Sampling  *int              `json:"sampling,omitempty"`
	Lightstep *LightstepTracing `json:"lightstep,omitempty"`
}

func LightstepSchema() *schema.Schema {

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		// ExactlyOneOf: []string{"lightstep_tracing"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sampling": {
					Type:     schema.TypeInt,
					Required: true,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(int)
						if v < 0 || v > 100 {
							errs = append(errs, fmt.Errorf("%q must be between 0 and 100 inclusive, got: %d", key, v))
							return
						}

						return
					},
				},
				"endpoint": {
					Type:     schema.TypeString,
					Required: true,
				},
				"credentials": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}
