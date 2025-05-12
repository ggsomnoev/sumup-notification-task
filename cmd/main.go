package main

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/messaging/rabbitmq"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/producer"
	"github.com/ggsomnoev/sumup-notification-task/internal/pg"
	"github.com/ggsomnoev/sumup-notification-task/internal/webapi"
	"github.com/sendgrid/sendgrid-go"

	// "github.com/twilio/twilio-go"
	"github.com/lateralusd/textbelt"
)

var texterTimeout = 2 * time.Second

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"local"`

	APIPort string `env:"API_PORT" envDefault:"8080"`

	DBConnectionURL   string        `env:"DB_CONNECTION_URL" envDefault:"postgres://notfuser:notfpass@notificationdb:5432/notificationdb"`
	DBMaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"30m"`
	DBMaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" envDefault:"5m"`
	DBHealthCheck     time.Duration `env:"DB_HEALTH_CHECK_PERIOD" envDefault:"1m"`
	DBMinConns        int32         `env:"DB_MIN_CONNS" envDefault:"1"`
	DBMaxConns        int32         `env:"DB_MAX_CONNS" envDefault:"2"`

	RabbitMQConnURL  string `env:"RABBITMQ_CONN_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	RabbitMQQueue    string `env:"RABBITMQ_QUEUE" envDefault:"notifications_queue"`
	RabbitMQCAFile   string `env:"RABBITMQ_CA_FILE" envDefault:"/etc/rabbitmq/ca-cert.pem"`
	RabbitMQCertFile string `env:"RABBITMQ_CERT_FILE" envDefault:"/etc/rabbitmq/client-cert.pem"`
	RabbitMQKeyFile  string `env:"RABBITMQ_KEY_FILE" envDefault:"/etc/rabbitmq/client-key.pem"`

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

	var tlsConfig *rabbitmq.TLSConfig
	if cfg.AppEnv != "local" {
		tlsConfig = &rabbitmq.TLSConfig{
			CAFile:   cfg.RabbitMQCAFile,
			CertFile: cfg.RabbitMQCertFile,
			KeyFile:  cfg.RabbitMQKeyFile,
		}
	}
	rmqClient, err := rabbitmq.NewClient(cfg.RabbitMQConnURL, cfg.RabbitMQQueue, tlsConfig)
	if err != nil {
		panic(fmt.Errorf("failed to initially connect to Rabbitmq, exiting - %w", err))
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
