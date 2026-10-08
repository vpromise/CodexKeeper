package cpa_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"cpa-usage-keeper/internal/cpa"
)

func TestAuthFilePriorityClientSendsNameAndPriorityOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v0/management/auth-files/fields" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		want := map[string]any{"name": "auth.json", "priority": float64(0)}
		if !reflect.DeepEqual(body, want) {
			t.Fatalf("PATCH body = %#v, want %#v", body, want)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client := cpa.NewClient(server.URL, "management-secret", 2*time.Second, false)
	if status, err := client.UpdateAuthFilePriority(context.Background(), "auth.json", 0); err != nil || status != http.StatusOK {
		t.Fatalf("priority patch status=%d err=%v", status, err)
	}
}
