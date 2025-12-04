package dto

// CreateProdutoRequest representa o payload para criação de um produto
// @Description Dados para criação de um novo produto
type CreateProdutoRequest struct {
	Codigo      string  `json:"codigo" validate:"required,min=1,max=50" example:"PROD001"`
	Descricao   string  `json:"descricao" validate:"required,min=3,max=255" example:"Produto de Exemplo"`
	Preco       float64 `json:"preco" validate:"required,gt=0" example:"99.90"`
	CategoriaID uint    `json:"categoria_id" validate:"required,gt=0" example:"1"`
}

// UpdateProdutoRequest representa o payload para atualização de um produto
// @Description Dados para atualização de um produto existente
type UpdateProdutoRequest struct {
	Codigo      string  `json:"codigo" validate:"omitempty,min=1,max=50" example:"PROD001"`
	Descricao   string  `json:"descricao" validate:"omitempty,min=3,max=255" example:"Produto Atualizado"`
	Preco       float64 `json:"preco" validate:"omitempty,gt=0" example:"149.90"`
	CategoriaID uint    `json:"categoria_id" validate:"omitempty,gt=0" example:"2"`
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

	// Dados da categoria associada
	CategoriaID uint               `json:"categoria_id" example:"1"`
	Categoria   *CategoriaResponse `json:"categoria,omitempty"`
}
