package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubParser struct{}

func (stubParser) Parse(tok string) (int64, error) {
	if tok == "good" {
		return 7, nil
	}
	return 0, errors.New("bad")
}

func TestRequireAuth(t *testing.T) {
	var got int64
	h := RequireAuth(stubParser{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = UserID(r.Context())
	}))

	for header, want := range map[string]int{"": 401, "Bearer bad": 401, "Basic good": 401, "Bearer good": 200} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if header != "" {
			req.Header.Set("Authorization", header)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Errorf("%q: got %d want %d", header, rec.Code, want)
		}
	}
	if got != 7 {
		t.Fatalf("user id not in context: %d", got)
	}
}

func TestRecover(t *testing.T) {
	h := Recover(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != 500 {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestRateLimit(t *testing.T) {
	l := NewRateLimiter(1, 2)
	h := l.Middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	codes := []int{}
	for range 3 {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		codes = append(codes, rec.Code)
	}
	if codes[0] != 200 || codes[1] != 200 || codes[2] != 429 {
		t.Fatalf("got %v", codes)
	}
}

func TestCORSPreflight(t *testing.T) {
	h := CORS([]string{"http://ok.test"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(418)
	}))
	req := httptest.NewRequest(http.MethodOptions, "/api/x", nil)
	req.Header.Set("Origin", "http://ok.test")
	req.Header.Set("Access-Control-Request-Method", "PUT")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 204 || rec.Header().Get("Access-Control-Allow-Origin") != "http://ok.test" {
		t.Fatalf("preflight: %d %v", rec.Code, rec.Header())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/x", nil)
	req.Header.Set("Origin", "http://evil.test")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("disallowed origin must not get CORS header")
	}
}
