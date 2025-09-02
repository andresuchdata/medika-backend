package seeder

import (
	"context"
	"fmt"

	"medika-backend/internal/infrastructure/persistence/models"

	"github.com/uptrace/bun"
)

// OrganizationSeeder seeds organizations data
type OrganizationSeeder struct{}

// NewOrganizationSeeder creates a new organization seeder
func NewOrganizationSeeder() *OrganizationSeeder {
	return &OrganizationSeeder{}
}

// Name returns the seeder name
func (s *OrganizationSeeder) Name() string {
	return "OrganizationSeeder"
}

// Seed creates organizations
func (s *OrganizationSeeder) Seed(ctx context.Context, db *bun.DB) error {
	organizations := []*models.Organization{
		{
			ID:          "01234567-89ab-cdef-0123-456789abcdef",
			Name:        "Medika General Hospital",
			Type:        "hospital",
			Address:     "123 Medical Center Dr, Healthcare City, HC 12345",
			Phone:       "+1-555-0123",
			Email:       "info@medikahospital.com",
			Website:     stringPtr("https://medikahospital.com"),
			IsActive:    true,
		},
		{
			ID:          "01234567-89ab-cdef-0123-456789abcde0",
			Name:        "Downtown Medical Clinic",
			Type:        "clinic",
			Address:     "456 Downtown Ave, Medical District, MD 67890",
			Phone:       "+1-555-0456",
			Email:       "contact@downtownmedical.com",
			Website:     stringPtr("https://downtownmedical.com"),
			IsActive:    true,
		},
		{
			ID:          "01234567-89ab-cdef-0123-456789abcde1",
			Name:        "Private Practice Center",
			Type:        "private_practice",
			Address:     "789 Health St, Wellness Plaza, WP 54321",
			Phone:       "+1-555-0789",
			Email:       "info@privatepractice.com",
			Website:     stringPtr("https://privatepractice.com"),
			IsActive:    true,
		},
	}

	for _, org := range organizations {
		_, err := db.NewInsert().
			Model(org).
			On("CONFLICT (id) DO UPDATE").
			Set("name = EXCLUDED.name").
			Set("type = EXCLUDED.type").
			Set("address = EXCLUDED.address").
			Set("phone = EXCLUDED.phone").
			Set("email = EXCLUDED.email").
			Set("website = EXCLUDED.website").
			Set("is_active = EXCLUDED.is_active").
			Set("updated_at = CURRENT_TIMESTAMP").
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to seed organization %s: %w", org.Name, err)
		}
	}

	fmt.Printf("✅ Seeded %d organizations\n", len(organizations))
	return nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
