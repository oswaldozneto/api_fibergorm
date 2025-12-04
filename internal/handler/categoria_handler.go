package handler

import (
	"context"

	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/service"
	arqhandler "api_fibergorm/pkg/arquitetura/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// CategoriaHandler gerencia as requisições HTTP relacionadas a categorias
// Herda todas as funcionalidades do BaseHandlerImpl e adiciona métodos específicos
type CategoriaHandler struct {
	*arqhandler.BaseHandlerImpl[dto.CreateCategoriaRequest, dto.UpdateCategoriaRequest, dto.CategoriaResponse]
	categoriaService service.CategoriaService
}

// NewCategoriaHandler cria uma nova instância do handler de categorias
func NewCategoriaHandler(s service.CategoriaService, log *logrus.Logger) *CategoriaHandler {
	config := arqhandler.DefaultHandlerConfig("Categoria")

	baseHandler := arqhandler.NewBaseHandler(s, log, config)

	return &CategoriaHandler{
		BaseHandlerImpl:  baseHandler,
		categoriaService: s,
	}
}

// GetByIDWithProdutos godoc
// @Summary Buscar categoria com produtos
// @Description Retorna uma categoria com a lista de seus produtos
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} dto.CategoriaWithProdutosResponse
// @Failure 400 {object} arqdto.ErrorResponse
// @Failure 404 {object} arqdto.ErrorResponse
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/categorias/{id}/produtos [get]
func (h *CategoriaHandler) GetByIDWithProdutos(c *fiber.Ctx) error {
	id, err := h.ParseID(c, "id")
	if err != nil {
		return err
	}

	ctx := context.Background()
	categoria, err := h.categoriaService.GetByIDWithProdutos(ctx, id)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(categoria)
}

// GetAllActive godoc
// @Summary Listar categorias ativas
// @Description Retorna uma lista de todas as categorias ativas (para seleção)
// @Tags Categorias
// @Accept json
// @Produce json
// @Success 200 {array} dto.CategoriaResponse
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/categorias/ativas [get]
func (h *CategoriaHandler) GetAllActive(c *fiber.Ctx) error {
	ctx := context.Background()
	response, err := h.categoriaService.GetAllActive(ctx)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(response)
}

// RegisterRoutes registra as rotas de categoria (sobrescreve para adicionar rotas específicas)
func (h *CategoriaHandler) RegisterRoutes(router fiber.Router) {
	// Rotas específicas primeiro (devem vir antes das rotas com parâmetros)
	router.Get("/ativas", h.GetAllActive)

	// Rotas base do CRUD
	router.Post("/", h.Create)
	router.Get("/", h.GetAll)
	router.Get("/:id", h.GetByID)
	router.Get("/:id/produtos", h.GetByIDWithProdutos)
	router.Put("/:id", h.Update)
	router.Delete("/:id", h.Delete)
}
