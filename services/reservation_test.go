package services

import (
	"context"
	"testing"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/models"
)

func newReservationFixture() (*ReservationService, *fakeReservations, *fakeUsers) {
	rest := &fakeRestaurants{byID: map[int64]*models.Restaurant{1: {ID: 1}, 2: {ID: 2}}}
	res := newFakeReservations()
	users := newFakeUsers()
	users.owners[1] = 100 // 사용자 100이 식당 1의 사장님
	return NewReservationService(res, rest, users), res, users
}

func TestReservationCreate_Validation(t *testing.T) {
	svc, _, _ := newReservationFixture()
	ctx := context.Background()

	cases := []struct {
		name string
		req  models.CreateReservationRequest
		kind apperr.Kind
	}{
		{"rest_id 없음", models.CreateReservationRequest{Count: 2, CustomerPhoneNum: "01012345678"}, apperr.KindInvalid},
		{"인원 0", models.CreateReservationRequest{RestID: 1, Count: 0, CustomerPhoneNum: "01012345678"}, apperr.KindInvalid},
		{"인원 초과", models.CreateReservationRequest{RestID: 1, Count: 21, CustomerPhoneNum: "01012345678"}, apperr.KindInvalid},
		{"전화번호 형식", models.CreateReservationRequest{RestID: 1, Count: 2, CustomerPhoneNum: "02-123-4567"}, apperr.KindInvalid},
		{"없는 식당", models.CreateReservationRequest{RestID: 99, Count: 2, CustomerPhoneNum: "010-1234-5678"}, apperr.KindNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Create(ctx, tc.req)
			if !apperr.Is(err, tc.kind) {
				t.Fatalf("want kind %v, got %v", tc.kind, err)
			}
		})
	}
}

func TestReservationCreate_NormalizesPhoneAndRejectsDuplicate(t *testing.T) {
	svc, _, _ := newReservationFixture()
	ctx := context.Background()

	r, err := svc.Create(ctx, models.CreateReservationRequest{RestID: 1, Count: 2, CustomerPhoneNum: "010-1234-5678"})
	if err != nil {
		t.Fatal(err)
	}
	if r.CustomerPhoneNum != "01012345678" || r.IsEntry != models.StatusWaiting {
		t.Fatalf("unexpected reservation: %+v", r)
	}

	_, err = svc.Create(ctx, models.CreateReservationRequest{RestID: 1, Count: 3, CustomerPhoneNum: "01012345678"})
	if !apperr.Is(err, apperr.KindConflict) {
		t.Fatalf("want conflict, got %v", err)
	}

	// 다른 식당에는 같은 번호로 등록 가능
	if _, err := svc.Create(ctx, models.CreateReservationRequest{RestID: 2, Count: 3, CustomerPhoneNum: "01012345678"}); err != nil {
		t.Fatalf("other restaurant should be allowed: %v", err)
	}
}

func TestReservationOwnerOnly(t *testing.T) {
	svc, res, _ := newReservationFixture()
	ctx := context.Background()
	r, _ := svc.Create(ctx, models.CreateReservationRequest{RestID: 1, Count: 2, CustomerPhoneNum: "01011112222"})

	if _, err := svc.ListByRestaurant(ctx, 200, 1); !apperr.Is(err, apperr.KindForbidden) {
		t.Fatalf("non-owner list: want forbidden, got %v", err)
	}
	if err := svc.UpdateStatus(ctx, 200, r.ID, models.StatusEntered); !apperr.Is(err, apperr.KindForbidden) {
		t.Fatalf("non-owner update: want forbidden, got %v", err)
	}
	if err := svc.UpdateStatus(ctx, 100, r.ID, 5); !apperr.Is(err, apperr.KindInvalid) {
		t.Fatalf("bad status: want invalid, got %v", err)
	}
	if err := svc.UpdateStatus(ctx, 100, 999, models.StatusEntered); !apperr.Is(err, apperr.KindNotFound) {
		t.Fatalf("missing reservation: want not found, got %v", err)
	}

	if err := svc.UpdateStatus(ctx, 100, r.ID, models.StatusEntered); err != nil {
		t.Fatal(err)
	}
	if res.updated[r.ID] != models.StatusEntered {
		t.Fatalf("status not updated")
	}
	list, err := svc.ListByRestaurant(ctx, 100, 1)
	if err != nil || len(list) != 1 {
		t.Fatalf("owner list: %v %v", list, err)
	}
}
