package repository

import (
	"app/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryStruct struct {
	pool *pgxpool.Pool
}

type Repository interface {
	Create(ctx context.Context, sub models.Subscription) (models.Subscription, error)
	GetById(ctx context.Context, id uuid.UUID) (models.Subscription, error)
	Update(ctx context.Context, sub models.Subscription) (models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]models.Subscription, error)

	CalculateTotal(ctx context.Context, filter models.TotalPriceFilter) (int, error)
}

func NewRepository(pool *pgxpool.Pool) *repositoryStruct {
	return &repositoryStruct{pool: pool}
}
