package frontend

import (
	"net/http"

	"encore.app/frontend/access"
	"encore.app/frontend/flash"
	"encore.app/frontend/form"
)

//encore:api public raw method=GET path=/login
func Login(w http.ResponseWriter, req *http.Request) {
	access.ApplySecurityHeaders(w, req)

	if access.CurrentUser(req) != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	loginPage(w, req, nil)
}

//encore:api public raw method=POST path=/login
func LoginPost(w http.ResponseWriter, req *http.Request) {
	access.ApplySecurityHeaders(w, req)

	f, err := form.Decode[struct {
		Username string `schema:"username,required"`
		Password string `schema:"password,required"`
	}](req)
	if err != nil {
		flash.Set(w, err.Error())
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	token := access.Login(f.Username, f.Password)
	if token == "" {
		flash.Set(w, "Invalid username or password.")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     access.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

//encore:api public raw method=POST path=/logout
func Logout(w http.ResponseWriter, req *http.Request) {
	access.ApplySecurityHeaders(w, req)

	if c, err := req.Cookie(access.SessionCookieName); err == nil {
		access.Logout(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   access.SessionCookieName,
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
