package services

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
)

const (
	minPasswordLen = 8
	maxPasswordLen = 72 // bcrypt 한계
	maxNameLen     = 100
	minPartySize   = 1
	maxPartySize   = 20
)

// 휴대전화: 010/011/016/017/018/019, 하이픈/공백 허용
var mobileRe = regexp.MustCompile(`^01[016789][0-9]{7,8}$`)

// normalizePhone은 하이픈/공백을 제거하고 휴대전화 번호 형식을 검증합니다.
func normalizePhone(raw string) (string, error) {
	p := strings.NewReplacer("-", "", " ", "").Replace(strings.TrimSpace(raw))
	if !mobileRe.MatchString(p) {
		return "", apperr.Invalid("휴대전화 번호 형식이 올바르지 않습니다. (예: 010-1234-5678)")
	}
	return p, nil
}

func normalizeEmail(raw string) (string, error) {
	e := strings.ToLower(strings.TrimSpace(raw))
	addr, err := mail.ParseAddress(e)
	if err != nil || addr.Address != e {
		return "", apperr.Invalid("이메일 형식이 올바르지 않습니다.")
	}
	return e, nil
}

func validateName(raw string) (string, error) {
	n := strings.TrimSpace(raw)
	if n == "" || len([]rune(n)) > maxNameLen {
		return "", apperr.Invalid("이름은 1~100자여야 합니다.")
	}
	return n, nil
}

func validatePassword(p string) error {
	if len(p) < minPasswordLen || len(p) > maxPasswordLen {
		return apperr.Invalid("비밀번호는 8~72자여야 합니다.")
	}
	return nil
}

// 사장님 전화번호는 선택값이고, 매장 번호(02-...)일 수도 있어 형식만 느슨하게 확인합니다.
var phoneLooseRe = regexp.MustCompile(`^[0-9]{8,11}$`)

func normalizeOptionalPhone(raw string) (string, error) {
	p := strings.NewReplacer("-", "", " ", "").Replace(strings.TrimSpace(raw))
	if p == "" {
		return "", nil
	}
	if !phoneLooseRe.MatchString(p) {
		return "", apperr.Invalid("전화번호 형식이 올바르지 않습니다.")
	}
	return p, nil
}
