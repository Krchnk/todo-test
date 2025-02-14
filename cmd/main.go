package main

import (
	"log"
	"os"

	"github.com/Krchnk/todo-test/internal/handlers"
	"github.com/Krchnk/todo-test/internal/storage"
	"github.com/gofiber/fiber/v2"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost:5432/todo?sslmode=disable"
	}

	storage, err := storage.New(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	app := fiber.New()

	taskHandler := handlers.NewTaskHandler(storage)

	app.Post("/tasks", taskHandler.CreateTask)
	app.Get("/tasks", taskHandler.GetAllTasks)
	app.Put("/tasks/:id", taskHandler.UpdateTask)
	app.Delete("/tasks/:id", taskHandler.DeleteTask)

	log.Fatal(app.Listen(":3000"))
}
