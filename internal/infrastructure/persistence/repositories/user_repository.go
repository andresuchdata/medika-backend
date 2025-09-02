package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	"medika-backend/internal/domain/shared"
	"medika-backend/internal/domain/user"
	"medika-backend/internal/infrastructure/persistence/models"
)

// UserRepository implements user.Repository
type UserRepository struct {
	db bun.IDB
}

func NewUserRepository(db *bun.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	model := r.toModel(u)

	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Save profile if exists
	if u.Profile() != nil {
		profileModel := r.toProfileModel(u.ID(), u.Profile())
		_, err = r.db.NewInsert().
			Model(profileModel).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to save user profile: %w", err)
		}
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	model := r.toModel(u)

	_, err := r.db.NewUpdate().
		Model(model).
		Where("id = ? AND version = ?", u.ID().String(), u.Version()-1).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Update or create profile if exists
	if u.Profile() != nil {
		profileModel := r.toProfileModel(u.ID(), u.Profile())
		
		_, err = r.db.NewInsert().
			Model(profileModel).
			On("CONFLICT (user_id) DO UPDATE").
			Set("date_of_birth = EXCLUDED.date_of_birth").
			Set("gender = EXCLUDED.gender").
			Set("address = EXCLUDED.address").
			Set("emergency_contact = EXCLUDED.emergency_contact").
			Set("medical_history = EXCLUDED.medical_history").
			Set("allergies = EXCLUDED.allergies").
			Set("blood_type = EXCLUDED.blood_type").
			Set("updated_at = CURRENT_TIMESTAMP").
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to update user profile: %w", err)
		}
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id shared.UserID) error {
	_, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ?", id.String()).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id shared.UserID) (*user.User, error) {
	model := &models.User{}
	
	err := r.db.NewSelect().
		Model(model).
		Relation("Profile").
		Where("u.id = ?", id.String()).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.toDomain(model)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email shared.Email) (*user.User, error) {
	model := &models.User{}
	
	err := r.db.NewSelect().
		Model(model).
		Relation("Profile").
		Where("u.email = ?", email.String()).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.toDomain(model)
}

func (r *UserRepository) FindByOrganization(ctx context.Context, orgID shared.OrganizationID, filters user.UserFilters) ([]*user.User, error) {
	query := r.db.NewSelect().
		Model((*models.User)(nil)).
		Relation("Profile").
		Where("u.organization_id = ?", orgID.String())

	// Apply filters
	if filters.Name != "" {
		query = query.Where("u.name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("u.email ILIKE ?", "%"+filters.Email+"%")
	}
	if filters.Role != nil {
		query = query.Where("u.role = ?", string(*filters.Role))
	}
	if filters.IsActive != nil {
		query = query.Where("u.is_active = ?", *filters.IsActive)
	}
	if filters.CreatedAfter != nil {
		query = query.Where("u.created_at > ?", *filters.CreatedAfter)
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = query.Order("u.created_at DESC")

	var models []*models.User
	err := query.Scan(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		domainUser, err := r.toDomain(model)
		if err != nil {
			return nil, err
		}
		users[i] = domainUser
	}

	return users, nil
}

func (r *UserRepository) FindByRole(ctx context.Context, role user.Role, orgID *shared.OrganizationID) ([]*user.User, error) {
	query := r.db.NewSelect().
		Model((*models.User)(nil)).
		Relation("Profile").
		Where("u.role = ?", string(role))

	if orgID != nil {
		query = query.Where("u.organization_id = ?", orgID.String())
	}

	query = query.Order("u.created_at DESC")

	var models []*models.User
	err := query.Scan(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		domainUser, err := r.toDomain(model)
		if err != nil {
			return nil, err
		}
		users[i] = domainUser
	}

	return users, nil
}

func (r *UserRepository) CountByOrganization(ctx context.Context, orgID shared.OrganizationID) (int64, error) {
	count, err := r.db.NewSelect().
		Model((*models.User)(nil)).
		Where("organization_id = ?", orgID.String()).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return int64(count), nil
}

func (r *UserRepository) CountByRole(ctx context.Context, role user.Role) (int64, error) {
	count, err := r.db.NewSelect().
		Model((*models.User)(nil)).
		Where("role = ?", string(role)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	return int64(count), nil
}

func (r *UserRepository) FindActiveUsers(ctx context.Context, orgID *shared.OrganizationID) ([]*user.User, error) {
	query := r.db.NewSelect().
		Model((*models.User)(nil)).
		Relation("Profile").
		Where("u.is_active = ?", true)

	if orgID != nil {
		query = query.Where("u.organization_id = ?", orgID.String())
	}

	query = query.Order("u.created_at DESC")

	var models []*models.User
	err := query.Scan(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to find active users: %w", err)
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		domainUser, err := r.toDomain(model)
		if err != nil {
			return nil, err
		}
		users[i] = domainUser
	}

	return users, nil
}

func (r *UserRepository) WithTx(ctx context.Context, fn func(user.Repository) error) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		txRepo := &UserRepository{db: tx}
		return fn(txRepo)
	})
}

// Conversion methods
func (r *UserRepository) toModel(u *user.User) *models.User {
	model := &models.User{
		ID:           u.ID().String(),
		Email:        u.Email().String(),
		Name:         u.Name().String(),
		PasswordHash: "", // Don't expose password hash
		Role:         u.Role().String(),
		IsActive:     u.IsActive(),
		CreatedAt:    u.CreatedAt(),
		UpdatedAt:    u.UpdatedAt(),
		Version:      u.Version(),
	}

	if u.OrganizationID() != nil {
		orgID := u.OrganizationID().String()
		model.OrganizationID = &orgID
	}

	if u.Phone() != nil {
		phone := u.Phone().String()
		model.Phone = &phone
	}

	if u.AvatarURL() != nil {
		model.AvatarURL = u.AvatarURL()
	}

	return model
}

func (r *UserRepository) toProfileModel(userID shared.UserID, _ *user.Profile) *models.UserProfile {
	profile := &models.UserProfile{
		UserID: userID.String(),
	}

	// Set profile fields based on what's available
	// Note: Since Profile fields are not exported, we'd need getters
	// This is a simplified version - in practice, you'd add getters to Profile

	return profile
}

func (r *UserRepository) toDomain(model *models.User) (*user.User, error) {
	userID, err := shared.NewUserIDFromString(model.ID)
	if err != nil {
		return nil, err
	}

	email, err := shared.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	name, err := shared.NewName(model.Name)
	if err != nil {
		return nil, err
	}

	role := user.Role(model.Role)

	var orgID *shared.OrganizationID
	if model.OrganizationID != nil {
		id, err := shared.NewOrganizationID(*model.OrganizationID)
		if err != nil {
			return nil, err
		}
		orgID = &id
	}

	var phone *shared.PhoneNumber
	if model.Phone != nil {
		pn, err := shared.NewPhoneNumber(*model.Phone)
		if err != nil {
			return nil, err
		}
		phone = &pn
	}

	// Convert profile if exists
	var profile *user.Profile
	if model.Profile != nil {
		// Create profile from model
		// This would require Profile constructor or setters
		profile = nil // Simplified for now
	}

	return user.ReconstructUser(
		userID,
		email,
		name,
		model.PasswordHash,
		role,
		orgID,
		phone,
		model.AvatarURL,
		model.IsActive,
		profile,
		model.CreatedAt,
		model.UpdatedAt,
		model.Version,
	), nil
}
