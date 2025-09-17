package handlers

import (
	"github.com/labstack/echo/v4"
	"rastro/cmd/api/requests"
	"rastro/cmd/api/services"
	"rastro/common"
	"rastro/internal/models"
)

func (h *Handler) ChangePassword(c echo.Context) error {
	userService := services.NewUserService(h.DB)
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	payload := new(requests.UpdatePasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	err := userService.ChangeUserPassword(payload.NewPassword, user)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Password updated successfully", nil)
}
