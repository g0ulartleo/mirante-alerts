//go:build !mysql

package stores

import (
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

func NewMySQLSignalRepository(cfg config.MySQLConfig) (signal.SignalRepository, error) {
	return nil, fmt.Errorf("mysql driver not included in build. Use -tags mysql when building")
}
