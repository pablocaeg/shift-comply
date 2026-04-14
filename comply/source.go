package comply

// Source is a legal citation for a scheduling regulation.
type Source struct {
	// Title is the name of the law, regulation, or standard.
	Title string `json:"title"`

	// Section identifies the specific section, article, or provision.
	Section string `json:"section,omitempty"`

	// URL links to the official text or authoritative source.
	URL string `json:"url,omitempty"`
}

// Citation returns a formatted string combining title and section.
func (s Source) Citation() string {
	if s.Section != "" {
		return s.Title + ", " + s.Section
	}
	return s.Title
}
