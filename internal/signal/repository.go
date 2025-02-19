package signal

type SignalRepository interface {
	Init() error
	Close() error
	Save(signal Signal) error
	GetAlarmLatestSignals(alarmID string, limit int) ([]Signal, error)
	GetAlarmHealth(alarmID string) (Status, error)
	CleanOldSignals() error
}
