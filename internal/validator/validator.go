package validator

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

func Validate(data any) (ValidationErrors, bool) {
	validationErrors := ValidationErrors{
		Errors: make(map[string][]string),
	}

	errs := validate.Struct(data)

	ok := true

	if errs != nil {
		ok = false
		for _, err := range errs.(validator.ValidationErrors) {
			validationErrors.Errors[err.Field()] =
				append(
					validationErrors.Errors[err.Field()],
					err.Error(),
				)
		}
	}

	return validationErrors, ok
}
