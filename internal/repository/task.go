package repo

import (
    "context"
    "time"

    "github.com/yosakoo/task-tracker-scheduler-/internal/domain/models"
    "github.com/yosakoo/task-tracker-scheduler-/pkg/postgres"
)

type TaskRepo struct {
    s *postgres.Storage
}

func NewTaskRepo(pg *postgres.Storage) *TaskRepo {
    return &TaskRepo{s: pg}
}

func (tr *TaskRepo) GetIncompleteTasksForUser(ctx context.Context, userID int) ([]*models.Task, error) {
    query := "SELECT id, status, user_id, title, text, time FROM tasks WHERE status != 'completed' AND user_id = $1"
    rows, err := tr.s.Pool.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*models.Task
    for rows.Next() {
        var task models.Task
        err := rows.Scan(&task.ID, &task.Status, &task.UserID, &task.Title, &task.Text, &task.Time)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, &task)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (tr *TaskRepo) GetCompletedTasksForUserSince(ctx context.Context, userID int, since time.Time) ([]*models.Task, error) {
    query := "SELECT id, status, user_id, title, text, time FROM tasks WHERE status = 'completed' AND user_id = $1 AND time >= $2"
    rows, err := tr.s.Pool.Query(ctx, query, userID, since)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*models.Task
    for rows.Next() {
        var task models.Task
        err := rows.Scan(&task.ID, &task.Status, &task.UserID, &task.Title, &task.Text, &task.Time)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, &task)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return tasks, nil
}
