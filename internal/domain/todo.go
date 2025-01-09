package domain

import (
	"context"
	"time"
)

// Todo เป็น entity หลักของระบบ
type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoRepository เป็น interface สำหรับจัดการข้อมูล Todo
type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	GetByID(ctx context.Context, id uint) (*Todo, error)
	GetAll(ctx context.Context) ([]Todo, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id uint) error
}
