package userDto

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email,min=10,max=150"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}