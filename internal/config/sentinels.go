package config

import (
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/builtin"
)

func InitSentinelFactory() {
	builtin.Register(sentinel.Factory)
}
