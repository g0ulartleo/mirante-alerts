package builtins

import (
	"context"
	"fmt"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/connections"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type MySQLCountCheckerSentinel struct {
	query      string
	expected   int64
	connection *connections.MySQLConnection
}

func NewMySQLCountCheckerSentinel() sentinel.Sentinel {
	return &MySQLCountCheckerSentinel{}
}

func (s *MySQLCountCheckerSentinel) Configure(config map[string]any) error {
	for _, field := range []string{"connection", "query", "expected"} {
		if _, ok := config[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	s.query = config["query"].(string)
	switch v := config["expected"].(type) {
	case int:
		s.expected = int64(v)
	case float64:
		s.expected = int64(v)
	default:
		return fmt.Errorf("cant convert `expected` to int64: %v", v)
	}
	connConfig, ok := config["connection"].(map[string]any)
	if !ok {
		return fmt.Errorf("connection config must be a map")
	}
	mysqlConfig, err := connections.NewMySQLConnectionConfig(connConfig)
	if err != nil {
		return fmt.Errorf("failed to create MySQL connection config: %v", err)
	}
	s.connection, err = connections.NewMySQLConnection(*mysqlConfig)
	if err != nil {
		return fmt.Errorf("failed to create MySQL connection: %v", err)
	}
	return nil
}

func (s *MySQLCountCheckerSentinel) Check(ctx context.Context, alarmID string) (signal.Signal, error) {
	defer s.connection.Close()
	rows, err := s.connection.DB.Query(s.query)
	if err != nil {
		return signal.Signal{
			AlarmID:   alarmID,
			Status:    signal.StatusUnknown,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("failed to execute query: %v", err),
		}, nil
	}
	defer rows.Close()

	var response int64
	if rows.Next() {
		if err = rows.Scan(&response); err != nil {
			return signal.Signal{
				AlarmID:   alarmID,
				Status:    signal.StatusUnknown,
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("failed to scan query result: %v", err),
			}, nil
		}
	} else {
		return signal.Signal{
			AlarmID:   alarmID,
			Status:    signal.StatusUnknown,
			Timestamp: time.Now(),
			Message:   "query returned no rows",
		}, nil
	}

	if response == s.expected {
		return signal.Signal{
			AlarmID:   alarmID,
			Status:    signal.StatusHealthy,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("query returned %v", response),
		}, nil
	}

	return signal.Signal{
		AlarmID:   alarmID,
		Status:    signal.StatusUnhealthy,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("query returned %v, expected %v", response, s.expected),
	}, nil
}
