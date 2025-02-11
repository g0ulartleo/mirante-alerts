package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/hibiken/asynq"
)

const (
	TypeSentinelRun = "sentinel:run"
)

type SentinelRunPayload struct {
	SentinelID int
}

func NewSentinelRunTask(sentinelID int) (*asynq.Task, error) {
	payload, err := json.Marshal(SentinelRunPayload{SentinelID: sentinelID})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	return asynq.NewTask(TypeSentinelRun, payload, asynq.MaxRetry(1)), nil
}

func HandleSentinelRunTask(ctx context.Context, t *asynq.Task, signalService *signal.Service) error {
	var payload SentinelRunPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	sentinelConfig, err := config.GetSentinelConfig(payload.SentinelID)
	if err != nil {
		return fmt.Errorf("failed to load sentinel config: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Running sentinel ID %d with config %v", payload.SentinelID, sentinelConfig.Name)
	sentinel, err := sentinel.Factory.GetSentinel(sentinelConfig.Type)
	if err != nil {
		return fmt.Errorf("failed to get sentinel from factory: %v: %w", err, asynq.SkipRetry)
	}
	err = sentinel.Configure(sentinelConfig.Config)
	if err != nil {
		return fmt.Errorf("failed to configure sentinel: %v: %w", err, asynq.SkipRetry)
	}
	signal, err := sentinel.Check(ctx)
	if err != nil {
		return fmt.Errorf("failed to check sentinel: %v", err)
	}
	log.Printf("Sentinel %d returned signal: %v", payload.SentinelID, signal)
	return signalService.WriteSignal(signal)
}
