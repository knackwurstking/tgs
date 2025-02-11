package opmanga

type TemplateData struct {
	PageTitle string
	Arcs      []Arc
}

func (td *TemplateData) Patterns() []string {
	return []string{
		"data/index.go.html",
		"data/opmangalist.go.html", // block: content
		"data/ui.css",              // block: ui
		"data/styles.css",          // block: styles
	}
}
