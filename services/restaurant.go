package services

import (
	"context"

	"github.com/GDG-KHU-Side/backend-side-project/models"
	"github.com/GDG-KHU-Side/backend-side-project/repository"
)

type RestaurantService struct {
	restaurants repository.RestaurantRepository
}

func NewRestaurantService(r repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{restaurants: r}
}

func (s *RestaurantService) List(ctx context.Context) ([]models.Restaurant, error) {
	return s.restaurants.List(ctx)
}

func (s *RestaurantService) Get(ctx context.Context, id int64) (*models.Restaurant, error) {
	return s.restaurants.GetByID(ctx, id)
}
