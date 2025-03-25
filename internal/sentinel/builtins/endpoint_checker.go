package builtins

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

const (
	DefaultExpectedStatus = 200
)

type EndpointCheckerSentinel struct {
	url            string
	expectedStatus int
	expectedBody   string
	client         *http.Client
}

func NewEndpointCheckerSentinel() sentinel.Sentinel {
	return &EndpointCheckerSentinel{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (e *EndpointCheckerSentinel) Configure(config map[string]interface{}) error {
	if url, ok := config["url"]; !ok {
		return fmt.Errorf("url is required")
	} else {
		e.url = url.(string)
	}
	if expectedStatus, ok := config["expected_status"]; ok {
		if status, ok := expectedStatus.(float64); ok {
			e.expectedStatus = int(status)
		} else if status, ok := expectedStatus.(int); ok {
			e.expectedStatus = status
		} else {
			return fmt.Errorf("expected_status must be a number")
		}
	} else {
		e.expectedStatus = DefaultExpectedStatus
	}
	if expectedBody, ok := config["expected_body"]; ok {
		e.expectedBody = expectedBody.(string)
	}
	return nil
}

func (e *EndpointCheckerSentinel) Check(ctx context.Context, alarmID string) (signal.Signal, error) {
	start := time.Now()
	response, err := e.client.Get(e.url)
	responseTime := time.Since(start)
	if err != nil {
		return signal.Signal{
			AlarmID:   alarmID,
			Status:    signal.StatusUnhealthy,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("error checking endpoint: %v", err),
		}, nil
	}
	defer response.Body.Close()

	if response.StatusCode != e.expectedStatus {
		return signal.Signal{
			AlarmID:   alarmID,
			Status:    signal.StatusUnhealthy,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("expected status %d, got %d", e.expectedStatus, response.StatusCode),
		}, nil
	}

	if e.expectedBody != "" {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return signal.Signal{
				AlarmID:   alarmID,
				Status:    signal.StatusUnhealthy,
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("error reading body: %v", err),
			}, nil
		}

		if string(body) != e.expectedBody {
			return signal.Signal{
				AlarmID:   alarmID,
				Status:    signal.StatusUnhealthy,
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("expected body %s, got %s", e.expectedBody, string(body)),
			}, nil
		}
	}

	return signal.Signal{
		AlarmID:   alarmID,
		Status:    signal.StatusHealthy,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("responded status %d in %vms", response.StatusCode, responseTime.Milliseconds()),
	}, nil
}
