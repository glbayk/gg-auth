package utils

import "github.com/go-playground/validator/v10"

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func TransformErrorMessage(fe validator.FieldError) ErrorMessage {
	switch fe.Tag() {
	case "required":
		return ErrorMessage{
			Field:   fe.Field(),
			Message: "This field is required",
		}
	case "email":
		return ErrorMessage{
			Field:   fe.Field(),
			Message: "Invalid email address",
		}
	case "min":
		return ErrorMessage{
			Field:   fe.Field(),
			Message: "Minimum length is " + fe.Param(),
		}
	case "max":
		return ErrorMessage{
			Field:   fe.Field(),
			Message: "Maximum length is " + fe.Param(),
		}
	}

	return ErrorMessage{
		Field:   fe.Field(),
		Message: fe.Error(),
	}
}
