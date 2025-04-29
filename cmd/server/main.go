package main

import (
	"friendsManagement/internal/config"
	"friendsManagement/internal/controller"
	"friendsManagement/internal/db"
	"friendsManagement/internal/handler"
	"friendsManagement/internal/repository"
	"friendsManagement/internal/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	config := config.LoadConfig()
	e := echo.New()
	db := db.InitDB(config)
	repo := repository.NewRepositoy(db)
	controller := controller.NewController(db, repo.UserRelationshipRepo)
	handler := handler.NewHandler(controller.UserRelationshipController)
	routes.RegisterUserRelationshipRoutes(e, handler.UserRelationshipHandler)
	e.Logger.Fatal(e.Start(config.PORT))
}
