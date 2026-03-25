package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/loov/watchrun/watchjs"
)

type compiled struct {
	base  *template.Template
	pages map[string]*template.Template
}

// Templates manages HTML template compilation and rendering.
// In dev mode, templates are loaded from disk and recompiled on changes.
type Templates struct {
	mu      sync.Mutex
	dev     bool
	dir     string
	embedFS embed.FS

	pagePaths []string
	compiled  *compiled
	err       error

	watchjs   *watchjs.Server
	errorPage *template.Template
}

func templatePage[T any](ts *Templates, name, path string) func(w http.ResponseWriter, data T) error {
	ts.pagePaths = append(ts.pagePaths, path)
	return func(w http.ResponseWriter, data T) error {
		c, err := ts.get()
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return c.pages[path].ExecuteTemplate(w, name, data)
	}
}

func templateFragment[T any](ts *Templates, name string) func(w http.ResponseWriter, data T) error {
	return func(w http.ResponseWriter, data T) error {
		c, err := ts.get()
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return c.base.ExecuteTemplate(w, name, data)
	}
}

// MustCompile compiles all registered templates.
// In production mode, panics on error to fail fast at startup.
func (t *Templates) MustCompile() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.compile()
	if t.err != nil && !t.dev {
		panic("failed to compile templates: " + t.err.Error())
	}
}

// WatchHandler returns the HTTP handler for serving watch.js, or nil in production.
func (t *Templates) WatchHandler() http.Handler {
	if t.watchjs == nil {
		return nil
	}
	return t.watchjs
}

// Recompile re-parses templates from disk and triggers a browser reload.
func (t *Templates) Recompile() {
	t.mu.Lock()
	t.compile()
	t.mu.Unlock()

	if t.watchjs != nil {
		t.watchjs.ReloadBrowser()
	}
}

func (t *Templates) get() (*compiled, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.compiled == nil && t.err == nil {
		t.compile()
	}
	return t.compiled, t.err
}

func (t *Templates) fsys() fs.FS {
	if t.dev {
		return os.DirFS(t.dir)
	}
	return t.embedFS
}

func (t *Templates) funcMap() template.FuncMap {
	return template.FuncMap{
		"devScripts": func() template.HTML {
			if t.dev {
				return `<script src="/~watch.js"></script>`
			}
			return ""
		},
	}
}

func (t *Templates) compile() {
	fsys := t.fsys()

	base, err := template.New("").Funcs(t.funcMap()).ParseFS(fsys,
		"templates/base.html",
		"templates/partials/*.html",
	)
	if err != nil {
		t.compiled = nil
		t.err = err
		return
	}

	pages := make(map[string]*template.Template, len(t.pagePaths))
	for _, path := range t.pagePaths {
		page, err := template.Must(base.Clone()).ParseFS(fsys, path)
		if err != nil {
			t.compiled = nil
			t.err = err
			return
		}
		pages[path] = page
	}

	t.compiled = &compiled{
		base:  base,
		pages: pages,
	}
	t.err = nil
}

// RenderError renders a template error page.
// In dev mode, includes watch.js so the browser auto-reloads when the error is fixed.
func (t *Templates) RenderError(w http.ResponseWriter, err error) {
	log.Println("template error:", err)
	if t.dev && t.errorPage != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.errorPage.Execute(w, err)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
