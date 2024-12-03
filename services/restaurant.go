package services

import (
	"github.com/GDG-KHU-Side/backend-side-project/db"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

type RestaurantService struct{}

func (s *RestaurantService) GetAllRestaurants() ([]models.Restaurant, error) {
	rows, err := db.DB.Query(`
        SELECT r.*, COUNT(res.id) as waiting_count 
        FROM restaurants r 
        LEFT JOIN reservation res ON r.id = res.rest_id AND res.is_entry = 0
        GROUP BY r.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var r models.Restaurant
		var waitingCount int
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Address,
			&r.Description,
			&r.PhoneNum,
			&r.WaitingTime,
			&r.CreatedAt,
			&r.UpdatedAt,
			&waitingCount,
		)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	return restaurants, nil
}

func (s *RestaurantService) GetRestaurant(id int64) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := db.DB.QueryRow("SELECT id, name, address, description, phone_num, waiting_time, created_at, updated_at FROM restaurants WHERE id = ?", id).
		Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.Description,
			&restaurant.PhoneNum,
			&restaurant.WaitingTime,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (s *RestaurantService) CreateReservation(r *models.Reservation) error {
	result, err := db.DB.Exec(`
        INSERT INTO reservation (rest_id, count, customer_phone_num) 
        VALUES (?, ?, ?)`,
		r.RestID, r.Count, r.CustomerPhoneNum)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	r.ID = id
	return nil
}

func (s *RestaurantService) UpdateReservationStatus(id int64, status int) error {
	_, err := db.DB.Exec("UPDATE reservation SET is_entry = ? WHERE id = ?", status, id)
	return err
}
