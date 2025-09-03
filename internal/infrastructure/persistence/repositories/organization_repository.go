package repositories

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"medika-backend/internal/domain/organization"
	"medika-backend/internal/infrastructure/persistence/models"
	"medika-backend/pkg/logger"
)

// OrganizationRepository implements organization.Repository
type OrganizationRepository struct {
	db     bun.IDB
	logger logger.Logger
}

func NewOrganizationRepository(db *bun.DB) organization.Repository {
	return &OrganizationRepository{
		db:     db,
		logger: logger.New(),
	}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *organization.Organization) error {
	model := r.toModel(org)

	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id string) (*organization.Organization, error) {
	model := &models.Organization{}
	
	err := r.db.NewSelect().
		Model(model).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return r.toDomain(model), nil
}

func (r *OrganizationRepository) GetAll(ctx context.Context, limit, offset int) ([]*organization.Organization, error) {
	var models []models.Organization
	
	err := r.db.NewSelect().
		Model(&models).
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	organizations := make([]*organization.Organization, len(models))
	for i, model := range models {
		organizations[i] = r.toDomain(&model)
	}

	return organizations, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, org *organization.Organization) error {
	model := r.toModel(org)

	_, err := r.db.NewUpdate().
		Model(model).
		Where("id = ?", org.ID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.Organization)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Organization)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count organizations: %w", err)
	}

	return count, nil
}

func (r *OrganizationRepository) toModel(org *organization.Organization) *models.Organization {
	return &models.Organization{
		ID:          org.ID,
		Name:        org.Name,
		Type:        org.Type,
		Address:     org.Address,
		Phone:       org.Phone,
		Email:       org.Email,
		Website:     org.Website,
		IsActive:    org.IsActive,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}

func (r *OrganizationRepository) toDomain(model *models.Organization) *organization.Organization {
	return &organization.Organization{
		ID:          model.ID,
		Name:        model.Name,
		Type:        model.Type,
		Address:     model.Address,
		Phone:       model.Phone,
		Email:       model.Email,
		Website:     model.Website,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
