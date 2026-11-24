package handlers

import (
	"net/http"

	"github.com/GDG-KHU-Side/backend-side-project/httpx"
	"github.com/GDG-KHU-Side/backend-side-project/services"
)

type RestaurantHandler struct {
	service *services.RestaurantService
}

func NewRestaurantHandler(s *services.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{service: s}
}

// GET /api/restaurant-list
func (h *RestaurantHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context())
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, list)
}

// GET /api/restaurant/{id}
func (h *RestaurantHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	rest, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, rest)
}
