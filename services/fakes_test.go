package services

import (
	"context"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

type fakeRestaurants struct{ byID map[int64]*models.Restaurant }

func (f *fakeRestaurants) List(context.Context) ([]models.Restaurant, error) { return nil, nil }
func (f *fakeRestaurants) GetByID(_ context.Context, id int64) (*models.Restaurant, error) {
	if r, ok := f.byID[id]; ok {
		return r, nil
	}
	return nil, apperr.NotFound("식당을 찾을 수 없습니다.")
}

type fakeReservations struct {
	byID    map[int64]*models.Reservation
	nextID  int64
	updated map[int64]int8
}

func newFakeReservations() *fakeReservations {
	return &fakeReservations{byID: map[int64]*models.Reservation{}, updated: map[int64]int8{}}
}

func (f *fakeReservations) Create(_ context.Context, r *models.Reservation) error {
	f.nextID++
	r.ID = f.nextID
	cp := *r
	f.byID[r.ID] = &cp
	return nil
}
func (f *fakeReservations) GetByID(_ context.Context, id int64) (*models.Reservation, error) {
	if r, ok := f.byID[id]; ok {
		return r, nil
	}
	return nil, apperr.NotFound("예약을 찾을 수 없습니다.")
}
func (f *fakeReservations) ListByRestaurant(_ context.Context, restID int64) ([]models.Reservation, error) {
	var out []models.Reservation
	for _, r := range f.byID {
		if r.RestID == restID {
			out = append(out, *r)
		}
	}
	return out, nil
}
func (f *fakeReservations) UpdateStatus(_ context.Context, id int64, s int8) error {
	f.updated[id] = s
	f.byID[id].IsEntry = s
	return nil
}
func (f *fakeReservations) HasActiveWaiting(_ context.Context, restID int64, phone string) (bool, error) {
	for _, r := range f.byID {
		if r.RestID == restID && r.CustomerPhoneNum == phone && r.IsEntry == models.StatusWaiting {
			return true, nil
		}
	}
	return false, nil
}

type fakeUsers struct {
	byID   map[int64]*models.User
	owners map[int64]int64 // restID -> userID
	nextID int64
}

func newFakeUsers() *fakeUsers {
	return &fakeUsers{byID: map[int64]*models.User{}, owners: map[int64]int64{}}
}

func (f *fakeUsers) Create(_ context.Context, u *models.User) error {
	for _, e := range f.byID {
		if e.Email == u.Email {
			return apperr.Conflict("이미 사용 중인 이메일입니다.")
		}
	}
	f.nextID++
	u.ID = f.nextID
	cp := *u
	f.byID[u.ID] = &cp
	return nil
}
func (f *fakeUsers) GetByID(_ context.Context, id int64) (*models.User, error) {
	if u, ok := f.byID[id]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, apperr.NotFound("사용자를 찾을 수 없습니다.")
}
func (f *fakeUsers) GetByEmail(_ context.Context, email string) (*models.User, error) {
	for _, u := range f.byID {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, apperr.NotFound("사용자를 찾을 수 없습니다.")
}
func (f *fakeUsers) Update(_ context.Context, u *models.User) error {
	e, ok := f.byID[u.ID]
	if !ok {
		return apperr.NotFound("사용자를 찾을 수 없습니다.")
	}
	e.Email, e.Name, e.PhoneNum = u.Email, u.Name, u.PhoneNum
	return nil
}
func (f *fakeUsers) Delete(_ context.Context, id int64) error {
	delete(f.byID, id)
	return nil
}
func (f *fakeUsers) ListRestaurantIDs(_ context.Context, userID int64) ([]int64, error) {
	ids := []int64{}
	for rest, u := range f.owners {
		if u == userID {
			ids = append(ids, rest)
		}
	}
	return ids, nil
}
func (f *fakeUsers) LinkRestaurant(_ context.Context, userID, restID int64) (*models.UserRest, error) {
	if _, taken := f.owners[restID]; taken {
		return nil, apperr.Conflict("이미 사장님이 등록된 식당입니다.")
	}
	f.owners[restID] = userID
	return &models.UserRest{UserID: userID, RestID: restID}, nil
}
func (f *fakeUsers) OwnsRestaurant(_ context.Context, userID, restID int64) (bool, error) {
	return f.owners[restID] == userID, nil
}

type fakeTokens struct{}

func (fakeTokens) Issue(int64) (string, error) { return "token", nil }
