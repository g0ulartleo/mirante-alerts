package main

import (
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmrepo "github.com/g0ulartleo/mirante-alerts/internal/alarm/repo"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	signalrepo "github.com/g0ulartleo/mirante-alerts/internal/signal/repo"
	"github.com/g0ulartleo/mirante-alerts/internal/web/api"
	"github.com/g0ulartleo/mirante-alerts/internal/web/dashboard"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	alarmRepo, err := alarmrepo.New()
	if err != nil {
		log.Fatalf("Error initializing alarm store: %v", err)
	}
	defer alarmRepo.Close()
	alarmService := alarm.NewAlarmService(alarmRepo)
	err = alarm.InitAlarms(alarmRepo)
	if err != nil {
		log.Fatalf("Error initializing alarm configs: %v", err)
	}
	signalRepo, err := signalrepo.New(config.LoadAppConfigFromEnv())
	if err != nil {
		log.Fatalf("Error initializing signal store: %v", err)
	}
	defer signalRepo.Close()
	signalService := signal.NewService(signalRepo)
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
