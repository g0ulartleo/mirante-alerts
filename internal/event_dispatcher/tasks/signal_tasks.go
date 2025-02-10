package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/hibiken/asynq"
)

const (
	TypeSignalWrite = "signal:write"
)

type SignalWritePayload struct {
	SentinelID int
	Signal     sentinel.Signal
}

func NewSignalWriteTask(sentinelID int, signal sentinel.Signal) (*asynq.Task, error) {
	payload, err := json.Marshal(SignalWritePayload{SentinelID: sentinelID, Signal: signal})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSignalWrite, payload, asynq.MaxRetry(3)), nil
}

func HandleSignalWriteTask(ctx context.Context, t *asynq.Task) error {
	var p SignalWritePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Writing signal: sentinel_id=%d, signal=%v", p.SentinelID, p.Signal)
	return nil
}
