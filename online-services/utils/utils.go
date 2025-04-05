package utils

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notblank", NotBlank)
	}
}

func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return strings.TrimSpace(field) != ""
}