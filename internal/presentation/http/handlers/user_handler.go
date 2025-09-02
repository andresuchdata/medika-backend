package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	userApp "medika-backend/internal/application/user"
	"medika-backend/internal/domain/user"
	"medika-backend/internal/presentation/http/dto"
	"medika-backend/pkg/logger"
)

type UserHandler struct {
	userService *userApp.Service
	validator   *validator.Validate
	logger      logger.Logger
}

func NewUserHandler(
	userService *userApp.Service,
	validator *validator.Validate,
	logger logger.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
		logger:      logger,
	}
}

// POST /api/v1/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid JSON format",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	cmd := userApp.CreateUserCommand{
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		Role:           req.Role,
		OrganizationID: req.OrganizationID,
		Phone:          req.Phone,
	}

	response, err := h.userService.CreateUser(c.Context(), cmd)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to create user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to create user",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "User created successfully",
	})
}

// POST /api/v1/auth/login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid JSON format",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	cmd := userApp.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.userService.Login(c.Context(), cmd)
	if err != nil {
		h.logger.Warn(c.Context(), "Login failed", "email", req.Email, "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:   "Login failed",
			Message: "Invalid credentials",
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "Login successful",
	})
}

// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "User ID is required",
			Message: "User ID parameter is missing",
		})
	}

	response, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get user", "user_id", userID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    response,
	})
}

// PUT /api/v1/users/:id/profile
func (h *UserHandler) UpdateUserProfile(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "User ID is required",
			Message: "User ID parameter is missing",
		})
	}

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid JSON format",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	cmd := userApp.UpdateUserProfileCommand{
		UserID:      userID,
		Name:        req.Name,
		Phone:       req.Phone,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
		Address:     req.Address,
	}

	response, err := h.userService.UpdateUserProfile(c.Context(), cmd)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to update user profile", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update profile",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "Profile updated successfully",
	})
}

// PUT /api/v1/users/:id/medical-info
func (h *UserHandler) UpdateMedicalInfo(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "User ID is required",
			Message: "User ID parameter is missing",
		})
	}

	var req dto.UpdateMedicalInfoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid JSON format",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	cmd := userApp.UpdateMedicalInfoCommand{
		UserID:           userID,
		EmergencyContact: req.EmergencyContact,
		MedicalHistory:   req.MedicalHistory,
		Allergies:        req.Allergies,
		BloodType:        req.BloodType,
	}

	response, err := h.userService.UpdateMedicalInfo(c.Context(), cmd)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to update medical info", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update medical info",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "Medical information updated successfully",
	})
}

// PUT /api/v1/users/:id/avatar
func (h *UserHandler) UpdateAvatar(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "User ID is required",
			Message: "user ID parameter is missing",
		})
	}

	var req dto.UpdateAvatarRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid JSON format",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	// Get current user ID from JWT token
	currentUserID := c.Locals("user_id")
	if currentUserID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	// Check if user is updating their own avatar or if they're an admin
	currentUser, err := h.userService.GetUserByID(c.Context(), currentUserID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get current user",
			Message: err.Error(),
		})
	}

	// Only allow users to update their own avatar, or admins to update any avatar
	if currentUser.ID != userID && currentUser.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
			Error:   "Forbidden",
			Message: "You can only update your own avatar",
		})
	}

	cmd := userApp.UpdateAvatarCommand{
		UserID:    userID,
		AvatarURL: req.AvatarURL,
	}

	response, err := h.userService.UpdateAvatar(c.Context(), cmd)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to update avatar", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to update avatar",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "Avatar updated successfully",
	})
}

// GET /api/v1/users?organization_id=:org_id
func (h *UserHandler) GetUsersByOrganization(c *fiber.Ctx) error {
	orgID := c.Query("organization_id")
	if orgID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Organization ID is required",
			Message: "organization_id query parameter is missing",
		})
	}

	// Parse filters
	filters := user.UserFilters{}
	
	if name := c.Query("name"); name != "" {
		filters.Name = name
	}
	
	if email := c.Query("email"); email != "" {
		filters.Email = email
	}
	
	if roleStr := c.Query("role"); roleStr != "" {
		role := user.Role(roleStr)
		filters.Role = &role
	}
	
	if activeStr := c.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filters.IsActive = &active
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	} else {
		filters.Limit = 50 // Default limit
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	responses, err := h.userService.GetUsersByOrganization(c.Context(), orgID, filters)
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get users by organization", "org_id", orgID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to get users",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    responses,
	})
}

// GET /api/v1/users/me (get current user)
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	// Get user ID from JWT token (set by auth middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	response, err := h.userService.GetUserByID(c.Context(), userID.(string))
	if err != nil {
		h.logger.Error(c.Context(), "Failed to get current user", "user_id", userID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    response,
	})
}
