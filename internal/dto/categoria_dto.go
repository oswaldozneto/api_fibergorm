package dto

// CreateCategoriaRequest representa o payload para criação de uma categoria
// @Description Dados para criação de uma nova categoria
type CreateCategoriaRequest struct {
	Nome      string `json:"nome" validate:"required,min=2,max=100" example:"Eletrônicos"`
	Descricao string `json:"descricao" validate:"omitempty,max=255" example:"Produtos eletrônicos em geral"`
	Ativo     *bool  `json:"ativo" validate:"omitempty" example:"true"`
}

// UpdateCategoriaRequest representa o payload para atualização de uma categoria
// @Description Dados para atualização de uma categoria existente
type UpdateCategoriaRequest struct {
	Nome      string `json:"nome" validate:"omitempty,min=2,max=100" example:"Eletrônicos"`
	Descricao string `json:"descricao" validate:"omitempty,max=255" example:"Produtos eletrônicos atualizados"`
	Ativo     *bool  `json:"ativo" validate:"omitempty" example:"true"`
}

// CategoriaResponse representa a resposta de uma categoria
// @Description Dados de resposta de uma categoria
type CategoriaResponse struct {
	ID        uint   `json:"id" example:"1"`
	Nome      string `json:"nome" example:"Eletrônicos"`
	Descricao string `json:"descricao" example:"Produtos eletrônicos em geral"`
	Ativo     bool   `json:"ativo" example:"true"`
	CreatedAt string `json:"created_at" example:"2024-01-01 10:00:00"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01 10:00:00"`
}

// CategoriaWithProdutosResponse representa uma categoria com seus produtos
// @Description Dados de resposta de uma categoria com lista de produtos
type CategoriaWithProdutosResponse struct {
	ID        uint                    `json:"id" example:"1"`
	Nome      string                  `json:"nome" example:"Eletrônicos"`
	Descricao string                  `json:"descricao" example:"Produtos eletrônicos em geral"`
	Ativo     bool                    `json:"ativo" example:"true"`
	CreatedAt string                  `json:"created_at" example:"2024-01-01 10:00:00"`
	UpdatedAt string                  `json:"updated_at" example:"2024-01-01 10:00:00"`
	Produtos  []ProdutoSimpleResponse `json:"produtos"`
}

// ProdutoSimpleResponse representa uma resposta simplificada de produto (sem categoria aninhada)
// @Description Dados simplificados de um produto
type ProdutoSimpleResponse struct {
	ID        uint    `json:"id" example:"1"`
	Codigo    string  `json:"codigo" example:"PROD001"`
	Descricao string  `json:"descricao" example:"Notebook Dell"`
	Preco     float64 `json:"preco" example:"3599.90"`
}

// CategoriaPaginatedResponse representa uma resposta paginada de categorias
// @Description Resposta paginada com lista de categorias
type CategoriaPaginatedResponse struct {
	Data       []CategoriaResponse `json:"data"`
	Total      int64               `json:"total" example:"50"`
	Page       int                 `json:"page" example:"1"`
	PageSize   int                 `json:"page_size" example:"10"`
	TotalPages int                 `json:"total_pages" example:"5"`
}

