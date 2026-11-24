package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

type reservationRepo struct{ base }

func NewReservationRepository(db *sql.DB, timeout time.Duration) ReservationRepository {
	return &reservationRepo{base{db: db, timeout: timeout}}
}

const reservationSelect = `SELECT id, rest_id, count, customer_phone_num, created_at, is_entry FROM reservation`

func scanReservation(s scanner) (*models.Reservation, error) {
	var r models.Reservation
	err := s.Scan(&r.ID, &r.RestID, &r.Count, &r.CustomerPhoneNum, &r.CreatedAt, &r.IsEntry)
	return &r, err
}

func (repo *reservationRepo) Create(ctx context.Context, r *models.Reservation) error {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO reservation (rest_id, count, customer_phone_num, is_entry) VALUES (?, ?, ?, ?)`,
		r.RestID, r.Count, r.CustomerPhoneNum, r.IsEntry)
	if mysqlErrNumber(err) == mysqlErrNoParent {
		return apperr.NotFound("식당을 찾을 수 없습니다.")
	}
	if err != nil {
		return err
	}
	if r.ID, err = res.LastInsertId(); err != nil {
		return err
	}
	r.CreatedAt = time.Now()
	return nil
}

func (repo *reservationRepo) GetByID(ctx context.Context, id int64) (*models.Reservation, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	r, err := scanReservation(repo.db.QueryRowContext(ctx, reservationSelect+` WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NotFound("예약을 찾을 수 없습니다.")
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ListByRestaurant는 대기 중인 예약을 먼저(먼저 온 순서), 그다음 처리된 예약을 최신순으로 반환합니다.
func (repo *reservationRepo) ListByRestaurant(ctx context.Context, restID int64) ([]models.Reservation, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx, reservationSelect+`
		WHERE rest_id = ?
		ORDER BY is_entry <> 0, CASE WHEN is_entry = 0 THEN id END ASC, id DESC`, restID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Reservation{}
	for rows.Next() {
		r, err := scanReservation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *r)
	}
	return list, rows.Err()
}

func (repo *reservationRepo) UpdateStatus(ctx context.Context, id int64, status int8) error {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	res, err := repo.db.ExecContext(ctx, `UPDATE reservation SET is_entry = ? WHERE id = ?`, status, id)
	if err != nil {
		return err
	}
	return requireAffected(res, "예약을 찾을 수 없습니다.")
}

func (repo *reservationRepo) HasActiveWaiting(ctx context.Context, restID int64, phone string) (bool, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	var ok bool
	err := repo.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM reservation WHERE rest_id = ? AND customer_phone_num = ? AND is_entry = 0
		)`, restID, phone).Scan(&ok)
	return ok, err
}
