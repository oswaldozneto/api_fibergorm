package handler

import (
	"context"

	"api_fibergorm/internal/dto"
	"api_fibergorm/internal/service"
	arqhandler "api_fibergorm/pkg/arquitetura/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ProdutoHandler gerencia as requisições HTTP relacionadas a produtos
// Herda todas as funcionalidades do BaseHandlerImpl e adiciona métodos específicos
type ProdutoHandler struct {
	*arqhandler.BaseHandlerImpl[dto.CreateProdutoRequest, dto.UpdateProdutoRequest, dto.ProdutoResponse]
	produtoService service.ProdutoService
}

// NewProdutoHandler cria uma nova instância do handler de produtos
func NewProdutoHandler(s service.ProdutoService, log *logrus.Logger) *ProdutoHandler {
	config := arqhandler.DefaultHandlerConfig("Produto")

	baseHandler := arqhandler.NewBaseHandler[
		dto.CreateProdutoRequest,
		dto.UpdateProdutoRequest,
		dto.ProdutoResponse,
	](s, log, config)

	return &ProdutoHandler{
		BaseHandlerImpl: baseHandler,
		produtoService:  s,
	}
}

// GetByCategoriaID godoc
// @Summary Listar produtos por categoria
// @Description Retorna uma lista paginada de produtos de uma categoria específica
// @Tags Produtos
// @Accept json
// @Produce json
// @Param categoria_id path int true "ID da categoria"
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {object} arqdto.PaginatedResponse[dto.ProdutoResponse]
// @Failure 400 {object} arqdto.ErrorResponse
// @Failure 404 {object} arqdto.ErrorResponse
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/produtos/categoria/{categoria_id} [get]
func (h *ProdutoHandler) GetByCategoriaID(c *fiber.Ctx) error {
	categoriaID, err := h.ParseID(c, "categoria_id")
	if err != nil {
		return err
	}

	page, pageSize := h.getPaginationParams(c)

	ctx := context.Background()
	response, err := h.produtoService.GetByCategoriaID(ctx, categoriaID, page, pageSize)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(response)
}

// getPaginationParams extrai os parâmetros de paginação da query
func (h *ProdutoHandler) getPaginationParams(c *fiber.Ctx) (int, int) {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	return page, pageSize
}

// RegisterRoutes registra as rotas de produto (sobrescreve para adicionar rotas específicas)
func (h *ProdutoHandler) RegisterRoutes(router fiber.Router) {
	// Rotas específicas primeiro (devem vir antes das rotas com parâmetros)
	router.Get("/categoria/:categoria_id", h.GetByCategoriaID)

	// Rotas base do CRUD
	router.Post("/", h.Create)
	router.Get("/", h.GetAll)
	router.Get("/:id", h.GetByID)
	router.Put("/:id", h.Update)
	router.Delete("/:id", h.Delete)
}

// Documentação Swagger para os métodos herdados do BaseHandler
// (necessário para o swag gerar a documentação corretamente)

// Create godoc
// @Summary Criar um novo produto
// @Description Cria um novo produto com os dados fornecidos
// @Tags Produtos
// @Accept json
// @Produce json
// @Param produto body dto.CreateProdutoRequest true "Dados do produto"
// @Success 201 {object} dto.ProdutoResponse
// @Failure 400 {object} arqdto.ErrorResponse
// @Failure 409 {object} arqdto.ErrorResponse
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/produtos [post]

// GetByID godoc
// @Summary Buscar produto por ID
// @Description Retorna um produto específico pelo seu ID
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "ID do produto"
// @Success 200 {object} dto.ProdutoResponse
// @Failure 400 {object} arqdto.ErrorResponse
// @Failure 404 {object} arqdto.ErrorResponse
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/produtos/{id} [get]

// GetAll godoc
// @Summary Listar todos os produtos
// @Description Retorna uma lista paginada de todos os produtos
// @Tags Produtos
// @Accept json
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {object} arqdto.PaginatedResponse[dto.ProdutoResponse]
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/produtos [get]

// Update godoc
// @Summary Atualizar produto
// @Description Atualiza os dados de um produto existente
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "ID do produto"
// @Param produto body dto.UpdateProdutoRequest true "Dados para atualização"
// @Success 200 {object} dto.ProdutoResponse
// @Failure 400 {object} arqdto.ErrorResponse
// @Failure 404 {object} arqdto.ErrorResponse
// @Failure 409 {object} arqdto.ErrorResponse
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/produtos/{id} [put]

// Delete godoc
// @Summary Excluir produto
// @Description Remove um produto pelo seu ID
// @Tags Produtos
// @Accept json
// @Produce json
// @Param id path int true "ID do produto"
// @Success 200 {object} arqdto.SuccessResponse
// @Failure 400 {object} arqdto.ErrorResponse
// @Failure 404 {object} arqdto.ErrorResponse
// @Failure 500 {object} arqdto.ErrorResponse
// @Router /api/v1/produtos/{id} [delete]
