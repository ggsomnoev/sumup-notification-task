package main

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/messaging/rabbitmq"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/producer"
	"github.com/ggsomnoev/sumup-notification-task/internal/pg"
	"github.com/ggsomnoev/sumup-notification-task/internal/webapi"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sendgrid/sendgrid-go"

	// "github.com/twilio/twilio-go"
	"github.com/lateralusd/textbelt"
)

var texterTimeout = 2 * time.Second

type Config struct {
	DBConnectionURL   string        `env:"DB_CONNECTION_URL" envDefault:"postgres://notfuser:notfpass@notificationdb:5432/notificationdb"`
	DBMaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"30m"`
	DBMaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" envDefault:"5m"`
	DBHealthCheck     time.Duration `env:"DB_HEALTH_CHECK_PERIOD" envDefault:"1m"`
	DBMinConns        int32         `env:"DB_MIN_CONNS" envDefault:"1"`
	DBMaxConns        int32         `env:"DB_MAX_CONNS" envDefault:"2"`

	APIPort         string `env:"API_PORT" envDefault:"8080"`
	RabbitMQConnURL string `env:"RABBITMQ_CONN_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	RabbitMQQueue   string `env:"RABBITMQ_QUEUE" envDefault:"notifications_queue"`

	TwilioAccountSSID string `env:"TWILIO_ACC_SSID"`
	TwilioAuthToken   string `env:"TWILIO_AUTH_TOKEN"`

	SendGridAPIKey              string `env:"SEND_GRID_API_KEY,required"`
	SendGridSenderIdentityEmail string `env:"SEND_GRID_SENDER_IDENTITY_EMAIL,required"`

	SlackWebhookURL string `env:"SLACK_WEBHOOK_URL,required"`
}

func main() {
	appController := lifecycle.NewController()
	appCtx, procSpawnFn := appController.Start()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Errorf("failed reading configuration, exiting - %w", err))
	}

	dbCfg := pg.PoolConfig{
		MinConns:          cfg.DBMinConns,
		MaxConns:          cfg.DBMaxConns,
		MaxConnLifetime:   cfg.DBMaxConnLifetime,
		MaxConnIdleTime:   cfg.DBMaxConnIdleTime,
		HealthCheckPeriod: cfg.DBHealthCheck,
	}

	pool, err := pg.InitPool(appCtx, cfg.DBConnectionURL, dbCfg)
	if err != nil {
		panic(fmt.Errorf("failed initializing db connection pool, exiting - %w", err))
	}

	defer pool.Close()

	// TODO: Change to DialTLS.
	rmqConn, err := amqp.Dial(cfg.RabbitMQConnURL)
	if err != nil {
		panic(fmt.Errorf("failed to dial, exiting - %w", err))
	}

	defer func() {
		if err := rmqConn.Close(); err != nil {
			logger.GetLogger().Errorf("failed to close RabbitMQ connection: %v", err)
		}
	}()

	rmqClient, err := rabbitmq.NewClient(rmqConn, cfg.RabbitMQQueue)
	if err != nil {
		panic(fmt.Errorf("failed to dial, exiting - %w", err))
	}

	srv := webapi.NewServer(appCtx)
	producer.Process(procSpawnFn, appCtx, srv, rmqClient)

	// Does not support PH or BG phone numbers. Easy to implement, but expensive.
	// twilioSMSClient := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: cfg.TwilioAccountSSID,
	// 	Password: cfg.TwilioAuthToken,
	// })

	textBeltSMSClient := textbelt.New(
		textbelt.WithKey("textbelt"),
		textbelt.WithTimeout(texterTimeout),
	)

	mailClient := sendgrid.NewSendClient(cfg.SendGridAPIKey)

	consumer.Process(
		procSpawnFn,
		appCtx,
		pool,
		rmqClient,
		textBeltSMSClient,
		mailClient,
		cfg.SlackWebhookURL,
		cfg.SendGridSenderIdentityEmail,
	)

	webapi.Start(procSpawnFn, srv, cfg.APIPort)

	appController.Wait()
}
