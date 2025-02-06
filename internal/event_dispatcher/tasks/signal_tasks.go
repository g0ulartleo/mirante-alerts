package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeSignalWrite = "signal:write"
)

type SignalWritePayload struct {
	SentinelID  int
	SignalValue bool
}

func NewSignalWriteTask(sentinelID int, signalValue bool) (*asynq.Task, error) {
	payload, err := json.Marshal(SignalWritePayload{SentinelID: sentinelID, SignalValue: signalValue})
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
	log.Printf("Writing signal: sentinel_id=%d, signal_value=%t", p.SentinelID, p.SignalValue)
	return nil
}
