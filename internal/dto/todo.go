package dto

type CreateTodoRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

type UpdateTodoRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Done  bool   `json:"done"`
}

type GetTodoResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Done      bool   `json:"done"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type TodoQuery struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type GetAllTodosResponse struct {
	Todos []GetTodoResponse `json:"todos"`
}

type DeleteTodoRequest struct {
	ID string `json:"id" validate:"required"`
}
