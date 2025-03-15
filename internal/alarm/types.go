package alarm

import "github.com/g0ulartleo/mirante-alerts/internal/signal"

type Alarm struct {
	ID            string                 `yaml:"id"`
	Name          string                 `yaml:"name"`
	Description   string                 `yaml:"description"`
	Path          []string               `yaml:"-"`
	Type          string                 `yaml:"type"`
	Config        map[string]any         `yaml:"config"`
	Cron          string                 `yaml:"cron"`
	Interval      string                 `yaml:"interval"`
	Notifications AlarmNotifications     `yaml:"notifications"`
}

func HasNotificationsEnabled(alarm *Alarm) bool {
	return len(alarm.Notifications.Email.To) > 0 || alarm.Notifications.Slack.WebhookURL != ""
}

type AlarmNotifications struct {
	Email                EmailNotificationConfig `yaml:"email"`
	Slack                SlackNotificationConfig `yaml:"slack"`
	NotifyMissingSignals bool                    `yaml:"notify_missing_signals"`
}

type EmailNotificationConfig struct {
	To []string `yaml:"to"`
}

type SlackNotificationConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type AlarmSignals struct {
	Alarm   Alarm
	Signals []signal.Signal
}
