package config

import (
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/builtin"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/custom"
)

func InitSentinelFactory() {
	builtin.Register(sentinel.Factory)
	custom.Register(sentinel.Factory)
}
