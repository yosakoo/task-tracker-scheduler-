package service

import (
	"context"

	"github.com/yosakoo/task-tracker-scheduler-/internal/repository"
	"github.com/yosakoo/task-tracker-scheduler-/pkg/logger"
	"github.com/yosakoo/task-tracker-scheduler-/pkg/rabbitmq"
)

type Email struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json"to"`
}

type Emails interface {
	SendEmail(ctx context.Context, email *Email) error
}

type Tasks interface {
	GenerateAndSendReports(ctx context.Context) error
}

type Services struct {
	Tasks  Tasks
	Emails Emails
}

type Deps struct {
	Repos        *repo.Repositories
	QueueConn    *rabbitmq.Connection
	Log          *logger.Logger
	EmailService Emails
}

func NewServices(deps Deps) *Services {
	emailService := NewEmailService(deps.QueueConn)
	TaskService  := NewTaskService(deps.Repos.Tasks,deps.Repos.Users,emailService)
	return &Services{
		Emails: emailService,
		Tasks: TaskService,
					}
}
