package notification

import (
	"context"
	"medika-backend/internal/domain/notification"
	"medika-backend/internal/domain/shared"
	"medika-backend/pkg/logger"
)

type Service struct {
	notificationRepo notification.Repository
	logger           logger.Logger
}

func NewService(notificationRepo notification.Repository, logger logger.Logger) *Service {
	return &Service{
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

// GetNotificationsByUserID retrieves notifications for a user with filters
func (s *Service) GetNotificationsByUserID(ctx context.Context, userID shared.UserID, filters notification.NotificationFilters) ([]*notification.Notification, error) {
	s.logger.Info(ctx, "Getting notifications for user", "user_id", userID.String())
	
	notifications, err := s.notificationRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		s.logger.Error(ctx, "Failed to get notifications", "error", err, "user_id", userID.String())
		return nil, err
	}
	
	s.logger.Info(ctx, "Successfully retrieved notifications", "count", len(notifications), "user_id", userID.String())
	return notifications, nil
}

// CountNotificationsByUserID counts notifications for a user with filters
func (s *Service) CountNotificationsByUserID(ctx context.Context, userID shared.UserID, filters notification.NotificationFilters) (int, error) {
	s.logger.Info(ctx, "Counting notifications for user", "user_id", userID.String())
	
	count, err := s.notificationRepo.CountByUserID(ctx, userID, filters)
	if err != nil {
		s.logger.Error(ctx, "Failed to count notifications", "error", err, "user_id", userID.String())
		return 0, err
	}
	
	s.logger.Info(ctx, "Successfully counted notifications", "count", count, "user_id", userID.String())
	return count, nil
}

// CreateNotification creates a new notification
func (s *Service) CreateNotification(
	ctx context.Context,
	userID shared.UserID,
	title, message string,
	notificationType notification.NotificationType,
	priority notification.Priority,
	channels []string,
	data map[string]interface{},
) (*notification.Notification, error) {
	s.logger.Info(ctx, "Creating notification", "user_id", userID.String(), "type", string(notificationType))
	
	notif := notification.NewNotification(
		userID,
		title,
		message,
		notificationType,
		priority,
		channels,
		data,
	)
	
	err := s.notificationRepo.Create(ctx, notif)
	if err != nil {
		s.logger.Error(ctx, "Failed to create notification", "error", err, "user_id", userID.String())
		return nil, err
	}
	
	s.logger.Info(ctx, "Successfully created notification", "notification_id", notif.ID().String(), "user_id", userID.String())
	return notif, nil
}

// MarkAsRead marks a notification as read
func (s *Service) MarkAsRead(ctx context.Context, notificationID notification.NotificationID) error {
	s.logger.Info(ctx, "Marking notification as read", "notification_id", notificationID.String())
	
	err := s.notificationRepo.MarkAsRead(ctx, notificationID)
	if err != nil {
		s.logger.Error(ctx, "Failed to mark notification as read", "error", err, "notification_id", notificationID.String())
		return err
	}
	
	s.logger.Info(ctx, "Successfully marked notification as read", "notification_id", notificationID.String())
	return nil
}

// MarkAsUnread marks a notification as unread
func (s *Service) MarkAsUnread(ctx context.Context, notificationID notification.NotificationID) error {
	s.logger.Info(ctx, "Marking notification as unread", "notification_id", notificationID.String())
	
	err := s.notificationRepo.MarkAsUnread(ctx, notificationID)
	if err != nil {
		s.logger.Error(ctx, "Failed to mark notification as unread", "error", err, "notification_id", notificationID.String())
		return err
	}
	
	s.logger.Info(ctx, "Successfully marked notification as unread", "notification_id", notificationID.String())
	return nil
}

// MarkAllAsRead marks all notifications for a user as read
func (s *Service) MarkAllAsRead(ctx context.Context, userID shared.UserID) error {
	s.logger.Info(ctx, "Marking all notifications as read", "user_id", userID.String())
	
	err := s.notificationRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, "Failed to mark all notifications as read", "error", err, "user_id", userID.String())
		return err
	}
	
	s.logger.Info(ctx, "Successfully marked all notifications as read", "user_id", userID.String())
	return nil
}

// DeleteNotification deletes a notification
func (s *Service) DeleteNotification(ctx context.Context, notificationID notification.NotificationID) error {
	s.logger.Info(ctx, "Deleting notification", "notification_id", notificationID.String())
	
	err := s.notificationRepo.Delete(ctx, notificationID)
	if err != nil {
		s.logger.Error(ctx, "Failed to delete notification", "error", err, "notification_id", notificationID.String())
		return err
	}
	
	s.logger.Info(ctx, "Successfully deleted notification", "notification_id", notificationID.String())
	return nil
}
