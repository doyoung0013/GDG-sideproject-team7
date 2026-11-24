package handlers

import "net/http"

// 정적 HTML 페이지. 데이터는 각 페이지가 API로 가져옵니다.
func page(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/"+file)
	}
}

var (
	IndexPage            = page("index.html")
	LoginPage            = page("login.html")
	RestaurantListPage   = page("restaurant_list.html")
	RestaurantDetailPage = page("restaurant_detail.html")
	ReservationListPage  = page("reservation_list.html")
)
