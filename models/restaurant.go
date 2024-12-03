package models

import (
	"time"
)

type Restaurant struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Description string    `json:"description"`
	PhoneNum    string    `json:"phone_num"`
	WaitingTime int       `json:"waiting_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
