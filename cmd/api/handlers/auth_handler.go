package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"os"
	"rastro/cmd/api/requests"
	"rastro/cmd/api/services"
	"rastro/common"
	"rastro/internal/mailer"
	"rastro/internal/models"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (h *Handler) RegisterHandler(c echo.Context) error {
	// bind request body
	payload := new(requests.RegisterUserRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	// Check if Username Exists
	_, err := userService.GetUserByUsername(payload.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) == false {
		return common.SendBadRequestResponse(c, "Username is already taken")
	}
	// Create the User
	registerdUser, err := userService.Register(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	// Send a Welcome message to the user
	mailData := mailer.EmailData{Subject: "Welcome To " + os.Getenv("APP_NAME") + " Signup", Meta: struct {
		FirstName string
		LoginLink string
	}{
		FirstName: *registerdUser.FirstName,
		LoginLink: "#",
	}}
	err = h.Mailer.Send(payload.Email, "welcome.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}
	// Send Respnse
	return common.SendSuccessResponse(c, "User registration successful", registerdUser)
}

func (h *Handler) LoginHandler(c echo.Context) error {
	userService := services.NewUserService(h.DB)

	// bind data / retrieving data sent by the client
	payload := new(requests.LoginUserRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate data sent by the client
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	// if the user with the supplied username exists
	userRetrieved, err := userService.GetUserByUsername(payload.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendBadRequestResponse(c, "Invalid credentials")
	}

	// compare password with the hashed password in the database
	if !common.ComparePasswordHash(payload.Password, userRetrieved.PasswordHash) {
		return common.SendBadRequestResponse(c, "Invalid credentials")
	}

	// if the password matches, generate a JWT token
	signedToken, refreshToken, err := common.GenerateJWT(*userRetrieved)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "User logged in", map[string]interface{}{
		"accessToken":  signedToken,
		"refreshToken": refreshToken,
		"userId":       userRetrieved.ID,
	})
}

func (h *Handler) GetAuthenticatedUser(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	return common.SendSuccessResponse(c, "Authenticated user retrieved", user)
}
