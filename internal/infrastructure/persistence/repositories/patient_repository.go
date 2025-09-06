package repositories

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"medika-backend/internal/domain/patient"
	"medika-backend/internal/infrastructure/persistence/models"
	"medika-backend/pkg/logger"
)

// PatientRepository implements patient.Repository
type PatientRepository struct {
	db     bun.IDB
	logger logger.Logger
}

func NewPatientRepository(db *bun.DB) patient.Repository {
	return &PatientRepository{
		db:     db,
		logger: logger.New(),
	}
}

func (r *PatientRepository) Create(ctx context.Context, p *patient.Patient) error {
	// For now, we'll use the User model since patients are users
	// This is a simplified implementation
	userModel := &models.User{
		ID:             p.ID,
		Email:          p.Email,
		Name:           p.Name,
		Role:           "patient",
		OrganizationID: &p.OrganizationID,
		Phone:          &p.Phone,
		IsActive:       true,
	}

	_, err := r.db.NewInsert().
		Model(userModel).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}

	return nil
}

func (r *PatientRepository) GetByID(ctx context.Context, id string) (*patient.Patient, error) {
	userModel := &models.User{}
	
	err := r.db.NewSelect().
		Model(userModel).
		Where("id = ? AND role = ?", id, "patient").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	return r.toDomain(userModel), nil
}

func (r *PatientRepository) GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*patient.Patient, error) {
	var userModels []models.User
	
	query := r.db.NewSelect().
		Model(&userModels).
		Where("role = ?", "patient")
	
	// Only filter by organization if organizationID is provided
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	}
	
	err := query.
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get patients by organization: %w", err)
	}

	patients := make([]*patient.Patient, len(userModels))
	for i, userModel := range userModels {
		patients[i] = r.toDomain(&userModel)
	}

	return patients, nil
}

func (r *PatientRepository) Update(ctx context.Context, p *patient.Patient) error {
	userModel := &models.User{
		ID:             p.ID,
		Email:          p.Email,
		Name:           p.Name,
		OrganizationID: &p.OrganizationID,
		Phone:          &p.Phone,
	}

	_, err := r.db.NewUpdate().
		Model(userModel).
		Where("id = ? AND role = ?", p.ID, "patient").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	return nil
}

func (r *PatientRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ? AND role = ?", id, "patient").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	return nil
}

func (r *PatientRepository) CountByOrganization(ctx context.Context, organizationID string) (int, error) {
	query := r.db.NewSelect().
		Model((*models.User)(nil)).
		Where("role = ?", "patient")
	
	// Only filter by organization if organizationID is provided
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	}
	
	count, err := query.Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count patients: %w", err)
	}

	return count, nil
}

func (r *PatientRepository) toDomain(userModel *models.User) *patient.Patient {
	// This is a simplified conversion - in a real app you'd have more complex logic
	return &patient.Patient{
		ID:             userModel.ID,
		Name:           userModel.Name,
		Email:          userModel.Email,
		Phone:          *userModel.Phone,
		OrganizationID: *userModel.OrganizationID,
		Status:         "active",
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
	}
}
