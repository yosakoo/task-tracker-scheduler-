package app

import (
	"context"
	"fmt"
	"time"

	"github.com/yosakoo/task-tracker-scheduler-/internal/config"
	"github.com/yosakoo/task-tracker-scheduler-/internal/repository"
	"github.com/yosakoo/task-tracker-scheduler-/internal/service"
	"github.com/yosakoo/task-tracker-scheduler-/pkg/logger"
	"github.com/yosakoo/task-tracker-scheduler-/pkg/postgres"
	"github.com/yosakoo/task-tracker-scheduler-/pkg/rabbitmq"

	"github.com/robfig/cron/v3"
)

func Run(cfg *config.Config) {

	l := logger.New(cfg.Log.Level)
	l.Info("start server")
	pg, err := postgres.New(cfg.PG.URL, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	rmqConn, err := rabbitmq.New(rabbitmq.Config{
		URL:          cfg.RabbitMQ.URL,
		WaitTime:     5 * time.Second,
		Attempts:     10,
		Exchange:     cfg.RabbitMQ.Exchange,
		ExchangeType: cfg.RabbitMQ.ExchangeType,
		Queue:        cfg.RabbitMQ.Queue,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("failed to create RabbitMQ connection: %w", err))
	}
	defer rmqConn.Close()

	l.Info("RabbitMQ connected")

	repos := repo.NewRepositories(pg)
	services := service.NewServices(service.Deps{
		Repos:     repos,
		Log:       l,
		QueueConn: rmqConn,
	})

	moscowLocation, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        l.Fatal(fmt.Errorf("failed to load Moscow location: %w", err))
    }
    c := cron.New(cron.WithLocation(moscowLocation))
    
	_, err = c.AddFunc("0 0 * * *", func() {
		err := services.Tasks.GenerateAndSendReports(context.Background())
		if err != nil {
			l.Error(fmt.Errorf("failed to send reports: %w", err))
		}
	})
	if err != nil {
		l.Fatal(fmt.Errorf("failed to add cron job: %w", err))
	}
	c.Start()

	select {}
}
