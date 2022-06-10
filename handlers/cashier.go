package handlers

import (
	"go-qchang/models"
	"go-qchang/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CashierHandler interface {
	ChangeMoney(e echo.Context) error
}

type cashierHandler struct {
	svc services.CashierService
}

func NewCashierHandler(svc services.CashierService) CashierHandler {
	return &cashierHandler{
		svc: svc,
	}
}

func (h *cashierHandler) ChangeMoney(c echo.Context) error {
	var req models.ChangeMoneyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Bad request",
		})
	}

	changes, err := h.svc.ChangeMoney(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "No change available",
		})
	}

	return c.JSON(http.StatusOK, changes)
}
