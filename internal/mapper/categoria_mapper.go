package mapper

import (
	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/models"
)

// CategoriaMapper implementa o mapeamento entre Categoria e seus DTOs
type CategoriaMapper struct{}

// NewCategoriaMapper cria uma nova instância do mapper
func NewCategoriaMapper() *CategoriaMapper {
	return &CategoriaMapper{}
}

// ToEntity converte CreateCategoriaRequest para Categoria
func (m *CategoriaMapper) ToEntity(req *dto.CreateCategoriaRequest) *models.Categoria {
	ativo := true
	if req.Ativo != nil {
		ativo = *req.Ativo
	}

	return &models.Categoria{
		Nome:      req.Nome,
		Descricao: req.Descricao,
		Ativo:     ativo,
	}
}

// ToResponse converte Categoria para CategoriaResponse
func (m *CategoriaMapper) ToResponse(entity *models.Categoria) *dto.CategoriaResponse {
	return &dto.CategoriaResponse{
		ID:        entity.ID,
		Nome:      entity.Nome,
		Descricao: entity.Descricao,
		Ativo:     entity.Ativo,
		CreatedAt: entity.GetCreatedAt(),
		UpdatedAt: entity.GetUpdatedAt(),
	}
}

// ApplyUpdate aplica as alterações do UpdateCategoriaRequest na entidade
func (m *CategoriaMapper) ApplyUpdate(entity *models.Categoria, req *dto.UpdateCategoriaRequest) {
	if req.Nome != "" {
		entity.Nome = req.Nome
	}
	if req.Descricao != "" {
		entity.Descricao = req.Descricao
	}
	if req.Ativo != nil {
		entity.Ativo = *req.Ativo
	}
}

// ToResponseWithProdutos converte Categoria para CategoriaWithProdutosResponse
func (m *CategoriaMapper) ToResponseWithProdutos(entity *models.Categoria) *dto.CategoriaWithProdutosResponse {
	produtos := make([]dto.ProdutoSimpleResponse, len(entity.Produtos))
	for i, p := range entity.Produtos {
		produtos[i] = dto.ProdutoSimpleResponse{
			ID:        p.ID,
			Codigo:    p.Codigo,
			Descricao: p.Descricao,
			Preco:     p.Preco,
		}
	}

	return &dto.CategoriaWithProdutosResponse{
		ID:        entity.ID,
		Nome:      entity.Nome,
		Descricao: entity.Descricao,
		Ativo:     entity.Ativo,
		CreatedAt: entity.GetCreatedAt(),
		UpdatedAt: entity.GetUpdatedAt(),
		Produtos:  produtos,
	}
}
