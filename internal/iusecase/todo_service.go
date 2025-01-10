package usecase

import (
	"context"
	"todo-app/internal/domain"

	"go.opentelemetry.io/otel"
)

type TodoService struct {
	repo domain.TodoRepository
}

func NewTodoService(repo domain.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	tracer := otel.Tracer("todo-service")
	_, span := tracer.Start(ctx, "TodoService.CreateTodo")
	defer span.End()

	return s.repo.Create(ctx, todo)
}

func (s *TodoService) GetTodoByID(ctx context.Context, id uint) (*domain.Todo, error) {
	tracer := otel.Tracer("todo-service")
	_, span := tracer.Start(ctx, "TodoService.GetTodoByID")
	defer span.End()

	return s.repo.GetByID(ctx, id)
}

func (s *TodoService) GetAllTodos(ctx context.Context) ([]domain.Todo, error) {
	tracer := otel.Tracer("todo-service")
	_, span := tracer.Start(ctx, "TodoService.GetAllTodos")
	defer span.End()

	return s.repo.GetAll(ctx)
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo *domain.Todo) error {
	tracer := otel.Tracer("todo-service")
	_, span := tracer.Start(ctx, "TodoService.UpdateTodo")
	defer span.End()

	return s.repo.Update(ctx, todo)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id uint) error {
	tracer := otel.Tracer("todo-service")
	_, span := tracer.Start(ctx, "TodoService.DeleteTodo")
	defer span.End()

	return s.repo.Delete(ctx, id)
}
