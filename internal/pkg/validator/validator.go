package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	en := en.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	_ = ent.RegisterDefaultTranslations(validate, trans)
}

// ValidateStruct validates the strcut 'v' using the validator instance and
// returns an error if the validation fails. The error message includes all the
// validation errors separated by commas.
func ValidateStruct(v any) error {
	err := validate.Struct(v)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		errMsgs := make([]string, len(errs))
		for i, e := range errs {
			errMsgs[i] = e.Translate(trans)
		}
		return errors.New(strings.Join(errMsgs, ", "))
	}
	return nil
}

// RegisterRules registers the custom validation rules for the given 'v'
func RegisterRules(v any, rules map[string]string) {
	validate.RegisterStructValidationMapRules(rules, v)
}
