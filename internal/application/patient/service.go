package patient

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/patient"
	"medika-backend/pkg/logger"
)

type Service struct {
	patientRepo patient.Repository
	logger      logger.Logger
}

func NewService(patientRepo patient.Repository, logger logger.Logger) *Service {
	return &Service{
		patientRepo: patientRepo,
		logger:      logger,
	}
}

// Methods that match the PatientService interface
func (s *Service) GetPatientsByOrganization(ctx *fiber.Ctx, organizationID string, limit, offset int) ([]*patient.Patient, error) {
	return s.patientRepo.GetByOrganization(ctx.Context(), organizationID, limit, offset)
}

func (s *Service) CountPatientsByOrganization(ctx *fiber.Ctx, organizationID string) (int, error) {
	return s.patientRepo.CountByOrganization(ctx.Context(), organizationID)
}

// Additional methods for future use
func (s *Service) CreatePatient(ctx context.Context, p *patient.Patient) error {
	return s.patientRepo.Create(ctx, p)
}

func (s *Service) GetPatient(ctx context.Context, id string) (*patient.Patient, error) {
	return s.patientRepo.GetByID(ctx, id)
}

func (s *Service) UpdatePatient(ctx context.Context, p *patient.Patient) error {
	return s.patientRepo.Update(ctx, p)
}

func (s *Service) DeletePatient(ctx context.Context, id string) error {
	return s.patientRepo.Delete(ctx, id)
}
