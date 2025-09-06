package routers

import (
	"time"
	"todolist-auth-fiber/handlers"
	rate "todolist-auth-fiber/middleware/rateLimiting"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App, userHandler handlers.UserHandler)  {
	user := app.Group("/api/v1/users")

	user.Get("", rate.GetRate(), userHandler.Me)
	user.Post("/register", rate.CreateRate(), userHandler.Create)
	user.Post("/login", rate.CustomRate(50, 15 * time.Second), userHandler.Login)
	user.Delete("", rate.DeleteRate(), userHandler.Delete)
	user.Put("", rate.UpdateRate(), userHandler.Update)
	user.Put("/revoke", rate.CustomRate(40, 10 * time.Second), userHandler.Revoke)
}