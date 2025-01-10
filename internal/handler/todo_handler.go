package handler

import (
	"context"
	"strconv"
	"todo-app/internal/domain"
	usecase "todo-app/internal/iusecase"
	"todo-app/utils"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoHandler struct {
	todoService *usecase.TodoService
}

func NewTodoHandler(todoService *usecase.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(context.Context)
	ctx, span := h.tr
	_, span := tracer.Start(ctx, "TodoHandler.CreateTodo")
	defer span.End()

	todo := new(domain.Todo)
	if err := c.BodyParser(todo); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	span.AddEvent("todo_data", trace.WithAttributes(
		attribute.String("input", utils.ToJSONString(todo)),
	))

	if err := h.todoService.CreateTodo(ctx, todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(todo)),
	))
	return c.Status(fiber.StatusCreated).JSON(todo)
}

func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	ctx := c.Context()
	tracer := otel.Tracer("todo-handler")
	_, span := tracer.Start(ctx, "TodoHandler.GetTodoByID")
	defer span.End()

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	todo, err := h.todoService.GetTodoByID(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	return c.JSON(todo)
}

func (h *TodoHandler) GetAllTodos(c *fiber.Ctx) error {
	ctx := c.Context()
	tracer := otel.Tracer("todo-handler")
	_, span := tracer.Start(ctx, "TodoHandler.GetAllTodos")
	defer span.End()

	todos, err := h.todoService.GetAllTodos(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(todos)
}

func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	ctx := c.Context()
	tracer := otel.Tracer("todo-handler")
	_, span := tracer.Start(ctx, "TodoHandler.UpdateTodo")
	defer span.End()

	todo := new(domain.Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	todo.ID = uint(id)

	if err := h.todoService.UpdateTodo(ctx, todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(todo)
}

func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	ctx := c.Context()
	tracer := otel.Tracer("todo-handler")
	_, span := tracer.Start(ctx, "TodoHandler.DeleteTodo")
	defer span.End()

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.todoService.DeleteTodo(ctx, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
