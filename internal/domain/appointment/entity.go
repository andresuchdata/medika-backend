package appointment

import (
	"context"
	"time"
)

type Appointment struct {
	ID             string            `json:"id"`
	PatientID      string            `json:"patientId"`
	DoctorID       string            `json:"doctorId"`
	OrganizationID string            `json:"organizationId"`
	RoomID         *string           `json:"roomId,omitempty"`
	Date           time.Time         `json:"date"`
	StartTime      string            `json:"startTime"`
	EndTime        string            `json:"endTime"`
	Duration       int               `json:"duration"` // in minutes
	Status         AppointmentStatus `json:"status"`
	Type           string            `json:"type"`
	Notes          *string           `json:"notes,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "pending"
	StatusConfirmed AppointmentStatus = "confirmed"
	StatusInProgress AppointmentStatus = "in_progress"
	StatusCompleted AppointmentStatus = "completed"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusNoShow    AppointmentStatus = "no_show"
)

// Repository interface
type Repository interface {
	Create(ctx context.Context, appointment *Appointment) error
	GetByID(ctx context.Context, id string) (*Appointment, error)
	GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*Appointment, error)
	GetByPatient(ctx context.Context, patientID string, limit, offset int) ([]*Appointment, error)
	GetByDoctor(ctx context.Context, doctorID string, limit, offset int) ([]*Appointment, error)
	Update(ctx context.Context, appointment *Appointment) error
	UpdateStatus(ctx context.Context, id string, status AppointmentStatus) error
	Delete(ctx context.Context, id string) error
	CountByOrganization(ctx context.Context, organizationID string) (int, error)
	GetAppointmentsByDate(ctx context.Context, organizationID, date string, limit int) ([]*Appointment, error)
	CountAppointmentsByDate(ctx context.Context, organizationID, date string) (int, error)
}
