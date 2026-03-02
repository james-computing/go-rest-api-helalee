package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userId string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO todos (title, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, completed, created_at, updated_at, user_id
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, title, completed, userId).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetAllTodos(pool *pgxpool.Pool, userId string) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []models.Todo = []models.Todo{}

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(
			&todo.Id,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserId,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoById(pool *pgxpool.Pool, id int, userId string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sql string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todos
		WHERE id = $1 AND user_id = $2
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, sql, id, userId).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool, userId string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sql string = `
		UPDATE todos
		SET title = $2,
			completed = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $4
		RETURNING id, title, completed, created_at, updated_at, user_id
	`

	var todo models.Todo
	var err error

	err = pool.QueryRow(ctx, sql, id, title, completed, userId).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int, userId string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sql string = `
		DELETE
		FROM todos
		WHERE id = $1 AND user_id = $2
	`

	commandTag, err := pool.Exec(ctx, sql, id, userId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}
