package mappers

import (
	"todolist-auth-fiber/dtos/userDto"
	"todolist-auth-fiber/models"
)

func UserToUserDTO(user *models.User) userDto.UserDTO {
	return userDto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}