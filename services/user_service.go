package services

import (
	"context"
	"fmt"
	"todolist-auth-fiber/dtos/userDto"
	"todolist-auth-fiber/models"
	"todolist-auth-fiber/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetById(ctx context.Context, id primitive.ObjectID) (*models.User, int, error)
	GetByEmail(ctx context.Context, email string) (*models.User, int, error)
	Delete(ctx context.Context, user *models.User) (int, error)
	Save(ctx context.Context, dto userDto.CreateUserDTO) (*models.User, int, error)
	ExistsByEmail(ctx context.Context, email string) (bool, int, error)
	Update(ctx context.Context, user *models.User, dto userDto.UpdateUserDTO) (*models.User, uint, error)
	ExistsByUserName(ctx context.Context, UserName string) (bool, int, error)
	SetRefreshToken(ctx context.Context, user *models.User, refreshToken string) (*models.User, int, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService {
		repo: repo,
	}
}

func (u *userService) GetById(ctx context.Context, id primitive.ObjectID) (*models.User, int, error) {
	user, code, err := u.repo.GetId(ctx, id)
	if err != nil {
		return nil, code, err
	}

	if user == nil {
		return nil, code, fmt.Errorf("User not found")
	}

	return user, code, nil
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*models.User, int, error) {
	user, code, err := u.repo.GetEmail(ctx, email)
	if err != nil {
		return nil, code, err
	}

	if user == nil {
		return nil, 404, fmt.Errorf("User not found")
	}

	return user, 200, nil
}

func (u *userService) Delete(ctx context.Context, user *models.User) (int, error) { 
	code, err := u.repo.Delete(ctx, user.ID)
	if err != nil {
		return code, err
	}

	return code, nil
}

func (u *userService) Save(ctx context.Context, dto userDto.CreateUserDTO) (*models.User, int, error) {
	var user models.User

	user.Username = dto.Username
	user.Email = dto.Email
	user.Password = dto.Password

	saved, code, err := u.repo.Save(ctx, &user)
	if err != nil {
		return nil, code, err
	}

	return saved, code, nil
} 

func (u *userService) ExistsByEmail(ctx context.Context, email string) (bool, int, error) {
	if email == "" {
		return false, 400, fmt.Errorf("Email is required")
	}

	check, code, err := u.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, code, err
	}

	return check, code, nil
}

func (u *userService) ExistsByUserName(ctx context.Context, UserName string) (bool, int, error) {
	if UserName == "" {
		return false, 400, fmt.Errorf("UserName is required")
	}

	check, code, err := u.repo.ExistsByUserName(ctx, UserName)
	if err != nil {
		return false, code, err
	}

	return check, code, nil
}

func (u *userService) Update(ctx context.Context, user *models.User, dto userDto.UpdateUserDTO) (*models.User, uint, error) {
	updated, code, err := u.repo.Update(ctx, user.ID, dto)
	if err != nil {
		return nil, code, err
	}

	return updated, code, nil
}

func (u *userService) SetRefreshToken(ctx context.Context, user *models.User, refreshToken string) (*models.User, int, error) {
	userSaved, code, err := u.repo.SetRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return nil, code, err
	}

	return userSaved, 200, nil
}