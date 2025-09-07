package dto

import "time"

// NotificationResponse represents a notification in API responses
type NotificationResponse struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Message      string                 `json:"message"`
	Type         string                 `json:"type"`
	Priority     string                 `json:"priority"`
	IsRead       bool                   `json:"is_read"`
	Channels     []string               `json:"channels,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	ScheduledFor *time.Time             `json:"scheduled_for,omitempty"`
	SentAt       *time.Time             `json:"sent_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// NotificationListResponse represents a paginated list of notifications
type NotificationListResponse struct {
	Data       []NotificationResponse `json:"data"`
	Pagination Pagination             `json:"pagination"`
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	Title          string                 `json:"title" validate:"required"`
	Message        string                 `json:"message" validate:"required"`
	Type           string                 `json:"type" validate:"required,oneof=appointment patient alert message system schedule lab emergency"`
	Priority       string                 `json:"priority" validate:"required,oneof=low medium high critical"`
	ActionRequired bool                   `json:"action_required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateNotificationRequest represents a request to update a notification
type UpdateNotificationRequest struct {
	Title          *string                `json:"title,omitempty"`
	Message        *string                `json:"message,omitempty"`
	Type           *string                `json:"type,omitempty" validate:"omitempty,oneof=appointment patient alert message system schedule lab emergency"`
	Priority       *string                `json:"priority,omitempty" validate:"omitempty,oneof=low medium high critical"`
	ActionRequired *bool                  `json:"action_required,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
