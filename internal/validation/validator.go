package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Validator provides validation functions
type Validator struct {
	errors ValidationErrors
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns all validation errors
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// Required validates that a field is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field, value string) *Validator {
	if value != "" {
		if _, err := mail.ParseAddress(value); err != nil {
			v.AddError(field, "must be a valid email address")
		}
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("must be at least %d characters long", min))
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("must be no more than %d characters long", max))
	}
	return v
}

// Range validates that an integer is within a range
func (v *Validator) Range(field string, value, min, max int) *Validator {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf("must be between %d and %d", min, max))
	}
	return v
}

// OneOf validates that a value is one of the allowed values
func (v *Validator) OneOf(field, value string, allowed []string) *Validator {
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.AddError(field, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")))
	return v
}

// Pattern validates that a string matches a regex pattern
func (v *Validator) Pattern(field, value, pattern string) *Validator {
	if value != "" {
		matched, err := regexp.MatchString(pattern, value)
		if err != nil {
			v.AddError(field, "invalid pattern validation")
		} else if !matched {
			v.AddError(field, "format is invalid")
		}
	}
	return v
}

// ArrayNotEmpty validates that an array is not empty
func (v *Validator) ArrayNotEmpty(field string, value []string) *Validator {
	if len(value) == 0 {
		v.AddError(field, "must contain at least one item")
	}
	return v
}

// Validate runs all validations and returns error if any fail
func (v *Validator) Validate() error {
	if v.HasErrors() {
		return v.errors
	}
	return nil
}

// Quick validation functions for common use cases

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	validator := NewValidator()
	validator.Required("email", email).Email("email", email)
	return validator.Validate()
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	validator := NewValidator()
	validator.Required("password", password).MinLength("password", password, 8)
	return validator.Validate()
}

// ValidateRequired validates that all required fields are present
func ValidateRequired(fields map[string]string) error {
	validator := NewValidator()
	for field, value := range fields {
		validator.Required(field, value)
	}
	return validator.Validate()
}
