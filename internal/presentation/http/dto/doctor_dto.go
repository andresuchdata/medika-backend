package dto

// DoctorResponse represents a doctor in API responses
type DoctorResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	Specialty      string   `json:"specialty"`
	LicenseNumber  string   `json:"licenseNumber"`
	Status         string   `json:"status"`
	OrganizationID string   `json:"organizationId"`
	Avatar         *string  `json:"avatar,omitempty"`
	Bio            *string  `json:"bio,omitempty"`
	Experience     *int     `json:"experience,omitempty"`
	Education      []string `json:"education,omitempty"`
	Certifications []string `json:"certifications,omitempty"`
	NextAvailable  *string  `json:"nextAvailable,omitempty"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

// CreateDoctorRequest represents the request to create a doctor
type CreateDoctorRequest struct {
	Name           string   `json:"name" validate:"required,min=2,max=100"`
	Email          string   `json:"email" validate:"required,email"`
	Password       string   `json:"password" validate:"required,min=8"`
	Phone          string   `json:"phone" validate:"required"`
	Specialty      string   `json:"specialty" validate:"required"`
	LicenseNumber  string   `json:"license_number" validate:"required"`
	OrganizationID string   `json:"organization_id" validate:"required,uuid"`
	Bio            *string  `json:"bio,omitempty"`
	Experience     *int     `json:"experience,omitempty"`
	Education      []string `json:"education,omitempty"`
	Certifications []string `json:"certifications,omitempty"`
	NextAvailable  *string  `json:"nextAvailable,omitempty"`
}

// UpdateDoctorRequest represents the request to update a doctor
type UpdateDoctorRequest struct {
	Name           *string   `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone          *string   `json:"phone,omitempty"`
	Specialty      *string   `json:"specialty,omitempty"`
	LicenseNumber  *string   `json:"license_number,omitempty"`
	Status         *string   `json:"status,omitempty" validate:"omitempty,oneof=active inactive on-leave"`
	Bio            *string   `json:"bio,omitempty"`
	Experience     *int      `json:"experience,omitempty"`
	Education      []string  `json:"education,omitempty"`
	Certifications []string  `json:"certifications,omitempty"`
	NextAvailable  *string   `json:"nextAvailable,omitempty"`
}

// DoctorStats represents statistics about doctors
type DoctorStats struct {
	Total       int            `json:"total"`
	Active      int            `json:"active"`
	Inactive    int            `json:"inactive"`
	OnLeave     int            `json:"onLeave"`
	Specialties map[string]int `json:"specialties"`
}

// DoctorsData represents the data structure for doctors list response
type DoctorsData struct {
	Doctors    []DoctorResponse `json:"doctors"`
	Pagination Pagination       `json:"pagination"`
	Stats      DoctorStats      `json:"stats"`
}

// DoctorsResponse represents the response for doctors list
type DoctorsResponse struct {
	Success bool        `json:"success"`
	Data    DoctorsData `json:"data"`
	Message string      `json:"message"`
}
