package repository

import (
	"api_fibergorm/internal/models"

	"gorm.io/gorm"
)

// CategoriaRepository define a interface para operações de repositório de categorias
type CategoriaRepository interface {
	Create(categoria *models.Categoria) error
	FindByID(id uint) (*models.Categoria, error)
	FindByIDWithProdutos(id uint) (*models.Categoria, error)
	FindByNome(nome string) (*models.Categoria, error)
	FindAll(page, pageSize int) ([]models.Categoria, int64, error)
	FindAllActive() ([]models.Categoria, error)
	Update(categoria *models.Categoria) error
	Delete(id uint) error
	ExistsByNome(nome string) (bool, error)
	ExistsByNomeExcludingID(nome string, id uint) (bool, error)
	ExistsByID(id uint) (bool, error)
	HasProdutos(id uint) (bool, error)
}

type categoriaRepository struct {
	db *gorm.DB
}

// NewCategoriaRepository cria uma nova instância do repositório de categorias
func NewCategoriaRepository(db *gorm.DB) CategoriaRepository {
	return &categoriaRepository{db: db}
}

// Create insere uma nova categoria no banco de dados
func (r *categoriaRepository) Create(categoria *models.Categoria) error {
	return r.db.Create(categoria).Error
}

// FindByID busca uma categoria pelo ID
func (r *categoriaRepository) FindByID(id uint) (*models.Categoria, error) {
	var categoria models.Categoria
	err := r.db.First(&categoria, id).Error
	if err != nil {
		return nil, err
	}
	return &categoria, nil
}

// FindByIDWithProdutos busca uma categoria pelo ID incluindo seus produtos (eager loading)
func (r *categoriaRepository) FindByIDWithProdutos(id uint) (*models.Categoria, error) {
	var categoria models.Categoria
	err := r.db.Preload("Produtos").First(&categoria, id).Error
	if err != nil {
		return nil, err
	}
	return &categoria, nil
}

// FindByNome busca uma categoria pelo nome
func (r *categoriaRepository) FindByNome(nome string) (*models.Categoria, error) {
	var categoria models.Categoria
	err := r.db.Where("nome = ?", nome).First(&categoria).Error
	if err != nil {
		return nil, err
	}
	return &categoria, nil
}

// FindAll retorna todas as categorias com paginação
func (r *categoriaRepository) FindAll(page, pageSize int) ([]models.Categoria, int64, error) {
	var categorias []models.Categoria
	var total int64

	if err := r.db.Model(&models.Categoria{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := r.db.Offset(offset).Limit(pageSize).Order("nome ASC").Find(&categorias).Error
	if err != nil {
		return nil, 0, err
	}

	return categorias, total, nil
}

// FindAllActive retorna todas as categorias ativas (para seleção em formulários)
func (r *categoriaRepository) FindAllActive() ([]models.Categoria, error) {
	var categorias []models.Categoria
	err := r.db.Where("ativo = ?", true).Order("nome ASC").Find(&categorias).Error
	return categorias, err
}

// Update atualiza uma categoria existente
func (r *categoriaRepository) Update(categoria *models.Categoria) error {
	return r.db.Save(categoria).Error
}

// Delete remove uma categoria pelo ID (soft delete)
func (r *categoriaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Categoria{}, id).Error
}

// ExistsByNome verifica se existe uma categoria com o nome informado
func (r *categoriaRepository) ExistsByNome(nome string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Categoria{}).Where("nome = ?", nome).Count(&count).Error
	return count > 0, err
}

// ExistsByNomeExcludingID verifica se existe outra categoria com o nome informado
func (r *categoriaRepository) ExistsByNomeExcludingID(nome string, id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Categoria{}).Where("nome = ? AND id != ?", nome, id).Count(&count).Error
	return count > 0, err
}

// ExistsByID verifica se existe uma categoria com o ID informado
func (r *categoriaRepository) ExistsByID(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Categoria{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// HasProdutos verifica se a categoria possui produtos associados
func (r *categoriaRepository) HasProdutos(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).Where("categoria_id = ?", id).Count(&count).Error
	return count > 0, err
}

