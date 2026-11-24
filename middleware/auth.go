package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/GDG-KHU-Side/backend-side-project/httpx"
)

type ctxKey struct{}

// TokenParser는 액세스 토큰을 검증하고 사용자 ID를 돌려줍니다.
type TokenParser interface {
	Parse(token string) (int64, error)
}

// RequireAuth는 "Authorization: Bearer <token>" 헤더를 검증하고 사용자 ID를 컨텍스트에 넣습니다.
func RequireAuth(tokens TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(h, "Bearer ")
			if !ok || token == "" {
				httpx.WriteErrorMessage(w, http.StatusUnauthorized, "로그인이 필요합니다.")
				return
			}
			userID, err := tokens.Parse(token)
			if err != nil {
				httpx.WriteErrorMessage(w, http.StatusUnauthorized, "인증이 만료되었거나 올바르지 않습니다. 다시 로그인해주세요.")
				return
			}
			next.ServeHTTP(w, r.WithContext(WithUserID(r.Context(), userID)))
		})
	}
}

func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// UserID는 RequireAuth가 넣어 둔 사용자 ID를 꺼냅니다.
func UserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ctxKey{}).(int64)
	return id, ok
}
