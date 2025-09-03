package dto

import (
	"time"
)

// Patient DTOs
type PatientResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	Avatar    *string   `json:"avatar,omitempty"`
	Address   AddressResponse `json:"address"`
	EmergencyContact EmergencyContactResponse `json:"emergencyContact"`
	MedicalHistory []MedicalConditionResponse `json:"medicalHistory"`
	Allergies []string `json:"allergies"`
	Medications []MedicationResponse `json:"medications"`
	LastVisit *time.Time `json:"lastVisit,omitempty"`
	NextAppointment *time.Time `json:"nextAppointment,omitempty"`
	Status    string    `json:"status"`
	OrganizationID string `json:"organizationId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AddressResponse struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

type EmergencyContactResponse struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Phone        string `json:"phone"`
}

type MedicalConditionResponse struct {
	Condition     string     `json:"condition"`
	DiagnosedDate time.Time  `json:"diagnosedDate"`
	Status        string     `json:"status"`
	Notes         *string    `json:"notes,omitempty"`
}

type MedicationResponse struct {
	Name           string     `json:"name"`
	Dosage         string     `json:"dosage"`
	PrescribedDate time.Time  `json:"prescribedDate"`
	Status         string     `json:"status"`
}

// Patient list response
type PatientsResponse struct {
	Success bool         `json:"success"`
	Data    PatientsData `json:"data"`
	Message string       `json:"message"`
}

type PatientsData struct {
	Patients   []PatientResponse `json:"patients"`
	Pagination Pagination        `json:"pagination"`
	Stats      PatientStats      `json:"stats"`
}

type PatientStats struct {
	Total        int                `json:"total"`
	Active       int                `json:"active"`
	Inactive     int                `json:"inactive"`
	AverageAge   int                `json:"averageAge"`
	GenderDistribution GenderDistribution `json:"genderDistribution"`
}

type GenderDistribution struct {
	Male   int `json:"male"`
	Female int `json:"female"`
}
