package services

import (
	"context"
	"testing"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

func TestRegisterAndLogin(t *testing.T) {
	users := newFakeUsers()
	svc := NewUserService(users, fakeTokens{})
	ctx := context.Background()

	u, err := svc.Register(ctx, models.RegisterRequest{Email: " Owner@Example.com ", Password: "password123", Name: "사장님"})
	if err != nil {
		t.Fatal(err)
	}
	if u.Email != "owner@example.com" {
		t.Fatalf("email not normalized: %q", u.Email)
	}
	if users.byID[u.ID].PasswordHash == "password123" {
		t.Fatal("password stored in plain text")
	}

	if _, err := svc.Register(ctx, models.RegisterRequest{Email: "owner@example.com", Password: "password123", Name: "x"}); !apperr.Is(err, apperr.KindConflict) {
		t.Fatalf("duplicate email: want conflict, got %v", err)
	}
	if _, err := svc.Register(ctx, models.RegisterRequest{Email: "a@b.com", Password: "short", Name: "x"}); !apperr.Is(err, apperr.KindInvalid) {
		t.Fatalf("short password: want invalid, got %v", err)
	}

	// 식당이 연결되지 않은 사용자도 로그인 가능해야 함 (기존 NULL rest_id 버그)
	res, err := svc.Login(ctx, models.LoginRequest{Email: "owner@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("login without restaurant failed: %v", err)
	}
	if res.Token == "" || len(res.User.RestIDs) != 0 {
		t.Fatalf("unexpected login response: %+v", res)
	}

	if _, err := svc.Login(ctx, models.LoginRequest{Email: "owner@example.com", Password: "wrong-password"}); !apperr.Is(err, apperr.KindUnauthorized) {
		t.Fatalf("wrong password: want unauthorized, got %v", err)
	}
	if _, err := svc.Login(ctx, models.LoginRequest{Email: "nobody@example.com", Password: "password123"}); !apperr.Is(err, apperr.KindUnauthorized) {
		t.Fatalf("unknown email: want unauthorized, got %v", err)
	}

	// 연결 후 로그인하면 rest_ids 포함
	if _, err := svc.LinkRestaurant(ctx, u.ID, 7); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.LinkRestaurant(ctx, 999, 7); !apperr.Is(err, apperr.KindConflict) {
		t.Fatalf("second owner: want conflict, got %v", err)
	}
	res, _ = svc.Login(ctx, models.LoginRequest{Email: "owner@example.com", Password: "password123"})
	if len(res.User.RestIDs) != 1 || res.User.RestIDs[0] != 7 {
		t.Fatalf("rest_ids: %v", res.User.RestIDs)
	}
}
