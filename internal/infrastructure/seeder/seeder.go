package seeder

import (
	"context"

	"github.com/uptrace/bun"
)

// Seeder defines the interface for all seeders
type Seeder interface {
	Seed(ctx context.Context, db *bun.DB) error
	Name() string
}

// Manager manages and runs multiple seeders
type Manager struct {
	seeders []Seeder
}

// NewManager creates a new seeder manager
func NewManager() *Manager {
	return &Manager{
		seeders: make([]Seeder, 0),
	}
}

// Register adds a seeder to the manager
func (m *Manager) Register(seeder Seeder) {
	m.seeders = append(m.seeders, seeder)
}

// Run executes all registered seeders
func (m *Manager) Run(ctx context.Context, db *bun.DB) error {
	for _, seeder := range m.seeders {
		if err := seeder.Seed(ctx, db); err != nil {
			return err
		}
	}
	return nil
}

// GetSeeders returns all registered seeders
func (m *Manager) GetSeeders() []Seeder {
	return m.seeders
}
