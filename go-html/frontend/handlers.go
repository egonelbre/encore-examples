package frontend

import (
	"net/http"
	"sync/atomic"
	"time"

	"encore.app/frontend/access"
	"encore.app/frontend/flash"
	"encore.app/frontend/form"
)

var message atomic.Value

func init() {
	message.Store(&Message{
		Text:       "",
		ModifiedAt: time.Now(),
	})
}

type Message struct {
	Text       string
	ModifiedAt time.Time
}

//encore:api public raw method=GET path=/!rest
func Home(w http.ResponseWriter, req *http.Request) {
	if !access.Public(w, req) {
		return
	}

	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	homePage(w, req, message.Load().(*Message))
}

//encore:api public raw method=GET path=/edit
func Edit(w http.ResponseWriter, req *http.Request) {
	if !access.Editor(w, req) {
		return
	}

	editPage(w, req, message.Load().(*Message))
}

//encore:api public raw method=POST path=/edit
func EditPost(w http.ResponseWriter, req *http.Request) {
	if !access.Editor(w, req) {
		return
	}

	f, err := form.Decode[struct {
		Text string `schema:"text,required"`
	}](req)
	if err != nil {
		flash.Set(w, err.Error())
		http.Redirect(w, req, "/edit", http.StatusSeeOther)
		return
	}

	message.Store(&Message{
		Text:       f.Text,
		ModifiedAt: time.Now(),
	})

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

//encore:api public raw method=GET path=/admin
func Admin(w http.ResponseWriter, req *http.Request) {
	if !access.Admin(w, req) {
		return
	}

	adminPage(w, req, message.Load().(*Message))
}
