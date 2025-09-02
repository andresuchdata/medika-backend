package repositories

import (
	"medika-backend/internal/domain/shared"
	"medika-backend/internal/infrastructure/persistence/errors"
)

// safeValidatePhoneNumber validates a phone number and returns a result that allows graceful degradation
func safeValidatePhoneNumber(phone string) errors.ValidationResult[*shared.PhoneNumber] {
	if phone == "" {
		return errors.NewValidationResult[*shared.PhoneNumber](nil)
	}
	
	pn, err := shared.NewPhoneNumber(phone)
	if err != nil {
		validationErr := errors.NewValidationError(
			"phone",
			phone,
			"invalid phone number format",
			err,
		)
		return errors.NewValidationResultWithError[*shared.PhoneNumber](validationErr)
	}
	
	return errors.NewValidationResult(&pn)
}

// safeValidateEmail validates an email and returns a result
func safeValidateEmail(email string) errors.ValidationResult[shared.Email] {
	em, err := shared.NewEmail(email)
	if err != nil {
		validationErr := errors.NewValidationError(
			"email",
			email,
			"invalid email format",
			err,
		)
		return errors.NewValidationResultWithError[shared.Email](validationErr)
	}
	
	return errors.NewValidationResult(em)
}

// safeValidateName validates a name and returns a result
func safeValidateName(name string) errors.ValidationResult[shared.Name] {
	n, err := shared.NewName(name)
	if err != nil {
		validationErr := errors.NewValidationError(
			"name",
			name,
			"invalid name format",
			err,
		)
		return errors.NewValidationResultWithError[shared.Name](validationErr)
	}
	
	return errors.NewValidationResult(n)
}

// safeValidateUserID validates a user ID and returns a result
func safeValidateUserID(id string) errors.ValidationResult[shared.UserID] {
	userID, err := shared.NewUserIDFromString(id)
	if err != nil {
		validationErr := errors.NewValidationError(
			"user_id",
			id,
			"invalid user ID format",
			err,
		)
		return errors.NewValidationResultWithError[shared.UserID](validationErr)
	}
	
	return errors.NewValidationResult(userID)
}

// safeValidateOrganizationID validates an organization ID and returns a result
func safeValidateOrganizationID(id string) errors.ValidationResult[shared.OrganizationID] {
	orgID, err := shared.NewOrganizationID(id)
	if err != nil {
		validationErr := errors.NewValidationError(
			"organization_id",
			id,
			"invalid organization ID format",
			err,
		)
		return errors.NewValidationResultWithError[shared.OrganizationID](validationErr)
	}
	
	return errors.NewValidationResult(orgID)
}
