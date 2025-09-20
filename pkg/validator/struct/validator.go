package structvalidator

import (
	"github.com/asaskevich/govalidator"
)

type Validator struct{}

func NewValidator() Validator {
	return Validator{}
}

// ValidateStruct is a method to validate a struct using govalidator.
func (v Validator) ValidateStruct(i any) (bool, error) {
	return govalidator.ValidateStruct(i)
}
