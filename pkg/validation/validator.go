package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

func NewValidator(_ *do.Injector) (*validator.Validate, error) {
	return validator.New(), nil
}
