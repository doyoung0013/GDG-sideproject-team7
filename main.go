package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GDG-KHU-Side/backend-side-project/config"
	"github.com/GDG-KHU-Side/backend-side-project/db"
	"github.com/GDG-KHU-Side/backend-side-project/handlers"
	"github.com/GDG-KHU-Side/backend-side-project/services"
)

func main() {
	conf := config.GetDBConfig()

	err := db.InitDB(conf)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.DB.Close()

	r := mux.NewRouter()

	restaurantService := &services.RestaurantService{}
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService)

	userService := &services.UserService{}
	userHandler := handlers.NewUserHandler(userService)

	reservationService := &services.ReservationService{}
	reservationHandler := handlers.NewReservationHandler(reservationService)

	//page
	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/restaurant-list", handlers.RestaurantListHandler)
	r.HandleFunc("/restaurant/{id:[0-9]+}", handlers.RestaurantDetailHandler)

	r.HandleFunc("/login", handlers.UserLoginHandler)

	r.HandleFunc("/reservation-list/{id:[0-9]+}", handlers.ReservationListHandler)

	//api
	r.HandleFunc("/api/restaurant-list", restaurantHandler.GetRestaurantList).Methods("GET")
	r.HandleFunc("/api/restaurant/{id}", restaurantHandler.GetRestaurant).Methods("GET")
	r.HandleFunc("/api/reservation", restaurantHandler.CreateReservation).Methods("POST")
	r.HandleFunc("/api/reservation/{id}/status", restaurantHandler.UpdateReservationStatus).Methods("PUT")

	r.HandleFunc("/api/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/api/register", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/user", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/user/{id}", userHandler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/api/user/reservation-list/{id}", reservationHandler.GetReservationList).Methods("GET")

	r.HandleFunc("/api/link", userHandler.LinkRestaurant).Methods("POST")

	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
