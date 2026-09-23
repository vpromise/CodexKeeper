package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	. "cpa-usage-keeper/internal/api"
	"cpa-usage-keeper/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthBcryptPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("shared-admin-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, hash, password string
		want                 int
	}{
		{"correct", string(hash), "shared-admin-password", http.StatusNoContent},
		{"wrong", string(hash), "wrong-password", http.StatusUnauthorized},
		{"empty", string(hash), "", http.StatusUnauthorized},
		{"hash-is-not-password", string(hash), string(hash), http.StatusUnauthorized},
		{"invalid-hash-no-fallback", "invalid", "old-password", http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := AuthConfig{Enabled: true, LoginPassword: "old-password", LoginPasswordHash: tc.hash}
			router := NewRouter(nil, nil, nil, nil, cfg, NewAuthHandler(cfg, auth.NewSessionManager(time.Hour)), "")
			body, err := json.Marshal(map[string]string{"password": tc.password})
			if err != nil {
				t.Fatal(err)
			}
			response := serveCredentialMutation(router, http.MethodPost, "/api/v1/auth/login", string(body))
			if response.Code != tc.want {
				t.Fatalf("login status = %d, want %d", response.Code, tc.want)
			}
			if tc.want == http.StatusNoContent {
				cookies := response.Result().Cookies()
				if len(cookies) != 1 || serveAPIGet(router, "/api/v1/usage/overview", cookies[0]).Code != http.StatusOK {
					t.Fatal("successful hash login did not unlock admin access")
				}
			} else if len(response.Result().Cookies()) != 0 {
				t.Fatal("failed login must not issue a session")
			}
		})
	}
}
