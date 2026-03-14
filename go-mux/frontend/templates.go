package frontend

import (
	"embed"
	"html/template"
	"net/http"

	"encore.app/frontend/access"
	"encore.app/frontend/flash"
)

var (
	//go:embed all:templates
	templateFS embed.FS

	base = template.Must(template.ParseFS(templateFS,
		"templates/base.html",
		"templates/partials/*.html",
	))

	homePage  = templatePage[*Message]("base", "templates/home.html")
	editPage  = templatePage[*Message]("base", "templates/edit.html")
	adminPage = templatePage[*Message]("base", "templates/admin.html")
	loginPage = templatePage[any]("base", "templates/login.html")
)

// PageData wraps page-specific data with flash messages and user info for base template rendering.
type PageData[T any] struct {
	Data     T
	Flash    string
	Username string
}

func templatePage[T any](name, path string) func(w http.ResponseWriter, req *http.Request, data T) error {
	t := template.Must(base.Clone())
	t = template.Must(t.ParseFS(templateFS, path))
	return func(w http.ResponseWriter, req *http.Request, data T) error {
		var username string
		if u := access.CurrentUser(req); u != nil {
			username = u.Username
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return t.ExecuteTemplate(w, name, &PageData[T]{
			Data:     data,
			Flash:    flash.Pop(w, req),
			Username: username,
		})
	}
}
