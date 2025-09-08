package queue

import (
	"context"
	"fmt"

	"medika-backend/internal/domain/queue"
	"medika-backend/pkg/logger"
)

type Service struct {
	queueRepo queue.Repository
	logger    logger.Logger
}

func NewService(queueRepo queue.Repository, logger logger.Logger) *Service {
	return &Service{
		queueRepo: queueRepo,
		logger:    logger,
	}
}

// GetQueuesByOrganization retrieves queues for a specific organization
func (s *Service) GetQueuesByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*queue.PatientQueue, error) {
	return s.queueRepo.GetByOrganization(ctx, organizationID, limit, offset)
}

// CountQueuesByOrganization counts total queues for a specific organization
func (s *Service) CountQueuesByOrganization(ctx context.Context, organizationID string) (int, error) {
	return s.queueRepo.CountByOrganization(ctx, organizationID)
}

// CreateQueue creates a new patient queue entry
func (s *Service) CreateQueue(ctx context.Context, q *queue.PatientQueue) error {
	// Set initial status
	q.Status = queue.QueueStatusWaiting
	
	// Get next position for the organization
	nextQueue, err := s.queueRepo.GetNextInQueue(ctx, q.OrganizationID)
	if err != nil {
		// If no existing queues, start with position 1
		q.Position = 1
	} else {
		q.Position = nextQueue.Position + 1
	}
	
	// Estimate wait time (15 minutes per position)
	q.EstimatedWaitTime = q.Position * 15
	
	return s.queueRepo.Create(ctx, q)
}

// GetQueue retrieves a specific queue by ID
func (s *Service) GetQueue(ctx context.Context, id string) (*queue.PatientQueue, error) {
	return s.queueRepo.GetByID(ctx, id)
}

// GetPatientQueue retrieves a patient's queue with enriched details
func (s *Service) GetPatientQueue(ctx context.Context, patientID string) (*queue.PatientQueueWithDetails, error) {
	patientQueue, err := s.queueRepo.GetPatientQueueWithDetails(ctx, patientID)
	if err != nil {
		return nil, err
	}
	
	// If no queue found (patient has no appointment today), return nil without error
	if patientQueue == nil {
		return nil, nil
	}
	
	return patientQueue, nil
}

// UpdateQueue updates an existing queue
func (s *Service) UpdateQueue(ctx context.Context, q *queue.PatientQueue) error {
	// If status changed to completed or cancelled, update positions
	if q.Status == queue.QueueStatusCompleted || q.Status == queue.QueueStatusCancelled {
		if err := s.queueRepo.UpdatePosition(ctx, q.OrganizationID); err != nil {
			s.logger.Error(ctx, "Failed to update queue positions", "error", err)
		}
	}
	
	return s.queueRepo.Update(ctx, q)
}

// DeleteQueue removes a queue entry
func (s *Service) DeleteQueue(ctx context.Context, id string) error {
	// Get queue to update positions after deletion
	q, err := s.queueRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get queue for deletion: %w", err)
	}
	
	// Delete the queue
	if err := s.queueRepo.Delete(ctx, id); err != nil {
		return err
	}
	
	// Update positions for the organization
	if err := s.queueRepo.UpdatePosition(ctx, q.OrganizationID); err != nil {
		s.logger.Error(ctx, "Failed to update queue positions after deletion", "error", err)
	}
	
	return nil
}

// CallNextPatient calls the next patient in the queue
func (s *Service) CallNextPatient(ctx context.Context, organizationID string) (*queue.PatientQueue, error) {
	nextQueue, err := s.queueRepo.GetNextInQueue(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("no patients in queue: %w", err)
	}
	
	// Update status to called
	nextQueue.Status = queue.QueueStatusCalled
	if err := s.queueRepo.Update(ctx, nextQueue); err != nil {
		return nil, fmt.Errorf("failed to update queue status: %w", err)
	}
	
	return nextQueue, nil
}

// StartPatientConsultation starts consultation for a patient
func (s *Service) StartPatientConsultation(ctx context.Context, queueID string) (*queue.PatientQueue, error) {
	q, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}
	
	q.Status = queue.QueueStatusInProgress
	if err := s.queueRepo.Update(ctx, q); err != nil {
		return nil, fmt.Errorf("failed to update queue status: %w", err)
	}
	
	return q, nil
}

// CompletePatientConsultation completes consultation for a patient
func (s *Service) CompletePatientConsultation(ctx context.Context, queueID string) (*queue.PatientQueue, error) {
	q, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}
	
	q.Status = queue.QueueStatusCompleted
	if err := s.queueRepo.Update(ctx, q); err != nil {
		return nil, fmt.Errorf("failed to update queue status: %w", err)
	}
	
	// Update positions for the organization
	if err := s.queueRepo.UpdatePosition(ctx, q.OrganizationID); err != nil {
		s.logger.Error(ctx, "Failed to update queue positions", "error", err)
	}
	
	return q, nil
}
