package main

import (
	"github.com/labstack/echo/v4"
	"github.com/quanluong166/friends_management/internal/config"
	"github.com/quanluong166/friends_management/internal/controller"
	"github.com/quanluong166/friends_management/internal/db"
	"github.com/quanluong166/friends_management/internal/handler"
	"github.com/quanluong166/friends_management/internal/repository"
	"github.com/quanluong166/friends_management/internal/routes"

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
