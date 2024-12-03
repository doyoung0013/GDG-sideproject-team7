package services

import (
	"github.com/GDG-KHU-Side/backend-side-project/db"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

type ReservationService struct{}

func (s *ReservationService) GetAllRestaurants(id int64) ([]models.Reservation, error) {
	rows, err := db.DB.Query(`
        SELECT r.*
        FROM reservation r
        WHERE r.rest_id = ?
        ORDER BY is_entry, id DESC
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var r models.Reservation

		err := rows.Scan(
			&r.ID,
			&r.RestID,
			&r.Count,
			&r.CustomerPhoneNum,
			&r.CreatedAt,
			&r.IsEntry,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}
