package validator

import (
	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/models"
	"api_fibergorm/internal/repository"
	"api_fibergorm/pkg/arquitetura/service"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ProdutoValidator implementa validações específicas para Produto
type ProdutoValidator struct {
	repo *repository.ProdutoRepository
	db   *gorm.DB
	log  *logrus.Logger
}

// NewProdutoValidator cria um novo validador de produto
func NewProdutoValidator(
	repo *repository.ProdutoRepository,
	db *gorm.DB,
	log *logrus.Logger,
) *ProdutoValidator {
	return &ProdutoValidator{
		repo: repo,
		db:   db,
		log:  log,
	}
}

// ValidateCreate valida a criação de um produto
func (v *ProdutoValidator) ValidateCreate(ctx *service.ValidationContext, req *dto.CreateProdutoRequest) *service.ValidationResult {
	result := service.NewValidationResult()

	// Validação: código obrigatório
	if req.Codigo == "" {
		v.log.Warn("Tentativa de criar produto sem código")
		result.AddError("codigo", "O código do produto é obrigatório")
		return result
	}

	// Validação: código único
	exists, err := v.repo.ExistsWhere("codigo = ?", req.Codigo)
	if err != nil {
		v.log.WithError(err).Error("Erro ao verificar código duplicado")
		result.AddError("codigo", "Erro ao verificar código")
		return result
	}
	if exists {
		v.log.WithField("codigo", req.Codigo).Warn("Tentativa de criar produto com código duplicado")
		result.AddError("codigo", "Já existe um produto com este código")
		return result
	}

	// Validação: preço positivo
	if req.Preco <= 0 {
		v.log.WithField("preco", req.Preco).Warn("Tentativa de criar produto com preço inválido")
		result.AddError("preco", "O preço deve ser maior que zero")
		return result
	}

	// Validação: descrição mínima
	if len(req.Descricao) < 3 {
		v.log.WithField("descricao", req.Descricao).Warn("Descrição muito curta")
		result.AddError("descricao", "A descrição deve ter pelo menos 3 caracteres")
		return result
	}

	// Validação: categoria obrigatória
	if req.CategoriaID == 0 {
		v.log.Warn("Tentativa de criar produto sem categoria")
		result.AddError("categoria_id", "A categoria é obrigatória")
		return result
	}

	// Validação: categoria deve existir e estar ativa
	categoria, err := v.findCategoria(req.CategoriaID)
	if err != nil {
		v.log.WithField("categoria_id", req.CategoriaID).Warn("Categoria não encontrada")
		result.AddError("categoria_id", "Categoria não encontrada")
		return result
	}
	if !categoria.Ativo {
		v.log.WithField("categoria_id", req.CategoriaID).Warn("Categoria inativa")
		result.AddError("categoria_id", "Categoria inativa não pode ser utilizada")
	}

	return result
}

// ValidateUpdate valida a atualização de um produto
func (v *ProdutoValidator) ValidateUpdate(ctx *service.ValidationContext, entity *models.Produto, req *dto.UpdateProdutoRequest) *service.ValidationResult {
	result := service.NewValidationResult()

	// Validação: código único (se alterado)
	if req.Codigo != "" && req.Codigo != entity.Codigo {
		exists, err := v.repo.ExistsWhereExcludingID(ctx.EntityID, "codigo = ?", req.Codigo)
		if err != nil {
			v.log.WithError(err).Error("Erro ao verificar código duplicado")
			result.AddError("codigo", "Erro ao verificar código")
			return result
		}
		if exists {
			v.log.WithField("codigo", req.Codigo).Warn("Tentativa de atualizar para código duplicado")
			result.AddError("codigo", "Já existe outro produto com este código")
			return result
		}
	}

	// Validação: preço positivo (se informado)
	if req.Preco != 0 && req.Preco <= 0 {
		v.log.WithField("preco", req.Preco).Warn("Tentativa de atualizar com preço inválido")
		result.AddError("preco", "O preço deve ser maior que zero")
		return result
	}

	// Validação: descrição mínima (se informada)
	if req.Descricao != "" && len(req.Descricao) < 3 {
		v.log.WithField("descricao", req.Descricao).Warn("Descrição muito curta")
		result.AddError("descricao", "A descrição deve ter pelo menos 3 caracteres")
		return result
	}

	// Validação: categoria (se informada)
	if req.CategoriaID != 0 && req.CategoriaID != entity.CategoriaID {
		categoria, err := v.findCategoria(req.CategoriaID)
		if err != nil {
			v.log.WithField("categoria_id", req.CategoriaID).Warn("Categoria não encontrada")
			result.AddError("categoria_id", "Categoria não encontrada")
			return result
		}
		if !categoria.Ativo {
			v.log.WithField("categoria_id", req.CategoriaID).Warn("Categoria inativa")
			result.AddError("categoria_id", "Categoria inativa não pode ser utilizada")
		}
	}

	return result
}

// ValidateDelete valida a exclusão de um produto
func (v *ProdutoValidator) ValidateDelete(ctx *service.ValidationContext, entity *models.Produto) *service.ValidationResult {
	// Produto não tem validações especiais para exclusão
	return nil
}

// findCategoria busca uma categoria pelo ID
func (v *ProdutoValidator) findCategoria(id uint) (*models.Categoria, error) {
	var categoria models.Categoria
	err := v.db.First(&categoria, id).Error
	if err != nil {
		return nil, err
	}
	return &categoria, nil
}
