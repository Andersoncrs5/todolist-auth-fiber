package routers

import (
	"todolist-auth-fiber/handlers"
	rate "todolist-auth-fiber/middleware/rateLimiting"

	"github.com/gofiber/fiber/v2"
)

func TaskRouter(app *fiber.App, taskHandler handlers.TaskHandler) {
	router := app.Group("/api/v1/tasks")

	router.Get("/:id", rate.GetRate(), taskHandler.GetById)
	router.Post("", rate.CreateRate(), taskHandler.Create)
	router.Delete("/:id", rate.DeleteRate(), taskHandler.Delete)
	router.Put("/:id", rate.UpdateRate(), taskHandler.Update)
	router.Put("/:id/status/done", rate.UpdateRate(), taskHandler.ChangeStatus)
	router.Get("", rate.GetRate(), taskHandler.GetAll)
}