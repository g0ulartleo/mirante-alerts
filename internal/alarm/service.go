package alarm

type AlarmService struct {
	alarmRepository AlarmRepository
}

func NewAlarmService(alarmRepository AlarmRepository) *AlarmService {
	return &AlarmService{alarmRepository: alarmRepository}
}

func (s *AlarmService) InitAlarms() error {
	return InitAlarms(s.alarmRepository)
}

func (s *AlarmService) GetAlarm(id string) (*Alarm, error) {
	return s.alarmRepository.GetAlarm(id)
}

func (s *AlarmService) GetAlarms() ([]*Alarm, error) {
	return s.alarmRepository.GetAlarms()
}

func (s *AlarmService) SetAlarm(alarm *Alarm) error {
	return s.alarmRepository.SetAlarm(alarm)
}

func (s *AlarmService) DeleteAlarm(id string) error {
	return s.alarmRepository.DeleteAlarm(id)
}
