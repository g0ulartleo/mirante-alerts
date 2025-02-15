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

func (s *Service) GetSentinelLatestSignals(sentinelID string, limit int) ([]Signal, error) {
	return s.repo.GetSentinelLatestSignals(sentinelID, limit)
}

func (s *Service) GetSentinelHealth(sentinelID string) (Status, error) {
	return s.repo.GetSentinelHealth(sentinelID)
}

func (s *Service) CleanOldSignals() error {
	return s.repo.CleanOldSignals()
}
