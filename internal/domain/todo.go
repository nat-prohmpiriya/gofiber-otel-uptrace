package domain

import (
	"context"
)

// Todo เป็น entity หลักของระบบ
type Todo struct {
	Base
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// TodoRepository เป็น interface สำหรับจัดการข้อมูล Todo
type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	GetByID(ctx context.Context, id uint) (*Todo, error)
	GetAll(ctx context.Context) ([]Todo, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id uint) error
}
