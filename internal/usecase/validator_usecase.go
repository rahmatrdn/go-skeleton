package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	"github.com/rahmatrdn/go-skeleton/entity"
)

// ValidatorUsecase is a usecase for custom validating data
type Validator struct {
}

func NewValidatorUsecase() *Validator {
	return &Validator{}
}

type ValidatorUsecase interface {
	ValidateWithMessage(data interface{}) string
	Validate(data interface{}) []entity.ErrorResponse
}

var validation = validator.New()

func init() {
	validation.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	})
}

func (v *Validator) ValidateWithMessage(data interface{}) string {
	validate := v.Validate(data)
	jsonString, _ := json.Marshal(validate)

	return string(jsonString) + "XX"
}

func (v *Validator) Validate(data interface{}) []entity.ErrorResponse {
	id := id.New()
	uni = ut.New(id, id)

	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validation, trans)

	var errors []entity.ErrorResponse
	err := validation.Struct(data)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element entity.ErrorResponse
			element.FailedField = err.StructField()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Message = err.Translate(trans)

			// Available responses
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println(err.Translate(trans))

			errors = append(errors, element)
		}
	}

	return errors
}
