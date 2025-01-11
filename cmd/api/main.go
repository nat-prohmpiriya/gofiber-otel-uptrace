package main

import (
	"fmt"
	"log"
	"todo-app/internal/domain"
	"todo-app/internal/handler"
	usecase "todo-app/internal/iusecase"
	"todo-app/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	// "go.opentelemetry.io/otel"
	"todo-app/pkg/otel"

	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var tracer = otel.Tracer("fiber-server")

func main() {
	fmt.Println("Starting todo-app")
	// Initialize tracer
	tp, err := otel.TraceProvider()
	if err != nil {
		log.Fatal("Failed to initialize tracer:", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// เชื่อมต่อกับฐานข้อมูล
	dsn := "host=postgres_app user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Connected to database")
	}

	// Auto Migrate
	err = db.AutoMigrate(&domain.Todo{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	} else {
		log.Println("Database migrated")
	}

	tracer := tp.Tracer("todo-service")

	// สร้าง Dependencies
	todoRepo := repository.NewTodoRepository(db, tracer)
	todoService := usecase.NewTodoService(todoRepo, tracer)
	todoHandler := handler.NewTodoHandler(todoService, tracer)

	// สร้าง Fiber App
	app := fiber.New()

	// Add middlewares
	app.Use(logger.New())
	app.Use(otel.OtelMiddleware("todo-service"))

	// Routes
	api := app.Group("/api")
	todos := api.Group("/todos")

	todos.Post("/", todoHandler.CreateTodo)
	todos.Get("/", todoHandler.GetAllTodos)
	todos.Get("/:id", todoHandler.GetTodoByID)
	todos.Put("/:id", todoHandler.UpdateTodo)
	todos.Delete("/:id", todoHandler.DeleteTodo)

	fmt.Println("Server is running on port 4000")
	log.Fatal(app.Listen(":4000"))
}
