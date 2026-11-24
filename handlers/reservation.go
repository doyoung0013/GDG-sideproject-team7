package handlers

import (
	"net/http"

	"github.com/GDG-KHU-Side/backend-side-project/apperr"
	"github.com/GDG-KHU-Side/backend-side-project/httpx"
	"github.com/GDG-KHU-Side/backend-side-project/models"
	"github.com/GDG-KHU-Side/backend-side-project/services"
)

type ReservationHandler struct {
	service *services.ReservationService
}

func NewReservationHandler(s *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: s}
}

// POST /api/reservation  (손님, 인증 불필요)
func (h *ReservationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateReservationRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	res, err := h.service.Create(r.Context(), req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, res)
}

// GET /api/restaurant/{id}/reservations  (사장님)
func (h *ReservationHandler) ListByRestaurant(w http.ResponseWriter, r *http.Request) {
	userID, err := currentUserID(r)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	restID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	list, err := h.service.ListByRestaurant(r.Context(), userID, restID)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, list)
}

// PUT /api/reservation/{id}/status  (사장님)
func (h *ReservationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := currentUserID(r)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	var req models.UpdateReservationStatusRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	if req.IsEntry == nil {
		httpx.WriteError(w, r, apperr.Invalid("is_entry 가 필요합니다."))
		return
	}
	if err := h.service.UpdateStatus(r.Context(), userID, id, *req.IsEntry); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"id": id, "is_entry": *req.IsEntry})
}
