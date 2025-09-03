package dto

// OrganizationResponse represents an organization in API responses
type OrganizationResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Address     string  `json:"address"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Website     *string `json:"website,omitempty"`
	Status      string  `json:"status"`
	StaffCount  int     `json:"staffCount"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// CreateOrganizationRequest represents the request to create an organization
type CreateOrganizationRequest struct {
	Name    string  `json:"name" validate:"required,min=2,max=100"`
	Type    string  `json:"type" validate:"required,oneof=hospital clinic urgent_care private_practice laboratory"`
	Address string  `json:"address" validate:"required"`
	Phone   string  `json:"phone" validate:"required"`
	Email   string  `json:"email" validate:"required,email"`
	Website *string `json:"website,omitempty"`
}

// UpdateOrganizationRequest represents the request to update an organization
type UpdateOrganizationRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Type    *string `json:"type,omitempty" validate:"omitempty,oneof=hospital clinic urgent_care private_practice laboratory"`
	Address *string `json:"address,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
	Website *string `json:"website,omitempty"`
	Status  *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

// OrganizationStats represents statistics about organizations
type OrganizationStats struct {
	Total      int            `json:"total"`
	Active     int            `json:"active"`
	Inactive   int            `json:"inactive"`
	Types      map[string]int `json:"types"`
	TotalStaff int            `json:"totalStaff"`
}

// OrganizationsData represents the data structure for organizations list response
type OrganizationsData struct {
	Organizations []OrganizationResponse `json:"organizations"`
	Pagination    Pagination            `json:"pagination"`
	Stats         OrganizationStats     `json:"stats"`
}

// OrganizationsResponse represents the response for organizations list
type OrganizationsResponse struct {
	Success bool              `json:"success"`
	Data    OrganizationsData `json:"data"`
	Message string            `json:"message"`
}
