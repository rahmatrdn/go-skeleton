package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rahmatrdn/go-skeleton/entity"
)

var validate = validator.New()

func ValidateStruct(data interface{}) string {
	var errors []entity.ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element entity.ErrorResponse
			element.FailedField = err.StructField()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)

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
		}
	}
	jsonString, _ := json.Marshal(errors)

	return string(jsonString)
}
