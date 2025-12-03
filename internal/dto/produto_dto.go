package dto

// CreateProdutoRequest representa o payload para criação de um produto
// @Description Dados para criação de um novo produto
type CreateProdutoRequest struct {
	Codigo    string  `json:"codigo" validate:"required,min=1,max=50" example:"PROD001"`
	Descricao string  `json:"descricao" validate:"required,min=3,max=255" example:"Produto de Exemplo"`
	Preco     float64 `json:"preco" validate:"required,gt=0" example:"99.90"`
}

// UpdateProdutoRequest representa o payload para atualização de um produto
// @Description Dados para atualização de um produto existente
type UpdateProdutoRequest struct {
	Codigo    string  `json:"codigo" validate:"omitempty,min=1,max=50" example:"PROD001"`
	Descricao string  `json:"descricao" validate:"omitempty,min=3,max=255" example:"Produto Atualizado"`
	Preco     float64 `json:"preco" validate:"omitempty,gt=0" example:"149.90"`
}

// ProdutoResponse representa a resposta de um produto
// @Description Dados de resposta de um produto
type ProdutoResponse struct {
	ID        uint    `json:"id" example:"1"`
	Codigo    string  `json:"codigo" example:"PROD001"`
	Descricao string  `json:"descricao" example:"Produto de Exemplo"`
	Preco     float64 `json:"preco" example:"99.90"`
	CreatedAt string  `json:"created_at" example:"2024-01-01 10:00:00"`
	UpdatedAt string  `json:"updated_at" example:"2024-01-01 10:00:00"`
}

// ErrorResponse representa uma resposta de erro da API
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

// PaginatedResponse representa uma resposta paginada
// @Description Resposta paginada com lista de produtos
type PaginatedResponse struct {
	Data       []ProdutoResponse `json:"data"`
	Total      int64             `json:"total" example:"100"`
	Page       int               `json:"page" example:"1"`
	PageSize   int               `json:"page_size" example:"10"`
	TotalPages int               `json:"total_pages" example:"10"`
}
