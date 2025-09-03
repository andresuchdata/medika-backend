package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/organization"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type OrganizationHandler struct {
	organizationService OrganizationService
	validator          *validator.Validate
	logger             logger.Logger
}

// OrganizationService interface for dependency injection
type OrganizationService interface {
	GetOrganizations(ctx *fiber.Ctx, limit, offset int) ([]*organization.Organization, error)
	CountOrganizations(ctx *fiber.Ctx) (int, error)
	GetOrganizationByID(ctx *fiber.Ctx, organizationID string) (*organization.Organization, error)
	CreateOrganization(ctx *fiber.Ctx, orgData *dto.CreateOrganizationRequest) (*organization.Organization, error)
	UpdateOrganization(ctx *fiber.Ctx, organizationID string, orgData *dto.UpdateOrganizationRequest) (*organization.Organization, error)
	DeleteOrganization(ctx *fiber.Ctx, organizationID string) error
}

func NewOrganizationHandler(
	organizationService OrganizationService,
	validator *validator.Validate,
	logger logger.Logger,
) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: organizationService,
		validator:          validator,
		logger:             logger,
	}
}

// GET /api/v1/organizations
func (h *OrganizationHandler) GetOrganizations(c *fiber.Ctx) error {
	// Get query parameters
	limitStr := c.Query("limit", "10")
	pageStr := c.Query("page", "1")
	
	// Parse pagination parameters
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	
	offset := (page - 1) * limit
	
	// Get organizations from service
	organizations, err := h.organizationService.GetOrganizations(c, limit, offset)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get organizations", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get organizations",
			Message: err.Error(),
		})
	}
	
	// Get total count for pagination
	total, err := h.organizationService.CountOrganizations(c)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to count organizations", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to count organizations",
			Message: err.Error(),
		})
	}
	
	// Convert domain organizations to DTO organizations
	orgResponses := make([]dto.OrganizationResponse, len(organizations))
	for i, org := range organizations {
		orgResponses[i] = dto.OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			Type:      org.Type,
			Address:   org.Address,
			Phone:     org.Phone,
			Email:     org.Email,
			Website:   org.Website,
			Status:    getOrganizationStatus(org.IsActive),
			StaffCount: 0, // TODO: Implement staff count calculation
			CreatedAt: org.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: org.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Build response
	response := dto.OrganizationsResponse{
		Success: true,
		Data: dto.OrganizationsData{
			Organizations: orgResponses,
			Pagination: dto.Pagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: (total + limit - 1) / limit,
			},
			Stats: dto.OrganizationStats{
				Total:      total,
				Active:     countOrganizationsByStatus(organizations, true),
				Inactive:   countOrganizationsByStatus(organizations, false),
				Types:      buildTypeStats(organizations),
				TotalStaff: calculateTotalStaff(organizations),
			},
		},
		Message: "Organizations retrieved successfully",
	}
	
	return c.JSON(response)
}

// GET /api/v1/organizations/:id
func (h *OrganizationHandler) GetOrganization(c *fiber.Ctx) error {
	organizationID := c.Params("id")
	if organizationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Organization ID is required",
			Message: "Please provide a valid organization ID",
		})
	}

	org, err := h.organizationService.GetOrganizationByID(c, organizationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get organization", "error", err, "organizationID", organizationID)
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "Organization not found",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			Type:      org.Type,
			Address:   org.Address,
			Phone:     org.Phone,
			Email:     org.Email,
			Website:   org.Website,
			Status:    getOrganizationStatus(org.IsActive),
			StaffCount: 0, // TODO: Implement staff count calculation
			CreatedAt: org.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: org.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Organization retrieved successfully",
	}

	return c.JSON(response)
}

// POST /api/v1/organizations
func (h *OrganizationHandler) CreateOrganization(c *fiber.Ctx) error {
	var req dto.CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid JSON format",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	org, err := h.organizationService.CreateOrganization(c, &req)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to create organization", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to create organization",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			Type:      org.Type,
			Address:   org.Address,
			Phone:     org.Phone,
			Email:     org.Email,
			Website:   org.Website,
			Status:    getOrganizationStatus(org.IsActive),
			StaffCount: 0, // TODO: Implement staff count calculation
			CreatedAt: org.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: org.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Organization created successfully",
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// PUT /api/v1/organizations/:id
func (h *OrganizationHandler) UpdateOrganization(c *fiber.Ctx) error {
	organizationID := c.Params("id")
	if organizationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Organization ID is required",
			Message: "Please provide a valid organization ID",
		})
	}

	var req dto.UpdateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid JSON format",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	org, err := h.organizationService.UpdateOrganization(c, organizationID, &req)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to update organization", "error", err, "organizationID", organizationID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update organization",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			Type:      org.Type,
			Address:   org.Address,
			Phone:     org.Phone,
			Email:     org.Email,
			Website:   org.Website,
			Status:    getOrganizationStatus(org.IsActive),
			StaffCount: 0, // TODO: Implement staff count calculation
			CreatedAt: org.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: org.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Organization updated successfully",
	}

	return c.JSON(response)
}

// DELETE /api/v1/organizations/:id
func (h *OrganizationHandler) DeleteOrganization(c *fiber.Ctx) error {
	organizationID := c.Params("id")
	if organizationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Organization ID is required",
			Message: "Please provide a valid organization ID",
		})
	}

	if err := h.organizationService.DeleteOrganization(c, organizationID); err != nil {
		h.logger.Error(c.Context(), "Failed to delete organization", "error", err, "organizationID", organizationID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to delete organization",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Organization deleted successfully",
	})
}

// Helper functions
func countOrganizationsByStatus(organizations []*organization.Organization, status bool) int {
	count := 0
	for _, org := range organizations {
		if org.IsActive == status {
			count++
		}
	}
	return count
}

func getOrganizationStatus(isActive bool) string {
	if isActive {
		return "active"
	}
	return "inactive"
}

func buildTypeStats(organizations []*organization.Organization) map[string]int {
	stats := make(map[string]int)
	for _, org := range organizations {
		orgType := org.Type
		stats[orgType]++
	}

	return stats
}

func calculateTotalStaff(_ []*organization.Organization) int {
	// TODO: Implement staff count calculation
	return 0
}
