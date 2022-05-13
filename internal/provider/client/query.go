package cpln

// Query - Query
type Query struct {
	Kind    *string      `json:"kind,omitempty"`
	Context *interface{} `json:"context,omitempty"`
	Fetch   *string      `json:"fetch,omitempty"`
	Spec    *Spec        `json:"spec,omitempty"`
}

// Spec - Spec
type Spec struct {
	Match *string `json:"match,omitempty"`
	Terms *[]Term `json:"terms,omitempty"`
}

// Term - Term
type Term struct {
	Op       *string `json:"op,omitempty"`
	Property *string `json:"property,omitempty"`
	Rel      *string `json:"rel,omitempty"`
	Tag      *string `json:"tag,omitempty"`
	Value    *string `json:"value,omitempty"`
}
