package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"rastro/common"
	"rastro/internal/models"
	"strings"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return common.SendUnauthorizedResponse(c, "Please provide Bearer token")
		}
		authHeaderSplit := strings.Split(authHeader, " ")
		authToken := authHeaderSplit[1]

		claims, err := common.ParseJWTSignedToken(authToken)
		if err != nil {
			return common.SendUnauthorizedResponse(c, err.Error())
		}

		if claims.IsExpired() {
			return common.SendUnauthorizedResponse(c, "Token has expired")
		}

		var user models.UserModel
		result := appMiddleware.DB.First(&user, claims.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.SendUnauthorizedResponse(c, "Invalid access token")
		}

		if result.Error != nil {
			return common.SendInternalServerErrorResponse(c, "Invalid access token")
		}

		c.Set("user", user)
		return next(c)
	}
}

// supply jwt
// middleware intercepts and validates the JWT token
// if jwt is not valid, bounce back user
// middleware attaches current user to the context
