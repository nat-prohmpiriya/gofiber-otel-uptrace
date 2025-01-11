package handler

import (
	"todo-app/internal/domain"
	usecase "todo-app/internal/iusecase"
	"todo-app/utils"

	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"todo-app/pkg/otel"
)

type TodoHandler struct {
	todoService *usecase.TodoService
	tracer      trace.Tracer
}

func NewTodoHandler(todoService *usecase.TodoService, tracer trace.Tracer) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
		tracer:      tracer,
	}
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(context.Context)
	ctx, span := h.tracer.Start(ctx, "TodoHandler.CreateTodo")
	defer span.End()
	logger := otel.NewTraceLogger(span)

	todo := new(domain.Todo)
	if err := c.BodyParser(todo); err != nil {
		logger.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	logger.Input(todo)
	if err := h.todoService.CreateTodo(ctx, todo); err != nil {
		logger.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	logger.Output(todo)
	return c.Status(fiber.StatusCreated).JSON(todo)
}

func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(context.Context)
	ctx, span := h.tracer.Start(ctx, "TodoHandler.GetTodoByID")
	defer span.End()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	todo, err := h.todoService.GetTodoByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if todo == nil {
		span.AddEvent("output", trace.WithAttributes(
			attribute.String("output", "Todo not found"),
		))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(todo)),
	))
	return c.JSON(todo)
}

func (h *TodoHandler) GetAllTodos(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(context.Context)
	ctx, span := h.tracer.Start(ctx, "TodoHandler.GetAllTodos")
	defer span.End()

	todos, err := h.todoService.GetAllTodos(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", utils.ToJSONString(todos)),
	))
	return c.JSON(todos)
}

func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(context.Context)
	ctx, span := h.tracer.Start(ctx, "TodoHandler.UpdateTodo")
	defer span.End()

	todo := new(domain.Todo)
	if err := c.BodyParser(todo); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	if err := h.todoService.UpdateTodo(ctx, todo); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(todo)
}

func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(context.Context)
	ctx, span := h.tracer.Start(ctx, "TodoHandler.UpdateTodo")
	defer span.End()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	span.AddEvent("input", trace.WithAttributes(
		attribute.String("id", id.String()),
	))
	if err := h.todoService.DeleteTodo(ctx, id); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	span.AddEvent("output", trace.WithAttributes(
		attribute.String("output", "Todo deleted"),
	))
	return c.SendStatus(fiber.StatusNoContent)
}
