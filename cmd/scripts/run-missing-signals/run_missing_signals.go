package main

import (
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/signal/stores"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	"github.com/hibiken/asynq"
)

func RunMissingSignals() {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	err := alarm.InitAlarms()
	if err != nil {
		log.Fatalf("failed to initialize alarm configs: %v", err)
	}
	signalStore, err := stores.NewStore(config.LoadAppConfigFromEnv())
	if err != nil {
		log.Fatalf("failed to initialize signal store: %v", err)
	}
	defer signalStore.Close()
	signalService := signal.NewService(signalStore)

	for _, a := range alarm.Alarms {
		signals, err := signalService.GetAlarmLatestSignals(a.ID, 1)
		if err != nil {
			log.Fatalf("failed to get latest signals for alarm %s: %v", a.ID, err)
		}
		if len(signals) == 0 {
			runSentinel(conn, a)
		}
	}
}

func runSentinel(conn *asynq.Client, a *alarm.Alarm) {
	log.Printf("running sentinel %s", a.ID)
	task, err := alarmTasks.NewCheckAlarmTask(a.ID)
	if err != nil {
		log.Fatalf("failed to create check alarm task for alarm %s: %v", a.ID, err)
	}
	if _, err := conn.Enqueue(task); err != nil {
		log.Fatalf("failed to enqueue check alarm task for alarm %s: %v", a.ID, err)
	}
	log.Printf("enqueued check alarm task for alarm %s", a.ID)
}

func main() {
	RunMissingSignals()
}
