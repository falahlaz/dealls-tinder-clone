package csvalidator

import (
	"errors"
	"reflect"

	"github.com/go-playground/locales/en"
	translator "github.com/go-playground/universal-translator"
	valid "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func addTranslationValidator(cv *valid.Validate) {
	cv.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
}

type customValidator struct {
	validator  *valid.Validate
	translator translator.Translator
}

func NewCustomValidator() *customValidator {
	newValidator := valid.New()

	en := en.New()

	uni := translator.New(en, en)

	// translate into language
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(newValidator, trans)

	// add custom translation
	addTranslationValidator(newValidator)

	return &customValidator{
		validator:  newValidator,
		translator: trans,
	}
}

func (cv *customValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)

	if err != nil {
		object, _ := err.(valid.ValidationErrors)

		for _, key := range object {
			return errors.New(key.Translate(cv.translator))
		}
	}
	return cv.validator.Struct(i)
}
