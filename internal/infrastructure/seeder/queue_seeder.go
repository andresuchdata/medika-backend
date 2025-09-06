package seeder

import (
	"context"
	"fmt"

	"medika-backend/internal/domain/queue"
	"medika-backend/internal/infrastructure/persistence/models"

	"github.com/uptrace/bun"
)

// QueueSeeder seeds patient queues data
type QueueSeeder struct{}

// NewQueueSeeder creates a new queue seeder
func NewQueueSeeder() *QueueSeeder {
	return &QueueSeeder{}
}

// Name returns the seeder name
func (s *QueueSeeder) Name() string {
	return "QueueSeeder"
}

// Seed creates patient queues
func (s *QueueSeeder) Seed(ctx context.Context, db *bun.DB) error {
	queues := []*models.PatientQueue{
		{
			ID:                "11180001-1111-1111-1111-111111111111",
			AppointmentID:     "11170001-1111-1111-1111-111111111111", // John Doe's appointment
			OrganizationID:    "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			Position:          1,
			EstimatedWaitTime: 15,
			Status:            queue.QueueStatusWaiting,
		},
		{
			ID:                "11180002-1111-1111-1111-111111111111",
			AppointmentID:     "11170002-1111-1111-1111-111111111111", // Jane Smith's appointment
			OrganizationID:    "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			Position:          2,
			EstimatedWaitTime: 30,
			Status:            queue.QueueStatusWaiting,
		},
		{
			ID:                "11180003-1111-1111-1111-111111111111",
			AppointmentID:     "11170003-1111-1111-1111-111111111111", // Bob Johnson's appointment
			OrganizationID:    "01234567-89ab-cdef-0123-456789abcde0", // Downtown Medical Clinic
			Position:          1,
			EstimatedWaitTime: 15,
			Status:            queue.QueueStatusCalled,
		},
		{
			ID:                "11180004-1111-1111-1111-111111111111",
			AppointmentID:     "11170004-1111-1111-1111-111111111111", // John Doe's tomorrow appointment
			OrganizationID:    "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			Position:          3,
			EstimatedWaitTime: 45,
			Status:            queue.QueueStatusWaiting,
		},
		{
			ID:                "11180005-1111-1111-1111-111111111111",
			AppointmentID:     "11170005-1111-1111-1111-111111111111", // Jane Smith's future appointment
			OrganizationID:    "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			Position:          4,
			EstimatedWaitTime: 60,
			Status:            queue.QueueStatusWaiting,
		},
	}

	for _, queue := range queues {
		_, err := db.NewInsert().
			Model(queue).
			On("CONFLICT (id) DO UPDATE").
			Set("appointment_id = EXCLUDED.appointment_id").
			Set("organization_id = EXCLUDED.organization_id").
			Set("position = EXCLUDED.position").
			Set("estimated_wait_time = EXCLUDED.estimated_wait_time").
			Set("status = EXCLUDED.status").
			Set("updated_at = CURRENT_TIMESTAMP").
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to seed queue %s: %w", queue.ID, err)
		}
	}

	fmt.Printf("✅ Seeded %d patient queues\n", len(queues))
	return nil
}
