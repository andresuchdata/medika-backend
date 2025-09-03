package organization

import (
	"context"
	"time"
)

type Organization struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Address       string          `json:"address"`
	Phone         string          `json:"phone"`
	Email         string          `json:"email"`
	Website       *string         `json:"website,omitempty"`
	BusinessHours []BusinessHours `json:"businessHours"`
	IsActive      bool            `json:"isActive"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type BusinessHours struct {
	Day    int    `json:"day"`    // 0-6 (Sunday-Saturday)
	Open   string `json:"open"`   // "09:00"
	Close  string `json:"close"`  // "17:00"
	IsOpen bool   `json:"isOpen"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, organization *Organization) error
	GetByID(ctx context.Context, id string) (*Organization, error)
	GetAll(ctx context.Context, limit, offset int) ([]*Organization, error)
	Update(ctx context.Context, organization *Organization) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}
