package main

import (
	"GoLearn/homework_2/server/account"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	accountsHandler := account.New()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/account/create", accountsHandler.CreateAccount)
	e.GET("/account", accountsHandler.GetAccount)
	e.DELETE("/account", accountsHandler.DeleteAccount)
	e.PATCH("/account", accountsHandler.PatchAccount)
	e.POST("/account/rename", accountsHandler.ChangeAccount)

	e.GET("/actuator", accountsHandler.Actuator)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
