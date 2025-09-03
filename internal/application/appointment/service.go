package appointment

import (
	"context"
	"fmt"

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
	// For now, return empty list - implement filtering logic later
	return []*appointment.Appointment{}, nil
}

func (s *Service) CountAppointments(ctx *fiber.Ctx, organizationID, doctorID, patientID string) (int, error) {
	// For now, return 0 - implement counting logic later
	return 0, nil
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
