package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

const (
	defaultLimit = 20
)

type Validator interface {
	Validate(i interface{}) error
}

type valid struct {
	v *validator.Validate
}

func NewValidator() (Validator, error) {
	v := validator.New()
	if err := v.RegisterValidation("value", valueValidate); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("limit", limitValidate); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("category", categoryValidate); err != nil {
		return nil, err
	}

	return &valid{v: v}, nil
}

func (v *valid) Validate(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		return validateError(err.(validator.ValidationErrors)[0])
	}
	return nil
}

func validateError(err validator.FieldError) error {
	switch err.Tag() {
	case "value":
		return fmt.Errorf("field %s is incorrect. Value must be >= 0, got %d", err.Field(), err.Value())
	case "limit":
		return fmt.Errorf("field %s is incorrect. Value must be >= 0 and <= 20, got %d", err.Field(), err.Value())
	case "category":
		return fmt.Errorf("field %s is incorrect. All id's must be > 0", err.Field())
	default:
		return fmt.Errorf("field %s is required", err.Field())
	}
}

func valueValidate(fl validator.FieldLevel) bool {
	return fl.Field().Int() >= 0
}

func limitValidate(fl validator.FieldLevel) bool {
	v := fl.Field().Int()
	return v > 0 && v <= defaultLimit
}

func categoryValidate(fl validator.FieldLevel) bool {
	categories, ok := fl.Field().Interface().([]int)
	if !ok {
		return false
	}
	for _, c := range categories {
		if c <= 0 {
			return false
		}
	}
	return true
}
