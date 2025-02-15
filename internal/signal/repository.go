package signal

type SignalRepository interface {
	Init() error
	Save(signal Signal) error
	GetSentinelLatestSignals(sentinelID string, limit int) ([]Signal, error)
	GetSentinelHealth(sentinelID string) (Status, error)
}
