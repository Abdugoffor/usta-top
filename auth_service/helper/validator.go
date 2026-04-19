package helper

import "github.com/go-playground/validator/v10"

var _v = validator.New()

func ValidateStruct(s interface{}) map[string]string {
	err := _v.Struct(s)
	if err == nil {
		return nil
	}
	errs := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errs[e.Field()] = e.Tag()
	}
	return errs
}
