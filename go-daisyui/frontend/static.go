package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFS embed.FS

var staticFiles, _ = fs.Sub(staticFS, "static")
var staticHandler = http.StripPrefix("/static", http.FileServer(http.FS(staticFiles)))

// Static serves static files (daisyUI, htmx).
//
//encore:api public raw method=GET path=/static/*path
func Static(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	staticHandler.ServeHTTP(w, req)
}
