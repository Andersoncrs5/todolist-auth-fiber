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

	token, err := utils.GenerateAccessToken(saved)
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

func (h *userHandler) Login(c *fiber.Ctx) error {
	var req userDto.LoginUserDTO

	if err := c.BodyParser(&req); err != nil {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusBadRequest,
			Status:    false,
			Message:   "Inputs invalid",
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	user, err := h.service.GetByEmail(c.Context(), req.Email)
	if err != nil {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "Login invalid",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	check := crypto.Compare(req.Password, user.Password)
	if check == false {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "Login invalid",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	token, err := utils.GenerateAccessToken(user)
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

	refreshToken, err := utils.GenerateRefreshToken(user)
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
		Message:   "Welcome again",
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
