package repo

import (
	"context"
	"github.com/yosakoo/task-tracker-scheduler-/internal/domain/models"
	"github.com/yosakoo/task-tracker-scheduler-/pkg/postgres"
)

type UserRepo struct {
	s *postgres.Storage
}

func NewUserRepo(pg *postgres.Storage) *UserRepo {
	return &UserRepo{s: pg}
}

func (ur *UserRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
    var user models.User
    err := ur.s.Pool.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (ur *UserRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
    rows, err := ur.s.Pool.Query(ctx, "SELECT id, name, email FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}