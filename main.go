package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GDG-KHU-Side/backend-side-project/auth"
	"github.com/GDG-KHU-Side/backend-side-project/config"
	"github.com/GDG-KHU-Side/backend-side-project/db"
	"github.com/GDG-KHU-Side/backend-side-project/handlers"
	"github.com/GDG-KHU-Side/backend-side-project/middleware"
	"github.com/GDG-KHU-Side/backend-side-project/repository"
	"github.com/GDG-KHU-Side/backend-side-project/services"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn, err := db.Open(ctx, cfg.DB)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 의존성 조립: repository -> service -> handler
	restaurantRepo := repository.NewRestaurantRepository(conn, cfg.QueryTimeout)
	userRepo := repository.NewUserRepository(conn, cfg.QueryTimeout)
	reservationRepo := repository.NewReservationRepository(conn, cfg.QueryTimeout)

	tokens := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTTTL)

	router := newRouter(routerDeps{
		restaurants:  handlers.NewRestaurantHandler(services.NewRestaurantService(restaurantRepo)),
		users:        handlers.NewUserHandler(services.NewUserService(userRepo, tokens)),
		reservations: handlers.NewReservationHandler(services.NewReservationService(reservationRepo, restaurantRepo, userRepo)),
		tokens:       tokens,
	})

	limiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	go limiter.Cleanup(ctx)

	srv := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: middleware.Chain(router,
			middleware.Recover,
			middleware.Logging,
			middleware.CORS(cfg.CORSAllowedOrigins),
			limiter.Middleware,
		),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Println("server stopped")
	return nil
}
