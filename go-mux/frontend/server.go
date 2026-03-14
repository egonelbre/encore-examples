package frontend

import (
	"net/http"
	"sync/atomic"
	"time"
)

// Server is the main HTTP server for the frontend.
//
//encore:service
type Server struct {
	router  *http.ServeMux
	message atomic.Value // stores *Message
}

func initServer() (*Server, error) {
	s := &Server{}
	s.message.Store(&Message{Text: "", ModifiedAt: time.Now()})
	s.router = http.NewServeMux()

	public := SecurityHeaders
	editor := RequireEditor
	admin := RequireAdmin

	s.router.Handle("GET /{$}", public(http.HandlerFunc(s.Home)))
	s.router.Handle("GET /edit", editor(http.HandlerFunc(s.Edit)))
	s.router.Handle("POST /edit", editor(http.HandlerFunc(s.EditPost)))
	s.router.Handle("GET /admin", admin(http.HandlerFunc(s.Admin)))
	s.router.Handle("GET /login", public(http.HandlerFunc(s.Login)))
	s.router.Handle("POST /login", public(http.HandlerFunc(s.LoginPost)))
	s.router.Handle("POST /logout", public(http.HandlerFunc(s.Logout)))
	s.router.Handle("GET /static/", public(http.HandlerFunc(Static)))

	return s, nil
}

// Handler serves all HTTP requests via the ServeMux router.
//
//encore:api public raw path=/!fallback
func (s *Server) Handler(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
