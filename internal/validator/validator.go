package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wrapper do validator para uso no Fiber
type CustomValidator struct {
	Validator *validator.Validate
}

// New cria uma nova instância do validador customizado
func New() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

// Validate executa a validação e retorna os erros formatados
func (cv *CustomValidator) Validate(i interface{}) map[string]string {
	errors := make(map[string]string)

	if err := cv.Validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()

			switch tag {
			case "required":
				errors[field] = "O campo " + field + " é obrigatório"
			case "min":
				errors[field] = "O campo " + field + " deve ter no mínimo " + err.Param() + " caracteres"
			case "max":
				errors[field] = "O campo " + field + " deve ter no máximo " + err.Param() + " caracteres"
			case "gt":
				errors[field] = "O campo " + field + " deve ser maior que " + err.Param()
			case "gte":
				errors[field] = "O campo " + field + " deve ser maior ou igual a " + err.Param()
			case "lt":
				errors[field] = "O campo " + field + " deve ser menor que " + err.Param()
			case "lte":
				errors[field] = "O campo " + field + " deve ser menor ou igual a " + err.Param()
			case "email":
				errors[field] = "O campo " + field + " deve ser um email válido"
			default:
				errors[field] = "O campo " + field + " é inválido"
			}
		}
	}

	return errors
}
