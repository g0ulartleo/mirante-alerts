package main

import (
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alert"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/signal/stores"
	"github.com/hibiken/asynq"
)

func RunMissingSignals() {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	err := config.InitAlerts()
	if err != nil {
		log.Fatalf("failed to initialize alert configs: %v", err)
	}
	signalStore, err := stores.NewStore(config.LoadSignalsDatabaseConfigFromEnv())
	if err != nil {
		log.Fatalf("failed to initialize signal store: %v", err)
	}
	defer signalStore.Close()
	signalService := signal.NewService(signalStore)

	for _, a := range config.Alerts {
		signals, err := signalService.GetAlertLatestSignals(a.ID, 1)
		if err != nil {
			log.Fatalf("failed to get latest signals for alert %s: %v", a.ID, err)
		}
		if len(signals) == 0 {
			runSentinel(conn, a)
		}
	}
}

func runSentinel(conn *asynq.Client, a *alert.Alert) {
	log.Printf("running sentinel %s", a.ID)
	task, err := tasks.NewSentinelCheckAlertTask(a.ID)
	if err != nil {
		log.Fatalf("failed to create sentinel check alert task for alert %s: %v", a.ID, err)
	}
	if _, err := conn.Enqueue(task); err != nil {
		log.Fatalf("failed to enqueue sentinel check alert task for alert %s: %v", a.ID, err)
	}
	log.Printf("enqueued sentinel check alert task for alert %s", a.ID)
}

func main() {
	RunMissingSignals()
}
