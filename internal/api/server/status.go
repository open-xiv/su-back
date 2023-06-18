package server

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/pkg/model"
	"net/http"
)

func Status(c echo.Context) error {
	status := model.ServerStatus{
		Status:  "Working",
		Version: "v0.1.0",
	}
	return c.JSONPretty(http.StatusOK, status, "  ")
}
