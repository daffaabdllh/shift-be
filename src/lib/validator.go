package lib

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator mendaftarkan fungsi penamaan kustom agar validator mengembalikan nama tag JSON, bukan nama field Struct Go.
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name == "" {
				return fld.Name
			}
			return name
		})
	}
}

func init() {
	InitValidator()
}

type ValidationErrorDetail struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Status string                  `json:"status"`
	Error  []ValidationErrorDetail `json:"error"`
}

// FormatValidationError mengubah error validator.ValidationErrors menjadi format kustom yang rapi.
func FormatValidationError(err error) (int, interface{}) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ValidationErrorDetail, len(ve))
		for i, fe := range ve {
			out[i] = ValidationErrorDetail{
				Key:     fe.Field(),
				Message: getErrorMsg(fe),
			}
		}
		return 400, ValidationErrorResponse{
			Status: "validation error",
			Error:  out,
		}
	}

	// Fallback jika terjadi error parsing JSON lainnya (misal payload JSON rusak)
	return 400, map[string]any{
		"error": err.Error(),
	}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Must be at least " + fe.Param() + " characters"
	case "max":
		return "Must be at most " + fe.Param() + " characters"
	case "eqfield":
		return "Must match the password field"
	}
	return "Field validation failed on '" + fe.Tag() + "' tag"
}
