package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProvisioningAuthNoTokenPassthrough(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	h := ProvisioningAuth("", next)

	req := httptest.NewRequest(http.MethodGet, "/api/provisioning/status", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called || rec.Code != http.StatusOK {
		t.Fatalf("expected passthrough, code=%d called=%v", rec.Code, called)
	}
}

func TestProvisioningAuthRejectsMissingHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next should not run")
	})
	h := ProvisioningAuth("secret", next)

	req := httptest.NewRequest(http.MethodGet, "/api/provisioning/status", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestProvisioningAuthBearerOK(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	h := ProvisioningAuth("mytoken", next)

	req := httptest.NewRequest(http.MethodGet, "/api/provisioning/status", nil)
	req.Header.Set("Authorization", "Bearer mytoken")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called || rec.Code != http.StatusOK {
		t.Fatalf("expected OK, code=%d called=%v", rec.Code, called)
	}
}

func TestProvisioningAuthCustomHeaderOK(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	h := ProvisioningAuth("x", next)

	req := httptest.NewRequest(http.MethodGet, "/api/provisioning/status", nil)
	req.Header.Set(HeaderProvisioningToken, "x")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called || rec.Code != http.StatusOK {
		t.Fatalf("expected OK")
	}
}

func TestProvisioningAuthNonProvisioningPath(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	h := ProvisioningAuth("secret", next)

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Fatal("non-provisioning path should skip auth")
	}
}
