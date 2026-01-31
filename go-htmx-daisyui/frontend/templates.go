package frontend

import (
	"embed"
	"html/template"
	"io"

	"encore.app/url"
)

var (
	//go:embed all:templates
	templateFS embed.FS

	base = template.Must(template.ParseFS(templateFS,
		"templates/base.html",
		"templates/partials/*.html",
	))

	dashboardPage = templatePage[*DashboardData]("base", "templates/dashboard.html")
	urlsPage  = templatePage[*url.ListResponse]("base", "templates/urls.html")

	urlListFragment = templateFragment[*url.ListResponse]("url-list-fragment")
	urlRowFragment  = templateFragment[*url.URL]("url-row-fragment")
)

func parseTemplate(name string) *template.Template {
	clone := template.Must(base.Clone())
	return template.Must(clone.ParseFS(templateFS, name))
}

func templatePage[T any](name, path string) func(w io.Writer, data T) error {
	t := template.Must(base.Clone())
	t = template.Must(t.ParseFS(templateFS, path))
	return func(w io.Writer, data T) error {
		return t.ExecuteTemplate(w, name, data)
	}
}

func templateFragment[T any](name string) func(w io.Writer, data T) error {
	return func(w io.Writer, data T) error {
		return base.ExecuteTemplate(w, name, data)
	}
}
