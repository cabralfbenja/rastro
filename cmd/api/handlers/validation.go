package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"rastro/common"
	"reflect"
	"strconv"
	"strings"
)

func (h *Handler) ValidateBodyRequest(c echo.Context, payload interface{}) []*common.ValidationError {
	var errors []*common.ValidationError
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)

	if ok {
		reflected := reflect.ValueOf(payload)
		for _, validationErr := range validationErrors {
			field, _ := reflected.Type().FieldByName(validationErr.StructField())

			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(validationErr.StructField())
			}

			condition := validationErr.Tag()

			errMessage := key + " field is " + condition

			keyToTitleCase := strings.Replace(key, "_", " ", -1)
			param := validationErr.Param()

			switch condition {
			case "required":
				errMessage = keyToTitleCase + " is required"
			case "email":
				errMessage = keyToTitleCase + " must be a valid email address"
			case "min":
				if _, err := strconv.Atoi(param); err == nil {
					errMessage = fmt.Sprintf("%s must be at least %s characters", key, param)
				}
			case "eqfield":
				errMessage = fmt.Sprintf("%s must match %s", keyToTitleCase, param)
			}

			currentValidationError := common.ValidationError{
				Error:     errMessage,
				Key:       keyToTitleCase,
				Condition: condition,
			}
			errors = append(errors, &currentValidationError)
		}
	}

	return errors
}
