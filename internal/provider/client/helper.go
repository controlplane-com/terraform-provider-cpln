package cpln

// GetLinkHref returns the Href of the first link matching the given rel attribute.
func GetLinkHref(links []Link, rel string) *string {
	// Loop through all provided links
	for _, link := range links {
		// Check if current link's Rel matches the target
		if link.Rel == rel {
			// Return pointer to the Href of the matching link
			return &link.Href
		}
	}

	// Return nil if no matching link is found
	return nil
}
