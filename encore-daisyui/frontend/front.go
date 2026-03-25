// Service frontend serves the HTMX frontend for the URL shortener.
package frontend

import (
	"net/http"

	"encore.app/url"
)

type Page struct {
	Icon     string
	Title    string
	Subtitle string
	URL      string
}

type DashboardData struct {
	Pages []Page
}

var dashboardData = &DashboardData{
	Pages: []Page{
		{Icon: "\U0001F517", Title: "URLs", Subtitle: "Manage shortened URLs", URL: "/urls"},
		{Icon: "\U0001F4C8", Title: "Analytics", Subtitle: "View click statistics", URL: "/urls"},
		{Icon: "\U0001F465", Title: "Users", Subtitle: "Manage user accounts", URL: "/urls"},
		{Icon: "\u2699\uFE0F", Title: "Settings", Subtitle: "Configure your app", URL: "/urls"},
		{Icon: "\U0001F511", Title: "API Keys", Subtitle: "Manage API access", URL: "/urls"},
		{Icon: "\U0001F4DA", Title: "Docs", Subtitle: "API documentation", URL: "/urls"},
		{Icon: "\U0001F514", Title: "Webhooks", Subtitle: "Event notifications", URL: "/urls"},
		{Icon: "\U0001F4E6", Title: "Integrations", Subtitle: "Third-party services", URL: "/urls"},
	},
}

// Serve serves the main page.
//
//encore:api public raw method=GET path=/!rest
func Serve(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	dashboardPage(w, dashboardData)
}

// URLs serves the main page.
//
//encore:api public raw method=GET path=/urls
func URLs(w http.ResponseWriter, req *http.Request) {
	resp, err := url.List(req.Context())
	if err != nil {
		http.Error(w, "Failed to list URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	urlsPage(w, resp)
}

// ShortenURL handles the form submission and returns an HTML fragment.
//
//encore:api public raw method=POST path=/url
func ShortenURL(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	rawURL := req.FormValue("url")
	if rawURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	_, err := url.Shorten(req.Context(), &url.ShortenParams{URL: rawURL})
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/urls", http.StatusSeeOther)
}
