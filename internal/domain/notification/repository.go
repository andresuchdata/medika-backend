package notification

import (
	"context"
	"medika-backend/internal/domain/shared"
)

// Repository defines the interface for notification data operations
type Repository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *Notification) error

	// FindByID finds a notification by ID
	FindByID(ctx context.Context, id NotificationID) (*Notification, error)

	// FindByUserID finds notifications for a specific user
	FindByUserID(ctx context.Context, userID shared.UserID, filters NotificationFilters) ([]*Notification, error)

	// CountByUserID counts notifications for a specific user
	CountByUserID(ctx context.Context, userID shared.UserID, filters NotificationFilters) (int, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id NotificationID) error

	// MarkAsUnread marks a notification as unread
	MarkAsUnread(ctx context.Context, id NotificationID) error

	// MarkAllAsRead marks all notifications for a user as read
	MarkAllAsRead(ctx context.Context, userID shared.UserID) error

	// Delete deletes a notification
	Delete(ctx context.Context, id NotificationID) error

	// DeleteByUserID deletes all notifications for a user
	DeleteByUserID(ctx context.Context, userID shared.UserID) error
}

// NotificationFilters represents filters for notification queries
type NotificationFilters struct {
	Type            *NotificationType `json:"type,omitempty"`
	Priority        *Priority         `json:"priority,omitempty"`
	IsRead          *bool             `json:"is_read,omitempty"`
	ActionRequired  *bool             `json:"action_required,omitempty"`
	Limit           int               `json:"limit,omitempty"`
	Offset          int               `json:"offset,omitempty"`
	OrderBy         string            `json:"order_by,omitempty"`
	Order           string            `json:"order,omitempty"`
}
