package cpln

type Tracing struct {
	Sampling   *float64                     `json:"sampling,omitempty"`
	Provider   *TracingProvider             `json:"provider,omitempty"`
	CustomTags *map[string]TracingCustomTag `json:"customTags,omitempty"`
}

type TracingProvider struct {
	Lightstep    *TracingProviderLightstep    `json:"lightstep,omitempty"`
	Otel         *TracingProviderOtel         `json:"otel,omitempty"`
	ControlPlane *TracingProviderControlPlane `json:"controlplane,omitempty"`
}

type TracingProviderLightstep struct {
	Endpoint    *string `json:"endpoint,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
}

type TracingProviderOtel struct {
	Endpoint *string `json:"endpoint,omitempty"`
}

type TracingProviderControlPlane struct{}

type TracingCustomTag struct {
	Literal *TracingCustomTagValue `json:"literal,omitempty"`
}

type TracingCustomTagValue struct {
	Value *string `json:"value,omitempty"`
}
