package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Validator handles struct validation
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Validate validates a struct and returns formatted errors
func (v *Validator) Validate(s any) error {
	if err := v.validate.Struct(s); err != nil {
		return v.formatValidationErr(err)
	}
	return nil
}

// formatValidationErr formats validation errors into a readable format
func (v *Validator) formatValidationErr(err error) error {
	var errMsg string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			// Custom error messages based on tag
			switch tag {
			case "required":
				errMsg += fmt.Sprintf("%s is required; ", field)
			case "email":
				errMsg += fmt.Sprintf("%s must be a valid email; ", field)
			case "min":
				errMsg += fmt.Sprintf("%s must be at least %s characters; ", field, e.Param())
			case "max":
				errMsg += fmt.Sprintf("%s must be at most %s characters; ", field, e.Param())
			case "alphanum":
				errMsg += fmt.Sprintf("%s must contain only alphanumeric characters; ", field)
			case "oneof":
				errMsg += fmt.Sprintf("%s must be one of: %s; ", field, e.Param())
			case "uuid":
				errMsg += fmt.Sprintf("%s must be a valid UUID; ", field)
			default:
				errMsg += fmt.Sprintf("%s failed validation: %s; ", field, tag)
			}
		}
		return errors.New(errMsg)
	}

	return err
}

// ValidateStruct is a convenience function to validate a struct
func ValidateStruct(s any) error {
	v := NewValidator()
	return v.Validate(s)
}

// GetField gets a field value by name from a struct
func GetField(s any, field string) any {
	r := reflect.ValueOf(s)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil
	}
	return f.Interface()
}
