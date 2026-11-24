package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // 절대 응답에 포함하지 않음
	PhoneNum     string    `json:"phone_num"`
	Name         string    `json:"name"`
	RestIDs      []int64   `json:"rest_ids"` // 연결된 식당 ID 목록
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PhoneNum string `json:"phone_num"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	PhoneNum string `json:"phone_num"`
	Name     string `json:"name"`
}
