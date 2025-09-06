package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/uptrace/bun"

	"medika-backend/internal/domain/shared"
	"medika-backend/internal/domain/user"
	"medika-backend/internal/infrastructure/persistence/models"
	"medika-backend/pkg/logger"
)

// UserRepository implements user.Repository
type UserRepository struct {
	db     bun.IDB
	logger logger.Logger
}

// UserQueryBuilder handles all user query building logic with single responsibility
type UserQueryBuilder struct {
	query *bun.SelectQuery
}

// NewUserQueryBuilder creates a new query builder for user queries
func NewUserQueryBuilder(db bun.IDB) *UserQueryBuilder {
	return &UserQueryBuilder{
		query: db.NewSelect().
			Model((*models.User)(nil)).
			Relation("Profile"),
	}
}

// ApplyFilters applies all filters to the query
func (qb *UserQueryBuilder) ApplyFilters(filters user.UserFilters) *UserQueryBuilder {
	// Basic filters
	if filters.Name != "" {
		qb.query = qb.query.Where("user.name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		qb.query = qb.query.Where("user.email ILIKE ?", "%"+filters.Email+"%")
	}
	if filters.Role != nil {
		qb.query = qb.query.Where("user.role = ?", string(*filters.Role))
	}
	if filters.IsActive != nil {
		qb.query = qb.query.Where("user.is_active = ?", *filters.IsActive)
	}
	if filters.OrganizationID != nil {
		qb.query = qb.query.Where("user.organization_id = ?", filters.OrganizationID.String())
	}

	// Date filters
	if filters.CreatedAfter != nil {
		qb.query = qb.query.Where("user.created_at > ?", *filters.CreatedAfter)
	}
	if filters.CreatedBefore != nil {
		qb.query = qb.query.Where("user.created_at < ?", *filters.CreatedBefore)
	}
	if filters.UpdatedAfter != nil {
		qb.query = qb.query.Where("user.updated_at > ?", *filters.UpdatedAfter)
	}
	if filters.UpdatedBefore != nil {
		qb.query = qb.query.Where("user.updated_at < ?", *filters.UpdatedBefore)
	}

	return qb
}

// ApplyPagination applies pagination to the query
func (qb *UserQueryBuilder) ApplyPagination(filters user.UserFilters) *UserQueryBuilder {
	if filters.Limit > 0 {
		qb.query = qb.query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		qb.query = qb.query.Offset(filters.Offset)
	}
	return qb
}

// ApplySorting applies sorting to the query
func (qb *UserQueryBuilder) ApplySorting(filters user.UserFilters) *UserQueryBuilder {
	orderBy := "created_at" // default column name
	if filters.OrderBy != "" {
		orderBy = filters.OrderBy
	}
	
	order := "DESC" // default
	if filters.Order != "" {
		order = strings.ToUpper(filters.Order)
		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}
	}
	
	// Use Order method with just the column name - Bun will handle the table alias automatically
	qb.query = qb.query.Order(fmt.Sprintf("%s %s", orderBy, order))
	return qb
}

// GetQuery returns the built query
func (qb *UserQueryBuilder) GetQuery() *bun.SelectQuery { 
	return qb.query
}

func NewUserRepository(db *bun.DB) user.Repository {
	return &UserRepository{
		db:     db,
		logger: logger.New(),
	}
}

// FindAll implements the unified user query with comprehensive filtering
func (r *UserRepository) FindAll(ctx context.Context, filters user.UserFilters) ([]*user.User, error) {
	query := r.db.NewSelect().
		Model((*models.User)(nil)).
		Relation("Profile")

	// Apply filters
	if filters.Name != "" {
		query = query.Where("\"user\".name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("\"user\".email ILIKE ?", "%"+filters.Email+"%")
	}
	if filters.Role != nil {
		query = query.Where("\"user\".role = ?", string(*filters.Role))
	}
	if filters.IsActive != nil {
		query = query.Where("\"user\".is_active = ?", *filters.IsActive)
	}
	if filters.OrganizationID != nil {
		query = query.Where("\"user\".organization_id = ?", filters.OrganizationID.String())
	}

	// Date filters
	if filters.CreatedAfter != nil {
		query = query.Where("\"user\".created_at > ?", *filters.CreatedAfter)
	}
	if filters.CreatedBefore != nil {
		query = query.Where("\"user\".created_at < ?", *filters.CreatedBefore)
	}
	if filters.UpdatedAfter != nil {
		query = query.Where("\"user\".updated_at > ?", *filters.UpdatedAfter)
	}
	if filters.UpdatedBefore != nil {
		query = query.Where("\"user\".updated_at < ?", *filters.UpdatedBefore)
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Apply sorting
	query = query.Order("created_at DESC")

	var models []*models.User
	err := query.Scan(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		domainUser, err := r.toDomainWithContext(ctx, model)
		if err != nil {
			return nil, err
		}
		users[i] = domainUser
	}

	return users, nil
}

// Count implements the unified count query with comprehensive filtering
func (r *UserRepository) Count(ctx context.Context, filters user.UserFilters) (int64, error) {
	query := r.db.NewSelect().
		Model((*models.User)(nil))

	// Apply filters
	if filters.Name != "" {
		query = query.Where("\"user\".name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("\"user\".email ILIKE ?", "%"+filters.Email+"%")
	}
	if filters.Role != nil {
		query = query.Where("\"user\".role = ?", string(*filters.Role))
	}
	if filters.IsActive != nil {
		query = query.Where("\"user\".is_active = ?", *filters.IsActive)
	}
	if filters.OrganizationID != nil {
		query = query.Where("\"user\".organization_id = ?", filters.OrganizationID.String())
	}

	// Date filters
	if filters.CreatedAfter != nil {
		query = query.Where("\"user\".created_at > ?", *filters.CreatedAfter)
	}
	if filters.CreatedBefore != nil {
		query = query.Where("\"user\".created_at < ?", *filters.CreatedBefore)
	}
	if filters.UpdatedAfter != nil {
		query = query.Where("\"user\".updated_at > ?", *filters.UpdatedAfter)
	}
	if filters.UpdatedBefore != nil {
		query = query.Where("\"user\".updated_at < ?", *filters.UpdatedBefore)
	}
	
	count, err := query.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return int64(count), nil
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
		Where("id = ?", id.String()).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.toDomainWithContext(ctx, model)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email shared.Email) (*user.User, error) {
	model := &models.User{}
	
	err := r.db.NewSelect().
		Model(model).
		Where("email = ?", email.String()).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.toDomainWithContext(ctx, model)
}

func (r *UserRepository) FindByOrganization(ctx context.Context, orgID shared.OrganizationID, filters user.UserFilters) ([]*user.User, error) {
	query := r.db.NewSelect().
		Model((*models.User)(nil)).
		Relation("Profile").
		Where("user.organization_id = ?", orgID.String())

	// Apply filters
	if filters.Name != "" {
		query = query.Where("user.name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("user.email ILIKE ?", "%"+filters.Email+"%")
	}
	if filters.Role != nil {
		query = query.Where("user.role = ?", string(*filters.Role))
	}
	if filters.IsActive != nil {
		query = query.Where("user.is_active = ?", *filters.IsActive)
	}
	if filters.CreatedAfter != nil {
		query = query.Where("user.created_at > ?", *filters.CreatedAfter)
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = query.Order("created_at DESC")

	var models []*models.User
	err := query.Scan(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		domainUser, err := r.toDomainWithContext(ctx, model)
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
		Where("user.role = ?", string(role))

	if orgID != nil {
		query = query.Where("user.organization_id = ?", orgID.String())
	}

	query = query.Order("created_at DESC")

	var models []*models.User
	err := query.Scan(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		domainUser, err := r.toDomainWithContext(ctx, model)
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
		Where("user.is_active = ?", true)

	if orgID != nil {
		query = query.Where("user.organization_id = ?", orgID.String())
	}

	query = query.Order("created_at DESC")

	var models []*models.User
	err := query.Scan(ctx, &models)
	if err != nil {
		return nil, fmt.Errorf("failed to find active users: %w", err)
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		domainUser, err := r.toDomainWithContext(ctx, model)
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
		PasswordHash: u.PasswordHash(),
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

func (r *UserRepository) toDomainWithContext(ctx context.Context, model *models.User) (*user.User, error) {
	// Validate critical fields that must be present
	userIDResult := safeValidateUserID(model.ID)
	if !userIDResult.IsValid() {
		r.logger.Error(ctx, "Critical validation failed: invalid user ID", 
			"user_id", model.ID, 
			"error", userIDResult.GetError())
		return nil, fmt.Errorf("invalid user ID: %w", userIDResult.GetError())
	}

	emailResult := safeValidateEmail(model.Email)
	if !emailResult.IsValid() {
		r.logger.Error(ctx, "Critical validation failed: invalid email", 
			"user_id", model.ID, 
			"email", model.Email, 
			"error", emailResult.GetError())
		return nil, fmt.Errorf("invalid email: %w", emailResult.GetError())
	}

	nameResult := safeValidateName(model.Name)
	if !nameResult.IsValid() {
		r.logger.Error(ctx, "Critical validation failed: invalid name", 
			"user_id", model.ID, 
			"name", model.Name, 
			"error", nameResult.GetError())
		return nil, fmt.Errorf("invalid name: %w", nameResult.GetError())
	}

	role := user.Role(model.Role)

	// Handle optional fields with graceful degradation
	var orgID *shared.OrganizationID
	if model.OrganizationID != nil {
		orgIDResult := safeValidateOrganizationID(*model.OrganizationID)
		if !orgIDResult.IsValid() {
			r.logger.Warn(ctx, "Non-critical validation failed: invalid organization ID", 
				"user_id", model.ID, 
				"organization_id", *model.OrganizationID, 
				"error", orgIDResult.GetError())
			// Continue with nil orgID instead of failing
			orgID = nil
		} else {
			id := orgIDResult.GetValue()
			orgID = &id
		}
	}

	var phone *shared.PhoneNumber
	if model.Phone != nil {
		phoneResult := safeValidatePhoneNumber(*model.Phone)
		if !phoneResult.IsValid() {
			r.logger.Warn(ctx, "Non-critical validation failed: invalid phone number", 
				"user_id", model.ID, 
				"email", model.Email, 
				"phone", *model.Phone, 
				"error", phoneResult.GetError())
			// Continue with nil phone instead of failing
			phone = nil
		} else {
			phone = phoneResult.GetValue()
		}
	}

	// Convert profile - always create one for doctors to avoid nil pointer issues
	var profile *user.Profile
	if model.Profile != nil {
		// Create a basic profile with default values for doctor-specific fields
		// TODO: Implement proper profile conversion when database schema is updated
		profile = user.NewProfile(userIDResult.GetValue())
		// Profile will have default values for doctor-specific fields
	} else {
		// Create a default profile even if none exists in database
		// This prevents nil pointer errors in the handler
		profile = user.NewProfile(userIDResult.GetValue())
	}

	return user.ReconstructUser(
		userIDResult.GetValue(),
		emailResult.GetValue(),
		nameResult.GetValue(),
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

