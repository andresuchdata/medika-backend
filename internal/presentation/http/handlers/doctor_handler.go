package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/shared"
	"medika-backend/internal/domain/user"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type DoctorHandler struct {
	doctorService DoctorService
	validator     *validator.Validate
	logger        logger.Logger
}

// DoctorService interface for dependency injection
type DoctorService interface {
	GetDoctorsByOrganization(ctx *fiber.Ctx, organizationID string, limit, offset int) ([]*user.User, error)
	CountDoctorsByOrganization(ctx *fiber.Ctx, organizationID string) (int, error)
	GetDoctorByID(ctx *fiber.Ctx, doctorID string) (*user.User, error)
	CreateDoctor(ctx *fiber.Ctx, doctorData *dto.CreateDoctorRequest) (*user.User, error)
	UpdateDoctor(ctx *fiber.Ctx, doctorID string, doctorData *dto.UpdateDoctorRequest) (*user.User, error)
	DeleteDoctor(ctx *fiber.Ctx, doctorID string) error
}

func NewDoctorHandler(
	doctorService DoctorService,
	validator *validator.Validate,
	logger logger.Logger,
) *DoctorHandler {
	return &DoctorHandler{
		doctorService: doctorService,
		validator:     validator,
		logger:        logger,
	}
}

// GET /api/v1/doctors
func (h *DoctorHandler) GetDoctors(c *fiber.Ctx) error {
	// Get query parameters
	limitStr := c.Query("limit", "10")
	pageStr := c.Query("page", "1")
	organizationID := c.Query("organizationId", "")
	
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
	
	// Get doctors from service
	doctors, err := h.doctorService.GetDoctorsByOrganization(c, organizationID, limit, offset)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get doctors", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get doctors",
			Message: err.Error(),
		})
	}
	
	// Get total count for pagination
	total, err := h.doctorService.CountDoctorsByOrganization(c, organizationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to count doctors", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to count doctors",
			Message: err.Error(),
		})
	}
	
	// Convert domain doctors to DTO doctors
	doctorResponses := make([]dto.DoctorResponse, len(doctors))
	for i, d := range doctors {
		doctorResponses[i] = dto.DoctorResponse{
			ID:             d.ID().String(),
			Name:           d.Name().String(),
			Email:          d.Email().String(),
			Phone:          getPhoneValue(d.Phone()),
			Specialty:      getStringValue(d.Profile().Specialty()),
			LicenseNumber:  getStringValue(d.Profile().LicenseNumber()),
			Status:         d.Status(),
			OrganizationID: getOrganizationIDString(d.OrganizationID()),
			Avatar:         d.AvatarURL(),
			Bio:            d.Profile().Bio(),
			Experience:     d.Profile().Experience(),
			Education:      d.Profile().Education(),
			Certifications: d.Profile().Certifications(),
			NextAvailable:  d.Profile().NextAvailable(),
			CreatedAt:      d.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      d.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Build response
	response := dto.DoctorsResponse{
		Success: true,
		Data: dto.DoctorsData{
			Doctors: doctorResponses,
			Pagination: dto.Pagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: (total + limit - 1) / limit,
			},
			Stats: dto.DoctorStats{
				Total:      total,
				Active:     countByStatus(doctors, "active"),
				Inactive:   countByStatus(doctors, "inactive"),
				OnLeave:    countByStatus(doctors, "on-leave"),
				Specialties: buildSpecialtyStats(doctors),
			},
		},
		Message: "Doctors retrieved successfully",
	}
	
	return c.JSON(response)
}

// GET /api/v1/doctors/:id
func (h *DoctorHandler) GetDoctor(c *fiber.Ctx) error {
	doctorID := c.Params("id")
	if doctorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Doctor ID is required",
			Message: "Please provide a valid doctor ID",
		})
	}

	doctor, err := h.doctorService.GetDoctorByID(c, doctorID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get doctor", "error", err, "doctorID", doctorID)
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "Doctor not found",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.DoctorResponse{
			ID:             doctor.ID().String(),
			Name:           doctor.Name().String(),
			Email:          doctor.Email().String(),
			Phone:          getPhoneValue(doctor.Phone()),
			Specialty:      getStringValue(doctor.Profile().Specialty()),
			LicenseNumber:  getStringValue(doctor.Profile().LicenseNumber()),
			Status:         doctor.Status(),
			OrganizationID: getOrganizationIDString(doctor.OrganizationID()),
			Avatar:         doctor.AvatarURL(),
			Bio:            doctor.Profile().Bio(),
			Experience:     doctor.Profile().Experience(),
			Education:      doctor.Profile().Education(),
			Certifications: doctor.Profile().Certifications(),
			NextAvailable:  doctor.Profile().NextAvailable(),
			CreatedAt:      doctor.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      doctor.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Doctor retrieved successfully",
	}

	return c.JSON(response)
}

// POST /api/v1/doctors
func (h *DoctorHandler) CreateDoctor(c *fiber.Ctx) error {
	var req dto.CreateDoctorRequest
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

	doctor, err := h.doctorService.CreateDoctor(c, &req)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to create doctor", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to create doctor",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.DoctorResponse{
			ID:             doctor.ID().String(),
			Name:           doctor.Name().String(),
			Email:          doctor.Email().String(),
			Phone:          getPhoneValue(doctor.Phone()),
			Specialty:      getStringValue(doctor.Profile().Specialty()),
			LicenseNumber:  getStringValue(doctor.Profile().LicenseNumber()),
			Status:         doctor.Status(),
			OrganizationID: getOrganizationIDString(doctor.OrganizationID()),
			Avatar:         doctor.AvatarURL(),
			Bio:            doctor.Profile().Bio(),
			Experience:     doctor.Profile().Experience(),
			Education:      doctor.Profile().Education(),
			Certifications: doctor.Profile().Certifications(),
			NextAvailable:  doctor.Profile().NextAvailable(),
			CreatedAt:      doctor.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      doctor.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Doctor created successfully",
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// PUT /api/v1/doctors/:id
func (h *DoctorHandler) UpdateDoctor(c *fiber.Ctx) error {
	doctorID := c.Params("id")
	if doctorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Doctor ID is required",
			Message: "Please provide a valid doctor ID",
		})
	}

	var req dto.UpdateDoctorRequest
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

	doctor, err := h.doctorService.UpdateDoctor(c, doctorID, &req)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to update doctor", "error", err, "doctorID", doctorID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update doctor",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.DoctorResponse{
			ID:             doctor.ID().String(),
			Name:           doctor.Name().String(),
			Email:          doctor.Email().String(),
			Phone:          getPhoneValue(doctor.Phone()),
			Specialty:      getStringValue(doctor.Profile().Specialty()),
			LicenseNumber:  getStringValue(doctor.Profile().LicenseNumber()),
			Status:         doctor.Status(),
			OrganizationID: getOrganizationIDString(doctor.OrganizationID()),
			Avatar:         doctor.AvatarURL(),
			Bio:            doctor.Profile().Bio(),
			Experience:     doctor.Profile().Experience(),
			Education:      doctor.Profile().Education(),
			Certifications: doctor.Profile().Certifications(),
			NextAvailable:  doctor.Profile().NextAvailable(),
			CreatedAt:      doctor.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      doctor.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Doctor updated successfully",
	}

	return c.JSON(response)
}

// DELETE /api/v1/doctors/:id
func (h *DoctorHandler) DeleteDoctor(c *fiber.Ctx) error {
	doctorID := c.Params("id")
	if doctorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Doctor ID is required",
			Message: "Please provide a valid doctor ID",
		})
	}

	if err := h.doctorService.DeleteDoctor(c, doctorID); err != nil {
		h.logger.Error(c.Context(), "Failed to delete doctor", "error", err, "doctorID", doctorID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to delete doctor",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Doctor deleted successfully",
	})
}

// Helper functions
func countByStatus(doctors []*user.User, status string) int {
	count := 0
	for _, d := range doctors {
		if d.Status() == status {
			count++
		}
	}
	return count
}

func buildSpecialtyStats(doctors []*user.User) map[string]int {
	stats := make(map[string]int)
	for _, d := range doctors {
		specialty := d.Profile().Specialty()
		if specialty != nil && *specialty != "" {
			stats[*specialty]++
		}
	}
	return stats
}

func getStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func getPhoneValue(phone *shared.PhoneNumber) string {
	if phone == nil {
		return ""
	}
	return phone.String()
}

func getOrganizationIDString(orgID *shared.OrganizationID) string {
	if orgID == nil {
		return ""
	}
	return orgID.String()
}
