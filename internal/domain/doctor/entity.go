package doctor

import (
	"context"
	"time"
)

type Doctor struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Specialization string    `json:"specialization"`
	LicenseNumber  string    `json:"licenseNumber"`
	Experience     int       `json:"experience"`
	Education      string    `json:"education"`
	Avatar         *string   `json:"avatar,omitempty"`
	Bio            string    `json:"bio"`
	Status         string    `json:"status"`
	OrganizationID string    `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, doctor *Doctor) error
	GetByID(ctx context.Context, id string) (*Doctor, error)
	GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*Doctor, error)
	Update(ctx context.Context, doctor *Doctor) error
	Delete(ctx context.Context, id string) error
	CountByOrganization(ctx context.Context, organizationID string) (int, error)
}
