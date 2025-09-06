package queue

import (
	"context"
	"time"
)

// PatientQueue represents a patient in the queue
type PatientQueue struct {
	ID                string    `json:"id"`
	AppointmentID     string    `json:"appointment_id"`
	OrganizationID    string    `json:"organization_id"`
	Position          int       `json:"position"`
	EstimatedWaitTime int       `json:"estimated_wait_time"` // minutes
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// QueueStatus constants
const (
	QueueStatusWaiting    = "waiting"
	QueueStatusCalled     = "called"
	QueueStatusInProgress = "in_progress"
	QueueStatusCompleted  = "completed"
	QueueStatusCancelled  = "cancelled"
)

// Repository defines the interface for queue data operations
type Repository interface {
	Create(ctx context.Context, queue *PatientQueue) error
	GetByID(ctx context.Context, id string) (*PatientQueue, error)
	GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*PatientQueue, error)
	GetByAppointment(ctx context.Context, appointmentID string) (*PatientQueue, error)
	Update(ctx context.Context, queue *PatientQueue) error
	Delete(ctx context.Context, id string) error
	CountByOrganization(ctx context.Context, organizationID string) (int, error)
	GetNextInQueue(ctx context.Context, organizationID string) (*PatientQueue, error)
	UpdatePosition(ctx context.Context, organizationID string) error
	GetQueueStats(ctx context.Context, organizationID string) (*QueueStats, error)
}

// QueueStats represents aggregated queue statistics
type QueueStats struct {
	Total           int    `json:"total"`
	Waiting         int    `json:"waiting"`
	InProgress      int    `json:"in_progress"`
	AverageWaitTime string `json:"average_wait_time"`
}
