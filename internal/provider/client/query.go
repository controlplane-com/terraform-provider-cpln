package cpln

// Query - Query
type Query struct {
	Kind    *string      `json:"kind,omitempty"`
	Context *interface{} `json:"context,omitempty"`
	Fetch   *string      `json:"fetch,omitempty"`
	Spec    *QuerySpec   `json:"spec,omitempty"`
}

// QuerySpec - QuerySpec
type QuerySpec struct {
	Match *string      `json:"match,omitempty"`
	Terms *[]QueryTerm `json:"terms,omitempty"`
}

// QueryTerm - QueryTerm
type QueryTerm struct {
	Op       *string `json:"op,omitempty"`
	Property *string `json:"property,omitempty"`
	Rel      *string `json:"rel,omitempty"`
	Tag      *string `json:"tag,omitempty"`
	Value    *string `json:"value,omitempty"`
}
