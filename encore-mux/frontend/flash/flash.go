package flash

import "net/http"

const cookieName = "flash"

// Set sets a flash message cookie that will be shown on the next page load.
func Set(w http.ResponseWriter, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    msg,
		Path:     "/",
		MaxAge:   10,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// Pop reads and clears the flash message cookie. Returns empty string if none.
func Pop(w http.ResponseWriter, req *http.Request) string {
	c, err := req.Cookie(cookieName)
	if err != nil {
		return ""
	}
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Path:   "/",
		MaxAge: -1,
	})
	return c.Value
}
