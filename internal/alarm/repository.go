package alarm

type AlarmRepository interface {
	Init() error
	GetAlarms() ([]*Alarm, error)
	GetAlarm(alarmID string) (*Alarm, error)
	SetAlarm(alarm *Alarm) error
	DeleteAlarm(alarmID string) error
	Close() error
}
