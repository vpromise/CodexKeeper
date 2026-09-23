package test

import (
	"strings"
	"testing"

	"cpa-usage-keeper/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func TestLoadBcryptLoginPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("shared-admin-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	path := writeAuthConfig(t, "LOGIN_PASSWORD_HASH='"+string(hash)+"'\n")
	cfg, err := config.Load(config.LoadOptions{EnvFile: path})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.AuthEnabled || cfg.LoginPassword != "" || cfg.LoginPasswordHash != string(hash) {
		t.Fatal("quoted bcrypt hash was not preserved as the sole login credential")
	}
}

func TestLoadRejectsInvalidLoginPasswordHash(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("shared-admin-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, value, want string }{
		{"malformed", "LOGIN_PASSWORD_HASH=invalid\n", "valid bcrypt hash"},
		{"truncated", "LOGIN_PASSWORD_HASH='" + string(hash[:59]) + "'\n", "valid bcrypt hash"},
		{"both", "LOGIN_PASSWORD=private-password\nLOGIN_PASSWORD_HASH='" + string(hash) + "'\n", "set only one"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := config.Load(config.LoadOptions{EnvFile: writeAuthConfig(t, tc.value)})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}
