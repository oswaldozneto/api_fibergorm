package service

import (
	"context"

	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/mapper"
	"api_fibergorm/internal/models"
	"api_fibergorm/internal/repository"
	"api_fibergorm/internal/validator"
	arqdto "api_fibergorm/pkg/arquitetura/dto"
	arqerrors "api_fibergorm/pkg/arquitetura/errors"
	"api_fibergorm/pkg/arquitetura/service"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ProdutoService define a interface para os serviços de produto
type ProdutoService interface {
	Create(ctx context.Context, req *dto.CreateProdutoRequest) (*dto.ProdutoResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.ProdutoResponse, error)
	GetAll(ctx context.Context, page, pageSize int) (*arqdto.PaginatedResponse[dto.ProdutoResponse], error)
	GetByCategoriaID(ctx context.Context, categoriaID uint, page, pageSize int) (*arqdto.PaginatedResponse[dto.ProdutoResponse], error)
	Update(ctx context.Context, id uint, req *dto.UpdateProdutoRequest) (*dto.ProdutoResponse, error)
	Delete(ctx context.Context, id uint) error
}

// produtoService é a implementação do serviço usando a arquitetura base
type produtoService struct {
	*service.BaseServiceImpl[models.Produto, dto.CreateProdutoRequest, dto.UpdateProdutoRequest, dto.ProdutoResponse]
	repo   *repository.ProdutoRepository
	mapper *mapper.ProdutoMapper
	db     *gorm.DB
	log    *logrus.Logger
}

// NewProdutoService cria uma nova instância do serviço de produtos
func NewProdutoService(db *gorm.DB, log *logrus.Logger) ProdutoService {
	// Cria o repositório específico de produto
	repo := repository.NewProdutoRepository(db)

	// Cria o mapper
	produtoMapper := mapper.NewProdutoMapper()

	// Configuração do serviço
	config := service.DefaultServiceConfig("Produto")

	// Cria o serviço base usando o repositório base embutido
	baseService := service.NewBaseService[
		models.Produto,
		dto.CreateProdutoRequest,
		dto.UpdateProdutoRequest,
		dto.ProdutoResponse,
	](repo.BaseRepositoryImpl, produtoMapper, log, config)

	// Cria o validador específico
	produtoValidator := validator.NewProdutoValidator(repo, db, log)

	// Configura o validador no serviço
	baseService.WithValidator(produtoValidator)

	return &produtoService{
		BaseServiceImpl: baseService,
		repo:            repo,
		mapper:          produtoMapper,
		db:              db,
		log:             log,
	}
}

// Create sobrescreve o Create base para recarregar com categoria
func (s *produtoService) Create(ctx context.Context, req *dto.CreateProdutoRequest) (*dto.ProdutoResponse, error) {
	// Chama o Create base
	response, err := s.BaseServiceImpl.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// Recarrega com categoria para garantir dados completos
	produto, err := s.repo.FindByID(response.ID)
	if err != nil {
		return response, nil // Retorna o response original se falhar
	}

	return s.mapper.ToResponse(produto), nil
}

// Update sobrescreve o Update base para recarregar com categoria
func (s *produtoService) Update(ctx context.Context, id uint, req *dto.UpdateProdutoRequest) (*dto.ProdutoResponse, error) {
	// Chama o Update base
	response, err := s.BaseServiceImpl.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Recarrega com categoria para garantir dados completos
	produto, err := s.repo.FindByID(response.ID)
	if err != nil {
		return response, nil // Retorna o response original se falhar
	}

	return s.mapper.ToResponse(produto), nil
}

// GetByCategoriaID retorna produtos de uma categoria específica
func (s *produtoService) GetByCategoriaID(ctx context.Context, categoriaID uint, page, pageSize int) (*arqdto.PaginatedResponse[dto.ProdutoResponse], error) {
	s.log.WithFields(logrus.Fields{
		"categoria_id": categoriaID,
		"page":         page,
		"pageSize":     pageSize,
	}).Info("Listando produtos por categoria")

	// Normaliza paginação
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Verifica se a categoria existe
	var count int64
	if err := s.db.Model(&models.Categoria{}).Where("id = ?", categoriaID).Count(&count).Error; err != nil {
		s.log.WithError(err).Error("Erro ao verificar categoria")
		return nil, err
	}
	if count == 0 {
		return nil, arqerrors.NewBusinessError("NOT_FOUND", "Categoria não encontrada")
	}

	// Busca produtos da categoria
	produtos, total, err := s.repo.FindAllWhere(page, pageSize, "id ASC", "categoria_id = ?", categoriaID)
	if err != nil {
		s.log.WithError(err).Error("Erro ao listar produtos por categoria")
		return nil, err
	}

	// Converte para responses
	responses := make([]dto.ProdutoResponse, len(produtos))
	for i := range produtos {
		responses[i] = *s.mapper.ToResponse(&produtos[i])
	}

	return arqdto.NewPaginatedResponse(responses, total, page, pageSize), nil
}
