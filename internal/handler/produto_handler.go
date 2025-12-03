package handler

import (
	"errors"
	"strconv"

	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/service"
	"api_fibergorm/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ProdutoHandler gerencia as requisições HTTP relacionadas a produtos
type ProdutoHandler struct {
	service   service.ProdutoService
	validator *validator.CustomValidator
	log       *logrus.Logger
}

// NewProdutoHandler cria uma nova instância do handler de produtos
func NewProdutoHandler(s service.ProdutoService, log *logrus.Logger) *ProdutoHandler {
	return &ProdutoHandler{
		service:   s,
		validator: validator.New(),
		log:       log,
	}
}

// Create godoc
// @Summary Criar um novo produto
// @Description Cria um novo produto com os dados fornecidos
// @Tags Produtos
// @Accept json
// @Produce json
// @Param produto body dto.CreateProdutoRequest true "Dados do produto"
// @Success 201 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/produtos [post]
func (h *ProdutoHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateProdutoRequest

	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Warn("Erro ao fazer parse do body")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Erro ao processar requisição",
		})
	}

	// Validação dos campos com validator
	if validationErrors := h.validator.Validate(req); len(validationErrors) > 0 {
		h.log.WithField("errors", validationErrors).Warn("Erro de validação na criação")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Erro de validação",
			Details: validationErrors,
		})
	}

	produto, err := h.service.Create(&req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(produto)
}

// GetByID godoc
// @Summary Buscar produto por ID
// @Description Retorna um produto específico pelo seu ID
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "ID do produto"
// @Success 200 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/produtos/{id} [get]
func (h *ProdutoHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		h.log.WithError(err).Warn("ID inválido")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "ID inválido",
		})
	}

	produto, err := h.service.GetByID(uint(id))
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(produto)
}

// GetAll godoc
// @Summary Listar todos os produtos
// @Description Retorna uma lista paginada de todos os produtos
// @Tags Produtos
// @Accept json
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/produtos [get]
func (h *ProdutoHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.GetAll(page, pageSize)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(response)
}

// Update godoc
// @Summary Atualizar produto
// @Description Atualiza os dados de um produto existente
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "ID do produto"
// @Param produto body dto.UpdateProdutoRequest true "Dados para atualização"
// @Success 200 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/produtos/{id} [put]
func (h *ProdutoHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		h.log.WithError(err).Warn("ID inválido")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "ID inválido",
		})
	}

	var req dto.UpdateProdutoRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Warn("Erro ao fazer parse do body")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Erro ao processar requisição",
		})
	}

	// Validação dos campos com validator
	if validationErrors := h.validator.Validate(req); len(validationErrors) > 0 {
		h.log.WithField("errors", validationErrors).Warn("Erro de validação na atualização")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Erro de validação",
			Details: validationErrors,
		})
	}

	produto, err := h.service.Update(uint(id), &req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(produto)
}

// Delete godoc
// @Summary Excluir produto
// @Description Remove um produto pelo seu ID
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "ID do produto"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/produtos/{id} [delete]
func (h *ProdutoHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		h.log.WithError(err).Warn("ID inválido")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "ID inválido",
		})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Produto excluído com sucesso",
	})
}

// handleServiceError trata os erros da camada de serviço
func (h *ProdutoHandler) handleServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrProdutoNaoEncontrado):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	case errors.Is(err, service.ErrCodigoDuplicado):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	case errors.Is(err, service.ErrPrecoInvalido),
		errors.Is(err, service.ErrDescricaoMuitoCurta),
		errors.Is(err, service.ErrCodigoVazio):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	default:
		h.log.WithError(err).Error("Erro interno do servidor")
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Erro interno do servidor",
		})
	}
}
