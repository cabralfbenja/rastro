package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type healtCheck struct {
	Health bool `json:"health"`
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, healtCheck{Health: true})
}
