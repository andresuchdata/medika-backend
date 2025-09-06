package handlers

import (
	"context"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"medika-backend/internal/domain/queue"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type QueueHandler struct {
	queueService QueueService
	validator    *validator.Validate
	logger       logger.Logger
}

// QueueService interface for dependency injection
type QueueService interface {
	GetQueuesByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*queue.PatientQueue, error)
	CountQueuesByOrganization(ctx context.Context, organizationID string) (int, error)
	CreateQueue(ctx context.Context, q *queue.PatientQueue) error
	GetQueue(ctx context.Context, id string) (*queue.PatientQueue, error)
	UpdateQueue(ctx context.Context, q *queue.PatientQueue) error
	DeleteQueue(ctx context.Context, id string) error
	CallNextPatient(ctx context.Context, organizationID string) (*queue.PatientQueue, error)
	StartPatientConsultation(ctx context.Context, queueID string) (*queue.PatientQueue, error)
	CompletePatientConsultation(ctx context.Context, queueID string) (*queue.PatientQueue, error)
}

func NewQueueHandler(
	queueService QueueService,
	validator *validator.Validate,
	logger logger.Logger,
) *QueueHandler {
	return &QueueHandler{
		queueService: queueService,
		validator:    validator,
		logger:      logger,
	}
}

// GET /api/v1/queues
func (h *QueueHandler) GetQueues(c *fiber.Ctx) error {
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
	
	// Get queues from service
	queues, err := h.queueService.GetQueuesByOrganization(c.Context(), organizationID, limit, offset)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get queues", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get queues",
			Message: err.Error(),
		})
	}
	
	// Get total count for pagination
	total, err := h.queueService.CountQueuesByOrganization(c.Context(), organizationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to count queues", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to count queues",
			Message: err.Error(),
		})
	}
	
	// Convert domain queues to DTO queues
	queueResponses := make([]dto.QueueResponse, len(queues))
	for i, q := range queues {
		queueResponses[i] = dto.QueueResponse{
			ID:                q.ID,
			AppointmentID:     q.AppointmentID,
			OrganizationID:    q.OrganizationID,
			Position:          q.Position,
			EstimatedWaitTime: q.EstimatedWaitTime,
			Status:            q.Status,
			CreatedAt:         q.CreatedAt,
			UpdatedAt:         q.UpdatedAt,
		}
	}

	// Build response
	response := dto.QueuesResponse{
		Success: true,
		Data: struct {
			Queues     []dto.QueueResponse `json:"queues"`
			Pagination dto.Pagination      `json:"pagination"`
		}{
			Queues: queueResponses,
			Pagination: dto.Pagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: (total + limit - 1) / limit,
			},
		},
		Message: "Queues retrieved successfully",
	}
	
	return c.JSON(response)
}

// GET /api/v1/queues/:id
func (h *QueueHandler) GetQueue(c *fiber.Ctx) error {
	id := c.Params("id")
	
	// Validate ID
	if err := h.validator.Var(id, "uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid queue ID",
			Message: "Queue ID must be a valid UUID",
		})
	}
	
	// Get queue from service
	queue, err := h.queueService.GetQueue(c.Context(), id)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get queue", "error", err)
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "Queue not found",
			Message: err.Error(),
		})
	}
	
	// Convert to DTO
	queueResponse := dto.QueueResponse{
		ID:                queue.ID,
		AppointmentID:     queue.AppointmentID,
		OrganizationID:    queue.OrganizationID,
		Position:          queue.Position,
		EstimatedWaitTime: queue.EstimatedWaitTime,
		Status:            queue.Status,
		CreatedAt:         queue.CreatedAt,
		UpdatedAt:         queue.UpdatedAt,
	}
	
	response := dto.QueueDetailResponse{
		Success: true,
		Data:    queueResponse,
		Message: "Queue retrieved successfully",
	}
	
	return c.JSON(response)
}

// POST /api/v1/queues
func (h *QueueHandler) CreateQueue(c *fiber.Ctx) error {
	var request dto.QueueRequest
	
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}
	
	// Validate request
	if err := h.validator.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}
	
	// Convert to domain model
	queue := &queue.PatientQueue{
		AppointmentID:     request.AppointmentID,
		OrganizationID:    request.OrganizationID,
		EstimatedWaitTime: request.EstimatedWaitTime,
	}
	
	// Create queue
	if err := h.queueService.CreateQueue(c.Context(), queue); err != nil {
		h.logger.Error(c.Context(), "Failed to create queue", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to create queue",
			Message: err.Error(),
		})
	}
	
	// Convert to DTO
	queueResponse := dto.QueueResponse{
		ID:                queue.ID,
		AppointmentID:     queue.AppointmentID,
		OrganizationID:    queue.OrganizationID,
		Position:          queue.Position,
		EstimatedWaitTime: queue.EstimatedWaitTime,
		Status:            queue.Status,
		CreatedAt:         queue.CreatedAt,
		UpdatedAt:         queue.UpdatedAt,
	}
	
	response := dto.QueueDetailResponse{
		Success: true,
		Data:    queueResponse,
		Message: "Queue created successfully",
	}
	
	return c.Status(fiber.StatusCreated).JSON(response)
}

// PUT /api/v1/queues/:id
func (h *QueueHandler) UpdateQueue(c *fiber.Ctx) error {
	queueID := c.Params("id")
	if queueID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid queue ID",
			Message: "Queue ID is required",
		})
	}

	var req dto.QueueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	// Get existing queue
	existingQueue, err := h.queueService.GetQueue(c.Context(), queueID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "Queue not found",
			Message: err.Error(),
		})
	}

	// Update queue fields
	existingQueue.Position = req.Position
	existingQueue.EstimatedWaitTime = req.EstimatedWaitTime
	existingQueue.Status = req.Status

	// Update queue
	if err := h.queueService.UpdateQueue(c.Context(), existingQueue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update queue",
			Message: err.Error(),
		})
	}

	// Convert to response DTO
	queueResponse := dto.QueueResponse{
		ID:                existingQueue.ID,
		AppointmentID:     existingQueue.AppointmentID,
		OrganizationID:    existingQueue.OrganizationID,
		Position:          existingQueue.Position,
		EstimatedWaitTime: existingQueue.EstimatedWaitTime,
		Status:            existingQueue.Status,
		CreatedAt:         existingQueue.CreatedAt,
		UpdatedAt:         existingQueue.UpdatedAt,
	}

	response := dto.QueueDetailResponse{
		Success: true,
		Data:    queueResponse,
		Message: "Queue updated successfully",
	}

	return c.JSON(response)
}

// DELETE /api/v1/queues/:id
func (h *QueueHandler) DeleteQueue(c *fiber.Ctx) error {
	queueID := c.Params("id")
	if queueID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid queue ID",
			Message: "Queue ID is required",
		})
	}

	// Check if queue exists
	_, err := h.queueService.GetQueue(c.Context(), queueID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "Queue not found",
			Message: err.Error(),
		})
	}

	// Delete queue
	if err := h.queueService.DeleteQueue(c.Context(), queueID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to delete queue",
			Message: err.Error(),
		})
	}

	response := dto.SuccessResponse{
		Success: true,
		Message: "Queue deleted successfully",
	}

	return c.JSON(response)
}

// POST /api/v1/queues/:id/actions
func (h *QueueHandler) QueueAction(c *fiber.Ctx) error {
	queueID := c.Params("id")
	if queueID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid queue ID",
			Message: "Queue ID is required",
		})
	}

	var req dto.QueueActionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	var updatedQueue *queue.PatientQueue
	var err error

	// Perform action based on type
	switch req.Action {
	case "call_next":
		// Get organization ID from queue
		existingQueue, getErr := h.queueService.GetQueue(c.Context(), queueID)
		if getErr != nil {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error:   "Queue not found",
				Message: getErr.Error(),
			})
		}
		updatedQueue, err = h.queueService.CallNextPatient(c.Context(), existingQueue.OrganizationID)
	case "start_consultation":
		updatedQueue, err = h.queueService.StartPatientConsultation(c.Context(), queueID)
	case "complete_consultation":
		updatedQueue, err = h.queueService.CompletePatientConsultation(c.Context(), queueID)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid action",
			Message: "Action must be one of: call_next, start_consultation, complete_consultation",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to perform action",
			Message: err.Error(),
		})
	}

	// Convert to response DTO
	queueResponse := dto.QueueResponse{
		ID:                updatedQueue.ID,
		AppointmentID:     updatedQueue.AppointmentID,
		OrganizationID:    updatedQueue.OrganizationID,
		Position:          updatedQueue.Position,
		EstimatedWaitTime: updatedQueue.EstimatedWaitTime,
		Status:            updatedQueue.Status,
		CreatedAt:         updatedQueue.CreatedAt,
		UpdatedAt:         updatedQueue.UpdatedAt,
	}

	response := dto.QueueDetailResponse{
		Success: true,
		Data:    queueResponse,
		Message: "Action performed successfully",
	}

	return c.JSON(response)
}
