package patient

import (
	"context"
	"time"
)

type Patient struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	Avatar    *string   `json:"avatar,omitempty"`
	Address   Address   `json:"address"`
	EmergencyContact EmergencyContact `json:"emergencyContact"`
	MedicalHistory []MedicalCondition `json:"medicalHistory"`
	Allergies []string `json:"allergies"`
	Medications []Medication `json:"medications"`
	LastVisit *time.Time `json:"lastVisit,omitempty"`
	NextAppointment *time.Time `json:"nextAppointment,omitempty"`
	Status    string    `json:"status"`
	OrganizationID string `json:"organizationId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

type EmergencyContact struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Phone        string `json:"phone"`
}

type MedicalCondition struct {
	Condition     string     `json:"condition"`
	DiagnosedDate time.Time  `json:"diagnosedDate"`
	Status        string     `json:"status"`
	Notes         *string    `json:"notes,omitempty"`
}

type Medication struct {
	Name           string     `json:"name"`
	Dosage         string     `json:"dosage"`
	PrescribedDate time.Time  `json:"prescribedDate"`
	Status         string     `json:"status"`
}

// Repository interface
type Repository interface {
	Create(ctx context.Context, patient *Patient) error
	GetByID(ctx context.Context, id string) (*Patient, error)
	GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*Patient, error)
	Update(ctx context.Context, patient *Patient) error
	Delete(ctx context.Context, id string) error
	CountByOrganization(ctx context.Context, organizationID string) (int, error)
}
