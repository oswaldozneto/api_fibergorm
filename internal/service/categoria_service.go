package service

import (
	"errors"
	"math"

	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/models"
	"api_fibergorm/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Erros de negócio customizados para categoria
var (
	ErrCategoriaNaoEncontrada = errors.New("categoria não encontrada")
	ErrNomeDuplicado          = errors.New("já existe uma categoria com este nome")
	ErrNomeVazio              = errors.New("o nome da categoria é obrigatório")
	ErrNomeMuitoCurto         = errors.New("o nome deve ter pelo menos 2 caracteres")
	ErrCategoriaComProdutos   = errors.New("não é possível excluir uma categoria que possui produtos")
	ErrCategoriaInativa       = errors.New("categoria inativa não pode ser utilizada")
)

// CategoriaService define a interface para os serviços de categoria
type CategoriaService interface {
	Create(req *dto.CreateCategoriaRequest) (*dto.CategoriaResponse, error)
	GetByID(id uint) (*dto.CategoriaResponse, error)
	GetByIDWithProdutos(id uint) (*dto.CategoriaWithProdutosResponse, error)
	GetAll(page, pageSize int) (*dto.CategoriaPaginatedResponse, error)
	GetAllActive() ([]dto.CategoriaResponse, error)
	Update(id uint, req *dto.UpdateCategoriaRequest) (*dto.CategoriaResponse, error)
	Delete(id uint) error
	ValidateCategoriaExists(id uint) error
}

type categoriaService struct {
	repo repository.CategoriaRepository
	log  *logrus.Logger
}

// NewCategoriaService cria uma nova instância do serviço de categorias
func NewCategoriaService(repo repository.CategoriaRepository, log *logrus.Logger) CategoriaService {
	return &categoriaService{
		repo: repo,
		log:  log,
	}
}

// Create cria uma nova categoria aplicando validações de negócio
func (s *categoriaService) Create(req *dto.CreateCategoriaRequest) (*dto.CategoriaResponse, error) {
	s.log.WithFields(logrus.Fields{
		"nome":      req.Nome,
		"descricao": req.Descricao,
	}).Info("Iniciando criação de categoria")

	// Validação de negócio: nome obrigatório
	if req.Nome == "" {
		s.log.Warn("Tentativa de criar categoria sem nome")
		return nil, ErrNomeVazio
	}

	// Validação de negócio: nome mínimo
	if len(req.Nome) < 2 {
		s.log.WithField("nome", req.Nome).Warn("Nome muito curto")
		return nil, ErrNomeMuitoCurto
	}

	// Validação de negócio: nome único
	exists, err := s.repo.ExistsByNome(req.Nome)
	if err != nil {
		s.log.WithError(err).Error("Erro ao verificar nome duplicado")
		return nil, err
	}
	if exists {
		s.log.WithField("nome", req.Nome).Warn("Tentativa de criar categoria com nome duplicado")
		return nil, ErrNomeDuplicado
	}

	// Define ativo como true por padrão se não informado
	ativo := true
	if req.Ativo != nil {
		ativo = *req.Ativo
	}

	categoria := &models.Categoria{
		Nome:      req.Nome,
		Descricao: req.Descricao,
		Ativo:     ativo,
	}

	if err := s.repo.Create(categoria); err != nil {
		s.log.WithError(err).Error("Erro ao criar categoria no banco de dados")
		return nil, err
	}

	s.log.WithField("id", categoria.ID).Info("Categoria criada com sucesso")
	return s.toResponse(categoria), nil
}

// GetByID busca uma categoria pelo ID
func (s *categoriaService) GetByID(id uint) (*dto.CategoriaResponse, error) {
	s.log.WithField("id", id).Info("Buscando categoria por ID")

	categoria, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("id", id).Warn("Categoria não encontrada")
			return nil, ErrCategoriaNaoEncontrada
		}
		s.log.WithError(err).Error("Erro ao buscar categoria")
		return nil, err
	}

	return s.toResponse(categoria), nil
}

// GetByIDWithProdutos busca uma categoria pelo ID incluindo seus produtos
func (s *categoriaService) GetByIDWithProdutos(id uint) (*dto.CategoriaWithProdutosResponse, error) {
	s.log.WithField("id", id).Info("Buscando categoria com produtos por ID")

	categoria, err := s.repo.FindByIDWithProdutos(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("id", id).Warn("Categoria não encontrada")
			return nil, ErrCategoriaNaoEncontrada
		}
		s.log.WithError(err).Error("Erro ao buscar categoria com produtos")
		return nil, err
	}

	return s.toResponseWithProdutos(categoria), nil
}

// GetAll retorna todas as categorias com paginação
func (s *categoriaService) GetAll(page, pageSize int) (*dto.CategoriaPaginatedResponse, error) {
	s.log.WithFields(logrus.Fields{
		"page":     page,
		"pageSize": pageSize,
	}).Info("Listando categorias")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	categorias, total, err := s.repo.FindAll(page, pageSize)
	if err != nil {
		s.log.WithError(err).Error("Erro ao listar categorias")
		return nil, err
	}

	categoriasResponse := make([]dto.CategoriaResponse, len(categorias))
	for i, c := range categorias {
		categoriasResponse[i] = *s.toResponse(&c)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.CategoriaPaginatedResponse{
		Data:       categoriasResponse,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllActive retorna todas as categorias ativas
func (s *categoriaService) GetAllActive() ([]dto.CategoriaResponse, error) {
	s.log.Info("Listando categorias ativas")

	categorias, err := s.repo.FindAllActive()
	if err != nil {
		s.log.WithError(err).Error("Erro ao listar categorias ativas")
		return nil, err
	}

	categoriasResponse := make([]dto.CategoriaResponse, len(categorias))
	for i, c := range categorias {
		categoriasResponse[i] = *s.toResponse(&c)
	}

	return categoriasResponse, nil
}

// Update atualiza uma categoria existente
func (s *categoriaService) Update(id uint, req *dto.UpdateCategoriaRequest) (*dto.CategoriaResponse, error) {
	s.log.WithFields(logrus.Fields{
		"id":        id,
		"nome":      req.Nome,
		"descricao": req.Descricao,
	}).Info("Iniciando atualização de categoria")

	categoria, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("id", id).Warn("Categoria não encontrada para atualização")
			return nil, ErrCategoriaNaoEncontrada
		}
		s.log.WithError(err).Error("Erro ao buscar categoria para atualização")
		return nil, err
	}

	// Validação de negócio: nome único (se alterado)
	if req.Nome != "" && req.Nome != categoria.Nome {
		if len(req.Nome) < 2 {
			s.log.WithField("nome", req.Nome).Warn("Nome muito curto")
			return nil, ErrNomeMuitoCurto
		}

		exists, err := s.repo.ExistsByNomeExcludingID(req.Nome, id)
		if err != nil {
			s.log.WithError(err).Error("Erro ao verificar nome duplicado")
			return nil, err
		}
		if exists {
			s.log.WithField("nome", req.Nome).Warn("Tentativa de atualizar para nome duplicado")
			return nil, ErrNomeDuplicado
		}
		categoria.Nome = req.Nome
	}

	if req.Descricao != "" {
		categoria.Descricao = req.Descricao
	}

	if req.Ativo != nil {
		categoria.Ativo = *req.Ativo
	}

	if err := s.repo.Update(categoria); err != nil {
		s.log.WithError(err).Error("Erro ao atualizar categoria no banco de dados")
		return nil, err
	}

	s.log.WithField("id", id).Info("Categoria atualizada com sucesso")
	return s.toResponse(categoria), nil
}

// Delete remove uma categoria pelo ID
func (s *categoriaService) Delete(id uint) error {
	s.log.WithField("id", id).Info("Iniciando exclusão de categoria")

	// Verifica se a categoria existe
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("id", id).Warn("Categoria não encontrada para exclusão")
			return ErrCategoriaNaoEncontrada
		}
		s.log.WithError(err).Error("Erro ao buscar categoria para exclusão")
		return err
	}

	// Validação de negócio: não permitir exclusão se houver produtos
	hasProdutos, err := s.repo.HasProdutos(id)
	if err != nil {
		s.log.WithError(err).Error("Erro ao verificar produtos da categoria")
		return err
	}
	if hasProdutos {
		s.log.WithField("id", id).Warn("Tentativa de excluir categoria com produtos")
		return ErrCategoriaComProdutos
	}

	if err := s.repo.Delete(id); err != nil {
		s.log.WithError(err).Error("Erro ao excluir categoria do banco de dados")
		return err
	}

	s.log.WithField("id", id).Info("Categoria excluída com sucesso")
	return nil
}

// ValidateCategoriaExists valida se uma categoria existe e está ativa
func (s *categoriaService) ValidateCategoriaExists(id uint) error {
	categoria, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoriaNaoEncontrada
		}
		return err
	}

	if !categoria.Ativo {
		return ErrCategoriaInativa
	}

	return nil
}

// toResponse converte um modelo Categoria para CategoriaResponse
func (s *categoriaService) toResponse(categoria *models.Categoria) *dto.CategoriaResponse {
	return &dto.CategoriaResponse{
		ID:        categoria.ID,
		Nome:      categoria.Nome,
		Descricao: categoria.Descricao,
		Ativo:     categoria.Ativo,
		CreatedAt: categoria.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: categoria.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// toResponseWithProdutos converte um modelo Categoria para CategoriaWithProdutosResponse
func (s *categoriaService) toResponseWithProdutos(categoria *models.Categoria) *dto.CategoriaWithProdutosResponse {
	produtos := make([]dto.ProdutoSimpleResponse, len(categoria.Produtos))
	for i, p := range categoria.Produtos {
		produtos[i] = dto.ProdutoSimpleResponse{
			ID:        p.ID,
			Codigo:    p.Codigo,
			Descricao: p.Descricao,
			Preco:     p.Preco,
		}
	}

	return &dto.CategoriaWithProdutosResponse{
		ID:        categoria.ID,
		Nome:      categoria.Nome,
		Descricao: categoria.Descricao,
		Ativo:     categoria.Ativo,
		CreatedAt: categoria.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: categoria.UpdatedAt.Format("2006-01-02 15:04:05"),
		Produtos:  produtos,
	}
}

