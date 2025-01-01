package journal

type TemplateData struct {
	PageTitle   string
	SystemUnits []Unit
	UserUnits   []Unit
}

func (td *TemplateData) Patterns() []string {
	return []string{
		"data/index.go.html",
		"data/journallist.go.html", // block: content
		"data/ui.css",              // block: ui
		"data/styles.css",          // block: styles
	}
}
