package service

import (
	"errors"
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"log"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {
	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	})
	return err
}

func (s CatalogService) EditCategory(id int, input dto.CreateCategoryRequest) (*domain.Category, error) {
	existingCategory, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}

	if len(input.Name) > 0 {
		existingCategory.Name = input.Name
	}

	if input.ParentId > 0 {
		existingCategory.ParentId = input.ParentId
	}

	if len(input.ImageUrl) > 0 {
		existingCategory.ImageUrl = input.ImageUrl
	}

	if input.DisplayOrder > 0 {
		existingCategory.DisplayOrder = input.DisplayOrder
	}

	updadtedCategory, err := s.Repo.EditCategory(existingCategory)

	return updadtedCategory, err
}

func (s CatalogService) DeleteCategory(id int) error {
	err := s.Repo.DeleteCategory(id)
	if err != nil {
		log.Println("delete category error:", err)
		return errors.New("error deleting category")
	}
	return nil
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (s CatalogService) GetCategory(id int) (*domain.Category, error) {
	category, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}
	return category, nil
}

func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	product, err := s.Repo.FindProducts()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s CatalogService) GetProduct(id int) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(id)

	if err != nil {
		return nil, errors.New("product does not exist")
	}

	return product, nil
}

func (s CatalogService) GetSellerProduct(id int) ([]*domain.Product, error) {
	products, err := s.Repo.FindSellerProducts(id)

	if err != nil {
		return nil, errors.New("product does not exist")
	}

	return products, nil
}

func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {

	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		CategoryId:  input.CategoryId,
		Price:       input.Price,
		UserId:      int(user.ID),
		Stock:       input.Stock,
		ImageUrl:    input.ImageUrl,
	})

	return err
}
func (s CatalogService) EditProduct(id int, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {
	existingProduct, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product does not exist")
	}

	if existingProduct.UserId != int(user.ID) {
		return nil, errors.New("you do not have manage rights of this product")
	}

	if len(input.Name) > 0 {
		existingProduct.Name = input.Name
	}

	if input.Price > 0 {
		existingProduct.Price = input.Price
	}

	if len(input.Description) > 0 {
		existingProduct.Description = input.Description
	}

	if input.CategoryId > 0 {
		existingProduct.CategoryId = input.CategoryId
	}

	updatedProduct, err := s.Repo.EditProduct(existingProduct)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s CatalogService) UpdateProductStock(e *domain.Product) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(int(e.ID))
	if err != nil {
		return nil, errors.New("product does not exist")
	}

	if product.UserId != e.UserId {
		return nil, errors.New("you do not have manage rights of this product")
	}

	product.Stock = e.Stock
	updatedProduct, err := s.Repo.EditProduct(product)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s CatalogService) DeleteProduct(id int, user domain.User) error {
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return errors.New("product not found")
	}

	if product.UserId != int(user.ID) {
		return errors.New("you are not allowed to delete this product")
	}

	err = s.Repo.DeleteProduct(id)
	if err != nil {
		log.Println("delete product error:", err)
		return errors.New("error deleting product")
	}
	return nil
}
