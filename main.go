package main

import (
	"todolist-auth-fiber/config"
	"todolist-auth-fiber/handlers"
	repository "todolist-auth-fiber/repositories"
	"todolist-auth-fiber/routers"
	"todolist-auth-fiber/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.ConnectDB()
	db := config.GetDB()
	
	taskRepository := repository.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepository)
	taskHandler := handlers.NewTaskHandler(taskService)

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService, taskService)

	routers.UserRouter(app, userHandler)
	routers.TaskRouter(app, taskHandler)

	app.Listen(":8080")
}