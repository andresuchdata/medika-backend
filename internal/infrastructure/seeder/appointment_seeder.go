package seeder

import (
	"context"
	"fmt"
	"time"

	"medika-backend/internal/infrastructure/persistence/models"

	"github.com/uptrace/bun"
)

// AppointmentSeeder seeds appointments data
type AppointmentSeeder struct{}

// NewAppointmentSeeder creates a new appointment seeder
func NewAppointmentSeeder() *AppointmentSeeder {
	return &AppointmentSeeder{}
}

// Name returns the seeder name
func (s *AppointmentSeeder) Name() string {
	return "AppointmentSeeder"
}

// Seed creates appointments
func (s *AppointmentSeeder) Seed(ctx context.Context, db *bun.DB) error {
	// Get current date
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	
	appointments := []*models.Appointment{
		{
			ID:             "11170001-1111-1111-1111-111111111111",
			PatientID:      "11150001-1111-1111-1111-111111111111", // John Doe
			DoctorID:       "11130001-1111-1111-1111-111111111111", // Dr. Michael Smith
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			RoomID:         stringPtr("11111001-1111-1111-1111-111111111111"),
			Date:           today,
			StartTime:      "09:00:00",
			EndTime:        "09:30:00",
			Duration:       30,
			Status:         "confirmed",
			Type:           "consultation",
			Notes:          stringPtr("Regular checkup"),
		},
		{
			ID:             "11170002-1111-1111-1111-111111111111",
			PatientID:      "11150002-1111-1111-1111-111111111111", // Jane Smith
			DoctorID:       "11130002-1111-1111-1111-111111111111", // Dr. Jennifer Jones
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			RoomID:         stringPtr("11111002-1111-1111-1111-111111111111"),
			Date:           today,
			StartTime:      "10:00:00",
			EndTime:        "10:30:00",
			Duration:       30,
			Status:         "confirmed",
			Type:           "follow_up",
			Notes:          stringPtr("Follow-up appointment"),
		},
		{
			ID:             "11170003-1111-1111-1111-111111111111",
			PatientID:      "11150003-1111-1111-1111-111111111111", // Bob Johnson
			DoctorID:       "11130003-1111-1111-1111-111111111111", // Dr. Robert Brown
			OrganizationID: "01234567-89ab-cdef-0123-456789abcde0", // Downtown Medical Clinic
			RoomID:         stringPtr("11111006-1111-1111-1111-111111111111"),
			Date:           today,
			StartTime:      "14:00:00",
			EndTime:        "14:45:00",
			Duration:       45,
			Status:         "confirmed",
			Type:           "consultation",
			Notes:          stringPtr("New patient consultation"),
		},
		{
			ID:             "11170004-1111-1111-1111-111111111111",
			PatientID:      "11150001-1111-1111-1111-111111111111", // John Doe
			DoctorID:       "11130001-1111-1111-1111-111111111111", // Dr. Michael Smith
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			RoomID:         stringPtr("11111001-1111-1111-1111-111111111111"),
			Date:           today.AddDate(0, 0, 1), // Tomorrow
			StartTime:      "11:00:00",
			EndTime:        "11:30:00",
			Duration:       30,
			Status:         "pending",
			Type:           "consultation",
			Notes:          stringPtr("Follow-up appointment"),
		},
		{
			ID:             "11170005-1111-1111-1111-111111111111",
			PatientID:      "11150002-1111-1111-1111-111111111111", // Jane Smith
			DoctorID:       "11130002-1111-1111-1111-111111111111", // Dr. Jennifer Jones
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef", // Medika General Hospital
			RoomID:         stringPtr("11111002-1111-1111-1111-111111111111"),
			Date:           today.AddDate(0, 0, 2), // Day after tomorrow
			StartTime:      "15:00:00",
			EndTime:        "15:30:00",
			Duration:       30,
			Status:         "pending",
			Type:           "routine_checkup",
			Notes:          stringPtr("Annual physical examination"),
		},
	}

	for _, appointment := range appointments {
		_, err := db.NewInsert().
			Model(appointment).
			On("CONFLICT (id) DO UPDATE").
			Set("patient_id = EXCLUDED.patient_id").
			Set("doctor_id = EXCLUDED.doctor_id").
			Set("organization_id = EXCLUDED.organization_id").
			Set("room_id = EXCLUDED.room_id").
			Set("date = EXCLUDED.date").
			Set("start_time = EXCLUDED.start_time").
			Set("end_time = EXCLUDED.end_time").
			Set("duration = EXCLUDED.duration").
			Set("status = EXCLUDED.status").
			Set("type = EXCLUDED.type").
			Set("notes = EXCLUDED.notes").
			Set("updated_at = CURRENT_TIMESTAMP").
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to seed appointment %s: %w", appointment.ID, err)
		}
	}

	fmt.Printf("✅ Seeded %d appointments\n", len(appointments))
	return nil
}
