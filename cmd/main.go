package main

import (
	"context"
	"fmt"
	"log"
	"todo-app/internal/domain"
	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/service"
	"todo-app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tracer = otel.Tracer("fiber-server")

func main() {
	fmt.Println("Starting todo-app")
	// Initialize tracer
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// เชื่อมต่อกับฐานข้อมูล
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5435 sslmode=disable TimeZone=Asia/Bangkok"
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

	// สร้าง Dependencies
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	// สร้าง Fiber App
	app := fiber.New()

	// Add middlewares
	app.Use(logger.New())
	app.Use(middleware.OtelMiddleware("todo-service"))

	// Routes
	api := app.Group("/api")
	todos := api.Group("/todos")

	todos.Post("/", todoHandler.CreateTodo)
	todos.Get("/", todoHandler.GetAllTodos)
	todos.Get("/:id", todoHandler.GetTodoByID)
	todos.Put("/:id", todoHandler.UpdateTodo)
	todos.Delete("/:id", todoHandler.DeleteTodo)

	fmt.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
