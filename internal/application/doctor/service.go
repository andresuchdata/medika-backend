package doctor

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/doctor"
	"medika-backend/internal/domain/user"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type Service struct {
	doctorRepo doctor.Repository
	logger     logger.Logger
}

func NewService(doctorRepo doctor.Repository, logger logger.Logger) *Service {
	return &Service{
		doctorRepo: doctorRepo,
		logger:     logger,
	}
}

// Methods that match the DoctorService interface
func (s *Service) GetDoctorsByOrganization(ctx *fiber.Ctx, organizationID string, limit, offset int) ([]*user.User, error) {
	// For now, return empty list - implement conversion logic later
	return []*user.User{}, nil
}

func (s *Service) CountDoctorsByOrganization(ctx *fiber.Ctx, organizationID string) (int, error) {
	return s.doctorRepo.CountByOrganization(ctx.Context(), organizationID)
}

func (s *Service) GetDoctorByID(ctx *fiber.Ctx, doctorID string) (*user.User, error) {
	// For now, return error - implement conversion logic later
	return nil, fmt.Errorf("not implemented yet")
}

func (s *Service) CreateDoctor(ctx *fiber.Ctx, doctorData *dto.CreateDoctorRequest) (*user.User, error) {
	// For now, return error - implement creation logic later
	return nil, fmt.Errorf("not implemented yet")
}

func (s *Service) UpdateDoctor(ctx *fiber.Ctx, doctorID string, doctorData *dto.UpdateDoctorRequest) (*user.User, error) {
	// For now, return error - implement update logic later
	return nil, fmt.Errorf("not implemented yet")
}

func (s *Service) DeleteDoctor(ctx *fiber.Ctx, doctorID string) error {
	return s.doctorRepo.Delete(ctx.Context(), doctorID)
}

// Additional methods for future use
func (s *Service) CreateDoctorDirect(ctx context.Context, d *doctor.Doctor) error {
	return s.doctorRepo.Create(ctx, d)
}

func (s *Service) GetDoctorDirect(ctx context.Context, id string) (*doctor.Doctor, error) {
	return s.doctorRepo.GetByID(ctx, id)
}

func (s *Service) UpdateDoctorDirect(ctx context.Context, d *doctor.Doctor) error {
	return s.doctorRepo.Update(ctx, d)
}

func (s *Service) GetDoctorsByOrganizationDirect(ctx context.Context, organizationID string, limit, offset int) ([]*doctor.Doctor, error) {
	return s.doctorRepo.GetByOrganization(ctx, organizationID, limit, offset)
}
