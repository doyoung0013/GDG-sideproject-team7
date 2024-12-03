package models

import (
	"time"
)

type UserRest struct {
	UserID    int64     `json:"user_id"`
	RestID    int64     `json:"rest_id"`
	CreatedAt time.Time `json:"created_at"`
}
