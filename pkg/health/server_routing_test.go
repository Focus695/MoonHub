package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutingDecisionsHandler_invalidTier(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/routing/decisions?tier=not-a-tier&limit=10", nil)
	rec := httptest.NewRecorder()
	s.routingDecisionsHandler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["error"] == "" {
		t.Fatalf("expected error field in body: %v", body)
	}
}

func TestRoutingDecisionsHandler_validTier(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/routing/decisions?tier=simple&limit=5", nil)
	rec := httptest.NewRecorder()
	s.routingDecisionsHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want 200", rec.Code)
	}
}
