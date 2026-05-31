package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	_ "app/docs"
	"app/internal/models"

	"github.com/google/uuid"
)

var (
	ErrInvalidRequest = errors.New("bad request")
	ErrNotFound       = errors.New("not found")
)

type CreateSubscriptionInput struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserId      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type UpdateSubscriptionInput struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

const dateFormatLayout = "01-2006"

// CreateSubscription
// @Summary      Создать запись о подписке
// @Description  Создает новую подписку для конкретного пользователя. Даты передаются в формате "MM-YYYY".
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body CreateSubscriptionInput true "Данные подписки"
// @Success      201  {object}  models.Subscription
// @Failure      400  {string}  string "Bad Request"
// @Router       /subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input CreateSubscriptionInput
	if err := ReadRequestBody(r, &input); err != nil {
		slog.Error("CreateSubscription: failed to decode request body", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	uID, err := uuid.Parse(input.UserId)
	if err != nil {
		slog.Error("CreateSubscription: invalid user_id format", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	startDate, err := time.Parse(dateFormatLayout, input.StartDate)
	if err != nil {
		slog.Error("CreateSubscription: invalid start_date format, expected MM-YYYY", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	var endDate *time.Time
	if input.EndDate != nil && *input.EndDate != "" {
		parsedEnd, err := time.Parse(dateFormatLayout, *input.EndDate)
		if err != nil {
			slog.Error("CreateSubscription: invalid end_date format", slog.Any("err", err))
			RespondError(w, ErrInvalidRequest)
			return
		}
		endDate = &parsedEnd
	}

	sub := models.Subscription{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserId:      uID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	sub, err = h.service.Create(r.Context(), sub)
	if err != nil {
		slog.Error("CreateSubscription: service failed to create subscription", slog.Any("err", err))
		RespondError(w, err)
		return
	}

	slog.Info("Subscription created successfully", slog.String("id", sub.Id.String()))
	RespondJSON(w, http.StatusCreated, sub)
}

// GetSubscription
// @Summary      Получить подписку по ID
// @Description  Возвращает данные конкретной подписки по её UUID первичному ключу.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID подписки"
// @Success      200  {object}  models.Subscription
// @Failure      400  {string}  string "Bad Request"
// @Failure      404  {string}  string "Not Found"
// @Router       /subscriptions/{id} [get]
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("GetSubscription: invalid id in path", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	sub, err := h.service.GetById(r.Context(), id)
	if err != nil {
		slog.Error("GetSubscription: service failed to find subscription", slog.String("id", id.String()), slog.Any("err", err))
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, sub)
}

// UpdateSubscription
// @Summary      Обновить запись о подписке
// @Description  Обновляет изменяемые поля подписки (название, цену, даты) по её ID.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "UUID подписки"
// @Param        input body      UpdateSubscriptionInput  true  "Новые данные подписки"
// @Success      200  {object}  models.Subscription
// @Failure      400  {string}  string "Bad Request"
// @Failure      404  {string}  string "Not Found"
// @Router       /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("UpdateSubscription: invalid id in path", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	var input UpdateSubscriptionInput
	if err := ReadRequestBody(r, &input); err != nil {
		slog.Error("UpdateSubscription: failed to decode request body", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	startDate, err := time.Parse(dateFormatLayout, input.StartDate)
	if err != nil {
		slog.Error("UpdateSubscription: invalid start_date format", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	var endDate *time.Time
	if input.EndDate != nil && *input.EndDate != "" {
		parsedEnd, err := time.Parse(dateFormatLayout, *input.EndDate)
		if err != nil {
			slog.Error("UpdateSubscription: invalid end_date format", slog.Any("err", err))
			RespondError(w, ErrInvalidRequest)
			return
		}
		endDate = &parsedEnd
	}

	sub := models.Subscription{
		Id:          id,
		ServiceName: input.ServiceName,
		Price:       input.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	sub, err = h.service.Update(r.Context(), sub)
	if err != nil {
		slog.Error("UpdateSubscription: service failed to update", slog.String("id", id.String()), slog.Any("err", err))
		RespondError(w, err)
		return
	}

	slog.Info("Subscription updated successfully", slog.String("id", id.String()))
	RespondJSON(w, http.StatusOK, sub)
}

// DeleteSubscription
// @Summary      Удалить запись о подписке
// @Description  Удаляет запись о подписке по переданному ID.
// @Tags         subscriptions
// @Param        id   path      string  true  "UUID подписки"
// @Success      24   "No Content"
// @Failure      400  {string}  string "Bad Request"
// @Failure      404  {string}  string "Not Found"
// @Router       /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("DeleteSubscription: invalid id in path", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		slog.Error("DeleteSubscription: service failed to delete", slog.String("id", id.String()), slog.Any("err", err))
		RespondError(w, err)
		return
	}

	slog.Info("Subscription deleted successfully", slog.String("id", id.String()))
	w.WriteHeader(http.StatusNoContent)
}

// ListSubscriptions
// @Summary      Получить список подписок
// @Description  Возвращает список подписок с поддержкой постраничной пагинации.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        page      query    int  false  "Номер страницы (по умолчанию 1)"
// @Param        pageSize  query    int  false  "Размер страницы (по умолчанию 20)"
// @Success      200  {object}  map[string]interface{} "Возвращает массив subscriptions и метаданные пагинации"
// @Failure      400  {string}  string "Bad Request"
// @Router       /subscriptions [get]
func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	if pageStr == "" {
		pageStr = "1"
	}
	if pageSizeStr == "" {
		pageSizeStr = "20"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		RespondError(w, ErrInvalidRequest)
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		RespondError(w, ErrInvalidRequest)
		return
	}

	subs, err := h.service.List(r.Context(), page, pageSize)
	if err != nil {
		slog.Error("ListSubscriptions: service failed", slog.Any("err", err))
		RespondError(w, err)
		return
	}

	type Response struct {
		Subscriptions []models.Subscription `json:"subscriptions"`
		Page          int                   `json:"page"`
		PageSize      int                   `json:"page_size"`
	}

	RespondJSON(w, http.StatusOK, Response{
		Subscriptions: subs,
		Page:          page,
		PageSize:      pageSize,
	})
}

// GetTotalCost
// @Summary      Подсчет суммарной стоимости подписок
// @Description  Вычисляет сумму всех подписок пользователя за выбранный период с фильтрацией по имени сервиса.
// @Tags         analytics
// @Accept       json
// @Produce      json
// @Param        user_id       query    string  true   "UUID пользователя"
// @Param        start_date    query    string  true   "Начало периода (MM-YYYY)"
// @Param        end_date      query    string  true   "Конец периода (MM-YYYY)"
// @Param        service_name  query    string  false  "Название подписки для фильтрации"
// @Success      200  {object}  map[string]int "Пример: {"total_cost": 1200}"
// @Failure      400  {string}  string "Bad Request"
// @Router       /subscriptions/total [get]
func (h *Handler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userIDStr := query.Get("user_id")
	serviceName := query.Get("service_name")
	startDateStr := query.Get("start_date")
	endDateStr := query.Get("end_date")

	if userIDStr == "" || startDateStr == "" || endDateStr == "" {
		slog.Error("GetTotalCost: missing required query parameters")
		RespondError(w, ErrInvalidRequest)
		return
	}

	uID, err := uuid.Parse(userIDStr)
	if err != nil {
		slog.Error("GetTotalCost: invalid user_id query format", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	startDate, err := time.Parse(dateFormatLayout, startDateStr)
	if err != nil {
		slog.Error("GetTotalCost: invalid start_date query format", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	endDate, err := time.Parse(dateFormatLayout, endDateStr)
	if err != nil {
		slog.Error("GetTotalCost: invalid end_date query format", slog.Any("err", err))
		RespondError(w, ErrInvalidRequest)
		return
	}

	filter := models.TotalPriceFilter{
		UserId:      uID,
		ServiceName: serviceName,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	total, err := h.service.GetTotalPrice(r.Context(), filter)
	if err != nil {
		slog.Error("GetTotalCost: service failed to calculate total", slog.Any("err", err))
		RespondError(w, err)
		return
	}

	type Response struct {
		TotalCost int `json:"total_cost"`
	}

	RespondJSON(w, http.StatusOK, Response{TotalCost: total})
}
