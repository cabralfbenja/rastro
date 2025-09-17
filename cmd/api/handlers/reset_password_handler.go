package handlers

import (
	"encoding/base64"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/url"
	"rastro/cmd/api/requests"
	"rastro/cmd/api/services"
	"rastro/common"
	"rastro/internal/mailer"
)

func (h *Handler) ForgotPassword(c echo.Context) error {
	payload := new(requests.ForgotPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	// get user by email
	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	retrievedUser, err := userService.GetUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "User not found, register this email")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred while retrieving user")
	}

	// generate reset token
	resetToken, err := appTokenService.ResetPasswordToken(*retrievedUser)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	// send email with reset token
	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(retrievedUser.Email))
	frontendURL, err := url.Parse(payload.FrontendURL)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid frontend URL")
	}

	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", resetToken.Token)

	frontendURL.RawQuery = query.Encode()

	mailData := mailer.EmailData{
		Subject: "Request for Password Reset",
		Meta: struct {
			Token       string
			FrontendUrl string
		}{
			Token:       resetToken.Token,
			FrontendUrl: frontendURL.String(),
		}}
	err = h.Mailer.Send(payload.Email, "forgot_password.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}

	return common.SendSuccessResponse(c, "Password reset email sent successfully, please check your inbox", nil)
}

func (h *Handler) ResetPassword(c echo.Context) error {
	payload := new(requests.ResetPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	email, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	retrievedUser, err := userService.GetUserByEmail(string(email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Invalid password reset token")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	appToken, err := appTokenService.ValidateResetPasswordToken(*retrievedUser, payload.Token)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	err = userService.ChangeUserPassword(payload.Password, *retrievedUser)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	appTokenService.InvalidateToken(retrievedUser.ID, *appToken)
	return common.SendSuccessResponse(c, "Password reset successfully", nil)
}
