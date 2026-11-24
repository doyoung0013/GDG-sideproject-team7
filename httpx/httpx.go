// Package httpx는 JSON 응답/에러 형식을 통일하는 헬퍼를 제공합니다.
// 에러 응답은 항상 {"error": "메시지"} 형태입니다.
package httpx

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
)

const maxBodyBytes = 1 << 20 // 1MB

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func WriteErrorMessage(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}

// WriteError는 apperr 종류에 맞는 상태 코드로 응답합니다.
// 알 수 없는 에러는 로그로만 남기고 클라이언트에는 일반 메시지를 보냅니다.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := apperr.As(err); ok {
		WriteErrorMessage(w, statusFor(e.Kind), e.Message)
		return
	}
	log.Printf("internal error: %s %s: %v", r.Method, r.URL.Path, err)
	WriteErrorMessage(w, http.StatusInternalServerError, "서버 내부 오류가 발생했습니다.")
}

func statusFor(k apperr.Kind) int {
	switch k {
	case apperr.KindInvalid:
		return http.StatusBadRequest
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// DecodeJSON은 요청 본문을 v로 디코딩합니다(최대 1MB).
func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return apperr.Invalid("요청 본문이 올바른 JSON 형식이 아닙니다.")
	}
	return nil
}

// PathID는 URL 경로 변수를 양의 정수 ID로 파싱합니다.
func PathID(r *http.Request, name string) (int64, error) {
	id, err := strconv.ParseInt(mux.Vars(r)[name], 10, 64)
	if err != nil || id <= 0 {
		return 0, apperr.Invalid("잘못된 ID입니다.")
	}
	return id, nil
}
