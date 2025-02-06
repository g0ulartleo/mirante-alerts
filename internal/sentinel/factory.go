package sentinel

import (
	"fmt"
	"log"
)

var Factory = &SentinelFactory{
	sentinels: make(map[string]func() Sentinel),
}

func RegisterSentinelType(sentinelType string, factory func() Sentinel) {
	Factory.Register(sentinelType, factory)
}

type SentinelFactory struct {
	sentinels map[string]func() Sentinel
}

func (f *SentinelFactory) Register(sentinelType string, factory func() Sentinel) {
	log.Printf("Registering sentinel type: %s", sentinelType)
	f.sentinels[sentinelType] = factory
}

func (f *SentinelFactory) GetSentinel(sentinelType string) (Sentinel, error) {
	factory, exists := f.sentinels[sentinelType]
	if !exists {
		return nil, fmt.Errorf("unknown sentinel type: %s", sentinelType)
	}
	return factory(), nil
}
