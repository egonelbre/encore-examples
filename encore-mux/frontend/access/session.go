package access

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"encore.dev/storage/cache"
)

const SessionCookieName = "session"

// User represents a logged-in user.
type User struct {
	Username string
	Level    Level
}

// users is the hardcoded user database.
var users = map[string]struct {
	Password string
	Level    Level
}{
	"admin":  {Password: "admin", Level: LevelAdmin},
	"editor": {Password: "editor", Level: LevelEditor},
}

var sessionCluster = cache.NewCluster("sessions", cache.ClusterConfig{
	EvictionPolicy: cache.AllKeysLRU,
})

var sessionStore = cache.NewStructKeyspace[string, User](sessionCluster, cache.KeyspaceConfig{
	KeyPattern:    "session/:key",
	DefaultExpiry: cache.ExpireIn(24 * time.Hour),
})

// Login validates credentials and creates a session. Returns the session token or empty string.
func Login(username, password string) string {
	u, ok := users[username]
	if !ok || u.Password != password {
		return ""
	}

	token := newToken()
	user := User{Username: username, Level: u.Level}
	if err := sessionStore.Set(context.Background(), token, user); err != nil {
		return ""
	}
	return token
}

// Logout removes a session.
func Logout(token string) {
	sessionStore.Delete(context.Background(), token)
}

// CurrentUser returns the user for the request, or nil if not logged in.
func CurrentUser(req *http.Request) *User {
	c, err := req.Cookie(SessionCookieName)
	if err != nil {
		return nil
	}
	u, err := sessionStore.Get(req.Context(), c.Value)
	if errors.Is(err, cache.Miss) || err != nil {
		return nil
	}
	return &u
}

func newToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
