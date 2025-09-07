package handlers

import (
	"strconv"
	"strings"
	"time"
	taskdto "todolist-auth-fiber/dtos/taskDto"
	"todolist-auth-fiber/models"
	"todolist-auth-fiber/services"
	"todolist-auth-fiber/utils"
	"todolist-auth-fiber/utils/pagination"
	"todolist-auth-fiber/utils/res"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validater = validator.New()

type TaskHandler interface {
	GetById(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	ChangeStatus(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type taskHandler struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) TaskHandler {
	return &taskHandler{service: service}
}

func (h *taskHandler) GetById(c *fiber.Ctx) error {
	id := c.Params("id")

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

	oid, errParseId := primitive.ObjectIDFromHex(id)
	if errParseId != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errParseId.Error(),
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Id invalid",
			},
		)

	}

	task, code, errGet := h.service.GetById(c.Context(), oid)
	if errGet != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errGet.Error(),
				Code:      code,
				Status:    false,
				Message:   "Error the get tasks",
			},
		)
	}

	if task.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusForbidden,
				Status:    false,
				Message:   "You are not authorized to see this task",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[*models.Todo]{
			Timestamp: time.Now(),
			Body:      task,
			Code:      fiber.StatusOK,
			Status:    false,
			Message:   "Task found with successfully!",
		},
	)
}

func (h *taskHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

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

	if id == "" {
		response := res.ResponseHttp[string]{
			Timestamp: time.Now(),
			Body:      "",
			Code:      fiber.StatusUnauthorized,
			Status:    false,
			Message:   "Id is required",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	oid, errParseId := primitive.ObjectIDFromHex(id)
	if errParseId != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errParseId.Error(),
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Id invalid",
			},
		)

	}

	task, code, errGet := h.service.GetById(c.Context(), oid)
	if errGet != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errGet.Error(),
				Code:      code,
				Status:    false,
				Message:   "Error the get tasks",
			},
		)
	}

	if task.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusForbidden,
				Status:    false,
				Message:   "You are not authorized to delete this task",
			},
		)
	}

	if _, err := h.service.Delete(c.Context(), oid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusInternalServerError,
				Status:    false,
				Message:   "Error the to delete task",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[*models.Todo]{
			Timestamp: time.Now(),
			Body:      task,
			Code:      fiber.StatusOK,
			Status:    true,
			Message:   "Task deleted with successfully!",
		},
	)
}

func (h *taskHandler) Create(c *fiber.Ctx) error {
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

	var req taskdto.CreateTaskDTO

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Inputs invalids",
			},
		)
	}

	if err := validater.Struct(req); err != nil {
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

	saved, code, err := h.service.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    false,
				Message:   "Inputs invalids",
			},
		)
	}

	return c.Status(201).JSON(
		res.ResponseHttp[*models.Todo]{
			Timestamp: time.Now(),
			Body:      saved,
			Code:      201,
			Status:    true,
			Message:   "Task created with successfully!",
		},
	)
}

func (h *taskHandler) ChangeStatus(c *fiber.Ctx) error {
	id := c.Params("id")

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

	if id == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "Id is required",
			},
		)
	}

	oid, errParseId := primitive.ObjectIDFromHex(id)
	if errParseId != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errParseId.Error(),
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Id invalid",
			},
		)

	}

	task, code, errGet := h.service.GetById(c.Context(), oid)
	if errGet != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errGet.Error(),
				Code:      code,
				Status:    false,
				Message:   "Error the get tasks",
			},
		)
	}

	if task.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusForbidden,
				Status:    false,
				Message:   "You are not authorized to change status this task",
			},
		)
	}

	taskChanged, code, err := h.service.ChangeStatus(c.Context(), oid, task)
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[*models.Todo]{
				Timestamp: time.Now(),
				Body:      task,
				Code:      code,
				Status:    false,
				Message:   "Error the change task status.",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[*models.Todo]{
			Timestamp: time.Now(),
			Body:      taskChanged,
			Code:      fiber.StatusOK,
			Status:    true,
			Message:   "Task status changed with successfully!",
		},
	)
}

func (h *taskHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

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

	oid, errParseId := primitive.ObjectIDFromHex(id)
	if errParseId != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errParseId.Error(),
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Id invalid",
			},
		)

	}

	var req taskdto.UpdateTaskDTO

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusBadRequest,
				Status:    false,
				Message:   "Inputs invalids",
			},
		)
	}

	if err := validater.Struct(req); err != nil {
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

	task, code, errGet := h.service.GetById(c.Context(), oid)
	if errGet != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      errGet.Error(),
				Code:      code,
				Status:    false,
				Message:   "Error the get tasks",
			},
		)
	}

	if task.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusForbidden,
				Status:    false,
				Message:   "You are not authorized to updated this task",
			},
		)
	}

	taskUpdated, code, err := h.service.Update(c.Context(), oid, req)
	if err != nil {
		return c.Status(code).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      code,
				Status:    false,
				Message:   "Error the to update task",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[*models.Todo]{
			Timestamp: time.Now(),
			Body:      taskUpdated,
			Code:      fiber.StatusOK,
			Status:    false,
			Message:   "Task updated with successfully!",
		},
	)
}

func (h *taskHandler) GetAll(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      "",
				Code:      fiber.StatusUnauthorized,
				Status:    false,
				Message:   "Missing Authorization header",
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
				Message:   "You are not authorized",
			},
		)
	}

	title := c.Query("title", "")
	doneParam := c.Query("done", "")
	var done *bool
	if doneParam != "" {
		val := doneParam == "true"
		done = &val
	}

	var createdAtBefore, createdAtAfter time.Time
	if beforeStr := c.Query("created_before"); beforeStr != "" {
		createdAtBefore, _ = time.Parse(time.RFC3339, beforeStr)
	}
	if afterStr := c.Query("created_after"); afterStr != "" {
		createdAtAfter, _ = time.Parse(time.RFC3339, afterStr)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	tasks, total, err := h.service.GetAll(c.Context(), userID, title, done, createdAtBefore, createdAtAfter, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			res.ResponseHttp[string]{
				Timestamp: time.Now(),
				Body:      err.Error(),
				Code:      fiber.StatusInternalServerError,
				Status:    false,
				Message:   "Error while fetching tasks",
			},
		)
	}

	response := pagination.Page[models.Todo]{
		Items:     tasks,
		Total:     total,
		PageIndex: page,
		PageSize:  pageSize,
	}

	return c.Status(fiber.StatusOK).JSON(
		res.ResponseHttp[pagination.Page[models.Todo]]{
			Timestamp: time.Now(),
			Body:      response,
			Code:      fiber.StatusOK,
			Status:    true,
			Message:   "Tasks retrieved successfully",
		},
	)
}
