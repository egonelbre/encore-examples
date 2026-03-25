// Service frontend serves the HTMX frontend for the URL shortener.
package frontend

import (
	"net/http"

	"github.com/egonelbre/web-examples/htmx/url"
)

type Render struct {
	DashboardPage   func(w http.ResponseWriter, data *DashboardData) error
	UrlsPage        func(w http.ResponseWriter, data *url.ListResponse) error
	UrlListFragment func(w http.ResponseWriter, data *url.ListResponse) error
	UrlRowFragment  func(w http.ResponseWriter, data *url.URL) error
}

type Server struct {
	urls      *url.Service
	templates *Templates
	render     Render
	router    *http.ServeMux
}

func NewServer(service *url.Service, templates *Templates) *Server {
	server := &Server{urls: service, templates: templates}

	server.render = newRender(templates)
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

func newRender(templates *Templates) Render {
	return Render{
		DashboardPage:   templatePage[*DashboardData](templates, "base", "templates/dashboard.html"),
		UrlsPage:        templatePage[*url.ListResponse](templates, "base", "templates/urls.html"),
		UrlListFragment: templateFragment[*url.ListResponse](templates, "url-list-fragment"),
		UrlRowFragment:  templateFragment[*url.URL](templates, "url-row-fragment"),
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
	err := s.render.DashboardPage(w, dashboardData)
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

	err = s.render.UrlsPage(w, resp)
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

	s.render.UrlRowFragment(w, result)
}

// HtmxListURLs returns the URL list as HTML fragments.
func (s *Server) HtmxListURLs(w http.ResponseWriter, req *http.Request) {
	resp, err := s.urls.List(req.Context())
	if err != nil {
		http.Error(w, "Failed to list URLs", http.StatusInternalServerError)
		return
	}

	s.render.UrlListFragment(w, resp)
}
