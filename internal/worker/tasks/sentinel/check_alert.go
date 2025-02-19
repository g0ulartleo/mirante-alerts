package sentinel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alert"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/hibiken/asynq"
)

const (
	TypeSentinelCheckAlert = "sentinel:check_alert"
)

type SentinelCheckAlertPayload struct {
	AlertID string
}

func NewSentinelCheckAlertTask(alertID string) (*asynq.Task, error) {
	payload, err := json.Marshal(SentinelCheckAlertPayload{AlertID: alertID})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	return asynq.NewTask(TypeSentinelCheckAlert, payload, asynq.MaxRetry(1)), nil
}

func HandleSentinelCheckAlertTask(ctx context.Context, t *asynq.Task, signalService *signal.Service) error {
	var payload SentinelCheckAlertPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	alertConfig, err := alert.GetAlertConfig(payload.AlertID)
	if err != nil {
		return fmt.Errorf("failed to load alert config: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sentinel checking alert ID %s", payload.AlertID)
	sentinel, err := sentinel.Factory.GetSentinel(alertConfig.Type)
	if err != nil {
		return fmt.Errorf("failed to get sentinel from factory: %v: %w", err, asynq.SkipRetry)
	}
	err = sentinel.Configure(alertConfig.Config)
	if err != nil {
		return fmt.Errorf("failed to configure sentinel: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sentinel configured for alert ID %s", payload.AlertID)
	signal, err := sentinel.Check(ctx, payload.AlertID)
	if err != nil {
		return fmt.Errorf("failed to check sentinel: %v", err)
	}
	log.Printf("Alert %s returned signal: %v", payload.AlertID, signal)
	return signalService.WriteSignal(signal)
}
