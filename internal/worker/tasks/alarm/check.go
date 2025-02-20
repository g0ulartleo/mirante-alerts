package alarm

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
	TypeCheckAlarm = "alarm:check"
)

type CheckAlarmPayload struct {
	AlarmID string
}

func NewCheckAlarmTask(alarmID string) (*asynq.Task, error) {
	payload, err := json.Marshal(CheckAlarmPayload{AlarmID: alarmID})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	return asynq.NewTask(TypeCheckAlarm, payload, asynq.MaxRetry(1)), nil
}

func HandleCheckAlarmTask(ctx context.Context, t *asynq.Task, signalService *signal.Service, asyncClient *asynq.Client) error {
	var payload CheckAlarmPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	alarmConfig, err := alarm.GetAlarmConfig(payload.AlarmID)
	if err != nil {
		return fmt.Errorf("failed to load alarm config: %v: %w", err, asynq.SkipRetry)
	}
	sentinel, err := initializeSentinel(alarmConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize sentinel: %w", err)
	}
	log.Printf("Sentinel checking alarm ID %s", payload.AlarmID)
	sig, err := sentinel.Check(ctx, payload.AlarmID)
	if err != nil {
		return fmt.Errorf("failed to check sentinel: %w", err)
	}
	log.Printf("Alarm %s returned signal: %v", payload.AlarmID, sig)
	err = signalService.WriteSignal(sig)
	if err != nil {
		return fmt.Errorf("failed to write signal: %w", err)
	}
	if alarm.HasNotificationsEnabled(alarmConfig) {
		changed, err := signalService.AlarmHasChangedStatus(payload.AlarmID)
		if err != nil {
			return fmt.Errorf("failed to get alarm latest signals: %w", err)
		}
		if !changed {
			return nil
		}
		if sig.Status == signal.StatusUnknown && !alarmConfig.Notifications.NotifyMissingSignals {
			return nil
		}
		task, err := NewNotifyTask(payload.AlarmID, sig)
		if err != nil {
			return fmt.Errorf("failed to create notify task: %w", err)
		}
		if _, err := asyncClient.Enqueue(task); err != nil {
			return fmt.Errorf("failed to enqueue task: %w", err)
		}
	}
	return nil
}

func initializeSentinel(alarmConfig *alarm.Alarm) (sentinel.Sentinel, error) {
	sentinel, err := sentinel.Factory.GetSentinel(alarmConfig.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to get sentinel from factory: %v: %w", err, asynq.SkipRetry)
	}
	err = sentinel.Configure(alarmConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure sentinel: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sentinel configured for alarm ID %s", alarmConfig.ID)
	return sentinel, nil
}
