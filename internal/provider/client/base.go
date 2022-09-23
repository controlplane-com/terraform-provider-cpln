package cpln

// Base - Control Plane Base Struct
type Base struct {
	ID           *string                 `json:"id,omitempty"`
	Name         *string                 `json:"name,omitempty"`
	Kind         *string                 `json:"kind,omitempty"`
	Version      *int                    `json:"version,omitempty"`
	Description  *string                 `json:"description,omitempty"`
	Tags         *map[string]interface{} `json:"tags,omitempty"`
	Created      *string                 `json:"created,omitempty"`
	LastModified *string                 `json:"lastModified,omitempty"`
	Links        *[]Link                 `json:"links,omitempty"`
}

// Link - Link
type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}
