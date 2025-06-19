package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	eng := en.New()
	uni = ut.New(eng, eng)

	var found bool
	trans, found = uni.GetTranslator("en")
	if !found {
		panic("translator for 'en' not found")
	}

	validate = validator.New()

	// Use JSON tag as field name in errors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(fmt.Sprintf("failed to register translations: %v", err))
	}
}

// ValidateStruct validates struct and returns translated errors
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Return only the first error as a simple message
	if errs, ok := err.(validator.ValidationErrors); ok {
		return fmt.Errorf(errs[0].Translate(trans))
	}

	return err
}

// GetTranslator returns a translator by language tag
func GetTranslator(lang string) (ut.Translator, error) {
	trans, _ := uni.GetTranslator(lang)
	return trans, nil
}
