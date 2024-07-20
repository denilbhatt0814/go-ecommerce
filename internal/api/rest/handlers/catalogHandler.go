package handlers

import (
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CatalogHandler struct {
	svc service.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	// create in instance of user service and inject to handler
	svc := service.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}
	handler := CatalogHandler{
		svc: svc,
	}

	// Public - Listing Products and categories
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)
	app.Get("/categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategory)

	// Private - manage Products and categories
	selRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)
	selRoutes.Post("/categories", handler.CreateCategories)
	selRoutes.Patch("/categories/:id", handler.EditCategory)
	selRoutes.Delete("/categories/:id", handler.DeleteCategory)

	selRoutes.Get("/products", handler.GetSellerProducts)
	selRoutes.Get("/products/:id", handler.GetProduct)
	selRoutes.Post("/products", handler.CreateProducts)
	selRoutes.Patch("/products/:id", handler.UpdateStock)
	selRoutes.Put("/products/:id", handler.EditProduct)
	selRoutes.Delete("/products/:id", handler.DeleteProduct)

}

func (h *CatalogHandler) GetCategories(ctx *fiber.Ctx) error {

	categories, err := h.svc.GetCategories()
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}

	return rest.SuccessResponse(ctx, "categories", categories)
}

func (h *CatalogHandler) GetCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	category, err := h.svc.GetCategory(id)
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}
	return rest.SuccessResponse(ctx, "category", category)
}

func (h *CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "create category request is not valid")
	}

	err = h.svc.CreateCategory(req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "CreateCategory", nil)
}

func (h *CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "update category request is not valid")
	}

	id, _ := strconv.Atoi(ctx.Params("id"))

	updatedCategory, err := h.svc.EditCategory(id, req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "category edited successfully", updatedCategory)
}

func (h *CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {

	id, _ := strconv.Atoi(ctx.Params("id"))
	err := h.svc.DeleteCategory(id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "category deleted succesfully", nil)
}

func (h *CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {

	var req dto.CreateProductRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "create product request is not valid")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	err = h.svc.CreateProduct(req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "CreateProduct", nil)
}

func (h *CatalogHandler) GetProducts(ctx *fiber.Ctx) error {

	products, err := h.svc.GetProducts()
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}

	return rest.SuccessResponse(ctx, "GetProducts", products)
}

func (h *CatalogHandler) GetSellerProducts(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	products, err := h.svc.GetSellerProduct(int(user.ID))
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}

	return rest.SuccessResponse(ctx, "GetSellerProduct", products)
}

func (h *CatalogHandler) GetProduct(ctx *fiber.Ctx) error {

	id, _ := strconv.Atoi(ctx.Params("id"))

	product, err := h.svc.Repo.FindProductById(id)
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}

	return rest.SuccessResponse(ctx, "product", product)
}

func (h *CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	var req dto.CreateProductRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "edit product request is not valid")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	product, err := h.svc.EditProduct(id, req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "EditProduct", product)
}

func (h *CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	var req dto.UpdateStockRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "update product stock request is not valid")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	product := &domain.Product{
		ID:     uint(id),
		Stock:  uint(req.Stock),
		UserId: int(user.ID),
	}

	updatedProduct, err := h.svc.UpdateProductStock(product)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "UpdateStock", updatedProduct)
}

func (h *CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	user := h.svc.Auth.GetCurrentUser(ctx)

	err := h.svc.DeleteProduct(id, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "DeleteProduct", nil)
}
