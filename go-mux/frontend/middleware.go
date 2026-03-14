package frontend

import (
	"net/http"

	"encore.app/frontend/access"
	"encore.app/frontend/flash"
)

// SecurityHeaders sets standard security headers on the response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access.ApplySecurityHeaders(w, r)
		next.ServeHTTP(w, r)
	})
}

// RequireEditor requires editor-level access.
func RequireEditor(next http.Handler) http.Handler {
	return SecurityHeaders(requireLevel(access.LevelEditor, next))
}

// RequireAdmin requires admin-level access.
func RequireAdmin(next http.Handler) http.Handler {
	return SecurityHeaders(requireLevel(access.LevelAdmin, next))
}

func requireLevel(required access.Level, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := access.CurrentUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if user.Level&required != required {
			flash.Set(w, "You don't have permission to access that page.")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
