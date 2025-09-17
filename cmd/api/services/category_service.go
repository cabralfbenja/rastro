package services

import (
	"errors"
	"rastro/cmd/api/requests"
	"rastro/internal/app_errors"
	"rastro/internal/models"
	"strings"

	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

func (c CategoryService) List() ([]*models.CategoryModel, error) {
	var categories []*models.CategoryModel
	result := c.DB.Find(&categories)

	if result.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}

	return categories, nil
}

func (c CategoryService) Create(data *requests.CreateCategoryRequest) (*models.CategoryModel, error) {
	slug := strings.ToLower(data.Name)
	slug = strings.Replace(slug, " ", "-", -1)
	category := &models.CategoryModel{
		Name:     data.Name,
		Slug:     slug,
		IsCustom: data.IsCustom,
	}

	result := c.DB.Where(models.CategoryModel{Slug: slug, Name: data.Name}).FirstOrCreate(category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return category, nil
		}
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

func (c CategoryService) GetByID(id uint) (*models.CategoryModel, error) {
	var category *models.CategoryModel
	result := c.DB.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_errors.NewNotFoundError("Category not found")
		}
		return nil, errors.New("failed to find category")
	}

	return category, nil
}

func (c CategoryService) DeleteById(id uint) error {
	var category *models.CategoryModel
	category, err := c.GetByID(id)
	if err != nil {
		return err
	}

	if !category.IsCustom {
		return errors.New("only custom categories can be deleted")
	}

	c.DB.Delete(category)

	return nil
}
