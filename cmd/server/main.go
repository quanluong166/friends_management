package main

import (
	"friendsManagement/internal/db"
	"friendsManagement/internal/repository"
	"friendsManagement/internal/routes"
	"friendsManagement/internal/services"
	"friendsManagement/internal/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db.InitDB()
	initRepo := repository.NewRepositoy(db.DB)
	initUC := usecase.NewUsecase(initRepo.UserRelationshipRepo)
	initService := services.NewService(initUC.UserRelationshipUC)
	group := e.Group("/api/user")
	routes.RegisterUserRelationshipRoutes(group, initService.UserRelationshipService)
	e.Logger.Fatal(e.Start(":8080"))
}
