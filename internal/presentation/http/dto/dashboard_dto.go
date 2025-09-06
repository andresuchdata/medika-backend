package dto

// DashboardSummaryResponse represents the dashboard summary data
type DashboardSummaryResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Stats struct {
			TotalPatients      int    `json:"total_patients"`
			TodaysAppointments int    `json:"todays_appointments"`
			QueueLength        int    `json:"queue_length"`
			AverageWaitTime    string `json:"average_wait_time"`
			MonthlyGrowth      string `json:"monthly_growth"`
			Revenue            string `json:"revenue"`
		} `json:"stats"`
		RecentAppointments []RecentAppointment `json:"recent_appointments"`
		SystemStatus       SystemStatus        `json:"system_status"`
	} `json:"data"`
	Message string `json:"message"`
}

// RecentAppointment represents a recent appointment for dashboard
type RecentAppointment struct {
	ID          string `json:"id"`
	PatientName string `json:"patient_name"`
	DoctorName  string `json:"doctor_name"`
	Time        string `json:"time"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

// SystemStatus represents system health status
type SystemStatus struct {
	Database     string `json:"database"`
	FileStorage  string `json:"file_storage"`
	Notifications string `json:"notifications"`
}

// DashboardStatsRequest represents request parameters for dashboard stats
type DashboardStatsRequest struct {
	OrganizationID string `json:"organization_id,omitempty"`
	Date          string `json:"date,omitempty"` // YYYY-MM-DD format
}
