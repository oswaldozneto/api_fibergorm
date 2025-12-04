package service

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// ValidationContext contém o contexto para validações
type ValidationContext struct {
	Context   context.Context
	Operation OperationType
	EntityID  uint // ID da entidade (para updates)
}

// OperationType define o tipo de operação sendo validada
type OperationType string

const (
	OperationCreate OperationType = "create"
	OperationUpdate OperationType = "update"
	OperationDelete OperationType = "delete"
)

// ValidationResult representa o resultado de uma validação
type ValidationResult struct {
	Errors map[string]string
}

// NewValidationResult cria um novo resultado de validação
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Errors: make(map[string]string),
	}
}

// AddError adiciona um erro de validação
func (v *ValidationResult) AddError(field, message string) {
	v.Errors[field] = message
}

// HasErrors retorna true se há erros
func (v *ValidationResult) HasErrors() bool {
	return len(v.Errors) > 0
}

// Merge combina dois resultados de validação
func (v *ValidationResult) Merge(other *ValidationResult) {
	if other == nil {
		return
	}
	for field, message := range other.Errors {
		v.Errors[field] = message
	}
}

// EntityValidator é a interface para validações específicas de uma entidade
// As entidades podem implementar esta interface para fornecer validações customizadas
type EntityValidator[E any, CreateReq any, UpdateReq any] interface {
	// ValidateCreate valida a criação de uma entidade
	ValidateCreate(ctx *ValidationContext, req *CreateReq) *ValidationResult

	// ValidateUpdate valida a atualização de uma entidade
	ValidateUpdate(ctx *ValidationContext, entity *E, req *UpdateReq) *ValidationResult

	// ValidateDelete valida a exclusão de uma entidade
	ValidateDelete(ctx *ValidationContext, entity *E) *ValidationResult
}

// NoOpValidator é um validador que não faz nada (para entidades sem validações customizadas)
type NoOpValidator[E any, CreateReq any, UpdateReq any] struct{}

func (n *NoOpValidator[E, CreateReq, UpdateReq]) ValidateCreate(ctx *ValidationContext, req *CreateReq) *ValidationResult {
	return nil
}

func (n *NoOpValidator[E, CreateReq, UpdateReq]) ValidateUpdate(ctx *ValidationContext, entity *E, req *UpdateReq) *ValidationResult {
	return nil
}

func (n *NoOpValidator[E, CreateReq, UpdateReq]) ValidateDelete(ctx *ValidationContext, entity *E) *ValidationResult {
	return nil
}

// StructValidator valida structs usando tags de validação
type StructValidator struct {
	validate *validator.Validate
}

// NewStructValidator cria um novo validador de structs
func NewStructValidator() *StructValidator {
	return &StructValidator{
		validate: validator.New(),
	}
}

// Validate valida uma struct e retorna os erros formatados
func (sv *StructValidator) Validate(i interface{}) map[string]string {
	errors := make(map[string]string)

	if err := sv.validate.Struct(i); err != nil {
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

// ToValidationResult converte erros de validação para ValidationResult
func (sv *StructValidator) ToValidationResult(i interface{}) *ValidationResult {
	errors := sv.Validate(i)
	if len(errors) == 0 {
		return nil
	}

	result := NewValidationResult()
	result.Errors = errors
	return result
}
