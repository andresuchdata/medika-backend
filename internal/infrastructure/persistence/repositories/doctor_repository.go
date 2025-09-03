package repositories

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"medika-backend/internal/domain/doctor"
	"medika-backend/internal/infrastructure/persistence/models"
	"medika-backend/pkg/logger"
)

// DoctorRepository implements doctor.Repository
type DoctorRepository struct {
	db     bun.IDB
	logger logger.Logger
}

func NewDoctorRepository(db *bun.DB) doctor.Repository {
	return &DoctorRepository{
		db:     db,
		logger: logger.New(),
	}
}

func (r *DoctorRepository) Create(ctx context.Context, d *doctor.Doctor) error {
	// For now, we'll use the User model since doctors are users
	// This is a simplified implementation
	userModel := &models.User{
		ID:             d.ID,
		Email:          d.Email,
		Name:           d.Name,
		Role:           "doctor",
		OrganizationID: &d.OrganizationID,
		Phone:          &d.Phone,
		IsActive:       true,
	}

	_, err := r.db.NewInsert().
		Model(userModel).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create doctor: %w", err)
	}

	return nil
}

func (r *DoctorRepository) GetByID(ctx context.Context, id string) (*doctor.Doctor, error) {
	userModel := &models.User{}
	
	err := r.db.NewSelect().
		Model(userModel).
		Where("id = ? AND role = ?", id, "doctor").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}

	return r.toDomain(userModel), nil
}

func (r *DoctorRepository) GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*doctor.Doctor, error) {
	var userModels []models.User
	
	err := r.db.NewSelect().
		Model(&userModels).
		Where("organization_id = ? AND role = ?", organizationID, "doctor").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get doctors by organization: %w", err)
	}

	doctors := make([]*doctor.Doctor, len(userModels))
	for i, userModel := range userModels {
		doctors[i] = r.toDomain(&userModel)
	}

	return doctors, nil
}

func (r *DoctorRepository) Update(ctx context.Context, d *doctor.Doctor) error {
	userModel := &models.User{
		ID:             d.ID,
		Email:          d.Email,
		Name:           d.Name,
		OrganizationID: &d.OrganizationID,
		Phone:          &d.Phone,
	}

	_, err := r.db.NewUpdate().
		Model(userModel).
		Where("id = ? AND role = ?", d.ID, "doctor").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update doctor: %w", err)
	}

	return nil
}

func (r *DoctorRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ? AND role = ?", id, "doctor").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete doctor: %w", err)
	}

	return nil
}

func (r *DoctorRepository) CountByOrganization(ctx context.Context, organizationID string) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.User)(nil)).
		Where("organization_id = ? AND role = ?", organizationID, "doctor").
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count doctors: %w", err)
	}

	return count, nil
}

func (r *DoctorRepository) toDomain(userModel *models.User) *doctor.Doctor {
	// This is a simplified conversion - in a real app you'd have more complex logic
	return &doctor.Doctor{
		ID:             userModel.ID,
		UserID:         userModel.ID,
		Name:           userModel.Name,
		Email:          userModel.Email,
		Phone:          *userModel.Phone,
		OrganizationID: *userModel.OrganizationID,
		Status:         "active",
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
	}
}
