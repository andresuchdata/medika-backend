package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"

	"medika-backend/internal/domain/queue"
	"medika-backend/internal/infrastructure/persistence/models"
	"medika-backend/pkg/logger"
)

// QueueRepository implements queue.Repository
type QueueRepository struct {
	db     bun.IDB
	logger logger.Logger
}

func NewQueueRepository(db *bun.DB) queue.Repository {
	return &QueueRepository{
		db:     db,
		logger: logger.New(),
	}
}

func (r *QueueRepository) Create(ctx context.Context, q *queue.PatientQueue) error {
	queueModel := &models.PatientQueue{
		ID:                q.ID,
		AppointmentID:     q.AppointmentID,
		OrganizationID:    q.OrganizationID,
		Position:          q.Position,
		EstimatedWaitTime: q.EstimatedWaitTime,
		Status:            q.Status,
	}

	_, err := r.db.NewInsert().
		Model(queueModel).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}

	return nil
}

func (r *QueueRepository) GetByID(ctx context.Context, id string) (*queue.PatientQueue, error) {
	queueModel := &models.PatientQueue{}
	
	err := r.db.NewSelect().
		Model(queueModel).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	return r.toDomain(queueModel), nil
}

func (r *QueueRepository) GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*queue.PatientQueue, error) {
	var queueModels []models.PatientQueue
	
	query := r.db.NewSelect().
		Model(&queueModels)
	
	// Only filter by organization if organizationID is provided
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	}
	
	err := query.
		Order("position ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get queues by organization: %w", err)
	}

	queues := make([]*queue.PatientQueue, len(queueModels))
	for i, queueModel := range queueModels {
		queues[i] = r.toDomain(&queueModel)
	}

	return queues, nil
}

func (r *QueueRepository) GetByAppointment(ctx context.Context, appointmentID string) (*queue.PatientQueue, error) {
	queueModel := &models.PatientQueue{}
	
	err := r.db.NewSelect().
		Model(queueModel).
		Where("appointment_id = ?", appointmentID).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get queue by appointment: %w", err)
	}

	return r.toDomain(queueModel), nil
}

func (r *QueueRepository) GetByPatient(ctx context.Context, patientID string) (*queue.PatientQueue, error) {
	queueModel := &models.PatientQueue{}
	
	// Join with appointments to get the patient's queue entry
	err := r.db.NewSelect().
		Model(queueModel).
		Join("JOIN appointments ON patient_queues.appointment_id = appointments.id").
		Where("appointments.patient_id = ?", patientID).
		Where("appointments.date = CURRENT_DATE"). // Only today's appointments
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get queue by patient: %w", err)
	}

	return r.toDomain(queueModel), nil
}

func (r *QueueRepository) GetPatientQueueWithDetails(ctx context.Context, patientID string) (*queue.PatientQueueWithDetails, error) {
	// Use a raw query to get all the data we need in one go
	var result struct {
		// Queue fields
		QueueID                string    `bun:"queue_id"`
		AppointmentID          string    `bun:"appointment_id"`
		OrganizationID         string    `bun:"organization_id"`
		Position               int       `bun:"position"`
		EstimatedWaitTime      int       `bun:"estimated_wait_time"`
		Status                 string    `bun:"status"`
		QueueCreatedAt         time.Time `bun:"queue_created_at"`
		QueueUpdatedAt         time.Time `bun:"queue_updated_at"`
		
		// Patient fields
		PatientName            string    `bun:"patient_name"`
		PatientID              string    `bun:"patient_id"`
		
		// Doctor fields
		DoctorName             string    `bun:"doctor_name"`
		DoctorID               string    `bun:"doctor_id"`
		
		// Appointment fields
		AppointmentDate        string    `bun:"appointment_date"`
		AppointmentTime        string    `bun:"appointment_time"`
		AppointmentType        string    `bun:"appointment_type"`
		AppointmentStatus      string    `bun:"appointment_status"`
	}
	
	err := r.db.NewRaw(`
		SELECT 
			pq.id as queue_id,
			pq.appointment_id,
			pq.organization_id,
			pq.position,
			pq.estimated_wait_time,
			pq.status,
			pq.created_at as queue_created_at,
			pq.updated_at as queue_updated_at,
			u.name as patient_name,
			a.patient_id,
			d.name as doctor_name,
			a.doctor_id,
			a.date as appointment_date,
			a.start_time as appointment_time,
			a.type as appointment_type,
			a.status as appointment_status
		FROM patient_queues pq
		JOIN appointments a ON pq.appointment_id = a.id
		JOIN users u ON a.patient_id = u.id
		JOIN users d ON a.doctor_id = d.id
		WHERE a.patient_id = ? 
		AND a.date = CURRENT_DATE
		LIMIT 1
	`, patientID).Scan(ctx, &result)

	if err != nil {
		// Check if it's a "no rows" error (which is expected when patient has no appointment today)
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // Return nil, nil to indicate no queue found (not an error)
		}

		return nil, fmt.Errorf("failed to get patient queue with details: %w", err)
	}

	// Create the enriched queue object
	queueWithDetails := &queue.PatientQueueWithDetails{
		PatientQueue: &queue.PatientQueue{
			ID:                result.QueueID,
			AppointmentID:     result.AppointmentID,
			OrganizationID:    result.OrganizationID,
			Position:          result.Position,
			EstimatedWaitTime: result.EstimatedWaitTime,
			Status:            result.Status,
			CreatedAt:         result.QueueCreatedAt,
			UpdatedAt:         result.QueueUpdatedAt,
		},
		PatientName:        result.PatientName,
		PatientID:          result.PatientID,
		DoctorName:         result.DoctorName,
		DoctorID:           result.DoctorID,
		AppointmentDate:    result.AppointmentDate,
		AppointmentTime:    result.AppointmentTime,
		AppointmentType:    result.AppointmentType,
		AppointmentStatus:  result.AppointmentStatus,
	}

	return queueWithDetails, nil
}

func (r *QueueRepository) Update(ctx context.Context, q *queue.PatientQueue) error {
	queueModel := &models.PatientQueue{
		ID:                q.ID,
		AppointmentID:     q.AppointmentID,
		OrganizationID:    q.OrganizationID,
		Position:          q.Position,
		EstimatedWaitTime: q.EstimatedWaitTime,
		Status:            q.Status,
	}

	_, err := r.db.NewUpdate().
		Model(queueModel).
		Where("id = ?", q.ID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update queue: %w", err)
	}

	return nil
}

func (r *QueueRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.PatientQueue)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete queue: %w", err)
	}

	return nil
}

func (r *QueueRepository) CountByOrganization(ctx context.Context, organizationID string) (int, error) {
	query := r.db.NewSelect().
		Model((*models.PatientQueue)(nil))
	
	// Only filter by organization if organizationID is provided
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	}
	
	count, err := query.Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count queues: %w", err)
	}

	return count, nil
}

func (r *QueueRepository) GetNextInQueue(ctx context.Context, organizationID string) (*queue.PatientQueue, error) {
	queueModel := &models.PatientQueue{}
	
	err := r.db.NewSelect().
		Model(queueModel).
		Where("organization_id = ? AND status = ?", organizationID, queue.QueueStatusWaiting).
		Order("position DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get next in queue: %w", err)
	}

	return r.toDomain(queueModel), nil
}

func (r *QueueRepository) UpdatePosition(ctx context.Context, organizationID string) error {
	// This is a simplified implementation
	// In a real application, you might want to use a more sophisticated approach
	// like using database transactions and proper position recalculation
	
	// For now, we'll just update the positions of waiting patients
	_, err := r.db.NewRaw(`
		UPDATE patient_queues 
		SET position = subquery.new_position, updated_at = CURRENT_TIMESTAMP
		FROM (
			SELECT id, ROW_NUMBER() OVER (ORDER BY created_at ASC) as new_position
			FROM patient_queues 
			WHERE organization_id = ? AND status = ?
		) as subquery
		WHERE patient_queues.id = subquery.id
	`, organizationID, queue.QueueStatusWaiting).Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update queue positions: %w", err)
	}

	return nil
}

// GetQueueStats returns aggregated queue statistics for dashboard
func (r *QueueRepository) GetQueueStats(ctx context.Context, organizationID string) (*queue.QueueStats, error) {
	query := r.db.NewSelect().
		Model((*models.PatientQueue)(nil))
	
	// Only filter by organization if organizationID is provided
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	}
	
	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total queues: %w", err)
	}
	
	// Get waiting count
	waitingQuery := r.db.NewSelect().
		Model((*models.PatientQueue)(nil))
	if organizationID != "" {
		waitingQuery = waitingQuery.Where("organization_id = ?", organizationID)
	}
	waiting, err := waitingQuery.Where("status = ?", queue.QueueStatusWaiting).Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count waiting queues: %w", err)
	}
	
	// Get in progress count
	inProgressQuery := r.db.NewSelect().
		Model((*models.PatientQueue)(nil))
	if organizationID != "" {
		inProgressQuery = inProgressQuery.Where("organization_id = ?", organizationID)
	}
	inProgress, err := inProgressQuery.Where("status = ?", queue.QueueStatusInProgress).Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count in progress queues: %w", err)
	}
	
	// Calculate average wait time (simplified - in real app, this would be more sophisticated)
	avgWaitTime := "15 min" // TODO: Calculate from actual data
	
	return &queue.QueueStats{
		Total:           total,
		Waiting:         waiting,
		InProgress:      inProgress,
		AverageWaitTime: avgWaitTime,
	}, nil
}

func (r *QueueRepository) toDomain(queueModel *models.PatientQueue) *queue.PatientQueue {
	return &queue.PatientQueue{
		ID:                queueModel.ID,
		AppointmentID:     queueModel.AppointmentID,
		OrganizationID:    queueModel.OrganizationID,
		Position:          queueModel.Position,
		EstimatedWaitTime: queueModel.EstimatedWaitTime,
		Status:            queueModel.Status,
		CreatedAt:         queueModel.CreatedAt,
		UpdatedAt:         queueModel.UpdatedAt,
	}
}
