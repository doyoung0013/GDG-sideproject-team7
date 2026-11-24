package models

import "time"

// 웨이팅 상태 (reservation.is_entry)
const (
	StatusWaiting  int8 = 0
	StatusEntered  int8 = 1
	StatusCanceled int8 = 2
)

type Reservation struct {
	ID               int64     `json:"id"`
	RestID           int64     `json:"rest_id"`
	Count            int       `json:"count"`
	CustomerPhoneNum string    `json:"customer_phone_num"`
	CreatedAt        time.Time `json:"created_at"`
	IsEntry          int8      `json:"is_entry"`
}

type CreateReservationRequest struct {
	RestID           int64  `json:"rest_id"`
	Count            int    `json:"count"`
	CustomerPhoneNum string `json:"customer_phone_num"`
}

type UpdateReservationStatusRequest struct {
	IsEntry *int8 `json:"is_entry"`
}
