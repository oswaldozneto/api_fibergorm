package mapper

import (
	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/models"
)

// ProdutoMapper implementa o mapeamento entre Produto e seus DTOs
type ProdutoMapper struct{}

// NewProdutoMapper cria uma nova instância do mapper
func NewProdutoMapper() *ProdutoMapper {
	return &ProdutoMapper{}
}

// ToEntity converte CreateProdutoRequest para Produto
func (m *ProdutoMapper) ToEntity(req *dto.CreateProdutoRequest) *models.Produto {
	return &models.Produto{
		Codigo:      req.Codigo,
		Descricao:   req.Descricao,
		Preco:       req.Preco,
		CategoriaID: req.CategoriaID,
	}
}

// ToResponse converte Produto para ProdutoResponse
func (m *ProdutoMapper) ToResponse(entity *models.Produto) *dto.ProdutoResponse {
	response := &dto.ProdutoResponse{
		ID:          entity.ID,
		Codigo:      entity.Codigo,
		Descricao:   entity.Descricao,
		Preco:       entity.Preco,
		CategoriaID: entity.CategoriaID,
		CreatedAt:   entity.GetCreatedAt(),
		UpdatedAt:   entity.GetUpdatedAt(),
	}

	// Se a categoria foi carregada (eager loading), inclui os dados
	if entity.Categoria.ID != 0 {
		response.Categoria = &dto.CategoriaResponse{
			ID:        entity.Categoria.ID,
			Nome:      entity.Categoria.Nome,
			Descricao: entity.Categoria.Descricao,
			Ativo:     entity.Categoria.Ativo,
			CreatedAt: entity.Categoria.GetCreatedAt(),
			UpdatedAt: entity.Categoria.GetUpdatedAt(),
		}
	}

	return response
}

// ApplyUpdate aplica as alterações do UpdateProdutoRequest na entidade
func (m *ProdutoMapper) ApplyUpdate(entity *models.Produto, req *dto.UpdateProdutoRequest) {
	if req.Codigo != "" {
		entity.Codigo = req.Codigo
	}
	if req.Descricao != "" {
		entity.Descricao = req.Descricao
	}
	if req.Preco != 0 {
		entity.Preco = req.Preco
	}
	if req.CategoriaID != 0 {
		entity.CategoriaID = req.CategoriaID
	}
}
