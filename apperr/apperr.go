// Package apperr는 계층 간에 전달되는 "예상된" 에러 종류를 정의합니다.
// HTTP 계층은 Kind를 보고 상태 코드를 결정하고, Message를 그대로 클라이언트에 보여줍니다.
// apperr가 아닌 에러는 내부 오류(500)로 취급되며 메시지는 노출되지 않습니다.
package apperr

import "errors"

type Kind int

const (
	KindInvalid      Kind = iota + 1 // 400
	KindUnauthorized                 // 401
	KindForbidden                    // 403
	KindNotFound                     // 404
	KindConflict                     // 409
)

type Error struct {
	Kind    Kind
	Message string
}

func (e *Error) Error() string { return e.Message }

func Invalid(msg string) error      { return &Error{Kind: KindInvalid, Message: msg} }
func Unauthorized(msg string) error { return &Error{Kind: KindUnauthorized, Message: msg} }
func Forbidden(msg string) error    { return &Error{Kind: KindForbidden, Message: msg} }
func NotFound(msg string) error     { return &Error{Kind: KindNotFound, Message: msg} }
func Conflict(msg string) error     { return &Error{Kind: KindConflict, Message: msg} }

// As는 err가 apperr.Error이면 그것을 반환합니다.
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// Is는 err가 주어진 Kind의 apperr.Error인지 확인합니다.
func Is(err error, k Kind) bool {
	e, ok := As(err)
	return ok && e.Kind == k
}
