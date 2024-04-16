package cpln

func GetLinkHref(links []Link, rel string) *string {

	for _, link := range links {
		if link.Rel == rel {
			return &link.Href
		}
	}

	return nil
}
