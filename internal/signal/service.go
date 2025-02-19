package signal

type Service struct {
	repo SignalRepository
}

func NewService(repo SignalRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) WriteSignal(signal Signal) error {
	return s.repo.Save(signal)
}

func (s *Service) GetAlarmLatestSignals(alarmID string, limit int) ([]Signal, error) {
	return s.repo.GetAlarmLatestSignals(alarmID, limit)
}

func (s *Service) GetAlarmHealth(alarmID string) (Status, error) {
	return s.repo.GetAlarmHealth(alarmID)
}

func (s *Service) CleanOldSignals() error {
	return s.repo.CleanOldSignals()
}
