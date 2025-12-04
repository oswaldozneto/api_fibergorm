package dto

// ErrorResponse representa uma resposta de erro padrão da API
// @Description Resposta de erro padrão da API
type ErrorResponse struct {
	Error   string            `json:"error" example:"Erro de validação"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse representa uma resposta de sucesso genérica
// @Description Resposta de sucesso padrão da API
type SuccessResponse struct {
	Message string `json:"message" example:"Operação realizada com sucesso"`
}

// PaginatedResponse representa uma resposta paginada genérica
// @Description Resposta paginada com lista de itens
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total" example:"100"`
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalPages int   `json:"total_pages" example:"10"`
}

// NewPaginatedResponse cria uma resposta paginada
func NewPaginatedResponse[T any](data []T, total int64, page, pageSize int) *PaginatedResponse[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	return &PaginatedResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// CreateRequest é a interface que todos os DTOs de criação devem implementar
type CreateRequest interface {
	ToEntity() interface{}
}

// UpdateRequest é a interface que todos os DTOs de atualização devem implementar
type UpdateRequest interface {
	ApplyTo(entity interface{}) error
}

// Response é a interface base para respostas
type Response interface {
	FromEntity(entity interface{}) Response
}

// Mapper é uma interface genérica para mapeamento entre entidades e DTOs
type Mapper[E any, CreateReq any, UpdateReq any, Resp any] interface {
	ToEntity(req *CreateReq) *E
	ToResponse(entity *E) *Resp
	ApplyUpdate(entity *E, req *UpdateReq)
}

