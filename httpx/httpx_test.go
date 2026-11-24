package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
)

func TestWriteError(t *testing.T) {
	cases := []struct {
		err    error
		status int
		msg    string
	}{
		{apperr.Invalid("bad"), 400, "bad"},
		{apperr.Unauthorized("login"), 401, "login"},
		{apperr.Forbidden("no"), 403, "no"},
		{apperr.NotFound("missing"), 404, "missing"},
		{apperr.Conflict("dup"), 409, "dup"},
		{errors.New("sql: secret details"), 500, "서버 내부 오류가 발생했습니다."},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		WriteError(rec, httptest.NewRequest(http.MethodGet, "/x", nil), tc.err)
		var body ErrorResponse
		_ = json.NewDecoder(rec.Body).Decode(&body)
		if rec.Code != tc.status || body.Error != tc.msg {
			t.Errorf("%v: got %d %q", tc.err, rec.Code, body.Error)
		}
	}
}
