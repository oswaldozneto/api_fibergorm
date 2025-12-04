package repository

import (
	"api_fibergorm/internal/models"
	"api_fibergorm/pkg/arquitetura/repository"

	"gorm.io/gorm"
)

// CategoriaRepository é o repositório específico para Categoria
// Herda todas as funcionalidades do BaseRepositoryImpl e pode ser estendido
// com métodos específicos conforme necessidade
type CategoriaRepository struct {
	*repository.BaseRepositoryImpl[*models.Categoria]
}

// NewCategoriaRepository cria uma nova instância do repositório de categorias
func NewCategoriaRepository(db *gorm.DB) *CategoriaRepository {
	baseRepo := repository.NewBaseRepository[*models.Categoria](db).
		WithDefaultOrder("nome ASC")

	return &CategoriaRepository{
		BaseRepositoryImpl: baseRepo,
	}
}

// Métodos específicos de Categoria podem ser adicionados aqui conforme necessidade
// Exemplo:
// func (r *CategoriaRepository) FindAllActive() ([]*models.Categoria, error) { ... }
