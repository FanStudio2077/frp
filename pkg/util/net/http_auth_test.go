package net

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPAuthRoles(t *testing.T) {
	mw := NewHTTPAuthMiddleware("admin", "pw").
		SetUsers(map[string]string{"viewer": "pw2"}, map[string]string{"viewer": "viewer"}).
		Middleware
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(RoleFromContext(r.Context())))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("admin", "pw")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 || rec.Body.String() != "admin" {
		t.Fatalf("admin: code=%d body=%s", rec.Code, rec.Body.String())
	}

	req2 := httptest.NewRequest("GET", "/", nil)
	req2.SetBasicAuth("viewer", "pw2")
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	if rec2.Code != 200 || rec2.Body.String() != "viewer" {
		t.Fatalf("viewer: code=%d body=%s", rec2.Code, rec2.Body.String())
	}

	req3 := httptest.NewRequest("GET", "/", nil)
	req3.SetBasicAuth("x", "y")
	rec3 := httptest.NewRecorder()
	h.ServeHTTP(rec3, req3)
	if rec3.Code != 401 {
		t.Fatalf("bad creds expected 401, got %d", rec3.Code)
	}
}
