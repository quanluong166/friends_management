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
	db := db.InitDB()
	repo := repository.NewRepositoy(db)
	usecase := usecase.NewUsecase(db, repo.UserRelationshipRepo)
	service := services.NewService(usecase.UserRelationshipUC)
	group := e.Group("/api/user")
	routes.RegisterUserRelationshipRoutes(group, service.UserRelationshipService)
	e.Logger.Fatal(e.Start(":8080"))
}
