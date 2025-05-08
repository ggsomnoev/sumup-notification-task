package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer"
	"github.com/ggsomnoev/sumup-notification-task/internal/rabbitmq"
	"github.com/ggsomnoev/sumup-notification-task/internal/webapi"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	APIPort         string `env:"API_PORT" envDefault:"8080"`
	RabbitMQConnURL string `env:"RABBITMQ_CONN_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	RabbitMQQueue   string `env:"RABBITMQ_QUEUE" envDefault:"notifications_queue"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Errorf("failed reading configuration, exiting - %w", err))
	}

	appController := lifecycle.NewController()
	appCtx, procSpawnFn := appController.Start()

	srv := webapi.NewServer(appCtx)

	// TODO: Change to DialTLS.
	rmqConn, err := amqp.Dial(cfg.RabbitMQConnURL)
	if err != nil {
		panic(fmt.Errorf("failed to dial, exiting - %w", err))
	}

	client, err := rabbitmq.NewClient(rmqConn, cfg.RabbitMQQueue)
	if err != nil {
		panic(fmt.Errorf("failed to dial, exiting - %w", err))
	}

	notificationproducer.Process(procSpawnFn, appCtx, srv, client)

	webapi.Start(procSpawnFn, srv, cfg.APIPort)

	appController.Wait()
}
