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

// Erros de negócio customizados
var (
	ErrProdutoNaoEncontrado = errors.New("produto não encontrado")
	ErrCodigoDuplicado      = errors.New("já existe um produto com este código")
	ErrPrecoInvalido        = errors.New("o preço deve ser maior que zero")
	ErrDescricaoMuitoCurta  = errors.New("a descrição deve ter pelo menos 3 caracteres")
	ErrCodigoVazio          = errors.New("o código do produto é obrigatório")
)

// ProdutoService define a interface para os serviços de produto
type ProdutoService interface {
	Create(req *dto.CreateProdutoRequest) (*dto.ProdutoResponse, error)
	GetByID(id uint) (*dto.ProdutoResponse, error)
	GetAll(page, pageSize int) (*dto.PaginatedResponse, error)
	Update(id uint, req *dto.UpdateProdutoRequest) (*dto.ProdutoResponse, error)
	Delete(id uint) error
}

type produtoService struct {
	repo repository.ProdutoRepository
	log  *logrus.Logger
}

// NewProdutoService cria uma nova instância do serviço de produtos
func NewProdutoService(repo repository.ProdutoRepository, log *logrus.Logger) ProdutoService {
	return &produtoService{
		repo: repo,
		log:  log,
	}
}

// Create cria um novo produto aplicando validações de negócio
func (s *produtoService) Create(req *dto.CreateProdutoRequest) (*dto.ProdutoResponse, error) {
	s.log.WithFields(logrus.Fields{
		"codigo":    req.Codigo,
		"descricao": req.Descricao,
		"preco":     req.Preco,
	}).Info("Iniciando criação de produto")

	// Validação de negócio: código único
	exists, err := s.repo.ExistsByCodigo(req.Codigo)
	if err != nil {
		s.log.WithError(err).Error("Erro ao verificar código duplicado")
		return nil, err
	}
	if exists {
		s.log.WithField("codigo", req.Codigo).Warn("Tentativa de criar produto com código duplicado")
		return nil, ErrCodigoDuplicado
	}

	// Validação de negócio: preço positivo
	if req.Preco <= 0 {
		s.log.WithField("preco", req.Preco).Warn("Tentativa de criar produto com preço inválido")
		return nil, ErrPrecoInvalido
	}

	// Validação de negócio: descrição mínima
	if len(req.Descricao) < 3 {
		s.log.WithField("descricao", req.Descricao).Warn("Descrição muito curta")
		return nil, ErrDescricaoMuitoCurta
	}

	// Validação de negócio: código obrigatório
	if req.Codigo == "" {
		s.log.Warn("Tentativa de criar produto sem código")
		return nil, ErrCodigoVazio
	}

	produto := &models.Produto{
		Codigo:    req.Codigo,
		Descricao: req.Descricao,
		Preco:     req.Preco,
	}

	if err := s.repo.Create(produto); err != nil {
		s.log.WithError(err).Error("Erro ao criar produto no banco de dados")
		return nil, err
	}

	s.log.WithField("id", produto.ID).Info("Produto criado com sucesso")
	return s.toResponse(produto), nil
}

// GetByID busca um produto pelo ID
func (s *produtoService) GetByID(id uint) (*dto.ProdutoResponse, error) {
	s.log.WithField("id", id).Info("Buscando produto por ID")

	produto, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("id", id).Warn("Produto não encontrado")
			return nil, ErrProdutoNaoEncontrado
		}
		s.log.WithError(err).Error("Erro ao buscar produto")
		return nil, err
	}

	return s.toResponse(produto), nil
}

// GetAll retorna todos os produtos com paginação
func (s *produtoService) GetAll(page, pageSize int) (*dto.PaginatedResponse, error) {
	s.log.WithFields(logrus.Fields{
		"page":     page,
		"pageSize": pageSize,
	}).Info("Listando produtos")

	// Validação de paginação
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Limite máximo
	}

	produtos, total, err := s.repo.FindAll(page, pageSize)
	if err != nil {
		s.log.WithError(err).Error("Erro ao listar produtos")
		return nil, err
	}

	// Converte para response
	produtosResponse := make([]dto.ProdutoResponse, len(produtos))
	for i, p := range produtos {
		produtosResponse[i] = *s.toResponse(&p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PaginatedResponse{
		Data:       produtosResponse,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update atualiza um produto existente
func (s *produtoService) Update(id uint, req *dto.UpdateProdutoRequest) (*dto.ProdutoResponse, error) {
	s.log.WithFields(logrus.Fields{
		"id":        id,
		"codigo":    req.Codigo,
		"descricao": req.Descricao,
		"preco":     req.Preco,
	}).Info("Iniciando atualização de produto")

	// Busca o produto existente
	produto, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("id", id).Warn("Produto não encontrado para atualização")
			return nil, ErrProdutoNaoEncontrado
		}
		s.log.WithError(err).Error("Erro ao buscar produto para atualização")
		return nil, err
	}

	// Validação de negócio: código único (se alterado)
	if req.Codigo != "" && req.Codigo != produto.Codigo {
		exists, err := s.repo.ExistsByCodigoExcludingID(req.Codigo, id)
		if err != nil {
			s.log.WithError(err).Error("Erro ao verificar código duplicado")
			return nil, err
		}
		if exists {
			s.log.WithField("codigo", req.Codigo).Warn("Tentativa de atualizar para código duplicado")
			return nil, ErrCodigoDuplicado
		}
		produto.Codigo = req.Codigo
	}

	// Validação de negócio: preço positivo (se informado)
	if req.Preco != 0 {
		if req.Preco <= 0 {
			s.log.WithField("preco", req.Preco).Warn("Tentativa de atualizar com preço inválido")
			return nil, ErrPrecoInvalido
		}
		produto.Preco = req.Preco
	}

	// Validação de negócio: descrição mínima (se informada)
	if req.Descricao != "" {
		if len(req.Descricao) < 3 {
			s.log.WithField("descricao", req.Descricao).Warn("Descrição muito curta")
			return nil, ErrDescricaoMuitoCurta
		}
		produto.Descricao = req.Descricao
	}

	if err := s.repo.Update(produto); err != nil {
		s.log.WithError(err).Error("Erro ao atualizar produto no banco de dados")
		return nil, err
	}

	s.log.WithField("id", id).Info("Produto atualizado com sucesso")
	return s.toResponse(produto), nil
}

// Delete remove um produto pelo ID
func (s *produtoService) Delete(id uint) error {
	s.log.WithField("id", id).Info("Iniciando exclusão de produto")

	// Verifica se o produto existe
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("id", id).Warn("Produto não encontrado para exclusão")
			return ErrProdutoNaoEncontrado
		}
		s.log.WithError(err).Error("Erro ao buscar produto para exclusão")
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		s.log.WithError(err).Error("Erro ao excluir produto do banco de dados")
		return err
	}

	s.log.WithField("id", id).Info("Produto excluído com sucesso")
	return nil
}

// toResponse converte um modelo Produto para ProdutoResponse
func (s *produtoService) toResponse(produto *models.Produto) *dto.ProdutoResponse {
	return &dto.ProdutoResponse{
		ID:        produto.ID,
		Codigo:    produto.Codigo,
		Descricao: produto.Descricao,
		Preco:     produto.Preco,
		CreatedAt: produto.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: produto.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
