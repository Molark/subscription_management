package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"

	"github.com/google/uuid"
)

type ServiceStruct struct {
	repo repository.Repository
}

type Service interface {
	Create(ctx context.Context, sub models.Subscription) (models.Subscription, error)
	GetById(ctx context.Context, id uuid.UUID) (models.Subscription, error)
	Update(ctx context.Context, sub models.Subscription) (models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]models.Subscription, error)

	GetTotalPrice(ctx context.Context, filter models.TotalPriceFilter) (int, error)
}

func NewService(repo repository.Repository) *ServiceStruct {
	return &ServiceStruct{repo: repo}
}
