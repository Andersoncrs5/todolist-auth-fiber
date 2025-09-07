package userDto

type CreateUserDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,min=10,max=150"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}