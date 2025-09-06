package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
)

func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://meudominio.com",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length, Authorization",
		AllowCredentials: true,
	})
}
