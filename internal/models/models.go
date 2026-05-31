package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Id          uuid.UUID `json:"id" db:"id"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price" db:"price"`
	UserId      uuid.UUID `json:"user_id" db:"user_id"`

	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" db:"end_date"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type TotalPriceFilter struct {
	UserId      uuid.UUID
	ServiceName string
	StartDate   time.Time
	EndDate     time.Time
}
