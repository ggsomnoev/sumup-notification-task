package main

import (
	"fmt"
	"time"

	"github.com/ggsomnoev/sumup-notification-task/internal/config"
	"github.com/ggsomnoev/sumup-notification-task/internal/healthcheck"
	"github.com/ggsomnoev/sumup-notification-task/internal/healthcheck/service"
	"github.com/ggsomnoev/sumup-notification-task/internal/healthcheck/service/component"
	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/messaging/rabbitmq"
	"github.com/ggsomnoev/sumup-notification-task/internal/pg"
	"github.com/ggsomnoev/sumup-notification-task/internal/webapi"
	"github.com/sendgrid/sendgrid-go"

	// "github.com/twilio/twilio-go"
	"github.com/lateralusd/textbelt"
)

func main() {
	appController := lifecycle.NewController()
	appCtx, procSpawnFn := appController.Start()

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed reading configuration: %w", err))
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
		panic(fmt.Errorf("failed initializing DB pool: %w", err))
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
		panic(fmt.Errorf("failed to connect to RabbitMQ: %w", err))
	}

	// Does not support PH or BG phone numbers. Easy to implement, but expensive.
	// twilioSMSClient := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: cfg.TwilioAccountSSID,
	// 	Password: cfg.TwilioAuthToken,
	// })

	textBeltSMSClient := textbelt.New(
		textbelt.WithKey("textbelt"),
		textbelt.WithTimeout(2*time.Second),
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

	srv := webapi.NewServer(appCtx)

	rmqConn := component.NewRabbitMQChecker(rmqClient.Connection())
	dbConn := component.NewDBChecker(pool)
	healthCheckService := service.NewHealthCheckService(rmqConn, dbConn)
	healthcheck.Process(procSpawnFn, appCtx, srv, healthCheckService)

	webapi.Start(procSpawnFn, srv, cfg.APIPort)

	appController.Wait()
}
