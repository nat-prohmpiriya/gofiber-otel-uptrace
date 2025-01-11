package repository

import (
	"context"
	"todo-app/internal/domain"
	"todo-app/pkg/otel"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
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
	logger := otel.NewTraceLogger(span)

	logger.Input(todo)

	result := r.db.Create(todo)
	if result.Error != nil {
		logger.Error(result.Error)
		return result.Error
	}
	logger.Output(todo)
	return result.Error
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	_, span := r.tracer.Start(ctx, "TodoRepository.GetByID")
	defer span.End()
	logger := otel.NewTraceLogger(span)

	logger.Input(map[string]interface{}{
		"id": id.String(),
	})
	todo := &domain.Todo{}
	err := r.db.First(&todo, id).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Output(todo)

	return todo, nil
}

func (r *todoRepository) GetAll(ctx context.Context) ([]domain.Todo, error) {
	_, span := r.tracer.Start(ctx, "TodoRepository.GetAll")
	defer span.End()
	logger := otel.NewTraceLogger(span)

	var todos []domain.Todo
	err := r.db.Find(&todos).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Output(map[string]interface{}{
		"todos": todos,
		"err":   err,
	})
	return todos, err
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	_, span := r.tracer.Start(ctx, "TodoRepository.Update")
	defer span.End()
	logger := otel.NewTraceLogger(span)

	logger.Input(todo)

	result := r.db.Save(todo)
	if result.Error != nil {
		logger.Error(result.Error)
		return result.Error
	}
	logger.Output(map[string]interface{}{
		"result":        todo,
		"rows_affected": result.RowsAffected,
	})
	return result.Error
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, span := r.tracer.Start(ctx, "TodoRepository.Delete")
	defer span.End()

	logger := otel.NewTraceLogger(span)

	logger.Input(map[string]interface{}{
		"id": id.String(),
	})

	result := r.db.Delete(&domain.Todo{}, id)
	if result.Error != nil {
		logger.Error(result.Error)
		return result.Error
	}
	logger.Output(map[string]interface{}{
		"result":        id,
		"rows_affected": result.RowsAffected,
	})
	return nil
}
