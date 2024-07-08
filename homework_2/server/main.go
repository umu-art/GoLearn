package main

import (
	"GoLearn/homework_2/server/account"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var secretKey = uuid.New().String()

func main() {
	println("Секретный ключ:", secretKey)

	accountsHandler := account.New(secretKey)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/api/account/create", accountsHandler.CreateAccount)
	e.GET("/api/account", accountsHandler.GetAccount)
	e.DELETE("/api/account", accountsHandler.DeleteAccount)
	e.PATCH("/api/account", accountsHandler.PatchAccount)
	e.POST("/api/account/rename", accountsHandler.ChangeAccount)

	e.GET("/api/accounts", accountsHandler.GetAll)

	e.GET("/api/actuator", accountsHandler.Actuator)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
