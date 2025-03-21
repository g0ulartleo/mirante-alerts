package api

import (
	"log"
	"net/http"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/web/dashboard"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func APIKeyAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey != config.Env().APIKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key")
			}
			return next(c)
		}
	}
}

func RegisterRoutes(e *echo.Echo, signalService *signal.Service, alarmService *alarm.AlarmService, asyncClient *asynq.Client) {
	if config.Env().APIKey == "" {
		log.Println("API_KEY is not set")
		return
	}

	api := e.Group("/api")
	api.Use(APIKeyAuthMiddleware())

	api.GET("/alarm/signals", func(c echo.Context) error {
		alarmSignals, err := dashboard.GetAlarmSignals(signalService, alarmService)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, alarmSignals)
	})

	api.GET("/alarm/:alarm_id/signals", func(c echo.Context) error {
		alarmID := c.Param("alarm_id")
		alarmSignals, err := signalService.GetAlarmLatestSignals(alarmID, 10)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, alarmSignals)
	})

	api.GET("/list_alarms", func(c echo.Context) error {
		alarms, err := alarmService.GetAlarms()
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, alarms)
	})

	api.GET("/alarm/:alarm_id", func(c echo.Context) error {
		alarmID := c.Param("alarm_id")
		alarm, err := alarmService.GetAlarm(alarmID)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, alarm)
	})

	api.DELETE("/alarm/:alarm_id", func(c echo.Context) error {
		alarmID := c.Param("alarm_id")
		if err := alarmService.DeleteAlarm(alarmID); err != nil {
			log.Printf("Error deleting alarm: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Alarm deleted"})
	})

	api.POST("/alarm", func(c echo.Context) error {
		alarm := new(alarm.Alarm)
		if err := c.Bind(alarm); err != nil {
			log.Printf("Error binding alarm: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := alarmService.SetAlarm(alarm); err != nil {
			log.Printf("Error setting alarm: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, alarm)
	})

	api.POST("/alarm/:alarm_id/check", func(c echo.Context) error {
		alarmID := c.Param("alarm_id")
		task, err := alarmTasks.NewCheckAlarmTask(alarmID)
		if err != nil {
			log.Printf("Error creating check alarm task: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if _, err := asyncClient.Enqueue(task); err != nil {
			log.Printf("Error enqueueing task: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Task enqueued"})
	})
}
