package taskdto

type CreateTaskDTO struct {
	Title  string             `json:"title" validate:"required,min=8,max=60"`
	Discription  string       `json:"discription" validate:"max=200"`
}