package dto

import "time"

// QueueRequest represents the request to create/update a queue
type QueueRequest struct {
	AppointmentID     string `json:"appointment_id" validate:"required,uuid"`
	OrganizationID    string `json:"organization_id" validate:"required,uuid"`
	Position          int    `json:"position,omitempty"`
	EstimatedWaitTime int    `json:"estimated_wait_time,omitempty"`
	Status            string `json:"status,omitempty"`
}

// QueueResponse represents a queue response
type QueueResponse struct {
	ID                string    `json:"id"`
	AppointmentID     string    `json:"appointment_id"`
	OrganizationID    string    `json:"organization_id"`
	Position          int       `json:"position"`
	EstimatedWaitTime int       `json:"estimated_wait_time"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// QueuesResponse represents the response for multiple queues
type QueuesResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Queues     []QueueResponse `json:"queues"`
		Pagination Pagination      `json:"pagination"`
	} `json:"data"`
	Message string `json:"message"`
}

// QueueDetailResponse represents a detailed queue response
type QueueDetailResponse struct {
	Success bool          `json:"success"`
	Data    QueueResponse `json:"data"`
	Message string        `json:"message"`
}

// QueueActionRequest represents a request for queue actions
type QueueActionRequest struct {
	Action string `json:"action" validate:"required,oneof=call start complete cancel"`
}

// QueuePositionUpdate represents a queue position update
type QueuePositionUpdate struct {
	ID       string `json:"id" validate:"required,uuid"`
	Position int    `json:"position" validate:"required,min=1"`
}

// PatientQueueResponse represents a queue response with enriched patient and appointment data
type PatientQueueResponse struct {
	ID                string    `json:"id"`
	AppointmentID     string    `json:"appointment_id"`
	OrganizationID    string    `json:"organization_id"`
	Position          int       `json:"position"`
	EstimatedWaitTime int       `json:"estimated_wait_time"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	
	// Enriched data
	PatientName    string `json:"patient_name"`
	PatientID      string `json:"patient_id"`
	DoctorName     string `json:"doctor_name"`
	DoctorID       string `json:"doctor_id"`
	AppointmentDate string `json:"appointment_date"`
	AppointmentTime string `json:"appointment_time"`
	AppointmentType string `json:"appointment_type"`
	AppointmentStatus string `json:"appointment_status"`
}

// PatientQueueDetailResponse represents a detailed patient queue response
type PatientQueueDetailResponse struct {
	Success bool                  `json:"success"`
	Data    *PatientQueueResponse `json:"data"`
	Message string                `json:"message"`
}
