package usecase

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	"github.com/rahmatrdn/go-skeleton/entity"
)

var validate = validator.New()
var uni *ut.UniversalTranslator

// func init() {
// 	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
// 		return field.Tag.Get("name")
// 	})
// }

func ValidateStructProcess(data interface{}) []entity.ErrorResponse {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	})

	id := id.New()
	uni = ut.New(id, id)

	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validate, trans)

	var errors []entity.ErrorResponse
	err := validate.Struct(data)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element entity.ErrorResponse
			element.FailedField = err.StructField()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Message = err.Translate(trans)

			// Available responses
			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println(err.Translate(trans))

			errors = append(errors, element)
		}
	}

	return errors
}

func ValidateStruct(data interface{}) string {
	validate := ValidateStructProcess(data)
	jsonString, _ := json.Marshal(validate)

	if string(jsonString) != "null" {
		return string(jsonString) + "XX"
	}

	return ""
}
