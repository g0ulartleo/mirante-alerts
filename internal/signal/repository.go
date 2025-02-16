package signal

type SignalRepository interface {
	Init() error
	Close() error
	Save(signal Signal) error
	GetAlertLatestSignals(alertID string, limit int) ([]Signal, error)
	GetAlertHealth(alertID string) (Status, error)
	CleanOldSignals() error
}
