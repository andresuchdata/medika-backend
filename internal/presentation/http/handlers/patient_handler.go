package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/patient"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type PatientHandler struct {
	patientService PatientService
	validator      *validator.Validate
	logger         logger.Logger
}

// PatientService interface for dependency injection
type PatientService interface {
	GetPatientsByOrganization(ctx *fiber.Ctx, organizationID string, limit, offset int) ([]*patient.Patient, error)
	CountPatientsByOrganization(ctx *fiber.Ctx, organizationID string) (int, error)
}

func NewPatientHandler(
	patientService PatientService,
	validator *validator.Validate,
	logger logger.Logger,
) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
		validator:      validator,
		logger:         logger,
	}
}

// GET /api/v1/patients
func (h *PatientHandler) GetPatients(c *fiber.Ctx) error {
	// Get query parameters
	limitStr := c.Query("limit", "10")
	pageStr := c.Query("page", "1")
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
	
	// Get patients from service
	patients, err := h.patientService.GetPatientsByOrganization(c, organizationID, limit, offset)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get patients", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get patients",
			Message: err.Error(),
		})
	}
	
	// Get total count for pagination
	total, err := h.patientService.CountPatientsByOrganization(c, organizationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to count patients", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to count patients",
			Message: err.Error(),
		})
	}
	
	// Convert domain patients to DTO patients
	patientResponses := make([]dto.PatientResponse, len(patients))
	for i, p := range patients {
		patientResponses[i] = dto.PatientResponse{
			ID:        p.ID,
			Name:      p.Name,
			Email:     p.Email,
			Phone:     p.Phone,
			DateOfBirth: p.DateOfBirth,
			Age:       p.Age,
			Gender:    p.Gender,
			Avatar:    p.Avatar,
			Address: dto.AddressResponse{
				Street:  p.Address.Street,
				City:    p.Address.City,
				State:   p.Address.State,
				ZipCode: p.Address.ZipCode,
				Country: p.Address.Country,
			},
			EmergencyContact: dto.EmergencyContactResponse{
				Name:         p.EmergencyContact.Name,
				Relationship: p.EmergencyContact.Relationship,
				Phone:        p.EmergencyContact.Phone,
			},
			MedicalHistory: []dto.MedicalConditionResponse{},
			Allergies:      p.Allergies,
			Medications:    []dto.MedicationResponse{},
			LastVisit:      p.LastVisit,
			NextAppointment: p.NextAppointment,
			Status:         p.Status,
			OrganizationID: p.OrganizationID,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		}
	}

	// Build response
	response := dto.PatientsResponse{
		Success: true,
		Data: dto.PatientsData{
			Patients: patientResponses,
			Pagination: dto.Pagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: (total + limit - 1) / limit,
			},
		},
		Message: "Patients retrieved successfully",
	}
	
	return c.JSON(response)
}
