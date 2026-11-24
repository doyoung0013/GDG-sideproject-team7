package services

import (
	"context"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/models"
	"github.com/GDG-KHU-Side/backend-side-project/repository"
)

// OwnershipChecker는 사용자가 식당의 사장님인지 확인합니다.
type OwnershipChecker interface {
	OwnsRestaurant(ctx context.Context, userID, restID int64) (bool, error)
}

type ReservationService struct {
	reservations repository.ReservationRepository
	restaurants  repository.RestaurantRepository
	owners       OwnershipChecker
}

func NewReservationService(res repository.ReservationRepository, rest repository.RestaurantRepository, owners OwnershipChecker) *ReservationService {
	return &ReservationService{reservations: res, restaurants: rest, owners: owners}
}

var errNotOwner = apperr.Forbidden("이 식당의 예약에 접근할 권한이 없습니다.")

// Create는 손님의 웨이팅을 등록합니다(로그인 불필요).
func (s *ReservationService) Create(ctx context.Context, req models.CreateReservationRequest) (*models.Reservation, error) {
	if req.RestID <= 0 {
		return nil, apperr.Invalid("rest_id 가 필요합니다.")
	}
	if req.Count < minPartySize || req.Count > maxPartySize {
		return nil, apperr.Invalid("인원수는 1~20명이어야 합니다.")
	}
	phone, err := normalizePhone(req.CustomerPhoneNum)
	if err != nil {
		return nil, err
	}

	if _, err := s.restaurants.GetByID(ctx, req.RestID); err != nil {
		return nil, err
	}
	dup, err := s.reservations.HasActiveWaiting(ctx, req.RestID, phone)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, apperr.Conflict("이미 이 식당에 대기 중인 예약이 있습니다.")
	}

	r := &models.Reservation{RestID: req.RestID, Count: req.Count, CustomerPhoneNum: phone, IsEntry: models.StatusWaiting}
	if err := s.reservations.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

// ListByRestaurant는 사장님 본인 식당의 예약 목록을 반환합니다.
func (s *ReservationService) ListByRestaurant(ctx context.Context, userID, restID int64) ([]models.Reservation, error) {
	if err := s.authorize(ctx, userID, restID); err != nil {
		return nil, err
	}
	return s.reservations.ListByRestaurant(ctx, restID)
}

// UpdateStatus는 사장님 본인 식당의 예약 상태를 바꿉니다.
func (s *ReservationService) UpdateStatus(ctx context.Context, userID, reservationID int64, status int8) error {
	switch status {
	case models.StatusWaiting, models.StatusEntered, models.StatusCanceled:
	default:
		return apperr.Invalid("is_entry 는 0(대기), 1(입장), 2(취소) 중 하나여야 합니다.")
	}

	r, err := s.reservations.GetByID(ctx, reservationID)
	if err != nil {
		return err
	}
	if err := s.authorize(ctx, userID, r.RestID); err != nil {
		return err
	}
	return s.reservations.UpdateStatus(ctx, reservationID, status)
}

func (s *ReservationService) authorize(ctx context.Context, userID, restID int64) error {
	ok, err := s.owners.OwnsRestaurant(ctx, userID, restID)
	if err != nil {
		return err
	}
	if !ok {
		return errNotOwner
	}
	return nil
}
