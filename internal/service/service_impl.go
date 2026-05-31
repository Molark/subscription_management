package service

import (
	"app/internal/models"
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var (
	ErrInvalidDates    = errors.New("start date cannot be after end date")
	ErrInvalidPageArgs = errors.New("page and pageSize must be greater than zero")
)

func (s *ServiceStruct) Create(ctx context.Context, sub models.Subscription) (models.Subscription, error) {
	slog.Info("Service: creating subscription", slog.String("service_name", sub.ServiceName), slog.String("user_id", sub.UserId.String()))

	if sub.EndDate != nil && sub.StartDate.After(*sub.EndDate) {
		slog.Warn("Service: failed to create subscription, invalid dates",
			slog.Time("start", sub.StartDate),
			slog.Time("end", *sub.EndDate),
		)
		return models.Subscription{}, ErrInvalidDates
	}

	createdSub, err := s.repo.Create(ctx, sub)
	if err != nil {
		slog.Error("Service: repo failed to create subscription", slog.Any("err", err))
		return models.Subscription{}, err
	}

	slog.Info("Service: subscription created successfully", slog.String("id", createdSub.Id.String()))
	return createdSub, nil
}

func (s *ServiceStruct) GetById(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	slog.Info("Service: fetching subscription by id", slog.String("id", id.String()))

	sub, err := s.repo.GetById(ctx, id)
	if err != nil {
		slog.Error("Service: repo failed to get subscription", slog.String("id", id.String()), slog.Any("err", err))
		return models.Subscription{}, err
	}

	return sub, nil
}

func (s *ServiceStruct) Update(ctx context.Context, sub models.Subscription) (models.Subscription, error) {
	slog.Info("Service: updating subscription", slog.String("id", sub.Id.String()))

	if sub.EndDate != nil && sub.StartDate.After(*sub.EndDate) {
		slog.Warn("Service: failed to update subscription, invalid dates",
			slog.Time("start", sub.StartDate),
			slog.Time("end", *sub.EndDate),
		)
		return models.Subscription{}, ErrInvalidDates
	}

	updatedSub, err := s.repo.Update(ctx, sub)
	if err != nil {
		slog.Error("Service: repo failed to update subscription", slog.String("id", sub.Id.String()), slog.Any("err", err))
		return models.Subscription{}, err
	}

	slog.Info("Service: subscription updated successfully", slog.String("id", updatedSub.Id.String()))
	return updatedSub, nil
}

func (s *ServiceStruct) Delete(ctx context.Context, id uuid.UUID) error {
	slog.Info("Service: deleting subscription", slog.String("id", id.String()))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		slog.Error("Service: repo failed to delete subscription", slog.String("id", id.String()), slog.Any("err", err))
		return err
	}

	slog.Info("Service: subscription deleted successfully", slog.String("id", id.String()))
	return nil
}

func (s *ServiceStruct) List(ctx context.Context, page, pageSize int) ([]models.Subscription, error) {
	slog.Info("Service: listing subscriptions", slog.Int("page", page), slog.Int("page_size", pageSize))

	if page <= 0 || pageSize <= 0 {
		slog.Warn("Service: invalid pagination arguments", slog.Int("page", page), slog.Int("page_size", pageSize))
		return nil, ErrInvalidPageArgs
	}

	limit := pageSize
	offset := (page - 1) * pageSize

	subs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		slog.Error("Service: repo failed to list subscriptions", slog.Any("err", err))
		return nil, err
	}

	return subs, nil
}

func (s *ServiceStruct) GetTotalPrice(ctx context.Context, filter models.TotalPriceFilter) (int, error) {
	slog.Info("Service: calculating total price",
		slog.String("user_id", filter.UserId.String()),
		slog.String("service_name", filter.ServiceName),
		slog.Time("start_period", filter.StartDate),
		slog.Time("end_period", filter.EndDate),
	)

	if filter.StartDate.After(filter.EndDate) {
		slog.Warn("Service: total price calculation failed, invalid period filter")
		return 0, ErrInvalidDates
	}

	total, err := s.repo.CalculateTotal(ctx, filter)
	if err != nil {
		slog.Error("Service: repo failed to calculate total price", slog.Any("err", err))
		return 0, err
	}

	slog.Info("Service: total price calculated successfully", slog.Int("total", total))
	return total, nil
}
