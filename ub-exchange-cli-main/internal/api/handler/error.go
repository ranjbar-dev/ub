package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)


type validationError struct {
	field   string
	message string
}

func customizeValidationError(validationErrors validator.ValidationErrors) []validationError {
	var errs []validationError
	for _, e := range validationErrors {
		err := ""
		switch e.Tag() {
		case "required":
			err = fmt.Sprintf("%s is required", e.Field())
		case "oneof":
			err = fmt.Sprintf("%s:value %s is not valid ", e.Field(), e.Value())
		case "gt":
			err = fmt.Sprintf("%s:must be greater than %s", e.Field(), e.Param())
		case "email":
			err = fmt.Sprintf("%s must be a valid email address", e.Field())
		case "min":
			err = fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
		case "max":
			err = fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
		case "len":
			err = fmt.Sprintf("%s must be exactly %s characters", e.Field(), e.Param())
		case "gte":
			err = fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
		case "lte":
			err = fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
		case "numeric":
			err = fmt.Sprintf("%s must be a number", e.Field())
		case "alphanum":
			err = fmt.Sprintf("%s must contain only letters and numbers", e.Field())
		case "uuid":
			err = fmt.Sprintf("%s must be a valid UUID", e.Field())
		case "url":
			err = fmt.Sprintf("%s must be a valid URL", e.Field())
		case "eqfield":
			err = fmt.Sprintf("%s must match %s", e.Field(), e.Param())
		default:
			err = fmt.Sprintf("field must be type  of %s", e.Tag())
		}

		ve := validationError{
			field:   e.Field(),
			message: err,
		}
		errs = append(errs, ve)
	}

	return errs
}

func HandleValidationError(err error) (errResponse APIResponse, StatusCode int) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := customizeValidationError(validationErrors)
		firstErrorMessage := errorMessages[0]
		result := APIResponse{
			Status:  false,
			Message: firstErrorMessage.message,
			Data:    make(map[string]string),
		}
		return result, http.StatusUnprocessableEntity
	}

	result := APIResponse{
		Status:  false,
		Message: "something went wrong",
		Data:    map[string]string{},
	}
	return result, http.StatusInternalServerError
}
