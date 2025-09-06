package ratelimiting

import (
	"time"
	"todolist-auth-fiber/utils/res"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func CreateRate() fiber.Handler {
	return limiter.New(limiter.Config{
		Max: 70,
		Expiration: 10 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()               
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				res.ResponseHttp[string]{
					Timestamp: time.Now(),
					Body: "",
					Code: fiber.StatusTooManyRequests,
					Status: false,
					Message: "Too many requests, please try again later.",
				},
			)
		},
	})
}

func GetRate() fiber.Handler {
	return limiter.New(limiter.Config{
		Max: 150,
		Expiration: 15 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()               
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				res.ResponseHttp[string]{
					Timestamp: time.Now(),
					Body: "",
					Code: fiber.StatusTooManyRequests,
					Status: false,
					Message: "Too many requests, please try again later.",
				},
			)
		},
	})
}

func DeleteRate() fiber.Handler {
	return limiter.New(limiter.Config{
		Max: 60,
		Expiration: 15 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()               
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				res.ResponseHttp[string]{
					Timestamp: time.Now(),
					Body: "",
					Code: fiber.StatusTooManyRequests,
					Status: false,
					Message: "Too many requests, please try again later.",
				},
			)
		},
	})
}

func UpdateRate() fiber.Handler {
	return limiter.New(limiter.Config{
		Max: 60,
		Expiration: 15 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()               
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				res.ResponseHttp[string]{
					Timestamp: time.Now(),
					Body: "",
					Code: fiber.StatusTooManyRequests,
					Status: false,
					Message: "Too many requests, please try again later.",
				},
			)
		},
	})
}

func CustomRate(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() 
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				res.ResponseHttp[string]{
					Timestamp: time.Now(),
					Body: "",
					Code: fiber.StatusTooManyRequests,
					Status: false,
					Message: "Too many requests, please try again later.",
				},
			)
		},
	})
}