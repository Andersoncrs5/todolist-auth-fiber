package services

import (
	"context"
	"fmt"
	"time"
	taskdto "todolist-auth-fiber/dtos/taskDto"
	"todolist-auth-fiber/models"
	repository "todolist-auth-fiber/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService interface {
	GetById(ctx context.Context, id primitive.ObjectID) (*models.Todo, int, error)
	Delete(ctx context.Context, id primitive.ObjectID) (int, error)
	Create(ctx context.Context, userID primitive.ObjectID, dto taskdto.CreateTaskDTO) (*models.Todo, int, error)
	ChangeStatus(ctx context.Context, id primitive.ObjectID, task models.Todo) (*models.Todo, int, error)
	Update(ctx context.Context, id primitive.ObjectID, dto taskdto.UpdateTaskDTO) (*models.Todo, int, error)
	DeleteAllByUserId(ctx context.Context, userId primitive.ObjectID) (int64, error)
	GetAll(
		ctx context.Context,
		userID primitive.ObjectID,
		title string,
		done *bool,
		createdAtBefore, createdAtAfter time.Time,
		page, pageSize int,
	) ([]models.Todo, int64, error)
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (s *taskService) GetById(ctx context.Context, id primitive.ObjectID) (*models.Todo, int, error) {
	task, code, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, code, err
	}

	if task == nil {
		return nil, 404, fmt.Errorf("Task not found")
	}

	return task, 200, nil
}

func (s *taskService) Delete(ctx context.Context, id primitive.ObjectID) (int, error) {
	code, err := s.repo.Delete(ctx, id)
	if err != nil {
		return code, err
	}

	return 200, nil
}

func (s *taskService) Create(ctx context.Context, userID primitive.ObjectID, dto taskdto.CreateTaskDTO) (*models.Todo, int, error) {
	var task models.Todo

	task.Title = dto.Title
	task.Discription = dto.Discription

	saved, code, err := s.repo.Create(ctx, userID, task)
	if err != nil {
		return nil, code, err
	}

	return saved, code, err
}

func (s *taskService) ChangeStatus(ctx context.Context, id primitive.ObjectID, task models.Todo) (*models.Todo, int, error) {
	taskChanged, code, err := s.repo.ChangeStatus(ctx, id, task)
	if err != nil {
		return nil, code, err
	}

	return taskChanged, code, nil
}

func (s *taskService) Update(ctx context.Context, id primitive.ObjectID, dto taskdto.UpdateTaskDTO) (*models.Todo, int, error) {
	updated, code, err := s.repo.Update(ctx, id, dto)
	if err != nil {
		return nil, code, err
	}

	return updated, code, nil
}

func (s *taskService) GetAll(
	ctx context.Context,
	userID primitive.ObjectID,
	title string,
	done *bool,
	createdAtBefore, createdAtAfter time.Time,
	page, pageSize int,
) ([]models.Todo, int64, error) {

	tasks, total, err := s.repo.GetAll(ctx, userID, title, done, createdAtBefore, createdAtAfter, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (s *taskService) DeleteAllByUserId(ctx context.Context, userId primitive.ObjectID) (int64, error) {
	result, err := s.repo.DeleteAllByUserId(ctx, userId);
	if err != nil {
		return 0, err
	}

	return result, nil	
}