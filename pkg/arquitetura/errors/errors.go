package errors

import (
	"errors"
	"fmt"
)

// Erros padrão da arquitetura
var (
	// ErrNotFound indica que o registro não foi encontrado
	ErrNotFound = errors.New("registro não encontrado")

	// ErrDuplicateKey indica violação de chave única
	ErrDuplicateKey = errors.New("registro já existe com esta chave")

	// ErrValidation indica erro de validação
	ErrValidation = errors.New("erro de validação")

	// ErrForeignKeyViolation indica violação de chave estrangeira
	ErrForeignKeyViolation = errors.New("violação de chave estrangeira")

	// ErrHasRelatedRecords indica que existem registros relacionados
	ErrHasRelatedRecords = errors.New("existem registros relacionados que impedem a operação")

	// ErrInactiveRecord indica que o registro está inativo
	ErrInactiveRecord = errors.New("registro está inativo")

	// ErrInvalidID indica ID inválido
	ErrInvalidID = errors.New("ID inválido")

	// ErrInternalServer indica erro interno do servidor
	ErrInternalServer = errors.New("erro interno do servidor")
)

// BusinessError representa um erro de negócio customizado
type BusinessError struct {
	Code    string
	Message string
	Field   string
	Err     error
}

// Error implementa a interface error
func (e *BusinessError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Field, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap permite que errors.Is e errors.As funcionem
func (e *BusinessError) Unwrap() error {
	return e.Err
}

// NewBusinessError cria um novo erro de negócio
func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// NewFieldError cria um erro de negócio associado a um campo
func NewFieldError(code, field, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Field:   field,
	}
}

// WrapError encapsula um erro existente
func WrapError(err error, code, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ValidationErrors representa uma coleção de erros de validação
type ValidationErrors struct {
	Errors map[string]string
}

// Error implementa a interface error
func (v *ValidationErrors) Error() string {
	return "erros de validação"
}

// Add adiciona um erro de validação
func (v *ValidationErrors) Add(field, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	v.Errors[field] = message
}

// HasErrors retorna true se há erros
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// NewValidationErrors cria uma nova instância de ValidationErrors
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make(map[string]string),
	}
}

// IsNotFound verifica se o erro é do tipo "não encontrado"
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsDuplicateKey verifica se o erro é do tipo "chave duplicada"
func IsDuplicateKey(err error) bool {
	return errors.Is(err, ErrDuplicateKey)
}

// IsValidation verifica se o erro é do tipo "validação"
func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

// IsBusinessError verifica se é um erro de negócio
func IsBusinessError(err error) bool {
	var be *BusinessError
	return errors.As(err, &be)
}

// GetBusinessError extrai o BusinessError se existir
func GetBusinessError(err error) (*BusinessError, bool) {
	var be *BusinessError
	if errors.As(err, &be) {
		return be, true
	}
	return nil, false
}

