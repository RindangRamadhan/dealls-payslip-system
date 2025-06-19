package utils

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	Body struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
		Error   interface{} `json:"errors,omitempty"`
	}

	ValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	ValidationErrors []ValidationError

	// Used for swagger response body
	BodySuccess struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"-"`
	}

	// Used for swagger response body
	BodyFailure struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Errors  []ValidationError `json:"-"`
	}
)

// Default general error code
const ErrorCodeGeneralError = "GENERAL_ERROR"

func WriteSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func WriteError(c *fiber.Ctx, code int, message string, err error) error {
	var errVal ValidationErrors

	payload := Body{
		Success: false,
		Message: message,
	}

	switch e := err.(type) {
	case validator.ValidationErrors:
		trans, _ := GetTranslator("en")
		for _, v := range e {
			errVal = append(errVal, ValidationError{
				Field:   strings.ToLower(v.Field()),
				Message: v.Translate(trans),
			})
		}

		if len(errVal) > 0 {
			payload.Error = errVal
			payload.Message = "Bad request"
		}

	default:
		if err != nil {
			payload.Error = err.Error()
		}
	}

	return c.Status(code).JSON(payload)
}
