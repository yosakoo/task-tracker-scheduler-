package repo

import (
    "context"
	"time"

    "github.com/yosakoo/task-tracker-scheduler-/internal/domain/models"
    "github.com/yosakoo/task-tracker-scheduler-/pkg/postgres"
)

type Users interface {
    GetUserByID(ctx context.Context, userID int) (*models.User, error)
    GetAllUsers(ctx context.Context) ([]*models.User, error)
}

type Tasks interface {
    GetIncompleteTasksForUser(ctx context.Context, userID int) ([]*models.Task, error)
    GetCompletedTasksForUserSince(ctx context.Context, userID int, since time.Time) ([]*models.Task, error)
}


type Repositories struct {

	Users Users
	Tasks Tasks

}



func NewRepositories(pool *postgres.Storage) *Repositories {
    return &Repositories{
        Tasks: NewTaskRepo(pool),
		Users: NewUserRepo(pool),
    }
}
