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
		//"data/pico.min.css", // block: style
		"data/ui.min.css",   // block: style
		"data/original.css", // block: theme
		"data/styles.css",   // block: custom-style
		//"data/ui.min.umd.cjs", // block: script
	}
}
