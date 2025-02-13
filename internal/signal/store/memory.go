package store

import (
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type MemorySignalRepository struct {
	signals map[string][]signal.Signal
}

func NewMemorySignalRepository() *MemorySignalRepository {
	return &MemorySignalRepository{signals: make(map[string][]signal.Signal)}
}

func (r *MemorySignalRepository) Save(signal signal.Signal) error {
	r.signals[signal.SentinelID] = append(r.signals[signal.SentinelID], signal)
	return nil
}

func (r *MemorySignalRepository) GetSentinelLatestSignals(sentinelID string, limit int) ([]signal.Signal, error) {
	signals := r.signals[sentinelID]
	if len(signals) == 0 {
		return nil, nil
	}
	return signals[len(signals)-limit:], nil
}

func (r *MemorySignalRepository) GetSentinelHealth(sentinelID string) (signal.Status, error) {
	signals, err := r.GetSentinelLatestSignals(sentinelID, 1)
	if err != nil {
		return signal.StatusUnknown, err
	}
	if len(signals) == 0 {
		return signal.StatusUnknown, nil
	}
	return signals[0].Status, nil
}
