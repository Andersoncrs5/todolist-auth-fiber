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
	GetById(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Delete(ctx context.Context, user *models.User) error
	Save(ctx context.Context, dto userDto.CreateUserDTO) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *models.User, dto userDto.UpdateUserDTO) (*models.User, uint, error)
	ExistsByUserName(ctx context.Context, UserName string) (bool, error)
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

func (u *userService) GetById(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	user, err := u.repo.GetId(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.repo.GetEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}

func (u *userService) Delete(ctx context.Context, user *models.User) error { 
	err := u.repo.Delete(ctx, user.ID)
	if err != nil {
		return err
	}

	return err
}

func (u *userService) Save(ctx context.Context, dto userDto.CreateUserDTO) (*models.User, error) {
	var user models.User

	user.Username = dto.Username
	user.Email = dto.Email
	user.Password = dto.Password

	saved, err := u.repo.Save(ctx, &user)
	if err != nil {
		return nil, err
	}

	return saved, nil
} 

func (u *userService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, fmt.Errorf("Email is required")
	}

	check, err := u.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return check, nil
}

func (u *userService) ExistsByUserName(ctx context.Context, UserName string) (bool, error) {
	if UserName == "" {
		return false, fmt.Errorf("UserName is required")
	}

	check, err := u.repo.ExistsByUserName(ctx, UserName)
	if err != nil {
		return false, err
	}

	return check, nil
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