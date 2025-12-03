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

// CategoriaHandler gerencia as requisições HTTP relacionadas a categorias
type CategoriaHandler struct {
	service   service.CategoriaService
	validator *validator.CustomValidator
	log       *logrus.Logger
}

// NewCategoriaHandler cria uma nova instância do handler de categorias
func NewCategoriaHandler(s service.CategoriaService, log *logrus.Logger) *CategoriaHandler {
	return &CategoriaHandler{
		service:   s,
		validator: validator.New(),
		log:       log,
	}
}

// Create godoc
// @Summary Criar uma nova categoria
// @Description Cria uma nova categoria com os dados fornecidos
// @Tags Categorias
// @Accept json
// @Produce json
// @Param categoria body dto.CreateCategoriaRequest true "Dados da categoria"
// @Success 201 {object} dto.CategoriaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias [post]
func (h *CategoriaHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCategoriaRequest

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

	categoria, err := h.service.Create(&req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(categoria)
}

// GetByID godoc
// @Summary Buscar categoria por ID
// @Description Retorna uma categoria específica pelo seu ID
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} dto.CategoriaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias/{id} [get]
func (h *CategoriaHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		h.log.WithError(err).Warn("ID inválido")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "ID inválido",
		})
	}

	categoria, err := h.service.GetByID(uint(id))
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(categoria)
}

// GetByIDWithProdutos godoc
// @Summary Buscar categoria com produtos
// @Description Retorna uma categoria com a lista de seus produtos
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} dto.CategoriaWithProdutosResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias/{id}/produtos [get]
func (h *CategoriaHandler) GetByIDWithProdutos(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		h.log.WithError(err).Warn("ID inválido")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "ID inválido",
		})
	}

	categoria, err := h.service.GetByIDWithProdutos(uint(id))
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(categoria)
}

// GetAll godoc
// @Summary Listar todas as categorias
// @Description Retorna uma lista paginada de todas as categorias
// @Tags Categorias
// @Accept json
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {object} dto.CategoriaPaginatedResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias [get]
func (h *CategoriaHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	response, err := h.service.GetAll(page, pageSize)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(response)
}

// GetAllActive godoc
// @Summary Listar categorias ativas
// @Description Retorna uma lista de todas as categorias ativas (para seleção)
// @Tags Categorias
// @Accept json
// @Produce json
// @Success 200 {array} dto.CategoriaResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias/ativas [get]
func (h *CategoriaHandler) GetAllActive(c *fiber.Ctx) error {
	response, err := h.service.GetAllActive()
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(response)
}

// Update godoc
// @Summary Atualizar categoria
// @Description Atualiza os dados de uma categoria existente
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Param categoria body dto.UpdateCategoriaRequest true "Dados para atualização"
// @Success 200 {object} dto.CategoriaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias/{id} [put]
func (h *CategoriaHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		h.log.WithError(err).Warn("ID inválido")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "ID inválido",
		})
	}

	var req dto.UpdateCategoriaRequest
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

	categoria, err := h.service.Update(uint(id), &req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(categoria)
}

// Delete godoc
// @Summary Excluir categoria
// @Description Remove uma categoria pelo seu ID
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Categoria possui produtos"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias/{id} [delete]
func (h *CategoriaHandler) Delete(c *fiber.Ctx) error {
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
		Message: "Categoria excluída com sucesso",
	})
}

// handleServiceError trata os erros da camada de serviço
func (h *CategoriaHandler) handleServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrCategoriaNaoEncontrada):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	case errors.Is(err, service.ErrNomeDuplicado):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	case errors.Is(err, service.ErrCategoriaComProdutos):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	case errors.Is(err, service.ErrNomeVazio),
		errors.Is(err, service.ErrNomeMuitoCurto):
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

