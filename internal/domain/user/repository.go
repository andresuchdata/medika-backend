package user

import (
	"context"
	"time"

	"medika-backend/internal/domain/shared"
)

// Repository interface (port)
type Repository interface {
	// Command operations
	Save(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id shared.UserID) error

	// Query operations
	FindByID(ctx context.Context, id shared.UserID) (*User, error)
	FindByEmail(ctx context.Context, email shared.Email) (*User, error)
	FindAll(ctx context.Context, filters UserFilters) ([]*User, error)
	Count(ctx context.Context, filters UserFilters) (int64, error)

	// Legacy methods for backward compatibility
	FindByOrganization(ctx context.Context, orgID shared.OrganizationID, filters UserFilters) ([]*User, error)
	FindByRole(ctx context.Context, role Role, orgID *shared.OrganizationID) ([]*User, error)
	CountByOrganization(ctx context.Context, orgID shared.OrganizationID) (int64, error)
	CountByRole(ctx context.Context, role Role) (int64, error)
	FindActiveUsers(ctx context.Context, orgID *shared.OrganizationID) ([]*User, error)

	// Transaction support
	WithTx(ctx context.Context, fn func(Repository) error) error
}

// UserFilters for query filtering
type UserFilters struct {
	// Basic filters
	Name         string
	Email        string
	Role         *Role
	IsActive     *bool
	OrganizationID *shared.OrganizationID
	
	// Date filters
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
	
	// Pagination
	Limit  int
	Offset int
	
	// Sorting
	OrderBy string // e.g., "created_at", "name", "email"
	Order   string // "ASC" or "DESC"
}

// Events
type UserCreatedEvent struct {
	UserID         shared.UserID
	Email          shared.Email
	Role           Role
	OrganizationID *shared.OrganizationID
	CreatedAt      time.Time
}

func (e UserCreatedEvent) EventType() string {
	return "user.created"
}

func (e UserCreatedEvent) EventData() map[string]interface{} {
	data := map[string]interface{}{
		"user_id":    e.UserID.String(),
		"email":      e.Email.String(),
		"role":       e.Role.String(),
		"created_at": e.CreatedAt,
	}
	if e.OrganizationID != nil {
		data["organization_id"] = e.OrganizationID.String()
	}
	return data
}

type UserUpdatedEvent struct {
	UserID    shared.UserID
	Changes   map[string]interface{}
	UpdatedAt time.Time
}

type UserDeactivatedEvent struct {
	UserID         shared.UserID
	DeactivatedAt  time.Time
	DeactivatedBy  shared.UserID
}
