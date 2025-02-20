package notification

import (
	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

func DispatchAlarmNotifications(alarmConfig *alarm.Alarm, sig signal.Signal) []error {
	notifications := []Notification{}
	if len(alarmConfig.Notifications.Email.To) > 0 {
		notifications = append(notifications, NewEmailNotification())
	}
	if alarmConfig.Notifications.Slack.WebhookURL != "" {
		notifications = append(notifications, NewSlackNotification())
	}

	channel := make(chan error)

	for _, notification := range notifications {
		go func(notification Notification) {
			err := notification.Build(alarmConfig, sig)
			if err != nil {
				channel <- err
			}
			err = notification.Send()
			if err != nil {
				channel <- err
			}
		}(notification)
	}

	errors := []error{}
	for err := range channel {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
