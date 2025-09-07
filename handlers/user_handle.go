package handlers

import (
	"strings"
	"time"
	dto "todolist-auth-fiber/dtos/userDto"
	"todolist-auth-fiber/services"
	"todolist-auth-fiber/utils"
	"todolist-auth-fiber/utils/crypto"
	mappers "todolist-auth-fiber/utils/mappers/user"
	"todolist-auth-fiber/utils/res"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validaterUser = validator.New()

type UserHandler interface {
	Create(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Revoke(c *fiber.Ctx) error
}

type userHandler struct {
	service     services.UserService
	taskService services.TaskService
}

func NewUserHandler(service services.UserService, taskService services.TaskService) UserHandler {
	return &userHandler{
		service:     service,
		taskService: taskService,
	}
}

func (h *userHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserDTO

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

	if err := validaterUser.Struct(req); err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" failed on "+err.Tag())
		}

		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[[]string]{
				Timestamp: time.Now(),
				Body:      errors,
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Inputs invalids",
			},
		)
	}

	checkEmail, code, err := h.service.ExistsByEmail(c.Context(), req.Email)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      code,
			Status:    false,
			Message:   "Error the check if email already exists!",
		}

		return c.Status(code).JSON(res)
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

	checkUserName, code, err := h.service.ExistsByUserName(c.Context(), req.Username)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      code,
			Status:    false,
			Message:   "Error the check if username exists!",
		}

		return c.Status(code).JSON(res)
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

	saved, code, err := h.service.Save(c.Context(), req)
	if err != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      code,
			Status:    false,
			Message:   "Error the save new user! Please try again later",
		}

		return c.Status(code).JSON(res)
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

	_, code, errRefreshToken := h.service.SetRefreshToken(c.Context(), saved, refreshToken)
	if errRefreshToken != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errRefreshToken.Error(),
				Code:      code,
				Status:    true,
				Message:   "Error internal in server! Please try again later",
			},
		)
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
	var req dto.LoginUserDTO

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

	if err := validaterUser.Struct(req); err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" failed on "+err.Tag())
		}

		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[[]string]{
				Timestamp: time.Now(),
				Body:      errors,
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Inputs invalids",
			},
		)
	}

	user, code, err := h.service.GetByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    false,
				Message:   "Login invalid",
			},
		)
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

	_, code, errRefreshToken := h.service.SetRefreshToken(c.Context(), user, refreshToken)
	if errRefreshToken != nil {
		res := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      errRefreshToken.Error(),
			Code:      code,
			Status:    true,
			Message:   "Error the set refresh token",
		}

		return c.Status(code).JSON(res)
	}

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[res.ResponseToken]{
			Timestamp: time.Now(),
			Body:      tokens,
			Code:      fiber.StatusCreated,
			Status:    true,
			Message:   "Welcome again",
		},
	)
}

func (h *userHandler) Me(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "Invalid Authorization header format",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	userID, err := utils.ExtractUserID(tokenString)
	if err != nil {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "You are not Authorization",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	user, code, err := h.service.GetById(c.Context(), userID)
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    false,
				Message:   "You are not Authorization",
			},
		)
	}

	userDto := mappers.UserToUserDTO(user)

	response := res.ResponseHttp[dto.UserDTO]{
		Timestamp: time.Now(),
		Body:      userDto,
		Code:      fiber.StatusOK,
		Status:    true,
		Message:   "Me",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *userHandler) Delete(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "You are not authorized",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "Invalid Authorization header format",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	userID, err := utils.ExtractUserID(token)
	if err != nil {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      err.Error(),
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "You are not Authorization",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	user, code, err := h.service.GetById(c.Context(), userID)
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    false,
				Message:   "You are not Authorization",
			},
		)
	}

	if code, err := h.service.Delete(c.Context(), user); err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    false,
				Message:   "Error the delete the user",
			},
		)
	}

	if _, err := h.taskService.DeleteAllByUserId(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusInternalServerError,
				Status:    false,
				Message:   "Error the delete all task of user",
			},
		)
	}

	response := res.ResponseHttp[string]{
		Timestamp: time.Now(),
		Body:      "",
		Code:      fiber.StatusOK,
		Status:    false,
		Message:   "Bye Bye",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *userHandler) Update(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "",
			},
		)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "Invalid Authorization header format",
			},
		)

	}

	userID, err := utils.ExtractUserID(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "You are not Authorization",
			},
		)
	}

	var req dto.UpdateUserDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Inputs invalds",
			},
		)
	}

	if err := validaterUser.Struct(req); err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" failed on "+err.Tag())
		}

		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[[]string]{
				Timestamp: time.Now(),
				Body:      errors,
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Inputs invalids",
			},
		)
	}

	user, code, err := h.service.GetById(c.Context(), userID)
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    false,
				Message:   "You are not Authorization",
			},
		)
	}

	if req.Username != user.Username {
		checkUserName, code, err := h.service.ExistsByUserName(c.Context(), req.Username)
		if err != nil {
			return c.Status(code).JSON(
				res.ResponseHttp[string]{
					Timestamp: time.Now(),
					Body:      err.Error(),
					Code:      code,
					Status:    false,
					Message:   "Error the check if username exists!",
				},
			)
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
	}

	newPasswordHash, err := crypto.Encoder(req.Password)
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

	req.Password = newPasswordHash

	userUpdated, codeUpdate, err := h.service.Update(c.Context(), user, req)
	if err != nil {
		return c.Status(int(codeUpdate)).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      int(codeUpdate),
				Status:    false,
				Message:   "Error the update user",
			},
		)
	}

	userDto := mappers.UserToUserDTO(userUpdated)

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[dto.UserDTO]{
			Timestamp: time.Now(),
			Body:      userDto,
			Code:      fiber.StatusOK,
			Status:    true,
			Message:   "User updated with successfully",
		},
	)
}

func (h *userHandler) Revoke(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "",
			},
		)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "Invalid Authorization header format",
			},
		)

	}

	userID, err := utils.ExtractUserID(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "You are not Authorization",
			},
		)
	}

	user, codeGet, err := h.service.GetById(c.Context(), userID)
	if err != nil {
		return c.Status(codeGet).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      codeGet,
				Status:    false,
				Message:   "You are not Authorization",
			},
		)
	}

	_, code, err := h.service.SetRefreshToken(c.Context(), user, "")
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    true,
				Message:   "Error internal in server! Please try again later",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusOK,
			Status:    true,
			Message:   "See you later",
		},
	)
}
