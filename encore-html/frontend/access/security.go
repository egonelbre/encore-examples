package access

import (
	"net/http"
)

// ApplySecurityHeaders sets standard security headers on the response.
func ApplySecurityHeaders(w http.ResponseWriter, req *http.Request) {
	h := w.Header()

	h.Set("X-Frame-Options", "DENY")
	h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; frame-ancestors 'none'")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
	// h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	h.Set("Cross-Origin-Opener-Policy", "same-origin")
	h.Set("Cross-Origin-Resource-Policy", "same-origin")
}
