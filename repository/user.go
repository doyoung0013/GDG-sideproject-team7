package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

type userRepo struct{ base }

func NewUserRepository(db *sql.DB, timeout time.Duration) UserRepository {
	return &userRepo{base{db: db, timeout: timeout}}
}

var errEmailTaken = apperr.Conflict("이미 사용 중인 이메일입니다.")

const userSelect = `SELECT id, email, password, phone_num, name, created_at, updated_at FROM users`

func (repo *userRepo) scanOne(ctx context.Context, where string, arg any) (*models.User, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	var u models.User
	err := repo.db.QueryRowContext(ctx, userSelect+" WHERE "+where, arg).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.PhoneNum, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NotFound("사용자를 찾을 수 없습니다.")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (repo *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return repo.scanOne(ctx, "id = ?", id)
}

func (repo *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return repo.scanOne(ctx, "email = ?", email)
}

func (repo *userRepo) Create(ctx context.Context, u *models.User) error {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO users (email, password, phone_num, name) VALUES (?, ?, ?, ?)`,
		u.Email, u.PasswordHash, u.PhoneNum, u.Name)
	if mysqlErrNumber(err) == mysqlErrDuplicate {
		return errEmailTaken
	}
	if err != nil {
		return err
	}
	u.ID, err = res.LastInsertId()
	return err
}

func (repo *userRepo) Update(ctx context.Context, u *models.User) error {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	res, err := repo.db.ExecContext(ctx,
		`UPDATE users SET name = ?, email = ?, phone_num = ? WHERE id = ?`,
		u.Name, u.Email, u.PhoneNum, u.ID)
	if mysqlErrNumber(err) == mysqlErrDuplicate {
		return errEmailTaken
	}
	if err != nil {
		return err
	}
	return requireAffected(res, "사용자를 찾을 수 없습니다.")
}

func (repo *userRepo) Delete(ctx context.Context, id int64) error {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	res, err := repo.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return requireAffected(res, "사용자를 찾을 수 없습니다.")
}

func (repo *userRepo) ListRestaurantIDs(ctx context.Context, userID int64) ([]int64, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx,
		`SELECT rest_id FROM user_rest WHERE user_id = ? ORDER BY rest_id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (repo *userRepo) LinkRestaurant(ctx context.Context, userID, restID int64) (*models.UserRest, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	_, err := repo.db.ExecContext(ctx, `INSERT INTO user_rest (user_id, rest_id) VALUES (?, ?)`, userID, restID)
	switch mysqlErrNumber(err) {
	case mysqlErrDuplicate:
		// PK(user_id, rest_id) 또는 UNIQUE(rest_id) 위반: 이미 연결되어 있거나 다른 사장님이 소유
		return nil, apperr.Conflict("이미 사장님이 등록된 식당입니다.")
	case mysqlErrNoParent:
		return nil, apperr.NotFound("식당을 찾을 수 없습니다.")
	}
	if err != nil {
		return nil, err
	}
	return &models.UserRest{UserID: userID, RestID: restID, CreatedAt: time.Now()}, nil
}

func (repo *userRepo) OwnsRestaurant(ctx context.Context, userID, restID int64) (bool, error) {
	ctx, cancel := repo.ctx(ctx)
	defer cancel()

	var ok bool
	err := repo.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM user_rest WHERE user_id = ? AND rest_id = ?)`, userID, restID).Scan(&ok)
	return ok, err
}

func requireAffected(res sql.Result, notFoundMsg string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return apperr.NotFound(notFoundMsg)
	}
	return nil
}
