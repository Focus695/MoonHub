package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

const provisioningAPIPrefix = "/api/provisioning/"

// HeaderProvisioningToken is an alternate header for environments that strip Authorization.
const HeaderProvisioningToken = "X-Moonhub-Provisioning-Token"

// ProvisioningAuth enforces a shared secret on /api/provisioning/* when token is non-empty.
func ProvisioningAuth(token string, next http.Handler) http.Handler {
	if token == "" {
		return next
	}
	want := []byte(token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, provisioningAPIPrefix) {
			next.ServeHTTP(w, r)
			return
		}
		got := provisioningTokenFromRequest(r)
		if len(got) != len(want) || subtle.ConstantTimeCompare(got, want) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"provisioning authentication required"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func provisioningTokenFromRequest(r *http.Request) []byte {
	if h := r.Header.Get(HeaderProvisioningToken); h != "" {
		return []byte(strings.TrimSpace(h))
	}
	const prefix = "Bearer "
	auth := r.Header.Get("Authorization")
	if len(auth) > len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		return []byte(strings.TrimSpace(auth[len(prefix):]))
	}
	return nil
}
