package wallet

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var IsPositiveDecimal validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if ok {
		match1, _ := regexp.MatchString(`^\d+(\.\d{1,2})?$`, data)
		match2, _ := regexp.MatchString(`^(\d+)?\.\d{1,2}?$`, data)
		return match1 || match2
	}
	return false
}
