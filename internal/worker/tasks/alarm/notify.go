package alarm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/notification"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/hibiken/asynq"
)

const (
	TypeNotify = "alarm:notify"
)

type NotifyPayload struct {
	AlarmID string
	Signal  signal.Signal
}

func NewNotifyTask(alarmID string, sig signal.Signal) (*asynq.Task, error) {
	payload, err := json.Marshal(NotifyPayload{AlarmID: alarmID, Signal: sig})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %w", err)
	}
	return asynq.NewTask(TypeNotify, payload, asynq.MaxRetry(1)), nil
}

func HandleNotifyTask(ctx context.Context, t *asynq.Task) error {
	var payload NotifyPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	alarmConfig, err := alarm.GetAlarmConfig(payload.AlarmID)
	if err != nil {
		return fmt.Errorf("failed to get alarm config: %w", err)
	}
	errors := notification.DispatchAlarmNotifications(alarmConfig, payload.Signal)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		return fmt.Errorf("failed to dispatch alarm notifications: %v", errors)
	}
	return nil
}
