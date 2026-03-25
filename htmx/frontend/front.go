// Service frontend serves the HTMX frontend for the URL shortener.
package frontend

import (
	"io"
	"net/http"

	"github.com/egonelbre/encore-example/go-pure/url"
)

type Pages struct {
	dashboardPage   func(w io.Writer, data *DashboardData) error
	urlsPage        func(w io.Writer, data *url.ListResponse) error
	urlListFragment func(w io.Writer, data *url.ListResponse) error
	urlRowFragment  func(w io.Writer, data *url.URL) error
}

type Server struct {
	urls      *url.Service
	templates *Templates
	render     Pages
	router    *http.ServeMux
}

func NewServer(service *url.Service, templates *Templates) *Server {
	server := &Server{urls: service, templates: templates}

	server.render = newPages(templates)
	templates.MustCompile()

	server.router = http.NewServeMux()

	server.router.HandleFunc("GET  /", server.Dashboard)
	server.router.HandleFunc("GET  /urls", server.URLs)
	server.router.HandleFunc("POST /htmx/shorten", server.HtmxShortenURL)
	server.router.HandleFunc("GET  /htmx/urls", server.HtmxListURLs)

	server.router.HandleFunc("GET  /static/", Static)

	if handler := templates.WatchHandler(); handler != nil {
		server.router.Handle("GET /~watch.js", handler)
	}

	return server
}

func newPages(templates *Templates) Pages {
	return Pages{
		dashboardPage:   templatePage[*DashboardData](templates, "base", "templates/dashboard.html"),
		urlsPage:        templatePage[*url.ListResponse](templates, "base", "templates/urls.html"),
		urlListFragment: templateFragment[*url.ListResponse](templates, "url-list-fragment"),
		urlRowFragment:  templateFragment[*url.URL](templates, "url-row-fragment"),
	}
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

// Dashboard serves the main page.
func (s *Server) Dashboard(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := s.render.dashboardPage(w, dashboardData)
	if err != nil {
		s.templates.RenderError(w, err)
	}
}

// URLs serves the URL management page.
func (s *Server) URLs(w http.ResponseWriter, req *http.Request) {
	resp, err := s.urls.List(req.Context())
	if err != nil {
		http.Error(w, "Failed to list URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = s.render.urlsPage(w, resp)
	if err != nil {
		s.templates.RenderError(w, err)
	}
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
	s.render.urlRowFragment(w, result)
}

// HtmxListURLs returns the URL list as HTML fragments.
func (s *Server) HtmxListURLs(w http.ResponseWriter, req *http.Request) {
	resp, err := s.urls.List(req.Context())
	if err != nil {
		http.Error(w, "Failed to list URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.render.urlListFragment(w, resp)
}
