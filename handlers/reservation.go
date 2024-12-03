package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GDG-KHU-Side/backend-side-project/services"
	"github.com/gorilla/mux"
)

type ReservationHandler struct {
	service *services.ReservationService
}

func NewReservationHandler(s *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: s}
}

func ReservationListHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/reservation_list.html")
}

func (h *ReservationHandler) GetReservationList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	reservations, err := h.service.GetAllRestaurants(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reservations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
