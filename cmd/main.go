package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/ggsomnoev/sumup-notification-task/internal/lifecycle"
	"github.com/ggsomnoev/sumup-notification-task/internal/webapi"
	"github.com/ggsomnoev/sumup-notification-task/internal/notificationproducer"
)

type Config struct {
	APIPort string `env:"API_PORT" envDefault:"8080"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Errorf("failed reading configuration, exiting - %w", err))
	}

	appController := lifecycle.NewController()
	appCtx, procSpawnFn := appController.Start()

	srv := webapi.NewServer(appCtx)
	notificationproducer.Process(appCtx, srv)

	webapi.Start(procSpawnFn, srv, cfg.APIPort)

	appController.Wait()
}
