package main

import (
	"github.com/labstack/echo/v4"

	"github.com/zumikiti/echo-test/handler"
)

func main() {
	e := echo.New()

	// routes
	e.GET("/users", handler.GetUsers)
	e.POST("/users", handler.StoreUser)
	e.GET("/users/:id", handler.GetUser)
	e.PUT("/users/:id", handler.UpdateUser)
	e.DELETE("/users/:id", handler.DeleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
