// Service frontend serves the HTMX frontend for the URL shortener.
package frontend

import (
	"net/http"

	"github.com/egonelbre/encore-example/go-pure/url"
)

type Server struct {
	urls   *url.Service
	router *http.ServeMux
}

func NewServer(service *url.Service) *Server {
	server := &Server{urls: service}
	server.router = http.NewServeMux()

	server.router.HandleFunc("GET  /", server.Dashboard)
	server.router.HandleFunc("GET  /urls", server.URLs)
	server.router.HandleFunc("POST /htmx/shorten", server.HtmxShortenURL)
	server.router.HandleFunc("GET  /htmx/urls", server.HtmxListURLs)

	server.router.HandleFunc("GET  /static/", Static)

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

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
func (s *Server) Dashboard(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	dashboardPage(w, dashboardData)
}

// URLs serves the main page.
func (s *Server) URLs(w http.ResponseWriter, req *http.Request) {
	resp, err := s.urls.List(req.Context())
	if err != nil {
		http.Error(w, "Failed to list URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	urlsPage(w, resp)
}

// HtmxShortenURL handles the form submission and returns an HTML fragment.
func (s *Server) HtmxShortenURL(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	rawURL := req.FormValue("url")
	if rawURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	result, err := s.urls.Shorten(req.Context(), &url.ShortenParams{URL: rawURL})
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	urlRowFragment(w, result)
}

// HtmxListURLs returns the URL list as HTML fragments.
func (s *Server) HtmxListURLs(w http.ResponseWriter, req *http.Request) {
	resp, err := s.urls.List(req.Context())
	if err != nil {
		http.Error(w, "Failed to list URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	urlListFragment(w, resp)
}
