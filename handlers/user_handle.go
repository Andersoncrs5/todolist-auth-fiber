package handlers

import (
	"time"
	"todolist-auth-fiber/dtos/userDto"
	"todolist-auth-fiber/services"
	"todolist-auth-fiber/utils"
	"todolist-auth-fiber/utils/crypto"
	"todolist-auth-fiber/utils/res"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return &userHandler{service: service}
}

func (h *userHandler) Create(c *fiber.Ctx) error {
	var req userDto.CreateUserDTO

	if err := c.BodyParser(&req); err != nil {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusBadRequest,
			Status:    false,
			Message:   "Inputs inválidos",
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	checkEmail, err := h.service.ExistsByEmail(c.Context(), req.Email)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusInternalServerError,
			Status:    false,
			Message:   "Error the check if email already exists!",
		}

		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	if checkEmail == true {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusConflict,
			Status:    false,
			Message:   "Email already exists",
		}

		return c.Status(fiber.StatusConflict).JSON(res)
	}

	checkUserName, err := h.service.ExistsByUserName(c.Context(), req.Username)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusInternalServerError,
			Status:    false,
			Message:   "Error the check if username exists!",
		}

		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	if checkUserName == true {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusConflict,
			Status:    false,
			Message:   "Username already exists",
		}

		return c.Status(fiber.StatusConflict).JSON(res)
	}

	password, err := crypto.Encoder(req.Password)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusInternalServerError,
			Status:    false,
			Message:   "Error the encoder password!",
		}

		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	req.Password = password

	saved, err := h.service.Save(c.Context(), req)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusInternalServerError,
			Status:    false,
			Message:   "Error the save new user! Please try again later",
		}

		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	token, err := utils.GenerateToken(saved)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusInternalServerError,
			Status:    false,
			Message:   "Error in server! Please try again later",
		}

		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	refreshToken, err := utils.GenerateRefreshToken(saved)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusInternalServerError,
			Status:    false,
			Message:   "Error in server! Please try again later",
		}

		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	tokens := res.ResponseToken{
		Token:        token,
		RefreshToken: refreshToken,
	}

	res := res.ResponseHttp[res.ResponseToken]{
		Timestamp: time.Now(),
		Body:      tokens,
		Code:      fiber.StatusCreated,
		Status:    true,
		Message:   "Welcome",
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}
