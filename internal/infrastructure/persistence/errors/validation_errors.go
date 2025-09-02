package errors

import (
	"fmt"
)

// ValidationError represents a validation error that should not break critical operations
type ValidationError struct {
	Field   string
	Value   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s' with value '%s': %s", e.Field, e.Value, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(field, value, message string, err error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Err:     err,
	}
}

// ValidationResult holds the result of a validation operation
type ValidationResult[T any] struct {
	Value T
	Error *ValidationError
}

// NewValidationResult creates a successful validation result
func NewValidationResult[T any](value T) ValidationResult[T] {
	return ValidationResult[T]{Value: value}
}

// NewValidationResultWithError creates a failed validation result
func NewValidationResultWithError[T any](err *ValidationError) ValidationResult[T] {
	return ValidationResult[T]{Error: err}
}

// IsValid returns true if the validation was successful
func (vr ValidationResult[T]) IsValid() bool {
	return vr.Error == nil
}

// GetValue returns the value if validation was successful, or the zero value if not
func (vr ValidationResult[T]) GetValue() T {
	return vr.Value
}

// GetError returns the validation error if any
func (vr ValidationResult[T]) GetError() *ValidationError {
	return vr.Error
}
