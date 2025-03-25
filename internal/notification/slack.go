package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type SlackNotification struct {
	WebhookURL string
	Message    string
}

func (s *SlackNotification) Build(alarmConfig *alarm.Alarm, sig signal.Signal) error {
	s.WebhookURL = alarmConfig.Notifications.Slack.WebhookURL
	s.Message = fmt.Sprintf("*Alert:* %s (*%s*)\n*Signal:* %v", alarmConfig.Name, sig.Status, sig)
	return nil
}

func (s *SlackNotification) Send() error {
	payload := map[string]string{
		"text": s.Message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", s.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non-success response from Slack: %s", resp.Status)
	}

	return nil
}

func NewSlackNotification() *SlackNotification {
	return &SlackNotification{}
}
