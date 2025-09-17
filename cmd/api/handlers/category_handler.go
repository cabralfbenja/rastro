package handlers

import (
	"errors"
	"rastro/cmd/api/requests"
	"rastro/cmd/api/services"
	"rastro/common"
	"rastro/internal/app_errors"
	"rastro/internal/models"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ListCategories(c echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)

	categories, err := categoryService.List()
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "ok", categories)
}

func (h *Handler) CreateCategory(c echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)
	_, ok := c.Get("user").(models.UserModel)

	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	payload := new(requests.CreateCategoryRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	payload.IsCustom = true
	category, err := categoryService.Create(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Category created successfully", category)
}

func (h *Handler) DeleteCategory(c echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)
	_, ok := c.Get("user").(models.UserModel)

	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	var categoryId requests.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &categoryId)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid category ID")
	}

	err = categoryService.DeleteById(categoryId.ID)
	if err != nil {
		if errors.Is(err, app_errors.NewNotFoundError(err.Error())) {
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendBadRequestResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "category deleted", nil)
}
