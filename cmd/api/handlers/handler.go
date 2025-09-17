package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"rastro/internal/mailer"
)

type Handler struct {
	DB     *gorm.DB
	Logger echo.Logger
	Mailer *mailer.Mailer
}

func (handler *Handler) BindBodyRequest(c echo.Context, payload interface{}) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return errors.New("Failed to bind request body")
	}
	return nil
}
