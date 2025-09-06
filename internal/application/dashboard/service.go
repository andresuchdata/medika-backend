package dashboard

import (
	"context"
	"fmt"
	"time"

	"medika-backend/internal/domain/appointment"
	"medika-backend/internal/domain/doctor"
	"medika-backend/internal/domain/patient"
	"medika-backend/internal/domain/queue"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type Service struct {
	patientRepo     PatientRepository
	appointmentRepo AppointmentRepository
	queueRepo       QueueRepository
	doctorRepo      DoctorRepository
	logger          logger.Logger
}

// Repository interfaces for dependency injection
type PatientRepository interface {
	CountByOrganization(ctx context.Context, organizationID string) (int, error)
	GetByID(ctx context.Context, id string) (*patient.Patient, error)
}

type AppointmentRepository interface {
	GetAppointmentsByDate(ctx context.Context, organizationID, date string, limit int) ([]*appointment.Appointment, error)
	CountAppointmentsByDate(ctx context.Context, organizationID, date string) (int, error)
}

type QueueRepository interface {
	CountByOrganization(ctx context.Context, organizationID string) (int, error)
	GetQueueStats(ctx context.Context, organizationID string) (*queue.QueueStats, error)
}

type DoctorRepository interface {
	GetByID(ctx context.Context, id string) (*doctor.Doctor, error)
}

func NewService(
	patientRepo PatientRepository,
	appointmentRepo AppointmentRepository,
	queueRepo QueueRepository,
	doctorRepo DoctorRepository,
	logger logger.Logger,
) *Service {
	return &Service{
		patientRepo:     patientRepo,
		appointmentRepo: appointmentRepo,
		queueRepo:       queueRepo,
		doctorRepo:      doctorRepo,
		logger:          logger,
	}
}

// GetDashboardSummary returns aggregated dashboard data
func (s *Service) GetDashboardSummary(ctx context.Context, organizationID string) (*dto.DashboardSummaryResponse, error) {
	// Get today's date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")

	// Fetch all data concurrently for better performance
	type statsResult struct {
		totalPatients      int
		todaysAppointments int
		queueLength        int
		averageWaitTime    string
		err                error
	}

	statsChan := make(chan statsResult, 1)
	recentAppointmentsChan := make(chan []dto.RecentAppointment, 1)

	// Fetch statistics concurrently
	go func() {
		var result statsResult

		// Get total patients count
		if totalPatients, err := s.patientRepo.CountByOrganization(ctx, organizationID); err != nil {
			result.err = fmt.Errorf("failed to count patients: %w", err)
			statsChan <- result
			return
		} else {
			result.totalPatients = totalPatients
		}

		// Get today's appointments count
		if todaysAppointments, err := s.appointmentRepo.CountAppointmentsByDate(ctx, organizationID, today); err != nil {
			result.err = fmt.Errorf("failed to count today's appointments: %w", err)
			statsChan <- result
			return
		} else {
			result.todaysAppointments = todaysAppointments
		}

		// Get queue statistics
		if queueStats, err := s.queueRepo.GetQueueStats(ctx, organizationID); err != nil {
			result.err = fmt.Errorf("failed to get queue stats: %w", err)
			statsChan <- result
			return
		} else {
			result.queueLength = queueStats.Total
			result.averageWaitTime = queueStats.AverageWaitTime
		}

		statsChan <- result
	}()

	// Fetch recent appointments concurrently
	go func() {
		recentAppointments, err := s.getRecentAppointments(ctx, organizationID, today, 5)
		if err != nil {
			s.logger.Error(ctx, "Failed to get recent appointments", "error", err)
			recentAppointmentsChan <- []dto.RecentAppointment{}
			return
		}
		recentAppointmentsChan <- recentAppointments
	}()

	// Wait for statistics
	stats := <-statsChan
	if stats.err != nil {
		return nil, stats.err
	}

	// Wait for recent appointments
	recentAppointments := <-recentAppointmentsChan

	// Build response
	response := &dto.DashboardSummaryResponse{
		Success: true,
		Data: struct {
			Stats struct {
				TotalPatients      int    `json:"total_patients"`
				TodaysAppointments int    `json:"todays_appointments"`
				QueueLength        int    `json:"queue_length"`
				AverageWaitTime    string `json:"average_wait_time"`
				MonthlyGrowth      string `json:"monthly_growth"`
				Revenue            string `json:"revenue"`
			} `json:"stats"`
			RecentAppointments []dto.RecentAppointment `json:"recent_appointments"`
			SystemStatus       dto.SystemStatus        `json:"system_status"`
		}{
			Stats: struct {
				TotalPatients      int    `json:"total_patients"`
				TodaysAppointments int    `json:"todays_appointments"`
				QueueLength        int    `json:"queue_length"`
				AverageWaitTime    string `json:"average_wait_time"`
				MonthlyGrowth      string `json:"monthly_growth"`
				Revenue            string `json:"revenue"`
			}{
				TotalPatients:      stats.totalPatients,
				TodaysAppointments: stats.todaysAppointments,
				QueueLength:        stats.queueLength,
				AverageWaitTime:    stats.averageWaitTime,
				MonthlyGrowth:      "+12%", // TODO: Calculate from historical data
				Revenue:            "$12,450", // TODO: Calculate from billing data
			},
			RecentAppointments: recentAppointments,
			SystemStatus: dto.SystemStatus{
				Database:      "operational",
				FileStorage:   "operational",
				Notifications: "operational",
			},
		},
		Message: "Dashboard summary retrieved successfully",
	}

	return response, nil
}

// getRecentAppointments fetches recent appointments with patient and doctor names
func (s *Service) getRecentAppointments(ctx context.Context, organizationID, date string, limit int) ([]dto.RecentAppointment, error) {
	appointments, err := s.appointmentRepo.GetAppointmentsByDate(ctx, organizationID, date, limit)
	if err != nil {
		return nil, err
	}

	recentAppointments := make([]dto.RecentAppointment, len(appointments))
	for i, apt := range appointments {
		// Format time to 12-hour format
		timeFormatted := s.formatTime(apt.StartTime)
		
		recentAppointments[i] = dto.RecentAppointment{
			ID:          apt.ID,
			PatientName: s.getPatientName(ctx, apt.PatientID),
			DoctorName:  s.getDoctorName(ctx, apt.DoctorID),
			Time:        timeFormatted,
			Status:      string(apt.Status),
			Type:        apt.Type,
		}
	}

	return recentAppointments, nil
}

// formatTime converts 24-hour format to 12-hour format
func (s *Service) formatTime(timeStr string) string {
	// Parse the time string (assuming HH:MM format)
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return timeStr // Return original if parsing fails
	}

	// Format to 12-hour format
	return t.Format("3:04 PM")
}

// getPatientName retrieves patient name by ID
func (s *Service) getPatientName(ctx context.Context, patientID string) string {
	patient, err := s.patientRepo.GetByID(ctx, patientID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get patient", "error", err, "patientID", patientID)
		return "Unknown Patient"
	}

	return patient.Name
}

// getDoctorName retrieves doctor name by ID
func (s *Service) getDoctorName(ctx context.Context, doctorID string) string {
	doctor, err := s.doctorRepo.GetByID(ctx, doctorID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get doctor", "error", err, "doctorID", doctorID)
		return "Unknown Doctor"
	}

	return doctor.Name
}
