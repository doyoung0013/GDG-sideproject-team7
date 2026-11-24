package models

import "time"

type Restaurant struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Description  string    `json:"description"`
	PhoneNum     string    `json:"phone_num"`
	WaitingTime  int       `json:"waiting_time"`
	WaitingCount int       `json:"waiting_count"` // 현재 대기 중인 팀 수
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
