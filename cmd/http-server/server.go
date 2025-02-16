package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/g0ulartleo/mirante-alerts/internal/alert"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/signal/stores"
	"github.com/g0ulartleo/mirante-alerts/internal/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getAlertSignals(signalService *signal.Service) ([]alert.AlertSignals, error) {
	alertsSignals := make([]alert.AlertSignals, 0)
	for _, a := range config.Alerts {
		signals, err := signalService.GetAlertLatestSignals(a.ID, 1)
		if err != nil {
			log.Printf("Error fetching signals for alert %s: %v", a.ID, err)
			signals = []signal.Signal{}
		}
		alertsSignals = append(alertsSignals, alert.AlertSignals{
			Alert:   *a,
			Signals: signals,
		})
	}
	return alertsSignals, nil
}

func main() {
	err := config.InitAlerts()
	if err != nil {
		log.Fatalf("Error initializing alert configs: %v", err)
	}
	signalStore, err := stores.NewStore(config.LoadSignalsDatabaseConfigFromEnv())
	if err != nil {
		log.Fatalf("Error initializing signal store: %v", err)
	}
	defer signalStore.Close()
	signalService := signal.NewService(signalStore)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		alertSignals, err := getAlertSignals(signalService)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return RenderError(c, http.StatusInternalServerError, err)
		}
		return Render(c, http.StatusOK, templates.Alerts(alertSignals))
	})

	e.GET("/*", func(c echo.Context) error {
		pathParam := c.Param("*")
		var level int
		var baseURL string
		if pathParam == "" {
			level = 0
			baseURL = "/"
		} else {
			segments := strings.Split(pathParam, "/")
			level = len(segments)
			baseURL = "/" + pathParam
		}
		alertSignals, err := getAlertSignals(signalService)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return RenderError(c, http.StatusInternalServerError, err)
		}

		return Render(c, http.StatusOK, templates.Treemap(alertSignals, level, baseURL))
	})

	e.Logger.Fatal(e.Start(config.Env().HTTPAddr + ":" + config.Env().HTTPPort))
}

func Render(ctx echo.Context, statusCode int, template templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)
	template = templates.Base(template)
	if err := template.Render(ctx.Request().Context(), buf); err != nil {
		return RenderError(ctx, http.StatusInternalServerError, err)
	}
	return ctx.HTML(statusCode, buf.String())
}

func RenderError(ctx echo.Context, statusCode int, err error) error {
	return ctx.HTML(statusCode, err.Error())
}
