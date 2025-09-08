package appointment

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/appointment"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type Service struct {
	appointmentRepo appointment.Repository
	logger          logger.Logger
}

func NewService(appointmentRepo appointment.Repository, logger logger.Logger) *Service {
	return &Service{
		appointmentRepo: appointmentRepo,
		logger:          logger,
	}
}

// Methods that match the AppointmentService interface
func (s *Service) GetAppointments(ctx *fiber.Ctx, organizationID, doctorID, patientID string, limit, offset int) ([]*appointment.Appointment, error) {
	// If specific filters are provided, use them
	if patientID != "" {
		return s.appointmentRepo.GetByPatient(ctx.Context(), patientID, limit, offset)
	}
	if doctorID != "" {
		return s.appointmentRepo.GetByDoctor(ctx.Context(), doctorID, limit, offset)
	}
	if organizationID != "" {
		return s.appointmentRepo.GetByOrganization(ctx.Context(), organizationID, limit, offset)
	}
	
	// If no specific filters, get appointments for the user's organization from auth context
	userOrgID := ctx.Locals("organization_id")
	if userOrgID != nil && userOrgID != "" {
		return s.appointmentRepo.GetByOrganization(ctx.Context(), userOrgID.(string), limit, offset)
	}
	
	// Fallback: if no auth context, return empty (should not happen in production)
	return []*appointment.Appointment{}, nil
}

func (s *Service) CountAppointments(ctx *fiber.Ctx, organizationID, doctorID, patientID string) (int, error) {
	// If specific filters are provided, we need to implement specific counting methods
	// For now, use organization-based counting as the primary method
	if organizationID != "" {
		return s.appointmentRepo.CountByOrganization(ctx.Context(), organizationID)
	}
	
	// If no organization filter, count appointments for the user's organization from auth context
	userOrgID := ctx.Locals("organization_id")
	if userOrgID != nil && userOrgID != "" {
		return s.appointmentRepo.CountByOrganization(ctx.Context(), userOrgID.(string))
	}
	
	// Fallback: if no auth context, return 0 (should not happen in production)
	return 0, nil
}

func (s *Service) GetAppointmentByID(ctx *fiber.Ctx, appointmentID string) (*appointment.Appointment, error) {
	return s.appointmentRepo.GetByID(ctx.Context(), appointmentID)
}

func (s *Service) CreateAppointment(ctx *fiber.Ctx, appointmentData *dto.CreateAppointmentRequest) (*appointment.Appointment, error) {
	// Parse the date from string to time.Time
	// Try ISO format first, then fallback to date-only format
	var date time.Time
	var err error
	
	// Try parsing as ISO format (e.g., "2025-09-09T00:00:00.000Z")
	date, err = time.Parse(time.RFC3339, appointmentData.Date)
	if err != nil {
		// Fallback to date-only format (e.g., "2025-09-09")
		date, err = time.Parse("2006-01-02", appointmentData.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, expected ISO format or YYYY-MM-DD: %w", err)
		}
	}

	// Organization ID is required for appointments
	// It should come from the request (set by the staff member creating the appointment)
	organizationID := appointmentData.OrganizationID
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required for appointments")
	}

	// Create domain entity from DTO
	apt := &appointment.Appointment{
		PatientID:      appointmentData.PatientID,
		DoctorID:       appointmentData.DoctorID,
		OrganizationID: organizationID,
		RoomID:         appointmentData.RoomID,
		Date:           date,
		StartTime:      appointmentData.StartTime,
		EndTime:        appointmentData.EndTime,
		Duration:       appointmentData.Duration,
		Status:         appointment.StatusPending, // Default status
		Type:           appointmentData.Type,
		Notes:          appointmentData.Notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Create appointment in repository
	if err := s.appointmentRepo.Create(ctx.Context(), apt); err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return apt, nil
}

func (s *Service) UpdateAppointment(ctx *fiber.Ctx, appointmentID string, appointmentData *dto.UpdateAppointmentRequest) (*appointment.Appointment, error) {
	// For now, return error - implement update logic later
	return nil, fmt.Errorf("not implemented yet")
}

func (s *Service) DeleteAppointment(ctx *fiber.Ctx, appointmentID string) error {
	return s.appointmentRepo.Delete(ctx.Context(), appointmentID)
}

func (s *Service) UpdateAppointmentStatus(ctx *fiber.Ctx, appointmentID string, status string) (*appointment.Appointment, error) {
	// Convert string status to AppointmentStatus and update
	// For now, return error - implement status update logic later
	return nil, fmt.Errorf("not implemented yet")
}

// Additional methods for future use
func (s *Service) CreateAppointmentDirect(ctx context.Context, apt *appointment.Appointment) error {
	return s.appointmentRepo.Create(ctx, apt)
}

func (s *Service) GetAppointmentDirect(ctx context.Context, id string) (*appointment.Appointment, error) {
	return s.appointmentRepo.GetByID(ctx, id)
}

func (s *Service) UpdateAppointmentDirect(ctx context.Context, apt *appointment.Appointment) error {
	return s.appointmentRepo.Update(ctx, apt)
}

func (s *Service) UpdateAppointmentStatusDirect(ctx context.Context, id string, status appointment.AppointmentStatus) error {
	return s.appointmentRepo.UpdateStatus(ctx, id, status)
}
