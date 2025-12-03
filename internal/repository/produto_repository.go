package repository

import (
	"api_fibergorm/internal/models"

	"gorm.io/gorm"
)

// ProdutoRepository define a interface para operações de repositório de produtos
type ProdutoRepository interface {
	Create(produto *models.Produto) error
	FindByID(id uint) (*models.Produto, error)
	FindByIDWithCategoria(id uint) (*models.Produto, error)
	FindByCodigo(codigo string) (*models.Produto, error)
	FindAll(page, pageSize int) ([]models.Produto, int64, error)
	FindAllWithCategoria(page, pageSize int) ([]models.Produto, int64, error)
	FindByCategoriaID(categoriaID uint, page, pageSize int) ([]models.Produto, int64, error)
	Update(produto *models.Produto) error
	Delete(id uint) error
	ExistsByCodigo(codigo string) (bool, error)
	ExistsByCodigoExcludingID(codigo string, id uint) (bool, error)
}

type produtoRepository struct {
	db *gorm.DB
}

// NewProdutoRepository cria uma nova instância do repositório de produtos
func NewProdutoRepository(db *gorm.DB) ProdutoRepository {
	return &produtoRepository{db: db}
}

// Create insere um novo produto no banco de dados
func (r *produtoRepository) Create(produto *models.Produto) error {
	return r.db.Create(produto).Error
}

// FindByID busca um produto pelo ID
func (r *produtoRepository) FindByID(id uint) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.First(&produto, id).Error
	if err != nil {
		return nil, err
	}
	return &produto, nil
}

// FindByIDWithCategoria busca um produto pelo ID incluindo a categoria (eager loading)
func (r *produtoRepository) FindByIDWithCategoria(id uint) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.Preload("Categoria").First(&produto, id).Error
	if err != nil {
		return nil, err
	}
	return &produto, nil
}

// FindByCodigo busca um produto pelo código
func (r *produtoRepository) FindByCodigo(codigo string) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.Where("codigo = ?", codigo).First(&produto).Error
	if err != nil {
		return nil, err
	}
	return &produto, nil
}

// FindAll retorna todos os produtos com paginação
func (r *produtoRepository) FindAll(page, pageSize int) ([]models.Produto, int64, error) {
	var produtos []models.Produto
	var total int64

	// Conta o total de registros
	if err := r.db.Model(&models.Produto{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calcula o offset
	offset := (page - 1) * pageSize

	// Busca os produtos com paginação
	err := r.db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&produtos).Error
	if err != nil {
		return nil, 0, err
	}

	return produtos, total, nil
}

// FindAllWithCategoria retorna todos os produtos com categoria (eager loading)
func (r *produtoRepository) FindAllWithCategoria(page, pageSize int) ([]models.Produto, int64, error) {
	var produtos []models.Produto
	var total int64

	if err := r.db.Model(&models.Produto{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := r.db.Preload("Categoria").Offset(offset).Limit(pageSize).Order("id ASC").Find(&produtos).Error
	if err != nil {
		return nil, 0, err
	}

	return produtos, total, nil
}

// FindByCategoriaID retorna produtos de uma categoria específica
func (r *produtoRepository) FindByCategoriaID(categoriaID uint, page, pageSize int) ([]models.Produto, int64, error) {
	var produtos []models.Produto
	var total int64

	if err := r.db.Model(&models.Produto{}).Where("categoria_id = ?", categoriaID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := r.db.Preload("Categoria").Where("categoria_id = ?", categoriaID).
		Offset(offset).Limit(pageSize).Order("id ASC").Find(&produtos).Error
	if err != nil {
		return nil, 0, err
	}

	return produtos, total, nil
}

// Update atualiza um produto existente
func (r *produtoRepository) Update(produto *models.Produto) error {
	return r.db.Save(produto).Error
}

// Delete remove um produto pelo ID (soft delete)
func (r *produtoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Produto{}, id).Error
}

// ExistsByCodigo verifica se existe um produto com o código informado
func (r *produtoRepository) ExistsByCodigo(codigo string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).Where("codigo = ?", codigo).Count(&count).Error
	return count > 0, err
}

// ExistsByCodigoExcludingID verifica se existe outro produto com o código informado
func (r *produtoRepository) ExistsByCodigoExcludingID(codigo string, id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).Where("codigo = ? AND id != ?", codigo, id).Count(&count).Error
	return count > 0, err
}
