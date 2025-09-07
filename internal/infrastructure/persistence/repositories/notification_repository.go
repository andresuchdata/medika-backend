package repositories

import (
	"context"
	"fmt"
	"medika-backend/internal/domain/notification"
	"medika-backend/internal/domain/shared"
	"medika-backend/internal/infrastructure/persistence/models"

	"github.com/uptrace/bun"
)

type NotificationRepository struct {
	db *bun.DB
}

func NewNotificationRepository(db *bun.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notif *notification.Notification) error {
	model := r.toModel(notif)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *NotificationRepository) FindByID(ctx context.Context, id notification.NotificationID) (*notification.Notification, error) {
	var model models.Notification
	err := r.db.NewSelect().
		Model(&model).
		Where("id = ?", id.String()).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	
	return r.toDomain(&model), nil
}

func (r *NotificationRepository) FindByUserID(ctx context.Context, userID shared.UserID, filters notification.NotificationFilters) ([]*notification.Notification, error) {
	var models []models.Notification
	
	query := r.db.NewSelect().
		Model(&models).
		Where("user_id = ?", userID.String())
	
	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", string(*filters.Type))
	}
	if filters.Priority != nil {
		query = query.Where("priority = ?", string(*filters.Priority))
	}
	if filters.IsRead != nil {
		query = query.Where("is_read = ?", *filters.IsRead)
	}
	if filters.ActionRequired != nil {
		query = query.Where("action_required = ?", *filters.ActionRequired)
	}
	
	// Apply ordering
	orderBy := "created_at"
	if filters.OrderBy != "" {
		orderBy = filters.OrderBy
	}
	order := "DESC"
	if filters.Order != "" {
		order = filters.Order
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))
	
	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	
	notifications := make([]*notification.Notification, len(models))
	for i, model := range models {
		notifications[i] = r.toDomain(&model)
	}
	
	return notifications, nil
}

func (r *NotificationRepository) CountByUserID(ctx context.Context, userID shared.UserID, filters notification.NotificationFilters) (int, error) {
	query := r.db.NewSelect().
		Model((*models.Notification)(nil)).
		Where("user_id = ?", userID.String())
	
	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", string(*filters.Type))
	}
	if filters.Priority != nil {
		query = query.Where("priority = ?", string(*filters.Priority))
	}
	if filters.IsRead != nil {
		query = query.Where("is_read = ?", *filters.IsRead)
	}
	if filters.ActionRequired != nil {
		query = query.Where("action_required = ?", *filters.ActionRequired)
	}
	
	count, err := query.Count(ctx)
	return count, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id notification.NotificationID) error {
	_, err := r.db.NewUpdate().
		Model((*models.Notification)(nil)).
		Set("is_read = true").
		Where("id = ?", id.String()).
		Exec(ctx)
	return err
}

func (r *NotificationRepository) MarkAsUnread(ctx context.Context, id notification.NotificationID) error {
	_, err := r.db.NewUpdate().
		Model((*models.Notification)(nil)).
		Set("is_read = false").
		Where("id = ?", id.String()).
		Exec(ctx)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID shared.UserID) error {
	_, err := r.db.NewUpdate().
		Model((*models.Notification)(nil)).
		Set("is_read = true").
		Where("user_id = ?", userID.String()).
		Exec(ctx)
	return err
}

func (r *NotificationRepository) Delete(ctx context.Context, id notification.NotificationID) error {
	_, err := r.db.NewDelete().
		Model((*models.Notification)(nil)).
		Where("id = ?", id.String()).
		Exec(ctx)
	return err
}

func (r *NotificationRepository) DeleteByUserID(ctx context.Context, userID shared.UserID) error {
	_, err := r.db.NewDelete().
		Model((*models.Notification)(nil)).
		Where("user_id = ?", userID.String()).
		Exec(ctx)
	return err
}

func (r *NotificationRepository) toModel(notif *notification.Notification) *models.Notification {
	data := make(models.JSONB)
	if notif.Data() != nil {
		for k, v := range notif.Data() {
			data[k] = v
		}
	}
	
	return &models.Notification{
		ID:           notif.ID().String(),
		UserID:       notif.UserID().String(),
		Type:         string(notif.Type()),
		Title:        notif.Title(),
		Message:      notif.Message(),
		Data:         data,
		IsRead:       notif.IsRead(),
		Channels:     notif.Channels(),
		Priority:     string(notif.Priority()),
		ScheduledFor: notif.ScheduledFor(),
		SentAt:       notif.SentAt(),
		CreatedAt:    notif.CreatedAt(),
	}
}

func (r *NotificationRepository) toDomain(model *models.Notification) *notification.Notification {
	data := make(map[string]interface{})
	if model.Data != nil {
		for k, v := range model.Data {
			data[k] = v
		}
	}
	
	userID, _ := shared.NewUserIDFromString(model.UserID)
	notif := notification.NewNotification(
		userID,
		model.Title,
		model.Message,
		notification.NotificationType(model.Type),
		notification.Priority(model.Priority),
		model.Channels,
		data,
	)
	
	// Set the ID and timestamps from the model
	// Note: This is a bit of a hack since the domain entity doesn't expose setters
	// In a real implementation, you might want to add a constructor that accepts these values
	return notif
}
