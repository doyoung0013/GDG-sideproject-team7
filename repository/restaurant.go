package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

type restaurantRepo struct{ base }

func NewRestaurantRepository(db *sql.DB, timeout time.Duration) RestaurantRepository {
	return &restaurantRepo{base{db: db, timeout: timeout}}
}

const restaurantSelect = `
	SELECT r.id, r.name, r.address, r.description, r.phone_num, r.waiting_time,
	       r.created_at, r.updated_at, COUNT(res.id) AS waiting_count
	FROM restaurants r
	LEFT JOIN reservation res ON res.rest_id = r.id AND res.is_entry = 0`

type scanner interface{ Scan(dest ...any) error }

func scanRestaurant(s scanner) (*models.Restaurant, error) {
	var r models.Restaurant
	err := s.Scan(&r.ID, &r.Name, &r.Address, &r.Description, &r.PhoneNum, &r.WaitingTime,
		&r.CreatedAt, &r.UpdatedAt, &r.WaitingCount)
	return &r, err
}

func (repo *restaurantRepo) List(ctx context.Context) ([]models.Restaurant, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx, restaurantSelect+` GROUP BY r.id ORDER BY r.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := []models.Restaurant{}
	for rows.Next() {
		r, err := scanRestaurant(rows)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, *r)
	}
	return restaurants, rows.Err()
}

func (repo *restaurantRepo) GetByID(ctx context.Context, id int64) (*models.Restaurant, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	r, err := scanRestaurant(repo.db.QueryRowContext(ctx, restaurantSelect+` WHERE r.id = ? GROUP BY r.id`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NotFound("식당을 찾을 수 없습니다.")
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}
