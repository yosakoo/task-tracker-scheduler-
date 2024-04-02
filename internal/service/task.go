package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yosakoo/task-tracker-scheduler-/internal/domain/models"
	"github.com/yosakoo/task-tracker-scheduler-/internal/repository"
)

const subject = "Отчет о задачах за сутки"

type TaskService struct {
	taskRepo repo.Tasks
	userRepo repo.Users
	emailer  Emails
}

func NewTaskService(taskRepo repo.Tasks, userRepo repo.Users, emailer Emails) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
		emailer:  emailer,
	}
}

func (ts *TaskService) GenerateAndSendReports(ctx context.Context) error {
	users, err := ts.userRepo.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		incompleteTasks, err := ts.taskRepo.GetIncompleteTasksForUser(ctx, user.ID)
		if err != nil {
			return err
		}
		completedTasks, err := ts.taskRepo.GetCompletedTasksForUserSince(ctx, user.ID, time.Now().Add(-24*time.Hour))
		if err != nil {
			return err
		}

		if len(incompleteTasks) > 0 {
			if len(completedTasks) > 0 {
				err := ts.sendEmailWithBothReports(user, incompleteTasks, completedTasks)
				if err != nil {
					return err
				}
			} else {
				err := ts.sendEmailWithIncompleteTasks(user, incompleteTasks)
				if err != nil {
					return err
				}
			}
		} else if len(completedTasks) > 0 {
			err := ts.sendEmailWithCompletedTasks(user, completedTasks)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ts *TaskService) sendEmailWithBothReports(user *models.User, incompleteTasks, completedTasks []*models.Task) error {
	email := &Email{
		Subject: subject,
		To:      user.Email,
	}

	emailBody := fmt.Sprintf("У вас осталось %d несделанных задач:\n", len(incompleteTasks))
	for _, task := range incompleteTasks {
		emailBody += fmt.Sprintf("- %s\n", task.Title)
	}
	emailBody += "\n"
	emailBody += fmt.Sprintf("Вы выполнили %d задач:\n", len(completedTasks))
	for _, task := range completedTasks {
		emailBody += fmt.Sprintf("- %s\n", task.Title)
	}
	email.Body = emailBody

	if err := ts.emailer.SendEmail(context.Background(), email); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) sendEmailWithIncompleteTasks(user *models.User, incompleteTasks []*models.Task) error {
	email := &Email{
		Subject: subject,
		To:      user.Email,
	}

	emailBody := fmt.Sprintf("У вас осталось %d несделанных задач:\n", len(incompleteTasks))
	for _, task := range incompleteTasks {
		emailBody += fmt.Sprintf("- %s\n", task.Title)
	}
	email.Body = emailBody
	if err := ts.emailer.SendEmail(context.Background(), email); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) sendEmailWithCompletedTasks(user *models.User, completedTasks []*models.Task) error {
	email := &Email{
		Subject: subject,
		To:      user.Email,
	}

	emailBody := fmt.Sprintf("Вы выполнили %d задач:\n", len(completedTasks))
	for _, task := range completedTasks {
		emailBody += fmt.Sprintf("- %s\n", task.Title)
	}
	email.Body = emailBody
	if err := ts.emailer.SendEmail(context.Background(), email); err != nil {
		return err
	}
	return nil
}
