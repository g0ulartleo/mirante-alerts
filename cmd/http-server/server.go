package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmStores "github.com/g0ulartleo/mirante-alerts/internal/alarm/stores"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	signalStores "github.com/g0ulartleo/mirante-alerts/internal/signal/stores"
	"github.com/g0ulartleo/mirante-alerts/internal/web"
	"github.com/g0ulartleo/mirante-alerts/internal/web/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	alarmStore, err := alarmStores.NewAlarmStore()
	if err != nil {
		log.Fatalf("Error initializing alarm store: %v", err)
	}
	defer alarmStore.Close()
	alarmService := alarm.NewAlarmService(alarmStore)
	err = alarm.InitAlarms(alarmStore)
	if err != nil {
		log.Fatalf("Error initializing alarm configs: %v", err)
	}
	signalStore, err := signalStores.NewStore(config.LoadAppConfigFromEnv())
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
		alarmSignals, err := web.GetAlarmSignals(signalService, alarmService)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return RenderError(c, http.StatusInternalServerError, err)
		}
		return Render(c, http.StatusOK, templates.Alarms(alarmSignals))
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
		alarmSignals, err := web.GetAlarmSignals(signalService, alarmService)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return RenderError(c, http.StatusInternalServerError, err)
		}

		return Render(c, http.StatusOK, templates.Treemap(alarmSignals, level, baseURL))
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
