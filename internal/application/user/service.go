package user

import (
	"context"
	"fmt"
	"time"

	"medika-backend/internal/application/shared/events"
	"medika-backend/internal/domain/shared"
	"medika-backend/internal/domain/user"
	"medika-backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

// Application Service
type Service struct {
	userRepo user.Repository
	eventBus events.Bus
	logger   logger.Logger
}

func NewService(
	userRepo user.Repository,
	eventBus events.Bus,
	logger logger.Logger,
) *Service {
	return &Service{
		userRepo: userRepo,
		eventBus: eventBus,
		logger:   logger,
	}
}

// Commands
type CreateUserCommand struct {
	Name           string `json:"name" validate:"required,min=2,max=100"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Role           string `json:"role" validate:"required,oneof=admin doctor nurse patient cashier"`
	OrganizationID string `json:"organization_id,omitempty" validate:"omitempty,uuid"`
	Phone          string `json:"phone,omitempty" validate:"omitempty"`
}

type UpdateUserProfileCommand struct {
	UserID      string     `json:"user_id" validate:"required,uuid"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone       *string    `json:"phone,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Address     *string    `json:"address,omitempty"`
}

type UpdateMedicalInfoCommand struct {
	UserID           string   `json:"user_id" validate:"required,uuid"`
	EmergencyContact *string  `json:"emergency_contact,omitempty"`
	MedicalHistory   *string  `json:"medical_history,omitempty"`
	Allergies        []string `json:"allergies,omitempty"`
	BloodType        *string  `json:"blood_type,omitempty" validate:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
}

type UpdateAvatarCommand struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	AvatarURL string `json:"avatar_url" validate:"required,url"`
}

type LoginCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Responses
type UserResponse struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	OrganizationID *string   `json:"organizationId,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	AvatarURL      *string   `json:"avatar,omitempty"`
	IsActive       bool      `json:"isActive"`
	Profile        *ProfileResponse `json:"profile,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ProfileResponse struct {
	DateOfBirth      *time.Time `json:"dateOfBirth,omitempty"`
	Gender           *string    `json:"gender,omitempty"`
	Address          *string    `json:"address,omitempty"`
	EmergencyContact *string    `json:"emergencyContact,omitempty"`
	MedicalHistory   *string    `json:"medicalHistory,omitempty"`
	Allergies        []string   `json:"allergies,omitempty"`
	BloodType        *string    `json:"bloodType,omitempty"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// Use Cases
func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCommand) (*UserResponse, error) {
	s.logger.Info(ctx, "Creating new user", "email", cmd.Email, "role", cmd.Role)

	// Check if user already exists
	email, err := shared.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", cmd.Email)
	}

	// Create domain entity
	var orgID *string
	if cmd.OrganizationID != "" {
		orgID = &cmd.OrganizationID
	}

	newUser, err := user.NewUser(
		cmd.Email,
		cmd.Name,
		cmd.Password,
		user.Role(cmd.Role),
		orgID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Update phone if provided
	if cmd.Phone != "" {
		if err := newUser.UpdateProfile(nil, nil, nil, &cmd.Phone); err != nil {
			return nil, fmt.Errorf("failed to update phone: %w", err)
		}
	}

	// Save to repository
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Publish domain event
	event := user.UserCreatedEvent{
		UserID:         newUser.ID(),
		Email:          newUser.Email(),
		Role:           newUser.Role(),
		OrganizationID: newUser.OrganizationID(),
		CreatedAt:      newUser.CreatedAt(),
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error(ctx, "Failed to publish user created event", "error", err)
	}

	s.logger.Info(ctx, "User created successfully", "user_id", newUser.ID().String())

	return s.toResponse(newUser), nil
}

func (s *Service) UpdateUserProfile(ctx context.Context, cmd UpdateUserProfileCommand) (*UserResponse, error) {
	userID, err := shared.NewUserIDFromString(cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update profile
	if err := existingUser.UpdateProfile(
		cmd.DateOfBirth,
		cmd.Gender,
		cmd.Address,
		cmd.Phone,
	); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Save changes
	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info(ctx, "User profile updated", "user_id", userID.String())

	return s.toResponse(existingUser), nil
}

func (s *Service) UpdateMedicalInfo(ctx context.Context, cmd UpdateMedicalInfoCommand) (*UserResponse, error) {
	userID, err := shared.NewUserIDFromString(cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update medical info
	if err := existingUser.UpdateMedicalInfo(
		cmd.EmergencyContact,
		cmd.MedicalHistory,
		cmd.Allergies,
		cmd.BloodType,
	); err != nil {
		return nil, fmt.Errorf("failed to update medical info: %w", err)
	}

	// Save changes
	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info(ctx, "User medical info updated", "user_id", userID.String())

	return s.toResponse(existingUser), nil
}

func (s *Service) UpdateAvatar(ctx context.Context, cmd UpdateAvatarCommand) (*UserResponse, error) {
	s.logger.Info(ctx, "Updating user avatar", "user_id", cmd.UserID)

	// Parse user ID
	userID, err := shared.NewUserIDFromString(cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Find existing user
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update avatar
	if err := existingUser.UpdateAvatar(cmd.AvatarURL); err != nil {
		return nil, fmt.Errorf("failed to update avatar: %w", err)
	}

	// Save changes
	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info(ctx, "User avatar updated", "user_id", userID.String())

	return s.toResponse(existingUser), nil
}

func (s *Service) Login(ctx context.Context, cmd LoginCommand) (*LoginResponse, error) {
	s.logger.Info(ctx, "User login attempt", "email", cmd.Email)

	// Find user by email
	email, err := shared.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Warn(ctx, "Login failed", "email", cmd.Email)

		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !user.VerifyPassword(cmd.Password) {
		s.logger.Warn(ctx, "Login failed - invalid password", "email", cmd.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		s.logger.Warn(ctx, "Login failed - user inactive", "email", cmd.Email)
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate JWT token (this would be handled by auth service)
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info(ctx, "User logged in successfully", "user_id", user.ID().String())

	return &LoginResponse{
		User:  s.toResponse(user),
		Token: token,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (*UserResponse, error) {
	id, err := shared.NewUserIDFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.toResponse(user), nil
}

func (s *Service) GetUsersByOrganization(ctx context.Context, orgID string, filters user.UserFilters) ([]*UserResponse, error) {
	organizationID, err := shared.NewOrganizationID(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	users, err := s.userRepo.FindByOrganization(ctx, organizationID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = s.toResponse(u)
	}

	return responses, nil
}

// Helper methods
func (s *Service) toResponse(u *user.User) *UserResponse {
	response := &UserResponse{
		ID:        u.ID().String(),
		Email:     u.Email().String(),
		Name:      u.Name().String(),
		Role:      u.Role().String(),
		IsActive:  u.IsActive(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}

	if u.OrganizationID() != nil {
		orgID := u.OrganizationID().String()
		response.OrganizationID = &orgID
	}

	if u.Phone() != nil {
		phone := u.Phone().String()
		response.Phone = &phone
	}

	if u.AvatarURL() != nil {
		response.AvatarURL = u.AvatarURL()
	}

	if u.Profile() != nil {
		response.Profile = s.toProfileResponse(u.Profile())
	}

	return response
}

func (s *Service) toProfileResponse(p *user.Profile) *ProfileResponse {
	profile := &ProfileResponse{}

	if p.DateOfBirth() != nil {
		profile.DateOfBirth = p.DateOfBirth()
	}

	if p.Gender() != nil {
		gender := p.Gender().String()
		profile.Gender = &gender
	}

	if p.Address() != nil {
		profile.Address = p.Address()
	}

	if p.EmergencyContact() != nil {
		profile.EmergencyContact = p.EmergencyContact()
	}

	if p.MedicalHistory() != nil {
		profile.MedicalHistory = p.MedicalHistory()
	}

	if p.Allergies() != nil {
		profile.Allergies = p.Allergies()
	}

	if p.BloodType() != nil {
		bloodType := p.BloodType().String()
		profile.BloodType = &bloodType
	}

	return profile
}

func (s *Service) generateJWT(u *user.User) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id":         u.ID().String(),
		"email":           u.Email().String(),
		"name":            u.Name().String(),
		"role":            u.Role().String(),
		"organization_id": func() *string {
			if u.OrganizationID() != nil {
				orgID := u.OrganizationID().String()
				return &orgID
			}
			return nil
		}(),
		"is_active": u.IsActive(),
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	// In production, this should come from environment variables
	jwtSecret := []byte("your-super-secret-jwt-key-change-in-production")
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}
