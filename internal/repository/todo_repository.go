package repository

import (
	"context"
	"todo-app/internal/domain"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"todo-app/utils"
)

type todoRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewTodoRepository(db *gorm.DB, tracer trace.Tracer) domain.TodoRepository {
	return &todoRepository{db: db, tracer: tracer}
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	_, span := r.tracer.Start(ctx, "TodoRepository.Create")
	defer span.End()

	span.AddEvent("todo_data", trace.WithAttributes(
		attribute.String("input", utils.ToJSONString(todo)),
	))

	result := r.db.Create(todo)
	if result.Error != nil {
		span.RecordError(result.Error)
		return result.Error
	}
	span.AddEvent("created_todo", trace.WithAttributes(
		attribute.String("result", todo.ID),
		attribute.Int64("rows_affected", result.RowsAffected),
	))
	return result.Error
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	_, span := r.tracer.Start(ctx, "TodoRepository.GetByID")
	defer span.End()

	todo := &domain.Todo{}
	err := r.db.First(&todo, id).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(todo)),
	))
	return todo, nil
}

func (r *todoRepository) GetAll(ctx context.Context) ([]domain.Todo, error) {
	_, span := r.tracer.Start(ctx, "TodoRepository.GetAll")
	defer span.End()

	var todos []domain.Todo
	err := r.db.Find(&todos).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(todos)),
	))
	return todos, err
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	_, span := r.tracer.Start(ctx, "TodoRepository.Update")
	defer span.End()

	result := r.db.Save(todo)
	if result.Error != nil {
		span.RecordError(result.Error)
		return result.Error
	}
	span.AddEvent("updated_todo", trace.WithAttributes(
		attribute.String("result", todo.ID),
		attribute.Int64("rows_affected", result.RowsAffected),
	))
	return result.Error
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, span := r.tracer.Start(ctx, "TodoRepository.Delete")
	defer span.End()

	span.AddEvent("id", trace.WithAttributes(
		attribute.String("id", id.String()),
	))

	result := r.db.Delete(&domain.Todo{}, id)
	if result.Error != nil {
		span.RecordError(result.Error)
		return result.Error
	}
	span.AddEvent("deleted", trace.WithAttributes(
		attribute.String("result", id.String()),
		attribute.Int64("rows_affected", result.RowsAffected),
	))
	return nil
}
