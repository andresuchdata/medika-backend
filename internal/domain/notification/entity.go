package notification

import (
	"time"

	"medika-backend/internal/domain/shared"

	"github.com/google/uuid"
)

// Notification represents a notification in the domain
type Notification struct {
	id          NotificationID
	userID      shared.UserID
	title       string
	message     string
	notificationType NotificationType
	priority    Priority
	isRead      bool
	channels    []string
	data        map[string]interface{}
	scheduledFor *time.Time
	sentAt      *time.Time
	createdAt   time.Time
}

// NotificationID represents a unique notification identifier
type NotificationID string

// NewNotificationID creates a new notification ID
func NewNotificationID() NotificationID {
	return NotificationID(uuid.New().String())
}

// String returns the string representation of NotificationID
func (id NotificationID) String() string {
	return string(id)
}

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeAppointment NotificationType = "appointment"
	NotificationTypePatient     NotificationType = "patient"
	NotificationTypeAlert       NotificationType = "alert"
	NotificationTypeMessage     NotificationType = "message"
	NotificationTypeSystem      NotificationType = "system"
	NotificationTypeSchedule    NotificationType = "schedule"
	NotificationTypeLab         NotificationType = "lab"
	NotificationTypeEmergency   NotificationType = "emergency"
)

// Priority represents the priority level of a notification
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// NewNotification creates a new notification
func NewNotification(
	userID shared.UserID,
	title, message string,
	notificationType NotificationType,
	priority Priority,
	channels []string,
	data map[string]interface{},
) *Notification {
	now := time.Now()
	return &Notification{
		id:               NewNotificationID(),
		userID:           userID,
		title:            title,
		message:          message,
		notificationType: notificationType,
		priority:         priority,
		isRead:           false,
		channels:         channels,
		data:             data,
		createdAt:        now,
	}
}

// Getters
func (n *Notification) ID() NotificationID {
	return n.id
}

func (n *Notification) UserID() shared.UserID {
	return n.userID
}

func (n *Notification) Title() string {
	return n.title
}

func (n *Notification) Message() string {
	return n.message
}

func (n *Notification) Type() NotificationType {
	return n.notificationType
}

func (n *Notification) Priority() Priority {
	return n.priority
}

func (n *Notification) IsRead() bool {
	return n.isRead
}

func (n *Notification) Channels() []string {
	return n.channels
}

func (n *Notification) Data() map[string]interface{} {
	return n.data
}

func (n *Notification) ScheduledFor() *time.Time {
	return n.scheduledFor
}

func (n *Notification) SentAt() *time.Time {
	return n.sentAt
}

func (n *Notification) CreatedAt() time.Time {
	return n.createdAt
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.isRead = true
}

// MarkAsUnread marks the notification as unread
func (n *Notification) MarkAsUnread() {
	n.isRead = false
}

// UpdateData updates the notification data
func (n *Notification) UpdateData(data map[string]interface{}) {
	n.data = data
}
