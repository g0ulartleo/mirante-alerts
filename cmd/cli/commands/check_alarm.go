package commands

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	"github.com/hibiken/asynq"
)

type CheckAlarmCommand struct{}

func (c *CheckAlarmCommand) Name() string {
	return "check-alarm"
}

func (c *CheckAlarmCommand) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: ./cli %s <alarm-id>", c.Name())
	}

	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	alarmID := args[0]

	task, err := alarmTasks.NewCheckAlarmTask(alarmID)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	if _, err := conn.Enqueue(task); err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}

	log.Printf("Task enqueued")
	return nil
}

func init() {
	c := &CheckAlarmCommand{}
	Register(c.Name(), c)
}
