package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"slices"
	"time"

	"github.com/GDG-KHU-Side/backend-side-project/httpx"
)

// Chain은 미들웨어를 바깥쪽부터 순서대로 적용합니다. Chain(h, a, b) == a(b(h))
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Logging은 요청마다 메서드/경로/상태/소요시간을 기록합니다. (요청 본문은 기록하지 않음)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %d %s %s", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond), clientIP(r))
	})
}

// Recover는 핸들러의 panic을 잡아 500으로 응답하고 서버가 죽지 않게 합니다.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}
				log.Printf("panic: %s %s: %v\n%s", r.Method, r.URL.Path, v, debug.Stack())
				httpx.WriteErrorMessage(w, http.StatusInternalServerError, "서버 내부 오류가 발생했습니다.")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS는 허용된 Origin에서 온 브라우저 요청을 허용합니다. "*" 는 전체 허용입니다.
func CORS(allowed []string) func(http.Handler) http.Handler {
	all := slices.Contains(allowed, "*")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && (all || slices.Contains(allowed, origin)) {
				h := w.Header()
				h.Set("Access-Control-Allow-Origin", origin)
				h.Add("Vary", "Origin")
				h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				h.Set("Access-Control-Max-Age", "600")
			}
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
