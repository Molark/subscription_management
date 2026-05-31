package repository

import (
	"app/internal/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPool(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("pgxpool.ParseConfig", slog.Any("error", err))
		return nil, err
	}
	config.MaxConns = 10
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.MinConns = 2
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("pgxpool.NewWithConfig", slog.Any("error", err))
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		slog.Error("pool.Ping", slog.Any("error", err))
		return nil, err
	}
	slog.Info("Pool established")
	return pool, nil
}

var ErrNotFound = errors.New("subscription not found")

var _ Repository = (*repositoryStruct)(nil)

func (r *repositoryStruct) Create(ctx context.Context, sub models.Subscription) (models.Subscription, error) {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;
	`

	var created models.Subscription
	err := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserId,
		sub.StartDate,
		sub.EndDate,
	).Scan(
		&created.Id,
		&created.ServiceName,
		&created.Price,
		&created.UserId,
		&created.StartDate,
		&created.EndDate,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		slog.Error("Repo: failed to execute Create query", slog.Any("error", err))
		return models.Subscription{}, fmt.Errorf("failed to create subscription: %w", err)
	}

	return created, nil
}

func (r *repositoryStruct) GetById(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1;
	`

	var sub models.Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.Id,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserId,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Subscription{}, ErrNotFound
		}
		slog.Error("Repo: failed to execute GetById query", slog.String("id", id.String()), slog.Any("error", err))
		return models.Subscription{}, fmt.Errorf("failed to get subscription by id: %w", err)
	}

	return sub, nil
}

func (r *repositoryStruct) Update(ctx context.Context, sub models.Subscription) (models.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;
	`

	var updated models.Subscription
	err := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.Id,
	).Scan(
		&updated.Id,
		&updated.ServiceName,
		&updated.Price,
		&updated.UserId,
		&updated.StartDate,
		&updated.EndDate,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Subscription{}, ErrNotFound
		}
		slog.Error("Repo: failed to execute Update query", slog.String("id", sub.Id.String()), slog.Any("error", err))
		return models.Subscription{}, fmt.Errorf("failed to update subscription: %w", err)
	}

	return updated, nil
}

func (r *repositoryStruct) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1;`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Repo: failed to execute Delete query", slog.String("id", id.String()), slog.Any("error", err))
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *repositoryStruct) List(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Repo: failed to execute List query", slog.Any("error", err))
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.Id,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserId,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			slog.Error("Repo: failed to scan subscription row during List", slog.Any("error", err))
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Repo: rows error during List execution", slog.Any("error", err))
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return subs, nil
}

func (r *repositoryStruct) CalculateTotal(ctx context.Context, filter models.TotalPriceFilter) (int, error) {

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1 
		  AND start_date >= $2 
		  AND start_date <= $3
	`

	var total int
	var err error

	if filter.ServiceName != "" {
		query += " AND service_name = $4;"
		err = r.pool.QueryRow(ctx, query,
			filter.UserId,
			filter.StartDate,
			filter.EndDate,
			filter.ServiceName,
		).Scan(&total)
	} else {
		query += ";"
		err = r.pool.QueryRow(ctx, query,
			filter.UserId,
			filter.StartDate,
			filter.EndDate,
		).Scan(&total)
	}

	if err != nil {
		slog.Error("Repo: failed to calculate total price", slog.Any("error", err))
		return 0, fmt.Errorf("failed to calculate total price: %w", err)
	}

	return total, nil
}
