package frontend

import (
	"net/http"
	"time"

	"encore.app/frontend/flash"
	"encore.app/frontend/form"
)

// Message represents the current message on the board.
type Message struct {
	Text       string
	ModifiedAt time.Time
}

// Home renders the home page.
func (s *Server) Home(w http.ResponseWriter, req *http.Request) {
	homePage(w, req, s.message.Load().(*Message))
}

// Edit renders the edit page.
func (s *Server) Edit(w http.ResponseWriter, req *http.Request) {
	editPage(w, req, s.message.Load().(*Message))
}

// EditPost handles the edit form submission.
func (s *Server) EditPost(w http.ResponseWriter, req *http.Request) {
	f, err := form.Decode[struct {
		Text string `schema:"text,required"`
	}](req)
	if err != nil {
		flash.Set(w, err.Error())
		http.Redirect(w, req, "/edit", http.StatusSeeOther)
		return
	}

	s.message.Store(&Message{
		Text:       f.Text,
		ModifiedAt: time.Now(),
	})

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// Admin renders the admin page.
func (s *Server) Admin(w http.ResponseWriter, req *http.Request) {
	adminPage(w, req, s.message.Load().(*Message))
}
