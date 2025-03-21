package main

import (
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmfactory "github.com/g0ulartleo/mirante-alerts/internal/alarm/factory"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	signalfactory "github.com/g0ulartleo/mirante-alerts/internal/signal/factory"
	"github.com/g0ulartleo/mirante-alerts/internal/web/api"
	"github.com/g0ulartleo/mirante-alerts/internal/web/dashboard"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	alarmStore, err := alarmfactory.New()
	if err != nil {
		log.Fatalf("Error initializing alarm store: %v", err)
	}
	defer alarmStore.Close()
	alarmService := alarm.NewAlarmService(alarmStore)
	err = alarm.InitAlarms(alarmStore)
	if err != nil {
		log.Fatalf("Error initializing alarm configs: %v", err)
	}
	signalStore, err := signalfactory.New(config.LoadAppConfigFromEnv())
	if err != nil {
		log.Fatalf("Error initializing signal store: %v", err)
	}
	defer signalStore.Close()
	signalService := signal.NewService(signalStore)
	asyncClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: config.Env().RedisAddr,
	})
	defer asyncClient.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	api.RegisterRoutes(e, signalService, alarmService, asyncClient)
	dashboard.RegisterRoutes(e, signalService, alarmService)

	e.Logger.Fatal(e.Start(config.Env().HTTPAddr + ":" + config.Env().HTTPPort))
}
