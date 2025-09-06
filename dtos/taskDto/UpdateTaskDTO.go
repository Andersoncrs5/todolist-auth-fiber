package taskdto

type UpdateTaskDTO struct {
	Title        string       `json:"title"`
	Discription  string       `json:"discription"`
	Done         bool		  `json:"done"`	
}