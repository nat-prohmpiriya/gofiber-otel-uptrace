package repository

import (
	"context"
	"todo-app/internal/domain"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) domain.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	tracer := otel.Tracer("todo-repository")
	_, span := tracer.Start(ctx, "TodoRepository.Create")
	defer span.End()

	return r.db.Create(todo).Error
}

func (r *todoRepository) GetByID(ctx context.Context, id uint) (*domain.Todo, error) {
	tracer := otel.Tracer("todo-repository")
	_, span := tracer.Start(ctx, "TodoRepository.GetByID")
	defer span.End()

	var todo domain.Todo
	err := r.db.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) GetAll(ctx context.Context) ([]domain.Todo, error) {
	tracer := otel.Tracer("todo-repository")
	_, span := tracer.Start(ctx, "TodoRepository.GetAll")
	defer span.End()

	var todos []domain.Todo
	err := r.db.Find(&todos).Error
	return todos, err
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	tracer := otel.Tracer("todo-repository")
	_, span := tracer.Start(ctx, "TodoRepository.Update")
	defer span.End()

	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	tracer := otel.Tracer("todo-repository")
	_, span := tracer.Start(ctx, "TodoRepository.Delete")
	defer span.End()

	return r.db.Delete(&domain.Todo{}, id).Error
}
