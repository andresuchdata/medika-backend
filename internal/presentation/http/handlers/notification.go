package handlers

import (
	"strconv"

	"medika-backend/internal/application/notification"
	notificationDomain "medika-backend/internal/domain/notification"
	"medika-backend/internal/domain/shared"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationService *notification.Service
	logger              logger.Logger
}

func NewNotificationHandler(notificationService *notification.Service, logger logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		logger:              logger,
	}
}

// GetNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("user_id").(string)
	userID, _ := shared.NewUserIDFromString(userIDStr)

	// Parse query parameters
	filters := h.parseNotificationFilters(c)

	// Get notifications
	notifications, err := h.notificationService.GetNotificationsByUserID(c.Context(), userID, filters)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get notifications", "error", err, "user_id", userIDStr)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to get notifications",
		})
	}

	// Count total notifications for pagination
	totalCount, err := h.notificationService.CountNotificationsByUserID(c.Context(), userID, filters)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to count notifications", "error", err, "user_id", userIDStr)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to count notifications",
		})
	}

	// Convert to DTOs
	notificationDTOs := make([]dto.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		notificationDTOs[i] = dto.NotificationResponse{
			ID:             notif.ID().String(),
			Title:          notif.Title(),
			Message:        notif.Message(),
			Type:           string(notif.Type()),
			Priority:       string(notif.Priority()),
			IsRead:         notif.IsRead(),
			Channels:       notif.Channels(),
			Data:           notif.Data(),
			ScheduledFor:   notif.ScheduledFor(),
			SentAt:         notif.SentAt(),
			CreatedAt:      notif.CreatedAt(),
		}
	}

	// Calculate pagination info
	limit := filters.Limit
	if limit == 0 {
		limit = 10 // default limit
	}
	offset := filters.Offset
	totalPages := (totalCount + limit - 1) / limit
	currentPage := (offset / limit) + 1

	response := dto.NotificationListResponse{
		Data: notificationDTOs,
		Pagination: dto.Pagination{
			Page:       currentPage,
			TotalPages: totalPages,
			Total:      totalCount,
			Limit:      limit,
		},
	}

	return c.JSON(response)
}

// GetUnreadCount handles GET /api/v1/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("user_id").(string)
	userID, _ := shared.NewUserIDFromString(userIDStr)

	// Create filters for unread notifications only
	filters := notificationDomain.NotificationFilters{
		IsRead: &[]bool{false}[0], // false means unread
	}

	// Count unread notifications
	count, err := h.notificationService.CountNotificationsByUserID(c.Context(), userID, filters)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to count unread notifications", "error", err, "user_id", userIDStr)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to count unread notifications",
		})
	}

	return c.JSON(dto.SuccessResponse{
		Data: map[string]interface{}{
			"unread_count": count,
		},
	})
}

// MarkAsRead handles PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notificationID := notificationDomain.NotificationID(c.Params("id"))

	err := h.notificationService.MarkAsRead(c.Context(), notificationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to mark notification as read", "error", err, "notification_id", notificationID.String())
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to mark notification as read",
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Notification marked as read",
	})
}

// MarkAsUnread handles PUT /api/v1/notifications/:id/unread
func (h *NotificationHandler) MarkAsUnread(c *fiber.Ctx) error {
	notificationID := notificationDomain.NotificationID(c.Params("id"))

	err := h.notificationService.MarkAsUnread(c.Context(), notificationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to mark notification as unread", "error", err, "notification_id", notificationID.String())
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to mark notification as unread",
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Notification marked as unread",
	})
}

// MarkAllAsRead handles PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := shared.NewUserIDFromString(userIDStr)

	err := h.notificationService.MarkAllAsRead(c.Context(), userID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to mark all notifications as read", "error", err, "user_id", userIDStr)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to mark all notifications as read",
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "All notifications marked as read",
	})
}

// DeleteNotification handles DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	notificationID := notificationDomain.NotificationID(c.Params("id"))

	err := h.notificationService.DeleteNotification(c.Context(), notificationID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to delete notification", "error", err, "notification_id", notificationID.String())
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to delete notification",
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Notification deleted",
	})
}

// parseNotificationFilters parses query parameters into notification filters
func (h *NotificationHandler) parseNotificationFilters(c *fiber.Ctx) notificationDomain.NotificationFilters {
	filters := notificationDomain.NotificationFilters{}

	// Parse type filter
	if typeStr := c.Query("type"); typeStr != "" {
		notificationType := notificationDomain.NotificationType(typeStr)
		filters.Type = &notificationType
	}

	// Parse priority filter
	if priorityStr := c.Query("priority"); priorityStr != "" {
		priority := notificationDomain.Priority(priorityStr)
		filters.Priority = &priority
	}

	// Parse read status filter
	if readStr := c.Query("is_read"); readStr != "" {
		if isRead, err := strconv.ParseBool(readStr); err == nil {
			filters.IsRead = &isRead
		}
	}

	// Parse action required filter
	if actionStr := c.Query("action_required"); actionStr != "" {
		if actionRequired, err := strconv.ParseBool(actionStr); err == nil {
			filters.ActionRequired = &actionRequired
		}
	}

	// Parse pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Parse ordering
	if orderBy := c.Query("order_by"); orderBy != "" {
		filters.OrderBy = orderBy
	}

	if order := c.Query("order"); order != "" {
		filters.Order = order
	}

	return filters
}
