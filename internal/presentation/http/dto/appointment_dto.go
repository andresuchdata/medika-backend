package dto

// AppointmentResponse represents an appointment in API responses
type AppointmentResponse struct {
	ID             string  `json:"id"`
	PatientID      string  `json:"patientId"`
	PatientName    string  `json:"patientName"`
	DoctorID       string  `json:"doctorId"`
	DoctorName     string  `json:"doctorName"`
	OrganizationID string  `json:"organizationId"`
	RoomID         *string `json:"roomId,omitempty"`
	Date           string  `json:"date"`
	StartTime      string  `json:"startTime"`
	EndTime        string  `json:"endTime"`
	Duration       int     `json:"duration"`
	Status         string  `json:"status"`
	Type           string  `json:"type"`
	Notes          *string `json:"notes,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// CreateAppointmentRequest represents the request to create an appointment
type CreateAppointmentRequest struct {
	PatientID      string  `json:"patientId" validate:"required,uuid"`
	DoctorID       string  `json:"doctorId" validate:"required,uuid"`
	OrganizationID string  `json:"organizationId" validate:"required,uuid"`
	RoomID         *string `json:"roomId,omitempty" validate:"omitempty,uuid"`
	Date           string  `json:"date" validate:"required"`
	StartTime      string  `json:"startTime" validate:"required"`
	EndTime        string  `json:"endTime" validate:"required"`
	Duration       int     `json:"duration" validate:"required,min=15,max=480"`
	Type           string  `json:"type" validate:"required,oneof=consultation follow_up emergency routine_checkup"`
	Notes          *string `json:"notes,omitempty"`
}

// UpdateAppointmentRequest represents the request to update an appointment
type UpdateAppointmentRequest struct {
	Date      *string `json:"date,omitempty"`
	StartTime *string `json:"startTime,omitempty"`
	EndTime   *string `json:"endTime,omitempty"`
	Duration  *int    `json:"duration,omitempty" validate:"omitempty,min=15,max=480"`
	Type      *string `json:"type,omitempty" validate:"omitempty,oneof=consultation follow_up emergency routine_checkup"`
	Notes     *string `json:"notes,omitempty"`
	RoomID    *string `json:"roomId,omitempty" validate:"omitempty,uuid"`
}

// UpdateAppointmentStatusRequest represents the request to update appointment status
type UpdateAppointmentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed in_progress completed cancelled no_show"`
}

// AppointmentStats represents statistics about appointments
type AppointmentStats struct {
	Total      int `json:"total"`
	Pending    int `json:"pending"`
	Confirmed  int `json:"confirmed"`
	InProgress int `json:"inProgress"`
	Completed  int `json:"completed"`
	Cancelled  int `json:"cancelled"`
	NoShow     int `json:"noShow"`
}

// AppointmentsData represents the data structure for appointments list response
type AppointmentsData struct {
	Appointments []AppointmentResponse `json:"appointments"`
	Pagination   Pagination           `json:"pagination"`
	Stats        AppointmentStats      `json:"stats"`
}

// AppointmentsResponse represents the response for appointments list
type AppointmentsResponse struct {
	Success bool              `json:"success"`
	Data    AppointmentsData  `json:"data"`
	Message string            `json:"message"`
}
