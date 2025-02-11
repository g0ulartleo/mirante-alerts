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

func (s *Service) GetSentinelLatestSignals(sentinelID int, limit int) ([]Signal, error) {
	return s.repo.GetSentinelLatestSignals(sentinelID, limit)
}

func (s *Service) GetSentinelHealth(sentinelID int) (Status, error) {
	return s.repo.GetSentinelHealth(sentinelID)
}
