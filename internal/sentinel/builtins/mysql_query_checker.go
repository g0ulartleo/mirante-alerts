package builtins

import (
	"context"
	"fmt"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/connections"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type MySQLQueryCheckerSentinel struct {
	query      string
	expected   string
	connection *connections.MySQLConnection
}

func NewMySQLQueryCheckerSentinel() sentinel.Sentinel {
	return &MySQLQueryCheckerSentinel{}
}

func (s *MySQLQueryCheckerSentinel) Configure(config map[string]any) error {
	for _, field := range []string{"connection", "query", "expected"} {
		if _, ok := config[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	s.query = config["query"].(string)
	if _, ok := config["expected"].(string); ok {
		s.expected = config["expected"].(string)
	} else if _, ok := config["expected"].(int); ok {
		s.expected = fmt.Sprintf("%d", config["expected"].(int))
	} else if _, ok := config["expected"].(float64); ok {
		s.expected = fmt.Sprintf("%f", config["expected"].(float64))
	} else {
		return fmt.Errorf("'expected' must be a string, int or float")
	}
	if s.expected == "" {
		return fmt.Errorf("'expected' cannot be empty")
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

func (s *MySQLQueryCheckerSentinel) Check(ctx context.Context, alarmID string) (signal.Signal, error) {
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

	var response any
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
			Status:    signal.StatusUnhealthy,
			Timestamp: time.Now(),
			Message:   "query returned no rows",
		}, nil
	}
	match := s.expected == fmt.Sprintf("%v", response)
	if match {
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
