package handlers

import (
	"net/http"

	"github.com/GDG-KHU-Side/backend-side-project/httpx"
	"github.com/GDG-KHU-Side/backend-side-project/models"
	"github.com/GDG-KHU-Side/backend-side-project/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// POST /api/register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	u, err := h.service.Register(r.Context(), req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, u)
}

// POST /api/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// GET /api/user  (본인)
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, err := currentUserID(r)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	u, err := h.service.Get(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

// PUT /api/user  (본인)
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := currentUserID(r)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	var req models.UpdateUserRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	u, err := h.service.Update(r.Context(), userID, req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

// DELETE /api/user  (본인 탈퇴)
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := currentUserID(r)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	if err := h.service.Delete(r.Context(), userID); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /api/link  (본인을 식당 사장님으로 등록)
func (h *UserHandler) LinkRestaurant(w http.ResponseWriter, r *http.Request) {
	userID, err := currentUserID(r)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	var req models.LinkRestaurantRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	ur, err := h.service.LinkRestaurant(r.Context(), userID, req.RestID)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, ur)
}
