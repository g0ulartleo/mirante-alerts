package signal

type SignalRepository interface {
	Save(signal Signal) error
	GetSentinelLatestSignals(sentinelID int, limit int) ([]Signal, error)
	GetSentinelHealth(sentinelID int) (Status, error)
}
