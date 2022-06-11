package router

import (
	"go-qchang/batch"
	"go-qchang/datasource"
	"go-qchang/handlers"
	"go-qchang/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func New() *echo.Echo {
	e := echo.New()

	e.Logger.SetLevel(log.DEBUG)

	desk := datasource.New()
	service := services.NewCashierService(desk)
	handler := handlers.NewCashierHandler(service)
	batch.New(desk)

	// e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	c.Logger().Debugf("request: %s", reqBody)
	// 	c.Logger().Debugf("response: %s", resBody)
	// }))
	e.Use(middleware.Recover())

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"data": "pong",
		})
	})

	e.POST("/cashier/make_change", handler.ChangeMoney)
	e.POST("/cashier/transfer_in", handler.TransferMoneyIn)
	e.POST("/cashier/transfer_out", handler.TransferMoneyOut)

	return e
}
