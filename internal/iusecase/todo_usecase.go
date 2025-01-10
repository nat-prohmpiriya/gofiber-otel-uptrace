package usecase

import (
	"context"
	"todo-app/internal/domain"

	"todo-app/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoService struct {
	repo   domain.TodoUsecase
	tracer trace.Tracer
}

func NewTodoService(repo domain.TodoUsecase, tracer trace.Tracer) *TodoService {
	return &TodoService{repo: repo, tracer: tracer}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	ctx, span := s.tracer.Start(ctx, "TodoService.CreateTodo")
	defer span.End()

	span.AddEvent("todo_data", trace.WithAttributes(
		attribute.String("input", utils.ToJSONString(todo)),
	))

	err := s.repo.Create(ctx, todo)
	if err != nil {
		span.RecordError(err)
		return err
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(todo)),
	))
	return nil
}

func (s *TodoService) GetTodoByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	_, span := s.tracer.Start(ctx, "TodoService.GetTodoByID")
	defer span.End()

	span.AddEvent("id", trace.WithAttributes(
		attribute.String("id", id.String()),
	))

	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(result)),
	))
	return result, nil
}

func (s *TodoService) GetAllTodos(ctx context.Context) ([]domain.Todo, error) {
	_, span := s.tracer.Start(ctx, "TodoService.GetAllTodos")
	defer span.End()

	result, err := s.repo.GetAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(result)),
	))
	return result, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo *domain.Todo) error {
	_, span := s.tracer.Start(ctx, "TodoService.UpdateTodo")
	defer span.End()

	span.AddEvent("todo_data", trace.WithAttributes(
		attribute.String("input", utils.ToJSONString(todo)),
	))

	err := s.repo.Update(ctx, todo)
	if err != nil {
		span.RecordError(err)
		return err
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(todo)),
	))
	return nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	_, span := s.tracer.Start(ctx, "TodoService.DeleteTodo")
	defer span.End()

	span.AddEvent("input", trace.WithAttributes(
		attribute.String("id", id.String()),
	))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		return err
	}

	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(id)),
	))
	return nil
}
