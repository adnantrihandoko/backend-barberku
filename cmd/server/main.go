package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/barberku/backend-barber/internal/handler"
	"github.com/barberku/backend-barber/internal/middleware"
	"github.com/barberku/backend-barber/internal/service"
	"github.com/barberku/backend-barber/pkg/config"
	"github.com/barberku/backend-barber/pkg/websocket"
)

func main() {
	cfg := config.Load()

	slog.Info("starting BarberKu server", "port", cfg.ServerPort)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.WsOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	hub := websocket.NewHub()
	go hub.Run()

	authService := service.NewAuthService(nil, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	broadcaster := func(event string, data interface{}) {
		hub.Broadcast(event, data)
	}

	fcmService := service.NewFCMService(nil, cfg.FCMServerKey)
	fcmHandler := handler.NewFCMHandler(fcmService)

	queueService := service.NewQueueService(nil, nil, broadcaster, fcmService)
	queueHandler := handler.NewQueueHandler(queueService)

	serviceService := service.NewServiceService(nil)
	serviceHandler := handler.NewServiceHandler(serviceService)

	barberService := service.NewBarberService(nil)
	barberHandler := handler.NewBarberHandler(barberService)

	historyService := service.NewHistoryService(nil)
	historyHandler := handler.NewHistoryHandler(historyService)

	storeSettingsService := service.NewStoreSettingsService(nil)
	storeSettingsHandler := handler.NewStoreSettingsHandler(storeSettingsService)

	statsService := service.NewStatsService(nil)
	statsHandler := handler.NewStatsHandler(statsService)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		authHandler.RegisterRoutes(r)
	})

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("failed to upgrade connection", "error", err)
			return
		}

		client := websocket.NewClient(hub, conn)
		hub.register <- client

		go client.WritePump()
		go client.ReadPump()
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))

		r.Route("/queue", func(r chi.Router) {
			queueHandler.RegisterRoutes(r)
		})

		r.Route("/services", func(r chi.Router) {
			serviceHandler.RegisterRoutes(r)
		})

		r.Route("/barbers", func(r chi.Router) {
			barberHandler.RegisterRoutes(r)
		})

		r.Route("/history", func(r chi.Router) {
			historyHandler.RegisterRoutes(r)
		})

		r.Route("/settings", func(r chi.Router) {
			storeSettingsHandler.RegisterRoutes(r)
		})

		r.Route("/stats", func(r chi.Router) {
			statsHandler.RegisterRoutes(r)
		})

		r.Route("/fcm", func(r chi.Router) {
			fcmHandler.RegisterRoutes(r)
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		slog.Info("server listening on", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
