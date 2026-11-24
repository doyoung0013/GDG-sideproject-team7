package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/GDG-KHU-Side/backend-side-project/handlers"
	"github.com/GDG-KHU-Side/backend-side-project/httpx"
	"github.com/GDG-KHU-Side/backend-side-project/middleware"
)

type routerDeps struct {
	restaurants  *handlers.RestaurantHandler
	users        *handlers.UserHandler
	reservations *handlers.ReservationHandler
	tokens       middleware.TokenParser
}

// apiRoutes는 "/api" 접두사를 붙여 루트 라우터에 직접 등록합니다.
// (gorilla/mux 서브라우터는 메서드 불일치를 405가 아닌 404로 처리하는 문제가 있음)
type apiRoutes struct{ r *mux.Router }

func (a apiRoutes) HandleFunc(path string, f http.HandlerFunc) *mux.Route {
	return a.r.HandleFunc("/api"+path, f)
}

func (a apiRoutes) Handle(path string, h http.Handler) *mux.Route {
	return a.r.Handle("/api"+path, h)
}

func newRouter(d routerDeps) *mux.Router {
	r := mux.NewRouter()
	authed := middleware.RequireAuth(d.tokens)
	protect := func(h http.HandlerFunc) http.Handler { return authed(h) }

	// 페이지
	r.HandleFunc("/", handlers.IndexPage).Methods(http.MethodGet)
	r.HandleFunc("/restaurant-list", handlers.RestaurantListPage).Methods(http.MethodGet)
	r.HandleFunc("/restaurant/{id:[0-9]+}", handlers.RestaurantDetailPage).Methods(http.MethodGet)
	r.HandleFunc("/login", handlers.LoginPage).Methods(http.MethodGet)
	r.HandleFunc("/reservation-list/{id:[0-9]+}", handlers.ReservationListPage).Methods(http.MethodGet)

	api := apiRoutes{r}

	// 공개 API
	api.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)
	api.HandleFunc("/restaurant-list", d.restaurants.List).Methods(http.MethodGet)
	api.HandleFunc("/restaurant/{id}", d.restaurants.Get).Methods(http.MethodGet)
	api.HandleFunc("/reservation", d.reservations.Create).Methods(http.MethodPost)
	api.HandleFunc("/register", d.users.Register).Methods(http.MethodPost)
	api.HandleFunc("/login", d.users.Login).Methods(http.MethodPost)

	// 로그인 필요 API (사용자 ID는 토큰에서 가져옴)
	api.Handle("/user", protect(d.users.Me)).Methods(http.MethodGet)
	api.Handle("/user", protect(d.users.Update)).Methods(http.MethodPut)
	api.Handle("/user", protect(d.users.Delete)).Methods(http.MethodDelete)
	api.Handle("/link", protect(d.users.LinkRestaurant)).Methods(http.MethodPost)
	api.Handle("/restaurant/{id}/reservations", protect(d.reservations.ListByRestaurant)).Methods(http.MethodGet)
	api.Handle("/reservation/{id}/status", protect(d.reservations.UpdateStatus)).Methods(http.MethodPut)

	// 404/405는 루트 라우터에 설정해야 서브라우터의 405가 404로 덮이지 않습니다.
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/api/") {
			httpx.WriteErrorMessage(w, http.StatusNotFound, "존재하지 않는 API입니다.")
			return
		}
		http.NotFound(w, req)
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteErrorMessage(w, http.StatusMethodNotAllowed, "허용되지 않은 메서드입니다.")
	})
	return r
}
