package service

import (
	"context"

	"api_fibergorm/pkg/arquitetura/dto"
	arqerrors "api_fibergorm/pkg/arquitetura/errors"
	"api_fibergorm/pkg/arquitetura/repository"

	"github.com/sirupsen/logrus"
)

// BaseService define a interface base para serviços
type BaseService[E any, CreateReq any, UpdateReq any, Resp any] interface {
	Create(ctx context.Context, req *CreateReq) (*Resp, error)
	GetByID(ctx context.Context, id uint) (*Resp, error)
	GetAll(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse[Resp], error)
	Update(ctx context.Context, id uint, req *UpdateReq) (*Resp, error)
	Delete(ctx context.Context, id uint) error
}

// ServiceConfig contém as configurações do serviço
type ServiceConfig struct {
	EntityName   string // Nome da entidade para logs e mensagens
	DefaultOrder string // Ordenação padrão
	MaxPageSize  int    // Tamanho máximo da página
}

// DefaultServiceConfig retorna configuração padrão
func DefaultServiceConfig(entityName string) *ServiceConfig {
	return &ServiceConfig{
		EntityName:   entityName,
		DefaultOrder: "id ASC",
		MaxPageSize:  100,
	}
}

// BaseServiceImpl é a implementação base do serviço genérico
type BaseServiceImpl[E any, CreateReq any, UpdateReq any, Resp any] struct {
	repo            *repository.BaseRepositoryImpl[E]
	mapper          dto.Mapper[E, CreateReq, UpdateReq, Resp]
	validator       EntityValidator[E, CreateReq, UpdateReq]
	structValidator *StructValidator
	log             *logrus.Logger
	config          *ServiceConfig
}

// NewBaseService cria uma nova instância do serviço base
func NewBaseService[E any, CreateReq any, UpdateReq any, Resp any](
	repo *repository.BaseRepositoryImpl[E],
	mapper dto.Mapper[E, CreateReq, UpdateReq, Resp],
	log *logrus.Logger,
	config *ServiceConfig,
) *BaseServiceImpl[E, CreateReq, UpdateReq, Resp] {
	return &BaseServiceImpl[E, CreateReq, UpdateReq, Resp]{
		repo:            repo,
		mapper:          mapper,
		validator:       &NoOpValidator[E, CreateReq, UpdateReq]{},
		structValidator: NewStructValidator(),
		log:             log,
		config:          config,
	}
}

// WithValidator configura um validador customizado
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) WithValidator(v EntityValidator[E, CreateReq, UpdateReq]) *BaseServiceImpl[E, CreateReq, UpdateReq, Resp] {
	s.validator = v
	return s
}

// GetRepository retorna o repositório para uso em validações
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) GetRepository() *repository.BaseRepositoryImpl[E] {
	return s.repo
}

// GetLogger retorna o logger para uso em validações
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) GetLogger() *logrus.Logger {
	return s.log
}

// Create cria uma nova entidade
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) Create(ctx context.Context, req *CreateReq) (*Resp, error) {
	s.log.WithField("entity", s.config.EntityName).Info("Iniciando criação")

	// Validação de struct (tags de validação)
	if structErrors := s.structValidator.ToValidationResult(req); structErrors != nil && structErrors.HasErrors() {
		s.log.WithField("errors", structErrors.Errors).Warn("Erro de validação de struct na criação")
		return nil, &arqerrors.ValidationErrors{Errors: structErrors.Errors}
	}

	// Validação customizada da entidade
	validationCtx := &ValidationContext{
		Context:   ctx,
		Operation: OperationCreate,
	}

	if customErrors := s.validator.ValidateCreate(validationCtx, req); customErrors != nil && customErrors.HasErrors() {
		s.log.WithField("errors", customErrors.Errors).Warn("Erro de validação customizada na criação")
		return nil, &arqerrors.ValidationErrors{Errors: customErrors.Errors}
	}

	// Converte request para entidade
	entity := s.mapper.ToEntity(req)

	// Persiste no banco
	if err := s.repo.Create(entity); err != nil {
		s.log.WithError(err).Error("Erro ao criar no banco de dados")
		return nil, err
	}

	s.log.WithField("entity", s.config.EntityName).Info("Criado com sucesso")

	// Converte para response
	response := s.mapper.ToResponse(entity)
	return response, nil
}

// GetByID busca uma entidade pelo ID
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) GetByID(ctx context.Context, id uint) (*Resp, error) {
	s.log.WithFields(logrus.Fields{
		"entity": s.config.EntityName,
		"id":     id,
	}).Info("Buscando por ID")

	entity, err := s.repo.FindByID(id)
	if err != nil {
		if arqerrors.IsNotFound(err) {
			s.log.WithField("id", id).Warn("Não encontrado")
			return nil, arqerrors.NewBusinessError("NOT_FOUND", s.config.EntityName+" não encontrado(a)")
		}
		s.log.WithError(err).Error("Erro ao buscar")
		return nil, err
	}

	response := s.mapper.ToResponse(entity)
	return response, nil
}

// GetAll retorna todas as entidades com paginação
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) GetAll(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse[Resp], error) {
	s.log.WithFields(logrus.Fields{
		"entity":   s.config.EntityName,
		"page":     page,
		"pageSize": pageSize,
	}).Info("Listando")

	// Normaliza paginação
	page, pageSize = s.normalizePagination(page, pageSize)

	entities, total, err := s.repo.FindAll(page, pageSize, s.config.DefaultOrder)
	if err != nil {
		s.log.WithError(err).Error("Erro ao listar")
		return nil, err
	}

	// Converte para responses
	responses := make([]Resp, len(entities))
	for i := range entities {
		responses[i] = *s.mapper.ToResponse(&entities[i])
	}

	return dto.NewPaginatedResponse(responses, total, page, pageSize), nil
}

// Update atualiza uma entidade existente
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) Update(ctx context.Context, id uint, req *UpdateReq) (*Resp, error) {
	s.log.WithFields(logrus.Fields{
		"entity": s.config.EntityName,
		"id":     id,
	}).Info("Iniciando atualização")

	// Busca a entidade existente
	entity, err := s.repo.FindByID(id)
	if err != nil {
		if arqerrors.IsNotFound(err) {
			s.log.WithField("id", id).Warn("Não encontrado para atualização")
			return nil, arqerrors.NewBusinessError("NOT_FOUND", s.config.EntityName+" não encontrado(a)")
		}
		s.log.WithError(err).Error("Erro ao buscar para atualização")
		return nil, err
	}

	// Validação de struct (tags de validação)
	if structErrors := s.structValidator.ToValidationResult(req); structErrors != nil && structErrors.HasErrors() {
		s.log.WithField("errors", structErrors.Errors).Warn("Erro de validação de struct na atualização")
		return nil, &arqerrors.ValidationErrors{Errors: structErrors.Errors}
	}

	// Validação customizada da entidade
	validationCtx := &ValidationContext{
		Context:   ctx,
		Operation: OperationUpdate,
		EntityID:  id,
	}

	if customErrors := s.validator.ValidateUpdate(validationCtx, entity, req); customErrors != nil && customErrors.HasErrors() {
		s.log.WithField("errors", customErrors.Errors).Warn("Erro de validação customizada na atualização")
		return nil, &arqerrors.ValidationErrors{Errors: customErrors.Errors}
	}

	// Aplica as alterações
	s.mapper.ApplyUpdate(entity, req)

	// Persiste no banco
	if err := s.repo.Update(entity); err != nil {
		s.log.WithError(err).Error("Erro ao atualizar no banco de dados")
		return nil, err
	}

	s.log.WithField("id", id).Info("Atualizado com sucesso")

	// Converte para response
	response := s.mapper.ToResponse(entity)
	return response, nil
}

// Delete remove uma entidade pelo ID
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) Delete(ctx context.Context, id uint) error {
	s.log.WithFields(logrus.Fields{
		"entity": s.config.EntityName,
		"id":     id,
	}).Info("Iniciando exclusão")

	// Busca a entidade existente
	entity, err := s.repo.FindByID(id)
	if err != nil {
		if arqerrors.IsNotFound(err) {
			s.log.WithField("id", id).Warn("Não encontrado para exclusão")
			return arqerrors.NewBusinessError("NOT_FOUND", s.config.EntityName+" não encontrado(a)")
		}
		s.log.WithError(err).Error("Erro ao buscar para exclusão")
		return err
	}

	// Validação customizada da entidade
	validationCtx := &ValidationContext{
		Context:   ctx,
		Operation: OperationDelete,
		EntityID:  id,
	}

	if customErrors := s.validator.ValidateDelete(validationCtx, entity); customErrors != nil && customErrors.HasErrors() {
		s.log.WithField("errors", customErrors.Errors).Warn("Erro de validação customizada na exclusão")
		return &arqerrors.ValidationErrors{Errors: customErrors.Errors}
	}

	// Remove do banco
	if err := s.repo.Delete(id); err != nil {
		s.log.WithError(err).Error("Erro ao excluir do banco de dados")
		return err
	}

	s.log.WithField("id", id).Info("Excluído com sucesso")
	return nil
}

// normalizePagination normaliza os valores de paginação
func (s *BaseServiceImpl[E, CreateReq, UpdateReq, Resp]) normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > s.config.MaxPageSize {
		pageSize = s.config.MaxPageSize
	}
	return page, pageSize
}
