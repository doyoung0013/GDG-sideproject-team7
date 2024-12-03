package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GDG-KHU-Side/backend-side-project/models"
	"github.com/GDG-KHU-Side/backend-side-project/services"
	"github.com/gorilla/mux"
)

type RestaurantHandler struct {
	service *services.RestaurantService
}

func NewRestaurantHandler(s *services.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{service: s}
}

func RestaurantListHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/restaurant_list.html")
}

func RestaurantDetailHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/restaurant_detail.html")
}

func (h *RestaurantHandler) GetRestaurantList(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.service.GetAllRestaurants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(restaurants)
}

func (h *RestaurantHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	restaurant, err := h.service.GetRestaurant(id)
	if err != nil {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(restaurant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *RestaurantHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateReservation(&reservation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reservation)
}

func (h *RestaurantHandler) UpdateReservationStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	var status struct {
		IsEntry int `json:"is_entry"`
	}
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateReservationStatus(id, status.IsEntry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
