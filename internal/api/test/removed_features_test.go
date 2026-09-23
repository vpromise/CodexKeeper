package test

import (
	. "cpa-usage-keeper/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemovedFeaturesHaveNoAPIRoutes(t *testing.T) {
	router := NewRouter(nil, nil, nil, nil, AuthConfig{}, nil, "")
	for _, path := range []string{
		"/api/v1/auth/api-key-login", "/api/v1/key-overview", "/api/v1/key-analysis",
		"/api/v1/key-activity", "/api/v1/ranking", "/api/v1/local-ranking",
	} {
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(method, path, nil))
			if response.Code != http.StatusNotFound {
				t.Errorf("%s %s: got %d", method, path, response.Code)
			}
		}
	}
}
