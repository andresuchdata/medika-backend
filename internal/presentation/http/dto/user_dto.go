package dto

import (
	"time"
)

// Request DTOs
type CreateUserRequest struct {
	Name           string `json:"name" validate:"required,min=2,max=100"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Role           string `json:"role" validate:"required,oneof=admin doctor nurse patient cashier"`
	OrganizationID string `json:"organization_id,omitempty" validate:"omitempty,uuid"`
	Phone          string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone       *string    `json:"phone,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Address     *string    `json:"address,omitempty"`
}

type UpdateMedicalInfoRequest struct {
	EmergencyContact *string  `json:"emergency_contact,omitempty"`
	MedicalHistory   *string  `json:"medical_history,omitempty"`
	Allergies        []string `json:"allergies,omitempty"`
	BloodType        *string  `json:"blood_type,omitempty" validate:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
}

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" validate:"required,url"`
}

// Response DTOs
type UserResponse struct {
	ID             string           `json:"id"`
	Email          string           `json:"email"`
	Name           string           `json:"name"`
	Role           string           `json:"role"`
	OrganizationID *string          `json:"organization_id,omitempty"`
	Phone          *string          `json:"phone,omitempty"`
	AvatarURL      *string          `json:"avatar_url,omitempty"`
	IsActive       bool             `json:"is_active"`
	Profile        *ProfileResponse `json:"profile,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type ProfileResponse struct {
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Gender           *string    `json:"gender,omitempty"`
	Address          *string    `json:"address,omitempty"`
	EmergencyContact *string    `json:"emergency_contact,omitempty"`
	MedicalHistory   *string    `json:"medical_history,omitempty"`
	Allergies        []string   `json:"allergies,omitempty"`
	BloodType        *string    `json:"blood_type,omitempty"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// Common response structures
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Version   string                 `json:"version"`
	Services  map[string]interface{} `json:"services,omitempty"`
}
