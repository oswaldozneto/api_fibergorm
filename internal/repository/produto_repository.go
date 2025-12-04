package repository

import (
	"api_fibergorm/internal/models"
	"api_fibergorm/pkg/arquitetura/repository"

	"gorm.io/gorm"
)

// ProdutoRepository é o repositório específico para Produto
// Herda todas as funcionalidades do BaseRepositoryImpl e pode ser estendido
// com métodos específicos conforme necessidade
type ProdutoRepository struct {
	*repository.BaseRepositoryImpl[models.Produto]
}

// NewProdutoRepository cria uma nova instância do repositório de produtos
func NewProdutoRepository(db *gorm.DB) *ProdutoRepository {
	baseRepo := repository.NewBaseRepository[models.Produto](db).
		WithPreloads("Categoria").
		WithDefaultOrder("id ASC")

	return &ProdutoRepository{
		BaseRepositoryImpl: baseRepo,
	}
}

// Métodos específicos de Produto podem ser adicionados aqui conforme necessidade
// Exemplo:
// func (r *ProdutoRepository) FindByCategoria(categoriaID uint) ([]models.Produto, error) { ... }
