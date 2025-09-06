package seeder

import (
	"context"
	"fmt"

	"medika-backend/internal/infrastructure/persistence/models"
	"medika-backend/pkg/crypto"

	"github.com/uptrace/bun"
)

// UserSeeder seeds users data
type UserSeeder struct{}

// NewUserSeeder creates a new user seeder
func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

// Name returns the seeder name
func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

// Seed creates users
func (s *UserSeeder) Seed(ctx context.Context, db *bun.DB) error {
	// Default password for all test users
	defaultPassword := crypto.MustHashPassword("password123")
	
	users := []*models.User{
		// Admin users
		{
			ID:             "11120001-1111-1111-1111-111111111111",
			Email:          "admin@medika.com",
			Name:           "Dr. Sarah Johnson",
			PasswordHash:   defaultPassword,
			Role:           "admin",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550001"),
			IsActive:       true,
		},
		{
			ID:             "11120002-1111-1111-1111-111111111111",
			Name:           "System Administrator",
			Email:          "sysadmin@medika.com",
			PasswordHash:   defaultPassword,
			Role:           "admin",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550011"),
			IsActive:       true,
		},

		// Doctor users
		{
			ID:             "11130001-1111-1111-1111-111111111111",
			Email:          "doctor.smith@medika.com",
			Name:           "Dr. Michael Smith",
			PasswordHash:   defaultPassword,
			Role:           "doctor",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550002"),
			IsActive:       true,
		},
		{
			ID:             "11130002-1111-1111-1111-111111111111",
			Email:          "doctor.jones@medika.com",
			Name:           "Dr. Jennifer Jones",
			PasswordHash:   defaultPassword,
			Role:           "doctor",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550022"),
			IsActive:       true,
		},
		{
			ID:             "11130003-1111-1111-1111-111111111111",
			Email:          "doctor.brown@downtownmedical.com",
			Name:           "Dr. Robert Brown",
			PasswordHash:   defaultPassword,
			Role:           "doctor",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcde0"),
			Phone:          stringPtr("+15555550023"),
			IsActive:       true,
		},

		// Nurse users
		{
			ID:             "11140001-1111-1111-1111-111111111111",
			Email:          "nurse.wilson@medika.com",
			Name:           "Nurse Emily Wilson",
			PasswordHash:   defaultPassword,
			Role:           "nurse",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550003"),
			IsActive:       true,
		},
		{
			ID:             "11140002-1111-1111-1111-111111111111",
			Email:          "nurse.davis@medika.com",
			Name:           "Nurse Lisa Davis",
			PasswordHash:   defaultPassword,
			Role:           "nurse",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550033"),
			IsActive:       true,
		},

		// Patient users
		{
			ID:             "11150001-1111-1111-1111-111111111111",
			Email:          "patient.john@email.com",
			Name:           "John Doe",
			PasswordHash:   defaultPassword,
			Role:           "patient",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550004"),
			IsActive:       true,
		},
		{
			ID:             "11150002-1111-1111-1111-111111111111",
			Email:          "patient.jane@email.com",
			Name:           "Jane Smith",
			PasswordHash:   defaultPassword,
			Role:           "patient",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550044"),
			IsActive:       true,
		},
		{
			ID:             "11150003-1111-1111-1111-111111111111",
			Email:          "patient.bob@email.com",
			Name:           "Bob Johnson",
			PasswordHash:   defaultPassword,
			Role:           "patient",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcde0"),
			Phone:          stringPtr("+15555550045"),
			IsActive:       true,
		},

		// Cashier users
		{
			ID:             "11160001-1111-1111-1111-111111111111",
			Email:          "cashier@medika.com",
			Name:           "Maria Rodriguez",
			PasswordHash:   defaultPassword,
			Role:           "cashier",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550005"),
			IsActive:       true,
		},
		{
			ID:             "11160002-1111-1111-1111-111111111111",
			Email:          "cashier2@medika.com",
			Name:           "Carlos Martinez",
			PasswordHash:   defaultPassword,
			Role:           "cashier",
			OrganizationID: stringPtr("01234567-89ab-cdef-0123-456789abcdef"),
			Phone:          stringPtr("+15555550055"),
			IsActive:       true,
		},
	}

	for _, user := range users {
		_, err := db.NewInsert().
			Model(user).
			On("CONFLICT (email) DO UPDATE").
			Set("name = EXCLUDED.name").
			Set("password_hash = EXCLUDED.password_hash").
			Set("role = EXCLUDED.role").
			Set("organization_id = EXCLUDED.organization_id").
			Set("phone = EXCLUDED.phone").
			Set("is_active = EXCLUDED.is_active").
			Set("updated_at = CURRENT_TIMESTAMP").
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Email, err)
		}
	}

	fmt.Printf("✅ Seeded %d users\n", len(users))
	return nil
}
