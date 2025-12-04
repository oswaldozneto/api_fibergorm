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

// CategoriaService define a interface para os serviços de categoria
type CategoriaService interface {
	Create(ctx context.Context, req *dto.CreateCategoriaRequest) (*dto.CategoriaResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.CategoriaResponse, error)
	GetByIDWithProdutos(ctx context.Context, id uint) (*dto.CategoriaWithProdutosResponse, error)
	GetAll(ctx context.Context, page, pageSize int) (*arqdto.PaginatedResponse[dto.CategoriaResponse], error)
	GetAllActive(ctx context.Context) ([]dto.CategoriaResponse, error)
	Update(ctx context.Context, id uint, req *dto.UpdateCategoriaRequest) (*dto.CategoriaResponse, error)
	Delete(ctx context.Context, id uint) error
}

// categoriaService é a implementação do serviço usando a arquitetura base
type categoriaService struct {
	*service.BaseServiceImpl[models.Categoria, dto.CreateCategoriaRequest, dto.UpdateCategoriaRequest, dto.CategoriaResponse]
	repo   *repository.CategoriaRepository
	mapper *mapper.CategoriaMapper
	log    *logrus.Logger
}

// NewCategoriaService cria uma nova instância do serviço de categorias
func NewCategoriaService(db *gorm.DB, log *logrus.Logger) CategoriaService {
	// Cria o repositório específico de categoria
	repo := repository.NewCategoriaRepository(db)

	// Cria o mapper
	categoriaMapper := mapper.NewCategoriaMapper()

	// Configuração do serviço
	config := service.DefaultServiceConfig("Categoria")
	config.DefaultOrder = "nome ASC"

	// Cria o serviço base usando o repositório base embutido
	baseService := service.NewBaseService[
		models.Categoria,
		dto.CreateCategoriaRequest,
		dto.UpdateCategoriaRequest,
		dto.CategoriaResponse,
	](repo.BaseRepositoryImpl, categoriaMapper, log, config)

	// Cria o validador específico
	categoriaValidator := validator.NewCategoriaValidator(repo, log)

	// Configura o validador no serviço
	baseService.WithValidator(categoriaValidator)

	return &categoriaService{
		BaseServiceImpl: baseService,
		repo:            repo,
		mapper:          categoriaMapper,
		log:             log,
	}
}

// GetByIDWithProdutos busca uma categoria pelo ID incluindo seus produtos
func (s *categoriaService) GetByIDWithProdutos(ctx context.Context, id uint) (*dto.CategoriaWithProdutosResponse, error) {
	s.log.WithField("id", id).Info("Buscando categoria com produtos por ID")

	categoria, err := s.repo.FindByIDWithPreloads(id, "Produtos")
	if err != nil {
		if arqerrors.IsNotFound(err) {
			s.log.WithField("id", id).Warn("Categoria não encontrada")
			return nil, arqerrors.NewBusinessError("NOT_FOUND", "Categoria não encontrada")
		}
		s.log.WithError(err).Error("Erro ao buscar categoria com produtos")
		return nil, err
	}

	return s.mapper.ToResponseWithProdutos(categoria), nil
}

// GetAllActive retorna todas as categorias ativas
func (s *categoriaService) GetAllActive(ctx context.Context) ([]dto.CategoriaResponse, error) {
	s.log.Info("Listando categorias ativas")

	categorias, _, err := s.repo.FindAllWhere(1, 1000, "nome ASC", "ativo = ?", true)
	if err != nil {
		s.log.WithError(err).Error("Erro ao listar categorias ativas")
		return nil, err
	}

	responses := make([]dto.CategoriaResponse, len(categorias))
	for i := range categorias {
		responses[i] = *s.mapper.ToResponse(&categorias[i])
	}

	return responses, nil
}
