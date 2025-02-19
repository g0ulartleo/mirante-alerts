package sentinel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/hibiken/asynq"
)

const (
	TypeSentinelCheckAlarm = "sentinel:check_alarm"
)

type SentinelCheckAlarmPayload struct {
	AlarmID string
}

func NewSentinelCheckAlarmTask(alarmID string) (*asynq.Task, error) {
	payload, err := json.Marshal(SentinelCheckAlarmPayload{AlarmID: alarmID})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	return asynq.NewTask(TypeSentinelCheckAlarm, payload, asynq.MaxRetry(1)), nil
}

func HandleSentinelCheckAlarmTask(ctx context.Context, t *asynq.Task, signalService *signal.Service) error {
	var payload SentinelCheckAlarmPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	alarmConfig, err := alarm.GetAlarmConfig(payload.AlarmID)
	if err != nil {
		return fmt.Errorf("failed to load alarm config: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sentinel checking alarm ID %s", payload.AlarmID)
	sentinel, err := sentinel.Factory.GetSentinel(alarmConfig.Type)
	if err != nil {
		return fmt.Errorf("failed to get sentinel from factory: %v: %w", err, asynq.SkipRetry)
	}
	err = sentinel.Configure(alarmConfig.Config)
	if err != nil {
		return fmt.Errorf("failed to configure sentinel: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sentinel configured for alarm ID %s", payload.AlarmID)
	signal, err := sentinel.Check(ctx, payload.AlarmID)
	if err != nil {
		return fmt.Errorf("failed to check sentinel: %v", err)
	}
	log.Printf("Alarm %s returned signal: %v", payload.AlarmID, signal)
	return signalService.WriteSignal(signal)
}
