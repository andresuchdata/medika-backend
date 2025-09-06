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
	
	// If no specific filters, get all appointments (similar to dashboard behavior)
	// We need to implement a method to get all appointments, or use a default organization
	// For now, let's get appointments for today's date without organization filter
	return s.appointmentRepo.GetAppointmentsByDate(ctx.Context(), "", time.Now().Format("2006-01-02"), limit)
}

func (s *Service) CountAppointments(ctx *fiber.Ctx, organizationID, doctorID, patientID string) (int, error) {
	// If specific filters are provided, we need to implement specific counting methods
	// For now, use organization-based counting as the primary method
	if organizationID != "" {
		return s.appointmentRepo.CountByOrganization(ctx.Context(), organizationID)
	}
	
	// If no organization filter, count all appointments for today (similar to dashboard behavior)
	return s.appointmentRepo.CountAppointmentsByDate(ctx.Context(), "", time.Now().Format("2006-01-02"))
}

func (s *Service) GetAppointmentByID(ctx *fiber.Ctx, appointmentID string) (*appointment.Appointment, error) {
	return s.appointmentRepo.GetByID(ctx.Context(), appointmentID)
}

func (s *Service) CreateAppointment(ctx *fiber.Ctx, appointmentData *dto.CreateAppointmentRequest) (*appointment.Appointment, error) {
	// Convert DTO to domain entity and create
	// For now, return error - implement conversion logic later
	return nil, fmt.Errorf("not implemented yet")
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
