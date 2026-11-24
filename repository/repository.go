// Package repository는 DB 접근을 담당합니다. 서비스 계층은 아래 인터페이스에만 의존하므로
// 테스트에서는 가짜 구현으로 바꿔 끼울 수 있습니다.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/GDG-KHU-Side/backend-side-project/models"
)

type RestaurantRepository interface {
	List(ctx context.Context) ([]models.Restaurant, error)
	GetByID(ctx context.Context, id int64) (*models.Restaurant, error)
}

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id int64) error

	ListRestaurantIDs(ctx context.Context, userID int64) ([]int64, error)
	LinkRestaurant(ctx context.Context, userID, restID int64) (*models.UserRest, error)
	OwnsRestaurant(ctx context.Context, userID, restID int64) (bool, error)
}

type ReservationRepository interface {
	Create(ctx context.Context, r *models.Reservation) error
	GetByID(ctx context.Context, id int64) (*models.Reservation, error)
	ListByRestaurant(ctx context.Context, restID int64) ([]models.Reservation, error)
	UpdateStatus(ctx context.Context, id int64, status int8) error
	HasActiveWaiting(ctx context.Context, restID int64, phone string) (bool, error)
}

// base는 MySQL 구현체들이 공유하는 DB 핸들과 쿼리 타임아웃입니다.
type base struct {
	db      *sql.DB
	timeout time.Duration
}

func (b base) ctx(parent context.Context) (context.Context, context.CancelFunc) {
	if b.timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, b.timeout)
}

const (
	mysqlErrDuplicate = 1062
	mysqlErrNoParent  = 1452
)

func mysqlErrNumber(err error) uint16 {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number
	}
	return 0
}
