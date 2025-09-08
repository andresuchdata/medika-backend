package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/appointment"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type AppointmentHandler struct {
	appointmentService AppointmentService
	validator         *validator.Validate
	logger            logger.Logger
}

// AppointmentService interface for dependency injection
type AppointmentService interface {
	GetAppointments(ctx *fiber.Ctx, organizationID, doctorID, patientID string, limit, offset int) ([]*appointment.Appointment, error)
	CountAppointments(ctx *fiber.Ctx, organizationID, doctorID, patientID string) (int, error)
	GetAppointmentByID(ctx *fiber.Ctx, appointmentID string) (*appointment.Appointment, error)
	CreateAppointment(ctx *fiber.Ctx, appointmentData *dto.CreateAppointmentRequest) (*appointment.Appointment, error)
	UpdateAppointment(ctx *fiber.Ctx, appointmentID string, appointmentData *dto.UpdateAppointmentRequest) (*appointment.Appointment, error)
	DeleteAppointment(ctx *fiber.Ctx, appointmentID string) error
	UpdateAppointmentStatus(ctx *fiber.Ctx, appointmentID string, status string) (*appointment.Appointment, error)
}

func NewAppointmentHandler(
	appointmentService AppointmentService,
	validator *validator.Validate,
	logger logger.Logger,
) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
		validator:         validator,
		logger:            logger,
	}
}

// GET /api/v1/appointments
func (h *AppointmentHandler) GetAppointments(c *fiber.Ctx) error {
	// Get query parameters
	limitStr := c.Query("limit", "10")
	pageStr := c.Query("page", "1")
	organizationID := c.Query("organizationId", "")
	doctorID := c.Query("doctorId", "")
	patientID := c.Query("patientId", "")
	
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
	
	// Get appointments from service
	appointments, err := h.appointmentService.GetAppointments(c, organizationID, doctorID, patientID, limit, offset)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get appointments", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get appointments",
			Message: err.Error(),
		})
	}
	
	// Get total count for pagination
	total, err := h.appointmentService.CountAppointments(c, organizationID, doctorID, patientID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to count appointments", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to count appointments",
			Message: err.Error(),
		})
	}
	
	// Convert domain appointments to DTO appointments
	appointmentResponses := make([]dto.AppointmentResponse, len(appointments))
	for i, apt := range appointments {
		appointmentResponses[i] = dto.AppointmentResponse{
			ID:             apt.ID,
			PatientID:      apt.PatientID,
			PatientName:    "", // TODO: Fetch from patient entity
			DoctorID:       apt.DoctorID,
			DoctorName:     "", // TODO: Fetch from doctor entity
			OrganizationID: apt.OrganizationID,
			RoomID:         apt.RoomID,
			Date:           apt.Date.Format("2006-01-02"),
			StartTime:      apt.StartTime,
			EndTime:        apt.EndTime,
			Duration:       apt.Duration,
			Status:         string(apt.Status),
			Type:           apt.Type,
			Notes:          apt.Notes,
			CreatedAt:      apt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      apt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Build response
	response := dto.AppointmentsResponse{
		Success: true,
		Data: dto.AppointmentsData{
			Appointments: appointmentResponses,
			Pagination: dto.Pagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: (total + limit - 1) / limit,
			},
			Stats: dto.AppointmentStats{
				Total:      total,
				Pending:    countAppointmentsByStatus(appointments, "pending"),
				Confirmed:  countAppointmentsByStatus(appointments, "confirmed"),
				InProgress: countAppointmentsByStatus(appointments, "in_progress"),
				Completed:  countAppointmentsByStatus(appointments, "completed"),
				Cancelled:  countAppointmentsByStatus(appointments, "cancelled"),
				NoShow:     countAppointmentsByStatus(appointments, "no_show"),
			},
		},
		Message: "Appointments retrieved successfully",
	}
	
	return c.JSON(response)
}

// GET /api/v1/appointments/:id
func (h *AppointmentHandler) GetAppointment(c *fiber.Ctx) error {
	appointmentID := c.Params("id")
	if appointmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Appointment ID is required",
			Message: "Please provide a valid appointment ID",
		})
	}

	apt, err := h.appointmentService.GetAppointmentByID(c, appointmentID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get appointment", "error", err, "appointmentID", appointmentID)
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "Appointment not found",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.AppointmentResponse{
			ID:             apt.ID,
			PatientID:      apt.PatientID,
			PatientName:    "", // TODO: Fetch from patient entity
			DoctorID:       apt.DoctorID,
			DoctorName:     "", // TODO: Fetch from doctor entity
			OrganizationID: apt.OrganizationID,
			RoomID:         apt.RoomID,
			Date:           apt.Date.Format("2006-01-02"),
			StartTime:      apt.StartTime,
			EndTime:        apt.EndTime,
			Duration:       apt.Duration,
			Status:         string(apt.Status),
			Type:           apt.Type,
			Notes:          apt.Notes,
			CreatedAt:      apt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      apt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Appointment retrieved successfully",
	}

	return c.JSON(response)
}

// POST /api/v1/appointments
func (h *AppointmentHandler) CreateAppointment(c *fiber.Ctx) error {
	var req dto.CreateAppointmentRequest
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

	apt, err := h.appointmentService.CreateAppointment(c, &req)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to create appointment", "error", err)
		
		// Return 401 if organization ID is missing
		if err.Error() == "organization ID is required" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Organization ID is required",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to create appointment",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.AppointmentResponse{
			ID:             apt.ID,
			PatientID:      apt.PatientID,
			PatientName:    "", // TODO: Fetch from patient entity
			DoctorID:       apt.DoctorID,
			DoctorName:     "", // TODO: Fetch from doctor entity
			OrganizationID: apt.OrganizationID,
			RoomID:         apt.RoomID,
			Date:           apt.Date.Format("2006-01-02"),
			StartTime:      apt.StartTime,
			EndTime:        apt.EndTime,
			Duration:       apt.Duration,
			Status:         string(apt.Status),
			Type:           apt.Type,
			Notes:          apt.Notes,
			CreatedAt:      apt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      apt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Appointment created successfully",
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// PUT /api/v1/appointments/:id
func (h *AppointmentHandler) UpdateAppointment(c *fiber.Ctx) error {
	appointmentID := c.Params("id")
	if appointmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Appointment ID is required",
			Message: "Please provide a valid appointment ID",
		})
	}

	var req dto.UpdateAppointmentRequest
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

	apt, err := h.appointmentService.UpdateAppointment(c, appointmentID, &req)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to update appointment", "error", err, "appointmentID", appointmentID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update appointment",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.AppointmentResponse{
			ID:             apt.ID,
			PatientID:      apt.PatientID,
			PatientName:    "", // TODO: Fetch from patient entity
			DoctorID:       apt.DoctorID,
			DoctorName:     "", // TODO: Fetch from doctor entity
			OrganizationID: apt.OrganizationID,
			RoomID:         apt.RoomID,
			Date:           apt.Date.Format("2006-01-02"),
			StartTime:      apt.StartTime,
			EndTime:        apt.EndTime,
			Duration:       apt.Duration,
			Status:         string(apt.Status),
			Type:           apt.Type,
			Notes:          apt.Notes,
			CreatedAt:      apt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      apt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Appointment updated successfully",
	}

	return c.JSON(response)
}

// DELETE /api/v1/appointments/:id
func (h *AppointmentHandler) DeleteAppointment(c *fiber.Ctx) error {
	appointmentID := c.Params("id")
	if appointmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Appointment ID is required",
			Message: "Please provide a valid appointment ID",
		})
	}

	if err := h.appointmentService.DeleteAppointment(c, appointmentID); err != nil {
		h.logger.Error(c.Context(), "Failed to delete appointment", "error", err, "appointmentID", appointmentID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to delete appointment",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Appointment deleted successfully",
	})
}

// PUT /api/v1/appointments/:id/status
func (h *AppointmentHandler) UpdateAppointmentStatus(c *fiber.Ctx) error {
	appointmentID := c.Params("id")
	if appointmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Appointment ID is required",
			Message: "Please provide a valid appointment ID",
		})
	}

	var req dto.UpdateAppointmentStatusRequest
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

	apt, err := h.appointmentService.UpdateAppointmentStatus(c, appointmentID, req.Status)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to update appointment status", "error", err, "appointmentID", appointmentID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update appointment status",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Data: dto.AppointmentResponse{
			ID:             apt.ID,
			PatientID:      apt.PatientID,
			PatientName:    "", // TODO: Fetch from patient entity
			DoctorID:       apt.DoctorID,
			DoctorName:     "", // TODO: Fetch from doctor entity
			OrganizationID: apt.OrganizationID,
			RoomID:         apt.RoomID,
			Date:           apt.Date.Format("2006-01-02"),
			StartTime:      apt.StartTime,
			EndTime:        apt.EndTime,
			Duration:       apt.Duration,
			Status:         string(apt.Status),
			Type:           apt.Type,
			Notes:          apt.Notes,
			CreatedAt:      apt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      apt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Appointment status updated successfully",
	}

	return c.JSON(response)
}

// Helper functions
func countAppointmentsByStatus(appointments []*appointment.Appointment, status string) int {
	count := 0
	for _, apt := range appointments {
		if string(apt.Status) == status {
			count++
		}
	}
	return count
}
