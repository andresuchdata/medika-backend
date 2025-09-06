package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type DashboardHandler struct {
	dashboardService DashboardService
	validator        *validator.Validate
	logger           logger.Logger
}

// DashboardService interface for dependency injection
type DashboardService interface {
	GetDashboardSummary(ctx context.Context, organizationID string) (*dto.DashboardSummaryResponse, error)
}

func NewDashboardHandler(
	dashboardService DashboardService,
	validator *validator.Validate,
	logger logger.Logger,
) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		validator:        validator,
		logger:           logger,
	}
}

// GET /api/v1/dashboard/summary
func (h *DashboardHandler) GetDashboardSummary(c *fiber.Ctx) error {
	// Get organization ID from query parameter or user context
	organizationID := c.Query("organizationId", "")
	
	// Validate organizationId if provided
	if organizationID != "" {
		if err := h.validator.Var(organizationID, "uuid"); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error:   "Invalid organization ID",
				Message: "Organization ID must be a valid UUID",
			})
		}
	}

	// Get dashboard summary
	summary, err := h.dashboardService.GetDashboardSummary(c.Context(), organizationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get dashboard summary", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get dashboard summary",
			Message: err.Error(),
		})
	}

	return c.JSON(summary)
}
