package builtin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
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
		e.expectedStatus = expectedStatus.(int)
	} else {
		e.expectedStatus = DefaultExpectedStatus
	}
	if expectedBody, ok := config["expected_body"]; ok {
		e.expectedBody = expectedBody.(string)
	}
	return nil
}

func (e *EndpointCheckerSentinel) Check(ctx context.Context) (sentinel.Signal, error) {
	response, err := e.client.Get(e.url)
	if err != nil {
		return sentinel.Signal{
			Status:    sentinel.StatusUnhealthy,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("error checking endpoint: %v", err),
		}, nil
	}
	defer response.Body.Close()

	if response.StatusCode != e.expectedStatus {
		return sentinel.Signal{
			Status:    sentinel.StatusUnhealthy,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("expected status %d, got %d", e.expectedStatus, response.StatusCode),
		}, nil
	}

	if e.expectedBody != "" {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return sentinel.Signal{
				Status:    sentinel.StatusUnhealthy,
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("error reading body: %v", err),
			}, nil
		}

		if string(body) != e.expectedBody {
			return sentinel.Signal{
				Status:    sentinel.StatusUnhealthy,
				Timestamp: time.Now(),
				Message:   fmt.Sprintf("expected body %s, got %s", e.expectedBody, string(body)),
			}, nil
		}
	}

	return sentinel.Signal{
		Status:    sentinel.StatusHealthy,
		Timestamp: time.Now(),
	}, nil
}
