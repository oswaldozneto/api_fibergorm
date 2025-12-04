package handler

import (
	"context"
	"errors"
	"strconv"

	"api_fibergorm/pkg/arquitetura/dto"
	arqerrors "api_fibergorm/pkg/arquitetura/errors"
	"api_fibergorm/pkg/arquitetura/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// BaseService define a interface que os serviços devem implementar para o handler base
type BaseService[CreateReq any, UpdateReq any, Resp any] interface {
	Create(ctx context.Context, req *CreateReq) (*Resp, error)
	GetByID(ctx context.Context, id uint) (*Resp, error)
	GetAll(ctx context.Context, page, pageSize int) (*dto.PaginatedResponse[Resp], error)
	Update(ctx context.Context, id uint, req *UpdateReq) (*Resp, error)
	Delete(ctx context.Context, id uint) error
}

// HandlerConfig contém configurações do handler
type HandlerConfig struct {
	EntityName     string // Nome da entidade para mensagens
	SuccessMessage string // Mensagem de sucesso para delete
}

// DefaultHandlerConfig retorna configuração padrão
func DefaultHandlerConfig(entityName string) *HandlerConfig {
	return &HandlerConfig{
		EntityName:     entityName,
		SuccessMessage: entityName + " excluído(a) com sucesso",
	}
}

// BaseHandlerImpl é a implementação base do handler genérico
type BaseHandlerImpl[CreateReq any, UpdateReq any, Resp any] struct {
	Service         BaseService[CreateReq, UpdateReq, Resp]
	StructValidator *service.StructValidator
	Log             *logrus.Logger
	Config          *HandlerConfig
}

// NewBaseHandler cria uma nova instância do handler base
func NewBaseHandler[CreateReq any, UpdateReq any, Resp any](
	svc BaseService[CreateReq, UpdateReq, Resp],
	log *logrus.Logger,
	config *HandlerConfig,
) *BaseHandlerImpl[CreateReq, UpdateReq, Resp] {
	return &BaseHandlerImpl[CreateReq, UpdateReq, Resp]{
		Service:         svc,
		StructValidator: service.NewStructValidator(),
		Log:             log,
		Config:          config,
	}
}

// Create cria uma nova entidade
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) Create(c *fiber.Ctx) error {
	var req CreateReq

	if err := c.BodyParser(&req); err != nil {
		h.Log.WithError(err).Warn("Erro ao fazer parse do body")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Erro ao processar requisição",
		})
	}

	// Validação dos campos com validator (tags)
	if validationErrors := h.StructValidator.Validate(req); len(validationErrors) > 0 {
		h.Log.WithField("errors", validationErrors).Warn("Erro de validação na criação")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Erro de validação",
			Details: validationErrors,
		})
	}

	ctx := context.Background()
	result, err := h.Service.Create(ctx, &req)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetByID busca uma entidade pelo ID
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) GetByID(c *fiber.Ctx) error {
	id, err := h.ParseID(c, "id")
	if err != nil {
		return err
	}

	ctx := context.Background()
	result, err := h.Service.GetByID(ctx, id)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(result)
}

// GetAll retorna todas as entidades com paginação
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	ctx := context.Background()
	result, err := h.Service.GetAll(ctx, page, pageSize)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(result)
}

// Update atualiza uma entidade existente
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) Update(c *fiber.Ctx) error {
	id, err := h.ParseID(c, "id")
	if err != nil {
		return err
	}

	var req UpdateReq
	if err := c.BodyParser(&req); err != nil {
		h.Log.WithError(err).Warn("Erro ao fazer parse do body")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Erro ao processar requisição",
		})
	}

	// Validação dos campos com validator (tags)
	if validationErrors := h.StructValidator.Validate(req); len(validationErrors) > 0 {
		h.Log.WithField("errors", validationErrors).Warn("Erro de validação na atualização")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Erro de validação",
			Details: validationErrors,
		})
	}

	ctx := context.Background()
	result, err := h.Service.Update(ctx, id, &req)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(result)
}

// Delete remove uma entidade pelo ID
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) Delete(c *fiber.Ctx) error {
	id, err := h.ParseID(c, "id")
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := h.Service.Delete(ctx, id); err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(dto.SuccessResponse{
		Message: h.Config.SuccessMessage,
	})
}

// ParseID extrai e valida um ID da URL (exportado para uso em handlers filhos)
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) ParseID(c *fiber.Ctx, param string) (uint, error) {
	id, err := strconv.ParseUint(c.Params(param), 10, 32)
	if err != nil {
		h.Log.WithError(err).Warn("ID inválido")
		return 0, c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "ID inválido",
		})
	}
	return uint(id), nil
}

// HandleError trata os erros retornados pelo serviço (exportado para uso em handlers filhos)
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) HandleError(c *fiber.Ctx, err error) error {
	// Erros de validação
	var validationErrors *arqerrors.ValidationErrors
	if errors.As(err, &validationErrors) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Erro de validação",
			Details: validationErrors.Errors,
		})
	}

	// Erros de negócio
	if businessErr, ok := arqerrors.GetBusinessError(err); ok {
		switch businessErr.Code {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: businessErr.Message,
			})
		case "DUPLICATE":
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
				Error: businessErr.Message,
			})
		case "FORBIDDEN":
			return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
				Error: businessErr.Message,
			})
		case "HAS_RELATIONS":
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
				Error: businessErr.Message,
			})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: businessErr.Message,
			})
		}
	}

	// Erro genérico
	h.Log.WithError(err).Error("Erro interno do servidor")
	return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
		Error: "Erro interno do servidor",
	})
}

// RegisterRoutes registra as rotas CRUD padrão
func (h *BaseHandlerImpl[CreateReq, UpdateReq, Resp]) RegisterRoutes(router fiber.Router) {
	router.Post("/", h.Create)
	router.Get("/", h.GetAll)
	router.Get("/:id", h.GetByID)
	router.Put("/:id", h.Update)
	router.Delete("/:id", h.Delete)
}
