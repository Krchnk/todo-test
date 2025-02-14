package handlers

import (
	"errors"
	"github.com/Krchnk/todo-test/internal/models"
	"github.com/Krchnk/todo-test/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	storage *storage.Storage
}

func NewTaskHandler(s *storage.Storage) *TaskHandler {
	return &TaskHandler{storage: s}
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var input models.TaskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      "new",
	}

	if input.Status != "" {
		if !isValidStatus(input.Status) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status value",
			})
		}
		task.Status = input.Status
	}

	if err := h.storage.CreateTask(&task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *TaskHandler) GetAllTasks(c *fiber.Ctx) error {
	tasks, err := h.storage.GetAllTasks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tasks",
		})
	}

	return c.JSON(tasks)
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	var input models.TaskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Status != "" && !isValidStatus(input.Status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status value",
		})
	}

	task, err := h.storage.UpdateTask(id, input)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}

	return c.JSON(task)
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	if err := h.storage.DeleteTask(id); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete task",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func isValidStatus(status string) bool {
	switch status {
	case "new", "in_progress", "done":
		return true
	default:
		return false
	}
}
