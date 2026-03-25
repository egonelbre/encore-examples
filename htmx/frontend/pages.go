package frontend

import (
	"embed"
	"html/template"
	"log/slog"
	"path/filepath"

	"github.com/loov/watchrun/watch"
	"github.com/loov/watchrun/watchjs"
)

//go:embed all:templates
var templateFS embed.FS

// errorPage is parsed from the embedded FS so it works even when
// disk templates fail to compile. It includes watch.js so the
// browser auto-reloads when the error is fixed.
var errorPage = template.Must(template.New("error.html").ParseFS(templateFS, "templates/error.html"))

// NewTemplates creates templates compiled from the embedded filesystem.
func NewTemplates(log *slog.Logger) *Templates {
	return &Templates{
		log:       log,
		embedFS:   templateFS,
		errorPage: errorPage,
	}
}

// NewDevTemplates creates templates that reload from disk on file changes.
// The dir parameter is the path to the frontend package directory (e.g. "frontend").
func NewDevTemplates(log *slog.Logger, dir string) *Templates {
	t := &Templates{
		log:       log,
		dev:       true,
		dir:       dir,
		embedFS:   templateFS,
		errorPage: errorPage,
	}

	t.watchjs = watchjs.NewServer(watchjs.Config{
		Monitor: []string{
			filepath.Join(dir, "templates", "**"),
		},
		Ignore: watchjs.DefaultIgnore,
		OnChange: func(change watch.Change) (string, watchjs.Action) {
			if filepath.Ext(change.Path) == ".html" {
				t.Recompile()
			}
			return "", watchjs.IgnoreChanges
		},
	})

	return t
}
