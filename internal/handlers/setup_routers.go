package handlers

import (
	_ "app/docs"
	"app/internal/service"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	slog.Info("NewHandler created")
	return &Handler{service: service}
}
func (h Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /subscriptions", h.CreateSubscription)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetSubscription)
	mux.HandleFunc("PUT /subscriptions/{id}", h.UpdateSubscription)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.DeleteSubscription)
	mux.HandleFunc("GET /subscriptions", h.ListSubscriptions)

	mux.HandleFunc("GET /subscriptions/total", h.GetTotalCost)

	mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
