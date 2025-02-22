package commands

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/signal/stores"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	"github.com/hibiken/asynq"
)

type RunMissingSignalsCommand struct{}

func (c *RunMissingSignalsCommand) Name() string {
	return "run-missing-signals"
}

func (c *RunMissingSignalsCommand) Run(args []string) error {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	err := alarm.InitAlarms()
	if err != nil {
		return fmt.Errorf("failed to initialize alarm configs: %v", err)
	}
	signalStore, err := stores.NewStore(config.LoadAppConfigFromEnv())
	if err != nil {
		return fmt.Errorf("failed to initialize signal store: %v", err)
	}
	defer signalStore.Close()
	signalService := signal.NewService(signalStore)

	for _, a := range alarm.Alarms {
		signals, err := signalService.GetAlarmLatestSignals(a.ID, 1)
		if err != nil {
			return fmt.Errorf("failed to get latest signals for alarm %s: %v", a.ID, err)
		}
		if len(signals) == 0 {
			runSentinel(conn, a)
		}
	}
	return nil
}

func runSentinel(conn *asynq.Client, a *alarm.Alarm) error {
	log.Printf("running sentinel %s", a.ID)
	task, err := alarmTasks.NewCheckAlarmTask(a.ID)
	if err != nil {
		return fmt.Errorf("failed to create check alarm task for alarm %s: %v", a.ID, err)
	}
	if _, err := conn.Enqueue(task); err != nil {
		return fmt.Errorf("failed to enqueue check alarm task for alarm %s: %v", a.ID, err)
	}
	log.Printf("enqueued check alarm task for alarm %s", a.ID)
	return nil
}

func init() {
	c := &RunMissingSignalsCommand{}
	Register(c.Name(), c)
}
