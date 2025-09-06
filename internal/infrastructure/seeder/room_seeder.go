package seeder

import (
	"context"
	"fmt"

	"medika-backend/internal/infrastructure/persistence/models"

	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

// RoomSeeder seeds rooms data
type RoomSeeder struct{}

// NewRoomSeeder creates a new room seeder
func NewRoomSeeder() *RoomSeeder {
	return &RoomSeeder{}
}

// Name returns the seeder name
func (s *RoomSeeder) Name() string {
	return "RoomSeeder"
}

// Seed creates rooms
func (s *RoomSeeder) Seed(ctx context.Context, db *bun.DB) error {
	rooms := []*models.Room{
		// Medika General Hospital rooms
		{
			ID:             "11111001-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef",
			Name:           "Consultation Room 1",
			Type:           "consultation",
			Capacity:       4,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Examination table", "Blood pressure monitor", "Stethoscope", "Digital thermometer"},
		},
		{
			ID:             "11111002-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef",
			Name:           "Consultation Room 2",
			Type:           "consultation",
			Capacity:       4,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Examination table", "Blood pressure monitor", "Stethoscope", "Digital thermometer"},
		},
		{
			ID:             "11111003-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef",
			Name:           "Emergency Room",
			Type:           "examination",
			Capacity:       8,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Trauma bed", "Defibrillator", "Oxygen supply", "IV equipment", "Monitors"},
		},
		{
			ID:             "11111004-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef",
			Name:           "Surgery Room A",
			Type:           "procedure",
			Capacity:       6,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Operating table", "Anesthesia machine", "Surgical lights", "Monitors", "Ventilator"},
		},
		{
			ID:             "11111005-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcdef",
			Name:           "Waiting Area Main",
			Type:           "waiting",
			Capacity:       30,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Chairs", "TV", "Magazines", "Water dispenser", "Reception desk"},
		},

		// Downtown Medical Clinic rooms
		{
			ID:             "11111006-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcde0",
			Name:           "Clinic Room 1",
			Type:           "consultation",
			Capacity:       3,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Examination table", "Blood pressure monitor", "Stethoscope"},
		},
		{
			ID:             "11111007-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcde0",
			Name:           "Clinic Room 2",
			Type:           "consultation",
			Capacity:       3,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Examination table", "Blood pressure monitor", "Stethoscope"},
		},
		{
			ID:             "11111008-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcde0",
			Name:           "Clinic Waiting Area",
			Type:           "waiting",
			Capacity:       15,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Chairs", "TV", "Magazines"},
		},

		// Private Practice Center rooms
		{
			ID:             "11111009-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcde1",
			Name:           "Private Office 1",
			Type:           "office",
			Capacity:       3,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Desk", "Examination table", "Medical equipment"},
		},
		{
			ID:             "11111010-1111-1111-1111-111111111111",
			OrganizationID: "01234567-89ab-cdef-0123-456789abcde1",
			Name:           "Private Waiting Room",
			Type:           "waiting",
			Capacity:       8,
			IsAvailable:    true,
			Equipment:      pq.StringArray{"Comfortable chairs", "Coffee table", "Magazines"},
		},
	}

	for _, room := range rooms {
		_, err := db.NewInsert().
			Model(room).
			On("CONFLICT (id) DO UPDATE").
			Set("organization_id = EXCLUDED.organization_id").
			Set("name = EXCLUDED.name").
			Set("type = EXCLUDED.type").
			Set("capacity = EXCLUDED.capacity").
			Set("is_available = EXCLUDED.is_available").
			Set("equipment = EXCLUDED.equipment").
			Set("updated_at = CURRENT_TIMESTAMP").
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to seed room %s: %w", room.Name, err)
		}
	}

	fmt.Printf("✅ Seeded %d rooms\n", len(rooms))
	return nil
}
