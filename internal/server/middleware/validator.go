package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 4 || len(username) > 24 {
		return false
	}

	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}

	return true
}

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	validate := validator.New()
	validate.RegisterValidation("username", validateUsername)

	return &CustomValidator{
		validator: validate,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)

	switch err.(type) {
	case validator.ValidationErrors:
		return err
	case nil:
		return nil
	default:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
}
