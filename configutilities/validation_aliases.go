package configutilities

import "github.com/go-playground/validator/v10"

func registerValidationAliases(validate *validator.Validate) *validator.Validate {
	validate.RegisterAlias("sslmode", "oneof=disable allow prefer require verify-ca verify-full")
	validate.RegisterAlias("inhouseports", "gte=1024,lte=49151")
	validate.RegisterAlias("loglevel", "oneof=DEBUG INFO WARN ERROR")
	validate.RegisterAlias("ginmode", "oneof=debug release test")

	return validate
}
