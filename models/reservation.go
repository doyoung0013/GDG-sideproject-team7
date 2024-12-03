package models

import (
	"time"
)

type Reservation struct {
	ID               int64     `json:"id"`
	RestID           int64     `json:"rest_id"`
	Count            int       `json:"count"`
	CustomerPhoneNum string    `json:"customer_phone_num"`
	CreatedAt        time.Time `json:"created_at"`
	IsEntry          int8      `json:"is_entry"`
}

/*
func (r *Reservation) Create() error {
	query := `
        INSERT INTO reservation (rest_id, count, customer_phone_num, is_entry)
        VALUES (?, ?, ?, ?)
    `
	result, err := db.DB.Exec(query, r.RestID, r.Count, r.CustomerPhoneNum, r.IsEntry)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	r.ID = id
	return nil
}
*/
