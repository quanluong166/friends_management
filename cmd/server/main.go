package main

import (
	"friendsManagement/internal/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	repository.InitDB()
	e.Logger.Fatal(e.Start(":8080"))
}
