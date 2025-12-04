package validator

import (
	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/models"
	"api_fibergorm/internal/repository"
	"api_fibergorm/pkg/arquitetura/service"

	"github.com/sirupsen/logrus"
)

// CategoriaValidator implementa validações específicas para Categoria
type CategoriaValidator struct {
	repo *repository.CategoriaRepository
	log  *logrus.Logger
}

// NewCategoriaValidator cria um novo validador de categoria
func NewCategoriaValidator(
	repo *repository.CategoriaRepository,
	log *logrus.Logger,
) *CategoriaValidator {
	return &CategoriaValidator{
		repo: repo,
		log:  log,
	}
}

// ValidateCreate valida a criação de uma categoria
func (v *CategoriaValidator) ValidateCreate(ctx *service.ValidationContext, req *dto.CreateCategoriaRequest) *service.ValidationResult {
	result := service.NewValidationResult()

	// Validação: nome obrigatório
	if req.Nome == "" {
		v.log.Warn("Tentativa de criar categoria sem nome")
		result.AddError("nome", "O nome da categoria é obrigatório")
		return result
	}

	// Validação: nome mínimo
	if len(req.Nome) < 2 {
		v.log.WithField("nome", req.Nome).Warn("Nome muito curto")
		result.AddError("nome", "O nome deve ter pelo menos 2 caracteres")
		return result
	}

	// Validação: nome único
	exists, err := v.repo.ExistsWhere("nome = ?", req.Nome)
	if err != nil {
		v.log.WithError(err).Error("Erro ao verificar nome duplicado")
		result.AddError("nome", "Erro ao verificar nome")
		return result
	}
	if exists {
		v.log.WithField("nome", req.Nome).Warn("Tentativa de criar categoria com nome duplicado")
		result.AddError("nome", "Já existe uma categoria com este nome")
	}

	return result
}

// ValidateUpdate valida a atualização de uma categoria
func (v *CategoriaValidator) ValidateUpdate(ctx *service.ValidationContext, entity *models.Categoria, req *dto.UpdateCategoriaRequest) *service.ValidationResult {
	result := service.NewValidationResult()

	// Validação: nome único (se alterado)
	if req.Nome != "" && req.Nome != entity.Nome {
		if len(req.Nome) < 2 {
			v.log.WithField("nome", req.Nome).Warn("Nome muito curto")
			result.AddError("nome", "O nome deve ter pelo menos 2 caracteres")
			return result
		}

		exists, err := v.repo.ExistsWhereExcludingID(ctx.EntityID, "nome = ?", req.Nome)
		if err != nil {
			v.log.WithError(err).Error("Erro ao verificar nome duplicado")
			result.AddError("nome", "Erro ao verificar nome")
			return result
		}
		if exists {
			v.log.WithField("nome", req.Nome).Warn("Tentativa de atualizar para nome duplicado")
			result.AddError("nome", "Já existe outra categoria com este nome")
		}
	}

	return result
}

// ValidateDelete valida a exclusão de uma categoria
func (v *CategoriaValidator) ValidateDelete(ctx *service.ValidationContext, entity *models.Categoria) *service.ValidationResult {
	result := service.NewValidationResult()

	// Validação: não permitir exclusão se houver produtos
	count, err := v.countProdutos(entity.ID)
	if err != nil {
		v.log.WithError(err).Error("Erro ao verificar produtos da categoria")
		result.AddError("categoria", "Erro ao verificar produtos relacionados")
		return result
	}
	if count > 0 {
		v.log.WithField("id", entity.ID).Warn("Tentativa de excluir categoria com produtos")
		result.AddError("categoria", "Não é possível excluir uma categoria que possui produtos")
	}

	return result
}

// countProdutos conta os produtos de uma categoria
func (v *CategoriaValidator) countProdutos(categoriaID uint) (int64, error) {
	var count int64
	err := v.repo.GetDB().Model(&models.Produto{}).Where("categoria_id = ?", categoriaID).Count(&count).Error
	return count, err
}
