package middleware

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/GDG-KHU-Side/backend-side-project/httpx"
)

// RateLimiter는 클라이언트 IP별 토큰 버킷으로 요청 속도를 제한합니다(인메모리, 단일 인스턴스용).
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    rate.Limit
	burst    int
	idleTTL  time.Duration
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    rate.Limit(rps),
		burst:    burst,
		idleTTL:  3 * time.Minute,
	}
}

func (l *RateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.visitors[ip]
	if !ok {
		v = &visitor{limiter: rate.NewLimiter(l.limit, l.burst)}
		l.visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v.limiter.Allow()
}

func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.allow(clientIP(r)) {
			w.Header().Set("Retry-After", "1")
			httpx.WriteErrorMessage(w, http.StatusTooManyRequests, "요청이 너무 많습니다. 잠시 후 다시 시도해주세요.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Cleanup은 오래 요청이 없던 IP 항목을 주기적으로 지웁니다. ctx가 끝나면 멈춥니다.
func (l *RateLimiter) Cleanup(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			l.mu.Lock()
			for ip, v := range l.visitors {
				if time.Since(v.lastSeen) > l.idleTTL {
					delete(l.visitors, ip)
				}
			}
			l.mu.Unlock()
		}
	}
}

// clientIP는 TCP 연결의 원격 주소를 사용합니다. (프록시 뒤에 배포하면 X-Forwarded-For 처리를 추가해야 함)
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
