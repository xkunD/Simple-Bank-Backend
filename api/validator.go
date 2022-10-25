package api

import (
	"go-simple-bank/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if ok {
		// check if currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}
