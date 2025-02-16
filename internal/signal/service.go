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

func (s *Service) GetAlertLatestSignals(alertID string, limit int) ([]Signal, error) {
	return s.repo.GetAlertLatestSignals(alertID, limit)
}

func (s *Service) GetAlertHealth(alertID string) (Status, error) {
	return s.repo.GetAlertHealth(alertID)
}

func (s *Service) CleanOldSignals() error {
	return s.repo.CleanOldSignals()
}
