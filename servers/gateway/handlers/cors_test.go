package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {
	ctx := InitializeFake()
	req, err := http.NewRequest("GET", "/v1/users/1", nil)
	if err != nil {
		t.Errorf("Error with request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SpecificUserHandler)
	wrap := NewCors(&CorsHandler{Handler: handler})
	wrap.ServeHTTP(rr, req)
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" ||
		rr.Header().Get("Access-Control-Allow-Methods") != "GET, PUT, POST, PATCH, DELETE" ||
		rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" ||
		rr.Header().Get("Access-Control-Expose-Headers") != "Authorization" ||
		rr.Header().Get("Access-Control-Max-Age") != "600" {
		t.Errorf("Incorrect header value")
	}
}
