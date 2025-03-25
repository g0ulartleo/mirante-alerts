package sentinel

import (
	"fmt"
	"log"
)

type SentinelFactory struct {
	sentinels map[string]func() Sentinel
}

func (f *SentinelFactory) Register(sentinelType string, factory func() Sentinel) {
	log.Printf("Registering sentinel type: %s", sentinelType)
	f.sentinels[sentinelType] = factory
}

func (f *SentinelFactory) Create(sentinelType string) (Sentinel, error) {
	factory, exists := f.sentinels[sentinelType]
	if !exists {
		return nil, fmt.Errorf("unknown sentinel type: %s", sentinelType)
	}
	return factory(), nil
}

func NewFactory() *SentinelFactory {
	return &SentinelFactory{
		sentinels: make(map[string]func() Sentinel),
	}
}
