package stores

import (
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type MemorySignalRepository struct {
	signals map[string][]signal.Signal
}

func NewMemorySignalRepository() *MemorySignalRepository {
	return &MemorySignalRepository{signals: make(map[string][]signal.Signal)}
}

func (r *MemorySignalRepository) Init() error {
	r.signals = make(map[string][]signal.Signal)
	return nil
}

func (r *MemorySignalRepository) Save(signal signal.Signal) error {
	r.signals[signal.AlertID] = append(r.signals[signal.AlertID], signal)
	return nil
}

func (r *MemorySignalRepository) GetAlertLatestSignals(alertID string, limit int) ([]signal.Signal, error) {
	signals := r.signals[alertID]
	if len(signals) == 0 {
		return nil, nil
	}
	return signals[len(signals)-limit:], nil
}

func (r *MemorySignalRepository) GetAlertHealth(alertID string) (signal.Status, error) {
	signals, err := r.GetAlertLatestSignals(alertID, 1)
	if err != nil {
		return signal.StatusUnknown, err
	}
	if len(signals) == 0 {
		return signal.StatusUnknown, nil
	}
	return signals[0].Status, nil
}

func (r *MemorySignalRepository) CleanOldSignals() error {
	return nil
}

func (r *MemorySignalRepository) Close() error {
	return nil
}
