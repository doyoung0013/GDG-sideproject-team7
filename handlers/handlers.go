package handlers

import (
	"errors"
	"net/http"

	"github.com/GDG-KHU-Side/backend-side-project/middleware"
)

// currentUserID는 인증 미들웨어가 넣어 둔 사용자 ID를 꺼냅니다.
// RequireAuth 없이 라우팅된 경우를 잡기 위해 없으면 에러를 반환합니다.
func currentUserID(r *http.Request) (int64, error) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		return 0, errors.New("handler requires authentication middleware")
	}
	return id, nil
}
