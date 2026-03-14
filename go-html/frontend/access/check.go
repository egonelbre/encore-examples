package access

import (
	"net/http"

	"encore.app/frontend/flash"
)

type Level uint32

const (
	LevelPublic Level = 0
	LevelEditor Level = 1 << iota
	LevelAdmin  Level = LevelEditor | (1 << iota)
)

var (
	Public = checker(LevelPublic)
	Editor = checker(LevelEditor)
	Admin  = checker(LevelAdmin)
)

// checker returns func that can be used to check a specific access level.
func checker(required Level) func(w http.ResponseWriter, req *http.Request) bool {
	return func(w http.ResponseWriter, req *http.Request) bool {
		return ApplyAndCheck(w, req, required)
	}
}

// ApplyAndCheck sets security headers and checks access. Redirects to /login if insufficient.
func ApplyAndCheck(w http.ResponseWriter, req *http.Request, required Level) bool {
	ApplySecurityHeaders(w, req)

	if required == LevelPublic {
		return true
	}

	user := CurrentUser(req)
	if user == nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return false
	}

	if user.Level&required != required {
		flash.Set(w, "You don't have permission to access that page.")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return false
	}

	return true
}
