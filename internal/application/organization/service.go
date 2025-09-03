package organization

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/organization"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type Service struct {
	orgRepo organization.Repository
	logger  logger.Logger
}

func NewService(orgRepo organization.Repository, logger logger.Logger) *Service {
	return &Service{
		orgRepo: orgRepo,
		logger:  logger,
	}
}

// Methods that match the OrganizationService interface
func (s *Service) GetOrganizations(ctx *fiber.Ctx, limit, offset int) ([]*organization.Organization, error) {
	return s.orgRepo.GetAll(ctx.Context(), limit, offset)
}

func (s *Service) CountOrganizations(ctx *fiber.Ctx) (int, error) {
	return s.orgRepo.Count(ctx.Context())
}

func (s *Service) GetOrganizationByID(ctx *fiber.Ctx, organizationID string) (*organization.Organization, error) {
	return s.orgRepo.GetByID(ctx.Context(), organizationID)
}

func (s *Service) CreateOrganization(ctx *fiber.Ctx, orgData *dto.CreateOrganizationRequest) (*organization.Organization, error) {
	// For now, return error - implement conversion logic later
	return nil, fmt.Errorf("not implemented yet")
}

func (s *Service) UpdateOrganization(ctx *fiber.Ctx, organizationID string, orgData *dto.UpdateOrganizationRequest) (*organization.Organization, error) {
	// For now, return error - implement update logic later
	return nil, fmt.Errorf("not implemented yet")
}

func (s *Service) DeleteOrganization(ctx *fiber.Ctx, organizationID string) error {
	return s.orgRepo.Delete(ctx.Context(), organizationID)
}

// Additional methods for future use
func (s *Service) CreateOrganizationDirect(ctx context.Context, org *organization.Organization) error {
	return s.orgRepo.Create(ctx, org)
}

func (s *Service) UpdateOrganizationDirect(ctx context.Context, org *organization.Organization) error {
	return s.orgRepo.Update(ctx, org)
}
