package doctor

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/doctor"
	"medika-backend/internal/domain/shared"
	"medika-backend/internal/domain/user"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

// UserRepository interface for dependency injection
type UserRepository interface {
	FindAll(ctx context.Context, filters user.UserFilters) ([]*user.User, error)
	Count(ctx context.Context, filters user.UserFilters) (int64, error)
}

type Service struct {
	doctorRepo doctor.Repository
	userRepo   UserRepository
	logger     logger.Logger
}

func NewService(doctorRepo doctor.Repository, userRepo UserRepository, logger logger.Logger) *Service {
	return &Service{
		doctorRepo: doctorRepo,
		userRepo:   userRepo,
		logger:     logger,
	}
}

// Methods that match the DoctorService interface
func (s *Service) GetDoctorsByOrganization(ctx *fiber.Ctx, organizationID string, limit, offset int) ([]*user.User, error) {
	// Build filters for doctors
	filters := user.UserFilters{
		Role:   &[]user.Role{user.RoleDoctor}[0], // Convert to pointer
		Limit:  limit,
		Offset: offset,
		OrderBy: "created_at",
		Order:   "DESC",
	}

	// Add organization filter if provided
	if organizationID != "" {
		orgID, err := shared.NewOrganizationID(organizationID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization ID: %w", err)
		}
		filters.OrganizationID = &orgID
	}

	// Get doctors using unified FindAll method
	doctors, err := s.userRepo.FindAll(ctx.Context(), filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctors: %w", err)
	}

	return doctors, nil
}

func (s *Service) CountDoctorsByOrganization(ctx *fiber.Ctx, organizationID string) (int, error) {
	// Build filters for doctors count
	filters := user.UserFilters{
		Role: &[]user.Role{user.RoleDoctor}[0], // Convert to pointer
	}

	// Add organization filter if provided
	if organizationID != "" {
		orgID, err := shared.NewOrganizationID(organizationID)
		if err != nil {
			return 0, fmt.Errorf("invalid organization ID: %w", err)
		}
		filters.OrganizationID = &orgID
	}

	// Get count using unified Count method
	count, err := s.userRepo.Count(ctx.Context(), filters)
	if err != nil {
		return 0, fmt.Errorf("failed to count doctors: %w", err)
	}

	return int(count), nil
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
