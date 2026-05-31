package main

import (
	"app/internal/config"
	"app/internal/handlers"
	"app/internal/repository"
	"app/internal/service"
	"log/slog"
	"net/http"
	"time"
)

// @title           Subscription Tracker API
// @version         1.0
// @description     Микросервис для управления подписками пользователей.

// @host      localhost:8080
// @BasePath  /
func main() {

	cfg := config.Load()
	dsn := "postgres://" + cfg.DBUser + ":" + cfg.DBPassword + "@postgres:" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=disable"
	pool, err := repository.ConnectPool(dsn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	repo := repository.NewRepository(pool)
	serviceS := service.NewService(repo)
	handler := handlers.NewHandler(serviceS)
	mux := http.NewServeMux()
	handler.SetupRoutes(mux)
	server := http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("server started at localhost:", cfg.AppPort)
}
