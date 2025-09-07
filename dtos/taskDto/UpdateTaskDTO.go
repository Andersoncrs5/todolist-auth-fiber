package taskdto

type UpdateTaskDTO struct {
	Title  string             `json:"title" validate:"required,min=8,max=60"`
	Discription  string       `json:"discription" validate:"max=200"`
	Done         bool		  `json:"done"`	
}