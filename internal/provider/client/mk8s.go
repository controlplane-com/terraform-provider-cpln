package cpln

type Mk8s struct {
	Base
	Alias  *string     `json:"alias,omitempty"`
	Spec   *Mk8sSpec   `json:"spec,omitempty"`
	Status *Mk8sStatus `json:"status,omitempty"`
}

type Mk8sSpec struct {
}

type Mk8sStatus struct {
}
