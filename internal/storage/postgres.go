package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Krchnk/todo-test/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return s.db.QueryRow(context.Background(), query,
		task.Title,
		task.Description,
		task.Status,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (s *Storage) GetAllTasks() ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		task.CreatedAt = createdAt
		task.UpdatedAt = updatedAt
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(id int, input models.TaskInput) (*models.Task, error) {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var task models.Task
	query := `
		UPDATE tasks
		SET 
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			status = COALESCE($3, status),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, title, description, status, created_at, updated_at`

	err = tx.QueryRow(context.Background(), query,
		input.Title,
		input.Description,
		input.Status,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *Storage) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	res, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}
